package entities

import "time"

type Customer struct {
	ID        int
	CreatedAt time.Time
	Name      string
	IsActive  bool
	Balance   int
	Password  string
	Email     string
}
