package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateOccupancyRequest struct {
	VesselUUID string     `json:"vessel_uuid"`
	VolumeUUID string     `json:"volume_uuid"`
	InAt       *time.Time `json:"in_at"`
	Status     *string    `json:"status"`
}

func (r CreateOccupancyRequest) Validate() error {
	if err := validate.Required(r.VesselUUID, "vessel_uuid"); err != nil {
		return err
	}
	if err := validate.Required(r.VolumeUUID, "volume_uuid"); err != nil {
		return err
	}
	if r.Status != nil {
		if err := validateOccupancyStatus(*r.Status); err != nil {
			return err
		}
	}

	return nil
}

type UpdateOccupancyStatusRequest struct {
	Status *string `json:"status"`
}

func (r UpdateOccupancyStatusRequest) Validate() error {
	if r.Status != nil {
		if err := validateOccupancyStatus(*r.Status); err != nil {
			return err
		}
	}
	return nil
}

type OccupancyResponse struct {
	UUID       string     `json:"uuid"`
	VesselUUID string     `json:"vessel_uuid"`
	VolumeUUID string     `json:"volume_uuid"`
	BatchUUID  *string    `json:"batch_uuid,omitempty"`
	InAt       time.Time  `json:"in_at"`
	OutAt      *time.Time `json:"out_at,omitempty"`
	Status     *string    `json:"status,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

func NewOccupancyResponse(occupancy storage.Occupancy) OccupancyResponse {
	return OccupancyResponse{
		UUID:       occupancy.UUID.String(),
		VesselUUID: occupancy.VesselUUID,
		VolumeUUID: occupancy.VolumeUUID,
		BatchUUID:  occupancy.BatchUUID,
		InAt:       occupancy.InAt,
		OutAt:      occupancy.OutAt,
		Status:     occupancy.Status,
		CreatedAt:  occupancy.CreatedAt,
		UpdatedAt:  occupancy.UpdatedAt,
		DeletedAt:  occupancy.DeletedAt,
	}
}

func NewOccupanciesResponse(occupancies []storage.Occupancy) []OccupancyResponse {
	resp := make([]OccupancyResponse, 0, len(occupancies))
	for _, occupancy := range occupancies {
		resp = append(resp, NewOccupancyResponse(occupancy))
	}
	return resp
}
