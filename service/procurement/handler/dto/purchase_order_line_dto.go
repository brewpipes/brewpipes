package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type CreatePurchaseOrderLineRequest struct {
	PurchaseOrderID   int64   `json:"purchase_order_id"`
	LineNumber        int     `json:"line_number"`
	ItemType          string  `json:"item_type"`
	ItemName          string  `json:"item_name"`
	InventoryItemUUID *string `json:"inventory_item_uuid"`
	Quantity          int64   `json:"quantity"`
	QuantityUnit      string  `json:"quantity_unit"`
	UnitCostCents     int64   `json:"unit_cost_cents"`
	Currency          string  `json:"currency"`
}

func (r CreatePurchaseOrderLineRequest) Validate() error {
	if r.PurchaseOrderID <= 0 {
		return fmt.Errorf("purchase_order_id is required")
	}
	if r.LineNumber <= 0 {
		return fmt.Errorf("line_number must be greater than zero")
	}
	if err := validateLineItemType(r.ItemType); err != nil {
		return err
	}
	if err := validateRequired(r.ItemName, "item_name"); err != nil {
		return err
	}
	if r.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}
	if err := validateRequired(r.QuantityUnit, "quantity_unit"); err != nil {
		return err
	}
	if r.UnitCostCents < 0 {
		return fmt.Errorf("unit_cost_cents must be zero or greater")
	}
	if err := validateCurrency(r.Currency); err != nil {
		return err
	}

	return nil
}

type PurchaseOrderLineResponse struct {
	ID                int64      `json:"id"`
	UUID              string     `json:"uuid"`
	PurchaseOrderID   int64      `json:"purchase_order_id"`
	LineNumber        int        `json:"line_number"`
	ItemType          string     `json:"item_type"`
	ItemName          string     `json:"item_name"`
	InventoryItemUUID *string    `json:"inventory_item_uuid,omitempty"`
	Quantity          int64      `json:"quantity"`
	QuantityUnit      string     `json:"quantity_unit"`
	UnitCostCents     int64      `json:"unit_cost_cents"`
	Currency          string     `json:"currency"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewPurchaseOrderLineResponse(line storage.PurchaseOrderLine) PurchaseOrderLineResponse {
	return PurchaseOrderLineResponse{
		ID:                line.ID,
		UUID:              line.UUID.String(),
		PurchaseOrderID:   line.PurchaseOrderID,
		LineNumber:        line.LineNumber,
		ItemType:          line.ItemType,
		ItemName:          line.ItemName,
		InventoryItemUUID: uuidToStringPointer(line.InventoryItemUUID),
		Quantity:          line.Quantity,
		QuantityUnit:      line.QuantityUnit,
		UnitCostCents:     line.UnitCostCents,
		Currency:          line.Currency,
		CreatedAt:         line.CreatedAt,
		UpdatedAt:         line.UpdatedAt,
		DeletedAt:         line.DeletedAt,
	}
}

func NewPurchaseOrderLinesResponse(lines []storage.PurchaseOrderLine) []PurchaseOrderLineResponse {
	resp := make([]PurchaseOrderLineResponse, 0, len(lines))
	for _, line := range lines {
		resp = append(resp, NewPurchaseOrderLineResponse(line))
	}
	return resp
}
