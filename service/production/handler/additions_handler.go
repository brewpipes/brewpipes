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
	"github.com/gofrs/uuid/v5"
)

type AdditionStore interface {
	CreateAddition(context.Context, storage.Addition) (storage.Addition, error)
	GetAddition(context.Context, int64) (storage.Addition, error)
	ListAdditionsByBatch(context.Context, int64) ([]storage.Addition, error)
	ListAdditionsByOccupancy(context.Context, int64) ([]storage.Addition, error)
	ListAdditionsByVolume(context.Context, int64) ([]storage.Addition, error)
}

// HandleAdditions handles [GET /additions] and [POST /additions].
func HandleAdditions(db AdditionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if batchValue := r.URL.Query().Get("batch_id"); batchValue != "" {
				batchID, err := parseInt64Param(batchValue)
				if err != nil {
					http.Error(w, "invalid batch_id", http.StatusBadRequest)
					return
				}

				additions, err := db.ListAdditionsByBatch(r.Context(), batchID)
				if err != nil {
					slog.Error("error listing additions", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewAdditionsResponse(additions))
				return
			}

			if occupancyValue := r.URL.Query().Get("occupancy_id"); occupancyValue != "" {
				occupancyID, err := parseInt64Param(occupancyValue)
				if err != nil {
					http.Error(w, "invalid occupancy_id", http.StatusBadRequest)
					return
				}

				additions, err := db.ListAdditionsByOccupancy(r.Context(), occupancyID)
				if err != nil {
					slog.Error("error listing additions", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewAdditionsResponse(additions))
				return
			}

			if volumeValue := r.URL.Query().Get("volume_id"); volumeValue != "" {
				volumeID, err := parseInt64Param(volumeValue)
				if err != nil {
					http.Error(w, "invalid volume_id", http.StatusBadRequest)
					return
				}

				additions, err := db.ListAdditionsByVolume(r.Context(), volumeID)
				if err != nil {
					slog.Error("error listing additions", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewAdditionsResponse(additions))
				return
			}

			http.Error(w, "batch_id, occupancy_id, or volume_id is required", http.StatusBadRequest)
		case http.MethodPost:
			var req dto.CreateAdditionRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			addedAt := time.Time{}
			if req.AddedAt != nil {
				addedAt = *req.AddedAt
			}

			var inventoryUUID *uuid.UUID
			if req.InventoryLotUUID != nil {
				parsed, err := uuid.FromString(*req.InventoryLotUUID)
				if err != nil {
					http.Error(w, "invalid inventory_lot_uuid", http.StatusBadRequest)
					return
				}
				inventoryUUID = &parsed
			}

			addition := storage.Addition{
				BatchID:          req.BatchID,
				OccupancyID:      req.OccupancyID,
				VolumeID:         req.VolumeID,
				AdditionType:     req.AdditionType,
				Stage:            req.Stage,
				InventoryLotUUID: inventoryUUID,
				Amount:           req.Amount,
				AmountUnit:       req.AmountUnit,
				AddedAt:          addedAt,
				Notes:            req.Notes,
			}

			created, err := db.CreateAddition(r.Context(), addition)
			if err != nil {
				slog.Error("error creating addition", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewAdditionResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleAdditionByID handles [GET /additions/{id}].
func HandleAdditionByID(db AdditionStore) http.HandlerFunc {
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
		additionID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		addition, err := db.GetAddition(r.Context(), additionID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "addition not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting addition", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewAdditionResponse(addition))
	}
}
