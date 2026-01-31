package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type CreatePurchaseOrderRequest struct {
	SupplierID  int64      `json:"supplier_id"`
	OrderNumber string     `json:"order_number"`
	Status      string     `json:"status"`
	OrderedAt   *time.Time `json:"ordered_at"`
	ExpectedAt  *time.Time `json:"expected_at"`
	Notes       *string    `json:"notes"`
}

func (r CreatePurchaseOrderRequest) Validate() error {
	if r.SupplierID <= 0 {
		return fmt.Errorf("supplier_id is required")
	}
	if err := validateRequired(r.OrderNumber, "order_number"); err != nil {
		return err
	}
	status := strings.TrimSpace(r.Status)
	if status != "" {
		if err := validatePurchaseOrderStatus(status); err != nil {
			return err
		}
	}

	return nil
}

type PurchaseOrderResponse struct {
	ID          int64      `json:"id"`
	UUID        string     `json:"uuid"`
	SupplierID  int64      `json:"supplier_id"`
	OrderNumber string     `json:"order_number"`
	Status      string     `json:"status"`
	OrderedAt   *time.Time `json:"ordered_at,omitempty"`
	ExpectedAt  *time.Time `json:"expected_at,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewPurchaseOrderResponse(order storage.PurchaseOrder) PurchaseOrderResponse {
	return PurchaseOrderResponse{
		ID:          order.ID,
		UUID:        order.UUID.String(),
		SupplierID:  order.SupplierID,
		OrderNumber: order.OrderNumber,
		Status:      order.Status,
		OrderedAt:   order.OrderedAt,
		ExpectedAt:  order.ExpectedAt,
		Notes:       order.Notes,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		DeletedAt:   order.DeletedAt,
	}
}

func NewPurchaseOrdersResponse(orders []storage.PurchaseOrder) []PurchaseOrderResponse {
	resp := make([]PurchaseOrderResponse, 0, len(orders))
	for _, order := range orders {
		resp = append(resp, NewPurchaseOrderResponse(order))
	}
	return resp
}
