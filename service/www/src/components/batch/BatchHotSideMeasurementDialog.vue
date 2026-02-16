<template>
  <v-dialog max-width="520" :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title class="text-h6">Add hot-side measurement</v-card-title>
      <v-card-text>
        <v-select
          density="comfortable"
          :items="hotSideMeasurementKinds"
          label="Kind"
          :model-value="form.kind"
          @update:model-value="updateForm('kind', $event)"
        />
        <v-row>
          <v-col cols="8">
            <v-text-field
              density="comfortable"
              label="Value"
              :model-value="form.value"
              type="number"
              @update:model-value="updateForm('value', $event)"
            />
          </v-col>
          <v-col cols="4">
            <v-text-field
              density="comfortable"
              label="Unit"
              :model-value="form.unit"
              :placeholder="getDefaultUnitForKind(form.kind)"
              @update:model-value="updateForm('unit', $event)"
            />
          </v-col>
        </v-row>
        <v-text-field
          density="comfortable"
          label="Observed at"
          :model-value="form.observed_at"
          type="datetime-local"
          @update:model-value="updateForm('observed_at', $event)"
        />
        <v-textarea
          auto-grow
          density="comfortable"
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
          :disabled="!form.kind || !form.value"
          :loading="saving"
          @click="emit('submit')"
        >
          Add measurement
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { getDefaultUnitForKind, hotSideMeasurementKinds } from './types'

  export type HotSideMeasurementForm = {
    kind: string
    value: string
    unit: string
    observed_at: string
    notes: string
  }

  const props = defineProps<{
    modelValue: boolean
    form: HotSideMeasurementForm
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'update:form': [form: HotSideMeasurementForm]
    'submit': []
  }>()

  function updateForm<K extends keyof HotSideMeasurementForm> (key: K, value: HotSideMeasurementForm[K]) {
    emit('update:form', { ...props.form, [key]: value })
  }
</script>
