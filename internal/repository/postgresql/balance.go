package postgresql

import (
	"context"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
	"github.com/abdullinmm/user-management-api/internal/domain/interfaces"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BalanceRepository handles balance-related database operations
type BalanceRepository struct {
	db *pgxpool.Pool
}

// NewBalanceRepository creates a new balance repository
func NewBalanceRepository(db *pgxpool.Pool) interfaces.BalanceRepository {
	return &BalanceRepository{db: db}
}

// GetByUserID retrieves a user's balance by user ID
func (r *BalanceRepository) GetByUserID(ctx context.Context, userID int64) (*entities.Balance, error) {
	query := `
			SELECT user_id, points, updated_at
			FROM balances
			WHERE user_id = $1`

	var balance entities.Balance
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&balance.UserID, &balance.Points, &balance.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &balance, nil
}

// UpdatePoints updates a user's balance by adding delta points
func (r *BalanceRepository) UpdatePoints(ctx context.Context, userID, delta int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update balance with delta, ensuring points don't go negative
	query := `
				UPDATE balances
				SET points = points + $2, updated_at = CURRENT_TIMESTAMP
				WHERE user_id = $1 AND (points + $2) >= 0`

	result, err := tx.Exec(ctx, query, userID, delta)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	if result.RowsAffected() == 0 {
		// Either user doesn't exist or would result in negative balance
		return &entities.InsufficientBalanceError{UserID: userID}
	}

	return tx.Commit(ctx)
}

// GetLeaderboard retrieves top users by balance with pagination
func (r *BalanceRepository) GetLeaderboard(ctx context.Context, limit, offset int) ([]*entities.UserWithBalance, error) {
	query := `
			SELECT b.user_id, u.username, b.points, b.updated_at
			FROM balances b
			JOIN users u ON b.user_id = u.id
			ORDER BY b.points DESC
			LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []*entities.UserWithBalance
	for rows.Next() {
		var entry entities.UserWithBalance
		err := rows.Scan(
			&entry.UserID, &entry.Username,
			&entry.Points, &entry.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		leaderboard = append(leaderboard, &entry)
	}

	return leaderboard, rows.Err()
}
