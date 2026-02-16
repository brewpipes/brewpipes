package storage

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/brewpipes/brewpipes/internal/database"
)

// BatchIngredientLot holds lot and ingredient data needed for batch cost calculation.
type BatchIngredientLot struct {
	IngredientLotUUID     string
	IngredientUUID        string
	IngredientName        string
	IngredientCategory    string
	BreweryLotCode        *string
	PurchaseOrderLineUUID *uuid.UUID
	ReceivedUnit          string
}

// ListBatchIngredientLots returns all ingredient lots consumed by a specific production batch.
func (c *Client) ListBatchIngredientLots(ctx context.Context, productionRefUUID string) ([]BatchIngredientLot, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT DISTINCT ON (il.uuid)
			il.uuid AS ingredient_lot_uuid,
			i.uuid AS ingredient_uuid,
			i.name AS ingredient_name,
			i.category AS ingredient_category,
			il.brewery_lot_code,
			il.purchase_order_line_uuid,
			il.received_unit
		FROM inventory_usage iu
		JOIN inventory_movement im ON im.usage_id = iu.id
		JOIN ingredient_lot il ON il.id = im.ingredient_lot_id
		JOIN ingredient i ON i.id = il.ingredient_id
		WHERE iu.production_ref_uuid = $1
		  AND iu.deleted_at IS NULL
		  AND im.deleted_at IS NULL
		  AND il.deleted_at IS NULL
		  AND im.direction = 'out'
		  AND im.reason = 'use'
		ORDER BY il.uuid, im.occurred_at ASC`,
		productionRefUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing batch ingredient lots: %w", err)
	}
	defer rows.Close()

	var lots []BatchIngredientLot
	for rows.Next() {
		var lot BatchIngredientLot
		var lotUUID, ingredientUUID pgtype.UUID
		var poLineUUID pgtype.UUID
		if err := rows.Scan(
			&lotUUID,
			&ingredientUUID,
			&lot.IngredientName,
			&lot.IngredientCategory,
			&lot.BreweryLotCode,
			&poLineUUID,
			&lot.ReceivedUnit,
		); err != nil {
			return nil, fmt.Errorf("scanning batch ingredient lot: %w", err)
		}
		if lotUUID.Valid {
			lot.IngredientLotUUID = uuid.UUID(lotUUID.Bytes).String()
		}
		if ingredientUUID.Valid {
			lot.IngredientUUID = uuid.UUID(ingredientUUID.Bytes).String()
		}
		database.AssignUUIDPointer(&lot.PurchaseOrderLineUUID, poLineUUID)
		lots = append(lots, lot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing batch ingredient lots: %w", err)
	}

	return lots, nil
}
