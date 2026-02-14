package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateBatchRelation(ctx context.Context, relation BatchRelation) (BatchRelation, error) {
	err := c.DB().QueryRow(ctx, `
		INSERT INTO batch_relation (
			parent_batch_id,
			child_batch_id,
			relation_type,
			volume_id
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, parent_batch_id, child_batch_id, relation_type, volume_id, created_at, updated_at, deleted_at`,
		relation.ParentBatchID,
		relation.ChildBatchID,
		relation.RelationType,
		relation.VolumeID,
	).Scan(
		&relation.ID,
		&relation.UUID,
		&relation.ParentBatchID,
		&relation.ChildBatchID,
		&relation.RelationType,
		&relation.VolumeID,
		&relation.CreatedAt,
		&relation.UpdatedAt,
		&relation.DeletedAt,
	)
	if err != nil {
		return BatchRelation{}, fmt.Errorf("creating batch relation: %w", err)
	}

	// Resolve UUIDs
	c.DB().QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, relation.ParentBatchID).Scan(&relation.ParentBatchUUID)
	c.DB().QueryRow(ctx, `SELECT uuid FROM batch WHERE id = $1`, relation.ChildBatchID).Scan(&relation.ChildBatchUUID)
	if relation.VolumeID != nil {
		var volUUID string
		if err := c.DB().QueryRow(ctx, `SELECT uuid FROM volume WHERE id = $1`, *relation.VolumeID).Scan(&volUUID); err == nil {
			relation.VolumeUUID = &volUUID
		}
	}

	return relation, nil
}

func (c *Client) GetBatchRelation(ctx context.Context, id int64) (BatchRelation, error) {
	var relation BatchRelation
	err := c.DB().QueryRow(ctx, `
		SELECT br.id, br.uuid, br.parent_batch_id, pb.uuid, br.child_batch_id, cb.uuid,
		       br.relation_type, br.volume_id, v.uuid,
		       br.created_at, br.updated_at, br.deleted_at
		FROM batch_relation br
		JOIN batch pb ON pb.id = br.parent_batch_id
		JOIN batch cb ON cb.id = br.child_batch_id
		LEFT JOIN volume v ON v.id = br.volume_id
		WHERE br.id = $1 AND br.deleted_at IS NULL`,
		id,
	).Scan(
		&relation.ID,
		&relation.UUID,
		&relation.ParentBatchID,
		&relation.ParentBatchUUID,
		&relation.ChildBatchID,
		&relation.ChildBatchUUID,
		&relation.RelationType,
		&relation.VolumeID,
		&relation.VolumeUUID,
		&relation.CreatedAt,
		&relation.UpdatedAt,
		&relation.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchRelation{}, service.ErrNotFound
		}
		return BatchRelation{}, fmt.Errorf("getting batch relation: %w", err)
	}

	return relation, nil
}

func (c *Client) GetBatchRelationByUUID(ctx context.Context, relationUUID string) (BatchRelation, error) {
	var relation BatchRelation
	err := c.DB().QueryRow(ctx, `
		SELECT br.id, br.uuid, br.parent_batch_id, pb.uuid, br.child_batch_id, cb.uuid,
		       br.relation_type, br.volume_id, v.uuid,
		       br.created_at, br.updated_at, br.deleted_at
		FROM batch_relation br
		JOIN batch pb ON pb.id = br.parent_batch_id
		JOIN batch cb ON cb.id = br.child_batch_id
		LEFT JOIN volume v ON v.id = br.volume_id
		WHERE br.uuid = $1 AND br.deleted_at IS NULL`,
		relationUUID,
	).Scan(
		&relation.ID,
		&relation.UUID,
		&relation.ParentBatchID,
		&relation.ParentBatchUUID,
		&relation.ChildBatchID,
		&relation.ChildBatchUUID,
		&relation.RelationType,
		&relation.VolumeID,
		&relation.VolumeUUID,
		&relation.CreatedAt,
		&relation.UpdatedAt,
		&relation.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchRelation{}, service.ErrNotFound
		}
		return BatchRelation{}, fmt.Errorf("getting batch relation by uuid: %w", err)
	}

	return relation, nil
}

func (c *Client) ListBatchRelations(ctx context.Context, batchID int64) ([]BatchRelation, error) {
	rows, err := c.DB().Query(ctx, `
		SELECT br.id, br.uuid, br.parent_batch_id, pb.uuid, br.child_batch_id, cb.uuid,
		       br.relation_type, br.volume_id, v.uuid,
		       br.created_at, br.updated_at, br.deleted_at
		FROM batch_relation br
		JOIN batch pb ON pb.id = br.parent_batch_id
		JOIN batch cb ON cb.id = br.child_batch_id
		LEFT JOIN volume v ON v.id = br.volume_id
		WHERE br.deleted_at IS NULL
		AND (br.parent_batch_id = $1 OR br.child_batch_id = $1)
		ORDER BY br.created_at ASC`,
		batchID,
	)
	if err != nil {
		return nil, fmt.Errorf("listing batch relations: %w", err)
	}
	defer rows.Close()

	var relations []BatchRelation
	for rows.Next() {
		var relation BatchRelation
		if err := rows.Scan(
			&relation.ID,
			&relation.UUID,
			&relation.ParentBatchID,
			&relation.ParentBatchUUID,
			&relation.ChildBatchID,
			&relation.ChildBatchUUID,
			&relation.RelationType,
			&relation.VolumeID,
			&relation.VolumeUUID,
			&relation.CreatedAt,
			&relation.UpdatedAt,
			&relation.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning batch relation: %w", err)
		}
		relations = append(relations, relation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing batch relations: %w", err)
	}

	return relations, nil
}

func (c *Client) ListBatchRelationsByBatchUUID(ctx context.Context, batchUUID string) ([]BatchRelation, error) {
	batch, err := c.GetBatchByUUID(ctx, batchUUID)
	if err != nil {
		return nil, fmt.Errorf("resolving batch uuid: %w", err)
	}

	return c.ListBatchRelations(ctx, batch.ID)
}
