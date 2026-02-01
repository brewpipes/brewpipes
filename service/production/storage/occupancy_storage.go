package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateOccupancy(ctx context.Context, occupancy Occupancy) (Occupancy, error) {
	inAt := occupancy.InAt
	if inAt.IsZero() {
		inAt = time.Now().UTC()
	}

	err := c.db.QueryRow(ctx, `
		INSERT INTO occupancy (
			vessel_id,
			volume_id,
			in_at,
			out_at
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, vessel_id, volume_id, in_at, out_at, created_at, updated_at, deleted_at`,
		occupancy.VesselID,
		occupancy.VolumeID,
		inAt,
		occupancy.OutAt,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VolumeID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.CreatedAt,
		&occupancy.UpdatedAt,
		&occupancy.DeletedAt,
	)
	if err != nil {
		return Occupancy{}, fmt.Errorf("creating occupancy: %w", err)
	}

	return occupancy, nil
}

func (c *Client) GetOccupancy(ctx context.Context, id int64) (Occupancy, error) {
	var occupancy Occupancy
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, vessel_id, volume_id, in_at, out_at, created_at, updated_at, deleted_at
		FROM occupancy
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VolumeID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.CreatedAt,
		&occupancy.UpdatedAt,
		&occupancy.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Occupancy{}, service.ErrNotFound
		}
		return Occupancy{}, fmt.Errorf("getting occupancy: %w", err)
	}

	return occupancy, nil
}

func (c *Client) GetActiveOccupancyByVessel(ctx context.Context, vesselID int64) (Occupancy, error) {
	var occupancy Occupancy
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, vessel_id, volume_id, in_at, out_at, created_at, updated_at, deleted_at
		FROM occupancy
		WHERE vessel_id = $1 AND out_at IS NULL AND deleted_at IS NULL`,
		vesselID,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VolumeID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.CreatedAt,
		&occupancy.UpdatedAt,
		&occupancy.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Occupancy{}, service.ErrNotFound
		}
		return Occupancy{}, fmt.Errorf("getting active occupancy by vessel: %w", err)
	}

	return occupancy, nil
}

func (c *Client) GetActiveOccupancyByVolume(ctx context.Context, volumeID int64) (Occupancy, error) {
	var occupancy Occupancy
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, vessel_id, volume_id, in_at, out_at, created_at, updated_at, deleted_at
		FROM occupancy
		WHERE volume_id = $1 AND out_at IS NULL AND deleted_at IS NULL`,
		volumeID,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VolumeID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.CreatedAt,
		&occupancy.UpdatedAt,
		&occupancy.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Occupancy{}, service.ErrNotFound
		}
		return Occupancy{}, fmt.Errorf("getting active occupancy by volume: %w", err)
	}

	return occupancy, nil
}

func (c *Client) ListActiveOccupancies(ctx context.Context) ([]Occupancy, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, vessel_id, volume_id, in_at, out_at, created_at, updated_at, deleted_at
		FROM occupancy
		WHERE out_at IS NULL AND deleted_at IS NULL
		ORDER BY in_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing active occupancies: %w", err)
	}
	defer rows.Close()

	var occupancies []Occupancy
	for rows.Next() {
		var occupancy Occupancy
		if err := rows.Scan(
			&occupancy.ID,
			&occupancy.UUID,
			&occupancy.VesselID,
			&occupancy.VolumeID,
			&occupancy.InAt,
			&occupancy.OutAt,
			&occupancy.CreatedAt,
			&occupancy.UpdatedAt,
			&occupancy.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning occupancy: %w", err)
		}
		occupancies = append(occupancies, occupancy)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing active occupancies: %w", err)
	}

	return occupancies, nil
}

func (c *Client) CloseOccupancy(ctx context.Context, occupancyID int64, outAt time.Time) error {
	var id int64
	err := c.db.QueryRow(ctx, `
		UPDATE occupancy
		SET out_at = $1, updated_at = timezone('utc', now())
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id`,
		outAt,
		occupancyID,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return service.ErrNotFound
		}
		return fmt.Errorf("closing occupancy: %w", err)
	}

	return nil
}
