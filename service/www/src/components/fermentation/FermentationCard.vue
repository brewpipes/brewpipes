<template>
  <v-card
    class="section-card fermentation-card"
    :class="attentionClass"
  >
    <!-- Header row -->
    <v-card-title class="d-flex align-center pb-1">
      <span class="text-subtitle-1 font-weight-bold text-truncate">
        {{ vesselName }}
      </span>
      <v-spacer />
      <v-chip
        class="ml-2 flex-shrink-0"
        color="primary"
        size="x-small"
        variant="tonal"
      >
        Day {{ daysInTank }}
      </v-chip>
      <v-menu location="bottom end" offset="4">
        <template #activator="{ props: menuProps }">
          <v-btn
            aria-label="Card actions"
            class="ml-1"
            density="compact"
            icon="mdi-dots-vertical"
            size="small"
            variant="text"
            v-bind="menuProps"
          />
        </template>
        <v-list density="compact">
          <v-list-item
            prepend-icon="mdi-eye"
            title="View Batch"
            @click="$router.push(`/batches/${occupancy.batch_uuid}`)"
          />
          <v-list-item
            prepend-icon="mdi-swap-horizontal"
            title="Change Status"
            @click="showStatusMenu = true"
          />
          <v-list-item
            prepend-icon="mdi-arrow-right-bold"
            title="Transfer"
            @click="emit('transfer')"
          />
          <v-list-item
            prepend-icon="mdi-call-split"
            title="Split"
            @click="emit('split')"
          />
          <v-list-item
            prepend-icon="mdi-call-merge"
            title="Blend"
            @click="emit('blend')"
          />
          <v-list-item
            prepend-icon="mdi-package-variant"
            title="Package"
            @click="emit('package')"
          />
          <v-list-item
            prepend-icon="mdi-flask-empty-outline"
            title="Mark Empty"
            @click="emit('markEmpty')"
          />
        </v-list>
      </v-menu>
    </v-card-title>

    <!-- Batch info row -->
    <v-card-subtitle class="pt-0 pb-2">
      <router-link
        class="batch-link text-body-2 font-weight-medium"
        :to="`/batches/${occupancy.batch_uuid}`"
      >
        {{ batchName }}
      </router-link>
      <span v-if="recipeName" class="text-caption text-medium-emphasis ml-1">
        · {{ recipeName }}
      </span>
    </v-card-subtitle>

    <v-card-text class="pt-0">
      <!-- Status chip -->
      <v-chip
        class="mb-3"
        :color="statusColor"
        size="small"
        variant="tonal"
      >
        {{ statusLabel }}
      </v-chip>

      <!-- Metrics section -->
      <v-row class="mb-2" dense>
        <v-col cols="6">
          <BatchSparklineCard
            color="primary"
            label="Gravity"
            :latest-label="latestGravityLabel"
            :values="gravityValues"
          />
        </v-col>
        <v-col cols="6">
          <BatchSparklineCard
            color="info"
            label="Temp"
            :latest-label="latestTemperatureLabel"
            :values="temperatureValues"
          />
        </v-col>
      </v-row>

      <v-row dense>
        <v-col cols="6">
          <div class="text-caption text-medium-emphasis">Attenuation</div>
          <div class="text-body-2 font-weight-medium">{{ attenuationLabel }}</div>
        </v-col>
        <v-col cols="6">
          <div class="text-caption text-medium-emphasis">Est. ABV</div>
          <div class="text-body-2 font-weight-medium">{{ abvLabel }}</div>
        </v-col>
      </v-row>

      <!-- Log Reading button -->
      <v-btn
        block
        class="mt-3"
        color="primary"
        min-height="44"
        variant="tonal"
        @click="emit('logReading')"
      >
        Log Reading
      </v-btn>
    </v-card-text>

    <!-- Change Status Dialog -->
    <v-dialog
      v-model="showStatusMenu"
      :fullscreen="xs"
      :max-width="xs ? '100%' : 360"
      persistent
    >
      <v-card>
        <v-card-title class="text-h6">Change Status</v-card-title>
        <v-card-text>
          <v-list density="compact">
            <v-list-item
              v-for="status in OCCUPANCY_STATUS_VALUES"
              :key="status"
              :active="status === occupancy.status"
              @click="handleStatusChange(status)"
            >
              <template #prepend>
                <v-chip
                  :color="getOccupancyStatusColor(status)"
                  size="x-small"
                  variant="tonal"
                >
                  {{ formatOccupancyStatus(status) }}
                </v-chip>
              </template>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions class="justify-end">
          <v-btn variant="text" @click="showStatusMenu = false">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script lang="ts" setup>
  import type { BatchSummary, Measurement, Occupancy, OccupancyStatus, Vessel } from '@/types'
  import { computed, ref } from 'vue'
  import { useDisplay } from 'vuetify'
  import { BatchSparklineCard } from '@/components/batch'
  import { useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { OCCUPANCY_STATUS_VALUES } from '@/types'

  const props = defineProps<{
    occupancy: Occupancy
    vessel: Vessel
    batchSummary: BatchSummary | null
    measurements: Measurement[]
  }>()

  const emit = defineEmits<{
    blend: []
    logReading: []
    markEmpty: []
    package: []
    split: []
    statusChanged: []
    transfer: []
  }>()

  const { xs } = useDisplay()
  const { formatGravityPreferred, formatTemperaturePreferred } = useUnitPreferences()
  const { formatOccupancyStatus, getOccupancyStatusColor } = useOccupancyStatusFormatters()
  const { updateOccupancyStatus } = useProductionApi()
  const { showNotice } = useSnackbar()

  const showStatusMenu = ref(false)

  // Derived data
  const vesselName = computed(() => props.vessel.name)
  const batchName = computed(() => props.batchSummary?.short_name ?? `Batch ${props.occupancy.batch_uuid?.slice(0, 8) ?? '—'}`)
  const recipeName = computed(() => props.batchSummary?.recipe_name ?? null)
  const statusLabel = computed(() => formatOccupancyStatus(props.occupancy.status))
  const statusColor = computed(() => getOccupancyStatusColor(props.occupancy.status))

  // Days in tank
  const daysInTank = computed(() => {
    const inAt = new Date(props.occupancy.in_at)
    const now = new Date()
    const diffMs = now.getTime() - inAt.getTime()
    return Math.max(0, Math.floor(diffMs / (1000 * 60 * 60 * 24)))
  })

  // Gravity measurements sorted by observed_at
  const gravityMeasurements = computed(() =>
    props.measurements
      .filter(m => m.kind === 'gravity')
      .sort((a, b) => new Date(a.observed_at || a.created_at || 0).getTime() - new Date(b.observed_at || b.created_at || 0).getTime()),
  )

  // Temperature measurements sorted by observed_at
  const temperatureMeasurements = computed(() =>
    props.measurements
      .filter(m => m.kind === 'temperature')
      .sort((a, b) => new Date(a.observed_at || a.created_at || 0).getTime() - new Date(b.observed_at || b.created_at || 0).getTime()),
  )

  // Sparkline values (raw numbers for the chart)
  const gravityValues = computed(() => gravityMeasurements.value.map(m => m.value))
  const temperatureValues = computed(() => temperatureMeasurements.value.map(m => m.value))

  // Latest values formatted with user preferences
  const latestGravity = computed(() => gravityMeasurements.value.at(-1) ?? null)
  const latestTemperature = computed(() => temperatureMeasurements.value.at(-1) ?? null)

  const latestGravityLabel = computed(() => {
    if (!latestGravity.value) return '—'
    return formatGravityPreferred(latestGravity.value.value, 'sg')
  })

  const latestTemperatureLabel = computed(() => {
    if (!latestTemperature.value) return '—'
    return formatTemperaturePreferred(latestTemperature.value.value, 'c')
  })

  // OG from batch summary, falling back to first gravity reading
  const og = computed(() => {
    if (props.batchSummary?.original_gravity != null) {
      return props.batchSummary.original_gravity
    }
    // Fall back to first gravity reading (same approach as FermentationCurve)
    const first = gravityMeasurements.value[0]
    if (first) {
      return first.value
    }
    return null
  })

  // Attenuation: (OG - currentGravity) / (OG - 1.0) * 100
  const attenuation = computed(() => {
    if (og.value === null || !latestGravity.value) return null
    const currentGravity = latestGravity.value.value
    const denominator = og.value - 1
    if (denominator <= 0) return null
    return ((og.value - currentGravity) / denominator) * 100
  })

  const attenuationLabel = computed(() => {
    if (attenuation.value === null) return '—'
    return `${attenuation.value.toFixed(1)}%`
  })

  // ABV: use batch summary ABV if available, otherwise calculate
  const abv = computed(() => {
    if (props.batchSummary?.abv !== null && props.batchSummary?.abv !== undefined) {
      return props.batchSummary.abv
    }
    if (og.value === null || !latestGravity.value) return null
    return (og.value - latestGravity.value.value) * 131.25
  })

  const abvLabel = computed(() => {
    if (abv.value === null) return '—'
    return `${abv.value.toFixed(1)}%`
  })

  // Attention indicators
  const STALE_HOURS = 24
  const STABLE_READINGS_COUNT = 3
  const STABLE_TOLERANCE = 0.001

  const isStaleGravity = computed(() => {
    if (!latestGravity.value) return true // No readings at all is stale
    const lastAt = new Date(latestGravity.value.observed_at || latestGravity.value.created_at || 0)
    const now = new Date()
    const diffHours = (now.getTime() - lastAt.getTime()) / (1000 * 60 * 60)
    return diffHours >= STALE_HOURS
  })

  const isStableGravity = computed(() => {
    const readings = gravityMeasurements.value
    if (readings.length < STABLE_READINGS_COUNT) return false
    const recent = readings.slice(-STABLE_READINGS_COUNT)
    const first = recent[0]!.value
    return recent.every(r => Math.abs(r.value - first) < STABLE_TOLERANCE)
  })

  const attentionLevel = computed<'warning' | 'info' | null>(() => {
    if (isStaleGravity.value) return 'warning'
    if (isStableGravity.value) return 'info'
    return null
  })

  const attentionClass = computed(() => {
    if (attentionLevel.value === 'warning') return 'attention-warning'
    if (attentionLevel.value === 'info') return 'attention-info'
    return ''
  })

  // Status change handler
  async function handleStatusChange (status: OccupancyStatus) {
    try {
      await updateOccupancyStatus(props.occupancy.uuid, status)
      showNotice(`Status changed to ${formatOccupancyStatus(status)}`)
      showStatusMenu.value = false
      emit('statusChanged')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to change status'
      showNotice(message, 'error')
    }
  }
</script>

<style scoped>
.fermentation-card {
  transition: border-left-color 0.2s ease;
}

.attention-warning {
  border-left: 4px solid rgb(var(--v-theme-warning));
}

.attention-info {
  border-left: 4px solid rgb(var(--v-theme-info));
}
</style>
