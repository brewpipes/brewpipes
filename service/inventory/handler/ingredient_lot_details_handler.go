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

type IngredientLotMaltDetailStore interface {
	CreateIngredientLotMaltDetail(context.Context, storage.IngredientLotMaltDetail) (storage.IngredientLotMaltDetail, error)
	GetIngredientLotMaltDetail(context.Context, int64) (storage.IngredientLotMaltDetail, error)
	GetIngredientLotMaltDetailByLot(context.Context, int64) (storage.IngredientLotMaltDetail, error)
}

// HandleIngredientLotMaltDetails handles [GET /ingredient-lot-malt-details] and [POST /ingredient-lot-malt-details].
func HandleIngredientLotMaltDetails(db IngredientLotMaltDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			lotValue := r.URL.Query().Get("ingredient_lot_id")
			if lotValue == "" {
				http.Error(w, "ingredient_lot_id is required", http.StatusBadRequest)
				return
			}
			lotID, err := parseInt64Param(lotValue)
			if err != nil {
				http.Error(w, "invalid ingredient_lot_id", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientLotMaltDetailByLot(r.Context(), lotID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient lot malt detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting ingredient lot malt detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientLotMaltDetailResponse(detail))
		case http.MethodPost:
			var req dto.CreateIngredientLotMaltDetailRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			detail := storage.IngredientLotMaltDetail{
				IngredientLotID: req.IngredientLotID,
				MoisturePercent: req.MoisturePercent,
			}

			created, err := db.CreateIngredientLotMaltDetail(r.Context(), detail)
			if err != nil {
				slog.Error("error creating ingredient lot malt detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientLotMaltDetailResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleIngredientLotMaltDetailByID handles [GET /ingredient-lot-malt-details/{id}].
func HandleIngredientLotMaltDetailByID(db IngredientLotMaltDetailStore) http.HandlerFunc {
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
		detailID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientLotMaltDetail(r.Context(), detailID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient lot malt detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting ingredient lot malt detail", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewIngredientLotMaltDetailResponse(detail))
	}
}

type IngredientLotHopDetailStore interface {
	CreateIngredientLotHopDetail(context.Context, storage.IngredientLotHopDetail) (storage.IngredientLotHopDetail, error)
	GetIngredientLotHopDetail(context.Context, int64) (storage.IngredientLotHopDetail, error)
	GetIngredientLotHopDetailByLot(context.Context, int64) (storage.IngredientLotHopDetail, error)
}

// HandleIngredientLotHopDetails handles [GET /ingredient-lot-hop-details] and [POST /ingredient-lot-hop-details].
func HandleIngredientLotHopDetails(db IngredientLotHopDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			lotValue := r.URL.Query().Get("ingredient_lot_id")
			if lotValue == "" {
				http.Error(w, "ingredient_lot_id is required", http.StatusBadRequest)
				return
			}
			lotID, err := parseInt64Param(lotValue)
			if err != nil {
				http.Error(w, "invalid ingredient_lot_id", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientLotHopDetailByLot(r.Context(), lotID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient lot hop detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting ingredient lot hop detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientLotHopDetailResponse(detail))
		case http.MethodPost:
			var req dto.CreateIngredientLotHopDetailRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			detail := storage.IngredientLotHopDetail{
				IngredientLotID: req.IngredientLotID,
				AlphaAcid:       req.AlphaAcid,
				BetaAcid:        req.BetaAcid,
			}

			created, err := db.CreateIngredientLotHopDetail(r.Context(), detail)
			if err != nil {
				slog.Error("error creating ingredient lot hop detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientLotHopDetailResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleIngredientLotHopDetailByID handles [GET /ingredient-lot-hop-details/{id}].
func HandleIngredientLotHopDetailByID(db IngredientLotHopDetailStore) http.HandlerFunc {
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
		detailID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientLotHopDetail(r.Context(), detailID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient lot hop detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting ingredient lot hop detail", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewIngredientLotHopDetailResponse(detail))
	}
}

type IngredientLotYeastDetailStore interface {
	CreateIngredientLotYeastDetail(context.Context, storage.IngredientLotYeastDetail) (storage.IngredientLotYeastDetail, error)
	GetIngredientLotYeastDetail(context.Context, int64) (storage.IngredientLotYeastDetail, error)
	GetIngredientLotYeastDetailByLot(context.Context, int64) (storage.IngredientLotYeastDetail, error)
}

// HandleIngredientLotYeastDetails handles [GET /ingredient-lot-yeast-details] and [POST /ingredient-lot-yeast-details].
func HandleIngredientLotYeastDetails(db IngredientLotYeastDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			lotValue := r.URL.Query().Get("ingredient_lot_id")
			if lotValue == "" {
				http.Error(w, "ingredient_lot_id is required", http.StatusBadRequest)
				return
			}
			lotID, err := parseInt64Param(lotValue)
			if err != nil {
				http.Error(w, "invalid ingredient_lot_id", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientLotYeastDetailByLot(r.Context(), lotID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient lot yeast detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting ingredient lot yeast detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientLotYeastDetailResponse(detail))
		case http.MethodPost:
			var req dto.CreateIngredientLotYeastDetailRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			detail := storage.IngredientLotYeastDetail{
				IngredientLotID:  req.IngredientLotID,
				ViabilityPercent: req.Viability,
				Generation:       req.Generation,
			}

			created, err := db.CreateIngredientLotYeastDetail(r.Context(), detail)
			if err != nil {
				slog.Error("error creating ingredient lot yeast detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientLotYeastDetailResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleIngredientLotYeastDetailByID handles [GET /ingredient-lot-yeast-details/{id}].
func HandleIngredientLotYeastDetailByID(db IngredientLotYeastDetailStore) http.HandlerFunc {
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
		detailID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientLotYeastDetail(r.Context(), detailID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient lot yeast detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting ingredient lot yeast detail", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewIngredientLotYeastDetailResponse(detail))
	}
}
