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

type BeerLotStore interface {
	CreateBeerLot(context.Context, storage.BeerLot) (storage.BeerLot, error)
	GetBeerLot(context.Context, int64) (storage.BeerLot, error)
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
					slog.Error("error listing beer lots", "error", err)
					service.InternalError(w, err.Error())
					return
				}

				service.JSON(w, dto.NewBeerLotsResponse(lots))
				return
			}

			lots, err := db.ListBeerLots(r.Context())
			if err != nil {
				slog.Error("error listing beer lots", "error", err)
				service.InternalError(w, err.Error())
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

			lot := storage.BeerLot{
				ProductionBatchUUID: batchUUID,
				LotCode:             req.LotCode,
				PackagedAt:          packagedAt,
				Notes:               req.Notes,
			}

			created, err := db.CreateBeerLot(r.Context(), lot)
			if err != nil {
				slog.Error("error creating beer lot", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewBeerLotResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleBeerLotByID handles [GET /beer-lots/{id}].
func HandleBeerLotByID(db BeerLotStore) http.HandlerFunc {
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
		lotID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		lot, err := db.GetBeerLot(r.Context(), lotID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "beer lot not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting beer lot", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewBeerLotResponse(lot))
	}
}
