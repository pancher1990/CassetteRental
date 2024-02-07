package storage

import "errors"

var (
	ErrCustomerNotFound          = errors.New("customer not found")
	ErrCustomerInvalidCredential = errors.New("invalid credentials")
	ErrCassetteNotFound          = errors.New("cassette not found")
	ErrFilmNotFound              = errors.New("film not found")
)
