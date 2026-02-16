import type { RemovalCategory, RemovalReason } from '@/types'

/** Human-readable labels for removal categories. */
export const categoryLabels: Record<RemovalCategory, string> = {
  dump: 'Batch Dump',
  waste: 'Waste / Spillage',
  sample: 'Sample Pull',
  expired: 'Expired / Destroyed',
  other: 'Other',
}

/** Vuetify color tokens for removal categories. */
export const categoryColors: Record<RemovalCategory, string> = {
  dump: 'error',
  waste: 'warning',
  sample: 'info',
  expired: 'amber',
  other: 'grey',
}

/** MDI icon names for removal categories. */
export const categoryIcons: Record<RemovalCategory, string> = {
  dump: 'mdi-delete-variant',
  waste: 'mdi-water-off',
  sample: 'mdi-test-tube',
  expired: 'mdi-calendar-remove',
  other: 'mdi-dots-horizontal',
}

/** Human-readable labels for removal reasons. */
export const reasonLabels: Record<RemovalReason, string> = {
  infection: 'Infection',
  off_flavor: 'Off-Flavor',
  failed_fermentation: 'Failed Fermentation',
  equipment_failure: 'Equipment Failure',
  quality_reject: 'Quality Reject',
  past_date: 'Past Date',
  damaged_package: 'Damaged Package',
  spillage: 'Spillage',
  cleaning: 'Cleaning Waste',
  qc_sample: 'QC Sample',
  tasting: 'Tasting',
  competition: 'Competition Entry',
  other: 'Other',
}

/** Valid reasons for each removal category. */
export const reasonsByCategory: Record<RemovalCategory, RemovalReason[]> = {
  dump: ['infection', 'off_flavor', 'failed_fermentation', 'equipment_failure', 'quality_reject', 'other'],
  waste: ['spillage', 'cleaning', 'other'],
  sample: ['qc_sample', 'tasting', 'competition', 'other'],
  expired: ['past_date', 'damaged_package', 'quality_reject', 'other'],
  other: ['other'],
}
