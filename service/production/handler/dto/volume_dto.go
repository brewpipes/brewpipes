package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateVolumeRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Amount      int64   `json:"amount"`
	AmountUnit  string  `json:"amount_unit"`
}

func (r CreateVolumeRequest) Validate() error {
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if err := validateVolumeUnit(r.AmountUnit); err != nil {
		return err
	}

	return nil
}

type VolumeResponse struct {
	ID          int64      `json:"id"`
	UUID        string     `json:"uuid"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Amount      int64      `json:"amount"`
	AmountUnit  string     `json:"amount_unit"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewVolumeResponse(volume storage.Volume) VolumeResponse {
	return VolumeResponse{
		ID:          volume.ID,
		UUID:        volume.UUID.String(),
		Name:        volume.Name,
		Description: volume.Description,
		Amount:      volume.Amount,
		AmountUnit:  volume.AmountUnit,
		CreatedAt:   volume.CreatedAt,
		UpdatedAt:   volume.UpdatedAt,
		DeletedAt:   volume.DeletedAt,
	}
}

func NewVolumesResponse(volumes []storage.Volume) []VolumeResponse {
	resp := make([]VolumeResponse, 0, len(volumes))
	for _, volume := range volumes {
		resp = append(resp, NewVolumeResponse(volume))
	}
	return resp
}
