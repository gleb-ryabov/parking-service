package session

import (
	"context"
	"time"
	"uuid"

	"github.com/rs/zerolog"

	"parking/internal/domain"
	svc "parking/internal/service"
)

// Repository defines data access operations for session.
type Repository interface {
	// CreateSession makes parking session in DB.
	CreateSession(ctx context.Context, session domain.Session) (*domain.Session, error)
	// CloseSession sets the exit time for a parking session.
	CloseSession(
		ctx context.Context,
		stateNumber string,
		closedAt time.Time,
		parkingID uuid.UUID,
	) (*domain.Session, error)
	// GetSessionsByCarNumber returns parking sessions by the car state number.
	GetSessionsByCarNumber(
		ctx context.Context,
		stateNumber string,
		parkingID uuid.UUID,
	) ([]domain.Session, error)
	// GetActiveSessionsByParking returns active parking sessions by the session ID.
	GetActiveSessionsByParking(ctx context.Context, parkingID uuid.UUID) ([]domain.Session, error)
}

// service provides business logic for session.
type service struct {
	logger      *zerolog.Logger
	parkingRepo Repository
}

// New creates service for session.
func New(
	logger *zerolog.Logger,
	parkingRepo Repository,
) *service {
	return &service{
		logger:      logger,
		parkingRepo: parkingRepo,
	}
}

// OpenSession creates a parking session when a car enters the session.
func (s *service) OpenSession(
	ctx context.Context,
	stateNumber string,
	openedAt time.Time,
	parkingID uuid.UUID,
) (*domain.Session, error) {
	const op = "service.session.OpenSession"

	log := s.logger.With().
		Str("op", op).
		Str("state_number", stateNumber).
		Str("parking_id", parkingID.String()).
		Str("opened_at", openedAt.Format(time.RFC3339)).
		Logger()

	session := domain.Session{
		ID:          uuid.New(),
		ParkingID:   parkingID,
		StateNumber: stateNumber,
		EnteredAt:   openedAt,
	}

	result, err := s.parkingRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to create parking session")
	}

	log.Debug().Str("sessionID", result.ID.String()).Msg("parking session created")

	return result, nil
}

// CloseSession finishes a parking session when a car exits the session.
func (s *service) CloseSession(
	ctx context.Context,
	stateNumber string,
	closedAt time.Time,
	parkingID uuid.UUID,
) (*domain.Session, error) {
	const op = "service.session.CloseSession"

	log := s.logger.With().
		Str("op", op).
		Str("state_number", stateNumber).
		Str("parking_id", parkingID.String()).
		Str("closed_at", closedAt.Format(time.RFC3339)).
		Logger()

	session, err := s.parkingRepo.CloseSession(ctx, stateNumber, closedAt, parkingID)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to close parking session")
	}

	log.Debug().Str("session_id", session.ID.String()).Msg("parking session closed")

	return session, nil
}

// GetSessionsByCarNumber finds parking sessions by the car state number.
func (s *service) GetSessionsByCarNumber(
	ctx context.Context,
	stateNumber string,
	parkingID uuid.UUID,
) ([]domain.Session, error) {
	const op = "service.session.GetSessionsByCarNumber"

	log := s.logger.With().
		Str("op", op).
		Str("state_number", stateNumber).
		Str("parking_id", parkingID.String()).
		Logger()

	sessions, err := s.parkingRepo.GetSessionsByCarNumber(ctx, stateNumber, parkingID)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to get parking sessions")
	}

	log.Debug().Msg("parking sessions retrieved")

	return sessions, nil
}

// GetActiveSessionsByParking finds active parking sessions by the session ID.
func (s *service) GetActiveSessionsByParking(
	ctx context.Context,
	parkingID uuid.UUID,
) ([]domain.Session, error) {
	const op = "service.session.GetActiveSessionsByParking"

	log := s.logger.With().
		Str("op", op).
		Str("parking_id", parkingID.String()).
		Logger()

	sessions, err := s.parkingRepo.GetActiveSessionsByParking(ctx, parkingID)
	if err != nil {
		return nil, svc.MapError(err, &log, "failed to get parking sessions")
	}

	log.Debug().Msg("parking sessions retrieved")

	return sessions, nil
}

// CalculatePrice calculates a session price.
func (s *service) CalculatePrice(ctx context.Context, stateNumber string, parkingID uuid.UUID) (int64, error) {
	const op = "service.session.CalculatePrice"

	// TODO: implement after schedule manager implementation

	return 0, nil
}
