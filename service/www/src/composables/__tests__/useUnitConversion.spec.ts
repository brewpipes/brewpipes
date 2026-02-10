import { describe, expect, it } from 'vitest'
import {
  colorLabels,
  convertColor,
  convertGravity,
  convertMass,
  convertPressure,
  convertTemperature,
  convertVolume,
  formatColor,
  formatGravity,
  formatMass,
  formatPressure,
  formatTemperature,
  formatVolume,
  getPrecision,
  getUnitLabel,
  gravityLabels,
  massLabels,
  pressureLabels,
  temperatureLabels,
  volumeLabels,
} from '../useUnitConversion'

describe('convertTemperature', () => {
  describe('Celsius to Fahrenheit', () => {
    it('converts 0°C to 32°F', () => {
      expect(convertTemperature(0, 'c', 'f')).toBe(32)
    })

    it('converts 100°C to 212°F', () => {
      expect(convertTemperature(100, 'c', 'f')).toBe(212)
    })

    it('converts -40°C to -40°F', () => {
      expect(convertTemperature(-40, 'c', 'f')).toBe(-40)
    })

    it('converts 20°C to 68°F', () => {
      expect(convertTemperature(20, 'c', 'f')).toBe(68)
    })
  })

  describe('Fahrenheit to Celsius', () => {
    it('converts 32°F to 0°C', () => {
      expect(convertTemperature(32, 'f', 'c')).toBe(0)
    })

    it('converts 212°F to 100°C', () => {
      expect(convertTemperature(212, 'f', 'c')).toBe(100)
    })

    it('converts -40°F to -40°C', () => {
      expect(convertTemperature(-40, 'f', 'c')).toBe(-40)
    })

    it('converts 68°F to 20°C', () => {
      expect(convertTemperature(68, 'f', 'c')).toBe(20)
    })
  })

  describe('same unit conversion', () => {
    it('returns the same value when converting C to C', () => {
      expect(convertTemperature(25, 'c', 'c')).toBe(25)
    })

    it('returns the same value when converting F to F', () => {
      expect(convertTemperature(77, 'f', 'f')).toBe(77)
    })
  })

  describe('edge cases', () => {
    it('returns null for null input', () => {
      expect(convertTemperature(null, 'c', 'f')).toBeNull()
    })

    it('returns null for undefined input', () => {
      expect(convertTemperature(undefined, 'c', 'f')).toBeNull()
    })

    it('returns null for NaN input', () => {
      expect(convertTemperature(Number.NaN, 'c', 'f')).toBeNull()
    })

    it('returns null for Infinity input', () => {
      expect(convertTemperature(Infinity, 'c', 'f')).toBeNull()
    })

    it('handles very large numbers', () => {
      const result = convertTemperature(1_000_000, 'c', 'f')
      expect(result).toBe(1_800_032)
    })
  })
})

describe('convertGravity', () => {
  describe('SG to Plato', () => {
    it('converts 1.000 SG to approximately 0°P', () => {
      const result = convertGravity(1, 'sg', 'plato')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(0, 0)
    })

    it('converts 1.040 SG to approximately 10°P', () => {
      const result = convertGravity(1.04, 'sg', 'plato')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(10, 0)
    })

    it('converts 1.080 SG to approximately 19.3°P', () => {
      const result = convertGravity(1.08, 'sg', 'plato')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(19.3, 0)
    })
  })

  describe('Plato to SG', () => {
    it('converts 0°P to approximately 1.000 SG', () => {
      const result = convertGravity(0, 'plato', 'sg')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(1, 2)
    })

    it('converts 10°P to approximately 1.040 SG', () => {
      const result = convertGravity(10, 'plato', 'sg')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(1.04, 2)
    })

    it('converts 20°P to approximately 1.083 SG', () => {
      const result = convertGravity(20, 'plato', 'sg')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(1.083, 2)
    })
  })

  describe('same unit conversion', () => {
    it('returns the same value when converting SG to SG', () => {
      expect(convertGravity(1.05, 'sg', 'sg')).toBe(1.05)
    })

    it('returns the same value when converting Plato to Plato', () => {
      expect(convertGravity(12.5, 'plato', 'plato')).toBe(12.5)
    })
  })

  describe('edge cases', () => {
    it('returns null for null input', () => {
      expect(convertGravity(null, 'sg', 'plato')).toBeNull()
    })

    it('returns null for undefined input', () => {
      expect(convertGravity(undefined, 'sg', 'plato')).toBeNull()
    })

    it('returns null for NaN input', () => {
      expect(convertGravity(Number.NaN, 'sg', 'plato')).toBeNull()
    })
  })
})

