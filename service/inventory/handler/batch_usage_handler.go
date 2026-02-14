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

// BatchUsageStore is the storage interface for batch usage operations.
type BatchUsageStore interface {
	CreateBatchUsage(context.Context, storage.BatchUsageRequest) (storage.BatchUsageResult, error)
}

// HandleCreateBatchUsage handles [POST /inventory-usage/batch].
func HandleCreateBatchUsage(db BatchUsageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			service.MethodNotAllowed(w)
			return
		}

		var req dto.CreateBatchUsageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		if err := req.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		usedAt, _ := time.Parse(time.RFC3339, req.UsedAt) // already validated

		var productionUUID *uuid.UUID
		if req.ProductionRefUUID != nil {
			parsed, err := uuid.FromString(*req.ProductionRefUUID)
			if err != nil {
				http.Error(w, "invalid production_ref_uuid", http.StatusBadRequest)
				return
			}
			productionUUID = &parsed
		}

		picks := make([]storage.BatchUsagePick, len(req.Picks))
		for i, p := range req.Picks {
			picks[i] = storage.BatchUsagePick{
				IngredientLotUUID: p.IngredientLotUUID,
				StockLocationUUID: p.StockLocationUUID,
				Amount:            p.Amount,
				AmountUnit:        p.AmountUnit,
			}
		}

		result, err := db.CreateBatchUsage(r.Context(), storage.BatchUsageRequest{
			ProductionRefUUID: productionUUID,
			UsedAt:            usedAt,
			Picks:             picks,
			Notes:             req.Notes,
		})
		if err != nil {
			var validationErr *storage.ErrBatchUsageValidation
			if errors.As(err, &validationErr) {
				http.Error(w, validationErr.Message, http.StatusBadRequest)
				return
			}
			service.InternalError(w, "error creating batch usage", "error", err)
			return
		}

		service.JSONCreated(w, dto.NewBatchUsageResponse(result.Usage, result.Movements))
	}
}
