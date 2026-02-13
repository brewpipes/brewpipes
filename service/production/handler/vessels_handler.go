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
	GetVesselByUUID(context.Context, string) (storage.Vessel, error)
	ListVessels(context.Context) ([]storage.Vessel, error)
	UpdateVesselByUUID(context.Context, string, storage.Vessel) (storage.Vessel, error)
	HasActiveOccupancyByVesselUUID(context.Context, string) (bool, error)
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

// HandleVesselByUUID handles [GET /vessels/{uuid}] and [PATCH /vessels/{uuid}].
func HandleVesselByUUID(db VesselStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vesselUUID := r.PathValue("uuid")
		if vesselUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			vessel, err := db.GetVesselByUUID(r.Context(), vesselUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "vessel not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting vessel", "error", err, "vessel_uuid", vesselUUID)
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
			currentVessel, err := db.GetVesselByUUID(r.Context(), vesselUUID)
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
				hasActive, err := db.HasActiveOccupancyByVesselUUID(r.Context(), vesselUUID)
				if err != nil {
					slog.Error("error checking active occupancy", "error", err, "vessel_uuid", vesselUUID)
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

			updated, err := db.UpdateVesselByUUID(r.Context(), vesselUUID, vessel)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "vessel not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating vessel", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("vessel updated", "vessel_uuid", vesselUUID, "name", updated.Name)

			service.JSON(w, dto.NewVesselResponse(updated))
		default:
			methodNotAllowed(w)
		}
	}
}
