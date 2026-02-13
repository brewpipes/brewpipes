package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type CreatePurchaseOrderLineRequest struct {
	PurchaseOrderUUID string  `json:"purchase_order_uuid"`
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
	if err := validateRequired(r.PurchaseOrderUUID, "purchase_order_uuid"); err != nil {
		return err
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

type UpdatePurchaseOrderLineRequest struct {
	LineNumber        *int    `json:"line_number"`
	ItemType          *string `json:"item_type"`
	ItemName          *string `json:"item_name"`
	InventoryItemUUID *string `json:"inventory_item_uuid"`
	Quantity          *int64  `json:"quantity"`
	QuantityUnit      *string `json:"quantity_unit"`
	UnitCostCents     *int64  `json:"unit_cost_cents"`
	Currency          *string `json:"currency"`
}

func (r UpdatePurchaseOrderLineRequest) Validate() error {
	if r.LineNumber == nil && r.ItemType == nil && r.ItemName == nil && r.InventoryItemUUID == nil && r.Quantity == nil && r.QuantityUnit == nil && r.UnitCostCents == nil && r.Currency == nil {
		return fmt.Errorf("at least one field must be provided")
	}
	if r.LineNumber != nil && *r.LineNumber <= 0 {
		return fmt.Errorf("line_number must be greater than zero")
	}
	if r.ItemType != nil {
		if err := validateLineItemType(*r.ItemType); err != nil {
			return err
		}
	}
	if r.ItemName != nil {
		if err := validateRequired(*r.ItemName, "item_name"); err != nil {
			return err
		}
	}
	if r.Quantity != nil && *r.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}
	if r.QuantityUnit != nil {
		if err := validateRequired(*r.QuantityUnit, "quantity_unit"); err != nil {
			return err
		}
	}
	if r.UnitCostCents != nil && *r.UnitCostCents < 0 {
		return fmt.Errorf("unit_cost_cents must be zero or greater")
	}
	if r.Currency != nil {
		if err := validateCurrency(*r.Currency); err != nil {
			return err
		}
	}

	return nil
}

type PurchaseOrderLineResponse struct {
	UUID              string     `json:"uuid"`
	PurchaseOrderUUID string     `json:"purchase_order_uuid"`
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
	var purchaseOrderUUID string
	if line.PurchaseOrderUUID != nil {
		purchaseOrderUUID = *line.PurchaseOrderUUID
	}
	return PurchaseOrderLineResponse{
		UUID:              line.UUID.String(),
		PurchaseOrderUUID: purchaseOrderUUID,
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
