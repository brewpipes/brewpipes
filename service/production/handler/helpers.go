package handler

import (
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"
)

func parseInt64Param(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func parseUUIDParam(value string) (uuid.UUID, error) {
	return uuid.FromString(value)
}

func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
