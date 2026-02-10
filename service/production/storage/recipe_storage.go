package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateRecipe(ctx context.Context, recipe Recipe) (Recipe, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO recipe (
			name, style_id, style_name, notes,
			batch_size, batch_size_unit,
			target_og, target_og_min, target_og_max,
			target_fg, target_fg_min, target_fg_max,
			target_ibu, target_ibu_min, target_ibu_max,
			target_srm, target_srm_min, target_srm_max,
			target_carbonation, ibu_method, brewhouse_efficiency
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id, uuid, name, style_id, style_name, notes,
		          batch_size, batch_size_unit,
		          target_og, target_og_min, target_og_max,
		          target_fg, target_fg_min, target_fg_max,
		          target_ibu, target_ibu_min, target_ibu_max,
		          target_srm, target_srm_min, target_srm_max,
		          target_carbonation, ibu_method, brewhouse_efficiency,
		          created_at, updated_at, deleted_at`,
		recipe.Name,
		recipe.StyleID,
		recipe.StyleName,
		recipe.Notes,
		recipe.BatchSize,
		recipe.BatchSizeUnit,
		recipe.TargetOG,
		recipe.TargetOGMin,
		recipe.TargetOGMax,
		recipe.TargetFG,
		recipe.TargetFGMin,
		recipe.TargetFGMax,
		recipe.TargetIBU,
		recipe.TargetIBUMin,
		recipe.TargetIBUMax,
		recipe.TargetSRM,
		recipe.TargetSRMMin,
		recipe.TargetSRMMax,
		recipe.TargetCarbonation,
		recipe.IBUMethod,
		recipe.BrewhouseEfficiency,
	).Scan(
		&recipe.ID,
		&recipe.UUID,
		&recipe.Name,
		&recipe.StyleID,
		&recipe.StyleName,
		&recipe.Notes,
		&recipe.BatchSize,
		&recipe.BatchSizeUnit,
		&recipe.TargetOG,
		&recipe.TargetOGMin,
		&recipe.TargetOGMax,
		&recipe.TargetFG,
		&recipe.TargetFGMin,
		&recipe.TargetFGMax,
		&recipe.TargetIBU,
		&recipe.TargetIBUMin,
		&recipe.TargetIBUMax,
		&recipe.TargetSRM,
		&recipe.TargetSRMMin,
		&recipe.TargetSRMMax,
		&recipe.TargetCarbonation,
		&recipe.IBUMethod,
		&recipe.BrewhouseEfficiency,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
		&recipe.DeletedAt,
	)
	if err != nil {
		return Recipe{}, fmt.Errorf("creating recipe: %w", err)
	}

	return recipe, nil
}

// RecipeQueryOpts controls optional query behavior for recipe retrieval.
type RecipeQueryOpts struct {
	// IncludeDeleted allows retrieval of soft-deleted recipes.
	// This is useful for historical references (e.g., batches that reference deleted recipes).
	IncludeDeleted bool
}

func (c *Client) GetRecipe(ctx context.Context, recipeUUID string, opts *RecipeQueryOpts) (Recipe, error) {
	includeDeleted := opts != nil && opts.IncludeDeleted

	query := `
		SELECT id, uuid, name, style_id, style_name, notes,
		       batch_size, batch_size_unit,
		       target_og, target_og_min, target_og_max,
		       target_fg, target_fg_min, target_fg_max,
		       target_ibu, target_ibu_min, target_ibu_max,
		       target_srm, target_srm_min, target_srm_max,
		       target_carbonation, ibu_method, brewhouse_efficiency,
		       created_at, updated_at, deleted_at
		FROM recipe
		WHERE uuid = $1`
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}

	var recipe Recipe
	err := c.db.QueryRow(ctx, query, recipeUUID).Scan(
		&recipe.ID,
		&recipe.UUID,
		&recipe.Name,
		&recipe.StyleID,
		&recipe.StyleName,
		&recipe.Notes,
		&recipe.BatchSize,
		&recipe.BatchSizeUnit,
		&recipe.TargetOG,
		&recipe.TargetOGMin,
		&recipe.TargetOGMax,
		&recipe.TargetFG,
		&recipe.TargetFGMin,
		&recipe.TargetFGMax,
		&recipe.TargetIBU,
		&recipe.TargetIBUMin,
		&recipe.TargetIBUMax,
		&recipe.TargetSRM,
		&recipe.TargetSRMMin,
		&recipe.TargetSRMMax,
		&recipe.TargetCarbonation,
		&recipe.IBUMethod,
		&recipe.BrewhouseEfficiency,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
		&recipe.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Recipe{}, service.ErrNotFound
		}
		return Recipe{}, fmt.Errorf("getting recipe: %w", err)
	}

	return recipe, nil
}

func (c *Client) ListRecipes(ctx context.Context) ([]Recipe, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, name, style_id, style_name, notes,
		       batch_size, batch_size_unit,
		       target_og, target_og_min, target_og_max,
		       target_fg, target_fg_min, target_fg_max,
		       target_ibu, target_ibu_min, target_ibu_max,
		       target_srm, target_srm_min, target_srm_max,
		       target_carbonation, ibu_method, brewhouse_efficiency,
		       created_at, updated_at, deleted_at
		FROM recipe
		WHERE deleted_at IS NULL
		ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing recipes: %w", err)
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		if err := rows.Scan(
			&recipe.ID,
			&recipe.UUID,
			&recipe.Name,
			&recipe.StyleID,
			&recipe.StyleName,
			&recipe.Notes,
			&recipe.BatchSize,
			&recipe.BatchSizeUnit,
			&recipe.TargetOG,
			&recipe.TargetOGMin,
			&recipe.TargetOGMax,
			&recipe.TargetFG,
			&recipe.TargetFGMin,
			&recipe.TargetFGMax,
			&recipe.TargetIBU,
			&recipe.TargetIBUMin,
			&recipe.TargetIBUMax,
			&recipe.TargetSRM,
			&recipe.TargetSRMMin,
			&recipe.TargetSRMMax,
			&recipe.TargetCarbonation,
			&recipe.IBUMethod,
			&recipe.BrewhouseEfficiency,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
			&recipe.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning recipe: %w", err)
		}
		recipes = append(recipes, recipe)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing recipes: %w", err)
	}

	return recipes, nil
}

func (c *Client) UpdateRecipe(ctx context.Context, recipeUUID string, recipe Recipe) (Recipe, error) {
	err := c.db.QueryRow(ctx, `
		UPDATE recipe
		SET name = $1, style_id = $2, style_name = $3, notes = $4,
		    batch_size = $5, batch_size_unit = $6,
		    target_og = $7, target_og_min = $8, target_og_max = $9,
		    target_fg = $10, target_fg_min = $11, target_fg_max = $12,
		    target_ibu = $13, target_ibu_min = $14, target_ibu_max = $15,
		    target_srm = $16, target_srm_min = $17, target_srm_max = $18,
		    target_carbonation = $19, ibu_method = $20, brewhouse_efficiency = $21,
		    updated_at = timezone('utc', now())
		WHERE uuid = $22 AND deleted_at IS NULL
		RETURNING id, uuid, name, style_id, style_name, notes,
		          batch_size, batch_size_unit,
		          target_og, target_og_min, target_og_max,
		          target_fg, target_fg_min, target_fg_max,
		          target_ibu, target_ibu_min, target_ibu_max,
		          target_srm, target_srm_min, target_srm_max,
		          target_carbonation, ibu_method, brewhouse_efficiency,
		          created_at, updated_at, deleted_at`,
		recipe.Name,
		recipe.StyleID,
		recipe.StyleName,
		recipe.Notes,
		recipe.BatchSize,
		recipe.BatchSizeUnit,
		recipe.TargetOG,
		recipe.TargetOGMin,
		recipe.TargetOGMax,
		recipe.TargetFG,
		recipe.TargetFGMin,
		recipe.TargetFGMax,
		recipe.TargetIBU,
		recipe.TargetIBUMin,
		recipe.TargetIBUMax,
		recipe.TargetSRM,
		recipe.TargetSRMMin,
		recipe.TargetSRMMax,
		recipe.TargetCarbonation,
		recipe.IBUMethod,
		recipe.BrewhouseEfficiency,
		recipeUUID,
	).Scan(
		&recipe.ID,
		&recipe.UUID,
		&recipe.Name,
		&recipe.StyleID,
		&recipe.StyleName,
		&recipe.Notes,
		&recipe.BatchSize,
		&recipe.BatchSizeUnit,
		&recipe.TargetOG,
		&recipe.TargetOGMin,
		&recipe.TargetOGMax,
		&recipe.TargetFG,
		&recipe.TargetFGMin,
		&recipe.TargetFGMax,
		&recipe.TargetIBU,
		&recipe.TargetIBUMin,
		&recipe.TargetIBUMax,
		&recipe.TargetSRM,
		&recipe.TargetSRMMin,
		&recipe.TargetSRMMax,
		&recipe.TargetCarbonation,
		&recipe.IBUMethod,
		&recipe.BrewhouseEfficiency,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
		&recipe.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Recipe{}, service.ErrNotFound
		}
		return Recipe{}, fmt.Errorf("updating recipe: %w", err)
	}

	return recipe, nil
}

func (c *Client) DeleteRecipe(ctx context.Context, recipeUUID string) error {
	result, err := c.db.Exec(ctx, `
		UPDATE recipe
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE uuid = $1 AND deleted_at IS NULL`,
		recipeUUID,
	)
	if err != nil {
		return fmt.Errorf("deleting recipe: %w", err)
	}

	if result.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	return nil
}

// getRecipeByID retrieves a recipe by internal ID. This is for internal use only
// (e.g., batch summary lookups where we have the FK ID, not the UUID).
func (c *Client) getRecipeByID(ctx context.Context, id int64, opts *RecipeQueryOpts) (Recipe, error) {
	includeDeleted := opts != nil && opts.IncludeDeleted

	query := `
		SELECT id, uuid, name, style_id, style_name, notes,
		       batch_size, batch_size_unit,
		       target_og, target_og_min, target_og_max,
		       target_fg, target_fg_min, target_fg_max,
		       target_ibu, target_ibu_min, target_ibu_max,
		       target_srm, target_srm_min, target_srm_max,
		       target_carbonation, ibu_method, brewhouse_efficiency,
		       created_at, updated_at, deleted_at
		FROM recipe
		WHERE id = $1`
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}

	var recipe Recipe
	err := c.db.QueryRow(ctx, query, id).Scan(
		&recipe.ID,
		&recipe.UUID,
		&recipe.Name,
		&recipe.StyleID,
		&recipe.StyleName,
		&recipe.Notes,
		&recipe.BatchSize,
		&recipe.BatchSizeUnit,
		&recipe.TargetOG,
		&recipe.TargetOGMin,
		&recipe.TargetOGMax,
		&recipe.TargetFG,
		&recipe.TargetFGMin,
		&recipe.TargetFGMax,
		&recipe.TargetIBU,
		&recipe.TargetIBUMin,
		&recipe.TargetIBUMax,
		&recipe.TargetSRM,
		&recipe.TargetSRMMin,
		&recipe.TargetSRMMax,
		&recipe.TargetCarbonation,
		&recipe.IBUMethod,
		&recipe.BrewhouseEfficiency,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
		&recipe.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Recipe{}, service.ErrNotFound
		}
		return Recipe{}, fmt.Errorf("getting recipe by id: %w", err)
	}

	return recipe, nil
}
