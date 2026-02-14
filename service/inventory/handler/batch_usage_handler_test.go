package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brewpipes/brewpipes/internal/database/entity"
	"github.com/brewpipes/brewpipes/service/inventory/handler"
	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
	"github.com/gofrs/uuid/v5"
)

// mockBatchUsageStore implements handler.BatchUsageStore for testing.
type mockBatchUsageStore struct {
	CreateBatchUsageFunc func(context.Context, storage.BatchUsageRequest) (storage.BatchUsageResult, error)
}

func (m *mockBatchUsageStore) CreateBatchUsage(ctx context.Context, req storage.BatchUsageRequest) (storage.BatchUsageResult, error) {
	if m.CreateBatchUsageFunc != nil {
		return m.CreateBatchUsageFunc(ctx, req)
	}
	return storage.BatchUsageResult{}, nil
}

func TestHandleCreateBatchUsage(t *testing.T) {
	usageUUID := uuid.Must(uuid.NewV4())
	movementUUID := uuid.Must(uuid.NewV4())
	now := time.Now().UTC().Truncate(time.Second)

	successStore := &mockBatchUsageStore{
		CreateBatchUsageFunc: func(_ context.Context, req storage.BatchUsageRequest) (storage.BatchUsageResult, error) {
			lotUUID := req.Picks[0].IngredientLotUUID
			locUUID := req.Picks[0].StockLocationUUID
			usageUUIDStr := usageUUID.String()
			return storage.BatchUsageResult{
				Usage: storage.InventoryUsage{
					Identifiers: entity.Identifiers{ID: 1, UUID: usageUUID},
					UsedAt:      req.UsedAt,
				},
				Movements: []storage.InventoryMovement{
					{
						Identifiers:       entity.Identifiers{ID: 1, UUID: movementUUID},
						IngredientLotUUID: &lotUUID,
						StockLocationUUID: locUUID,
						Direction:         "out",
						Reason:            "use",
						Amount:            req.Picks[0].Amount,
						AmountUnit:        req.Picks[0].AmountUnit,
						OccurredAt:        req.UsedAt,
						UsageUUID:         &usageUUIDStr,
						Timestamps:        entity.Timestamps{CreatedAt: now, UpdatedAt: now},
					},
				},
			}, nil
		},
	}

	tests := []struct {
		name       string
		store      *mockBatchUsageStore
		body       string
		wantStatus int
		wantBody   string // substring to check in response body
	}{
		{
			name:  "valid request returns 201",
			store: successStore,
			body: fmt.Sprintf(`{
				"used_at": "%s",
				"picks": [{
					"ingredient_lot_uuid": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					"stock_location_uuid": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					"amount": 100,
					"amount_unit": "kg"
				}]
			}`, now.Format(time.RFC3339)),
			wantStatus: http.StatusCreated,
			wantBody:   usageUUID.String(),
		},
		{
			name:       "invalid JSON returns 400",
			store:      successStore,
			body:       `{invalid`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "invalid request",
		},
		{
			name:  "validation failure returns 400",
			store: successStore,
			body: `{
				"used_at": "2026-01-15T10:00:00Z",
				"picks": []
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "picks must not be empty",
		},
		{
			name:  "invalid production_ref_uuid returns 400",
			store: successStore,
			body: `{
				"production_ref_uuid": "not-a-uuid",
				"used_at": "2026-01-15T10:00:00Z",
				"picks": [{
					"ingredient_lot_uuid": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					"stock_location_uuid": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					"amount": 100,
					"amount_unit": "kg"
				}]
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "production_ref_uuid",
		},
		{
			name: "storage validation error returns 400",
			store: &mockBatchUsageStore{
				CreateBatchUsageFunc: func(_ context.Context, _ storage.BatchUsageRequest) (storage.BatchUsageResult, error) {
					return storage.BatchUsageResult{}, &storage.ErrBatchUsageValidation{
						Message: "insufficient stock for lot ABC at Cold Room: available 50 kg, requested 100 kg",
					}
				},
			},
			body: fmt.Sprintf(`{
				"used_at": "%s",
				"picks": [{
					"ingredient_lot_uuid": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					"stock_location_uuid": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					"amount": 100,
					"amount_unit": "kg"
				}]
			}`, now.Format(time.RFC3339)),
			wantStatus: http.StatusBadRequest,
			wantBody:   "insufficient stock",
		},
		{
			name: "storage internal error returns 500",
			store: &mockBatchUsageStore{
				CreateBatchUsageFunc: func(_ context.Context, _ storage.BatchUsageRequest) (storage.BatchUsageResult, error) {
					return storage.BatchUsageResult{}, fmt.Errorf("database connection lost")
				},
			},
			body: fmt.Sprintf(`{
				"used_at": "%s",
				"picks": [{
					"ingredient_lot_uuid": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					"stock_location_uuid": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					"amount": 100,
					"amount_unit": "kg"
				}]
			}`, now.Format(time.RFC3339)),
			wantStatus: http.StatusInternalServerError,
			wantBody:   "server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.HandleCreateBatchUsage(tt.store)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/inventory-usage/batch", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")

			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.wantBody != "" && !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body does not contain %q; got: %s", tt.wantBody, rec.Body.String())
			}
		})
	}
}

