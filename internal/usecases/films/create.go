package films

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type FilmRepository interface {
	Create(context.Context, transaction.Querier, entities.Film) (*entities.Film, error)
	Find(context.Context, transaction.Querier, string) ([]entities.Film, error)
}

func Create(repository FilmRepository, tx transaction.TxFunc) func(context.Context, entities.Film) (*entities.Film, error) {
	return func(ctx context.Context, c entities.Film) (*entities.Film, error) {
		var created *entities.Film

		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				created, err = repository.Create(ctx, tx, c)
				if err != nil {
					return fmt.Errorf("failed to create film: %w", err)
				}

				return nil
			},
		)

		return created, err
	}
}
