package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBatchVolumeRequest struct {
	BatchID     int64      `json:"batch_id"`
	VolumeID    int64      `json:"volume_id"`
	LiquidPhase string     `json:"liquid_phase"`
	PhaseAt     *time.Time `json:"phase_at"`
}

func (r CreateBatchVolumeRequest) Validate() error {
	if r.BatchID <= 0 || r.VolumeID <= 0 {
		return fmt.Errorf("batch_id and volume_id are required")
	}
	if err := validateLiquidPhase(r.LiquidPhase); err != nil {
		return err
	}

	return nil
}

type BatchVolumeResponse struct {
	ID          int64      `json:"id"`
	UUID        string     `json:"uuid"`
	BatchID     int64      `json:"batch_id"`
	VolumeID    int64      `json:"volume_id"`
	LiquidPhase string     `json:"liquid_phase"`
	PhaseAt     time.Time  `json:"phase_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewBatchVolumeResponse(volume storage.BatchVolume) BatchVolumeResponse {
	return BatchVolumeResponse{
		ID:          volume.ID,
		UUID:        volume.UUID.String(),
		BatchID:     volume.BatchID,
		VolumeID:    volume.VolumeID,
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
