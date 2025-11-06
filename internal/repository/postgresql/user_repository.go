package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
	"github.com/abdullinmm/user-management-api/internal/repository"
)

// UserRepository PostgreSQL implementation
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
        INSERT INTO users (username, referrer_id, created_at)
        VALUES ($1, $2, $3)
        RETURNING id`

	err := r.db.QueryRowContext(ctx, query, user.Username, user.ReferrerID, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID gets user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*entities.User, error) {
	query := `
		SELECT id, username, referrer_id, created_at
        FROM users
        WHERE id = $1`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.ReferrerID,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetByUsername gets user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {

	query := `
		SELECT id, username, referrer_id, created_at
        FROM users
        WHERE username = $1`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.ReferrerID,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetWithReferrals gets user with referral information
func (r *UserRepository) GetWithReferrals(ctx context.Context, id int64) (*entities.UserWithReferrals, error) {
	query := `
        SELECT u.id, u.username, u.referrer_id, u.created_at,
               COALESCE(COUNT(ref.id), 0) as referral_count
        FROM users u
        LEFT JOIN users ref ON ref.referrer_id = u.id
        WHERE u.id = $1
        GROUP BY u.id, u.username, u.referrer_id, u.created_at`

	userWithRefs := &entities.UserWithReferrals{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&userWithRefs.ID,
		&userWithRefs.Username,
		&userWithRefs.ReferrerID,
		&userWithRefs.CreatedAt,
		&userWithRefs.ReferralCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user with referrals: %w", err)
	}
	return userWithRefs, nil
}

// SetReferrer sets referrer for a user
func (r *UserRepository) SetReferrer(ctx context.Context, userID, referrerID int64) error {
	query := `
        UPDATE users
        SET referrer_id = $1
        WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, referrerID, userID)
	if err != nil {
		return fmt.Errorf("failed to set referrer: %w", err)
	}

	return nil
}
