package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type BatchVolumeStore interface {
	CreateBatchVolume(context.Context, storage.BatchVolume) (storage.BatchVolume, error)
	GetBatchVolume(context.Context, int64) (storage.BatchVolume, error)
	ListBatchVolumes(context.Context, int64) ([]storage.BatchVolume, error)
}

// HandleBatchVolumes handles [GET /batch-volumes] and [POST /batch-volumes].
func HandleBatchVolumes(db BatchVolumeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchValue := r.URL.Query().Get("batch_id")
			if batchValue == "" {
				http.Error(w, "batch_id is required", http.StatusBadRequest)
				return
			}
			batchID, err := parseInt64Param(batchValue)
			if err != nil {
				http.Error(w, "invalid batch_id", http.StatusBadRequest)
				return
			}

			volumes, err := db.ListBatchVolumes(r.Context(), batchID)
			if err != nil {
				slog.Error("error listing batch volumes", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBatchVolumesResponse(volumes))
		case http.MethodPost:
			var req dto.CreateBatchVolumeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			phaseAt := time.Time{}
			if req.PhaseAt != nil {
				phaseAt = *req.PhaseAt
			}

			batchVolume := storage.BatchVolume{
				BatchID:     req.BatchID,
				VolumeID:    req.VolumeID,
				LiquidPhase: req.LiquidPhase,
				PhaseAt:     phaseAt,
			}

			created, err := db.CreateBatchVolume(r.Context(), batchVolume)
			if err != nil {
				slog.Error("error creating batch volume", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBatchVolumeResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleBatchVolumeByID handles [GET /batch-volumes/{id}].
func HandleBatchVolumeByID(db BatchVolumeStore) http.HandlerFunc {
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
		batchVolumeID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		batchVolume, err := db.GetBatchVolume(r.Context(), batchVolumeID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch volume not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting batch volume", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewBatchVolumeResponse(batchVolume))
	}
}
