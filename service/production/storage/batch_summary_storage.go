package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
)

// BatchSummary aggregates batch data for the summary endpoint.
type BatchSummary struct {
	// Core batch info
	Batch Batch

	// Recipe and style (nullable if no recipe assigned)
	Recipe *Recipe

	// Brew sessions for this batch
	BrewSessions []BrewSession

	// Current process phase (most recent by phase_at)
	CurrentPhase *BatchProcessPhase

	// Current vessel (from active occupancy of a batch volume)
	CurrentVessel *Vessel

	// Current occupancy status
	CurrentOccupancyStatus *string

	// All measurements for ABV/OG/FG/IBU extraction
	Measurements []Measurement

	// Volumes for loss calculations
	Volumes []BatchVolumeWithAmount

	// Transfers for loss calculations
	Transfers []Transfer

	// Process phases for duration calculations
	Phases []BatchProcessPhase

	// Occupancies for duration calculations (fermenting/conditioning time)
	Occupancies []OccupancyWithVessel
}

// BatchVolumeWithAmount includes volume amount info for loss calculations.
type BatchVolumeWithAmount struct {
	BatchVolume BatchVolume
	Volume      Volume
}

// OccupancyWithVessel includes vessel info for duration and status tracking.
type OccupancyWithVessel struct {
	Occupancy Occupancy
	Vessel    Vessel
}

// GetBatchSummaryByUUID returns a batch summary by resolving the UUID to an internal ID.
func (c *Client) GetBatchSummaryByUUID(ctx context.Context, batchUUID string) (BatchSummary, error) {
	batch, err := c.GetBatchByUUID(ctx, batchUUID)
	if err != nil {
		return BatchSummary{}, fmt.Errorf("resolving batch uuid: %w", err)
	}

	return c.GetBatchSummary(ctx, batch.ID)
}

func (c *Client) GetBatchSummary(ctx context.Context, batchID int64) (BatchSummary, error) {
	var summary BatchSummary

	// Get batch
	batch, err := c.GetBatch(ctx, batchID)
	if err != nil {
		return BatchSummary{}, fmt.Errorf("getting batch: %w", err)
	}
	summary.Batch = batch

	// Get recipe if assigned (including soft-deleted recipes for historical reference)
	if batch.RecipeID != nil {
		recipe, err := c.getRecipeByID(ctx, *batch.RecipeID, &RecipeQueryOpts{IncludeDeleted: true})
		if err != nil && !errors.Is(err, service.ErrNotFound) {
			return BatchSummary{}, fmt.Errorf("getting recipe: %w", err)
		}
		if err == nil {
			summary.Recipe = &recipe
		}
	}

	// Get brew sessions
	sessions, err := c.ListBrewSessionsByBatch(ctx, batchID)
	if err != nil {
		return BatchSummary{}, fmt.Errorf("listing brew sessions: %w", err)
	}
	summary.BrewSessions = sessions

	// Get process phases
	phases, err := c.ListBatchProcessPhases(ctx, batchID)
	if err != nil {
		return BatchSummary{}, fmt.Errorf("listing batch process phases: %w", err)
	}
	summary.Phases = phases
	if len(phases) > 0 {
		// Most recent phase is the current one
		currentPhase := phases[len(phases)-1]
		summary.CurrentPhase = &currentPhase
	}

	// Get measurements
	measurements, err := c.ListMeasurementsByBatch(ctx, batchID)
	if err != nil {
		return BatchSummary{}, fmt.Errorf("listing measurements: %w", err)
	}
	summary.Measurements = measurements

	// Get batch volumes with amount info
	batchVolumes, err := c.ListBatchVolumes(ctx, batchID)
	if err != nil {
		return BatchSummary{}, fmt.Errorf("listing batch volumes: %w", err)
	}
	for _, bv := range batchVolumes {
		vol, err := c.GetVolume(ctx, bv.VolumeID)
		if err != nil && !errors.Is(err, service.ErrNotFound) {
			return BatchSummary{}, fmt.Errorf("getting volume: %w", err)
		}
		if err == nil {
			summary.Volumes = append(summary.Volumes, BatchVolumeWithAmount{
				BatchVolume: bv,
				Volume:      vol,
			})
		}
	}

	// Get transfers for loss calculations
	transfers, err := c.ListTransfersByBatch(ctx, batchID)
	if err != nil {
		return BatchSummary{}, fmt.Errorf("listing transfers: %w", err)
	}
	summary.Transfers = transfers

	// Get occupancies with vessel info for duration tracking
	occupancies, err := c.listOccupanciesByBatch(ctx, batchID)
	if err != nil {
		return BatchSummary{}, fmt.Errorf("listing occupancies: %w", err)
	}
	summary.Occupancies = occupancies

	// Find current vessel and occupancy status from active occupancy
	for _, occ := range occupancies {
		if occ.Occupancy.OutAt == nil {
			summary.CurrentVessel = &occ.Vessel
			summary.CurrentOccupancyStatus = occ.Occupancy.Status
			break
		}
	}

	return summary, nil
}

