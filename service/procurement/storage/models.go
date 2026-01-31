package storage

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/database/entity"
	"github.com/gofrs/uuid/v5"
)

const (
	PurchaseOrderStatusDraft             = "draft"
	PurchaseOrderStatusSubmitted         = "submitted"
	PurchaseOrderStatusConfirmed         = "confirmed"
	PurchaseOrderStatusPartiallyReceived = "partially_received"
	PurchaseOrderStatusReceived          = "received"
	PurchaseOrderStatusCancelled         = "cancelled"
)

const (
	PurchaseOrderItemTypeIngredient = "ingredient"
	PurchaseOrderItemTypePackaging  = "packaging"
	PurchaseOrderItemTypeService    = "service"
	PurchaseOrderItemTypeEquipment  = "equipment"
	PurchaseOrderItemTypeOther      = "other"
)

type Supplier struct {
	entity.Identifiers
	Name         string
	ContactName  *string
	Email        *string
	Phone        *string
	AddressLine1 *string
	AddressLine2 *string
	City         *string
	Region       *string
	PostalCode   *string
	Country      *string
	entity.Timestamps
}

type PurchaseOrder struct {
	entity.Identifiers
	SupplierID  int64
	OrderNumber string
	Status      string
	OrderedAt   *time.Time
	ExpectedAt  *time.Time
	Notes       *string
	entity.Timestamps
}

type PurchaseOrderLine struct {
	entity.Identifiers
	PurchaseOrderID   int64
	LineNumber        int
	ItemType          string
	ItemName          string
	InventoryItemUUID *uuid.UUID
	Quantity          int64
	QuantityUnit      string
	UnitCostCents     int64
	Currency          string
	entity.Timestamps
}

type PurchaseOrderFee struct {
	entity.Identifiers
	PurchaseOrderID int64
	FeeType         string
	AmountCents     int64
	Currency        string
	entity.Timestamps
}
