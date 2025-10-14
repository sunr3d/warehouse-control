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

type ItemHistory struct {
	ID        int
	ItemID    int
	UserID    int
	Action    string
	OldValue  *string
	NewValue  *string
	ChangedAt time.Time
}
