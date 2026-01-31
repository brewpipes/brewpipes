package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryReceiptRequest struct {
	SupplierUUID  *string    `json:"supplier_uuid"`
	ReferenceCode *string    `json:"reference_code"`
	ReceivedAt    *time.Time `json:"received_at"`
	Notes         *string    `json:"notes"`
}

func (r CreateInventoryReceiptRequest) Validate() error {
	return nil
}

type InventoryReceiptResponse struct {
	ID            int64      `json:"id"`
	UUID          string     `json:"uuid"`
	SupplierUUID  *string    `json:"supplier_uuid,omitempty"`
	ReferenceCode *string    `json:"reference_code,omitempty"`
	ReceivedAt    time.Time  `json:"received_at"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

func NewInventoryReceiptResponse(receipt storage.InventoryReceipt) InventoryReceiptResponse {
	return InventoryReceiptResponse{
		ID:            receipt.ID,
		UUID:          receipt.UUID.String(),
		SupplierUUID:  uuidToStringPointer(receipt.SupplierUUID),
		ReferenceCode: receipt.ReferenceCode,
		ReceivedAt:    receipt.ReceivedAt,
		Notes:         receipt.Notes,
		CreatedAt:     receipt.CreatedAt,
		UpdatedAt:     receipt.UpdatedAt,
		DeletedAt:     receipt.DeletedAt,
	}
}

func NewInventoryReceiptsResponse(receipts []storage.InventoryReceipt) []InventoryReceiptResponse {
	resp := make([]InventoryReceiptResponse, 0, len(receipts))
	for _, receipt := range receipts {
		resp = append(resp, NewInventoryReceiptResponse(receipt))
	}
	return resp
}
