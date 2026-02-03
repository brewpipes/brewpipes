package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type RecipeStore interface {
	CreateRecipe(context.Context, storage.Recipe) (storage.Recipe, error)
	GetRecipe(context.Context, int64) (storage.Recipe, error)
	ListRecipes(context.Context) ([]storage.Recipe, error)
	UpdateRecipe(context.Context, int64, storage.Recipe) (storage.Recipe, error)
	DeleteRecipe(context.Context, int64) error
}

type RecipeBatchCounter interface {
	CountBatchesByRecipe(context.Context, int64) (int, error)
}

// HandleRecipes handles [GET /recipes] and [POST /recipes].
func HandleRecipes(db RecipeStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			recipes, err := db.ListRecipes(r.Context())
			if err != nil {
				slog.Error("error listing recipes", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewRecipesResponse(recipes))
		case http.MethodPost:
			var req dto.CreateRecipeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			recipe := storage.Recipe{
				Name:      req.Name,
				StyleID:   req.StyleID,
				StyleName: req.StyleName,
				Notes:     req.Notes,
			}

			created, err := db.CreateRecipe(r.Context(), recipe)
			if err != nil {
				slog.Error("error creating recipe", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("recipe created", "recipe_id", created.ID, "name", created.Name)

			service.JSON(w, dto.NewRecipeResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleRecipeByID handles [GET /recipes/{id}], [PUT /recipes/{id}], and [DELETE /recipes/{id}].
func HandleRecipeByID(db RecipeStore, batches RecipeBatchCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		recipeID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			recipe, err := db.GetRecipe(r.Context(), recipeID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting recipe", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewRecipeResponse(recipe))
		case http.MethodPut:
			var req dto.UpdateRecipeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			recipe := storage.Recipe{
				Name:      req.Name,
				StyleID:   req.StyleID,
				StyleName: req.StyleName,
				Notes:     req.Notes,
			}

			updated, err := db.UpdateRecipe(r.Context(), recipeID, recipe)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating recipe", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("recipe updated", "recipe_id", updated.ID, "name", updated.Name)

			service.JSON(w, dto.NewRecipeResponse(updated))
		case http.MethodDelete:
			// Check if any batches reference this recipe
			batchCount, err := batches.CountBatchesByRecipe(r.Context(), recipeID)
			if err != nil {
				slog.Error("error counting batches by recipe", "error", err, "recipe_id", recipeID)
				service.InternalError(w, err.Error())
				return
			}
			if batchCount > 0 {
				msg := fmt.Sprintf("cannot delete recipe: %d batch(es) reference this recipe", batchCount)
				http.Error(w, msg, http.StatusConflict)
				return
			}

			err = db.DeleteRecipe(r.Context(), recipeID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error deleting recipe", "error", err, "recipe_id", recipeID)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("recipe deleted", "recipe_id", recipeID)

			w.WriteHeader(http.StatusNoContent)
		default:
			methodNotAllowed(w)
		}
	}
}
