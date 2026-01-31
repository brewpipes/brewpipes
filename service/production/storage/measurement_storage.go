package storage

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) CreateMeasurement(ctx context.Context, measurement Measurement) (Measurement, error) {
	if (measurement.BatchID == nil && measurement.OccupancyID == nil) || (measurement.BatchID != nil && measurement.OccupancyID != nil) {
		return Measurement{}, fmt.Errorf("measurement must reference batch or occupancy")
	}

	observedAt := measurement.ObservedAt
	if observedAt.IsZero() {
		observedAt = time.Now().UTC()
	}

	err := c.db.QueryRow(ctx, `
		INSERT INTO measurement (
			batch_id,
			occupancy_id,
			kind,
			value,
			unit,
			observed_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, batch_id, occupancy_id, kind, value, unit, observed_at, notes, created_at, updated_at, deleted_at`,
		measurement.BatchID,
		measurement.OccupancyID,
		measurement.Kind,
		measurement.Value,
		measurement.Unit,
		observedAt,
		measurement.Notes,
	).Scan(
		&measurement.ID,
		&measurement.UUID,
		&measurement.BatchID,
		&measurement.OccupancyID,
		&measurement.Kind,
		&measurement.Value,
		&measurement.Unit,
		&measurement.ObservedAt,
		&measurement.Notes,
		&measurement.CreatedAt,
		&measurement.UpdatedAt,
		&measurement.DeletedAt,
	)
	if err != nil {
		return Measurement{}, fmt.Errorf("creating measurement: %w", err)
	}

	return measurement, nil
}

func (c *Client) ListMeasurementsByBatch(ctx context.Context, batchID int64) ([]Measurement, error) {
	rows, err := c.db.Query(ctx, `
		SELECT m.id, m.uuid, m.batch_id, m.occupancy_id, m.kind, m.value, m.unit, m.observed_at, m.notes, m.created_at, m.updated_at, m.deleted_at
		FROM measurement m
		WHERE m.deleted_at IS NULL
		AND (
			m.batch_id = $1 OR
			EXISTS (
				SELECT 1
				FROM occupancy o
				JOIN batch_volume bv ON bv.volume_id = o.volume_id
				WHERE o.id = m.occupancy_id
				AND o.deleted_at IS NULL
				AND bv.deleted_at IS NULL
				AND bv.batch_id = $1
			)
		)
		ORDER BY m.observed_at ASC`,
		batchID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing measurements by batch: %w", err)
	}
	defer rows.Close()

	var measurements []Measurement
	for rows.Next() {
		var measurement Measurement
		if err := rows.Scan(
			&measurement.ID,
			&measurement.UUID,
			&measurement.BatchID,
			&measurement.OccupancyID,
			&measurement.Kind,
			&measurement.Value,
			&measurement.Unit,
			&measurement.ObservedAt,
			&measurement.Notes,
			&measurement.CreatedAt,
			&measurement.UpdatedAt,
			&measurement.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning measurement: %w", err)
		}
		measurements = append(measurements, measurement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing measurements by batch: %w", err)
	}

	return measurements, nil
}

func (c *Client) ListMeasurementsByOccupancy(ctx context.Context, occupancyID int64) ([]Measurement, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, batch_id, occupancy_id, kind, value, unit, observed_at, notes, created_at, updated_at, deleted_at
		FROM measurement
		WHERE occupancy_id = $1 AND deleted_at IS NULL
		ORDER BY observed_at ASC`,
		occupancyID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing measurements by occupancy: %w", err)
	}
	defer rows.Close()

	var measurements []Measurement
	for rows.Next() {
		var measurement Measurement
		if err := rows.Scan(
			&measurement.ID,
			&measurement.UUID,
			&measurement.BatchID,
			&measurement.OccupancyID,
			&measurement.Kind,
			&measurement.Value,
			&measurement.Unit,
			&measurement.ObservedAt,
			&measurement.Notes,
			&measurement.CreatedAt,
			&measurement.UpdatedAt,
			&measurement.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning measurement: %w", err)
		}
		measurements = append(measurements, measurement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing measurements by occupancy: %w", err)
	}

	return measurements, nil
}
