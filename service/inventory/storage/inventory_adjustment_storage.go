package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

// AdjustmentWithMovementRequest describes the parameters for creating an
// adjustment record and its corresponding inventory movement atomically.
type AdjustmentWithMovementRequest struct {
	IngredientLotID *int64
	BeerLotID       *int64
	StockLocationID int64
	Amount          int64 // positive = increase, negative = decrease
	AmountUnit      string
	Reason          string
	AdjustedAt      time.Time
	Notes           *string

	// UUID fields for populating the response without extra lookups.
	IngredientLotUUID *string
	BeerLotUUID       *string
	StockLocationUUID string
}

// AdjustmentWithMovementResult holds the created adjustment and movement.
type AdjustmentWithMovementResult struct {
	Adjustment InventoryAdjustment
	Movement   InventoryMovement
}

// CreateInventoryAdjustmentWithMovement atomically creates an inventory
// adjustment record and a corresponding inventory movement within a single
// transaction.
func (c *Client) CreateInventoryAdjustmentWithMovement(ctx context.Context, req AdjustmentWithMovementRequest) (AdjustmentWithMovementResult, error) {
	tx, err := c.DB().Begin(ctx)
	if err != nil {
		return AdjustmentWithMovementResult{}, fmt.Errorf("starting adjustment transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	adjustedAt := req.AdjustedAt
	if adjustedAt.IsZero() {
		adjustedAt = time.Now().UTC()
	}

	// Create the adjustment record.
	var adjustment InventoryAdjustment
	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_adjustment (
			reason,
			adjusted_at,
			notes
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, reason, adjusted_at, notes, created_at, updated_at, deleted_at`,
		req.Reason,
		adjustedAt,
		req.Notes,
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
		return AdjustmentWithMovementResult{}, fmt.Errorf("creating inventory adjustment in transaction: %w", err)
	}

	// Derive movement direction and absolute amount from the signed amount.
	direction := MovementDirectionIn
	absAmount := req.Amount
	if req.Amount < 0 {
		direction = MovementDirectionOut
		absAmount = -req.Amount
	}

	// Create the corresponding inventory movement.
	var movement InventoryMovement
	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_movement (
			ingredient_lot_id,
			beer_lot_id,
			stock_location_id,
			direction,
			reason,
			amount,
			amount_unit,
			occurred_at,
			adjustment_id,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, uuid, ingredient_lot_id, beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, adjustment_id, notes, created_at, updated_at, deleted_at`,
		req.IngredientLotID,
		req.BeerLotID,
		req.StockLocationID,
		direction,
		MovementReasonAdjust,
		absAmount,
		req.AmountUnit,
		adjustedAt,
		adjustment.ID,
		req.Notes,
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
		&movement.AdjustmentID,
		&movement.Notes,
		&movement.CreatedAt,
		&movement.UpdatedAt,
		&movement.DeletedAt,
	)
	if err != nil {
		return AdjustmentWithMovementResult{}, fmt.Errorf("creating adjustment movement in transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return AdjustmentWithMovementResult{}, fmt.Errorf("committing adjustment transaction: %w", err)
	}

	// Populate UUID fields for the response without extra lookups.
	movement.IngredientLotUUID = req.IngredientLotUUID
	if req.BeerLotID != nil {
		movement.BeerLotUUID = req.BeerLotUUID
	}
	movement.StockLocationUUID = req.StockLocationUUID
	adjUUID := adjustment.UUID.String()
	movement.AdjustmentUUID = &adjUUID

	return AdjustmentWithMovementResult{
		Adjustment: adjustment,
		Movement:   movement,
	}, nil
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
