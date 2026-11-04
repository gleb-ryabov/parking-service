package postgres

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Storage provides access to PostgreSQL database.
type Storage struct {
	DB     *gorm.DB
	logger *zerolog.Logger
}

// New creates a new PostgreSQL storage instance.
func New(ctx context.Context, logger *zerolog.Logger, dsn string) (*Storage, error) {
	const op = "storage.postgres.new"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open connect db - %s: %w", op, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("filed to get sql db - %s: %w", op, err)
	}

	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping db - %s: %w", op, err)
	}

	return &Storage{
		DB:     db,
		logger: logger,
	}, nil
}
