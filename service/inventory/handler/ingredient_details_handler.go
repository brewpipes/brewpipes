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

type IngredientMaltDetailStore interface {
	CreateIngredientMaltDetail(context.Context, storage.IngredientMaltDetail) (storage.IngredientMaltDetail, error)
	GetIngredientMaltDetail(context.Context, int64) (storage.IngredientMaltDetail, error)
	GetIngredientMaltDetailByIngredient(context.Context, int64) (storage.IngredientMaltDetail, error)
}

// HandleIngredientMaltDetails handles [GET /ingredient-malt-details] and [POST /ingredient-malt-details].
func HandleIngredientMaltDetails(db IngredientMaltDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_id")
			if ingredientValue == "" {
				http.Error(w, "ingredient_id is required", http.StatusBadRequest)
				return
			}
			ingredientID, err := parseInt64Param(ingredientValue)
			if err != nil {
				http.Error(w, "invalid ingredient_id", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientMaltDetailByIngredient(r.Context(), ingredientID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient malt detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting ingredient malt detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientMaltDetailResponse(detail))
		case http.MethodPost:
			var req dto.CreateIngredientMaltDetailRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			detail := storage.IngredientMaltDetail{
				IngredientID:   req.IngredientID,
				MaltsterName:   req.MaltsterName,
				Variety:        req.Variety,
				Lovibond:       req.Lovibond,
				SRM:            req.SRM,
				DiastaticPower: req.DiastaticPower,
			}

			created, err := db.CreateIngredientMaltDetail(r.Context(), detail)
			if err != nil {
				slog.Error("error creating ingredient malt detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientMaltDetailResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleIngredientMaltDetailByID handles [GET /ingredient-malt-details/{id}].
func HandleIngredientMaltDetailByID(db IngredientMaltDetailStore) http.HandlerFunc {
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

		detail, err := db.GetIngredientMaltDetail(r.Context(), detailID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient malt detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting ingredient malt detail", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewIngredientMaltDetailResponse(detail))
	}
}

type IngredientHopDetailStore interface {
	CreateIngredientHopDetail(context.Context, storage.IngredientHopDetail) (storage.IngredientHopDetail, error)
	GetIngredientHopDetail(context.Context, int64) (storage.IngredientHopDetail, error)
	GetIngredientHopDetailByIngredient(context.Context, int64) (storage.IngredientHopDetail, error)
}

// HandleIngredientHopDetails handles [GET /ingredient-hop-details] and [POST /ingredient-hop-details].
func HandleIngredientHopDetails(db IngredientHopDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_id")
			if ingredientValue == "" {
				http.Error(w, "ingredient_id is required", http.StatusBadRequest)
				return
			}
			ingredientID, err := parseInt64Param(ingredientValue)
			if err != nil {
				http.Error(w, "invalid ingredient_id", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientHopDetailByIngredient(r.Context(), ingredientID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient hop detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting ingredient hop detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientHopDetailResponse(detail))
		case http.MethodPost:
			var req dto.CreateIngredientHopDetailRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			detail := storage.IngredientHopDetail{
				IngredientID: req.IngredientID,
				ProducerName: req.ProducerName,
				Variety:      req.Variety,
				CropYear:     req.CropYear,
				Form:         req.Form,
				AlphaAcid:    req.AlphaAcid,
				BetaAcid:     req.BetaAcid,
			}

			created, err := db.CreateIngredientHopDetail(r.Context(), detail)
			if err != nil {
				slog.Error("error creating ingredient hop detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientHopDetailResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleIngredientHopDetailByID handles [GET /ingredient-hop-details/{id}].
func HandleIngredientHopDetailByID(db IngredientHopDetailStore) http.HandlerFunc {
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

		detail, err := db.GetIngredientHopDetail(r.Context(), detailID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient hop detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting ingredient hop detail", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewIngredientHopDetailResponse(detail))
	}
}

type IngredientYeastDetailStore interface {
	CreateIngredientYeastDetail(context.Context, storage.IngredientYeastDetail) (storage.IngredientYeastDetail, error)
	GetIngredientYeastDetail(context.Context, int64) (storage.IngredientYeastDetail, error)
	GetIngredientYeastDetailByIngredient(context.Context, int64) (storage.IngredientYeastDetail, error)
}

// HandleIngredientYeastDetails handles [GET /ingredient-yeast-details] and [POST /ingredient-yeast-details].
func HandleIngredientYeastDetails(db IngredientYeastDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_id")
			if ingredientValue == "" {
				http.Error(w, "ingredient_id is required", http.StatusBadRequest)
				return
			}
			ingredientID, err := parseInt64Param(ingredientValue)
			if err != nil {
				http.Error(w, "invalid ingredient_id", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientYeastDetailByIngredient(r.Context(), ingredientID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient yeast detail not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting ingredient yeast detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientYeastDetailResponse(detail))
		case http.MethodPost:
			var req dto.CreateIngredientYeastDetailRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			detail := storage.IngredientYeastDetail{
				IngredientID: req.IngredientID,
				LabName:      req.LabName,
				Strain:       req.Strain,
				Form:         req.Form,
			}

			created, err := db.CreateIngredientYeastDetail(r.Context(), detail)
			if err != nil {
				slog.Error("error creating ingredient yeast detail", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewIngredientYeastDetailResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleIngredientYeastDetailByID handles [GET /ingredient-yeast-details/{id}].
func HandleIngredientYeastDetailByID(db IngredientYeastDetailStore) http.HandlerFunc {
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

		detail, err := db.GetIngredientYeastDetail(r.Context(), detailID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient yeast detail not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting ingredient yeast detail", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewIngredientYeastDetailResponse(detail))
	}
}
