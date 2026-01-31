package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateBatchRelation(ctx context.Context, relation BatchRelation) (BatchRelation, error) {
	err := c.db.QueryRow(ctx, `
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

	return relation, nil
}

func (c *Client) GetBatchRelation(ctx context.Context, id int64) (BatchRelation, error) {
	var relation BatchRelation
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, parent_batch_id, child_batch_id, relation_type, volume_id, created_at, updated_at, deleted_at
		FROM batch_relation
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return BatchRelation{}, service.ErrNotFound
		}
		return BatchRelation{}, fmt.Errorf("getting batch relation: %w", err)
	}

	return relation, nil
}

func (c *Client) ListBatchRelations(ctx context.Context, batchID int64) ([]BatchRelation, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, parent_batch_id, child_batch_id, relation_type, volume_id, created_at, updated_at, deleted_at
		FROM batch_relation
		WHERE deleted_at IS NULL
		AND (parent_batch_id = $1 OR child_batch_id = $1)
		ORDER BY created_at ASC`,
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
			&relation.ChildBatchID,
			&relation.RelationType,
			&relation.VolumeID,
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
