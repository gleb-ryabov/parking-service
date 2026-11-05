package gorm

import (
	"context"
	"fmt"
	"time"
	"uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"parking/internal/domain"
	repo "parking/internal/repository"
	gormrepo "parking/internal/repository/gorm"
)

type repository struct {
	db *gorm.DB
}

// New creates repository for session.
func New(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

// CreateSession makes parking session in DB.
func (r *repository) CreateSession(ctx context.Context, session domain.Session) (*domain.Session, error) {
	const op = "repository.session.gorm.CreateSession"

	if err := r.db.WithContext(ctx).Create(&session).Error; err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	return &session, nil
}

// CloseSession sets the exit time for a parking session.
func (r *repository) CloseSession(
	ctx context.Context,
	stateNumber string,
	closedAt time.Time,
	parkingID uuid.UUID,
) (*domain.Session, error) {
	const op = "repository.session.gorm.CloseSession"

	var session domain.Session

	result := r.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Model(&session).
		Where("state_number = ? AND parking_id = ? AND exited_at IS NULL", stateNumber, parkingID).
		Updates(map[string]any{
			"exited_at": closedAt,
		})

	if err := result.Error; err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("%s: %w", op, repo.ErrNotFound)
	}

	return &session, nil
}

// GetSessionsByCarNumber returns parking sessions by the car state number.
func (r *repository) GetSessionsByCarNumber(
	ctx context.Context,
	stateNumber string,
	parkingID uuid.UUID,
) ([]domain.Session, error) {
	const op = "repository.session.gorm.GetSessionsByCarNumber"

	var sessions []domain.Session

	if err := r.db.WithContext(ctx).
		Where("state_number = ? AND parking_id = ?", stateNumber, parkingID).
		Find(&sessions).Error; err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	return sessions, nil
}

// GetActiveSessionsByParking returns active parking sessions by the session ID.
func (r *repository) GetActiveSessionsByParking(
	ctx context.Context,
	parkingID uuid.UUID,
) ([]domain.Session, error) {
	const op = "repository.session.gorm.GetActiveSessionsByParking"

	var sessions []domain.Session

	if err := r.db.WithContext(ctx).
		Where("parking_id = ? AND exited_at IS NULL", parkingID).
		Find(&sessions).Error; err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	return sessions, nil
}
