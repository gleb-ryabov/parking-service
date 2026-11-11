package dto

import (
	"parking/internal/domain"
)

// ParkingRequest request for create and update parking.
type ParkingRequest struct {
	Code     string `json:"code" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Capacity int    `json:"capacity" binding:"required,min=1"`
}

// ToDomain converts DTO to domain model.
func (r ParkingRequest) ToDomain() domain.Parking {
	return domain.Parking{
		Code:     r.Code,
		Address:  r.Address,
		Capacity: r.Capacity,
	}
}
