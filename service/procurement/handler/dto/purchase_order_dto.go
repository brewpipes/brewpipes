package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type CreatePurchaseOrderRequest struct {
	SupplierUUID string     `json:"supplier_uuid"`
	OrderNumber  string     `json:"order_number"`
	Status       string     `json:"status"`
	OrderedAt    *time.Time `json:"ordered_at"`
	ExpectedAt   *time.Time `json:"expected_at"`
	Notes        *string    `json:"notes"`
}

func (r CreatePurchaseOrderRequest) Validate() error {
	if strings.TrimSpace(r.SupplierUUID) == "" {
		return fmt.Errorf("supplier_uuid is required")
	}
	status := strings.TrimSpace(r.Status)
	if status != "" {
		if err := validatePurchaseOrderStatus(status); err != nil {
			return err
		}
	}

	return nil
}

type UpdatePurchaseOrderRequest struct {
	OrderNumber *string    `json:"order_number"`
	Status      *string    `json:"status"`
	OrderedAt   *time.Time `json:"ordered_at"`
	ExpectedAt  *time.Time `json:"expected_at"`
	Notes       *string    `json:"notes"`
}

func (r UpdatePurchaseOrderRequest) Validate() error {
	if r.OrderNumber == nil && r.Status == nil && r.OrderedAt == nil && r.ExpectedAt == nil && r.Notes == nil {
		return fmt.Errorf("at least one field must be provided")
	}
	if r.OrderNumber != nil {
		if err := validate.Required(*r.OrderNumber, "order_number"); err != nil {
			return err
		}
	}
	if r.Status != nil {
		status := strings.TrimSpace(*r.Status)
		if status == "" {
			return fmt.Errorf("status is required")
		}
		if err := validatePurchaseOrderStatus(status); err != nil {
			return err
		}
	}

	return nil
}

type PurchaseOrderResponse struct {
	UUID         string     `json:"uuid"`
	SupplierUUID string     `json:"supplier_uuid"`
	OrderNumber  string     `json:"order_number"`
	Status       string     `json:"status"`
	OrderedAt    *time.Time `json:"ordered_at,omitempty"`
	ExpectedAt   *time.Time `json:"expected_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewPurchaseOrderResponse(order storage.PurchaseOrder) PurchaseOrderResponse {
	var supplierUUID string
	if order.SupplierUUID != nil {
		supplierUUID = *order.SupplierUUID
	}
	return PurchaseOrderResponse{
		UUID:         order.UUID.String(),
		SupplierUUID: supplierUUID,
		OrderNumber:  order.OrderNumber,
		Status:       order.Status,
		OrderedAt:    order.OrderedAt,
		ExpectedAt:   order.ExpectedAt,
		Notes:        order.Notes,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		DeletedAt:    order.DeletedAt,
	}
}

func NewPurchaseOrdersResponse(orders []storage.PurchaseOrder) []PurchaseOrderResponse {
	resp := make([]PurchaseOrderResponse, 0, len(orders))
	for _, order := range orders {
		resp = append(resp, NewPurchaseOrderResponse(order))
	}
	return resp
}
