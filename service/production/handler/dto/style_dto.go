package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreateStyleRequest struct {
	Name string `json:"name"`
}

func (r CreateStyleRequest) Validate() error {
	return validateRequired(r.Name, "name")
}

type StyleResponse struct {
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewStyleResponse(style storage.Style) StyleResponse {
	return StyleResponse{
		UUID:      style.UUID.String(),
		Name:      style.Name,
		CreatedAt: style.CreatedAt,
		UpdatedAt: style.UpdatedAt,
		DeletedAt: style.DeletedAt,
	}
}

func NewStylesResponse(styles []storage.Style) []StyleResponse {
	resp := make([]StyleResponse, 0, len(styles))
	for _, style := range styles {
		resp = append(resp, NewStyleResponse(style))
	}
	return resp
}
