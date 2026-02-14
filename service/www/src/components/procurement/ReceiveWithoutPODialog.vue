<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 700"
    :model-value="modelValue"
    persistent
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <span class="text-h6">Receive Inventory</span>
        <v-spacer />
        <v-btn
          :disabled="saving"
          icon="mdi-close"
          size="small"
          variant="text"
          @click="handleClose"
        />
      </v-card-title>

      <v-divider />

      <v-card-text class="pa-4">
        <!-- Loading state for reference data -->
        <template v-if="loadingReferenceData">
          <div class="d-flex flex-column align-center py-8">
            <v-progress-circular color="primary" indeterminate size="48" />
            <p class="text-body-2 text-medium-emphasis mt-4">Loading reference data...</p>
          </div>
        </template>

        <template v-else>
          <!-- Reference data load warning -->
          <v-alert
            v-if="referenceDataWarning"
            class="mb-4"
            closable
            density="compact"
            type="warning"
            variant="tonal"
            @click:close="referenceDataWarning = ''"
          >
            {{ referenceDataWarning }}
          </v-alert>

          <!-- Error alert -->
          <v-alert
            v-if="saveError"
            class="mb-4"
            closable
            density="compact"
            type="error"
            variant="tonal"
            @click:close="saveError = ''"
          >
            {{ saveError }}
          </v-alert>

          <!-- Receipt info -->
          <v-row>
            <v-col cols="12" sm="6">
              <v-select
                v-model="form.supplierUuid"
                clearable
                density="comfortable"
                hint="Optional"
                :items="supplierItems"
                label="Supplier"
                persistent-hint
              />
            </v-col>
            <v-col cols="12" sm="6">
              <v-text-field
                v-model="form.referenceCode"
                density="comfortable"
                hint="Optional invoice or delivery number"
                label="Reference code"
                persistent-hint
              />
            </v-col>
          </v-row>

          <v-divider class="my-4" />

          <!-- Items section -->
          <div class="d-flex align-center mb-3">
            <span class="text-subtitle-1 font-weight-medium">Items</span>
            <v-spacer />
            <v-btn
              color="primary"
              prepend-icon="mdi-plus"
              size="small"
              variant="text"
              @click="addItem"
            >
              Add item
            </v-btn>
          </div>

          <v-alert
            v-if="form.items.length === 0"
            density="compact"
            type="info"
            variant="tonal"
          >
            Add at least one item to receive.
          </v-alert>

          <!-- Item rows -->
          <div v-for="(item, index) in form.items" :key="index" class="mb-4">
            <v-card variant="outlined">
              <v-card-text class="pa-3">
                <div class="d-flex align-center mb-2">
                  <span class="text-body-2 font-weight-medium">Item {{ index + 1 }}</span>
                  <v-spacer />
                  <v-btn
                    color="error"
                    density="compact"
                    icon="mdi-delete"
                    size="small"
                    variant="text"
                    @click="removeItem(index)"
                  />
                </div>

                <v-row dense>
                  <v-col cols="12">
                    <v-select
                      v-model="item.ingredientUuid"
                      density="comfortable"
                      :items="ingredientItems"
                      label="Ingredient *"
                      :rules="[v => !!v || 'Required']"
                    />
                  </v-col>
                  <v-col cols="6" sm="4">
                    <v-text-field
                      v-model.number="item.quantity"
                      density="comfortable"
                      inputmode="decimal"
                      label="Quantity *"
                      min="0"
                      :rules="[v => v > 0 || 'Required']"
                      step="any"
                      type="number"
                    />
                  </v-col>
                  <v-col cols="6" sm="4">
                    <v-combobox
                      v-model="item.unit"
                      density="comfortable"
                      :items="unitOptions"
                      label="Unit *"
                      :rules="[v => !!v || 'Required']"
                    />
                  </v-col>
                  <v-col cols="12" sm="4">
                    <v-select
                      v-model="item.locationUuid"
                      density="comfortable"
                      :items="locationItems"
                      label="Location *"
                      :rules="[v => !!v || 'Required']"
                    />
                  </v-col>
                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model="item.breweryLotCode"
                      density="comfortable"
                      hint="Optional"
                      label="Brewery lot code"
                      persistent-hint
                    />
                  </v-col>
                  <v-col cols="12" sm="6">
                    <v-text-field
                      v-model="item.supplierLotCode"
                      density="comfortable"
                      hint="Optional"
                      label="Supplier lot code"
                      persistent-hint
                    />
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>
          </div>

          <v-divider class="my-4" />

          <!-- Notes -->
          <v-textarea
            v-model="form.notes"
            auto-grow
            density="comfortable"
            hint="Optional notes about this receipt"
            label="Notes"
            persistent-hint
            rows="2"
          />
        </template>
      </v-card-text>

      <v-divider />

      <v-card-actions class="justify-end pa-4">
        <v-btn :disabled="saving" variant="text" @click="handleClose">
          Cancel
        </v-btn>
        <v-btn
          color="success"
          :disabled="!isValid || loadingReferenceData"
          :loading="saving"
          @click="handleSave"
        >
          Receive Items
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type {
    CreateIngredientLotRequest,
    CreateInventoryMovementRequest,
    CreateInventoryReceiptRequest,
    Ingredient,
    StockLocation,
    Supplier,
  } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProcurementApi } from '@/composables/useProcurementApi'

  interface ReceiveItem {
    ingredientUuid: string
    quantity: number
    unit: string
    locationUuid: string
    breweryLotCode: string
    supplierLotCode: string
  }

  const props = defineProps<{
    modelValue: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'received': []
  }>()

  const {
    getIngredients,
    getStockLocations,
    createInventoryReceipt,
    createIngredientLot,
    createInventoryMovement,
  } = useInventoryApi()

  const { getSuppliers } = useProcurementApi()

  // Form state
  const form = reactive({
    supplierUuid: '' as string,
    referenceCode: '',
    items: [] as ReceiveItem[],
    notes: '',
  })

  // Reference data
  const suppliers = ref<Supplier[]>([])
  const ingredients = ref<Ingredient[]>([])
  const stockLocations = ref<StockLocation[]>([])
  const loadingReferenceData = ref(false)
  const referenceDataWarning = ref('')

  // State
  const saving = ref(false)
  const saveError = ref('')

  // Options
  const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']

  // Computed
  const supplierItems = computed(() =>
    suppliers.value.map(s => ({ title: s.name, value: s.uuid })),
  )

  const ingredientItems = computed(() =>
    ingredients.value.map(i => ({ title: `${i.name} (${i.category})`, value: i.uuid })),
  )

  const locationItems = computed(() =>
    stockLocations.value.map(l => ({ title: l.name, value: l.uuid })),
  )

  const isValid = computed(() => {
    if (form.items.length === 0) return false
    return form.items.every(
      item =>
        item.ingredientUuid
        && item.quantity > 0
        && item.unit
        && item.locationUuid,
    )
  })

  // Watch for dialog open
  watch(() => props.modelValue, async open => {
    if (open) {
      resetForm()
      await loadReferenceData()
    }
  })

  async function loadReferenceData () {
    loadingReferenceData.value = true
    referenceDataWarning.value = ''

    const warnings: string[] = []

    try {
      const [suppliersResult, ingredientsResult, locationsResult] = await Promise.allSettled([
        getSuppliers(),
        getIngredients(),
        getStockLocations(),
      ])

      if (suppliersResult.status === 'fulfilled') {
        suppliers.value = suppliersResult.value
      } else {
        warnings.push('suppliers')
      }

      if (ingredientsResult.status === 'fulfilled') {
        ingredients.value = ingredientsResult.value
      } else {
        warnings.push('ingredients')
      }

      if (locationsResult.status === 'fulfilled') {
        stockLocations.value = locationsResult.value
      } else {
        warnings.push('stock locations')
      }

      if (warnings.length > 0) {
        referenceDataWarning.value = `Failed to load ${warnings.join(', ')}. Some options may be unavailable.`
      }
    } catch {
      referenceDataWarning.value = 'Failed to load reference data. Some options may be unavailable.'
    } finally {
      loadingReferenceData.value = false
    }
  }

  function resetForm () {
    form.supplierUuid = ''
    form.referenceCode = ''
    form.items = []
    form.notes = ''
    saveError.value = ''
    saving.value = false
  }

  function addItem () {
    form.items.push({
      ingredientUuid: '',
      quantity: 0,
      unit: 'kg',
      locationUuid: '',
      breweryLotCode: '',
      supplierLotCode: '',
    })
  }

  function removeItem (index: number) {
    form.items.splice(index, 1)
  }

  function handleClose () {
    emit('update:modelValue', false)
  }

  async function handleSave () {
    if (!isValid.value) return

    saving.value = true
    saveError.value = ''

    try {
      // 1. Create inventory receipt
      const receiptPayload: CreateInventoryReceiptRequest = {
        supplier_uuid: form.supplierUuid || null,
        reference_code: form.referenceCode.trim() || null,
        received_at: new Date().toISOString(),
        notes: form.notes.trim() || null,
      }
      const receipt = await createInventoryReceipt(receiptPayload)

      // 2. For each item, create lot and movement
      // Track successes and failures for partial completion reporting
      const createdLots: string[] = []
      const failedItems: string[] = []

      for (const item of form.items) {
        // Get ingredient name for error reporting
        const ingredient = ingredients.value.find(i => i.uuid === item.ingredientUuid)
        const itemName = ingredient?.name ?? 'Unknown ingredient'

        try {
          const lotPayload: CreateIngredientLotRequest = {
            ingredient_uuid: item.ingredientUuid,
            receipt_uuid: receipt.uuid,
            supplier_uuid: form.supplierUuid || null,
            brewery_lot_code: item.breweryLotCode.trim() || null,
            originator_lot_code: item.supplierLotCode.trim() || null,
            received_at: new Date().toISOString(),
            received_amount: item.quantity,
            received_unit: item.unit,
          }
          const lot = await createIngredientLot(lotPayload)
          createdLots.push(lot.uuid)

          const movementPayload: CreateInventoryMovementRequest = {
            ingredient_lot_uuid: lot.uuid,
            stock_location_uuid: item.locationUuid,
            direction: 'in',
            reason: 'receive',
            amount: item.quantity,
            amount_unit: item.unit,
            receipt_uuid: receipt.uuid,
          }
          await createInventoryMovement(movementPayload)
        } catch {
          failedItems.push(itemName)
        }
      }

      // Check for partial failures
      if (failedItems.length > 0) {
        saveError.value = createdLots.length > 0 ? `Partially completed. Created ${createdLots.length} lot(s). Failed: ${failedItems.join(', ')}. Please review inventory and retry failed items.` : `Failed to create lots for: ${failedItems.join(', ')}. Receipt was created but no inventory was added.`
        // Don't close dialog on partial failure - let user see the error
        return
      }

      emit('received')
      emit('update:modelValue', false)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save receipt'
      saveError.value = message
    } finally {
      saving.value = false
    }
  }
</script>
