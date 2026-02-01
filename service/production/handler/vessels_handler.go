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

type VesselStore interface {
	CreateVessel(context.Context, storage.Vessel) (storage.Vessel, error)
	GetVessel(context.Context, int64) (storage.Vessel, error)
	GetVesselByUUID(context.Context, string) (storage.Vessel, error)
	ListVessels(context.Context) ([]storage.Vessel, error)
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

// HandleVesselByID handles [GET /vessels/{id}].
func HandleVesselByID(db VesselStore) http.HandlerFunc {
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
		vesselID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

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
