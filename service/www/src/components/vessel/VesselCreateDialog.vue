<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 640"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Register vessel</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-select
              v-model="form.type"
              density="comfortable"
              item-title="title"
              item-value="value"
              :items="vesselTypeOptions"
              label="Type"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="form.name"
              density="comfortable"
              label="Name"
              placeholder="FV-01"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="form.capacity"
              density="comfortable"
              label="Capacity"
              :min="0"
              type="number"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-select
              v-model="form.capacity_unit"
              density="comfortable"
              :items="unitOptions"
              label="Capacity unit"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-select
              v-model="form.status"
              density="comfortable"
              item-title="title"
              item-value="value"
              :items="vesselStatusOptions"
              label="Status"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="form.make"
              density="comfortable"
              label="Make"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="form.model"
              density="comfortable"
              label="Model"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSubmit"
        >
          Add vessel
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateVesselRequest, VesselStatus, VesselType, VolumeUnit } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import {
    useVesselStatusFormatters,
    useVesselTypeFormatters,
  } from '@/composables/useFormatters'
  import { useUnitPreferences, volumeOptions } from '@/composables/useUnitPreferences'
  import { VESSEL_STATUS_VALUES, VESSEL_TYPE_VALUES } from '@/types'
  import { normalizeText, toNumber } from '@/utils/normalize'

  const props = defineProps<{
    modelValue: boolean
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateVesselRequest]
  }>()

  const { preferences } = useUnitPreferences()
  const { formatVesselStatus } = useVesselStatusFormatters()
  const { formatVesselType } = useVesselTypeFormatters()

  const unitOptions = volumeOptions.map(opt => opt.value)

  const vesselTypeOptions = computed(() =>
    VESSEL_TYPE_VALUES.map(type => ({
      value: type,
      title: formatVesselType(type),
    })),
  )

  const vesselStatusOptions = computed(() =>
    VESSEL_STATUS_VALUES.map(status => ({
      value: status,
      title: formatVesselStatus(status),
    })),
  )

  const form = reactive({
    type: '' as VesselType | '',
    name: '',
    capacity: '',
    capacity_unit: preferences.value.volume as VolumeUnit,
    status: 'active' as VesselStatus,
    make: '',
    model: '',
  })

  const isFormValid = computed(() => {
    const capacity = Number(form.capacity)
    return form.type.trim().length > 0
      && form.name.trim().length > 0
      && form.capacity !== ''
      && Number.isFinite(capacity)
      && capacity > 0
  })

  // Reset form when dialog opens
  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        form.type = '' as VesselType | ''
        form.name = ''
        form.capacity = ''
        form.capacity_unit = preferences.value.volume
        form.status = 'active' as VesselStatus
        form.make = ''
        form.model = ''
      }
    },
  )

  function handleSubmit () {
    if (!isFormValid.value) return

    const payload: CreateVesselRequest = {
      type: form.type as VesselType,
      name: form.name.trim(),
      capacity: toNumber(form.capacity) ?? 0,
      capacity_unit: form.capacity_unit,
      status: form.status,
      make: normalizeText(form.make),
      model: normalizeText(form.model),
    }

    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
