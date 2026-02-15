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
	ShortName    string
	BrewDate     *time.Time
	Notes        *string
	RecipeID     *int64
	RecipeUUID   *string // Joined from recipe table
	RecipeName   *string // Joined from recipe table
	CurrentPhase *string // Derived from most recent batch_process_phase
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
	ParentVolumeID   int64
	ParentVolumeUUID string // Joined from volume table
	ChildVolumeID    int64
	ChildVolumeUUID  string // Joined from volume table
	RelationType     string
	Amount           int64
	AmountUnit       string
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
	VesselID   int64
	VesselUUID string // Joined from vessel table
	VolumeID   int64
	VolumeUUID string // Joined from volume table
	InAt       time.Time
	OutAt      *time.Time
	Status     *string
	BatchID    *int64  // Derived from batch_volume join, not stored in occupancy table
	BatchUUID  *string // Derived from batch join, not stored in occupancy table
	entity.Timestamps
}

type Transfer struct {
	entity.Identifiers
	SourceOccupancyID   int64
	SourceOccupancyUUID string // Joined from occupancy table
	DestOccupancyID     int64
	DestOccupancyUUID   string // Joined from occupancy table
	Amount              int64
	AmountUnit          string
	LossAmount          *int64
	LossUnit            *string
	StartedAt           time.Time
	EndedAt             *time.Time
	entity.Timestamps
}

type BatchVolume struct {
	entity.Identifiers
	BatchID     int64
	BatchUUID   string // Joined from batch table
	VolumeID    int64
	VolumeUUID  string // Joined from volume table
	LiquidPhase string
	PhaseAt     time.Time
	entity.Timestamps
}

type BatchProcessPhase struct {
	entity.Identifiers
	BatchID      int64
	BatchUUID    string // Joined from batch table
	ProcessPhase string
	PhaseAt      time.Time
	entity.Timestamps
}

type BatchRelation struct {
	entity.Identifiers
	ParentBatchID   int64
	ParentBatchUUID string // Joined from batch table
	ChildBatchID    int64
	ChildBatchUUID  string // Joined from batch table
	RelationType    string
	VolumeID        *int64
	VolumeUUID      *string // Joined from volume table
	entity.Timestamps
}

type Addition struct {
	entity.Identifiers
	BatchID          *int64
	BatchUUID        *string // Joined from batch table
	OccupancyID      *int64
	OccupancyUUID    *string // Joined from occupancy table
	VolumeID         *int64
	VolumeUUID       *string // Joined from volume table
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
	BatchID       *int64
	BatchUUID     *string // Joined from batch table
	OccupancyID   *int64
	OccupancyUUID *string // Joined from occupancy table
	VolumeID      *int64
	VolumeUUID    *string // Joined from volume table
	Kind          string
	Value         float64
	Unit          *string
	ObservedAt    time.Time
	Notes         *string
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
	StyleUUID           *string // Joined from style table
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
	BatchID        *int64
	BatchUUID      *string // Joined from batch table
	WortVolumeID   *int64
	WortVolumeUUID *string // Joined from volume table
	MashVesselID   *int64
	MashVesselUUID *string // Joined from vessel table
	BoilVesselID   *int64
	BoilVesselUUID *string // Joined from vessel table
	BrewedAt       time.Time
	Notes          *string
	entity.Timestamps
}

// RecipeIngredient represents an ingredient in a recipe's ingredient bill.
type RecipeIngredient struct {
	entity.Identifiers
	RecipeID              int64
	Name                  string     // Human-readable ingredient name
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

// PackageFormat represents a container type for packaged beer (e.g., 1/2 BBL Keg, 16oz Can).
type PackageFormat struct {
	entity.Identifiers
	Name              string
	Container         string
	VolumePerUnit     int64
	VolumePerUnitUnit string
	IsActive          bool
	entity.Timestamps
}

// PackagingRun represents a packaging event for a batch.
type PackagingRun struct {
	entity.Identifiers
	BatchID       int64
	BatchUUID     string // Joined from batch table
	OccupancyID   int64
	OccupancyUUID string // Joined from occupancy table
	StartedAt     time.Time
	EndedAt       *time.Time
	LossAmount    *int64
	LossUnit      *string
	Notes         *string
	entity.Timestamps
}

// PackagingRunLine represents one format line within a packaging run.
type PackagingRunLine struct {
	entity.Identifiers
	PackagingRunID                 int64
	PackagingRunUUID               string // Joined from packaging_run table
	PackageFormatID                int64
	PackageFormatUUID              string // Joined from package_format table
	PackageFormatName              string // Joined from package_format table
	PackageFormatVolumePerUnit     int64  // Joined from package_format table
	PackageFormatVolumePerUnitUnit string // Joined from package_format table
	Quantity                       int
	entity.Timestamps
}
