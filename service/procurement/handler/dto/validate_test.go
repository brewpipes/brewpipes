package dto_test

import (
	"strings"
	"testing"

	"github.com/brewpipes/brewpipes/service/procurement/handler/dto"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

func TestValidatePurchaseOrderStatusTransition(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		newStatus     string
		wantErr       bool
		errContains   string
	}{
		// Same status (no-op) - always valid
		{
			name:          "draft to draft is valid",
			currentStatus: storage.PurchaseOrderStatusDraft,
			newStatus:     storage.PurchaseOrderStatusDraft,
			wantErr:       false,
		},
		{
			name:          "received to received is valid",
			currentStatus: storage.PurchaseOrderStatusReceived,
			newStatus:     storage.PurchaseOrderStatusReceived,
			wantErr:       false,
		},

		// Valid transitions from draft
		{
			name:          "draft to submitted is valid",
			currentStatus: storage.PurchaseOrderStatusDraft,
			newStatus:     storage.PurchaseOrderStatusSubmitted,
			wantErr:       false,
		},
		{
			name:          "draft to cancelled is valid",
			currentStatus: storage.PurchaseOrderStatusDraft,
			newStatus:     storage.PurchaseOrderStatusCancelled,
			wantErr:       false,
		},

		// Valid transitions from submitted
		{
			name:          "submitted to confirmed is valid",
			currentStatus: storage.PurchaseOrderStatusSubmitted,
			newStatus:     storage.PurchaseOrderStatusConfirmed,
			wantErr:       false,
		},
		{
			name:          "submitted to cancelled is valid",
			currentStatus: storage.PurchaseOrderStatusSubmitted,
			newStatus:     storage.PurchaseOrderStatusCancelled,
			wantErr:       false,
		},

		// Valid transitions from confirmed
		{
			name:          "confirmed to partially_received is valid",
			currentStatus: storage.PurchaseOrderStatusConfirmed,
			newStatus:     storage.PurchaseOrderStatusPartiallyReceived,
			wantErr:       false,
		},
		{
			name:          "confirmed to received is valid",
			currentStatus: storage.PurchaseOrderStatusConfirmed,
			newStatus:     storage.PurchaseOrderStatusReceived,
			wantErr:       false,
		},
		{
			name:          "confirmed to cancelled is valid",
			currentStatus: storage.PurchaseOrderStatusConfirmed,
			newStatus:     storage.PurchaseOrderStatusCancelled,
			wantErr:       false,
		},

		// Valid transitions from partially_received
		{
			name:          "partially_received to received is valid",
			currentStatus: storage.PurchaseOrderStatusPartiallyReceived,
			newStatus:     storage.PurchaseOrderStatusReceived,
			wantErr:       false,
		},
		{
			name:          "partially_received to cancelled is valid",
			currentStatus: storage.PurchaseOrderStatusPartiallyReceived,
			newStatus:     storage.PurchaseOrderStatusCancelled,
			wantErr:       false,
		},

		// Invalid transitions from draft
		{
			name:          "draft to confirmed is invalid",
			currentStatus: storage.PurchaseOrderStatusDraft,
			newStatus:     storage.PurchaseOrderStatusConfirmed,
			wantErr:       true,
			errContains:   "transition not allowed",
		},
		{
			name:          "draft to received is invalid",
			currentStatus: storage.PurchaseOrderStatusDraft,
			newStatus:     storage.PurchaseOrderStatusReceived,
			wantErr:       true,
			errContains:   "transition not allowed",
		},

		// Invalid transitions from submitted
		{
			name:          "submitted to draft is invalid",
			currentStatus: storage.PurchaseOrderStatusSubmitted,
			newStatus:     storage.PurchaseOrderStatusDraft,
			wantErr:       true,
			errContains:   "transition not allowed",
		},
		{
			name:          "submitted to received is invalid",
			currentStatus: storage.PurchaseOrderStatusSubmitted,
			newStatus:     storage.PurchaseOrderStatusReceived,
			wantErr:       true,
			errContains:   "transition not allowed",
		},

		// Invalid transitions from confirmed
		{
			name:          "confirmed to draft is invalid",
			currentStatus: storage.PurchaseOrderStatusConfirmed,
			newStatus:     storage.PurchaseOrderStatusDraft,
			wantErr:       true,
			errContains:   "transition not allowed",
		},
		{
			name:          "confirmed to submitted is invalid",
			currentStatus: storage.PurchaseOrderStatusConfirmed,
			newStatus:     storage.PurchaseOrderStatusSubmitted,
			wantErr:       true,
			errContains:   "transition not allowed",
		},

		// Terminal state: received
		{
			name:          "received to draft is invalid",
			currentStatus: storage.PurchaseOrderStatusReceived,
			newStatus:     storage.PurchaseOrderStatusDraft,
			wantErr:       true,
			errContains:   "purchase order is already complete",
		},
		{
			name:          "received to submitted is invalid",
			currentStatus: storage.PurchaseOrderStatusReceived,
			newStatus:     storage.PurchaseOrderStatusSubmitted,
			wantErr:       true,
			errContains:   "purchase order is already complete",
		},
		{
			name:          "received to cancelled is invalid",
			currentStatus: storage.PurchaseOrderStatusReceived,
			newStatus:     storage.PurchaseOrderStatusCancelled,
			wantErr:       true,
			errContains:   "purchase order is already complete",
		},

		// Terminal state: cancelled
		{
			name:          "cancelled to draft is invalid",
			currentStatus: storage.PurchaseOrderStatusCancelled,
			newStatus:     storage.PurchaseOrderStatusDraft,
			wantErr:       true,
			errContains:   "purchase order is cancelled",
		},
		{
			name:          "cancelled to submitted is invalid",
			currentStatus: storage.PurchaseOrderStatusCancelled,
			newStatus:     storage.PurchaseOrderStatusSubmitted,
			wantErr:       true,
			errContains:   "purchase order is cancelled",
		},
		{
			name:          "cancelled to received is invalid",
			currentStatus: storage.PurchaseOrderStatusCancelled,
			newStatus:     storage.PurchaseOrderStatusReceived,
			wantErr:       true,
			errContains:   "purchase order is cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dto.ValidatePurchaseOrderStatusTransition(tt.currentStatus, tt.newStatus)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}
