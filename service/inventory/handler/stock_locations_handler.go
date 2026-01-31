package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type StockLocationStore interface {
	ListStockLocations(context.Context) ([]storage.StockLocation, error)
	GetStockLocation(context.Context, int64) (storage.StockLocation, error)
	CreateStockLocation(context.Context, storage.StockLocation) (storage.StockLocation, error)
}

// HandleStockLocations handles [GET /stock-locations] and [POST /stock-locations].
func HandleStockLocations(db StockLocationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			locations, err := db.ListStockLocations(r.Context())
			if err != nil {
				slog.Error("error listing stock locations", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewStockLocationsResponse(locations))
		case http.MethodPost:
			var req dto.CreateStockLocationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			location := storage.StockLocation{
				Name:         req.Name,
				LocationType: req.LocationType,
				Description:  req.Description,
			}

			created, err := db.CreateStockLocation(r.Context(), location)
			if err != nil {
				slog.Error("error creating stock location", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewStockLocationResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleStockLocationByID handles [GET /stock-locations/{id}].
func HandleStockLocationByID(db StockLocationStore) http.HandlerFunc {
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
		locationID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		location, err := db.GetStockLocation(r.Context(), locationID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "stock location not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting stock location", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewStockLocationResponse(location))
	}
}
