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

type BatchStore interface {
	CreateBatch(context.Context, storage.Batch) (storage.Batch, error)
	GetBatch(context.Context, int64) (storage.Batch, error)
	ListBatches(context.Context) ([]storage.Batch, error)
}

// HandleBatches handles [GET /batches] and [POST /batches].
func HandleBatches(db BatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batches, err := db.ListBatches(r.Context())
			if err != nil {
				slog.Error("error listing batches", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBatchesResponse(batches))
		case http.MethodPost:
			slog.Info("create batch request", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)
			var req dto.CreateBatchRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				slog.Warn("invalid batch request body", "error", err)
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				slog.Warn("batch validation failed", "error", err, "short_name", req.ShortName)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			batch := storage.Batch{
				ShortName: req.ShortName,
				BrewDate:  req.BrewDate,
				Notes:     req.Notes,
			}

			created, err := db.CreateBatch(r.Context(), batch)
			if err != nil {
				slog.Error("error creating batch", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("batch created", "batch_id", created.ID, "short_name", created.ShortName)

			service.JSON(w, dto.NewBatchResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleBatchByID handles [GET /batches/{id}].
func HandleBatchByID(db BatchStore) http.HandlerFunc {
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
		batchID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		batch, err := db.GetBatch(r.Context(), batchID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting batch", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewBatchResponse(batch))
	}
}
