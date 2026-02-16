<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 600"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Create ingredient lot</v-card-title>
      <v-card-text>
        <v-select
          v-model="form.ingredient_uuid"
          density="comfortable"
          :items="filteredIngredientSelectItems"
          label="Ingredient"
          :rules="[rules.required]"
        />
        <v-select
          v-model="form.receipt_uuid"
          clearable
          density="comfortable"
          :items="receiptSelectItems"
          label="Receipt (optional)"
        />
        <v-select
          v-model="form.supplier_uuid"
          clearable
          density="comfortable"
          :items="supplierSelectItems"
          label="Supplier"
        />
        <v-text-field
          v-model="form.brewery_lot_code"
          density="comfortable"
          label="Brewery lot code"
        />
        <v-text-field
          v-model="form.originator_lot_code"
          density="comfortable"
          label="Originator lot code"
        />
        <v-text-field
          v-model="form.originator_name"
          density="comfortable"
          label="Originator name"
        />
        <v-text-field
          v-model="form.originator_type"
          density="comfortable"
          label="Originator type"
        />
        <v-text-field
          v-model="form.received_at"
          density="comfortable"
          label="Received at"
          type="datetime-local"
        />
        <v-text-field
          v-model="form.received_amount"
          density="comfortable"
          label="Received amount"
          :rules="[rules.required]"
          type="number"
        />
        <v-combobox
          v-model="form.received_unit"
          density="comfortable"
          :items="unitOptions"
          label="Received unit"
          :rules="[rules.required]"
        />
        <v-text-field
          v-model="form.best_by_at"
          density="comfortable"
          label="Best by"
          type="datetime-local"
        />
        <v-text-field
          v-model="form.expires_at"
          density="comfortable"
          label="Expires at"
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
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSubmit"
        >
          Create lot
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateIngredientLotRequest, Ingredient, InventoryReceipt, Supplier } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { normalizeDateTime, normalizeText, toNumber } from '@/utils/normalize'

  const otherCategories = ['adjunct', 'salt', 'chemical', 'gas', 'other']
  const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']

  const props = defineProps<{
    modelValue: boolean
    ingredients: Ingredient[]
    receipts: InventoryReceipt[]
    suppliers: Supplier[]
    category: string
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateIngredientLotRequest]
  }>()

  const form = reactive({
    ingredient_uuid: null as string | null,
    receipt_uuid: null as string | null,
    supplier_uuid: null as string | null,
    brewery_lot_code: '',
    originator_lot_code: '',
    originator_name: '',
    originator_type: '',
    received_at: '',
    received_amount: '',
    received_unit: '',
    best_by_at: '',
    expires_at: '',
    notes: '',
  })

  const rules = {
    required: (v: string | null) => (v !== null && v !== '' && String(v).trim() !== '') || 'Required',
  }

  const filteredIngredientSelectItems = computed(() => {
    const categoriesToInclude = props.category === 'other'
      ? otherCategories
      : [props.category]

    return props.ingredients
      .filter(ingredient => categoriesToInclude.includes(ingredient.category))
      .map(ingredient => ({
        title: ingredient.name,
        value: ingredient.uuid,
      }))
  })

  const receiptSelectItems = computed(() =>
    props.receipts.map(receipt => ({
      title: receipt.reference_code || 'Unknown Receipt',
      value: receipt.uuid,
    })),
  )

  const supplierSelectItems = computed(() =>
    props.suppliers.map(supplier => ({
      title: supplier.name,
      value: supplier.uuid,
    })),
  )

  const isFormValid = computed(() => {
    return form.ingredient_uuid && form.received_amount && form.received_unit
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
    form.ingredient_uuid = null
    form.receipt_uuid = null
    form.supplier_uuid = null
    form.brewery_lot_code = ''
    form.originator_lot_code = ''
    form.originator_name = ''
    form.originator_type = ''
    form.received_at = ''
    form.received_amount = ''
    form.received_unit = ''
    form.best_by_at = ''
    form.expires_at = ''
    form.notes = ''
  }

  function handleSubmit () {
    if (!isFormValid.value) return

    const payload: CreateIngredientLotRequest = {
      ingredient_uuid: form.ingredient_uuid ?? '',
      receipt_uuid: form.receipt_uuid,
      supplier_uuid: form.supplier_uuid,
      brewery_lot_code: normalizeText(form.brewery_lot_code),
      originator_lot_code: normalizeText(form.originator_lot_code),
      originator_name: normalizeText(form.originator_name),
      originator_type: normalizeText(form.originator_type),
      received_at: normalizeDateTime(form.received_at),
      received_amount: toNumber(form.received_amount) ?? 0,
      received_unit: form.received_unit,
      best_by_at: normalizeDateTime(form.best_by_at),
      expires_at: normalizeDateTime(form.expires_at),
      notes: normalizeText(form.notes),
    }
    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
