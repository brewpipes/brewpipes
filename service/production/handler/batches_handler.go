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
	UpdateBatch(context.Context, int64, storage.Batch) (storage.Batch, error)
	GetBatchDependencies(context.Context, int64) (storage.BatchDependencies, error)
	DeleteBatch(context.Context, int64) error
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
				RecipeID:  req.RecipeID,
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

// HandleBatchByID handles [GET /batches/{id}], [PATCH /batches/{id}], and [DELETE /batches/{id}].
func HandleBatchByID(db BatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		switch r.Method {
		case http.MethodGet:
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
		case http.MethodPatch:
			var req dto.UpdateBatchRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			batch := storage.Batch{
				ShortName: req.ShortName,
				BrewDate:  req.BrewDate,
				Notes:     req.Notes,
				RecipeID:  req.RecipeID,
			}

			updated, err := db.UpdateBatch(r.Context(), batchID, batch)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating batch", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("batch updated", "batch_id", updated.ID, "short_name", updated.ShortName)

			service.JSON(w, dto.NewBatchResponse(updated))
		case http.MethodDelete:
			// Log what will be deleted for audit purposes
			deps, err := db.GetBatchDependencies(r.Context(), batchID)
			if err != nil {
				slog.Warn("could not check batch dependencies before delete", "batch_id", batchID, "error", err)
				// Continue with deletion anyway - this is just for logging
			} else if deps.HasDependencies() {
				slog.Info("deleting batch with related records",
					"batch_id", batchID,
					"batch_volumes", deps.BatchVolumeCount,
					"process_phases", deps.BatchProcessPhaseCount,
					"brew_sessions", deps.BrewSessionCount,
					"additions", deps.AdditionCount,
					"measurements", deps.MeasurementCount,
				)
			}

			// Cascade soft-delete the batch and all related records
			err = db.DeleteBatch(r.Context(), batchID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error deleting batch", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("batch deleted", "batch_id", batchID)

			w.WriteHeader(http.StatusNoContent)
		default:
			methodNotAllowed(w)
		}
	}
}
