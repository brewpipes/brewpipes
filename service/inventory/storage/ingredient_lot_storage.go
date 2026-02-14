package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Client) CreateIngredientLot(ctx context.Context, lot IngredientLot) (IngredientLot, error) {
	receivedAt := lot.ReceivedAt
	if receivedAt.IsZero() {
		receivedAt = time.Now().UTC()
	}

	var supplierUUID pgtype.UUID
	var purchaseOrderLineUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		INSERT INTO ingredient_lot (
			ingredient_id,
			receipt_id,
			supplier_uuid,
			purchase_order_line_uuid,
			brewery_lot_code,
			originator_lot_code,
			originator_name,
			originator_type,
			received_at,
			received_amount,
			received_unit,
			best_by_at,
			expires_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, uuid, ingredient_id, receipt_id, supplier_uuid, purchase_order_line_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, expires_at, notes, created_at, updated_at, deleted_at`,
		lot.IngredientID,
		lot.ReceiptID,
		database.UUIDParam(lot.SupplierUUID),
		database.UUIDParam(lot.PurchaseOrderLineUUID),
		lot.BreweryLotCode,
		lot.OriginatorLotCode,
		lot.OriginatorName,
		lot.OriginatorType,
		receivedAt,
		lot.ReceivedAmount,
		lot.ReceivedUnit,
		lot.BestByAt,
		lot.ExpiresAt,
		lot.Notes,
	).Scan(
		&lot.ID,
		&lot.UUID,
		&lot.IngredientID,
		&lot.ReceiptID,
		&supplierUUID,
		&purchaseOrderLineUUID,
		&lot.BreweryLotCode,
		&lot.OriginatorLotCode,
		&lot.OriginatorName,
		&lot.OriginatorType,
		&lot.ReceivedAt,
		&lot.ReceivedAmount,
		&lot.ReceivedUnit,
		&lot.BestByAt,
		&lot.ExpiresAt,
		&lot.Notes,
		&lot.CreatedAt,
		&lot.UpdatedAt,
		&lot.DeletedAt,
	)
	if err != nil {
		return IngredientLot{}, fmt.Errorf("creating ingredient lot: %w", err)
	}

	database.AssignUUIDPointer(&lot.SupplierUUID, supplierUUID)
	database.AssignUUIDPointer(&lot.PurchaseOrderLineUUID, purchaseOrderLineUUID)

	// Resolve ingredient UUID
	c.DB().QueryRow(ctx, `SELECT uuid FROM ingredient WHERE id = $1`, lot.IngredientID).Scan(&lot.IngredientUUID)

	// Resolve receipt UUID if set
	if lot.ReceiptID != nil {
		var receiptUUID string
		err := c.DB().QueryRow(ctx, `SELECT uuid FROM inventory_receipt WHERE id = $1`, *lot.ReceiptID).Scan(&receiptUUID)
		if err == nil {
			lot.ReceiptUUID = &receiptUUID
		}
	}

	return lot, nil
}

func (c *Client) GetIngredientLot(ctx context.Context, id int64) (IngredientLot, error) {
	return c.scanIngredientLotRow(c.DB().QueryRow(ctx, ingredientLotSelectSQL+`
		WHERE il.id = $1 AND il.deleted_at IS NULL`, id))
}

func (c *Client) GetIngredientLotByUUID(ctx context.Context, lotUUID string) (IngredientLot, error) {
	return c.scanIngredientLotRow(c.DB().QueryRow(ctx, ingredientLotSelectSQL+`
		WHERE il.uuid = $1 AND il.deleted_at IS NULL`, lotUUID))
}

func (c *Client) ListIngredientLots(ctx context.Context) ([]IngredientLot, error) {
	rows, err := c.DB().Query(ctx, ingredientLotSelectSQL+`
		WHERE il.deleted_at IS NULL
		ORDER BY il.received_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("listing ingredient lots: %w", err)
	}
	defer rows.Close()

	return c.scanIngredientLotRows(rows)
}

func (c *Client) ListIngredientLotsByIngredient(ctx context.Context, ingredientUUID string) ([]IngredientLot, error) {
	rows, err := c.DB().Query(ctx, ingredientLotSelectSQL+`
		WHERE i.uuid = $1 AND il.deleted_at IS NULL
		ORDER BY il.received_at DESC`,
		ingredientUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing ingredient lots by ingredient: %w", err)
	}
	defer rows.Close()

	return c.scanIngredientLotRows(rows)
}

