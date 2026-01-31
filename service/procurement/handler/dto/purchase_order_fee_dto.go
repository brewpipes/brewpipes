package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type CreatePurchaseOrderFeeRequest struct {
	PurchaseOrderID int64  `json:"purchase_order_id"`
	FeeType         string `json:"fee_type"`
	AmountCents     int64  `json:"amount_cents"`
	Currency        string `json:"currency"`
}

func (r CreatePurchaseOrderFeeRequest) Validate() error {
	if r.PurchaseOrderID <= 0 {
		return fmt.Errorf("purchase_order_id is required")
	}
	if err := validateRequired(r.FeeType, "fee_type"); err != nil {
		return err
	}
	if r.AmountCents < 0 {
		return fmt.Errorf("amount_cents must be zero or greater")
	}
	if err := validateCurrency(r.Currency); err != nil {
		return err
	}

	return nil
}

type PurchaseOrderFeeResponse struct {
	ID              int64      `json:"id"`
	UUID            string     `json:"uuid"`
	PurchaseOrderID int64      `json:"purchase_order_id"`
	FeeType         string     `json:"fee_type"`
	AmountCents     int64      `json:"amount_cents"`
	Currency        string     `json:"currency"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewPurchaseOrderFeeResponse(fee storage.PurchaseOrderFee) PurchaseOrderFeeResponse {
	return PurchaseOrderFeeResponse{
		ID:              fee.ID,
		UUID:            fee.UUID.String(),
		PurchaseOrderID: fee.PurchaseOrderID,
		FeeType:         fee.FeeType,
		AmountCents:     fee.AmountCents,
		Currency:        fee.Currency,
		CreatedAt:       fee.CreatedAt,
		UpdatedAt:       fee.UpdatedAt,
		DeletedAt:       fee.DeletedAt,
	}
}

func NewPurchaseOrderFeesResponse(fees []storage.PurchaseOrderFee) []PurchaseOrderFeeResponse {
	resp := make([]PurchaseOrderFeeResponse, 0, len(fees))
	for _, fee := range fees {
		resp = append(resp, NewPurchaseOrderFeeResponse(fee))
	}
	return resp
}
