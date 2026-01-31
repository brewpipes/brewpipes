package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type InventoryTransferStore interface {
	CreateInventoryTransfer(context.Context, storage.InventoryTransfer) (storage.InventoryTransfer, error)
	GetInventoryTransfer(context.Context, int64) (storage.InventoryTransfer, error)
	ListInventoryTransfers(context.Context) ([]storage.InventoryTransfer, error)
}

// HandleInventoryTransfers handles [GET /inventory-transfers] and [POST /inventory-transfers].
func HandleInventoryTransfers(db InventoryTransferStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			transfers, err := db.ListInventoryTransfers(r.Context())
			if err != nil {
				slog.Error("error listing inventory transfers", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewInventoryTransfersResponse(transfers))
		case http.MethodPost:
			var req dto.CreateInventoryTransferRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			transferredAt := time.Time{}
			if req.TransferredAt != nil {
				transferredAt = *req.TransferredAt
			}

			transfer := storage.InventoryTransfer{
				SourceLocationID: req.SourceLocationID,
				DestLocationID:   req.DestLocationID,
				TransferredAt:    transferredAt,
				Notes:            req.Notes,
			}

			created, err := db.CreateInventoryTransfer(r.Context(), transfer)
			if err != nil {
				slog.Error("error creating inventory transfer", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewInventoryTransferResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleInventoryTransferByID handles [GET /inventory-transfers/{id}].
func HandleInventoryTransferByID(db InventoryTransferStore) http.HandlerFunc {
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

		transfer, err := db.GetInventoryTransfer(r.Context(), transferID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "inventory transfer not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting inventory transfer", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewInventoryTransferResponse(transfer))
	}
}
