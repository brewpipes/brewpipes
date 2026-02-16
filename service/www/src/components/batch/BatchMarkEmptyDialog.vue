<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 520"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Mark Vessel as Empty</v-card-title>
      <v-card-text>
        <v-alert
          v-if="errorMessage"
          class="mb-4"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ errorMessage }}
        </v-alert>

        <!-- Context info (read-only) -->
        <v-list class="mb-4" density="compact">
          <v-list-item>
            <template #prepend>
              <v-icon icon="mdi-flask-round-bottom" size="small" />
            </template>
            <v-list-item-title>{{ vesselName }}</v-list-item-title>
            <v-list-item-subtitle>Vessel</v-list-item-subtitle>
          </v-list-item>
          <v-list-item>
            <template #prepend>
              <v-icon icon="mdi-beaker-outline" size="small" />
            </template>
            <v-list-item-title>{{ batchName }}</v-list-item-title>
            <v-list-item-subtitle>Batch</v-list-item-subtitle>
          </v-list-item>
          <v-list-item v-if="occupancy?.status">
            <template #prepend>
              <v-icon icon="mdi-information-outline" size="small" />
            </template>
            <v-list-item-title>
              <v-chip
                :color="getOccupancyStatusColor(occupancy.status)"
                size="small"
                variant="tonal"
              >
                {{ formatOccupancyStatus(occupancy.status) }}
              </v-chip>
            </v-list-item-title>
            <v-list-item-subtitle>Beer status</v-list-item-subtitle>
          </v-list-item>
          <v-list-item v-if="durationLabel">
            <template #prepend>
              <v-icon icon="mdi-clock-outline" size="small" />
            </template>
            <v-list-item-title>{{ durationLabel }}</v-list-item-title>
            <v-list-item-subtitle>Duration in vessel</v-list-item-subtitle>
          </v-list-item>
        </v-list>

        <v-form :disabled="saving" @submit.prevent="handleSubmit">
          <v-text-field
            v-model="form.out_at"
            density="comfortable"
            label="Emptied at"
            type="datetime-local"
          />
        </v-form>

        <div class="text-body-2 text-medium-emphasis mt-2">
          This will mark the vessel as empty and available.
        </div>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="saving"
          @click="handleSubmit"
        >
          Mark Empty
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Occupancy } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { nowInputValue } from '@/utils/normalize'

  const props = defineProps<{
    modelValue: boolean
    occupancy: Occupancy | null
    vesselName: string
    batchName: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'emptied': []
  }>()

  const { closeOccupancy } = useProductionApi()
  const {
    formatOccupancyStatus,
    getOccupancyStatusColor,
  } = useOccupancyStatusFormatters()

  const saving = ref(false)
  const errorMessage = ref('')

  const form = reactive({
    out_at: '',
  })

  /** Compute a human-readable duration from in_at to now */
  const durationLabel = computed(() => {
    if (!props.occupancy?.in_at) return ''
    const inAt = new Date(props.occupancy.in_at)
    const now = new Date()
    const diffMs = now.getTime() - inAt.getTime()
    if (diffMs < 0) return ''

    const totalHours = diffMs / (1000 * 60 * 60)
    if (totalHours < 1) {
      const minutes = Math.floor(diffMs / (1000 * 60))
      return `${minutes}m`
    }
    if (totalHours < 24) {
      const hours = Math.floor(totalHours)
      return `${hours}h`
    }
    const days = totalHours / 24
    return `${days.toFixed(1)} days`
  })

  // Reset form when dialog opens
  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        errorMessage.value = ''
        form.out_at = nowInputValue()
      }
    },
  )

  async function handleSubmit () {
    if (!props.occupancy) return

    saving.value = true
    errorMessage.value = ''

    try {
      const outAt = form.out_at ? new Date(form.out_at).toISOString() : undefined
      await closeOccupancy(props.occupancy.uuid, outAt)
      emit('update:modelValue', false)
      emit('emptied')
    } catch (error) {
      if (error instanceof Error && error.message.toLowerCase().includes('already closed')) {
        errorMessage.value = 'This vessel has already been marked as empty.'
      } else {
        errorMessage.value = error instanceof Error ? error.message : 'Failed to mark vessel as empty'
      }
    } finally {
      saving.value = false
    }
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
