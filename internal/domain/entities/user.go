package entities

import "time"

// User represents a user in the system
type User struct {
	ID         int64     `json:"id"`
	ReferrerID *int64    `json:"referrer_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	Username   string    `json:"username"`
}

// UserWithReferrals represents a user with referral statistics
type UserWithReferrals struct {
	ID            int64     `json:"id"`
	ReferralCount int64     `json:"referral_count"`
	ReferrerID    *int64    `json:"referrer_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Username      string    `json:"username"`
}

// UserWithBalance represents a user with balance information for leaderboard
type UserWithBalance struct {
	UserID    int64     `json:"user_id"`
	Points    int64     `json:"points"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}
