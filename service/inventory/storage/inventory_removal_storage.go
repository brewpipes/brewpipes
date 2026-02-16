package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

// RemovalWithMovementRequest describes the parameters for creating a removal
// record and an optional inventory movement atomically.
type RemovalWithMovementRequest struct {
	Category      string
	Reason        string
	BatchUUID     *string
	BeerLotID     *int64
	OccupancyUUID *string
	Amount        int64
	AmountUnit    string
	AmountBBL     *float64
	IsTaxable     bool
	ReferenceCode *string
	PerformedBy   *string
	RemovedAt     time.Time
	Destination   *string
	Notes         *string

	// For creating the movement (required when BeerLotID is set).
	StockLocationID *int64

	// UUID fields for populating the response without extra lookups.
	BeerLotUUID       *string
	StockLocationUUID *string
}

// RemovalWithMovementResult holds the created removal and optional movement UUID.
type RemovalWithMovementResult struct {
	Removal      InventoryRemoval
	MovementUUID *string
}

// RemovalListFilter describes optional filters for listing removals.
type RemovalListFilter struct {
	BatchUUID   *string
	BeerLotUUID *string
	Category    *string
	From        *time.Time
	To          *time.Time
}

// RemovalSummary holds aggregated removal data for reporting.
type RemovalSummary struct {
	TotalBBL   float64
	TaxableBBL float64
	TaxFreeBBL float64
	TotalCount int
	ByCategory []RemovalCategorySummary
}

// RemovalCategorySummary holds per-category aggregated removal data.
type RemovalCategorySummary struct {
	Category string
	TotalBBL float64
	Count    int
}

// removalColumns is the column list shared by removal queries.
const removalColumns = `r.id, r.uuid, r.category, r.reason,
	r.batch_uuid, r.beer_lot_id, bl.uuid,
	r.occupancy_uuid, r.amount, r.amount_unit, r.amount_bbl,
	r.is_taxable, r.reference_code, r.performed_by, r.removed_at,
	r.destination, r.notes,
	r.movement_id, mv.uuid,
	sl.uuid,
	r.created_at, r.updated_at, r.deleted_at`

// removalJoins is the JOIN clause shared by removal queries.
const removalJoins = `
	FROM inventory_removal r
	LEFT JOIN beer_lot bl ON bl.id = r.beer_lot_id
	LEFT JOIN inventory_movement mv ON mv.id = r.movement_id
	LEFT JOIN stock_location sl ON sl.id = mv.stock_location_id`

func scanRemoval(row pgx.Row) (InventoryRemoval, error) {
	var removal InventoryRemoval
	err := row.Scan(
		&removal.ID,
		&removal.UUID,
		&removal.Category,
		&removal.Reason,
		&removal.BatchUUID,
		&removal.BeerLotID,
		&removal.BeerLotUUID,
		&removal.OccupancyUUID,
		&removal.Amount,
		&removal.AmountUnit,
		&removal.AmountBBL,
		&removal.IsTaxable,
		&removal.ReferenceCode,
		&removal.PerformedBy,
		&removal.RemovedAt,
		&removal.Destination,
		&removal.Notes,
		&removal.MovementID,
		&removal.MovementUUID,
		&removal.StockLocationUUID,
		&removal.CreatedAt,
		&removal.UpdatedAt,
		&removal.DeletedAt,
	)
	return removal, err
}

