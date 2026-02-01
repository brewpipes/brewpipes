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
)

type PurchaseOrderFeeStore interface {
	ListPurchaseOrderFees(context.Context) ([]storage.PurchaseOrderFee, error)
	ListPurchaseOrderFeesByOrder(context.Context, int64) ([]storage.PurchaseOrderFee, error)
	GetPurchaseOrderFee(context.Context, int64) (storage.PurchaseOrderFee, error)
	CreatePurchaseOrderFee(context.Context, storage.PurchaseOrderFee) (storage.PurchaseOrderFee, error)
	UpdatePurchaseOrderFee(context.Context, int64, storage.PurchaseOrderFeeUpdate) (storage.PurchaseOrderFee, error)
	DeletePurchaseOrderFee(context.Context, int64) (storage.PurchaseOrderFee, error)
}

// HandlePurchaseOrderFees handles [GET /purchase-order-fees] and [POST /purchase-order-fees].
func HandlePurchaseOrderFees(db PurchaseOrderFeeStore) http.HandlerFunc {
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

				fees, err := db.ListPurchaseOrderFeesByOrder(r.Context(), orderID)
				if err != nil {
					slog.Error("error listing purchase order fees by order", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewPurchaseOrderFeesResponse(fees))
				return
			}

			fees, err := db.ListPurchaseOrderFees(r.Context())
			if err != nil {
				slog.Error("error listing purchase order fees", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderFeesResponse(fees))
		case http.MethodPost:
			var req dto.CreatePurchaseOrderFeeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fee := storage.PurchaseOrderFee{
				PurchaseOrderID: req.PurchaseOrderID,
				FeeType:         req.FeeType,
				AmountCents:     req.AmountCents,
				Currency:        req.Currency,
			}

			created, err := db.CreatePurchaseOrderFee(r.Context(), fee)
			if err != nil {
				slog.Error("error creating purchase order fee", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderFeeResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandlePurchaseOrderFeeByID handles [GET /purchase-order-fees/{id}], [PATCH /purchase-order-fees/{id}], and [DELETE /purchase-order-fees/{id}].
func HandlePurchaseOrderFeeByID(db PurchaseOrderFeeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		feeID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			fee, err := db.GetPurchaseOrderFee(r.Context(), feeID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order fee not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting purchase order fee", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderFeeResponse(fee))
		case http.MethodPatch:
			var req dto.UpdatePurchaseOrderFeeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			update := storage.PurchaseOrderFeeUpdate{
				FeeType:     req.FeeType,
				AmountCents: req.AmountCents,
				Currency:    req.Currency,
			}

			fee, err := db.UpdatePurchaseOrderFee(r.Context(), feeID, update)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order fee not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating purchase order fee", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderFeeResponse(fee))
		case http.MethodDelete:
			fee, err := db.DeletePurchaseOrderFee(r.Context(), feeID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order fee not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error deleting purchase order fee", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewPurchaseOrderFeeResponse(fee))
		default:
			methodNotAllowed(w)
		}
	}
}
