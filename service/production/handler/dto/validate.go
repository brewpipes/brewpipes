package dto

import (
	"fmt"
	"strings"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

func validateRequired(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}

	return nil
}

func validateVolumeUnit(unit string) error {
	switch unit {
	case storage.VolumeUnitML, storage.VolumeUnitUSFlOz, storage.VolumeUnitUKFlOz, storage.VolumeUnitBBL:
		return nil
	default:
		return fmt.Errorf("invalid amount_unit")
	}
}

func validateRelationType(relationType string) error {
	switch relationType {
	case storage.RelationTypeSplit, storage.RelationTypeBlend:
		return nil
	default:
		return fmt.Errorf("invalid relation_type")
	}
}

func validateVesselStatus(status string) error {
	switch status {
	case storage.VesselStatusActive, storage.VesselStatusInactive, storage.VesselStatusRetired:
		return nil
	default:
		return fmt.Errorf("invalid status")
	}
}

func validateLiquidPhase(phase string) error {
	switch phase {
	case storage.LiquidPhaseWater, storage.LiquidPhaseWort, storage.LiquidPhaseBeer:
		return nil
	default:
		return fmt.Errorf("invalid liquid_phase")
	}
}

func validateProcessPhase(phase string) error {
	switch phase {
	case storage.ProcessPhasePlanning,
		storage.ProcessPhaseMashing,
		storage.ProcessPhaseHeating,
		storage.ProcessPhaseBoiling,
		storage.ProcessPhaseCooling,
		storage.ProcessPhaseFermenting,
		storage.ProcessPhaseConditioning,
		storage.ProcessPhasePackaging,
		storage.ProcessPhaseFinished:
		return nil
	default:
		return fmt.Errorf("invalid process_phase")
	}
}

func validateAdditionType(additionType string) error {
	switch additionType {
	case storage.AdditionTypeMalt,
		storage.AdditionTypeHop,
		storage.AdditionTypeYeast,
		storage.AdditionTypeAdjunct,
		storage.AdditionTypeWaterChem,
		storage.AdditionTypeGas,
		storage.AdditionTypeOther:
		return nil
	default:
		return fmt.Errorf("invalid addition_type")
	}
}

func additionTypeRequiresInventory(additionType string) bool {
	switch additionType {
	case storage.AdditionTypeMalt,
		storage.AdditionTypeHop,
		storage.AdditionTypeYeast,
		storage.AdditionTypeAdjunct:
		return true
	default:
		return false
	}
}

func validateOccupancyStatus(status string) error {
	switch status {
	case storage.OccupancyStatusFermenting,
		storage.OccupancyStatusConditioning,
		storage.OccupancyStatusColdCrashing,
		storage.OccupancyStatusDryHopping,
		storage.OccupancyStatusCarbonating,
		storage.OccupancyStatusHolding,
		storage.OccupancyStatusPackaging:
		return nil
	default:
		return fmt.Errorf("invalid status")
	}
}
