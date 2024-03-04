package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pancher1990/cassette-rental/internal/entities"
)

type newCustomer struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (n *newCustomer) Validate() error {
	var errs []error

	if n.Name = strings.TrimSpace(n.Name); n.Name == "" {
		errs = append(errs, errors.New("name required"))
	}

	if n.Password = strings.TrimSpace(n.Password); n.Password == "" {
		errs = append(errs, errors.New("password required"))
	}

	if n.Email = strings.TrimSpace(n.Email); n.Email == "" {
		errs = append(errs, errors.New("email required"))
	}

	return errors.Join(errs...)
}

type customer struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"isActive"`
	Balance   int       `json:"balance"`
}

func (c *Controller) createCustomer(w http.ResponseWriter, r *http.Request) {
	var n newCustomer
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := n.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid customer: %s", err.Error()), http.StatusBadRequest)

		return
	}

	created, err := c.CustomerCreater(r.Context(), entities.Customer{
		Name:     n.Name,
		Password: n.Password,
		Email:    n.Email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(customer{
		ID:        created.ID,
		CreatedAt: created.CreatedAt,
		Name:      created.Name,
		Email:     created.Email,
		IsActive:  created.IsActive,
		Balance:   created.Balance,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
