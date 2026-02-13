package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryTransferRequest struct {
	SourceLocationUUID string     `json:"source_location_uuid"`
	DestLocationUUID   string     `json:"dest_location_uuid"`
	TransferredAt      *time.Time `json:"transferred_at"`
	Notes              *string    `json:"notes"`
}

func (r CreateInventoryTransferRequest) Validate() error {
	if err := validateRequired(r.SourceLocationUUID, "source_location_uuid"); err != nil {
		return err
	}
	if err := validateRequired(r.DestLocationUUID, "dest_location_uuid"); err != nil {
		return err
	}
	if r.SourceLocationUUID == r.DestLocationUUID {
		return fmt.Errorf("source_location_uuid and dest_location_uuid must differ")
	}

	return nil
}

type InventoryTransferResponse struct {
	UUID               string     `json:"uuid"`
	SourceLocationUUID string     `json:"source_location_uuid"`
	DestLocationUUID   string     `json:"dest_location_uuid"`
	TransferredAt      time.Time  `json:"transferred_at"`
	Notes              *string    `json:"notes,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

func NewInventoryTransferResponse(transfer storage.InventoryTransfer) InventoryTransferResponse {
	return InventoryTransferResponse{
		UUID:               transfer.UUID.String(),
		SourceLocationUUID: transfer.SourceLocationUUID,
		DestLocationUUID:   transfer.DestLocationUUID,
		TransferredAt:      transfer.TransferredAt,
		Notes:              transfer.Notes,
		CreatedAt:          transfer.CreatedAt,
		UpdatedAt:          transfer.UpdatedAt,
		DeletedAt:          transfer.DeletedAt,
	}
}

func NewInventoryTransfersResponse(transfers []storage.InventoryTransfer) []InventoryTransferResponse {
	resp := make([]InventoryTransferResponse, 0, len(transfers))
	for _, transfer := range transfers {
		resp = append(resp, NewInventoryTransferResponse(transfer))
	}
	return resp
}
