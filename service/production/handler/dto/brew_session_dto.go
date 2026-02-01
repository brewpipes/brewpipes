package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBrewSessionRequest struct {
	BatchID      *int64     `json:"batch_id"`
	WortVolumeID *int64     `json:"wort_volume_id"`
	MashVesselID *int64     `json:"mash_vessel_id"`
	BoilVesselID *int64     `json:"boil_vessel_id"`
	BrewedAt     *time.Time `json:"brewed_at"`
	Notes        *string    `json:"notes"`
}

func (r CreateBrewSessionRequest) Validate() error {
	if r.BrewedAt == nil {
		return fmt.Errorf("brewed_at is required")
	}
	return nil
}

type UpdateBrewSessionRequest struct {
	BatchID      *int64     `json:"batch_id"`
	WortVolumeID *int64     `json:"wort_volume_id"`
	MashVesselID *int64     `json:"mash_vessel_id"`
	BoilVesselID *int64     `json:"boil_vessel_id"`
	BrewedAt     *time.Time `json:"brewed_at"`
	Notes        *string    `json:"notes"`
}

func (r UpdateBrewSessionRequest) Validate() error {
	if r.BrewedAt == nil {
		return fmt.Errorf("brewed_at is required")
	}
	return nil
}

type BrewSessionResponse struct {
	ID           int64      `json:"id"`
	UUID         string     `json:"uuid"`
	BatchID      *int64     `json:"batch_id,omitempty"`
	WortVolumeID *int64     `json:"wort_volume_id,omitempty"`
	MashVesselID *int64     `json:"mash_vessel_id,omitempty"`
	BoilVesselID *int64     `json:"boil_vessel_id,omitempty"`
	BrewedAt     time.Time  `json:"brewed_at"`
	Notes        *string    `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewBrewSessionResponse(session storage.BrewSession) BrewSessionResponse {
	return BrewSessionResponse{
		ID:           session.ID,
		UUID:         session.UUID.String(),
		BatchID:      session.BatchID,
		WortVolumeID: session.WortVolumeID,
		MashVesselID: session.MashVesselID,
		BoilVesselID: session.BoilVesselID,
		BrewedAt:     session.BrewedAt,
		Notes:        session.Notes,
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
		DeletedAt:    session.DeletedAt,
	}
}

func NewBrewSessionsResponse(sessions []storage.BrewSession) []BrewSessionResponse {
	resp := make([]BrewSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		resp = append(resp, NewBrewSessionResponse(session))
	}
	return resp
}
