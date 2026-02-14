package handler

import (
	"context"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// StockLevelStore defines the storage interface for stock level queries.
type StockLevelStore interface {
	GetStockLevels(context.Context) ([]storage.StockLevel, error)
}

// HandleStockLevels handles [GET /stock-levels].
func HandleStockLevels(db StockLevelStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			service.MethodNotAllowed(w)
			return
		}

		levels, err := db.GetStockLevels(r.Context())
		if err != nil {
			service.InternalError(w, "error getting stock levels", "error", err)
			return
		}

		service.JSON(w, dto.NewStockLevelsResponse(levels))
	}
}
