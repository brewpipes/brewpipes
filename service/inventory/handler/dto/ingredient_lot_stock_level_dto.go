package dto

import (
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// IngredientLotStockLevelResponse represents the current stock level for an
// ingredient lot at a specific location.
type IngredientLotStockLevelResponse struct {
	IngredientLotUUID  string  `json:"ingredient_lot_uuid"`
	IngredientUUID     string  `json:"ingredient_uuid"`
	IngredientName     string  `json:"ingredient_name"`
	IngredientCategory string  `json:"ingredient_category"`
	BreweryLotCode     *string `json:"brewery_lot_code,omitempty"`
	ReceivedAt         string  `json:"received_at"`
	ReceivedAmount     int64   `json:"received_amount"`
	ReceivedUnit       string  `json:"received_unit"`
	StockLocationUUID  string  `json:"stock_location_uuid"`
	StockLocationName  string  `json:"stock_location_name"`
	CurrentAmount      int64   `json:"current_amount"`
	CurrentUnit        string  `json:"current_unit"`
}

// NewIngredientLotStockLevelResponse converts a storage IngredientLotStockLevel to a response DTO.
func NewIngredientLotStockLevelResponse(level storage.IngredientLotStockLevel) IngredientLotStockLevelResponse {
	return IngredientLotStockLevelResponse{
		IngredientLotUUID:  level.IngredientLotUUID,
		IngredientUUID:     level.IngredientUUID,
		IngredientName:     level.IngredientName,
		IngredientCategory: level.IngredientCategory,
		BreweryLotCode:     level.BreweryLotCode,
		ReceivedAt:         level.ReceivedAt,
		ReceivedAmount:     level.ReceivedAmount,
		ReceivedUnit:       level.ReceivedUnit,
		StockLocationUUID:  level.LocationUUID,
		StockLocationName:  level.LocationName,
		CurrentAmount:      level.CurrentAmount,
		CurrentUnit:        level.CurrentUnit,
	}
}

// NewIngredientLotStockLevelsResponse converts a slice of storage IngredientLotStockLevels to response DTOs.
func NewIngredientLotStockLevelsResponse(levels []storage.IngredientLotStockLevel) []IngredientLotStockLevelResponse {
	resp := make([]IngredientLotStockLevelResponse, 0, len(levels))
	for _, level := range levels {
		resp = append(resp, NewIngredientLotStockLevelResponse(level))
	}
	return resp
}
