/**
 * Shared types for batch-related components.
 *
 * Entity types (Batch, Volume, Addition, Measurement, etc.) are imported from
 * '@/types'. This file contains only component-specific types used by batch
 * UI components.
 */

import type { Measurement } from '@/types'

// Re-export entity types from canonical source for barrel-file consumers
export type {
  Addition,
  AdditionType,
  Batch,
  Measurement,
  Volume,
} from '@/types'

// Aliases for VolumeUnit used in batch component forms
export type Unit = import('@/types').VolumeUnit

// ============================================================================
// Component-specific types (not duplicated in @/types)
// ============================================================================

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
