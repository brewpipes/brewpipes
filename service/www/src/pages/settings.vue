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
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="preferences.temperature"
              :items="temperatureOptions"
              item-title="label"
              item-value="value"
              label="Temperature"
              variant="outlined"
              density="comfortable"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="preferences.gravity"
              :items="gravityOptions"
              item-title="label"
              item-value="value"
              label="Gravity"
              variant="outlined"
              density="comfortable"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="preferences.volume"
              :items="volumeOptions"
              item-title="label"
              item-value="value"
              label="Volume"
              variant="outlined"
              density="comfortable"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="preferences.mass"
              :items="massOptions"
              item-title="label"
              item-value="value"
              label="Mass"
              variant="outlined"
              density="comfortable"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="preferences.pressure"
              :items="pressureOptions"
              item-title="label"
              item-value="value"
              label="Pressure"
              variant="outlined"
              density="comfortable"
              hide-details
            />
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="preferences.color"
              :items="colorOptions"
              item-title="label"
              item-value="value"
              label="Color"
              variant="outlined"
              density="comfortable"
              hide-details
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
