import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Import after mocking localStorage
import {
  colorOptions,
  gravityOptions,
  isMassUnit,
  isVolumeUnit,
  massOptions,
  normalizeMassUnit,
  normalizeVolumeUnit,
  pressureOptions,
  temperatureOptions,
  useUnitPreferences,
  volumeOptions,
} from '../useUnitPreferences'

// Mock localStorage before importing the module
const localStorageMock = (() => {
  let store: Record<string, string> = {}
  return {
    getItem: vi.fn((key: string) => store[key] ?? null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key]
    }),
    clear: vi.fn(() => {
      store = {}
    }),
    get length () {
      return Object.keys(store).length
    },
    key: vi.fn((index: number) => Object.keys(store)[index] ?? null),
  }
})()

Object.defineProperty(globalThis, 'localStorage', {
  value: localStorageMock,
  writable: true,
})

describe('useUnitPreferences', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  afterEach(() => {
    // Reset preferences to defaults after each test
    const { resetToDefaults } = useUnitPreferences()
    resetToDefaults()
  })

  describe('default preferences', () => {
    it('has US-centric defaults', () => {
      const { preferences } = useUnitPreferences()
      expect(preferences.value.temperature).toBe('f')
      expect(preferences.value.gravity).toBe('sg')
      expect(preferences.value.volume).toBe('bbl')
      expect(preferences.value.mass).toBe('lb')
      expect(preferences.value.pressure).toBe('psi')
      expect(preferences.value.color).toBe('srm')
    })
  })

  describe('setters', () => {
    it('setTemperatureUnit updates temperature preference', () => {
      const { preferences, setTemperatureUnit } = useUnitPreferences()
      setTemperatureUnit('c')
      expect(preferences.value.temperature).toBe('c')
    })

    it('setGravityUnit updates gravity preference', () => {
      const { preferences, setGravityUnit } = useUnitPreferences()
      setGravityUnit('plato')
      expect(preferences.value.gravity).toBe('plato')
    })

    it('setVolumeUnit updates volume preference', () => {
      const { preferences, setVolumeUnit } = useUnitPreferences()
      setVolumeUnit('l')
      expect(preferences.value.volume).toBe('l')
    })

    it('setMassUnit updates mass preference', () => {
      const { preferences, setMassUnit } = useUnitPreferences()
      setMassUnit('kg')
      expect(preferences.value.mass).toBe('kg')
    })

    it('setPressureUnit updates pressure preference', () => {
      const { preferences, setPressureUnit } = useUnitPreferences()
      setPressureUnit('bar')
      expect(preferences.value.pressure).toBe('bar')
    })

    it('setColorUnit updates color preference', () => {
      const { preferences, setColorUnit } = useUnitPreferences()
      setColorUnit('ebc')
      expect(preferences.value.color).toBe('ebc')
    })
  })

  describe('resetToDefaults', () => {
    it('resets all preferences to defaults', () => {
      const {
        preferences,
        setTemperatureUnit,
        setGravityUnit,
        setVolumeUnit,
        setMassUnit,
        setPressureUnit,
        setColorUnit,
        resetToDefaults,
      } = useUnitPreferences()

      // Change all preferences
      setTemperatureUnit('c')
      setGravityUnit('plato')
      setVolumeUnit('l')
      setMassUnit('kg')
      setPressureUnit('bar')
      setColorUnit('ebc')

      // Reset
      resetToDefaults()

      // Verify defaults
      expect(preferences.value.temperature).toBe('f')
      expect(preferences.value.gravity).toBe('sg')
      expect(preferences.value.volume).toBe('bbl')
      expect(preferences.value.mass).toBe('lb')
      expect(preferences.value.pressure).toBe('psi')
      expect(preferences.value.color).toBe('srm')
    })
  })

  describe('formatTemperaturePreferred', () => {
    it('formats temperature to preferred unit', () => {
      const { formatTemperaturePreferred, setTemperatureUnit } = useUnitPreferences()

      // Default is Fahrenheit
      expect(formatTemperaturePreferred(0, 'c')).toBe('32.0°F')

      // Change to Celsius
      setTemperatureUnit('c')
      expect(formatTemperaturePreferred(32, 'f')).toBe('0.0°C')
    })

    it('returns em dash for null value', () => {
      const { formatTemperaturePreferred } = useUnitPreferences()
      expect(formatTemperaturePreferred(null, 'c')).toBe('—')
    })
  })

  describe('formatGravityPreferred', () => {
    it('formats gravity to preferred unit', () => {
      const { formatGravityPreferred, setGravityUnit } = useUnitPreferences()

      // Default is SG
      expect(formatGravityPreferred(1.05, 'sg')).toBe('1.050 SG')

      // Change to Plato
      setGravityUnit('plato')
      const result = formatGravityPreferred(1.05, 'sg')
      expect(result).toMatch(/°P$/)
    })
  })

  describe('formatVolumePreferred', () => {
    it('formats volume to preferred unit', () => {
      const { formatVolumePreferred, setVolumeUnit } = useUnitPreferences()

      // Default is barrels
      const result = formatVolumePreferred(117_347.77, 'ml')
      expect(result).toMatch(/bbl$/)

      // Change to liters
      setVolumeUnit('l')
      expect(formatVolumePreferred(1000, 'ml')).toBe('1.00 L')
    })
  })

  describe('formatMassPreferred', () => {
    it('formats mass to preferred unit', () => {
      const { formatMassPreferred, setMassUnit } = useUnitPreferences()

      // Default is pounds
      const result = formatMassPreferred(453.592, 'g')
      expect(result).toMatch(/lb$/)

      // Change to kilograms
      setMassUnit('kg')
      expect(formatMassPreferred(1000, 'g')).toBe('1.00 kg')
    })
  })

  describe('formatPressurePreferred', () => {
    it('formats pressure to preferred unit', () => {
      const { formatPressurePreferred, setPressureUnit } = useUnitPreferences()

      // Default is psi
      const result = formatPressurePreferred(100, 'kpa')
      expect(result).toMatch(/psi$/)

      // Change to bar
      setPressureUnit('bar')
      expect(formatPressurePreferred(100, 'kpa')).toBe('1.0 bar')
    })
  })

  describe('formatColorPreferred', () => {
    it('formats color to preferred unit', () => {
      const { formatColorPreferred, setColorUnit } = useUnitPreferences()

      // Default is SRM
      expect(formatColorPreferred(10, 'srm')).toBe('10.0 SRM')

      // Change to EBC
      setColorUnit('ebc')
      expect(formatColorPreferred(10, 'srm')).toBe('19.7 EBC')
    })
  })

  describe('formatAmountPreferred', () => {
    it('formats volume units to preferred volume', () => {
      const { formatAmountPreferred, setVolumeUnit } = useUnitPreferences()
      setVolumeUnit('l')
      expect(formatAmountPreferred(1000, 'ml')).toBe('1.00 L')
    })

    it('formats mass units to preferred mass', () => {
      const { formatAmountPreferred, setMassUnit } = useUnitPreferences()
      setMassUnit('kg')
      expect(formatAmountPreferred(1000, 'g')).toBe('1.00 kg')
    })

    it('normalizes gal to usgal', () => {
      const { formatAmountPreferred, setVolumeUnit } = useUnitPreferences()
      setVolumeUnit('l')
      const result = formatAmountPreferred(1, 'gal')
      expect(result).toMatch(/L$/)
    })

    it('returns em dash for null value', () => {
      const { formatAmountPreferred } = useUnitPreferences()
      expect(formatAmountPreferred(null, 'kg')).toBe('—')
    })

    it('returns value only for null unit', () => {
      const { formatAmountPreferred } = useUnitPreferences()
      expect(formatAmountPreferred(100, null)).toBe('100')
    })

    it('returns value with unit for unknown units', () => {
      const { formatAmountPreferred } = useUnitPreferences()
      expect(formatAmountPreferred(100, 'unknown')).toBe('100 unknown')
    })
  })
})

