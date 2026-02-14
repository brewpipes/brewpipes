package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBatchRelationRequest struct {
	ParentBatchUUID string  `json:"parent_batch_uuid"`
	ChildBatchUUID  string  `json:"child_batch_uuid"`
	RelationType    string  `json:"relation_type"`
	VolumeUUID      *string `json:"volume_uuid"`
}

func (r CreateBatchRelationRequest) Validate() error {
	if err := validate.Required(r.ParentBatchUUID, "parent_batch_uuid"); err != nil {
		return err
	}
	if err := validate.Required(r.ChildBatchUUID, "child_batch_uuid"); err != nil {
		return err
	}
	if r.ParentBatchUUID == r.ChildBatchUUID {
		return fmt.Errorf("parent_batch_uuid and child_batch_uuid must differ")
	}
	if err := validateRelationType(r.RelationType); err != nil {
		return err
	}

	return nil
}

type BatchRelationResponse struct {
	UUID            string     `json:"uuid"`
	ParentBatchUUID string     `json:"parent_batch_uuid"`
	ChildBatchUUID  string     `json:"child_batch_uuid"`
	RelationType    string     `json:"relation_type"`
	VolumeUUID      *string    `json:"volume_uuid,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func NewBatchRelationResponse(relation storage.BatchRelation) BatchRelationResponse {
	return BatchRelationResponse{
		UUID:            relation.UUID.String(),
		ParentBatchUUID: relation.ParentBatchUUID,
		ChildBatchUUID:  relation.ChildBatchUUID,
		RelationType:    relation.RelationType,
		VolumeUUID:      relation.VolumeUUID,
		CreatedAt:       relation.CreatedAt,
		UpdatedAt:       relation.UpdatedAt,
		DeletedAt:       relation.DeletedAt,
	}
}

func NewBatchRelationsResponse(relations []storage.BatchRelation) []BatchRelationResponse {
	resp := make([]BatchRelationResponse, 0, len(relations))
	for _, relation := range relations {
		resp = append(resp, NewBatchRelationResponse(relation))
	}
	return resp
}
