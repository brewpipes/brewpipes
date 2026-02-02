import { ref, watch } from 'vue'
import {
  colorLabels,
  type ColorUnit,
  formatColor,
  formatGravity,
  formatMass,
  formatPressure,
  formatTemperature,
  formatVolume,
  gravityLabels,
  type GravityUnit,
  massLabels,
  type MassUnit,
  pressureLabels,
  type PressureUnit,
  temperatureLabels,
  type TemperatureUnit,
  volumeLabels,
  type VolumeUnit,
} from './useUnitConversion'

// ==================== Unit Type Detection ====================
// Sets for detecting unit categories from string values (e.g., from backend data)

const VOLUME_UNITS = new Set<string>([
  'ml', 'l', 'hl', 'usfloz', 'ukfloz', 'usgal', 'ukgal', 'bbl', 'ukbbl',
  // Common shorthand variants from backend
  'gal',
])

const MASS_UNITS = new Set<string>(['g', 'kg', 'oz', 'lb'])

/**
 * Check if a unit string represents a volume unit.
 */
export function isVolumeUnit (unit: string): boolean {
  return VOLUME_UNITS.has(unit.toLowerCase())
}

/**
 * Check if a unit string represents a mass unit.
 */
export function isMassUnit (unit: string): boolean {
  return MASS_UNITS.has(unit.toLowerCase())
}

/**
 * Normalize a unit string to a typed VolumeUnit.
 * Maps common shorthand variants (e.g., 'gal' -> 'usgal').
 */
export function normalizeVolumeUnit (unit: string): VolumeUnit {
  const lower = unit.toLowerCase()
  if (lower === 'gal') {
    return 'usgal'
  }
  return lower as VolumeUnit
}

/**
 * Normalize a unit string to a typed MassUnit.
 */
export function normalizeMassUnit (unit: string): MassUnit {
  return unit.toLowerCase() as MassUnit
}

// Re-export unit types for convenience

export interface UnitPreferences {
  temperature: TemperatureUnit
  gravity: GravityUnit
  volume: VolumeUnit
  mass: MassUnit
  pressure: PressureUnit
  color: ColorUnit
}

const STORAGE_KEY = 'brewpipes:unitPreferences'

// US-centric defaults (first customer is US-based)
const DEFAULT_PREFERENCES: UnitPreferences = {
  temperature: 'f',
  gravity: 'sg',
  volume: 'bbl',
  mass: 'lb',
  pressure: 'psi',
  color: 'srm',
}

// Options arrays for building select dropdowns
export const temperatureOptions: Array<{ value: TemperatureUnit, label: string }> = [
  { value: 'c', label: 'Celsius (°C)' },
  { value: 'f', label: 'Fahrenheit (°F)' },
]

export const gravityOptions: Array<{ value: GravityUnit, label: string }> = [
  { value: 'sg', label: 'Specific Gravity (SG)' },
  { value: 'plato', label: 'Degrees Plato (°P)' },
]

export const volumeOptions: Array<{ value: VolumeUnit, label: string }> = [
  { value: 'ml', label: 'Milliliters (mL)' },
  { value: 'l', label: 'Liters (L)' },
  { value: 'hl', label: 'Hectoliters (hL)' },
  { value: 'usfloz', label: 'Fluid Ounces (US)' },
  { value: 'ukfloz', label: 'Fluid Ounces (UK)' },
  { value: 'usgal', label: 'Gallons (US)' },
  { value: 'ukgal', label: 'Gallons (UK)' },
  { value: 'bbl', label: 'Barrels (US)' },
  { value: 'ukbbl', label: 'Barrels (UK)' },
]

export const massOptions: Array<{ value: MassUnit, label: string }> = [
  { value: 'g', label: 'Grams (g)' },
  { value: 'kg', label: 'Kilograms (kg)' },
  { value: 'oz', label: 'Ounces (oz)' },
  { value: 'lb', label: 'Pounds (lb)' },
]

