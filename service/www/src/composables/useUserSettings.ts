import type { UserSettings } from '@/types'
import { computed, ref, watch } from 'vue'

// Re-export UserSettings for backward compatibility
export type { UserSettings } from '@/types'

const STORAGE_KEY = 'brewpipes:userSettings'

const DEFAULT_SETTINGS: UserSettings = {
  breweryName: 'Acme Brewing',
}

// Module-level singleton state
const settings = ref<UserSettings>(loadSettings())

/**
 * Load settings from localStorage, falling back to defaults on error.
 */
function loadSettings (): UserSettings {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      const parsed = JSON.parse(stored) as Partial<UserSettings>
      return {
        breweryName: typeof parsed.breweryName === 'string' && parsed.breweryName.trim()
          ? parsed.breweryName
          : DEFAULT_SETTINGS.breweryName,
      }
    }
  } catch {
    // localStorage unavailable or JSON parse error - use defaults
  }
  return { ...DEFAULT_SETTINGS }
}

/**
 * Save settings to localStorage.
 */
function saveSettings (userSettings: UserSettings): void {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(userSettings))
  } catch {
    // localStorage unavailable - silently ignore
  }
}

// Watch for changes and persist to localStorage
watch(
  settings,
  newSettings => {
    saveSettings(newSettings)
  },
  { deep: true },
)

/**
 * Composable for managing user settings with localStorage persistence.
 * Returns singleton state shared across all component instances.
 */
export function useUserSettings () {
  /**
   * Set the brewery name. Empty or whitespace-only values are ignored.
   */
  function setBreweryName (name: string): void {
    const trimmed = name.trim()
    if (!trimmed) {
      return
    }
    settings.value = { ...settings.value, breweryName: trimmed }
  }

  /**
   * Reset all settings to defaults.
   */
  function resetToDefaults (): void {
    settings.value = { ...DEFAULT_SETTINGS }
  }

  // Computed ref for reactive brewery name access
  const breweryName = computed(() => settings.value.breweryName)

  return {
    // Reactive settings (singleton)
    settings,

    // Convenience accessor for brewery name (computed for reactivity)
    breweryName,

    // Setters
    setBreweryName,
    resetToDefaults,
  }
}
