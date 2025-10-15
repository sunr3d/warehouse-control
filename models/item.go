package models

import "time"

type Item struct {
	ID          int
	Quantity    int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
