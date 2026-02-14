package service

import (
	"context"
	"encoding/json"
	"errors"
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

// JSONCreated writes a 201 Created response with a JSON body.
func JSONCreated(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("error encoding JSON response payload after 201", "error", err)
	}
}

// InternalError accepts a string (not error) as its second parameter so callers
// provide a human-readable description for the log rather than forwarding raw
// errors, which could leak internal details into structured log messages.
func InternalError(w http.ResponseWriter, err string, logAttrs ...any) {
	cid := uuid.Must(uuid.NewV4())
	slog.Error(err, append([]any{"cid", cid}, logAttrs...)...)
	http.Error(w, fmt.Sprintf("server error (cid: %s)", cid), http.StatusInternalServerError)
}

// MethodNotAllowed writes a 405 Method Not Allowed response.
func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

var ErrNotFound = fmt.Errorf("not found")

// ResolveFK resolves a UUID to an entity using the provided getter function.
// Returns the entity and true on success, or writes an error response and returns false.
func ResolveFK[T any](ctx context.Context, w http.ResponseWriter, uuid string, label string, getter func(context.Context, string) (T, error)) (T, bool) {
	entity, err := getter(ctx, uuid)
	if errors.Is(err, ErrNotFound) {
		http.Error(w, label+" not found", http.StatusBadRequest)
		var zero T
		return zero, false
	} else if err != nil {
		InternalError(w, "error resolving "+label, "error", err)
		var zero T
		return zero, false
	}
	return entity, true
}

// ResolveFKOptional resolves an optional UUID pointer to an entity.
// If the pointer is nil, returns zero value and true (no-op).
// Otherwise behaves like ResolveFK.
func ResolveFKOptional[T any](ctx context.Context, w http.ResponseWriter, uuidPtr *string, label string, getter func(context.Context, string) (T, error)) (T, bool) {
	var zero T
	if uuidPtr == nil {
		return zero, true
	}
	return ResolveFK[T](ctx, w, *uuidPtr, label, getter)
}
