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
	"github.com/gofrs/uuid/v5"
)

type IngredientLotStore interface {
	CreateIngredientLot(context.Context, storage.IngredientLot) (storage.IngredientLot, error)
	GetIngredientLotByUUID(context.Context, string) (storage.IngredientLot, error)
	GetIngredientByUUID(context.Context, string) (storage.Ingredient, error)
	GetInventoryReceiptByUUID(context.Context, string) (storage.InventoryReceipt, error)
	ListIngredientLots(context.Context) ([]storage.IngredientLot, error)
	ListIngredientLotsIncludingDeleted(context.Context) ([]storage.IngredientLot, error)
	ListIngredientLotsByIngredient(context.Context, string) ([]storage.IngredientLot, error)
	ListIngredientLotsByReceipt(context.Context, string) ([]storage.IngredientLot, error)
	ListIngredientLotsByPurchaseOrderLineUUID(context.Context, string) ([]storage.IngredientLot, error)
}

// HandleIngredientLots handles [GET /ingredient-lots] and [POST /ingredient-lots].
func HandleIngredientLots(db IngredientLotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_uuid")
			receiptValue := r.URL.Query().Get("receipt_uuid")
			purchaseOrderLineValue := r.URL.Query().Get("purchase_order_line_uuid")

			// Count how many filter params are set
			filterCount := 0
			if ingredientValue != "" {
				filterCount++
			}
			if receiptValue != "" {
				filterCount++
			}
			if purchaseOrderLineValue != "" {
				filterCount++
			}
			if filterCount > 1 {
				http.Error(w, "only one filter parameter allowed: ingredient_uuid, receipt_uuid, or purchase_order_line_uuid", http.StatusBadRequest)
				return
			}

			if ingredientValue != "" {
				lots, err := db.ListIngredientLotsByIngredient(r.Context(), ingredientValue)
				if err != nil {
					service.InternalError(w, "error listing ingredient lots by ingredient", "error", err)
					return
				}

				service.JSON(w, dto.NewIngredientLotsResponse(lots))
				return
			}
			if receiptValue != "" {
				lots, err := db.ListIngredientLotsByReceipt(r.Context(), receiptValue)
				if err != nil {
					service.InternalError(w, "error listing ingredient lots by receipt", "error", err)
					return
				}

				service.JSON(w, dto.NewIngredientLotsResponse(lots))
				return
			}
			if purchaseOrderLineValue != "" {
				lots, err := db.ListIngredientLotsByPurchaseOrderLineUUID(r.Context(), purchaseOrderLineValue)
				if err != nil {
					service.InternalError(w, "error listing ingredient lots by purchase order line", "error", err)
					return
				}

				service.JSON(w, dto.NewIngredientLotsResponse(lots))
				return
			}

			includeDeleted := r.URL.Query().Get("include_deleted") == "true"
			var lots []storage.IngredientLot
			var err error
			if includeDeleted {
				lots, err = db.ListIngredientLotsIncludingDeleted(r.Context())
			} else {
				lots, err = db.ListIngredientLots(r.Context())
			}
			if err != nil {
				service.InternalError(w, "error listing ingredient lots", "error", err)
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

			var purchaseOrderLineUUID *uuid.UUID
			if req.PurchaseOrderLineUUID != nil {
				parsed, err := uuid.FromString(*req.PurchaseOrderLineUUID)
				if err != nil {
					http.Error(w, "invalid purchase_order_line_uuid", http.StatusBadRequest)
					return
				}
				purchaseOrderLineUUID = &parsed
			}

			// Resolve ingredient UUID to internal ID
			ingredient, ok := service.ResolveFK(r.Context(), w, req.IngredientUUID, "ingredient", db.GetIngredientByUUID)
			if !ok {
				return
			}

			// Resolve receipt UUID to internal ID if provided
			var receiptID *int64
			if receipt, ok := service.ResolveFKOptional(r.Context(), w, req.ReceiptUUID, "receipt", db.GetInventoryReceiptByUUID); !ok {
				return
			} else if req.ReceiptUUID != nil {
				receiptID = &receipt.ID
			}

			lot := storage.IngredientLot{
				IngredientID:          ingredient.ID,
				ReceiptID:             receiptID,
				SupplierUUID:          supplierUUID,
				PurchaseOrderLineUUID: purchaseOrderLineUUID,
				BreweryLotCode:        req.BreweryLotCode,
				OriginatorLotCode:     req.OriginatorLotCode,
				OriginatorName:        req.OriginatorName,
				OriginatorType:        req.OriginatorType,
				ReceivedAt:            receivedAt,
				ReceivedAmount:        req.ReceivedAmount,
				ReceivedUnit:          req.ReceivedUnit,
				BestByAt:              req.BestByAt,
				ExpiresAt:             req.ExpiresAt,
				Notes:                 req.Notes,
			}

			created, err := db.CreateIngredientLot(r.Context(), lot)
			if err != nil {
				service.InternalError(w, "error creating ingredient lot", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewIngredientLotResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleIngredientLotByUUID handles [GET /ingredient-lots/{uuid}].
func HandleIngredientLotByUUID(db IngredientLotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			service.InternalError(w, "error getting ingredient lot", "error", err)
			return
		}

		service.JSON(w, dto.NewIngredientLotResponse(lot))
	}
}
