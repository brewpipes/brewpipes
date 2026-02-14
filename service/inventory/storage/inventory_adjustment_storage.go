package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateInventoryAdjustment(ctx context.Context, adjustment InventoryAdjustment) (InventoryAdjustment, error) {
	adjustedAt := adjustment.AdjustedAt
	if adjustedAt.IsZero() {
		adjustedAt = time.Now().UTC()
	}

	err := c.DB().QueryRow(ctx, `
		INSERT INTO inventory_adjustment (
			reason,
			adjusted_at,
			notes
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, reason, adjusted_at, notes, created_at, updated_at, deleted_at`,
		adjustment.Reason,
		adjustedAt,
		adjustment.Notes,
	).Scan(
		&adjustment.ID,
		&adjustment.UUID,
		&adjustment.Reason,
		&adjustment.AdjustedAt,
		&adjustment.Notes,
		&adjustment.CreatedAt,
		&adjustment.UpdatedAt,
		&adjustment.DeletedAt,
	)
	if err != nil {
		return InventoryAdjustment{}, fmt.Errorf("creating inventory adjustment: %w", err)
	}

	return adjustment, nil
}

func (c *Client) GetInventoryAdjustment(ctx context.Context, id int64) (InventoryAdjustment, error) {
	var adjustment InventoryAdjustment
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, reason, adjusted_at, notes, created_at, updated_at, deleted_at
		FROM inventory_adjustment
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&adjustment.ID,
		&adjustment.UUID,
		&adjustment.Reason,
		&adjustment.AdjustedAt,
		&adjustment.Notes,
		&adjustment.CreatedAt,
		&adjustment.UpdatedAt,
		&adjustment.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryAdjustment{}, service.ErrNotFound
		}
		return InventoryAdjustment{}, fmt.Errorf("getting inventory adjustment: %w", err)
	}

	return adjustment, nil
}

func (c *Client) GetInventoryAdjustmentByUUID(ctx context.Context, adjustmentUUID string) (InventoryAdjustment, error) {
	var adjustment InventoryAdjustment
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, reason, adjusted_at, notes, created_at, updated_at, deleted_at
		FROM inventory_adjustment
		WHERE uuid = $1 AND deleted_at IS NULL`,
		adjustmentUUID,
	).Scan(
		&adjustment.ID,
		&adjustment.UUID,
		&adjustment.Reason,
		&adjustment.AdjustedAt,
		&adjustment.Notes,
		&adjustment.CreatedAt,
		&adjustment.UpdatedAt,
		&adjustment.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryAdjustment{}, service.ErrNotFound
		}
		return InventoryAdjustment{}, fmt.Errorf("getting inventory adjustment by uuid: %w", err)
	}

	return adjustment, nil
}

func (c *Client) ListInventoryAdjustments(ctx context.Context) ([]InventoryAdjustment, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, reason, adjusted_at, notes, created_at, updated_at, deleted_at
		FROM inventory_adjustment
		WHERE deleted_at IS NULL
		ORDER BY adjusted_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory adjustments: %w", err)
	}
	defer rows.Close()

	var adjustments []InventoryAdjustment
	for rows.Next() {
		var adjustment InventoryAdjustment
		if err := rows.Scan(
			&adjustment.ID,
			&adjustment.UUID,
			&adjustment.Reason,
			&adjustment.AdjustedAt,
			&adjustment.Notes,
			&adjustment.CreatedAt,
			&adjustment.UpdatedAt,
			&adjustment.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning inventory adjustment: %w", err)
		}
		adjustments = append(adjustments, adjustment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing inventory adjustments: %w", err)
	}

	return adjustments, nil
}
