<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 500"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Create receipt</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="form.reference_code"
          density="comfortable"
          label="Reference code"
        />
        <v-text-field
          v-model="form.supplier_uuid"
          density="comfortable"
          label="Supplier UUID"
        />
        <v-text-field
          v-model="form.received_at"
          density="comfortable"
          label="Received at"
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
          Create receipt
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateInventoryReceiptRequest } from '@/types'
  import { reactive, watch } from 'vue'
  import { normalizeDateTime, normalizeText } from '@/utils/normalize'

  const props = defineProps<{
    modelValue: boolean
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateInventoryReceiptRequest]
  }>()

  const form = reactive({
    reference_code: '',
    supplier_uuid: '',
    received_at: '',
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
    form.reference_code = ''
    form.supplier_uuid = ''
    form.received_at = ''
    form.notes = ''
  }

  function handleSubmit () {
    const payload: CreateInventoryReceiptRequest = {
      supplier_uuid: normalizeText(form.supplier_uuid),
      reference_code: normalizeText(form.reference_code),
      received_at: normalizeDateTime(form.received_at),
      notes: normalizeText(form.notes),
    }
    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
