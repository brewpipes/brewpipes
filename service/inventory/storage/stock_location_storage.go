package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateStockLocation(ctx context.Context, location StockLocation) (StockLocation, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO stock_location (
			name,
			location_type,
			description
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, name, location_type, description, created_at, updated_at, deleted_at`,
		location.Name,
		location.LocationType,
		location.Description,
	).Scan(
		&location.ID,
		&location.UUID,
		&location.Name,
		&location.LocationType,
		&location.Description,
		&location.CreatedAt,
		&location.UpdatedAt,
		&location.DeletedAt,
	)
	if err != nil {
		return StockLocation{}, fmt.Errorf("creating stock location: %w", err)
	}

	return location, nil
}

func (c *Client) GetStockLocation(ctx context.Context, id int64) (StockLocation, error) {
	var location StockLocation
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, name, location_type, description, created_at, updated_at, deleted_at
		FROM stock_location
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&location.ID,
		&location.UUID,
		&location.Name,
		&location.LocationType,
		&location.Description,
		&location.CreatedAt,
		&location.UpdatedAt,
		&location.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return StockLocation{}, service.ErrNotFound
		}
		return StockLocation{}, fmt.Errorf("getting stock location: %w", err)
	}

	return location, nil
}

func (c *Client) ListStockLocations(ctx context.Context) ([]StockLocation, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, name, location_type, description, created_at, updated_at, deleted_at
		FROM stock_location
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing stock locations: %w", err)
	}
	defer rows.Close()

	var locations []StockLocation
	for rows.Next() {
		var location StockLocation
		if err := rows.Scan(
			&location.ID,
			&location.UUID,
			&location.Name,
			&location.LocationType,
			&location.Description,
			&location.CreatedAt,
			&location.UpdatedAt,
			&location.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning stock location: %w", err)
		}
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing stock locations: %w", err)
	}

	return locations, nil
}
