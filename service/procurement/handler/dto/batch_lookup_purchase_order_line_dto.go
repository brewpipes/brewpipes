package dto

import "fmt"

// BatchLookupPurchaseOrderLinesRequest is the request body for batch-looking up PO lines by UUID.
type BatchLookupPurchaseOrderLinesRequest struct {
	UUIDs []string `json:"uuids"`
}

// Validate checks that the request is well-formed.
func (r BatchLookupPurchaseOrderLinesRequest) Validate() error {
	if len(r.UUIDs) == 0 {
		return fmt.Errorf("uuids must not be empty")
	}
	if len(r.UUIDs) > 100 {
		return fmt.Errorf("uuids must not exceed 100 items")
	}
	return nil
}
