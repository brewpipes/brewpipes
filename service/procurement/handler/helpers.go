package handler

import (
	"net/http"

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

func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
