/**
 * Common types shared across the application.
 */

/**
 * Standard identifier fields present on most entities.
 */
export interface EntityIdentifiers {
  id: number
  uuid: string
}

/**
 * Standard timestamp fields present on most entities.
 */
export interface EntityTimestamps {
  created_at: string
  updated_at: string
}

/**
 * Soft-delete timestamp field for entities that support soft deletion.
 */
export interface SoftDeletable {
  deleted_at: string | null
}

/**
 * Base entity type combining identifiers and timestamps.
 */
export type BaseEntity = EntityIdentifiers & EntityTimestamps
