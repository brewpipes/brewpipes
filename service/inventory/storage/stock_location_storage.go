package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

// ErrStockLocationHasInventory is returned when a stock location cannot be
// deleted because it has inventory movements referencing it.
var ErrStockLocationHasInventory = fmt.Errorf("stock location has inventory and cannot be deleted")

func (c *Client) CreateStockLocation(ctx context.Context, location StockLocation) (StockLocation, error) {
	err := c.DB().QueryRow(ctx, `
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
	err := c.DB().QueryRow(ctx, `
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

func (c *Client) GetStockLocationByUUID(ctx context.Context, locationUUID string) (StockLocation, error) {
	var location StockLocation
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, name, location_type, description, created_at, updated_at, deleted_at
		FROM stock_location
		WHERE uuid = $1 AND deleted_at IS NULL`,
		locationUUID,
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
		return StockLocation{}, fmt.Errorf("getting stock location by uuid: %w", err)
	}

	return location, nil
}

// UpdateStockLocationRequest describes the optional fields for a partial update.
type UpdateStockLocationRequest struct {
	Name         *string
	LocationType *string
	Description  *string
}

func (c *Client) UpdateStockLocation(ctx context.Context, locationUUID string, req UpdateStockLocationRequest) (StockLocation, error) {
	setClauses := []string{"updated_at = timezone('utc', now())"}
	args := []any{}
	argIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.LocationType != nil {
		setClauses = append(setClauses, fmt.Sprintf("location_type = $%d", argIdx))
		args = append(args, *req.LocationType)
		argIdx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *req.Description)
		argIdx++
	}

	args = append(args, locationUUID)

	query := fmt.Sprintf(`
		UPDATE stock_location
		SET %s
		WHERE uuid = $%d AND deleted_at IS NULL
		RETURNING id, uuid, name, location_type, description, created_at, updated_at, deleted_at`,
		strings.Join(setClauses, ", "), argIdx)

	var location StockLocation
	err := c.DB().QueryRow(ctx, query, args...).Scan(
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
		return StockLocation{}, fmt.Errorf("updating stock location: %w", err)
	}

	return location, nil
}

func (c *Client) DeleteStockLocation(ctx context.Context, locationUUID string) error {
	// Check if any non-deleted inventory movements reference this location.
	var location StockLocation
	err := c.DB().QueryRow(ctx, `
		SELECT id FROM stock_location
		WHERE uuid = $1 AND deleted_at IS NULL`,
		locationUUID,
	).Scan(&location.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return service.ErrNotFound
		}
		return fmt.Errorf("looking up stock location for delete: %w", err)
	}

	var refCount int
	err = c.DB().QueryRow(ctx, `
		SELECT COUNT(*)
		FROM inventory_movement
		WHERE stock_location_id = $1 AND deleted_at IS NULL`,
		location.ID,
	).Scan(&refCount)
	if err != nil {
		return fmt.Errorf("checking stock location references: %w", err)
	}
	if refCount > 0 {
		return ErrStockLocationHasInventory
	}

	tag, err := c.DB().Exec(ctx, `
		UPDATE stock_location
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE id = $1 AND deleted_at IS NULL`,
		location.ID,
	)
	if err != nil {
		return fmt.Errorf("deleting stock location: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	return nil
}

func (c *Client) ListStockLocations(ctx context.Context) ([]StockLocation, error) {
	rows, err := c.DB().Query(ctx, `
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
