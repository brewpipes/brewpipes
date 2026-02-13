package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBatchRequest struct {
	ShortName  string     `json:"short_name"`
	BrewDate   *time.Time `json:"brew_date"`
	Notes      *string    `json:"notes"`
	RecipeUUID *string    `json:"recipe_uuid"`
}

func (r CreateBatchRequest) Validate() error {
	return validateRequired(r.ShortName, "short_name")
}

type UpdateBatchRequest struct {
	ShortName  string     `json:"short_name"`
	BrewDate   *time.Time `json:"brew_date"`
	Notes      *string    `json:"notes"`
	RecipeUUID *string    `json:"recipe_uuid"`
}

func (r UpdateBatchRequest) Validate() error {
	return validateRequired(r.ShortName, "short_name")
}

type BatchResponse struct {
	UUID         string     `json:"uuid"`
	ShortName    string     `json:"short_name"`
	BrewDate     *time.Time `json:"brew_date,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	RecipeUUID   *string    `json:"recipe_uuid,omitempty"`
	RecipeName   *string    `json:"recipe_name,omitempty"`
	CurrentPhase *string    `json:"current_phase,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewBatchResponse(batch storage.Batch) BatchResponse {
	return BatchResponse{
		UUID:         batch.UUID.String(),
		ShortName:    batch.ShortName,
		BrewDate:     batch.BrewDate,
		Notes:        batch.Notes,
		RecipeUUID:   batch.RecipeUUID,
		RecipeName:   batch.RecipeName,
		CurrentPhase: batch.CurrentPhase,
		CreatedAt:    batch.CreatedAt,
		UpdatedAt:    batch.UpdatedAt,
		DeletedAt:    batch.DeletedAt,
	}
}

func NewBatchesResponse(batches []storage.Batch) []BatchResponse {
	resp := make([]BatchResponse, 0, len(batches))
	for _, batch := range batches {
		resp = append(resp, NewBatchResponse(batch))
	}
	return resp
}
