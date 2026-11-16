package gorm

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"gorm.io/gorm"

	"parking/internal/domain"
	repo "parking/internal/repository"
	gormrepo "parking/internal/repository/gorm"
)

type repository struct {
	db *gorm.DB
}

// New creates repository for schedule.
func New(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

// Create creates a new schedule in DB.
func (r *repository) Create(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error) {
	const op = "repository.schedule.gorm.Create"

	m := toModel(schedule)

	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	result, err := toDomain(m)
	if err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	return &result, nil
}

// Update updates schedule fields in DB.
func (r *repository) Update(ctx context.Context, schedule domain.Schedule) error {
	const op = "repository.schedule.gorm.Update"

	m := toModel(schedule)

	result := r.db.WithContext(ctx).
		Model(&m).
		Where("id = ?", m.ID).
		Updates(map[string]any{
			"parking_id": m.ParkingID,
			"active":     m.Active,
			"valid_from": m.ValidFrom,
			"valid_to":   m.ValidTo,
			"start_time": m.StartTime,
			"end_time":   m.EndTime,
			"week_days":  m.WeekDays,
			"cost":       m.Cost,
		})

	if err := result.Error; err != nil {
		return gormrepo.MapError(err, op)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, repo.ErrNotFound)
	}

	return nil
}

// Delete removes schedule from DB by ID.
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.schedule.gorm.Delete"

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model{})

	if err := result.Error; err != nil {
		return gormrepo.MapError(err, op)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, repo.ErrNotFound)
	}

	return nil
}

// GetByParkingID returns schedules for a specific parking with optional active filter.
func (r *repository) GetByParkingID(ctx context.Context, parkingID uuid.UUID, active *bool) ([]domain.Schedule, error) {
	const op = "repository.schedule.gorm.GetByParkingID"

	var schedules []model

	query := r.db.WithContext(ctx).Where("parking_id = ?", parkingID)

	if active != nil {
		query = query.Where("active = ?", *active)
	}

	if err := query.Find(&schedules).Error; err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	result := make([]domain.Schedule, 0, len(schedules))
	for _, schedule := range schedules {
		s, err := toDomain(schedule)
		if err != nil {
			return nil, gormrepo.MapError(err, op)
		}

		result = append(result, s)
	}

	return result, nil
}

// CheckIntersections checks whether the schedule overlaps with any existing active schedule.
func (r *repository) CheckIntersections(ctx context.Context, schedule domain.Schedule) error {
	const op = "repository.schedule.gorm.CheckIntersections"

	m := toModel(schedule)

	filter := r.db.Model(&m).WithContext(ctx).
		Where("parking_id = ?", m.ParkingID).
		Where("active = true").
		Where("? = true", m.Active).
		Where("start_time <= ?", m.EndTime).
		Where("end_time >= ?", m.StartTime).
		Where("valid_from <= ?", m.ValidTo).
		Where("week_days && ?", m.WeekDays).
		Where("(valid_to IS NULL OR valid_to >= ?)", m.ValidFrom)

	if m.ValidTo != nil {
		filter = filter.Where("valid_from <= ?", *m.ValidTo)
	}

	var count int64
	if err := filter.Count(&count).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return gormrepo.MapError(err, op)
	}

	if count > 0 {
		return repo.ErrScheduleIntersection
	}

	return nil
}
