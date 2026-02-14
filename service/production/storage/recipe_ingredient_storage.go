package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/internal/database"
	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ListRecipeIngredients returns all ingredients for a recipe, ordered by sort_order.
func (c *Client) ListRecipeIngredients(ctx context.Context, recipeUUID string) ([]RecipeIngredient, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT ri.id, ri.uuid, ri.recipe_id, ri.ingredient_uuid, ri.ingredient_type, ri.amount, ri.amount_unit,
		       ri.use_stage, ri.use_type, ri.timing_duration_minutes, ri.timing_temperature_c,
		       ri.alpha_acid_assumed, ri.scaling_factor, ri.sort_order, ri.notes,
		       ri.created_at, ri.updated_at, ri.deleted_at
		FROM recipe_ingredient ri
		JOIN recipe r ON ri.recipe_id = r.id
		WHERE r.uuid = $1 AND ri.deleted_at IS NULL
		ORDER BY ri.sort_order ASC, ri.id ASC`,
		recipeUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing recipe ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []RecipeIngredient
	for rows.Next() {
		ri, err := scanRecipeIngredient(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning recipe ingredient: %w", err)
		}
		ingredients = append(ingredients, ri)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing recipe ingredients: %w", err)
	}

	return ingredients, nil
}

// GetRecipeIngredient returns a recipe ingredient by UUID.
func (c *Client) GetRecipeIngredient(ctx context.Context, ingredientUUID string) (RecipeIngredient, error) {
	row := c.DB().QueryRow(ctx, `
		SELECT id, uuid, recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit,
		       use_stage, use_type, timing_duration_minutes, timing_temperature_c,
		       alpha_acid_assumed, scaling_factor, sort_order, notes,
		       created_at, updated_at, deleted_at
		FROM recipe_ingredient
		WHERE uuid = $1 AND deleted_at IS NULL`,
		ingredientUUID,
	)

	ri, err := scanRecipeIngredientRow(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RecipeIngredient{}, service.ErrNotFound
		}
		return RecipeIngredient{}, fmt.Errorf("getting recipe ingredient: %w", err)
	}

	return ri, nil
}

// CreateRecipeIngredient creates a new recipe ingredient.
func (c *Client) CreateRecipeIngredient(ctx context.Context, ri RecipeIngredient) (RecipeIngredient, error) {
	var ingredientUUID any
	if ri.IngredientUUID != nil {
		ingredientUUID = *ri.IngredientUUID
	}

	var ingredientUUIDResult pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		INSERT INTO recipe_ingredient (
			recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit,
			use_stage, use_type, timing_duration_minutes, timing_temperature_c,
			alpha_acid_assumed, scaling_factor, sort_order, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, uuid, recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit,
		          use_stage, use_type, timing_duration_minutes, timing_temperature_c,
		          alpha_acid_assumed, scaling_factor, sort_order, notes,
		          created_at, updated_at, deleted_at`,
		ri.RecipeID,
		ingredientUUID,
		ri.IngredientType,
		ri.Amount,
		ri.AmountUnit,
		ri.UseStage,
		ri.UseType,
		ri.TimingDurationMinutes,
		ri.TimingTemperatureC,
		ri.AlphaAcidAssumed,
		ri.ScalingFactor,
		ri.SortOrder,
		ri.Notes,
	).Scan(
		&ri.ID,
		&ri.UUID,
		&ri.RecipeID,
		&ingredientUUIDResult,
		&ri.IngredientType,
		&ri.Amount,
		&ri.AmountUnit,
		&ri.UseStage,
		&ri.UseType,
		&ri.TimingDurationMinutes,
		&ri.TimingTemperatureC,
		&ri.AlphaAcidAssumed,
		&ri.ScalingFactor,
		&ri.SortOrder,
		&ri.Notes,
		&ri.CreatedAt,
		&ri.UpdatedAt,
		&ri.DeletedAt,
	)
	if err != nil {
		return RecipeIngredient{}, fmt.Errorf("creating recipe ingredient: %w", err)
	}

	database.AssignUUIDPointer(&ri.IngredientUUID, ingredientUUIDResult)
	return ri, nil
}

