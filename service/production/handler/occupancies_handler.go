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
	GetOccupancy(context.Context, int64) (storage.Occupancy, error)
	GetActiveOccupancyByVessel(context.Context, int64) (storage.Occupancy, error)
	GetActiveOccupancyByVolume(context.Context, int64) (storage.Occupancy, error)
	ListActiveOccupancies(context.Context) ([]storage.Occupancy, error)
	UpdateOccupancyStatus(context.Context, int64, *string) (storage.Occupancy, error)
}

// HandleOccupancies handles [GET /occupancies].
func HandleOccupancies(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

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
			slog.Error("error listing active occupancies", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewOccupanciesResponse(occupancies))
	}
}

// HandleCreateOccupancy handles [POST /occupancies].
func HandleCreateOccupancy(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}

		var req dto.CreateOccupancyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		inAt := time.Time{}
		if req.InAt != nil {
			inAt = *req.InAt
		}

		occupancy := storage.Occupancy{
			VesselID: req.VesselID,
			VolumeID: req.VolumeID,
			InAt:     inAt,
			Status:   req.Status,
		}

		created, err := db.CreateOccupancy(r.Context(), occupancy)
		if err != nil {
			slog.Error("error creating occupancy", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewOccupancyResponse(created))
	}
}

// HandleOccupancyByID handles [GET /occupancies/{id}].
func HandleOccupancyByID(db OccupancyStore) http.HandlerFunc {
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
		occupancyID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		occupancy, err := db.GetOccupancy(r.Context(), occupancyID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "occupancy not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting occupancy", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewOccupancyResponse(occupancy))
	}
}

// HandleActiveOccupancy handles [GET /occupancies/active].
func HandleActiveOccupancy(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		vesselValue := r.URL.Query().Get("active_vessel_id")
		volumeValue := r.URL.Query().Get("active_volume_id")
		if vesselValue != "" && volumeValue != "" {
			http.Error(w, "active_vessel_id or active_volume_id is required", http.StatusBadRequest)
			return
		}
		if vesselValue == "" && volumeValue == "" {
			http.Error(w, "active_vessel_id or active_volume_id is required", http.StatusBadRequest)
			return
		}

		if vesselValue != "" {
			vesselID, err := parseInt64Param(vesselValue)
			if err != nil {
				http.Error(w, "invalid active_vessel_id", http.StatusBadRequest)
				return
			}

			occupancy, err := db.GetActiveOccupancyByVessel(r.Context(), vesselID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "occupancy not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting active occupancy by vessel", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewOccupancyResponse(occupancy))
			return
		}

		volumeID, err := parseInt64Param(volumeValue)
		if err != nil {
			http.Error(w, "invalid active_volume_id", http.StatusBadRequest)
			return
		}

		occupancy, err := db.GetActiveOccupancyByVolume(r.Context(), volumeID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "occupancy not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting active occupancy by volume", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewOccupancyResponse(occupancy))
	}
}

// HandleOccupancyStatus handles [PATCH /occupancies/{id}/status].
func HandleOccupancyStatus(db OccupancyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			methodNotAllowed(w)
			return
		}

		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		occupancyID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
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

		occupancy, err := db.UpdateOccupancyStatus(r.Context(), occupancyID, req.Status)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "occupancy not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error updating occupancy status", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		slog.Info("occupancy status updated", "occupancy_id", occupancyID, "status", req.Status)

		service.JSON(w, dto.NewOccupancyResponse(occupancy))
	}
}
