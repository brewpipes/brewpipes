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

type InventoryUsageStore interface {
	CreateInventoryUsage(context.Context, storage.InventoryUsage) (storage.InventoryUsage, error)
	GetInventoryUsageByUUID(context.Context, string) (storage.InventoryUsage, error)
	ListInventoryUsage(context.Context) ([]storage.InventoryUsage, error)
}

// HandleInventoryUsage handles [GET /inventory-usage] and [POST /inventory-usage].
func HandleInventoryUsage(db InventoryUsageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			usageRecords, err := db.ListInventoryUsage(r.Context())
			if err != nil {
				service.InternalError(w, "error listing inventory usage", "error", err)
				return
			}

			service.JSON(w, dto.NewInventoryUsageRecordsResponse(usageRecords))
		case http.MethodPost:
			var req dto.CreateInventoryUsageRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			usedAt := time.Time{}
			if req.UsedAt != nil {
				usedAt = *req.UsedAt
			}

			var productionUUID *uuid.UUID
			if req.ProductionRefUUID != nil {
				parsed, err := uuid.FromString(*req.ProductionRefUUID)
				if err != nil {
					http.Error(w, "invalid production_ref_uuid", http.StatusBadRequest)
					return
				}
				productionUUID = &parsed
			}

			usage := storage.InventoryUsage{
				ProductionRefUUID: productionUUID,
				UsedAt:            usedAt,
				Notes:             req.Notes,
			}

			created, err := db.CreateInventoryUsage(r.Context(), usage)
			if err != nil {
				service.InternalError(w, "error creating inventory usage", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewInventoryUsageResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleInventoryUsageByUUID handles [GET /inventory-usage/{uuid}].
func HandleInventoryUsageByUUID(db InventoryUsageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usageUUID := r.PathValue("uuid")
		if usageUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		usage, err := db.GetInventoryUsageByUUID(r.Context(), usageUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "inventory usage not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting inventory usage", "error", err)
			return
		}

		service.JSON(w, dto.NewInventoryUsageResponse(usage))
	}
}