func TestHandleCreateBatchUsage_ResponseShape(t *testing.T) {
	usageUUID := uuid.Must(uuid.NewV4())
	movementUUID := uuid.Must(uuid.NewV4())
	now := time.Now().UTC().Truncate(time.Second)

	store := &mockBatchUsageStore{
		CreateBatchUsageFunc: func(_ context.Context, req storage.BatchUsageRequest) (storage.BatchUsageResult, error) {
			lotUUID := req.Picks[0].IngredientLotUUID
			locUUID := req.Picks[0].StockLocationUUID
			usageUUIDStr := usageUUID.String()
			return storage.BatchUsageResult{
				Usage: storage.InventoryUsage{
					Identifiers: entity.Identifiers{ID: 1, UUID: usageUUID},
					UsedAt:      req.UsedAt,
				},
				Movements: []storage.InventoryMovement{
					{
						Identifiers:       entity.Identifiers{ID: 1, UUID: movementUUID},
						IngredientLotUUID: &lotUUID,
						StockLocationUUID: locUUID,
						Direction:         "out",
						Reason:            "use",
						Amount:            100,
						AmountUnit:        "kg",
						OccurredAt:        req.UsedAt,
						UsageUUID:         &usageUUIDStr,
						Timestamps:        entity.Timestamps{CreatedAt: now, UpdatedAt: now},
					},
				},
			}, nil
		},
	}

	h := handler.HandleCreateBatchUsage(store)
	body := fmt.Sprintf(`{
		"used_at": "%s",
		"picks": [{
			"ingredient_lot_uuid": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			"stock_location_uuid": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			"amount": 100,
			"amount_unit": "kg"
		}]
	}`, now.Format(time.RFC3339))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inventory-usage/batch", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rec.Code, rec.Body.String())
	}

	var resp dto.BatchUsageResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.UsageUUID != usageUUID.String() {
		t.Errorf("usage_uuid = %q, want %q", resp.UsageUUID, usageUUID.String())
	}
	if len(resp.Movements) != 1 {
		t.Fatalf("movements count = %d, want 1", len(resp.Movements))
	}
	m := resp.Movements[0]
	if m.UUID != movementUUID.String() {
		t.Errorf("movement uuid = %q, want %q", m.UUID, movementUUID.String())
	}
	if m.Direction != "out" {
		t.Errorf("movement direction = %q, want %q", m.Direction, "out")
	}
	if m.Reason != "use" {
		t.Errorf("movement reason = %q, want %q", m.Reason, "use")
	}
	if m.Amount != 100 {
		t.Errorf("movement amount = %d, want 100", m.Amount)
	}
}
