import type {
  ColorUnit,
  GravityUnit,
  MassUnit,
  PressureUnit,
  TemperatureUnit,
  VolumeUnit,
} from '@/types'

// Re-export unit types for backward compatibility
export type {
  ColorUnit,
  GravityUnit,
  MassUnit,
  PressureUnit,
  TemperatureUnit,
  VolumeUnit,
} from '@/types'

// Unit labels for display
export const temperatureLabels: Record<TemperatureUnit, string> = {
  c: '°C',
  f: '°F',
}

export const gravityLabels: Record<GravityUnit, string> = {
  sg: 'SG',
  plato: '°P',
}

export const volumeLabels: Record<VolumeUnit, string> = {
  ml: 'mL',
  l: 'L',
  hl: 'hL',
  usfloz: 'fl oz',
  ukfloz: 'fl oz (UK)',
  usgal: 'gal',
  ukgal: 'gal (UK)',
  bbl: 'bbl',
  ukbbl: 'bbl (UK)',
}

export const massLabels: Record<MassUnit, string> = {
  g: 'g',
  kg: 'kg',
  oz: 'oz',
  lb: 'lb',
}

export const pressureLabels: Record<PressureUnit, string> = {
  kpa: 'kPa',
  psi: 'psi',
  bar: 'bar',
}

export const colorLabels: Record<ColorUnit, string> = {
  srm: 'SRM',
  ebc: 'EBC',
}

// Base conversion factors (to base unit)
// Volume: base unit is mL
const volumeToMl: Record<VolumeUnit, number> = {
  ml: 1,
  l: 1000,
  hl: 100_000,
  usfloz: 29.5735,
  ukfloz: 28.4131,
  usgal: 3785.41,
  ukgal: 4546.09,
  bbl: 117_347.77, // 31 * 3785.41
  ukbbl: 163_659.24, // 36 * 4546.09
}

// Mass: base unit is grams
const massToGrams: Record<MassUnit, number> = {
  g: 1,
  kg: 1000,
  oz: 28.3495,
  lb: 453.592,
}

// Pressure: base unit is kPa
const pressureToKpa: Record<PressureUnit, number> = {
  kpa: 1,
  psi: 6.894_76,
  bar: 100,
}

// Display precision by unit
export const temperaturePrecision: Record<TemperatureUnit, number> = {
  c: 1,
  f: 1,
}

export const gravityPrecision: Record<GravityUnit, number> = {
  sg: 3,
  plato: 1,
}

export const volumePrecision: Record<VolumeUnit, number> = {
  ml: 0,
  l: 2,
  hl: 2,
  usfloz: 2,
  ukfloz: 2,
  usgal: 2,
  ukgal: 2,
  bbl: 2,
  ukbbl: 2,
}

export const massPrecision: Record<MassUnit, number> = {
  g: 0,
  kg: 2,
  oz: 2,
  lb: 2,
}

export const pressurePrecision: Record<PressureUnit, number> = {
  kpa: 1,
  psi: 1,
  bar: 1,
}

export const colorPrecision: Record<ColorUnit, number> = {
  srm: 1,
  ebc: 1,
}

// Helper to get precision for any unit type
export function getPrecision (
  unit: TemperatureUnit | GravityUnit | VolumeUnit | MassUnit | PressureUnit | ColorUnit,
): number {
  if (unit in temperaturePrecision) {
    return temperaturePrecision[unit as TemperatureUnit]
  }
  if (unit in gravityPrecision) {
    return gravityPrecision[unit as GravityUnit]
  }
  if (unit in volumePrecision) {
    return volumePrecision[unit as VolumeUnit]
  }
  if (unit in massPrecision) {
    return massPrecision[unit as MassUnit]
  }
  if (unit in pressurePrecision) {
    return pressurePrecision[unit as PressureUnit]
  }
  if (unit in colorPrecision) {
    return colorPrecision[unit as ColorUnit]
  }
  return 2 // default
}

// Helper to get label for any unit type
export function getUnitLabel (
  unit: TemperatureUnit | GravityUnit | VolumeUnit | MassUnit | PressureUnit | ColorUnit,
): string {
  if (unit in temperatureLabels) {
    return temperatureLabels[unit as TemperatureUnit]
  }
  if (unit in gravityLabels) {
    return gravityLabels[unit as GravityUnit]
  }
  if (unit in volumeLabels) {
    return volumeLabels[unit as VolumeUnit]
  }
  if (unit in massLabels) {
    return massLabels[unit as MassUnit]
  }
  if (unit in pressureLabels) {
    return pressureLabels[unit as PressureUnit]
  }
  if (unit in colorLabels) {
    return colorLabels[unit as ColorUnit]
  }
  return unit
}

// Low-level conversion functions

/**
 * Convert temperature between Celsius and Fahrenheit.
 */
export function convertTemperature (
  value: number | null | undefined,
  from: TemperatureUnit,
  to: TemperatureUnit,
): number | null {
  if (value === null || value === undefined || !Number.isFinite(value)) {
    return null
  }
  if (from === to) {
    return value
  }

  if (from === 'c' && to === 'f') {
    // C → F: F = (C × 9/5) + 32
    return (value * 9) / 5 + 32
  }
  // F → C: C = (F - 32) × 5/9
  return ((value - 32) * 5) / 9
}

/**
 * Convert gravity between Specific Gravity and Degrees Plato.
 */
