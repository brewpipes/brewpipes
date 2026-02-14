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

type IngredientLotMaltDetailStore interface {
	CreateIngredientLotMaltDetail(context.Context, storage.IngredientLotMaltDetail) (storage.IngredientLotMaltDetail, error)
	GetIngredientLotMaltDetailByUUID(context.Context, string) (storage.IngredientLotMaltDetail, error)
	GetIngredientLotMaltDetailByLot(context.Context, string) (storage.IngredientLotMaltDetail, error)
	GetIngredientLotByUUID(context.Context, string) (storage.IngredientLot, error)
}

// HandleIngredientLotMaltDetails handles [GET /ingredient-lot-malt-details] and [POST /ingredient-lot-malt-details].
func HandleIngredientLotMaltDetails(db IngredientLotMaltDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			lotValue := r.URL.Query().Get("ingredient_lot_uuid")
			if lotValue == "" {
				http.Error(w, "ingredient_lot_uuid is required", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientLotMaltDetailByLot(r.Context(), lotValue)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient lot malt detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting ingredient lot malt detail", "error", err)
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

			// Resolve ingredient lot UUID to internal ID
			lot, ok := service.ResolveFK(r.Context(), w, req.IngredientLotUUID, "ingredient lot", db.GetIngredientLotByUUID)
			if !ok {
				return
			}

			detail := storage.IngredientLotMaltDetail{
				IngredientLotID: lot.ID,
				MoisturePercent: req.MoisturePercent,
			}

			created, err := db.CreateIngredientLotMaltDetail(r.Context(), detail)
			if err != nil {
				service.InternalError(w, "error creating ingredient lot malt detail", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewIngredientLotMaltDetailResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleIngredientLotMaltDetailByUUID handles [GET /ingredient-lot-malt-details/{uuid}].
func HandleIngredientLotMaltDetailByUUID(db IngredientLotMaltDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		detailUUID := r.PathValue("uuid")
		if detailUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientLotMaltDetailByUUID(r.Context(), detailUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient lot malt detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting ingredient lot malt detail", "error", err)
			return
		}

		service.JSON(w, dto.NewIngredientLotMaltDetailResponse(detail))
	}
}

type IngredientLotHopDetailStore interface {
	CreateIngredientLotHopDetail(context.Context, storage.IngredientLotHopDetail) (storage.IngredientLotHopDetail, error)
	GetIngredientLotHopDetailByUUID(context.Context, string) (storage.IngredientLotHopDetail, error)
	GetIngredientLotHopDetailByLot(context.Context, string) (storage.IngredientLotHopDetail, error)
	GetIngredientLotByUUID(context.Context, string) (storage.IngredientLot, error)
}

// HandleIngredientLotHopDetails handles [GET /ingredient-lot-hop-details] and [POST /ingredient-lot-hop-details].
func HandleIngredientLotHopDetails(db IngredientLotHopDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			lotValue := r.URL.Query().Get("ingredient_lot_uuid")
			if lotValue == "" {
				http.Error(w, "ingredient_lot_uuid is required", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientLotHopDetailByLot(r.Context(), lotValue)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient lot hop detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting ingredient lot hop detail", "error", err)
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

			// Resolve ingredient lot UUID to internal ID
			lot, ok := service.ResolveFK(r.Context(), w, req.IngredientLotUUID, "ingredient lot", db.GetIngredientLotByUUID)
			if !ok {
				return
			}

			detail := storage.IngredientLotHopDetail{
				IngredientLotID: lot.ID,
				AlphaAcid:       req.AlphaAcid,
				BetaAcid:        req.BetaAcid,
			}

			created, err := db.CreateIngredientLotHopDetail(r.Context(), detail)
			if err != nil {
				service.InternalError(w, "error creating ingredient lot hop detail", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewIngredientLotHopDetailResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleIngredientLotHopDetailByUUID handles [GET /ingredient-lot-hop-details/{uuid}].
func HandleIngredientLotHopDetailByUUID(db IngredientLotHopDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		detailUUID := r.PathValue("uuid")
		if detailUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientLotHopDetailByUUID(r.Context(), detailUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient lot hop detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting ingredient lot hop detail", "error", err)
			return
		}

		service.JSON(w, dto.NewIngredientLotHopDetailResponse(detail))
	}
}

type IngredientLotYeastDetailStore interface {
	CreateIngredientLotYeastDetail(context.Context, storage.IngredientLotYeastDetail) (storage.IngredientLotYeastDetail, error)
	GetIngredientLotYeastDetailByUUID(context.Context, string) (storage.IngredientLotYeastDetail, error)
	GetIngredientLotYeastDetailByLot(context.Context, string) (storage.IngredientLotYeastDetail, error)
	GetIngredientLotByUUID(context.Context, string) (storage.IngredientLot, error)
}

// HandleIngredientLotYeastDetails handles [GET /ingredient-lot-yeast-details] and [POST /ingredient-lot-yeast-details].
func HandleIngredientLotYeastDetails(db IngredientLotYeastDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			lotValue := r.URL.Query().Get("ingredient_lot_uuid")
			if lotValue == "" {
				http.Error(w, "ingredient_lot_uuid is required", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientLotYeastDetailByLot(r.Context(), lotValue)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient lot yeast detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting ingredient lot yeast detail", "error", err)
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

			// Resolve ingredient lot UUID to internal ID
			lot, ok := service.ResolveFK(r.Context(), w, req.IngredientLotUUID, "ingredient lot", db.GetIngredientLotByUUID)
			if !ok {
				return
			}

			detail := storage.IngredientLotYeastDetail{
				IngredientLotID:  lot.ID,
				ViabilityPercent: req.Viability,
				Generation:       req.Generation,
			}

			created, err := db.CreateIngredientLotYeastDetail(r.Context(), detail)
			if err != nil {
				service.InternalError(w, "error creating ingredient lot yeast detail", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewIngredientLotYeastDetailResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleIngredientLotYeastDetailByUUID handles [GET /ingredient-lot-yeast-details/{uuid}].
func HandleIngredientLotYeastDetailByUUID(db IngredientLotYeastDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		detailUUID := r.PathValue("uuid")
		if detailUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientLotYeastDetailByUUID(r.Context(), detailUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient lot yeast detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting ingredient lot yeast detail", "error", err)
			return
		}

		service.JSON(w, dto.NewIngredientLotYeastDetailResponse(detail))
	}
}
