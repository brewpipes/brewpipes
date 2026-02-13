package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Client) CreateAddition(ctx context.Context, addition Addition) (Addition, error) {
	targetCount := 0
	if addition.BatchID != nil {
		targetCount++
	}
	if addition.OccupancyID != nil {
		targetCount++
	}
	if addition.VolumeID != nil {
		targetCount++
	}
	if targetCount != 1 {
		return Addition{}, fmt.Errorf("addition must reference exactly one of batch, occupancy, or volume")
	}

	addedAt := addition.AddedAt
	if addedAt.IsZero() {
		addedAt = time.Now().UTC()
	}

	var inventoryLot any
	if addition.InventoryLotUUID != nil {
		inventoryLot = *addition.InventoryLotUUID
	}

	var inventoryLotUUID pgtype.UUID
	err := c.db.QueryRow(ctx, `
		INSERT INTO addition (
			batch_id,
			occupancy_id,
			volume_id,
			addition_type,
			stage,
			inventory_lot_uuid,
			amount,
			amount_unit,
			added_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, uuid, batch_id, occupancy_id, volume_id, addition_type, stage, inventory_lot_uuid, amount, amount_unit, added_at, notes, created_at, updated_at, deleted_at`,
		addition.BatchID,
		addition.OccupancyID,
		addition.VolumeID,
		addition.AdditionType,
		addition.Stage,
		inventoryLot,
		addition.Amount,
		addition.AmountUnit,
		addedAt,
		addition.Notes,
	).Scan(
		&addition.ID,
		&addition.UUID,
		&addition.BatchID,
		&addition.OccupancyID,
		&addition.VolumeID,
		&addition.AdditionType,
		&addition.Stage,
		&inventoryLotUUID,
		&addition.Amount,
		&addition.AmountUnit,
		&addition.AddedAt,
		&addition.Notes,
		&addition.CreatedAt,
		&addition.UpdatedAt,
		&addition.DeletedAt,
	)
	if err != nil {
		return Addition{}, fmt.Errorf("creating addition: %w", err)
	}

	assignUUIDPointer(&addition.InventoryLotUUID, inventoryLotUUID)

	// Resolve FK UUIDs
	if addition.BatchID != nil {
		var batchUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, *addition.BatchID).Scan(&batchUUID); err == nil {
			addition.BatchUUID = &batchUUID
		}
	}
	if addition.OccupancyID != nil {
		var occUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM occupancy WHERE id = $1`, *addition.OccupancyID).Scan(&occUUID); err == nil {
			addition.OccupancyUUID = &occUUID
		}
	}
	if addition.VolumeID != nil {
		var volUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM volume WHERE id = $1`, *addition.VolumeID).Scan(&volUUID); err == nil {
			addition.VolumeUUID = &volUUID
		}
	}

	return addition, nil
}

func scanAddition(row pgx.Row) (Addition, error) {
	var addition Addition
	var inventoryLotUUID pgtype.UUID
	err := row.Scan(
		&addition.ID,
		&addition.UUID,
		&addition.BatchID,
		&addition.BatchUUID,
		&addition.OccupancyID,
		&addition.OccupancyUUID,
		&addition.VolumeID,
		&addition.VolumeUUID,
		&addition.AdditionType,
		&addition.Stage,
		&inventoryLotUUID,
		&addition.Amount,
		&addition.AmountUnit,
		&addition.AddedAt,
		&addition.Notes,
		&addition.CreatedAt,
		&addition.UpdatedAt,
		&addition.DeletedAt,
	)
	if err != nil {
		return Addition{}, err
	}
	assignUUIDPointer(&addition.InventoryLotUUID, inventoryLotUUID)
	return addition, nil
}

const additionSelectWithJoins = `
	SELECT a.id, a.uuid, a.batch_id, b.uuid, a.occupancy_id, o.uuid, a.volume_id, v.uuid,
	       a.addition_type, a.stage, a.inventory_lot_uuid, a.amount, a.amount_unit,
	       a.added_at, a.notes, a.created_at, a.updated_at, a.deleted_at
	FROM addition a
	LEFT JOIN batch b ON b.id = a.batch_id
	LEFT JOIN occupancy o ON o.id = a.occupancy_id
	LEFT JOIN volume v ON v.id = a.volume_id`

func (c *Client) GetAddition(ctx context.Context, id int64) (Addition, error) {
	addition, err := scanAddition(c.db.QueryRow(ctx,
		additionSelectWithJoins+` WHERE a.id = $1 AND a.deleted_at IS NULL`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Addition{}, service.ErrNotFound
		}
		return Addition{}, fmt.Errorf("getting addition: %w", err)
	}
	return addition, nil
}

