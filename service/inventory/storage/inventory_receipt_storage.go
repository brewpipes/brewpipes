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

func (c *Client) CreateInventoryReceipt(ctx context.Context, receipt InventoryReceipt) (InventoryReceipt, error) {
	receivedAt := receipt.ReceivedAt
	if receivedAt.IsZero() {
		receivedAt = time.Now().UTC()
	}

	var supplierUUID pgtype.UUID
	var purchaseOrderUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		INSERT INTO inventory_receipt (
			supplier_uuid,
			purchase_order_uuid,
			reference_code,
			received_at,
			notes
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, uuid, supplier_uuid, purchase_order_uuid, reference_code, received_at, notes, created_at, updated_at, deleted_at`,
		database.UUIDParam(receipt.SupplierUUID),
		database.UUIDParam(receipt.PurchaseOrderUUID),
		receipt.ReferenceCode,
		receivedAt,
		receipt.Notes,
	).Scan(
		&receipt.ID,
		&receipt.UUID,
		&supplierUUID,
		&purchaseOrderUUID,
		&receipt.ReferenceCode,
		&receipt.ReceivedAt,
		&receipt.Notes,
		&receipt.CreatedAt,
		&receipt.UpdatedAt,
		&receipt.DeletedAt,
	)
	if err != nil {
		return InventoryReceipt{}, fmt.Errorf("creating inventory receipt: %w", err)
	}

	database.AssignUUIDPointer(&receipt.SupplierUUID, supplierUUID)
	database.AssignUUIDPointer(&receipt.PurchaseOrderUUID, purchaseOrderUUID)
	return receipt, nil
}

func (c *Client) GetInventoryReceipt(ctx context.Context, id int64) (InventoryReceipt, error) {
	var receipt InventoryReceipt
	var supplierUUID pgtype.UUID
	var purchaseOrderUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, supplier_uuid, purchase_order_uuid, reference_code, received_at, notes, created_at, updated_at, deleted_at
		FROM inventory_receipt
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&receipt.ID,
		&receipt.UUID,
		&supplierUUID,
		&purchaseOrderUUID,
		&receipt.ReferenceCode,
		&receipt.ReceivedAt,
		&receipt.Notes,
		&receipt.CreatedAt,
		&receipt.UpdatedAt,
		&receipt.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryReceipt{}, service.ErrNotFound
		}
		return InventoryReceipt{}, fmt.Errorf("getting inventory receipt: %w", err)
	}

	database.AssignUUIDPointer(&receipt.SupplierUUID, supplierUUID)
	database.AssignUUIDPointer(&receipt.PurchaseOrderUUID, purchaseOrderUUID)
	return receipt, nil
}

func (c *Client) GetInventoryReceiptByUUID(ctx context.Context, receiptUUID string) (InventoryReceipt, error) {
	var receipt InventoryReceipt
	var supplierUUID pgtype.UUID
	var purchaseOrderUUID pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, supplier_uuid, purchase_order_uuid, reference_code, received_at, notes, created_at, updated_at, deleted_at
		FROM inventory_receipt
		WHERE uuid = $1 AND deleted_at IS NULL`,
		receiptUUID,
	).Scan(
		&receipt.ID,
		&receipt.UUID,
		&supplierUUID,
		&purchaseOrderUUID,
		&receipt.ReferenceCode,
		&receipt.ReceivedAt,
		&receipt.Notes,
		&receipt.CreatedAt,
		&receipt.UpdatedAt,
		&receipt.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return InventoryReceipt{}, service.ErrNotFound
		}
		return InventoryReceipt{}, fmt.Errorf("getting inventory receipt by uuid: %w", err)
	}

	database.AssignUUIDPointer(&receipt.SupplierUUID, supplierUUID)
	database.AssignUUIDPointer(&receipt.PurchaseOrderUUID, purchaseOrderUUID)
	return receipt, nil
}

func (c *Client) ListInventoryReceipts(ctx context.Context) ([]InventoryReceipt, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, supplier_uuid, purchase_order_uuid, reference_code, received_at, notes, created_at, updated_at, deleted_at
		FROM inventory_receipt
		WHERE deleted_at IS NULL
		ORDER BY received_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory receipts: %w", err)
	}
	defer rows.Close()

	var receipts []InventoryReceipt
	for rows.Next() {
		var receipt InventoryReceipt
		var supplierUUID pgtype.UUID
		var purchaseOrderUUID pgtype.UUID
		if err := rows.Scan(
			&receipt.ID,
			&receipt.UUID,
			&supplierUUID,
			&purchaseOrderUUID,
			&receipt.ReferenceCode,
			&receipt.ReceivedAt,
			&receipt.Notes,
			&receipt.CreatedAt,
			&receipt.UpdatedAt,
			&receipt.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning inventory receipt: %w", err)
		}
		database.AssignUUIDPointer(&receipt.SupplierUUID, supplierUUID)
		database.AssignUUIDPointer(&receipt.PurchaseOrderUUID, purchaseOrderUUID)
		receipts = append(receipts, receipt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing inventory receipts: %w", err)
	}

	return receipts, nil
}

func (c *Client) ListInventoryReceiptsByPurchaseOrderUUID(ctx context.Context, purchaseOrderUUID string) ([]InventoryReceipt, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, supplier_uuid, purchase_order_uuid, reference_code, received_at, notes, created_at, updated_at, deleted_at
		FROM inventory_receipt
		WHERE purchase_order_uuid = $1 AND deleted_at IS NULL
		ORDER BY received_at DESC`,
		purchaseOrderUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing inventory receipts by purchase order: %w", err)
	}
	defer rows.Close()

	var receipts []InventoryReceipt
	for rows.Next() {
		var receipt InventoryReceipt
		var supplierUUID pgtype.UUID
		var purchaseOrderUUID pgtype.UUID
		if err := rows.Scan(
			&receipt.ID,
			&receipt.UUID,
			&supplierUUID,
			&purchaseOrderUUID,
			&receipt.ReferenceCode,
			&receipt.ReceivedAt,
			&receipt.Notes,
			&receipt.CreatedAt,
			&receipt.UpdatedAt,
			&receipt.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning inventory receipt: %w", err)
		}
		database.AssignUUIDPointer(&receipt.SupplierUUID, supplierUUID)
		database.AssignUUIDPointer(&receipt.PurchaseOrderUUID, purchaseOrderUUID)
		receipts = append(receipts, receipt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing inventory receipts by purchase order: %w", err)
	}

	return receipts, nil
}
