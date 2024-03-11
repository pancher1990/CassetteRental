package customers

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

// select *
// from customer
// where
// 	email = 'valentin@mail.ru' and
// 	crypt('qwerty', password) = password

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrCustomerNotFound = errors.New("customer not found")

func (r *Repository) Create(ctx context.Context, tx transaction.Querier, c entities.Customer) (*entities.Customer, error) {
	sql, args, err := squirrel.
		Insert("customer as c").
		Columns("name", "password", "email").
		Values(
			c.Name,
			squirrel.Expr("crypt(?, gen_salt('bf'))", c.Password),
			c.Email,
		).
		Suffix(`returning
					c.id,
					c.created_at,
					c.name,
					c.is_active,
					c.balance,
					c.email
		`).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}

	var created entities.Customer
	if err := tx.QueryRow(ctx, sql, args...).
		Scan(
			&created.ID,
			&created.CreatedAt,
			&created.Name,
			&created.IsActive,
			&created.Balance,
			&created.Email,
		); err != nil {
		if casted, ok := err.(*pgconn.PgError); ok {
			if casted.ConstraintName == "customer_email_key" {
				return nil, ErrEmailAlreadyExists
			}
		}

		return nil, fmt.Errorf("failed to insert customer: %w", err)
	}

	return &created, nil
}

func (r *Repository) UpdateBalance(ctx context.Context, tx transaction.Querier, customerID, balance int) (int, error) {
	sql, args, err := squirrel.
		Update("customer as c").
		Set("balance", balance).
		Where("c.id = ?", customerID).
		Suffix("returning c.balance").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to compile sql: %w", err)
	}

	var resultBalance int
	if err := tx.QueryRow(ctx, sql, args...).Scan(&resultBalance); err != nil {
		return 0, fmt.Errorf("failed to update customer balance: %w", err)
	}

	return resultBalance, nil
}

func (r *Repository) Get(ctx context.Context, tx transaction.Querier, customerID int) (*entities.Customer, error) {
	sql, args, err := squirrel.
		Select(
			"c.id",
			"c.created_at",
			"c.name",
			"c.is_active",
			"c.balance",
			"c.email",
		).
		From("customer c").
		Where("c.id = ?", customerID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to compile sql: %w", err)
	}

	var customer entities.Customer

	if err := tx.QueryRow(ctx, sql, args...).Scan(
		&customer.ID,
		&customer.CreatedAt,
		&customer.Name,
		&customer.IsActive,
		&customer.Balance,
		&customer.Email,
	); err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrCustomerNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}
