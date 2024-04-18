package api

import (
	"errors"
	"net/http"
	"strings"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *loginRequest) Validate() error {
	var errs []error

	if l.Email = strings.TrimSpace(l.Email); l.Email == "" {
		errs = append(errs, errors.New("email required"))
	}

	if l.Password = strings.TrimSpace(l.Password); l.Password == "" {
		errs = append(errs, errors.New("password required"))
	}

	return errors.Join(errs...)
}

type loginResponse struct {
	Token string `json:"token"`
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	token := getAuthToken(r)
	if token != "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

		return
	}

	var n loginRequest
	if err := c.decodeAndValidateBody(r, &n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	token, err := c.Login(r.Context(), n.Email, n.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	c.writeOK(w, loginResponse{Token: token})
}
