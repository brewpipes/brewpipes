<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 600"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit line item' : 'Add line item' }}
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" sm="6">
            <v-text-field
              v-model.number="form.line_number"
              label="Line number"
              :rules="[v => !!v || 'Required']"
              type="number"
            />
          </v-col>
          <v-col cols="12" sm="6">
            <v-select
              v-model="form.item_type"
              :items="itemTypeOptions"
              label="Item type"
              :rules="[v => !!v || 'Required']"
            />
          </v-col>
          <v-col cols="12">
            <v-text-field
              v-model="form.item_name"
              label="Item name"
              :rules="[v => !!v?.trim() || 'Required']"
            />
          </v-col>
          <v-col cols="12">
            <v-autocomplete
              v-model="form.inventory_item_uuid"
              clearable
              hint="Optional - link to inventory ingredient"
              item-title="title"
              item-value="value"
              :items="ingredientSelectItems"
              label="Inventory ingredient"
              :loading="ingredientsLoading"
              persistent-hint
            >
              <template #item="{ props: itemProps, item }">
                <v-list-item v-bind="itemProps">
                  <template #subtitle>
                    <span>{{ item.raw.category }}</span>
                  </template>
                </v-list-item>
              </template>
              <template #no-data>
                <v-list-item>
                  <v-list-item-title>No ingredients found</v-list-item-title>
                  <v-list-item-subtitle>Create ingredients in the Inventory section first</v-list-item-subtitle>
                </v-list-item>
              </template>
            </v-autocomplete>
          </v-col>
          <v-col cols="6" sm="4">
            <v-text-field
              v-model.number="form.quantity"
              label="Quantity"
              :rules="[v => v > 0 || 'Must be positive']"
              type="number"
            />
          </v-col>
          <v-col cols="6" sm="4">
            <v-combobox
              v-model="form.quantity_unit"
              :items="unitOptions"
              label="Unit"
              :rules="[v => !!v || 'Required']"
            />
          </v-col>
          <v-col cols="6" sm="4">
            <v-text-field
              v-model.number="form.unit_cost_dollars"
              label="Unit cost"
              prefix="$"
              :rules="[v => v >= 0 || 'Must be non-negative']"
              step="0.01"
              type="number"
            />
          </v-col>
          <v-col cols="6" sm="4">
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
          {{ isEditing ? 'Save changes' : 'Add line' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Ingredient, PurchaseOrderLine } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { useLineItemTypeFormatters } from '@/composables/useFormatters'
  import { useProcurementApi } from '@/composables/useProcurementApi'

  const props = defineProps<{
    modelValue: boolean
    line?: PurchaseOrderLine | null
    saving?: boolean
    ingredients: Ingredient[]
    ingredientsLoading?: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'save': [form: LineForm]
    'cancel': []
  }>()

  export interface LineForm {
    line_number: number | null
    item_type: string
    item_name: string
    inventory_item_uuid: string | null
    quantity: number | null
    quantity_unit: string
    unit_cost_cents: number | null
    currency: string
  }

  interface InternalForm {
    line_number: number | null
    item_type: string
    item_name: string
    inventory_item_uuid: string | null
    quantity: number | null
    quantity_unit: string
    unit_cost_dollars: number | null
    currency: string
  }

  const { lineItemTypeOptions: itemTypeOptions } = useLineItemTypeFormatters()
  const { dollarsToCents, centsToDollars } = useProcurementApi()
  const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']
  const currencyOptions = ['USD', 'CAD', 'EUR', 'GBP']

  const ingredientSelectItems = computed(() =>
    props.ingredients.map(ingredient => ({
      title: ingredient.name,
      value: ingredient.uuid,
      category: ingredient.category.charAt(0).toUpperCase() + ingredient.category.slice(1).replace(/_/g, ' '),
    })),
  )

  const form = reactive<InternalForm>({
    line_number: null,
    item_type: '',
    item_name: '',
    inventory_item_uuid: null,
    quantity: null,
    quantity_unit: '',
    unit_cost_dollars: null,
    currency: 'USD',
  })

  const isEditing = computed(() => !!props.line)

  const isValid = computed(() => {
    return (
      form.line_number !== null
      && form.line_number > 0
      && form.item_type.trim().length > 0
      && form.item_name.trim().length > 0
      && form.quantity !== null
      && form.quantity > 0
      && form.quantity_unit.trim().length > 0
      && form.unit_cost_dollars !== null
      && form.unit_cost_dollars >= 0
      && form.currency.trim().length > 0
    )
  })

  watch(() => props.modelValue, open => {
    if (open) {
      if (props.line) {
        form.line_number = props.line.line_number
        form.item_type = props.line.item_type
        form.item_name = props.line.item_name
        form.inventory_item_uuid = props.line.inventory_item_uuid ?? null
        form.quantity = props.line.quantity
        form.quantity_unit = props.line.quantity_unit
        form.unit_cost_dollars = centsToDollars(props.line.unit_cost_cents)
        form.currency = props.line.currency
      } else {
        resetForm()
      }
    }
  })

  function resetForm () {
    form.line_number = null
    form.item_type = ''
    form.item_name = ''
    form.inventory_item_uuid = null
    form.quantity = null
    form.quantity_unit = ''
    form.unit_cost_dollars = null
    form.currency = 'USD'
  }

  function handleSave () {
    emit('save', {
      line_number: form.line_number,
      item_type: form.item_type,
      item_name: form.item_name,
      inventory_item_uuid: form.inventory_item_uuid,
      quantity: form.quantity,
      quantity_unit: form.quantity_unit,
      unit_cost_cents: dollarsToCents(form.unit_cost_dollars),
      currency: form.currency,
    })
  }

  function handleCancel () {
    emit('cancel')
    emit('update:modelValue', false)
  }
</script>
