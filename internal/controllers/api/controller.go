package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/pancher1990/cassette-rental/internal/entities"
)

type (
	CreateCustomer        func(context.Context, entities.Customer) (*entities.Customer, error)
	UpdateCustomerBalance func(ctx context.Context, customerID int, balance int) (resultBalance int, err error)
	CreateFilm            func(context.Context, entities.Film) (*entities.Film, error)
	FindFilms             func(context.Context, string) ([]entities.Film, error)
	CreateOrder           func(context.Context, int, int, int) (*entities.Order, *entities.OrderCassette, error)
	CreateCassettes       func(context.Context, string, int) ([]entities.Cassette, error)
	Login                 func(context.Context, string, string) (string, error)
	Logout                func(context.Context, string) error
)

type Controller struct {
	CreateCustomer
	UpdateCustomerBalance
	CreateFilm
	FindFilms
	CreateOrder
	CreateCassettes
	Login
	Logout
}

type option interface {
	apply(*Controller)
}

type optionFunc func(*Controller)

func (o optionFunc) apply(c *Controller) {
	o(c)
}

func WithCustomerCreater(customerCreater CreateCustomer) option {
	return optionFunc(func(c *Controller) {
		c.CreateCustomer = customerCreater
	})
}

func WithOrderCreater(orderCreater CreateOrder) option {
	return optionFunc(func(c *Controller) {
		c.CreateOrder = orderCreater
	})
}

func WithCustomerBalanceUpdater(customerBalanceUpdater UpdateCustomerBalance) option {
	return optionFunc(func(c *Controller) {
		c.UpdateCustomerBalance = customerBalanceUpdater
	})
}

func WithFilmCreater(filmCreater CreateFilm) option {
	return optionFunc(func(c *Controller) {
		c.CreateFilm = filmCreater
	})
}

func WithFilmFinder(filmFinder FindFilms) option {
	return optionFunc(func(c *Controller) {
		c.FindFilms = filmFinder
	})
}

func WithCassetteCreater(cassetteCreater CreateCassettes) option {
	return optionFunc(func(c *Controller) {
		c.CreateCassettes = cassetteCreater
	})
}

func WithLogin(login Login) option {
	return optionFunc(func(c *Controller) {
		c.Login = login
	})
}

func WithLogout(logout Logout) option {
	return optionFunc(func(c *Controller) {
		c.Logout = logout
	})
}

func New(opts ...option) (*Controller, error) {
	var c Controller

	for _, o := range opts {
		o.apply(&c)
	}

	if c.CreateCustomer == nil {
		return nil, errors.New("customer creater required")
	}

	if c.UpdateCustomerBalance == nil {
		return nil, errors.New("customer balance updater required")
	}

	if c.CreateFilm == nil {
		return nil, errors.New("film creater required")
	}

	if c.FindFilms == nil {
		return nil, errors.New("film finder required")
	}

	if c.CreateOrder == nil {
		return nil, errors.New("order creater required")
	}

	if c.CreateCassettes == nil {
		return nil, errors.New("cassette creater required")
	}

	if c.Login == nil {
		return nil, errors.New("login required")
	}

	if c.Logout == nil {
		return nil, errors.New("logout required")
	}

	return &c, nil
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /customers", c.createCustomer)
	mux.HandleFunc("PUT /customer/balance", c.updateCustomerBalance)
	mux.HandleFunc("POST /films", c.createFilm)
	mux.HandleFunc("GET /films", c.getFilms)
	mux.HandleFunc("POST /cassettes", c.createCassettes)
	mux.HandleFunc("POST /order", c.createOrder)
	mux.HandleFunc("POST /login", c.login)
	mux.ServeHTTP(w, r)
}

func (c *Controller) decodeAndValidateBody(r *http.Request, d Validator) error {
	if err := json.NewDecoder(r.Body).Decode(d); err != nil {
		return fmt.Errorf("failed to decode request: %w", err)
	}

	if err := d.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	return nil
}

func (c *Controller) writeOK(w http.ResponseWriter, data any) {
	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (c *Controller) getAuthToken(r *http.Request) string {
	username, password, ok := r.BasicAuth()
	if !ok {
		return ""
	}

	if username != "Bearer" {
		return ""
	}

	return password
}
