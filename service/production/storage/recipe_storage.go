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
			name,
			style_id,
			style_name,
			notes
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, name, style_id, style_name, notes, created_at, updated_at, deleted_at`,
		recipe.Name,
		recipe.StyleID,
		recipe.StyleName,
		recipe.Notes,
	).Scan(
		&recipe.ID,
		&recipe.UUID,
		&recipe.Name,
		&recipe.StyleID,
		&recipe.StyleName,
		&recipe.Notes,
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

func (c *Client) GetRecipe(ctx context.Context, id int64, opts *RecipeQueryOpts) (Recipe, error) {
	includeDeleted := opts != nil && opts.IncludeDeleted

	query := `
		SELECT id, uuid, name, style_id, style_name, notes, created_at, updated_at, deleted_at
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
		SELECT id, uuid, name, style_id, style_name, notes, created_at, updated_at, deleted_at
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

func (c *Client) UpdateRecipe(ctx context.Context, id int64, recipe Recipe) (Recipe, error) {
	err := c.db.QueryRow(ctx, `
		UPDATE recipe
		SET name = $1, style_id = $2, style_name = $3, notes = $4, updated_at = timezone('utc', now())
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING id, uuid, name, style_id, style_name, notes, created_at, updated_at, deleted_at`,
		recipe.Name,
		recipe.StyleID,
		recipe.StyleName,
		recipe.Notes,
		id,
	).Scan(
		&recipe.ID,
		&recipe.UUID,
		&recipe.Name,
		&recipe.StyleID,
		&recipe.StyleName,
		&recipe.Notes,
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

func (c *Client) DeleteRecipe(ctx context.Context, id int64) error {
	result, err := c.db.Exec(ctx, `
		UPDATE recipe
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("deleting recipe: %w", err)
	}

	if result.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	return nil
}
