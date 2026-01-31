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
	GetMeasurement(context.Context, int64) (storage.Measurement, error)
	ListMeasurementsByBatch(context.Context, int64) ([]storage.Measurement, error)
	ListMeasurementsByOccupancy(context.Context, int64) ([]storage.Measurement, error)
}

// HandleMeasurements handles [GET /measurements] and [POST /measurements].
func HandleMeasurements(db MeasurementStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if batchValue := r.URL.Query().Get("batch_id"); batchValue != "" {
				batchID, err := parseInt64Param(batchValue)
				if err != nil {
					http.Error(w, "invalid batch_id", http.StatusBadRequest)
					return
				}

				measurements, err := db.ListMeasurementsByBatch(r.Context(), batchID)
				if err != nil {
					slog.Error("error listing measurements", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewMeasurementsResponse(measurements))
				return
			}

			if occupancyValue := r.URL.Query().Get("occupancy_id"); occupancyValue != "" {
				occupancyID, err := parseInt64Param(occupancyValue)
				if err != nil {
					http.Error(w, "invalid occupancy_id", http.StatusBadRequest)
					return
				}

				measurements, err := db.ListMeasurementsByOccupancy(r.Context(), occupancyID)
				if err != nil {
					slog.Error("error listing measurements", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewMeasurementsResponse(measurements))
				return
			}

			http.Error(w, "batch_id or occupancy_id is required", http.StatusBadRequest)
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
				BatchID:     req.BatchID,
				OccupancyID: req.OccupancyID,
				Kind:        req.Kind,
				Value:       req.Value,
				Unit:        req.Unit,
				ObservedAt:  observedAt,
				Notes:       req.Notes,
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

// HandleMeasurementByID handles [GET /measurements/{id}].
func HandleMeasurementByID(db MeasurementStore) http.HandlerFunc {
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
		measurementID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		measurement, err := db.GetMeasurement(r.Context(), measurementID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "measurement not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting measurement", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewMeasurementResponse(measurement))
	}
}
