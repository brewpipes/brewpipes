package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

// ErrOccupancyAlreadyClosed is returned when attempting to close an occupancy
// that already has an out_at timestamp set.
var ErrOccupancyAlreadyClosed = fmt.Errorf("occupancy is already closed")

func (c *Client) CreateOccupancy(ctx context.Context, occupancy Occupancy) (Occupancy, error) {
	inAt := occupancy.InAt
	if inAt.IsZero() {
		inAt = time.Now().UTC()
	}

	err := c.DB().QueryRow(ctx, `
		INSERT INTO occupancy (
			vessel_id,
			volume_id,
			in_at,
			out_at,
			status
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, uuid, vessel_id, volume_id, in_at, out_at, status, created_at, updated_at, deleted_at`,
		occupancy.VesselID,
		occupancy.VolumeID,
		inAt,
		occupancy.OutAt,
		occupancy.Status,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VolumeID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.Status,
		&occupancy.CreatedAt,
		&occupancy.UpdatedAt,
		&occupancy.DeletedAt,
	)
	if err != nil {
		return Occupancy{}, fmt.Errorf("creating occupancy: %w", err)
	}

	// Resolve vessel and volume UUIDs
	c.DB().QueryRow(ctx, `SELECT uuid FROM vessel WHERE id = $1`, occupancy.VesselID).Scan(&occupancy.VesselUUID)
	c.DB().QueryRow(ctx, `SELECT uuid FROM volume WHERE id = $1`, occupancy.VolumeID).Scan(&occupancy.VolumeUUID)

	return occupancy, nil
}

