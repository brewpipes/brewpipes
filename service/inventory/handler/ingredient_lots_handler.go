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
	GetIngredientLotByUUID(context.Context, string) (storage.IngredientLot, error)
	GetIngredientByUUID(context.Context, string) (storage.Ingredient, error)
	GetInventoryReceiptByUUID(context.Context, string) (storage.InventoryReceipt, error)
	ListIngredientLots(context.Context) ([]storage.IngredientLot, error)
	ListIngredientLotsByIngredient(context.Context, string) ([]storage.IngredientLot, error)
	ListIngredientLotsByReceipt(context.Context, string) ([]storage.IngredientLot, error)
}

// HandleIngredientLots handles [GET /ingredient-lots] and [POST /ingredient-lots].
func HandleIngredientLots(db IngredientLotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_uuid")
			receiptValue := r.URL.Query().Get("receipt_uuid")
			if ingredientValue != "" && receiptValue != "" {
				http.Error(w, "ingredient_uuid and receipt_uuid cannot be combined", http.StatusBadRequest)
				return
			}
			if ingredientValue != "" {
				lots, err := db.ListIngredientLotsByIngredient(r.Context(), ingredientValue)
				if err != nil {
					slog.Error("error listing ingredient lots", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewIngredientLotsResponse(lots))
				return
			}
			if receiptValue != "" {
				lots, err := db.ListIngredientLotsByReceipt(r.Context(), receiptValue)
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

			// Resolve ingredient UUID to internal ID
			ingredient, err := db.GetIngredientByUUID(r.Context(), req.IngredientUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving ingredient uuid", "error", err)
				service.InternalError(w, err.Error())
				return
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

			lot := storage.IngredientLot{
				IngredientID:      ingredient.ID,
				ReceiptID:         receiptID,
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

// HandleIngredientLotByUUID handles [GET /ingredient-lots/{uuid}].
func HandleIngredientLotByUUID(db IngredientLotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		lotUUID := r.PathValue("uuid")
		if lotUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		lot, err := db.GetIngredientLotByUUID(r.Context(), lotUUID)
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
