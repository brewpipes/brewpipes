/**
 * Shared types for batch-related components.
 */

export type Unit = 'ml' | 'usfloz' | 'ukfloz' | 'bbl'
export type LiquidPhase = 'water' | 'wort' | 'beer'
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
export type RelationType = 'split' | 'blend'
export type AdditionType
  = | 'malt'
    | 'hop'
    | 'yeast'
    | 'adjunct'
    | 'water_chem'
    | 'gas'
    | 'other'

export type Batch = {
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

export type Volume = {
  uuid: string
  name: string | null
  description: string | null
  amount: number
  amount_unit: Unit
  created_at: string
  updated_at: string
}

export type VolumeRelation = {
  uuid: string
  parent_volume_uuid: string
  child_volume_uuid: string
  relation_type: RelationType
  amount: number
  amount_unit: Unit
  created_at: string
  updated_at: string
}

export type BatchVolume = {
  uuid: string
  batch_uuid: string
  volume_uuid: string
  liquid_phase: LiquidPhase
  phase_at: string
  created_at: string
  updated_at: string
}

export type BatchProcessPhase = {
  uuid: string
  batch_uuid: string
  process_phase: ProcessPhase
  phase_at: string
  created_at: string
  updated_at: string
}

export type Addition = {
  uuid: string
  batch_uuid: string | null
  occupancy_uuid: string | null
  addition_type: AdditionType
  stage: string | null
  inventory_lot_uuid: string | null
  amount: number
  amount_unit: Unit
  added_at: string
  notes: string | null
  created_at: string
  updated_at: string
}

export type Measurement = {
  uuid: string
  batch_uuid: string | null
  occupancy_uuid: string | null
  kind: string
  value: number
  unit: string | null
  observed_at: string
  notes: string | null
  created_at: string
  updated_at: string
}

export type TimelineEvent = {
  id: string
  title: string
  subtitle: string
  at: string
  color: string
  icon: string
}

export type FlowNode = {
  id: string
  label: string
}

export type FlowLink = {
  source: string
  target: string
  value: number
  label: string
}

export type SparklineSeries = {
  values: number[]
  latest: Measurement | null
  latestLabel: string
  linePath: string
  areaPath: string
}
