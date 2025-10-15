package models

import "time"

type ItemHistory struct {
	ID        int
	ItemID    int
	UserID    int
	Operation string
	OldValue  *string
	NewValue  *string
	ChangedAt time.Time
}
