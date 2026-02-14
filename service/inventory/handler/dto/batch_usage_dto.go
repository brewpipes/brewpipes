package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
	"github.com/gofrs/uuid/v5"
)

// BatchUsagePick represents a single ingredient pick for batch usage deduction.
type BatchUsagePick struct {
	IngredientLotUUID string `json:"ingredient_lot_uuid"`
	StockLocationUUID string `json:"stock_location_uuid"`
	Amount            int64  `json:"amount"`
	AmountUnit        string `json:"amount_unit"`
}

// CreateBatchUsageRequest is the request body for creating a batch usage deduction.
type CreateBatchUsageRequest struct {
	ProductionRefUUID *string          `json:"production_ref_uuid"`
	UsedAt            string           `json:"used_at"`
	Picks             []BatchUsagePick `json:"picks"`
	Notes             *string          `json:"notes"`
}

// Validate checks that the request is well-formed.
func (r CreateBatchUsageRequest) Validate() error {
	if r.ProductionRefUUID != nil {
		if _, err := uuid.FromString(*r.ProductionRefUUID); err != nil {
			return fmt.Errorf("production_ref_uuid must be a valid UUID")
		}
	}
	if _, err := time.Parse(time.RFC3339, r.UsedAt); err != nil {
		return fmt.Errorf("used_at must be a valid RFC3339 timestamp")
	}
	if len(r.Picks) == 0 {
		return fmt.Errorf("picks must not be empty")
	}
	seen := make(map[string]int) // "lot_uuid|location_uuid" -> first pick index
	for i, pick := range r.Picks {
		if err := validate.Required(pick.IngredientLotUUID, fmt.Sprintf("picks[%d].ingredient_lot_uuid", i)); err != nil {
			return err
		}
		if err := validate.Required(pick.StockLocationUUID, fmt.Sprintf("picks[%d].stock_location_uuid", i)); err != nil {
			return err
		}
		if pick.Amount <= 0 {
			return fmt.Errorf("picks[%d].amount must be greater than zero", i)
		}
		if err := validate.Required(pick.AmountUnit, fmt.Sprintf("picks[%d].amount_unit", i)); err != nil {
			return err
		}
		key := pick.IngredientLotUUID + "|" + pick.StockLocationUUID
		if first, exists := seen[key]; exists {
			return fmt.Errorf("picks[%d] duplicates lot+location from picks[%d]; combine into a single pick", i, first)
		}
		seen[key] = i
	}
	return nil
}

// BatchUsageResponse is the response body for a batch usage deduction.
type BatchUsageResponse struct {
	UsageUUID string                      `json:"usage_uuid"`
	Movements []InventoryMovementResponse `json:"movements"`
}

// NewBatchUsageResponse converts storage results to a response DTO.
func NewBatchUsageResponse(usage storage.InventoryUsage, movements []storage.InventoryMovement) BatchUsageResponse {
	return BatchUsageResponse{
		UsageUUID: usage.UUID.String(),
		Movements: NewInventoryMovementsResponse(movements),
	}
}
