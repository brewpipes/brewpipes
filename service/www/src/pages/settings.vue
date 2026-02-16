<template>
  <v-container>
    <h1 class="text-h4 mb-6">Settings</h1>

    <v-card class="mb-6">
      <v-card-title>Brewery Settings</v-card-title>
      <v-card-subtitle>
        Configure your brewery identity and branding.
      </v-card-subtitle>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="breweryNameInput"
              density="comfortable"
              hide-details
              label="Brewery Name"
              variant="outlined"
              @blur="handleBreweryNameBlur"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-btn variant="text" @click="resetBrewerySettings">Reset to Defaults</v-btn>
      </v-card-actions>
    </v-card>

    <v-card class="mb-6">
      <v-card-title>Display Units</v-card-title>
      <v-card-subtitle>
        Choose your preferred units for displaying measurements throughout the app.
        Values are converted automatically.
      </v-card-subtitle>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="4" sm="6">
            <v-select
              v-model="preferences.temperature"
              density="comfortable"
              hide-details
              item-title="label"
              item-value="value"
              :items="temperatureOptions"
              label="Temperature"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="4" sm="6">
            <v-select
              v-model="preferences.gravity"
              density="comfortable"
              hide-details
              item-title="label"
              item-value="value"
              :items="gravityOptions"
              label="Gravity"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="4" sm="6">
            <v-select
              v-model="preferences.volume"
              density="comfortable"
              hide-details
              item-title="label"
              item-value="value"
              :items="volumeOptions"
              label="Volume"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="4" sm="6">
            <v-select
              v-model="preferences.mass"
              density="comfortable"
              hide-details
              item-title="label"
              item-value="value"
              :items="massOptions"
              label="Mass"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="4" sm="6">
            <v-select
              v-model="preferences.pressure"
              density="comfortable"
              hide-details
              item-title="label"
              item-value="value"
              :items="pressureOptions"
              label="Pressure"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="4" sm="6">
            <v-select
              v-model="preferences.color"
              density="comfortable"
              hide-details
              item-title="label"
              item-value="value"
              :items="colorOptions"
              label="Color"
              variant="outlined"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-btn variant="text" @click="resetToDefaults">Reset to Defaults</v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script lang="ts" setup>
  import { nextTick, onMounted, ref, watch } from 'vue'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { useUserSettings } from '@/composables/useUserSettings'

  const {
    preferences,
    resetToDefaults,
    temperatureOptions,
    gravityOptions,
    volumeOptions,
    massOptions,
    pressureOptions,
    colorOptions,
  } = useUnitPreferences()

  const {
    breweryName,
    setBreweryName,
    resetToDefaults: resetBrewerySettings,
  } = useUserSettings()

  const { showNotice } = useSnackbar()

  // Track whether initial load is complete to avoid showing notification on mount
  const initialized = ref(false)

  onMounted(() => {
    nextTick(() => {
      initialized.value = true
    })
  })

  // Show save confirmation when unit preferences change (user-initiated only)
  watch(
    preferences,
    () => {
      if (initialized.value) {
        showNotice('Settings saved')
      }
    },
    { deep: true },
  )

  // Local ref for text field editing; syncs on blur to allow validation
  const breweryNameInput = ref(breweryName.value)

  // Keep local input in sync if external changes occur (e.g., reset)
  watch(breweryName, value => {
    breweryNameInput.value = value
  })

  function handleBreweryNameBlur () {
    const trimmed = breweryNameInput.value.trim()
    if (trimmed) {
      setBreweryName(trimmed)
      if (initialized.value) {
        showNotice('Settings saved')
      }
    } else {
      // Revert to current valid value if input is empty
      breweryNameInput.value = breweryName.value
    }
  }
</script>

<style scoped>
  :deep(.v-card-subtitle) {
    overflow: visible;
    text-overflow: unset;
    white-space: normal;
  }
</style>
