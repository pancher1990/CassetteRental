package storage

import "errors"

var (
	// ErrURLNotFound TODO изменить ошибки

	ErrCustomerNotFound = errors.New("customer not found")
	ErrCassetteNotFound = errors.New("film not found")
	ErrFilmNotFound     = errors.New("film not found")
	ErrURLExists        = errors.New("url exists")
)
