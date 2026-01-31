package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type InventoryAdjustmentStore interface {
	CreateInventoryAdjustment(context.Context, storage.InventoryAdjustment) (storage.InventoryAdjustment, error)
	GetInventoryAdjustment(context.Context, int64) (storage.InventoryAdjustment, error)
	ListInventoryAdjustments(context.Context) ([]storage.InventoryAdjustment, error)
}

// HandleInventoryAdjustments handles [GET /inventory-adjustments] and [POST /inventory-adjustments].
func HandleInventoryAdjustments(db InventoryAdjustmentStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			adjustments, err := db.ListInventoryAdjustments(r.Context())
			if err != nil {
				slog.Error("error listing inventory adjustments", "error", err)
				service.InternalError(w, err.Error())
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
				slog.Error("error creating inventory adjustment", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewInventoryAdjustmentResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleInventoryAdjustmentByID handles [GET /inventory-adjustments/{id}].
func HandleInventoryAdjustmentByID(db InventoryAdjustmentStore) http.HandlerFunc {
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
		adjustmentID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		adjustment, err := db.GetInventoryAdjustment(r.Context(), adjustmentID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "inventory adjustment not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting inventory adjustment", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewInventoryAdjustmentResponse(adjustment))
	}
}
