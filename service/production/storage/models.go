package storage

import (
	"time"

	"github.com/brewpipes/brewpipes/internal/database/entity"
	"github.com/gofrs/uuid/v5"
)

const (
	VolumeUnitML     = "ml"
	VolumeUnitUSFlOz = "usfloz"
	VolumeUnitUKFlOz = "ukfloz"
	VolumeUnitBBL    = "bbl"
)

const (
	RelationTypeSplit = "split"
	RelationTypeBlend = "blend"
)

const (
	VesselStatusActive   = "active"
	VesselStatusInactive = "inactive"
	VesselStatusRetired  = "retired"
)

const (
	VesselTypeMashTun     = "mash_tun"
	VesselTypeLauterTun   = "lauter_tun"
	VesselTypeKettle      = "kettle"
	VesselTypeWhirlpool   = "whirlpool"
	VesselTypeFermenter   = "fermenter"
	VesselTypeBriteTank   = "brite_tank"
	VesselTypeServingTank = "serving_tank"
	VesselTypeOther       = "other"
)

const (
	LiquidPhaseWater = "water"
	LiquidPhaseWort  = "wort"
	LiquidPhaseBeer  = "beer"
)

const (
	ProcessPhasePlanning     = "planning"
	ProcessPhaseMashing      = "mashing"
	ProcessPhaseHeating      = "heating"
	ProcessPhaseBoiling      = "boiling"
	ProcessPhaseCooling      = "cooling"
	ProcessPhaseFermenting   = "fermenting"
	ProcessPhaseConditioning = "conditioning"
	ProcessPhasePackaging    = "packaging"
	ProcessPhaseFinished     = "finished"
)

const (
	AdditionTypeMalt      = "malt"
	AdditionTypeHop       = "hop"
	AdditionTypeYeast     = "yeast"
	AdditionTypeAdjunct   = "adjunct"
	AdditionTypeWaterChem = "water_chem"
	AdditionTypeGas       = "gas"
	AdditionTypeOther     = "other"
)

// Recipe ingredient types
const (
	IngredientTypeFermentable = "fermentable"
	IngredientTypeHop         = "hop"
	IngredientTypeYeast       = "yeast"
	IngredientTypeAdjunct     = "adjunct"
	IngredientTypeSalt        = "salt"
	IngredientTypeChemical    = "chemical"
	IngredientTypeGas         = "gas"
	IngredientTypeOther       = "other"
)

// Recipe ingredient use stages
const (
	UseStageMash         = "mash"
	UseStageBoil         = "boil"
	UseStageWhirlpool    = "whirlpool"
	UseStageFermentation = "fermentation"
	UseStagePackaging    = "packaging"
)

// Recipe ingredient use types
const (
	UseTypeBittering = "bittering"
	UseTypeFlavor    = "flavor"
	UseTypeAroma     = "aroma"
	UseTypeDryHop    = "dry_hop"
	UseTypeBase      = "base"
	UseTypeSpecialty = "specialty"
	UseTypeAdjunct   = "adjunct"
	UseTypeSugar     = "sugar"
	UseTypePrimary   = "primary"
	UseTypeSecondary = "secondary"
	UseTypeBottle    = "bottle"
	UseTypeOther     = "other"
)

// IBU calculation methods
const (
	IBUMethodTinseth = "tinseth"
	IBUMethodRager   = "rager"
	IBUMethodGaretz  = "garetz"
	IBUMethodDaniels = "daniels"
)

const (
	OccupancyStatusFermenting   = "fermenting"
	OccupancyStatusConditioning = "conditioning"
	OccupancyStatusColdCrashing = "cold_crashing"
	OccupancyStatusDryHopping   = "dry_hopping"
	OccupancyStatusCarbonating  = "carbonating"
	OccupancyStatusHolding      = "holding"
	OccupancyStatusPackaging    = "packaging"
)

type Batch struct {
	entity.Identifiers
	ShortName string
	BrewDate  *time.Time
	Notes     *string
	RecipeID  *int64
	entity.Timestamps
}

