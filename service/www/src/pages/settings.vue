<template>
  <v-container>
    <h1 class="text-h4 mb-6">Settings</h1>

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
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

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
</script>
