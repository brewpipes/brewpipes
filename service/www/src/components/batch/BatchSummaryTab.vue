<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-text>
      <v-progress-linear
        v-if="loading"
        class="mb-4"
        color="primary"
        indeterminate
      />

      <template v-if="summary">
        <!-- Recipe & Style -->
        <v-row class="mb-4">
          <v-col cols="12" md="6">
            <div class="text-overline text-medium-emphasis">Recipe</div>
            <div class="text-h6">{{ summary.recipe_name ?? 'Not specified' }}</div>
          </v-col>
          <v-col cols="12" md="6">
            <div class="text-overline text-medium-emphasis">Style</div>
            <div class="text-h6">{{ summary.style_name ?? 'Not specified' }}</div>
          </v-col>
        </v-row>

        <v-divider class="mb-4" />

        <!-- Status & Location -->
        <v-row class="mb-4">
          <v-col cols="12" md="4">
            <div class="text-overline text-medium-emphasis">Current Phase</div>
            <v-chip
              v-if="summary.current_phase"
              class="mt-1"
              :color="getPhaseColor(summary.current_phase)"
              variant="tonal"
            >
              {{ formatPhase(summary.current_phase) }}
            </v-chip>
            <div v-else class="text-body-2 text-medium-emphasis mt-1">Not set</div>
          </v-col>
          <v-col cols="12" md="4">
            <div class="text-overline text-medium-emphasis">Current Vessel</div>
            <div class="text-body-1 font-weight-medium mt-1">
              {{ summary.current_vessel ?? 'Not assigned' }}
            </div>
          </v-col>
          <v-col cols="12" md="4">
            <div class="text-overline text-medium-emphasis">Occupancy Status</div>
            <v-menu v-if="summary.current_occupancy_uuid" location="bottom">
              <template #activator="{ props }">
                <v-chip
                  v-bind="props"
                  append-icon="mdi-menu-down"
                  class="mt-1 cursor-pointer"
                  :color="getOccupancyStatusColor(summary.current_occupancy_status)"
                  size="small"
                  variant="tonal"
                >
                  {{ formatOccupancyStatus(summary.current_occupancy_status) }}
                </v-chip>
              </template>
              <v-list density="compact" nav>
                <v-list-subheader>Change status</v-list-subheader>
                <v-list-item
                  v-for="statusOption in occupancyStatusOptions"
                  :key="statusOption.value"
                  :active="statusOption.value === summary.current_occupancy_status"
                    @click="emit('occupancy-status-change', summary.current_occupancy_uuid!, statusOption.value)"
                >
                  <template #prepend>
                    <v-avatar
                      class="mr-2"
                      :color="getOccupancyStatusColor(statusOption.value)"
                      size="24"
                    >
                      <v-icon :icon="getOccupancyStatusIcon(statusOption.value)" size="14" />
                    </v-avatar>
                  </template>
                  <v-list-item-title>{{ statusOption.title }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
            <div v-else class="text-body-2 text-medium-emphasis mt-1">Not set</div>
          </v-col>
        </v-row>

        <v-divider class="mb-4" />

        <!-- Gravity & ABV Metrics -->
        <div class="text-overline text-medium-emphasis mb-2">Gravity & ABV</div>
        <v-row class="mb-4">
          <v-col cols="6" md="3">
            <div class="metric-card">
              <div class="metric-label">OG</div>
              <div class="metric-value">
                {{ formatGravityPreferred(summary.original_gravity, 'sg') }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card">
              <div class="metric-label">FG</div>
              <div class="metric-value">
                {{ formatGravityPreferred(summary.final_gravity, 'sg') }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card">
              <div class="metric-label">ABV</div>
              <div class="metric-value d-flex align-center ga-1">
                {{ formatPercent(summary.abv, 1) }}
                <v-chip
                  v-if="summary.abv !== null && summary.abv !== undefined && summary.abv_calculated"
                  color="info"
                  size="x-small"
                  variant="tonal"
                >
                  calc
                </v-chip>
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card">
              <div class="metric-label">IBU</div>
              <div class="metric-value">
                {{ summary.ibu ?? '—' }}
              </div>
            </div>
          </v-col>
        </v-row>

        <!-- Duration Metrics -->
        <div class="text-overline text-medium-emphasis mb-2">Duration</div>
        <v-row class="mb-4">
          <v-col cols="6" md="4">
            <div class="metric-card">
              <div class="metric-label">Days in FV</div>
              <div class="metric-value">
                {{ formatDays(summary.days_in_fermenter) }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="4">
            <div class="metric-card">
              <div class="metric-label">Days in Brite</div>
              <div class="metric-value">
                {{ formatDays(summary.days_in_brite) }}
              </div>
            </div>
          </v-col>
          <v-col cols="12" md="4">
            <div class="metric-card metric-card--highlight">
              <div class="metric-label">Grain to Glass</div>
              <div class="metric-value">
                {{ formatDays(summary.days_grain_to_glass) }}
              </div>
            </div>
          </v-col>
        </v-row>

        <!-- Volume Metrics -->
        <div class="text-overline text-medium-emphasis mb-2">Volume & Loss</div>
        <v-row class="mb-4">
          <v-col cols="6" md="3">
            <div class="metric-card">
              <div class="metric-label">Starting</div>
              <div class="metric-value">
                {{ formatVolumePreferred(summary.starting_volume_bbl, 'bbl') }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card">
              <div class="metric-label">Current</div>
              <div class="metric-value">
                {{ formatVolumePreferred(summary.current_volume_bbl, 'bbl') }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card">
              <div class="metric-label">Total Loss</div>
              <div class="metric-value">
                {{ formatVolumePreferred(summary.total_loss_bbl, 'bbl') }}
              </div>
            </div>
          </v-col>
          <v-col cols="6" md="3">
            <div class="metric-card" :class="{ 'metric-card--warning': (summary.loss_percentage ?? 0) > 10 }">
              <div class="metric-label">Loss %</div>
              <div class="metric-value">
                {{ formatPercent(summary.loss_percentage, 1) }}
              </div>
            </div>
          </v-col>
        </v-row>

        <v-divider class="mb-4" />

        <!-- Brew Sessions List -->
        <div class="text-overline text-medium-emphasis mb-2">Brew Sessions</div>
        <v-list
          v-if="summary.brew_sessions && summary.brew_sessions.length > 0"
          class="brew-session-summary-list"
          density="compact"
          variant="tonal"
        >
          <v-list-item
            v-for="session in summary.brew_sessions"
            :key="session.uuid"
          >
            <template #prepend>
              <v-icon icon="mdi-kettle-steam" size="small" />
            </template>
            <v-list-item-title>
              {{ formatDateTime(session.brewed_at) }}
            </v-list-item-title>
            <v-list-item-subtitle v-if="session.notes">
              {{ session.notes }}
            </v-list-item-subtitle>
          </v-list-item>
        </v-list>
        <v-alert
          v-else
          density="compact"
          type="info"
          variant="tonal"
        >
          No brew sessions recorded yet.
        </v-alert>

        <!-- Notes -->
        <template v-if="summary.notes">
          <v-divider class="my-4" />
          <div class="text-overline text-medium-emphasis mb-2">Notes</div>
          <div class="text-body-2">{{ summary.notes }}</div>
        </template>
      </template>

      <v-alert
        v-else-if="!loading"
        density="comfortable"
        type="info"
        variant="tonal"
      >
        Summary data not available.
      </v-alert>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useFormatters, useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import {
    type BatchSummary,
    OCCUPANCY_STATUS_VALUES,
    type OccupancyStatus,
  } from '@/composables/useProductionApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  defineProps<{
    summary: BatchSummary | null
    loading: boolean
  }>()

  const emit = defineEmits<{
    'occupancy-status-change': [occupancyUuid: string, status: OccupancyStatus]
  }>()

  const { formatDateTime } = useFormatters()
  const { formatGravityPreferred, formatVolumePreferred } = useUnitPreferences()
  const {
    formatOccupancyStatus,
    getOccupancyStatusColor,
    getOccupancyStatusIcon,
  } = useOccupancyStatusFormatters()

  const occupancyStatusOptions = computed(() =>
    OCCUPANCY_STATUS_VALUES.map(status => ({
      value: status,
      title: formatOccupancyStatus(status),
    })),
  )

  function formatPercent (value: number | null | undefined, decimals = 1) {
    if (value === null || value === undefined) {
      return '—'
    }
    return `${value.toFixed(decimals)}%`
  }

  function formatDays (days: number | null | undefined) {
    if (days === null || days === undefined) {
      return '—'
    }
    if (days < 1) {
      const hours = Math.round(days * 24)
      return `${hours}h`
    }
    return `${days.toFixed(1)} days`
  }

  function formatPhase (phase: string) {
    const phaseLabels: Record<string, string> = {
      planning: 'Planning',
      mashing: 'Mashing',
      heating: 'Heating',
      boiling: 'Boiling',
      cooling: 'Cooling',
      fermenting: 'Fermenting',
      conditioning: 'Conditioning',
      packaging: 'Packaging',
      finished: 'Finished',
    }
    return phaseLabels[phase] ?? phase.charAt(0).toUpperCase() + phase.slice(1)
  }

  function getPhaseColor (phase: string) {
    const phaseColors: Record<string, string> = {
      planning: 'grey',
      mashing: 'orange',
      heating: 'deep-orange',
      boiling: 'red',
      cooling: 'cyan',
      fermenting: 'primary',
      conditioning: 'teal',
      packaging: 'blue',
      finished: 'success',
    }
    return phaseColors[phase] ?? 'secondary'
  }
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.metric-card {
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(var(--v-theme-surface), 0.5);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  text-align: center;
}

.metric-card--highlight {
  background: rgba(var(--v-theme-primary), 0.08);
  border-color: rgba(var(--v-theme-primary), 0.2);
}

.metric-card--warning {
  background: rgba(var(--v-theme-warning), 0.08);
  border-color: rgba(var(--v-theme-warning), 0.2);
}

.metric-label {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-bottom: 4px;
}

.metric-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.87);
}

.brew-session-summary-list {
  border-radius: 8px;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
