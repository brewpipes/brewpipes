package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type InventoryAdjustmentStore interface {
	CreateInventoryAdjustment(context.Context, storage.InventoryAdjustment) (storage.InventoryAdjustment, error)
	GetInventoryAdjustmentByUUID(context.Context, string) (storage.InventoryAdjustment, error)
	ListInventoryAdjustments(context.Context) ([]storage.InventoryAdjustment, error)
}

// HandleInventoryAdjustments handles [GET /inventory-adjustments] and [POST /inventory-adjustments].
func HandleInventoryAdjustments(db InventoryAdjustmentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			adjustments, err := db.ListInventoryAdjustments(r.Context())
			if err != nil {
				service.InternalError(w, "error listing inventory adjustments", "error", err)
				return
			}

			service.JSON(w, dto.NewInventoryAdjustmentsResponse(adjustments))
		case http.MethodPost:
			var req dto.CreateInventoryAdjustmentRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			adjustedAt := time.Time{}
			if req.AdjustedAt != nil {
				adjustedAt = *req.AdjustedAt
			}

			adjustment := storage.InventoryAdjustment{
				Reason:     req.Reason,
				AdjustedAt: adjustedAt,
				Notes:      req.Notes,
			}

			created, err := db.CreateInventoryAdjustment(r.Context(), adjustment)
			if err != nil {
				service.InternalError(w, "error creating inventory adjustment", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewInventoryAdjustmentResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleInventoryAdjustmentByUUID handles [GET /inventory-adjustments/{uuid}].
func HandleInventoryAdjustmentByUUID(db InventoryAdjustmentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adjustmentUUID := r.PathValue("uuid")
		if adjustmentUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		adjustment, err := db.GetInventoryAdjustmentByUUID(r.Context(), adjustmentUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "inventory adjustment not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting inventory adjustment", "error", err)
			return
		}

		service.JSON(w, dto.NewInventoryAdjustmentResponse(adjustment))
	}
}
