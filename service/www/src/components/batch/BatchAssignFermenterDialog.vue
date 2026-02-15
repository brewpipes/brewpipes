<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 520"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Assign to Fermenter</v-card-title>
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

        <v-form :disabled="saving" @submit.prevent="handleSubmit">
          <v-autocomplete
            v-model="form.vessel_uuid"
            class="mb-2"
            density="comfortable"
            item-title="title"
            item-value="value"
            :items="availableVesselOptions"
            label="Vessel"
            :rules="[rules.required]"
          >
            <template #item="{ props: itemProps, item }">
              <v-list-item v-bind="itemProps">
                <template #subtitle>
                  <span>{{ item.raw.subtitle }}</span>
                </template>
              </v-list-item>
            </template>
            <template #no-data>
              <v-list-item>
                <v-list-item-title>No available vessels</v-list-item-title>
                <v-list-item-subtitle>All vessels are currently occupied or inactive</v-list-item-subtitle>
              </v-list-item>
            </template>
          </v-autocomplete>

          <v-autocomplete
            v-model="form.volume_uuid"
            class="mb-2"
            density="comfortable"
            item-title="title"
            item-value="value"
            :items="volumeOptions"
            label="Volume"
            :rules="[rules.required]"
          >
            <template #no-data>
              <v-list-item>
                <v-list-item-title>No volumes available</v-list-item-title>
                <v-list-item-subtitle>This batch has no tracked volumes</v-list-item-subtitle>
              </v-list-item>
            </template>
          </v-autocomplete>

          <v-text-field
            v-model="form.in_at"
            class="mb-2"
            density="comfortable"
            label="Assigned at"
            type="datetime-local"
          />

          <v-select
            v-model="form.status"
            density="comfortable"
            item-title="title"
            item-value="value"
            :items="statusOptions"
            label="Initial Status"
          />
        </v-form>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSubmit"
        >
          Assign
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Occupancy, OccupancyStatus, Vessel, Volume } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { useVesselTypeFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { nowInputValue } from '@/utils/normalize'

  const props = defineProps<{
    modelValue: boolean
    batchVolumes: Volume[]
    vessels: Vessel[]
    activeOccupancies: Occupancy[]
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'assigned': []
  }>()

  const { createOccupancy } = useProductionApi()
  const { formatVesselType } = useVesselTypeFormatters()

  const saving = ref(false)
  const errorMessage = ref('')

  const form = reactive({
    vessel_uuid: '' as string,
    volume_uuid: '' as string,
    in_at: '',
    status: 'fermenting' as OccupancyStatus,
  })

  const rules = {
    required: (v: string) => !!v || 'Required',
  }

  const isFormValid = computed(() => {
    return !!form.vessel_uuid && !!form.volume_uuid
  })

  // Compute the set of vessel UUIDs that are currently occupied
  const occupiedVesselUuids = computed(() =>
    new Set(props.activeOccupancies.map(o => o.vessel_uuid)),
  )

  // Filter to active vessels that are NOT currently occupied
  const availableVesselOptions = computed(() =>
    props.vessels
      .filter(v => v.status === 'active' && !occupiedVesselUuids.value.has(v.uuid))
      .map(v => ({
        title: v.name,
        value: v.uuid,
        subtitle: `${formatVesselType(v.type)} · ${v.capacity} ${v.capacity_unit}`,
      })),
  )

  const volumeOptions = computed(() =>
    props.batchVolumes.map(v => ({
      title: v.name
        ? `${v.name} (${v.amount} ${v.amount_unit})`
        : `Unnamed Volume (${v.amount} ${v.amount_unit})`,
      value: v.uuid,
    })),
  )

  const statusOptions: { title: string, value: OccupancyStatus }[] = [
    { title: 'Fermenting', value: 'fermenting' },
    { title: 'Conditioning', value: 'conditioning' },
    { title: 'Cold Crashing', value: 'cold_crashing' },
    { title: 'Dry Hopping', value: 'dry_hopping' },
    { title: 'Carbonating', value: 'carbonating' },
    { title: 'Holding', value: 'holding' },
    { title: 'Packaging', value: 'packaging' },
  ]

  // Reset form when dialog opens
  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        errorMessage.value = ''
        form.vessel_uuid = ''
        form.in_at = nowInputValue()
        form.status = 'fermenting'

        // Auto-select volume if there's exactly one
        form.volume_uuid = props.batchVolumes.length === 1 && props.batchVolumes[0] ? props.batchVolumes[0].uuid : ''
      }
    },
  )

  async function handleSubmit () {
    if (!isFormValid.value) return

    saving.value = true
    errorMessage.value = ''

    try {
      const payload = {
        vessel_uuid: form.vessel_uuid,
        volume_uuid: form.volume_uuid,
        in_at: form.in_at ? new Date(form.in_at).toISOString() : undefined,
        status: form.status || undefined,
      }

      await createOccupancy(payload)
      emit('update:modelValue', false)
      emit('assigned')
    } catch (error) {
      errorMessage.value = error instanceof Error ? error.message : 'Failed to assign batch to fermenter'
    } finally {
      saving.value = false
    }
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
