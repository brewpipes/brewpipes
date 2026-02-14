package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateVolumeRelationRequest struct {
	ParentVolumeUUID string `json:"parent_volume_uuid"`
	ChildVolumeUUID  string `json:"child_volume_uuid"`
	RelationType     string `json:"relation_type"`
	Amount           int64  `json:"amount"`
	AmountUnit       string `json:"amount_unit"`
}

func (r CreateVolumeRelationRequest) Validate() error {
	if err := validate.Required(r.ParentVolumeUUID, "parent_volume_uuid"); err != nil {
		return err
	}
	if err := validate.Required(r.ChildVolumeUUID, "child_volume_uuid"); err != nil {
		return err
	}
	if r.ParentVolumeUUID == r.ChildVolumeUUID {
		return fmt.Errorf("parent_volume_uuid and child_volume_uuid must differ")
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
	UUID             string     `json:"uuid"`
	ParentVolumeUUID string     `json:"parent_volume_uuid"`
	ChildVolumeUUID  string     `json:"child_volume_uuid"`
	RelationType     string     `json:"relation_type"`
	Amount           int64      `json:"amount"`
	AmountUnit       string     `json:"amount_unit"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

func NewVolumeRelationResponse(relation storage.VolumeRelation) VolumeRelationResponse {
	return VolumeRelationResponse{
		UUID:             relation.UUID.String(),
		ParentVolumeUUID: relation.ParentVolumeUUID,
		ChildVolumeUUID:  relation.ChildVolumeUUID,
		RelationType:     relation.RelationType,
		Amount:           relation.Amount,
		AmountUnit:       relation.AmountUnit,
		CreatedAt:        relation.CreatedAt,
		UpdatedAt:        relation.UpdatedAt,
		DeletedAt:        relation.DeletedAt,
	}
}

func NewVolumeRelationsResponse(relations []storage.VolumeRelation) []VolumeRelationResponse {
	resp := make([]VolumeRelationResponse, 0, len(relations))
	for _, relation := range relations {
		resp = append(resp, NewVolumeRelationResponse(relation))
	}
	return resp
}