describe('convertVolume', () => {
  describe('liters to other units', () => {
    it('converts 1 L to 1000 mL', () => {
      expect(convertVolume(1, 'l', 'ml')).toBe(1000)
    })

    it('converts 100 L to 1 hL', () => {
      expect(convertVolume(100, 'l', 'hl')).toBe(1)
    })

    it('converts 1 L to approximately 0.264 US gallons', () => {
      const result = convertVolume(1, 'l', 'usgal')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(0.264, 2)
    })
  })

  describe('gallons to liters', () => {
    it('converts 1 US gallon to approximately 3.785 L', () => {
      const result = convertVolume(1, 'usgal', 'l')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(3.785, 2)
    })

    it('converts 1 UK gallon to approximately 4.546 L', () => {
      const result = convertVolume(1, 'ukgal', 'l')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(4.546, 2)
    })
  })

  describe('barrels', () => {
    it('converts 1 US barrel to approximately 31 US gallons', () => {
      const result = convertVolume(1, 'bbl', 'usgal')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(31, 0)
    })

    it('converts 1 UK barrel to approximately 36 UK gallons', () => {
      const result = convertVolume(1, 'ukbbl', 'ukgal')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(36, 0)
    })
  })

  describe('fluid ounces', () => {
    it('converts 1 US fl oz to approximately 29.57 mL', () => {
      const result = convertVolume(1, 'usfloz', 'ml')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(29.57, 1)
    })

    it('converts 1 UK fl oz to approximately 28.41 mL', () => {
      const result = convertVolume(1, 'ukfloz', 'ml')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(28.41, 1)
    })
  })

  describe('same unit conversion', () => {
    it('returns the same value when converting L to L', () => {
      expect(convertVolume(5, 'l', 'l')).toBe(5)
    })
  })

  describe('edge cases', () => {
    it('returns null for null input', () => {
      expect(convertVolume(null, 'l', 'ml')).toBeNull()
    })

    it('returns null for undefined input', () => {
      expect(convertVolume(undefined, 'l', 'ml')).toBeNull()
    })

    it('handles zero correctly', () => {
      expect(convertVolume(0, 'l', 'ml')).toBe(0)
    })

    it('handles negative values', () => {
      expect(convertVolume(-1, 'l', 'ml')).toBe(-1000)
    })
  })
})

describe('convertMass', () => {
  describe('kilograms to other units', () => {
    it('converts 1 kg to 1000 g', () => {
      expect(convertMass(1, 'kg', 'g')).toBe(1000)
    })

    it('converts 1 kg to approximately 2.205 lb', () => {
      const result = convertMass(1, 'kg', 'lb')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(2.205, 2)
    })

    it('converts 1 kg to approximately 35.27 oz', () => {
      const result = convertMass(1, 'kg', 'oz')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(35.27, 1)
    })
  })

  describe('pounds to other units', () => {
    it('converts 1 lb to approximately 453.6 g', () => {
      const result = convertMass(1, 'lb', 'g')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(453.6, 0)
    })

    it('converts 1 lb to 16 oz', () => {
      const result = convertMass(1, 'lb', 'oz')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(16, 0)
    })
  })

  describe('same unit conversion', () => {
    it('returns the same value when converting kg to kg', () => {
      expect(convertMass(5, 'kg', 'kg')).toBe(5)
    })
  })

  describe('edge cases', () => {
    it('returns null for null input', () => {
      expect(convertMass(null, 'kg', 'g')).toBeNull()
    })

    it('returns null for undefined input', () => {
      expect(convertMass(undefined, 'kg', 'g')).toBeNull()
    })

    it('handles zero correctly', () => {
      expect(convertMass(0, 'kg', 'g')).toBe(0)
    })
  })
})

describe('convertPressure', () => {
  describe('kPa to other units', () => {
    it('converts 100 kPa to 1 bar', () => {
      expect(convertPressure(100, 'kpa', 'bar')).toBe(1)
    })

    it('converts 1 kPa to approximately 0.145 psi', () => {
      const result = convertPressure(1, 'kpa', 'psi')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(0.145, 2)
    })
  })

  describe('psi to other units', () => {
    it('converts 1 psi to approximately 6.895 kPa', () => {
      const result = convertPressure(1, 'psi', 'kpa')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(6.895, 2)
    })

    it('converts 14.5 psi to approximately 1 bar', () => {
      const result = convertPressure(14.5, 'psi', 'bar')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(1, 0)
    })
  })

  describe('bar to other units', () => {
    it('converts 1 bar to 100 kPa', () => {
      expect(convertPressure(1, 'bar', 'kpa')).toBe(100)
    })

    it('converts 1 bar to approximately 14.5 psi', () => {
      const result = convertPressure(1, 'bar', 'psi')
      expect(result).not.toBeNull()
      expect(result!).toBeCloseTo(14.5, 0)
    })
  })

  describe('same unit conversion', () => {
    it('returns the same value when converting kPa to kPa', () => {
      expect(convertPressure(101.3, 'kpa', 'kpa')).toBe(101.3)
    })
  })

  describe('edge cases', () => {
    it('returns null for null input', () => {
      expect(convertPressure(null, 'kpa', 'psi')).toBeNull()
    })

    it('returns null for undefined input', () => {
      expect(convertPressure(undefined, 'kpa', 'psi')).toBeNull()
    })

    it('handles zero correctly', () => {
      expect(convertPressure(0, 'kpa', 'psi')).toBe(0)
    })
  })
})

