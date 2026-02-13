package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateIngredient(ctx context.Context, ingredient Ingredient) (Ingredient, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO ingredient (
			name,
			category,
			default_unit,
			description
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, name, category, default_unit, description, created_at, updated_at, deleted_at`,
		ingredient.Name,
		ingredient.Category,
		ingredient.DefaultUnit,
		ingredient.Description,
	).Scan(
		&ingredient.ID,
		&ingredient.UUID,
		&ingredient.Name,
		&ingredient.Category,
		&ingredient.DefaultUnit,
		&ingredient.Description,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
		&ingredient.DeletedAt,
	)
	if err != nil {
		return Ingredient{}, fmt.Errorf("creating ingredient: %w", err)
	}

	return ingredient, nil
}

func (c *Client) GetIngredient(ctx context.Context, id int64) (Ingredient, error) {
	var ingredient Ingredient
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, name, category, default_unit, description, created_at, updated_at, deleted_at
		FROM ingredient
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&ingredient.ID,
		&ingredient.UUID,
		&ingredient.Name,
		&ingredient.Category,
		&ingredient.DefaultUnit,
		&ingredient.Description,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
		&ingredient.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Ingredient{}, service.ErrNotFound
		}
		return Ingredient{}, fmt.Errorf("getting ingredient: %w", err)
	}

	return ingredient, nil
}

func (c *Client) GetIngredientByUUID(ctx context.Context, ingredientUUID string) (Ingredient, error) {
	var ingredient Ingredient
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, name, category, default_unit, description, created_at, updated_at, deleted_at
		FROM ingredient
		WHERE uuid = $1 AND deleted_at IS NULL`,
		ingredientUUID,
	).Scan(
		&ingredient.ID,
		&ingredient.UUID,
		&ingredient.Name,
		&ingredient.Category,
		&ingredient.DefaultUnit,
		&ingredient.Description,
		&ingredient.CreatedAt,
		&ingredient.UpdatedAt,
		&ingredient.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Ingredient{}, service.ErrNotFound
		}
		return Ingredient{}, fmt.Errorf("getting ingredient by uuid: %w", err)
	}

	return ingredient, nil
}

func (c *Client) ListIngredients(ctx context.Context) ([]Ingredient, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, name, category, default_unit, description, created_at, updated_at, deleted_at
		FROM ingredient
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing ingredients: %w", err)
	}
	defer rows.Close()

	var ingredients []Ingredient
	for rows.Next() {
		var ingredient Ingredient
		if err := rows.Scan(
			&ingredient.ID,
			&ingredient.UUID,
			&ingredient.Name,
			&ingredient.Category,
			&ingredient.DefaultUnit,
			&ingredient.Description,
			&ingredient.CreatedAt,
			&ingredient.UpdatedAt,
			&ingredient.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning ingredient: %w", err)
		}
		ingredients = append(ingredients, ingredient)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing ingredients: %w", err)
	}

	return ingredients, nil
}
