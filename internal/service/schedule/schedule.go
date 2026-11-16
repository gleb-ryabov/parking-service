package schedule

import (
	"context"
	"uuid"

	"github.com/rs/zerolog"

	"parking/internal/domain"
	svc "parking/internal/service"
)

// Repository defines data access operations for schedule.
type Repository interface {
	// Create makes schedule in DB.
	Create(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error)
	// Update updates schedule fields in DB.
	Update(ctx context.Context, schedule domain.Schedule) error
	// Delete removes schedule from DB by ID.
	Delete(ctx context.Context, id uuid.UUID) error
	// GetByParkingID returns schedules for a specific parking with optional active filter.
	GetByParkingID(ctx context.Context, parkingID uuid.UUID, active *bool) ([]domain.Schedule, error)
	// CheckIntersections checks whether the schedule overlaps with any existing active schedule.
	CheckIntersections(ctx context.Context, schedule domain.Schedule) error
}

// service provides business logic for schedule.
type service struct {
	logger       *zerolog.Logger
	scheduleRepo Repository
}

// New creates service for schedule.
func New(
	logger *zerolog.Logger,
	scheduleRepo Repository,
) *service {
	return &service{
		logger:       logger,
		scheduleRepo: scheduleRepo,
	}
}

// CreateSchedule creates a new schedule.
func (s *service) CreateSchedule(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error) {
	const op = "service.schedule.CreateSchedule"

	schedule.ID = uuid.New()

	log := s.logger.With().
		Str("op", op).
		Any("schedule", schedule).
		Logger()

	err := s.scheduleRepo.CheckIntersections(ctx, schedule)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to check intersections")
	}

	result, err := s.scheduleRepo.Create(ctx, schedule)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to create schedule")
	}

	log.Debug().Str("schedule_id", result.ID.String()).Msg("schedule created")

	return result, nil
}

// UpdateSchedule updates schedule fields.
func (s *service) UpdateSchedule(ctx context.Context, schedule domain.Schedule) error {
	const op = "service.schedule.UpdateSchedule"

	log := s.logger.With().
		Str("op", op).
		Any("schedule", schedule).
		Logger()

	err := s.scheduleRepo.CheckIntersections(ctx, schedule)
	if err != nil {
		return svc.MapError(err, &log, "failed to check intersections")
	}

	err = s.scheduleRepo.Update(ctx, schedule)
	if err != nil {
		return svc.MapError(err, &log, "failed to update schedule")
	}

	log.Debug().Msg("schedule updated")

	return nil
}

// DeleteSchedule removes schedule by ID.
func (s *service) DeleteSchedule(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "service.schedule.DeleteSchedule"

	log := s.logger.With().
		Str("op", op).
		Str("schedule_id", id.String()).
		Logger()

	err := s.scheduleRepo.Delete(ctx, id)
	if err != nil {
		return svc.MapError(err, &log, "failed to delete schedule")
	}

	log.Debug().Msg("schedule deleted")

	return nil
}

// GetSchedulesByParkingID retrieves schedules for the specified parking,.
func (s *service) GetSchedulesByParkingID(
	ctx context.Context,
	parkingID uuid.UUID,
	active *bool,
) ([]domain.Schedule, error) {
	const op = "service.schedule.GetSchedulesByParkingID"

	log := s.logger.With().
		Str("op", op).
		Str("parking_id", parkingID.String()).
		Logger()

	schedules, err := s.scheduleRepo.GetByParkingID(ctx, parkingID, active)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to get schedules")
	}

	log.Debug().Int("count", len(schedules)).Msg("schedules retrieved")

	return schedules, nil
}
