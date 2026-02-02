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

type InventoryReceiptStore interface {
	CreateInventoryReceipt(context.Context, storage.InventoryReceipt) (storage.InventoryReceipt, error)
	GetInventoryReceipt(context.Context, int64) (storage.InventoryReceipt, error)
	ListInventoryReceipts(context.Context) ([]storage.InventoryReceipt, error)
}

// HandleInventoryReceipts handles [GET /inventory-receipts] and [POST /inventory-receipts].
func HandleInventoryReceipts(db InventoryReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			receipts, err := db.ListInventoryReceipts(r.Context())
			if err != nil {
				slog.Error("error listing inventory receipts", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewInventoryReceiptsResponse(receipts))
		case http.MethodPost:
			var req dto.CreateInventoryReceiptRequest
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

			receipt := storage.InventoryReceipt{
				SupplierUUID:  supplierUUID,
				ReferenceCode: req.ReferenceCode,
				ReceivedAt:    receivedAt,
				Notes:         req.Notes,
			}

			created, err := db.CreateInventoryReceipt(r.Context(), receipt)
			if err != nil {
				slog.Error("error creating inventory receipt", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewInventoryReceiptResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleInventoryReceiptByID handles [GET /inventory-receipts/{id}].
func HandleInventoryReceiptByID(db InventoryReceiptStore) http.HandlerFunc {
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
		receiptID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		receipt, err := db.GetInventoryReceipt(r.Context(), receiptID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "inventory receipt not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting inventory receipt", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewInventoryReceiptResponse(receipt))
	}
}
