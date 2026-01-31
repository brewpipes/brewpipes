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
	MaltsterName   *string
	Variety        *string
	Lovibond       *float64
	SRM            *float64
	DiastaticPower *float64
	entity.Timestamps
}

type IngredientHopDetail struct {
	entity.Identifiers
	IngredientID int64
	ProducerName *string
	Variety      *string
	CropYear     *int
	Form         *string
	AlphaAcid    *float64
	BetaAcid     *float64
	entity.Timestamps
}

type IngredientYeastDetail struct {
	entity.Identifiers
	IngredientID int64
	LabName      *string
	Strain       *string
	Form         *string
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
	SupplierUUID  *uuid.UUID
	ReferenceCode *string
	ReceivedAt    time.Time
	Notes         *string
	entity.Timestamps
}

type IngredientLot struct {
	entity.Identifiers
	IngredientID      int64
	ReceiptID         *int64
	SupplierUUID      *uuid.UUID
	BreweryLotCode    *string
	OriginatorLotCode *string
	OriginatorName    *string
	OriginatorType    *string
	ReceivedAt        time.Time
	ReceivedAmount    int64
	ReceivedUnit      string
	BestByAt          *time.Time
	ExpiresAt         *time.Time
	Notes             *string
	entity.Timestamps
}

type IngredientLotMaltDetail struct {
	entity.Identifiers
	IngredientLotID int64
	MoisturePercent *float64
	entity.Timestamps
}

type IngredientLotHopDetail struct {
	entity.Identifiers
	IngredientLotID int64
	AlphaAcid       *float64
	BetaAcid        *float64
	entity.Timestamps
}

type IngredientLotYeastDetail struct {
	entity.Identifiers
	IngredientLotID  int64
	ViabilityPercent *float64
	Generation       *int
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
	SourceLocationID int64
	DestLocationID   int64
	TransferredAt    time.Time
	Notes            *string
	entity.Timestamps
}

type BeerLot struct {
	entity.Identifiers
	ProductionBatchUUID uuid.UUID
	LotCode             *string
	PackagedAt          time.Time
	Notes               *string
	entity.Timestamps
}

type InventoryMovement struct {
	entity.Identifiers
	IngredientLotID *int64
	BeerLotID       *int64
	StockLocationID int64
	Direction       string
	Reason          string
	Amount          int64
	AmountUnit      string
	OccurredAt      time.Time
	ReceiptID       *int64
	UsageID         *int64
	AdjustmentID    *int64
	TransferID      *int64
	Notes           *string
	entity.Timestamps
}
