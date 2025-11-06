package postgresql

import (
	"context"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
	"github.com/abdullinmm/user-management-api/internal/domain/interfaces"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserTaskRepository handles user task completion operations
type UserTaskRepository struct {
	db *pgxpool.Pool
}

// NewUserTaskRepository creates a new user task repository
func NewUserTaskRepository(db *pgxpool.Pool) interfaces.UserTaskRepository {
	return &UserTaskRepository{db: db}
}

// Create records a task completion for a user
func (r *UserTaskRepository) Create(ctx context.Context, userTask *entities.UserTask) error {
	query := `
		INSERT INTO user_tasks (user_id, task_id, completed_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, userTask.UserID, userTask.TaskID, userTask.CompletedAt)
	return err
}

// IsCompleted checks if a user has already completed a specific task
func (r *UserTaskRepository) IsCompleted(ctx context.Context, userID int64, taskID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM user_tasks
			WHERE user_id = $1 AND task_id = $2
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, taskID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetByUserID retrieves all completed tasks for a user with task details
func (r *UserTaskRepository) GetByUserID(ctx context.Context, userID int64) ([]*entities.UserTaskWithDetails, error) {
	query := `
		SELECT
			ut.user_id,
			ut.task_id,
			t.code as task_code,
			t.title as task_title,
			t.reward_points,
			ut.completed_at
		FROM user_tasks ut
		JOIN tasks t ON ut.task_id = t.id
		WHERE ut.user_id = $1
		ORDER BY ut.completed_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTasks []*entities.UserTaskWithDetails
	for rows.Next() {
		var ut entities.UserTaskWithDetails
		err := rows.Scan(
			&ut.UserID,
			&ut.TaskID,
			&ut.TaskCode,
			&ut.TaskTitle,
			&ut.RewardPoints,
			&ut.CompletedAt,
		)
		if err != nil {
			return nil, err
		}
		userTasks = append(userTasks, &ut)
	}

	return userTasks, rows.Err()
}

// GetAvailableTasksForUser retrieves all active tasks that the user has not completed yet
func (r *UserTaskRepository) GetAvailableTasksForUser(ctx context.Context, userID int64) ([]*entities.Task, error) {
	query := `
		SELECT t.id, t.code, t.title, t.reward_points, t.is_active, t.created_at
		FROM tasks t
		WHERE t.is_active = true
		  AND NOT EXISTS (
			SELECT 1 FROM user_tasks ut
			WHERE ut.task_id = t.id AND ut.user_id = $1
		  )
		ORDER BY t.created_at ASC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*entities.Task
	for rows.Next() {
		var task entities.Task
		err := rows.Scan(
			&task.ID,
			&task.Code,
			&task.Title,
			&task.RewardPoints,
			&task.IsActive,
			&task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}
