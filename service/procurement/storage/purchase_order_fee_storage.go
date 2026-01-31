package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreatePurchaseOrderFee(ctx context.Context, fee PurchaseOrderFee) (PurchaseOrderFee, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO purchase_order_fee (
			purchase_order_id,
			fee_type,
			amount_cents,
			currency
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, purchase_order_id, fee_type, amount_cents, currency, created_at, updated_at, deleted_at`,
		fee.PurchaseOrderID,
		fee.FeeType,
		fee.AmountCents,
		fee.Currency,
	).Scan(
		&fee.ID,
		&fee.UUID,
		&fee.PurchaseOrderID,
		&fee.FeeType,
		&fee.AmountCents,
		&fee.Currency,
		&fee.CreatedAt,
		&fee.UpdatedAt,
		&fee.DeletedAt,
	)
	if err != nil {
		return PurchaseOrderFee{}, fmt.Errorf("creating purchase order fee: %w", err)
	}

	return fee, nil
}

func (c *Client) GetPurchaseOrderFee(ctx context.Context, id int64) (PurchaseOrderFee, error) {
	var fee PurchaseOrderFee
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, purchase_order_id, fee_type, amount_cents, currency, created_at, updated_at, deleted_at
		FROM purchase_order_fee
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&fee.ID,
		&fee.UUID,
		&fee.PurchaseOrderID,
		&fee.FeeType,
		&fee.AmountCents,
		&fee.Currency,
		&fee.CreatedAt,
		&fee.UpdatedAt,
		&fee.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PurchaseOrderFee{}, service.ErrNotFound
		}
		return PurchaseOrderFee{}, fmt.Errorf("getting purchase order fee: %w", err)
	}

	return fee, nil
}

func (c *Client) ListPurchaseOrderFees(ctx context.Context) ([]PurchaseOrderFee, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, purchase_order_id, fee_type, amount_cents, currency, created_at, updated_at, deleted_at
		FROM purchase_order_fee
		WHERE deleted_at IS NULL
		ORDER BY purchase_order_id, id`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase order fees: %w", err)
	}
	defer rows.Close()

	var fees []PurchaseOrderFee
	for rows.Next() {
		var fee PurchaseOrderFee
		if err := rows.Scan(
			&fee.ID,
			&fee.UUID,
			&fee.PurchaseOrderID,
			&fee.FeeType,
			&fee.AmountCents,
			&fee.Currency,
			&fee.CreatedAt,
			&fee.UpdatedAt,
			&fee.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning purchase order fee: %w", err)
		}
		fees = append(fees, fee)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing purchase order fees: %w", err)
	}

	return fees, nil
}

func (c *Client) ListPurchaseOrderFeesByOrder(ctx context.Context, orderID int64) ([]PurchaseOrderFee, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, purchase_order_id, fee_type, amount_cents, currency, created_at, updated_at, deleted_at
		FROM purchase_order_fee
		WHERE purchase_order_id = $1 AND deleted_at IS NULL
		ORDER BY id`,
		orderID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase order fees by order: %w", err)
	}
	defer rows.Close()

	var fees []PurchaseOrderFee
	for rows.Next() {
		var fee PurchaseOrderFee
		if err := rows.Scan(
			&fee.ID,
			&fee.UUID,
			&fee.PurchaseOrderID,
			&fee.FeeType,
			&fee.AmountCents,
			&fee.Currency,
			&fee.CreatedAt,
			&fee.UpdatedAt,
			&fee.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning purchase order fee: %w", err)
		}
		fees = append(fees, fee)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing purchase order fees by order: %w", err)
	}

	return fees, nil
}
