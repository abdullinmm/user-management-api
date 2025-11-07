package http

import (
	"net/http"
	"strconv"

	"github.com/abdullinmm/user-management-api/internal/usecase"
)

// BalanceHandler handles balance-related HTTP requests
type BalanceHandler struct {
	balanceUC *usecase.BalanceUseCase
}

// NewBalanceHandler creates a new balance handler
func NewBalanceHandler(balanceUC *usecase.BalanceUseCase) *BalanceHandler {
	return &BalanceHandler{
		balanceUC: balanceUC,
	}
}

// Leaderboard returns top users by balance
// GET /users/leaderboard
func (h *BalanceHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	// Get limit from query params (default 10)
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Get offset from query params (default 0)
	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	leaderboard, err := h.balanceUC.GetLeaderboard(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch leaderboard")
		return
	}

	respondJSON(w, http.StatusOK, leaderboard)
}