type Volume struct {
	entity.Identifiers
	Name        *string
	Description *string
	Amount      int64
	AmountUnit  string
	entity.Timestamps
}

type VolumeRelation struct {
	entity.Identifiers
	ParentVolumeID int64
	ChildVolumeID  int64
	RelationType   string
	Amount         int64
	AmountUnit     string
	entity.Timestamps
}

type Vessel struct {
	entity.Identifiers
	Type         string
	Name         string
	Capacity     int64
	CapacityUnit string
	Make         *string
	Model        *string
	Status       string
	entity.Timestamps
}

type Occupancy struct {
	entity.Identifiers
	VesselID int64
	VolumeID int64
	InAt     time.Time
	OutAt    *time.Time
	Status   *string
	BatchID  *int64 // Derived from batch_volume join, not stored in occupancy table
	entity.Timestamps
}

type Transfer struct {
	entity.Identifiers
	SourceOccupancyID int64
	DestOccupancyID   int64
	Amount            int64
	AmountUnit        string
	LossAmount        *int64
	LossUnit          *string
	StartedAt         time.Time
	EndedAt           *time.Time
	entity.Timestamps
}

type BatchVolume struct {
	entity.Identifiers
	BatchID     int64
	VolumeID    int64
	LiquidPhase string
	PhaseAt     time.Time
	entity.Timestamps
}

type BatchProcessPhase struct {
	entity.Identifiers
	BatchID      int64
	ProcessPhase string
	PhaseAt      time.Time
	entity.Timestamps
}

type BatchRelation struct {
	entity.Identifiers
	ParentBatchID int64
	ChildBatchID  int64
	RelationType  string
	VolumeID      *int64
	entity.Timestamps
}

type Addition struct {
	entity.Identifiers
	BatchID          *int64
	OccupancyID      *int64
	VolumeID         *int64
	AdditionType     string
	Stage            *string
	InventoryLotUUID *uuid.UUID
	Amount           int64
	AmountUnit       string
	AddedAt          time.Time
	Notes            *string
	entity.Timestamps
}

type Measurement struct {
	entity.Identifiers
	BatchID     *int64
	OccupancyID *int64
	VolumeID    *int64
	Kind        string
	Value       float64
	Unit        *string
	ObservedAt  time.Time
	Notes       *string
	entity.Timestamps
}

type Style struct {
	entity.Identifiers
	Name string
	entity.Timestamps
}

type Recipe struct {
	entity.Identifiers
	Name                string
	StyleID             *int64
	StyleName           *string
	Notes               *string
	BatchSize           *float64
	BatchSizeUnit       *string
	TargetOG            *float64
	TargetOGMin         *float64
	TargetOGMax         *float64
	TargetFG            *float64
	TargetFGMin         *float64
	TargetFGMax         *float64
	TargetIBU           *float64
	TargetIBUMin        *float64
	TargetIBUMax        *float64
	TargetSRM           *float64
	TargetSRMMin        *float64
	TargetSRMMax        *float64
	TargetCarbonation   *float64
	IBUMethod           *string
	BrewhouseEfficiency *float64
	entity.Timestamps
}

type BrewSession struct {
	entity.Identifiers
	BatchID      *int64
	WortVolumeID *int64
	MashVesselID *int64
	BoilVesselID *int64
	BrewedAt     time.Time
	Notes        *string
	entity.Timestamps
}

// RecipeIngredient represents an ingredient in a recipe's ingredient bill.
type RecipeIngredient struct {
	entity.Identifiers
	RecipeID              int64
	IngredientUUID        *uuid.UUID // Cross-service ref to inventory.ingredient
	IngredientType        string     // fermentable, hop, yeast, adjunct, salt, chemical, gas, other
	Amount                float64
	AmountUnit            string
	UseStage              string  // mash, boil, whirlpool, fermentation, packaging
	UseType               *string // bittering, flavor, aroma, dry_hop, base, specialty, etc.
	TimingDurationMinutes *int
	TimingTemperatureC    *float64
	AlphaAcidAssumed      *float64 // For hops only
	ScalingFactor         float64  // Default 1.0
	SortOrder             int
	Notes                 *string
	entity.Timestamps
}
