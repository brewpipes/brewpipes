package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type StyleStore interface {
	CreateStyle(context.Context, storage.Style) (storage.Style, error)
	GetStyleByUUID(context.Context, string) (storage.Style, error)
	ListStyles(context.Context) ([]storage.Style, error)
}

// HandleStyles handles [GET /styles] and [POST /styles].
func HandleStyles(db StyleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			styles, err := db.ListStyles(r.Context())
			if err != nil {
				service.InternalError(w, "error listing styles", "error", err)
				return
			}

			service.JSON(w, dto.NewStylesResponse(styles))
		case http.MethodPost:
			var req dto.CreateStyleRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			style := storage.Style{
				Name: req.Name,
			}

			created, err := db.CreateStyle(r.Context(), style)
			if err != nil {
				service.InternalError(w, "error creating style", "error", err)
				return
			}

			slog.Info("style created", "style_uuid", created.UUID, "name", created.Name)

			service.JSONCreated(w, dto.NewStyleResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleStyleByUUID handles [GET /styles/{uuid}].
func HandleStyleByUUID(db StyleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		styleUUID := r.PathValue("uuid")
		if styleUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		style, err := db.GetStyleByUUID(r.Context(), styleUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "style not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting style", "error", err, "style_uuid", styleUUID)
			return
		}

		service.JSON(w, dto.NewStyleResponse(style))
	}
}
