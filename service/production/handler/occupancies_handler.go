package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type OccupancyStore interface {
	CreateOccupancy(context.Context, storage.Occupancy) (storage.Occupancy, error)
	GetOccupancyByUUID(context.Context, string) (storage.Occupancy, error)
	GetActiveOccupancyByVesselUUID(context.Context, string) (storage.Occupancy, error)
	GetActiveOccupancyByVolumeUUID(context.Context, string) (storage.Occupancy, error)
	ListActiveOccupancies(context.Context) ([]storage.Occupancy, error)
	UpdateOccupancyStatusByUUID(context.Context, string, *string) (storage.Occupancy, error)
	GetVesselByUUID(context.Context, string) (storage.Vessel, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
}

// HandleOccupancies handles [GET /occupancies].
func HandleOccupancies(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activeValue := r.URL.Query().Get("active")
		if activeValue == "" {
			http.Error(w, "active is required", http.StatusBadRequest)
			return
		}
		active, err := strconv.ParseBool(activeValue)
		if err != nil {
			http.Error(w, "invalid active", http.StatusBadRequest)
			return
		}
		if !active {
			http.Error(w, "active must be true", http.StatusBadRequest)
			return
		}

		occupancies, err := db.ListActiveOccupancies(r.Context())
		if err != nil {
			service.InternalError(w, "error listing active occupancies", "error", err)
			return
		}

		service.JSON(w, dto.NewOccupanciesResponse(occupancies))
	}
}

// HandleCreateOccupancy handles [POST /occupancies].
func HandleCreateOccupancy(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateOccupancyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Resolve vessel UUID to internal ID
		vessel, ok := service.ResolveFK(r.Context(), w, req.VesselUUID, "vessel", db.GetVesselByUUID)
		if !ok {
			return
		}

		// Resolve volume UUID to internal ID
		volume, ok := service.ResolveFK(r.Context(), w, req.VolumeUUID, "volume", db.GetVolumeByUUID)
		if !ok {
			return
		}

		inAt := time.Time{}
		if req.InAt != nil {
			inAt = *req.InAt
		}

		occupancy := storage.Occupancy{
			VesselID: vessel.ID,
			VolumeID: volume.ID,
			InAt:     inAt,
			Status:   req.Status,
		}

		created, err := db.CreateOccupancy(r.Context(), occupancy)
		if err != nil {
			service.InternalError(w, "error creating occupancy", "error", err)
			return
		}

		service.JSONCreated(w, dto.NewOccupancyResponse(created))
	}
}

// HandleOccupancyByUUID handles [GET /occupancies/{uuid}].
func HandleOccupancyByUUID(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		occupancyUUID := r.PathValue("uuid")
		if occupancyUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		occupancy, err := db.GetOccupancyByUUID(r.Context(), occupancyUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "occupancy not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting occupancy", "error", err, "occupancy_uuid", occupancyUUID)
			return
		}

		service.JSON(w, dto.NewOccupancyResponse(occupancy))
	}
}

// HandleActiveOccupancy handles [GET /occupancies/active].
func HandleActiveOccupancy(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vesselUUID := r.URL.Query().Get("active_vessel_uuid")
		volumeUUID := r.URL.Query().Get("active_volume_uuid")
		if vesselUUID != "" && volumeUUID != "" {
			http.Error(w, "active_vessel_uuid or active_volume_uuid is required", http.StatusBadRequest)
			return
		}
		if vesselUUID == "" && volumeUUID == "" {
			http.Error(w, "active_vessel_uuid or active_volume_uuid is required", http.StatusBadRequest)
			return
		}

		if vesselUUID != "" {
			occupancy, err := db.GetActiveOccupancyByVesselUUID(r.Context(), vesselUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "occupancy not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting active occupancy by vessel", "error", err, "vessel_uuid", vesselUUID)
				return
			}

			service.JSON(w, dto.NewOccupancyResponse(occupancy))
			return
		}

		occupancy, err := db.GetActiveOccupancyByVolumeUUID(r.Context(), volumeUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "occupancy not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting active occupancy by volume", "error", err, "volume_uuid", volumeUUID)
			return
		}

		service.JSON(w, dto.NewOccupancyResponse(occupancy))
	}
}

// HandleOccupancyStatus handles [PATCH /occupancies/{uuid}/status].
func HandleOccupancyStatus(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		occupancyUUID := r.PathValue("uuid")
		if occupancyUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		var req dto.UpdateOccupancyStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		occupancy, err := db.UpdateOccupancyStatusByUUID(r.Context(), occupancyUUID, req.Status)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "occupancy not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error updating occupancy status", "error", err, "occupancy_uuid", occupancyUUID)
			return
		}

		slog.Info("occupancy status updated", "occupancy_uuid", occupancyUUID, "status", req.Status)

		service.JSON(w, dto.NewOccupancyResponse(occupancy))
	}
}
