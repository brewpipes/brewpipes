package dto_test

import (
	"strings"
	"testing"

	"github.com/brewpipes/brewpipes/service/inventory/handler/dto"
)

func TestCreateBatchUsageRequest_Validate(t *testing.T) {
	validPick := dto.BatchUsagePick{
		IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		Amount:            100,
		AmountUnit:        "kg",
	}

	validUsedAt := "2026-01-15T10:00:00Z"

	tests := []struct {
		name    string
		req     dto.CreateBatchUsageRequest
		wantErr string // substring expected in error; empty means no error
	}{
		{
			name: "valid request with one pick",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks:  []dto.BatchUsagePick{validPick},
			},
		},
		{
			name: "valid request with production_ref_uuid",
			req: dto.CreateBatchUsageRequest{
				ProductionRefUUID: strPtr("cccccccc-cccc-cccc-cccc-cccccccccccc"),
				UsedAt:            validUsedAt,
				Picks:             []dto.BatchUsagePick{validPick},
			},
		},
		{
			name: "invalid production_ref_uuid",
			req: dto.CreateBatchUsageRequest{
				ProductionRefUUID: strPtr("not-a-uuid"),
				UsedAt:            validUsedAt,
				Picks:             []dto.BatchUsagePick{validPick},
			},
			wantErr: "production_ref_uuid must be a valid UUID",
		},
		{
			name: "invalid used_at",
			req: dto.CreateBatchUsageRequest{
				UsedAt: "not-a-timestamp",
				Picks:  []dto.BatchUsagePick{validPick},
			},
			wantErr: "used_at must be a valid RFC3339 timestamp",
		},
		{
			name: "empty picks",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks:  []dto.BatchUsagePick{},
			},
			wantErr: "picks must not be empty",
		},
		{
			name: "nil picks",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks:  nil,
			},
			wantErr: "picks must not be empty",
		},
		{
			name: "pick missing ingredient_lot_uuid",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks: []dto.BatchUsagePick{{
					StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					Amount:            100,
					AmountUnit:        "kg",
				}},
			},
			wantErr: "picks[0].ingredient_lot_uuid is required",
		},
		{
			name: "pick missing stock_location_uuid",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks: []dto.BatchUsagePick{{
					IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					Amount:            100,
					AmountUnit:        "kg",
				}},
			},
			wantErr: "picks[0].stock_location_uuid is required",
		},
		{
			name: "pick with zero amount",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks: []dto.BatchUsagePick{{
					IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					Amount:            0,
					AmountUnit:        "kg",
				}},
			},
			wantErr: "picks[0].amount must be greater than zero",
		},
		{
			name: "pick with negative amount",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks: []dto.BatchUsagePick{{
					IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					Amount:            -5,
					AmountUnit:        "kg",
				}},
			},
			wantErr: "picks[0].amount must be greater than zero",
		},
		{
			name: "pick missing amount_unit",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks: []dto.BatchUsagePick{{
					IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
					StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
					Amount:            100,
				}},
			},
			wantErr: "picks[0].amount_unit is required",
		},
		{
			name: "duplicate lot+location picks",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks: []dto.BatchUsagePick{
					{
						IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
						StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
						Amount:            50,
						AmountUnit:        "kg",
					},
					{
						IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
						StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
						Amount:            30,
						AmountUnit:        "kg",
					},
				},
			},
			wantErr: "picks[1] duplicates lot+location from picks[0]",
		},
		{
			name: "same lot different locations is valid",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks: []dto.BatchUsagePick{
					{
						IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
						StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
						Amount:            50,
						AmountUnit:        "kg",
					},
					{
						IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
						StockLocationUUID: "cccccccc-cccc-cccc-cccc-cccccccccccc",
						Amount:            30,
						AmountUnit:        "kg",
					},
				},
			},
		},
		{
			name: "multiple valid picks",
			req: dto.CreateBatchUsageRequest{
				UsedAt: validUsedAt,
				Picks: []dto.BatchUsagePick{
					{
						IngredientLotUUID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
						StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
						Amount:            50,
						AmountUnit:        "kg",
					},
					{
						IngredientLotUUID: "dddddddd-dddd-dddd-dddd-dddddddddddd",
						StockLocationUUID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
						Amount:            10,
						AmountUnit:        "kg",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
