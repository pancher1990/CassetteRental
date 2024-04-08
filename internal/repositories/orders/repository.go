package orders

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) Create(ctx context.Context, tx transaction.Querier, o entities.Order) (*entities.Order, error) {
	sql, args, err := squirrel.
		Insert("\"order\" as o").
		Columns("customer_id", "return_deadline").
		Values(
			o.CustomerId,
			o.ReturnDeadline,
		).
		Suffix(`returning
					o.id,
					o.customer_id,
					o.created_at,
					o.return_deadline,
					o.is_active
		`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}

	var created entities.Order
	if err := tx.QueryRow(ctx, sql, args...).
		Scan(
			&created.ID,
			&created.CustomerId,
			&created.CreatedAt,
			&created.ReturnDeadline,
			&created.IsActive,
		); err != nil {

		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &created, nil
}
