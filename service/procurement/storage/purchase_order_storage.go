package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreatePurchaseOrder(ctx context.Context, order PurchaseOrder) (PurchaseOrder, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO purchase_order (
			supplier_id,
			order_number,
			status,
			ordered_at,
			expected_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, uuid, supplier_id, order_number, status, ordered_at, expected_at, notes, created_at, updated_at, deleted_at`,
		order.SupplierID,
		order.OrderNumber,
		order.Status,
		order.OrderedAt,
		order.ExpectedAt,
		order.Notes,
	).Scan(
		&order.ID,
		&order.UUID,
		&order.SupplierID,
		&order.OrderNumber,
		&order.Status,
		&order.OrderedAt,
		&order.ExpectedAt,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
	)
	if err != nil {
		return PurchaseOrder{}, fmt.Errorf("creating purchase order: %w", err)
	}

	return order, nil
}

func (c *Client) GetPurchaseOrder(ctx context.Context, id int64) (PurchaseOrder, error) {
	var order PurchaseOrder
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, supplier_id, order_number, status, ordered_at, expected_at, notes, created_at, updated_at, deleted_at
		FROM purchase_order
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&order.ID,
		&order.UUID,
		&order.SupplierID,
		&order.OrderNumber,
		&order.Status,
		&order.OrderedAt,
		&order.ExpectedAt,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PurchaseOrder{}, service.ErrNotFound
		}
		return PurchaseOrder{}, fmt.Errorf("getting purchase order: %w", err)
	}

	return order, nil
}

func (c *Client) ListPurchaseOrders(ctx context.Context) ([]PurchaseOrder, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, supplier_id, order_number, status, ordered_at, expected_at, notes, created_at, updated_at, deleted_at
		FROM purchase_order
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase orders: %w", err)
	}
	defer rows.Close()

	var orders []PurchaseOrder
	for rows.Next() {
		var order PurchaseOrder
		if err := rows.Scan(
			&order.ID,
			&order.UUID,
			&order.SupplierID,
			&order.OrderNumber,
			&order.Status,
			&order.OrderedAt,
			&order.ExpectedAt,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning purchase order: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing purchase orders: %w", err)
	}

	return orders, nil
}

func (c *Client) ListPurchaseOrdersBySupplier(ctx context.Context, supplierID int64) ([]PurchaseOrder, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, supplier_id, order_number, status, ordered_at, expected_at, notes, created_at, updated_at, deleted_at
		FROM purchase_order
		WHERE supplier_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`,
		supplierID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase orders by supplier: %w", err)
	}
	defer rows.Close()

	var orders []PurchaseOrder
	for rows.Next() {
		var order PurchaseOrder
		if err := rows.Scan(
			&order.ID,
			&order.UUID,
			&order.SupplierID,
			&order.OrderNumber,
			&order.Status,
			&order.OrderedAt,
			&order.ExpectedAt,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning purchase order: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing purchase orders by supplier: %w", err)
	}

	return orders, nil
}
