package postgresql

import (
	"context"

	"github.com/abdullinmm/user-management-api/internal/domain/entities"
	"github.com/abdullinmm/user-management-api/internal/domain/interfaces"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionRepository handles transaction-related database operations
type TransactionRepository struct {
	db *pgxpool.Pool
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *pgxpool.Pool) interfaces.TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction record
func (r *TransactionRepository) Create(ctx context.Context, transaction *entities.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, delta, reason, reference_type, reference_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRow(
		ctx,
		query,
		transaction.UserID,
		transaction.Delta,
		transaction.Reason,
		transaction.ReferenceType,
		transaction.ReferenceID,
		transaction.CreatedAt,
	).Scan(&transaction.ID)

	return err
}

// GetByUserID retrieves all transactions for a specific user with pagination
func (r *TransactionRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entities.Transaction, error) {
	query := `
		SELECT id, user_id, delta, reason, reference_type, reference_id, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*entities.Transaction
	for rows.Next() {
		var tx entities.Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Delta,
			&tx.Reason,
			&tx.ReferenceType,
			&tx.ReferenceID,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}

	return transactions, rows.Err()
}
