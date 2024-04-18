package api

import (
	"net/http"
)

func (c *Controller) logout(w http.ResponseWriter, r *http.Request) {
	token := getAuthToken(r)
	if token == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

		return
	}

	if err := c.Logout(r.Context(), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	c.writeOK(w, nil)
}
