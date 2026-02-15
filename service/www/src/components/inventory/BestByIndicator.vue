<template>
  <template v-if="chipOnly">
    <v-chip
      v-if="isExpired"
      color="error"
      density="compact"
      size="small"
      variant="tonal"
    >
      Expired
    </v-chip>
    <v-chip
      v-else-if="isApproaching"
      color="warning"
      density="compact"
      size="small"
      variant="tonal"
    >
      {{ formattedDate }}
    </v-chip>
    <span v-else-if="bestBy" class="text-body-2">{{ formattedDate }}</span>
  </template>
  <template v-else>
    <span v-if="isExpired" class="d-inline-flex align-center">
      <span class="text-error mr-1">{{ formattedDate }}</span>
      <v-chip color="error" density="compact" size="small" variant="tonal">
        Expired
      </v-chip>
    </span>
    <span v-else-if="isApproaching" class="text-warning">
      {{ formattedDate }}
    </span>
    <span v-else>
      {{ formattedDate }}
    </span>
  </template>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { formatDate } from '@/composables/useFormatters'

  const props = defineProps<{
    bestBy?: string | null
    chipOnly?: boolean
  }>()

  const THIRTY_DAYS_MS = 30 * 24 * 60 * 60 * 1000

  const bestByDate = computed(() => {
    if (!props.bestBy) return null
    const d = new Date(props.bestBy)
    return Number.isNaN(d.getTime()) ? null : d
  })

  const isExpired = computed(() => {
    if (!bestByDate.value) return false
    return bestByDate.value.getTime() < Date.now()
  })

  const isApproaching = computed(() => {
    if (!bestByDate.value || isExpired.value) return false
    return bestByDate.value.getTime() - Date.now() < THIRTY_DAYS_MS
  })

  const formattedDate = computed(() => {
    if (!props.bestBy) return '\u2014'
    return formatDate(props.bestBy)
  })
</script>
