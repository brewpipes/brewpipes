<template>
  <v-dialog max-width="720" :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title class="text-h6">Extended reading</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="Observed at"
              :model-value="form.observed_at"
              type="datetime-local"
              @update:model-value="updateForm('observed_at', $event)"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="Temperature"
              :model-value="form.temperature"
              type="number"
              @update:model-value="updateForm('temperature', $event)"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="Temp unit"
              :model-value="form.temperature_unit"
              :placeholder="temperatureUnitPlaceholder"
              @update:model-value="updateForm('temperature_unit', $event)"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="Gravity"
              :model-value="form.gravity"
              type="number"
              @update:model-value="updateForm('gravity', $event)"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="Gravity unit"
              :model-value="form.gravity_unit"
              :placeholder="gravityUnitPlaceholder"
              @update:model-value="updateForm('gravity_unit', $event)"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="pH"
              :model-value="form.ph"
              type="number"
              @update:model-value="updateForm('ph', $event)"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="pH unit"
              :model-value="form.ph_unit"
              placeholder="pH"
              @update:model-value="updateForm('ph_unit', $event)"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="Other kind"
              :model-value="form.extra_kind"
              placeholder="CO2"
              @update:model-value="updateForm('extra_kind', $event)"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="Other value"
              :model-value="form.extra_value"
              type="number"
              @update:model-value="updateForm('extra_value', $event)"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              density="comfortable"
              label="Other unit"
              :model-value="form.extra_unit"
              @update:model-value="updateForm('extra_unit', $event)"
            />
          </v-col>
          <v-col cols="12" md="8">
            <v-text-field
              density="comfortable"
              label="Notes"
              :model-value="form.notes"
              placeholder="Aroma, flavor, observations"
              @update:model-value="updateForm('notes', $event)"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="emit('update:modelValue', false)">Cancel</v-btn>
        <v-btn color="primary" :disabled="!isReady" @click="emit('submit')">
          Record
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'

  export type TimelineExtendedForm = {
    observed_at: string
    temperature: string
    temperature_unit: string
    gravity: string
    gravity_unit: string
    ph: string
    ph_unit: string
    extra_kind: string
    extra_value: string
    extra_unit: string
    notes: string
  }

  const props = defineProps<{
    modelValue: boolean
    form: TimelineExtendedForm
    temperatureUnit: 'f' | 'c'
    gravityUnit: 'sg' | 'plato'
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'update:form': [form: TimelineExtendedForm]
    'submit': []
  }>()

  const temperatureUnitPlaceholder = computed(() =>
    props.temperatureUnit === 'f' ? 'F' : 'C',
  )

  const gravityUnitPlaceholder = computed(() =>
    props.gravityUnit === 'sg' ? 'SG' : 'Plato',
  )

  const isReady = computed(() => {
    const hasTemperature = toNumber(props.form.temperature) !== null
    const hasGravity = toNumber(props.form.gravity) !== null
    const hasPh = toNumber(props.form.ph) !== null
    const hasNotes = props.form.notes.trim().length > 0
    const hasExtraKind = props.form.extra_kind.trim().length > 0
    const extraValue = toNumber(props.form.extra_value)
    if (hasExtraKind && extraValue === null) {
      return false
    }
    const hasExtra = hasExtraKind && extraValue !== null
    return hasTemperature || hasGravity || hasPh || hasExtra || hasNotes
  })

  function toNumber (value: string | number | null) {
    if (value === null || value === undefined || value === '') {
      return null
    }
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : null
  }

  function updateForm<K extends keyof TimelineExtendedForm> (key: K, value: TimelineExtendedForm[K]) {
    emit('update:form', { ...props.form, [key]: value })
  }
</script>
