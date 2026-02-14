package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/uuidutil"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryUsageRequest struct {
	ProductionRefUUID *string    `json:"production_ref_uuid"`
	UsedAt            *time.Time `json:"used_at"`
	Notes             *string    `json:"notes"`
}

func (r CreateInventoryUsageRequest) Validate() error {
	return nil
}

type InventoryUsageResponse struct {
	UUID              string     `json:"uuid"`
	ProductionRefUUID *string    `json:"production_ref_uuid,omitempty"`
	UsedAt            time.Time  `json:"used_at"`
	Notes             *string    `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewInventoryUsageResponse(usage storage.InventoryUsage) InventoryUsageResponse {
	return InventoryUsageResponse{
		UUID:              usage.UUID.String(),
		ProductionRefUUID: uuidutil.ToStringPointer(usage.ProductionRefUUID),
		UsedAt:            usage.UsedAt,
		Notes:             usage.Notes,
		CreatedAt:         usage.CreatedAt,
		UpdatedAt:         usage.UpdatedAt,
		DeletedAt:         usage.DeletedAt,
	}
}

func NewInventoryUsageRecordsResponse(usageRecords []storage.InventoryUsage) []InventoryUsageResponse {
	resp := make([]InventoryUsageResponse, 0, len(usageRecords))
	for _, usage := range usageRecords {
		resp = append(resp, NewInventoryUsageResponse(usage))
	}
	return resp
}
