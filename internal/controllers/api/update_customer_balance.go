package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type updateCustomerBalance struct {
	CustomerID int  `json:"customerId"`
	Balance    *int `json:"balance"`
}

func (n *updateCustomerBalance) Validate() error {
	var errs []error

	if n.CustomerID == 0 {
		errs = append(errs, errors.New("customer id required"))
	} else if n.CustomerID < 0 {
		errs = append(errs, errors.New("customer id must be posistive"))
	}

	if n.Balance == nil {
		errs = append(errs, errors.New("balance required"))
	} else if *n.Balance < 0 {
		errs = append(errs, errors.New("balance must be posistive"))
	}

	return errors.Join(errs...)
}

func (c *Controller) updateCustomerBalance(w http.ResponseWriter, r *http.Request) {
	var n updateCustomerBalance
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := n.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid customer balance: %s", err.Error()), http.StatusBadRequest)

		return
	}

	newBalance, err := c.CustomerBalanceUpdater(r.Context(), n.CustomerID, *n.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(newBalance); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
