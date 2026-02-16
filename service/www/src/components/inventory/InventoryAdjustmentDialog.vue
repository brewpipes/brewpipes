<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? undefined : 500"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Adjust inventory</v-card-title>
      <v-card-text>
        <div v-if="lot" class="mb-4 pa-3 selected-item-summary rounded">
          <div class="text-subtitle-2 font-weight-bold">{{ lot.name }}</div>
          <div class="text-caption text-medium-emphasis">
            {{ lot.type === 'ingredient' ? 'Ingredient lot' : 'Beer lot' }}
            <span class="mx-1">|</span>
            {{ lot.locationName }}
          </div>
          <div class="text-body-2 mt-1">
            Current quantity: <strong>{{ formatAmountPreferred(lot.quantity, lot.unit) }}</strong>
          </div>
        </div>

        <v-text-field
          v-model="form.amount"
          density="comfortable"
          hint="Use negative values to decrease inventory"
          label="Adjustment amount"
          persistent-hint
          type="number"
        />
        <v-select
          v-model="form.reason"
          class="mt-2"
          density="comfortable"
          :items="adjustmentReasonOptions"
          label="Reason"
          :rules="[rules.required]"
        />
        <v-textarea
          v-model="form.notes"
          auto-grow
          class="mt-2"
          density="comfortable"
          label="Notes (optional)"
          rows="2"
        />
        <v-text-field
          v-model="form.adjusted_at"
          class="mt-2"
          density="comfortable"
          label="Adjusted at"
          type="datetime-local"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSubmit"
        >
          Save Adjustment
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateInventoryAdjustmentRequest } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { normalizeDateTime, normalizeText, nowInputValue } from '@/utils/normalize'

  export interface InventoryLotInfo {
    type: 'ingredient' | 'beer'
    lotUuid: string
    name: string
    quantity: number
    unit: string
    locationUuid: string
    locationName: string
  }

  const adjustmentReasonOptions = [
    { title: 'Cycle Count', value: 'cycle_count' },
    { title: 'Spoilage', value: 'spoilage' },
    { title: 'Shrink', value: 'shrink' },
    { title: 'Damage', value: 'damage' },
    { title: 'Correction', value: 'correction' },
    { title: 'Other', value: 'other' },
  ]

  const props = defineProps<{
    modelValue: boolean
    lot: InventoryLotInfo | null
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateInventoryAdjustmentRequest]
  }>()

  const { formatAmountPreferred } = useUnitPreferences()

  const form = reactive({
    amount: '',
    reason: '',
    notes: '',
    adjusted_at: '',
  })

  const rules = {
    required: (v: string | null) => (v !== null && v !== '' && String(v).trim() !== '') || 'Required',
  }

  const isFormValid = computed(() => {
    return form.amount !== '' && form.reason.trim().length > 0
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
    form.amount = ''
    form.reason = ''
    form.notes = ''
    form.adjusted_at = nowInputValue()
  }

  function handleSubmit () {
    if (!isFormValid.value || !props.lot) return

    const payload: CreateInventoryAdjustmentRequest = {
      ingredient_lot_uuid: props.lot.type === 'ingredient' ? props.lot.lotUuid : null,
      beer_lot_uuid: props.lot.type === 'beer' ? props.lot.lotUuid : null,
      stock_location_uuid: props.lot.locationUuid,
      amount: Number(form.amount),
      amount_unit: props.lot.unit,
      reason: form.reason.trim(),
      notes: normalizeText(form.notes),
      adjusted_at: normalizeDateTime(form.adjusted_at),
    }
    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>

<style scoped>
.selected-item-summary {
  background: rgba(var(--v-theme-on-surface), 0.05);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
}
</style>
