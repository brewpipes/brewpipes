package storage

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) CreateBatchVolume(ctx context.Context, batchVolume BatchVolume) (BatchVolume, error) {
	phaseAt := batchVolume.PhaseAt
	if phaseAt.IsZero() {
		phaseAt = time.Now().UTC()
	}

	err := c.db.QueryRow(ctx, `
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

	return batchVolume, nil
}

func (c *Client) ListBatchVolumes(ctx context.Context, batchID int64) ([]BatchVolume, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, batch_id, volume_id, liquid_phase, phase_at, created_at, updated_at, deleted_at
		FROM batch_volume
		WHERE batch_id = $1 AND deleted_at IS NULL
		ORDER BY phase_at ASC`,
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
			&batchVolume.VolumeID,
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
