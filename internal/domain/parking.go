package domain

import (
	"uuid"
)

// Parking represents a session facility.
type Parking struct {
	ID       uuid.UUID
	Code     string
	Address  string
	Capacity int
}
