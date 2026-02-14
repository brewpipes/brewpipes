package handler

import (
	"github.com/gofrs/uuid/v5"
)

func parseUUIDPointer(value *string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}

	parsed, err := uuid.FromString(*value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
