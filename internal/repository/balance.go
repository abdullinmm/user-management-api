package repository

import (
	"context"

	"github.com/abdullinmm/grok-lean-go/user-management-api/internal/domain/entities"
)

type BalanceRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*entities.Balance, error)
	UpdatePoints(ctx context.Context, userID, delta int64) error
	GetLeaderboard(ctx context.Context) ([]entities.LeaderboardEntry, error)
}

type TransactionRepository interface{}
