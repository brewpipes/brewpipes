package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateBrewSession(ctx context.Context, session BrewSession) (BrewSession, error) {
	brewedAt := session.BrewedAt
	if brewedAt.IsZero() {
		brewedAt = time.Now().UTC()
	}

	err := c.DB().QueryRow(ctx, `
		INSERT INTO brew_session (
			batch_id,
			wort_volume_id,
			mash_vessel_id,
			boil_vessel_id,
			brewed_at,
			notes
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, uuid, batch_id, wort_volume_id, mash_vessel_id, boil_vessel_id, brewed_at, notes, created_at, updated_at, deleted_at`,
		session.BatchID,
		session.WortVolumeID,
		session.MashVesselID,
		session.BoilVesselID,
		brewedAt,
		session.Notes,
	).Scan(
		&session.ID,
		&session.UUID,
		&session.BatchID,
		&session.WortVolumeID,
		&session.MashVesselID,
		&session.BoilVesselID,
		&session.BrewedAt,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.DeletedAt,
	)
	if err != nil {
		return BrewSession{}, fmt.Errorf("creating brew session: %w", err)
	}

	// Resolve FK UUIDs
	if session.BatchID != nil {
		var batchUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, *session.BatchID).Scan(&batchUUID); err == nil {
			session.BatchUUID = &batchUUID
		}
	}
	if session.WortVolumeID != nil {
		var volUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM volume WHERE id = $1`, *session.WortVolumeID).Scan(&volUUID); err == nil {
			session.WortVolumeUUID = &volUUID
		}
	}
	if session.MashVesselID != nil {
		var vesselUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM vessel WHERE id = $1`, *session.MashVesselID).Scan(&vesselUUID); err == nil {
			session.MashVesselUUID = &vesselUUID
		}
	}
	if session.BoilVesselID != nil {
		var vesselUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM vessel WHERE id = $1`, *session.BoilVesselID).Scan(&vesselUUID); err == nil {
			session.BoilVesselUUID = &vesselUUID
		}
	}

	return session, nil
}

const brewSessionSelectWithJoins = `
	SELECT bs.id, bs.uuid, bs.batch_id, ba.uuid, bs.wort_volume_id, wv.uuid,
	       bs.mash_vessel_id, mv.uuid, bs.boil_vessel_id, bv.uuid,
	       bs.brewed_at, bs.notes, bs.created_at, bs.updated_at, bs.deleted_at
	FROM brew_session bs
	LEFT JOIN batch ba ON ba.id = bs.batch_id
	LEFT JOIN volume wv ON wv.id = bs.wort_volume_id
	LEFT JOIN vessel mv ON mv.id = bs.mash_vessel_id
	LEFT JOIN vessel bv ON bv.id = bs.boil_vessel_id`

func scanBrewSession(row pgx.Row) (BrewSession, error) {
	var session BrewSession
	err := row.Scan(
		&session.ID,
		&session.UUID,
		&session.BatchID,
		&session.BatchUUID,
		&session.WortVolumeID,
		&session.WortVolumeUUID,
		&session.MashVesselID,
		&session.MashVesselUUID,
		&session.BoilVesselID,
		&session.BoilVesselUUID,
		&session.BrewedAt,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.DeletedAt,
	)
	if err != nil {
		return BrewSession{}, err
	}
	return session, nil
}

func (c *Client) GetBrewSession(ctx context.Context, id int64) (BrewSession, error) {
	session, err := scanBrewSession(c.DB().QueryRow(ctx,
		brewSessionSelectWithJoins+` WHERE bs.id = $1 AND bs.deleted_at IS NULL`, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BrewSession{}, service.ErrNotFound
		}
		return BrewSession{}, fmt.Errorf("getting brew session: %w", err)
	}
	return session, nil
}

func (c *Client) GetBrewSessionByUUID(ctx context.Context, sessionUUID string) (BrewSession, error) {
	session, err := scanBrewSession(c.DB().QueryRow(ctx,
		brewSessionSelectWithJoins+` WHERE bs.uuid = $1 AND bs.deleted_at IS NULL`, sessionUUID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BrewSession{}, service.ErrNotFound
		}
		return BrewSession{}, fmt.Errorf("getting brew session by uuid: %w", err)
	}
	return session, nil
}

func (c *Client) ListBrewSessionsByBatch(ctx context.Context, batchID int64) ([]BrewSession, error) {
	rows, err := c.DB().Query(ctx, brewSessionSelectWithJoins+`
		WHERE bs.batch_id = $1 AND bs.deleted_at IS NULL
		ORDER BY bs.brewed_at ASC`, batchID)
	if err != nil {
		return nil, fmt.Errorf("listing brew sessions by batch: %w", err)
	}
	defer rows.Close()

	var sessions []BrewSession
	for rows.Next() {
		var session BrewSession
		if err := rows.Scan(
			&session.ID,
			&session.UUID,
			&session.BatchID,
			&session.BatchUUID,
			&session.WortVolumeID,
			&session.WortVolumeUUID,
			&session.MashVesselID,
			&session.MashVesselUUID,
			&session.BoilVesselID,
			&session.BoilVesselUUID,
			&session.BrewedAt,
			&session.Notes,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning brew session: %w", err)
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing brew sessions by batch: %w", err)
	}

	return sessions, nil
}

func (c *Client) ListBrewSessionsByBatchUUID(ctx context.Context, batchUUID string) ([]BrewSession, error) {
	batch, err := c.GetBatchByUUID(ctx, batchUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving batch uuid: %w", err)
	}
	return c.ListBrewSessionsByBatch(ctx, batch.ID)
}

func (c *Client) UpdateBrewSession(ctx context.Context, id int64, session BrewSession) (BrewSession, error) {
	err := c.DB().QueryRow(ctx, `
		UPDATE brew_session
		SET batch_id = $1, wort_volume_id = $2, mash_vessel_id = $3, boil_vessel_id = $4, brewed_at = $5, notes = $6, updated_at = timezone('utc', now())
		WHERE id = $7 AND deleted_at IS NULL
		RETURNING id, uuid, batch_id, wort_volume_id, mash_vessel_id, boil_vessel_id, brewed_at, notes, created_at, updated_at, deleted_at`,
		session.BatchID,
		session.WortVolumeID,
		session.MashVesselID,
		session.BoilVesselID,
		session.BrewedAt,
		session.Notes,
		id,
	).Scan(
		&session.ID,
		&session.UUID,
		&session.BatchID,
		&session.WortVolumeID,
		&session.MashVesselID,
		&session.BoilVesselID,
		&session.BrewedAt,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BrewSession{}, service.ErrNotFound
		}
		return BrewSession{}, fmt.Errorf("updating brew session: %w", err)
	}

	// Resolve FK UUIDs
	if session.BatchID != nil {
		var batchUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, *session.BatchID).Scan(&batchUUID); err == nil {
			session.BatchUUID = &batchUUID
		}
	}
	if session.WortVolumeID != nil {
		var volUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM volume WHERE id = $1`, *session.WortVolumeID).Scan(&volUUID); err == nil {
			session.WortVolumeUUID = &volUUID
		}
	}
	if session.MashVesselID != nil {
		var vesselUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM vessel WHERE id = $1`, *session.MashVesselID).Scan(&vesselUUID); err == nil {
			session.MashVesselUUID = &vesselUUID
		}
	}
	if session.BoilVesselID != nil {
		var vesselUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM vessel WHERE id = $1`, *session.BoilVesselID).Scan(&vesselUUID); err == nil {
			session.BoilVesselUUID = &vesselUUID
		}
	}

	return session, nil
}
