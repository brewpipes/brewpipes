<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? undefined : 500"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Transfer inventory</v-card-title>
      <v-card-text>
        <div v-if="lot" class="mb-4 pa-3 selected-item-summary rounded">
          <div class="text-subtitle-2 font-weight-bold">{{ lot.name }}</div>
          <div class="text-caption text-medium-emphasis">
            {{ lot.type === 'ingredient' ? 'Ingredient lot' : 'Beer lot' }}
          </div>
          <div class="text-body-2 mt-1">
            Available: <strong>{{ formatAmountPreferred(lot.quantity, lot.unit) }}</strong>
          </div>
        </div>

        <v-text-field
          v-model="form.from_location"
          density="comfortable"
          disabled
          label="From location"
        />
        <v-select
          v-model="form.to_location_uuid"
          class="mt-2"
          density="comfortable"
          :items="destinationItems"
          label="To location"
          :rules="[rules.required]"
        />
        <v-text-field
          v-model="form.quantity"
          class="mt-2"
          density="comfortable"
          :hint="lot ? `Max: ${formatAmountPreferred(lot.quantity, lot.unit)}` : ''"
          label="Quantity to transfer"
          persistent-hint
          :rules="[rules.required, rules.positiveNumber]"
          type="number"
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
          v-model="form.transferred_at"
          class="mt-2"
          density="comfortable"
          label="Transferred at"
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
          Transfer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateInventoryTransferRequest, StockLocation } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { convertMass, convertVolume } from '@/composables/useUnitConversion'
  import { isMassUnit, isVolumeUnit, normalizeMassUnit, normalizeVolumeUnit, useUnitPreferences } from '@/composables/useUnitPreferences'
  import { normalizeDateTime, normalizeText, nowInputValue } from '@/utils/normalize'

  import type { InventoryLotInfo } from './InventoryAdjustmentDialog.vue'

  const props = defineProps<{
    modelValue: boolean
    lot: InventoryLotInfo | null
    locations: StockLocation[]
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateInventoryTransferRequest]
  }>()

  const { formatAmountPreferred, preferences } = useUnitPreferences()

  const form = reactive({
    from_location: '',
    to_location_uuid: null as string | null,
    quantity: '',
    notes: '',
    transferred_at: '',
  })

  const rules = {
    required: (v: string | null) => (v !== null && v !== '' && String(v).trim() !== '') || 'Required',
    positiveNumber: (v: string | null) => {
      if (v === null || v === '') return true // Let required handle empty
      const num = Number(v)
      return (Number.isFinite(num) && num > 0) || 'Must be a positive number'
    },
  }

  /**
   * Convert a value from the lot's backend unit to the user's preferred display unit.
   * Returns null if the unit type is unrecognized.
   */
  function toPreferredUnit (value: number, backendUnit: string): { value: number, unit: string } | null {
    const lower = backendUnit.toLowerCase()
    if (isMassUnit(lower)) {
      const from = normalizeMassUnit(lower)
      const to = preferences.value.mass
      const converted = convertMass(value, from, to)
      return converted !== null ? { value: converted, unit: to } : null
    }
    if (isVolumeUnit(lower)) {
      const from = normalizeVolumeUnit(lower)
      const to = preferences.value.volume
      const converted = convertVolume(value, from, to)
      return converted !== null ? { value: converted, unit: to } : null
    }
    return null
  }

  /**
   * Convert a value from the user's preferred display unit back to the lot's backend unit.
   * Returns the value unchanged if the unit type is unrecognized.
   */
  function toBackendUnit (value: number, backendUnit: string): number {
    const lower = backendUnit.toLowerCase()
    if (isMassUnit(lower)) {
      const from = preferences.value.mass
      const to = normalizeMassUnit(lower)
      return convertMass(value, from, to) ?? value
    }
    if (isVolumeUnit(lower)) {
      const from = preferences.value.volume
      const to = normalizeVolumeUnit(lower)
      return convertVolume(value, from, to) ?? value
    }
    return value
  }

  /** The lot's available quantity converted to the user's preferred unit. */
  const maxQuantityPreferred = computed(() => {
    if (!props.lot) return null
    return toPreferredUnit(props.lot.quantity, props.lot.unit)
  })

  const destinationItems = computed(() => {
    if (!props.lot) {
      return props.locations.map(loc => ({ title: loc.name, value: loc.uuid }))
    }
    return props.locations
      .filter(loc => loc.uuid !== props.lot?.locationUuid)
      .map(loc => ({ title: loc.name, value: loc.uuid }))
  })

  const isFormValid = computed(() => {
    if (form.to_location_uuid === null || form.quantity === '') {
      return false
    }
    const qty = Number(form.quantity)
    if (!Number.isFinite(qty) || qty <= 0) {
      return false
    }
    // Ensure transfer quantity doesn't exceed available (compared in preferred units)
    if (maxQuantityPreferred.value && qty > maxQuantityPreferred.value.value) {
      return false
    }
    return true
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
    form.from_location = props.lot?.locationName ?? ''
    form.to_location_uuid = null
    form.quantity = ''
    form.notes = ''
    form.transferred_at = nowInputValue()
  }

  function handleSubmit () {
    if (!isFormValid.value || !props.lot) return

    // Convert user input (in preferred display unit) back to backend unit
    const amountInBackendUnit = toBackendUnit(Number(form.quantity), props.lot.unit)

    const payload: CreateInventoryTransferRequest = {
      ingredient_lot_uuid: props.lot.type === 'ingredient' ? props.lot.lotUuid : null,
      beer_lot_uuid: props.lot.type === 'beer' ? props.lot.lotUuid : null,
      source_location_uuid: props.lot.locationUuid,
      dest_location_uuid: form.to_location_uuid,
      amount: amountInBackendUnit,
      amount_unit: props.lot.unit,
      notes: normalizeText(form.notes),
      transferred_at: normalizeDateTime(form.transferred_at),
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
