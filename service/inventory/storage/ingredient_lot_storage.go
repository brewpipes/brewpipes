package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	err := c.db.QueryRow(ctx, `
		INSERT INTO ingredient_lot (
			ingredient_id,
			receipt_id,
			supplier_uuid,
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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, expires_at, notes, created_at, updated_at, deleted_at`,
		lot.IngredientID,
		lot.ReceiptID,
		uuidParam(lot.SupplierUUID),
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

	assignUUIDPointer(&lot.SupplierUUID, supplierUUID)
	return lot, nil
}

func (c *Client) GetIngredientLot(ctx context.Context, id int64) (IngredientLot, error) {
	var lot IngredientLot
	var supplierUUID pgtype.UUID
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, expires_at, notes, created_at, updated_at, deleted_at
		FROM ingredient_lot
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&lot.ID,
		&lot.UUID,
		&lot.IngredientID,
		&lot.ReceiptID,
		&supplierUUID,
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
		return IngredientLot{}, fmt.Errorf("getting ingredient lot: %w", err)
	}

	assignUUIDPointer(&lot.SupplierUUID, supplierUUID)
	return lot, nil
}

func (c *Client) ListIngredientLots(ctx context.Context) ([]IngredientLot, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, expires_at, notes, created_at, updated_at, deleted_at
		FROM ingredient_lot
		WHERE deleted_at IS NULL
		ORDER BY received_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing ingredient lots: %w", err)
	}
	defer rows.Close()

	var lots []IngredientLot
	for rows.Next() {
		var lot IngredientLot
		var supplierUUID pgtype.UUID
		if err := rows.Scan(
			&lot.ID,
			&lot.UUID,
			&lot.IngredientID,
			&lot.ReceiptID,
			&supplierUUID,
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
		assignUUIDPointer(&lot.SupplierUUID, supplierUUID)
		lots = append(lots, lot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing ingredient lots: %w", err)
	}

	return lots, nil
}

func (c *Client) ListIngredientLotsByIngredient(ctx context.Context, ingredientID int64) ([]IngredientLot, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, expires_at, notes, created_at, updated_at, deleted_at
		FROM ingredient_lot
		WHERE ingredient_id = $1 AND deleted_at IS NULL
		ORDER BY received_at DESC`,
		ingredientID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing ingredient lots by ingredient: %w", err)
	}
	defer rows.Close()

	var lots []IngredientLot
	for rows.Next() {
		var lot IngredientLot
		var supplierUUID pgtype.UUID
		if err := rows.Scan(
			&lot.ID,
			&lot.UUID,
			&lot.IngredientID,
			&lot.ReceiptID,
			&supplierUUID,
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
		assignUUIDPointer(&lot.SupplierUUID, supplierUUID)
		lots = append(lots, lot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing ingredient lots by ingredient: %w", err)
	}

	return lots, nil
}

func (c *Client) ListIngredientLotsByReceipt(ctx context.Context, receiptID int64) ([]IngredientLot, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, ingredient_id, receipt_id, supplier_uuid, brewery_lot_code, originator_lot_code, originator_name, originator_type, received_at, received_amount, received_unit, best_by_at, expires_at, notes, created_at, updated_at, deleted_at
		FROM ingredient_lot
		WHERE receipt_id = $1 AND deleted_at IS NULL
		ORDER BY received_at DESC`,
		receiptID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing ingredient lots by receipt: %w", err)
	}
	defer rows.Close()

	var lots []IngredientLot
	for rows.Next() {
		var lot IngredientLot
		var supplierUUID pgtype.UUID
		if err := rows.Scan(
			&lot.ID,
			&lot.UUID,
			&lot.IngredientID,
			&lot.ReceiptID,
			&supplierUUID,
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
		assignUUIDPointer(&lot.SupplierUUID, supplierUUID)
		lots = append(lots, lot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing ingredient lots by receipt: %w", err)
	}

	return lots, nil
}
