package films

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

func Find(repository FilmRepository, tx transaction.TxFunc) func(context.Context, string) ([]entities.Film, error) {
	return func(ctx context.Context, title string) ([]entities.Film, error) {
		var films []entities.Film

		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				films, err = repository.Find(ctx, tx, title)
				if err != nil {
					return fmt.Errorf("failed to search films: %w", err)
				}

				return nil
			},
		)

		return films, err
	}
}
