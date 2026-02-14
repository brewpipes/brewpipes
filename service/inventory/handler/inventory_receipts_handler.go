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

type InventoryReceiptStore interface {
	CreateInventoryReceipt(context.Context, storage.InventoryReceipt) (storage.InventoryReceipt, error)
	GetInventoryReceiptByUUID(context.Context, string) (storage.InventoryReceipt, error)
	ListInventoryReceipts(context.Context) ([]storage.InventoryReceipt, error)
	ListInventoryReceiptsByPurchaseOrderUUID(context.Context, string) ([]storage.InventoryReceipt, error)
}

// HandleInventoryReceipts handles [GET /inventory-receipts] and [POST /inventory-receipts].
func HandleInventoryReceipts(db InventoryReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			purchaseOrderValue := r.URL.Query().Get("purchase_order_uuid")
			if purchaseOrderValue != "" {
				receipts, err := db.ListInventoryReceiptsByPurchaseOrderUUID(r.Context(), purchaseOrderValue)
				if err != nil {
					service.InternalError(w, "error listing inventory receipts by purchase order", "error", err)
					return
				}

				service.JSON(w, dto.NewInventoryReceiptsResponse(receipts))
				return
			}

			receipts, err := db.ListInventoryReceipts(r.Context())
			if err != nil {
				service.InternalError(w, "error listing inventory receipts", "error", err)
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

			var purchaseOrderUUID *uuid.UUID
			if req.PurchaseOrderUUID != nil {
				parsed, err := uuid.FromString(*req.PurchaseOrderUUID)
				if err != nil {
					http.Error(w, "invalid purchase_order_uuid", http.StatusBadRequest)
					return
				}
				purchaseOrderUUID = &parsed
			}

			receipt := storage.InventoryReceipt{
				SupplierUUID:      supplierUUID,
				PurchaseOrderUUID: purchaseOrderUUID,
				ReferenceCode:     req.ReferenceCode,
				ReceivedAt:        receivedAt,
				Notes:             req.Notes,
			}

			created, err := db.CreateInventoryReceipt(r.Context(), receipt)
			if err != nil {
				service.InternalError(w, "error creating inventory receipt", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewInventoryReceiptResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleInventoryReceiptByUUID handles [GET /inventory-receipts/{uuid}].
func HandleInventoryReceiptByUUID(db InventoryReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		receiptUUID := r.PathValue("uuid")
		if receiptUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		receipt, err := db.GetInventoryReceiptByUUID(r.Context(), receiptUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "inventory receipt not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting inventory receipt", "error", err)
			return
		}

		service.JSON(w, dto.NewInventoryReceiptResponse(receipt))
	}
}
