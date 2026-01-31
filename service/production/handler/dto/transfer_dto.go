package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateTransferRequest struct {
	SourceOccupancyID int64      `json:"source_occupancy_id"`
	DestVesselID      int64      `json:"dest_vessel_id"`
	VolumeID          int64      `json:"volume_id"`
	Amount            int64      `json:"amount"`
	AmountUnit        string     `json:"amount_unit"`
	LossAmount        *int64     `json:"loss_amount"`
	LossUnit          *string    `json:"loss_unit"`
	StartedAt         *time.Time `json:"started_at"`
	EndedAt           *time.Time `json:"ended_at"`
}

func (r CreateTransferRequest) Validate() error {
	if r.SourceOccupancyID <= 0 || r.DestVesselID <= 0 || r.VolumeID <= 0 {
		return fmt.Errorf("source_occupancy_id, dest_vessel_id, and volume_id are required")
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if err := validateVolumeUnit(r.AmountUnit); err != nil {
		return err
	}
	if (r.LossAmount == nil) != (r.LossUnit == nil) {
		return fmt.Errorf("loss_amount and loss_unit must be provided together")
	}
	if r.LossAmount != nil && *r.LossAmount <= 0 {
		return fmt.Errorf("loss_amount must be greater than zero")
	}
	if r.LossUnit != nil {
		if err := validateVolumeUnit(*r.LossUnit); err != nil {
			return err
		}
	}
	if r.StartedAt != nil && r.EndedAt != nil {
		if r.EndedAt.Before(*r.StartedAt) {
			return fmt.Errorf("ended_at must be after started_at")
		}
	}

	return nil
}

type TransferResponse struct {
	ID                int64      `json:"id"`
	UUID              string     `json:"uuid"`
	SourceOccupancyID int64      `json:"source_occupancy_id"`
	DestOccupancyID   int64      `json:"dest_occupancy_id"`
	Amount            int64      `json:"amount"`
	AmountUnit        string     `json:"amount_unit"`
	LossAmount        *int64     `json:"loss_amount,omitempty"`
	LossUnit          *string    `json:"loss_unit,omitempty"`
	StartedAt         time.Time  `json:"started_at"`
	EndedAt           *time.Time `json:"ended_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type TransferRecordResponse struct {
	Transfer      TransferResponse  `json:"transfer"`
	DestOccupancy OccupancyResponse `json:"dest_occupancy"`
}

func NewTransferResponse(transfer storage.Transfer) TransferResponse {
	return TransferResponse{
		ID:                transfer.ID,
		UUID:              transfer.UUID.String(),
		SourceOccupancyID: transfer.SourceOccupancyID,
		DestOccupancyID:   transfer.DestOccupancyID,
		Amount:            transfer.Amount,
		AmountUnit:        transfer.AmountUnit,
		LossAmount:        transfer.LossAmount,
		LossUnit:          transfer.LossUnit,
		StartedAt:         transfer.StartedAt,
		EndedAt:           transfer.EndedAt,
		CreatedAt:         transfer.CreatedAt,
		UpdatedAt:         transfer.UpdatedAt,
		DeletedAt:         transfer.DeletedAt,
	}
}

func NewTransfersResponse(transfers []storage.Transfer) []TransferResponse {
	resp := make([]TransferResponse, 0, len(transfers))
	for _, transfer := range transfers {
		resp = append(resp, NewTransferResponse(transfer))
	}
	return resp
}

func NewTransferRecordResponse(transfer storage.Transfer, occupancy storage.Occupancy) TransferRecordResponse {
	return TransferRecordResponse{
		Transfer:      NewTransferResponse(transfer),
		DestOccupancy: NewOccupancyResponse(occupancy),
	}
}
