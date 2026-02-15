package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// BeerLotStockLevelResponse represents the current stock level for a beer lot at a location.
type BeerLotStockLevelResponse struct {
	BeerLotUUID         string     `json:"beer_lot_uuid"`
	ProductionBatchUUID string     `json:"production_batch_uuid"`
	LotCode             *string    `json:"lot_code,omitempty"`
	PackagedAt          *time.Time `json:"packaged_at,omitempty"`
	BestBy              *string    `json:"best_by,omitempty"`
	PackageFormatName   *string    `json:"package_format_name,omitempty"`
	Container           *string    `json:"container,omitempty"`
	VolumePerUnit       *int64     `json:"volume_per_unit,omitempty"`
	VolumePerUnitUnit   *string    `json:"volume_per_unit_unit,omitempty"`
	InitialQuantity     *int       `json:"initial_quantity,omitempty"`
	StockLocationUUID   string     `json:"stock_location_uuid"`
	StockLocationName   string     `json:"stock_location_name"`
	CurrentVolume       int64      `json:"current_volume"`
	CurrentVolumeUnit   string     `json:"current_volume_unit"`
	CurrentQuantity     *int       `json:"current_quantity,omitempty"`
}

// NewBeerLotStockLevelResponse converts a storage BeerLotStockLevel to a response DTO.
func NewBeerLotStockLevelResponse(level storage.BeerLotStockLevel) BeerLotStockLevelResponse {
	resp := BeerLotStockLevelResponse{
		BeerLotUUID:         level.BeerLotUUID,
		ProductionBatchUUID: level.ProductionBatchUUID,
		LotCode:             level.LotCode,
		PackagedAt:          level.PackagedAt,
		PackageFormatName:   level.PackageFormatName,
		Container:           level.Container,
		VolumePerUnit:       level.VolumePerUnit,
		VolumePerUnitUnit:   level.VolumePerUnitUnit,
		InitialQuantity:     level.InitialQuantity,
		StockLocationUUID:   level.LocationUUID,
		StockLocationName:   level.LocationName,
		CurrentVolume:       level.CurrentVolume,
		CurrentVolumeUnit:   level.CurrentVolumeUnit,
		CurrentQuantity:     level.CurrentQuantity,
	}

	if level.BestBy != nil {
		s := level.BestBy.Format(time.RFC3339)
		resp.BestBy = &s
	}

	return resp
}

// NewBeerLotStockLevelsResponse converts a slice of storage BeerLotStockLevels to response DTOs.
func NewBeerLotStockLevelsResponse(levels []storage.BeerLotStockLevel) []BeerLotStockLevelResponse {
	resp := make([]BeerLotStockLevelResponse, 0, len(levels))
	for _, level := range levels {
		resp = append(resp, NewBeerLotStockLevelResponse(level))
	}
	return resp
}
