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

type IngredientStore interface {
	ListIngredients(context.Context) ([]storage.Ingredient, error)
	ListIngredientsIncludingDeleted(context.Context) ([]storage.Ingredient, error)
	GetIngredientByUUID(context.Context, string) (storage.Ingredient, error)
	CreateIngredient(context.Context, storage.Ingredient) (storage.Ingredient, error)
}

// HandleIngredients handles [GET /ingredients] and [POST /ingredients].
func HandleIngredients(db IngredientStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			includeDeleted := r.URL.Query().Get("include_deleted") == "true"
			var ingredients []storage.Ingredient
			var err error
			if includeDeleted {
				ingredients, err = db.ListIngredientsIncludingDeleted(r.Context())
			} else {
				ingredients, err = db.ListIngredients(r.Context())
			}
			if err != nil {
				service.InternalError(w, "error listing ingredients", "error", err)
				return
			}

			service.JSON(w, dto.NewIngredientsResponse(ingredients))
		case http.MethodPost:
			var req dto.CreateIngredientRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ingredient := storage.Ingredient{
				Name:        req.Name,
				Category:    req.Category,
				DefaultUnit: req.DefaultUnit,
				Description: req.Description,
			}

			created, err := db.CreateIngredient(r.Context(), ingredient)
			if err != nil {
				service.InternalError(w, "error creating ingredient", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewIngredientResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleIngredientByUUID handles [GET /ingredients/{uuid}].
func HandleIngredientByUUID(db IngredientStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredientUUID := r.PathValue("uuid")
		if ingredientUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		ingredient, err := db.GetIngredientByUUID(r.Context(), ingredientUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "ingredient not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting ingredient", "error", err)
			return
		}

		service.JSON(w, dto.NewIngredientResponse(ingredient))
	}
}
