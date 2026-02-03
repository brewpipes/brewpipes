package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type VesselStore interface {
	CreateVessel(context.Context, storage.Vessel) (storage.Vessel, error)
	GetVessel(context.Context, int64) (storage.Vessel, error)
	GetVesselByUUID(context.Context, string) (storage.Vessel, error)
	ListVessels(context.Context) ([]storage.Vessel, error)
	UpdateVessel(context.Context, int64, storage.Vessel) (storage.Vessel, error)
}

type VesselOccupancyChecker interface {
	HasActiveOccupancy(context.Context, int64) (bool, error)
}

// HandleVessels handles [GET /vessels] and [POST /vessels].
func HandleVessels(db VesselStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			vessels, err := db.ListVessels(r.Context())
			if err != nil {
				slog.Error("error listing vessels", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewVesselsResponse(vessels))
		case http.MethodPost:
			var req dto.CreateVesselRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			status := ""
			if req.Status != nil {
				status = *req.Status
			}

			vessel := storage.Vessel{
				Type:         req.Type,
				Name:         req.Name,
				Capacity:     req.Capacity,
				CapacityUnit: req.CapacityUnit,
				Make:         req.Make,
				Model:        req.Model,
				Status:       status,
			}

			created, err := db.CreateVessel(r.Context(), vessel)
			if err != nil {
				slog.Error("error creating vessel", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewVesselResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// VesselByIDStore combines vessel storage and occupancy checking for the by-ID handler.
type VesselByIDStore interface {
	VesselStore
	VesselOccupancyChecker
}

// HandleVesselByID handles [GET /vessels/{id}] and [PATCH /vessels/{id}].
func HandleVesselByID(db VesselByIDStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		vesselID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			vessel, err := db.GetVessel(r.Context(), vesselID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "vessel not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting vessel", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewVesselResponse(vessel))
		case http.MethodPatch:
			var req dto.UpdateVesselRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Get current vessel to check status change
			currentVessel, err := db.GetVessel(r.Context(), vesselID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "vessel not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting vessel for update", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			// Check if status is changing to inactive or retired
			if currentVessel.Status == storage.VesselStatusActive &&
				(req.Status == storage.VesselStatusInactive || req.Status == storage.VesselStatusRetired) {
				hasActive, err := db.HasActiveOccupancy(r.Context(), vesselID)
				if err != nil {
					slog.Error("error checking active occupancy", "error", err, "vessel_id", vesselID)
					service.InternalError(w, err.Error())
					return
				}
				if hasActive {
					statusWord := "retire"
					if req.Status == storage.VesselStatusInactive {
						statusWord = "deactivate"
					}
					http.Error(w, fmt.Sprintf("cannot %s vessel with active occupancy", statusWord), http.StatusConflict)
					return
				}
			}

			vessel := storage.Vessel{
				Type:         req.Type,
				Name:         req.Name,
				Capacity:     req.Capacity,
				CapacityUnit: req.CapacityUnit,
				Make:         req.Make,
				Model:        req.Model,
				Status:       req.Status,
			}

			updated, err := db.UpdateVessel(r.Context(), vesselID, vessel)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "vessel not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating vessel", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("vessel updated", "vessel_id", updated.ID, "name", updated.Name)

			service.JSON(w, dto.NewVesselResponse(updated))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleVesselByUUID handles [GET /vessels/uuid/{uuid}].
func HandleVesselByUUID(db VesselStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		uuidValue := r.PathValue("uuid")
		if uuidValue == "" {
			http.Error(w, "uuid is required", http.StatusBadRequest)
			return
		}

		parsedUUID, err := parseUUIDParam(uuidValue)
		if err != nil {
			http.Error(w, "invalid uuid format", http.StatusBadRequest)
			return
		}

		vessel, err := db.GetVesselByUUID(r.Context(), parsedUUID.String())
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "vessel not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting vessel by uuid", "error", err, "uuid", uuidValue)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewVesselResponse(vessel))
	}
}
