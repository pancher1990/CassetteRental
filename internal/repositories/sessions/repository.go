package sessions

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) Create(ctx context.Context, tx transaction.Querier, customerID int, token string) error {
	sql, args, err := squirrel.
		Insert("session").
		Columns("customer_id", "token").
		Values(customerID, token).
		Suffix("returning token").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to compile sql: %w", err)
	}

	if _, err := tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	return nil
}

var ErrNotFound = errors.New("not found")

func (r *Repository) Remove(ctx context.Context, tx transaction.Querier, token string) error {
	sql, args, err := squirrel.
		Delete("session as s").
		Where("s.token = ?", token).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to compile sql: %w", err)
	}

	ct, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if ct.RowsAffected() != 1 {
		return fmt.Errorf("failed to delete session: %w", ErrNotFound)
	}

	return nil
}
