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

// validPurchaseOrderTransitions defines the allowed status transitions for purchase orders.
// Terminal states (received, cancelled) have no valid transitions.
var validPurchaseOrderTransitions = map[string][]string{
	storage.PurchaseOrderStatusDraft:             {storage.PurchaseOrderStatusSubmitted, storage.PurchaseOrderStatusCancelled},
	storage.PurchaseOrderStatusSubmitted:         {storage.PurchaseOrderStatusConfirmed, storage.PurchaseOrderStatusCancelled},
	storage.PurchaseOrderStatusConfirmed:         {storage.PurchaseOrderStatusPartiallyReceived, storage.PurchaseOrderStatusReceived, storage.PurchaseOrderStatusCancelled},
	storage.PurchaseOrderStatusPartiallyReceived: {storage.PurchaseOrderStatusReceived, storage.PurchaseOrderStatusCancelled},
	storage.PurchaseOrderStatusReceived:          {},
	storage.PurchaseOrderStatusCancelled:         {},
}

// ValidatePurchaseOrderStatusTransition checks if a status transition is allowed.
// Returns nil if the transition is valid, or an error with a clear message if invalid.
func ValidatePurchaseOrderStatusTransition(currentStatus, newStatus string) error {
	// Same status is always allowed (no-op)
	if currentStatus == newStatus {
		return nil
	}

	allowedTransitions, exists := validPurchaseOrderTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("cannot transition from '%s': unknown status", currentStatus)
	}

	// Check if newStatus is in the allowed transitions
	for _, allowed := range allowedTransitions {
		if newStatus == allowed {
			return nil
		}
	}

	// Build a descriptive error message
	reason := terminalStatusReason(currentStatus)
	if reason != "" {
		return fmt.Errorf("cannot transition from '%s' to '%s': %s", currentStatus, newStatus, reason)
	}

	return fmt.Errorf("cannot transition from '%s' to '%s': transition not allowed", currentStatus, newStatus)
}

// terminalStatusReason returns a human-readable reason for terminal states.
func terminalStatusReason(status string) string {
	switch status {
	case storage.PurchaseOrderStatusReceived:
		return "purchase order is already complete"
	case storage.PurchaseOrderStatusCancelled:
		return "purchase order is cancelled"
	default:
		return ""
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
