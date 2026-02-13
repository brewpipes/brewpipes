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

	// Resolve recipe UUID and name if recipe_id is set
	if batch.RecipeID != nil {
		var recipeUUID string
		var recipeName string
		err := c.db.QueryRow(ctx, `SELECT uuid, name FROM recipe WHERE id = $1`, *batch.RecipeID).Scan(&recipeUUID, &recipeName)
		if err == nil {
			batch.RecipeUUID = &recipeUUID
			batch.RecipeName = &recipeName
		}
	}

	return batch, nil
}

func (c *Client) GetBatch(ctx context.Context, id int64) (Batch, error) {
	var batch Batch
	err := c.db.QueryRow(ctx, `
		SELECT b.id, b.uuid, b.short_name, b.brew_date, b.notes, b.recipe_id, r.uuid, r.name,
		       latest_phase.process_phase,
		       b.created_at, b.updated_at, b.deleted_at
		FROM batch b
		LEFT JOIN recipe r ON r.id = b.recipe_id
		LEFT JOIN LATERAL (
			SELECT bpp.process_phase
			FROM batch_process_phase bpp
			WHERE bpp.batch_id = b.id AND bpp.deleted_at IS NULL
			ORDER BY bpp.phase_at DESC
			LIMIT 1
		) latest_phase ON true
		WHERE b.id = $1 AND b.deleted_at IS NULL`,
		id,
	).Scan(
		&batch.ID,
		&batch.UUID,
		&batch.ShortName,
		&batch.BrewDate,
		&batch.Notes,
		&batch.RecipeID,
		&batch.RecipeUUID,
		&batch.RecipeName,
		&batch.CurrentPhase,
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

func (c *Client) GetBatchByUUID(ctx context.Context, batchUUID string) (Batch, error) {
	var batch Batch
	err := c.db.QueryRow(ctx, `
		SELECT b.id, b.uuid, b.short_name, b.brew_date, b.notes, b.recipe_id, r.uuid, r.name,
		       latest_phase.process_phase,
		       b.created_at, b.updated_at, b.deleted_at
		FROM batch b
		LEFT JOIN recipe r ON r.id = b.recipe_id
		LEFT JOIN LATERAL (
			SELECT bpp.process_phase
			FROM batch_process_phase bpp
			WHERE bpp.batch_id = b.id AND bpp.deleted_at IS NULL
			ORDER BY bpp.phase_at DESC
			LIMIT 1
		) latest_phase ON true
		WHERE b.uuid = $1 AND b.deleted_at IS NULL`,
		batchUUID,
	).Scan(
		&batch.ID,
		&batch.UUID,
		&batch.ShortName,
		&batch.BrewDate,
		&batch.Notes,
		&batch.RecipeID,
		&batch.RecipeUUID,
		&batch.RecipeName,
		&batch.CurrentPhase,
		&batch.CreatedAt,
		&batch.UpdatedAt,
		&batch.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Batch{}, service.ErrNotFound
		}
		return Batch{}, fmt.Errorf("getting batch by uuid: %w", err)
	}

	return batch, nil
}

func (c *Client) CountBatchesByRecipe(ctx context.Context, recipeUUID string) (int, error) {
	var count int
	err := c.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM batch b
		JOIN recipe r ON b.recipe_id = r.id
		WHERE r.uuid = $1 AND b.deleted_at IS NULL`,
		recipeUUID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting batches by recipe: %w", err)
	}

	return count, nil
}

func (c *Client) ListBatches(ctx context.Context) ([]Batch, error) {
	rows, err := c.db.Query(ctx, `
		SELECT b.id, b.uuid, b.short_name, b.brew_date, b.notes, b.recipe_id, r.uuid, r.name,
		       latest_phase.process_phase,
		       b.created_at, b.updated_at, b.deleted_at
		FROM batch b
		LEFT JOIN recipe r ON r.id = b.recipe_id
		LEFT JOIN LATERAL (
			SELECT bpp.process_phase
			FROM batch_process_phase bpp
			WHERE bpp.batch_id = b.id AND bpp.deleted_at IS NULL
			ORDER BY bpp.phase_at DESC
			LIMIT 1
		) latest_phase ON true
		WHERE b.deleted_at IS NULL
		ORDER BY b.created_at DESC`,
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
			&batch.RecipeUUID,
			&batch.RecipeName,
			&batch.CurrentPhase,
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

	// Resolve recipe UUID and name if recipe_id is set
	if batch.RecipeID != nil {
		var recipeUUID string
		var recipeName string
		err := c.db.QueryRow(ctx, `SELECT uuid, name FROM recipe WHERE id = $1`, *batch.RecipeID).Scan(&recipeUUID, &recipeName)
		if err == nil {
			batch.RecipeUUID = &recipeUUID
			batch.RecipeName = &recipeName
		}
	}

	return batch, nil
}

func (c *Client) UpdateBatchByUUID(ctx context.Context, batchUUID string, batch Batch) (Batch, error) {
	err := c.db.QueryRow(ctx, `
		UPDATE batch
		SET short_name = $1, brew_date = $2, notes = $3, recipe_id = $4, updated_at = timezone('utc', now())
		WHERE uuid = $5 AND deleted_at IS NULL
		RETURNING id, uuid, short_name, brew_date, notes, recipe_id, created_at, updated_at, deleted_at`,
		batch.ShortName,
		batch.BrewDate,
		batch.Notes,
		batch.RecipeID,
		batchUUID,
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
		return Batch{}, fmt.Errorf("updating batch by uuid: %w", err)
	}

	// Resolve recipe UUID and name if recipe_id is set
	if batch.RecipeID != nil {
		var recipeUUID string
		var recipeName string
		err := c.db.QueryRow(ctx, `SELECT uuid, name FROM recipe WHERE id = $1`, *batch.RecipeID).Scan(&recipeUUID, &recipeName)
		if err == nil {
			batch.RecipeUUID = &recipeUUID
			batch.RecipeName = &recipeName
		}
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

// DeleteBatchByUUID soft-deletes a batch by UUID and cascades the soft-delete to all related records.
func (c *Client) DeleteBatchByUUID(ctx context.Context, batchUUID string) error {
	// Resolve UUID to internal ID for cascade operations
	var id int64
	err := c.db.QueryRow(ctx, `SELECT id FROM batch WHERE uuid = $1 AND deleted_at IS NULL`, batchUUID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return service.ErrNotFound
		}
		return fmt.Errorf("resolving batch uuid: %w", err)
	}

	return c.DeleteBatch(ctx, id)
}

// GetBatchDependenciesByUUID returns dependency counts for a batch identified by UUID.
func (c *Client) GetBatchDependenciesByUUID(ctx context.Context, batchUUID string) (BatchDependencies, error) {
	var id int64
	err := c.db.QueryRow(ctx, `SELECT id FROM batch WHERE uuid = $1 AND deleted_at IS NULL`, batchUUID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchDependencies{}, service.ErrNotFound
		}
		return BatchDependencies{}, fmt.Errorf("resolving batch uuid: %w", err)
	}

	return c.GetBatchDependencies(ctx, id)
}
