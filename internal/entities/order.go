package entities

import "time"

type Order struct {
	ID             int
	CustomerID     int
	CreatedAt      time.Time
	ReturnDeadline time.Time
	IsActive       bool
}
