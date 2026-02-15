package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateTransferRequest struct {
	SourceOccupancyUUID string     `json:"source_occupancy_uuid"`
	DestVesselUUID      string     `json:"dest_vessel_uuid"`
	VolumeUUID          string     `json:"volume_uuid"`
	Amount              int64      `json:"amount"`
	AmountUnit          string     `json:"amount_unit"`
	LossAmount          *int64     `json:"loss_amount"`
	LossUnit            *string    `json:"loss_unit"`
	StartedAt           *time.Time `json:"started_at"`
	EndedAt             *time.Time `json:"ended_at"`
	CloseSource         *bool      `json:"close_source,omitempty"` // default true if nil
	DestStatus          *string    `json:"dest_status,omitempty"`  // status for destination occupancy
}

func (r CreateTransferRequest) Validate() error {
	if err := validate.Required(r.SourceOccupancyUUID, "source_occupancy_uuid"); err != nil {
		return err
	}
	if err := validate.Required(r.DestVesselUUID, "dest_vessel_uuid"); err != nil {
		return err
	}
	if err := validate.Required(r.VolumeUUID, "volume_uuid"); err != nil {
		return err
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
	if r.DestStatus != nil {
		if err := validateOccupancyStatus(*r.DestStatus); err != nil {
			return fmt.Errorf("invalid dest_status: %w", err)
		}
	}

	return nil
}

type TransferResponse struct {
	UUID                     string     `json:"uuid"`
	SourceOccupancyUUID      string     `json:"source_occupancy_uuid"`
	DestinationOccupancyUUID string     `json:"dest_occupancy_uuid"`
	Amount                   int64      `json:"amount"`
	AmountUnit               string     `json:"amount_unit"`
	LossAmount               *int64     `json:"loss_amount,omitempty"`
	LossUnit                 *string    `json:"loss_unit,omitempty"`
	StartedAt                time.Time  `json:"started_at"`
	EndedAt                  *time.Time `json:"ended_at,omitempty"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	DeletedAt                *time.Time `json:"deleted_at,omitempty"`
}

type TransferRecordResponse struct {
	Transfer      TransferResponse  `json:"transfer"`
	DestOccupancy OccupancyResponse `json:"dest_occupancy"`
}

func NewTransferResponse(transfer storage.Transfer) TransferResponse {
	return TransferResponse{
		UUID:                     transfer.UUID.String(),
		SourceOccupancyUUID:      transfer.SourceOccupancyUUID,
		DestinationOccupancyUUID: transfer.DestOccupancyUUID,
		Amount:                   transfer.Amount,
		AmountUnit:               transfer.AmountUnit,
		LossAmount:               transfer.LossAmount,
		LossUnit:                 transfer.LossUnit,
		StartedAt:                transfer.StartedAt,
		EndedAt:                  transfer.EndedAt,
		CreatedAt:                transfer.CreatedAt,
		UpdatedAt:                transfer.UpdatedAt,
		DeletedAt:                transfer.DeletedAt,
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