describe('isVolumeUnit', () => {
  it('returns true for volume units', () => {
    expect(isVolumeUnit('ml')).toBe(true)
    expect(isVolumeUnit('l')).toBe(true)
    expect(isVolumeUnit('hl')).toBe(true)
    expect(isVolumeUnit('usfloz')).toBe(true)
    expect(isVolumeUnit('ukfloz')).toBe(true)
    expect(isVolumeUnit('usgal')).toBe(true)
    expect(isVolumeUnit('ukgal')).toBe(true)
    expect(isVolumeUnit('bbl')).toBe(true)
    expect(isVolumeUnit('ukbbl')).toBe(true)
    expect(isVolumeUnit('gal')).toBe(true)
  })

  it('returns false for non-volume units', () => {
    expect(isVolumeUnit('kg')).toBe(false)
    expect(isVolumeUnit('lb')).toBe(false)
    expect(isVolumeUnit('psi')).toBe(false)
    expect(isVolumeUnit('unknown')).toBe(false)
  })

  it('is case-insensitive', () => {
    expect(isVolumeUnit('ML')).toBe(true)
    expect(isVolumeUnit('L')).toBe(true)
    expect(isVolumeUnit('BBL')).toBe(true)
  })
})

describe('isMassUnit', () => {
  it('returns true for mass units', () => {
    expect(isMassUnit('g')).toBe(true)
    expect(isMassUnit('kg')).toBe(true)
    expect(isMassUnit('oz')).toBe(true)
    expect(isMassUnit('lb')).toBe(true)
  })

  it('returns false for non-mass units', () => {
    expect(isMassUnit('ml')).toBe(false)
    expect(isMassUnit('l')).toBe(false)
    expect(isMassUnit('psi')).toBe(false)
    expect(isMassUnit('unknown')).toBe(false)
  })

  it('is case-insensitive', () => {
    expect(isMassUnit('G')).toBe(true)
    expect(isMassUnit('KG')).toBe(true)
    expect(isMassUnit('LB')).toBe(true)
  })
})

