package gorm

import (
	"context"
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

// New creates repository for parking.
func New(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

// Create creates a new parking in DB.
func (r *repository) Create(ctx context.Context, parking domain.Parking) (*domain.Parking, error) {
	const op = "repository.parking.gorm.Create"

	if err := r.db.WithContext(ctx).Create(&parking).Error; err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	return &parking, nil
}

// Update updates parking fields in DB.
func (r *repository) Update(ctx context.Context, parking domain.Parking) error {
	const op = "repository.parking.gorm.Update"

	result := r.db.WithContext(ctx).
		Model(&parking).
		Where("id = ?", parking.ID).
		Updates(map[string]any{
			"code":     parking.Code,
			"address":  parking.Address,
			"capacity": parking.Capacity,
		})

	if err := result.Error; err != nil {
		return gormrepo.MapError(err, op)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, repo.ErrNotFound)
	}

	return nil
}

// Delete removes parking from DB by ID.
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.parking.gorm.Delete"

	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&domain.Parking{})

	if err := result.Error; err != nil {
		return gormrepo.MapError(err, op)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, repo.ErrNotFound)
	}

	return nil
}

// Get returns parkings matching the specified filters.
func (r *repository) Get(
	ctx context.Context,
	code, address string,
) ([]domain.Parking, error) {
	const op = "repository.parking.gorm.Get"

	var parkings []domain.Parking

	query := r.db.WithContext(ctx)

	if code != "" {
		query = query.Where("code = ?", code)
	}

	if address != "" {
		query = query.Where("address ILIKE ?", "%"+address+"%")
	}

	if err := query.Find(&parkings).Error; err != nil {
		return nil, gormrepo.MapError(err, op)
	}

	return parkings, nil
}
