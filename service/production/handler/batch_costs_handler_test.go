package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
	"github.com/gofrs/uuid/v5"
)

// mockBatchCostsStore implements handler.BatchCostsStore for testing.
type mockBatchCostsStore struct {
	batch        storage.Batch
	batchErr     error
	summary      storage.BatchSummary
	summaryErr   error
	additions    []storage.Addition
	additionsErr error
}

func (m *mockBatchCostsStore) GetBatchByUUID(_ context.Context, _ string) (storage.Batch, error) {
	return m.batch, m.batchErr
}

func (m *mockBatchCostsStore) GetBatchSummaryByUUID(_ context.Context, _ string) (storage.BatchSummary, error) {
	return m.summary, m.summaryErr
}

func (m *mockBatchCostsStore) ListAdditionsByBatchUUID(_ context.Context, _ string) ([]storage.Addition, error) {
	return m.additions, m.additionsErr
}

// mockIngredientLotFetcher implements handler.IngredientLotFetcher for testing.
type mockIngredientLotFetcher struct {
	lots []handler.BatchIngredientLot
	err  error
}

func (m *mockIngredientLotFetcher) GetBatchIngredientLots(_ context.Context, _ string, _ string) ([]handler.BatchIngredientLot, error) {
	return m.lots, m.err
}

// mockPOLineFetcher implements handler.POLineFetcher for testing.
type mockPOLineFetcher struct {
	lines []handler.PurchaseOrderLineCost
	err   error
}

func (m *mockPOLineFetcher) BatchLookupPOLines(_ context.Context, _ string, _ []string) ([]handler.PurchaseOrderLineCost, error) {
	return m.lines, m.err
}

// helper to create a *uuid.UUID from a string.
func uuidPtr(s string) *uuid.UUID {
	u := uuid.Must(uuid.FromString(s))
	return &u
}

// helper to create a *string.
func sp(s string) *string {
	return &s
}

// helper to create a *int64.
func i64p(v int64) *int64 {
	return &v
}

