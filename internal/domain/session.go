package domain

import (
	"time"
	"uuid"
)

// Session represents a car's parking session.
type Session struct {
	ID          uuid.UUID
	ParkingID   uuid.UUID
	StateNumber string
	EnteredAt   time.Time
	ExitedAt    *time.Time
	Cost        *int64
	Paid        bool
}
