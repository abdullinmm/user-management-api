package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abdullinmm/user-management-api/internal/middleware"
	"github.com/abdullinmm/user-management-api/internal/usecase"
	"github.com/go-chi/chi/v5"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	taskUC *usecase.TaskUseCase
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskUC *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUC: taskUC,
	}
}

// CompleteTask marks a task as completed for a user
// POST /users/{id}/task/complete
func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Verify authenticated user matches requested user
	authUserID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || authUserID != userID {
		respondError(w, http.StatusForbidden, "access denied")
		return
	}

	var req struct {
		TaskID int64 `json:"task_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.TaskID <= 0 {
		respondError(w, http.StatusBadRequest, "task_id must be positive")
		return
	}

	if err := h.taskUC.CompleteTask(r.Context(), userID, req.TaskID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "task completed successfully",
	})
}

// ListActive returns all active tasks
// GET /tasks
func (h *TaskHandler) ListActive(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskUC.GetAvailableTasks(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}
