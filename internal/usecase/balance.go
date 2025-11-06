package usecases

import (
	"context"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
	"github.com/abdullinmm/user-management-api/internal/domain/interfaces"
)

// BalanceUseCase handles balance-related business logic
type BalanceUseCase struct {
	balanceRepo     interfaces.BalanceRepository
	transactionRepo interfaces.TransactionRepository
}

// NewBalanceUseCase creates a new BalanceUseCase instance
func NewBalanceUseCase(balanceRepo interfaces.BalanceRepository, transactionRepo interfaces.TransactionRepository) *BalanceUseCase {
	return &BalanceUseCase{
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}
}

// GetUserBalance returns user's current balance
func (b *BalanceUseCase) GetUserBalance(ctx context.Context, userID int64) (*entities.Balance, error) {
	return b.balanceRepo.GetByUserID(ctx, userID)
}

// GetLeaderboard returns top users by points
func (b *BalanceUseCase) GetLeaderboard(ctx context.Context, limit, offset int) ([]*entities.UserWithBalance, error) {
	return b.balanceRepo.GetLeaderboard(ctx, limit, offset)
}

// GetTransactionHistory returns user's transaction history
func (b *BalanceUseCase) GetTransactionHistory(ctx context.Context, userID int64, limit, offset int) ([]*entities.Transaction, error) {
	return b.transactionRepo.GetByUserID(ctx, userID, limit, offset)
}
