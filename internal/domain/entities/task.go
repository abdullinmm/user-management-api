package entities

import "time"

// Task represents a task that users can complete
type Task struct {
	ID           int64     `json:"id"`
	RewardPoints int64     `json:"reward_points"`
	CreatedAt    time.Time `json:"created_at"`
	Code         string    `json:"code"`
	Title        string    `json:"title"`
	IsActive     bool      `json:"is_active"`
}

// UserTask represents a completed task by a user
type UserTask struct {
	UserID      int64     `json:"user_id"`
	TaskID      int64     `json:"task_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// UserTaskWithDetails represents a completed task with task details
type UserTaskWithDetails struct {
	UserID       int64     `json:"user_id"`
	TaskID       int64     `json:"task_id"`
	RewardPoints int64     `json:"reward_points"`
	CompletedAt  time.Time `json:"completed_at"`
	TaskCode     string    `json:"task_code"`
	TaskTitle    string    `json:"task_title"`
}
