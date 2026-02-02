package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateRecipeRequest struct {
	Name      string  `json:"name"`
	StyleID   *int64  `json:"style_id"`
	StyleName *string `json:"style_name"`
	Notes     *string `json:"notes"`
}

func (r CreateRecipeRequest) Validate() error {
	return validateRequired(r.Name, "name")
}

type UpdateRecipeRequest struct {
	Name      string  `json:"name"`
	StyleID   *int64  `json:"style_id"`
	StyleName *string `json:"style_name"`
	Notes     *string `json:"notes"`
}

func (r UpdateRecipeRequest) Validate() error {
	return validateRequired(r.Name, "name")
}

type RecipeResponse struct {
	ID        int64      `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	StyleID   *int64     `json:"style_id,omitempty"`
	StyleName *string    `json:"style_name,omitempty"`
	Notes     *string    `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewRecipeResponse(recipe storage.Recipe) RecipeResponse {
	return RecipeResponse{
		ID:        recipe.ID,
		UUID:      recipe.UUID.String(),
		Name:      recipe.Name,
		StyleID:   recipe.StyleID,
		StyleName: recipe.StyleName,
		Notes:     recipe.Notes,
		CreatedAt: recipe.CreatedAt,
		UpdatedAt: recipe.UpdatedAt,
		DeletedAt: recipe.DeletedAt,
	}
}

func NewRecipesResponse(recipes []storage.Recipe) []RecipeResponse {
	resp := make([]RecipeResponse, 0, len(recipes))
	for _, recipe := range recipes {
		resp = append(resp, NewRecipeResponse(recipe))
	}
	return resp
}
