/**
 * Shared normalize/conversion utilities used across composables and components.
 */

/**
 * Trim a string and return null if empty.
 */
export function normalizeText (value: string): string | null {
  const trimmed = value.trim()
  return trimmed.length > 0 ? trimmed : null
}

/**
 * Convert a datetime string to ISO format, or return null if empty.
 */
export function normalizeDateTime (value: string): string | null {
  return value ? new Date(value).toISOString() : null
}

/**
 * Convert a date-only string (YYYY-MM-DD) to ISO format, or return null if empty.
 */
export function normalizeDateOnly (value: string): string | null {
  return value ? new Date(`${value}T00:00:00Z`).toISOString() : null
}

/**
 * Parse a string or number to a finite number, or return null.
 */
export function toNumber (value: string | number | null): number | null {
  if (value === null || value === undefined || value === '') {
    return null
  }
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : null
}

/**
 * Return the current local date/time as a value suitable for datetime-local inputs.
 */
export function nowInputValue (): string {
  const now = new Date()
  const pad = (v: number) => String(v).padStart(2, '0')
  const year = now.getFullYear()
  const month = pad(now.getMonth() + 1)
  const day = pad(now.getDate())
  const hours = pad(now.getHours())
  const minutes = pad(now.getMinutes())
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

/**
 * Convert an ISO date/time string to a value suitable for datetime-local inputs.
 */
export function toLocalDateTimeInput (isoString: string): string {
  if (!isoString) {
    return ''
  }
  const date = new Date(isoString)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}
