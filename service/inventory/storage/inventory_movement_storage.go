package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateInventoryMovement(ctx context.Context, movement InventoryMovement) (InventoryMovement, error) {
	if (movement.IngredientLotID == nil && movement.BeerLotID == nil) || (movement.IngredientLotID != nil && movement.BeerLotID != nil) {
		return InventoryMovement{}, fmt.Errorf("inventory movement must reference ingredient lot or beer lot")
	}

	referenceCount := 0
	if movement.ReceiptID != nil {
		referenceCount++
	}
	if movement.UsageID != nil {
		referenceCount++
	}
	if movement.AdjustmentID != nil {
		referenceCount++
	}
	if movement.TransferID != nil {
		referenceCount++
	}
	if referenceCount > 1 {
		return InventoryMovement{}, fmt.Errorf("inventory movement must have at most one reference")
	}

	switch movement.Reason {
	case MovementReasonReceive:
		if movement.ReceiptID == nil {
			return InventoryMovement{}, fmt.Errorf("receive movement must reference receipt")
		}
	case MovementReasonUse:
		if movement.UsageID == nil {
			return InventoryMovement{}, fmt.Errorf("use movement must reference usage")
		}
	case MovementReasonTransfer:
		if movement.TransferID == nil {
			return InventoryMovement{}, fmt.Errorf("transfer movement must reference transfer")
		}
	case MovementReasonAdjust, MovementReasonWaste:
		if movement.AdjustmentID == nil {
			return InventoryMovement{}, fmt.Errorf("adjustment movement must reference adjustment")
		}
	}

	occurredAt := movement.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	err := c.db.QueryRow(ctx, `
		INSERT INTO inventory_movement (
			ingredient_lot_id,
			beer_lot_id,
			stock_location_id,
			direction,
			reason,
			amount,
			amount_unit,
			occurred_at,
			receipt_id,
			usage_id,
			adjustment_id,
			transfer_id,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, uuid, ingredient_lot_id, beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, receipt_id, usage_id, adjustment_id, transfer_id, notes, created_at, updated_at, deleted_at`,
		movement.IngredientLotID,
		movement.BeerLotID,
		movement.StockLocationID,
		movement.Direction,
		movement.Reason,
		movement.Amount,
		movement.AmountUnit,
		occurredAt,
		movement.ReceiptID,
		movement.UsageID,
		movement.AdjustmentID,
		movement.TransferID,
		movement.Notes,
	).Scan(
		&movement.ID,
		&movement.UUID,
		&movement.IngredientLotID,
		&movement.BeerLotID,
		&movement.StockLocationID,
		&movement.Direction,
		&movement.Reason,
		&movement.Amount,
		&movement.AmountUnit,
		&movement.OccurredAt,
		&movement.ReceiptID,
		&movement.UsageID,
		&movement.AdjustmentID,
		&movement.TransferID,
		&movement.Notes,
		&movement.CreatedAt,
		&movement.UpdatedAt,
		&movement.DeletedAt,
	)
	if err != nil {
		return InventoryMovement{}, fmt.Errorf("creating inventory movement: %w", err)
	}

	// Resolve FK UUIDs after insert
	c.resolveMovementUUIDs(ctx, &movement)
	return movement, nil
}

func (c *Client) GetInventoryMovement(ctx context.Context, id int64) (InventoryMovement, error) {
	return c.scanMovementRow(c.db.QueryRow(ctx, movementSelectSQL+`
		WHERE m.id = $1 AND m.deleted_at IS NULL`, id))
}

func (c *Client) GetInventoryMovementByUUID(ctx context.Context, movementUUID string) (InventoryMovement, error) {
	return c.scanMovementRow(c.db.QueryRow(ctx, movementSelectSQL+`
		WHERE m.uuid = $1 AND m.deleted_at IS NULL`, movementUUID))
}