// listOccupanciesByBatch returns all occupancies for volumes linked to a batch.
func (c *Client) listOccupanciesByBatch(ctx context.Context, batchID int64) ([]OccupancyWithVessel, error) {
	rows, err := c.db.Query(ctx, `
		SELECT 
			o.id, o.uuid, o.vessel_id, o.volume_id, o.in_at, o.out_at, o.status, o.created_at, o.updated_at, o.deleted_at,
			v.id, v.uuid, v.type, v.name, v.capacity, v.capacity_unit, v.make, v.model, v.status, v.created_at, v.updated_at, v.deleted_at
		FROM occupancy o
		JOIN vessel v ON v.id = o.vessel_id
		WHERE o.deleted_at IS NULL
		AND v.deleted_at IS NULL
		AND EXISTS (
			SELECT 1
			FROM batch_volume bv
			WHERE bv.volume_id = o.volume_id
			AND bv.deleted_at IS NULL
			AND bv.batch_id = $1
		)
		ORDER BY o.in_at ASC`,
		batchID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing occupancies by batch: %w", err)
	}
	defer rows.Close()

	var result []OccupancyWithVessel
	for rows.Next() {
		var occ Occupancy
		var vessel Vessel
		if err := rows.Scan(
			&occ.ID, &occ.UUID, &occ.VesselID, &occ.VolumeID, &occ.InAt, &occ.OutAt, &occ.Status, &occ.CreatedAt, &occ.UpdatedAt, &occ.DeletedAt,
			&vessel.ID, &vessel.UUID, &vessel.Type, &vessel.Name, &vessel.Capacity, &vessel.CapacityUnit, &vessel.Make, &vessel.Model, &vessel.Status, &vessel.CreatedAt, &vessel.UpdatedAt, &vessel.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning occupancy with vessel: %w", err)
		}
		result = append(result, OccupancyWithVessel{
			Occupancy: occ,
			Vessel:    vessel,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing occupancies by batch: %w", err)
	}

	return result, nil
}

// GetFirstBrewDate returns the earliest brew session date for a batch, or nil if none.
func (s *BatchSummary) GetFirstBrewDate() *time.Time {
	if len(s.BrewSessions) == 0 {
		return nil
	}
	earliest := s.BrewSessions[0].BrewedAt
	for _, session := range s.BrewSessions[1:] {
		if session.BrewedAt.Before(earliest) {
			earliest = session.BrewedAt
		}
	}
	return &earliest
}

// GetMeasurementByKind returns the most recent measurement of the specified kind.
func (s *BatchSummary) GetMeasurementByKind(kind string) *Measurement {
	var latest *Measurement
	for i := range s.Measurements {
		m := &s.Measurements[i]
		if m.Kind == kind {
			if latest == nil || m.ObservedAt.After(latest.ObservedAt) {
				latest = m
			}
		}
	}
	return latest
}

// GetFirstMeasurementByKind returns the earliest measurement of the specified kind.
func (s *BatchSummary) GetFirstMeasurementByKind(kind string) *Measurement {
	var earliest *Measurement
	for i := range s.Measurements {
		m := &s.Measurements[i]
		if m.Kind == kind {
			if earliest == nil || m.ObservedAt.Before(earliest.ObservedAt) {
				earliest = m
			}
		}
	}
	return earliest
}