func (c *Client) GetOccupancy(ctx context.Context, id int64) (Occupancy, error) {
	var occupancy Occupancy
	err := c.DB().QueryRow(ctx, `
		SELECT o.id, o.uuid, o.vessel_id, ve.uuid, o.volume_id, vo.uuid,
		       o.in_at, o.out_at, o.status, o.created_at, o.updated_at, o.deleted_at
		FROM occupancy o
		JOIN vessel ve ON ve.id = o.vessel_id
		JOIN volume vo ON vo.id = o.volume_id
		WHERE o.id = $1 AND o.deleted_at IS NULL`,
		id,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VesselUUID,
		&occupancy.VolumeID,
		&occupancy.VolumeUUID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.Status,
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

func (c *Client) GetOccupancyByUUID(ctx context.Context, occupancyUUID string) (Occupancy, error) {
	var occupancy Occupancy
	err := c.DB().QueryRow(ctx, `
		SELECT o.id, o.uuid, o.vessel_id, ve.uuid, o.volume_id, vo.uuid,
		       o.in_at, o.out_at, o.status, o.created_at, o.updated_at, o.deleted_at
		FROM occupancy o
		JOIN vessel ve ON ve.id = o.vessel_id
		JOIN volume vo ON vo.id = o.volume_id
		WHERE o.uuid = $1 AND o.deleted_at IS NULL`,
		occupancyUUID,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VesselUUID,
		&occupancy.VolumeID,
		&occupancy.VolumeUUID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.Status,
		&occupancy.CreatedAt,
		&occupancy.UpdatedAt,
		&occupancy.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Occupancy{}, service.ErrNotFound
		}
		return Occupancy{}, fmt.Errorf("getting occupancy by uuid: %w", err)
	}

	return occupancy, nil
}

func (c *Client) GetActiveOccupancyByVessel(ctx context.Context, vesselID int64) (Occupancy, error) {
	var occupancy Occupancy
	err := c.DB().QueryRow(ctx, `
		SELECT o.id, o.uuid, o.vessel_id, ve.uuid, o.volume_id, vo.uuid,
		       o.in_at, o.out_at, o.status, o.created_at, o.updated_at, o.deleted_at
		FROM occupancy o
		JOIN vessel ve ON ve.id = o.vessel_id
		JOIN volume vo ON vo.id = o.volume_id
		WHERE o.vessel_id = $1 AND o.out_at IS NULL AND o.deleted_at IS NULL`,
		vesselID,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VesselUUID,
		&occupancy.VolumeID,
		&occupancy.VolumeUUID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.Status,
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

func (c *Client) GetActiveOccupancyByVesselUUID(ctx context.Context, vesselUUID string) (Occupancy, error) {
	vessel, err := c.GetVesselByUUID(ctx, vesselUUID)
	if err != nil {
		return Occupancy{}, fmt.Errorf("resolving vessel uuid: %w", err)
	}

	return c.GetActiveOccupancyByVessel(ctx, vessel.ID)
}

func (c *Client) GetActiveOccupancyByVolume(ctx context.Context, volumeID int64) (Occupancy, error) {
	var occupancy Occupancy
	err := c.DB().QueryRow(ctx, `
		SELECT o.id, o.uuid, o.vessel_id, ve.uuid, o.volume_id, vo.uuid,
		       o.in_at, o.out_at, o.status, o.created_at, o.updated_at, o.deleted_at
		FROM occupancy o
		JOIN vessel ve ON ve.id = o.vessel_id
		JOIN volume vo ON vo.id = o.volume_id
		WHERE o.volume_id = $1 AND o.out_at IS NULL AND o.deleted_at IS NULL`,
		volumeID,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VesselUUID,
		&occupancy.VolumeID,
		&occupancy.VolumeUUID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.Status,
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

func (c *Client) GetActiveOccupancyByVolumeUUID(ctx context.Context, volumeUUID string) (Occupancy, error) {
	vol, err := c.GetVolumeByUUID(ctx, volumeUUID)
	if err != nil {
		return Occupancy{}, fmt.Errorf("resolving volume uuid: %w", err)
	}

	return c.GetActiveOccupancyByVolume(ctx, vol.ID)
}

func (c *Client) ListActiveOccupancies(ctx context.Context) ([]Occupancy, error) {
	// Use a subquery to get the most recent batch_id for each volume.
	// A volume can have multiple batch_volume records (e.g., wort -> beer phase transitions),
	// so we select the one with the latest phase_at to avoid duplicate occupancy rows.
	rows, err := c.DB().Query(ctx, `
		SELECT o.id, o.uuid, o.vessel_id, ve.uuid, o.volume_id, vo.uuid,
		       o.in_at, o.out_at, o.status,
		       o.created_at, o.updated_at, o.deleted_at,
		       (SELECT bv.batch_id
		        FROM batch_volume bv
		        WHERE bv.volume_id = o.volume_id AND bv.deleted_at IS NULL
		        ORDER BY bv.phase_at DESC
		        LIMIT 1) AS batch_id,
		       (SELECT b.uuid
		        FROM batch_volume bv
		        JOIN batch b ON b.id = bv.batch_id
		        WHERE bv.volume_id = o.volume_id AND bv.deleted_at IS NULL
		        ORDER BY bv.phase_at DESC
		        LIMIT 1) AS batch_uuid
		FROM occupancy o
		JOIN vessel ve ON ve.id = o.vessel_id
		JOIN volume vo ON vo.id = o.volume_id
		WHERE o.out_at IS NULL AND o.deleted_at IS NULL
		ORDER BY o.in_at DESC`,
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
			&occupancy.VesselUUID,
			&occupancy.VolumeID,
			&occupancy.VolumeUUID,
			&occupancy.InAt,
			&occupancy.OutAt,
			&occupancy.Status,
			&occupancy.CreatedAt,
			&occupancy.UpdatedAt,
			&occupancy.DeletedAt,
			&occupancy.BatchID,
			&occupancy.BatchUUID,
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
	err := c.DB().QueryRow(ctx, `
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

func (c *Client) CloseOccupancyByUUID(ctx context.Context, occupancyUUID string, outAt time.Time) (Occupancy, error) {
	occ, err := c.GetOccupancyByUUID(ctx, occupancyUUID)
	if err != nil {
		return Occupancy{}, err
	}

	if occ.OutAt != nil {
		return Occupancy{}, ErrOccupancyAlreadyClosed
	}

	if err := c.CloseOccupancy(ctx, occ.ID, outAt); err != nil {
		return Occupancy{}, fmt.Errorf("closing occupancy by uuid: %w", err)
	}

	closed, err := c.GetOccupancyByUUID(ctx, occupancyUUID)
	if err != nil {
		return Occupancy{}, fmt.Errorf("re-fetching closed occupancy: %w", err)
	}

	return closed, nil
}

func (c *Client) HasActiveOccupancy(ctx context.Context, vesselID int64) (bool, error) {
	var exists bool
	err := c.DB().QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM occupancy
			WHERE vessel_id = $1 AND out_at IS NULL AND deleted_at IS NULL
		)`,
		vesselID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking active occupancy: %w", err)
	}

	return exists, nil
}

func (c *Client) HasActiveOccupancyByVesselUUID(ctx context.Context, vesselUUID string) (bool, error) {
	vessel, err := c.GetVesselByUUID(ctx, vesselUUID)
	if err != nil {
		return false, fmt.Errorf("resolving vessel uuid: %w", err)
	}

	return c.HasActiveOccupancy(ctx, vessel.ID)
}

func (c *Client) UpdateOccupancyStatus(ctx context.Context, occupancyID int64, status *string) (Occupancy, error) {
	var occupancy Occupancy
	err := c.DB().QueryRow(ctx, `
		UPDATE occupancy
		SET status = $1, updated_at = timezone('utc', now())
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING id, uuid, vessel_id, volume_id, in_at, out_at, status, created_at, updated_at, deleted_at`,
		status,
		occupancyID,
	).Scan(
		&occupancy.ID,
		&occupancy.UUID,
		&occupancy.VesselID,
		&occupancy.VolumeID,
		&occupancy.InAt,
		&occupancy.OutAt,
		&occupancy.Status,
		&occupancy.CreatedAt,
		&occupancy.UpdatedAt,
		&occupancy.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Occupancy{}, service.ErrNotFound
		}
		return Occupancy{}, fmt.Errorf("updating occupancy status: %w", err)
	}

	// Resolve vessel and volume UUIDs
	c.DB().QueryRow(ctx, `SELECT uuid FROM vessel WHERE id = $1`, occupancy.VesselID).Scan(&occupancy.VesselUUID)
	c.DB().QueryRow(ctx, `SELECT uuid FROM volume WHERE id = $1`, occupancy.VolumeID).Scan(&occupancy.VolumeUUID)

	return occupancy, nil
}

func (c *Client) UpdateOccupancyStatusByUUID(ctx context.Context, occupancyUUID string, status *string) (Occupancy, error) {
	occ, err := c.GetOccupancyByUUID(ctx, occupancyUUID)
	if err != nil {
		return Occupancy{}, err
	}

	return c.UpdateOccupancyStatus(ctx, occ.ID, status)
}
