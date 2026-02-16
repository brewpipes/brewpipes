package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryAdjustmentRequest struct {
	IngredientLotUUID *string    `json:"ingredient_lot_uuid"`
	BeerLotUUID       *string    `json:"beer_lot_uuid"`
	StockLocationUUID string     `json:"stock_location_uuid"`
	Amount            int64      `json:"amount"`
	AmountUnit        string     `json:"amount_unit"`
	Reason            string     `json:"reason"`
	AdjustedAt        *time.Time `json:"adjusted_at"`
	Notes             *string    `json:"notes"`
}

func (r CreateInventoryAdjustmentRequest) Validate() error {
	if (r.IngredientLotUUID == nil && r.BeerLotUUID == nil) || (r.IngredientLotUUID != nil && r.BeerLotUUID != nil) {
		return fmt.Errorf("exactly one of ingredient_lot_uuid or beer_lot_uuid is required")
	}
	if r.IngredientLotUUID != nil {
		if err := validate.Required(*r.IngredientLotUUID, "ingredient_lot_uuid"); err != nil {
			return err
		}
	}
	if r.BeerLotUUID != nil {
		if err := validate.Required(*r.BeerLotUUID, "beer_lot_uuid"); err != nil {
			return err
		}
	}
	if err := validate.Required(r.StockLocationUUID, "stock_location_uuid"); err != nil {
		return err
	}
	if r.Amount == 0 {
		return fmt.Errorf("amount must not be zero")
	}
	if err := validate.Required(r.AmountUnit, "amount_unit"); err != nil {
		return err
	}
	if err := validate.Required(r.Reason, "reason"); err != nil {
		return err
	}
	if err := validateAdjustmentReason(r.Reason); err != nil {
		return err
	}

	return nil
}

type InventoryAdjustmentResponse struct {
	UUID       string     `json:"uuid"`
	Reason     string     `json:"reason"`
	AdjustedAt time.Time  `json:"adjusted_at"`
	Notes      *string    `json:"notes,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

func NewInventoryAdjustmentResponse(adjustment storage.InventoryAdjustment) InventoryAdjustmentResponse {
	return InventoryAdjustmentResponse{
		UUID:       adjustment.UUID.String(),
		Reason:     adjustment.Reason,
		AdjustedAt: adjustment.AdjustedAt,
		Notes:      adjustment.Notes,
		CreatedAt:  adjustment.CreatedAt,
		UpdatedAt:  adjustment.UpdatedAt,
		DeletedAt:  adjustment.DeletedAt,
	}
}

func NewInventoryAdjustmentsResponse(adjustments []storage.InventoryAdjustment) []InventoryAdjustmentResponse {
	resp := make([]InventoryAdjustmentResponse, 0, len(adjustments))
	for _, adjustment := range adjustments {
		resp = append(resp, NewInventoryAdjustmentResponse(adjustment))
	}
	return resp
}
