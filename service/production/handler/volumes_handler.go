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
	GetVolume(context.Context, int64) (storage.Volume, error)
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

// HandleVolumeByID handles [GET /volumes/{id}].
func HandleVolumeByID(db VolumeStore) http.HandlerFunc {
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
		volumeID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		volume, err := db.GetVolume(r.Context(), volumeID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "volume not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting volume", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewVolumeResponse(volume))
	}
}

// HandleGetVolumes is kept for backward compatibility.
func HandleGetVolumes(db VolumeStore) http.HandlerFunc {
	return HandleVolumes(db)
}
