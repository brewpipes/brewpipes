package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateMeasurementRequest struct {
	BatchID     *int64     `json:"batch_id"`
	OccupancyID *int64     `json:"occupancy_id"`
	VolumeID    *int64     `json:"volume_id"`
	Kind        string     `json:"kind"`
	Value       float64    `json:"value"`
	Unit        *string    `json:"unit"`
	ObservedAt  *time.Time `json:"observed_at"`
	Notes       *string    `json:"notes"`
}

func (r CreateMeasurementRequest) Validate() error {
	targetCount := 0
	if r.BatchID != nil {
		targetCount++
	}
	if r.OccupancyID != nil {
		targetCount++
	}
	if r.VolumeID != nil {
		targetCount++
	}
	if targetCount != 1 {
		return fmt.Errorf("exactly one of batch_id, occupancy_id, or volume_id is required")
	}
	if err := validateRequired(r.Kind, "kind"); err != nil {
		return err
	}

	return nil
}

type MeasurementResponse struct {
	ID          int64      `json:"id"`
	UUID        string     `json:"uuid"`
	BatchID     *int64     `json:"batch_id,omitempty"`
	OccupancyID *int64     `json:"occupancy_id,omitempty"`
	VolumeID    *int64     `json:"volume_id,omitempty"`
	Kind        string     `json:"kind"`
	Value       float64    `json:"value"`
	Unit        *string    `json:"unit,omitempty"`
	ObservedAt  time.Time  `json:"observed_at"`
	Notes       *string    `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func NewMeasurementResponse(measurement storage.Measurement) MeasurementResponse {
	return MeasurementResponse{
		ID:          measurement.ID,
		UUID:        measurement.UUID.String(),
		BatchID:     measurement.BatchID,
		OccupancyID: measurement.OccupancyID,
		VolumeID:    measurement.VolumeID,
		Kind:        measurement.Kind,
		Value:       measurement.Value,
		Unit:        measurement.Unit,
		ObservedAt:  measurement.ObservedAt,
		Notes:       measurement.Notes,
		CreatedAt:   measurement.CreatedAt,
		UpdatedAt:   measurement.UpdatedAt,
		DeletedAt:   measurement.DeletedAt,
	}
}

func NewMeasurementsResponse(measurements []storage.Measurement) []MeasurementResponse {
	resp := make([]MeasurementResponse, 0, len(measurements))
	for _, measurement := range measurements {
		resp = append(resp, NewMeasurementResponse(measurement))
	}
	return resp
}
