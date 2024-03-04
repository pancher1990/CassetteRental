package customers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pancher1990/cassette-rental/internal/transaction"
)

func UpdateBalance(repository CustomerRepository, tx transaction.TxFunc) func(ctx context.Context, customerID int, balance int) (resultBalance int, err error) {
	return func(ctx context.Context, customerID int, balance int) (int, error) {
		var resultBalance int

		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				if c, err := repository.Get(ctx, tx, customerID); err != nil {
					return fmt.Errorf("failed to get customer: %w", err)
				} else if c == nil {
					return fmt.Errorf("customer not found: %d", customerID)
				}

				resultBalance, err = repository.UpdateBalance(ctx, tx, customerID, balance)
				if err != nil {
					return fmt.Errorf("failed to update customer balance: %w", err)
				}

				return nil
			},
		)

		return resultBalance, err
	}
}
