package entities

import "time"

// Task represents a task entity
type Task struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	Title        string    `json:"title"`
	RewardPoints int64     `json:"reward_points"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserTask represents a user task entity
type UserTask struct {
	UserId      int64     `json:"user_id"`
	TaskId      int64     `json:"task_id"`
	CompletedAt time.Time `json:"completed_at"`
	Task        *Task     `json:"task,omitempty"`
}
