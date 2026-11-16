package domain

import (
	"time"
	"uuid"
)

// Schedule represents a parking schedule.
type Schedule struct {
	ID        uuid.UUID
	ParkingID uuid.UUID
	Active    bool

	ValidFrom time.Time  // Start date of the schedule validity period.
	ValidTo   *time.Time // End date of the schedule validity period.

	StartTime time.Time // Start time of the day.
	EndTime   time.Time // End time of the day.

	WeekDays []int
	Cost     int64
}
