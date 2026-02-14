package handler

import (
	"context"
	"encoding/json"
	"errors"
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
					service.InternalError(w, "error listing measurements by batch", "error", err)
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
					service.InternalError(w, "error listing measurements by occupancy", "error", err)
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
					service.InternalError(w, "error listing measurements by volume", "error", err)
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
			if batch, ok := service.ResolveFKOptional(r.Context(), w, req.BatchUUID, "batch", db.GetBatchByUUID); !ok {
				return
			} else if req.BatchUUID != nil {
				measurement.BatchID = &batch.ID
			}
			if occ, ok := service.ResolveFKOptional(r.Context(), w, req.OccupancyUUID, "occupancy", db.GetOccupancyByUUID); !ok {
				return
			} else if req.OccupancyUUID != nil {
				measurement.OccupancyID = &occ.ID
			}
			if vol, ok := service.ResolveFKOptional(r.Context(), w, req.VolumeUUID, "volume", db.GetVolumeByUUID); !ok {
				return
			} else if req.VolumeUUID != nil {
				measurement.VolumeID = &vol.ID
			}

			created, err := db.CreateMeasurement(r.Context(), measurement)
			if err != nil {
				service.InternalError(w, "error creating measurement", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewMeasurementResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleMeasurementByUUID handles [GET /measurements/{uuid}].
func HandleMeasurementByUUID(db MeasurementStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			service.InternalError(w, "error getting measurement", "error", err, "measurement_uuid", measurementUUID)
			return
		}

		service.JSON(w, dto.NewMeasurementResponse(measurement))
	}
}
