package dto

import (
	"fmt"
	"strings"

	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

func validatePurchaseOrderStatus(status string) error {
	switch status {
	case storage.PurchaseOrderStatusDraft,
		storage.PurchaseOrderStatusSubmitted,
		storage.PurchaseOrderStatusConfirmed,
		storage.PurchaseOrderStatusPartiallyReceived,
		storage.PurchaseOrderStatusReceived,
		storage.PurchaseOrderStatusCancelled:
		return nil
	default:
		return fmt.Errorf("invalid status")
	}
}

func validateLineItemType(itemType string) error {
	switch itemType {
	case storage.PurchaseOrderItemTypeIngredient,
		storage.PurchaseOrderItemTypePackaging,
		storage.PurchaseOrderItemTypeService,
		storage.PurchaseOrderItemTypeEquipment,
		storage.PurchaseOrderItemTypeOther:
		return nil
	default:
		return fmt.Errorf("invalid item_type")
	}
}

func validateCurrency(code string) error {
	value := strings.TrimSpace(code)
	if value == "" {
		return fmt.Errorf("currency is required")
	}
	if len(value) != 3 {
		return fmt.Errorf("currency must be a 3-letter code")
	}

	return nil
}
