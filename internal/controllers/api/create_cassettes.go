package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/usecases/cassettes"
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

type cassetteList []entities.Cassette

func (c cassetteList) ids() []int {
	ids := make([]int, 0, len(c))

	for _, cassette := range c {
		ids = append(ids, cassette.ID)
	}

	return ids
}

func (c *Controller) createCassettes(w http.ResponseWriter, r *http.Request) {
	var n newCassettes
	if err := c.decodeAndValidateBody(r, &n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	created, err := c.CreateCassettes(r.Context(), n.FilmTitle, n.Count)
	switch {
	case errors.Is(err, cassettes.ErrFilmNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	c.writeOK(w, cassettesGroup{IDs: cassetteList(created).ids()})
}
