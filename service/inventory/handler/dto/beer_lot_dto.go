package dto

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateBeerLotRequest struct {
	ProductionBatchUUID string     `json:"production_batch_uuid"`
	LotCode             *string    `json:"lot_code"`
	PackagedAt          *time.Time `json:"packaged_at"`
	Notes               *string    `json:"notes"`
}

func (r CreateBeerLotRequest) Validate() error {
	return validate.Required(r.ProductionBatchUUID, "production_batch_uuid")
}

type BeerLotResponse struct {
	UUID                string     `json:"uuid"`
	ProductionBatchUUID string     `json:"production_batch_uuid"`
	LotCode             *string    `json:"lot_code,omitempty"`
	PackagedAt          time.Time  `json:"packaged_at"`
	Notes               *string    `json:"notes,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

func NewBeerLotResponse(lot storage.BeerLot) BeerLotResponse {
	return BeerLotResponse{
		UUID:                lot.UUID.String(),
		ProductionBatchUUID: lot.ProductionBatchUUID.String(),
		LotCode:             lot.LotCode,
		PackagedAt:          lot.PackagedAt,
		Notes:               lot.Notes,
		CreatedAt:           lot.CreatedAt,
		UpdatedAt:           lot.UpdatedAt,
		DeletedAt:           lot.DeletedAt,
	}
}

func NewBeerLotsResponse(lots []storage.BeerLot) []BeerLotResponse {
	resp := make([]BeerLotResponse, 0, len(lots))
	for _, lot := range lots {
		resp = append(resp, NewBeerLotResponse(lot))
	}
	return resp
}
