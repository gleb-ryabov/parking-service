package gorm

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	repo "parking/internal/repository"
)

// MapError handles the GORM error and returns a new repository error.
func MapError(err error, op string) error {
	if err == nil {
		return nil
	}

	mappedErr := err

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		mappedErr = repo.ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		mappedErr = repo.ErrAlreadyExists
	case errors.Is(err, gorm.ErrForeignKeyViolated),
		errors.Is(err, gorm.ErrCheckConstraintViolated):
		mappedErr = repo.ErrConstraintViolated
	case errors.Is(err, gorm.ErrInvalidData):
		mappedErr = repo.ErrInvalidData
	case errors.Is(err, gorm.ErrInvalidTransaction):
		mappedErr = repo.ErrInvalidTransaction
	case errors.Is(err, gorm.ErrMissingWhereClause):
		mappedErr = repo.ErrMissingWhereClause
	case errors.Is(err, gorm.ErrUnsupportedRelation):
		mappedErr = repo.ErrUnsupportedRelation
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		mappedErr = repo.ErrPrimaryKeyRequired
	case errors.Is(err, gorm.ErrModelValueRequired):
		mappedErr = repo.ErrModelValueRequired
	case errors.Is(err, gorm.ErrNotImplemented):
		mappedErr = repo.ErrNotImplemented
	case errors.Is(err, gorm.ErrUnsupportedDriver):
		mappedErr = repo.ErrUnsupportedDriver
	}

	return fmt.Errorf("%s: %w", op, mappedErr)
}
