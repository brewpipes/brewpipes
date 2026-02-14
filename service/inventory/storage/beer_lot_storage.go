package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateBeerLot(ctx context.Context, lot BeerLot) (BeerLot, error) {
	packagedAt := lot.PackagedAt
	if packagedAt.IsZero() {
		packagedAt = time.Now().UTC()
	}

	err := c.DB().QueryRow(ctx, `
		INSERT INTO beer_lot (
			production_batch_uuid,
			lot_code,
			packaged_at,
			notes
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, production_batch_uuid, lot_code, packaged_at, notes, created_at, updated_at, deleted_at`,
		lot.ProductionBatchUUID,
		lot.LotCode,
		packagedAt,
		lot.Notes,
	).Scan(
		&lot.ID,
		&lot.UUID,
		&lot.ProductionBatchUUID,
		&lot.LotCode,
		&lot.PackagedAt,
		&lot.Notes,
		&lot.CreatedAt,
		&lot.UpdatedAt,
		&lot.DeletedAt,
	)
	if err != nil {
		return BeerLot{}, fmt.Errorf("creating beer lot: %w", err)
	}

	return lot, nil
}

func (c *Client) GetBeerLot(ctx context.Context, id int64) (BeerLot, error) {
	var lot BeerLot
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, production_batch_uuid, lot_code, packaged_at, notes, created_at, updated_at, deleted_at
		FROM beer_lot
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&lot.ID,
		&lot.UUID,
		&lot.ProductionBatchUUID,
		&lot.LotCode,
		&lot.PackagedAt,
		&lot.Notes,
		&lot.CreatedAt,
		&lot.UpdatedAt,
		&lot.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BeerLot{}, service.ErrNotFound
		}
		return BeerLot{}, fmt.Errorf("getting beer lot: %w", err)
	}

	return lot, nil
}

func (c *Client) GetBeerLotByUUID(ctx context.Context, lotUUID string) (BeerLot, error) {
	var lot BeerLot
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, production_batch_uuid, lot_code, packaged_at, notes, created_at, updated_at, deleted_at
		FROM beer_lot
		WHERE uuid = $1 AND deleted_at IS NULL`,
		lotUUID,
	).Scan(
		&lot.ID,
		&lot.UUID,
		&lot.ProductionBatchUUID,
		&lot.LotCode,
		&lot.PackagedAt,
		&lot.Notes,
		&lot.CreatedAt,
		&lot.UpdatedAt,
		&lot.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BeerLot{}, service.ErrNotFound
		}
		return BeerLot{}, fmt.Errorf("getting beer lot by uuid: %w", err)
	}

	return lot, nil
}

func (c *Client) ListBeerLots(ctx context.Context) ([]BeerLot, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, production_batch_uuid, lot_code, packaged_at, notes, created_at, updated_at, deleted_at
		FROM beer_lot
		WHERE deleted_at IS NULL
		ORDER BY packaged_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing beer lots: %w", err)
	}
	defer rows.Close()

	var lots []BeerLot
	for rows.Next() {
		var lot BeerLot
		if err := rows.Scan(
			&lot.ID,
			&lot.UUID,
			&lot.ProductionBatchUUID,
			&lot.LotCode,
			&lot.PackagedAt,
			&lot.Notes,
			&lot.CreatedAt,
			&lot.UpdatedAt,
			&lot.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning beer lot: %w", err)
		}
		lots = append(lots, lot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing beer lots: %w", err)
	}

	return lots, nil
}

func (c *Client) ListBeerLotsByBatchUUID(ctx context.Context, batchUUID uuid.UUID) ([]BeerLot, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, production_batch_uuid, lot_code, packaged_at, notes, created_at, updated_at, deleted_at
		FROM beer_lot
		WHERE production_batch_uuid = $1 AND deleted_at IS NULL
		ORDER BY packaged_at DESC`,
		batchUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing beer lots by batch UUID: %w", err)
	}
	defer rows.Close()

	var lots []BeerLot
	for rows.Next() {
		var lot BeerLot
		if err := rows.Scan(
			&lot.ID,
			&lot.UUID,
			&lot.ProductionBatchUUID,
			&lot.LotCode,
			&lot.PackagedAt,
			&lot.Notes,
			&lot.CreatedAt,
			&lot.UpdatedAt,
			&lot.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning beer lot: %w", err)
		}
		lots = append(lots, lot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing beer lots by batch UUID: %w", err)
	}

	return lots, nil
}
