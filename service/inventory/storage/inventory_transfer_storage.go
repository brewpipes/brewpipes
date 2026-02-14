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

	err := c.DB().QueryRow(ctx, `
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

	// Resolve location UUIDs
	c.resolveTransferLocationUUIDs(ctx, &transfer)
	return transfer, nil
}

func (c *Client) GetInventoryTransfer(ctx context.Context, id int64) (InventoryTransfer, error) {
	var transfer InventoryTransfer
	err := c.DB().QueryRow(ctx, `
		SELECT t.id, t.uuid, t.source_location_id, sl.uuid, t.dest_location_id, dl.uuid,
		       t.transferred_at, t.notes, t.created_at, t.updated_at, t.deleted_at
		FROM inventory_transfer t
		JOIN stock_location sl ON sl.id = t.source_location_id
		JOIN stock_location dl ON dl.id = t.dest_location_id
		WHERE t.id = $1 AND t.deleted_at IS NULL`,
		id,
	).Scan(
		&transfer.ID,
		&transfer.UUID,
		&transfer.SourceLocationID,
		&transfer.SourceLocationUUID,
		&transfer.DestLocationID,
		&transfer.DestLocationUUID,
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

func (c *Client) GetInventoryTransferByUUID(ctx context.Context, transferUUID string) (InventoryTransfer, error) {
	var transfer InventoryTransfer
	err := c.DB().QueryRow(ctx, `
		SELECT t.id, t.uuid, t.source_location_id, sl.uuid, t.dest_location_id, dl.uuid,
		       t.transferred_at, t.notes, t.created_at, t.updated_at, t.deleted_at
		FROM inventory_transfer t
		JOIN stock_location sl ON sl.id = t.source_location_id
		JOIN stock_location dl ON dl.id = t.dest_location_id
		WHERE t.uuid = $1 AND t.deleted_at IS NULL`,
		transferUUID,
	).Scan(
		&transfer.ID,
		&transfer.UUID,
		&transfer.SourceLocationID,
		&transfer.SourceLocationUUID,
		&transfer.DestLocationID,
		&transfer.DestLocationUUID,
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
		return InventoryTransfer{}, fmt.Errorf("getting inventory transfer by uuid: %w", err)
	}

	return transfer, nil
}

func (c *Client) ListInventoryTransfers(ctx context.Context) ([]InventoryTransfer, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT t.id, t.uuid, t.source_location_id, sl.uuid, t.dest_location_id, dl.uuid,
		       t.transferred_at, t.notes, t.created_at, t.updated_at, t.deleted_at
		FROM inventory_transfer t
		JOIN stock_location sl ON sl.id = t.source_location_id
		JOIN stock_location dl ON dl.id = t.dest_location_id
		WHERE t.deleted_at IS NULL
		ORDER BY t.transferred_at DESC`,
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
			&transfer.SourceLocationUUID,
			&transfer.DestLocationID,
			&transfer.DestLocationUUID,
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

// resolveTransferLocationUUIDs resolves source and dest location UUIDs after INSERT.
func (c *Client) resolveTransferLocationUUIDs(ctx context.Context, transfer *InventoryTransfer) {
	var srcUUID, dstUUID string
	err := c.DB().QueryRow(ctx, `SELECT uuid FROM stock_location WHERE id = $1`, transfer.SourceLocationID).Scan(&srcUUID)
	if err == nil {
		transfer.SourceLocationUUID = srcUUID
	}
	err = c.DB().QueryRow(ctx, `SELECT uuid FROM stock_location WHERE id = $1`, transfer.DestLocationID).Scan(&dstUUID)
	if err == nil {
		transfer.DestLocationUUID = dstUUID
	}
}
