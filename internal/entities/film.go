package entities

import "time"

type Film struct {
	ID        int
	CreatedAt time.Time
	Price     int
	Title     string
}
