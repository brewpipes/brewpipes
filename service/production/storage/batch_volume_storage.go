package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateBatchVolume(ctx context.Context, batchVolume BatchVolume) (BatchVolume, error) {
	phaseAt := batchVolume.PhaseAt
	if phaseAt.IsZero() {
		phaseAt = time.Now().UTC()
	}

	err := c.DB().QueryRow(ctx, `
		INSERT INTO batch_volume (
			batch_id,
			volume_id,
			liquid_phase,
			phase_at
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, batch_id, volume_id, liquid_phase, phase_at, created_at, updated_at, deleted_at`,
		batchVolume.BatchID,
		batchVolume.VolumeID,
		batchVolume.LiquidPhase,
		phaseAt,
	).Scan(
		&batchVolume.ID,
		&batchVolume.UUID,
		&batchVolume.BatchID,
		&batchVolume.VolumeID,
		&batchVolume.LiquidPhase,
		&batchVolume.PhaseAt,
		&batchVolume.CreatedAt,
		&batchVolume.UpdatedAt,
		&batchVolume.DeletedAt,
	)
	if err != nil {
		return BatchVolume{}, fmt.Errorf("creating batch volume: %w", err)
	}

	// Resolve batch and volume UUIDs
	c.DB().QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, batchVolume.BatchID).Scan(&batchVolume.BatchUUID)
	c.DB().QueryRow(ctx, `SELECT uuid FROM volume WHERE id = $1`, batchVolume.VolumeID).Scan(&batchVolume.VolumeUUID)

	return batchVolume, nil
}

func (c *Client) GetBatchVolume(ctx context.Context, id int64) (BatchVolume, error) {
	var batchVolume BatchVolume
	err := c.DB().QueryRow(ctx, `
		SELECT bv.id, bv.uuid, bv.batch_id, b.uuid, bv.volume_id, v.uuid,
		       bv.liquid_phase, bv.phase_at, bv.created_at, bv.updated_at, bv.deleted_at
		FROM batch_volume bv
		JOIN batch b ON b.id = bv.batch_id
		JOIN volume v ON v.id = bv.volume_id
		WHERE bv.id = $1 AND bv.deleted_at IS NULL`,
		id,
	).Scan(
		&batchVolume.ID,
		&batchVolume.UUID,
		&batchVolume.BatchID,
		&batchVolume.BatchUUID,
		&batchVolume.VolumeID,
		&batchVolume.VolumeUUID,
		&batchVolume.LiquidPhase,
		&batchVolume.PhaseAt,
		&batchVolume.CreatedAt,
		&batchVolume.UpdatedAt,
		&batchVolume.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchVolume{}, service.ErrNotFound
		}
		return BatchVolume{}, fmt.Errorf("getting batch volume: %w", err)
	}

	return batchVolume, nil
}

func (c *Client) GetBatchVolumeByUUID(ctx context.Context, bvUUID string) (BatchVolume, error) {
	var batchVolume BatchVolume
	err := c.DB().QueryRow(ctx, `
		SELECT bv.id, bv.uuid, bv.batch_id, b.uuid, bv.volume_id, v.uuid,
		       bv.liquid_phase, bv.phase_at, bv.created_at, bv.updated_at, bv.deleted_at
		FROM batch_volume bv
		JOIN batch b ON b.id = bv.batch_id
		JOIN volume v ON v.id = bv.volume_id
		WHERE bv.uuid = $1 AND bv.deleted_at IS NULL`,
		bvUUID,
	).Scan(
		&batchVolume.ID,
		&batchVolume.UUID,
		&batchVolume.BatchID,
		&batchVolume.BatchUUID,
		&batchVolume.VolumeID,
		&batchVolume.VolumeUUID,
		&batchVolume.LiquidPhase,
		&batchVolume.PhaseAt,
		&batchVolume.CreatedAt,
		&batchVolume.UpdatedAt,
		&batchVolume.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchVolume{}, service.ErrNotFound
		}
		return BatchVolume{}, fmt.Errorf("getting batch volume by uuid: %w", err)
	}

	return batchVolume, nil
}

func (c *Client) ListBatchVolumes(ctx context.Context, batchID int64) ([]BatchVolume, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT bv.id, bv.uuid, bv.batch_id, b.uuid, bv.volume_id, v.uuid,
		       bv.liquid_phase, bv.phase_at, bv.created_at, bv.updated_at, bv.deleted_at
		FROM batch_volume bv
		JOIN batch b ON b.id = bv.batch_id
		JOIN volume v ON v.id = bv.volume_id
		WHERE bv.batch_id = $1 AND bv.deleted_at IS NULL
		ORDER BY bv.phase_at ASC`,
		batchID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing batch volumes: %w", err)
	}
	defer rows.Close()

	var batchVolumes []BatchVolume
	for rows.Next() {
		var batchVolume BatchVolume
		if err := rows.Scan(
			&batchVolume.ID,
			&batchVolume.UUID,
			&batchVolume.BatchID,
			&batchVolume.BatchUUID,
			&batchVolume.VolumeID,
			&batchVolume.VolumeUUID,
			&batchVolume.LiquidPhase,
			&batchVolume.PhaseAt,
			&batchVolume.CreatedAt,
			&batchVolume.UpdatedAt,
			&batchVolume.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning batch volume: %w", err)
		}
		batchVolumes = append(batchVolumes, batchVolume)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing batch volumes: %w", err)
	}

	return batchVolumes, nil
}

func (c *Client) ListBatchVolumesByBatchUUID(ctx context.Context, batchUUID string) ([]BatchVolume, error) {
	batch, err := c.GetBatchByUUID(ctx, batchUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving batch uuid: %w", err)
	}

	return c.ListBatchVolumes(ctx, batch.ID)
}
