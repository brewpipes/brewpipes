package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brewpipes/brewpipes/service/procurement/handler"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
	"github.com/gofrs/uuid/v5"
)

type PurchaseOrderStore struct {
	ListPurchaseOrdersFunc               func(context.Context) ([]storage.PurchaseOrder, error)
	ListPurchaseOrdersBySupplierUUIDFunc func(context.Context, string) ([]storage.PurchaseOrder, error)
	GetPurchaseOrderByUUIDFunc           func(context.Context, string) (storage.PurchaseOrder, error)
	GetSupplierByUUIDFunc                func(context.Context, string) (storage.Supplier, error)
	CreatePurchaseOrderFunc              func(context.Context, storage.PurchaseOrder) (storage.PurchaseOrder, error)
	UpdatePurchaseOrderByUUIDFunc        func(context.Context, string, storage.PurchaseOrderUpdate) (storage.PurchaseOrder, error)
}

func (s PurchaseOrderStore) ListPurchaseOrders(ctx context.Context) ([]storage.PurchaseOrder, error) {
	if s.ListPurchaseOrdersFunc == nil {
		return nil, nil
	}
	return s.ListPurchaseOrdersFunc(ctx)
}

func (s PurchaseOrderStore) ListPurchaseOrdersBySupplierUUID(ctx context.Context, supplierUUID string) ([]storage.PurchaseOrder, error) {
	if s.ListPurchaseOrdersBySupplierUUIDFunc == nil {
		return nil, nil
	}
	return s.ListPurchaseOrdersBySupplierUUIDFunc(ctx, supplierUUID)
}

func (s PurchaseOrderStore) GetPurchaseOrderByUUID(ctx context.Context, orderUUID string) (storage.PurchaseOrder, error) {
	if s.GetPurchaseOrderByUUIDFunc == nil {
		return storage.PurchaseOrder{}, nil
	}
	return s.GetPurchaseOrderByUUIDFunc(ctx, orderUUID)
}

func (s PurchaseOrderStore) GetSupplierByUUID(ctx context.Context, supplierUUID string) (storage.Supplier, error) {
	if s.GetSupplierByUUIDFunc == nil {
		return storage.Supplier{}, nil
	}
	return s.GetSupplierByUUIDFunc(ctx, supplierUUID)
}

func (s PurchaseOrderStore) CreatePurchaseOrder(ctx context.Context, order storage.PurchaseOrder) (storage.PurchaseOrder, error) {
	if s.CreatePurchaseOrderFunc == nil {
		return order, nil
	}
	return s.CreatePurchaseOrderFunc(ctx, order)
}

func (s PurchaseOrderStore) UpdatePurchaseOrderByUUID(ctx context.Context, orderUUID string, update storage.PurchaseOrderUpdate) (storage.PurchaseOrder, error) {
	if s.UpdatePurchaseOrderByUUIDFunc == nil {
		return storage.PurchaseOrder{}, nil
	}
	return s.UpdatePurchaseOrderByUUIDFunc(ctx, orderUUID, update)
}

