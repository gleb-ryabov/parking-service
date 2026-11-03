package service

import (
	"errors"
)

var (
	// ErrNotFound is returned when item not found.
	ErrNotFound = errors.New("not found")
	// ErrInternal is returned when not typed error.
	ErrInternal = errors.New("internal error")
)

// MapError handles the original error and returns a new service error.
func MapError(err error) error {
	switch {
	default:
		return ErrInternal
	}
}
