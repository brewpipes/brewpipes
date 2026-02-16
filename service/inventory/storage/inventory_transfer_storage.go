package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

// TransferWithMovementsRequest describes the parameters for creating a
// transfer record and its two corresponding inventory movements atomically.
type TransferWithMovementsRequest struct {
	IngredientLotID  *int64
	BeerLotID        *int64
	SourceLocationID int64
	DestLocationID   int64
	Amount           int64
	AmountUnit       string
	TransferredAt    time.Time
	Notes            *string

	// UUID fields for populating the response without extra lookups.
	IngredientLotUUID  *string
	BeerLotUUID        *string
	SourceLocationUUID string
	DestLocationUUID   string
}

// TransferWithMovementsResult holds the created transfer and its two movements.
type TransferWithMovementsResult struct {
	Transfer    InventoryTransfer
	MovementOut InventoryMovement
	MovementIn  InventoryMovement
}

// CreateInventoryTransferWithMovements atomically creates an inventory transfer
// record and two corresponding inventory movements (out from source, in to
// destination) within a single transaction.
func (c *Client) CreateInventoryTransferWithMovements(ctx context.Context, req TransferWithMovementsRequest) (TransferWithMovementsResult, error) {
	tx, err := c.DB().Begin(ctx)
	if err != nil {
		return TransferWithMovementsResult{}, fmt.Errorf("starting transfer transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	transferredAt := req.TransferredAt
	if transferredAt.IsZero() {
		transferredAt = time.Now().UTC()
	}

	// Create the transfer record.
	var transfer InventoryTransfer
	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_transfer (
			source_location_id,
			dest_location_id,
			transferred_at,
			notes
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, source_location_id, dest_location_id, transferred_at, notes, created_at, updated_at, deleted_at`,
		req.SourceLocationID,
		req.DestLocationID,
		transferredAt,
		req.Notes,
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
		return TransferWithMovementsResult{}, fmt.Errorf("creating inventory transfer in transaction: %w", err)
	}

	// Create the OUT movement from the source location.
	var movementOut InventoryMovement
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
			transfer_id,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, uuid, ingredient_lot_id, beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, transfer_id, notes, created_at, updated_at, deleted_at`,
		req.IngredientLotID,
		req.BeerLotID,
		req.SourceLocationID,
		MovementDirectionOut,
		MovementReasonTransfer,
		req.Amount,
		req.AmountUnit,
		transferredAt,
		transfer.ID,
		req.Notes,
	).Scan(
		&movementOut.ID,
		&movementOut.UUID,
		&movementOut.IngredientLotID,
		&movementOut.BeerLotID,
		&movementOut.StockLocationID,
		&movementOut.Direction,
		&movementOut.Reason,
		&movementOut.Amount,
		&movementOut.AmountUnit,
		&movementOut.OccurredAt,
		&movementOut.TransferID,
		&movementOut.Notes,
		&movementOut.CreatedAt,
		&movementOut.UpdatedAt,
		&movementOut.DeletedAt,
	)
	if err != nil {
		return TransferWithMovementsResult{}, fmt.Errorf("creating transfer out-movement in transaction: %w", err)
	}

	// Create the IN movement to the destination location.
	var movementIn InventoryMovement
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
			transfer_id,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, uuid, ingredient_lot_id, beer_lot_id, stock_location_id, direction, reason, amount, amount_unit, occurred_at, transfer_id, notes, created_at, updated_at, deleted_at`,
		req.IngredientLotID,
		req.BeerLotID,
		req.DestLocationID,
		MovementDirectionIn,
		MovementReasonTransfer,
		req.Amount,
		req.AmountUnit,
		transferredAt,
		transfer.ID,
		req.Notes,
	).Scan(
		&movementIn.ID,
		&movementIn.UUID,
		&movementIn.IngredientLotID,
		&movementIn.BeerLotID,
		&movementIn.StockLocationID,
		&movementIn.Direction,
		&movementIn.Reason,
		&movementIn.Amount,
		&movementIn.AmountUnit,
		&movementIn.OccurredAt,
		&movementIn.TransferID,
		&movementIn.Notes,
		&movementIn.CreatedAt,
		&movementIn.UpdatedAt,
		&movementIn.DeletedAt,
	)
	if err != nil {
		return TransferWithMovementsResult{}, fmt.Errorf("creating transfer in-movement in transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return TransferWithMovementsResult{}, fmt.Errorf("committing transfer transaction: %w", err)
	}

	// Populate UUID fields for the response without extra lookups.
	transfer.SourceLocationUUID = req.SourceLocationUUID
	transfer.DestLocationUUID = req.DestLocationUUID

	trUUID := transfer.UUID.String()
	movementOut.IngredientLotUUID = req.IngredientLotUUID
	movementOut.BeerLotUUID = req.BeerLotUUID
	movementOut.StockLocationUUID = req.SourceLocationUUID
	movementOut.TransferUUID = &trUUID

	movementIn.IngredientLotUUID = req.IngredientLotUUID
	movementIn.BeerLotUUID = req.BeerLotUUID
	movementIn.StockLocationUUID = req.DestLocationUUID
	movementIn.TransferUUID = &trUUID

	return TransferWithMovementsResult{
		Transfer:    transfer,
		MovementOut: movementOut,
		MovementIn:  movementIn,
	}, nil
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
