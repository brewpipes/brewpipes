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

	return phase, nil
}

func (c *Client) GetBatchProcessPhase(ctx context.Context, id int64) (BatchProcessPhase, error) {
	var phase BatchProcessPhase
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, batch_id, process_phase, phase_at, created_at, updated_at, deleted_at
		FROM batch_process_phase
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchProcessPhase{}, service.ErrNotFound
		}
		return BatchProcessPhase{}, fmt.Errorf("getting batch process phase: %w", err)
	}

	return phase, nil
}

func (c *Client) ListBatchProcessPhases(ctx context.Context, batchID int64) ([]BatchProcessPhase, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, batch_id, process_phase, phase_at, created_at, updated_at, deleted_at
		FROM batch_process_phase
		WHERE batch_id = $1 AND deleted_at IS NULL
		ORDER BY phase_at ASC`,
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
