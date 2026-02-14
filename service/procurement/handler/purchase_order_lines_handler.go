package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/procurement/handler/dto"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
	"github.com/gofrs/uuid/v5"
)

type PurchaseOrderLineStore interface {
	ListPurchaseOrderLines(context.Context) ([]storage.PurchaseOrderLine, error)
	ListPurchaseOrderLinesByOrderUUID(context.Context, string) ([]storage.PurchaseOrderLine, error)
	GetPurchaseOrderLineByUUID(context.Context, string) (storage.PurchaseOrderLine, error)
	GetPurchaseOrderByUUID(context.Context, string) (storage.PurchaseOrder, error)
	CreatePurchaseOrderLine(context.Context, storage.PurchaseOrderLine) (storage.PurchaseOrderLine, error)
	UpdatePurchaseOrderLineByUUID(context.Context, string, storage.PurchaseOrderLineUpdate) (storage.PurchaseOrderLine, error)
	DeletePurchaseOrderLineByUUID(context.Context, string) (storage.PurchaseOrderLine, error)
}

// HandlePurchaseOrderLines handles [GET /purchase-order-lines] and [POST /purchase-order-lines].
func HandlePurchaseOrderLines(db PurchaseOrderLineStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderUUID := r.URL.Query().Get("purchase_order_uuid")
			if orderUUID != "" {
				lines, err := db.ListPurchaseOrderLinesByOrderUUID(r.Context(), orderUUID)
				if err != nil {
					service.InternalError(w, "error listing purchase order lines by order", "error", err)
					return
				}

				service.JSON(w, dto.NewPurchaseOrderLinesResponse(lines))
				return
			}

			lines, err := db.ListPurchaseOrderLines(r.Context())
			if err != nil {
				service.InternalError(w, "error listing purchase order lines", "error", err)
				return
			}

			service.JSON(w, dto.NewPurchaseOrderLinesResponse(lines))
		case http.MethodPost:
			var req dto.CreatePurchaseOrderLineRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Resolve purchase order UUID to internal ID
			order, ok := service.ResolveFK(r.Context(), w, req.PurchaseOrderUUID, "purchase order", db.GetPurchaseOrderByUUID)
			if !ok {
				return
			}

			inventoryItemUUID, err := parseUUIDPointer(req.InventoryItemUUID)
			if err != nil {
				http.Error(w, "invalid inventory_item_uuid", http.StatusBadRequest)
				return
			}

			line := storage.PurchaseOrderLine{
				PurchaseOrderID:   order.ID,
				LineNumber:        req.LineNumber,
				ItemType:          req.ItemType,
				ItemName:          req.ItemName,
				InventoryItemUUID: inventoryItemUUID,
				Quantity:          req.Quantity,
				QuantityUnit:      req.QuantityUnit,
				UnitCostCents:     req.UnitCostCents,
				Currency:          req.Currency,
			}

			created, err := db.CreatePurchaseOrderLine(r.Context(), line)
			if err != nil {
				service.InternalError(w, "error creating purchase order line", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewPurchaseOrderLineResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandlePurchaseOrderLineByUUID handles [GET /purchase-order-lines/{uuid}], [PATCH /purchase-order-lines/{uuid}], and [DELETE /purchase-order-lines/{uuid}].
func HandlePurchaseOrderLineByUUID(db PurchaseOrderLineStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lineUUID := r.PathValue("uuid")
		if lineUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			line, err := db.GetPurchaseOrderLineByUUID(r.Context(), lineUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order line not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting purchase order line", "error", err)
				return
			}

			service.JSON(w, dto.NewPurchaseOrderLineResponse(line))
		case http.MethodPatch:
			var req dto.UpdatePurchaseOrderLineRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var inventoryItemUUID *uuid.UUID
			if req.InventoryItemUUID != nil {
				parsed, err := parseUUIDPointer(req.InventoryItemUUID)
				if err != nil {
					http.Error(w, "invalid inventory_item_uuid", http.StatusBadRequest)
					return
				}
				inventoryItemUUID = parsed
			}

			update := storage.PurchaseOrderLineUpdate{
				LineNumber:        req.LineNumber,
				ItemType:          req.ItemType,
				ItemName:          req.ItemName,
				InventoryItemUUID: inventoryItemUUID,
				Quantity:          req.Quantity,
				QuantityUnit:      req.QuantityUnit,
				UnitCostCents:     req.UnitCostCents,
				Currency:          req.Currency,
			}

			line, err := db.UpdatePurchaseOrderLineByUUID(r.Context(), lineUUID, update)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order line not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error updating purchase order line", "error", err)
				return
			}

			service.JSON(w, dto.NewPurchaseOrderLineResponse(line))
		case http.MethodDelete:
			line, err := db.DeletePurchaseOrderLineByUUID(r.Context(), lineUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order line not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error deleting purchase order line", "error", err)
				return
			}

			service.JSON(w, dto.NewPurchaseOrderLineResponse(line))
		default:
			service.MethodNotAllowed(w)
		}
	}
}
