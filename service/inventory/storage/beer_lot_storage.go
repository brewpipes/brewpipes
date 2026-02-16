package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/brewpipes/brewpipes/service"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
)

// beerLotColumns is the column list shared by all beer lot queries.
const beerLotColumns = `id, uuid, production_batch_uuid, packaging_run_uuid, lot_code, best_by,
		package_format_name, container, volume_per_unit, volume_per_unit_unit, quantity,
		packaged_at, notes, created_at, updated_at, deleted_at`

func scanBeerLot(row pgx.Row) (BeerLot, error) {
	var lot BeerLot
	err := row.Scan(
		&lot.ID,
		&lot.UUID,
		&lot.ProductionBatchUUID,
		&lot.PackagingRunUUID,
		&lot.LotCode,
		&lot.BestBy,
		&lot.PackageFormatName,
		&lot.Container,
		&lot.VolumePerUnit,
		&lot.VolumePerUnitUnit,
		&lot.Quantity,
		&lot.PackagedAt,
		&lot.Notes,
		&lot.CreatedAt,
		&lot.UpdatedAt,
		&lot.DeletedAt,
	)
	return lot, err
}

func scanBeerLotRows(rows pgx.Rows) ([]BeerLot, error) {
	var lots []BeerLot
	for rows.Next() {
		var lot BeerLot
		if err := rows.Scan(
			&lot.ID,
			&lot.UUID,
			&lot.ProductionBatchUUID,
			&lot.PackagingRunUUID,
			&lot.LotCode,
			&lot.BestBy,
			&lot.PackageFormatName,
			&lot.Container,
			&lot.VolumePerUnit,
			&lot.VolumePerUnitUnit,
			&lot.Quantity,
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
		return nil, err
	}
	return lots, nil
}

func (c *Client) CreateBeerLot(ctx context.Context, lot BeerLot) (BeerLot, error) {
	packagedAt := lot.PackagedAt
	if packagedAt.IsZero() {
		packagedAt = time.Now().UTC()
	}

	err := c.DB().QueryRow(ctx, `
		INSERT INTO beer_lot (
			production_batch_uuid,
			packaging_run_uuid,
			lot_code,
			best_by,
			package_format_name,
			container,
			volume_per_unit,
			volume_per_unit_unit,
			quantity,
			packaged_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING `+beerLotColumns,
		lot.ProductionBatchUUID,
		database.UUIDParam(lot.PackagingRunUUID),
		lot.LotCode,
		lot.BestBy,
		lot.PackageFormatName,
		lot.Container,
		lot.VolumePerUnit,
		lot.VolumePerUnitUnit,
		lot.Quantity,
		packagedAt,
		lot.Notes,
	).Scan(
		&lot.ID,
		&lot.UUID,
		&lot.ProductionBatchUUID,
		&lot.PackagingRunUUID,
		&lot.LotCode,
		&lot.BestBy,
		&lot.PackageFormatName,
		&lot.Container,
		&lot.VolumePerUnit,
		&lot.VolumePerUnitUnit,
		&lot.Quantity,
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

// CreateBeerLotWithMovement atomically creates a beer lot and an initial
// inventory movement within a single transaction. It returns the created lot
// and the UUID of the movement.
func (c *Client) CreateBeerLotWithMovement(ctx context.Context, lot BeerLot, stockLocationID int64, movementAmount int64, movementAmountUnit string) (BeerLot, uuid.UUID, error) {
	tx, err := c.DB().Begin(ctx)
	if err != nil {
		return BeerLot{}, uuid.UUID{}, fmt.Errorf("starting beer lot transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	packagedAt := lot.PackagedAt
	if packagedAt.IsZero() {
		packagedAt = time.Now().UTC()
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO beer_lot (
			production_batch_uuid,
			packaging_run_uuid,
			lot_code,
			best_by,
			package_format_name,
			container,
			volume_per_unit,
			volume_per_unit_unit,
			quantity,
			packaged_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING `+beerLotColumns,
		lot.ProductionBatchUUID,
		database.UUIDParam(lot.PackagingRunUUID),
		lot.LotCode,
		lot.BestBy,
		lot.PackageFormatName,
		lot.Container,
		lot.VolumePerUnit,
		lot.VolumePerUnitUnit,
		lot.Quantity,
		packagedAt,
		lot.Notes,
	).Scan(
		&lot.ID,
		&lot.UUID,
		&lot.ProductionBatchUUID,
		&lot.PackagingRunUUID,
		&lot.LotCode,
		&lot.BestBy,
		&lot.PackageFormatName,
		&lot.Container,
		&lot.VolumePerUnit,
		&lot.VolumePerUnitUnit,
		&lot.Quantity,
		&lot.PackagedAt,
		&lot.Notes,
		&lot.CreatedAt,
		&lot.UpdatedAt,
		&lot.DeletedAt,
	)
	if err != nil {
		return BeerLot{}, uuid.UUID{}, fmt.Errorf("creating beer lot in transaction: %w", err)
	}

	var movementUUID uuid.UUID
	err = tx.QueryRow(ctx, `
		INSERT INTO inventory_movement (
			beer_lot_id,
			stock_location_id,
			direction,
			reason,
			amount,
			amount_unit,
			occurred_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid`,
		lot.ID,
		stockLocationID,
		MovementDirectionIn,
		MovementReasonPackage,
		movementAmount,
		movementAmountUnit,
		lot.PackagedAt,
	).Scan(&movementUUID)
	if err != nil {
		return BeerLot{}, uuid.UUID{}, fmt.Errorf("creating package movement in transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return BeerLot{}, uuid.UUID{}, fmt.Errorf("committing beer lot transaction: %w", err)
	}

	return lot, movementUUID, nil
}

func (c *Client) GetBeerLot(ctx context.Context, id int64) (BeerLot, error) {
	lot, err := scanBeerLot(c.DB().QueryRow(ctx, `
		SELECT `+beerLotColumns+`
		FROM beer_lot
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BeerLot{}, service.ErrNotFound
		}
		return BeerLot{}, fmt.Errorf("getting beer lot: %w", err)
	}

	return lot, nil
}

func (c *Client) GetBeerLotByUUID(ctx context.Context, lotUUID string) (BeerLot, error) {
	lot, err := scanBeerLot(c.DB().QueryRow(ctx, `
		SELECT `+beerLotColumns+`
		FROM beer_lot
		WHERE uuid = $1 AND deleted_at IS NULL`,
		lotUUID,
	))
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
		SELECT `+beerLotColumns+`
		FROM beer_lot
		WHERE deleted_at IS NULL
		ORDER BY packaged_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing beer lots: %w", err)
	}
	defer rows.Close()

	lots, err := scanBeerLotRows(rows)
	if err != nil {
		return nil, fmt.Errorf("listing beer lots: %w", err)
	}

	return lots, nil
}

func (c *Client) ListBeerLotsByBatchUUID(ctx context.Context, batchUUID uuid.UUID) ([]BeerLot, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT `+beerLotColumns+`
		FROM beer_lot
		WHERE production_batch_uuid = $1 AND deleted_at IS NULL
		ORDER BY packaged_at DESC`,
		batchUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing beer lots by batch UUID: %w", err)
	}
	defer rows.Close()

	lots, err := scanBeerLotRows(rows)
	if err != nil {
		return nil, fmt.Errorf("listing beer lots by batch UUID: %w", err)
	}

	return lots, nil
}
