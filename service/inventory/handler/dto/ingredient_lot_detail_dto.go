package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateIngredientLotMaltDetailRequest struct {
	IngredientLotID int64    `json:"ingredient_lot_id"`
	MoisturePercent *float64 `json:"moisture_percent"`
}

func (r CreateIngredientLotMaltDetailRequest) Validate() error {
	if r.IngredientLotID <= 0 {
		return fmt.Errorf("ingredient_lot_id is required")
	}
	if r.MoisturePercent != nil && (*r.MoisturePercent < 0 || *r.MoisturePercent > 100) {
		return fmt.Errorf("moisture_percent must be between 0 and 100")
	}

	return nil
}

type IngredientLotMaltDetailResponse struct {
	ID              int64      `json:"id"`
	UUID            string     `json:"uuid"`
	IngredientLotID int64      `json:"ingredient_lot_id"`
	MoisturePercent *float64   `json:"moisture_percent,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientLotMaltDetailResponse(detail storage.IngredientLotMaltDetail) IngredientLotMaltDetailResponse {
	return IngredientLotMaltDetailResponse{
		ID:              detail.ID,
		UUID:            detail.UUID.String(),
		IngredientLotID: detail.IngredientLotID,
		MoisturePercent: detail.MoisturePercent,
		CreatedAt:       detail.CreatedAt,
		UpdatedAt:       detail.UpdatedAt,
		DeletedAt:       detail.DeletedAt,
	}
}

type CreateIngredientLotHopDetailRequest struct {
	IngredientLotID int64    `json:"ingredient_lot_id"`
	AlphaAcid       *float64 `json:"alpha_acid"`
	BetaAcid        *float64 `json:"beta_acid"`
}

func (r CreateIngredientLotHopDetailRequest) Validate() error {
	if r.IngredientLotID <= 0 {
		return fmt.Errorf("ingredient_lot_id is required")
	}
	if r.AlphaAcid != nil && (*r.AlphaAcid < 0 || *r.AlphaAcid > 100) {
		return fmt.Errorf("alpha_acid must be between 0 and 100")
	}
	if r.BetaAcid != nil && (*r.BetaAcid < 0 || *r.BetaAcid > 100) {
		return fmt.Errorf("beta_acid must be between 0 and 100")
	}

	return nil
}

type IngredientLotHopDetailResponse struct {
	ID              int64      `json:"id"`
	UUID            string     `json:"uuid"`
	IngredientLotID int64      `json:"ingredient_lot_id"`
	AlphaAcid       *float64   `json:"alpha_acid,omitempty"`
	BetaAcid        *float64   `json:"beta_acid,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientLotHopDetailResponse(detail storage.IngredientLotHopDetail) IngredientLotHopDetailResponse {
	return IngredientLotHopDetailResponse{
		ID:              detail.ID,
		UUID:            detail.UUID.String(),
		IngredientLotID: detail.IngredientLotID,
		AlphaAcid:       detail.AlphaAcid,
		BetaAcid:        detail.BetaAcid,
		CreatedAt:       detail.CreatedAt,
		UpdatedAt:       detail.UpdatedAt,
		DeletedAt:       detail.DeletedAt,
	}
}

type CreateIngredientLotYeastDetailRequest struct {
	IngredientLotID int64    `json:"ingredient_lot_id"`
	Viability       *float64 `json:"viability_percent"`
	Generation      *int     `json:"generation"`
}

func (r CreateIngredientLotYeastDetailRequest) Validate() error {
	if r.IngredientLotID <= 0 {
		return fmt.Errorf("ingredient_lot_id is required")
	}
	if r.Viability != nil && (*r.Viability < 0 || *r.Viability > 100) {
		return fmt.Errorf("viability_percent must be between 0 and 100")
	}
	if r.Generation != nil && *r.Generation < 0 {
		return fmt.Errorf("generation must be greater than or equal to zero")
	}

	return nil
}

type IngredientLotYeastDetailResponse struct {
	ID              int64      `json:"id"`
	UUID            string     `json:"uuid"`
	IngredientLotID int64      `json:"ingredient_lot_id"`
	Viability       *float64   `json:"viability_percent,omitempty"`
	Generation      *int       `json:"generation,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientLotYeastDetailResponse(detail storage.IngredientLotYeastDetail) IngredientLotYeastDetailResponse {
	return IngredientLotYeastDetailResponse{
		ID:              detail.ID,
		UUID:            detail.UUID.String(),
		IngredientLotID: detail.IngredientLotID,
		Viability:       detail.ViabilityPercent,
		Generation:      detail.Generation,
		CreatedAt:       detail.CreatedAt,
		UpdatedAt:       detail.UpdatedAt,
		DeletedAt:       detail.DeletedAt,
	}
}