describe('convertColor', () => {
  describe('SRM to EBC', () => {
    it('converts 1 SRM to 1.97 EBC', () => {
      expect(convertColor(1, 'srm', 'ebc')).toBe(1.97)
    })

    it('converts 10 SRM to 19.7 EBC', () => {
      expect(convertColor(10, 'srm', 'ebc')).toBe(19.7)
    })

    it('converts 40 SRM to 78.8 EBC', () => {
      expect(convertColor(40, 'srm', 'ebc')).toBe(78.8)
    })
  })

  describe('EBC to SRM', () => {
    it('converts 1.97 EBC to 1 SRM', () => {
      expect(convertColor(1.97, 'ebc', 'srm')).toBe(1)
    })

    it('converts 19.7 EBC to 10 SRM', () => {
      expect(convertColor(19.7, 'ebc', 'srm')).toBe(10)
    })

    it('converts 78.8 EBC to 40 SRM', () => {
      expect(convertColor(78.8, 'ebc', 'srm')).toBe(40)
    })
  })

  describe('same unit conversion', () => {
    it('returns the same value when converting SRM to SRM', () => {
      expect(convertColor(15, 'srm', 'srm')).toBe(15)
    })

    it('returns the same value when converting EBC to EBC', () => {
      expect(convertColor(30, 'ebc', 'ebc')).toBe(30)
    })
  })

  describe('edge cases', () => {
    it('returns null for null input', () => {
      expect(convertColor(null, 'srm', 'ebc')).toBeNull()
    })

    it('returns null for undefined input', () => {
      expect(convertColor(undefined, 'srm', 'ebc')).toBeNull()
    })

    it('handles zero correctly', () => {
      expect(convertColor(0, 'srm', 'ebc')).toBe(0)
    })
  })
})

describe('formatTemperature', () => {
  it('formats temperature with correct precision and label', () => {
    expect(formatTemperature(20, 'c', 'c')).toBe('20.0°C')
    expect(formatTemperature(68, 'f', 'f')).toBe('68.0°F')
  })

  it('converts and formats temperature', () => {
    expect(formatTemperature(0, 'c', 'f')).toBe('32.0°F')
    expect(formatTemperature(32, 'f', 'c')).toBe('0.0°C')
  })

  it('returns em dash for null input', () => {
    expect(formatTemperature(null, 'c', 'f')).toBe('—')
  })

  it('returns em dash for undefined input', () => {
    expect(formatTemperature(undefined, 'c', 'f')).toBe('—')
  })
})

describe('formatGravity', () => {
  it('formats gravity with correct precision and label', () => {
    expect(formatGravity(1.05, 'sg', 'sg')).toBe('1.050 SG')
  })

  it('returns em dash for null input', () => {
    expect(formatGravity(null, 'sg', 'plato')).toBe('—')
  })
})

describe('formatVolume', () => {
  it('formats volume with correct precision and label', () => {
    expect(formatVolume(1000, 'ml', 'ml')).toBe('1000 mL')
    expect(formatVolume(1, 'l', 'l')).toBe('1.00 L')
    expect(formatVolume(10, 'bbl', 'bbl')).toBe('10.00 bbl')
  })

  it('converts and formats volume', () => {
    expect(formatVolume(1, 'l', 'ml')).toBe('1000 mL')
  })

  it('returns em dash for null input', () => {
    expect(formatVolume(null, 'l', 'ml')).toBe('—')
  })
})

describe('formatMass', () => {
  it('formats mass with correct precision and label', () => {
    expect(formatMass(500, 'g', 'g')).toBe('500 g')
    expect(formatMass(1, 'kg', 'kg')).toBe('1.00 kg')
    expect(formatMass(10, 'lb', 'lb')).toBe('10.00 lb')
  })

  it('returns em dash for null input', () => {
    expect(formatMass(null, 'kg', 'g')).toBe('—')
  })
})

