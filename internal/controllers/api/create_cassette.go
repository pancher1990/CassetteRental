package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type newCassettes struct {
	FilmTitle string `json:"filmTitle"`
	Count     int    `json:"count"`
}

func (n *newCassettes) Validate() error {
	var errs []error

	if n.FilmTitle = strings.TrimSpace(n.FilmTitle); n.FilmTitle == "" {
		errs = append(errs, errors.New("film title required"))
	}

	if n.Count <= 0 {
		errs = append(errs, errors.New("cassette count must be positive"))
	}

	return errors.Join(errs...)
}

type cassettesGroup struct {
	IDs []int `json:"id"`
}

func (c *Controller) createCassette(w http.ResponseWriter, r *http.Request) {
	var n newCassettes
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := n.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid film: %s", err.Error()), http.StatusBadRequest)

		return
	}

	created, err := c.CassetteCreater(r.Context(), n.FilmTitle, n.Count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	var ids []int
	for _, cassette := range created {
		ids = append(ids, cassette.ID)
	}

	if err := json.NewEncoder(w).Encode(cassettesGroup{
		IDs: ids,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
