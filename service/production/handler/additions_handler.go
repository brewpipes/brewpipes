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
	"github.com/gofrs/uuid/v5"
)

type AdditionStore interface {
	CreateAddition(context.Context, storage.Addition) (storage.Addition, error)
	GetAdditionByUUID(context.Context, string) (storage.Addition, error)
	ListAdditionsByBatchUUID(context.Context, string) ([]storage.Addition, error)
	ListAdditionsByOccupancyUUID(context.Context, string) ([]storage.Addition, error)
	ListAdditionsByVolumeUUID(context.Context, string) ([]storage.Addition, error)
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
	GetOccupancyByUUID(context.Context, string) (storage.Occupancy, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
}

// HandleAdditions handles [GET /additions] and [POST /additions].
func HandleAdditions(db AdditionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if batchUUID := r.URL.Query().Get("batch_uuid"); batchUUID != "" {
				additions, err := db.ListAdditionsByBatchUUID(r.Context(), batchUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "batch not found", http.StatusNotFound)
					return
				} else if err != nil {
					service.InternalError(w, "error listing additions by batch", "error", err)
					return
				}

				service.JSON(w, dto.NewAdditionsResponse(additions))
				return
			}

			if occupancyUUID := r.URL.Query().Get("occupancy_uuid"); occupancyUUID != "" {
				additions, err := db.ListAdditionsByOccupancyUUID(r.Context(), occupancyUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "occupancy not found", http.StatusNotFound)
					return
				} else if err != nil {
					service.InternalError(w, "error listing additions by occupancy", "error", err)
					return
				}

				service.JSON(w, dto.NewAdditionsResponse(additions))
				return
			}

			if volumeUUID := r.URL.Query().Get("volume_uuid"); volumeUUID != "" {
				additions, err := db.ListAdditionsByVolumeUUID(r.Context(), volumeUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "volume not found", http.StatusNotFound)
					return
				} else if err != nil {
					service.InternalError(w, "error listing additions by volume", "error", err)
					return
				}

				service.JSON(w, dto.NewAdditionsResponse(additions))
				return
			}

			http.Error(w, "batch_uuid, occupancy_uuid, or volume_uuid is required", http.StatusBadRequest)
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
				AdditionType:     req.AdditionType,
				Stage:            req.Stage,
				InventoryLotUUID: inventoryUUID,
				Amount:           req.Amount,
				AmountUnit:       req.AmountUnit,
				AddedAt:          addedAt,
				Notes:            req.Notes,
			}

			// Resolve FK UUIDs to internal IDs
			if batch, ok := service.ResolveFKOptional(r.Context(), w, req.BatchUUID, "batch", db.GetBatchByUUID); !ok {
				return
			} else if req.BatchUUID != nil {
				addition.BatchID = &batch.ID
			}
			if occ, ok := service.ResolveFKOptional(r.Context(), w, req.OccupancyUUID, "occupancy", db.GetOccupancyByUUID); !ok {
				return
			} else if req.OccupancyUUID != nil {
				addition.OccupancyID = &occ.ID
			}
			if vol, ok := service.ResolveFKOptional(r.Context(), w, req.VolumeUUID, "volume", db.GetVolumeByUUID); !ok {
				return
			} else if req.VolumeUUID != nil {
				addition.VolumeID = &vol.ID
			}

			created, err := db.CreateAddition(r.Context(), addition)
			if err != nil {
				service.InternalError(w, "error creating addition", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewAdditionResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleAdditionByUUID handles [GET /additions/{uuid}].
func HandleAdditionByUUID(db AdditionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		additionUUID := r.PathValue("uuid")
		if additionUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		addition, err := db.GetAdditionByUUID(r.Context(), additionUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "addition not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting addition", "error", err, "addition_uuid", additionUUID)
			return
		}

		service.JSON(w, dto.NewAdditionResponse(addition))
	}
}