func scanRemovalRows(rows pgx.Rows) ([]InventoryRemoval, error) {
	var removals []InventoryRemoval
	for rows.Next() {
		var removal InventoryRemoval
		if err := rows.Scan(
			&removal.ID,
			&removal.UUID,
			&removal.Category,
			&removal.Reason,
			&removal.BatchUUID,
			&removal.BeerLotID,
			&removal.BeerLotUUID,
			&removal.OccupancyUUID,
			&removal.Amount,
			&removal.AmountUnit,
			&removal.AmountBBL,
			&removal.IsTaxable,
			&removal.ReferenceCode,
			&removal.PerformedBy,
			&removal.RemovedAt,
			&removal.Destination,
			&removal.Notes,
			&removal.MovementID,
			&removal.MovementUUID,
			&removal.StockLocationUUID,
			&removal.CreatedAt,
			&removal.UpdatedAt,
			&removal.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning inventory removal: %w", err)
		}
		removals = append(removals, removal)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return removals, nil
}

// CreateRemovalWithMovement atomically creates an inventory removal record and,
// when a beer lot is referenced, a corresponding inventory movement within a
// single transaction.
func (c *Client) CreateRemovalWithMovement(ctx context.Context, req RemovalWithMovementRequest) (RemovalWithMovementResult, error) {
	tx, err := c.DB().Begin(ctx)
	if err != nil {
		return RemovalWithMovementResult{}, fmt.Errorf("starting removal transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	removedAt := req.RemovedAt
	if removedAt.IsZero() {
		removedAt = time.Now().UTC()
	}

	// Insert the removal record.
	var removal InventoryRemoval
	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_removal (
			category, reason, batch_uuid, beer_lot_id, occupancy_uuid,
			amount, amount_unit, amount_bbl, is_taxable,
			reference_code, performed_by, removed_at, destination, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, uuid, category, reason, batch_uuid, beer_lot_id, occupancy_uuid,
			amount, amount_unit, amount_bbl, is_taxable, reference_code, performed_by,
			removed_at, destination, notes, movement_id, created_at, updated_at, deleted_at`,
		req.Category,
		req.Reason,
		req.BatchUUID,
		req.BeerLotID,
		req.OccupancyUUID,
		req.Amount,
		req.AmountUnit,
		req.AmountBBL,
		req.IsTaxable,
		req.ReferenceCode,
		req.PerformedBy,
		removedAt,
		req.Destination,
		req.Notes,
	).Scan(
		&removal.ID,
		&removal.UUID,
		&removal.Category,
		&removal.Reason,
		&removal.BatchUUID,
		&removal.BeerLotID,
		&removal.OccupancyUUID,
		&removal.Amount,
		&removal.AmountUnit,
		&removal.AmountBBL,
		&removal.IsTaxable,
		&removal.ReferenceCode,
		&removal.PerformedBy,
		&removal.RemovedAt,
		&removal.Destination,
		&removal.Notes,
		&removal.MovementID,
		&removal.CreatedAt,
		&removal.UpdatedAt,
		&removal.DeletedAt,
	)
	if err != nil {
		return RemovalWithMovementResult{}, fmt.Errorf("creating inventory removal: %w", err)
	}

	var movementUUID *string

	// If a beer lot is referenced, create the corresponding movement.
	if req.BeerLotID != nil && req.StockLocationID != nil {
		var mvUUID string
		var mvID int64
		err = tx.QueryRow(ctx, `
			INSERT INTO inventory_movement (
				beer_lot_id, stock_location_id, direction, reason,
				amount, amount_unit, occurred_at, removal_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, uuid`,
			req.BeerLotID,
			req.StockLocationID,
			MovementDirectionOut,
			MovementReasonRemoval,
			req.Amount,
			req.AmountUnit,
			removedAt,
			removal.ID,
		).Scan(&mvID, &mvUUID)
		if err != nil {
			return RemovalWithMovementResult{}, fmt.Errorf("creating removal movement: %w", err)
		}

		// Back-reference the movement on the removal.
		_, err = tx.Exec(ctx, `
			UPDATE inventory_removal SET movement_id = $1, updated_at = timezone('utc', now())
			WHERE id = $2`,
			mvID,
			removal.ID,
		)
		if err != nil {
			return RemovalWithMovementResult{}, fmt.Errorf("updating removal movement reference: %w", err)
		}

		removal.MovementID = &mvID
		movementUUID = &mvUUID
		removal.MovementUUID = &mvUUID
		removal.StockLocationUUID = req.StockLocationUUID
	}

	if err := tx.Commit(ctx); err != nil {
		return RemovalWithMovementResult{}, fmt.Errorf("committing removal transaction: %w", err)
	}

	// Populate UUID fields for the response without extra lookups.
	removal.BeerLotUUID = req.BeerLotUUID

	return RemovalWithMovementResult{
		Removal:      removal,
		MovementUUID: movementUUID,
	}, nil
}

// GetRemovalByUUID returns a single removal by UUID.
func (c *Client) GetRemovalByUUID(ctx context.Context, removalUUID string) (InventoryRemoval, error) {
	removal, err := scanRemoval(c.DB().QueryRow(ctx, `
		SELECT `+removalColumns+removalJoins+`
		WHERE r.uuid = $1 AND r.deleted_at IS NULL`,
		removalUUID,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryRemoval{}, service.ErrNotFound
		}
		return InventoryRemoval{}, fmt.Errorf("getting removal by uuid: %w", err)
	}
	return removal, nil
}

// ListRemovals returns removals matching the given filters.
func (c *Client) ListRemovals(ctx context.Context, filter RemovalListFilter) ([]InventoryRemoval, error) {
	query := `SELECT ` + removalColumns + removalJoins + `
		WHERE r.deleted_at IS NULL`
	args := []any{}
	argIdx := 1

	if filter.BatchUUID != nil {
		query += fmt.Sprintf(` AND r.batch_uuid = $%d`, argIdx)
		args = append(args, *filter.BatchUUID)
		argIdx++
	}
	if filter.BeerLotUUID != nil {
		query += fmt.Sprintf(` AND bl.uuid = $%d`, argIdx)
		args = append(args, *filter.BeerLotUUID)
		argIdx++
	}
	if filter.Category != nil {
		query += fmt.Sprintf(` AND r.category = $%d`, argIdx)
		args = append(args, *filter.Category)
		argIdx++
	}
	if filter.From != nil {
		query += fmt.Sprintf(` AND r.removed_at >= $%d`, argIdx)
		args = append(args, *filter.From)
		argIdx++
	}
	if filter.To != nil {
		query += fmt.Sprintf(` AND r.removed_at <= $%d`, argIdx)
		args = append(args, *filter.To)
		argIdx++
	}

	query += ` ORDER BY r.removed_at DESC`

	rows, err := c.DB().Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing removals: %w", err)
	}
	defer rows.Close()

	removals, err := scanRemovalRows(rows)
	if err != nil {
		return nil, fmt.Errorf("listing removals: %w", err)
	}
	return removals, nil
}

// UpdateRemoval updates mutable fields on a removal record.
func (c *Client) UpdateRemoval(ctx context.Context, removalUUID string, req UpdateRemovalRequest) (InventoryRemoval, error) {
	// Fetch the existing removal first.
	existing, err := c.GetRemovalByUUID(ctx, removalUUID)
	if err != nil {
		return InventoryRemoval{}, err
	}

	// Apply PATCH semantics: only update fields that are provided.
	if req.Category != nil {
		existing.Category = *req.Category
	}
	if req.Reason != nil {
		existing.Reason = *req.Reason
	}
	if req.Amount != nil {
		existing.Amount = *req.Amount
	}
	if req.AmountUnit != nil {
		existing.AmountUnit = *req.AmountUnit
	}
	if req.AmountBBL != nil {
		existing.AmountBBL = req.AmountBBL
	}
	if req.IsTaxable != nil {
		existing.IsTaxable = *req.IsTaxable
	}
	if req.ReferenceCode != nil {
		existing.ReferenceCode = req.ReferenceCode
	}
	if req.PerformedBy != nil {
		existing.PerformedBy = req.PerformedBy
	}
	if req.RemovedAt != nil {
		existing.RemovedAt = *req.RemovedAt
	}
	if req.Destination != nil {
		existing.Destination = req.Destination
	}
	if req.Notes != nil {
		existing.Notes = req.Notes
	}

	err = c.DB().QueryRow(ctx, `
		UPDATE inventory_removal SET
			category = $1, reason = $2, amount = $3, amount_unit = $4,
			amount_bbl = $5, is_taxable = $6, reference_code = $7,
			performed_by = $8, removed_at = $9, destination = $10,
			notes = $11, updated_at = timezone('utc', now())
		WHERE uuid = $12 AND deleted_at IS NULL
		RETURNING id, uuid, category, reason, batch_uuid, beer_lot_id, occupancy_uuid,
			amount, amount_unit, amount_bbl, is_taxable, reference_code, performed_by,
			removed_at, destination, notes, movement_id, created_at, updated_at, deleted_at`,
		existing.Category,
		existing.Reason,
		existing.Amount,
		existing.AmountUnit,
		existing.AmountBBL,
		existing.IsTaxable,
		existing.ReferenceCode,
		existing.PerformedBy,
		existing.RemovedAt,
		existing.Destination,
		existing.Notes,
		removalUUID,
	).Scan(
		&existing.ID,
		&existing.UUID,
		&existing.Category,
		&existing.Reason,
		&existing.BatchUUID,
		&existing.BeerLotID,
		&existing.OccupancyUUID,
		&existing.Amount,
		&existing.AmountUnit,
		&existing.AmountBBL,
		&existing.IsTaxable,
		&existing.ReferenceCode,
		&existing.PerformedBy,
		&existing.RemovedAt,
		&existing.Destination,
		&existing.Notes,
		&existing.MovementID,
		&existing.CreatedAt,
		&existing.UpdatedAt,
		&existing.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryRemoval{}, service.ErrNotFound
		}
		return InventoryRemoval{}, fmt.Errorf("updating removal: %w", err)
	}

	// Re-fetch to get joined fields (beer_lot UUID, movement UUID, stock location UUID).
	return c.GetRemovalByUUID(ctx, removalUUID)
}

// UpdateRemovalRequest describes the mutable fields for a PATCH update.
type UpdateRemovalRequest struct {
	Category      *string
	Reason        *string
	Amount        *int64
	AmountUnit    *string
	AmountBBL     *float64
	IsTaxable     *bool
	ReferenceCode *string
	PerformedBy   *string
	RemovedAt     *time.Time
	Destination   *string
	Notes         *string
}

// SoftDeleteRemoval soft-deletes a removal and its linked movement.
func (c *Client) SoftDeleteRemoval(ctx context.Context, removalUUID string) error {
	tx, err := c.DB().Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting removal delete transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Soft-delete any linked movement first.
	_, err = tx.Exec(ctx, `
		UPDATE inventory_movement SET deleted_at = timezone('utc', now())
		WHERE removal_id = (
			SELECT id FROM inventory_removal WHERE uuid = $1 AND deleted_at IS NULL
		) AND deleted_at IS NULL`,
		removalUUID,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting removal movement: %w", err)
	}

	// Soft-delete the removal.
	tag, err := tx.Exec(ctx, `
		UPDATE inventory_removal SET deleted_at = timezone('utc', now())
		WHERE uuid = $1 AND deleted_at IS NULL`,
		removalUUID,
	)
	if err != nil {
		return fmt.Errorf("soft-deleting removal: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing removal delete transaction: %w", err)
	}
	return nil
}

// GetRemovalSummary returns aggregated removal data for reporting.
func (c *Client) GetRemovalSummary(ctx context.Context, filter RemovalListFilter) (RemovalSummary, error) {
	query := `
		SELECT
			COALESCE(SUM(amount_bbl), 0) AS total_bbl,
			COALESCE(SUM(CASE WHEN is_taxable THEN amount_bbl ELSE 0 END), 0) AS taxable_bbl,
			COALESCE(SUM(CASE WHEN NOT is_taxable THEN amount_bbl ELSE 0 END), 0) AS tax_free_bbl,
			COUNT(*) AS total_count
		FROM inventory_removal r
		LEFT JOIN beer_lot bl ON bl.id = r.beer_lot_id
		WHERE r.deleted_at IS NULL`
	args := []any{}
	argIdx := 1

	if filter.BatchUUID != nil {
		query += fmt.Sprintf(` AND r.batch_uuid = $%d`, argIdx)
		args = append(args, *filter.BatchUUID)
		argIdx++
	}
	if filter.BeerLotUUID != nil {
		query += fmt.Sprintf(` AND bl.uuid = $%d`, argIdx)
		args = append(args, *filter.BeerLotUUID)
		argIdx++
	}
	if filter.Category != nil {
		query += fmt.Sprintf(` AND r.category = $%d`, argIdx)
		args = append(args, *filter.Category)
		argIdx++
	}
	if filter.From != nil {
		query += fmt.Sprintf(` AND r.removed_at >= $%d`, argIdx)
		args = append(args, *filter.From)
		argIdx++
	}
	if filter.To != nil {
		query += fmt.Sprintf(` AND r.removed_at <= $%d`, argIdx)
		args = append(args, *filter.To)
		argIdx++
	}

	var summary RemovalSummary
	err := c.DB().QueryRow(ctx, query, args...).Scan(
		&summary.TotalBBL,
		&summary.TaxableBBL,
		&summary.TaxFreeBBL,
		&summary.TotalCount,
	)
	if err != nil {
		return RemovalSummary{}, fmt.Errorf("getting removal summary: %w", err)
	}

	// Get per-category breakdown.
	catQuery := `
		SELECT category,
			COALESCE(SUM(amount_bbl), 0) AS total_bbl,
			COUNT(*) AS count
		FROM inventory_removal r
		LEFT JOIN beer_lot bl ON bl.id = r.beer_lot_id
		WHERE r.deleted_at IS NULL`
	catArgs := []any{}
	catArgIdx := 1

	if filter.BatchUUID != nil {
		catQuery += fmt.Sprintf(` AND r.batch_uuid = $%d`, catArgIdx)
		catArgs = append(catArgs, *filter.BatchUUID)
		catArgIdx++
	}
	if filter.BeerLotUUID != nil {
		catQuery += fmt.Sprintf(` AND bl.uuid = $%d`, catArgIdx)
		catArgs = append(catArgs, *filter.BeerLotUUID)
		catArgIdx++
	}
	if filter.Category != nil {
		catQuery += fmt.Sprintf(` AND r.category = $%d`, catArgIdx)
		catArgs = append(catArgs, *filter.Category)
		catArgIdx++
	}
	if filter.From != nil {
		catQuery += fmt.Sprintf(` AND r.removed_at >= $%d`, catArgIdx)
		catArgs = append(catArgs, *filter.From)
		catArgIdx++
	}
	if filter.To != nil {
		catQuery += fmt.Sprintf(` AND r.removed_at <= $%d`, catArgIdx)
		catArgs = append(catArgs, *filter.To)
		catArgIdx++
	}

	catQuery += ` GROUP BY category ORDER BY category`

	rows, err := c.DB().Query(ctx, catQuery, catArgs...)
	if err != nil {
		return RemovalSummary{}, fmt.Errorf("getting removal category summary: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cs RemovalCategorySummary
		if err := rows.Scan(&cs.Category, &cs.TotalBBL, &cs.Count); err != nil {
			return RemovalSummary{}, fmt.Errorf("scanning removal category summary: %w", err)
		}
		summary.ByCategory = append(summary.ByCategory, cs)
	}
	if err := rows.Err(); err != nil {
		return RemovalSummary{}, fmt.Errorf("getting removal category summary: %w", err)
	}

	return summary, nil
}
