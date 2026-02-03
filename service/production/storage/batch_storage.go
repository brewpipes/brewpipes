package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateBatch(ctx context.Context, batch Batch) (Batch, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO batch (
			short_name,
			brew_date,
			notes,
			recipe_id
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, short_name, brew_date, notes, recipe_id, created_at, updated_at, deleted_at`,
		batch.ShortName,
		batch.BrewDate,
		batch.Notes,
		batch.RecipeID,
	).Scan(
		&batch.ID,
		&batch.UUID,
		&batch.ShortName,
		&batch.BrewDate,
		&batch.Notes,
		&batch.RecipeID,
		&batch.CreatedAt,
		&batch.UpdatedAt,
		&batch.DeletedAt,
	)
	if err != nil {
		return Batch{}, fmt.Errorf("creating batch: %w", err)
	}

	return batch, nil
}

func (c *Client) GetBatch(ctx context.Context, id int64) (Batch, error) {
	var batch Batch
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, short_name, brew_date, notes, recipe_id, created_at, updated_at, deleted_at
		FROM batch
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&batch.ID,
		&batch.UUID,
		&batch.ShortName,
		&batch.BrewDate,
		&batch.Notes,
		&batch.RecipeID,
		&batch.CreatedAt,
		&batch.UpdatedAt,
		&batch.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Batch{}, service.ErrNotFound
		}
		return Batch{}, fmt.Errorf("getting batch: %w", err)
	}

	return batch, nil
}

func (c *Client) CountBatchesByRecipe(ctx context.Context, recipeID int64) (int, error) {
	var count int
	err := c.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM batch
		WHERE recipe_id = $1 AND deleted_at IS NULL`,
		recipeID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting batches by recipe: %w", err)
	}

	return count, nil
}

func (c *Client) ListBatches(ctx context.Context) ([]Batch, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, short_name, brew_date, notes, recipe_id, created_at, updated_at, deleted_at
		FROM batch
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing batches: %w", err)
	}
	defer rows.Close()

	var batches []Batch
	for rows.Next() {
		var batch Batch
		if err := rows.Scan(
			&batch.ID,
			&batch.UUID,
			&batch.ShortName,
			&batch.BrewDate,
			&batch.Notes,
			&batch.RecipeID,
			&batch.CreatedAt,
			&batch.UpdatedAt,
			&batch.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning batch: %w", err)
		}
		batches = append(batches, batch)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing batches: %w", err)
	}

	return batches, nil
}

func (c *Client) UpdateBatch(ctx context.Context, id int64, batch Batch) (Batch, error) {
	err := c.db.QueryRow(ctx, `
		UPDATE batch
		SET short_name = $1, brew_date = $2, notes = $3, recipe_id = $4, updated_at = timezone('utc', now())
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id, uuid, short_name, brew_date, notes, recipe_id, created_at, updated_at, deleted_at`,
		batch.ShortName,
		batch.BrewDate,
		batch.Notes,
		batch.RecipeID,
		id,
	).Scan(
		&batch.ID,
		&batch.UUID,
		&batch.ShortName,
		&batch.BrewDate,
		&batch.Notes,
		&batch.RecipeID,
		&batch.CreatedAt,
		&batch.UpdatedAt,
		&batch.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Batch{}, service.ErrNotFound
		}
		return Batch{}, fmt.Errorf("updating batch: %w", err)
	}

	return batch, nil
}

// BatchDependencies holds counts of records that depend on a batch.
type BatchDependencies struct {
	BatchVolumeCount       int
	BatchProcessPhaseCount int
	BrewSessionCount       int
	AdditionCount          int
	MeasurementCount       int
}

// HasDependencies returns true if any dependency count is greater than zero.
func (d BatchDependencies) HasDependencies() bool {
	return d.BatchVolumeCount > 0 ||
		d.BatchProcessPhaseCount > 0 ||
		d.BrewSessionCount > 0 ||
		d.AdditionCount > 0 ||
		d.MeasurementCount > 0
}

func (c *Client) GetBatchDependencies(ctx context.Context, id int64) (BatchDependencies, error) {
	var deps BatchDependencies
	err := c.db.QueryRow(ctx, `
		SELECT
			(SELECT COUNT(*) FROM batch_volume WHERE batch_id = $1 AND deleted_at IS NULL),
			(SELECT COUNT(*) FROM batch_process_phase WHERE batch_id = $1 AND deleted_at IS NULL),
			(SELECT COUNT(*) FROM brew_session WHERE batch_id = $1 AND deleted_at IS NULL),
			(SELECT COUNT(*) FROM addition WHERE batch_id = $1 AND deleted_at IS NULL),
			(SELECT COUNT(*) FROM measurement WHERE batch_id = $1 AND deleted_at IS NULL)`,
		id,
	).Scan(
		&deps.BatchVolumeCount,
		&deps.BatchProcessPhaseCount,
		&deps.BrewSessionCount,
		&deps.AdditionCount,
		&deps.MeasurementCount,
	)
	if err != nil {
		return BatchDependencies{}, fmt.Errorf("getting batch dependencies: %w", err)
	}

	return deps, nil
}

// DeleteBatch soft-deletes a batch and cascades the soft-delete to all related records.
func (c *Client) DeleteBatch(ctx context.Context, id int64) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Check if batch exists and is not already deleted
	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM batch WHERE id = $1 AND deleted_at IS NULL)`, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("checking batch existence: %w", err)
	}
	if !exists {
		return service.ErrNotFound
	}

	// Cascade soft-delete to related records
	// Order matters for foreign key relationships, but since we're soft-deleting, order is less critical

	// Soft-delete measurements linked to this batch
	_, err = tx.Exec(ctx, `
		UPDATE measurement
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE batch_id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting measurements: %w", err)
	}

	// Soft-delete additions linked to this batch
	_, err = tx.Exec(ctx, `
		UPDATE addition
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE batch_id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting additions: %w", err)
	}

	// Soft-delete brew sessions linked to this batch
	_, err = tx.Exec(ctx, `
		UPDATE brew_session
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE batch_id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting brew sessions: %w", err)
	}

	// Soft-delete batch process phases linked to this batch
	_, err = tx.Exec(ctx, `
		UPDATE batch_process_phase
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE batch_id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting batch process phases: %w", err)
	}

	// Soft-delete batch volumes linked to this batch
	_, err = tx.Exec(ctx, `
		UPDATE batch_volume
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE batch_id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting batch volumes: %w", err)
	}

	// Finally, soft-delete the batch itself
	_, err = tx.Exec(ctx, `
		UPDATE batch
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
