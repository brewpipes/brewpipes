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
)

type RecipeStore interface {
	CreateRecipe(context.Context, storage.Recipe) (storage.Recipe, error)
	GetRecipe(context.Context, string, *storage.RecipeQueryOpts) (storage.Recipe, error)
	ListRecipes(context.Context) ([]storage.Recipe, error)
	UpdateRecipe(context.Context, string, storage.Recipe) (storage.Recipe, error)
	DeleteRecipe(context.Context, string) error
}

type RecipeBatchCounter interface {
	CountBatchesByRecipe(context.Context, string) (int, error)
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
				Name:                req.Name,
				StyleID:             req.StyleID,
				StyleName:           req.StyleName,
				Notes:               req.Notes,
				BatchSize:           req.BatchSize,
				BatchSizeUnit:       req.BatchSizeUnit,
				TargetOG:            req.TargetOG,
				TargetOGMin:         req.TargetOGMin,
				TargetOGMax:         req.TargetOGMax,
				TargetFG:            req.TargetFG,
				TargetFGMin:         req.TargetFGMin,
				TargetFGMax:         req.TargetFGMax,
				TargetIBU:           req.TargetIBU,
				TargetIBUMin:        req.TargetIBUMin,
				TargetIBUMax:        req.TargetIBUMax,
				TargetSRM:           req.TargetSRM,
				TargetSRMMin:        req.TargetSRMMin,
				TargetSRMMax:        req.TargetSRMMax,
				TargetCarbonation:   req.TargetCarbonation,
				IBUMethod:           req.IBUMethod,
				BrewhouseEfficiency: req.BrewhouseEfficiency,
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

// HandleRecipe handles [GET /recipes/{uuid}], [PUT /recipes/{uuid}], [PATCH /recipes/{uuid}], and [DELETE /recipes/{uuid}].
func HandleRecipe(db RecipeStore, batches RecipeBatchCounter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recipeUUID := r.PathValue("uuid")
		if recipeUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			recipe, err := db.GetRecipe(r.Context(), recipeUUID, nil)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error getting recipe", "error", err, "recipe_uuid", recipeUUID)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewRecipeResponse(recipe))
		case http.MethodPut, http.MethodPatch:
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
				Name:                req.Name,
				StyleID:             req.StyleID,
				StyleName:           req.StyleName,
				Notes:               req.Notes,
				BatchSize:           req.BatchSize,
				BatchSizeUnit:       req.BatchSizeUnit,
				TargetOG:            req.TargetOG,
				TargetOGMin:         req.TargetOGMin,
				TargetOGMax:         req.TargetOGMax,
				TargetFG:            req.TargetFG,
				TargetFGMin:         req.TargetFGMin,
				TargetFGMax:         req.TargetFGMax,
				TargetIBU:           req.TargetIBU,
				TargetIBUMin:        req.TargetIBUMin,
				TargetIBUMax:        req.TargetIBUMax,
				TargetSRM:           req.TargetSRM,
				TargetSRMMin:        req.TargetSRMMin,
				TargetSRMMax:        req.TargetSRMMax,
				TargetCarbonation:   req.TargetCarbonation,
				IBUMethod:           req.IBUMethod,
				BrewhouseEfficiency: req.BrewhouseEfficiency,
			}

			updated, err := db.UpdateRecipe(r.Context(), recipeUUID, recipe)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error updating recipe", "error", err, "recipe_uuid", recipeUUID)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("recipe updated", "recipe_uuid", recipeUUID, "name", updated.Name)

			service.JSON(w, dto.NewRecipeResponse(updated))
		case http.MethodDelete:
			// Log batch references for audit purposes (but don't block deletion)
			batchCount, err := batches.CountBatchesByRecipe(r.Context(), recipeUUID)
			if err != nil {
				slog.Warn("could not count batches referencing recipe before delete", "recipe_uuid", recipeUUID, "error", err)
				// Continue with deletion anyway - this is just for logging
			} else if batchCount > 0 {
				slog.Info("deleting recipe with batch references", "recipe_uuid", recipeUUID, "batch_count", batchCount)
			}

			// Soft-delete the recipe; batches retain their recipe_id reference for historical traceability
			err = db.DeleteRecipe(r.Context(), recipeUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "recipe not found", http.StatusNotFound)
				return
			} else if err != nil {
				slog.Error("error deleting recipe", "error", err, "recipe_uuid", recipeUUID)
				service.InternalError(w, err.Error())
				return
			}

			slog.Info("recipe deleted", "recipe_uuid", recipeUUID)

			w.WriteHeader(http.StatusNoContent)
		default:
			methodNotAllowed(w)
		}
	}
}
