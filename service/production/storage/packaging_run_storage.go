package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

// CreatePackagingRunWithLines atomically creates a packaging run and its lines
// in a single transaction. It returns the created run and lines.
func (c *Client) CreatePackagingRunWithLines(ctx context.Context, run PackagingRun, lines []PackagingRunLine) (PackagingRun, []PackagingRunLine, error) {
	tx, err := c.DB().Begin(ctx)
	if err != nil {
		return PackagingRun{}, nil, fmt.Errorf("starting packaging run transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	startedAt := run.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now().UTC()
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO packaging_run (
			batch_id,
			occupancy_id,
			started_at,
			ended_at,
			loss_amount,
			loss_unit,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, batch_id, occupancy_id, started_at, ended_at, loss_amount, loss_unit, notes, created_at, updated_at, deleted_at`,
		run.BatchID,
		run.OccupancyID,
		startedAt,
		run.EndedAt,
		run.LossAmount,
		run.LossUnit,
		run.Notes,
	).Scan(
		&run.ID,
		&run.UUID,
		&run.BatchID,
		&run.OccupancyID,
		&run.StartedAt,
		&run.EndedAt,
		&run.LossAmount,
		&run.LossUnit,
		&run.Notes,
		&run.CreatedAt,
		&run.UpdatedAt,
		&run.DeletedAt,
	)
	if err != nil {
		return PackagingRun{}, nil, fmt.Errorf("creating packaging run: %w", err)
	}

	// Resolve batch and occupancy UUIDs
	tx.QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, run.BatchID).Scan(&run.BatchUUID)
	tx.QueryRow(ctx, `SELECT uuid FROM occupancy WHERE id = $1`, run.OccupancyID).Scan(&run.OccupancyUUID)

	createdLines := make([]PackagingRunLine, 0, len(lines))
	for _, line := range lines {
		line.PackagingRunID = run.ID
		err = tx.QueryRow(ctx, `
			INSERT INTO packaging_run_line (
				packaging_run_id,
				package_format_id,
				quantity
			) VALUES ($1, $2, $3)
			RETURNING id, uuid, packaging_run_id, package_format_id, quantity, created_at, updated_at, deleted_at`,
			line.PackagingRunID,
			line.PackageFormatID,
			line.Quantity,
		).Scan(
			&line.ID,
			&line.UUID,
			&line.PackagingRunID,
			&line.PackageFormatID,
			&line.Quantity,
			&line.CreatedAt,
			&line.UpdatedAt,
			&line.DeletedAt,
		)
		if err != nil {
			return PackagingRun{}, nil, fmt.Errorf("creating packaging run line: %w", err)
		}

		// Resolve joined fields
		line.PackagingRunUUID = run.UUID.String()
		tx.QueryRow(ctx, `
			SELECT uuid, name, volume_per_unit, volume_per_unit_unit
			FROM package_format WHERE id = $1`,
			line.PackageFormatID,
		).Scan(
			&line.PackageFormatUUID,
			&line.PackageFormatName,
			&line.PackageFormatVolumePerUnit,
			&line.PackageFormatVolumePerUnitUnit,
		)

		createdLines = append(createdLines, line)
	}

	if err := tx.Commit(ctx); err != nil {
		return PackagingRun{}, nil, fmt.Errorf("committing packaging run transaction: %w", err)
	}

	return run, createdLines, nil
}

func (c *Client) GetPackagingRunByUUID(ctx context.Context, runUUID string) (PackagingRun, error) {
	var run PackagingRun
	err := c.DB().QueryRow(ctx, `
		SELECT pr.id, pr.uuid, pr.batch_id, b.uuid, pr.occupancy_id, o.uuid,
		       pr.started_at, pr.ended_at, pr.loss_amount, pr.loss_unit, pr.notes,
		       pr.created_at, pr.updated_at, pr.deleted_at
		FROM packaging_run pr
		JOIN batch b ON b.id = pr.batch_id
		JOIN occupancy o ON o.id = pr.occupancy_id
		WHERE pr.uuid = $1 AND pr.deleted_at IS NULL`,
		runUUID,
	).Scan(
		&run.ID,
		&run.UUID,
		&run.BatchID,
		&run.BatchUUID,
		&run.OccupancyID,
		&run.OccupancyUUID,
		&run.StartedAt,
		&run.EndedAt,
		&run.LossAmount,
		&run.LossUnit,
		&run.Notes,
		&run.CreatedAt,
		&run.UpdatedAt,
		&run.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PackagingRun{}, service.ErrNotFound
		}
		return PackagingRun{}, fmt.Errorf("getting packaging run by uuid: %w", err)
	}

	return run, nil
}

func (c *Client) ListPackagingRuns(ctx context.Context) ([]PackagingRun, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT pr.id, pr.uuid, pr.batch_id, b.uuid, pr.occupancy_id, o.uuid,
		       pr.started_at, pr.ended_at, pr.loss_amount, pr.loss_unit, pr.notes,
		       pr.created_at, pr.updated_at, pr.deleted_at
		FROM packaging_run pr
		JOIN batch b ON b.id = pr.batch_id
		JOIN occupancy o ON o.id = pr.occupancy_id
		WHERE pr.deleted_at IS NULL
		ORDER BY pr.started_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing packaging runs: %w", err)
	}
	defer rows.Close()

	var runs []PackagingRun
	for rows.Next() {
		var run PackagingRun
		if err := rows.Scan(
			&run.ID,
			&run.UUID,
			&run.BatchID,
			&run.BatchUUID,
			&run.OccupancyID,
			&run.OccupancyUUID,
			&run.StartedAt,
			&run.EndedAt,
			&run.LossAmount,
			&run.LossUnit,
			&run.Notes,
			&run.CreatedAt,
			&run.UpdatedAt,
			&run.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning packaging run: %w", err)
		}
		runs = append(runs, run)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing packaging runs: %w", err)
	}

	return runs, nil
}

func (c *Client) ListPackagingRunsByBatchUUID(ctx context.Context, batchUUID string) ([]PackagingRun, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT pr.id, pr.uuid, pr.batch_id, b.uuid, pr.occupancy_id, o.uuid,
		       pr.started_at, pr.ended_at, pr.loss_amount, pr.loss_unit, pr.notes,
		       pr.created_at, pr.updated_at, pr.deleted_at
		FROM packaging_run pr
		JOIN batch b ON b.id = pr.batch_id
		JOIN occupancy o ON o.id = pr.occupancy_id
		WHERE b.uuid = $1 AND pr.deleted_at IS NULL
		ORDER BY pr.started_at DESC`,
		batchUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing packaging runs by batch uuid: %w", err)
	}
	defer rows.Close()

	var runs []PackagingRun
	for rows.Next() {
		var run PackagingRun
		if err := rows.Scan(
			&run.ID,
			&run.UUID,
			&run.BatchID,
			&run.BatchUUID,
			&run.OccupancyID,
			&run.OccupancyUUID,
			&run.StartedAt,
			&run.EndedAt,
			&run.LossAmount,
			&run.LossUnit,
			&run.Notes,
			&run.CreatedAt,
			&run.UpdatedAt,
			&run.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning packaging run: %w", err)
		}
		runs = append(runs, run)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing packaging runs by batch uuid: %w", err)
	}

	return runs, nil
}

func (c *Client) DeletePackagingRun(ctx context.Context, id int64) error {
	tx, err := c.DB().Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting packaging run delete transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Check if packaging run exists and is not already deleted
	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM packaging_run WHERE id = $1 AND deleted_at IS NULL)`, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking packaging run existence: %w", err)
	}
	if !exists {
		return service.ErrNotFound
	}

	// Soft-delete packaging run lines first
	_, err = tx.Exec(ctx, `
		UPDATE packaging_run_line
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE packaging_run_id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting packaging run lines: %w", err)
	}

	// Soft-delete the packaging run
	_, err = tx.Exec(ctx, `
		UPDATE packaging_run
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting packaging run: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing packaging run delete transaction: %w", err)
	}

	return nil
}