export function convertGravity (
  value: number | null | undefined,
  from: GravityUnit,
  to: GravityUnit,
): number | null {
  if (value === null || value === undefined || !Number.isFinite(value)) {
    return null
  }
  if (from === to) {
    return value
  }

  if (from === 'sg' && to === 'plato') {
    // SG → Plato: P = (-1 × 616.868) + (1111.14 × SG) - (630.272 × SG²) + (135.997 × SG³)
    const sg = value
    return -616.868 + 1111.14 * sg - 630.272 * sg * sg + 135.997 * sg * sg * sg
  }
  // Plato → SG: SG = 1 + (P / (258.6 - ((P / 258.2) × 227.1)))
  const p = value
  return 1 + p / (258.6 - (p / 258.2) * 227.1)
}

/**
 * Convert volume between supported units.
 */
export function convertVolume (
  value: number | null | undefined,
  from: VolumeUnit,
  to: VolumeUnit,
): number | null {
  if (value === null || value === undefined || !Number.isFinite(value)) {
    return null
  }
  if (from === to) {
    return value
  }

  // Convert to mL first, then to target unit
  const ml = value * volumeToMl[from]
  return ml / volumeToMl[to]
}

/**
 * Convert mass/weight between supported units.
 */
export function convertMass (
  value: number | null | undefined,
  from: MassUnit,
  to: MassUnit,
): number | null {
  if (value === null || value === undefined || !Number.isFinite(value)) {
    return null
  }
  if (from === to) {
    return value
  }

  // Convert to grams first, then to target unit
  const grams = value * massToGrams[from]
  return grams / massToGrams[to]
}

/**
 * Convert pressure between supported units.
 */
export function convertPressure (
  value: number | null | undefined,
  from: PressureUnit,
  to: PressureUnit,
): number | null {
  if (value === null || value === undefined || !Number.isFinite(value)) {
    return null
  }
  if (from === to) {
    return value
  }

  // Convert to kPa first, then to target unit
  const kpa = value * pressureToKpa[from]
  return kpa / pressureToKpa[to]
}

/**
 * Convert color between SRM and EBC.
 */
export function convertColor (
  value: number | null | undefined,
  from: ColorUnit,
  to: ColorUnit,
): number | null {
  if (value === null || value === undefined || !Number.isFinite(value)) {
    return null
  }
  if (from === to) {
    return value
  }

  if (from === 'srm' && to === 'ebc') {
    // SRM → EBC: EBC = SRM × 1.97
    return value * 1.97
  }
  // EBC → SRM: SRM = EBC / 1.97
  return value / 1.97
}

// Formatting functions

/**
 * Convert temperature and format with unit label.
 */
export function formatTemperature (
  value: number | null | undefined,
  fromUnit: TemperatureUnit,
  toUnit: TemperatureUnit,
): string {
  const converted = convertTemperature(value, fromUnit, toUnit)
  if (converted === null) {
    return '—'
  }
  const precision = temperaturePrecision[toUnit]
  return `${converted.toFixed(precision)}${temperatureLabels[toUnit]}`
}

/**
 * Convert gravity and format with unit label.
 */
export function formatGravity (
  value: number | null | undefined,
  fromUnit: GravityUnit,
  toUnit: GravityUnit,
): string {
  const converted = convertGravity(value, fromUnit, toUnit)
  if (converted === null) {
    return '—'
  }
  const precision = gravityPrecision[toUnit]
  return `${converted.toFixed(precision)} ${gravityLabels[toUnit]}`
}

/**
 * Convert volume and format with unit label.
 */
export function formatVolume (
  value: number | null | undefined,
  fromUnit: VolumeUnit,
  toUnit: VolumeUnit,
): string {
  const converted = convertVolume(value, fromUnit, toUnit)
  if (converted === null) {
    return '—'
  }
  const precision = volumePrecision[toUnit]
  return `${converted.toFixed(precision)} ${volumeLabels[toUnit]}`
}

/**
 * Convert mass and format with unit label.
 */
export function formatMass (
  value: number | null | undefined,
  fromUnit: MassUnit,
  toUnit: MassUnit,
): string {
  const converted = convertMass(value, fromUnit, toUnit)
  if (converted === null) {
    return '—'
  }
  const precision = massPrecision[toUnit]
  return `${converted.toFixed(precision)} ${massLabels[toUnit]}`
}

/**
 * Convert pressure and format with unit label.
 */
export function formatPressure (
  value: number | null | undefined,
  fromUnit: PressureUnit,
  toUnit: PressureUnit,
): string {
  const converted = convertPressure(value, fromUnit, toUnit)
  if (converted === null) {
    return '—'
  }
  const precision = pressurePrecision[toUnit]
  return `${converted.toFixed(precision)} ${pressureLabels[toUnit]}`
}

/**
 * Convert color and format with unit label.
 */
export function formatColor (
  value: number | null | undefined,
  fromUnit: ColorUnit,
  toUnit: ColorUnit,
): string {
  const converted = convertColor(value, fromUnit, toUnit)
  if (converted === null) {
    return '—'
  }
  const precision = colorPrecision[toUnit]
  return `${converted.toFixed(precision)} ${colorLabels[toUnit]}`
}

/**
 * Composable providing unit conversion utilities.
 */
export function useUnitConversion () {
  return {
    // Labels
    temperatureLabels,
    gravityLabels,
    volumeLabels,
    massLabels,
    pressureLabels,
    colorLabels,

    // Precision
    temperaturePrecision,
    gravityPrecision,
    volumePrecision,
    massPrecision,
    pressurePrecision,
    colorPrecision,
    getPrecision,
    getUnitLabel,

    // Conversion functions
    convertTemperature,
    convertGravity,
    convertVolume,
    convertMass,
    convertPressure,
    convertColor,

    // Formatting functions
    formatTemperature,
    formatGravity,
    formatVolume,
    formatMass,
    formatPressure,
    formatColor,
  }
}
