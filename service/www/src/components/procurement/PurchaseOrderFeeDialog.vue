<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 480"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit fee' : 'Add fee' }}
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12">
            <v-combobox
              v-model="form.fee_type"
              :items="feeTypeOptions"
              label="Fee type"
              :rules="[v => !!v?.trim() || 'Required']"
            />
          </v-col>
          <v-col cols="12" sm="6">
            <v-text-field
              v-model.number="form.amount_dollars"
              label="Amount"
              prefix="$"
              :rules="[v => v >= 0 || 'Must be non-negative']"
              step="0.01"
              type="number"
            />
          </v-col>
          <v-col cols="12" sm="6">
            <v-combobox
              v-model="form.currency"
              :items="currencyOptions"
              label="Currency"
              :rules="[v => !!v || 'Required']"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isValid"
          :loading="saving"
          @click="handleSave"
        >
          {{ isEditing ? 'Save changes' : 'Add fee' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { PurchaseOrderFee } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { useProcurementApi } from '@/composables/useProcurementApi'

  const props = defineProps<{
    modelValue: boolean
    fee?: PurchaseOrderFee | null
    saving?: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'save': [form: FeeForm]
    'cancel': []
  }>()

  export interface FeeForm {
    fee_type: string
    amount_cents: number | null
    currency: string
  }

  interface InternalForm {
    fee_type: string
    amount_dollars: number | null
    currency: string
  }

  const { dollarsToCents, centsToDollars } = useProcurementApi()

  // v-combobox requires plain string items — {title,value} objects cause the
  // model to be set to the full object on selection, breaking .trim() calls
  // and API payloads. Use raw enum values as suggestions; display formatting
  // is handled by formatFeeType in list/table contexts.
  const feeTypeOptions = ['shipping', 'handling', 'tax', 'insurance', 'customs', 'freight', 'hazmat', 'other']
  const currencyOptions = ['USD', 'CAD', 'EUR', 'GBP']

  const form = reactive<InternalForm>({
    fee_type: '',
    amount_dollars: null,
    currency: 'USD',
  })

  const isEditing = computed(() => !!props.fee)

  const isValid = computed(() => {
    return (
      form.fee_type.trim().length > 0
      && form.amount_dollars !== null
      && form.amount_dollars >= 0
      && form.currency.trim().length > 0
    )
  })

  watch(() => props.modelValue, open => {
    if (open) {
      if (props.fee) {
        form.fee_type = props.fee.fee_type
        form.amount_dollars = centsToDollars(props.fee.amount_cents)
        form.currency = props.fee.currency
      } else {
        resetForm()
      }
    }
  })

  function resetForm () {
    form.fee_type = ''
    form.amount_dollars = null
    form.currency = 'USD'
  }

  function handleSave () {
    emit('save', {
      fee_type: form.fee_type,
      amount_cents: dollarsToCents(form.amount_dollars),
      currency: form.currency,
    })
  }

  function handleCancel () {
    emit('cancel')
    emit('update:modelValue', false)
  }
</script>
