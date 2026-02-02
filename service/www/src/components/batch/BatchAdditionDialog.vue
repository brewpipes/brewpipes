<template>
  <v-dialog max-width="520" :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title class="text-h6">Record addition</v-card-title>
      <v-card-text>
        <v-select
          :items="targetOptions"
          label="Target"
          :model-value="form.target"
          @update:model-value="updateForm('target', $event)"
        />
        <v-text-field
          v-if="form.target === 'occupancy'"
          label="Occupancy ID"
          :model-value="form.occupancy_id"
          type="number"
          @update:model-value="updateForm('occupancy_id', $event)"
        />
        <v-select
          :items="additionTypeOptions"
          label="Addition type"
          :model-value="form.addition_type"
          @update:model-value="updateForm('addition_type', $event)"
        />
        <v-text-field
          label="Stage"
          :model-value="form.stage"
          @update:model-value="updateForm('stage', $event)"
        />
        <v-text-field
          label="Inventory lot UUID"
          :model-value="form.inventory_lot_uuid"
          @update:model-value="updateForm('inventory_lot_uuid', $event)"
        />
        <v-text-field
          label="Amount"
          :model-value="form.amount"
          type="number"
          @update:model-value="updateForm('amount', $event)"
        />
        <v-select
          :items="unitOptions"
          label="Unit"
          :model-value="form.amount_unit"
          @update:model-value="updateForm('amount_unit', $event)"
        />
        <v-text-field
          label="Added at"
          :model-value="form.added_at"
          type="datetime-local"
          @update:model-value="updateForm('added_at', $event)"
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
          color="primary"
          :disabled="!isValid"
          @click="emit('submit')"
        >
          Add addition
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { AdditionType, Unit } from './types'
  import { computed } from 'vue'

  export type AdditionForm = {
    target: 'batch' | 'occupancy'
    occupancy_id: string
    addition_type: AdditionType
    stage: string
    inventory_lot_uuid: string
    amount: string
    amount_unit: Unit
    added_at: string
    notes: string
  }

  const props = defineProps<{
    modelValue: boolean
    form: AdditionForm
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'update:form': [form: AdditionForm]
    'submit': []
  }>()

  const unitOptions: Unit[] = ['ml', 'usfloz', 'ukfloz', 'bbl']
  const additionTypeOptions: AdditionType[] = [
    'malt',
    'hop',
    'yeast',
    'adjunct',
    'water_chem',
    'gas',
    'other',
  ]
  const targetOptions = [
    { title: 'Batch', value: 'batch' },
    { title: 'Occupancy', value: 'occupancy' },
  ]

  const isValid = computed(() => {
    if (!props.form.amount) return false
    if (props.form.target === 'occupancy' && !props.form.occupancy_id) return false
    return true
  })

  function updateForm<K extends keyof AdditionForm> (key: K, value: AdditionForm[K]) {
    emit('update:form', { ...props.form, [key]: value })
  }
</script>
