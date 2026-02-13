package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateBatchProcessPhase(ctx context.Context, phase BatchProcessPhase) (BatchProcessPhase, error) {
	phaseAt := phase.PhaseAt
	if phaseAt.IsZero() {
		phaseAt = time.Now().UTC()
	}

	err := c.db.QueryRow(ctx, `
		INSERT INTO batch_process_phase (
			batch_id,
			process_phase,
			phase_at
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, batch_id, process_phase, phase_at, created_at, updated_at, deleted_at`,
		phase.BatchID,
		phase.ProcessPhase,
		phaseAt,
	).Scan(
		&phase.ID,
		&phase.UUID,
		&phase.BatchID,
		&phase.ProcessPhase,
		&phase.PhaseAt,
		&phase.CreatedAt,
		&phase.UpdatedAt,
		&phase.DeletedAt,
	)
	if err != nil {
		return BatchProcessPhase{}, fmt.Errorf("creating batch process phase: %w", err)
	}

	// Resolve batch UUID
	c.db.QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, phase.BatchID).Scan(&phase.BatchUUID)

	return phase, nil
}

func (c *Client) GetBatchProcessPhase(ctx context.Context, id int64) (BatchProcessPhase, error) {
	var phase BatchProcessPhase
	err := c.db.QueryRow(ctx, `
		SELECT bpp.id, bpp.uuid, bpp.batch_id, b.uuid, bpp.process_phase, bpp.phase_at,
		       bpp.created_at, bpp.updated_at, bpp.deleted_at
		FROM batch_process_phase bpp
		JOIN batch b ON b.id = bpp.batch_id
		WHERE bpp.id = $1 AND bpp.deleted_at IS NULL`,
		id,
	).Scan(
		&phase.ID,
		&phase.UUID,
		&phase.BatchID,
		&phase.BatchUUID,
		&phase.ProcessPhase,
		&phase.PhaseAt,
		&phase.CreatedAt,
		&phase.UpdatedAt,
		&phase.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchProcessPhase{}, service.ErrNotFound
		}
		return BatchProcessPhase{}, fmt.Errorf("getting batch process phase: %w", err)
	}

	return phase, nil
}

func (c *Client) GetBatchProcessPhaseByUUID(ctx context.Context, phaseUUID string) (BatchProcessPhase, error) {
	var phase BatchProcessPhase
	err := c.db.QueryRow(ctx, `
		SELECT bpp.id, bpp.uuid, bpp.batch_id, b.uuid, bpp.process_phase, bpp.phase_at,
		       bpp.created_at, bpp.updated_at, bpp.deleted_at
		FROM batch_process_phase bpp
		JOIN batch b ON b.id = bpp.batch_id
		WHERE bpp.uuid = $1 AND bpp.deleted_at IS NULL`,
		phaseUUID,
	).Scan(
		&phase.ID,
		&phase.UUID,
		&phase.BatchID,
		&phase.BatchUUID,
		&phase.ProcessPhase,
		&phase.PhaseAt,
		&phase.CreatedAt,
		&phase.UpdatedAt,
		&phase.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchProcessPhase{}, service.ErrNotFound
		}
		return BatchProcessPhase{}, fmt.Errorf("getting batch process phase by uuid: %w", err)
	}

	return phase, nil
}

func (c *Client) ListBatchProcessPhases(ctx context.Context, batchID int64) ([]BatchProcessPhase, error) {
	rows, err := c.db.Query(ctx, `
		SELECT bpp.id, bpp.uuid, bpp.batch_id, b.uuid, bpp.process_phase, bpp.phase_at,
		       bpp.created_at, bpp.updated_at, bpp.deleted_at
		FROM batch_process_phase bpp
		JOIN batch b ON b.id = bpp.batch_id
		WHERE bpp.batch_id = $1 AND bpp.deleted_at IS NULL
		ORDER BY bpp.phase_at ASC`,
		batchID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing batch process phases: %w", err)
	}
	defer rows.Close()

	var phases []BatchProcessPhase
	for rows.Next() {
		var phase BatchProcessPhase
		if err := rows.Scan(
			&phase.ID,
			&phase.UUID,
			&phase.BatchID,
			&phase.BatchUUID,
			&phase.ProcessPhase,
			&phase.PhaseAt,
			&phase.CreatedAt,
			&phase.UpdatedAt,
			&phase.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning batch process phase: %w", err)
		}
		phases = append(phases, phase)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing batch process phases: %w", err)
	}

	return phases, nil
}

func (c *Client) ListBatchProcessPhasesByBatchUUID(ctx context.Context, batchUUID string) ([]BatchProcessPhase, error) {
	batch, err := c.GetBatchByUUID(ctx, batchUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving batch uuid: %w", err)
	}

	return c.ListBatchProcessPhases(ctx, batch.ID)
}
