package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateAdditionRequest struct {
	BatchID          *int64     `json:"batch_id"`
	OccupancyID      *int64     `json:"occupancy_id"`
	AdditionType     string     `json:"addition_type"`
	Stage            *string    `json:"stage"`
	InventoryLotUUID *string    `json:"inventory_lot_uuid"`
	Amount           int64      `json:"amount"`
	AmountUnit       string     `json:"amount_unit"`
	AddedAt          *time.Time `json:"added_at"`
	Notes            *string    `json:"notes"`
}

func (r CreateAdditionRequest) Validate() error {
	if (r.BatchID == nil && r.OccupancyID == nil) || (r.BatchID != nil && r.OccupancyID != nil) {
		return fmt.Errorf("batch_id or occupancy_id is required")
	}
	if err := validateAdditionType(r.AdditionType); err != nil {
		return err
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if err := validateVolumeUnit(r.AmountUnit); err != nil {
		return err
	}
	if additionTypeRequiresInventory(r.AdditionType) && r.InventoryLotUUID == nil {
		return fmt.Errorf("inventory_lot_uuid is required for ingredient additions")
	}

	return nil
}

type AdditionResponse struct {
	ID               int64      `json:"id"`
	UUID             string     `json:"uuid"`
	BatchID          *int64     `json:"batch_id,omitempty"`
	OccupancyID      *int64     `json:"occupancy_id,omitempty"`
	AdditionType     string     `json:"addition_type"`
	Stage            *string    `json:"stage,omitempty"`
	InventoryLotUUID *string    `json:"inventory_lot_uuid,omitempty"`
	Amount           int64      `json:"amount"`
	AmountUnit       string     `json:"amount_unit"`
	AddedAt          time.Time  `json:"added_at"`
	Notes            *string    `json:"notes,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

func NewAdditionResponse(addition storage.Addition) AdditionResponse {
	return AdditionResponse{
		ID:               addition.ID,
		UUID:             addition.UUID.String(),
		BatchID:          addition.BatchID,
		OccupancyID:      addition.OccupancyID,
		AdditionType:     addition.AdditionType,
		Stage:            addition.Stage,
		InventoryLotUUID: uuidToStringPointer(addition.InventoryLotUUID),
		Amount:           addition.Amount,
		AmountUnit:       addition.AmountUnit,
		AddedAt:          addition.AddedAt,
		Notes:            addition.Notes,
		CreatedAt:        addition.CreatedAt,
		UpdatedAt:        addition.UpdatedAt,
		DeletedAt:        addition.DeletedAt,
	}
}

func NewAdditionsResponse(additions []storage.Addition) []AdditionResponse {
	resp := make([]AdditionResponse, 0, len(additions))
	for _, addition := range additions {
		resp = append(resp, NewAdditionResponse(addition))
	}
	return resp
}
