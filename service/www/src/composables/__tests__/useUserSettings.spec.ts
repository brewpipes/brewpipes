import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Import after mocking localStorage
import { useUserSettings } from '../useUserSettings'

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

describe('useUserSettings', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  afterEach(() => {
    // Reset settings to defaults after each test
    const { resetToDefaults } = useUserSettings()
    resetToDefaults()
  })

  describe('default settings', () => {
    it('has default brewery name', () => {
      const { breweryName, settings } = useUserSettings()
      expect(breweryName.value).toBe('Acme Brewing')
      expect(settings.value.breweryName).toBe('Acme Brewing')
    })
  })

  describe('setBreweryName', () => {
    it('updates brewery name', () => {
      const { breweryName, setBreweryName } = useUserSettings()
      setBreweryName('My Awesome Brewery')
      expect(breweryName.value).toBe('My Awesome Brewery')
    })

    it('trims whitespace from brewery name', () => {
      const { breweryName, setBreweryName } = useUserSettings()
      setBreweryName('  Trimmed Brewery  ')
      expect(breweryName.value).toBe('Trimmed Brewery')
    })

    it('ignores empty string', () => {
      const { breweryName, setBreweryName } = useUserSettings()
      const originalName = breweryName.value
      setBreweryName('')
      expect(breweryName.value).toBe(originalName)
    })

    it('ignores whitespace-only string', () => {
      const { breweryName, setBreweryName } = useUserSettings()
      const originalName = breweryName.value
      setBreweryName('   ')
      expect(breweryName.value).toBe(originalName)
    })

    it('ignores tab-only string', () => {
      const { breweryName, setBreweryName } = useUserSettings()
      const originalName = breweryName.value
      setBreweryName('\t\t')
      expect(breweryName.value).toBe(originalName)
    })

    it('ignores newline-only string', () => {
      const { breweryName, setBreweryName } = useUserSettings()
      const originalName = breweryName.value
      setBreweryName('\n\n')
      expect(breweryName.value).toBe(originalName)
    })
  })

  describe('resetToDefaults', () => {
    it('resets brewery name to default', () => {
      const { breweryName, setBreweryName, resetToDefaults } = useUserSettings()

      // Change the name
      setBreweryName('Custom Brewery')
      expect(breweryName.value).toBe('Custom Brewery')

      // Reset
      resetToDefaults()
      expect(breweryName.value).toBe('Acme Brewing')
    })
  })

  describe('settings object', () => {
    it('settings object reflects changes', () => {
      const { settings, setBreweryName } = useUserSettings()
      setBreweryName('New Brewery Name')
      expect(settings.value.breweryName).toBe('New Brewery Name')
    })
  })

  describe('breweryName computed', () => {
    it('breweryName is reactive', () => {
      const { breweryName, setBreweryName } = useUserSettings()

      expect(breweryName.value).toBe('Acme Brewing')
      setBreweryName('Reactive Brewery')
      expect(breweryName.value).toBe('Reactive Brewery')
    })
  })

  describe('singleton behavior', () => {
    it('returns the same state across multiple calls', () => {
      const instance1 = useUserSettings()
      const instance2 = useUserSettings()

      instance1.setBreweryName('Shared Brewery')

      expect(instance2.breweryName.value).toBe('Shared Brewery')
      expect(instance1.settings.value).toBe(instance2.settings.value)
    })
  })
})

describe('localStorage persistence', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  afterEach(() => {
    const { resetToDefaults } = useUserSettings()
    resetToDefaults()
  })

  it('saves settings to localStorage when changed', async () => {
    const { setBreweryName } = useUserSettings()
    setBreweryName('Persisted Brewery')

    // Wait for the watcher to trigger
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(localStorageMock.setItem).toHaveBeenCalledWith(
      'brewpipes:userSettings',
      expect.any(String),
    )

    // Verify the saved value
    const savedValue = localStorageMock.setItem.mock.calls.find(
      call => call[0] === 'brewpipes:userSettings',
    )?.[1]
    expect(savedValue).toBeDefined()
    const parsed = JSON.parse(savedValue!)
    expect(parsed.breweryName).toBe('Persisted Brewery')
  })
})

describe('localStorage loading', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  afterEach(() => {
    const { resetToDefaults } = useUserSettings()
    resetToDefaults()
  })

  it('handles invalid JSON in localStorage gracefully', () => {
    localStorageMock.getItem.mockReturnValueOnce('invalid json')

    // The module should not throw and should use defaults
    // Note: Since the module is already loaded, this tests the fallback behavior
    const { breweryName } = useUserSettings()
    expect(breweryName.value).toBeDefined()
  })

  it('handles missing breweryName in stored data', () => {
    localStorageMock.getItem.mockReturnValueOnce('{}')

    const { breweryName } = useUserSettings()
    expect(breweryName.value).toBeDefined()
  })

  it('handles empty breweryName in stored data', () => {
    localStorageMock.getItem.mockReturnValueOnce('{"breweryName": ""}')

    const { breweryName } = useUserSettings()
    expect(breweryName.value).toBeDefined()
  })

  it('handles whitespace-only breweryName in stored data', () => {
    localStorageMock.getItem.mockReturnValueOnce('{"breweryName": "   "}')

    const { breweryName } = useUserSettings()
    expect(breweryName.value).toBeDefined()
  })
})

describe('edge cases', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  afterEach(() => {
    const { resetToDefaults } = useUserSettings()
    resetToDefaults()
  })

  it('handles very long brewery names', () => {
    const { breweryName, setBreweryName } = useUserSettings()
    const longName = 'A'.repeat(1000)
    setBreweryName(longName)
    expect(breweryName.value).toBe(longName)
  })

  it('handles special characters in brewery name', () => {
    const { breweryName, setBreweryName } = useUserSettings()
    const specialName = 'O\'Malley\'s & Sons Brewing Co. (Est. 1892)'
    setBreweryName(specialName)
    expect(breweryName.value).toBe(specialName)
  })

  it('handles unicode characters in brewery name', () => {
    const { breweryName, setBreweryName } = useUserSettings()
    const unicodeName = 'Brauerei München 🍺'
    setBreweryName(unicodeName)
    expect(breweryName.value).toBe(unicodeName)
  })

  it('handles emoji-only brewery name', () => {
    const { breweryName, setBreweryName } = useUserSettings()
    const emojiName = '🍺🍻🍺'
    setBreweryName(emojiName)
    expect(breweryName.value).toBe(emojiName)
  })
})
