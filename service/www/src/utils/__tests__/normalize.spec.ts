import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import {
  normalizeDateTime,
  normalizeDateOnly,
  normalizeText,
  nowInputValue,
  toLocalDateTimeInput,
  toNumber,
} from '../normalize'

describe('normalizeText', () => {
  it('returns trimmed string for non-empty input', () => {
    expect(normalizeText('hello')).toBe('hello')
  })

  it('trims leading and trailing whitespace', () => {
    expect(normalizeText('  hello  ')).toBe('hello')
  })

  it('trims tabs and newlines', () => {
    expect(normalizeText('\thello\n')).toBe('hello')
  })

  it('returns null for empty string', () => {
    expect(normalizeText('')).toBeNull()
  })

  it('returns null for whitespace-only string', () => {
    expect(normalizeText('   ')).toBeNull()
  })

  it('returns null for tab-only string', () => {
    expect(normalizeText('\t\t')).toBeNull()
  })

  it('returns null for newline-only string', () => {
    expect(normalizeText('\n\n')).toBeNull()
  })

  it('preserves internal whitespace', () => {
    expect(normalizeText('  hello world  ')).toBe('hello world')
  })

  it('handles single character', () => {
    expect(normalizeText('a')).toBe('a')
  })

  it('handles unicode characters', () => {
    expect(normalizeText('  Brauerei Munchen  ')).toBe('Brauerei Munchen')
  })
})

describe('normalizeDateTime', () => {
  it('converts a datetime string to ISO format', () => {
    const result = normalizeDateTime('2024-06-15T14:30:00')
    expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
    expect(result).toContain('2024')
  })

  it('converts an ISO datetime string', () => {
    const result = normalizeDateTime('2024-06-15T14:30:00Z')
    expect(result).toBe('2024-06-15T14:30:00.000Z')
  })

  it('returns null for empty string', () => {
    expect(normalizeDateTime('')).toBeNull()
  })

  it('handles date-only string', () => {
    const result = normalizeDateTime('2024-06-15')
    expect(result).toBeDefined()
    expect(result).toContain('2024')
  })

  it('returns ISO string for valid date input', () => {
    const result = normalizeDateTime('2024-01-01T00:00:00Z')
    expect(result).toBe('2024-01-01T00:00:00.000Z')
  })
})

describe('normalizeDateOnly', () => {
  it('converts a date string to ISO format at midnight UTC', () => {
    const result = normalizeDateOnly('2024-06-15')
    expect(result).toBe('2024-06-15T00:00:00.000Z')
  })

  it('returns null for empty string', () => {
    expect(normalizeDateOnly('')).toBeNull()
  })

  it('handles first day of year', () => {
    const result = normalizeDateOnly('2024-01-01')
    expect(result).toBe('2024-01-01T00:00:00.000Z')
  })

  it('handles last day of year', () => {
    const result = normalizeDateOnly('2024-12-31')
    expect(result).toBe('2024-12-31T00:00:00.000Z')
  })

  it('handles leap year date', () => {
    const result = normalizeDateOnly('2024-02-29')
    expect(result).toBe('2024-02-29T00:00:00.000Z')
  })
})

describe('toNumber', () => {
  it('parses a numeric string to a number', () => {
    expect(toNumber('42')).toBe(42)
  })

  it('parses a decimal string', () => {
    expect(toNumber('3.14')).toBe(3.14)
  })

  it('parses negative numbers', () => {
    expect(toNumber('-10')).toBe(-10)
  })

  it('returns the number itself when given a number', () => {
    expect(toNumber(42)).toBe(42)
  })

  it('returns zero for zero string', () => {
    expect(toNumber('0')).toBe(0)
  })

  it('returns zero for number zero', () => {
    expect(toNumber(0)).toBe(0)
  })

  it('returns null for null', () => {
    expect(toNumber(null)).toBeNull()
  })

  it('returns null for empty string', () => {
    expect(toNumber('')).toBeNull()
  })

  it('returns null for non-numeric string', () => {
    expect(toNumber('abc')).toBeNull()
  })

  it('returns null for NaN-producing input', () => {
    expect(toNumber('not a number')).toBeNull()
  })

  it('returns null for Infinity string', () => {
    expect(toNumber('Infinity')).toBeNull()
  })

  it('returns null for -Infinity string', () => {
    expect(toNumber('-Infinity')).toBeNull()
  })

  it('parses string with leading/trailing whitespace', () => {
    expect(toNumber('  42  ')).toBe(42)
  })

  it('parses very small decimal', () => {
    expect(toNumber('0.001')).toBe(0.001)
  })

  it('parses scientific notation', () => {
    expect(toNumber('1e3')).toBe(1000)
  })
})

describe('nowInputValue', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('returns a string in datetime-local format', () => {
    vi.setSystemTime(new Date(2024, 5, 15, 14, 30)) // June 15, 2024 14:30
    const result = nowInputValue()
    expect(result).toBe('2024-06-15T14:30')
  })

  it('pads single-digit month and day', () => {
    vi.setSystemTime(new Date(2024, 0, 5, 9, 5)) // Jan 5, 2024 09:05
    const result = nowInputValue()
    expect(result).toBe('2024-01-05T09:05')
  })

  it('handles midnight', () => {
    vi.setSystemTime(new Date(2024, 0, 1, 0, 0)) // Jan 1, 2024 00:00
    const result = nowInputValue()
    expect(result).toBe('2024-01-01T00:00')
  })

  it('handles end of day', () => {
    vi.setSystemTime(new Date(2024, 11, 31, 23, 59)) // Dec 31, 2024 23:59
    const result = nowInputValue()
    expect(result).toBe('2024-12-31T23:59')
  })

  it('matches datetime-local input format pattern', () => {
    vi.setSystemTime(new Date(2024, 5, 15, 14, 30))
    const result = nowInputValue()
    expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/)
  })
})

describe('toLocalDateTimeInput', () => {
  it('converts an ISO string to datetime-local format', () => {
    // Use a UTC time and check the local conversion
    const result = toLocalDateTimeInput('2024-06-15T14:30:00.000Z')
    expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/)
    expect(result).toContain('2024')
  })

  it('returns empty string for empty input', () => {
    expect(toLocalDateTimeInput('')).toBe('')
  })

  it('matches datetime-local input format pattern', () => {
    const result = toLocalDateTimeInput('2024-01-01T00:00:00.000Z')
    expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/)
  })

  it('pads single-digit values', () => {
    // January 5 at midnight UTC
    const result = toLocalDateTimeInput('2024-01-05T00:00:00.000Z')
    expect(result).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/)
    // Verify no single-digit segments
    const parts = result.split(/[-T:]/)
    expect(parts[0]).toHaveLength(4) // year
    expect(parts[1]).toHaveLength(2) // month
    expect(parts[2]).toHaveLength(2) // day
    expect(parts[3]).toHaveLength(2) // hours
    expect(parts[4]).toHaveLength(2) // minutes
  })
})
