<template>
  <v-dialog max-width="520" :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title class="text-h6">Record measurement</v-card-title>
      <v-card-text>
        <v-select
          :items="targetOptions"
          label="Target"
          :model-value="form.target"
          @update:model-value="updateForm('target', $event)"
        />
        <v-text-field
          v-if="form.target === 'occupancy'"
          label="Occupancy"
          :model-value="form.occupancy_uuid"
          @update:model-value="updateForm('occupancy_uuid', $event)"
        />
        <v-text-field
          label="Kind"
          :model-value="form.kind"
          placeholder="gravity"
          @update:model-value="updateForm('kind', $event)"
        />
        <v-text-field
          label="Value"
          :model-value="form.value"
          type="number"
          @update:model-value="updateForm('value', $event)"
        />
        <v-text-field
          label="Unit"
          :model-value="form.unit"
          @update:model-value="updateForm('unit', $event)"
        />
        <v-text-field
          label="Observed at"
          :model-value="form.observed_at"
          type="datetime-local"
          @update:model-value="updateForm('observed_at', $event)"
        />
        <v-textarea
          auto-grow
          label="Notes"
          :model-value="form.notes"
          rows="2"
          @update:model-value="updateForm('notes', $event)"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="emit('update:modelValue', false)">Cancel</v-btn>
        <v-btn
          color="secondary"
          :disabled="!isValid"
          @click="emit('submit')"
        >
          Add measurement
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'

  export type MeasurementForm = {
    target: 'batch' | 'occupancy'
    occupancy_uuid: string
    kind: string
    value: string
    unit: string
    observed_at: string
    notes: string
  }

  const props = defineProps<{
    modelValue: boolean
    form: MeasurementForm
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'update:form': [form: MeasurementForm]
    'submit': []
  }>()

  const targetOptions = [
    { title: 'Batch', value: 'batch' },
    { title: 'Occupancy', value: 'occupancy' },
  ]

  const isValid = computed(() => {
    if (!props.form.kind.trim()) return false
    if (!props.form.value) return false
    if (props.form.target === 'occupancy' && !props.form.occupancy_uuid) return false
    return true
  })

  function updateForm<K extends keyof MeasurementForm> (key: K, value: MeasurementForm[K]) {
    emit('update:form', { ...props.form, [key]: value })
  }
</script>
