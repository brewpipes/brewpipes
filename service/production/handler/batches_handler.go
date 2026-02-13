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
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
	ListBatches(context.Context) ([]storage.Batch, error)
	UpdateBatchByUUID(context.Context, string, storage.Batch) (storage.Batch, error)
	GetBatchDependenciesByUUID(context.Context, string) (storage.BatchDependencies, error)
	DeleteBatchByUUID(context.Context, string) error
	GetRecipe(context.Context, string, *storage.RecipeQueryOpts) (storage.Recipe, error)
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

			// Resolve recipe UUID to internal ID if provided
			if req.RecipeUUID != nil {
				recipe, err := db.GetRecipe(r.Context(), *req.RecipeUUID, nil)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "recipe not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving recipe uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				batch.RecipeID = &recipe.ID
			}

			created, err := db.CreateBatch(r.Context(), batch)
			if err != nil {
				slog.Error("error creating batch", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("batch created", "batch_uuid", created.UUID, "short_name", created.ShortName)

			service.JSON(w, dto.NewBatchResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleBatchByUUID handles [GET /batches/{uuid}], [PATCH /batches/{uuid}], and [DELETE /batches/{uuid}].
func HandleBatchByUUID(db BatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		batchUUID := r.PathValue("uuid")
		if batchUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			batch, err := db.GetBatchByUUID(r.Context(), batchUUID)
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
			}

			// Resolve recipe UUID to internal ID if provided
			if req.RecipeUUID != nil {
				recipe, err := db.GetRecipe(r.Context(), *req.RecipeUUID, nil)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "recipe not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving recipe uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				batch.RecipeID = &recipe.ID
			}

			updated, err := db.UpdateBatchByUUID(r.Context(), batchUUID, batch)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating batch", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("batch updated", "batch_uuid", batchUUID, "short_name", updated.ShortName)

			service.JSON(w, dto.NewBatchResponse(updated))
		case http.MethodDelete:
			// Log what will be deleted for audit purposes
			deps, err := db.GetBatchDependenciesByUUID(r.Context(), batchUUID)
			if err != nil {
				slog.Warn("could not check batch dependencies before delete", "batch_uuid", batchUUID, "error", err)
				// Continue with deletion anyway - this is just for logging
			} else if deps.HasDependencies() {
				slog.Info("deleting batch with related records",
					"batch_uuid", batchUUID,
					"batch_volumes", deps.BatchVolumeCount,
					"process_phases", deps.BatchProcessPhaseCount,
					"brew_sessions", deps.BrewSessionCount,
					"additions", deps.AdditionCount,
					"measurements", deps.MeasurementCount,
				)
			}

			// Cascade soft-delete the batch and all related records
			err = db.DeleteBatchByUUID(r.Context(), batchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error deleting batch", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("batch deleted", "batch_uuid", batchUUID)

			w.WriteHeader(http.StatusNoContent)
		default:
			methodNotAllowed(w)
		}
	}
}
