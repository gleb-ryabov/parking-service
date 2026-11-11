package parking

import (
	"context"
	"uuid"

	"github.com/rs/zerolog"

	"parking/internal/domain"
	svc "parking/internal/service"
)

// Repository defines data access operations for parking.
type Repository interface {
	// Create makes parking in DB.
	Create(ctx context.Context, parking domain.Parking) (*domain.Parking, error)
	// Update updates parking fields in DB.
	Update(ctx context.Context, parking domain.Parking) error
	// Delete removes parking from DB by ID.
	Delete(ctx context.Context, id uuid.UUID) error
	// Get returns parkings matching the specified filters.
	Get(ctx context.Context, code, address string) ([]domain.Parking, error)
}

// service provides business logic for parking.
type service struct {
	logger      *zerolog.Logger
	parkingRepo Repository
}

// New creates service for parking.
func New(
	logger *zerolog.Logger,
	parkingRepo Repository,
) *service {
	return &service{
		logger:      logger,
		parkingRepo: parkingRepo,
	}
}

// CreateParking creates a new parking.
func (s *service) CreateParking(ctx context.Context, parking domain.Parking) (*domain.Parking, error) {
	const op = "service.parking.CreateParking"

	parking.ID = uuid.New()

	log := s.logger.With().
		Str("op", op).
		Any("parking", parking).
		Logger()

	result, err := s.parkingRepo.Create(ctx, parking)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to create parking")
	}

	log.Debug().Str("parking_id", result.ID.String()).Msg("parking created")

	return result, nil
}

// UpdateParking updates parking fields.
func (s *service) UpdateParking(ctx context.Context, parking domain.Parking) error {
	const op = "service.parking.UpdateParking"

	log := s.logger.With().
		Str("op", op).
		Any("parking", parking).
		Logger()

	err := s.parkingRepo.Update(ctx, parking)
	if err != nil {
		return svc.MapError(err, &log, "failed to update parking")
	}

	log.Debug().Msg("parking updated")

	return nil
}

// DeleteParking removes parking by ID.
func (s *service) DeleteParking(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "service.parking.DeleteParking"

	log := s.logger.With().
		Str("op", op).
		Str("parking_id", id.String()).
		Logger()

	err := s.parkingRepo.Delete(ctx, id)
	if err != nil {
		return svc.MapError(err, &log, "failed to delete parking")
	}

	log.Debug().Msg("parking deleted")

	return nil
}

func (s *service) GetParkings(
	ctx context.Context,
	code, address string,
) ([]domain.Parking, error) {
	const op = "service.parking.GetParkings"

	log := s.logger.With().Str("op", op).Logger()

	parkings, err := s.parkingRepo.Get(ctx, code, address)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to get parkings")
	}

	log.Debug().Int("count", len(parkings)).Msg("parkings retrieved")

	return parkings, nil
}
