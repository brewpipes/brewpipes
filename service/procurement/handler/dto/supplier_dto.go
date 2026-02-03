package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type CreateSupplierRequest struct {
	Name         string  `json:"name"`
	ContactName  *string `json:"contact_name"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	City         *string `json:"city"`
	Region       *string `json:"region"`
	PostalCode   *string `json:"postal_code"`
	Country      *string `json:"country"`
}

func (r CreateSupplierRequest) Validate() error {
	return validateRequired(r.Name, "name")
}

type UpdateSupplierRequest struct {
	Name         *string `json:"name"`
	ContactName  *string `json:"contact_name"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	City         *string `json:"city"`
	Region       *string `json:"region"`
	PostalCode   *string `json:"postal_code"`
	Country      *string `json:"country"`
}

func (r UpdateSupplierRequest) Validate() error {
	if r.Name == nil && r.ContactName == nil && r.Email == nil && r.Phone == nil &&
		r.AddressLine1 == nil && r.AddressLine2 == nil && r.City == nil &&
		r.Region == nil && r.PostalCode == nil && r.Country == nil {
		return fmt.Errorf("at least one field must be provided")
	}
	if r.Name != nil {
		if strings.TrimSpace(*r.Name) == "" {
			return fmt.Errorf("name cannot be empty")
		}
	}

	return nil
}

type SupplierResponse struct {
	ID           int64      `json:"id"`
	UUID         string     `json:"uuid"`
	Name         string     `json:"name"`
	ContactName  *string    `json:"contact_name,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	AddressLine1 *string    `json:"address_line1,omitempty"`
	AddressLine2 *string    `json:"address_line2,omitempty"`
	City         *string    `json:"city,omitempty"`
	Region       *string    `json:"region,omitempty"`
	PostalCode   *string    `json:"postal_code,omitempty"`
	Country      *string    `json:"country,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewSupplierResponse(supplier storage.Supplier) SupplierResponse {
	return SupplierResponse{
		ID:           supplier.ID,
		UUID:         supplier.UUID.String(),
		Name:         supplier.Name,
		ContactName:  supplier.ContactName,
		Email:        supplier.Email,
		Phone:        supplier.Phone,
		AddressLine1: supplier.AddressLine1,
		AddressLine2: supplier.AddressLine2,
		City:         supplier.City,
		Region:       supplier.Region,
		PostalCode:   supplier.PostalCode,
		Country:      supplier.Country,
		CreatedAt:    supplier.CreatedAt,
		UpdatedAt:    supplier.UpdatedAt,
		DeletedAt:    supplier.DeletedAt,
	}
}

func NewSuppliersResponse(suppliers []storage.Supplier) []SupplierResponse {
	resp := make([]SupplierResponse, 0, len(suppliers))
	for _, supplier := range suppliers {
		resp = append(resp, NewSupplierResponse(supplier))
	}
	return resp
}
