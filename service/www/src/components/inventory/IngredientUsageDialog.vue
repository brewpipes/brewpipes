<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 500"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Log ingredient usage</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="form.production_ref_uuid"
          density="comfortable"
          label="Batch reference UUID"
        />
        <v-text-field
          v-model="form.used_at"
          density="comfortable"
          label="Used at"
          type="datetime-local"
        />
        <v-textarea
          v-model="form.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="saving"
          @click="handleSubmit"
        >
          Log usage
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateInventoryUsageRequest } from '@/types'
  import { reactive, watch } from 'vue'
  import { normalizeDateTime, normalizeText } from '@/utils/normalize'

  const props = defineProps<{
    modelValue: boolean
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateInventoryUsageRequest]
  }>()

  const form = reactive({
    production_ref_uuid: '',
    used_at: '',
    notes: '',
  })

  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        resetForm()
      }
    },
  )

  function resetForm () {
    form.production_ref_uuid = ''
    form.used_at = ''
    form.notes = ''
  }

  function handleSubmit () {
    const payload: CreateInventoryUsageRequest = {
      production_ref_uuid: normalizeText(form.production_ref_uuid),
      used_at: normalizeDateTime(form.used_at),
      notes: normalizeText(form.notes),
    }
    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
