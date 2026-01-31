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
	ListPurchaseOrdersBySupplier(context.Context, int64) ([]storage.PurchaseOrder, error)
	GetPurchaseOrder(context.Context, int64) (storage.PurchaseOrder, error)
	CreatePurchaseOrder(context.Context, storage.PurchaseOrder) (storage.PurchaseOrder, error)
}

// HandlePurchaseOrders handles [GET /purchase-orders] and [POST /purchase-orders].
func HandlePurchaseOrders(db PurchaseOrderStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			supplierValue := r.URL.Query().Get("supplier_id")
			if supplierValue != "" {
				supplierID, err := parseInt64Param(supplierValue)
				if err != nil {
					http.Error(w, "invalid supplier_id", http.StatusBadRequest)
					return
				}

				orders, err := db.ListPurchaseOrdersBySupplier(r.Context(), supplierID)
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

			status := strings.TrimSpace(req.Status)
			if status == "" {
				status = storage.PurchaseOrderStatusDraft
			}

			order := storage.PurchaseOrder{
				SupplierID:  req.SupplierID,
				OrderNumber: req.OrderNumber,
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

// HandlePurchaseOrderByID handles [GET /purchase-orders/{id}].
func HandlePurchaseOrderByID(db PurchaseOrderStore) http.HandlerFunc {
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
		orderID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		order, err := db.GetPurchaseOrder(r.Context(), orderID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "purchase order not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting purchase order", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewPurchaseOrderResponse(order))
	}
}
