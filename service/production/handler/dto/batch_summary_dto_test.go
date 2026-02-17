package dto

import (
	"math"
	"testing"
	"time"

	"github.com/brewpipes/brewpipes/internal/database/entity"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

func TestPopulateMeasurements(t *testing.T) {
	now := time.Now().UTC()

	makeMeasurement := func(id int64, kind string, value float64, observedAt time.Time) storage.Measurement {
		return storage.Measurement{
			Identifiers: entity.Identifiers{ID: id},
			Kind:        kind,
			Value:       value,
			ObservedAt:  observedAt,
		}
	}

	tests := []struct {
		name         string
		measurements []storage.Measurement
		wantOG       *float64
		wantFG       *float64
		wantABV      *float64
		wantABVCalc  bool
		wantIBU      *float64
	}{
		{
			name:         "no measurements",
			measurements: nil,
			wantOG:       nil,
			wantFG:       nil,
			wantABV:      nil,
			wantABVCalc:  false,
			wantIBU:      nil,
		},
		{
			name: "explicit og and fg kinds",
			measurements: []storage.Measurement{
				makeMeasurement(1, "og", 1.054, now.Add(-7*24*time.Hour)),
				makeMeasurement(2, "fg", 1.012, now),
			},
			wantOG:      floatPtr(1.054),
			wantFG:      floatPtr(1.012),
			wantABV:     floatPtr(math.Round((1.054-1.012)*131.25*100) / 100),
			wantABVCalc: true,
			wantIBU:     nil,
		},
		{
			name: "gravity kind time series — OG is first, FG is last",
			measurements: []storage.Measurement{
				makeMeasurement(1, "gravity", 1.054, now.Add(-7*24*time.Hour)),
				makeMeasurement(2, "gravity", 1.038, now.Add(-5*24*time.Hour)),
				makeMeasurement(3, "gravity", 1.020, now.Add(-3*24*time.Hour)),
				makeMeasurement(4, "gravity", 1.012, now),
			},
			wantOG:      floatPtr(1.054),
			wantFG:      floatPtr(1.012),
			wantABV:     floatPtr(math.Round((1.054-1.012)*131.25*100) / 100),
			wantABVCalc: true,
			wantIBU:     nil,
		},
		{
			name: "single gravity measurement — OG only, no FG",
			measurements: []storage.Measurement{
				makeMeasurement(1, "gravity", 1.054, now),
			},
			wantOG:      floatPtr(1.054),
			wantFG:      nil,
			wantABV:     nil,
			wantABVCalc: false,
			wantIBU:     nil,
		},
		{
			name: "explicit og/fg take precedence over gravity",
			measurements: []storage.Measurement{
				makeMeasurement(1, "og", 1.060, now.Add(-7*24*time.Hour)),
				makeMeasurement(2, "gravity", 1.054, now.Add(-6*24*time.Hour)),
				makeMeasurement(3, "gravity", 1.020, now.Add(-1*24*time.Hour)),
				makeMeasurement(4, "fg", 1.014, now),
			},
			wantOG:      floatPtr(1.060),
			wantFG:      floatPtr(1.014),
			wantABV:     floatPtr(math.Round((1.060-1.014)*131.25*100) / 100),
			wantABVCalc: true,
			wantIBU:     nil,
		},
		{
			name: "manual ABV overrides calculated",
			measurements: []storage.Measurement{
				makeMeasurement(1, "gravity", 1.054, now.Add(-7*24*time.Hour)),
				makeMeasurement(2, "gravity", 1.012, now),
				makeMeasurement(3, "abv", 5.5, now),
			},
			wantOG:      floatPtr(1.054),
			wantFG:      floatPtr(1.012),
			wantABV:     floatPtr(5.5),
			wantABVCalc: false,
			wantIBU:     nil,
		},
		{
			name: "IBU measurement",
			measurements: []storage.Measurement{
				makeMeasurement(1, "ibu", 45.0, now),
			},
			wantOG:      nil,
			wantFG:      nil,
			wantABV:     nil,
			wantABVCalc: false,
			wantIBU:     floatPtr(45.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &BatchSummaryResponse{}
			summary := &storage.BatchSummary{
				Measurements: tt.measurements,
			}

			populateMeasurements(resp, summary)

			assertFloat64Ptr(t, "OriginalGravity", tt.wantOG, resp.OriginalGravity)
			assertFloat64Ptr(t, "FinalGravity", tt.wantFG, resp.FinalGravity)
			assertFloat64Ptr(t, "ABV", tt.wantABV, resp.ABV)
			if resp.ABVCalculated != tt.wantABVCalc {
				t.Errorf("ABVCalculated = %v, want %v", resp.ABVCalculated, tt.wantABVCalc)
			}
			assertFloat64Ptr(t, "IBU", tt.wantIBU, resp.IBU)
		})
	}
}

func floatPtr(v float64) *float64 {
	return &v
}

func assertFloat64Ptr(t *testing.T, field string, want, got *float64) {
	t.Helper()
	if want == nil && got == nil {
		return
	}
	if want == nil && got != nil {
		t.Errorf("%s = %v, want nil", field, *got)
		return
	}
	if want != nil && got == nil {
		t.Errorf("%s = nil, want %v", field, *want)
		return
	}
	if math.Abs(*want-*got) > 0.0001 {
		t.Errorf("%s = %v, want %v", field, *got, *want)
	}
}
