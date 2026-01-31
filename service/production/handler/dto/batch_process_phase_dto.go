package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateBatchProcessPhaseRequest struct {
	BatchID      int64      `json:"batch_id"`
	ProcessPhase string     `json:"process_phase"`
	PhaseAt      *time.Time `json:"phase_at"`
}

func (r CreateBatchProcessPhaseRequest) Validate() error {
	if r.BatchID <= 0 {
		return fmt.Errorf("batch_id is required")
	}
	if err := validateProcessPhase(r.ProcessPhase); err != nil {
		return err
	}

	return nil
}

type BatchProcessPhaseResponse struct {
	ID           int64      `json:"id"`
	UUID         string     `json:"uuid"`
	BatchID      int64      `json:"batch_id"`
	ProcessPhase string     `json:"process_phase"`
	PhaseAt      time.Time  `json:"phase_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func NewBatchProcessPhaseResponse(phase storage.BatchProcessPhase) BatchProcessPhaseResponse {
	return BatchProcessPhaseResponse{
		ID:           phase.ID,
		UUID:         phase.UUID.String(),
		BatchID:      phase.BatchID,
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
