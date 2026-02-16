package dto

import (
	"github.com/brewpipes/brewpipes/internal/uuidutil"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// BatchIngredientLotResponse is the API response for a batch ingredient lot.
type BatchIngredientLotResponse struct {
	IngredientLotUUID     string  `json:"ingredient_lot_uuid"`
	IngredientUUID        string  `json:"ingredient_uuid"`
	IngredientName        string  `json:"ingredient_name"`
	IngredientCategory    string  `json:"ingredient_category"`
	BreweryLotCode        *string `json:"brewery_lot_code,omitempty"`
	PurchaseOrderLineUUID *string `json:"purchase_order_line_uuid,omitempty"`
	ReceivedUnit          string  `json:"received_unit"`
}

// NewBatchIngredientLotResponse creates a BatchIngredientLotResponse from a storage model.
func NewBatchIngredientLotResponse(lot storage.BatchIngredientLot) BatchIngredientLotResponse {
	return BatchIngredientLotResponse{
		IngredientLotUUID:     lot.IngredientLotUUID,
		IngredientUUID:        lot.IngredientUUID,
		IngredientName:        lot.IngredientName,
		IngredientCategory:    lot.IngredientCategory,
		BreweryLotCode:        lot.BreweryLotCode,
		PurchaseOrderLineUUID: uuidutil.ToStringPointer(lot.PurchaseOrderLineUUID),
		ReceivedUnit:          lot.ReceivedUnit,
	}
}

// NewBatchIngredientLotsResponse creates a slice of BatchIngredientLotResponse from storage models.
func NewBatchIngredientLotsResponse(lots []storage.BatchIngredientLot) []BatchIngredientLotResponse {
	resp := make([]BatchIngredientLotResponse, 0, len(lots))
	for _, lot := range lots {
		resp = append(resp, NewBatchIngredientLotResponse(lot))
	}
	return resp
}
