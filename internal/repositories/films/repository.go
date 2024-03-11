package films

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

var (
	ErrFilmAlreadyExists = errors.New("film already exists")
	ErrFilmNotFound      = errors.New("film not found")
	// Добавьте другие ошибки по мере необходимости
)

func (r *Repository) Create(ctx context.Context, tx transaction.Querier, f entities.Film) (*entities.Film, error) {
	sql, args, err := squirrel.
		Insert("film as f").
		Columns("title", "price").
		Values(f.Title, f.Price).
		Suffix(`returning
					f.id,
					f.created_at,
					f.price,
					f.title
		`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}

	var created entities.Film
	if err := tx.QueryRow(ctx, sql, args...).
		Scan(
			&created.ID,
			&created.CreatedAt,
			&created.Price,
			&created.Title,
		); err != nil {
		if casted, ok := err.(*pgconn.PgError); ok {
			if casted.ConstraintName == "film_title_key" {
				return nil, ErrFilmAlreadyExists
			}
		}

		return nil, fmt.Errorf("failed to insert film: %w", err)
	}

	return &created, nil
}

func (r *Repository) Find(ctx context.Context, tx transaction.Querier, title string) ([]entities.Film, error) {
	sql, args, err := squirrel.
		Select(
			"f.id",
			"f.created_at",
			"f.price",
			"f.title",
		).
		From("film f").
		Where("lower(f.title) like '%?%'", strings.ToLower(title)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rows: %w", err)
	}
	defer rows.Close()

	var films []entities.Film
	for rows.Next() {
		var film entities.Film
		if err := rows.Scan(
			&film.ID,
			&film.CreatedAt,
			&film.Price,
			&film.Title,
		); err != nil {
			return nil, fmt.Errorf("failed to scan film: %w", err)
		}

		films = append(films, film)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to exec query: %w", err)
	}

	return films, nil
}

func (r *Repository) GetById(ctx context.Context, tx transaction.Querier, id int) (*entities.Film, error) {
	sql, args, err := squirrel.
		Select(
			"f.id",
			"f.created_at",
			"f.price",
			"f.title",
		).
		From("film f").
		Where("f.id = ?", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}
	var film entities.Film

	if err := tx.QueryRow(ctx, sql, args...).Scan(
		&film.ID,
		&film.CreatedAt,
		&film.Price,
		&film.Title,
	); err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrFilmNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get film: %w", err)
	}

	return &film, nil
}
