package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type BatchSummaryGetter interface {
	GetBatchSummaryByUUID(context.Context, string) (storage.BatchSummary, error)
}

// HandleBatchSummaryByUUID handles [GET /batches/{uuid}/summary].
func HandleBatchSummaryByUUID(db BatchSummaryGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		batchUUID := r.PathValue("uuid")
		if batchUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		summary, err := db.GetBatchSummaryByUUID(r.Context(), batchUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting batch summary", "error", err, "batch_uuid", batchUUID)
			return
		}

		service.JSON(w, dto.NewBatchSummaryResponse(summary))
	}
}
