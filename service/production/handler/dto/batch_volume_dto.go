package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBatchVolumeRequest struct {
	BatchUUID   string     `json:"batch_uuid"`
	VolumeUUID  string     `json:"volume_uuid"`
	LiquidPhase string     `json:"liquid_phase"`
	PhaseAt     *time.Time `json:"phase_at"`
}

func (r CreateBatchVolumeRequest) Validate() error {
	if err := validate.Required(r.BatchUUID, "batch_uuid"); err != nil {
		return err
	}
	if err := validate.Required(r.VolumeUUID, "volume_uuid"); err != nil {
		return err
	}
	if err := validateLiquidPhase(r.LiquidPhase); err != nil {
		return err
	}

	return nil
}

type BatchVolumeResponse struct {
	UUID        string     `json:"uuid"`
	BatchUUID   string     `json:"batch_uuid"`
	VolumeUUID  string     `json:"volume_uuid"`
	LiquidPhase string     `json:"liquid_phase"`
	PhaseAt     time.Time  `json:"phase_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewBatchVolumeResponse(volume storage.BatchVolume) BatchVolumeResponse {
	return BatchVolumeResponse{
		UUID:        volume.UUID.String(),
		BatchUUID:   volume.BatchUUID,
		VolumeUUID:  volume.VolumeUUID,
		LiquidPhase: volume.LiquidPhase,
		PhaseAt:     volume.PhaseAt,
		CreatedAt:   volume.CreatedAt,
		UpdatedAt:   volume.UpdatedAt,
		DeletedAt:   volume.DeletedAt,
	}
}

func NewBatchVolumesResponse(volumes []storage.BatchVolume) []BatchVolumeResponse {
	resp := make([]BatchVolumeResponse, 0, len(volumes))
	for _, volume := range volumes {
		resp = append(resp, NewBatchVolumeResponse(volume))
	}
	return resp
}
