<template>
  <!-- Empty state -->
  <v-card v-if="gravityMeasurements.length === 0" class="sub-card" variant="outlined">
    <v-card-text class="text-center py-8">
      <v-icon class="mb-4" color="medium-emphasis" icon="mdi-flask-round-bottom-empty" size="64" />
      <div class="text-h6 text-medium-emphasis mb-2">No fermentation data yet</div>
      <div class="text-body-2 text-medium-emphasis mb-4">
        Record gravity readings on the Timeline tab to see your fermentation curve.
      </div>
      <v-btn
        color="primary"
        min-height="44"
        prepend-icon="mdi-timeline-text"
        variant="tonal"
        @click="emit('go-to-timeline')"
      >
        Go to Timeline
      </v-btn>
    </v-card-text>
  </v-card>

  <!-- Chart and stats -->
  <div v-else>
    <!-- Rotate hint for mobile -->
    <v-alert
      v-if="showRotateHint && xs"
      class="mb-3"
      closable
      density="compact"
      prepend-icon="mdi-screen-rotation"
      type="info"
      variant="tonal"
      @click:close="dismissRotateHint"
    >
      Rotate for a better view
    </v-alert>

    <!-- Chart -->
    <v-card class="sub-card mb-4" variant="outlined">
      <v-card-text>
        <div class="chart-container">
          <Line
            :data="chartData"
            :options="chartOptions"
          />
        </div>
      </v-card-text>
    </v-card>

    <!-- Toggle controls -->
    <div class="d-flex justify-center mb-4">
      <v-btn-toggle
        v-model="visibleDatasets"
        density="compact"
        multiple
        variant="outlined"
      >
        <v-btn value="gravity" min-height="44">
          <v-icon icon="mdi-flask" start />
          <span class="d-none d-sm-inline">Gravity</span>
        </v-btn>
        <v-btn value="temperature" min-height="44">
          <v-icon icon="mdi-thermometer" start />
          <span class="d-none d-sm-inline">Temperature</span>
        </v-btn>
      </v-btn-toggle>
    </div>

    <!-- Stats summary -->
    <v-row dense>
      <v-col cols="6" md="4">
        <div class="metric-card">
          <div class="metric-label">OG</div>
          <div class="metric-value">{{ ogDisplay }}</div>
        </div>
      </v-col>
      <v-col cols="6" md="4">
        <div class="metric-card">
          <div class="metric-label">Current</div>
          <div class="metric-value">{{ currentGravityDisplay }}</div>
        </div>
      </v-col>
      <v-col cols="6" md="4">
        <div class="metric-card">
          <div class="metric-label">Target FG</div>
          <div class="metric-value">{{ targetFgDisplay }}</div>
        </div>
      </v-col>
      <v-col cols="6" md="4">
        <div class="metric-card">
          <div class="metric-label">Attenuation</div>
          <div class="metric-value">{{ attenuationDisplay }}</div>
        </div>
      </v-col>
      <v-col cols="6" md="4">
        <div class="metric-card">
          <div class="metric-label">Days</div>
          <div class="metric-value">{{ daysDisplay }}</div>
        </div>
      </v-col>
      <v-col cols="6" md="4">
        <div class="metric-card">
          <div class="metric-label">Temp</div>
          <div class="metric-value">{{ latestTempDisplay }}</div>
        </div>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts" setup>
  import type { ChartData, ChartOptions } from 'chart.js'
  import type { AnnotationOptions } from 'chartjs-plugin-annotation'
  import type { BatchSummary, GravityUnit, Measurement, TemperatureUnit } from '@/types'
  import {
    Chart as ChartJS,
    Legend,
    LinearScale,
    LineElement,
    PointElement,
    Tooltip,
  } from 'chart.js'
  import annotationPlugin from 'chartjs-plugin-annotation'
  import { computed, ref } from 'vue'
  import { useDisplay } from 'vuetify'
  import {
    convertGravity,
    convertTemperature,
    gravityLabels,
    gravityPrecision,
    temperatureLabels,
    temperaturePrecision,
  } from '@/composables/useUnitConversion'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { Line } from 'vue-chartjs'

  ChartJS.register(
    Legend,
    LinearScale,
    LineElement,
    PointElement,
    Tooltip,
    annotationPlugin,
  )

  const props = defineProps<{
    measurements: Measurement[]
    targetOg?: number | null
    targetFg?: number | null
    batchSummary?: BatchSummary | null
  }>()

  const emit = defineEmits<{
    'go-to-timeline': []
  }>()

  const { xs } = useDisplay()
  const { preferences, formatGravityPreferred, formatTemperaturePreferred } = useUnitPreferences()

  // Toggle state
  const visibleDatasets = ref<string[]>(['gravity', 'temperature'])

  // Rotate hint (show once, dismiss via localStorage)
  const ROTATE_HINT_KEY = 'brewpipes:fermentationRotateHintDismissed'
  const showRotateHint = ref(!localStorage.getItem(ROTATE_HINT_KEY))

  function dismissRotateHint () {
    showRotateHint.value = false
    try {
      localStorage.setItem(ROTATE_HINT_KEY, '1')
    } catch {
      // localStorage unavailable
    }
  }

  // ==================== Data Filtering & Sorting ====================

  const gravityMeasurements = computed(() =>
    props.measurements
      .filter(m => m.kind === 'gravity')
      .sort((a, b) => new Date(a.observed_at).getTime() - new Date(b.observed_at).getTime()),
  )

  const temperatureMeasurements = computed(() =>
    props.measurements
      .filter(m => m.kind === 'temperature')
      .sort((a, b) => new Date(a.observed_at).getTime() - new Date(b.observed_at).getTime()),
  )

  // ==================== Time Helpers ====================

  /** Earliest measurement timestamp across gravity and temperature */
  const firstTimestamp = computed(() => {
    const allSorted = [...gravityMeasurements.value, ...temperatureMeasurements.value]
      .sort((a, b) => new Date(a.observed_at).getTime() - new Date(b.observed_at).getTime())
    if (allSorted.length === 0) return Date.now()
    return new Date(allSorted[0]!.observed_at).getTime()
  })

  /** Convert a timestamp to relative days from the first measurement */
  function toDays (observedAt: string): number {
    const ms = new Date(observedAt).getTime() - firstTimestamp.value
    return ms / (1000 * 60 * 60 * 24)
  }

  /** Format a date for tooltip display */
  function formatDateForTooltip (observedAt: string): string {
    const d = new Date(observedAt)
    return d.toLocaleString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit',
    })
  }

  // ==================== Unit Helpers ====================

  function normalizeTemperatureUnit (unit: string | null | undefined): TemperatureUnit {
    if (!unit) return 'c' // Storage default is Celsius
    const normalized = unit.trim().toLowerCase()
    if (normalized === 'f' || normalized === 'fahrenheit' || normalized === '°f') return 'f'
    return 'c'
  }

  function normalizeGravityUnit (unit: string | null | undefined): GravityUnit {
    if (!unit) return 'sg' // Storage default is SG
    const normalized = unit.trim().toLowerCase()
    if (normalized === 'plato' || normalized === '°p' || normalized === 'p') return 'plato'
    return 'sg'
  }

  const gravityUnit = computed(() => preferences.value.gravity)
  const tempUnit = computed(() => preferences.value.temperature)

  // ==================== Chart Data ====================

  const gravityDataPoints = computed(() =>
    gravityMeasurements.value.map(m => ({
      x: toDays(m.observed_at),
      y: convertGravity(m.value, normalizeGravityUnit(m.unit), gravityUnit.value) ?? m.value,
      raw: m,
    })),
  )

  const temperatureDataPoints = computed(() =>
    temperatureMeasurements.value.map(m => ({
      x: toDays(m.observed_at),
      y: convertTemperature(m.value, normalizeTemperatureUnit(m.unit), tempUnit.value) ?? m.value,
      raw: m,
    })),
  )

  const showGravity = computed(() => visibleDatasets.value.includes('gravity'))
  const showTemperature = computed(() => visibleDatasets.value.includes('temperature'))

  const chartData = computed<ChartData<'line'>>(() => ({
    datasets: [
      {
        label: `Gravity (${gravityLabels[gravityUnit.value]})`,
        data: gravityDataPoints.value.map(p => ({ x: p.x, y: p.y })),
        borderColor: '#1976D2',
        backgroundColor: 'rgba(25, 118, 210, 0.1)',
        pointBackgroundColor: '#1976D2',
        pointRadius: 5,
        pointHoverRadius: 7,
        tension: 0.3,
        yAxisID: 'yGravity',
        hidden: !showGravity.value,
      },
      {
        label: `Temperature (${temperatureLabels[tempUnit.value]})`,
        data: temperatureDataPoints.value.map(p => ({ x: p.x, y: p.y })),
        borderColor: '#FF9800',
        backgroundColor: 'rgba(255, 152, 0, 0.1)',
        pointBackgroundColor: '#FF9800',
        pointRadius: 5,
        pointHoverRadius: 7,
        tension: 0.3,
        yAxisID: 'yTemp',
        hidden: !showTemperature.value,
      },
    ],
  }))

  // ==================== Annotation Lines ====================

  const ogValue = computed(() => {
    if (props.targetOg !== null && props.targetOg !== undefined) {
      return convertGravity(props.targetOg, 'sg', gravityUnit.value)
    }
    if (props.batchSummary?.original_gravity !== null && props.batchSummary?.original_gravity !== undefined) {
      return convertGravity(props.batchSummary.original_gravity, 'sg', gravityUnit.value)
    }
    return null
  })

  const targetFgValue = computed(() => {
    if (props.targetFg !== null && props.targetFg !== undefined) {
      return convertGravity(props.targetFg, 'sg', gravityUnit.value)
    }
    return null
  })

  const annotations = computed(() => {
    const result: Record<string, AnnotationOptions> = {}
    if (ogValue.value !== null && showGravity.value) {
      result.ogLine = {
        type: 'line' as const,
        yMin: ogValue.value,
        yMax: ogValue.value,
        yScaleID: 'yGravity',
        borderColor: 'rgba(0, 0, 0, 0.25)',
        borderWidth: 1,
        borderDash: [6, 4],
        label: {
          display: true,
          content: 'OG',
          position: 'start' as const,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          color: '#fff',
          font: { size: 10 },
          padding: 3,
        },
      }
    }
    if (targetFgValue.value !== null && showGravity.value) {
      result.fgLine = {
        type: 'line' as const,
        yMin: targetFgValue.value,
        yMax: targetFgValue.value,
        yScaleID: 'yGravity',
        borderColor: 'rgba(0, 0, 0, 0.25)',
        borderWidth: 1,
        borderDash: [6, 4],
        label: {
          display: true,
          content: 'Target FG',
          position: 'end' as const,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          color: '#fff',
          font: { size: 10 },
          padding: 3,
        },
      }
    }
    return result
  })

  // ==================== Chart Options ====================

  // Build tooltip metadata maps keyed by dataset+index
  const gravityTooltipDates = computed(() =>
    gravityMeasurements.value.map(m => formatDateForTooltip(m.observed_at)),
  )

  const temperatureTooltipDates = computed(() =>
    temperatureMeasurements.value.map(m => formatDateForTooltip(m.observed_at)),
  )

  const chartOptions = computed<ChartOptions<'line'>>(() => ({
    responsive: true,
    maintainAspectRatio: true,
    aspectRatio: xs.value ? 1.5 : 16 / 9,
    animation: {
      duration: 300,
    },
    interaction: {
      mode: 'nearest' as const,
      intersect: true,
    },
    plugins: {
      legend: {
        display: true,
        position: 'bottom' as const,
        onClick: () => {
          // Disable legend click toggling — use the toggle buttons instead
        },
      },
      tooltip: {
        callbacks: {
          title: (items) => {
            if (items.length === 0) return ''
            const item = items[0]!
            const dates = item.datasetIndex === 0 ? gravityTooltipDates.value : temperatureTooltipDates.value
            return dates[item.dataIndex] ?? `Day ${(item.parsed.x ?? 0).toFixed(1)}`
          },
          label: (item) => {
            const precision = item.datasetIndex === 0
              ? gravityPrecision[gravityUnit.value]
              : temperaturePrecision[tempUnit.value]
            const unit = item.datasetIndex === 0
              ? gravityLabels[gravityUnit.value]
              : temperatureLabels[tempUnit.value]
            const yValue = item.parsed.y ?? 0
            return `${item.dataset.label}: ${yValue.toFixed(precision)} ${unit}`
          },
        },
      },
      annotation: {
        annotations: annotations.value,
      },
    },
    scales: {
      x: {
        type: 'linear' as const,
        title: {
          display: true,
          text: 'Days',
        },
        ticks: {
          callback: (value) => `Day ${Math.round(Number(value))}`,
          stepSize: 1,
        },
        min: 0,
      },
      yGravity: {
        type: 'linear' as const,
        position: 'left' as const,
        display: showGravity.value,
        title: {
          display: true,
          text: gravityLabels[gravityUnit.value],
        },
        ticks: {
          callback: (value) => Number(value).toFixed(gravityPrecision[gravityUnit.value]),
        },
        grace: '5%',
      },
      yTemp: {
        type: 'linear' as const,
        position: 'right' as const,
        display: showTemperature.value,
        title: {
          display: true,
          text: temperatureLabels[tempUnit.value],
        },
        grid: {
          drawOnChartArea: false,
        },
        grace: '5%',
      },
    },
  }))

  // ==================== Stats ====================

  const ogRaw = computed(() => {
    if (props.batchSummary?.original_gravity !== null && props.batchSummary?.original_gravity !== undefined) {
      return props.batchSummary.original_gravity
    }
    // Fall back to first gravity measurement
    return gravityMeasurements.value.length > 0 ? gravityMeasurements.value[0]!.value : null
  })

  const currentGravityRaw = computed(() => {
    if (gravityMeasurements.value.length === 0) return null
    return gravityMeasurements.value.at(-1)!.value
  })

  const ogDisplay = computed(() => {
    if (ogRaw.value === null) return '—'
    return formatGravityPreferred(ogRaw.value, 'sg')
  })

  const currentGravityDisplay = computed(() => {
    if (currentGravityRaw.value === null) return '—'
    return formatGravityPreferred(currentGravityRaw.value, 'sg')
  })

  const targetFgDisplay = computed(() => {
    if (props.targetFg === null || props.targetFg === undefined) return '—'
    return formatGravityPreferred(props.targetFg, 'sg')
  })

  const attenuationDisplay = computed(() => {
    if (ogRaw.value === null || currentGravityRaw.value === null) return '—'
    const denominator = ogRaw.value - 1.0
    if (denominator <= 0) return '—'
    const att = ((ogRaw.value - currentGravityRaw.value) / denominator) * 100
    return `${att.toFixed(1)}%`
  })

  const daysDisplay = computed(() => {
    if (gravityMeasurements.value.length < 1) return '—'
    const first = new Date(gravityMeasurements.value[0]!.observed_at).getTime()
    const last = new Date(gravityMeasurements.value.at(-1)!.observed_at).getTime()
    const days = Math.round((last - first) / (1000 * 60 * 60 * 24))
    return String(days)
  })

  const latestTempDisplay = computed(() => {
    if (temperatureMeasurements.value.length === 0) return '—'
    const latest = temperatureMeasurements.value.at(-1)!
    const sourceUnit = normalizeTemperatureUnit(latest.unit)
    return formatTemperaturePreferred(latest.value, sourceUnit)
  })
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.chart-container {
  position: relative;
  width: 100%;
  min-height: 200px;
}

.metric-card {
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(var(--v-theme-surface), 0.5);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  text-align: center;
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
</style>
