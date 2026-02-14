package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateStyle(ctx context.Context, style Style) (Style, error) {
	err := c.DB().QueryRow(ctx, `
		INSERT INTO style (
			name
		) VALUES ($1)
		RETURNING id, uuid, name, created_at, updated_at, deleted_at`,
		style.Name,
	).Scan(
		&style.ID,
		&style.UUID,
		&style.Name,
		&style.CreatedAt,
		&style.UpdatedAt,
		&style.DeletedAt,
	)
	if err != nil {
		return Style{}, fmt.Errorf("creating style: %w", err)
	}

	return style, nil
}

func (c *Client) GetStyle(ctx context.Context, id int64) (Style, error) {
	var style Style
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, name, created_at, updated_at, deleted_at
		FROM style
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&style.ID,
		&style.UUID,
		&style.Name,
		&style.CreatedAt,
		&style.UpdatedAt,
		&style.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Style{}, service.ErrNotFound
		}
		return Style{}, fmt.Errorf("getting style: %w", err)
	}

	return style, nil
}

func (c *Client) GetStyleByUUID(ctx context.Context, styleUUID string) (Style, error) {
	var style Style
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, name, created_at, updated_at, deleted_at
		FROM style
		WHERE uuid = $1 AND deleted_at IS NULL`,
		styleUUID,
	).Scan(
		&style.ID,
		&style.UUID,
		&style.Name,
		&style.CreatedAt,
		&style.UpdatedAt,
		&style.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Style{}, service.ErrNotFound
		}
		return Style{}, fmt.Errorf("getting style by uuid: %w", err)
	}

	return style, nil
}

func (c *Client) GetStyleByName(ctx context.Context, name string) (Style, error) {
	var style Style
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, name, created_at, updated_at, deleted_at
		FROM style
		WHERE lower(name) = lower($1) AND deleted_at IS NULL`,
		name,
	).Scan(
		&style.ID,
		&style.UUID,
		&style.Name,
		&style.CreatedAt,
		&style.UpdatedAt,
		&style.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Style{}, service.ErrNotFound
		}
		return Style{}, fmt.Errorf("getting style by name: %w", err)
	}

	return style, nil
}

func (c *Client) ListStyles(ctx context.Context) ([]Style, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, name, created_at, updated_at, deleted_at
		FROM style
		WHERE deleted_at IS NULL
		ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing styles: %w", err)
	}
	defer rows.Close()

	var styles []Style
	for rows.Next() {
		var style Style
		if err := rows.Scan(
			&style.ID,
			&style.UUID,
			&style.Name,
			&style.CreatedAt,
			&style.UpdatedAt,
			&style.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning style: %w", err)
		}
		styles = append(styles, style)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing styles: %w", err)
	}

	return styles, nil
}
