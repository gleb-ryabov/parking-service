package repository

import "errors"

var (
	// ErrNotFound is returned when item not found in DB.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when item already exists in DB.
	ErrAlreadyExists = errors.New("already exists")
	// ErrConstraintViolated is returned when a database constraint is violated.
	ErrConstraintViolated = errors.New("constraint violated")
	// ErrInvalidData is returned when data is invalid for DB operation.
	ErrInvalidData = errors.New("invalid data")
	// ErrInvalidTransaction is returned when transaction is invalid.
	ErrInvalidTransaction = errors.New("invalid transaction")
	// ErrMissingWhereClause is returned when update or delete has no WHERE condition.
	ErrMissingWhereClause = errors.New("missing where clause")
	// ErrUnsupportedRelation is returned when relation is unsupported.
	ErrUnsupportedRelation = errors.New("unsupported relation")
	// ErrPrimaryKeyRequired is returned when primary key is required.
	ErrPrimaryKeyRequired = errors.New("primary key required")
	// ErrModelValueRequired is returned when model value is required.
	ErrModelValueRequired = errors.New("model value required")
	// ErrUnsupportedDriver is returned when database driver is unsupported.
	ErrUnsupportedDriver = errors.New("unsupported driver")
	// ErrNotImplemented is returned when operation is not implemented.
	ErrNotImplemented = errors.New("not implemented")
)
