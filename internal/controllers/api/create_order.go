package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/pancher1990/cassette-rental/internal/usecases/orders"
)

type newOrder struct {
	CustomerID int `json:"customerId"`
	FilmID     int `json:"filmId"`
	RentDays   int `json:"rentDays"`
}

func (n *newOrder) Validate() error {
	var errs []error

	if n.FilmID <= 0 {
		errs = append(errs, errors.New("film id must be positive"))
	}

	if n.RentDays <= 0 {
		errs = append(errs, errors.New("rent days must be positive"))
	}

	return errors.Join(errs...)
}

type order struct {
	OrderID        int       `json:"orderId"`
	CassetteID     int       `json:"cassetteId"`
	CreatedAt      time.Time `json:"createdAt"`
	ReturnDeadline time.Time `json:"returnDeadline"`
	IsActive       bool      `json:"isActive"`
}

func (c *Controller) createOrder(w http.ResponseWriter, r *http.Request) {
	var n newOrder
	if err := c.decodeAndValidateBody(r, &n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	createdOrder, createdOrderCassette, err := c.CreateOrder(r.Context(), n.FilmID, n.RentDays, n.CustomerID)
	if err != nil {
		var status int
		switch err {
		case orders.ErrFilmNotFound, orders.ErrCustomerNotFound, orders.ErrCassetteNotFound:
			status = http.StatusNotFound
		case orders.RentPossibilityErrStatusConflict:
			status = http.StatusConflict
		default:
			status = http.StatusInternalServerError
		}
		http.Error(w, err.Error(), status)

		return
	}

	c.writeOK(w, order{
		OrderID:        createdOrderCassette.CassetteID,
		CassetteID:     createdOrder.ID,
		CreatedAt:      createdOrder.CreatedAt,
		ReturnDeadline: createdOrder.ReturnDeadline,
		IsActive:       createdOrder.IsActive,
	})
}
