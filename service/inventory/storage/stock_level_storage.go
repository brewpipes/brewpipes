package storage

import (
	"context"
	"fmt"
)

// StockLevelLocation represents the stock level at a specific location.
type StockLevelLocation struct {
	LocationUUID string
	LocationName string
	Quantity     float64
}

// StockLevel represents the aggregated stock level for an ingredient.
type StockLevel struct {
	IngredientUUID string
	IngredientName string
	Category       string
	DefaultUnit    string
	TotalOnHand    float64
	Locations      []StockLevelLocation
}

// stockLevelRow is the raw row returned from the stock levels query.
type stockLevelRow struct {
	IngredientUUID string
	IngredientName string
	Category       string
	DefaultUnit    string
	LocationUUID   string
	LocationName   string
	Quantity       float64
}

// GetStockLevels returns current stock levels for all ingredients, computed
// from the inventory movement ledger. Stock is grouped by ingredient and
// location. Only ingredient lots are included (not beer lots).
func (c *Client) GetStockLevels(ctx context.Context) ([]StockLevel, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT 
			i.uuid as ingredient_uuid,
			i.name as ingredient_name,
			i.category,
			i.default_unit,
			sl.uuid as location_uuid,
			sl.name as location_name,
			SUM(CASE m.direction 
				WHEN 'in' THEN m.amount 
				WHEN 'out' THEN -m.amount 
			END) as quantity
		FROM inventory_movement m
		JOIN ingredient_lot il ON il.id = m.ingredient_lot_id
		JOIN ingredient i ON i.id = il.ingredient_id
		JOIN stock_location sl ON sl.id = m.stock_location_id
		WHERE m.deleted_at IS NULL
		  AND il.deleted_at IS NULL
		  AND i.deleted_at IS NULL
		  AND sl.deleted_at IS NULL
		GROUP BY i.id, i.uuid, i.name, i.category, i.default_unit, sl.id, sl.uuid, sl.name
		ORDER BY i.category, i.name, sl.name
	`)
	if err != nil {
		return nil, fmt.Errorf("querying stock levels: %w", err)
	}
	defer rows.Close()

	// Scan rows into flat structure
	var rawRows []stockLevelRow
	for rows.Next() {
		var row stockLevelRow
		if err := rows.Scan(
			&row.IngredientUUID,
			&row.IngredientName,
			&row.Category,
			&row.DefaultUnit,
			&row.LocationUUID,
			&row.LocationName,
			&row.Quantity,
		); err != nil {
			return nil, fmt.Errorf("scanning stock level row: %w", err)
		}
		rawRows = append(rawRows, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating stock level rows: %w", err)
	}

	// Aggregate rows by ingredient
	return aggregateStockLevels(rawRows), nil
}

// aggregateStockLevels groups flat rows by ingredient and computes totals.
func aggregateStockLevels(rows []stockLevelRow) []StockLevel {
	// Use a map to group by ingredient UUID
	ingredientMap := make(map[string]*StockLevel)
	// Track order for stable output
	var order []string

	for _, row := range rows {
		sl, exists := ingredientMap[row.IngredientUUID]
		if !exists {
			sl = &StockLevel{
				IngredientUUID: row.IngredientUUID,
				IngredientName: row.IngredientName,
				Category:       row.Category,
				DefaultUnit:    row.DefaultUnit,
				TotalOnHand:    0,
				Locations:      []StockLevelLocation{},
			}
			ingredientMap[row.IngredientUUID] = sl
			order = append(order, row.IngredientUUID)
		}

		sl.Locations = append(sl.Locations, StockLevelLocation{
			LocationUUID: row.LocationUUID,
			LocationName: row.LocationName,
			Quantity:     row.Quantity,
		})
		sl.TotalOnHand += row.Quantity
	}

	// Build result in original order (preserves ORDER BY from SQL)
	result := make([]StockLevel, 0, len(order))
	for _, uuid := range order {
		result = append(result, *ingredientMap[uuid])
	}

	return result
}
