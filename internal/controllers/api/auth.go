package api

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"net/http"
	"os"
)

func AuthMiddleware(next http.HandlerFunc, userType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getAuthToken(r)
		tokenClaims, err := Parse(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}
		if (tokenClaims.UserType != userType) && (tokenClaims.UserType != "admin") {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		next(w, r)
	}
}

func Parse(accessToken string) (*entities.TokenCustomClaims, error) {
	sign := os.Getenv("SIGN")

	token, err := jwt.ParseWithClaims(accessToken, &entities.TokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}
		return []byte(sign), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*entities.TokenCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token or cannot convert claims")
	}
}
