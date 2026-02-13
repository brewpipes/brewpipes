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
	ListPurchaseOrderFeesByOrderUUID(context.Context, string) ([]storage.PurchaseOrderFee, error)
	GetPurchaseOrderFeeByUUID(context.Context, string) (storage.PurchaseOrderFee, error)
	GetPurchaseOrderByUUID(context.Context, string) (storage.PurchaseOrder, error)
	CreatePurchaseOrderFee(context.Context, storage.PurchaseOrderFee) (storage.PurchaseOrderFee, error)
	UpdatePurchaseOrderFeeByUUID(context.Context, string, storage.PurchaseOrderFeeUpdate) (storage.PurchaseOrderFee, error)
	DeletePurchaseOrderFeeByUUID(context.Context, string) (storage.PurchaseOrderFee, error)
}

// HandlePurchaseOrderFees handles [GET /purchase-order-fees] and [POST /purchase-order-fees].
func HandlePurchaseOrderFees(db PurchaseOrderFeeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			orderUUID := r.URL.Query().Get("purchase_order_uuid")
			if orderUUID != "" {
				fees, err := db.ListPurchaseOrderFeesByOrderUUID(r.Context(), orderUUID)
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

			// Resolve purchase order UUID to internal ID
			order, err := db.GetPurchaseOrderByUUID(r.Context(), req.PurchaseOrderUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "purchase order not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving purchase order uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			fee := storage.PurchaseOrderFee{
				PurchaseOrderID: order.ID,
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

// HandlePurchaseOrderFeeByUUID handles [GET /purchase-order-fees/{uuid}], [PATCH /purchase-order-fees/{uuid}], and [DELETE /purchase-order-fees/{uuid}].
func HandlePurchaseOrderFeeByUUID(db PurchaseOrderFeeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		feeUUID := r.PathValue("uuid")
		if feeUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			fee, err := db.GetPurchaseOrderFeeByUUID(r.Context(), feeUUID)
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

			fee, err := db.UpdatePurchaseOrderFeeByUUID(r.Context(), feeUUID, update)
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
			fee, err := db.DeletePurchaseOrderFeeByUUID(r.Context(), feeUUID)
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