describe('formatPressure', () => {
  it('formats pressure with correct precision and label', () => {
    expect(formatPressure(100, 'kpa', 'kpa')).toBe('100.0 kPa')
    expect(formatPressure(14.7, 'psi', 'psi')).toBe('14.7 psi')
    expect(formatPressure(1, 'bar', 'bar')).toBe('1.0 bar')
  })

  it('returns em dash for null input', () => {
    expect(formatPressure(null, 'kpa', 'psi')).toBe('—')
  })
})

describe('formatColor', () => {
  it('formats color with correct precision and label', () => {
    expect(formatColor(10, 'srm', 'srm')).toBe('10.0 SRM')
    expect(formatColor(20, 'ebc', 'ebc')).toBe('20.0 EBC')
  })

  it('converts and formats color', () => {
    expect(formatColor(10, 'srm', 'ebc')).toBe('19.7 EBC')
  })

  it('returns em dash for null input', () => {
    expect(formatColor(null, 'srm', 'ebc')).toBe('—')
  })
})

describe('getPrecision', () => {
  it('returns correct precision for temperature units', () => {
    expect(getPrecision('c')).toBe(1)
    expect(getPrecision('f')).toBe(1)
  })

  it('returns correct precision for gravity units', () => {
    expect(getPrecision('sg')).toBe(3)
    expect(getPrecision('plato')).toBe(1)
  })

  it('returns correct precision for volume units', () => {
    expect(getPrecision('ml')).toBe(0)
    expect(getPrecision('l')).toBe(2)
    expect(getPrecision('bbl')).toBe(2)
  })

  it('returns correct precision for mass units', () => {
    expect(getPrecision('g')).toBe(0)
    expect(getPrecision('kg')).toBe(2)
    expect(getPrecision('lb')).toBe(2)
  })

  it('returns correct precision for pressure units', () => {
    expect(getPrecision('kpa')).toBe(1)
    expect(getPrecision('psi')).toBe(1)
    expect(getPrecision('bar')).toBe(1)
  })

  it('returns correct precision for color units', () => {
    expect(getPrecision('srm')).toBe(1)
    expect(getPrecision('ebc')).toBe(1)
  })
})

describe('getUnitLabel', () => {
  it('returns correct labels for temperature units', () => {
    expect(getUnitLabel('c')).toBe('°C')
    expect(getUnitLabel('f')).toBe('°F')
  })

  it('returns correct labels for gravity units', () => {
    expect(getUnitLabel('sg')).toBe('SG')
    expect(getUnitLabel('plato')).toBe('°P')
  })

  it('returns correct labels for volume units', () => {
    expect(getUnitLabel('ml')).toBe('mL')
    expect(getUnitLabel('l')).toBe('L')
    expect(getUnitLabel('bbl')).toBe('bbl')
    expect(getUnitLabel('usgal')).toBe('gal')
  })

  it('returns correct labels for mass units', () => {
    expect(getUnitLabel('g')).toBe('g')
    expect(getUnitLabel('kg')).toBe('kg')
    expect(getUnitLabel('lb')).toBe('lb')
  })

  it('returns correct labels for pressure units', () => {
    expect(getUnitLabel('kpa')).toBe('kPa')
    expect(getUnitLabel('psi')).toBe('psi')
    expect(getUnitLabel('bar')).toBe('bar')
  })

  it('returns correct labels for color units', () => {
    expect(getUnitLabel('srm')).toBe('SRM')
    expect(getUnitLabel('ebc')).toBe('EBC')
  })
})

describe('unit label constants', () => {
  it('temperatureLabels has correct values', () => {
    expect(temperatureLabels.c).toBe('°C')
    expect(temperatureLabels.f).toBe('°F')
  })

  it('gravityLabels has correct values', () => {
    expect(gravityLabels.sg).toBe('SG')
    expect(gravityLabels.plato).toBe('°P')
  })

  it('volumeLabels has correct values', () => {
    expect(volumeLabels.ml).toBe('mL')
    expect(volumeLabels.l).toBe('L')
    expect(volumeLabels.hl).toBe('hL')
    expect(volumeLabels.bbl).toBe('bbl')
  })

  it('massLabels has correct values', () => {
    expect(massLabels.g).toBe('g')
    expect(massLabels.kg).toBe('kg')
    expect(massLabels.oz).toBe('oz')
    expect(massLabels.lb).toBe('lb')
  })

  it('pressureLabels has correct values', () => {
    expect(pressureLabels.kpa).toBe('kPa')
    expect(pressureLabels.psi).toBe('psi')
    expect(pressureLabels.bar).toBe('bar')
  })

  it('colorLabels has correct values', () => {
    expect(colorLabels.srm).toBe('SRM')
    expect(colorLabels.ebc).toBe('EBC')
  })
})
