package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
	"github.com/gofrs/uuid/v5"
)

// RecipeIngredientStore defines the storage interface for recipe ingredients.
type RecipeIngredientStore interface {
	ListRecipeIngredients(ctx context.Context, recipeUUID string) ([]storage.RecipeIngredient, error)
	GetRecipeIngredient(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error)
	CreateRecipeIngredient(ctx context.Context, ri storage.RecipeIngredient) (storage.RecipeIngredient, error)
	UpdateRecipeIngredient(ctx context.Context, ingredientUUID string, ri storage.RecipeIngredient) (storage.RecipeIngredient, error)
	DeleteRecipeIngredient(ctx context.Context, ingredientUUID string) error
}

// RecipeExistenceChecker checks if a recipe exists and returns the recipe.
type RecipeExistenceChecker interface {
	GetRecipe(ctx context.Context, recipeUUID string, opts *storage.RecipeQueryOpts) (storage.Recipe, error)
}

// HandleRecipeIngredients handles [GET /recipes/{uuid}/ingredients] and [POST /recipes/{uuid}/ingredients].
func HandleRecipeIngredients(db RecipeIngredientStore, recipes RecipeExistenceChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeUUID := r.PathValue("uuid")
		if recipeUUID == "" {
			http.Error(w, "invalid recipe uuid", http.StatusBadRequest)
			return
		}

		// Verify recipe exists and get internal ID for FK
		recipe, err := recipes.GetRecipe(r.Context(), recipeUUID, nil)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "recipe not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error checking recipe existence", "error", err, "recipe_uuid", recipeUUID)
			return
		}

		switch r.Method {
		case http.MethodGet:
			ingredients, err := db.ListRecipeIngredients(r.Context(), recipeUUID)
			if err != nil {
				service.InternalError(w, "error listing recipe ingredients", "error", err, "recipe_uuid", recipeUUID)
				return
			}

			service.JSON(w, dto.NewRecipeIngredientsResponse(ingredients, recipeUUID))
		case http.MethodPost:
			var req dto.RecipeIngredientRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var ingredientUUID *uuid.UUID
			if req.IngredientUUID != nil {
				parsed, err := uuid.FromString(*req.IngredientUUID)
				if err != nil {
					http.Error(w, "invalid ingredient_uuid", http.StatusBadRequest)
					return
				}
				ingredientUUID = &parsed
			}

			// Apply defaults
			scalingFactor := 1.0
			if req.ScalingFactor != nil {
				scalingFactor = *req.ScalingFactor
			}
			sortOrder := 0
			if req.SortOrder != nil {
				sortOrder = *req.SortOrder
			}

			ri := storage.RecipeIngredient{
				RecipeID:              recipe.ID, // Use internal ID for FK
				Name:                  req.Name,
				IngredientUUID:        ingredientUUID,
				IngredientType:        req.IngredientType,
				Amount:                req.Amount,
				AmountUnit:            req.AmountUnit,
				UseStage:              req.UseStage,
				UseType:               req.UseType,
				TimingDurationMinutes: req.TimingDurationMinutes,
				TimingTemperatureC:    req.TimingTemperatureC,
				AlphaAcidAssumed:      req.AlphaAcidAssumed,
				ScalingFactor:         scalingFactor,
				SortOrder:             sortOrder,
				Notes:                 req.Notes,
			}

			created, err := db.CreateRecipeIngredient(r.Context(), ri)
			if err != nil {
				service.InternalError(w, "error creating recipe ingredient", "error", err, "recipe_uuid", recipeUUID)
				return
			}

			slog.Info("recipe ingredient created",
				"recipe_ingredient_uuid", created.UUID.String(),
				"recipe_uuid", recipeUUID,
				"ingredient_type", created.IngredientType)

			service.JSONCreated(w, dto.NewRecipeIngredientResponse(created, recipeUUID))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleRecipeIngredient handles [GET /recipes/{uuid}/ingredients/{ingredient_uuid}],
// [PATCH /recipes/{uuid}/ingredients/{ingredient_uuid}], and [DELETE /recipes/{uuid}/ingredients/{ingredient_uuid}].
func HandleRecipeIngredient(db RecipeIngredientStore, recipes RecipeExistenceChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeUUID := r.PathValue("uuid")
		if recipeUUID == "" {
			http.Error(w, "invalid recipe uuid", http.StatusBadRequest)
			return
		}

		ingredientUUID := r.PathValue("ingredient_uuid")
		if ingredientUUID == "" {
			http.Error(w, "invalid ingredient uuid", http.StatusBadRequest)
			return
		}

		// Verify recipe exists and get internal ID for FK comparison
		recipe, err := recipes.GetRecipe(r.Context(), recipeUUID, nil)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "recipe not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error checking recipe existence", "error", err, "recipe_uuid", recipeUUID)
			return
		}

		switch r.Method {
		case http.MethodGet:
			ri, err := db.GetRecipeIngredient(r.Context(), ingredientUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe ingredient not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting recipe ingredient", "error", err, "ingredient_uuid", ingredientUUID)
				return
			}

			// Verify ingredient belongs to the recipe
			if ri.RecipeID != recipe.ID {
				http.Error(w, "recipe ingredient not found", http.StatusNotFound)
				return
			}

			service.JSON(w, dto.NewRecipeIngredientResponse(ri, recipeUUID))
		case http.MethodPatch:
			// Get existing ingredient
			existing, err := db.GetRecipeIngredient(r.Context(), ingredientUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe ingredient not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting recipe ingredient", "error", err, "ingredient_uuid", ingredientUUID)
				return
			}

			// Verify ingredient belongs to the recipe
			if existing.RecipeID != recipe.ID {
				http.Error(w, "recipe ingredient not found", http.StatusNotFound)
				return
			}

			var req dto.RecipeIngredientRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var invIngredientUUID *uuid.UUID
			if req.IngredientUUID != nil {
				parsed, err := uuid.FromString(*req.IngredientUUID)
				if err != nil {
					http.Error(w, "invalid ingredient_uuid", http.StatusBadRequest)
					return
				}
				invIngredientUUID = &parsed
			}

			// Apply defaults from existing if not provided
			scalingFactor := existing.ScalingFactor
			if req.ScalingFactor != nil {
				scalingFactor = *req.ScalingFactor
			}
			sortOrder := existing.SortOrder
			if req.SortOrder != nil {
				sortOrder = *req.SortOrder
			}

			existing.Name = req.Name
			existing.IngredientUUID = invIngredientUUID
			existing.IngredientType = req.IngredientType
			existing.Amount = req.Amount
			existing.AmountUnit = req.AmountUnit
			existing.UseStage = req.UseStage
			existing.UseType = req.UseType
			existing.TimingDurationMinutes = req.TimingDurationMinutes
			existing.TimingTemperatureC = req.TimingTemperatureC
			existing.AlphaAcidAssumed = req.AlphaAcidAssumed
			existing.ScalingFactor = scalingFactor
			existing.SortOrder = sortOrder
			existing.Notes = req.Notes

			updated, err := db.UpdateRecipeIngredient(r.Context(), ingredientUUID, existing)
			if err != nil {
				service.InternalError(w, "error updating recipe ingredient", "error", err, "ingredient_uuid", ingredientUUID)
				return
			}

			slog.Info("recipe ingredient updated",
				"recipe_ingredient_uuid", ingredientUUID,
				"recipe_uuid", recipeUUID)

			service.JSON(w, dto.NewRecipeIngredientResponse(updated, recipeUUID))
		case http.MethodDelete:
			// Verify ingredient exists and belongs to recipe before deleting
			existing, err := db.GetRecipeIngredient(r.Context(), ingredientUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe ingredient not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting recipe ingredient for delete", "error", err, "ingredient_uuid", ingredientUUID)
				return
			}

			if existing.RecipeID != recipe.ID {
				http.Error(w, "recipe ingredient not found", http.StatusNotFound)
				return
			}

			if err := db.DeleteRecipeIngredient(r.Context(), ingredientUUID); err != nil {
				if errors.Is(err, service.ErrNotFound) {
					http.Error(w, "recipe ingredient not found", http.StatusNotFound)
					return
				}
				service.InternalError(w, "error deleting recipe ingredient", "error", err, "ingredient_uuid", ingredientUUID)
				return
			}

			slog.Info("recipe ingredient deleted",
				"recipe_ingredient_uuid", ingredientUUID,
				"recipe_uuid", recipeUUID)

			w.WriteHeader(http.StatusNoContent)
		default:
			service.MethodNotAllowed(w)
		}
	}
}
