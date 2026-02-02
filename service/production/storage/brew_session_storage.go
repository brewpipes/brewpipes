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

	err := c.db.QueryRow(ctx, `
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

	return session, nil
}

func (c *Client) GetBrewSession(ctx context.Context, id int64) (BrewSession, error) {
	var session BrewSession
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, batch_id, wort_volume_id, mash_vessel_id, boil_vessel_id, brewed_at, notes, created_at, updated_at, deleted_at
		FROM brew_session
		WHERE id = $1 AND deleted_at IS NULL`,
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
		return BrewSession{}, fmt.Errorf("getting brew session: %w", err)
	}

	return session, nil
}

func (c *Client) ListBrewSessionsByBatch(ctx context.Context, batchID int64) ([]BrewSession, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, batch_id, wort_volume_id, mash_vessel_id, boil_vessel_id, brewed_at, notes, created_at, updated_at, deleted_at
		FROM brew_session
		WHERE batch_id = $1 AND deleted_at IS NULL
		ORDER BY brewed_at ASC`,
		batchID,
	)
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
			&session.WortVolumeID,
			&session.MashVesselID,
			&session.BoilVesselID,
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

func (c *Client) UpdateBrewSession(ctx context.Context, id int64, session BrewSession) (BrewSession, error) {
	err := c.db.QueryRow(ctx, `
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

	return session, nil
}
