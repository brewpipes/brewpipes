package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateRecipeRequest struct {
	Name                string   `json:"name"`
	StyleID             *int64   `json:"style_id,omitempty"`
	StyleName           *string  `json:"style_name,omitempty"`
	Notes               *string  `json:"notes,omitempty"`
	BatchSize           *float64 `json:"batch_size,omitempty"`
	BatchSizeUnit       *string  `json:"batch_size_unit,omitempty"`
	TargetOG            *float64 `json:"target_og,omitempty"`
	TargetOGMin         *float64 `json:"target_og_min,omitempty"`
	TargetOGMax         *float64 `json:"target_og_max,omitempty"`
	TargetFG            *float64 `json:"target_fg,omitempty"`
	TargetFGMin         *float64 `json:"target_fg_min,omitempty"`
	TargetFGMax         *float64 `json:"target_fg_max,omitempty"`
	TargetIBU           *float64 `json:"target_ibu,omitempty"`
	TargetIBUMin        *float64 `json:"target_ibu_min,omitempty"`
	TargetIBUMax        *float64 `json:"target_ibu_max,omitempty"`
	TargetSRM           *float64 `json:"target_srm,omitempty"`
	TargetSRMMin        *float64 `json:"target_srm_min,omitempty"`
	TargetSRMMax        *float64 `json:"target_srm_max,omitempty"`
	TargetCarbonation   *float64 `json:"target_carbonation,omitempty"`
	IBUMethod           *string  `json:"ibu_method,omitempty"`
	BrewhouseEfficiency *float64 `json:"brewhouse_efficiency,omitempty"`
}

func (r CreateRecipeRequest) Validate() error {
	if err := validateRequired(r.Name, "name"); err != nil {
		return err
	}
	if r.IBUMethod != nil {
		if err := validateIBUMethod(*r.IBUMethod); err != nil {
			return err
		}
	}
	if r.BatchSize != nil && *r.BatchSize <= 0 {
		return errPositiveRequired("batch_size")
	}
	if r.BrewhouseEfficiency != nil && (*r.BrewhouseEfficiency < 0 || *r.BrewhouseEfficiency > 100) {
		return errRangeRequired("brewhouse_efficiency", 0, 100)
	}
	if r.TargetCarbonation != nil && (*r.TargetCarbonation < 0 || *r.TargetCarbonation > 5) {
		return errRangeRequired("target_carbonation", 0, 5)
	}
	return nil
}

type UpdateRecipeRequest struct {
	Name                string   `json:"name"`
	StyleID             *int64   `json:"style_id,omitempty"`
	StyleName           *string  `json:"style_name,omitempty"`
	Notes               *string  `json:"notes,omitempty"`
	BatchSize           *float64 `json:"batch_size,omitempty"`
	BatchSizeUnit       *string  `json:"batch_size_unit,omitempty"`
	TargetOG            *float64 `json:"target_og,omitempty"`
	TargetOGMin         *float64 `json:"target_og_min,omitempty"`
	TargetOGMax         *float64 `json:"target_og_max,omitempty"`
	TargetFG            *float64 `json:"target_fg,omitempty"`
	TargetFGMin         *float64 `json:"target_fg_min,omitempty"`
	TargetFGMax         *float64 `json:"target_fg_max,omitempty"`
	TargetIBU           *float64 `json:"target_ibu,omitempty"`
	TargetIBUMin        *float64 `json:"target_ibu_min,omitempty"`
	TargetIBUMax        *float64 `json:"target_ibu_max,omitempty"`
	TargetSRM           *float64 `json:"target_srm,omitempty"`
	TargetSRMMin        *float64 `json:"target_srm_min,omitempty"`
	TargetSRMMax        *float64 `json:"target_srm_max,omitempty"`
	TargetCarbonation   *float64 `json:"target_carbonation,omitempty"`
	IBUMethod           *string  `json:"ibu_method,omitempty"`
	BrewhouseEfficiency *float64 `json:"brewhouse_efficiency,omitempty"`
}

