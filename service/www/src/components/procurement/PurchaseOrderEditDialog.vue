<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 640"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Edit purchase order</v-card-title>
      <v-card-text>
        <v-alert
          v-if="errorMessage"
          class="mb-3"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ errorMessage }}
        </v-alert>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.order_number" label="Order number" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.ordered_at" label="Ordered at" type="datetime-local" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.expected_at" label="Expected at" type="datetime-local" />
          </v-col>
          <v-col cols="12">
            <v-textarea v-model="form.notes" auto-grow label="Notes" rows="3" />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!form.order_number.trim()"
          :loading="saving"
          @click="handleSubmit"
        >
          Save changes
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { PurchaseOrder } from '@/types'
  import { reactive, watch } from 'vue'
  import { normalizeDateTime, normalizeText, toLocalDateTimeInput } from '@/utils/normalize'

  export interface PurchaseOrderEditFormData {
    order_number: string
    ordered_at: string | null
    expected_at: string | null
    notes: string | null
  }

  const props = defineProps<{
    modelValue: boolean
    purchaseOrder: PurchaseOrder
    saving: boolean
    errorMessage?: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: PurchaseOrderEditFormData]
  }>()

  const form = reactive({
    order_number: '',
    ordered_at: '',
    expected_at: '',
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
    form.order_number = props.purchaseOrder.order_number
    form.ordered_at = props.purchaseOrder.ordered_at ? toLocalDateTimeInput(props.purchaseOrder.ordered_at) : ''
    form.expected_at = props.purchaseOrder.expected_at ? toLocalDateTimeInput(props.purchaseOrder.expected_at) : ''
    form.notes = props.purchaseOrder.notes ?? ''
  }

  function handleSubmit () {
    if (!form.order_number.trim()) return

    const payload: PurchaseOrderEditFormData = {
      order_number: form.order_number.trim(),
      ordered_at: normalizeDateTime(form.ordered_at),
      expected_at: normalizeDateTime(form.expected_at),
      notes: normalizeText(form.notes),
    }
    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
