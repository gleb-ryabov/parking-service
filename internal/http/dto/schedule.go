package dto

import (
	"fmt"
	"time"
	"uuid"

	"parking/internal/domain"
	handl "parking/internal/http/handler"
)

// ScheduleRequest request for create and update schedule.
type ScheduleRequest struct {
	ParkingID string  `json:"parking_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Active    bool    `json:"active" example:"true"`
	ValidFrom string  `json:"valid_from" binding:"required" example:"2024-01-01"`
	ValidTo   *string `json:"valid_to" example:"2024-12-31"`
	StartTime string  `json:"start_time" binding:"required" example:"08:00:00"`
	EndTime   string  `json:"end_time" binding:"required" example:"20:00:00"`
	WeekDays  []int   `json:"week_days" binding:"required,min=1,max=7,dive,min=0,max=6" example:"1,2,3,4,5"`
	Cost      int64   `json:"cost" binding:"required,min=0" example:"50000"`
}

// ToDomain converts DTO to domain model.
func (r ScheduleRequest) ToDomain() (domain.Schedule, error) {
	validFrom, err := time.Parse("2006-01-02", r.ValidFrom)
	if err != nil {
		return domain.Schedule{}, fmt.Errorf("%w: %s", handl.ErrInvalidDate, "valid_from")
	}

	var validTo *time.Time
	if r.ValidTo != nil {
		var t time.Time
		t, err = time.Parse("2006-01-02", *r.ValidTo)
		if err != nil {
			return domain.Schedule{}, fmt.Errorf("%w: %s", handl.ErrInvalidDate, "valid_to")
		}
		validTo = &t
	}

	startTime, err := time.Parse("15:04:05", r.StartTime)
	if err != nil {
		return domain.Schedule{}, fmt.Errorf("%w: %s", handl.ErrInvalidTime, "start_time")
	}

	endTime, err := time.Parse("15:04:05", r.EndTime)
	if err != nil {
		return domain.Schedule{}, fmt.Errorf("%w: %s", handl.ErrInvalidTime, "end_time")
	}

	parkingID, err := uuid.Parse(r.ParkingID)
	if err != nil {
		return domain.Schedule{}, fmt.Errorf("%w: %s", handl.ErrInvalidUUID, "parking_id")
	}

	return domain.Schedule{
		ParkingID: parkingID,
		Active:    r.Active,
		ValidFrom: validFrom,
		ValidTo:   validTo,
		StartTime: startTime,
		EndTime:   endTime,
		WeekDays:  r.WeekDays,
		Cost:      r.Cost,
	}, nil
}
