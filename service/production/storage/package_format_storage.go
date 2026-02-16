package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

// ErrPackageFormatInUse is returned when a package format cannot be deleted
// because it is referenced by one or more non-deleted packaging run lines.
var ErrPackageFormatInUse = fmt.Errorf("package format is in use and cannot be deleted")

func (c *Client) CreatePackageFormat(ctx context.Context, format PackageFormat) (PackageFormat, error) {
	err := c.DB().QueryRow(ctx, `
		INSERT INTO package_format (
			name,
			container,
			volume_per_unit,
			volume_per_unit_unit,
			is_active
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, uuid, name, container, volume_per_unit, volume_per_unit_unit, is_active, created_at, updated_at, deleted_at`,
		format.Name,
		format.Container,
		format.VolumePerUnit,
		format.VolumePerUnitUnit,
		format.IsActive,
	).Scan(
		&format.ID,
		&format.UUID,
		&format.Name,
		&format.Container,
		&format.VolumePerUnit,
		&format.VolumePerUnitUnit,
		&format.IsActive,
		&format.CreatedAt,
		&format.UpdatedAt,
		&format.DeletedAt,
	)
	if err != nil {
		return PackageFormat{}, fmt.Errorf("creating package format: %w", err)
	}

	return format, nil
}

func (c *Client) GetPackageFormatByUUID(ctx context.Context, formatUUID string) (PackageFormat, error) {
	var format PackageFormat
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, name, container, volume_per_unit, volume_per_unit_unit, is_active, created_at, updated_at, deleted_at
		FROM package_format
		WHERE uuid = $1 AND deleted_at IS NULL`,
		formatUUID,
	).Scan(
		&format.ID,
		&format.UUID,
		&format.Name,
		&format.Container,
		&format.VolumePerUnit,
		&format.VolumePerUnitUnit,
		&format.IsActive,
		&format.CreatedAt,
		&format.UpdatedAt,
		&format.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PackageFormat{}, service.ErrNotFound
		}
		return PackageFormat{}, fmt.Errorf("getting package format by uuid: %w", err)
	}

	return format, nil
}

func (c *Client) ListPackageFormats(ctx context.Context) ([]PackageFormat, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, name, container, volume_per_unit, volume_per_unit_unit, is_active, created_at, updated_at, deleted_at
		FROM package_format
		WHERE deleted_at IS NULL
		ORDER BY container ASC, name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing package formats: %w", err)
	}
	defer rows.Close()

	var formats []PackageFormat
	for rows.Next() {
		var format PackageFormat
		if err := rows.Scan(
			&format.ID,
			&format.UUID,
			&format.Name,
			&format.Container,
			&format.VolumePerUnit,
			&format.VolumePerUnitUnit,
			&format.IsActive,
			&format.CreatedAt,
			&format.UpdatedAt,
			&format.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning package format: %w", err)
		}
		formats = append(formats, format)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing package formats: %w", err)
	}

	return formats, nil
}

func (c *Client) UpdatePackageFormat(ctx context.Context, id int64, format PackageFormat) (PackageFormat, error) {
	err := c.DB().QueryRow(ctx, `
		UPDATE package_format
		SET name = $1, container = $2, volume_per_unit = $3, volume_per_unit_unit = $4, is_active = $5, updated_at = timezone('utc', now())
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, uuid, name, container, volume_per_unit, volume_per_unit_unit, is_active, created_at, updated_at, deleted_at`,
		format.Name,
		format.Container,
		format.VolumePerUnit,
		format.VolumePerUnitUnit,
		format.IsActive,
		id,
	).Scan(
		&format.ID,
		&format.UUID,
		&format.Name,
		&format.Container,
		&format.VolumePerUnit,
		&format.VolumePerUnitUnit,
		&format.IsActive,
		&format.CreatedAt,
		&format.UpdatedAt,
		&format.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PackageFormat{}, service.ErrNotFound
		}
		return PackageFormat{}, fmt.Errorf("updating package format: %w", err)
	}

	return format, nil
}

func (c *Client) DeletePackageFormat(ctx context.Context, id int64) error {
	// Check if any non-deleted packaging_run_line references this format
	var refCount int
	err := c.DB().QueryRow(ctx, `
		SELECT COUNT(*)
		FROM packaging_run_line
		WHERE package_format_id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&refCount)
	if err != nil {
		return fmt.Errorf("checking package format references: %w", err)
	}
	if refCount > 0 {
		return ErrPackageFormatInUse
	}

	tag, err := c.DB().Exec(ctx, `
		UPDATE package_format
		SET deleted_at = timezone('utc', now()), updated_at = timezone('utc', now())
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("deleting package format: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	return nil
}
