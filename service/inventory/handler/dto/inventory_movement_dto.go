package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryMovementRequest struct {
	IngredientLotUUID *string    `json:"ingredient_lot_uuid"`
	BeerLotUUID       *string    `json:"beer_lot_uuid"`
	StockLocationUUID string     `json:"stock_location_uuid"`
	Direction         string     `json:"direction"`
	Reason            string     `json:"reason"`
	Amount            int64      `json:"amount"`
	AmountUnit        string     `json:"amount_unit"`
	OccurredAt        *time.Time `json:"occurred_at"`
	ReceiptUUID       *string    `json:"receipt_uuid"`
	UsageUUID         *string    `json:"usage_uuid"`
	AdjustmentUUID    *string    `json:"adjustment_uuid"`
	TransferUUID      *string    `json:"transfer_uuid"`
	Notes             *string    `json:"notes"`
}

func (r CreateInventoryMovementRequest) Validate() error {
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
	if err := validateMovementDirection(r.Direction); err != nil {
		return err
	}
	if err := validateMovementReason(r.Reason); err != nil {
		return err
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if err := validate.Required(r.AmountUnit, "amount_unit"); err != nil {
		return err
	}

	referenceCount := 0
	if r.ReceiptUUID != nil {
		if err := validate.Required(*r.ReceiptUUID, "receipt_uuid"); err != nil {
			return err
		}
		referenceCount++
	}
	if r.UsageUUID != nil {
		if err := validate.Required(*r.UsageUUID, "usage_uuid"); err != nil {
			return err
		}
		referenceCount++
	}
	if r.AdjustmentUUID != nil {
		if err := validate.Required(*r.AdjustmentUUID, "adjustment_uuid"); err != nil {
			return err
		}
		referenceCount++
	}
	if r.TransferUUID != nil {
		if err := validate.Required(*r.TransferUUID, "transfer_uuid"); err != nil {
			return err
		}
		referenceCount++
	}
	if referenceCount > 1 {
		return fmt.Errorf("only one of receipt_uuid, usage_uuid, adjustment_uuid, transfer_uuid may be set")
	}

	switch r.Reason {
	case storage.MovementReasonReceive:
		if r.ReceiptUUID == nil {
			return fmt.Errorf("receipt_uuid is required for receive movements")
		}
	case storage.MovementReasonUse:
		if r.UsageUUID == nil {
			return fmt.Errorf("usage_uuid is required for use movements")
		}
	case storage.MovementReasonTransfer:
		if r.TransferUUID == nil {
			return fmt.Errorf("transfer_uuid is required for transfer movements")
		}
	case storage.MovementReasonAdjust, storage.MovementReasonWaste:
		if r.AdjustmentUUID == nil {
			return fmt.Errorf("adjustment_uuid is required for adjust or waste movements")
		}
	}

	return nil
}

type InventoryMovementResponse struct {
	UUID              string     `json:"uuid"`
	IngredientLotUUID *string    `json:"ingredient_lot_uuid,omitempty"`
	BeerLotUUID       *string    `json:"beer_lot_uuid,omitempty"`
	StockLocationUUID string     `json:"stock_location_uuid"`
	Direction         string     `json:"direction"`
	Reason            string     `json:"reason"`
	Amount            int64      `json:"amount"`
	AmountUnit        string     `json:"amount_unit"`
	OccurredAt        time.Time  `json:"occurred_at"`
	ReceiptUUID       *string    `json:"receipt_uuid,omitempty"`
	UsageUUID         *string    `json:"usage_uuid,omitempty"`
	AdjustmentUUID    *string    `json:"adjustment_uuid,omitempty"`
	TransferUUID      *string    `json:"transfer_uuid,omitempty"`
	RemovalUUID       *string    `json:"removal_uuid,omitempty"`
	Notes             *string    `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewInventoryMovementResponse(movement storage.InventoryMovement) InventoryMovementResponse {
	return InventoryMovementResponse{
		UUID:              movement.UUID.String(),
		IngredientLotUUID: movement.IngredientLotUUID,
		BeerLotUUID:       movement.BeerLotUUID,
		StockLocationUUID: movement.StockLocationUUID,
		Direction:         movement.Direction,
		Reason:            movement.Reason,
		Amount:            movement.Amount,
		AmountUnit:        movement.AmountUnit,
		OccurredAt:        movement.OccurredAt,
		ReceiptUUID:       movement.ReceiptUUID,
		UsageUUID:         movement.UsageUUID,
		AdjustmentUUID:    movement.AdjustmentUUID,
		TransferUUID:      movement.TransferUUID,
		RemovalUUID:       movement.RemovalUUID,
		Notes:             movement.Notes,
		CreatedAt:         movement.CreatedAt,
		UpdatedAt:         movement.UpdatedAt,
		DeletedAt:         movement.DeletedAt,
	}
}

func NewInventoryMovementsResponse(movements []storage.InventoryMovement) []InventoryMovementResponse {
	resp := make([]InventoryMovementResponse, 0, len(movements))
	for _, movement := range movements {
		resp = append(resp, NewInventoryMovementResponse(movement))
	}
	return resp
}
