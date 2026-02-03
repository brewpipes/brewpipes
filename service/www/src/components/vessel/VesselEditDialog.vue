<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 640"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">
        Edit vessel
      </v-card-title>
      <v-card-text>
        <v-alert
          v-if="errorMessage"
          class="mb-4"
          closable
          density="compact"
          type="error"
          variant="tonal"
          @click:close="errorMessage = ''"
        >
          {{ errorMessage }}
        </v-alert>

        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="form.name"
              density="comfortable"
              label="Name"
              :rules="[rules.required]"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-select
              v-model="form.type"
              density="comfortable"
              :items="vesselTypeOptions"
              label="Type"
              :rules="[rules.required]"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="form.capacity"
              density="comfortable"
              label="Capacity"
              :rules="[rules.required, rules.positiveNumber]"
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
              :items="vesselStatusOptions"
              label="Status"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="form.make"
              density="comfortable"
              label="Make"
              placeholder="Optional"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="form.model"
              density="comfortable"
              label="Model"
              placeholder="Optional"
            />
          </v-col>
        </v-row>

        <!-- Retirement warning -->
        <v-alert
          v-if="showRetirementWarning"
          class="mt-2"
          density="compact"
          icon="mdi-alert"
          type="warning"
          variant="tonal"
        >
          <div class="font-weight-medium">Retiring this vessel</div>
          <div class="text-body-2">
            Retired vessels cannot be used for new occupancies. This action can be reversed by changing the status back to Active.
          </div>
        </v-alert>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSave"
        >
          Save changes
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { UpdateVesselRequest, Vessel, VesselStatus, VolumeUnit } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { useVesselTypeFormatters } from '@/composables/useFormatters'
  import { VESSEL_STATUS_VALUES, VESSEL_TYPE_VALUES } from '@/composables/useProductionApi'
  import { volumeOptions } from '@/composables/useUnitPreferences'

  const props = defineProps<{
    modelValue: boolean
    vessel: Vessel | null
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'save': [data: UpdateVesselRequest]
  }>()

  const { formatVesselType } = useVesselTypeFormatters()

  // Form state
  const form = reactive({
    name: '',
    type: '',
    capacity: '',
    capacity_unit: 'ml' as VolumeUnit,
    status: 'active' as VesselStatus,
    make: '',
    model: '',
  })

  const saving = ref(false)
  const errorMessage = ref('')

  // Options for dropdowns
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
      title: status.charAt(0).toUpperCase() + status.slice(1),
    })),
  )

  // Validation rules
  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
    positiveNumber: (v: string) => {
      const num = Number(v)
      return (Number.isFinite(num) && num > 0) || 'Must be a positive number'
    },
  }

  // Computed
  const isFormValid = computed(() => {
    return form.name.trim().length > 0
      && form.type.trim().length > 0
      && form.capacity !== ''
      && Number(form.capacity) > 0
  })

  const showRetirementWarning = computed(() => {
    if (!props.vessel) return false
    return props.vessel.status !== 'retired' && form.status === 'retired'
  })

  // Reset form when dialog opens
  watch(
    () => props.modelValue,
    open => {
      if (open && props.vessel) {
        form.name = props.vessel.name
        form.type = props.vessel.type
        form.capacity = String(props.vessel.capacity)
        form.capacity_unit = props.vessel.capacity_unit
        form.status = props.vessel.status
        form.make = props.vessel.make ?? ''
        form.model = props.vessel.model ?? ''
        errorMessage.value = ''
      }
    },
  )

  function handleCancel () {
    emit('update:modelValue', false)
  }

  function handleSave () {
    if (!isFormValid.value) return

    const data: UpdateVesselRequest = {
      name: form.name.trim(),
      type: form.type.trim(),
      capacity: Number(form.capacity),
      capacity_unit: form.capacity_unit,
      status: form.status,
      make: form.make.trim() || null,
      model: form.model.trim() || null,
    }

    emit('save', data)
  }

  // Expose methods for parent to control state
  defineExpose({
    setSaving: (value: boolean) => { saving.value = value },
    setError: (message: string) => { errorMessage.value = message },
    clearError: () => { errorMessage.value = '' },
  })
</script>
