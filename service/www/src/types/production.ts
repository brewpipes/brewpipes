/**
 * Production domain types for batches, vessels, recipes, and brewing operations.
 */

import type { VolumeUnit } from './units'

// ============================================================================
// Style Types
// ============================================================================

/** Beer style definition */
export interface Style {
  uuid: string
  name: string
  created_at: string
  updated_at: string
}

/** Request payload for creating a new style */
export interface CreateStyleRequest {
  name: string
}

// ============================================================================
// Recipe Types
// ============================================================================

/** Recipe definition with optional style association and target specifications */
export interface Recipe {
  uuid: string
  name: string
  style_uuid: string | null
  style_name: string | null
  notes: string | null
  // Target specifications
  batch_size: number | null
  batch_size_unit: string | null
  target_og: number | null
  target_og_min: number | null
  target_og_max: number | null
  target_fg: number | null
  target_fg_min: number | null
  target_fg_max: number | null
  target_ibu: number | null
  target_ibu_min: number | null
  target_ibu_max: number | null
  target_srm: number | null
  target_srm_min: number | null
  target_srm_max: number | null
  target_carbonation: number | null
  ibu_method: string | null
  brewhouse_efficiency: number | null
  target_abv: number | null // calculated by backend from OG/FG
  created_at: string
  updated_at: string
}

/** Request payload for creating a new recipe */
export interface CreateRecipeRequest {
  name: string
  style_uuid?: string | null
  style_name?: string | null
  notes?: string | null
}

/** Request payload for updating an existing recipe */
export interface UpdateRecipeRequest {
  name: string
  style_uuid?: string | null
  style_name?: string | null
  notes?: string | null
  // Target specifications
  batch_size?: number | null
  batch_size_unit?: string | null
  target_og?: number | null
  target_og_min?: number | null
  target_og_max?: number | null
  target_fg?: number | null
  target_fg_min?: number | null
  target_fg_max?: number | null
  target_ibu?: number | null
  target_ibu_min?: number | null
  target_ibu_max?: number | null
  target_srm?: number | null
  target_srm_min?: number | null
  target_srm_max?: number | null
  target_carbonation?: number | null
  ibu_method?: string | null
  brewhouse_efficiency?: number | null
}

// ============================================================================
// Recipe Ingredient Types
// ============================================================================

/** Type of ingredient in a recipe */
export type RecipeIngredientType
  = | 'fermentable'
    | 'hop'
    | 'yeast'
    | 'adjunct'
    | 'salt'
    | 'chemical'
    | 'gas'
    | 'other'

/** All valid recipe ingredient type values */
export const RECIPE_INGREDIENT_TYPE_VALUES: RecipeIngredientType[] = [
  'fermentable',
  'hop',
  'yeast',
  'adjunct',
  'salt',
  'chemical',
  'gas',
  'other',
]

/** Stage at which an ingredient is used */
export type RecipeUseStage
  = | 'mash'
    | 'boil'
    | 'whirlpool'
    | 'fermentation'
    | 'packaging'

/** All valid recipe use stage values */
export const RECIPE_USE_STAGE_VALUES: RecipeUseStage[] = [
  'mash',
  'boil',
  'whirlpool',
  'fermentation',
  'packaging',
]

/** How an ingredient is used (purpose) */
export type RecipeUseType
  = | 'bittering'
    | 'flavor'
    | 'aroma'
    | 'dry_hop'
    | 'base'
    | 'specialty'
    | 'adjunct'
    | 'sugar'
    | 'primary'
    | 'secondary'
    | 'bottle'
    | 'other'

/** All valid recipe use type values */
export const RECIPE_USE_TYPE_VALUES: RecipeUseType[] = [
  'bittering',
  'flavor',
  'aroma',
  'dry_hop',
  'base',
  'specialty',
  'adjunct',
  'sugar',
  'primary',
  'secondary',
  'bottle',
  'other',
]

/** A recipe ingredient entry */
export interface RecipeIngredient {
  uuid: string
  recipe_uuid: string
  ingredient_uuid: string | null
  ingredient_type: RecipeIngredientType
  name: string
  amount: number
  amount_unit: string
  use_stage: RecipeUseStage
  use_type: RecipeUseType | null
  timing_duration_minutes: number | null
  timing_temperature_c: number | null
  alpha_acid_assumed: number | null
  scaling_factor: number
  sort_order: number
  notes: string | null
  created_at: string
  updated_at: string
}

