package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Client) CreatePurchaseOrderLine(ctx context.Context, line PurchaseOrderLine) (PurchaseOrderLine, error) {
	var inventoryItemUUID pgtype.UUID
	err := c.db.QueryRow(ctx, `
		INSERT INTO purchase_order_line (
			purchase_order_id,
			line_number,
			item_type,
			item_name,
			inventory_item_uuid,
			quantity,
			quantity_unit,
			unit_cost_cents,
			currency
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency, created_at, updated_at, deleted_at`,
		line.PurchaseOrderID,
		line.LineNumber,
		line.ItemType,
		line.ItemName,
		uuidParam(line.InventoryItemUUID),
		line.Quantity,
		line.QuantityUnit,
		line.UnitCostCents,
		line.Currency,
	).Scan(
		&line.ID,
		&line.UUID,
		&line.PurchaseOrderID,
		&line.LineNumber,
		&line.ItemType,
		&line.ItemName,
		&inventoryItemUUID,
		&line.Quantity,
		&line.QuantityUnit,
		&line.UnitCostCents,
		&line.Currency,
		&line.CreatedAt,
		&line.UpdatedAt,
		&line.DeletedAt,
	)
	if err != nil {
		return PurchaseOrderLine{}, fmt.Errorf("creating purchase order line: %w", err)
	}

	assignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)
	return line, nil
}

func (c *Client) GetPurchaseOrderLine(ctx context.Context, id int64) (PurchaseOrderLine, error) {
	var line PurchaseOrderLine
	var inventoryItemUUID pgtype.UUID
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency, created_at, updated_at, deleted_at
		FROM purchase_order_line
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&line.ID,
		&line.UUID,
		&line.PurchaseOrderID,
		&line.LineNumber,
		&line.ItemType,
		&line.ItemName,
		&inventoryItemUUID,
		&line.Quantity,
		&line.QuantityUnit,
		&line.UnitCostCents,
		&line.Currency,
		&line.CreatedAt,
		&line.UpdatedAt,
		&line.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PurchaseOrderLine{}, service.ErrNotFound
		}
		return PurchaseOrderLine{}, fmt.Errorf("getting purchase order line: %w", err)
	}

	assignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)
	return line, nil
}

func (c *Client) ListPurchaseOrderLines(ctx context.Context) ([]PurchaseOrderLine, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency, created_at, updated_at, deleted_at
		FROM purchase_order_line
		WHERE deleted_at IS NULL
		ORDER BY purchase_order_id, line_number`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase order lines: %w", err)
	}
	defer rows.Close()

	var lines []PurchaseOrderLine
	for rows.Next() {
		var line PurchaseOrderLine
		var inventoryItemUUID pgtype.UUID
		if err := rows.Scan(
			&line.ID,
			&line.UUID,
			&line.PurchaseOrderID,
			&line.LineNumber,
			&line.ItemType,
			&line.ItemName,
			&inventoryItemUUID,
			&line.Quantity,
			&line.QuantityUnit,
			&line.UnitCostCents,
			&line.Currency,
			&line.CreatedAt,
			&line.UpdatedAt,
			&line.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning purchase order line: %w", err)
		}
		assignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing purchase order lines: %w", err)
	}

	return lines, nil
}

func (c *Client) ListPurchaseOrderLinesByOrder(ctx context.Context, orderID int64) ([]PurchaseOrderLine, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency, created_at, updated_at, deleted_at
		FROM purchase_order_line
		WHERE purchase_order_id = $1 AND deleted_at IS NULL
		ORDER BY line_number`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase order lines by order: %w", err)
	}
	defer rows.Close()

	var lines []PurchaseOrderLine
	for rows.Next() {
		var line PurchaseOrderLine
		var inventoryItemUUID pgtype.UUID
		if err := rows.Scan(
			&line.ID,
			&line.UUID,
			&line.PurchaseOrderID,
			&line.LineNumber,
			&line.ItemType,
			&line.ItemName,
			&inventoryItemUUID,
			&line.Quantity,
			&line.QuantityUnit,
			&line.UnitCostCents,
			&line.Currency,
			&line.CreatedAt,
			&line.UpdatedAt,
			&line.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning purchase order line: %w", err)
		}
		assignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing purchase order lines by order: %w", err)
	}

	return lines, nil
}