func (c *Client) GetAdditionByUUID(ctx context.Context, additionUUID string) (Addition, error) {
	addition, err := scanAddition(c.db.QueryRow(ctx,
		additionSelectWithJoins+` WHERE a.uuid = $1 AND a.deleted_at IS NULL`, additionUUID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Addition{}, service.ErrNotFound
		}
		return Addition{}, fmt.Errorf("getting addition by uuid: %w", err)
	}
	return addition, nil
}

func (c *Client) listAdditions(ctx context.Context, query string, args ...any) ([]Addition, error) {
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var additions []Addition
	for rows.Next() {
		var addition Addition
		var inventoryLotUUID pgtype.UUID
		if err := rows.Scan(
			&addition.ID,
			&addition.UUID,
			&addition.BatchID,
			&addition.BatchUUID,
			&addition.OccupancyID,
			&addition.OccupancyUUID,
			&addition.VolumeID,
			&addition.VolumeUUID,
			&addition.AdditionType,
			&addition.Stage,
			&inventoryLotUUID,
			&addition.Amount,
			&addition.AmountUnit,
			&addition.AddedAt,
			&addition.Notes,
			&addition.CreatedAt,
			&addition.UpdatedAt,
			&addition.DeletedAt,
		); err != nil {
			return nil, err
		}
		assignUUIDPointer(&addition.InventoryLotUUID, inventoryLotUUID)
		additions = append(additions, addition)
	}
	return additions, rows.Err()
}

func (c *Client) ListAdditionsByBatch(ctx context.Context, batchID int64) ([]Addition, error) {
	additions, err := c.listAdditions(ctx, additionSelectWithJoins+`
		WHERE a.deleted_at IS NULL
		AND (
			a.batch_id = $1 OR
			EXISTS (
				SELECT 1
				FROM occupancy oo
				JOIN batch_volume bv ON bv.volume_id = oo.volume_id
				WHERE oo.id = a.occupancy_id
				AND oo.deleted_at IS NULL
				AND bv.deleted_at IS NULL
				AND bv.batch_id = $1
			) OR
			EXISTS (
				SELECT 1
				FROM batch_volume bv
				WHERE bv.volume_id = a.volume_id
				AND bv.deleted_at IS NULL
				AND bv.batch_id = $1
			)
		)
		ORDER BY a.added_at ASC`, batchID)
	if err != nil {
		return nil, fmt.Errorf("listing additions by batch: %w", err)
	}
	return additions, nil
}

func (c *Client) ListAdditionsByBatchUUID(ctx context.Context, batchUUID string) ([]Addition, error) {
	batch, err := c.GetBatchByUUID(ctx, batchUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving batch uuid: %w", err)
	}
	return c.ListAdditionsByBatch(ctx, batch.ID)
}

func (c *Client) ListAdditionsByOccupancy(ctx context.Context, occupancyID int64) ([]Addition, error) {
	additions, err := c.listAdditions(ctx, additionSelectWithJoins+`
		WHERE a.occupancy_id = $1 AND a.deleted_at IS NULL
		ORDER BY a.added_at ASC`, occupancyID)
	if err != nil {
		return nil, fmt.Errorf("listing additions by occupancy: %w", err)
	}
	return additions, nil
}

func (c *Client) ListAdditionsByOccupancyUUID(ctx context.Context, occupancyUUID string) ([]Addition, error) {
	occ, err := c.GetOccupancyByUUID(ctx, occupancyUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving occupancy uuid: %w", err)
	}
	return c.ListAdditionsByOccupancy(ctx, occ.ID)
}

func (c *Client) ListAdditionsByVolume(ctx context.Context, volumeID int64) ([]Addition, error) {
	additions, err := c.listAdditions(ctx, additionSelectWithJoins+`
		WHERE a.volume_id = $1 AND a.deleted_at IS NULL
		ORDER BY a.added_at ASC`, volumeID)
	if err != nil {
		return nil, fmt.Errorf("listing additions by volume: %w", err)
	}
	return additions, nil
}

func (c *Client) ListAdditionsByVolumeUUID(ctx context.Context, volumeUUID string) ([]Addition, error) {
	vol, err := c.GetVolumeByUUID(ctx, volumeUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving volume uuid: %w", err)
	}
	return c.ListAdditionsByVolume(ctx, vol.ID)
}

func assignUUIDPointer(destination **uuid.UUID, value pgtype.UUID) {
	if value.Valid {
		uuidValue := uuid.UUID(value.Bytes)
		*destination = &uuidValue
		return
	}

	*destination = nil
}
