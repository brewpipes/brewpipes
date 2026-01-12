package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
)

type HTTPRoute struct {
	Method  string
	Path    string
	Handler http.Handler
}

func JSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		InternalError(w, "error encoding JSON response payload")
	}
}

func InternalError(w http.ResponseWriter, err string, logAttrs ...any) {
	cid := uuid.Must(uuid.NewV4())
	slog.Error(err, append([]any{"cid", cid}, logAttrs...)...)
	http.Error(w, fmt.Sprintf("server error (cid: %s)", cid), http.StatusInternalServerError)
}

var ErrNotFound = fmt.Errorf("not found")
