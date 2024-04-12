package orders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/repositories/cassettes"
	"github.com/pancher1990/cassette-rental/internal/repositories/customers"
	"github.com/pancher1990/cassette-rental/internal/repositories/films"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type FilmRepository interface {
	Create(context.Context, transaction.Querier, entities.Film) (*entities.Film, error)
	GetByID(ctx context.Context, tx transaction.Querier, id int) (*entities.Film, error)
}

type CustomerRepository interface {
	Get(ctx context.Context, tx transaction.Querier, customerID int) (*entities.Customer, error)
	UpdateBalance(ctx context.Context, tx transaction.Querier, customerID, balance int) (int, error)
}

type CassetteRepository interface {
	GetAvailableByFilmID(ctx context.Context, tx transaction.Querier, filmID int) (*entities.Cassette, error)
	UpdateStatus(ctx context.Context, tx transaction.Querier, ID int, isAvailable bool) error
}

type OrderRepository interface {
	Create(ctx context.Context, tx transaction.Querier, o entities.Order) (*entities.Order, error)
}

type OrderCassetteRepository interface {
	Create(ctx context.Context, tx transaction.Querier, o entities.OrderCassette) (*entities.OrderCassette, error)
}

type Repositories struct {
	Film          FilmRepository
	Customer      CustomerRepository
	Cassette      CassetteRepository
	Order         OrderRepository
	OrderCassette OrderCassetteRepository
}

var (
	RentPossibilityErrStatusConflict = errors.New("insufficient funds")
	ErrFilmNotFound                  = errors.New("film not found")
	ErrCustomerNotFound              = errors.New("customer not found")
	ErrCassetteNotFound              = errors.New("cassette not found")
)

func Create(r Repositories, tx transaction.TxFunc) func(context.Context, int, int, int) (*entities.Order, *entities.OrderCassette, error) {
	return func(ctx context.Context, filmID int, rentDays int, customerID int) (*entities.Order, *entities.OrderCassette, error) {
		var (
			film             *entities.Film
			customer         *entities.Customer
			cassette         *entities.Cassette
			newOrder         *entities.Order
			newCassetteOrder *entities.OrderCassette
		)

		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				film, err = r.Film.GetByID(ctx, tx, filmID)
				if errors.Is(err, films.ErrNotFound) {
					return ErrFilmNotFound
				} else if err != nil {
					return fmt.Errorf("can not find film: %w", err)
				}

				customer, err = r.Customer.Get(ctx, tx, customerID)
				if errors.Is(err, customers.ErrCustomerNotFound) {
					return ErrCustomerNotFound
				} else if err != nil {
					return fmt.Errorf("can not get user: %w", err)
				}

				if err = checkRentPossibility(film, customer, rentDays); err != nil {
					return RentPossibilityErrStatusConflict
				}

				cassette, err = r.Cassette.GetAvailableByFilmID(ctx, tx, filmID)
				if errors.Is(err, cassettes.ErrCassetteNotFound) {
					return ErrCassetteNotFound
				} else if err != nil {
					return fmt.Errorf("can not get active cassette: %w", err)
				}

				if err = r.Cassette.UpdateStatus(ctx, tx, cassette.ID, false); err != nil {
					return fmt.Errorf("can not update cassette status: %w", err)
				}

				order := entities.Order{
					CustomerID:     customerID,
					ReturnDeadline: time.Now().AddDate(0, 0, 3),
					IsActive:       true,
				}

				newOrder, err = r.Order.Create(ctx, tx, order)
				if err != nil {
					return fmt.Errorf("can create order: %w", err)
				}

				orderCassette := entities.OrderCassette{
					OrderID:    newOrder.ID,
					CassetteID: cassette.ID,
				}

				if newCassetteOrder, err = r.OrderCassette.Create(ctx, tx, orderCassette); err != nil {
					return fmt.Errorf("can not create order-cassette relation: %w", err)
				}

				newBalance := customer.Balance - rentDays*film.Price

				_, err = r.Customer.UpdateBalance(ctx, tx, customerID, newBalance)
				if err != nil {
					return fmt.Errorf("failed to update customer balance: %w", err)
				}

				return nil
			},
		)

		return newOrder, newCassetteOrder, err
	}
}
func checkRentPossibility(f *entities.Film, c *entities.Customer, rentDays int) error {
	rentCost := rentDays * f.Price
	if rentCost > c.Balance {
		return errors.New("insufficient funds")
	}
	return nil
}