func (c *Client) ListIngredientLotsByReceipt(ctx context.Context, receiptUUID string) ([]IngredientLot, error) {
	rows, err := c.DB().Query(ctx, ingredientLotSelectSQL+`
		WHERE r.uuid = $1 AND il.deleted_at IS NULL
		ORDER BY il.received_at DESC`,
		receiptUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing ingredient lots by receipt: %w", err)
	}
	defer rows.Close()

	return c.scanIngredientLotRows(rows)
}

func (c *Client) ListIngredientLotsByPurchaseOrderLineUUID(ctx context.Context, purchaseOrderLineUUID string) ([]IngredientLot, error) {
	rows, err := c.DB().Query(ctx, ingredientLotSelectSQL+`
		WHERE il.purchase_order_line_uuid = $1 AND il.deleted_at IS NULL
		ORDER BY il.received_at DESC`,
		purchaseOrderLineUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing ingredient lots by purchase order line: %w", err)
	}
	defer rows.Close()

	return c.scanIngredientLotRows(rows)
}

const ingredientLotSelectSQL = `
	SELECT il.id, il.uuid, il.ingredient_id, i.uuid, il.receipt_id, r.uuid,
	       il.supplier_uuid, il.purchase_order_line_uuid,
	       il.brewery_lot_code, il.originator_lot_code, il.originator_name, il.originator_type,
	       il.received_at, il.received_amount, il.received_unit,
	       il.best_by_at, il.expires_at, il.notes,
	       il.created_at, il.updated_at, il.deleted_at
	FROM ingredient_lot il
	JOIN ingredient i ON i.id = il.ingredient_id
	LEFT JOIN inventory_receipt r ON r.id = il.receipt_id
`

func (c *Client) scanIngredientLotRow(row pgx.Row) (IngredientLot, error) {
	var lot IngredientLot
	var supplierUUID pgtype.UUID
	var purchaseOrderLineUUID pgtype.UUID
	err := row.Scan(
		&lot.ID,
		&lot.UUID,
		&lot.IngredientID,
		&lot.IngredientUUID,
		&lot.ReceiptID,
		&lot.ReceiptUUID,
		&supplierUUID,
		&purchaseOrderLineUUID,
		&lot.BreweryLotCode,
		&lot.OriginatorLotCode,
		&lot.OriginatorName,
		&lot.OriginatorType,
		&lot.ReceivedAt,
		&lot.ReceivedAmount,
		&lot.ReceivedUnit,
		&lot.BestByAt,
		&lot.ExpiresAt,
		&lot.Notes,
		&lot.CreatedAt,
		&lot.UpdatedAt,
		&lot.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLot{}, service.ErrNotFound
		}
		return IngredientLot{}, fmt.Errorf("scanning ingredient lot: %w", err)
	}

	database.AssignUUIDPointer(&lot.SupplierUUID, supplierUUID)
	database.AssignUUIDPointer(&lot.PurchaseOrderLineUUID, purchaseOrderLineUUID)
	return lot, nil
}

func (c *Client) scanIngredientLotRows(rows pgx.Rows) ([]IngredientLot, error) {
	var lots []IngredientLot
	for rows.Next() {
		var lot IngredientLot
		var supplierUUID pgtype.UUID
		var purchaseOrderLineUUID pgtype.UUID
		if err := rows.Scan(
			&lot.ID,
			&lot.UUID,
			&lot.IngredientID,
			&lot.IngredientUUID,
			&lot.ReceiptID,
			&lot.ReceiptUUID,
			&supplierUUID,
			&purchaseOrderLineUUID,
			&lot.BreweryLotCode,
			&lot.OriginatorLotCode,
			&lot.OriginatorName,
			&lot.OriginatorType,
			&lot.ReceivedAt,
			&lot.ReceivedAmount,
			&lot.ReceivedUnit,
			&lot.BestByAt,
			&lot.ExpiresAt,
			&lot.Notes,
			&lot.CreatedAt,
			&lot.UpdatedAt,
			&lot.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning ingredient lot: %w", err)
		}
		database.AssignUUIDPointer(&lot.SupplierUUID, supplierUUID)
		database.AssignUUIDPointer(&lot.PurchaseOrderLineUUID, purchaseOrderLineUUID)
		lots = append(lots, lot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing ingredient lots: %w", err)
	}

	return lots, nil
}
