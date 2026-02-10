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

func errPositiveRequired(field string) error {
	return fmt.Errorf("%s must be greater than zero", field)
}

func errRangeRequired(field string, min, max float64) error {
	return fmt.Errorf("%s must be between %.0f and %.0f", field, min, max)
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

func validateVesselType(vesselType string) error {
	switch vesselType {
	case storage.VesselTypeMashTun,
		storage.VesselTypeLauterTun,
		storage.VesselTypeKettle,
		storage.VesselTypeWhirlpool,
		storage.VesselTypeFermenter,
		storage.VesselTypeBriteTank,
		storage.VesselTypeServingTank,
		storage.VesselTypeOther:
		return nil
	default:
		return fmt.Errorf("invalid type")
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

func validateIngredientType(ingredientType string) error {
	switch ingredientType {
	case storage.IngredientTypeFermentable,
		storage.IngredientTypeHop,
		storage.IngredientTypeYeast,
		storage.IngredientTypeAdjunct,
		storage.IngredientTypeSalt,
		storage.IngredientTypeChemical,
		storage.IngredientTypeGas,
		storage.IngredientTypeOther:
		return nil
	default:
		return fmt.Errorf("invalid ingredient_type")
	}
}

func validateUseStage(useStage string) error {
	switch useStage {
	case storage.UseStageMash,
		storage.UseStageBoil,
		storage.UseStageWhirlpool,
		storage.UseStageFermentation,
		storage.UseStagePackaging:
		return nil
	default:
		return fmt.Errorf("invalid use_stage")
	}
}

func validateUseType(useType string) error {
	switch useType {
	case storage.UseTypeBittering,
		storage.UseTypeFlavor,
		storage.UseTypeAroma,
		storage.UseTypeDryHop,
		storage.UseTypeBase,
		storage.UseTypeSpecialty,
		storage.UseTypeAdjunct,
		storage.UseTypeSugar,
		storage.UseTypePrimary,
		storage.UseTypeSecondary,
		storage.UseTypeBottle,
		storage.UseTypeOther:
		return nil
	default:
		return fmt.Errorf("invalid use_type")
	}
}

func validateIBUMethod(method string) error {
	switch method {
	case storage.IBUMethodTinseth,
		storage.IBUMethodRager,
		storage.IBUMethodGaretz,
		storage.IBUMethodDaniels:
		return nil
	default:
		return fmt.Errorf("invalid ibu_method")
	}
}
