package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/uuidutil"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateAdditionRequest struct {
	BatchUUID        *string    `json:"batch_uuid"`
	OccupancyUUID    *string    `json:"occupancy_uuid"`
	VolumeUUID       *string    `json:"volume_uuid"`
	AdditionType     string     `json:"addition_type"`
	Stage            *string    `json:"stage"`
	InventoryLotUUID *string    `json:"inventory_lot_uuid"`
	Amount           int64      `json:"amount"`
	AmountUnit       string     `json:"amount_unit"`
	AddedAt          *time.Time `json:"added_at"`
	Notes            *string    `json:"notes"`
}

func (r CreateAdditionRequest) Validate() error {
	targetCount := 0
	if r.BatchUUID != nil {
		targetCount++
	}
	if r.OccupancyUUID != nil {
		targetCount++
	}
	if r.VolumeUUID != nil {
		targetCount++
	}
	if targetCount != 1 {
		return fmt.Errorf("exactly one of batch_uuid, occupancy_uuid, or volume_uuid is required")
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
	UUID             string     `json:"uuid"`
	BatchUUID        *string    `json:"batch_uuid,omitempty"`
	OccupancyUUID    *string    `json:"occupancy_uuid,omitempty"`
	VolumeUUID       *string    `json:"volume_uuid,omitempty"`
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
		UUID:             addition.UUID.String(),
		BatchUUID:        addition.BatchUUID,
		OccupancyUUID:    addition.OccupancyUUID,
		VolumeUUID:       addition.VolumeUUID,
		AdditionType:     addition.AdditionType,
		Stage:            addition.Stage,
		InventoryLotUUID: uuidutil.ToStringPointer(addition.InventoryLotUUID),
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
