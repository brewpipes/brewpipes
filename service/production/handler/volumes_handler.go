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

type VolumeStore interface {
	ListVolumes(context.Context) ([]storage.Volume, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
	CreateVolume(context.Context, storage.Volume) (storage.Volume, error)
}

// HandleVolumes handles [GET /volumes] and [POST /volumes].
func HandleVolumes(db VolumeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			volumes, err := db.ListVolumes(r.Context())
			if err != nil {
				slog.Error("error listing volumes", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewVolumesResponse(volumes))
		case http.MethodPost:
			var req dto.CreateVolumeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			volume := storage.Volume{
				Name:        req.Name,
				Description: req.Description,
				Amount:      req.Amount,
				AmountUnit:  req.AmountUnit,
			}

			created, err := db.CreateVolume(r.Context(), volume)
			if err != nil {
				slog.Error("error creating volume", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewVolumeResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleVolumeByUUID handles [GET /volumes/{uuid}].
func HandleVolumeByUUID(db VolumeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		volumeUUID := r.PathValue("uuid")
		if volumeUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		volume, err := db.GetVolumeByUUID(r.Context(), volumeUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "volume not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting volume", "error", err, "volume_uuid", volumeUUID)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewVolumeResponse(volume))
	}
}
