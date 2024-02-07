package auth

import (
	"context"
	"net/http"
	"strings"
)

type Parser interface {
	Parse(accessToken string) (string, error)
}

func CheckToken(p Parser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")

			if token == "" {
				http.Error(w, "authorization token is missing", http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(token, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "error in header", http.StatusUnauthorized)

				return
			}

			if len(tokenParts[1]) == 0 {
				http.Error(w, "authorization token is missing", http.StatusUnauthorized)
				return
			}
			id, err := p.Parse(tokenParts[1])
			if err != nil {
				http.Error(w, "error with parsing authorization token", http.StatusUnauthorized)
				return

			}
			ctx := context.WithValue(r.Context(), "customerId", id)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
