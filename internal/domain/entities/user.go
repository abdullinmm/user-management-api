package entities

import "time"

type User struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	ReferrerID *int64    `json:"referrer_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Добавьте эту структуру, если её нет:
type UserWithReferrals struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	ReferrerID    *int64    `json:"referrer_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	ReferralCount int64     `json:"referral_count"`
}

// И эту для лидерборда:
type UserWithBalance struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Points    int64     `json:"points"`
	UpdatedAt time.Time `json:"updated_at"`
}
