package entities

import "time"

type Order struct {
	ID             int
	CustomerId     int
	CreatedAt      time.Time
	ReturnDeadline time.Time
	IsActive       bool
}
