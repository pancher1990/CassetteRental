package sessions

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"os"
	"strconv"
	"time"

	"github.com/pancher1990/cassette-rental/internal/entities"
	"github.com/pancher1990/cassette-rental/internal/transaction"
)

type Repositories struct {
	CustomerRepository
	SessionRepository
}

type CustomerRepository interface {
	GetByEmailAndPassword(ctx context.Context, tx transaction.Querier, email, password string) (*entities.Customer, error)
}

type SessionRepository interface {
	Create(ctx context.Context, tx transaction.Querier, customerID int, token string, expireTime time.Time) error
	Remove(ctx context.Context, tx transaction.Querier, token string) error
}

func Login(repositories Repositories, tx transaction.TxFunc, expireTime time.Duration) func(ctx context.Context, email, password string) (string, error) {
	return func(ctx context.Context, email, password string) (string, error) {
		var tokenString string
		err := tx(
			ctx,
			func(tx pgx.Tx) (err error) {
				customer, err := repositories.CustomerRepository.GetByEmailAndPassword(ctx, tx, email, password)
				if err != nil {
					return fmt.Errorf("failed to create session: %w", err)
				}

				expiredIs := time.Now().Add(expireTime)
				tokenString, err = createToken(*customer, expiredIs)

				if err = repositories.SessionRepository.Create(ctx, tx, customer.ID, tokenString, expiredIs); err != nil {
					return fmt.Errorf("failed to create session: %w", err)
				}

				return nil
			},
		)

		return tokenString, err
	}
}

func createToken(c entities.Customer, expireTime time.Time) (string, error) {
	sign := os.Getenv("SIGN")

	var userType string
	if c.IsAdmin {
		userType = "admin"
	} else {
		userType = "customer"
	}

	claims := entities.TokenCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Subject:   strconv.Itoa(c.ID)},
		UserType: userType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(sign))
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil

}
