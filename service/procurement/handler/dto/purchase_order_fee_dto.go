package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type CreatePurchaseOrderFeeRequest struct {
	PurchaseOrderUUID string `json:"purchase_order_uuid"`
	FeeType           string `json:"fee_type"`
	AmountCents       int64  `json:"amount_cents"`
	Currency          string `json:"currency"`
}

func (r CreatePurchaseOrderFeeRequest) Validate() error {
	if err := validateRequired(r.PurchaseOrderUUID, "purchase_order_uuid"); err != nil {
		return err
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

type UpdatePurchaseOrderFeeRequest struct {
	FeeType     *string `json:"fee_type"`
	AmountCents *int64  `json:"amount_cents"`
	Currency    *string `json:"currency"`
}

func (r UpdatePurchaseOrderFeeRequest) Validate() error {
	if r.FeeType == nil && r.AmountCents == nil && r.Currency == nil {
		return fmt.Errorf("at least one field must be provided")
	}
	if r.FeeType != nil {
		if err := validateRequired(*r.FeeType, "fee_type"); err != nil {
			return err
		}
	}
	if r.AmountCents != nil && *r.AmountCents < 0 {
		return fmt.Errorf("amount_cents must be zero or greater")
	}
	if r.Currency != nil {
		if err := validateCurrency(*r.Currency); err != nil {
			return err
		}
	}

	return nil
}

type PurchaseOrderFeeResponse struct {
	UUID              string     `json:"uuid"`
	PurchaseOrderUUID string     `json:"purchase_order_uuid"`
	FeeType           string     `json:"fee_type"`
	AmountCents       int64      `json:"amount_cents"`
	Currency          string     `json:"currency"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewPurchaseOrderFeeResponse(fee storage.PurchaseOrderFee) PurchaseOrderFeeResponse {
	var purchaseOrderUUID string
	if fee.PurchaseOrderUUID != nil {
		purchaseOrderUUID = *fee.PurchaseOrderUUID
	}
	return PurchaseOrderFeeResponse{
		UUID:              fee.UUID.String(),
		PurchaseOrderUUID: purchaseOrderUUID,
		FeeType:           fee.FeeType,
		AmountCents:       fee.AmountCents,
		Currency:          fee.Currency,
		CreatedAt:         fee.CreatedAt,
		UpdatedAt:         fee.UpdatedAt,
		DeletedAt:         fee.DeletedAt,
	}
}

func NewPurchaseOrderFeesResponse(fees []storage.PurchaseOrderFee) []PurchaseOrderFeeResponse {
	resp := make([]PurchaseOrderFeeResponse, 0, len(fees))
	for _, fee := range fees {
		resp = append(resp, NewPurchaseOrderFeeResponse(fee))
	}
	return resp
}
