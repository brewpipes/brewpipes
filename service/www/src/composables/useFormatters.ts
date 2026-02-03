import type { OccupancyStatus, VesselStatus, VesselType } from '@/types'

// Re-export VesselStatus for backward compatibility
export type { VesselStatus, VesselType } from '@/types'

/**
 * Shared formatting utilities for dates, times, and domain-specific values.
 */
export function useFormatters () {
  /**
   * Format a date/time string to a localized medium date and short time.
   */
  function formatDateTime (value: string | null | undefined): string {
    if (!value) {
      return 'Unknown'
    }
    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(new Date(value))
  }

  /**
   * Format a date string to a localized medium date.
   */
  function formatDate (value: string | null | undefined): string {
    if (!value) {
      return 'Unknown'
    }
    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
    }).format(new Date(value))
  }

  /**
   * Format a date/time string to a relative time (e.g., "5m ago", "2d ago").
   */
  function formatRelativeTime (value: string | null | undefined): string {
    if (!value) {
      return 'Unknown'
    }

    const date = new Date(value)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffSecs = Math.floor(diffMs / 1000)
    const diffMins = Math.floor(diffSecs / 60)
    const diffHours = Math.floor(diffMins / 60)
    const diffDays = Math.floor(diffHours / 24)

    if (diffSecs < 60) {
      return 'just now'
    }
    if (diffMins < 60) {
      return `${diffMins}m ago`
    }
    if (diffHours < 24) {
      return `${diffHours}h ago`
    }
    if (diffDays < 7) {
      return `${diffDays}d ago`
    }

    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
    }).format(date)
  }

  return {
    formatDateTime,
    formatDate,
    formatRelativeTime,
  }
}

// Vessel status formatting
const VESSEL_STATUS_LABELS: Record<VesselStatus, string> = {
  active: 'Active',
  inactive: 'Inactive',
  retired: 'Retired',
}

const VESSEL_STATUS_COLORS: Record<VesselStatus, string> = {
  active: 'success',
  inactive: 'grey',
  retired: 'error',
}

export function useVesselStatusFormatters () {
  function formatVesselStatus (status: VesselStatus): string {
    return VESSEL_STATUS_LABELS[status] ?? status
  }

  function getVesselStatusColor (status: VesselStatus): string {
    return VESSEL_STATUS_COLORS[status] ?? 'secondary'
  }

  return {
    formatVesselStatus,
    getVesselStatusColor,
  }
}

// Vessel type formatting
const VESSEL_TYPE_LABELS: Record<VesselType, string> = {
  mash_tun: 'Mash Tun',
  lauter_tun: 'Lauter Tun',
  kettle: 'Kettle',
  whirlpool: 'Whirlpool',
  fermenter: 'Fermenter',
  brite_tank: 'Brite Tank',
  serving_tank: 'Serving Tank',
  other: 'Other',
}

export function useVesselTypeFormatters () {
  function formatVesselType (type: VesselType | string): string {
    return VESSEL_TYPE_LABELS[type as VesselType] ?? type.charAt(0).toUpperCase() + type.slice(1).replace(/_/g, ' ')
  }

  return {
    formatVesselType,
  }
}

// Occupancy status formatting
const OCCUPANCY_STATUS_LABELS: Record<string, string> = {
  fermenting: 'Fermenting',
  conditioning: 'Conditioning',
  cold_crashing: 'Cold Crashing',
  dry_hopping: 'Dry Hopping',
  carbonating: 'Carbonating',
  holding: 'Holding',
  packaging: 'Packaging',
}

const OCCUPANCY_STATUS_COLORS: Record<string, string> = {
  fermenting: 'orange',
  conditioning: 'blue',
  cold_crashing: 'cyan',
  dry_hopping: 'green',
  carbonating: 'purple',
  holding: 'grey',
  packaging: 'teal',
}

const OCCUPANCY_STATUS_ICONS: Record<string, string> = {
  fermenting: 'mdi-molecule',
  conditioning: 'mdi-clock-outline',
  cold_crashing: 'mdi-snowflake',
  dry_hopping: 'mdi-leaf',
  carbonating: 'mdi-shimmer',
  holding: 'mdi-pause-circle-outline',
  packaging: 'mdi-package-variant',
}

export function useOccupancyStatusFormatters () {
  function formatOccupancyStatus (status: OccupancyStatus | string | null | undefined): string {
    if (!status) {
      return 'No status'
    }
    return OCCUPANCY_STATUS_LABELS[status] ?? status.charAt(0).toUpperCase() + status.slice(1).replace(/_/g, ' ')
  }

  function getOccupancyStatusColor (status: OccupancyStatus | string | null | undefined): string {
    if (!status) {
      return 'grey'
    }
    return OCCUPANCY_STATUS_COLORS[status] ?? 'secondary'
  }

  function getOccupancyStatusIcon (status: OccupancyStatus | string | null | undefined): string {
    if (!status) {
      return 'mdi-help-circle-outline'
    }
    return OCCUPANCY_STATUS_ICONS[status] ?? 'mdi-circle'
  }

  return {
    formatOccupancyStatus,
    getOccupancyStatusColor,
    getOccupancyStatusIcon,
  }
}
