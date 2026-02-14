package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/uuidutil"
	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/inventory/storage"
)

type CreateIngredientLotRequest struct {
	IngredientUUID    string     `json:"ingredient_uuid"`
	ReceiptUUID       *string    `json:"receipt_uuid"`
	SupplierUUID      *string    `json:"supplier_uuid"`
	BreweryLotCode    *string    `json:"brewery_lot_code"`
	OriginatorLotCode *string    `json:"originator_lot_code"`
	OriginatorName    *string    `json:"originator_name"`
	OriginatorType    *string    `json:"originator_type"`
	ReceivedAt        *time.Time `json:"received_at"`
	ReceivedAmount    int64      `json:"received_amount"`
	ReceivedUnit      string     `json:"received_unit"`
	BestByAt          *time.Time `json:"best_by_at"`
	ExpiresAt         *time.Time `json:"expires_at"`
	Notes             *string    `json:"notes"`
}

func (r CreateIngredientLotRequest) Validate() error {
	if err := validate.Required(r.IngredientUUID, "ingredient_uuid"); err != nil {
		return err
	}
	if r.ReceivedAmount <= 0 {
		return fmt.Errorf("received_amount must be greater than zero")
	}
	if err := validate.Required(r.ReceivedUnit, "received_unit"); err != nil {
		return err
	}
	if r.OriginatorType != nil {
		if err := validateOriginatorType(*r.OriginatorType); err != nil {
			return err
		}
	}
	if r.BestByAt != nil && r.ExpiresAt != nil {
		if r.ExpiresAt.Before(*r.BestByAt) {
			return fmt.Errorf("expires_at must be after best_by_at")
		}
	}

	return nil
}

type IngredientLotResponse struct {
	UUID              string     `json:"uuid"`
	IngredientUUID    string     `json:"ingredient_uuid"`
	ReceiptUUID       *string    `json:"receipt_uuid,omitempty"`
	SupplierUUID      *string    `json:"supplier_uuid,omitempty"`
	BreweryLotCode    *string    `json:"brewery_lot_code,omitempty"`
	OriginatorLotCode *string    `json:"originator_lot_code,omitempty"`
	OriginatorName    *string    `json:"originator_name,omitempty"`
	OriginatorType    *string    `json:"originator_type,omitempty"`
	ReceivedAt        time.Time  `json:"received_at"`
	ReceivedAmount    int64      `json:"received_amount"`
	ReceivedUnit      string     `json:"received_unit"`
	BestByAt          *time.Time `json:"best_by_at,omitempty"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	Notes             *string    `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func NewIngredientLotResponse(lot storage.IngredientLot) IngredientLotResponse {
	return IngredientLotResponse{
		UUID:              lot.UUID.String(),
		IngredientUUID:    lot.IngredientUUID,
		ReceiptUUID:       lot.ReceiptUUID,
		SupplierUUID:      uuidutil.ToStringPointer(lot.SupplierUUID),
		BreweryLotCode:    lot.BreweryLotCode,
		OriginatorLotCode: lot.OriginatorLotCode,
		OriginatorName:    lot.OriginatorName,
		OriginatorType:    lot.OriginatorType,
		ReceivedAt:        lot.ReceivedAt,
		ReceivedAmount:    lot.ReceivedAmount,
		ReceivedUnit:      lot.ReceivedUnit,
		BestByAt:          lot.BestByAt,
		ExpiresAt:         lot.ExpiresAt,
		Notes:             lot.Notes,
		CreatedAt:         lot.CreatedAt,
		UpdatedAt:         lot.UpdatedAt,
		DeletedAt:         lot.DeletedAt,
	}
}

func NewIngredientLotsResponse(lots []storage.IngredientLot) []IngredientLotResponse {
	resp := make([]IngredientLotResponse, 0, len(lots))
	for _, lot := range lots {
		resp = append(resp, NewIngredientLotResponse(lot))
	}
	return resp
}
