package storage

import (
	"context"
	"fmt"
)

func (c *Client) CreatePackagingRunLine(ctx context.Context, line PackagingRunLine) (PackagingRunLine, error) {
	err := c.DB().QueryRow(ctx, `
		INSERT INTO packaging_run_line (
			packaging_run_id,
			package_format_id,
			quantity
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, packaging_run_id, package_format_id, quantity, created_at, updated_at, deleted_at`,
		line.PackagingRunID,
		line.PackageFormatID,
		line.Quantity,
	).Scan(
		&line.ID,
		&line.UUID,
		&line.PackagingRunID,
		&line.PackageFormatID,
		&line.Quantity,
		&line.CreatedAt,
		&line.UpdatedAt,
		&line.DeletedAt,
	)
	if err != nil {
		return PackagingRunLine{}, fmt.Errorf("creating packaging run line: %w", err)
	}

	// Resolve joined fields
	c.DB().QueryRow(ctx, `SELECT uuid FROM packaging_run WHERE id = $1`, line.PackagingRunID).Scan(&line.PackagingRunUUID)
	c.DB().QueryRow(ctx, `
		SELECT uuid, name, volume_per_unit, volume_per_unit_unit
		FROM package_format WHERE id = $1`,
		line.PackageFormatID,
	).Scan(
		&line.PackageFormatUUID,
		&line.PackageFormatName,
		&line.PackageFormatVolumePerUnit,
		&line.PackageFormatVolumePerUnitUnit,
	)

	return line, nil
}

func (c *Client) ListPackagingRunLinesByRunID(ctx context.Context, runID int64) ([]PackagingRunLine, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT prl.id, prl.uuid, prl.packaging_run_id, pr.uuid,
		       prl.package_format_id, pf.uuid, pf.name, pf.volume_per_unit, pf.volume_per_unit_unit,
		       prl.quantity, prl.created_at, prl.updated_at, prl.deleted_at
		FROM packaging_run_line prl
		JOIN packaging_run pr ON pr.id = prl.packaging_run_id
		JOIN package_format pf ON pf.id = prl.package_format_id
		WHERE prl.packaging_run_id = $1 AND prl.deleted_at IS NULL
		ORDER BY prl.id ASC`,
		runID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing packaging run lines by run id: %w", err)
	}
	defer rows.Close()

	var lines []PackagingRunLine
	for rows.Next() {
		var line PackagingRunLine
		if err := rows.Scan(
			&line.ID,
			&line.UUID,
			&line.PackagingRunID,
			&line.PackagingRunUUID,
			&line.PackageFormatID,
			&line.PackageFormatUUID,
			&line.PackageFormatName,
			&line.PackageFormatVolumePerUnit,
			&line.PackageFormatVolumePerUnitUnit,
			&line.Quantity,
			&line.CreatedAt,
			&line.UpdatedAt,
			&line.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning packaging run line: %w", err)
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing packaging run lines by run id: %w", err)
	}

	return lines, nil
}

func (c *Client) ListPackagingRunLinesByRunUUID(ctx context.Context, runUUID string) ([]PackagingRunLine, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT prl.id, prl.uuid, prl.packaging_run_id, pr.uuid,
		       prl.package_format_id, pf.uuid, pf.name, pf.volume_per_unit, pf.volume_per_unit_unit,
		       prl.quantity, prl.created_at, prl.updated_at, prl.deleted_at
		FROM packaging_run_line prl
		JOIN packaging_run pr ON pr.id = prl.packaging_run_id
		JOIN package_format pf ON pf.id = prl.package_format_id
		WHERE pr.uuid = $1 AND prl.deleted_at IS NULL
		ORDER BY prl.id ASC`,
		runUUID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing packaging run lines by run uuid: %w", err)
	}
	defer rows.Close()

	var lines []PackagingRunLine
	for rows.Next() {
		var line PackagingRunLine
		if err := rows.Scan(
			&line.ID,
			&line.UUID,
			&line.PackagingRunID,
			&line.PackagingRunUUID,
			&line.PackageFormatID,
			&line.PackageFormatUUID,
			&line.PackageFormatName,
			&line.PackageFormatVolumePerUnit,
			&line.PackageFormatVolumePerUnitUnit,
			&line.Quantity,
			&line.CreatedAt,
			&line.UpdatedAt,
			&line.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning packaging run line: %w", err)
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing packaging run lines by run uuid: %w", err)
	}

	return lines, nil
}
