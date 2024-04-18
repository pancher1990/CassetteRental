package authorisation

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
	"os"
	"strconv"
)

type CustomerRepository interface {
	Get(context.Context, transaction.Querier, int) (*entities.Customer, error)
}

var ParseTokenError = errors.New("error with parse token")
var UserTypeError = errors.New("forbidden for this user type")

func Authorize(r CustomerRepository, tx transaction.TxFunc) func(ctx context.Context, token string, userType string) (*entities.Customer, error) {
	return func(ctx context.Context, token string, userType string) (*entities.Customer, error) {
		var c *entities.Customer
		tokenClaims, err := Parse(token)
		if err != nil {
			return c, fmt.Errorf("%w: %w", ParseTokenError, err)
		}

		if (tokenClaims.UserType != userType) && (tokenClaims.UserType != "admin") {
			return c, fmt.Errorf("%w: %w", UserTypeError, err)
		}
		userId, _ := strconv.Atoi(tokenClaims.Subject)
		err = tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				c, err = r.Get(ctx, tx, userId)
				if err != nil {
					return fmt.Errorf("failed to get customer: %w", err)
				}

				return nil
			},
		)

		return c, err

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
