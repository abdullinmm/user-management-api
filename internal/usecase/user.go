package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
	"github.com/abdullinmm/user-management-api/internal/domain/interfaces"
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo        interfaces.UserRepository
	balanceRepo     interfaces.BalanceRepository
	transactionRepo interfaces.TransactionRepository
	referralBonus   int64
	refereeBonus    int64
}

// NewUserUseCase creates a new UserUseCase instance
func NewUserUseCase(userRepo interfaces.UserRepository, balanceRepo interfaces.BalanceRepository, transactionRepo interfaces.TransactionRepository, referralBonus int64, refereeBonus int64) *UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
		referralBonus:   referralBonus,
		refereeBonus:    refereeBonus,
	}
}

// CreateUser creates a new user and handles referral bonuses
func (u *UserUseCase) CreateUser(ctx context.Context, username string, referrerID *int64) (*entities.User, error) {
	user := &entities.User{
		Username:   username,
		CreatedAt:  time.Now(),
		ReferrerID: referrerID,
	}

	// Create user record
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Initialize user balance with 0 points
	if err := u.balanceRepo.UpdatePoints(ctx, user.ID, 0); err != nil {
		return nil, err
	}

	// Handle referral bonuses if user was referred
	if referrerID != nil {
		// Give bonus to referee (new user)
		if u.refereeBonus > 0 {
			if err := u.givePoints(ctx, user.ID, u.refereeBonus, "Referral signup bonus", nil, nil); err != nil {
				return nil, err
			}
		}

		// Give bonus to referrer
		if u.referralBonus > 0 {
			if err := u.givePoints(ctx, *referrerID, u.referralBonus, "Referral reward", &user.ID, stringPtr("referral")); err != nil {
				return nil, err
			}
		}
	}

	return user, nil
}

// GetUser returns user information by ID
func (u *UserUseCase) GetUser(ctx context.Context, userID int64) (*entities.User, error) {
	return u.userRepo.GetByID(ctx, userID)
}

// GetUserWithReferrals returns user with referral information
func (u *UserUseCase) GetUserWithReferrals(ctx context.Context, userID int64) (*entities.UserWithReferrals, error) {
	return u.userRepo.GetWithReferrals(ctx, userID)
}

// GetUserByUsername returns user by username
func (u *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	return u.userRepo.GetByUsername(ctx, username)
}

// SetReferrer sets a referrer for existing user (if they don't have one)
func (u *UserUseCase) SetReferrer(ctx context.Context, userID, referrerID int64) error {
	// Get current user to check if they already have a referrer
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Don't allow setting referrer if user already has one
	if user.ReferrerID != nil {
		return fmt.Errorf("user already has a referrer")
	}

	// Set the referrer
	if err := u.userRepo.SetReferrer(ctx, userID, referrerID); err != nil {
		return err
	}

	// Give referral bonuses
	if u.refereeBonus > 0 {
		if err := u.givePoints(ctx, userID, u.refereeBonus, "Referral signup bonus", nil, nil); err != nil {
			return err
		}
	}

	if u.referralBonus > 0 {
		if err := u.givePoints(ctx, referrerID, u.referralBonus, "Referral reward", &userID, stringPtr("referral")); err != nil {
			return err
		}
	}

	return nil
}

// Helper function to give points and create transaction
func (u *UserUseCase) givePoints(ctx context.Context, userID, points int64, reason string, refID *int64, refType *string) error {
	// Update balance
	if err := u.balanceRepo.UpdatePoints(ctx, userID, points); err != nil {
		return err
	}

	// Create transaction record
	transaction := &entities.Transaction{
		UserID:        userID,
		Delta:         points,
		Reason:        reason,
		ReferenceID:   refID,
		ReferenceType: refType,
		CreatedAt:     time.Now(),
	}

	return u.transactionRepo.Create(ctx, transaction)
}

// Create is an alias for CreateUser with simpler signature
func (u *UserUseCase) Create(ctx context.Context, username string) (*entities.User, error) {
	return u.CreateUser(ctx, username, nil)
}

// GetByID retrieves a user by ID with balance
func (u *UserUseCase) GetByID(ctx context.Context, userID int64) (*entities.User, error) {
	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Get balance
	balance, err := u.balanceRepo.GetByUserID(ctx, userID)
	if err != nil {
		// If balance doesn't exist, set to 0
		user.Balance = 0
	} else {
		user.Balance = balance.Points
	}

	return user, nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