func (c *Client) ListInventoryMovements(ctx context.Context) ([]InventoryMovement, error) {
	rows, err := c.db.Query(ctx, movementSelectSQL+`
		WHERE m.deleted_at IS NULL
		ORDER BY m.occurred_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("listing inventory movements: %w", err)
	}
	defer rows.Close()

	return c.scanMovementRows(rows)
}

func (c *Client) ListInventoryMovementsByIngredientLot(ctx context.Context, lotUUID string) ([]InventoryMovement, error) {
	rows, err := c.db.Query(ctx, movementSelectSQL+`
		WHERE il.uuid = $1 AND m.deleted_at IS NULL
		ORDER BY m.occurred_at ASC`,
		lotUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory movements by lot: %w", err)
	}
	defer rows.Close()

	return c.scanMovementRows(rows)
}

func (c *Client) ListInventoryMovementsByBeerLot(ctx context.Context, lotUUID string) ([]InventoryMovement, error) {
	rows, err := c.db.Query(ctx, movementSelectSQL+`
		WHERE bl.uuid = $1 AND m.deleted_at IS NULL
		ORDER BY m.occurred_at ASC`,
		lotUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory movements by beer lot: %w", err)
	}
	defer rows.Close()

	return c.scanMovementRows(rows)
}

const movementSelectSQL = `
	SELECT m.id, m.uuid,
	       m.ingredient_lot_id, il.uuid,
	       m.beer_lot_id, bl.uuid,
	       m.stock_location_id, sl.uuid,
	       m.direction, m.reason, m.amount, m.amount_unit, m.occurred_at,
	       m.receipt_id, rc.uuid,
	       m.usage_id, us.uuid,
	       m.adjustment_id, aj.uuid,
	       m.transfer_id, tr.uuid,
	       m.notes, m.created_at, m.updated_at, m.deleted_at
	FROM inventory_movement m
	LEFT JOIN ingredient_lot il ON il.id = m.ingredient_lot_id
	LEFT JOIN beer_lot bl ON bl.id = m.beer_lot_id
	JOIN stock_location sl ON sl.id = m.stock_location_id
	LEFT JOIN inventory_receipt rc ON rc.id = m.receipt_id
	LEFT JOIN inventory_usage us ON us.id = m.usage_id
	LEFT JOIN inventory_adjustment aj ON aj.id = m.adjustment_id
	LEFT JOIN inventory_transfer tr ON tr.id = m.transfer_id
`

func (c *Client) scanMovementRow(row pgx.Row) (InventoryMovement, error) {
	var movement InventoryMovement
	err := row.Scan(
		&movement.ID,
		&movement.UUID,
		&movement.IngredientLotID,
		&movement.IngredientLotUUID,
		&movement.BeerLotID,
		&movement.BeerLotUUID,
		&movement.StockLocationID,
		&movement.StockLocationUUID,
		&movement.Direction,
		&movement.Reason,
		&movement.Amount,
		&movement.AmountUnit,
		&movement.OccurredAt,
		&movement.ReceiptID,
		&movement.ReceiptUUID,
		&movement.UsageID,
		&movement.UsageUUID,
		&movement.AdjustmentID,
		&movement.AdjustmentUUID,
		&movement.TransferID,
		&movement.TransferUUID,
		&movement.Notes,
		&movement.CreatedAt,
		&movement.UpdatedAt,
		&movement.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryMovement{}, service.ErrNotFound
		}
		return InventoryMovement{}, fmt.Errorf("scanning inventory movement: %w", err)
	}

	return movement, nil
}

func (c *Client) scanMovementRows(rows pgx.Rows) ([]InventoryMovement, error) {
	var movements []InventoryMovement
	for rows.Next() {
		var movement InventoryMovement
		if err := rows.Scan(
			&movement.ID,
			&movement.UUID,
			&movement.IngredientLotID,
			&movement.IngredientLotUUID,
			&movement.BeerLotID,
			&movement.BeerLotUUID,
			&movement.StockLocationID,
			&movement.StockLocationUUID,
			&movement.Direction,
			&movement.Reason,
			&movement.Amount,
			&movement.AmountUnit,
			&movement.OccurredAt,
			&movement.ReceiptID,
			&movement.ReceiptUUID,
			&movement.UsageID,
			&movement.UsageUUID,
			&movement.AdjustmentID,
			&movement.AdjustmentUUID,
			&movement.TransferID,
			&movement.TransferUUID,
			&movement.Notes,
			&movement.CreatedAt,
			&movement.UpdatedAt,
			&movement.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning inventory movement: %w", err)
		}
		movements = append(movements, movement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing inventory movements: %w", err)
	}

	return movements, nil
}

// resolveMovementUUIDs resolves FK UUIDs after an INSERT (which only returns int IDs).
func (c *Client) resolveMovementUUIDs(ctx context.Context, m *InventoryMovement) {
	// Stock location UUID (required)
	var slUUID string
	if err := c.db.QueryRow(ctx, `SELECT uuid FROM stock_location WHERE id = $1`, m.StockLocationID).Scan(&slUUID); err == nil {
		m.StockLocationUUID = slUUID
	}

	// Ingredient lot UUID (optional)
	if m.IngredientLotID != nil {
		var ilUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM ingredient_lot WHERE id = $1`, *m.IngredientLotID).Scan(&ilUUID); err == nil {
			m.IngredientLotUUID = &ilUUID
		}
	}

	// Beer lot UUID (optional)
	if m.BeerLotID != nil {
		var blUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM beer_lot WHERE id = $1`, *m.BeerLotID).Scan(&blUUID); err == nil {
			m.BeerLotUUID = &blUUID
		}
	}

	// Receipt UUID (optional)
	if m.ReceiptID != nil {
		var rcUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM inventory_receipt WHERE id = $1`, *m.ReceiptID).Scan(&rcUUID); err == nil {
			m.ReceiptUUID = &rcUUID
		}
	}

	// Usage UUID (optional)
	if m.UsageID != nil {
		var usUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM inventory_usage WHERE id = $1`, *m.UsageID).Scan(&usUUID); err == nil {
			m.UsageUUID = &usUUID
		}
	}

	// Adjustment UUID (optional)
	if m.AdjustmentID != nil {
		var ajUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM inventory_adjustment WHERE id = $1`, *m.AdjustmentID).Scan(&ajUUID); err == nil {
			m.AdjustmentUUID = &ajUUID
		}
	}

	// Transfer UUID (optional)
	if m.TransferID != nil {
		var trUUID string
		if err := c.db.QueryRow(ctx, `SELECT uuid FROM inventory_transfer WHERE id = $1`, *m.TransferID).Scan(&trUUID); err == nil {
			m.TransferUUID = &trUUID
		}
	}
}
