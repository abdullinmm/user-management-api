package entities

import "time"

// Task represents a task entity.
type Task struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	Title        string    `json:"title"`
	RewardPoints int64     `json:"reward_points"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserTask represents a user-task relationship.
type UserTask struct {
	UserID      int64     `json:"user_id"`
	TaskID      int64     `json:"task_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// UserTaskWithDetails represents a user-task relationship with details.
type UserTaskWithDetails struct {
	UserID       int64     `json:"user_id"`
	TaskID       int64     `json:"task_id"`
	CompletedAt  time.Time `json:"completed_at"`
	TaskCode     string    `json:"task_code"`
	TaskTitle    string    `json:"task_title"`
	RewardPoints int64     `json:"reward_points"`
}
