package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateInventoryTransfer(ctx context.Context, transfer InventoryTransfer) (InventoryTransfer, error) {
	transferredAt := transfer.TransferredAt
	if transferredAt.IsZero() {
		transferredAt = time.Now().UTC()
	}

	err := c.db.QueryRow(ctx, `
		INSERT INTO inventory_transfer (
			source_location_id,
			dest_location_id,
			transferred_at,
			notes
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, source_location_id, dest_location_id, transferred_at, notes, created_at, updated_at, deleted_at`,
		transfer.SourceLocationID,
		transfer.DestLocationID,
		transferredAt,
		transfer.Notes,
	).Scan(
		&transfer.ID,
		&transfer.UUID,
		&transfer.SourceLocationID,
		&transfer.DestLocationID,
		&transfer.TransferredAt,
		&transfer.Notes,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
		&transfer.DeletedAt,
	)
	if err != nil {
		return InventoryTransfer{}, fmt.Errorf("creating inventory transfer: %w", err)
	}

	return transfer, nil
}

func (c *Client) GetInventoryTransfer(ctx context.Context, id int64) (InventoryTransfer, error) {
	var transfer InventoryTransfer
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, source_location_id, dest_location_id, transferred_at, notes, created_at, updated_at, deleted_at
		FROM inventory_transfer
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&transfer.ID,
		&transfer.UUID,
		&transfer.SourceLocationID,
		&transfer.DestLocationID,
		&transfer.TransferredAt,
		&transfer.Notes,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
		&transfer.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryTransfer{}, service.ErrNotFound
		}
		return InventoryTransfer{}, fmt.Errorf("getting inventory transfer: %w", err)
	}

	return transfer, nil
}

func (c *Client) ListInventoryTransfers(ctx context.Context) ([]InventoryTransfer, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, source_location_id, dest_location_id, transferred_at, notes, created_at, updated_at, deleted_at
		FROM inventory_transfer
		WHERE deleted_at IS NULL
		ORDER BY transferred_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory transfers: %w", err)
	}
	defer rows.Close()

	var transfers []InventoryTransfer
	for rows.Next() {
		var transfer InventoryTransfer
		if err := rows.Scan(
			&transfer.ID,
			&transfer.UUID,
			&transfer.SourceLocationID,
			&transfer.DestLocationID,
			&transfer.TransferredAt,
			&transfer.Notes,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
			&transfer.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning inventory transfer: %w", err)
		}
		transfers = append(transfers, transfer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing inventory transfers: %w", err)
	}

	return transfers, nil
}
