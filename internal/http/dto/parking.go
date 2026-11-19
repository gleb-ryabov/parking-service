package dto

import (
	"parking/internal/domain"
)

// ParkingRequest request for create and update parking.
type ParkingRequest struct {
	Code     string `json:"code" binding:"required" example:"PARK-001"`
	Address  string `json:"address" binding:"required" example:"123 Main St, City"`
	Capacity int    `json:"capacity" binding:"required,min=1" example:"100"`
}

// ToDomain converts DTO to domain model.
func (r ParkingRequest) ToDomain() domain.Parking {
	return domain.Parking{
		Code:     r.Code,
		Address:  r.Address,
		Capacity: r.Capacity,
	}
}
