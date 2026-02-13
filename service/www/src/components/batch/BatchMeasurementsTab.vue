<template>
  <v-row>
    <v-col cols="12">
      <v-card class="sub-card" variant="outlined">
        <v-card-title class="text-subtitle-1 d-flex align-center">
          Measurement log
          <v-spacer />
          <v-btn
            aria-label="Record measurement"
            icon="mdi-plus"
            size="small"
            variant="text"
            @click="emit('create')"
          />
        </v-card-title>
        <v-card-text>
          <v-table class="data-table" density="compact">
            <thead>
              <tr>
                <th>Kind</th>
                <th>Value</th>
                <th>Target</th>
                <th>Time</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="measurement in measurementsSorted" :key="measurement.uuid">
                <td>{{ measurement.kind }}</td>
                <td>
                  {{
                    isNoteMeasurement(measurement)
                      ? measurement.notes ?? 'Note'
                      : formatValue(measurement.value, measurement.unit)
                  }}
                </td>
                <td>{{ measurement.occupancy_uuid ? 'Occupancy' : 'Batch' }}</td>
                <td>{{ formatDateTime(measurement.observed_at) }}</td>
              </tr>
              <tr v-if="measurementsSorted.length === 0">
                <td colspan="4">No measurements recorded.</td>
              </tr>
            </tbody>
          </v-table>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
  import type { Measurement } from './types'
  import { computed } from 'vue'
  import { useFormatters } from '@/composables/useFormatters'

  const props = defineProps<{
    measurements: Measurement[]
  }>()

  const emit = defineEmits<{
    create: []
  }>()

  const { formatDateTime } = useFormatters()

  const measurementsSorted = computed(() =>
    sortByTime(props.measurements, item => item.observed_at),
  )

  function formatValue (value: number | null, unit: string | null | undefined) {
    if (value === null || value === undefined) {
      return 'Unknown'
    }
    return `${value}${unit ? ` ${unit}` : ''}`
  }

  function isNoteMeasurement (measurement: Measurement) {
    const normalized = measurement.kind.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
    return normalized === 'note'
  }

  function toTimestamp (value: string | null | undefined) {
    if (!value) {
      return 0
    }
    return new Date(value).getTime()
  }

  function sortByTime<T> (items: T[], selector: (item: T) => string | null | undefined) {
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...items].sort((a, b) => toTimestamp(selector(b)) - toTimestamp(selector(a)))
  }
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
}

.data-table :deep(td) {
  font-size: 0.85rem;
}
</style>
