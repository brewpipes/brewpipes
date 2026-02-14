package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreatePurchaseOrderFee(ctx context.Context, fee PurchaseOrderFee) (PurchaseOrderFee, error) {
	err := c.DB().QueryRow(ctx, `
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

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, fee.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		fee.PurchaseOrderUUID = &poUUID
	}

	return fee, nil
}

func (c *Client) UpdatePurchaseOrderFee(ctx context.Context, id int64, update PurchaseOrderFeeUpdate) (PurchaseOrderFee, error) {
	var fee PurchaseOrderFee
	err := c.DB().QueryRow(ctx, `
		UPDATE purchase_order_fee
		SET
			fee_type = COALESCE($1, fee_type),
			amount_cents = COALESCE($2, amount_cents),
			currency = COALESCE($3, currency),
			updated_at = timezone('utc', now())
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING id, uuid, purchase_order_id, fee_type, amount_cents, currency, created_at, updated_at, deleted_at`,
		update.FeeType,
		update.AmountCents,
		update.Currency,
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
		return PurchaseOrderFee{}, fmt.Errorf("updating purchase order fee: %w", err)
	}

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, fee.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		fee.PurchaseOrderUUID = &poUUID
	}

	return fee, nil
}

func (c *Client) UpdatePurchaseOrderFeeByUUID(ctx context.Context, feeUUID string, update PurchaseOrderFeeUpdate) (PurchaseOrderFee, error) {
	var fee PurchaseOrderFee
	err := c.DB().QueryRow(ctx, `
		UPDATE purchase_order_fee
		SET
			fee_type = COALESCE($1, fee_type),
			amount_cents = COALESCE($2, amount_cents),
			currency = COALESCE($3, currency),
			updated_at = timezone('utc', now())
		WHERE uuid = $4 AND deleted_at IS NULL
		RETURNING id, uuid, purchase_order_id, fee_type, amount_cents, currency, created_at, updated_at, deleted_at`,
		update.FeeType,
		update.AmountCents,
		update.Currency,
		feeUUID,
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
		return PurchaseOrderFee{}, fmt.Errorf("updating purchase order fee by uuid: %w", err)
	}

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, fee.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		fee.PurchaseOrderUUID = &poUUID
	}

	return fee, nil
}

func (c *Client) DeletePurchaseOrderFee(ctx context.Context, id int64) (PurchaseOrderFee, error) {
	var fee PurchaseOrderFee
	err := c.DB().QueryRow(ctx, `
		UPDATE purchase_order_fee
		SET deleted_at = timezone('utc', now()),
			updated_at = timezone('utc', now())
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, uuid, purchase_order_id, fee_type, amount_cents, currency, created_at, updated_at, deleted_at`,
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
		return PurchaseOrderFee{}, fmt.Errorf("deleting purchase order fee: %w", err)
	}

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, fee.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		fee.PurchaseOrderUUID = &poUUID
	}

	return fee, nil
}

func (c *Client) DeletePurchaseOrderFeeByUUID(ctx context.Context, feeUUID string) (PurchaseOrderFee, error) {
	var fee PurchaseOrderFee
	err := c.DB().QueryRow(ctx, `
		UPDATE purchase_order_fee
		SET deleted_at = timezone('utc', now()),
			updated_at = timezone('utc', now())
		WHERE uuid = $1 AND deleted_at IS NULL
		RETURNING id, uuid, purchase_order_id, fee_type, amount_cents, currency, created_at, updated_at, deleted_at`,
		feeUUID,
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
		return PurchaseOrderFee{}, fmt.Errorf("deleting purchase order fee by uuid: %w", err)
	}

	// Resolve purchase order UUID
	var poUUID string
	sErr := c.DB().QueryRow(ctx, `SELECT uuid FROM purchase_order WHERE id = $1`, fee.PurchaseOrderID).Scan(&poUUID)
	if sErr == nil {
		fee.PurchaseOrderUUID = &poUUID
	}

	return fee, nil
}

func (c *Client) GetPurchaseOrderFee(ctx context.Context, id int64) (PurchaseOrderFee, error) {
	var fee PurchaseOrderFee
	err := c.DB().QueryRow(ctx, `
		SELECT pof.id, pof.uuid, pof.purchase_order_id, po.uuid, pof.fee_type, pof.amount_cents, pof.currency, pof.created_at, pof.updated_at, pof.deleted_at
		FROM purchase_order_fee pof
		JOIN purchase_order po ON po.id = pof.purchase_order_id
		WHERE pof.id = $1 AND pof.deleted_at IS NULL`,
		id,
	).Scan(
		&fee.ID,
		&fee.UUID,
		&fee.PurchaseOrderID,
		&fee.PurchaseOrderUUID,
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

func (c *Client) GetPurchaseOrderFeeByUUID(ctx context.Context, feeUUID string) (PurchaseOrderFee, error) {
	var fee PurchaseOrderFee
	err := c.DB().QueryRow(ctx, `
		SELECT pof.id, pof.uuid, pof.purchase_order_id, po.uuid, pof.fee_type, pof.amount_cents, pof.currency, pof.created_at, pof.updated_at, pof.deleted_at
		FROM purchase_order_fee pof
		JOIN purchase_order po ON po.id = pof.purchase_order_id
		WHERE pof.uuid = $1 AND pof.deleted_at IS NULL`,
		feeUUID,
	).Scan(
		&fee.ID,
		&fee.UUID,
		&fee.PurchaseOrderID,
		&fee.PurchaseOrderUUID,
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
		return PurchaseOrderFee{}, fmt.Errorf("getting purchase order fee by uuid: %w", err)
	}

	return fee, nil
}

func (c *Client) ListPurchaseOrderFees(ctx context.Context) ([]PurchaseOrderFee, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT pof.id, pof.uuid, pof.purchase_order_id, po.uuid, pof.fee_type, pof.amount_cents, pof.currency, pof.created_at, pof.updated_at, pof.deleted_at
		FROM purchase_order_fee pof
		JOIN purchase_order po ON po.id = pof.purchase_order_id
		WHERE pof.deleted_at IS NULL
		ORDER BY pof.purchase_order_id, pof.id`,
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
			&fee.PurchaseOrderUUID,
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

func (c *Client) ListPurchaseOrderFeesByOrderUUID(ctx context.Context, orderUUID string) ([]PurchaseOrderFee, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT pof.id, pof.uuid, pof.purchase_order_id, po.uuid, pof.fee_type, pof.amount_cents, pof.currency, pof.created_at, pof.updated_at, pof.deleted_at
		FROM purchase_order_fee pof
		JOIN purchase_order po ON po.id = pof.purchase_order_id
		WHERE po.uuid = $1 AND pof.deleted_at IS NULL
		ORDER BY pof.id`,
		orderUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing purchase order fees by order uuid: %w", err)
	}
	defer rows.Close()

	var fees []PurchaseOrderFee
	for rows.Next() {
		var fee PurchaseOrderFee
		if err := rows.Scan(
			&fee.ID,
			&fee.UUID,
			&fee.PurchaseOrderID,
			&fee.PurchaseOrderUUID,
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
		return nil, fmt.Errorf("listing purchase order fees by order uuid: %w", err)
	}

	return fees, nil
}
