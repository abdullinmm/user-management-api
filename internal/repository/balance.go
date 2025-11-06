package repository

import (
	"context"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
)

// BalanceRepository is an interface for interacting with the balance data.
type BalanceRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*entities.Balance, error)
	UpdatePoints(ctx context.Context, userID, delta int64) error
	GetLeaderboard(ctx context.Context) ([]entities.LeaderboardEntry, error)
}

// TransactionRepository is an interface for interacting with the transaction data.
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entities.Transaction) error
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]entities.Transaction, error)
}
