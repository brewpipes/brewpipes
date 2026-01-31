package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

type TransferRecord struct {
	SourceOccupancyID int64
	DestVesselID      int64
	VolumeID          int64
	Amount            int64
	AmountUnit        string
	LossAmount        *int64
	LossUnit          *string
	StartedAt         time.Time
	EndedAt           *time.Time
}

func (c *Client) CreateTransfer(ctx context.Context, transfer Transfer) (Transfer, error) {
	startedAt := transfer.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now().UTC()
	}

	err := c.db.QueryRow(ctx, `
		INSERT INTO transfer (
			source_occupancy_id,
			dest_occupancy_id,
			amount,
			amount_unit,
			loss_amount,
			loss_unit,
			started_at,
			ended_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, uuid, source_occupancy_id, dest_occupancy_id, amount, amount_unit, loss_amount, loss_unit, started_at, ended_at, created_at, updated_at, deleted_at`,
		transfer.SourceOccupancyID,
		transfer.DestOccupancyID,
		transfer.Amount,
		transfer.AmountUnit,
		transfer.LossAmount,
		transfer.LossUnit,
		startedAt,
		transfer.EndedAt,
	).Scan(
		&transfer.ID,
		&transfer.UUID,
		&transfer.SourceOccupancyID,
		&transfer.DestOccupancyID,
		&transfer.Amount,
		&transfer.AmountUnit,
		&transfer.LossAmount,
		&transfer.LossUnit,
		&transfer.StartedAt,
		&transfer.EndedAt,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
		&transfer.DeletedAt,
	)
	if err != nil {
		return Transfer{}, fmt.Errorf("creating transfer: %w", err)
	}

	return transfer, nil
}

func (c *Client) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	var transfer Transfer
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, source_occupancy_id, dest_occupancy_id, amount, amount_unit, loss_amount, loss_unit, started_at, ended_at, created_at, updated_at, deleted_at
		FROM transfer
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&transfer.ID,
		&transfer.UUID,
		&transfer.SourceOccupancyID,
		&transfer.DestOccupancyID,
		&transfer.Amount,
		&transfer.AmountUnit,
		&transfer.LossAmount,
		&transfer.LossUnit,
		&transfer.StartedAt,
		&transfer.EndedAt,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
		&transfer.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Transfer{}, service.ErrNotFound
		}
		return Transfer{}, fmt.Errorf("getting transfer: %w", err)
	}

	return transfer, nil
}

func (c *Client) RecordTransfer(ctx context.Context, record TransferRecord) (Transfer, Occupancy, error) {
	startedAt := record.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now().UTC()
	}

	outAt := startedAt
	if record.EndedAt != nil {
		outAt = *record.EndedAt
	}

	tx, err := c.db.Begin(ctx)
	if err != nil {
		return Transfer{}, Occupancy{}, fmt.Errorf("starting transfer transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var updatedID int64
	if err := tx.QueryRow(ctx, `
		UPDATE occupancy
		SET out_at = $1, updated_at = timezone('utc', now())
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id`,
		outAt,
		record.SourceOccupancyID,
	).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Transfer{}, Occupancy{}, service.ErrNotFound
		}
		return Transfer{}, Occupancy{}, fmt.Errorf("updating source occupancy: %w", err)
	}

	var dest Occupancy
	if err := tx.QueryRow(ctx, `
		INSERT INTO occupancy (
			vessel_id,
			volume_id,
			in_at
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, vessel_id, volume_id, in_at, out_at, created_at, updated_at, deleted_at`,
		record.DestVesselID,
		record.VolumeID,
		startedAt,
	).Scan(
		&dest.ID,
		&dest.UUID,
		&dest.VesselID,
		&dest.VolumeID,
		&dest.InAt,
		&dest.OutAt,
		&dest.CreatedAt,
		&dest.UpdatedAt,
		&dest.DeletedAt,
	); err != nil {
		return Transfer{}, Occupancy{}, fmt.Errorf("creating destination occupancy: %w", err)
	}

	var transfer Transfer
	if err := tx.QueryRow(ctx, `
		INSERT INTO transfer (
			source_occupancy_id,
			dest_occupancy_id,
			amount,
			amount_unit,
			loss_amount,
			loss_unit,
			started_at,
			ended_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, uuid, source_occupancy_id, dest_occupancy_id, amount, amount_unit, loss_amount, loss_unit, started_at, ended_at, created_at, updated_at, deleted_at`,
		record.SourceOccupancyID,
		dest.ID,
		record.Amount,
		record.AmountUnit,
		record.LossAmount,
		record.LossUnit,
		startedAt,
		record.EndedAt,
	).Scan(
		&transfer.ID,
		&transfer.UUID,
		&transfer.SourceOccupancyID,
		&transfer.DestOccupancyID,
		&transfer.Amount,
		&transfer.AmountUnit,
		&transfer.LossAmount,
		&transfer.LossUnit,
		&transfer.StartedAt,
		&transfer.EndedAt,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
		&transfer.DeletedAt,
	); err != nil {
		return Transfer{}, Occupancy{}, fmt.Errorf("creating transfer: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Transfer{}, Occupancy{}, fmt.Errorf("committing transfer: %w", err)
	}

	return transfer, dest, nil
}

func (c *Client) ListTransfersByBatch(ctx context.Context, batchID int64) ([]Transfer, error) {
	rows, err := c.db.Query(ctx, `
		SELECT t.id, t.uuid, t.source_occupancy_id, t.dest_occupancy_id, t.amount, t.amount_unit, t.loss_amount, t.loss_unit, t.started_at, t.ended_at, t.created_at, t.updated_at, t.deleted_at
		FROM transfer t
		WHERE t.deleted_at IS NULL
		AND EXISTS (
			SELECT 1
			FROM occupancy o
			JOIN batch_volume bv ON bv.volume_id = o.volume_id
			WHERE o.id = t.source_occupancy_id
			AND o.deleted_at IS NULL
			AND bv.deleted_at IS NULL
			AND bv.batch_id = $1
		)
		ORDER BY t.started_at ASC`,
		batchID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing transfers by batch: %w", err)
	}
	defer rows.Close()

	var transfers []Transfer
	for rows.Next() {
		var transfer Transfer
		if err := rows.Scan(
			&transfer.ID,
			&transfer.UUID,
			&transfer.SourceOccupancyID,
			&transfer.DestOccupancyID,
			&transfer.Amount,
			&transfer.AmountUnit,
			&transfer.LossAmount,
			&transfer.LossUnit,
			&transfer.StartedAt,
			&transfer.EndedAt,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
			&transfer.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning transfer: %w", err)
		}
		transfers = append(transfers, transfer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing transfers by batch: %w", err)
	}

	return transfers, nil
}
