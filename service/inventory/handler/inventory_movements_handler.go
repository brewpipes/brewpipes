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
					slog.Error("error listing inventory movements by ingredient lot", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewInventoryMovementsResponse(movements))
				return
			}

			if beerValue != "" {
				movements, err := db.ListInventoryMovementsByBeerLot(r.Context(), beerValue)
				if err != nil {
					slog.Error("error listing inventory movements by beer lot", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewInventoryMovementsResponse(movements))
				return
			}

			movements, err := db.ListInventoryMovements(r.Context())
			if err != nil {
				slog.Error("error listing inventory movements", "error", err)
				service.InternalError(w, err.Error())
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
			stockLocation, err := db.GetStockLocationByUUID(r.Context(), req.StockLocationUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "stock location not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving stock location uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			// Resolve ingredient lot UUID to internal ID if provided
			var ingredientLotID *int64
			if req.IngredientLotUUID != nil {
				lot, err := db.GetIngredientLotByUUID(r.Context(), *req.IngredientLotUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "ingredient lot not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving ingredient lot uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				ingredientLotID = &lot.ID
			}

			// Resolve beer lot UUID to internal ID if provided
			var beerLotID *int64
			if req.BeerLotUUID != nil {
				lot, err := db.GetBeerLotByUUID(r.Context(), *req.BeerLotUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "beer lot not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving beer lot uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				beerLotID = &lot.ID
			}

			// Resolve receipt UUID to internal ID if provided
			var receiptID *int64
			if req.ReceiptUUID != nil {
				receipt, err := db.GetInventoryReceiptByUUID(r.Context(), *req.ReceiptUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "receipt not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving receipt uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				receiptID = &receipt.ID
			}

			// Resolve usage UUID to internal ID if provided
			var usageID *int64
			if req.UsageUUID != nil {
				usage, err := db.GetInventoryUsageByUUID(r.Context(), *req.UsageUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "usage not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving usage uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				usageID = &usage.ID
			}

			// Resolve adjustment UUID to internal ID if provided
			var adjustmentID *int64
			if req.AdjustmentUUID != nil {
				adjustment, err := db.GetInventoryAdjustmentByUUID(r.Context(), *req.AdjustmentUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "adjustment not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving adjustment uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
				adjustmentID = &adjustment.ID
			}

			// Resolve transfer UUID to internal ID if provided
			var transferID *int64
			if req.TransferUUID != nil {
				transfer, err := db.GetInventoryTransferByUUID(r.Context(), *req.TransferUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "transfer not found", http.StatusBadRequest)
					return
				} else if err != nil {
					slog.Error("error resolving transfer uuid", "error", err)
					service.InternalError(w, err.Error())
					return
				}
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
				slog.Error("error creating inventory movement", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewInventoryMovementResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleInventoryMovementByUUID handles [GET /inventory-movements/{uuid}].
func HandleInventoryMovementByUUID(db InventoryMovementStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

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
			slog.Error("error getting inventory movement", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewInventoryMovementResponse(movement))
	}
}
