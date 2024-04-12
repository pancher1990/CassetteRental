package sessions

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pancher1990/cassette-rental/internal/transaction"
)

func Logout(repository SessionRepository, tx transaction.TxFunc) func(ctx context.Context, token string) error {
	return func(ctx context.Context, token string) error {
		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				if err = repository.Remove(ctx, tx, token); err != nil {
					return fmt.Errorf("failed to remove session: %w", err)
				}

				return nil
			},
		)

		return err
	}
}
