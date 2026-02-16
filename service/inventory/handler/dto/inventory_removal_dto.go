package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

// CreateRemovalRequest is the request body for POST /removals.
type CreateRemovalRequest struct {
	Category          string     `json:"category"`
	Reason            string     `json:"reason"`
	Amount            int64      `json:"amount"`
	AmountUnit        string     `json:"amount_unit"`
	RemovedAt         *time.Time `json:"removed_at"`
	BatchUUID         *string    `json:"batch_uuid"`
	BeerLotUUID       *string    `json:"beer_lot_uuid"`
	OccupancyUUID     *string    `json:"occupancy_uuid"`
	StockLocationUUID *string    `json:"stock_location_uuid"`
	ReferenceCode     *string    `json:"reference_code"`
	PerformedBy       *string    `json:"performed_by"`
	Destination       *string    `json:"destination"`
	Notes             *string    `json:"notes"`
}

// Validate validates the create removal request.
func (r CreateRemovalRequest) Validate() error {
	if err := validate.Required(r.Category, "category"); err != nil {
		return err
	}
	if err := validateRemovalCategory(r.Category); err != nil {
		return err
	}
	if err := validate.Required(r.Reason, "reason"); err != nil {
		return err
	}
	if err := validateRemovalReason(r.Reason); err != nil {
		return err
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if err := validate.Required(r.AmountUnit, "amount_unit"); err != nil {
		return err
	}
	if r.BatchUUID == nil && r.BeerLotUUID == nil && r.OccupancyUUID == nil {
		return fmt.Errorf("at least one of batch_uuid, beer_lot_uuid, or occupancy_uuid is required")
	}
	if r.BeerLotUUID != nil {
		if err := validate.Required(*r.BeerLotUUID, "beer_lot_uuid"); err != nil {
			return err
		}
		if r.StockLocationUUID == nil {
			return fmt.Errorf("stock_location_uuid is required when beer_lot_uuid is provided")
		}
		if err := validate.Required(*r.StockLocationUUID, "stock_location_uuid"); err != nil {
			return err
		}
	}
	if r.BatchUUID != nil {
		if err := validate.Required(*r.BatchUUID, "batch_uuid"); err != nil {
			return err
		}
	}
	if r.OccupancyUUID != nil {
		if err := validate.Required(*r.OccupancyUUID, "occupancy_uuid"); err != nil {
			return err
		}
	}
	return nil
}

// UpdateRemovalRequest is the request body for PATCH /removals/{uuid}.
type UpdateRemovalRequest struct {
	Category      *string    `json:"category"`
	Reason        *string    `json:"reason"`
	Amount        *int64     `json:"amount"`
	AmountUnit    *string    `json:"amount_unit"`
	RemovedAt     *time.Time `json:"removed_at"`
	ReferenceCode *string    `json:"reference_code"`
	PerformedBy   *string    `json:"performed_by"`
	Destination   *string    `json:"destination"`
	Notes         *string    `json:"notes"`
}

// Validate validates the update removal request.
func (r UpdateRemovalRequest) Validate() error {
	if r.Category != nil {
		if err := validateRemovalCategory(*r.Category); err != nil {
			return err
		}
	}
	if r.Reason != nil {
		if err := validateRemovalReason(*r.Reason); err != nil {
			return err
		}
	}
	if r.Amount != nil && *r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if r.AmountUnit != nil {
		if err := validate.Required(*r.AmountUnit, "amount_unit"); err != nil {
			return err
		}
	}
	return nil
}

// RemovalResponse is the response body for a single removal.
type RemovalResponse struct {
	UUID              string     `json:"uuid"`
	Category          string     `json:"category"`
	Reason            string     `json:"reason"`
	BatchUUID         *string    `json:"batch_uuid,omitempty"`
	BeerLotUUID       *string    `json:"beer_lot_uuid,omitempty"`
	OccupancyUUID     *string    `json:"occupancy_uuid,omitempty"`
	Amount            int64      `json:"amount"`
	AmountUnit        string     `json:"amount_unit"`
	AmountBBL         *float64   `json:"amount_bbl,omitempty"`
	IsTaxable         bool       `json:"is_taxable"`
	ReferenceCode     *string    `json:"reference_code,omitempty"`
	PerformedBy       *string    `json:"performed_by,omitempty"`
	RemovedAt         time.Time  `json:"removed_at"`
	Destination       *string    `json:"destination,omitempty"`
	Notes             *string    `json:"notes,omitempty"`
	MovementUUID      *string    `json:"movement_uuid,omitempty"`
	StockLocationUUID *string    `json:"stock_location_uuid,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

// NewRemovalResponse creates a RemovalResponse from a storage model.
func NewRemovalResponse(removal storage.InventoryRemoval) RemovalResponse {
	return RemovalResponse{
		UUID:              removal.UUID.String(),
		Category:          removal.Category,
		Reason:            removal.Reason,
		BatchUUID:         removal.BatchUUID,
		BeerLotUUID:       removal.BeerLotUUID,
		OccupancyUUID:     removal.OccupancyUUID,
		Amount:            removal.Amount,
		AmountUnit:        removal.AmountUnit,
		AmountBBL:         removal.AmountBBL,
		IsTaxable:         removal.IsTaxable,
		ReferenceCode:     removal.ReferenceCode,
		PerformedBy:       removal.PerformedBy,
		RemovedAt:         removal.RemovedAt,
		Destination:       removal.Destination,
		Notes:             removal.Notes,
		MovementUUID:      removal.MovementUUID,
		StockLocationUUID: removal.StockLocationUUID,
		CreatedAt:         removal.CreatedAt,
		UpdatedAt:         removal.UpdatedAt,
		DeletedAt:         removal.DeletedAt,
	}
}

// NewRemovalsResponse creates a slice of RemovalResponse from storage models.
func NewRemovalsResponse(removals []storage.InventoryRemoval) []RemovalResponse {
	resp := make([]RemovalResponse, 0, len(removals))
	for _, removal := range removals {
		resp = append(resp, NewRemovalResponse(removal))
	}
	return resp
}

// RemovalSummaryResponse is the response body for GET /removal-summary.
type RemovalSummaryResponse struct {
	TotalBBL   float64                     `json:"total_bbl"`
	TaxableBBL float64                     `json:"taxable_bbl"`
	TaxFreeBBL float64                     `json:"tax_free_bbl"`
	TotalCount int                         `json:"total_count"`
	ByCategory []RemovalCategorySummaryDTO `json:"by_category"`
}

// RemovalCategorySummaryDTO is a per-category summary entry.
type RemovalCategorySummaryDTO struct {
	Category string  `json:"category"`
	TotalBBL float64 `json:"total_bbl"`
	Count    int     `json:"count"`
}

// NewRemovalSummaryResponse creates a RemovalSummaryResponse from a storage model.
func NewRemovalSummaryResponse(summary storage.RemovalSummary) RemovalSummaryResponse {
	cats := make([]RemovalCategorySummaryDTO, 0, len(summary.ByCategory))
	for _, c := range summary.ByCategory {
		cats = append(cats, RemovalCategorySummaryDTO{
			Category: c.Category,
			TotalBBL: c.TotalBBL,
			Count:    c.Count,
		})
	}
	return RemovalSummaryResponse{
		TotalBBL:   summary.TotalBBL,
		TaxableBBL: summary.TaxableBBL,
		TaxFreeBBL: summary.TaxFreeBBL,
		TotalCount: summary.TotalCount,
		ByCategory: cats,
	}
}
