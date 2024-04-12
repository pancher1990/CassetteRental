package cassettes

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/repositories/films"
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

var ErrFilmNotFound = errors.New("film not found")

func Create(r Repositories, tx transaction.TxFunc) func(context.Context, string, int) ([]entities.Cassette, error) {
	return func(ctx context.Context, title string, count int) ([]entities.Cassette, error) {
		var created []entities.Cassette

		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				findedFilms, err := r.Film.Find(ctx, tx, title)
				switch {
				case errors.Is(err, films.ErrNotFound):
					return fmt.Errorf("%w: %w", ErrFilmNotFound, err)
				case err != nil:
					return fmt.Errorf("failed to find film: %w", err)
				case len(findedFilms) > 1:
					return fmt.Errorf("failed to find unique film")
				}

				film := findedFilms[0]
				newCassettes := make([]entities.Cassette, count)

				for i := range newCassettes {
					newCassettes[i].FilmID = film.ID
				}

				created, err = r.Cassette.Create(ctx, tx, newCassettes)
				if err != nil {
					return fmt.Errorf("failed to create cassette: %w", err)
				}

				return nil
			},
		)

		return created, err
	}
}
