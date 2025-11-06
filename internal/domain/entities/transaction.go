package entities

import (
	"time"
)

// Transaction audit of all transactions with points
type Transaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Delta         int64     `json:"delta"`
	ReferenceID   *int64    `json:"reference_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Reason        string    `json:"reason"`
	ReferenceType *string   `json:"reference_type,omitempty"`
}
