package cassettes

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type CassetteRepository interface {
	Create(context.Context, transaction.Querier, []entities.Cassette) ([]entities.Cassette, error)
}
type FilmRepository interface {
	Find(context.Context, transaction.Querier, string) ([]entities.Film, error)
}

type Repositories struct {
	Film     FilmRepository
	Cassette CassetteRepository
}

func Create(r Repositories, tx transaction.TxFunc) func(context.Context, string, int) ([]entities.Cassette, error) {
	return func(ctx context.Context, title string, count int) ([]entities.Cassette, error) {
		created := make([]entities.Cassette, count)

		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				film, err := r.Film.Find(ctx, tx, title)
				if err != nil {
					return err
				}
				if film == nil {
					return fmt.Errorf("failed to find film")
				}
				if len(film) > 1 {
					return fmt.Errorf("failed to find unique film")
				}

				for i := range created {
					created[i] = entities.Cassette{FilmId: film[0].ID}
				}
				created, err = r.Cassette.Create(ctx, tx, created)
				if err != nil {
					return fmt.Errorf("failed to create cassette: %w", err)
				}

				return nil
			},
		)

		return created, err
	}
}
