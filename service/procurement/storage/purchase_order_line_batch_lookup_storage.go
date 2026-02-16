package storage

import (
	"context"
	"fmt"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListPurchaseOrderLinesByUUIDs returns purchase order lines matching any of the given UUIDs.
func (c *Client) ListPurchaseOrderLinesByUUIDs(ctx context.Context, uuids []string) ([]PurchaseOrderLine, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT pol.id, pol.uuid, pol.purchase_order_id, po.uuid, pol.line_number, pol.item_type, pol.item_name, pol.inventory_item_uuid, pol.quantity, pol.quantity_unit, pol.unit_cost_cents, pol.currency, pol.created_at, pol.updated_at, pol.deleted_at
		FROM purchase_order_line pol
		JOIN purchase_order po ON po.id = pol.purchase_order_id
		WHERE pol.uuid = ANY($1::uuid[]) AND pol.deleted_at IS NULL`,
		uuids,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase order lines by uuids: %w", err)
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
		return nil, fmt.Errorf("listing purchase order lines by uuids: %w", err)
	}

	return lines, nil
}
