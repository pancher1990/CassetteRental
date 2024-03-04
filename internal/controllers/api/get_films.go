package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type getFilms struct {
	Title string `json:"title"`
}

func (g *getFilms) Validate() error {
	g.Title = strings.TrimSpace(g.Title)

	return nil
}

func (c *Controller) getFilms(w http.ResponseWriter, r *http.Request) {
	var g getFilms
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := g.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid criteria: %s", err.Error()), http.StatusBadRequest)

		return
	}

	films, err := c.FilmFinder(r.Context(), g.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(films); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
