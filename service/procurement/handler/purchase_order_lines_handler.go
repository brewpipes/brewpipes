package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/procurement/handler/dto"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
	"github.com/gofrs/uuid/v5"
)

type PurchaseOrderLineStore interface {
	ListPurchaseOrderLines(context.Context) ([]storage.PurchaseOrderLine, error)
	ListPurchaseOrderLinesByOrder(context.Context, int64) ([]storage.PurchaseOrderLine, error)
	GetPurchaseOrderLine(context.Context, int64) (storage.PurchaseOrderLine, error)
	CreatePurchaseOrderLine(context.Context, storage.PurchaseOrderLine) (storage.PurchaseOrderLine, error)
	UpdatePurchaseOrderLine(context.Context, int64, storage.PurchaseOrderLineUpdate) (storage.PurchaseOrderLine, error)
	DeletePurchaseOrderLine(context.Context, int64) (storage.PurchaseOrderLine, error)
}

// HandlePurchaseOrderLines handles [GET /purchase-order-lines] and [POST /purchase-order-lines].
func HandlePurchaseOrderLines(db PurchaseOrderLineStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderValue := r.URL.Query().Get("purchase_order_id")
			if orderValue != "" {
				orderID, err := parseInt64Param(orderValue)
				if err != nil {
					http.Error(w, "invalid purchase_order_id", http.StatusBadRequest)
					return
				}

				lines, err := db.ListPurchaseOrderLinesByOrder(r.Context(), orderID)
				if err != nil {
					slog.Error("error listing purchase order lines by order", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewPurchaseOrderLinesResponse(lines))
				return
			}

			lines, err := db.ListPurchaseOrderLines(r.Context())
			if err != nil {
				slog.Error("error listing purchase order lines", "error", err)
				service.InternalError(w, err.Error())
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

			inventoryItemUUID, err := parseUUIDPointer(req.InventoryItemUUID)
			if err != nil {
				http.Error(w, "invalid inventory_item_uuid", http.StatusBadRequest)
				return
			}

			line := storage.PurchaseOrderLine{
				PurchaseOrderID:   req.PurchaseOrderID,
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
				slog.Error("error creating purchase order line", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderLineResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandlePurchaseOrderLineByID handles [GET /purchase-order-lines/{id}], [PATCH /purchase-order-lines/{id}], and [DELETE /purchase-order-lines/{id}].
func HandlePurchaseOrderLineByID(db PurchaseOrderLineStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		lineID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			line, err := db.GetPurchaseOrderLine(r.Context(), lineID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order line not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting purchase order line", "error", err)
				service.InternalError(w, err.Error())
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

			line, err := db.UpdatePurchaseOrderLine(r.Context(), lineID, update)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order line not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating purchase order line", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderLineResponse(line))
		case http.MethodDelete:
			line, err := db.DeletePurchaseOrderLine(r.Context(), lineID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order line not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error deleting purchase order line", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderLineResponse(line))
		default:
			methodNotAllowed(w)
		}
	}
}
