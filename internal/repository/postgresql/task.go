package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
)

// TaskRepository represents a repository for tasks
type TaskRepository struct {
	db *pgxpool.Pool
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

// GetByID retrieves a task by its ID
func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*entities.Task, error) {
	query := `
		SELECT id, code, title, reward_points, is_active, created_at
		FROM tasks
		WHERE id = $1`

	var task entities.Task
	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Code,
		&task.Title,
		&task.RewardPoints,
		&task.IsActive,
		&task.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Task not found
		}
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepository) GetByCode(ctx context.Context, code string) (*entities.Task, error) {
	query := `
	SELECT id, code,title, reward_points, is_active, created_at
	FROM tasks
	WHERE code = $1`

	var task entities.Task
	err := r.db.QueryRow(ctx, query, code).Scan(
		&task.ID,
		&task.Code,
		&task.Title,
		&task.RewardPoints,
		&task.IsActive,
		&task.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // ← Добавьте правильный return
		}
		return nil, err
	}

	return &task, nil
}
