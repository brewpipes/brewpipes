package handler

import (
	"net/http"

	"github.com/gofrs/uuid/v5"
)

func parseUUIDParam(value string) (uuid.UUID, error) {
	return uuid.FromString(value)
}

func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
