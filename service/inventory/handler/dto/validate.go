package dto

import (
	"fmt"
	"strings"

	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

func validateRequired(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}

	return nil
}

func validateIngredientCategory(category string) error {
	switch category {
	case storage.IngredientCategoryFermentable,
		storage.IngredientCategoryHop,
		storage.IngredientCategoryYeast,
		storage.IngredientCategoryAdjunct,
		storage.IngredientCategorySalt,
		storage.IngredientCategoryChemical,
		storage.IngredientCategoryGas,
		storage.IngredientCategoryOther:
		return nil
	default:
		return fmt.Errorf("invalid category")
	}
}

func validateHopForm(form string) error {
	switch form {
	case storage.HopFormPellet,
		storage.HopFormWholeLeaf,
		storage.HopFormCryo,
		storage.HopFormExtract,
		storage.HopFormOther:
		return nil
	default:
		return fmt.Errorf("invalid hop_form")
	}
}

func validateYeastForm(form string) error {
	switch form {
	case storage.YeastFormLiquid,
		storage.YeastFormDry,
		storage.YeastFormSlurry,
		storage.YeastFormPropagated,
		storage.YeastFormOther:
		return nil
	default:
		return fmt.Errorf("invalid yeast_form")
	}
}

func validateStockLocationType(locationType string) error {
	switch locationType {
	case storage.StockLocationTypeDry,
		storage.StockLocationTypeCold,
		storage.StockLocationTypeGas,
		storage.StockLocationTypeBulk,
		storage.StockLocationTypePackaging,
		storage.StockLocationTypeOther:
		return nil
	default:
		return fmt.Errorf("invalid location_type")
	}
}

func validateOriginatorType(originatorType string) error {
	switch originatorType {
	case storage.OriginatorTypeMaltster,
		storage.OriginatorTypeHopProducer,
		storage.OriginatorTypeYeastLab,
		storage.OriginatorTypeGasVendor,
		storage.OriginatorTypeOther:
		return nil
	default:
		return fmt.Errorf("invalid originator_type")
	}
}

func validateMovementDirection(direction string) error {
	switch direction {
	case storage.MovementDirectionIn, storage.MovementDirectionOut:
		return nil
	default:
		return fmt.Errorf("invalid direction")
	}
}

func validateMovementReason(reason string) error {
	switch reason {
	case storage.MovementReasonReceive,
		storage.MovementReasonUse,
		storage.MovementReasonTransfer,
		storage.MovementReasonAdjust,
		storage.MovementReasonWaste:
		return nil
	default:
		return fmt.Errorf("invalid reason")
	}
}

func validateAdjustmentReason(reason string) error {
	switch reason {
	case storage.AdjustmentReasonCycleCount,
		storage.AdjustmentReasonSpoilage,
		storage.AdjustmentReasonShrink,
		storage.AdjustmentReasonDamage,
		storage.AdjustmentReasonCorrection,
		storage.AdjustmentReasonOther:
		return nil
	default:
		return fmt.Errorf("invalid reason")
	}
}
