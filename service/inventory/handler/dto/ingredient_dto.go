package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateIngredientRequest struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	DefaultUnit string  `json:"default_unit"`
	Description *string `json:"description"`
}

func (r CreateIngredientRequest) Validate() error {
	if err := validateRequired(r.Name, "name"); err != nil {
		return err
	}
	if err := validateIngredientCategory(r.Category); err != nil {
		return err
	}
	if err := validateRequired(r.DefaultUnit, "default_unit"); err != nil {
		return err
	}

	return nil
}

type IngredientResponse struct {
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	DefaultUnit string     `json:"default_unit"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientResponse(ingredient storage.Ingredient) IngredientResponse {
	return IngredientResponse{
		UUID:        ingredient.UUID.String(),
		Name:        ingredient.Name,
		Category:    ingredient.Category,
		DefaultUnit: ingredient.DefaultUnit,
		Description: ingredient.Description,
		CreatedAt:   ingredient.CreatedAt,
		UpdatedAt:   ingredient.UpdatedAt,
		DeletedAt:   ingredient.DeletedAt,
	}
}

func NewIngredientsResponse(ingredients []storage.Ingredient) []IngredientResponse {
	resp := make([]IngredientResponse, 0, len(ingredients))
	for _, ingredient := range ingredients {
		resp = append(resp, NewIngredientResponse(ingredient))
	}
	return resp
}
