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

type newFilm struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

func (n *newFilm) Validate() error {
	var errs []error

	if n.Title = strings.TrimSpace(n.Title); n.Title == "" {
		errs = append(errs, errors.New("title required"))
	}

	if n.Price <= 0 {
		errs = append(errs, errors.New("price must be positive"))
	}

	return errors.Join(errs...)
}

type film struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Price     int       `json:"price"`
	Title     string    `json:"title"`
}

func (c *Controller) createFilm(w http.ResponseWriter, r *http.Request) {
	var n newFilm
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := n.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid film: %s", err.Error()), http.StatusBadRequest)

		return
	}

	created, err := c.FilmCreater(r.Context(), entities.Film{
		Title: n.Title,
		Price: n.Price,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(film{
		ID:        created.ID,
		CreatedAt: created.CreatedAt,
		Title:     created.Title,
		Price:     created.Price,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
