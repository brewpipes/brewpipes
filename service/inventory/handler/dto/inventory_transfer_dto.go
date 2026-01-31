package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryTransferRequest struct {
	SourceLocationID int64      `json:"source_location_id"`
	DestLocationID   int64      `json:"dest_location_id"`
	TransferredAt    *time.Time `json:"transferred_at"`
	Notes            *string    `json:"notes"`
}

func (r CreateInventoryTransferRequest) Validate() error {
	if r.SourceLocationID <= 0 || r.DestLocationID <= 0 {
		return fmt.Errorf("source_location_id and dest_location_id are required")
	}
	if r.SourceLocationID == r.DestLocationID {
		return fmt.Errorf("source_location_id and dest_location_id must differ")
	}

	return nil
}

type InventoryTransferResponse struct {
	ID               int64      `json:"id"`
	UUID             string     `json:"uuid"`
	SourceLocationID int64      `json:"source_location_id"`
	DestLocationID   int64      `json:"dest_location_id"`
	TransferredAt    time.Time  `json:"transferred_at"`
	Notes            *string    `json:"notes,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

func NewInventoryTransferResponse(transfer storage.InventoryTransfer) InventoryTransferResponse {
	return InventoryTransferResponse{
		ID:               transfer.ID,
		UUID:             transfer.UUID.String(),
		SourceLocationID: transfer.SourceLocationID,
		DestLocationID:   transfer.DestLocationID,
		TransferredAt:    transfer.TransferredAt,
		Notes:            transfer.Notes,
		CreatedAt:        transfer.CreatedAt,
		UpdatedAt:        transfer.UpdatedAt,
		DeletedAt:        transfer.DeletedAt,
	}
}

func NewInventoryTransfersResponse(transfers []storage.InventoryTransfer) []InventoryTransferResponse {
	resp := make([]InventoryTransferResponse, 0, len(transfers))
	for _, transfer := range transfers {
		resp = append(resp, NewInventoryTransferResponse(transfer))
	}
	return resp
}
