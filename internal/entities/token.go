package entities

import "github.com/golang-jwt/jwt/v5"

type TokenCustomClaims struct {
	jwt.RegisteredClaims
	UserType string `json:"user_type"`
}
