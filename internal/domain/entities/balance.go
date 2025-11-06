package entities

import (
	"time"
)

// Balance - storage of the user's point balance
type Balance struct {
	UserID    int64     `json:"user_id"`
	Points    int64     `json:"points"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LeaderboardEntry display in the leaderboard
type LeaderboardEntry struct {
	UserID   int64  `json:"user_id"`
	Points   int64  `json:"points"`
	Rank     int64  `json:"rank"`
	Username string `json:"username"`
}
