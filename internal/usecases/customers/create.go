package customers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type CustomerRepository interface {
	Get(context.Context, transaction.Querier, int) (*entities.Customer, error)
	Create(context.Context, transaction.Querier, entities.Customer) (*entities.Customer, error)
	UpdateBalance(ctx context.Context, tx transaction.Querier, customerID, balance int) (resultBalance int, err error)
}

func Create(repository CustomerRepository, tx transaction.TxFunc) func(context.Context, entities.Customer) (*entities.Customer, error) {
	return func(ctx context.Context, c entities.Customer) (*entities.Customer, error) {
		var created *entities.Customer

		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				created, err = repository.Create(ctx, tx, c)
				if err != nil {
					return fmt.Errorf("failed to create customer: %w", err)
				}

				return nil
			},
		)

		return created, err
	}
}
