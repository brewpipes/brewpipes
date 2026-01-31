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

	return movement, nil
}

func (c *Client) GetInventoryMovement(ctx context.Context, id int64) (InventoryMovement, error) {
	var movement InventoryMovement
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, ingredient_lot_id, beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, receipt_id, usage_id, adjustment_id, transfer_id, notes, created_at, updated_at, deleted_at
		FROM inventory_movement
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryMovement{}, service.ErrNotFound
		}
		return InventoryMovement{}, fmt.Errorf("getting inventory movement: %w", err)
	}

	return movement, nil
}

func (c *Client) ListInventoryMovements(ctx context.Context) ([]InventoryMovement, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, ingredient_lot_id, beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, receipt_id, usage_id, adjustment_id, transfer_id, notes, created_at, updated_at, deleted_at
		FROM inventory_movement
		WHERE deleted_at IS NULL
		ORDER BY occurred_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory movements: %w", err)
	}
	defer rows.Close()

	var movements []InventoryMovement
	for rows.Next() {
		var movement InventoryMovement
		if err := rows.Scan(
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

func (c *Client) ListInventoryMovementsByIngredientLot(ctx context.Context, lotID int64) ([]InventoryMovement, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, ingredient_lot_id, beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, receipt_id, usage_id, adjustment_id, transfer_id, notes, created_at, updated_at, deleted_at
		FROM inventory_movement
		WHERE ingredient_lot_id = $1 AND deleted_at IS NULL
		ORDER BY occurred_at ASC`,
		lotID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory movements by lot: %w", err)
	}
	defer rows.Close()

	var movements []InventoryMovement
	for rows.Next() {
		var movement InventoryMovement
		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("scanning inventory movement: %w", err)
		}
		movements = append(movements, movement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing inventory movements by lot: %w", err)
	}

	return movements, nil
}

func (c *Client) ListInventoryMovementsByBeerLot(ctx context.Context, lotID int64) ([]InventoryMovement, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, ingredient_lot_id, beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, receipt_id, usage_id, adjustment_id, transfer_id, notes, created_at, updated_at, deleted_at
		FROM inventory_movement
		WHERE beer_lot_id = $1 AND deleted_at IS NULL
		ORDER BY occurred_at ASC`,
		lotID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory movements by beer lot: %w", err)
	}
	defer rows.Close()

	var movements []InventoryMovement
	for rows.Next() {
		var movement InventoryMovement
		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("scanning inventory movement: %w", err)
		}
		movements = append(movements, movement)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing inventory movements by beer lot: %w", err)
	}

	return movements, nil
}
