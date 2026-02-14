package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateMeasurementRequest struct {
	BatchUUID     *string    `json:"batch_uuid"`
	OccupancyUUID *string    `json:"occupancy_uuid"`
	VolumeUUID    *string    `json:"volume_uuid"`
	Kind          string     `json:"kind"`
	Value         float64    `json:"value"`
	Unit          *string    `json:"unit"`
	ObservedAt    *time.Time `json:"observed_at"`
	Notes         *string    `json:"notes"`
}

func (r CreateMeasurementRequest) Validate() error {
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
	if err := validate.Required(r.Kind, "kind"); err != nil {
		return err
	}

	return nil
}

type MeasurementResponse struct {
	UUID          string     `json:"uuid"`
	BatchUUID     *string    `json:"batch_uuid,omitempty"`
	OccupancyUUID *string    `json:"occupancy_uuid,omitempty"`
	VolumeUUID    *string    `json:"volume_uuid,omitempty"`
	Kind          string     `json:"kind"`
	Value         float64    `json:"value"`
	Unit          *string    `json:"unit,omitempty"`
	ObservedAt    time.Time  `json:"observed_at"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

func NewMeasurementResponse(measurement storage.Measurement) MeasurementResponse {
	return MeasurementResponse{
		UUID:          measurement.UUID.String(),
		BatchUUID:     measurement.BatchUUID,
		OccupancyUUID: measurement.OccupancyUUID,
		VolumeUUID:    measurement.VolumeUUID,
		Kind:          measurement.Kind,
		Value:         measurement.Value,
		Unit:          measurement.Unit,
		ObservedAt:    measurement.ObservedAt,
		Notes:         measurement.Notes,
		CreatedAt:     measurement.CreatedAt,
		UpdatedAt:     measurement.UpdatedAt,
		DeletedAt:     measurement.DeletedAt,
	}
}

func NewMeasurementsResponse(measurements []storage.Measurement) []MeasurementResponse {
	resp := make([]MeasurementResponse, 0, len(measurements))
	for _, measurement := range measurements {
		resp = append(resp, NewMeasurementResponse(measurement))
	}
	return resp
}
