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

type MeasurementStore interface {
	CreateMeasurement(context.Context, storage.Measurement) (storage.Measurement, error)
	GetMeasurementByUUID(context.Context, string) (storage.Measurement, error)
	ListMeasurementsByBatchUUID(context.Context, string) ([]storage.Measurement, error)
	ListMeasurementsByOccupancyUUID(context.Context, string) ([]storage.Measurement, error)
	ListMeasurementsByVolumeUUID(context.Context, string) ([]storage.Measurement, error)
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
	GetOccupancyByUUID(context.Context, string) (storage.Occupancy, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
}

// HandleMeasurements handles [GET /measurements] and [POST /measurements].
func HandleMeasurements(db MeasurementStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if batchUUID := r.URL.Query().Get("batch_uuid"); batchUUID != "" {
				measurements, err := db.ListMeasurementsByBatchUUID(r.Context(), batchUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "batch not found", http.StatusNotFound)
					return
				} else if err != nil {
					slog.Error("error listing measurements", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewMeasurementsResponse(measurements))
				return
			}

			if occupancyUUID := r.URL.Query().Get("occupancy_uuid"); occupancyUUID != "" {
				measurements, err := db.ListMeasurementsByOccupancyUUID(r.Context(), occupancyUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "occupancy not found", http.StatusNotFound)
					return
				} else if err != nil {
					slog.Error("error listing measurements", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewMeasurementsResponse(measurements))
				return
			}

			if volumeUUID := r.URL.Query().Get("volume_uuid"); volumeUUID != "" {
				measurements, err := db.ListMeasurementsByVolumeUUID(r.Context(), volumeUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "volume not found", http.StatusNotFound)
					return
				} else if err != nil {
					slog.Error("error listing measurements", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewMeasurementsResponse(measurements))
				return
			}

			http.Error(w, "batch_uuid, occupancy_uuid, or volume_uuid is required", http.StatusBadRequest)
		case http.MethodPost:
			var req dto.CreateMeasurementRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			observedAt := time.Time{}
			if req.ObservedAt != nil {
				observedAt = *req.ObservedAt
			}

			measurement := storage.Measurement{
				Kind:       req.Kind,
				Value:      req.Value,
				Unit:       req.Unit,
				ObservedAt: observedAt,
				Notes:      req.Notes,
			}

			// Resolve FK UUIDs to internal IDs
			if req.BatchUUID != nil {
				batch, err := db.GetBatchByUUID(r.Context(), *req.BatchUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "batch not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving batch uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				measurement.BatchID = &batch.ID
			}
			if req.OccupancyUUID != nil {
				occ, err := db.GetOccupancyByUUID(r.Context(), *req.OccupancyUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "occupancy not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving occupancy uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				measurement.OccupancyID = &occ.ID
			}
			if req.VolumeUUID != nil {
				vol, err := db.GetVolumeByUUID(r.Context(), *req.VolumeUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "volume not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving volume uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				measurement.VolumeID = &vol.ID
			}

			created, err := db.CreateMeasurement(r.Context(), measurement)
			if err != nil {
				slog.Error("error creating measurement", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewMeasurementResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleMeasurementByUUID handles [GET /measurements/{uuid}].
func HandleMeasurementByUUID(db MeasurementStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		measurementUUID := r.PathValue("uuid")
		if measurementUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		measurement, err := db.GetMeasurementByUUID(r.Context(), measurementUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "measurement not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting measurement", "error", err, "measurement_uuid", measurementUUID)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewMeasurementResponse(measurement))
	}
}
