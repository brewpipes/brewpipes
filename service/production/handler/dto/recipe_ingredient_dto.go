package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/uuidutil"
	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

// RecipeIngredientRequest is the request body for creating or updating a recipe ingredient.
type RecipeIngredientRequest struct {
	IngredientUUID        *string  `json:"ingredient_uuid,omitempty"`
	IngredientType        string   `json:"ingredient_type"`
	Amount                float64  `json:"amount"`
	AmountUnit            string   `json:"amount_unit"`
	UseStage              string   `json:"use_stage"`
	UseType               *string  `json:"use_type,omitempty"`
	TimingDurationMinutes *int     `json:"timing_duration_minutes,omitempty"`
	TimingTemperatureC    *float64 `json:"timing_temperature_c,omitempty"`
	AlphaAcidAssumed      *float64 `json:"alpha_acid_assumed,omitempty"`
	ScalingFactor         *float64 `json:"scaling_factor,omitempty"`
	SortOrder             *int     `json:"sort_order,omitempty"`
	Notes                 *string  `json:"notes,omitempty"`
}

// Validate validates the recipe ingredient request.
func (r RecipeIngredientRequest) Validate() error {
	if err := validateIngredientType(r.IngredientType); err != nil {
		return err
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if err := validate.Required(r.AmountUnit, "amount_unit"); err != nil {
		return err
	}
	if err := validateUseStage(r.UseStage); err != nil {
		return err
	}
	if r.UseType != nil {
		if err := validateUseType(*r.UseType); err != nil {
			return err
		}
	}
	if r.AlphaAcidAssumed != nil && r.IngredientType != storage.IngredientTypeHop {
		return fmt.Errorf("alpha_acid_assumed can only be set for hop ingredients")
	}
	if r.ScalingFactor != nil && *r.ScalingFactor <= 0 {
		return fmt.Errorf("scaling_factor must be greater than zero")
	}
	return nil
}

// RecipeIngredientResponse is the response body for a recipe ingredient.
type RecipeIngredientResponse struct {
	UUID                  string    `json:"uuid"`
	RecipeUUID            string    `json:"recipe_uuid"`
	IngredientUUID        *string   `json:"ingredient_uuid,omitempty"`
	IngredientType        string    `json:"ingredient_type"`
	Amount                float64   `json:"amount"`
	AmountUnit            string    `json:"amount_unit"`
	UseStage              string    `json:"use_stage"`
	UseType               *string   `json:"use_type,omitempty"`
	TimingDurationMinutes *int      `json:"timing_duration_minutes,omitempty"`
	TimingTemperatureC    *float64  `json:"timing_temperature_c,omitempty"`
	AlphaAcidAssumed      *float64  `json:"alpha_acid_assumed,omitempty"`
	ScalingFactor         float64   `json:"scaling_factor"`
	SortOrder             int       `json:"sort_order"`
	Notes                 *string   `json:"notes,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// NewRecipeIngredientResponse creates a response from a storage model.
// recipeUUID is passed separately since the storage model only has the internal recipe_id.
func NewRecipeIngredientResponse(ri storage.RecipeIngredient, recipeUUID string) RecipeIngredientResponse {
	return RecipeIngredientResponse{
		UUID:                  ri.UUID.String(),
		RecipeUUID:            recipeUUID,
		IngredientUUID:        uuidutil.ToStringPointer(ri.IngredientUUID),
		IngredientType:        ri.IngredientType,
		Amount:                ri.Amount,
		AmountUnit:            ri.AmountUnit,
		UseStage:              ri.UseStage,
		UseType:               ri.UseType,
		TimingDurationMinutes: ri.TimingDurationMinutes,
		TimingTemperatureC:    ri.TimingTemperatureC,
		AlphaAcidAssumed:      ri.AlphaAcidAssumed,
		ScalingFactor:         ri.ScalingFactor,
		SortOrder:             ri.SortOrder,
		Notes:                 ri.Notes,
		CreatedAt:             ri.CreatedAt,
		UpdatedAt:             ri.UpdatedAt,
	}
}

// NewRecipeIngredientsResponse creates a slice of responses from storage models.
// recipeUUID is passed separately since the storage model only has the internal recipe_id.
func NewRecipeIngredientsResponse(ingredients []storage.RecipeIngredient, recipeUUID string) []RecipeIngredientResponse {
	resp := make([]RecipeIngredientResponse, 0, len(ingredients))
	for _, ri := range ingredients {
		resp = append(resp, NewRecipeIngredientResponse(ri, recipeUUID))
	}
	return resp
}
