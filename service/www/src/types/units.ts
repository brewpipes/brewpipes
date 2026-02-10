/**
 * Unit types for measurements and conversions.
 *
 * These types define the supported units for various measurement categories
 * used throughout the brewing process.
 */

/** Temperature units: Celsius and Fahrenheit */
export type TemperatureUnit = 'c' | 'f'

/** Gravity units: Specific Gravity and Degrees Plato */
export type GravityUnit = 'sg' | 'plato'

/**
 * Volume units for liquid measurements.
 * Includes metric, US customary, and UK imperial units.
 */
export type VolumeUnit
  = | 'ml' // milliliters
    | 'l' // liters
    | 'hl' // hectoliters
    | 'usfloz' // US fluid ounces
    | 'ukfloz' // UK fluid ounces
    | 'usgal' // US gallons
    | 'ukgal' // UK gallons
    | 'bbl' // US barrels (31 gallons)
    | 'ukbbl' // UK barrels (36 UK gallons)

/** Mass/weight units */
export type MassUnit = 'g' | 'kg' | 'oz' | 'lb'

/** Pressure units */
export type PressureUnit = 'kpa' | 'psi' | 'bar'

/** Beer color units: SRM and EBC */
export type ColorUnit = 'srm' | 'ebc'

/**
 * User preferences for display units.
 * Determines how values are converted and displayed throughout the UI.
 */
export interface UnitPreferences {
  temperature: TemperatureUnit
  gravity: GravityUnit
  volume: VolumeUnit
  mass: MassUnit
  pressure: PressureUnit
  color: ColorUnit
}
