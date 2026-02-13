package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBatchProcessPhaseRequest struct {
	BatchUUID    string     `json:"batch_uuid"`
	ProcessPhase string     `json:"process_phase"`
	PhaseAt      *time.Time `json:"phase_at"`
}

func (r CreateBatchProcessPhaseRequest) Validate() error {
	if err := validateRequired(r.BatchUUID, "batch_uuid"); err != nil {
		return err
	}
	if err := validateProcessPhase(r.ProcessPhase); err != nil {
		return err
	}

	return nil
}

type BatchProcessPhaseResponse struct {
	UUID         string     `json:"uuid"`
	BatchUUID    string     `json:"batch_uuid"`
	ProcessPhase string     `json:"process_phase"`
	PhaseAt      time.Time  `json:"phase_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewBatchProcessPhaseResponse(phase storage.BatchProcessPhase) BatchProcessPhaseResponse {
	return BatchProcessPhaseResponse{
		UUID:         phase.UUID.String(),
		BatchUUID:    phase.BatchUUID,
		ProcessPhase: phase.ProcessPhase,
		PhaseAt:      phase.PhaseAt,
		CreatedAt:    phase.CreatedAt,
		UpdatedAt:    phase.UpdatedAt,
		DeletedAt:    phase.DeletedAt,
	}
}

func NewBatchProcessPhasesResponse(phases []storage.BatchProcessPhase) []BatchProcessPhaseResponse {
	resp := make([]BatchProcessPhaseResponse, 0, len(phases))
	for _, phase := range phases {
		resp = append(resp, NewBatchProcessPhaseResponse(phase))
	}
	return resp
}
