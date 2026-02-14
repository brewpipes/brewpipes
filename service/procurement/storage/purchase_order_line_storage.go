package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (c *Client) CreatePurchaseOrderLine(ctx context.Context, line PurchaseOrderLine) (PurchaseOrderLine, error) {
	var inventoryItemUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
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
		database.UUIDParam(line.InventoryItemUUID),
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

	database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, line.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		line.PurchaseOrderUUID = &poUUID
	}

	return line, nil
}

func (c *Client) UpdatePurchaseOrderLine(ctx context.Context, id int64, update PurchaseOrderLineUpdate) (PurchaseOrderLine, error) {
	var line PurchaseOrderLine
	var inventoryItemUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		UPDATE purchase_order_line
		SET
			line_number = COALESCE($1, line_number),
			item_type = COALESCE($2, item_type),
			item_name = COALESCE($3, item_name),
			inventory_item_uuid = COALESCE($4, inventory_item_uuid),
			quantity = COALESCE($5, quantity),
			quantity_unit = COALESCE($6, quantity_unit),
			unit_cost_cents = COALESCE($7, unit_cost_cents),
			currency = COALESCE($8, currency),
			updated_at = timezone('utc', now())
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING id, uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency, created_at, updated_at, deleted_at`,
		update.LineNumber,
		update.ItemType,
		update.ItemName,
		database.UUIDParam(update.InventoryItemUUID),
		update.Quantity,
		update.QuantityUnit,
		update.UnitCostCents,
		update.Currency,
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
		return PurchaseOrderLine{}, fmt.Errorf("updating purchase order line: %w", err)
	}

	database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, line.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		line.PurchaseOrderUUID = &poUUID
	}

	return line, nil
}

func (c *Client) UpdatePurchaseOrderLineByUUID(ctx context.Context, lineUUID string, update PurchaseOrderLineUpdate) (PurchaseOrderLine, error) {
	var line PurchaseOrderLine
	var inventoryItemUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		UPDATE purchase_order_line
		SET
			line_number = COALESCE($1, line_number),
			item_type = COALESCE($2, item_type),
			item_name = COALESCE($3, item_name),
			inventory_item_uuid = COALESCE($4, inventory_item_uuid),
			quantity = COALESCE($5, quantity),
			quantity_unit = COALESCE($6, quantity_unit),
			unit_cost_cents = COALESCE($7, unit_cost_cents),
			currency = COALESCE($8, currency),
			updated_at = timezone('utc', now())
		WHERE uuid = $9 AND deleted_at IS NULL
		RETURNING id, uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency, created_at, updated_at, deleted_at`,
		update.LineNumber,
		update.ItemType,
		update.ItemName,
		database.UUIDParam(update.InventoryItemUUID),
		update.Quantity,
		update.QuantityUnit,
		update.UnitCostCents,
		update.Currency,
		lineUUID,
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
		return PurchaseOrderLine{}, fmt.Errorf("updating purchase order line by uuid: %w", err)
	}

	database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, line.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		line.PurchaseOrderUUID = &poUUID
	}

	return line, nil
}

func (c *Client) DeletePurchaseOrderLine(ctx context.Context, id int64) (PurchaseOrderLine, error) {
	var line PurchaseOrderLine
	var inventoryItemUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		UPDATE purchase_order_line
		SET deleted_at = timezone('utc', now()),
			updated_at = timezone('utc', now())
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency, created_at, updated_at, deleted_at`,
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
		return PurchaseOrderLine{}, fmt.Errorf("deleting purchase order line: %w", err)
	}

	database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, line.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		line.PurchaseOrderUUID = &poUUID
	}

	return line, nil
}

func (c *Client) DeletePurchaseOrderLineByUUID(ctx context.Context, lineUUID string) (PurchaseOrderLine, error) {
	var line PurchaseOrderLine
	var inventoryItemUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		UPDATE purchase_order_line
		SET deleted_at = timezone('utc', now()),
			updated_at = timezone('utc', now())
		WHERE uuid = $1 AND deleted_at IS NULL
		RETURNING id, uuid, purchase_order_id, line_number, item_type, item_name, inventory_item_uuid, quantity, quantity_unit, unit_cost_cents, currency, created_at, updated_at, deleted_at`,
		lineUUID,
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
		return PurchaseOrderLine{}, fmt.Errorf("deleting purchase order line by uuid: %w", err)
	}

	database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, line.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		line.PurchaseOrderUUID = &poUUID
	}

	return line, nil
}

