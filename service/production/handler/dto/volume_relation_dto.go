package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateVolumeRelationRequest struct {
	ParentVolumeID int64  `json:"parent_volume_id"`
	ChildVolumeID  int64  `json:"child_volume_id"`
	RelationType   string `json:"relation_type"`
	Amount         int64  `json:"amount"`
	AmountUnit     string `json:"amount_unit"`
}

func (r CreateVolumeRelationRequest) Validate() error {
	if r.ParentVolumeID <= 0 || r.ChildVolumeID <= 0 {
		return fmt.Errorf("parent_volume_id and child_volume_id are required")
	}
	if r.ParentVolumeID == r.ChildVolumeID {
		return fmt.Errorf("parent_volume_id and child_volume_id must differ")
	}
	if err := validateRelationType(r.RelationType); err != nil {
		return err
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if err := validateVolumeUnit(r.AmountUnit); err != nil {
		return err
	}

	return nil
}

type VolumeRelationResponse struct {
	ID             int64      `json:"id"`
	UUID           string     `json:"uuid"`
	ParentVolumeID int64      `json:"parent_volume_id"`
	ChildVolumeID  int64      `json:"child_volume_id"`
	RelationType   string     `json:"relation_type"`
	Amount         int64      `json:"amount"`
	AmountUnit     string     `json:"amount_unit"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func NewVolumeRelationResponse(relation storage.VolumeRelation) VolumeRelationResponse {
	return VolumeRelationResponse{
		ID:             relation.ID,
		UUID:           relation.UUID.String(),
		ParentVolumeID: relation.ParentVolumeID,
		ChildVolumeID:  relation.ChildVolumeID,
		RelationType:   relation.RelationType,
		Amount:         relation.Amount,
		AmountUnit:     relation.AmountUnit,
		CreatedAt:      relation.CreatedAt,
		UpdatedAt:      relation.UpdatedAt,
		DeletedAt:      relation.DeletedAt,
	}
}

func NewVolumeRelationsResponse(relations []storage.VolumeRelation) []VolumeRelationResponse {
	resp := make([]VolumeRelationResponse, 0, len(relations))
	for _, relation := range relations {
		resp = append(resp, NewVolumeRelationResponse(relation))
	}
	return resp
}
