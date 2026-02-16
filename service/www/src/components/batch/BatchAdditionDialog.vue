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
          label="Occupancy"
          :model-value="form.occupancy_uuid"
          @update:model-value="updateForm('occupancy_uuid', $event)"
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
        <v-autocomplete
          clearable
          item-title="title"
          item-value="value"
          :items="lotSelectItems"
          label="Inventory lot"
          :loading="ingredientLotsLoading"
          :model-value="form.inventory_lot_uuid || null"
          placeholder="Search by ingredient or lot code"
          @update:model-value="updateForm('inventory_lot_uuid', $event ?? '')"
        >
          <template #item="{ props: itemProps, item }">
            <v-list-item v-bind="itemProps">
              <template #subtitle>
                <span>{{ item.raw.subtitle }}</span>
              </template>
            </v-list-item>
          </template>
          <template #no-data>
            <v-list-item>
              <v-list-item-title>No lots available</v-list-item-title>
              <v-list-item-subtitle>No ingredient lots found in inventory</v-list-item-subtitle>
            </v-list-item>
          </template>
        </v-autocomplete>
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
  import type { Ingredient, IngredientLot } from '@/types'
  import type { AdditionType, Unit } from './types'
  import { computed } from 'vue'

  export type AdditionForm = {
    target: 'batch' | 'occupancy'
    occupancy_uuid: string
    addition_type: AdditionType
    stage: string
    inventory_lot_uuid: string
    amount: string
    amount_unit: Unit
    added_at: string
    notes: string
  }

  const props = withDefaults(defineProps<{
    modelValue: boolean
    form: AdditionForm
    ingredientLots?: IngredientLot[]
    ingredients?: Ingredient[]
    ingredientLotsLoading?: boolean
  }>(), {
    ingredientLots: () => [],
    ingredients: () => [],
    ingredientLotsLoading: false,
  })

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

  const ingredientMap = computed(() => {
    const map = new Map<string, Ingredient>()
    for (const ingredient of props.ingredients) {
      map.set(ingredient.uuid, ingredient)
    }
    return map
  })

  const lotSelectItems = computed(() =>
    props.ingredientLots.map(lot => {
      const ingredient = ingredientMap.value.get(lot.ingredient_uuid)
      const ingredientName = ingredient?.name ?? 'Unknown ingredient'
      const lotCode = lot.brewery_lot_code || lot.originator_lot_code || lot.uuid.slice(0, 8)
      const stock = `${lot.current_amount} ${lot.current_unit}`
      return {
        title: `${ingredientName} — Lot #${lotCode}`,
        value: lot.uuid,
        subtitle: `${stock} available`,
      }
    }),
  )

  const isValid = computed(() => {
    if (!props.form.amount) return false
    if (props.form.target === 'occupancy' && !props.form.occupancy_uuid) return false
    return true
  })

  function updateForm<K extends keyof AdditionForm> (key: K, value: AdditionForm[K]) {
    emit('update:form', { ...props.form, [key]: value })
  }
</script>
