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
	ListInventoryMovementsByIngredientLot(context.Context, int64) ([]storage.InventoryMovement, error)
	ListInventoryMovementsByBeerLot(context.Context, int64) ([]storage.InventoryMovement, error)
	GetInventoryMovement(context.Context, int64) (storage.InventoryMovement, error)
	CreateInventoryMovement(context.Context, storage.InventoryMovement) (storage.InventoryMovement, error)
}

// HandleInventoryMovements handles [GET /inventory-movements] and [POST /inventory-movements].
func HandleInventoryMovements(db InventoryMovementStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_lot_id")
			beerValue := r.URL.Query().Get("beer_lot_id")
			if ingredientValue != "" && beerValue != "" {
				http.Error(w, "ingredient_lot_id and beer_lot_id cannot both be set", http.StatusBadRequest)
				return
			}

			if ingredientValue != "" {
				lotID, err := parseInt64Param(ingredientValue)
				if err != nil {
					http.Error(w, "invalid ingredient_lot_id", http.StatusBadRequest)
					return
				}

				movements, err := db.ListInventoryMovementsByIngredientLot(r.Context(), lotID)
				if err != nil {
					slog.Error("error listing inventory movements by ingredient lot", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewInventoryMovementsResponse(movements))
				return
			}

			if beerValue != "" {
				lotID, err := parseInt64Param(beerValue)
				if err != nil {
					http.Error(w, "invalid beer_lot_id", http.StatusBadRequest)
					return
				}

				movements, err := db.ListInventoryMovementsByBeerLot(r.Context(), lotID)
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

			movement := storage.InventoryMovement{
				IngredientLotID: req.IngredientLotID,
				BeerLotID:       req.BeerLotID,
				StockLocationID: req.StockLocationID,
				Direction:       req.Direction,
				Reason:          req.Reason,
				Amount:          req.Amount,
				AmountUnit:      req.AmountUnit,
				OccurredAt:      occurredAt,
				ReceiptID:       req.ReceiptID,
				UsageID:         req.UsageID,
				AdjustmentID:    req.AdjustmentID,
				TransferID:      req.TransferID,
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

// HandleInventoryMovementByID handles [GET /inventory-movements/{id}].
func HandleInventoryMovementByID(db InventoryMovementStore) http.HandlerFunc {
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
		movementID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		movement, err := db.GetInventoryMovement(r.Context(), movementID)
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
