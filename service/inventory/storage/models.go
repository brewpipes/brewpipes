package storage

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/database/entity"
	"github.com/gofrs/uuid/v5"
)

const (
	IngredientCategoryFermentable = "fermentable"
	IngredientCategoryHop         = "hop"
	IngredientCategoryYeast       = "yeast"
	IngredientCategoryAdjunct     = "adjunct"
	IngredientCategorySalt        = "salt"
	IngredientCategoryChemical    = "chemical"
	IngredientCategoryGas         = "gas"
	IngredientCategoryOther       = "other"
)

const (
	HopFormPellet    = "pellet"
	HopFormWholeLeaf = "whole_leaf"
	HopFormCryo      = "cryo"
	HopFormExtract   = "extract"
	HopFormOther     = "other"
)

const (
	YeastFormLiquid     = "liquid"
	YeastFormDry        = "dry"
	YeastFormSlurry     = "slurry"
	YeastFormPropagated = "propagated"
	YeastFormOther      = "other"
)

const (
	StockLocationTypeDry       = "dry"
	StockLocationTypeCold      = "cold"
	StockLocationTypeGas       = "gas"
	StockLocationTypeBulk      = "bulk"
	StockLocationTypePackaging = "packaging"
	StockLocationTypeOther     = "other"
)

const (
	OriginatorTypeMaltster    = "maltster"
	OriginatorTypeHopProducer = "hop_producer"
	OriginatorTypeYeastLab    = "yeast_lab"
	OriginatorTypeGasVendor   = "gas_vendor"
	OriginatorTypeOther       = "other"
)

const (
	MovementDirectionIn  = "in"
	MovementDirectionOut = "out"
)

const (
	MovementReasonReceive  = "receive"
	MovementReasonUse      = "use"
	MovementReasonTransfer = "transfer"
	MovementReasonAdjust   = "adjust"
	MovementReasonWaste    = "waste"
	MovementReasonPackage  = "package"
)

// Beer lot container types.
const (
	BeerLotContainerKeg     = "keg"
	BeerLotContainerCan     = "can"
	BeerLotContainerBottle  = "bottle"
	BeerLotContainerCask    = "cask"
	BeerLotContainerGrowler = "growler"
	BeerLotContainerOther   = "other"
)

// Beer lot item statuses.
const (
	BeerLotItemStatusAvailable = "available"
	BeerLotItemStatusReserved  = "reserved"
	BeerLotItemStatusSold      = "sold"
	BeerLotItemStatusReturned  = "returned"
	BeerLotItemStatusDamaged   = "damaged"
	BeerLotItemStatusDestroyed = "destroyed"
)

const (
	AdjustmentReasonCycleCount = "cycle_count"
	AdjustmentReasonSpoilage   = "spoilage"
	AdjustmentReasonShrink     = "shrink"
	AdjustmentReasonDamage     = "damage"
	AdjustmentReasonCorrection = "correction"
	AdjustmentReasonOther      = "other"
)

type Ingredient struct {
	entity.Identifiers
	Name        string
	Category    string
	DefaultUnit string
	Description *string
	entity.Timestamps
}

type IngredientMaltDetail struct {
	entity.Identifiers
	IngredientID   int64
	IngredientUUID string // Joined from ingredient table
	MaltsterName   *string
	Variety        *string
	Lovibond       *float64
	SRM            *float64
	DiastaticPower *float64
	entity.Timestamps
}

type IngredientHopDetail struct {
	entity.Identifiers
	IngredientID   int64
	IngredientUUID string // Joined from ingredient table
	ProducerName   *string
	Variety        *string
	CropYear       *int
	Form           *string
	AlphaAcid      *float64
	BetaAcid       *float64
	entity.Timestamps
}

type IngredientYeastDetail struct {
	entity.Identifiers
	IngredientID   int64
	IngredientUUID string // Joined from ingredient table
	LabName        *string
	Strain         *string
	Form           *string
	entity.Timestamps
}

type StockLocation struct {
	entity.Identifiers
	Name         string
	LocationType *string
	Description  *string
	entity.Timestamps
}