describe('normalizeVolumeUnit', () => {
  it('normalizes gal to usgal', () => {
    expect(normalizeVolumeUnit('gal')).toBe('usgal')
    expect(normalizeVolumeUnit('GAL')).toBe('usgal')
  })

  it('returns lowercase for other units', () => {
    expect(normalizeVolumeUnit('ML')).toBe('ml')
    expect(normalizeVolumeUnit('L')).toBe('l')
    expect(normalizeVolumeUnit('BBL')).toBe('bbl')
  })
})

describe('normalizeMassUnit', () => {
  it('returns lowercase unit', () => {
    expect(normalizeMassUnit('KG')).toBe('kg')
    expect(normalizeMassUnit('LB')).toBe('lb')
    expect(normalizeMassUnit('G')).toBe('g')
    expect(normalizeMassUnit('OZ')).toBe('oz')
  })
})

describe('option arrays', () => {
  it('temperatureOptions has correct structure', () => {
    expect(temperatureOptions).toHaveLength(2)
    expect(temperatureOptions).toContainEqual({ value: 'c', label: 'Celsius (°C)' })
    expect(temperatureOptions).toContainEqual({ value: 'f', label: 'Fahrenheit (°F)' })
  })

  it('gravityOptions has correct structure', () => {
    expect(gravityOptions).toHaveLength(2)
    expect(gravityOptions).toContainEqual({ value: 'sg', label: 'Specific Gravity (SG)' })
    expect(gravityOptions).toContainEqual({ value: 'plato', label: 'Degrees Plato (°P)' })
  })

  it('volumeOptions has correct structure', () => {
    expect(volumeOptions).toHaveLength(9)
    expect(volumeOptions.map(o => o.value)).toContain('ml')
    expect(volumeOptions.map(o => o.value)).toContain('l')
    expect(volumeOptions.map(o => o.value)).toContain('bbl')
  })

  it('massOptions has correct structure', () => {
    expect(massOptions).toHaveLength(4)
    expect(massOptions.map(o => o.value)).toContain('g')
    expect(massOptions.map(o => o.value)).toContain('kg')
    expect(massOptions.map(o => o.value)).toContain('lb')
  })

  it('pressureOptions has correct structure', () => {
    expect(pressureOptions).toHaveLength(3)
    expect(pressureOptions.map(o => o.value)).toContain('kpa')
    expect(pressureOptions.map(o => o.value)).toContain('psi')
    expect(pressureOptions.map(o => o.value)).toContain('bar')
  })

  it('colorOptions has correct structure', () => {
    expect(colorOptions).toHaveLength(2)
    expect(colorOptions).toContainEqual({ value: 'srm', label: 'SRM' })
    expect(colorOptions).toContainEqual({ value: 'ebc', label: 'EBC' })
  })
})

describe('localStorage persistence', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  it('saves preferences to localStorage when changed', async () => {
    const { setTemperatureUnit } = useUnitPreferences()
    setTemperatureUnit('c')

    // Wait for the watcher to trigger
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      'brewpipes:unitPreferences',
      expect.any(String),
    )
  })
})
