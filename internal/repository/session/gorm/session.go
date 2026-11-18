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

// GetPriceByCarInParking returns parking cost by state number and parking ID.
//
//nolint:funlen
func (r *repository) GetPriceByCarInParking(
	ctx context.Context,
	stateNumber string,
	parkingID uuid.UUID,
) (int64, error) {
	const op = "repository.session.gorm.GetPriceByCarInParking"

	const query = `
        WITH prepared AS(
			SELECT
				ps.id AS session_id,
				sc.id AS schedule_id,
				sc.cost,
						SUM(
								EXTRACT(EPOCH FROM (
									LEAST(
											COALESCE(ps.exited_at, NOW()),
											d.day + sc.end_time
												+ CASE
													  WHEN sc.end_time <= sc.start_time
														  THEN INTERVAL '1 day'
													  ELSE INTERVAL '0'
												END
									)
										-
									GREATEST(
											ps.entered_at,
											d.day + sc.start_time
									)
									)) / 60
						) AS total_minutes
			FROM parking_sessions ps
					 JOIN schedules sc
						  ON sc.parking_id = ps.parking_id
					 CROSS JOIN LATERAL generate_series(
					GREATEST(ps.entered_at::date, sc.valid_from)::timestamp,
					LEAST(COALESCE(ps.exited_at, NOW())::date, sc.valid_to)::timestamp,
					INTERVAL '1 day'
										) AS d(day)
			WHERE
				EXTRACT(ISODOW FROM d.day)::integer = ANY(sc.week_days)
			  AND ps.entered_at < (
				d.day + sc.end_time
					+ CASE
						  WHEN sc.end_time <= sc.start_time
							  THEN INTERVAL '1 day'
						  ELSE INTERVAL '0'
					END
				)
			  AND COALESCE(ps.exited_at, NOW()) > (d.day + sc.start_time)
			  AND ps.paid = false
			  AND ps.parking_id = ?
			  AND ps.state_number = ?
			GROUP BY ps.id, sc.id)
			SELECT
				COALESCE(SUM(CEIL(total_minutes / 60) * cost), 0) AS price
			FROM prepared
    `
	var price int64

	if err := r.db.WithContext(ctx).
		Raw(query, parkingID, stateNumber).
		Scan(&price).Error; err != nil {
		return 0, gormrepo.MapError(err, op)
	}

	return price, nil
}
