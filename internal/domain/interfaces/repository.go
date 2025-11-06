package interfaces

import (
	"context"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
)

// UserRepository defines operations for users
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id int64) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	GetWithReferrals(ctx context.Context, id int64) (*entities.UserWithReferrals, error)
	SetReferrer(ctx context.Context, userID, referrerID int64) error
}

// TaskRepository defines operations for tasks
type TaskRepository interface {
	GetByID(ctx context.Context, id int64) (*entities.Task, error)
	GetByCode(ctx context.Context, code string) (*entities.Task, error)
	GetActive(ctx context.Context) ([]*entities.Task, error)
	GetAll(ctx context.Context) ([]*entities.Task, error) // ← Добавьте этот метод
}

// UserTaskRepository defines operations for user_tasks
type UserTaskRepository interface {
	Create(ctx context.Context, userTask *entities.UserTask) error
	IsCompleted(ctx context.Context, userID, taskID int64) (bool, error)
	GetByUserID(ctx context.Context, userID int64) ([]*entities.UserTaskWithDetails, error) // ← Этот тоже нужен
	GetAvailableTasksForUser(ctx context.Context, userID int64) ([]*entities.Task, error)
}

// BalanceRepository defines operations for balances
type BalanceRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*entities.Balance, error)
	UpdatePoints(ctx context.Context, userID int64, delta int64) error
	GetLeaderboard(ctx context.Context, limit, offset int) ([]*entities.UserWithBalance, error)
}

// TransactionRepository defines operations for transactions
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entities.Transaction) error
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entities.Transaction, error)
}
