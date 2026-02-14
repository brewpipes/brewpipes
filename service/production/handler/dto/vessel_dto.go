package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateVesselRequest struct {
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	Capacity     int64   `json:"capacity"`
	CapacityUnit string  `json:"capacity_unit"`
	Make         *string `json:"make"`
	Model        *string `json:"model"`
	Status       *string `json:"status"`
}

func (r CreateVesselRequest) Validate() error {
	if err := validate.Required(r.Type, "type"); err != nil {
		return err
	}
	if err := validateVesselType(r.Type); err != nil {
		return err
	}
	if err := validate.Required(r.Name, "name"); err != nil {
		return err
	}
	if r.Capacity <= 0 {
		return fmt.Errorf("capacity must be greater than zero")
	}
	if err := validateVolumeUnit(r.CapacityUnit); err != nil {
		return err
	}
	if r.Status != nil {
		if err := validateVesselStatus(*r.Status); err != nil {
			return err
		}
	}

	return nil
}

type UpdateVesselRequest struct {
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	Capacity     int64   `json:"capacity"`
	CapacityUnit string  `json:"capacity_unit"`
	Make         *string `json:"make"`
	Model        *string `json:"model"`
	Status       string  `json:"status"`
}

func (r UpdateVesselRequest) Validate() error {
	if err := validate.Required(r.Type, "type"); err != nil {
		return err
	}
	if err := validateVesselType(r.Type); err != nil {
		return err
	}
	if err := validate.Required(r.Name, "name"); err != nil {
		return err
	}
	if r.Capacity <= 0 {
		return fmt.Errorf("capacity must be greater than zero")
	}
	if err := validateVolumeUnit(r.CapacityUnit); err != nil {
		return err
	}
	if err := validateVesselStatus(r.Status); err != nil {
		return err
	}

	return nil
}

type VesselResponse struct {
	UUID         string     `json:"uuid"`
	Type         string     `json:"type"`
	Name         string     `json:"name"`
	Capacity     int64      `json:"capacity"`
	CapacityUnit string     `json:"capacity_unit"`
	Make         *string    `json:"make,omitempty"`
	Model        *string    `json:"model,omitempty"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewVesselResponse(vessel storage.Vessel) VesselResponse {
	return VesselResponse{
		UUID:         vessel.UUID.String(),
		Type:         vessel.Type,
		Name:         vessel.Name,
		Capacity:     vessel.Capacity,
		CapacityUnit: vessel.CapacityUnit,
		Make:         vessel.Make,
		Model:        vessel.Model,
		Status:       vessel.Status,
		CreatedAt:    vessel.CreatedAt,
		UpdatedAt:    vessel.UpdatedAt,
		DeletedAt:    vessel.DeletedAt,
	}
}

func NewVesselsResponse(vessels []storage.Vessel) []VesselResponse {
	resp := make([]VesselResponse, 0, len(vessels))
	for _, vessel := range vessels {
		resp = append(resp, NewVesselResponse(vessel))
	}
	return resp
}
