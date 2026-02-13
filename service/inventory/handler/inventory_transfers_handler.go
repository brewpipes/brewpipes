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
	GetInventoryTransferByUUID(context.Context, string) (storage.InventoryTransfer, error)
	ListInventoryTransfers(context.Context) ([]storage.InventoryTransfer, error)
	GetStockLocationByUUID(context.Context, string) (storage.StockLocation, error)
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

			// Resolve source location UUID to internal ID
			sourceLocation, err := db.GetStockLocationByUUID(r.Context(), req.SourceLocationUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "source location not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving source location uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			// Resolve dest location UUID to internal ID
			destLocation, err := db.GetStockLocationByUUID(r.Context(), req.DestLocationUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "dest location not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving dest location uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			transfer := storage.InventoryTransfer{
				SourceLocationID: sourceLocation.ID,
				DestLocationID:   destLocation.ID,
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

// HandleInventoryTransferByUUID handles [GET /inventory-transfers/{uuid}].
func HandleInventoryTransferByUUID(db InventoryTransferStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		transferUUID := r.PathValue("uuid")
		if transferUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		transfer, err := db.GetInventoryTransferByUUID(r.Context(), transferUUID)
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
