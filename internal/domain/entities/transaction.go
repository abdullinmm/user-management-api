package entities

import (
	"time"
)

// Transaction audit of all transactions with points
type Transaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Delta         int64     `json:"delta"`
	Reason        string    `json:"reason"`
	ReferenceID   *int64    `json:"reference_id,omitempty"`
	ReferenceType *string   `json:"reference_type,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
