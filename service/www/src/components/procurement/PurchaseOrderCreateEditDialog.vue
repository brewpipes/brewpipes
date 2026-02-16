<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 640"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit purchase order' : 'Create purchase order' }}
      </v-card-title>
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
          <v-col cols="12">
            <v-select
              v-model="form.supplier_uuid"
              :disabled="isEditing"
              :hint="isEditing ? 'Supplier cannot be changed after creation' : ''"
              :items="suppliers"
              label="Supplier"
              :persistent-hint="isEditing"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.order_number" label="Order number" />
          </v-col>
          <v-col cols="12" md="6">
            <v-select v-model="form.status" clearable :items="statusOptions" label="Status" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.ordered_at" label="Ordered at" type="datetime-local" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="form.expected_at" label="Expected at" type="datetime-local" />
          </v-col>
          <v-col cols="12">
            <v-textarea v-model="form.notes" auto-grow label="Notes" rows="2" />
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
          {{ isEditing ? 'Save changes' : 'Add purchase order' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { PurchaseOrder } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { usePurchaseOrderStatusFormatters } from '@/composables/useFormatters'
  import { normalizeDateTime, normalizeText, toLocalDateTimeInput } from '@/utils/normalize'

  export interface PurchaseOrderCreateSubmitData {
    supplier_uuid: string
    order_number: string
    status: string | null
    ordered_at: string | null
    expected_at: string | null
    notes: string | null
  }

  export interface PurchaseOrderEditSubmitData {
    order_number: string
    status: string | undefined
    ordered_at: string | null
    expected_at: string | null
    notes: string | null
  }

  interface SupplierSelectItem {
    title: string
    value: string
  }

  const props = defineProps<{
    modelValue: boolean
    editOrder?: PurchaseOrder | null
    suppliers: SupplierSelectItem[]
    saving: boolean
    errorMessage?: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: PurchaseOrderCreateSubmitData | PurchaseOrderEditSubmitData]
  }>()

  const { purchaseOrderStatusOptions } = usePurchaseOrderStatusFormatters()

  const form = reactive({
    supplier_uuid: null as string | null,
    order_number: '',
    status: '',
    ordered_at: '',
    expected_at: '',
    notes: '',
  })

  const statusOptions = computed(() => purchaseOrderStatusOptions)

  const isEditing = computed(() => !!props.editOrder)

  const isFormValid = computed(() => {
    return form.supplier_uuid !== null && form.order_number.trim().length > 0
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
    if (props.editOrder) {
      form.supplier_uuid = props.editOrder.supplier_uuid
      form.order_number = props.editOrder.order_number
      form.status = props.editOrder.status || ''
      form.ordered_at = props.editOrder.ordered_at ? toLocalDateTimeInput(props.editOrder.ordered_at) : ''
      form.expected_at = props.editOrder.expected_at ? toLocalDateTimeInput(props.editOrder.expected_at) : ''
      form.notes = props.editOrder.notes || ''
    } else {
      form.supplier_uuid = null
      form.order_number = ''
      form.status = ''
      form.ordered_at = ''
      form.expected_at = ''
      form.notes = ''
    }
  }

  function handleSubmit () {
    if (!isFormValid.value) return

    if (isEditing.value) {
      const payload: PurchaseOrderEditSubmitData = {
        order_number: form.order_number.trim(),
        status: normalizeText(form.status) ?? undefined,
        ordered_at: normalizeDateTime(form.ordered_at),
        expected_at: normalizeDateTime(form.expected_at),
        notes: normalizeText(form.notes),
      }
      emit('submit', payload)
    } else {
      const payload: PurchaseOrderCreateSubmitData = {
        supplier_uuid: form.supplier_uuid!,
        order_number: form.order_number.trim(),
        status: normalizeText(form.status),
        ordered_at: normalizeDateTime(form.ordered_at),
        expected_at: normalizeDateTime(form.expected_at),
        notes: normalizeText(form.notes),
      }
      emit('submit', payload)
    }
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
