package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// validBeerLotContainers is the set of allowed container values.
var validBeerLotContainers = map[string]bool{
	"keg":     true,
	"can":     true,
	"bottle":  true,
	"cask":    true,
	"growler": true,
	"other":   true,
}

type CreateBeerLotRequest struct {
	ProductionBatchUUID string     `json:"production_batch_uuid"`
	LotCode             *string    `json:"lot_code"`
	PackagedAt          *time.Time `json:"packaged_at"`
	Notes               *string    `json:"notes"`
	PackagingRunUUID    *string    `json:"packaging_run_uuid"`
	BestBy              *string    `json:"best_by"`
	PackageFormatName   *string    `json:"package_format_name"`
	Container           *string    `json:"container"`
	VolumePerUnit       *int64     `json:"volume_per_unit"`
	VolumePerUnitUnit   *string    `json:"volume_per_unit_unit"`
	Quantity            *int       `json:"quantity"`
	StockLocationUUID   *string    `json:"stock_location_uuid"`
}

func (r CreateBeerLotRequest) Validate() error {
	if err := validate.Required(r.ProductionBatchUUID, "production_batch_uuid"); err != nil {
		return err
	}

	if r.Container != nil {
		if !validBeerLotContainers[*r.Container] {
			return fmt.Errorf("container must be one of: keg, can, bottle, cask, growler, other")
		}
	}

	if r.VolumePerUnit != nil {
		if *r.VolumePerUnit <= 0 {
			return fmt.Errorf("volume_per_unit must be greater than 0")
		}
		if r.VolumePerUnitUnit == nil {
			return fmt.Errorf("volume_per_unit_unit is required when volume_per_unit is provided")
		}
	}

	if r.Quantity != nil {
		if *r.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0")
		}
	}

	if r.StockLocationUUID != nil {
		if r.VolumePerUnit == nil || r.VolumePerUnitUnit == nil || r.Quantity == nil {
			return fmt.Errorf("volume_per_unit, volume_per_unit_unit, and quantity are required when stock_location_uuid is provided")
		}
	}

	return nil
}

type BeerLotResponse struct {
	UUID                string     `json:"uuid"`
	ProductionBatchUUID string     `json:"production_batch_uuid"`
	LotCode             *string    `json:"lot_code,omitempty"`
	PackagedAt          time.Time  `json:"packaged_at"`
	Notes               *string    `json:"notes,omitempty"`
	PackagingRunUUID    *string    `json:"packaging_run_uuid,omitempty"`
	BestBy              *string    `json:"best_by,omitempty"`
	PackageFormatName   *string    `json:"package_format_name,omitempty"`
	Container           *string    `json:"container,omitempty"`
	VolumePerUnit       *int64     `json:"volume_per_unit,omitempty"`
	VolumePerUnitUnit   *string    `json:"volume_per_unit_unit,omitempty"`
	Quantity            *int       `json:"quantity,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

func NewBeerLotResponse(lot storage.BeerLot) BeerLotResponse {
	resp := BeerLotResponse{
		UUID:                lot.UUID.String(),
		ProductionBatchUUID: lot.ProductionBatchUUID.String(),
		LotCode:             lot.LotCode,
		PackagedAt:          lot.PackagedAt,
		Notes:               lot.Notes,
		PackageFormatName:   lot.PackageFormatName,
		Container:           lot.Container,
		VolumePerUnit:       lot.VolumePerUnit,
		VolumePerUnitUnit:   lot.VolumePerUnitUnit,
		Quantity:            lot.Quantity,
		CreatedAt:           lot.CreatedAt,
		UpdatedAt:           lot.UpdatedAt,
		DeletedAt:           lot.DeletedAt,
	}

	if lot.PackagingRunUUID != nil {
		s := lot.PackagingRunUUID.String()
		resp.PackagingRunUUID = &s
	}

	if lot.BestBy != nil {
		s := lot.BestBy.Format(time.RFC3339)
		resp.BestBy = &s
	}

	return resp
}

func NewBeerLotsResponse(lots []storage.BeerLot) []BeerLotResponse {
	resp := make([]BeerLotResponse, 0, len(lots))
	for _, lot := range lots {
		resp = append(resp, NewBeerLotResponse(lot))
	}
	return resp
}
