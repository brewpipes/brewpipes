package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBrewSessionRequest struct {
	BatchUUID      *string    `json:"batch_uuid"`
	WortVolumeUUID *string    `json:"wort_volume_uuid"`
	MashVesselUUID *string    `json:"mash_vessel_uuid"`
	BoilVesselUUID *string    `json:"boil_vessel_uuid"`
	BrewedAt       *time.Time `json:"brewed_at"`
	Notes          *string    `json:"notes"`
}

func (r CreateBrewSessionRequest) Validate() error {
	if r.BrewedAt == nil {
		return fmt.Errorf("brewed_at is required")
	}
	return nil
}

type UpdateBrewSessionRequest struct {
	BatchUUID      *string    `json:"batch_uuid"`
	WortVolumeUUID *string    `json:"wort_volume_uuid"`
	MashVesselUUID *string    `json:"mash_vessel_uuid"`
	BoilVesselUUID *string    `json:"boil_vessel_uuid"`
	BrewedAt       *time.Time `json:"brewed_at"`
	Notes          *string    `json:"notes"`
}

func (r UpdateBrewSessionRequest) Validate() error {
	if r.BrewedAt == nil {
		return fmt.Errorf("brewed_at is required")
	}
	return nil
}

type BrewSessionResponse struct {
	UUID           string     `json:"uuid"`
	BatchUUID      *string    `json:"batch_uuid,omitempty"`
	WortVolumeUUID *string    `json:"wort_volume_uuid,omitempty"`
	MashVesselUUID *string    `json:"mash_vessel_uuid,omitempty"`
	BoilVesselUUID *string    `json:"boil_vessel_uuid,omitempty"`
	BrewedAt       time.Time  `json:"brewed_at"`
	Notes          *string    `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

func NewBrewSessionResponse(session storage.BrewSession) BrewSessionResponse {
	return BrewSessionResponse{
		UUID:           session.UUID.String(),
		BatchUUID:      session.BatchUUID,
		WortVolumeUUID: session.WortVolumeUUID,
		MashVesselUUID: session.MashVesselUUID,
		BoilVesselUUID: session.BoilVesselUUID,
		BrewedAt:       session.BrewedAt,
		Notes:          session.Notes,
		CreatedAt:      session.CreatedAt,
		UpdatedAt:      session.UpdatedAt,
		DeletedAt:      session.DeletedAt,
	}
}

func NewBrewSessionsResponse(sessions []storage.BrewSession) []BrewSessionResponse {
	resp := make([]BrewSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		resp = append(resp, NewBrewSessionResponse(session))
	}
	return resp
}
