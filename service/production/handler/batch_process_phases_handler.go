package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type BatchProcessPhaseStore interface {
	CreateBatchProcessPhase(context.Context, storage.BatchProcessPhase) (storage.BatchProcessPhase, error)
	GetBatchProcessPhase(context.Context, int64) (storage.BatchProcessPhase, error)
	ListBatchProcessPhases(context.Context, int64) ([]storage.BatchProcessPhase, error)
}

// HandleBatchProcessPhases handles [GET /batch-process-phases] and [POST /batch-process-phases].
func HandleBatchProcessPhases(db BatchProcessPhaseStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchValue := r.URL.Query().Get("batch_id")
			if batchValue == "" {
				http.Error(w, "batch_id is required", http.StatusBadRequest)
				return
			}
			batchID, err := parseInt64Param(batchValue)
			if err != nil {
				http.Error(w, "invalid batch_id", http.StatusBadRequest)
				return
			}

			phases, err := db.ListBatchProcessPhases(r.Context(), batchID)
			if err != nil {
				slog.Error("error listing batch process phases", "error", err)
				service.InternalError(w, err.Error())
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

			phaseAt := time.Time{}
			if req.PhaseAt != nil {
				phaseAt = *req.PhaseAt
			}

			phase := storage.BatchProcessPhase{
				BatchID:      req.BatchID,
				ProcessPhase: req.ProcessPhase,
				PhaseAt:      phaseAt,
			}

			created, err := db.CreateBatchProcessPhase(r.Context(), phase)
			if err != nil {
				slog.Error("error creating batch process phase", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBatchProcessPhaseResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleBatchProcessPhaseByID handles [GET /batch-process-phases/{id}].
func HandleBatchProcessPhaseByID(db BatchProcessPhaseStore) http.HandlerFunc {
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
		phaseID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		phase, err := db.GetBatchProcessPhase(r.Context(), phaseID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch process phase not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting batch process phase", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewBatchProcessPhaseResponse(phase))
	}
}
