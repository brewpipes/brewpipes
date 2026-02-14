package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Client) CreateInventoryUsage(ctx context.Context, usage InventoryUsage) (InventoryUsage, error) {
	usedAt := usage.UsedAt
	if usedAt.IsZero() {
		usedAt = time.Now().UTC()
	}

	var productionUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		INSERT INTO inventory_usage (
			production_ref_uuid,
			used_at,
			notes
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, production_ref_uuid, used_at, notes, created_at, updated_at, deleted_at`,
		database.UUIDParam(usage.ProductionRefUUID),
		usedAt,
		usage.Notes,
	).Scan(
		&usage.ID,
		&usage.UUID,
		&productionUUID,
		&usage.UsedAt,
		&usage.Notes,
		&usage.CreatedAt,
		&usage.UpdatedAt,
		&usage.DeletedAt,
	)
	if err != nil {
		return InventoryUsage{}, fmt.Errorf("creating inventory usage: %w", err)
	}

	database.AssignUUIDPointer(&usage.ProductionRefUUID, productionUUID)
	return usage, nil
}

func (c *Client) GetInventoryUsage(ctx context.Context, id int64) (InventoryUsage, error) {
	var usage InventoryUsage
	var productionUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, production_ref_uuid, used_at, notes, created_at, updated_at, deleted_at
		FROM inventory_usage
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&usage.ID,
		&usage.UUID,
		&productionUUID,
		&usage.UsedAt,
		&usage.Notes,
		&usage.CreatedAt,
		&usage.UpdatedAt,
		&usage.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryUsage{}, service.ErrNotFound
		}
		return InventoryUsage{}, fmt.Errorf("getting inventory usage: %w", err)
	}

	database.AssignUUIDPointer(&usage.ProductionRefUUID, productionUUID)
	return usage, nil
}

func (c *Client) GetInventoryUsageByUUID(ctx context.Context, usageUUID string) (InventoryUsage, error) {
	var usage InventoryUsage
	var productionUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, production_ref_uuid, used_at, notes, created_at, updated_at, deleted_at
		FROM inventory_usage
		WHERE uuid = $1 AND deleted_at IS NULL`,
		usageUUID,
	).Scan(
		&usage.ID,
		&usage.UUID,
		&productionUUID,
		&usage.UsedAt,
		&usage.Notes,
		&usage.CreatedAt,
		&usage.UpdatedAt,
		&usage.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryUsage{}, service.ErrNotFound
		}
		return InventoryUsage{}, fmt.Errorf("getting inventory usage by uuid: %w", err)
	}

	database.AssignUUIDPointer(&usage.ProductionRefUUID, productionUUID)
	return usage, nil
}

func (c *Client) ListInventoryUsage(ctx context.Context) ([]InventoryUsage, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, production_ref_uuid, used_at, notes, created_at, updated_at, deleted_at
		FROM inventory_usage
		WHERE deleted_at IS NULL
		ORDER BY used_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory usage: %w", err)
	}
	defer rows.Close()

	var usageRecords []InventoryUsage
	for rows.Next() {
		var usage InventoryUsage
		var productionUUID pgtype.UUID
		if err := rows.Scan(
			&usage.ID,
			&usage.UUID,
			&productionUUID,
			&usage.UsedAt,
			&usage.Notes,
			&usage.CreatedAt,
			&usage.UpdatedAt,
			&usage.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning inventory usage: %w", err)
		}
		database.AssignUUIDPointer(&usage.ProductionRefUUID, productionUUID)
		usageRecords = append(usageRecords, usage)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing inventory usage: %w", err)
	}

	return usageRecords, nil
}