func TestHandleBatchCosts(t *testing.T) {
	batchUUID := "550e8400-e29b-41d4-a716-446655440000"
	additionUUID1 := uuid.Must(uuid.FromString("660e8400-e29b-41d4-a716-446655440001"))
	additionUUID2 := uuid.Must(uuid.FromString("660e8400-e29b-41d4-a716-446655440002"))
	lotUUID1 := "770e8400-e29b-41d4-a716-446655440001"
	lotUUID2 := "770e8400-e29b-41d4-a716-446655440002"
	poLineUUID1 := "880e8400-e29b-41d4-a716-446655440001"
	poLineUUID2 := "880e8400-e29b-41d4-a716-446655440002"

	baseBatch := storage.Batch{
		ShortName: "IPA 24-07",
	}
	baseBatch.UUID = uuid.Must(uuid.FromString(batchUUID))

	// Summary with a 10 BBL starting volume.
	summaryWith10BBL := storage.BatchSummary{
		Batch: baseBatch,
		Volumes: []storage.BatchVolumeWithAmount{
			{
				BatchVolume: storage.BatchVolume{
					LiquidPhase: "wort",
					PhaseAt:     time.Date(2024, 7, 1, 8, 0, 0, 0, time.UTC),
				},
				Volume: storage.Volume{
					Amount:     10,
					AmountUnit: "bbl",
				},
			},
		},
	}

	tests := []struct {
		name           string
		batchUUID      string
		store          *mockBatchCostsStore
		invClient      *mockIngredientLotFetcher
		procClient     *mockPOLineFetcher
		expectedStatus int
		validate       func(t *testing.T, resp dto.BatchCostsResponse)
	}{
		{
			name:      "batch not found",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batchErr: service.ErrNotFound,
			},
			invClient:      &mockIngredientLotFetcher{},
			procClient:     &mockPOLineFetcher{},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "batch with no additions",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:     baseBatch,
				summary:   storage.BatchSummary{Batch: baseBatch},
				additions: []storage.Addition{},
			},
			invClient:      &mockIngredientLotFetcher{},
			procClient:     &mockPOLineFetcher{},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, resp dto.BatchCostsResponse) {
				if len(resp.LineItems) != 0 {
					t.Errorf("expected 0 line items, got %d", len(resp.LineItems))
				}
				if len(resp.UncostedAdditions) != 0 {
					t.Errorf("expected 0 uncosted additions, got %d", len(resp.UncostedAdditions))
				}
				if !resp.Totals.CostComplete {
					t.Error("expected cost_complete=true for empty additions")
				}
			},
		},
		{
			name:      "happy path: all additions costed",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:   baseBatch,
				summary: summaryWith10BBL,
				additions: []storage.Addition{
					{
						AdditionType:     "malt",
						Amount:           100,
						AmountUnit:       "kg",
						InventoryLotUUID: uuidPtr(lotUUID1),
					},
					{
						AdditionType:     "hop",
						Amount:           5,
						AmountUnit:       "kg",
						InventoryLotUUID: uuidPtr(lotUUID2),
					},
				},
			},
			invClient: &mockIngredientLotFetcher{
				lots: []handler.BatchIngredientLot{
					{
						IngredientLotUUID:     lotUUID1,
						IngredientUUID:        "aaa00000-0000-0000-0000-000000000001",
						IngredientName:        "Pale Malt",
						IngredientCategory:    "fermentable",
						PurchaseOrderLineUUID: sp(poLineUUID1),
					},
					{
						IngredientLotUUID:     lotUUID2,
						IngredientUUID:        "aaa00000-0000-0000-0000-000000000002",
						IngredientName:        "Cascade",
						IngredientCategory:    "hop",
						PurchaseOrderLineUUID: sp(poLineUUID2),
					},
				},
			},
			procClient: &mockPOLineFetcher{
				lines: []handler.PurchaseOrderLineCost{
					{UUID: poLineUUID1, UnitCostCents: 50, Quantity: 1000, QuantityUnit: "kg", Currency: "USD"},
					{UUID: poLineUUID2, UnitCostCents: 200, Quantity: 50, QuantityUnit: "kg", Currency: "USD"},
				},
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, resp dto.BatchCostsResponse) {
				if len(resp.LineItems) != 2 {
					t.Fatalf("expected 2 line items, got %d", len(resp.LineItems))
				}
				if len(resp.UncostedAdditions) != 0 {
					t.Errorf("expected 0 uncosted additions, got %d", len(resp.UncostedAdditions))
				}
				if !resp.Totals.CostComplete {
					t.Error("expected cost_complete=true")
				}
				// 100*50 + 5*200 = 5000 + 1000 = 6000
				if resp.Totals.TotalCostCents != 6000 {
					t.Errorf("expected total_cost_cents=6000, got %d", resp.Totals.TotalCostCents)
				}
				if resp.Totals.CostedLineCount != 2 {
					t.Errorf("expected costed_line_count=2, got %d", resp.Totals.CostedLineCount)
				}
				if resp.Currency == nil || *resp.Currency != "USD" {
					t.Errorf("expected currency=USD, got %v", resp.Currency)
				}
				for _, item := range resp.LineItems {
					if item.CostSource != "purchase_order" {
						t.Errorf("expected cost_source=purchase_order, got %s", item.CostSource)
					}
				}
			},
		},
		{
			name:      "addition without inventory lot UUID",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:   baseBatch,
				summary: summaryWith10BBL,
				additions: []storage.Addition{
					{
						AdditionType: "gas",
						Amount:       10,
						AmountUnit:   "l",
						// InventoryLotUUID is nil
					},
				},
			},
			invClient:      &mockIngredientLotFetcher{lots: []handler.BatchIngredientLot{}},
			procClient:     &mockPOLineFetcher{},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, resp dto.BatchCostsResponse) {
				if len(resp.LineItems) != 0 {
					t.Errorf("expected 0 line items, got %d", len(resp.LineItems))
				}
				if len(resp.UncostedAdditions) != 1 {
					t.Fatalf("expected 1 uncosted addition, got %d", len(resp.UncostedAdditions))
				}
				if resp.UncostedAdditions[0].Reason != "no_inventory_lot" {
					t.Errorf("expected reason=no_inventory_lot, got %s", resp.UncostedAdditions[0].Reason)
				}
				if resp.Totals.CostComplete {
					t.Error("expected cost_complete=false")
				}
			},
		},
		{
			name:      "lot without PO line link",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:   baseBatch,
				summary: summaryWith10BBL,
				additions: []storage.Addition{
					{
						AdditionType:     "malt",
						Amount:           100,
						AmountUnit:       "kg",
						InventoryLotUUID: uuidPtr(lotUUID1),
					},
				},
			},
			invClient: &mockIngredientLotFetcher{
				lots: []handler.BatchIngredientLot{
					{
						IngredientLotUUID:     lotUUID1,
						IngredientUUID:        "aaa00000-0000-0000-0000-000000000001",
						IngredientName:        "Pale Malt",
						IngredientCategory:    "fermentable",
						PurchaseOrderLineUUID: nil, // No PO line
					},
				},
			},
			procClient:     &mockPOLineFetcher{},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, resp dto.BatchCostsResponse) {
				if len(resp.LineItems) != 1 {
					t.Fatalf("expected 1 line item, got %d", len(resp.LineItems))
				}
				if resp.LineItems[0].CostSource != "unavailable" {
					t.Errorf("expected cost_source=unavailable, got %s", resp.LineItems[0].CostSource)
				}
				if resp.LineItems[0].CostCents != nil {
					t.Errorf("expected nil cost_cents, got %d", *resp.LineItems[0].CostCents)
				}
			},
		},
		{
			name:      "unit mismatch between addition and PO line",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:   baseBatch,
				summary: summaryWith10BBL,
				additions: []storage.Addition{
					{
						AdditionType:     "malt",
						Amount:           500,
						AmountUnit:       "g", // grams
						InventoryLotUUID: uuidPtr(lotUUID1),
					},
				},
			},
			invClient: &mockIngredientLotFetcher{
				lots: []handler.BatchIngredientLot{
					{
						IngredientLotUUID:     lotUUID1,
						IngredientUUID:        "aaa00000-0000-0000-0000-000000000001",
						IngredientName:        "Pale Malt",
						IngredientCategory:    "fermentable",
						PurchaseOrderLineUUID: sp(poLineUUID1),
					},
				},
			},
			procClient: &mockPOLineFetcher{
				lines: []handler.PurchaseOrderLineCost{
					{UUID: poLineUUID1, UnitCostCents: 50, Quantity: 1000, QuantityUnit: "kg", Currency: "USD"}, // kg != g
				},
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, resp dto.BatchCostsResponse) {
				if len(resp.LineItems) != 1 {
					t.Fatalf("expected 1 line item, got %d", len(resp.LineItems))
				}
				if resp.LineItems[0].CostSource != "unavailable" {
					t.Errorf("expected cost_source=unavailable, got %s", resp.LineItems[0].CostSource)
				}
				if resp.LineItems[0].CostCents != nil {
					t.Errorf("expected nil cost_cents for unit mismatch")
				}
			},
		},
		{
			name:      "mixed currencies",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:   baseBatch,
				summary: summaryWith10BBL,
				additions: []storage.Addition{
					{
						AdditionType:     "malt",
						Amount:           100,
						AmountUnit:       "kg",
						InventoryLotUUID: uuidPtr(lotUUID1),
					},
					{
						AdditionType:     "hop",
						Amount:           5,
						AmountUnit:       "kg",
						InventoryLotUUID: uuidPtr(lotUUID2),
					},
				},
			},
			invClient: &mockIngredientLotFetcher{
				lots: []handler.BatchIngredientLot{
					{
						IngredientLotUUID:     lotUUID1,
						IngredientUUID:        "aaa00000-0000-0000-0000-000000000001",
						IngredientName:        "Pale Malt",
						IngredientCategory:    "fermentable",
						PurchaseOrderLineUUID: sp(poLineUUID1),
					},
					{
						IngredientLotUUID:     lotUUID2,
						IngredientUUID:        "aaa00000-0000-0000-0000-000000000002",
						IngredientName:        "Cascade",
						IngredientCategory:    "hop",
						PurchaseOrderLineUUID: sp(poLineUUID2),
					},
				},
			},
			procClient: &mockPOLineFetcher{
				lines: []handler.PurchaseOrderLineCost{
					{UUID: poLineUUID1, UnitCostCents: 50, Quantity: 1000, QuantityUnit: "kg", Currency: "USD"},
					{UUID: poLineUUID2, UnitCostCents: 200, Quantity: 50, QuantityUnit: "kg", Currency: "EUR"}, // Different currency
				},
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, resp dto.BatchCostsResponse) {
				if resp.Currency == nil || *resp.Currency != "MIXED" {
					t.Errorf("expected currency=MIXED, got %v", resp.Currency)
				}
				if resp.Totals.CostPerBBLCents != nil {
					t.Errorf("expected cost_per_bbl_cents=nil for mixed currencies, got %d", *resp.Totals.CostPerBBLCents)
				}
			},
		},
		{
			name:      "inventory service failure",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:   baseBatch,
				summary: summaryWith10BBL,
				additions: []storage.Addition{
					{
						AdditionType:     "malt",
						Amount:           100,
						AmountUnit:       "kg",
						InventoryLotUUID: uuidPtr(lotUUID1),
					},
				},
			},
			invClient: &mockIngredientLotFetcher{
				err: errors.New("inventory service unavailable"),
			},
			procClient:     &mockPOLineFetcher{},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "procurement service failure",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:   baseBatch,
				summary: summaryWith10BBL,
				additions: []storage.Addition{
					{
						AdditionType:     "malt",
						Amount:           100,
						AmountUnit:       "kg",
						InventoryLotUUID: uuidPtr(lotUUID1),
					},
				},
			},
			invClient: &mockIngredientLotFetcher{
				lots: []handler.BatchIngredientLot{
					{
						IngredientLotUUID:     lotUUID1,
						IngredientUUID:        "aaa00000-0000-0000-0000-000000000001",
						IngredientName:        "Pale Malt",
						IngredientCategory:    "fermentable",
						PurchaseOrderLineUUID: sp(poLineUUID1),
					},
				},
			},
			procClient: &mockPOLineFetcher{
				err: errors.New("procurement service unavailable"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "cost per BBL calculation",
			batchUUID: batchUUID,
			store: &mockBatchCostsStore{
				batch:   baseBatch,
				summary: summaryWith10BBL,
				additions: []storage.Addition{
					{
						AdditionType:     "malt",
						Amount:           100,
						AmountUnit:       "kg",
						InventoryLotUUID: uuidPtr(lotUUID1),
					},
				},
			},
			invClient: &mockIngredientLotFetcher{
				lots: []handler.BatchIngredientLot{
					{
						IngredientLotUUID:     lotUUID1,
						IngredientUUID:        "aaa00000-0000-0000-0000-000000000001",
						IngredientName:        "Pale Malt",
						IngredientCategory:    "fermentable",
						PurchaseOrderLineUUID: sp(poLineUUID1),
					},
				},
			},
			procClient: &mockPOLineFetcher{
				lines: []handler.PurchaseOrderLineCost{
					{UUID: poLineUUID1, UnitCostCents: 50, Quantity: 1000, QuantityUnit: "kg", Currency: "USD"},
				},
			},
			expectedStatus: http.StatusOK,
			validate: func(t *testing.T, resp dto.BatchCostsResponse) {
				// Total cost: 100 * 50 = 5000 cents
				// Batch volume: 10 BBL
				// Cost per BBL: round(5000 / 10) = 500
				if resp.Totals.TotalCostCents != 5000 {
					t.Errorf("expected total_cost_cents=5000, got %d", resp.Totals.TotalCostCents)
				}
				if resp.Totals.BatchVolumeBBL == nil {
					t.Fatal("expected batch_volume_bbl to be set")
				}
				if *resp.Totals.BatchVolumeBBL != 10.0 {
					t.Errorf("expected batch_volume_bbl=10.0, got %f", *resp.Totals.BatchVolumeBBL)
				}
				if resp.Totals.CostPerBBLCents == nil {
					t.Fatal("expected cost_per_bbl_cents to be set")
				}
				expected := int64(math.Round(5000.0 / 10.0))
				if *resp.Totals.CostPerBBLCents != expected {
					t.Errorf("expected cost_per_bbl_cents=%d, got %d", expected, *resp.Totals.CostPerBBLCents)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set UUIDs on additions that need them (the handler reads addition.UUID).
			for i := range tt.store.additions {
				if tt.store.additions[i].UUID == uuid.Nil {
					if i == 0 {
						tt.store.additions[i].UUID = additionUUID1
					} else {
						tt.store.additions[i].UUID = additionUUID2
					}
				}
			}

			h := handler.HandleBatchCosts(tt.store, tt.invClient, tt.procClient)

			req := httptest.NewRequest(http.MethodGet, "/batches/"+tt.batchUUID+"/costs", nil)
			req.SetPathValue("uuid", tt.batchUUID)
			req.Header.Set("Authorization", "Bearer test-token")
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}

			if tt.validate != nil && rec.Code == http.StatusOK {
				var resp dto.BatchCostsResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				tt.validate(t, resp)
			}
		})
	}
}
