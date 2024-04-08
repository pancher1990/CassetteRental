package cassettes

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/Masterminds/squirrel"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

var ErrCassetteNotFound = errors.New("cassette not found")

func (r *Repository) Create(ctx context.Context, tx transaction.Querier, c []entities.Cassette) ([]entities.Cassette, error) {

	queryBuilder := squirrel.
		Insert("cassette").
		Columns("film_id", "is_available")
	for _, cassette := range c {
		queryBuilder = queryBuilder.Values(cassette.FilmId, cassette.IsAvailable)
	}
	sql, args, err := queryBuilder.
		Suffix(`returning id, film_id, is_available`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to insert cassettes: %w", err)
	}

	var created []entities.Cassette
	for rows.Next() {
		var c entities.Cassette
		if err := rows.Scan(&c.ID, &c.FilmId, &c.IsAvailable); err != nil {
			return nil, fmt.Errorf("failed to scan cassette: %w", err)
		}
		created = append(created, c)
	}
	return created, nil
}

//
//func (r *Repository) Create(ctx context.Context, tx transaction.Querier, c entities.Cassette) (*entities.Cassette, error) {
//	sql, args, err := squirrel.
//		Insert("cassette as c").
//		Columns("film_id", "is_available").
//		Values(
//			c.FilmId,
//		).
//		Suffix(`returning
//					c.id,
//					c.film_id,
//					c.is_available,
//		`).
//		PlaceholderFormat(squirrel.Dollar).
//		ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("failed to compile sql: %w", err)
//	}
//
//	var created entities.Cassette
//	if err := tx.QueryRow(ctx, sql, args...).
//		Scan(
//			&created.ID,
//			&created.FilmId,
//			&created.IsAvailable,
//		); err != nil {
//
//		return nil, fmt.Errorf("failed to insert cassette: %w", err)
//	}
//
//	return &created, nil
//}

func (r *Repository) UpdateStatus(ctx context.Context, tx transaction.Querier, ID int, isAvailable bool) error {
	sql, args, err := squirrel.
		Update("cassette as c").
		Set("is_available", isAvailable).
		Where("c.id = ?", ID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to compile sql: %w", err)
	}

	if _, err := tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("failed to update cassette status: %w", err)
	}

	return nil
}

func (r *Repository) GetAvailableByFilmId(ctx context.Context, tx transaction.Querier, filmId int) (*entities.Cassette, error) {
	sql, args, err := squirrel.
		Select(
			"c.id",
			"c.film_id",
			"c.is_available",
		).
		From("cassette c").
		Where("c.film_id = ? AND c.is_available = true", filmId).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}
	var cassette entities.Cassette

	if err := tx.QueryRow(ctx, sql, args...).Scan(
		&cassette.ID,
		&cassette.FilmId,
		&cassette.IsAvailable,
	); err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCassetteNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get cassette: %w", err)
	}
	return &cassette, nil
}