func (r UpdateRecipeRequest) Validate() error {
	if err := validateRequired(r.Name, "name"); err != nil {
		return err
	}
	if r.IBUMethod != nil {
		if err := validateIBUMethod(*r.IBUMethod); err != nil {
			return err
		}
	}
	if r.BatchSize != nil && *r.BatchSize <= 0 {
		return errPositiveRequired("batch_size")
	}
	if r.BrewhouseEfficiency != nil && (*r.BrewhouseEfficiency < 0 || *r.BrewhouseEfficiency > 100) {
		return errRangeRequired("brewhouse_efficiency", 0, 100)
	}
	if r.TargetCarbonation != nil && (*r.TargetCarbonation < 0 || *r.TargetCarbonation > 5) {
		return errRangeRequired("target_carbonation", 0, 5)
	}
	return nil
}

type RecipeResponse struct {
	UUID                string     `json:"uuid"`
	Name                string     `json:"name"`
	StyleID             *int64     `json:"style_id,omitempty"`
	StyleName           *string    `json:"style_name,omitempty"`
	Notes               *string    `json:"notes,omitempty"`
	BatchSize           *float64   `json:"batch_size,omitempty"`
	BatchSizeUnit       *string    `json:"batch_size_unit,omitempty"`
	TargetOG            *float64   `json:"target_og,omitempty"`
	TargetOGMin         *float64   `json:"target_og_min,omitempty"`
	TargetOGMax         *float64   `json:"target_og_max,omitempty"`
	TargetFG            *float64   `json:"target_fg,omitempty"`
	TargetFGMin         *float64   `json:"target_fg_min,omitempty"`
	TargetFGMax         *float64   `json:"target_fg_max,omitempty"`
	TargetIBU           *float64   `json:"target_ibu,omitempty"`
	TargetIBUMin        *float64   `json:"target_ibu_min,omitempty"`
	TargetIBUMax        *float64   `json:"target_ibu_max,omitempty"`
	TargetSRM           *float64   `json:"target_srm,omitempty"`
	TargetSRMMin        *float64   `json:"target_srm_min,omitempty"`
	TargetSRMMax        *float64   `json:"target_srm_max,omitempty"`
	TargetCarbonation   *float64   `json:"target_carbonation,omitempty"`
	IBUMethod           *string    `json:"ibu_method,omitempty"`
	BrewhouseEfficiency *float64   `json:"brewhouse_efficiency,omitempty"`
	TargetABV           *float64   `json:"target_abv,omitempty"` // Calculated from OG/FG
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

// calculateTargetABV calculates ABV from OG and FG using (OG - FG) * 131.25
func calculateTargetABV(og, fg *float64) *float64 {
	if og == nil || fg == nil {
		return nil
	}
	abv := (*og - *fg) * 131.25
	return &abv
}

func NewRecipeResponse(recipe storage.Recipe) RecipeResponse {
	return RecipeResponse{
		UUID:                recipe.UUID.String(),
		Name:                recipe.Name,
		StyleID:             recipe.StyleID,
		StyleName:           recipe.StyleName,
		Notes:               recipe.Notes,
		BatchSize:           recipe.BatchSize,
		BatchSizeUnit:       recipe.BatchSizeUnit,
		TargetOG:            recipe.TargetOG,
		TargetOGMin:         recipe.TargetOGMin,
		TargetOGMax:         recipe.TargetOGMax,
		TargetFG:            recipe.TargetFG,
		TargetFGMin:         recipe.TargetFGMin,
		TargetFGMax:         recipe.TargetFGMax,
		TargetIBU:           recipe.TargetIBU,
		TargetIBUMin:        recipe.TargetIBUMin,
		TargetIBUMax:        recipe.TargetIBUMax,
		TargetSRM:           recipe.TargetSRM,
		TargetSRMMin:        recipe.TargetSRMMin,
		TargetSRMMax:        recipe.TargetSRMMax,
		TargetCarbonation:   recipe.TargetCarbonation,
		IBUMethod:           recipe.IBUMethod,
		BrewhouseEfficiency: recipe.BrewhouseEfficiency,
		TargetABV:           calculateTargetABV(recipe.TargetOG, recipe.TargetFG),
		CreatedAt:           recipe.CreatedAt,
		UpdatedAt:           recipe.UpdatedAt,
		DeletedAt:           recipe.DeletedAt,
	}
}

func NewRecipesResponse(recipes []storage.Recipe) []RecipeResponse {
	resp := make([]RecipeResponse, 0, len(recipes))
	for _, recipe := range recipes {
		resp = append(resp, NewRecipeResponse(recipe))
	}
	return resp
}
