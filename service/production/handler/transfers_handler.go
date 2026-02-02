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

type TransferStore interface {
	GetTransfer(context.Context, int64) (storage.Transfer, error)
	RecordTransfer(context.Context, storage.TransferRecord) (storage.Transfer, storage.Occupancy, error)
	ListTransfersByBatch(context.Context, int64) ([]storage.Transfer, error)
}

// HandleTransfers handles [GET /transfers] and [POST /transfers].
func HandleTransfers(db TransferStore) http.HandlerFunc {
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

			transfers, err := db.ListTransfersByBatch(r.Context(), batchID)
			if err != nil {
				slog.Error("error listing transfers", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewTransfersResponse(transfers))
		case http.MethodPost:
			var req dto.CreateTransferRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			startedAt := time.Time{}
			if req.StartedAt != nil {
				startedAt = *req.StartedAt
			}

			record := storage.TransferRecord{
				SourceOccupancyID: req.SourceOccupancyID,
				DestVesselID:      req.DestVesselID,
				VolumeID:          req.VolumeID,
				Amount:            req.Amount,
				AmountUnit:        req.AmountUnit,
				LossAmount:        req.LossAmount,
				LossUnit:          req.LossUnit,
				StartedAt:         startedAt,
				EndedAt:           req.EndedAt,
			}

			transfer, occupancy, err := db.RecordTransfer(r.Context(), record)
			if err != nil {
				slog.Error("error recording transfer", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewTransferRecordResponse(transfer, occupancy))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleTransferByID handles [GET /transfers/{id}].
func HandleTransferByID(db TransferStore) http.HandlerFunc {
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
		transferID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		transfer, err := db.GetTransfer(r.Context(), transferID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "transfer not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting transfer", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewTransferResponse(transfer))
	}
}
