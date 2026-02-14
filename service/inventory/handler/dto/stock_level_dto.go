package dto

import (
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// StockLevelLocationResponse represents stock at a specific location.
type StockLevelLocationResponse struct {
	LocationUUID string  `json:"location_uuid"`
	LocationName string  `json:"location_name"`
	Quantity     float64 `json:"quantity"`
}

// StockLevelResponse represents the aggregated stock level for an ingredient.
type StockLevelResponse struct {
	IngredientUUID string                       `json:"ingredient_uuid"`
	IngredientName string                       `json:"ingredient_name"`
	Category       string                       `json:"category"`
	DefaultUnit    string                       `json:"default_unit"`
	TotalOnHand    float64                      `json:"total_on_hand"`
	Locations      []StockLevelLocationResponse `json:"locations"`
}

// NewStockLevelResponse converts a storage StockLevel to a response DTO.
func NewStockLevelResponse(sl storage.StockLevel) StockLevelResponse {
	locations := make([]StockLevelLocationResponse, 0, len(sl.Locations))
	for _, loc := range sl.Locations {
		locations = append(locations, StockLevelLocationResponse{
			LocationUUID: loc.LocationUUID,
			LocationName: loc.LocationName,
			Quantity:     loc.Quantity,
		})
	}

	return StockLevelResponse{
		IngredientUUID: sl.IngredientUUID,
		IngredientName: sl.IngredientName,
		Category:       sl.Category,
		DefaultUnit:    sl.DefaultUnit,
		TotalOnHand:    sl.TotalOnHand,
		Locations:      locations,
	}
}

// NewStockLevelsResponse converts a slice of storage StockLevels to response DTOs.
func NewStockLevelsResponse(levels []storage.StockLevel) []StockLevelResponse {
	resp := make([]StockLevelResponse, 0, len(levels))
	for _, sl := range levels {
		resp = append(resp, NewStockLevelResponse(sl))
	}
	return resp
}
