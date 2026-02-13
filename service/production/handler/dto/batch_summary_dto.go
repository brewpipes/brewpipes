package dto

import (
	"math"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

// Measurement kinds used for ABV calculation
const (
	MeasurementKindOG  = "og"
	MeasurementKindFG  = "fg"
	MeasurementKindABV = "abv"
	MeasurementKindIBU = "ibu"
)

// BatchSummaryResponse provides an aggregated view of batch data with derived metrics.
type BatchSummaryResponse struct {
	// Core batch info
	UUID      string     `json:"uuid"`
	ShortName string     `json:"short_name"`
	BrewDate  *time.Time `json:"brew_date,omitempty"`
	Notes     *string    `json:"notes,omitempty"`

	// Recipe and style
	RecipeName *string `json:"recipe_name,omitempty"`
	StyleName  *string `json:"style_name,omitempty"`

	// Brew sessions
	BrewSessions []BrewSessionSummary `json:"brew_sessions"`

	// Current state
	CurrentPhase           *string `json:"current_phase,omitempty"`
	CurrentVessel          *string `json:"current_vessel,omitempty"`
	CurrentOccupancyStatus *string `json:"current_occupancy_status,omitempty"`

	// Key measurements
	OriginalGravity *float64 `json:"original_gravity,omitempty"`
	FinalGravity    *float64 `json:"final_gravity,omitempty"`
	ABV             *float64 `json:"abv,omitempty"`
	ABVCalculated   bool     `json:"abv_calculated"`
	IBU             *float64 `json:"ibu,omitempty"`

	// Duration metrics (in days)
	DaysInFermenter  *float64 `json:"days_in_fermenter,omitempty"`
	DaysInBrite      *float64 `json:"days_in_brite,omitempty"`
	DaysGrainToGlass *float64 `json:"days_grain_to_glass,omitempty"`

	// Loss metrics
	TotalLossBBL   *float64 `json:"total_loss_bbl,omitempty"`
	LossPercentage *float64 `json:"loss_percentage,omitempty"`

	// Volume tracking
	StartingVolumeBBL *float64 `json:"starting_volume_bbl,omitempty"`
	CurrentVolumeBBL  *float64 `json:"current_volume_bbl,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// BrewSessionSummary is a simplified view of a brew session for the batch summary.
type BrewSessionSummary struct {
	UUID     string    `json:"uuid"`
	BrewedAt time.Time `json:"brewed_at"`
	Notes    *string   `json:"notes,omitempty"`
}

// NewBatchSummaryResponse creates a BatchSummaryResponse from a BatchSummary.
func NewBatchSummaryResponse(summary storage.BatchSummary) BatchSummaryResponse {
	resp := BatchSummaryResponse{
		UUID:      summary.Batch.UUID.String(),
		ShortName: summary.Batch.ShortName,
		BrewDate:  summary.Batch.BrewDate,
		Notes:     summary.Batch.Notes,
		CreatedAt: summary.Batch.CreatedAt,
		UpdatedAt: summary.Batch.UpdatedAt,
		DeletedAt: summary.Batch.DeletedAt,
	}

	// Recipe and style
	if summary.Recipe != nil {
		resp.RecipeName = &summary.Recipe.Name
		resp.StyleName = summary.Recipe.StyleName
	}

	// Brew sessions
	resp.BrewSessions = make([]BrewSessionSummary, 0, len(summary.BrewSessions))
	for _, session := range summary.BrewSessions {
		resp.BrewSessions = append(resp.BrewSessions, BrewSessionSummary{
			UUID:     session.UUID.String(),
			BrewedAt: session.BrewedAt,
			Notes:    session.Notes,
		})
	}

	// Current state
	if summary.CurrentPhase != nil {
		resp.CurrentPhase = &summary.CurrentPhase.ProcessPhase
	}
	if summary.CurrentVessel != nil {
		resp.CurrentVessel = &summary.CurrentVessel.Name
	}
	resp.CurrentOccupancyStatus = summary.CurrentOccupancyStatus

	// Key measurements
	populateMeasurements(&resp, &summary)

	// Duration metrics
	populateDurations(&resp, &summary)

	// Loss and volume metrics
	populateVolumeMetrics(&resp, &summary)

	return resp
}

// populateMeasurements extracts OG, FG, ABV, and IBU from measurements.
func populateMeasurements(resp *BatchSummaryResponse, summary *storage.BatchSummary) {
	// Get OG (first measurement of kind "og")
	ogMeasurement := summary.GetFirstMeasurementByKind(MeasurementKindOG)
	if ogMeasurement != nil {
		resp.OriginalGravity = &ogMeasurement.Value
	}

	// Get FG (most recent measurement of kind "fg")
	fgMeasurement := summary.GetMeasurementByKind(MeasurementKindFG)
	if fgMeasurement != nil {
		resp.FinalGravity = &fgMeasurement.Value
	}

	// Get ABV - first check for manual measurement, then calculate
	abvMeasurement := summary.GetMeasurementByKind(MeasurementKindABV)
	if abvMeasurement != nil {
		resp.ABV = &abvMeasurement.Value
		resp.ABVCalculated = false
	} else if resp.OriginalGravity != nil && resp.FinalGravity != nil {
		// Calculate ABV using formula: (OG - FG) × 131.25
		abv := (*resp.OriginalGravity - *resp.FinalGravity) * 131.25
		abv = math.Round(abv*100) / 100 // Round to 2 decimal places
		resp.ABV = &abv
		resp.ABVCalculated = true
	}

	// Get IBU (most recent measurement of kind "ibu")
	ibuMeasurement := summary.GetMeasurementByKind(MeasurementKindIBU)
	if ibuMeasurement != nil {
		resp.IBU = &ibuMeasurement.Value
	}
}

// populateDurations calculates time spent in fermenters, brite tanks, and total grain-to-glass.
func populateDurations(resp *BatchSummaryResponse, summary *storage.BatchSummary) {
	now := time.Now().UTC()

	var totalFermenterHours float64
	var totalBriteHours float64
	var firstBrewDate *time.Time
	var lastActivityDate *time.Time

	// Get first brew date
	firstBrewDate = summary.GetFirstBrewDate()

	// Calculate time in fermenters and brite tanks
	for _, occ := range summary.Occupancies {
		startTime := occ.Occupancy.InAt
		endTime := now
		if occ.Occupancy.OutAt != nil {
			endTime = *occ.Occupancy.OutAt
		}

		duration := endTime.Sub(startTime).Hours()

		switch occ.Vessel.Type {
		case storage.VesselTypeFermenter:
			totalFermenterHours += duration
		case storage.VesselTypeBriteTank:
			totalBriteHours += duration
		}

		// Track last activity for grain-to-glass calculation
		if lastActivityDate == nil || endTime.After(*lastActivityDate) {
			activityTime := endTime
			lastActivityDate = &activityTime
		}
	}

	// Convert hours to days
	if totalFermenterHours > 0 {
		days := totalFermenterHours / 24
		days = math.Round(days*10) / 10 // Round to 1 decimal place
		resp.DaysInFermenter = &days
	}

	if totalBriteHours > 0 {
		days := totalBriteHours / 24
		days = math.Round(days*10) / 10
		resp.DaysInBrite = &days
	}

	// Calculate grain-to-glass (from first brew date to now or last activity if finished)
	if firstBrewDate != nil {
		// If the batch is finished, use the last activity date
		// Otherwise, use current time
		endTime := now
		for _, phase := range summary.Phases {
			if phase.ProcessPhase == storage.ProcessPhaseFinished {
				if lastActivityDate != nil {
					endTime = *lastActivityDate
				}
				break
			}
		}
		days := endTime.Sub(*firstBrewDate).Hours() / 24
		days = math.Round(days*10) / 10
		resp.DaysGrainToGlass = &days
	}
}

// populateVolumeMetrics calculates volume and loss metrics.
func populateVolumeMetrics(resp *BatchSummaryResponse, summary *storage.BatchSummary) {
	// Find starting volume (earliest batch volume, typically wort phase)
	var startingVolume *storage.BatchVolumeWithAmount
	var currentVolume *storage.BatchVolumeWithAmount
	for i := range summary.Volumes {
		v := &summary.Volumes[i]
		if startingVolume == nil || v.BatchVolume.PhaseAt.Before(startingVolume.BatchVolume.PhaseAt) {
			startingVolume = v
		}
		if currentVolume == nil || v.BatchVolume.PhaseAt.After(currentVolume.BatchVolume.PhaseAt) {
			currentVolume = v
		}
	}

	// Calculate volumes in BBL
	if startingVolume != nil {
		bbl := convertToBBL(startingVolume.Volume.Amount, startingVolume.Volume.AmountUnit)
		if bbl != nil {
			resp.StartingVolumeBBL = bbl
		}
	}

	if currentVolume != nil {
		bbl := convertToBBL(currentVolume.Volume.Amount, currentVolume.Volume.AmountUnit)
		if bbl != nil {
			resp.CurrentVolumeBBL = bbl
		}
	}

	// Calculate total transfer losses
	var totalLossBBL float64
	for _, transfer := range summary.Transfers {
		if transfer.LossAmount != nil && transfer.LossUnit != nil {
			loss := convertToBBL(*transfer.LossAmount, *transfer.LossUnit)
			if loss != nil {
				totalLossBBL += *loss
			}
		}
	}

	if totalLossBBL > 0 {
		totalLossBBL = math.Round(totalLossBBL*100) / 100
		resp.TotalLossBBL = &totalLossBBL
	}

	// Calculate loss percentage
	if resp.StartingVolumeBBL != nil && *resp.StartingVolumeBBL > 0 && resp.TotalLossBBL != nil {
		pct := (*resp.TotalLossBBL / *resp.StartingVolumeBBL) * 100
		pct = math.Round(pct*10) / 10 // Round to 1 decimal place
		resp.LossPercentage = &pct
	}
}

// convertToBBL converts a volume amount to US barrels.
// Returns nil if the unit is not supported.
func convertToBBL(amount int64, unit string) *float64 {
	const (
		mlPerBBL     = 117347.77 // 31 gallons * 3785.41 ml/gallon
		usFlOzPerBBL = 3968.0    // 31 gallons * 128 oz/gallon
		ukFlOzPerBBL = 4135.68   // 31 gallons * 133.41 uk oz/gallon (approx)
	)

	var bbl float64
	switch unit {
	case storage.VolumeUnitBBL:
		bbl = float64(amount)
	case storage.VolumeUnitML:
		bbl = float64(amount) / mlPerBBL
	case storage.VolumeUnitUSFlOz:
		bbl = float64(amount) / usFlOzPerBBL
	case storage.VolumeUnitUKFlOz:
		bbl = float64(amount) / ukFlOzPerBBL
	default:
		return nil
	}

	bbl = math.Round(bbl*100) / 100
	return &bbl
}
