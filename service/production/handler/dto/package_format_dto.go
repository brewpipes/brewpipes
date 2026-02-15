package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

// Container type constants for package formats.
const (
	ContainerKeg     = "keg"
	ContainerCan     = "can"
	ContainerBottle  = "bottle"
	ContainerCask    = "cask"
	ContainerGrowler = "growler"
	ContainerOther   = "other"
)

type CreatePackageFormatRequest struct {
	Name              string `json:"name"`
	Container         string `json:"container"`
	VolumePerUnit     int64  `json:"volume_per_unit"`
	VolumePerUnitUnit string `json:"volume_per_unit_unit"`
}

func (r CreatePackageFormatRequest) Validate() error {
	if err := validate.Required(r.Name, "name"); err != nil {
		return err
	}
	if err := validate.Required(r.Container, "container"); err != nil {
		return err
	}
	if err := validateContainerType(r.Container); err != nil {
		return err
	}
	if r.VolumePerUnit <= 0 {
		return fmt.Errorf("volume_per_unit must be greater than zero")
	}
	if err := validate.Required(r.VolumePerUnitUnit, "volume_per_unit_unit"); err != nil {
		return err
	}
	if err := validateVolumeUnit(r.VolumePerUnitUnit); err != nil {
		return fmt.Errorf("invalid volume_per_unit_unit: %w", err)
	}
	return nil
}

type UpdatePackageFormatRequest struct {
	Name              *string `json:"name"`
	Container         *string `json:"container"`
	VolumePerUnit     *int64  `json:"volume_per_unit"`
	VolumePerUnitUnit *string `json:"volume_per_unit_unit"`
	IsActive          *bool   `json:"is_active"`
}

func (r UpdatePackageFormatRequest) Validate() error {
	if r.Name != nil {
		if err := validate.Required(*r.Name, "name"); err != nil {
			return err
		}
	}
	if r.Container != nil {
		if err := validateContainerType(*r.Container); err != nil {
			return err
		}
	}
	if r.VolumePerUnit != nil && *r.VolumePerUnit <= 0 {
		return fmt.Errorf("volume_per_unit must be greater than zero")
	}
	if r.VolumePerUnitUnit != nil {
		if err := validate.Required(*r.VolumePerUnitUnit, "volume_per_unit_unit"); err != nil {
			return err
		}
		if err := validateVolumeUnit(*r.VolumePerUnitUnit); err != nil {
			return fmt.Errorf("invalid volume_per_unit_unit: %w", err)
		}
	}
	return nil
}

type PackageFormatResponse struct {
	UUID              string     `json:"uuid"`
	Name              string     `json:"name"`
	Container         string     `json:"container"`
	VolumePerUnit     int64      `json:"volume_per_unit"`
	VolumePerUnitUnit string     `json:"volume_per_unit_unit"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewPackageFormatResponse(format storage.PackageFormat) PackageFormatResponse {
	return PackageFormatResponse{
		UUID:              format.UUID.String(),
		Name:              format.Name,
		Container:         format.Container,
		VolumePerUnit:     format.VolumePerUnit,
		VolumePerUnitUnit: format.VolumePerUnitUnit,
		IsActive:          format.IsActive,
		CreatedAt:         format.CreatedAt,
		UpdatedAt:         format.UpdatedAt,
		DeletedAt:         format.DeletedAt,
	}
}

func NewPackageFormatsResponse(formats []storage.PackageFormat) []PackageFormatResponse {
	resp := make([]PackageFormatResponse, 0, len(formats))
	for _, format := range formats {
		resp = append(resp, NewPackageFormatResponse(format))
	}
	return resp
}

func validateContainerType(container string) error {
	switch container {
	case ContainerKeg, ContainerCan, ContainerBottle, ContainerCask, ContainerGrowler, ContainerOther:
		return nil
	default:
		return fmt.Errorf("invalid container: must be one of keg, can, bottle, cask, growler, other")
	}
}
