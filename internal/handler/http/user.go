package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abdullinmm/user-management-api/internal/middleware"
	"github.com/abdullinmm/user-management-api/internal/usecase"
	"github.com/go-chi/chi/v5"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUC *usecase.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUC *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUC: userUC,
	}
}

// GetStatus returns user status with balance and completed tasks
// GET /users/{id}/status
func (h *UserHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userUC.GetByID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":          user.ID,
		"username":    user.Username,
		"referrer_id": user.ReferrerID,
		"balance":     user.Balance,
		"created_at":  user.CreatedAt,
	})
}

// SetReferrer sets the referrer for a user
// POST /users/{id}/referrer
func (h *UserHandler) SetReferrer(w http.ResponseWriter, r *http.Request) {
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
		ReferrerID int64 `json:"referrer_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ReferrerID <= 0 {
		respondError(w, http.StatusBadRequest, "referrer_id must be positive")
		return
	}

	if err := h.userUC.SetReferrer(r.Context(), userID, req.ReferrerID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "referrer set successfully",
	})
}

// Register creates a new user
// POST /auth/register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" {
		respondError(w, http.StatusBadRequest, "username is required")
		return
	}

	user, err := h.userUC.Create(r.Context(), req.Username)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
}
