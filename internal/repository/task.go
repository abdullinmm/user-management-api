package repository

import (
	"context"

	"github.com/abdullinmm/grok-lean-go/user-management-api/internal/domain/entities"
)

// TaskRepository defines the interface for managing tasks.
type TaskRepository interface {
	GetAll(ctx context.Context) ([]entities.Task, error)
	GetByID(ctx context.Context, id int) (*entities.Task, error)
}

// UserTaskRepository defines the interface for managing user tasks.
type UserTaskRepository interface {
	Create(ctx context.Context, userTask *entities.UserTask) error
	GetByUserID(ctx context.Context, userID int64) ([]entities.UserTask, error)
	IsCompleted(ctx context.Context, userID, taskID int64) (bool, error)
}
