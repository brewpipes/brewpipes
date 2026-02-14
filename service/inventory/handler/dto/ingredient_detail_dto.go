package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateIngredientMaltDetailRequest struct {
	IngredientUUID string   `json:"ingredient_uuid"`
	MaltsterName   *string  `json:"maltster_name"`
	Variety        *string  `json:"variety"`
	Lovibond       *float64 `json:"lovibond"`
	SRM            *float64 `json:"srm"`
	DiastaticPower *float64 `json:"diastatic_power"`
}

func (r CreateIngredientMaltDetailRequest) Validate() error {
	if err := validate.Required(r.IngredientUUID, "ingredient_uuid"); err != nil {
		return err
	}
	if r.Lovibond != nil && *r.Lovibond < 0 {
		return fmt.Errorf("lovibond must be greater than or equal to zero")
	}
	if r.SRM != nil && *r.SRM < 0 {
		return fmt.Errorf("srm must be greater than or equal to zero")
	}
	if r.DiastaticPower != nil && *r.DiastaticPower < 0 {
		return fmt.Errorf("diastatic_power must be greater than or equal to zero")
	}

	return nil
}

type IngredientMaltDetailResponse struct {
	UUID           string     `json:"uuid"`
	IngredientUUID string     `json:"ingredient_uuid"`
	MaltsterName   *string    `json:"maltster_name,omitempty"`
	Variety        *string    `json:"variety,omitempty"`
	Lovibond       *float64   `json:"lovibond,omitempty"`
	SRM            *float64   `json:"srm,omitempty"`
	DiastaticPower *float64   `json:"diastatic_power,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientMaltDetailResponse(detail storage.IngredientMaltDetail) IngredientMaltDetailResponse {
	return IngredientMaltDetailResponse{
		UUID:           detail.UUID.String(),
		IngredientUUID: detail.IngredientUUID,
		MaltsterName:   detail.MaltsterName,
		Variety:        detail.Variety,
		Lovibond:       detail.Lovibond,
		SRM:            detail.SRM,
		DiastaticPower: detail.DiastaticPower,
		CreatedAt:      detail.CreatedAt,
		UpdatedAt:      detail.UpdatedAt,
		DeletedAt:      detail.DeletedAt,
	}
}

type CreateIngredientHopDetailRequest struct {
	IngredientUUID string   `json:"ingredient_uuid"`
	ProducerName   *string  `json:"producer_name"`
	Variety        *string  `json:"variety"`
	CropYear       *int     `json:"crop_year"`
	Form           *string  `json:"form"`
	AlphaAcid      *float64 `json:"alpha_acid"`
	BetaAcid       *float64 `json:"beta_acid"`
}

func (r CreateIngredientHopDetailRequest) Validate() error {
	if err := validate.Required(r.IngredientUUID, "ingredient_uuid"); err != nil {
		return err
	}
	if r.CropYear != nil && *r.CropYear < 1900 {
		return fmt.Errorf("crop_year must be 1900 or later")
	}
	if r.Form != nil {
		if err := validateHopForm(*r.Form); err != nil {
			return err
		}
	}
	if r.AlphaAcid != nil && (*r.AlphaAcid < 0 || *r.AlphaAcid > 100) {
		return fmt.Errorf("alpha_acid must be between 0 and 100")
	}
	if r.BetaAcid != nil && (*r.BetaAcid < 0 || *r.BetaAcid > 100) {
		return fmt.Errorf("beta_acid must be between 0 and 100")
	}

	return nil
}

type IngredientHopDetailResponse struct {
	UUID           string     `json:"uuid"`
	IngredientUUID string     `json:"ingredient_uuid"`
	ProducerName   *string    `json:"producer_name,omitempty"`
	Variety        *string    `json:"variety,omitempty"`
	CropYear       *int       `json:"crop_year,omitempty"`
	Form           *string    `json:"form,omitempty"`
	AlphaAcid      *float64   `json:"alpha_acid,omitempty"`
	BetaAcid       *float64   `json:"beta_acid,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientHopDetailResponse(detail storage.IngredientHopDetail) IngredientHopDetailResponse {
	return IngredientHopDetailResponse{
		UUID:           detail.UUID.String(),
		IngredientUUID: detail.IngredientUUID,
		ProducerName:   detail.ProducerName,
		Variety:        detail.Variety,
		CropYear:       detail.CropYear,
		Form:           detail.Form,
		AlphaAcid:      detail.AlphaAcid,
		BetaAcid:       detail.BetaAcid,
		CreatedAt:      detail.CreatedAt,
		UpdatedAt:      detail.UpdatedAt,
		DeletedAt:      detail.DeletedAt,
	}
}

type CreateIngredientYeastDetailRequest struct {
	IngredientUUID string  `json:"ingredient_uuid"`
	LabName        *string `json:"lab_name"`
	Strain         *string `json:"strain"`
	Form           *string `json:"form"`
}

func (r CreateIngredientYeastDetailRequest) Validate() error {
	if err := validate.Required(r.IngredientUUID, "ingredient_uuid"); err != nil {
		return err
	}
	if r.Form != nil {
		if err := validateYeastForm(*r.Form); err != nil {
			return err
		}
	}

	return nil
}

type IngredientYeastDetailResponse struct {
	UUID           string     `json:"uuid"`
	IngredientUUID string     `json:"ingredient_uuid"`
	LabName        *string    `json:"lab_name,omitempty"`
	Strain         *string    `json:"strain,omitempty"`
	Form           *string    `json:"form,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientYeastDetailResponse(detail storage.IngredientYeastDetail) IngredientYeastDetailResponse {
	return IngredientYeastDetailResponse{
		UUID:           detail.UUID.String(),
		IngredientUUID: detail.IngredientUUID,
		LabName:        detail.LabName,
		Strain:         detail.Strain,
		Form:           detail.Form,
		CreatedAt:      detail.CreatedAt,
		UpdatedAt:      detail.UpdatedAt,
		DeletedAt:      detail.DeletedAt,
	}
}
