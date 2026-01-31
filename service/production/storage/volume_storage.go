package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateVolume(ctx context.Context, volume Volume) (Volume, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO volume (
			name,
			description,
			amount,
			amount_unit
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, name, description, amount, amount_unit, created_at, updated_at, deleted_at`,
		volume.Name,
		volume.Description,
		volume.Amount,
		volume.AmountUnit,
	).Scan(
		&volume.ID,
		&volume.UUID,
		&volume.Name,
		&volume.Description,
		&volume.Amount,
		&volume.AmountUnit,
		&volume.CreatedAt,
		&volume.UpdatedAt,
		&volume.DeletedAt,
	)
	if err != nil {
		return Volume{}, fmt.Errorf("creating volume: %w", err)
	}

	return volume, nil
}

func (c *Client) GetVolume(ctx context.Context, id int64) (Volume, error) {
	var volume Volume
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, name, description, amount, amount_unit, created_at, updated_at, deleted_at
		FROM volume
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&volume.ID,
		&volume.UUID,
		&volume.Name,
		&volume.Description,
		&volume.Amount,
		&volume.AmountUnit,
		&volume.CreatedAt,
		&volume.UpdatedAt,
		&volume.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Volume{}, service.ErrNotFound
		}
		return Volume{}, fmt.Errorf("getting volume: %w", err)
	}

	return volume, nil
}

func (c *Client) ListVolumes(ctx context.Context) ([]Volume, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, name, description, amount, amount_unit, created_at, updated_at, deleted_at
		FROM volume
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing volumes: %w", err)
	}
	defer rows.Close()

	var volumes []Volume
	for rows.Next() {
		var volume Volume
		if err := rows.Scan(
			&volume.ID,
			&volume.UUID,
			&volume.Name,
			&volume.Description,
			&volume.Amount,
			&volume.AmountUnit,
			&volume.CreatedAt,
			&volume.UpdatedAt,
			&volume.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning volume: %w", err)
		}
		volumes = append(volumes, volume)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing volumes: %w", err)
	}

	return volumes, nil
}

func (c *Client) GetVolumes(ctx context.Context) ([]Volume, error) {
	return c.ListVolumes(ctx)
}

func (c *Client) CreateVolumeRelation(ctx context.Context, relation VolumeRelation) (VolumeRelation, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO volume_relation (
			parent_volume_id,
			child_volume_id,
			relation_type,
			amount,
			amount_unit
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, uuid, parent_volume_id, child_volume_id, relation_type, amount, amount_unit, created_at, updated_at, deleted_at`,
		relation.ParentVolumeID,
		relation.ChildVolumeID,
		relation.RelationType,
		relation.Amount,
		relation.AmountUnit,
	).Scan(
		&relation.ID,
		&relation.UUID,
		&relation.ParentVolumeID,
		&relation.ChildVolumeID,
		&relation.RelationType,
		&relation.Amount,
		&relation.AmountUnit,
		&relation.CreatedAt,
		&relation.UpdatedAt,
		&relation.DeletedAt,
	)
	if err != nil {
		return VolumeRelation{}, fmt.Errorf("creating volume relation: %w", err)
	}

	return relation, nil
}

func (c *Client) ListVolumeRelations(ctx context.Context, volumeID int64) ([]VolumeRelation, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, parent_volume_id, child_volume_id, relation_type, amount, amount_unit, created_at, updated_at, deleted_at
		FROM volume_relation
		WHERE deleted_at IS NULL
		AND (parent_volume_id = $1 OR child_volume_id = $1)
		ORDER BY created_at ASC`,
		volumeID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing volume relations: %w", err)
	}
	defer rows.Close()

	var relations []VolumeRelation
	for rows.Next() {
		var relation VolumeRelation
		if err := rows.Scan(
			&relation.ID,
			&relation.UUID,
			&relation.ParentVolumeID,
			&relation.ChildVolumeID,
			&relation.RelationType,
			&relation.Amount,
			&relation.AmountUnit,
			&relation.CreatedAt,
			&relation.UpdatedAt,
			&relation.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning volume relation: %w", err)
		}
		relations = append(relations, relation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing volume relations: %w", err)
	}

	return relations, nil
}
