package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
	"github.com/abdullinmm/user-management-api/internal/repository"
)

// TaskUseCase handles task-related business logic
type TaskUseCase struct {
	taskRepo        repository.TaskRepository
	userTaskRepo    repository.UserTaskRepository
	balanceRepo     repository.BalanceRepository
	transactionRepo repository.TransactionRepository
}

// NewTaskUseCase creates a new TaskUseCase instance
func NewTaskUseCase(
	taskRepo repository.TaskRepository,
	userTaskRepo repository.UserTaskRepository,
	balanceRepo repository.BalanceRepository,
	transactionRepo repository.TransactionRepository,
) *TaskUseCase {
	return &TaskUseCase{
		taskRepo:        taskRepo,
		userTaskRepo:    userTaskRepo,
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
	}
}

// GetAvailableTasks returns all active tasks
func (t *TaskUseCase) GetAvailableTasks(ctx context.Context) ([]entities.Task, error) {
	return t.taskRepo.GetAll(ctx)
}

// CompleteTask marks task as completed for user and awards points
func (t *TaskUseCase) CompleteTask(ctx context.Context, userID, taskID int64) error {
	// Check if task exists
	task, err := t.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	if !task.IsActive {
		return fmt.Errorf("task is not active")
	}

	// Check if user already completed this task
	completed, err := t.userTaskRepo.IsCompleted(ctx, userID, taskID)
	if err != nil {
		return err
	}

	if completed {
		return fmt.Errorf("task already completed")
	}

	// Record task completion
	userTask := &entities.UserTask{
		UserID:      userID,
		TaskID:      taskID,
		CompletedAt: time.Now(),
	}

	if err := t.userTaskRepo.Create(ctx, userTask); err != nil {
		return err
	}

	// Award points if task has rewards
	if task.RewardPoints > 0 {
		if err := t.balanceRepo.UpdatePoints(ctx, userID, task.RewardPoints); err != nil {
			return err
		}

		// Create transaction record
		transaction := &entities.Transaction{
			UserID:        userID,
			Delta:         task.RewardPoints,
			Reason:        fmt.Sprintf("Task completed: %s", task.Title),
			ReferenceID:   &taskID,
			ReferenceType: stringPtr("task"),
			CreatedAt:     time.Now(),
		}

		if err := t.transactionRepo.Create(ctx, transaction); err != nil {
			return err
		}
	}

	return nil
}

// GetUserTasks returns completed tasks for user
func (t *TaskUseCase) GetUserTasks(ctx context.Context, userID int64) ([]entities.UserTask, error) {
	return t.userTaskRepo.GetByUserID(ctx, userID)
}
