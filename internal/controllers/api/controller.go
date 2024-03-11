package api

import (
	"context"
	"errors"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"net/http"
)

type (
	CustomerCreater        func(context.Context, entities.Customer) (*entities.Customer, error)
	CustomerBalanceUpdater func(ctx context.Context, customerID int, balance int) (resultBalance int, err error)
	FilmCreater            func(context.Context, entities.Film) (*entities.Film, error)
	FilmFinder             func(context.Context, string) ([]entities.Film, error)
	OrderCreater           func(context.Context, int, int, int) (*entities.Order, error)
)

type Controller struct {
	CustomerCreater
	CustomerBalanceUpdater
	FilmCreater
	FilmFinder
	OrderCreater
}

type option interface {
	apply(*Controller)
}

type optionFunc func(*Controller)

func (o optionFunc) apply(c *Controller) {
	o(c)
}

func WithCustomerCreater(customerCreater CustomerCreater) option {
	return optionFunc(func(c *Controller) {
		c.CustomerCreater = customerCreater
	})
}

func WithOrderCreater(orderCreater OrderCreater) option {
	return optionFunc(func(c *Controller) {
		c.OrderCreater = orderCreater
	})
}

func WithCustomerBalanceUpdater(customerBalanceUpdater CustomerBalanceUpdater) option {
	return optionFunc(func(c *Controller) {
		c.CustomerBalanceUpdater = customerBalanceUpdater
	})
}

func WithFilmCreater(filmCreater FilmCreater) option {
	return optionFunc(func(c *Controller) {
		c.FilmCreater = filmCreater
	})
}

func WithFilmFinder(filmFinder FilmFinder) option {
	return optionFunc(func(c *Controller) {
		c.FilmFinder = filmFinder
	})
}

func New(opts ...option) (*Controller, error) {
	var c Controller

	for _, o := range opts {
		o.apply(&c)
	}

	if c.CustomerCreater == nil {
		return nil, errors.New("customer creater required")
	}

	if c.CustomerBalanceUpdater == nil {
		return nil, errors.New("customer balance updater required")
	}

	if c.FilmCreater == nil {
		return nil, errors.New("film creater required")
	}

	if c.FilmFinder == nil {
		return nil, errors.New("film finder required")
	}

	return &c, nil
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/customers", c.customers)
	mux.HandleFunc("/customer/balance", c.customerBalance)
	mux.HandleFunc("/films", c.films)
	mux.HandleFunc("/order", c.orders)
	mux.ServeHTTP(w, r)
}

func (c *Controller) customers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.createCustomer(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (c *Controller) customerBalance(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		c.updateCustomerBalance(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (c *Controller) films(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.createFilm(w, r)
	case http.MethodGet:
		c.getFilms(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (c *Controller) orders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.createOrder(w, r)
	default:
		http.NotFound(w, r)
	}
}
