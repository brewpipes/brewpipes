package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateVessel(ctx context.Context, vessel Vessel) (Vessel, error) {
	status := vessel.Status
	if status == "" {
		status = VesselStatusActive
	}

	err := c.DB().QueryRow(ctx, `
		INSERT INTO vessel (
			type,
			name,
			capacity,
			capacity_unit,
			make,
			model,
			status
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, type, name, capacity, capacity_unit, make, model, status, created_at, updated_at, deleted_at`,
		vessel.Type,
		vessel.Name,
		vessel.Capacity,
		vessel.CapacityUnit,
		vessel.Make,
		vessel.Model,
		status,
	).Scan(
		&vessel.ID,
		&vessel.UUID,
		&vessel.Type,
		&vessel.Name,
		&vessel.Capacity,
		&vessel.CapacityUnit,
		&vessel.Make,
		&vessel.Model,
		&vessel.Status,
		&vessel.CreatedAt,
		&vessel.UpdatedAt,
		&vessel.DeletedAt,
	)
	if err != nil {
		return Vessel{}, fmt.Errorf("creating vessel: %w", err)
	}

	return vessel, nil
}

func (c *Client) GetVessel(ctx context.Context, id int64) (Vessel, error) {
	var vessel Vessel
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, type, name, capacity, capacity_unit, make, model, status, created_at, updated_at, deleted_at
		FROM vessel
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&vessel.ID,
		&vessel.UUID,
		&vessel.Type,
		&vessel.Name,
		&vessel.Capacity,
		&vessel.CapacityUnit,
		&vessel.Make,
		&vessel.Model,
		&vessel.Status,
		&vessel.CreatedAt,
		&vessel.UpdatedAt,
		&vessel.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Vessel{}, service.ErrNotFound
		}
		return Vessel{}, fmt.Errorf("getting vessel: %w", err)
	}

	return vessel, nil
}

func (c *Client) ListVessels(ctx context.Context) ([]Vessel, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT id, uuid, type, name, capacity, capacity_unit, make, model, status, created_at, updated_at, deleted_at
		FROM vessel
		WHERE deleted_at IS NULL
		ORDER BY name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing vessels: %w", err)
	}
	defer rows.Close()

	var vessels []Vessel
	for rows.Next() {
		var vessel Vessel
		if err := rows.Scan(
			&vessel.ID,
			&vessel.UUID,
			&vessel.Type,
			&vessel.Name,
			&vessel.Capacity,
			&vessel.CapacityUnit,
			&vessel.Make,
			&vessel.Model,
			&vessel.Status,
			&vessel.CreatedAt,
			&vessel.UpdatedAt,
			&vessel.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning vessel: %w", err)
		}
		vessels = append(vessels, vessel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing vessels: %w", err)
	}

	return vessels, nil
}

func (c *Client) GetVesselByUUID(ctx context.Context, uuid string) (Vessel, error) {
	var vessel Vessel
	err := c.DB().QueryRow(ctx, `
		SELECT id, uuid, type, name, capacity, capacity_unit, make, model, status, created_at, updated_at, deleted_at
		FROM vessel
		WHERE uuid = $1 AND deleted_at IS NULL`,
		uuid,
	).Scan(
		&vessel.ID,
		&vessel.UUID,
		&vessel.Type,
		&vessel.Name,
		&vessel.Capacity,
		&vessel.CapacityUnit,
		&vessel.Make,
		&vessel.Model,
		&vessel.Status,
		&vessel.CreatedAt,
		&vessel.UpdatedAt,
		&vessel.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Vessel{}, service.ErrNotFound
		}
		return Vessel{}, fmt.Errorf("getting vessel by uuid: %w", err)
	}

	return vessel, nil
}

func (c *Client) UpdateVessel(ctx context.Context, id int64, vessel Vessel) (Vessel, error) {
	err := c.DB().QueryRow(ctx, `
		UPDATE vessel
		SET type = $1, name = $2, capacity = $3, capacity_unit = $4, make = $5, model = $6, status = $7, updated_at = timezone('utc', now())
		WHERE id = $8 AND deleted_at IS NULL
		RETURNING id, uuid, type, name, capacity, capacity_unit, make, model, status, created_at, updated_at, deleted_at`,
		vessel.Type,
		vessel.Name,
		vessel.Capacity,
		vessel.CapacityUnit,
		vessel.Make,
		vessel.Model,
		vessel.Status,
		id,
	).Scan(
		&vessel.ID,
		&vessel.UUID,
		&vessel.Type,
		&vessel.Name,
		&vessel.Capacity,
		&vessel.CapacityUnit,
		&vessel.Make,
		&vessel.Model,
		&vessel.Status,
		&vessel.CreatedAt,
		&vessel.UpdatedAt,
		&vessel.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Vessel{}, service.ErrNotFound
		}
		return Vessel{}, fmt.Errorf("updating vessel: %w", err)
	}

	return vessel, nil
}

func (c *Client) UpdateVesselByUUID(ctx context.Context, vesselUUID string, vessel Vessel) (Vessel, error) {
	err := c.DB().QueryRow(ctx, `
		UPDATE vessel
		SET type = $1, name = $2, capacity = $3, capacity_unit = $4, make = $5, model = $6, status = $7, updated_at = timezone('utc', now())
		WHERE uuid = $8 AND deleted_at IS NULL
		RETURNING id, uuid, type, name, capacity, capacity_unit, make, model, status, created_at, updated_at, deleted_at`,
		vessel.Type,
		vessel.Name,
		vessel.Capacity,
		vessel.CapacityUnit,
		vessel.Make,
		vessel.Model,
		vessel.Status,
		vesselUUID,
	).Scan(
		&vessel.ID,
		&vessel.UUID,
		&vessel.Type,
		&vessel.Name,
		&vessel.Capacity,
		&vessel.CapacityUnit,
		&vessel.Make,
		&vessel.Model,
		&vessel.Status,
		&vessel.CreatedAt,
		&vessel.UpdatedAt,
		&vessel.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Vessel{}, service.ErrNotFound
		}
		return Vessel{}, fmt.Errorf("updating vessel by uuid: %w", err)
	}

	return vessel, nil
}
