package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type StockLocationStore interface {
	ListStockLocations(context.Context) ([]storage.StockLocation, error)
	GetStockLocationByUUID(context.Context, string) (storage.StockLocation, error)
	CreateStockLocation(context.Context, storage.StockLocation) (storage.StockLocation, error)
}

// HandleStockLocations handles [GET /stock-locations] and [POST /stock-locations].
func HandleStockLocations(db StockLocationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			locations, err := db.ListStockLocations(r.Context())
			if err != nil {
				service.InternalError(w, "error listing stock locations", "error", err)
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
				service.InternalError(w, "error creating stock location", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewStockLocationResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleStockLocationByUUID handles [GET /stock-locations/{uuid}].
func HandleStockLocationByUUID(db StockLocationStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locationUUID := r.PathValue("uuid")
		if locationUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		location, err := db.GetStockLocationByUUID(r.Context(), locationUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "stock location not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting stock location", "error", err)
			return
		}

		service.JSON(w, dto.NewStockLocationResponse(location))
	}
}