export const pressureOptions: Array<{ value: PressureUnit, label: string }> = [
  { value: 'kpa', label: 'Kilopascals (kPa)' },
  { value: 'psi', label: 'PSI' },
  { value: 'bar', label: 'Bar' },
]

export const colorOptions: Array<{ value: ColorUnit, label: string }> = [
  { value: 'srm', label: 'SRM' },
  { value: 'ebc', label: 'EBC' },
]

// Module-level singleton state
const preferences = ref<UnitPreferences>(loadPreferences())

/**
 * Load preferences from localStorage, falling back to defaults on error.
 */
function loadPreferences (): UnitPreferences {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      const parsed = JSON.parse(stored) as Partial<UnitPreferences>
      // Merge with defaults to handle missing or invalid keys
      return {
        temperature: isValidTemperatureUnit(parsed.temperature)
          ? parsed.temperature
          : DEFAULT_PREFERENCES.temperature,
        gravity: isValidGravityUnit(parsed.gravity)
          ? parsed.gravity
          : DEFAULT_PREFERENCES.gravity,
        volume: isValidVolumeUnit(parsed.volume)
          ? parsed.volume
          : DEFAULT_PREFERENCES.volume,
        mass: isValidMassUnit(parsed.mass)
          ? parsed.mass
          : DEFAULT_PREFERENCES.mass,
        pressure: isValidPressureUnit(parsed.pressure)
          ? parsed.pressure
          : DEFAULT_PREFERENCES.pressure,
        color: isValidColorUnit(parsed.color)
          ? parsed.color
          : DEFAULT_PREFERENCES.color,
      }
    }
  } catch {
    // localStorage unavailable or JSON parse error - use defaults
  }
  return { ...DEFAULT_PREFERENCES }
}

/**
 * Save preferences to localStorage.
 */
function savePreferences (prefs: UnitPreferences): void {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(prefs))
  } catch {
    // localStorage unavailable - silently ignore
  }
}

// Validation helpers
function isValidTemperatureUnit (value: unknown): value is TemperatureUnit {
  return value === 'c' || value === 'f'
}

function isValidGravityUnit (value: unknown): value is GravityUnit {
  return value === 'sg' || value === 'plato'
}

function isValidVolumeUnit (value: unknown): value is VolumeUnit {
  return ['ml', 'l', 'hl', 'usfloz', 'ukfloz', 'usgal', 'ukgal', 'bbl', 'ukbbl'].includes(
    value as string,
  )
}

function isValidMassUnit (value: unknown): value is MassUnit {
  return ['g', 'kg', 'oz', 'lb'].includes(value as string)
}

function isValidPressureUnit (value: unknown): value is PressureUnit {
  return value === 'kpa' || value === 'psi' || value === 'bar'
}

function isValidColorUnit (value: unknown): value is ColorUnit {
  return value === 'srm' || value === 'ebc'
}

// Watch for changes and persist to localStorage
watch(
  preferences,
  newPrefs => {
    savePreferences(newPrefs)
  },
  { deep: true },
)

/**
 * Composable for managing user unit preferences with localStorage persistence.
 * Returns singleton state shared across all component instances.
 */
