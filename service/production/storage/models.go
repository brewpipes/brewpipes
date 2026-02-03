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
	Name      string
	StyleID   *int64
	StyleName *string
	Notes     *string
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
