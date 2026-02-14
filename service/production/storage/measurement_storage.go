package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateMeasurement(ctx context.Context, measurement Measurement) (Measurement, error) {
	targetCount := 0
	if measurement.BatchID != nil {
		targetCount++
	}
	if measurement.OccupancyID != nil {
		targetCount++
	}
	if measurement.VolumeID != nil {
		targetCount++
	}
	if targetCount != 1 {
		return Measurement{}, fmt.Errorf("measurement must reference exactly one of batch, occupancy, or volume")
	}

	observedAt := measurement.ObservedAt
	if observedAt.IsZero() {
		observedAt = time.Now().UTC()
	}

	err := c.DB().QueryRow(ctx, `
		INSERT INTO measurement (
			batch_id,
			occupancy_id,
			volume_id,
			kind,
			value,
			unit,
			observed_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, uuid, batch_id, occupancy_id, volume_id, kind, value, unit, observed_at, notes, created_at, updated_at, deleted_at`,
		measurement.BatchID,
		measurement.OccupancyID,
		measurement.VolumeID,
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
		&measurement.VolumeID,
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

	// Resolve FK UUIDs
	if measurement.BatchID != nil {
		var batchUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, *measurement.BatchID).Scan(&batchUUID); err == nil {
			measurement.BatchUUID = &batchUUID
		}
	}
	if measurement.OccupancyID != nil {
		var occUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM occupancy WHERE id = $1`, *measurement.OccupancyID).Scan(&occUUID); err == nil {
			measurement.OccupancyUUID = &occUUID
		}
	}
	if measurement.VolumeID != nil {
		var volUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM volume WHERE id = $1`, *measurement.VolumeID).Scan(&volUUID); err == nil {
			measurement.VolumeUUID = &volUUID
		}
	}

	return measurement, nil
}

const measurementSelectWithJoins = `
	SELECT m.id, m.uuid, m.batch_id, b.uuid, m.occupancy_id, o.uuid, m.volume_id, v.uuid,
	       m.kind, m.value, m.unit, m.observed_at, m.notes,
	       m.created_at, m.updated_at, m.deleted_at
	FROM measurement m
	LEFT JOIN batch b ON b.id = m.batch_id
	LEFT JOIN occupancy o ON o.id = m.occupancy_id
	LEFT JOIN volume v ON v.id = m.volume_id`

func scanMeasurement(row pgx.Row) (Measurement, error) {
	var measurement Measurement
	err := row.Scan(
		&measurement.ID,
		&measurement.UUID,
		&measurement.BatchID,
		&measurement.BatchUUID,
		&measurement.OccupancyID,
		&measurement.OccupancyUUID,
		&measurement.VolumeID,
		&measurement.VolumeUUID,
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
		return Measurement{}, err
	}
	return measurement, nil
}

func (c *Client) GetMeasurement(ctx context.Context, id int64) (Measurement, error) {
	measurement, err := scanMeasurement(c.DB().QueryRow(ctx,
		measurementSelectWithJoins+` WHERE m.id = $1 AND m.deleted_at IS NULL`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Measurement{}, service.ErrNotFound
		}
		return Measurement{}, fmt.Errorf("getting measurement: %w", err)
	}
	return measurement, nil
}

func (c *Client) GetMeasurementByUUID(ctx context.Context, measurementUUID string) (Measurement, error) {
	measurement, err := scanMeasurement(c.DB().QueryRow(ctx,
		measurementSelectWithJoins+` WHERE m.uuid = $1 AND m.deleted_at IS NULL`, measurementUUID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Measurement{}, service.ErrNotFound
		}
		return Measurement{}, fmt.Errorf("getting measurement by uuid: %w", err)
	}
	return measurement, nil
}

func (c *Client) listMeasurements(ctx context.Context, query string, args ...any) ([]Measurement, error) {
	rows, err := c.DB().Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var measurements []Measurement
	for rows.Next() {
		var measurement Measurement
		if err := rows.Scan(
			&measurement.ID,
			&measurement.UUID,
			&measurement.BatchID,
			&measurement.BatchUUID,
			&measurement.OccupancyID,
			&measurement.OccupancyUUID,
			&measurement.VolumeID,
			&measurement.VolumeUUID,
			&measurement.Kind,
			&measurement.Value,
			&measurement.Unit,
			&measurement.ObservedAt,
			&measurement.Notes,
			&measurement.CreatedAt,
			&measurement.UpdatedAt,
			&measurement.DeletedAt,
		); err != nil {
			return nil, err
		}
		measurements = append(measurements, measurement)
	}
	return measurements, rows.Err()
}

func (c *Client) ListMeasurementsByBatch(ctx context.Context, batchID int64) ([]Measurement, error) {
	measurements, err := c.listMeasurements(ctx, measurementSelectWithJoins+`
		WHERE m.deleted_at IS NULL
		AND (
			m.batch_id = $1 OR
			EXISTS (
				SELECT 1
				FROM occupancy oo
				JOIN batch_volume bv ON bv.volume_id = oo.volume_id
				WHERE oo.id = m.occupancy_id
				AND oo.deleted_at IS NULL
				AND bv.deleted_at IS NULL
				AND bv.batch_id = $1
			) OR
			EXISTS (
				SELECT 1
				FROM batch_volume bv
				WHERE bv.volume_id = m.volume_id
				AND bv.deleted_at IS NULL
				AND bv.batch_id = $1
			)
		)
		ORDER BY m.observed_at ASC`, batchID)
	if err != nil {
		return nil, fmt.Errorf("listing measurements by batch: %w", err)
	}
	return measurements, nil
}

func (c *Client) ListMeasurementsByBatchUUID(ctx context.Context, batchUUID string) ([]Measurement, error) {
	batch, err := c.GetBatchByUUID(ctx, batchUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving batch uuid: %w", err)
	}
	return c.ListMeasurementsByBatch(ctx, batch.ID)
}

func (c *Client) ListMeasurementsByOccupancy(ctx context.Context, occupancyID int64) ([]Measurement, error) {
	measurements, err := c.listMeasurements(ctx, measurementSelectWithJoins+`
		WHERE m.occupancy_id = $1 AND m.deleted_at IS NULL
		ORDER BY m.observed_at ASC`, occupancyID)
	if err != nil {
		return nil, fmt.Errorf("listing measurements by occupancy: %w", err)
	}
	return measurements, nil
}

func (c *Client) ListMeasurementsByOccupancyUUID(ctx context.Context, occupancyUUID string) ([]Measurement, error) {
	occ, err := c.GetOccupancyByUUID(ctx, occupancyUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving occupancy uuid: %w", err)
	}
	return c.ListMeasurementsByOccupancy(ctx, occ.ID)
}

func (c *Client) ListMeasurementsByVolume(ctx context.Context, volumeID int64) ([]Measurement, error) {
	measurements, err := c.listMeasurements(ctx, measurementSelectWithJoins+`
		WHERE m.volume_id = $1 AND m.deleted_at IS NULL
		ORDER BY m.observed_at ASC`, volumeID)
	if err != nil {
		return nil, fmt.Errorf("listing measurements by volume: %w", err)
	}
	return measurements, nil
}

func (c *Client) ListMeasurementsByVolumeUUID(ctx context.Context, volumeUUID string) ([]Measurement, error) {
	vol, err := c.GetVolumeByUUID(ctx, volumeUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving volume uuid: %w", err)
	}
	return c.ListMeasurementsByVolume(ctx, vol.ID)
}
