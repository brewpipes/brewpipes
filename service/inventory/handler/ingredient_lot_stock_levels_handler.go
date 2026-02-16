package handler

import (
	"context"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// IngredientLotStockLevelStore defines the storage interface for ingredient lot stock level queries.
type IngredientLotStockLevelStore interface {
	GetIngredientLotStockLevels(context.Context) ([]storage.IngredientLotStockLevel, error)
}

// HandleIngredientLotStockLevels handles [GET /ingredient-lot-stock-levels].
func HandleIngredientLotStockLevels(db IngredientLotStockLevelStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			service.MethodNotAllowed(w)
			return
		}

		levels, err := db.GetIngredientLotStockLevels(r.Context())
		if err != nil {
			service.InternalError(w, "error getting ingredient lot stock levels", "error", err)
			return
		}

		service.JSON(w, dto.NewIngredientLotStockLevelsResponse(levels))
	}
}
