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
	"github.com/gofrs/uuid/v5"
)

type InventoryUsageStore interface {
	CreateInventoryUsage(context.Context, storage.InventoryUsage) (storage.InventoryUsage, error)
	GetInventoryUsage(context.Context, int64) (storage.InventoryUsage, error)
	ListInventoryUsage(context.Context) ([]storage.InventoryUsage, error)
}

// HandleInventoryUsage handles [GET /inventory-usage] and [POST /inventory-usage].
func HandleInventoryUsage(db InventoryUsageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			usageRecords, err := db.ListInventoryUsage(r.Context())
			if err != nil {
				slog.Error("error listing inventory usage", "error", err)
				service.InternalError(w, err.Error())
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
				slog.Error("error creating inventory usage", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewInventoryUsageResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleInventoryUsageByID handles [GET /inventory-usage/{id}].
func HandleInventoryUsageByID(db InventoryUsageStore) http.HandlerFunc {
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
		usageID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		usage, err := db.GetInventoryUsage(r.Context(), usageID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "inventory usage not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting inventory usage", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewInventoryUsageResponse(usage))
	}
}
