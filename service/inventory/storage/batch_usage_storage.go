package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// BatchUsagePick describes a single ingredient pick for batch usage deduction.
type BatchUsagePick struct {
	IngredientLotUUID string
	StockLocationUUID string
	Amount            int64
	AmountUnit        string
}

// BatchUsageRequest describes the full batch usage deduction request.
type BatchUsageRequest struct {
	ProductionRefUUID *uuid.UUID
	UsedAt            time.Time
	Picks             []BatchUsagePick
	Notes             *string
}

// BatchUsageResult holds the created usage and movements.
type BatchUsageResult struct {
	Usage     InventoryUsage
	Movements []InventoryMovement
}

// ErrBatchUsageValidation is returned when a pick fails validation (lot not found,
// location not found, or insufficient stock). The message is safe to return to clients.
type ErrBatchUsageValidation struct {
	Message string
}

func (e *ErrBatchUsageValidation) Error() string {
	return e.Message
}

// CreateBatchUsage atomically creates an InventoryUsage record and one
// InventoryMovement per pick within a single transaction. It validates that
// each lot and location exist and that sufficient stock is available.
func (c *Client) CreateBatchUsage(ctx context.Context, req BatchUsageRequest) (BatchUsageResult, error) {
	tx, err := c.DB().Begin(ctx)
	if err != nil {
		return BatchUsageResult{}, fmt.Errorf("starting batch usage transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Resolve and validate each pick before any writes.
	type resolvedPick struct {
		lotID      int64
		lotCode    string
		locationID int64
		locName    string
		amount     int64
		amountUnit string
	}
	resolved := make([]resolvedPick, len(req.Picks))

	for i, pick := range req.Picks {
		// Verify ingredient lot exists and is not deleted.
		var lotID int64
		var lotCode *string
		err := tx.QueryRow(ctx, `
			SELECT id, brewery_lot_code
			FROM ingredient_lot
			WHERE uuid = $1 AND deleted_at IS NULL`,
			pick.IngredientLotUUID,
		).Scan(&lotID, &lotCode)
		if err != nil {
			return BatchUsageResult{}, &ErrBatchUsageValidation{
				Message: fmt.Sprintf("picks[%d]: ingredient lot %s not found", i, pick.IngredientLotUUID),
			}
		}

		lotLabel := pick.IngredientLotUUID
		if lotCode != nil && *lotCode != "" {
			lotLabel = *lotCode
		}

		// Verify stock location exists and is not deleted.
		var locationID int64
		var locName string
		err = tx.QueryRow(ctx, `
			SELECT id, name
			FROM stock_location
			WHERE uuid = $1 AND deleted_at IS NULL`,
			pick.StockLocationUUID,
		).Scan(&locationID, &locName)
		if err != nil {
			return BatchUsageResult{}, &ErrBatchUsageValidation{
				Message: fmt.Sprintf("picks[%d]: stock location %s not found", i, pick.StockLocationUUID),
			}
		}

		// Check available stock for this lot at this location.
		// NOTE: This read is not protected by row-level locks (e.g. SELECT FOR
		// UPDATE), so concurrent transactions could both pass the stock check and
		// overdraw inventory. Acceptable for V1 (single-user), but a future
		// multi-user version should add advisory locking or serializable isolation.
		var available int64
		err = tx.QueryRow(ctx, `
			SELECT COALESCE(SUM(CASE direction
				WHEN 'in' THEN amount
				WHEN 'out' THEN -amount
			END), 0)
			FROM inventory_movement
			WHERE ingredient_lot_id = $1
			  AND stock_location_id = $2
			  AND deleted_at IS NULL`,
			lotID,
			locationID,
		).Scan(&available)
		if err != nil {
			return BatchUsageResult{}, fmt.Errorf("checking stock for pick %d: %w", i, err)
		}

		if pick.Amount > available {
			return BatchUsageResult{}, &ErrBatchUsageValidation{
				Message: fmt.Sprintf("insufficient stock for lot %s at %s: available %d %s, requested %d %s",
					lotLabel, locName, available, pick.AmountUnit, pick.Amount, pick.AmountUnit),
			}
		}

		resolved[i] = resolvedPick{
			lotID:      lotID,
			lotCode:    lotLabel,
			locationID: locationID,
			locName:    locName,
			amount:     pick.Amount,
			amountUnit: pick.AmountUnit,
		}
	}

	// Create the usage record.
	var usage InventoryUsage
	var productionUUID pgtype.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_usage (
			production_ref_uuid,
			used_at,
			notes
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, production_ref_uuid, used_at, notes, created_at, updated_at, deleted_at`,
		database.UUIDParam(req.ProductionRefUUID),
		req.UsedAt,
		req.Notes,
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
		return BatchUsageResult{}, fmt.Errorf("creating inventory usage: %w", err)
	}
	database.AssignUUIDPointer(&usage.ProductionRefUUID, productionUUID)

	// Create one movement per pick.
	movements := make([]InventoryMovement, len(resolved))
	for i, rp := range resolved {
		var m InventoryMovement
		err = tx.QueryRow(ctx, `
			INSERT INTO inventory_movement (
				ingredient_lot_id,
				stock_location_id,
				direction,
				reason,
				amount,
				amount_unit,
				occurred_at,
				usage_id,
				notes
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id, uuid, ingredient_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, usage_id, notes, created_at, updated_at, deleted_at`,
			rp.lotID,
			rp.locationID,
			MovementDirectionOut,
			MovementReasonUse,
			rp.amount,
			rp.amountUnit,
			req.UsedAt,
			usage.ID,
			req.Notes,
		).Scan(
			&m.ID,
			&m.UUID,
			&m.IngredientLotID,
			&m.StockLocationID,
			&m.Direction,
			&m.Reason,
			&m.Amount,
			&m.AmountUnit,
			&m.OccurredAt,
			&m.UsageID,
			&m.Notes,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
		)
		if err != nil {
			return BatchUsageResult{}, fmt.Errorf("creating movement for pick %d: %w", i, err)
		}

		// Set the UUID fields that the JOIN-based queries normally resolve.
		lotUUID := req.Picks[i].IngredientLotUUID
		m.IngredientLotUUID = &lotUUID
		locUUID := req.Picks[i].StockLocationUUID
		m.StockLocationUUID = locUUID
		usageUUIDStr := usage.UUID.String()
		m.UsageUUID = &usageUUIDStr

		movements[i] = m
	}

	if err := tx.Commit(ctx); err != nil {
		return BatchUsageResult{}, fmt.Errorf("committing batch usage transaction: %w", err)
	}

	return BatchUsageResult{
		Usage:     usage,
		Movements: movements,
	}, nil
}
