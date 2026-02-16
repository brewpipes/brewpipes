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
)

// InventoryRemovalStore defines the storage interface for removal handlers.
type InventoryRemovalStore interface {
	CreateRemovalWithMovement(context.Context, storage.RemovalWithMovementRequest) (storage.RemovalWithMovementResult, error)
	GetRemovalByUUID(context.Context, string) (storage.InventoryRemoval, error)
	ListRemovals(context.Context, storage.RemovalListFilter) ([]storage.InventoryRemoval, error)
	UpdateRemoval(context.Context, string, storage.UpdateRemovalRequest) (storage.InventoryRemoval, error)
	SoftDeleteRemoval(context.Context, string) error
	GetRemovalSummary(context.Context, storage.RemovalListFilter) (storage.RemovalSummary, error)
	GetBeerLotByUUID(context.Context, string) (storage.BeerLot, error)
	GetStockLocationByUUID(context.Context, string) (storage.StockLocation, error)
}

// convertToBBL converts an amount in the given unit to barrels (BBL).
// Returns nil for unknown units.
func convertToBBL(amount int64, unit string) *float64 {
	var bbl float64
	switch unit {
	case "bbl":
		bbl = float64(amount)
	case "gal":
		bbl = float64(amount) / 31.0
	case "l":
		bbl = float64(amount) / 117.34777
	case "ml":
		bbl = float64(amount) / 117347.77
	case "usfloz":
		bbl = float64(amount) / 3968.0
	default:
		return nil
	}
	return &bbl
}

// isTaxable returns whether a removal category is taxable.
// All V1 categories are non-taxable.
func isTaxable(_ string) bool {
	return false
}

