/**
 * Shared types for batch-related components.
 *
 * Entity types (Batch, Volume, Addition, Measurement, BatchProcessPhase,
 * BatchVolume, VolumeRelation, etc.) are imported from '@/types'. This file
 * contains only component-specific types used by batch UI components and
 * re-exports canonical types for barrel-file consumers.
 */

import type { Measurement } from '@/types'

// Re-export entity types from canonical source for barrel-file consumers
export type {
  Addition,
  AdditionType,
  Batch,
  BatchProcessPhase,
  BatchVolume,
  LiquidPhase,
  Measurement,
  ProcessPhase,
  Volume,
  VolumeRelation,
  VolumeRelationType,
} from '@/types'

// Aliases for VolumeUnit used in batch component forms
export type Unit = import('@/types').VolumeUnit

// Re-export relation type alias for backwards compatibility
export type RelationType = import('@/types').VolumeRelationType

// ============================================================================
// Component-specific types (not duplicated in @/types)
// ============================================================================

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

// ============================================================================
// Measurement kind options (shared across measurement dialogs)
// ============================================================================

export type MeasurementKindOption = {
  title: string
  value: string
}

/**
 * Curated list of measurement kinds for the batch-level measurement dialog.
 * Covers fermentation tracking, quality checks, and general observations.
 */
export const measurementKinds: MeasurementKindOption[] = [
  { title: 'Temperature', value: 'temperature' },
  { title: 'Gravity', value: 'gravity' },
  { title: 'pH', value: 'ph' },
  { title: 'Pressure', value: 'pressure' },
  { title: 'Volume', value: 'volume' },
  { title: 'Dissolved Oxygen', value: 'dissolved_oxygen' },
  { title: 'ABV', value: 'abv' },
  { title: 'Note', value: 'note' },
  { title: 'Other', value: 'other' },
]

/**
 * Curated list of measurement kinds for the hot-side (brew day) measurement dialog.
 * These are more specific to mash, boil, and pre-fermentation stages.
 */
export const hotSideMeasurementKinds: MeasurementKindOption[] = [
  { title: 'Mash Temperature', value: 'mash_temp' },
  { title: 'Mash pH', value: 'mash_ph' },
  { title: 'Pre-Boil Gravity', value: 'pre_boil_gravity' },
  { title: 'Original Gravity', value: 'original_gravity' },
  { title: 'Boil Temperature', value: 'boil_temp' },
  { title: 'Post-Boil Volume', value: 'post_boil_volume' },
  { title: 'Other', value: 'other' },
]

/**
 * Get a sensible default unit for a given measurement kind.
 * Returns an empty string if no default is known.
 */
export function getDefaultUnitForKind (kind: string): string {
  switch (kind) {
    case 'temperature':
    case 'mash_temp':
    case 'boil_temp': {
      return 'F'
    }
    case 'ph':
    case 'mash_ph': {
      return 'pH'
    }
    case 'gravity':
    case 'pre_boil_gravity':
    case 'original_gravity': {
      return 'SG'
    }
    case 'volume':
    case 'post_boil_volume': {
      return 'bbl'
    }
    case 'pressure': {
      return 'PSI'
    }
    case 'dissolved_oxygen': {
      return 'ppb'
    }
    case 'abv': {
      return '%'
    }
    default: {
      return ''
    }
  }
}