type InventoryReceipt struct {
	entity.Identifiers
	SupplierUUID      *uuid.UUID
	PurchaseOrderUUID *uuid.UUID
	ReferenceCode     *string
	ReceivedAt        time.Time
	Notes             *string
	entity.Timestamps
}

type IngredientLot struct {
	entity.Identifiers
	IngredientID          int64
	IngredientUUID        string // Joined from ingredient table
	ReceiptID             *int64
	ReceiptUUID           *string // Joined from inventory_receipt table
	SupplierUUID          *uuid.UUID
	PurchaseOrderLineUUID *uuid.UUID
	BreweryLotCode        *string
	OriginatorLotCode     *string
	OriginatorName        *string
	OriginatorType        *string
	ReceivedAt            time.Time
	ReceivedAmount        int64
	ReceivedUnit          string
	CurrentAmount         int64 // Computed from inventory_movement ledger
	BestByAt              *time.Time
	ExpiresAt             *time.Time
	Notes                 *string
	entity.Timestamps
}

type IngredientLotMaltDetail struct {
	entity.Identifiers
	IngredientLotID   int64
	IngredientLotUUID string // Joined from ingredient_lot table
	MoisturePercent   *float64
	entity.Timestamps
}

type IngredientLotHopDetail struct {
	entity.Identifiers
	IngredientLotID   int64
	IngredientLotUUID string // Joined from ingredient_lot table
	AlphaAcid         *float64
	BetaAcid          *float64
	entity.Timestamps
}

type IngredientLotYeastDetail struct {
	entity.Identifiers
	IngredientLotID   int64
	IngredientLotUUID string // Joined from ingredient_lot table
	ViabilityPercent  *float64
	Generation        *int
	entity.Timestamps
}

type InventoryUsage struct {
	entity.Identifiers
	ProductionRefUUID *uuid.UUID
	UsedAt            time.Time
	Notes             *string
	entity.Timestamps
}

type InventoryAdjustment struct {
	entity.Identifiers
	Reason     string
	AdjustedAt time.Time
	Notes      *string
	entity.Timestamps
}

type InventoryTransfer struct {
	entity.Identifiers
	SourceLocationID   int64
	SourceLocationUUID string // Joined from stock_location table
	DestLocationID     int64
	DestLocationUUID   string // Joined from stock_location table
	TransferredAt      time.Time
	Notes              *string
	entity.Timestamps
}

type BeerLot struct {
	entity.Identifiers
	ProductionBatchUUID uuid.UUID
	PackagingRunUUID    *uuid.UUID
	LotCode             *string
	BestBy              *time.Time
	PackageFormatName   *string
	Container           *string
	VolumePerUnit       *int64
	VolumePerUnitUnit   *string
	Quantity            *int
	PackagedAt          time.Time
	Notes               *string
	entity.Timestamps
}

// BeerLotItem represents a single trackable unit within a beer lot.
type BeerLotItem struct {
	entity.Identifiers
	BeerLotID   int64
	BeerLotUUID string // Joined from beer_lot table
	Status      string
	Identifier  *string
	Notes       *string
	entity.Timestamps
}

type InventoryMovement struct {
	entity.Identifiers
	IngredientLotID   *int64
	IngredientLotUUID *string // Joined from ingredient_lot table
	BeerLotID         *int64
	BeerLotUUID       *string // Joined from beer_lot table
	StockLocationID   int64
	StockLocationUUID string // Joined from stock_location table
	Direction         string
	Reason            string
	Amount            int64
	AmountUnit        string
	OccurredAt        time.Time
	ReceiptID         *int64
	ReceiptUUID       *string // Joined from inventory_receipt table
	UsageID           *int64
	UsageUUID         *string // Joined from inventory_usage table
	AdjustmentID      *int64
	AdjustmentUUID    *string // Joined from inventory_adjustment table
	TransferID        *int64
	TransferUUID      *string // Joined from inventory_transfer table
	Notes             *string
	entity.Timestamps
}
