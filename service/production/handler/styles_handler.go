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
	GetStyle(context.Context, int64) (storage.Style, error)
	ListStyles(context.Context) ([]storage.Style, error)
}

// HandleStyles handles [GET /styles] and [POST /styles].
func HandleStyles(db StyleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			styles, err := db.ListStyles(r.Context())
			if err != nil {
				slog.Error("error listing styles", "error", err)
				service.InternalError(w, err.Error())
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
				slog.Error("error creating style", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("style created", "style_id", created.ID, "name", created.Name)

			service.JSON(w, dto.NewStyleResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleStyleByID handles [GET /styles/{id}].
func HandleStyleByID(db StyleStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		styleID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		style, err := db.GetStyle(r.Context(), styleID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "style not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting style", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewStyleResponse(style))
	}
}
