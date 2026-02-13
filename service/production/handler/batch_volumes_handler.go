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
	GetBatchVolumeByUUID(context.Context, string) (storage.BatchVolume, error)
	ListBatchVolumesByBatchUUID(context.Context, string) ([]storage.BatchVolume, error)
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
}

// HandleBatchVolumes handles [GET /batch-volumes] and [POST /batch-volumes].
func HandleBatchVolumes(db BatchVolumeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchUUID := r.URL.Query().Get("batch_uuid")
			if batchUUID == "" {
				http.Error(w, "batch_uuid is required", http.StatusBadRequest)
				return
			}

			volumes, err := db.ListBatchVolumesByBatchUUID(r.Context(), batchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
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

			// Resolve batch UUID to internal ID
			batch, err := db.GetBatchByUUID(r.Context(), req.BatchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving batch uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			// Resolve volume UUID to internal ID
			volume, err := db.GetVolumeByUUID(r.Context(), req.VolumeUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "volume not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving volume uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			phaseAt := time.Time{}
			if req.PhaseAt != nil {
				phaseAt = *req.PhaseAt
			}

			batchVolume := storage.BatchVolume{
				BatchID:     batch.ID,
				VolumeID:    volume.ID,
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

// HandleBatchVolumeByUUID handles [GET /batch-volumes/{uuid}].
func HandleBatchVolumeByUUID(db BatchVolumeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		bvUUID := r.PathValue("uuid")
		if bvUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		batchVolume, err := db.GetBatchVolumeByUUID(r.Context(), bvUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch volume not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting batch volume", "error", err, "batch_volume_uuid", bvUUID)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewBatchVolumeResponse(batchVolume))
	}
}
