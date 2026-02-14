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

type InventoryMovementStore interface {
	ListInventoryMovements(context.Context) ([]storage.InventoryMovement, error)
	ListInventoryMovementsByIngredientLot(context.Context, string) ([]storage.InventoryMovement, error)
	ListInventoryMovementsByBeerLot(context.Context, string) ([]storage.InventoryMovement, error)
	GetInventoryMovementByUUID(context.Context, string) (storage.InventoryMovement, error)
	CreateInventoryMovement(context.Context, storage.InventoryMovement) (storage.InventoryMovement, error)
	GetIngredientLotByUUID(context.Context, string) (storage.IngredientLot, error)
	GetBeerLotByUUID(context.Context, string) (storage.BeerLot, error)
	GetStockLocationByUUID(context.Context, string) (storage.StockLocation, error)
	GetInventoryReceiptByUUID(context.Context, string) (storage.InventoryReceipt, error)
	GetInventoryUsageByUUID(context.Context, string) (storage.InventoryUsage, error)
	GetInventoryAdjustmentByUUID(context.Context, string) (storage.InventoryAdjustment, error)
	GetInventoryTransferByUUID(context.Context, string) (storage.InventoryTransfer, error)
}

// HandleInventoryMovements handles [GET /inventory-movements] and [POST /inventory-movements].
func HandleInventoryMovements(db InventoryMovementStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_lot_uuid")
			beerValue := r.URL.Query().Get("beer_lot_uuid")
			if ingredientValue != "" && beerValue != "" {
				http.Error(w, "ingredient_lot_uuid and beer_lot_uuid cannot both be set", http.StatusBadRequest)
				return
			}

			if ingredientValue != "" {
				movements, err := db.ListInventoryMovementsByIngredientLot(r.Context(), ingredientValue)
				if err != nil {
					service.InternalError(w, "error listing inventory movements by ingredient lot", "error", err)
					return
				}

				service.JSON(w, dto.NewInventoryMovementsResponse(movements))
				return
			}

			if beerValue != "" {
				movements, err := db.ListInventoryMovementsByBeerLot(r.Context(), beerValue)
				if err != nil {
					service.InternalError(w, "error listing inventory movements by beer lot", "error", err)
					return
				}

				service.JSON(w, dto.NewInventoryMovementsResponse(movements))
				return
			}

			movements, err := db.ListInventoryMovements(r.Context())
			if err != nil {
				service.InternalError(w, "error listing inventory movements", "error", err)
				return
			}

			service.JSON(w, dto.NewInventoryMovementsResponse(movements))
		case http.MethodPost:
			var req dto.CreateInventoryMovementRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			occurredAt := time.Time{}
			if req.OccurredAt != nil {
				occurredAt = *req.OccurredAt
			}

			// Resolve stock location UUID to internal ID
			stockLocation, ok := service.ResolveFK(r.Context(), w, req.StockLocationUUID, "stock location", db.GetStockLocationByUUID)
			if !ok {
				return
			}

			// Resolve ingredient lot UUID to internal ID if provided
			var ingredientLotID *int64
			if lot, ok := service.ResolveFKOptional(r.Context(), w, req.IngredientLotUUID, "ingredient lot", db.GetIngredientLotByUUID); !ok {
				return
			} else if req.IngredientLotUUID != nil {
				ingredientLotID = &lot.ID
			}

			// Resolve beer lot UUID to internal ID if provided
			var beerLotID *int64
			if lot, ok := service.ResolveFKOptional(r.Context(), w, req.BeerLotUUID, "beer lot", db.GetBeerLotByUUID); !ok {
				return
			} else if req.BeerLotUUID != nil {
				beerLotID = &lot.ID
			}

			// Resolve receipt UUID to internal ID if provided
			var receiptID *int64
			if receipt, ok := service.ResolveFKOptional(r.Context(), w, req.ReceiptUUID, "receipt", db.GetInventoryReceiptByUUID); !ok {
				return
			} else if req.ReceiptUUID != nil {
				receiptID = &receipt.ID
			}

			// Resolve usage UUID to internal ID if provided
			var usageID *int64
			if usage, ok := service.ResolveFKOptional(r.Context(), w, req.UsageUUID, "usage", db.GetInventoryUsageByUUID); !ok {
				return
			} else if req.UsageUUID != nil {
				usageID = &usage.ID
			}

			// Resolve adjustment UUID to internal ID if provided
			var adjustmentID *int64
			if adjustment, ok := service.ResolveFKOptional(r.Context(), w, req.AdjustmentUUID, "adjustment", db.GetInventoryAdjustmentByUUID); !ok {
				return
			} else if req.AdjustmentUUID != nil {
				adjustmentID = &adjustment.ID
			}

			// Resolve transfer UUID to internal ID if provided
			var transferID *int64
			if transfer, ok := service.ResolveFKOptional(r.Context(), w, req.TransferUUID, "transfer", db.GetInventoryTransferByUUID); !ok {
				return
			} else if req.TransferUUID != nil {
				transferID = &transfer.ID
			}

			movement := storage.InventoryMovement{
				IngredientLotID: ingredientLotID,
				BeerLotID:       beerLotID,
				StockLocationID: stockLocation.ID,
				Direction:       req.Direction,
				Reason:          req.Reason,
				Amount:          req.Amount,
				AmountUnit:      req.AmountUnit,
				OccurredAt:      occurredAt,
				ReceiptID:       receiptID,
				UsageID:         usageID,
				AdjustmentID:    adjustmentID,
				TransferID:      transferID,
				Notes:           req.Notes,
			}

			created, err := db.CreateInventoryMovement(r.Context(), movement)
			if err != nil {
				service.InternalError(w, "error creating inventory movement", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewInventoryMovementResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleInventoryMovementByUUID handles [GET /inventory-movements/{uuid}].
func HandleInventoryMovementByUUID(db InventoryMovementStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movementUUID := r.PathValue("uuid")
		if movementUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		movement, err := db.GetInventoryMovementByUUID(r.Context(), movementUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "inventory movement not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting inventory movement", "error", err)
			return
		}

		service.JSON(w, dto.NewInventoryMovementResponse(movement))
	}
}
