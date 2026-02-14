package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type BatchProcessPhaseStore interface {
	CreateBatchProcessPhase(context.Context, storage.BatchProcessPhase) (storage.BatchProcessPhase, error)
	GetBatchProcessPhaseByUUID(context.Context, string) (storage.BatchProcessPhase, error)
	ListBatchProcessPhasesByBatchUUID(context.Context, string) ([]storage.BatchProcessPhase, error)
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
}

// HandleBatchProcessPhases handles [GET /batch-process-phases] and [POST /batch-process-phases].
func HandleBatchProcessPhases(db BatchProcessPhaseStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchUUID := r.URL.Query().Get("batch_uuid")
			if batchUUID == "" {
				http.Error(w, "batch_uuid is required", http.StatusBadRequest)
				return
			}

			phases, err := db.ListBatchProcessPhasesByBatchUUID(r.Context(), batchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error listing batch process phases", "error", err)
				return
			}

			service.JSON(w, dto.NewBatchProcessPhasesResponse(phases))
		case http.MethodPost:
			var req dto.CreateBatchProcessPhaseRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Resolve batch UUID to internal ID
			batch, ok := service.ResolveFK(r.Context(), w, req.BatchUUID, "batch", db.GetBatchByUUID)
			if !ok {
				return
			}

			phaseAt := time.Time{}
			if req.PhaseAt != nil {
				phaseAt = *req.PhaseAt
			}

			phase := storage.BatchProcessPhase{
				BatchID:      batch.ID,
				ProcessPhase: req.ProcessPhase,
				PhaseAt:      phaseAt,
			}

			created, err := db.CreateBatchProcessPhase(r.Context(), phase)
			if err != nil {
				service.InternalError(w, "error creating batch process phase", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewBatchProcessPhaseResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleBatchProcessPhaseByUUID handles [GET /batch-process-phases/{uuid}].
func HandleBatchProcessPhaseByUUID(db BatchProcessPhaseStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phaseUUID := r.PathValue("uuid")
		if phaseUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		phase, err := db.GetBatchProcessPhaseByUUID(r.Context(), phaseUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch process phase not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting batch process phase", "error", err, "phase_uuid", phaseUUID)
			return
		}

		service.JSON(w, dto.NewBatchProcessPhaseResponse(phase))
	}
}
