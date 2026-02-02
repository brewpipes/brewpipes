package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryMovementRequest struct {
	IngredientLotID *int64     `json:"ingredient_lot_id"`
	BeerLotID       *int64     `json:"beer_lot_id"`
	StockLocationID int64      `json:"stock_location_id"`
	Direction       string     `json:"direction"`
	Reason          string     `json:"reason"`
	Amount          int64      `json:"amount"`
	AmountUnit      string     `json:"amount_unit"`
	OccurredAt      *time.Time `json:"occurred_at"`
	ReceiptID       *int64     `json:"receipt_id"`
	UsageID         *int64     `json:"usage_id"`
	AdjustmentID    *int64     `json:"adjustment_id"`
	TransferID      *int64     `json:"transfer_id"`
	Notes           *string    `json:"notes"`
}

func (r CreateInventoryMovementRequest) Validate() error {
	if (r.IngredientLotID == nil && r.BeerLotID == nil) || (r.IngredientLotID != nil && r.BeerLotID != nil) {
		return fmt.Errorf("ingredient_lot_id or beer_lot_id is required")
	}
	if r.IngredientLotID != nil && *r.IngredientLotID <= 0 {
		return fmt.Errorf("ingredient_lot_id must be greater than zero")
	}
	if r.BeerLotID != nil && *r.BeerLotID <= 0 {
		return fmt.Errorf("beer_lot_id must be greater than zero")
	}
	if r.StockLocationID <= 0 {
		return fmt.Errorf("stock_location_id is required")
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
	if err := validateRequired(r.AmountUnit, "amount_unit"); err != nil {
		return err
	}

	referenceCount := 0
	if r.ReceiptID != nil {
		if *r.ReceiptID <= 0 {
			return fmt.Errorf("receipt_id must be greater than zero")
		}
		referenceCount++
	}
	if r.UsageID != nil {
		if *r.UsageID <= 0 {
			return fmt.Errorf("usage_id must be greater than zero")
		}
		referenceCount++
	}
	if r.AdjustmentID != nil {
		if *r.AdjustmentID <= 0 {
			return fmt.Errorf("adjustment_id must be greater than zero")
		}
		referenceCount++
	}
	if r.TransferID != nil {
		if *r.TransferID <= 0 {
			return fmt.Errorf("transfer_id must be greater than zero")
		}
		referenceCount++
	}
	if referenceCount > 1 {
		return fmt.Errorf("only one of receipt_id, usage_id, adjustment_id, transfer_id may be set")
	}

	switch r.Reason {
	case storage.MovementReasonReceive:
		if r.ReceiptID == nil {
			return fmt.Errorf("receipt_id is required for receive movements")
		}
	case storage.MovementReasonUse:
		if r.UsageID == nil {
			return fmt.Errorf("usage_id is required for use movements")
		}
	case storage.MovementReasonTransfer:
		if r.TransferID == nil {
			return fmt.Errorf("transfer_id is required for transfer movements")
		}
	case storage.MovementReasonAdjust, storage.MovementReasonWaste:
		if r.AdjustmentID == nil {
			return fmt.Errorf("adjustment_id is required for adjust or waste movements")
		}
	}

	return nil
}

type InventoryMovementResponse struct {
	ID              int64      `json:"id"`
	UUID            string     `json:"uuid"`
	IngredientLotID *int64     `json:"ingredient_lot_id,omitempty"`
	BeerLotID       *int64     `json:"beer_lot_id,omitempty"`
	StockLocationID int64      `json:"stock_location_id"`
	Direction       string     `json:"direction"`
	Reason          string     `json:"reason"`
	Amount          int64      `json:"amount"`
	AmountUnit      string     `json:"amount_unit"`
	OccurredAt      time.Time  `json:"occurred_at"`
	ReceiptID       *int64     `json:"receipt_id,omitempty"`
	UsageID         *int64     `json:"usage_id,omitempty"`
	AdjustmentID    *int64     `json:"adjustment_id,omitempty"`
	TransferID      *int64     `json:"transfer_id,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewInventoryMovementResponse(movement storage.InventoryMovement) InventoryMovementResponse {
	return InventoryMovementResponse{
		ID:              movement.ID,
		UUID:            movement.UUID.String(),
		IngredientLotID: movement.IngredientLotID,
		BeerLotID:       movement.BeerLotID,
		StockLocationID: movement.StockLocationID,
		Direction:       movement.Direction,
		Reason:          movement.Reason,
		Amount:          movement.Amount,
		AmountUnit:      movement.AmountUnit,
		OccurredAt:      movement.OccurredAt,
		ReceiptID:       movement.ReceiptID,
		UsageID:         movement.UsageID,
		AdjustmentID:    movement.AdjustmentID,
		TransferID:      movement.TransferID,
		Notes:           movement.Notes,
		CreatedAt:       movement.CreatedAt,
		UpdatedAt:       movement.UpdatedAt,
		DeletedAt:       movement.DeletedAt,
	}
}

func NewInventoryMovementsResponse(movements []storage.InventoryMovement) []InventoryMovementResponse {
	resp := make([]InventoryMovementResponse, 0, len(movements))
	for _, movement := range movements {
		resp = append(resp, NewInventoryMovementResponse(movement))
	}
	return resp
}