// HandleRemovals handles [GET /removals] and [POST /removals].
func HandleRemovals(db InventoryRemovalStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			filter := storage.RemovalListFilter{}
			q := r.URL.Query()

			if v := q.Get("batch_uuid"); v != "" {
				filter.BatchUUID = &v
			}
			if v := q.Get("beer_lot_uuid"); v != "" {
				filter.BeerLotUUID = &v
			}
			if v := q.Get("category"); v != "" {
				filter.Category = &v
			}
			if v := q.Get("from"); v != "" {
				t, err := time.Parse(time.RFC3339, v)
				if err != nil {
					http.Error(w, "invalid from: must be RFC3339 format", http.StatusBadRequest)
					return
				}
				filter.From = &t
			}
			if v := q.Get("to"); v != "" {
				t, err := time.Parse(time.RFC3339, v)
				if err != nil {
					http.Error(w, "invalid to: must be RFC3339 format", http.StatusBadRequest)
					return
				}
				filter.To = &t
			}

			removals, err := db.ListRemovals(r.Context(), filter)
			if err != nil {
				service.InternalError(w, "error listing removals", "error", err)
				return
			}

			service.JSON(w, dto.NewRemovalsResponse(removals))

		case http.MethodPost:
			var req dto.CreateRemovalRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			removedAt := time.Time{}
			if req.RemovedAt != nil {
				removedAt = *req.RemovedAt
			}

			// Resolve beer lot UUID to internal ID if provided.
			var beerLotID *int64
			if lot, ok := service.ResolveFKOptional(r.Context(), w, req.BeerLotUUID, "beer lot", db.GetBeerLotByUUID); !ok {
				return
			} else if req.BeerLotUUID != nil {
				beerLotID = &lot.ID
			}

			// Resolve stock location UUID to internal ID if provided.
			var stockLocationID *int64
			if loc, ok := service.ResolveFKOptional(r.Context(), w, req.StockLocationUUID, "stock location", db.GetStockLocationByUUID); !ok {
				return
			} else if req.StockLocationUUID != nil {
				stockLocationID = &loc.ID
			}

			amountBBL := convertToBBL(req.Amount, req.AmountUnit)
			taxable := isTaxable(req.Category)

			result, err := db.CreateRemovalWithMovement(r.Context(), storage.RemovalWithMovementRequest{
				Category:          req.Category,
				Reason:            req.Reason,
				BatchUUID:         req.BatchUUID,
				BeerLotID:         beerLotID,
				OccupancyUUID:     req.OccupancyUUID,
				Amount:            req.Amount,
				AmountUnit:        req.AmountUnit,
				AmountBBL:         amountBBL,
				IsTaxable:         taxable,
				ReferenceCode:     req.ReferenceCode,
				PerformedBy:       req.PerformedBy,
				RemovedAt:         removedAt,
				Destination:       req.Destination,
				Notes:             req.Notes,
				StockLocationID:   stockLocationID,
				BeerLotUUID:       req.BeerLotUUID,
				StockLocationUUID: req.StockLocationUUID,
			})
			if err != nil {
				service.InternalError(w, "error creating removal", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewRemovalResponse(result.Removal))

		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleRemovalByUUID handles [GET /removals/{uuid}], [PATCH /removals/{uuid}],
// and [DELETE /removals/{uuid}].
func HandleRemovalByUUID(db InventoryRemovalStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		removalUUID := r.PathValue("uuid")
		if removalUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			removal, err := db.GetRemovalByUUID(r.Context(), removalUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "removal not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting removal", "error", err)
				return
			}

			service.JSON(w, dto.NewRemovalResponse(removal))

		case http.MethodPatch:
			var req dto.UpdateRemovalRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Build the storage update request, recalculating derived fields.
			storageReq := storage.UpdateRemovalRequest{
				Category:      req.Category,
				Reason:        req.Reason,
				Amount:        req.Amount,
				AmountUnit:    req.AmountUnit,
				ReferenceCode: req.ReferenceCode,
				PerformedBy:   req.PerformedBy,
				RemovedAt:     req.RemovedAt,
				Destination:   req.Destination,
				Notes:         req.Notes,
			}

			// Recalculate amount_bbl if amount or amount_unit changes.
			if req.Amount != nil || req.AmountUnit != nil {
				// We need the current values to fill in whichever wasn't provided.
				existing, err := db.GetRemovalByUUID(r.Context(), removalUUID)
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "removal not found", http.StatusNotFound)
					return
				} else if err != nil {
					service.InternalError(w, "error getting removal for update", "error", err)
					return
				}

				amount := existing.Amount
				if req.Amount != nil {
					amount = *req.Amount
				}
				unit := existing.AmountUnit
				if req.AmountUnit != nil {
					unit = *req.AmountUnit
				}
				storageReq.AmountBBL = convertToBBL(amount, unit)
			}

			// Recalculate is_taxable if category changes.
			if req.Category != nil {
				taxable := isTaxable(*req.Category)
				storageReq.IsTaxable = &taxable
			}

			updated, err := db.UpdateRemoval(r.Context(), removalUUID, storageReq)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "removal not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error updating removal", "error", err)
				return
			}

			service.JSON(w, dto.NewRemovalResponse(updated))

		case http.MethodDelete:
			err := db.SoftDeleteRemoval(r.Context(), removalUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "removal not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error deleting removal", "error", err)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleRemovalSummary handles [GET /removal-summary].
func HandleRemovalSummary(db InventoryRemovalStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			service.MethodNotAllowed(w)
			return
		}

		filter := storage.RemovalListFilter{}
		q := r.URL.Query()

		if v := q.Get("batch_uuid"); v != "" {
			filter.BatchUUID = &v
		}
		if v := q.Get("beer_lot_uuid"); v != "" {
			filter.BeerLotUUID = &v
		}
		if v := q.Get("category"); v != "" {
			filter.Category = &v
		}
		if v := q.Get("from"); v != "" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				http.Error(w, "invalid from: must be RFC3339 format", http.StatusBadRequest)
				return
			}
			filter.From = &t
		}
		if v := q.Get("to"); v != "" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				http.Error(w, "invalid to: must be RFC3339 format", http.StatusBadRequest)
				return
			}
			filter.To = &t
		}

		summary, err := db.GetRemovalSummary(r.Context(), filter)
		if err != nil {
			service.InternalError(w, "error getting removal summary", "error", err)
			return
		}

		service.JSON(w, dto.NewRemovalSummaryResponse(summary))
	}
}
