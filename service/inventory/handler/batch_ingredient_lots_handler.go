package handler

import (
	"context"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// BatchIngredientLotStore defines the storage methods needed by the batch ingredient lots handler.
type BatchIngredientLotStore interface {
	ListBatchIngredientLots(context.Context, string) ([]storage.BatchIngredientLot, error)
}

// HandleBatchIngredientLots handles [GET /ingredient-lots/batch].
func HandleBatchIngredientLots(db BatchIngredientLotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			service.MethodNotAllowed(w)
			return
		}

		productionRefUUID := r.URL.Query().Get("production_ref_uuid")
		if productionRefUUID == "" {
			http.Error(w, "production_ref_uuid is required", http.StatusBadRequest)
			return
		}

		lots, err := db.ListBatchIngredientLots(r.Context(), productionRefUUID)
		if err != nil {
			service.InternalError(w, "error listing batch ingredient lots", "error", err, "production_ref_uuid", productionRefUUID)
			return
		}

		service.JSON(w, dto.NewBatchIngredientLotsResponse(lots))
	}
}
