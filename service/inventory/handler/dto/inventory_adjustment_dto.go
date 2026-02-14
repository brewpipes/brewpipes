package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryAdjustmentRequest struct {
	Reason     string     `json:"reason"`
	AdjustedAt *time.Time `json:"adjusted_at"`
	Notes      *string    `json:"notes"`
}

func (r CreateInventoryAdjustmentRequest) Validate() error {
	if err := validate.Required(r.Reason, "reason"); err != nil {
		return err
	}
	if err := validateAdjustmentReason(r.Reason); err != nil {
		return err
	}

	return nil
}

type InventoryAdjustmentResponse struct {
	UUID       string     `json:"uuid"`
	Reason     string     `json:"reason"`
	AdjustedAt time.Time  `json:"adjusted_at"`
	Notes      *string    `json:"notes,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

func NewInventoryAdjustmentResponse(adjustment storage.InventoryAdjustment) InventoryAdjustmentResponse {
	return InventoryAdjustmentResponse{
		UUID:       adjustment.UUID.String(),
		Reason:     adjustment.Reason,
		AdjustedAt: adjustment.AdjustedAt,
		Notes:      adjustment.Notes,
		CreatedAt:  adjustment.CreatedAt,
		UpdatedAt:  adjustment.UpdatedAt,
		DeletedAt:  adjustment.DeletedAt,
	}
}

func NewInventoryAdjustmentsResponse(adjustments []storage.InventoryAdjustment) []InventoryAdjustmentResponse {
	resp := make([]InventoryAdjustmentResponse, 0, len(adjustments))
	for _, adjustment := range adjustments {
		resp = append(resp, NewInventoryAdjustmentResponse(adjustment))
	}
	return resp
}
