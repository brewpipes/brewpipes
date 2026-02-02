package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBatchRelationRequest struct {
	ParentBatchID int64  `json:"parent_batch_id"`
	ChildBatchID  int64  `json:"child_batch_id"`
	RelationType  string `json:"relation_type"`
	VolumeID      *int64 `json:"volume_id"`
}

func (r CreateBatchRelationRequest) Validate() error {
	if r.ParentBatchID <= 0 || r.ChildBatchID <= 0 {
		return fmt.Errorf("parent_batch_id and child_batch_id are required")
	}
	if r.ParentBatchID == r.ChildBatchID {
		return fmt.Errorf("parent_batch_id and child_batch_id must differ")
	}
	if err := validateRelationType(r.RelationType); err != nil {
		return err
	}

	return nil
}

type BatchRelationResponse struct {
	ID            int64      `json:"id"`
	UUID          string     `json:"uuid"`
	ParentBatchID int64      `json:"parent_batch_id"`
	ChildBatchID  int64      `json:"child_batch_id"`
	RelationType  string     `json:"relation_type"`
	VolumeID      *int64     `json:"volume_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

func NewBatchRelationResponse(relation storage.BatchRelation) BatchRelationResponse {
	return BatchRelationResponse{
		ID:            relation.ID,
		UUID:          relation.UUID.String(),
		ParentBatchID: relation.ParentBatchID,
		ChildBatchID:  relation.ChildBatchID,
		RelationType:  relation.RelationType,
		VolumeID:      relation.VolumeID,
		CreatedAt:     relation.CreatedAt,
		UpdatedAt:     relation.UpdatedAt,
		DeletedAt:     relation.DeletedAt,
	}
}

func NewBatchRelationsResponse(relations []storage.BatchRelation) []BatchRelationResponse {
	resp := make([]BatchRelationResponse, 0, len(relations))
	for _, relation := range relations {
		resp = append(resp, NewBatchRelationResponse(relation))
	}
	return resp
}
