package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type BatchSummaryGetter interface {
	GetBatchSummary(context.Context, int64) (storage.BatchSummary, error)
}

// HandleBatchSummary handles [GET /batches/{id}/summary].
func HandleBatchSummary(db BatchSummaryGetter) http.HandlerFunc {
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
		batchID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		summary, err := db.GetBatchSummary(r.Context(), batchID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting batch summary", "error", err, "batch_id", batchID)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewBatchSummaryResponse(summary))
	}
}
