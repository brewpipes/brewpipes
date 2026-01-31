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
	"github.com/gofrs/uuid/v5"
)

type IngredientLotStore interface {
	CreateIngredientLot(context.Context, storage.IngredientLot) (storage.IngredientLot, error)
	GetIngredientLot(context.Context, int64) (storage.IngredientLot, error)
	ListIngredientLots(context.Context) ([]storage.IngredientLot, error)
	ListIngredientLotsByIngredient(context.Context, int64) ([]storage.IngredientLot, error)
	ListIngredientLotsByReceipt(context.Context, int64) ([]storage.IngredientLot, error)
}

// HandleIngredientLots handles [GET /ingredient-lots] and [POST /ingredient-lots].
func HandleIngredientLots(db IngredientLotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_id")
			receiptValue := r.URL.Query().Get("receipt_id")
			if ingredientValue != "" && receiptValue != "" {
				http.Error(w, "ingredient_id and receipt_id cannot be combined", http.StatusBadRequest)
				return
			}
			if ingredientValue != "" {
				ingredientID, err := parseInt64Param(ingredientValue)
				if err != nil {
					http.Error(w, "invalid ingredient_id", http.StatusBadRequest)
					return
				}

				lots, err := db.ListIngredientLotsByIngredient(r.Context(), ingredientID)
				if err != nil {
					slog.Error("error listing ingredient lots", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewIngredientLotsResponse(lots))
				return
			}
			if receiptValue != "" {
				receiptID, err := parseInt64Param(receiptValue)
				if err != nil {
					http.Error(w, "invalid receipt_id", http.StatusBadRequest)
					return
				}

				lots, err := db.ListIngredientLotsByReceipt(r.Context(), receiptID)
				if err != nil {
					slog.Error("error listing ingredient lots", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewIngredientLotsResponse(lots))
				return
			}

			lots, err := db.ListIngredientLots(r.Context())
			if err != nil {
				slog.Error("error listing ingredient lots", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientLotsResponse(lots))
		case http.MethodPost:
			var req dto.CreateIngredientLotRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			receivedAt := time.Time{}
			if req.ReceivedAt != nil {
				receivedAt = *req.ReceivedAt
			}

			var supplierUUID *uuid.UUID
			if req.SupplierUUID != nil {
				parsed, err := uuid.FromString(*req.SupplierUUID)
				if err != nil {
					http.Error(w, "invalid supplier_uuid", http.StatusBadRequest)
					return
				}
				supplierUUID = &parsed
			}

			lot := storage.IngredientLot{
				IngredientID:      req.IngredientID,
				ReceiptID:         req.ReceiptID,
				SupplierUUID:      supplierUUID,
				BreweryLotCode:    req.BreweryLotCode,
				OriginatorLotCode: req.OriginatorLotCode,
				OriginatorName:    req.OriginatorName,
				OriginatorType:    req.OriginatorType,
				ReceivedAt:        receivedAt,
				ReceivedAmount:    req.ReceivedAmount,
				ReceivedUnit:      req.ReceivedUnit,
				BestByAt:          req.BestByAt,
				ExpiresAt:         req.ExpiresAt,
				Notes:             req.Notes,
			}

			created, err := db.CreateIngredientLot(r.Context(), lot)
			if err != nil {
				slog.Error("error creating ingredient lot", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientLotResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleIngredientLotByID handles [GET /ingredient-lots/{id}].
func HandleIngredientLotByID(db IngredientLotStore) http.HandlerFunc {
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
		lotID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		lot, err := db.GetIngredientLot(r.Context(), lotID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient lot not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting ingredient lot", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewIngredientLotResponse(lot))
	}
}
