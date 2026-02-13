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
	GetIngredientMaltDetailByUUID(context.Context, string) (storage.IngredientMaltDetail, error)
	GetIngredientMaltDetailByIngredient(context.Context, string) (storage.IngredientMaltDetail, error)
	GetIngredientByUUID(context.Context, string) (storage.Ingredient, error)
}

// HandleIngredientMaltDetails handles [GET /ingredient-malt-details] and [POST /ingredient-malt-details].
func HandleIngredientMaltDetails(db IngredientMaltDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_uuid")
			if ingredientValue == "" {
				http.Error(w, "ingredient_uuid is required", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientMaltDetailByIngredient(r.Context(), ingredientValue)
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

			// Resolve ingredient UUID to internal ID
			ingredient, err := db.GetIngredientByUUID(r.Context(), req.IngredientUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving ingredient uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			detail := storage.IngredientMaltDetail{
				IngredientID:   ingredient.ID,
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

// HandleIngredientMaltDetailByUUID handles [GET /ingredient-malt-details/{uuid}].
func HandleIngredientMaltDetailByUUID(db IngredientMaltDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		detailUUID := r.PathValue("uuid")
		if detailUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientMaltDetailByUUID(r.Context(), detailUUID)
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
	GetIngredientHopDetailByUUID(context.Context, string) (storage.IngredientHopDetail, error)
	GetIngredientHopDetailByIngredient(context.Context, string) (storage.IngredientHopDetail, error)
	GetIngredientByUUID(context.Context, string) (storage.Ingredient, error)
}

// HandleIngredientHopDetails handles [GET /ingredient-hop-details] and [POST /ingredient-hop-details].
func HandleIngredientHopDetails(db IngredientHopDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_uuid")
			if ingredientValue == "" {
				http.Error(w, "ingredient_uuid is required", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientHopDetailByIngredient(r.Context(), ingredientValue)
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

			// Resolve ingredient UUID to internal ID
			ingredient, err := db.GetIngredientByUUID(r.Context(), req.IngredientUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving ingredient uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			detail := storage.IngredientHopDetail{
				IngredientID: ingredient.ID,
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

// HandleIngredientHopDetailByUUID handles [GET /ingredient-hop-details/{uuid}].
func HandleIngredientHopDetailByUUID(db IngredientHopDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		detailUUID := r.PathValue("uuid")
		if detailUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientHopDetailByUUID(r.Context(), detailUUID)
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
	GetIngredientYeastDetailByUUID(context.Context, string) (storage.IngredientYeastDetail, error)
	GetIngredientYeastDetailByIngredient(context.Context, string) (storage.IngredientYeastDetail, error)
	GetIngredientByUUID(context.Context, string) (storage.Ingredient, error)
}

// HandleIngredientYeastDetails handles [GET /ingredient-yeast-details] and [POST /ingredient-yeast-details].
func HandleIngredientYeastDetails(db IngredientYeastDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ingredientValue := r.URL.Query().Get("ingredient_uuid")
			if ingredientValue == "" {
				http.Error(w, "ingredient_uuid is required", http.StatusBadRequest)
				return
			}

			detail, err := db.GetIngredientYeastDetailByIngredient(r.Context(), ingredientValue)
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

			// Resolve ingredient UUID to internal ID
			ingredient, err := db.GetIngredientByUUID(r.Context(), req.IngredientUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "ingredient not found", http.StatusBadRequest)
				return
			} else if err != nil {
				slog.Error("error resolving ingredient uuid", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			detail := storage.IngredientYeastDetail{
				IngredientID: ingredient.ID,
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

// HandleIngredientYeastDetailByUUID handles [GET /ingredient-yeast-details/{uuid}].
func HandleIngredientYeastDetailByUUID(db IngredientYeastDetailStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		detailUUID := r.PathValue("uuid")
		if detailUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		detail, err := db.GetIngredientYeastDetailByUUID(r.Context(), detailUUID)
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
