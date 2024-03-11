package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pancher1990/cassette-rental/internal/repositories/cassettes"
	"github.com/pancher1990/cassette-rental/internal/repositories/customers"
	"github.com/pancher1990/cassette-rental/internal/repositories/films"
	"github.com/pancher1990/cassette-rental/internal/usecases/orders"
	"net/http"
	"time"
)

type newOrder struct {
	CustomerId int `json:"customerId"`
	FilmId     int `json:"filmId"`
	RentDays   int `json:"rentDays"`
}

func (n *newOrder) Validate() error {
	var errs []error

	if n.FilmId <= 0 {
		errs = append(errs, errors.New("film id must be positive"))
	}

	if n.RentDays <= 0 {
		errs = append(errs, errors.New("rent days must be positive"))
	}

	return errors.Join(errs...)
}

type order struct {
	ID             int       `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	ReturnDeadline time.Time `json:"returnDeadline"`
	IsActive       bool      `json:"isActive"`
}

func (c *Controller) createOrder(w http.ResponseWriter, r *http.Request) {
	var n newOrder
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := n.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid rent details: %s", err.Error()), http.StatusBadRequest)

		return
	}

	created, err := c.OrderCreater(r.Context(), n.FilmId, n.RentDays, n.CustomerId)
	if err != nil {
		var status int
		switch err {
		case films.ErrFilmNotFound, customers.ErrCustomerNotFound, cassettes.ErrCassetteNotFound:
			status = http.StatusNotFound
		case orders.RentPossibilityErrStatusConflict:
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
		http.Error(w, err.Error(), status)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(order{
		ID:             created.ID,
		CreatedAt:      created.CreatedAt,
		ReturnDeadline: created.ReturnDeadline,
		IsActive:       created.IsActive,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
