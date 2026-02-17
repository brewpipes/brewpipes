package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateStockLocationRequest struct {
	Name         string  `json:"name"`
	LocationType *string `json:"location_type"`
	Description  *string `json:"description"`
}

func (r CreateStockLocationRequest) Validate() error {
	if err := validate.Required(r.Name, "name"); err != nil {
		return err
	}
	if r.LocationType != nil {
		if err := validateStockLocationType(*r.LocationType); err != nil {
			return err
		}
	}

	return nil
}

// UpdateStockLocationRequest is the request payload for PATCH /stock-locations/{uuid}.
type UpdateStockLocationRequest struct {
	Name         *string `json:"name"`
	LocationType *string `json:"location_type"`
	Description  *string `json:"description"`
}

func (r UpdateStockLocationRequest) Validate() error {
	if r.Name != nil {
		if err := validate.Required(*r.Name, "name"); err != nil {
			return err
		}
	}
	if r.LocationType != nil {
		if err := validateStockLocationType(*r.LocationType); err != nil {
			return err
		}
	}
	return nil
}

type StockLocationResponse struct {
	UUID         string     `json:"uuid"`
	Name         string     `json:"name"`
	LocationType *string    `json:"location_type,omitempty"`
	Description  *string    `json:"description,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewStockLocationResponse(location storage.StockLocation) StockLocationResponse {
	return StockLocationResponse{
		UUID:         location.UUID.String(),
		Name:         location.Name,
		LocationType: location.LocationType,
		Description:  location.Description,
		CreatedAt:    location.CreatedAt,
		UpdatedAt:    location.UpdatedAt,
		DeletedAt:    location.DeletedAt,
	}
}

func NewStockLocationsResponse(locations []storage.StockLocation) []StockLocationResponse {
	resp := make([]StockLocationResponse, 0, len(locations))
	for _, location := range locations {
		resp = append(resp, NewStockLocationResponse(location))
	}
	return resp
}