func TestHandlePurchaseOrderByUUID_StatusTransition(t *testing.T) {
	testUUID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name           string
		currentStatus  string
		newStatus      string
		expectedStatus int
		errContains    string
	}{
		// Valid transitions
		{
			name:           "draft to submitted succeeds",
			currentStatus:  storage.PurchaseOrderStatusDraft,
			newStatus:      storage.PurchaseOrderStatusSubmitted,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "submitted to confirmed succeeds",
			currentStatus:  storage.PurchaseOrderStatusSubmitted,
			newStatus:      storage.PurchaseOrderStatusConfirmed,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "confirmed to received succeeds",
			currentStatus:  storage.PurchaseOrderStatusConfirmed,
			newStatus:      storage.PurchaseOrderStatusReceived,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "draft to cancelled succeeds",
			currentStatus:  storage.PurchaseOrderStatusDraft,
			newStatus:      storage.PurchaseOrderStatusCancelled,
			expectedStatus: http.StatusOK,
		},

		// Invalid transitions
		{
			name:           "draft to confirmed fails",
			currentStatus:  storage.PurchaseOrderStatusDraft,
			newStatus:      storage.PurchaseOrderStatusConfirmed,
			expectedStatus: http.StatusConflict,
			errContains:    "transition not allowed",
		},
		{
			name:           "submitted to draft fails",
			currentStatus:  storage.PurchaseOrderStatusSubmitted,
			newStatus:      storage.PurchaseOrderStatusDraft,
			expectedStatus: http.StatusConflict,
			errContains:    "transition not allowed",
		},

		// Terminal state: received
		{
			name:           "received to draft fails",
			currentStatus:  storage.PurchaseOrderStatusReceived,
			newStatus:      storage.PurchaseOrderStatusDraft,
			expectedStatus: http.StatusConflict,
			errContains:    "purchase order is already complete",
		},
		{
			name:           "received to cancelled fails",
			currentStatus:  storage.PurchaseOrderStatusReceived,
			newStatus:      storage.PurchaseOrderStatusCancelled,
			expectedStatus: http.StatusConflict,
			errContains:    "purchase order is already complete",
		},

		// Terminal state: cancelled
		{
			name:           "cancelled to draft fails",
			currentStatus:  storage.PurchaseOrderStatusCancelled,
			newStatus:      storage.PurchaseOrderStatusDraft,
			expectedStatus: http.StatusConflict,
			errContains:    "purchase order is cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := PurchaseOrderStore{
				GetPurchaseOrderByUUIDFunc: func(ctx context.Context, orderUUID string) (storage.PurchaseOrder, error) {
					return storage.PurchaseOrder{
						Status: tt.currentStatus,
					}, nil
				},
				UpdatePurchaseOrderByUUIDFunc: func(ctx context.Context, orderUUID string, update storage.PurchaseOrderUpdate) (storage.PurchaseOrder, error) {
					return storage.PurchaseOrder{
						Status: *update.Status,
					}, nil
				},
			}

			h := handler.HandlePurchaseOrderByUUID(store)

			body, _ := json.Marshal(map[string]string{"status": tt.newStatus})
			req := httptest.NewRequest(http.MethodPatch, "/purchase-orders/"+testUUID.String(), bytes.NewReader(body))
			req.SetPathValue("uuid", testUUID.String())
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.errContains != "" && !bytes.Contains(rec.Body.Bytes(), []byte(tt.errContains)) {
				t.Errorf("expected body to contain %q, got %q", tt.errContains, rec.Body.String())
			}
		})
	}
}

func TestHandlePurchaseOrderByUUID_NoStatusChange(t *testing.T) {
	testUUID := uuid.Must(uuid.NewV4())

	// When status is not in the request, no validation should occur
	store := PurchaseOrderStore{
		GetPurchaseOrderByUUIDFunc: func(ctx context.Context, orderUUID string) (storage.PurchaseOrder, error) {
			return storage.PurchaseOrder{
				Status:      storage.PurchaseOrderStatusReceived, // Terminal state
				OrderNumber: "PO-001",
			}, nil
		},
		UpdatePurchaseOrderByUUIDFunc: func(ctx context.Context, orderUUID string, update storage.PurchaseOrderUpdate) (storage.PurchaseOrder, error) {
			return storage.PurchaseOrder{
				Status:      storage.PurchaseOrderStatusReceived,
				OrderNumber: *update.OrderNumber,
			}, nil
		},
	}

	h := handler.HandlePurchaseOrderByUUID(store)

	// Update only order_number, not status
	body, _ := json.Marshal(map[string]string{"order_number": "PO-002"})
	req := httptest.NewRequest(http.MethodPatch, "/purchase-orders/"+testUUID.String(), bytes.NewReader(body))
	req.SetPathValue("uuid", testUUID.String())
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandlePurchaseOrderByUUID_SameStatus(t *testing.T) {
	testUUID := uuid.Must(uuid.NewV4())

	// When status is the same, no validation error should occur
	store := PurchaseOrderStore{
		GetPurchaseOrderByUUIDFunc: func(ctx context.Context, orderUUID string) (storage.PurchaseOrder, error) {
			return storage.PurchaseOrder{
				Status: storage.PurchaseOrderStatusReceived, // Terminal state
			}, nil
		},
		UpdatePurchaseOrderByUUIDFunc: func(ctx context.Context, orderUUID string, update storage.PurchaseOrderUpdate) (storage.PurchaseOrder, error) {
			return storage.PurchaseOrder{
				Status: *update.Status,
			}, nil
		},
	}

	h := handler.HandlePurchaseOrderByUUID(store)

	// Update to same status (no-op)
	body, _ := json.Marshal(map[string]string{"status": storage.PurchaseOrderStatusReceived})
	req := httptest.NewRequest(http.MethodPatch, "/purchase-orders/"+testUUID.String(), bytes.NewReader(body))
	req.SetPathValue("uuid", testUUID.String())
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
}
