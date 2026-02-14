<template>
  <v-dialog max-width="520" :model-value="modelValue" @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title class="text-h6">Add hot-side addition</v-card-title>
      <v-card-text>
        <v-select
          density="comfortable"
          :items="additionTypeOptions"
          label="Addition type"
          :model-value="form.addition_type"
          @update:model-value="updateForm('addition_type', $event)"
        />
        <v-text-field
          density="comfortable"
          label="Stage"
          :model-value="form.stage"
          placeholder="60 min, 15 min, whirlpool"
          @update:model-value="updateForm('stage', $event)"
        />
        <v-text-field
          density="comfortable"
          label="Inventory lot UUID"
          :model-value="form.inventory_lot_uuid"
          placeholder="Optional"
          @update:model-value="updateForm('inventory_lot_uuid', $event)"
        />
        <v-row>
          <v-col cols="8">
            <v-text-field
              density="comfortable"
              label="Amount"
              :model-value="form.amount"
              type="number"
              @update:model-value="updateForm('amount', $event)"
            />
          </v-col>
          <v-col cols="4">
            <v-select
              density="comfortable"
              :items="volumeUnitOptions"
              label="Unit"
              :model-value="form.amount_unit"
              @update:model-value="updateForm('amount_unit', $event)"
            />
          </v-col>
        </v-row>
        <v-text-field
          density="comfortable"
          label="Added at"
          :model-value="form.added_at"
          type="datetime-local"
          @update:model-value="updateForm('added_at', $event)"
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
          color="primary"
          :disabled="!form.amount"
          :loading="saving"
          @click="emit('submit')"
        >
          Add addition
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { AdditionType as ProductionAdditionType, VolumeUnit } from '@/types'

  export type HotSideAdditionForm = {
    addition_type: ProductionAdditionType
    stage: string
    inventory_lot_uuid: string
    amount: string
    amount_unit: VolumeUnit
    added_at: string
    notes: string
  }

  const props = defineProps<{
    modelValue: boolean
    form: HotSideAdditionForm
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'update:form': [form: HotSideAdditionForm]
    'submit': []
  }>()

  const volumeUnitOptions: VolumeUnit[] = ['ml', 'usfloz', 'ukfloz', 'bbl']
  const additionTypeOptions: ProductionAdditionType[] = [
    'malt',
    'hop',
    'yeast',
    'adjunct',
    'water_chem',
    'gas',
    'other',
  ]

  function updateForm<K extends keyof HotSideAdditionForm> (key: K, value: HotSideAdditionForm[K]) {
    emit('update:form', { ...props.form, [key]: value })
  }
</script>