export function useUnitPreferences () {
  // Setters for individual preferences
  function setTemperatureUnit (unit: TemperatureUnit): void {
    preferences.value = { ...preferences.value, temperature: unit }
  }

  function setGravityUnit (unit: GravityUnit): void {
    preferences.value = { ...preferences.value, gravity: unit }
  }

  function setVolumeUnit (unit: VolumeUnit): void {
    preferences.value = { ...preferences.value, volume: unit }
  }

  function setMassUnit (unit: MassUnit): void {
    preferences.value = { ...preferences.value, mass: unit }
  }

  function setPressureUnit (unit: PressureUnit): void {
    preferences.value = { ...preferences.value, pressure: unit }
  }

  function setColorUnit (unit: ColorUnit): void {
    preferences.value = { ...preferences.value, color: unit }
  }

  /**
   * Reset all preferences to US-centric defaults.
   */
  function resetToDefaults (): void {
    preferences.value = { ...DEFAULT_PREFERENCES }
  }

  // Convenience formatting functions that use current preferences

  /**
   * Format temperature value to user's preferred unit.
   */
  function formatTemperaturePreferred (
    value: number | null | undefined,
    fromUnit: TemperatureUnit,
  ): string {
    return formatTemperature(value, fromUnit, preferences.value.temperature)
  }

  /**
   * Format gravity value to user's preferred unit.
   */
  function formatGravityPreferred (
    value: number | null | undefined,
    fromUnit: GravityUnit,
  ): string {
    return formatGravity(value, fromUnit, preferences.value.gravity)
  }

  /**
   * Format volume value to user's preferred unit.
   */
  function formatVolumePreferred (
    value: number | null | undefined,
    fromUnit: VolumeUnit,
  ): string {
    return formatVolume(value, fromUnit, preferences.value.volume)
  }

  /**
   * Format mass value to user's preferred unit.
   */
  function formatMassPreferred (
    value: number | null | undefined,
    fromUnit: MassUnit,
  ): string {
    return formatMass(value, fromUnit, preferences.value.mass)
  }

  /**
   * Format pressure value to user's preferred unit.
   */
  function formatPressurePreferred (
    value: number | null | undefined,
    fromUnit: PressureUnit,
  ): string {
    return formatPressure(value, fromUnit, preferences.value.pressure)
  }

  /**
   * Format color value to user's preferred unit.
   */
  function formatColorPreferred (
    value: number | null | undefined,
    fromUnit: ColorUnit,
  ): string {
    return formatColor(value, fromUnit, preferences.value.color)
  }

  /**
   * Smart formatter that detects unit type and converts to user preferences.
   * Useful for inventory data where the unit type is stored as a string.
   *
   * @param value - The numeric value to format
   * @param unit - The source unit as a string (e.g., 'kg', 'gal', 'lb')
   * @returns Formatted string with value converted to preferred unit, or fallback
   */
  function formatAmountPreferred (
    value: number | null | undefined,
    unit: string | null | undefined,
  ): string {
    if (value === null || value === undefined) {
      return '—'
    }
    if (!unit) {
      return `${value}`
    }

    const lowerUnit = unit.toLowerCase()
    if (isVolumeUnit(lowerUnit)) {
      return formatVolumePreferred(value, normalizeVolumeUnit(lowerUnit))
    }
    if (isMassUnit(lowerUnit)) {
      return formatMassPreferred(value, normalizeMassUnit(lowerUnit))
    }
    // Fallback for unknown units - display as-is
    return `${value} ${unit}`
  }

  return {
    // Reactive preferences (singleton)
    preferences,

    // Setters
    setTemperatureUnit,
    setGravityUnit,
    setVolumeUnit,
    setMassUnit,
    setPressureUnit,
    setColorUnit,
    resetToDefaults,

    // Convenience formatting functions
    formatTemperaturePreferred,
    formatGravityPreferred,
    formatVolumePreferred,
    formatMassPreferred,
    formatPressurePreferred,
    formatColorPreferred,
    formatAmountPreferred,

    // Unit type detection helpers
    isVolumeUnit,
    isMassUnit,
    normalizeVolumeUnit,
    normalizeMassUnit,

    // Options for select dropdowns
    temperatureOptions,
    gravityOptions,
    volumeOptions,
    massOptions,
    pressureOptions,
    colorOptions,

    // Labels (re-exported from useUnitConversion)
    temperatureLabels,
    gravityLabels,
    volumeLabels,
    massLabels,
    pressureLabels,
    colorLabels,
  }
}

export { type ColorUnit, type GravityUnit, type MassUnit, type PressureUnit, type TemperatureUnit, type VolumeUnit } from './useUnitConversion'
