package handler

import (
	"context"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// BeerLotStockLevelStore defines the storage interface for beer lot stock level queries.
type BeerLotStockLevelStore interface {
	GetBeerLotStockLevels(context.Context) ([]storage.BeerLotStockLevel, error)
}

// HandleBeerLotStockLevels handles [GET /beer-lot-stock-levels].
func HandleBeerLotStockLevels(db BeerLotStockLevelStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			service.MethodNotAllowed(w)
			return
		}

		levels, err := db.GetBeerLotStockLevels(r.Context())
		if err != nil {
			service.InternalError(w, "error getting beer lot stock levels", "error", err)
			return
		}

		service.JSON(w, dto.NewBeerLotStockLevelsResponse(levels))
	}
}
