package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/procurement/handler/dto"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

// PurchaseOrderLineBatchLookupStore defines the storage methods needed by the batch lookup handler.
type PurchaseOrderLineBatchLookupStore interface {
	ListPurchaseOrderLinesByUUIDs(context.Context, []string) ([]storage.PurchaseOrderLine, error)
}

// HandleBatchLookupPurchaseOrderLines handles [POST /purchase-order-lines/batch-lookup].
func HandleBatchLookupPurchaseOrderLines(db PurchaseOrderLineBatchLookupStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			service.MethodNotAllowed(w)
			return
		}

		var req dto.BatchLookupPurchaseOrderLinesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		lines, err := db.ListPurchaseOrderLinesByUUIDs(r.Context(), req.UUIDs)
		if err != nil {
			service.InternalError(w, "error looking up purchase order lines", "error", err)
			return
		}

		service.JSON(w, dto.NewPurchaseOrderLinesResponse(lines))
	}
}
