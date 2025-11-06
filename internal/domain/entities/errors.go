package entities

import "fmt"

// InsufficientBalanceError represents an error when user doesn't have enough balance
type InsufficientBalanceError struct {
	UserID int64
}

func (e *InsufficientBalanceError) Error() string {
	return fmt.Sprintf("insufficient balance for user %d", e.UserID)
}

// TaskAlreadyCompletedError represents an error when user tries to complete already completed task
type TaskAlreadyCompletedError struct {
	UserID int64
	TaskID int64
}

func (e *TaskAlreadyCompletedError) Error() string {
	return fmt.Sprintf("task %d already completed by user %d", e.TaskID, e.UserID)
}

// UserNotFoundError represents an error when user is not found
type UserNotFoundError struct {
	ID int64
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user with id %d not found", e.ID)
}

// TaskNotFoundError represents an error when task is not found
type TaskNotFoundError struct {
	ID int64
}

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with id %d not found", e.ID)
}

// SelfReferralError represents an error when user tries to refer themselves
type SelfReferralError struct {
	UserID int64
}

func (e *SelfReferralError) Error() string {
	return fmt.Sprintf("user %d cannot refer themselves", e.UserID)
}

// AlreadyHasReferrerError represents an error when user already has a referrer
type AlreadyHasReferrerError struct {
	UserID int64
}

func (e *AlreadyHasReferrerError) Error() string {
	return fmt.Sprintf("user %d already has a referrer", e.UserID)
}
