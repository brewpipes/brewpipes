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
