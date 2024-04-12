package sessions

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type Repositories struct {
	CustomerRepository
	SessionRepository
}

type CustomerRepository interface {
	GetByEmailAndPassword(ctx context.Context, tx transaction.Querier, email, password string) (*entities.Customer, error)
}

type SessionRepository interface {
	Create(ctx context.Context, tx transaction.Querier, customerID int, token string) error
	Remove(ctx context.Context, tx transaction.Querier, token string) error
}

func Login(repositories Repositories, tx transaction.TxFunc) func(ctx context.Context, email, password string) (string, error) {
	return func(ctx context.Context, email, password string) (string, error) {
		var token string

		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				customer, err := repositories.CustomerRepository.GetByEmailAndPassword(ctx, tx, email, password)
				if err != nil {
					return fmt.Errorf("failed to create session: %w", err)
				}

				if err = repositories.SessionRepository.Create(ctx, tx, customer.ID, token); err != nil {
					return fmt.Errorf("failed to create session: %w", err)
				}

				return nil
			},
		)

		return token, err
	}
}