/** Request payload for creating a recipe ingredient */
export interface CreateRecipeIngredientRequest {
  ingredient_uuid?: string | null
  ingredient_type: RecipeIngredientType
  name: string
  amount: number
  amount_unit: string
  use_stage: RecipeUseStage
  use_type?: RecipeUseType | null
  timing_duration_minutes?: number | null
  timing_temperature_c?: number | null
  alpha_acid_assumed?: number | null
  scaling_factor?: number
  sort_order?: number
  notes?: string | null
}

/** Request payload for updating a recipe ingredient */
export interface UpdateRecipeIngredientRequest {
  ingredient_uuid?: string | null
  ingredient_type?: RecipeIngredientType
  name?: string
  amount?: number
  amount_unit?: string
  use_stage?: RecipeUseStage
  use_type?: RecipeUseType | null
  timing_duration_minutes?: number | null
  timing_temperature_c?: number | null
  alpha_acid_assumed?: number | null
  scaling_factor?: number
  sort_order?: number
  notes?: string | null
}

// ============================================================================
// Batch Types
// ============================================================================

/** A production batch of beer */
export interface Batch {
  uuid: string
  short_name: string
  brew_date: string | null
  recipe_uuid: string | null
  recipe_name: string | null
  current_phase: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

/** Request payload for updating an existing batch */
export interface UpdateBatchRequest {
  short_name: string
  brew_date?: string | null
  notes?: string | null
  recipe_uuid?: string | null
}

// ============================================================================
// Vessel Types
// ============================================================================

/** Vessel operational status */
export type VesselStatus = 'active' | 'inactive' | 'retired'

/** All valid vessel status values */
export const VESSEL_STATUS_VALUES: VesselStatus[] = [
  'active',
  'inactive',
  'retired',
]

/** Vessel type classification */
export type VesselType
  = | 'mash_tun'
    | 'lauter_tun'
    | 'kettle'
    | 'whirlpool'
    | 'fermenter'
    | 'brite_tank'
    | 'serving_tank'
    | 'other'

/** All valid vessel type values */
export const VESSEL_TYPE_VALUES: VesselType[] = [
  'mash_tun',
  'lauter_tun',
  'kettle',
  'whirlpool',
  'fermenter',
  'brite_tank',
  'serving_tank',
  'other',
]

/** A brewing vessel (fermenter, brite tank, kettle, etc.) */
export interface Vessel {
  uuid: string
  type: VesselType
  name: string
  capacity: number
  capacity_unit: VolumeUnit
  make: string | null
  model: string | null
  status: VesselStatus
  created_at: string
  updated_at: string
  deleted_at: string | null
}

/** Request payload for updating an existing vessel */
export interface UpdateVesselRequest {
  name: string
  type: VesselType
  capacity: number
  capacity_unit: VolumeUnit
  make?: string | null
  model?: string | null
  status: VesselStatus
}

// ============================================================================
// Volume Types
// ============================================================================

/** A tracked volume of liquid (wort, beer, etc.) */
export interface Volume {
  uuid: string
  name: string | null
  description: string | null
  amount: number
  amount_unit: VolumeUnit
  created_at: string
  updated_at: string
  deleted_at: string | null
}

/** Request payload for creating a new volume */
export interface CreateVolumeRequest {
  name?: string | null
  description?: string | null
  amount: number
  amount_unit: VolumeUnit
}

// ============================================================================
// Brew Session Types
// ============================================================================

/** A brew session representing a single brew day */
export interface BrewSession {
  uuid: string
  batch_uuid: string | null
  wort_volume_uuid: string | null
  mash_vessel_uuid: string | null
  boil_vessel_uuid: string | null
  brewed_at: string
  notes: string | null
  created_at: string
  updated_at: string
  deleted_at: string | null
}

/** Request payload for creating a new brew session */
export interface CreateBrewSessionRequest {
  batch_uuid?: string | null
  wort_volume_uuid?: string | null
  mash_vessel_uuid?: string | null
  boil_vessel_uuid?: string | null
  brewed_at: string
  notes?: string | null
}

/** Request payload for updating an existing brew session */
export interface UpdateBrewSessionRequest {
  batch_uuid?: string | null
  wort_volume_uuid?: string | null
  mash_vessel_uuid?: string | null
  boil_vessel_uuid?: string | null
  brewed_at: string
  notes?: string | null
}

// ============================================================================
// Addition Types
// ============================================================================

/** Types of additions that can be made during brewing */
export type AdditionType
  = | 'malt'
    | 'hop'
    | 'yeast'
    | 'adjunct'
    | 'water_chem'
    | 'gas'
    | 'other'

/** An ingredient or material addition to a batch or volume */
export interface Addition {
  uuid: string
  batch_uuid: string | null
  occupancy_uuid: string | null
  volume_uuid: string | null
  addition_type: AdditionType
  stage: string | null
  inventory_lot_uuid: string | null
  amount: number
  amount_unit: VolumeUnit
  added_at: string
  notes: string | null
  created_at: string
  updated_at: string
  deleted_at: string | null
}

/** Request payload for creating a new addition */
export interface CreateAdditionRequest {
  batch_uuid?: string | null
  occupancy_uuid?: string | null
  volume_uuid?: string | null
  addition_type: AdditionType
  stage?: string | null
  inventory_lot_uuid?: string | null
  amount: number
  amount_unit: VolumeUnit
  added_at?: string | null
  notes?: string | null
}

// ============================================================================
// Measurement Types
// ============================================================================

/** A measurement observation (temperature, gravity, pH, etc.) */
export interface Measurement {
  uuid: string
  batch_uuid: string | null
  occupancy_uuid: string | null
  volume_uuid: string | null
  kind: string
  value: number
  unit: string | null
  observed_at: string
  notes: string | null
  created_at: string
  updated_at: string
  deleted_at: string | null
}

/** Request payload for creating a new measurement */
export interface CreateMeasurementRequest {
  batch_uuid?: string | null
  occupancy_uuid?: string | null
  volume_uuid?: string | null
  kind: string
  value: number
  unit?: string | null
  observed_at?: string | null
  notes?: string | null
}

// ============================================================================
// Occupancy Types
// ============================================================================

/** Status of a vessel occupancy (what's happening to the beer in the vessel) */
export type OccupancyStatus
  = | 'fermenting'
    | 'conditioning'
    | 'cold_crashing'
    | 'dry_hopping'
    | 'carbonating'
    | 'holding'
    | 'packaging'

/** All valid occupancy status values */
export const OCCUPANCY_STATUS_VALUES: OccupancyStatus[] = [
  'fermenting',
  'conditioning',
  'cold_crashing',
  'dry_hopping',
  'carbonating',
  'holding',
  'packaging',
]

/** Request payload for creating a new occupancy (assigning beer to a vessel) */
export interface CreateOccupancyRequest {
  vessel_uuid: string
  volume_uuid: string
  in_at?: string | null
  status?: OccupancyStatus | null
}

/** A vessel occupancy record tracking what's in a vessel and when */
export interface Occupancy {
  uuid: string
  vessel_uuid: string
  volume_uuid: string
  batch_uuid: string | null
  status: OccupancyStatus | null
  in_at: string
  out_at: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// Transfer Types
// ============================================================================

/** A transfer of liquid between vessel occupancies */
export interface Transfer {
  uuid: string
  source_occupancy_uuid: string
  dest_occupancy_uuid: string
  amount: number
  amount_unit: string
  loss_amount: number | null
  loss_unit: string | null
  started_at: string
  ended_at: string | null
  created_at: string
  updated_at: string
}

/** Request payload for creating a new transfer */
export interface CreateTransferRequest {
  source_occupancy_uuid: string
  dest_vessel_uuid: string
  volume_uuid: string
  amount: number
  amount_unit: string
  loss_amount?: number
  loss_unit?: string
  started_at?: string
  ended_at?: string
  close_source?: boolean
  dest_status?: OccupancyStatus
}

/** Response from creating a transfer, includes the new destination occupancy */
export interface TransferRecordResponse {
  transfer: Transfer
  dest_occupancy: Occupancy
}

// ============================================================================
// Batch Process Phase Types
// ============================================================================

/** Process phase of a batch (e.g. mashing, fermenting, conditioning) */
export type ProcessPhase
  = | 'planning'
    | 'mashing'
    | 'heating'
    | 'boiling'
    | 'cooling'
    | 'fermenting'
    | 'conditioning'
    | 'packaging'
    | 'finished'

/** A batch process phase record tracking when a batch entered a phase */
export interface BatchProcessPhase {
  uuid: string
  batch_uuid: string
  process_phase: ProcessPhase
  phase_at: string
  created_at: string
  updated_at: string
}

/** Request payload for creating a new batch process phase */
export interface CreateBatchProcessPhaseRequest {
  batch_uuid: string
  process_phase: string
  phase_at?: string
}

// ============================================================================
// Batch Volume Types
// ============================================================================

/** Liquid phase classification */
export type LiquidPhase = 'water' | 'wort' | 'beer'

/** A batch volume record linking a batch to a volume with a liquid phase */
export interface BatchVolume {
  uuid: string
  batch_uuid: string
  volume_uuid: string
  liquid_phase: LiquidPhase
  phase_at: string
  created_at: string
  updated_at: string
}

/** Request payload for creating a new batch volume */
export interface CreateBatchVolumeRequest {
  batch_uuid: string
  volume_uuid: string
  liquid_phase: string
  phase_at?: string
}

// ============================================================================
// Volume Relation Types
// ============================================================================

/** Relation type between volumes */
export type VolumeRelationType = 'split' | 'blend'

/** A relation between two volumes (parent → child) */
export interface VolumeRelation {
  uuid: string
  parent_volume_uuid: string
  child_volume_uuid: string
  relation_type: VolumeRelationType
  amount: number
  amount_unit: string
  created_at: string
  updated_at: string
}

/** Request payload for creating a new volume relation */
export interface CreateVolumeRelationRequest {
  parent_volume_uuid: string
  child_volume_uuid: string
  relation_type: VolumeRelationType
  amount: number
  amount_unit: string
}

// ============================================================================
// Batch Relation Types
// ============================================================================

/** Relation type between batches */
export type BatchRelationType = 'split' | 'blend'

/** A relation between two batches (parent → child) */
export interface BatchRelation {
  uuid: string
  parent_batch_uuid: string
  child_batch_uuid: string
  relation_type: BatchRelationType
  volume_uuid: string | null
  created_at: string
  updated_at: string
}

/** Request payload for creating a new batch relation */
export interface CreateBatchRelationRequest {
  parent_batch_uuid: string
  child_batch_uuid: string
  relation_type: BatchRelationType
  volume_uuid?: string
}

// ============================================================================
// Package Format Types
// ============================================================================

/** Container type for package formats */
export type ContainerType = 'keg' | 'can' | 'bottle' | 'cask' | 'growler' | 'other'

/** All valid container type values */
export const CONTAINER_TYPE_VALUES: ContainerType[] = [
  'keg',
  'can',
  'bottle',
  'cask',
  'growler',
  'other',
]

/** A package format definition (container type for packaged beer) */
export interface PackageFormat {
  uuid: string
  name: string
  container: ContainerType
  volume_per_unit: number
  volume_per_unit_unit: string
  is_active: boolean
  created_at: string
  updated_at: string
}

/** Request payload for creating a new package format */
export interface CreatePackageFormatRequest {
  name: string
  container: ContainerType
  volume_per_unit: number
  volume_per_unit_unit: string
}

/** Request payload for updating an existing package format */
export interface UpdatePackageFormatRequest {
  name?: string
  container?: ContainerType
  volume_per_unit?: number
  volume_per_unit_unit?: string
  is_active?: boolean
}

// ============================================================================
// Packaging Run Types
// ============================================================================

/** A line item within a packaging run */
export interface PackagingRunLine {
  uuid: string
  packaging_run_uuid: string
  package_format_uuid: string
  package_format_name: string
  package_format_volume_per_unit: number
  package_format_volume_per_unit_unit: string
  quantity: number
  created_at: string
  updated_at: string
}

/** A packaging run recording a packaging event for a batch */
export interface PackagingRun {
  uuid: string
  batch_uuid: string
  occupancy_uuid: string
  started_at: string
  ended_at: string | null
  loss_amount: number | null
  loss_unit: string | null
  notes: string | null
  lines: PackagingRunLine[]
  created_at: string
  updated_at: string
}

/** Request payload for creating a packaging run line */
export interface CreatePackagingRunLineRequest {
  package_format_uuid: string
  quantity: number
}

/** Request payload for creating a new packaging run */
export interface CreatePackagingRunRequest {
  batch_uuid: string
  occupancy_uuid: string
  started_at?: string
  ended_at?: string
  loss_amount?: number
  loss_unit?: string
  notes?: string
  lines: CreatePackagingRunLineRequest[]
  close_source?: boolean
  stock_location_uuid?: string
  lot_code_prefix?: string
}

// ============================================================================
// Batch Summary Types
// ============================================================================

/** Brew session info within a batch summary */
export interface BatchSummaryBrewSession {
  uuid: string
  brewed_at: string
  notes: string | null
}

/** Comprehensive batch summary with computed metrics */
export interface BatchSummary {
  uuid: string
  short_name: string
  brew_date: string | null
  notes: string | null
  recipe_name: string | null
  style_name: string | null
  brew_sessions: BatchSummaryBrewSession[]
  current_phase: string | null
  current_vessel: string | null
  current_occupancy_status: string | null
  current_occupancy_uuid: string | null
  original_gravity: number | null
  final_gravity: number | null
  abv: number | null
  abv_calculated: boolean
  ibu: number | null
  days_in_fermenter: number | null
  days_in_brite: number | null
  days_grain_to_glass: number | null
  starting_volume_bbl: number | null
  current_volume_bbl: number | null
  total_loss_bbl: number | null
  loss_percentage: number | null
  created_at: string
  updated_at: string
}
