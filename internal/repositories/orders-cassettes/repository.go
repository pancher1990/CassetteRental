package orders_cassettes

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

func (r *Repository) Create(ctx context.Context, tx transaction.Querier, o entities.OrderCassette) (*entities.OrderCassette, error) {
	sql, args, err := squirrel.
		Insert("order_cassette as o").
		Columns("order_id", "cassette_id").
		Values(
			o.OrderID,
			o.CassetteID,
		).
		Suffix(`returning
					o.id,
					o.order_id,
					o.cassette_id
		`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}

	var created entities.OrderCassette
	if err := tx.QueryRow(ctx, sql, args...).
		Scan(
			&created.ID,
			&created.OrderID,
			&created.CassetteID,
		); err != nil {

		return nil, fmt.Errorf("failed to create order cassete relation: %w", err)
	}

	return &created, nil
}