// UpdateRecipeIngredient updates an existing recipe ingredient by UUID.
func (c *Client) UpdateRecipeIngredient(ctx context.Context, ingredientUUID string, ri RecipeIngredient) (RecipeIngredient, error) {
	var invIngredientUUID any
	if ri.IngredientUUID != nil {
		invIngredientUUID = *ri.IngredientUUID
	}

	var invIngredientUUIDResult pgtype.UUID
	err := c.DB().QueryRow(ctx, `
		UPDATE recipe_ingredient
		SET ingredient_uuid = $1,
		    ingredient_type = $2,
		    amount = $3,
		    amount_unit = $4,
		    use_stage = $5,
		    use_type = $6,
		    timing_duration_minutes = $7,
		    timing_temperature_c = $8,
		    alpha_acid_assumed = $9,
		    scaling_factor = $10,
		    sort_order = $11,
		    notes = $12,
		    updated_at = timezone('utc', now())
		WHERE uuid = $13 AND deleted_at IS NULL
		RETURNING id, uuid, recipe_id, ingredient_uuid, ingredient_type, amount, amount_unit,
		          use_stage, use_type, timing_duration_minutes, timing_temperature_c,
		          alpha_acid_assumed, scaling_factor, sort_order, notes,
		          created_at, updated_at, deleted_at`,
		invIngredientUUID,
		ri.IngredientType,
		ri.Amount,
		ri.AmountUnit,
		ri.UseStage,
		ri.UseType,
		ri.TimingDurationMinutes,
		ri.TimingTemperatureC,
		ri.AlphaAcidAssumed,
		ri.ScalingFactor,
		ri.SortOrder,
		ri.Notes,
		ingredientUUID,
	).Scan(
		&ri.ID,
		&ri.UUID,
		&ri.RecipeID,
		&invIngredientUUIDResult,
		&ri.IngredientType,
		&ri.Amount,
		&ri.AmountUnit,
		&ri.UseStage,
		&ri.UseType,
		&ri.TimingDurationMinutes,
		&ri.TimingTemperatureC,
		&ri.AlphaAcidAssumed,
		&ri.ScalingFactor,
		&ri.SortOrder,
		&ri.Notes,
		&ri.CreatedAt,
		&ri.UpdatedAt,
		&ri.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RecipeIngredient{}, service.ErrNotFound
		}
		return RecipeIngredient{}, fmt.Errorf("updating recipe ingredient: %w", err)
	}

	database.AssignUUIDPointer(&ri.IngredientUUID, invIngredientUUIDResult)
	return ri, nil
}

// DeleteRecipeIngredient soft-deletes a recipe ingredient by UUID.
func (c *Client) DeleteRecipeIngredient(ctx context.Context, ingredientUUID string) error {
	result, err := c.DB().Exec(ctx, `
		UPDATE recipe_ingredient
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE uuid = $1 AND deleted_at IS NULL`,
		ingredientUUID,
	)
	if err != nil {
		return fmt.Errorf("deleting recipe ingredient: %w", err)
	}

	if result.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	return nil
}

// scanRecipeIngredient scans a recipe ingredient from a pgx.Rows.
func scanRecipeIngredient(rows pgx.Rows) (RecipeIngredient, error) {
	var ri RecipeIngredient
	var ingredientUUID pgtype.UUID

	err := rows.Scan(
		&ri.ID,
		&ri.UUID,
		&ri.RecipeID,
		&ingredientUUID,
		&ri.IngredientType,
		&ri.Amount,
		&ri.AmountUnit,
		&ri.UseStage,
		&ri.UseType,
		&ri.TimingDurationMinutes,
		&ri.TimingTemperatureC,
		&ri.AlphaAcidAssumed,
		&ri.ScalingFactor,
		&ri.SortOrder,
		&ri.Notes,
		&ri.CreatedAt,
		&ri.UpdatedAt,
		&ri.DeletedAt,
	)
	if err != nil {
		return RecipeIngredient{}, err
	}

	database.AssignUUIDPointer(&ri.IngredientUUID, ingredientUUID)
	return ri, nil
}

// scanRecipeIngredientRow scans a recipe ingredient from a pgx.Row.
func scanRecipeIngredientRow(row pgx.Row) (RecipeIngredient, error) {
	var ri RecipeIngredient
	var ingredientUUID pgtype.UUID

	err := row.Scan(
		&ri.ID,
		&ri.UUID,
		&ri.RecipeID,
		&ingredientUUID,
		&ri.IngredientType,
		&ri.Amount,
		&ri.AmountUnit,
		&ri.UseStage,
		&ri.UseType,
		&ri.TimingDurationMinutes,
		&ri.TimingTemperatureC,
		&ri.AlphaAcidAssumed,
		&ri.ScalingFactor,
		&ri.SortOrder,
		&ri.Notes,
		&ri.CreatedAt,
		&ri.UpdatedAt,
		&ri.DeletedAt,
	)
	if err != nil {
		return RecipeIngredient{}, err
	}

	database.AssignUUIDPointer(&ri.IngredientUUID, ingredientUUID)
	return ri, nil
}
