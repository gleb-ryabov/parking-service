package parser

import (
	"errors"
	"fmt"
	"uuid"
)

// ParseBoolPtr parses "true" or "false" into a bool pointer.
// Returns nil if the input is empty.
func ParseBoolPtr(value string) (*bool, error) {
	switch value {
	case "":
		return nil, nil //nolint:nilnil
	case "true":
		v := true

		return &v, nil
	case "false":
		v := false

		return &v, nil
	default:
		return nil, fmt.Errorf("invalid boolean value %q", value)
	}
}

// ParseRequiredUUID parses a required UUID string.
func ParseRequiredUUID(value string) (uuid.UUID, error) {
	if value == "" {
		return uuid.Nil(), errors.New("value is required")
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil(), fmt.Errorf("invalid UUID: %w", err)
	}

	return id, nil
}
