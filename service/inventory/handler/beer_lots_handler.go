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

type BeerLotStore interface {
	CreateBeerLot(context.Context, storage.BeerLot) (storage.BeerLot, error)
	CreateBeerLotWithMovement(ctx context.Context, lot storage.BeerLot, stockLocationID int64, movementAmount int64, movementAmountUnit string) (storage.BeerLot, uuid.UUID, error)
	GetBeerLotByUUID(context.Context, string) (storage.BeerLot, error)
	GetStockLocationByUUID(context.Context, string) (storage.StockLocation, error)
	ListBeerLots(context.Context) ([]storage.BeerLot, error)
	ListBeerLotsByBatchUUID(context.Context, uuid.UUID) ([]storage.BeerLot, error)
}

// HandleBeerLots handles [GET /beer-lots] and [POST /beer-lots].
func HandleBeerLots(db BeerLotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			batchValue := r.URL.Query().Get("production_batch_uuid")
			if batchValue != "" {
				batchUUID, err := uuid.FromString(batchValue)
				if err != nil {
					http.Error(w, "invalid production_batch_uuid", http.StatusBadRequest)
					return
				}

				lots, err := db.ListBeerLotsByBatchUUID(r.Context(), batchUUID)
				if err != nil {
					service.InternalError(w, "error listing beer lots by batch", "error", err)
					return
				}

				service.JSON(w, dto.NewBeerLotsResponse(lots))
				return
			}

			lots, err := db.ListBeerLots(r.Context())
			if err != nil {
				service.InternalError(w, "error listing beer lots", "error", err)
				return
			}

			service.JSON(w, dto.NewBeerLotsResponse(lots))
		case http.MethodPost:
			var req dto.CreateBeerLotRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			batchUUID, err := uuid.FromString(req.ProductionBatchUUID)
			if err != nil {
				http.Error(w, "invalid production_batch_uuid", http.StatusBadRequest)
				return
			}

			packagedAt := time.Time{}
			if req.PackagedAt != nil {
				packagedAt = *req.PackagedAt
			}

			// Parse optional packaging_run_uuid.
			var packagingRunUUID *uuid.UUID
			if req.PackagingRunUUID != nil {
				parsed, err := uuid.FromString(*req.PackagingRunUUID)
				if err != nil {
					http.Error(w, "invalid packaging_run_uuid", http.StatusBadRequest)
					return
				}
				packagingRunUUID = &parsed
			}

			// Parse optional best_by.
			var bestBy *time.Time
			if req.BestBy != nil {
				parsed, err := time.Parse(time.RFC3339, *req.BestBy)
				if err != nil {
					http.Error(w, "invalid best_by: must be RFC3339 format", http.StatusBadRequest)
					return
				}
				bestBy = &parsed
			}

			lot := storage.BeerLot{
				ProductionBatchUUID: batchUUID,
				PackagingRunUUID:    packagingRunUUID,
				LotCode:             req.LotCode,
				BestBy:              bestBy,
				PackageFormatName:   req.PackageFormatName,
				Container:           req.Container,
				VolumePerUnit:       req.VolumePerUnit,
				VolumePerUnitUnit:   req.VolumePerUnitUnit,
				Quantity:            req.Quantity,
				PackagedAt:          packagedAt,
				Notes:               req.Notes,
			}

			// If stock_location_uuid is provided, create lot with movement atomically.
			if req.StockLocationUUID != nil {
				location, ok := service.ResolveFK(r.Context(), w, *req.StockLocationUUID, "stock location", db.GetStockLocationByUUID)
				if !ok {
					return
				}

				movementAmount := int64(*req.Quantity) * *req.VolumePerUnit
				movementAmountUnit := *req.VolumePerUnitUnit

				created, _, err := db.CreateBeerLotWithMovement(r.Context(), lot, location.ID, movementAmount, movementAmountUnit)
				if err != nil {
					service.InternalError(w, "error creating beer lot with movement", "error", err)
					return
				}

				service.JSONCreated(w, dto.NewBeerLotResponse(created))
				return
			}

			created, err := db.CreateBeerLot(r.Context(), lot)
			if err != nil {
				service.InternalError(w, "error creating beer lot", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewBeerLotResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleBeerLotByUUID handles [GET /beer-lots/{uuid}].
func HandleBeerLotByUUID(db BeerLotStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lotUUID := r.PathValue("uuid")
		if lotUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		lot, err := db.GetBeerLotByUUID(r.Context(), lotUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "beer lot not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting beer lot", "error", err)
			return
		}

		service.JSON(w, dto.NewBeerLotResponse(lot))
	}
}
