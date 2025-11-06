package entities

import (
	"time"
)

// User represents a user entity
type User struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
	ReferrerID    *int64    `json:"referrer_id,omitempty"`
	ReferralCount int       `json:"referral_count"`
}

// UserWithReferrals represents a user with their referrer and referrals
type UserWithReferrals struct {
	User
	Referrer  *User  `json:"referrer,omitempty"`
	Referrals []User `json:"referrals"`
}
