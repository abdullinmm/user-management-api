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

// GetByUserID retrieves all completed tasks for a user with task details
func (r *UserTaskRepository) GetByUserId(ctx context.Context, userID int64) ([]*entities.UserTaskWithDetails, error) {
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

// GetByTaskID retrieves all users who completed a specific task
func (r *UserTaskRepository) GetByTaskID(ctx context.Context, taskID int64) ([]*entities.UserTask, error) {
	query := `
		SELECT user_id, task_id, completed_at
		FROM user_tasks
		WHERE task_id = $1
		ORDER BY completed_at DESC`

	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTasks []*entities.UserTask
	for rows.Next() {
		var ut entities.UserTask
		err := rows.Scan(&ut.UserID, &ut.TaskID, &ut.CompletedAt)
		if err != nil {
			return nil, err
		}
		userTasks = append(userTasks, &ut)
	}

	return userTasks, rows.Err()
}
