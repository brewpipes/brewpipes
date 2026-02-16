package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateInventoryTransferRequest struct {
	IngredientLotUUID  *string    `json:"ingredient_lot_uuid"`
	BeerLotUUID        *string    `json:"beer_lot_uuid"`
	SourceLocationUUID string     `json:"source_location_uuid"`
	DestLocationUUID   string     `json:"dest_location_uuid"`
	Amount             int64      `json:"amount"`
	AmountUnit         string     `json:"amount_unit"`
	TransferredAt      *time.Time `json:"transferred_at"`
	Notes              *string    `json:"notes"`
}

func (r CreateInventoryTransferRequest) Validate() error {
	if (r.IngredientLotUUID == nil && r.BeerLotUUID == nil) || (r.IngredientLotUUID != nil && r.BeerLotUUID != nil) {
		return fmt.Errorf("exactly one of ingredient_lot_uuid or beer_lot_uuid is required")
	}
	if r.IngredientLotUUID != nil {
		if err := validate.Required(*r.IngredientLotUUID, "ingredient_lot_uuid"); err != nil {
			return err
		}
	}
	if r.BeerLotUUID != nil {
		if err := validate.Required(*r.BeerLotUUID, "beer_lot_uuid"); err != nil {
			return err
		}
	}
	if err := validate.Required(r.SourceLocationUUID, "source_location_uuid"); err != nil {
		return err
	}
	if err := validate.Required(r.DestLocationUUID, "dest_location_uuid"); err != nil {
		return err
	}
	if r.SourceLocationUUID == r.DestLocationUUID {
		return fmt.Errorf("source_location_uuid and dest_location_uuid must differ")
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if err := validate.Required(r.AmountUnit, "amount_unit"); err != nil {
		return err
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
