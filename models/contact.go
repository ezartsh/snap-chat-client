package models

import "time"

type Contact struct {
	ID        *uint
	UserID    uint
	ContactID uint
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
