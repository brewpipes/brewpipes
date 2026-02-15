package storage

import (
	"context"
	"fmt"
	"time"
)

// BeerLotStockLevel represents the current stock level for a beer lot at a location.
type BeerLotStockLevel struct {
	BeerLotUUID         string
	ProductionBatchUUID string
	LotCode             *string
	PackagedAt          *time.Time
	BestBy              *time.Time
	PackageFormatName   *string
	Container           *string
	VolumePerUnit       *int64
	VolumePerUnitUnit   *string
	InitialQuantity     *int
	LocationUUID        string
	LocationName        string
	CurrentVolume       int64
	CurrentVolumeUnit   string
	CurrentQuantity     *int
}

// GetBeerLotStockLevels returns current stock levels for all beer lots, computed
// from the inventory movement ledger. Stock is grouped by beer lot and location.
// Only lots with positive stock are included.
func (c *Client) GetBeerLotStockLevels(ctx context.Context) ([]BeerLotStockLevel, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT
			bl.uuid AS beer_lot_uuid,
			bl.production_batch_uuid,
			bl.lot_code,
			bl.packaged_at,
			bl.best_by,
			bl.package_format_name,
			bl.container,
			bl.volume_per_unit,
			bl.volume_per_unit_unit,
			bl.quantity AS initial_quantity,
			sl.uuid AS location_uuid,
			sl.name AS location_name,
			SUM(CASE m.direction
				WHEN 'in' THEN m.amount
				WHEN 'out' THEN -m.amount
			END) AS current_volume,
			m.amount_unit AS current_volume_unit
		FROM inventory_movement m
		JOIN beer_lot bl ON bl.id = m.beer_lot_id
		JOIN stock_location sl ON sl.id = m.stock_location_id
		WHERE m.deleted_at IS NULL
		  AND bl.deleted_at IS NULL
		  AND sl.deleted_at IS NULL
		GROUP BY bl.id, bl.uuid, bl.production_batch_uuid, bl.lot_code, bl.packaged_at, bl.best_by,
			bl.package_format_name, bl.container, bl.volume_per_unit, bl.volume_per_unit_unit,
			bl.quantity, sl.id, sl.uuid, sl.name, m.amount_unit
		HAVING SUM(CASE m.direction
			WHEN 'in' THEN m.amount
			WHEN 'out' THEN -m.amount
		END) > 0
		ORDER BY bl.packaged_at DESC, sl.name
	`)
	if err != nil {
		return nil, fmt.Errorf("querying beer lot stock levels: %w", err)
	}
	defer rows.Close()

	var levels []BeerLotStockLevel
	for rows.Next() {
		var level BeerLotStockLevel
		if err := rows.Scan(
			&level.BeerLotUUID,
			&level.ProductionBatchUUID,
			&level.LotCode,
			&level.PackagedAt,
			&level.BestBy,
			&level.PackageFormatName,
			&level.Container,
			&level.VolumePerUnit,
			&level.VolumePerUnitUnit,
			&level.InitialQuantity,
			&level.LocationUUID,
			&level.LocationName,
			&level.CurrentVolume,
			&level.CurrentVolumeUnit,
		); err != nil {
			return nil, fmt.Errorf("scanning beer lot stock level: %w", err)
		}

		// Derive current quantity from current volume / volume per unit when possible.
		if level.VolumePerUnit != nil && *level.VolumePerUnit > 0 {
			qty := int(level.CurrentVolume / *level.VolumePerUnit)
			level.CurrentQuantity = &qty
		}

		levels = append(levels, level)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating beer lot stock level rows: %w", err)
	}

	return levels, nil
}
