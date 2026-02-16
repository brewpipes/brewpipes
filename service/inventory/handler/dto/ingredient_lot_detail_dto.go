package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateIngredientLotMaltDetailRequest struct {
	IngredientLotUUID string   `json:"ingredient_lot_uuid"`
	MoisturePercent   *float64 `json:"moisture_percent"`
}

func (r CreateIngredientLotMaltDetailRequest) Validate() error {
	if err := validate.Required(r.IngredientLotUUID, "ingredient_lot_uuid"); err != nil {
		return err
	}
	if r.MoisturePercent != nil && (*r.MoisturePercent < 0 || *r.MoisturePercent > 100) {
		return fmt.Errorf("moisture_percent must be between 0 and 100")
	}

	return nil
}

type UpdateIngredientLotMaltDetailRequest struct {
	MoisturePercent *float64 `json:"moisture_percent"`
}

func (r UpdateIngredientLotMaltDetailRequest) Validate() error {
	if r.MoisturePercent != nil && (*r.MoisturePercent < 0 || *r.MoisturePercent > 100) {
		return fmt.Errorf("moisture_percent must be between 0 and 100")
	}

	return nil
}

type IngredientLotMaltDetailResponse struct {
	UUID              string     `json:"uuid"`
	IngredientLotUUID string     `json:"ingredient_lot_uuid"`
	MoisturePercent   *float64   `json:"moisture_percent,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientLotMaltDetailResponse(detail storage.IngredientLotMaltDetail) IngredientLotMaltDetailResponse {
	return IngredientLotMaltDetailResponse{
		UUID:              detail.UUID.String(),
		IngredientLotUUID: detail.IngredientLotUUID,
		MoisturePercent:   detail.MoisturePercent,
		CreatedAt:         detail.CreatedAt,
		UpdatedAt:         detail.UpdatedAt,
		DeletedAt:         detail.DeletedAt,
	}
}

type CreateIngredientLotHopDetailRequest struct {
	IngredientLotUUID string   `json:"ingredient_lot_uuid"`
	AlphaAcid         *float64 `json:"alpha_acid"`
	BetaAcid          *float64 `json:"beta_acid"`
}

func (r CreateIngredientLotHopDetailRequest) Validate() error {
	if err := validate.Required(r.IngredientLotUUID, "ingredient_lot_uuid"); err != nil {
		return err
	}
	if r.AlphaAcid != nil && (*r.AlphaAcid < 0 || *r.AlphaAcid > 100) {
		return fmt.Errorf("alpha_acid must be between 0 and 100")
	}
	if r.BetaAcid != nil && (*r.BetaAcid < 0 || *r.BetaAcid > 100) {
		return fmt.Errorf("beta_acid must be between 0 and 100")
	}

	return nil
}

type UpdateIngredientLotHopDetailRequest struct {
	AlphaAcid *float64 `json:"alpha_acid"`
	BetaAcid  *float64 `json:"beta_acid"`
}

func (r UpdateIngredientLotHopDetailRequest) Validate() error {
	if r.AlphaAcid != nil && (*r.AlphaAcid < 0 || *r.AlphaAcid > 100) {
		return fmt.Errorf("alpha_acid must be between 0 and 100")
	}
	if r.BetaAcid != nil && (*r.BetaAcid < 0 || *r.BetaAcid > 100) {
		return fmt.Errorf("beta_acid must be between 0 and 100")
	}

	return nil
}

type IngredientLotHopDetailResponse struct {
	UUID              string     `json:"uuid"`
	IngredientLotUUID string     `json:"ingredient_lot_uuid"`
	AlphaAcid         *float64   `json:"alpha_acid,omitempty"`
	BetaAcid          *float64   `json:"beta_acid,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientLotHopDetailResponse(detail storage.IngredientLotHopDetail) IngredientLotHopDetailResponse {
	return IngredientLotHopDetailResponse{
		UUID:              detail.UUID.String(),
		IngredientLotUUID: detail.IngredientLotUUID,
		AlphaAcid:         detail.AlphaAcid,
		BetaAcid:          detail.BetaAcid,
		CreatedAt:         detail.CreatedAt,
		UpdatedAt:         detail.UpdatedAt,
		DeletedAt:         detail.DeletedAt,
	}
}

type CreateIngredientLotYeastDetailRequest struct {
	IngredientLotUUID string   `json:"ingredient_lot_uuid"`
	Viability         *float64 `json:"viability_percent"`
	Generation        *int     `json:"generation"`
}

func (r CreateIngredientLotYeastDetailRequest) Validate() error {
	if err := validate.Required(r.IngredientLotUUID, "ingredient_lot_uuid"); err != nil {
		return err
	}
	if r.Viability != nil && (*r.Viability < 0 || *r.Viability > 100) {
		return fmt.Errorf("viability_percent must be between 0 and 100")
	}
	if r.Generation != nil && *r.Generation < 0 {
		return fmt.Errorf("generation must be greater than or equal to zero")
	}

	return nil
}

type UpdateIngredientLotYeastDetailRequest struct {
	Viability  *float64 `json:"viability_percent"`
	Generation *int     `json:"generation"`
}

func (r UpdateIngredientLotYeastDetailRequest) Validate() error {
	if r.Viability != nil && (*r.Viability < 0 || *r.Viability > 100) {
		return fmt.Errorf("viability_percent must be between 0 and 100")
	}
	if r.Generation != nil && *r.Generation < 0 {
		return fmt.Errorf("generation must be greater than or equal to zero")
	}

	return nil
}

type IngredientLotYeastDetailResponse struct {
	UUID              string     `json:"uuid"`
	IngredientLotUUID string     `json:"ingredient_lot_uuid"`
	Viability         *float64   `json:"viability_percent,omitempty"`
	Generation        *int       `json:"generation,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientLotYeastDetailResponse(detail storage.IngredientLotYeastDetail) IngredientLotYeastDetailResponse {
	return IngredientLotYeastDetailResponse{
		UUID:              detail.UUID.String(),
		IngredientLotUUID: detail.IngredientLotUUID,
		Viability:         detail.ViabilityPercent,
		Generation:        detail.Generation,
		CreatedAt:         detail.CreatedAt,
		UpdatedAt:         detail.UpdatedAt,
		DeletedAt:         detail.DeletedAt,
	}
}
