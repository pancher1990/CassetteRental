package api

import (
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

type Validator interface {
	Validate() error
}

func (c *Controller) getFilms(w http.ResponseWriter, r *http.Request) {
	var g getFilms
	if err := c.decodeAndValidateBody(r, &g); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	films, err := c.FindFilms(r.Context(), g.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	c.writeOK(w, films)
}
