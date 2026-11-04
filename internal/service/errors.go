package service

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"parking/internal/repository"
)

var (
	// ErrNotFound is returned when an entity is not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when an entity already exists.
	ErrAlreadyExists = errors.New("already exists")
	// ErrConstraintViolated is returned when a database constraint is violated.
	ErrConstraintViolated = errors.New("constraint violated")
	// ErrInvalidData is returned when input data is invalid.
	ErrInvalidData = errors.New("invalid data")
	// ErrInternal is returned when an unexpected internal error occurs.
	ErrInternal = errors.New("internal error")
)

// MapError maps repository errors to service errors.
func MapError(err error, l *zerolog.Logger, msg string) error {
	if err == nil {
		return nil
	}

	var serviceErr error

	switch {
	case errors.Is(err, repository.ErrNotFound):
		serviceErr = ErrNotFound
	case errors.Is(err, repository.ErrAlreadyExists):
		serviceErr = ErrAlreadyExists
	case errors.Is(err, repository.ErrConstraintViolated):
		serviceErr = ErrConstraintViolated
	case errors.Is(err, repository.ErrInvalidData):
		serviceErr = ErrInvalidData
	default:
		serviceErr = ErrInternal
	}

	return fmt.Errorf("%w: %w", serviceErr, err)
}
