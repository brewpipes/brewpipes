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
