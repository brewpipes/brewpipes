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

type TransferStore interface {
	GetTransferByUUID(context.Context, string) (storage.Transfer, error)
	RecordTransfer(context.Context, storage.TransferRecord) (storage.Transfer, storage.Occupancy, error)
	ListTransfersByBatchUUID(context.Context, string) ([]storage.Transfer, error)
	GetOccupancyByUUID(context.Context, string) (storage.Occupancy, error)
	GetVesselByUUID(context.Context, string) (storage.Vessel, error)
	GetVolumeByUUID(context.Context, string) (storage.Volume, error)
}

// HandleTransfers handles [GET /transfers] and [POST /transfers].
func HandleTransfers(db TransferStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchUUID := r.URL.Query().Get("batch_uuid")
			if batchUUID == "" {
				http.Error(w, "batch_uuid is required", http.StatusBadRequest)
				return
			}

			transfers, err := db.ListTransfersByBatchUUID(r.Context(), batchUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "batch not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error listing transfers", "error", err)
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

			// Resolve source occupancy UUID to internal ID
			sourceOcc, ok := service.ResolveFK(r.Context(), w, req.SourceOccupancyUUID, "source occupancy", db.GetOccupancyByUUID)
			if !ok {
				return
			}

			// Resolve dest vessel UUID to internal ID
			destVessel, ok := service.ResolveFK(r.Context(), w, req.DestVesselUUID, "destination vessel", db.GetVesselByUUID)
			if !ok {
				return
			}

			// Resolve volume UUID to internal ID
			volume, ok := service.ResolveFK(r.Context(), w, req.VolumeUUID, "volume", db.GetVolumeByUUID)
			if !ok {
				return
			}

			startedAt := time.Time{}
			if req.StartedAt != nil {
				startedAt = *req.StartedAt
			}

			record := storage.TransferRecord{
				SourceOccupancyID: sourceOcc.ID,
				DestVesselID:      destVessel.ID,
				VolumeID:          volume.ID,
				Amount:            req.Amount,
				AmountUnit:        req.AmountUnit,
				LossAmount:        req.LossAmount,
				LossUnit:          req.LossUnit,
				StartedAt:         startedAt,
				EndedAt:           req.EndedAt,
			}

			transfer, occupancy, err := db.RecordTransfer(r.Context(), record)
			if err != nil {
				service.InternalError(w, "error recording transfer", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewTransferRecordResponse(transfer, occupancy))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleTransferByUUID handles [GET /transfers/{uuid}].
func HandleTransferByUUID(db TransferStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transferUUID := r.PathValue("uuid")
		if transferUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		transfer, err := db.GetTransferByUUID(r.Context(), transferUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "transfer not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting transfer", "error", err, "transfer_uuid", transferUUID)
			return
		}

		service.JSON(w, dto.NewTransferResponse(transfer))
	}
}
