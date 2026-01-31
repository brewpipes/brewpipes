package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateOccupancyRequest struct {
	VesselID int64      `json:"vessel_id"`
	VolumeID int64      `json:"volume_id"`
	InAt     *time.Time `json:"in_at"`
}

func (r CreateOccupancyRequest) Validate() error {
	if r.VesselID <= 0 || r.VolumeID <= 0 {
		return fmt.Errorf("vessel_id and volume_id are required")
	}

	return nil
}

type OccupancyResponse struct {
	ID        int64      `json:"id"`
	UUID      string     `json:"uuid"`
	VesselID  int64      `json:"vessel_id"`
	VolumeID  int64      `json:"volume_id"`
	InAt      time.Time  `json:"in_at"`
	OutAt     *time.Time `json:"out_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewOccupancyResponse(occupancy storage.Occupancy) OccupancyResponse {
	return OccupancyResponse{
		ID:        occupancy.ID,
		UUID:      occupancy.UUID.String(),
		VesselID:  occupancy.VesselID,
		VolumeID:  occupancy.VolumeID,
		InAt:      occupancy.InAt,
		OutAt:     occupancy.OutAt,
		CreatedAt: occupancy.CreatedAt,
		UpdatedAt: occupancy.UpdatedAt,
		DeletedAt: occupancy.DeletedAt,
	}
}