func (c *Client) GetPurchaseOrderLine(ctx context.Context, id int64) (PurchaseOrderLine, error) {
	var line PurchaseOrderLine
	var inventoryItemUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		SELECT pol.id, pol.uuid, pol.purchase_order_id, po.uuid, pol.line_number, pol.item_type, pol.item_name, pol.inventory_item_uuid, pol.quantity, pol.quantity_unit, pol.unit_cost_cents, pol.currency, pol.created_at, pol.updated_at, pol.deleted_at
		FROM purchase_order_line pol
		JOIN purchase_order po ON po.id = pol.purchase_order_id
		WHERE pol.id = $1 AND pol.deleted_at IS NULL`,
		id,
	).Scan(
		&line.ID,
		&line.UUID,
		&line.PurchaseOrderID,
		&line.PurchaseOrderUUID,
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

	database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)
	return line, nil
}

func (c *Client) GetPurchaseOrderLineByUUID(ctx context.Context, lineUUID string) (PurchaseOrderLine, error) {
	var line PurchaseOrderLine
	var inventoryItemUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		SELECT pol.id, pol.uuid, pol.purchase_order_id, po.uuid, pol.line_number, pol.item_type, pol.item_name, pol.inventory_item_uuid, pol.quantity, pol.quantity_unit, pol.unit_cost_cents, pol.currency, pol.created_at, pol.updated_at, pol.deleted_at
		FROM purchase_order_line pol
		JOIN purchase_order po ON po.id = pol.purchase_order_id
		WHERE pol.uuid = $1 AND pol.deleted_at IS NULL`,
		lineUUID,
	).Scan(
		&line.ID,
		&line.UUID,
		&line.PurchaseOrderID,
		&line.PurchaseOrderUUID,
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
		return PurchaseOrderLine{}, fmt.Errorf("getting purchase order line by uuid: %w", err)
	}

	database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)
	return line, nil
}

func (c *Client) ListPurchaseOrderLines(ctx context.Context) ([]PurchaseOrderLine, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT pol.id, pol.uuid, pol.purchase_order_id, po.uuid, pol.line_number, pol.item_type, pol.item_name, pol.inventory_item_uuid, pol.quantity, pol.quantity_unit, pol.unit_cost_cents, pol.currency, pol.created_at, pol.updated_at, pol.deleted_at
		FROM purchase_order_line pol
		JOIN purchase_order po ON po.id = pol.purchase_order_id
		WHERE pol.deleted_at IS NULL
		ORDER BY pol.purchase_order_id, pol.line_number`,
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
			&line.PurchaseOrderUUID,
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
		database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing purchase order lines: %w", err)
	}

	return lines, nil
}

func (c *Client) ListPurchaseOrderLinesByOrderUUID(ctx context.Context, orderUUID string) ([]PurchaseOrderLine, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT pol.id, pol.uuid, pol.purchase_order_id, po.uuid, pol.line_number, pol.item_type, pol.item_name, pol.inventory_item_uuid, pol.quantity, pol.quantity_unit, pol.unit_cost_cents, pol.currency, pol.created_at, pol.updated_at, pol.deleted_at
		FROM purchase_order_line pol
		JOIN purchase_order po ON po.id = pol.purchase_order_id
		WHERE po.uuid = $1 AND pol.deleted_at IS NULL
		ORDER BY pol.line_number`,
		orderUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase order lines by order uuid: %w", err)
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
			&line.PurchaseOrderUUID,
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
		database.AssignUUIDPointer(&line.InventoryItemUUID, inventoryItemUUID)
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing purchase order lines by order uuid: %w", err)
	}

	return lines, nil
}
