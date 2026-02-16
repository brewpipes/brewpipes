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
	CreateInventoryAdjustmentWithMovement(context.Context, storage.AdjustmentWithMovementRequest) (storage.AdjustmentWithMovementResult, error)
	GetInventoryAdjustmentByUUID(context.Context, string) (storage.InventoryAdjustment, error)
	ListInventoryAdjustments(context.Context) ([]storage.InventoryAdjustment, error)
	GetIngredientLotByUUID(context.Context, string) (storage.IngredientLot, error)
	GetBeerLotByUUID(context.Context, string) (storage.BeerLot, error)
	GetStockLocationByUUID(context.Context, string) (storage.StockLocation, error)
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

			// Resolve stock location UUID to internal ID.
			stockLocation, ok := service.ResolveFK(r.Context(), w, req.StockLocationUUID, "stock location", db.GetStockLocationByUUID)
			if !ok {
				return
			}

			// Resolve ingredient lot UUID to internal ID if provided.
			var ingredientLotID *int64
			if lot, ok := service.ResolveFKOptional(r.Context(), w, req.IngredientLotUUID, "ingredient lot", db.GetIngredientLotByUUID); !ok {
				return
			} else if req.IngredientLotUUID != nil {
				ingredientLotID = &lot.ID
			}

			// Resolve beer lot UUID to internal ID if provided.
			var beerLotID *int64
			if lot, ok := service.ResolveFKOptional(r.Context(), w, req.BeerLotUUID, "beer lot", db.GetBeerLotByUUID); !ok {
				return
			} else if req.BeerLotUUID != nil {
				beerLotID = &lot.ID
			}

			result, err := db.CreateInventoryAdjustmentWithMovement(r.Context(), storage.AdjustmentWithMovementRequest{
				IngredientLotID:   ingredientLotID,
				BeerLotID:         beerLotID,
				StockLocationID:   stockLocation.ID,
				Amount:            req.Amount,
				AmountUnit:        req.AmountUnit,
				Reason:            req.Reason,
				AdjustedAt:        adjustedAt,
				Notes:             req.Notes,
				IngredientLotUUID: req.IngredientLotUUID,
				BeerLotUUID:       req.BeerLotUUID,
				StockLocationUUID: req.StockLocationUUID,
			})
			if err != nil {
				service.InternalError(w, "error creating inventory adjustment", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewInventoryAdjustmentResponse(result.Adjustment))
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
