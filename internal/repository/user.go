package repository

import (
	"context"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
)

// UserRepository is an interface for user repository.
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id int64) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	GetWithReferrals(ctx context.Context, id int64) (*entities.User, error)
	SetReferrer(ctx context.Context, userID, referrerID int64) error
}
