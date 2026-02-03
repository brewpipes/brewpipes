/**
 * Production domain types for batches, vessels, recipes, and brewing operations.
 */

import type { VolumeUnit } from './units'

// ============================================================================
// Style Types
// ============================================================================

/** Beer style definition */
export interface Style {
  id: number
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

/** Recipe definition with optional style association */
export interface Recipe {
  id: number
  uuid: string
  name: string
  style_id: number | null
  style_name: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

/** Request payload for creating a new recipe */
export interface CreateRecipeRequest {
  name: string
  style_id?: number | null
  style_name?: string | null
  notes?: string | null
}

/** Request payload for updating an existing recipe */
export interface UpdateRecipeRequest {
  name: string
  style_id?: number | null
  style_name?: string | null
  notes?: string | null
}

// ============================================================================
// Batch Types
// ============================================================================

/** A production batch of beer */
export interface Batch {
  id: number
  uuid: string
  short_name: string
  brew_date: string | null
  recipe_id: number | null
  notes: string | null
  created_at: string
  updated_at: string
}

/** Request payload for updating an existing batch */
export interface UpdateBatchRequest {
  short_name: string
  brew_date?: string | null
  notes?: string | null
  recipe_id?: number | null
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
export type VesselType =
  | 'mash_tun'
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
  id: number
  uuid: string
  type: string
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
  type: string
  capacity: number
  capacity_unit: string
  make?: string | null
  model?: string | null
  status: VesselStatus
}

// ============================================================================
// Volume Types
// ============================================================================

/** A tracked volume of liquid (wort, beer, etc.) */
export interface Volume {
  id: number
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
  id: number
  uuid: string
  batch_id: number | null
  wort_volume_id: number | null
  mash_vessel_id: number | null
  boil_vessel_id: number | null
  brewed_at: string
  notes: string | null
  created_at: string
  updated_at: string
  deleted_at: string | null
}

/** Request payload for creating a new brew session */
export interface CreateBrewSessionRequest {
  batch_id?: number | null
  wort_volume_id?: number | null
  mash_vessel_id?: number | null
  boil_vessel_id?: number | null
  brewed_at: string
  notes?: string | null
}

/** Request payload for updating an existing brew session */
export interface UpdateBrewSessionRequest {
  batch_id?: number | null
  wort_volume_id?: number | null
  mash_vessel_id?: number | null
  boil_vessel_id?: number | null
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
  id: number
  uuid: string
  batch_id: number | null
  occupancy_id: number | null
  volume_id: number | null
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
  batch_id?: number | null
  occupancy_id?: number | null
  volume_id?: number | null
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
  id: number
  uuid: string
  batch_id: number | null
  occupancy_id: number | null
  volume_id: number | null
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
  batch_id?: number | null
  occupancy_id?: number | null
  volume_id?: number | null
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

/** A vessel occupancy record tracking what's in a vessel and when */
export interface Occupancy {
  id: number
  uuid: string
  vessel_id: number
  volume_id: number
  batch_id: number | null
  status: OccupancyStatus | null
  in_at: string
  out_at: string | null
  created_at: string
  updated_at: string
}

// ============================================================================
// Batch Summary Types
// ============================================================================

/** Brew session info within a batch summary */
export interface BatchSummaryBrewSession {
  id: number
  brewed_at: string
  notes: string | null
}

/** Comprehensive batch summary with computed metrics */
export interface BatchSummary {
  id: number
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
  current_occupancy_id: number | null
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
