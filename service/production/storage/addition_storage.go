package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Client) CreateAddition(ctx context.Context, addition Addition) (Addition, error) {
	if (addition.BatchID == nil && addition.OccupancyID == nil) || (addition.BatchID != nil && addition.OccupancyID != nil) {
		return Addition{}, fmt.Errorf("addition must reference batch or occupancy")
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
			addition_type,
			stage,
			inventory_lot_uuid,
			amount,
			amount_unit,
			added_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, uuid, batch_id, occupancy_id, addition_type, stage, inventory_lot_uuid, amount, amount_unit, added_at, notes, created_at, updated_at, deleted_at`,
		addition.BatchID,
		addition.OccupancyID,
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
	return addition, nil
}

func (c *Client) ListAdditionsByBatch(ctx context.Context, batchID int64) ([]Addition, error) {
	rows, err := c.db.Query(ctx, `
		SELECT a.id, a.uuid, a.batch_id, a.occupancy_id, a.addition_type, a.stage, a.inventory_lot_uuid, a.amount, a.amount_unit, a.added_at, a.notes, a.created_at, a.updated_at, a.deleted_at
		FROM addition a
		WHERE a.deleted_at IS NULL
		AND (
			a.batch_id = $1 OR
			EXISTS (
				SELECT 1
				FROM occupancy o
				JOIN batch_volume bv ON bv.volume_id = o.volume_id
				WHERE o.id = a.occupancy_id
				AND o.deleted_at IS NULL
				AND bv.deleted_at IS NULL
				AND bv.batch_id = $1
			)
		)
		ORDER BY a.added_at ASC`,
		batchID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing additions by batch: %w", err)
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
			&addition.OccupancyID,
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
			return nil, fmt.Errorf("scanning addition: %w", err)
		}
		assignUUIDPointer(&addition.InventoryLotUUID, inventoryLotUUID)
		additions = append(additions, addition)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing additions by batch: %w", err)
	}

	return additions, nil
}

func (c *Client) ListAdditionsByOccupancy(ctx context.Context, occupancyID int64) ([]Addition, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, batch_id, occupancy_id, addition_type, stage, inventory_lot_uuid, amount, amount_unit, added_at, notes, created_at, updated_at, deleted_at
		FROM addition
		WHERE occupancy_id = $1 AND deleted_at IS NULL
		ORDER BY added_at ASC`,
		occupancyID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing additions by occupancy: %w", err)
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
			&addition.OccupancyID,
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
			return nil, fmt.Errorf("scanning addition: %w", err)
		}
		assignUUIDPointer(&addition.InventoryLotUUID, inventoryLotUUID)
		additions = append(additions, addition)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing additions by occupancy: %w", err)
	}

	return additions, nil
}

func assignUUIDPointer(destination **uuid.UUID, value pgtype.UUID) {
	if value.Valid {
		uuidValue := uuid.UUID(value.Bytes)
		*destination = &uuidValue
		return
	}

	*destination = nil
}
