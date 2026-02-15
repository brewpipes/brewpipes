<template>
  <v-row>
    <v-col cols="12">
      <v-card class="sub-card" variant="outlined">
        <v-card-title class="text-subtitle-1 d-flex align-center">
          Addition log
          <v-spacer />
          <v-btn
            aria-label="Record addition"
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
                <th>Type</th>
                <th>Amount</th>
                <th>Target</th>
                <th>Time</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="addition in additionsSorted" :key="addition.uuid">
                <td>{{ formatAdditionType(addition.addition_type) }}</td>
                <td>{{ formatAmount(addition.amount, addition.amount_unit) }}</td>
                <td>{{ addition.occupancy_uuid ? 'Occupancy' : 'Batch' }}</td>
                <td>{{ formatDateTime(addition.added_at) }}</td>
              </tr>
              <tr v-if="additionsSorted.length === 0">
                <td colspan="4">No additions recorded.</td>
              </tr>
            </tbody>
          </v-table>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
  import type { Addition } from './types'
  import { computed } from 'vue'
  import { useAdditionTypeFormatters, useFormatters } from '@/composables/useFormatters'

  const props = defineProps<{
    additions: Addition[]
  }>()

  const emit = defineEmits<{
    create: []
  }>()

  const { formatDateTime } = useFormatters()
  const { formatAdditionType } = useAdditionTypeFormatters()

  const additionsSorted = computed(() =>
    sortByTime(props.additions, item => item.added_at),
  )

  function formatAmount (amount: number | null, unit: string | null | undefined) {
    if (amount === null || amount === undefined) {
      return 'Unknown'
    }
    return `${amount} ${unit ?? ''}`.trim()
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
