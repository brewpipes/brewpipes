package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/procurement/handler/dto"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type PurchaseOrderStore interface {
	ListPurchaseOrders(context.Context) ([]storage.PurchaseOrder, error)
	ListPurchaseOrdersBySupplierUUID(context.Context, string) ([]storage.PurchaseOrder, error)
	GetPurchaseOrderByUUID(context.Context, string) (storage.PurchaseOrder, error)
	GetSupplierByUUID(context.Context, string) (storage.Supplier, error)
	CreatePurchaseOrder(context.Context, storage.PurchaseOrder) (storage.PurchaseOrder, error)
	UpdatePurchaseOrderByUUID(context.Context, string, storage.PurchaseOrderUpdate) (storage.PurchaseOrder, error)
}

// HandlePurchaseOrders handles [GET /purchase-orders] and [POST /purchase-orders].
func HandlePurchaseOrders(db PurchaseOrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			supplierUUID := r.URL.Query().Get("supplier_uuid")
			if supplierUUID != "" {
				orders, err := db.ListPurchaseOrdersBySupplierUUID(r.Context(), supplierUUID)
				if err != nil {
					slog.Error("error listing purchase orders by supplier", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewPurchaseOrdersResponse(orders))
				return
			}

			orders, err := db.ListPurchaseOrders(r.Context())
			if err != nil {
				slog.Error("error listing purchase orders", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrdersResponse(orders))
		case http.MethodPost:
			var req dto.CreatePurchaseOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Resolve supplier UUID to internal ID
			supplier, err := db.GetSupplierByUUID(r.Context(), req.SupplierUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "supplier not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving supplier uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			status := strings.TrimSpace(req.Status)
			if status == "" {
				status = storage.PurchaseOrderStatusDraft
			}
			orderNumber := strings.TrimSpace(req.OrderNumber)

			order := storage.PurchaseOrder{
				SupplierID:  supplier.ID,
				OrderNumber: orderNumber,
				Status:      status,
				OrderedAt:   req.OrderedAt,
				ExpectedAt:  req.ExpectedAt,
				Notes:       req.Notes,
			}

			created, err := db.CreatePurchaseOrder(r.Context(), order)
			if err != nil {
				slog.Error("error creating purchase order", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandlePurchaseOrderByUUID handles [GET /purchase-orders/{uuid}] and [PATCH /purchase-orders/{uuid}].
func HandlePurchaseOrderByUUID(db PurchaseOrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUUID := r.PathValue("uuid")
		if orderUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			order, err := db.GetPurchaseOrderByUUID(r.Context(), orderUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting purchase order", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderResponse(order))
		case http.MethodPatch:
			var req dto.UpdatePurchaseOrderRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			update := storage.PurchaseOrderUpdate{
				OrderNumber: req.OrderNumber,
				Status:      req.Status,
				OrderedAt:   req.OrderedAt,
				ExpectedAt:  req.ExpectedAt,
				Notes:       req.Notes,
			}
			if update.OrderNumber != nil {
				value := strings.TrimSpace(*update.OrderNumber)
				update.OrderNumber = &value
			}
			if update.Status != nil {
				value := strings.TrimSpace(*update.Status)
				update.Status = &value
			}

			order, err := db.UpdatePurchaseOrderByUUID(r.Context(), orderUUID, update)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating purchase order", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderResponse(order))
		default:
			methodNotAllowed(w)
		}
	}
}
