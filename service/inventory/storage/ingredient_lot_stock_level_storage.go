package storage

import (
	"context"
	"fmt"
	"time"
)

// IngredientLotStockLevel represents the current stock level for an ingredient
// lot at a specific location, computed from the inventory movement ledger.
type IngredientLotStockLevel struct {
	IngredientLotUUID  string
	IngredientUUID     string
	IngredientName     string
	IngredientCategory string
	BreweryLotCode     *string
	ReceivedAt         time.Time
	ReceivedAmount     int64
	ReceivedUnit       string
	LocationUUID       string
	LocationName       string
	CurrentAmount      int64
	CurrentUnit        string
}

// GetIngredientLotStockLevels returns current stock levels for all ingredient
// lots, computed from the inventory movement ledger. Stock is grouped by
// ingredient lot and location. Only lots with positive stock are included.
func (c *Client) GetIngredientLotStockLevels(ctx context.Context) ([]IngredientLotStockLevel, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT
			il.uuid AS ingredient_lot_uuid,
			i.uuid AS ingredient_uuid,
			i.name AS ingredient_name,
			i.category AS ingredient_category,
			il.brewery_lot_code,
			il.received_at,
			il.received_amount,
			il.received_unit,
			sl.uuid AS location_uuid,
			sl.name AS location_name,
			SUM(CASE m.direction
				WHEN 'in' THEN m.amount
				WHEN 'out' THEN -m.amount
			END) AS current_amount,
			m.amount_unit AS current_unit
		FROM inventory_movement m
		JOIN ingredient_lot il ON il.id = m.ingredient_lot_id
		JOIN ingredient i ON i.id = il.ingredient_id
		JOIN stock_location sl ON sl.id = m.stock_location_id
		WHERE m.deleted_at IS NULL
		  AND il.deleted_at IS NULL
		  AND i.deleted_at IS NULL
		  AND sl.deleted_at IS NULL
		GROUP BY il.id, il.uuid, i.uuid, i.name, i.category,
			il.brewery_lot_code, il.received_at, il.received_amount, il.received_unit,
			sl.id, sl.uuid, sl.name, m.amount_unit
		HAVING SUM(CASE m.direction
			WHEN 'in' THEN m.amount
			WHEN 'out' THEN -m.amount
		END) > 0
		ORDER BY i.name, il.received_at DESC, sl.name
	`)
	if err != nil {
		return nil, fmt.Errorf("querying ingredient lot stock levels: %w", err)
	}
	defer rows.Close()

	var levels []IngredientLotStockLevel
	for rows.Next() {
		var level IngredientLotStockLevel
		if err := rows.Scan(
			&level.IngredientLotUUID,
			&level.IngredientUUID,
			&level.IngredientName,
			&level.IngredientCategory,
			&level.BreweryLotCode,
			&level.ReceivedAt,
			&level.ReceivedAmount,
			&level.ReceivedUnit,
			&level.LocationUUID,
			&level.LocationName,
			&level.CurrentAmount,
			&level.CurrentUnit,
		); err != nil {
			return nil, fmt.Errorf("scanning ingredient lot stock level: %w", err)
		}
		levels = append(levels, level)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating ingredient lot stock level rows: %w", err)
	}

	return levels, nil
}
