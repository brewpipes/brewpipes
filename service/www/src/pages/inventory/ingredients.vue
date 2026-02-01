<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title>Ingredients</v-card-title>
      <v-card-text>
        <v-tabs v-model="activeTab" class="inventory-tabs" color="primary" show-arrows>
          <v-tab value="stock">Stock</v-tab>
          <v-tab value="usage">Usage</v-tab>
          <v-tab value="types">Types</v-tab>
        </v-tabs>

        <v-window v-model="activeTab" class="mt-4">
          <v-window-item value="stock">
            <v-row align="stretch">
              <v-col cols="12" md="7">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title class="d-flex align-center">
                    Receipts
                    <v-spacer />
                    <v-btn size="small" variant="text" :loading="receiptLoading" @click="loadReceipts">
                      Refresh
                    </v-btn>
                  </v-card-title>
                  <v-card-text>
                    <v-alert
                      v-if="receiptErrorMessage"
                      class="mb-3"
                      density="compact"
                      type="error"
                      variant="tonal"
                    >
                      {{ receiptErrorMessage }}
                    </v-alert>
                    <v-table class="data-table" density="compact">
                      <thead>
                        <tr>
                          <th>Reference</th>
                          <th>Supplier</th>
                          <th>Received at</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="receipt in receipts" :key="receipt.id">
                          <td>{{ receipt.reference_code || 'n/a' }}</td>
                          <td>{{ receipt.supplier_uuid || 'n/a' }}</td>
                          <td>{{ formatDateTime(receipt.received_at) }}</td>
                        </tr>
                        <tr v-if="receipts.length === 0">
                          <td colspan="3">No receipts yet.</td>
                        </tr>
                      </tbody>
                    </v-table>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="5">
                <v-card class="sub-card" variant="tonal">
                  <v-card-title>Create receipt</v-card-title>
                  <v-card-text>
                    <v-text-field v-model="receiptForm.reference_code" label="Reference code" />
                    <v-text-field v-model="receiptForm.supplier_uuid" label="Supplier UUID" />
                    <v-text-field v-model="receiptForm.received_at" label="Received at" type="datetime-local" />
                    <v-textarea
                      v-model="receiptForm.notes"
                      auto-grow
                      label="Notes"
                      rows="2"
                    />
                    <v-btn block color="primary" @click="createReceipt">
                      Add receipt
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-divider class="my-6" />

            <v-row align="stretch">
              <v-col cols="12" md="7">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title class="d-flex align-center">
                    Ingredient lots
                    <v-spacer />
                    <v-btn size="small" variant="text" :loading="lotLoading" @click="loadLots">
                      Apply filter
                    </v-btn>
                  </v-card-title>
                  <v-card-text>
                    <v-alert
                      v-if="lotErrorMessage"
                      class="mb-3"
                      density="compact"
                      type="error"
                      variant="tonal"
                    >
                      {{ lotErrorMessage }}
                    </v-alert>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-select
                          v-model="lotFilters.ingredient_id"
                          :items="ingredientSelectItems"
                          label="Filter by ingredient"
                          clearable
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-select
                          v-model="lotFilters.receipt_id"
                          :items="receiptSelectItems"
                          label="Filter by receipt"
                          clearable
                        />
                      </v-col>
                    </v-row>
                    <v-table class="data-table" density="compact">
                      <thead>
                        <tr>
                          <th>Ingredient</th>
                          <th>Received</th>
                          <th>Best by</th>
                          <th>Expires</th>
                          <th></th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="lot in lots" :key="lot.id">
                          <td>{{ ingredientName(lot.ingredient_id) }}</td>
                          <td>{{ formatAmountPreferred(lot.received_amount, lot.received_unit) }}</td>
                          <td>{{ formatDateTime(lot.best_by_at) }}</td>
                          <td>{{ formatDateTime(lot.expires_at) }}</td>
                          <td>
                            <v-btn size="x-small" variant="text" @click="openLotDetails(lot.id)">
                              Details
                            </v-btn>
                          </td>
                        </tr>
                        <tr v-if="lots.length === 0">
                          <td colspan="5">No lots yet.</td>
                        </tr>
                      </tbody>
                    </v-table>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="5">
                <v-card class="sub-card" variant="tonal">
                  <v-card-title>Create lot</v-card-title>
                  <v-card-text>
                    <v-select
                      v-model="lotForm.ingredient_id"
                      :items="ingredientSelectItems"
                      label="Ingredient"
                    />
                    <v-select
                      v-model="lotForm.receipt_id"
                      :items="receiptSelectItems"
                      label="Receipt (optional)"
                      clearable
                    />
                    <v-text-field v-model="lotForm.supplier_uuid" label="Supplier UUID" />
                    <v-text-field v-model="lotForm.brewery_lot_code" label="Brewery lot code" />
                    <v-text-field v-model="lotForm.originator_lot_code" label="Originator lot code" />
                    <v-text-field v-model="lotForm.originator_name" label="Originator name" />
                    <v-text-field v-model="lotForm.originator_type" label="Originator type" />
                    <v-text-field v-model="lotForm.received_at" label="Received at" type="datetime-local" />
                    <v-text-field v-model="lotForm.received_amount" label="Received amount" type="number" />
                    <v-combobox
                      v-model="lotForm.received_unit"
                      :items="unitOptions"
                      label="Received unit"
                    />
                    <v-text-field v-model="lotForm.best_by_at" label="Best by" type="datetime-local" />
                    <v-text-field v-model="lotForm.expires_at" label="Expires at" type="datetime-local" />
                    <v-textarea
                      v-model="lotForm.notes"
                      auto-grow
                      label="Notes"
                      rows="2"
                    />
                    <v-btn
                      block
                      color="primary"
                      :disabled="!lotForm.ingredient_id || !lotForm.received_amount || !lotForm.received_unit"
                      @click="createLot"
                    >
                      Add lot
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-window-item>

          <v-window-item value="usage">
            <v-row align="stretch">
              <v-col cols="12" md="7">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title class="d-flex align-center">
                    Usage log
                    <v-spacer />
                    <v-btn size="small" variant="text" :loading="usageLoading" @click="loadUsage">
                      Refresh
                    </v-btn>
                  </v-card-title>
                  <v-card-text>
                    <v-alert
                      v-if="usageErrorMessage"
                      class="mb-3"
                      density="compact"
                      type="error"
                      variant="tonal"
                    >
                      {{ usageErrorMessage }}
                    </v-alert>
                    <v-table class="data-table" density="compact">
                      <thead>
                        <tr>
                          <th>Production ref</th>
                          <th>Used at</th>
                          <th>Notes</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="usage in usages" :key="usage.id">
                          <td>{{ usage.production_ref_uuid || 'n/a' }}</td>
                          <td>{{ formatDateTime(usage.used_at) }}</td>
                          <td>{{ usage.notes || '' }}</td>
                        </tr>
                        <tr v-if="usages.length === 0">
                          <td colspan="3">No usage records yet.</td>
                        </tr>
                      </tbody>
                    </v-table>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="5">
                <v-card class="sub-card" variant="tonal">
                  <v-card-title>Create usage</v-card-title>
                  <v-card-text>
                    <v-text-field v-model="usageForm.production_ref_uuid" label="Production ref UUID" />
                    <v-text-field v-model="usageForm.used_at" label="Used at" type="datetime-local" />
                    <v-textarea
                      v-model="usageForm.notes"
                      auto-grow
                      label="Notes"
                      rows="2"
                    />
                    <v-btn block color="primary" @click="createUsage">
                      Add usage
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-window-item>

          <v-window-item value="types">
            <v-row align="stretch">
              <v-col cols="12" md="7">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title class="d-flex align-center">
                    Ingredient types
                    <v-spacer />
                    <v-btn size="small" variant="text" :loading="ingredientLoading" @click="loadIngredients">
                      Refresh
                    </v-btn>
                  </v-card-title>
                  <v-card-text>
                    <v-alert
                      v-if="ingredientErrorMessage"
                      class="mb-3"
                      density="compact"
                      type="error"
                      variant="tonal"
                    >
                      {{ ingredientErrorMessage }}
                    </v-alert>
                    <v-table class="data-table" density="compact">
                      <thead>
                        <tr>
                          <th>Name</th>
                          <th>Category</th>
                          <th>Unit</th>
                          <th>Updated</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="ingredient in ingredients" :key="ingredient.id">
                          <td>{{ ingredient.name }}</td>
                          <td>{{ ingredient.category }}</td>
                          <td>{{ ingredient.default_unit }}</td>
                          <td>{{ formatDateTime(ingredient.updated_at) }}</td>
                        </tr>
                        <tr v-if="ingredients.length === 0">
                          <td colspan="4">No ingredients yet.</td>
                        </tr>
                      </tbody>
                    </v-table>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="5">
                <v-card class="sub-card" variant="tonal">
                  <v-card-title>Create ingredient</v-card-title>
                  <v-card-text>
                    <v-text-field v-model="ingredientForm.name" label="Name" />
                    <v-combobox
                      v-model="ingredientForm.category"
                      :items="ingredientCategoryOptions"
                      label="Category"
                    />
                    <v-combobox
                      v-model="ingredientForm.default_unit"
                      :items="unitOptions"
                      label="Default unit"
                    />
                    <v-textarea
                      v-model="ingredientForm.description"
                      auto-grow
                      label="Description"
                      rows="2"
                    />
                    <v-btn
                      block
                      color="primary"
                      :disabled="
                        !ingredientForm.name.trim() ||
                        !ingredientForm.category ||
                        !ingredientForm.default_unit
                      "
                      @click="createIngredient"
                    >
                      Add ingredient
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useInventoryApi } from '@/composables/useInventoryApi'
import { useUnitPreferences } from '@/composables/useUnitPreferences'

type Ingredient = {
  id: number
  uuid: string
  name: string
  category: string
  default_unit: string
  description: string
  created_at: string
  updated_at: string
}

type InventoryReceipt = {
  id: number
  uuid: string
  supplier_uuid: string
  reference_code: string
  received_at: string
  notes: string
  created_at: string
  updated_at: string
}

type IngredientLot = {
  id: number
  ingredient_id: number
  receipt_id: number | null
  received_amount: number
  received_unit: string
  best_by_at: string
  expires_at: string
  supplier_uuid: string
  brewery_lot_code: string
  originator_lot_code: string
  originator_name: string
  originator_type: string
  received_at: string
  notes: string
}

type InventoryUsage = {
  id: number
  uuid: string
  production_ref_uuid: string
  used_at: string
  notes: string
  created_at: string
  updated_at: string
}

const { request, normalizeText, normalizeDateTime, toNumber, formatDateTime } = useInventoryApi()
const { formatAmountPreferred } = useUnitPreferences()
const router = useRouter()

const activeTab = ref('stock')

const ingredients = ref<Ingredient[]>([])
const receipts = ref<InventoryReceipt[]>([])
const lots = ref<IngredientLot[]>([])
const usages = ref<InventoryUsage[]>([])

const ingredientLoading = ref(false)
const receiptLoading = ref(false)
const lotLoading = ref(false)
const usageLoading = ref(false)

const ingredientErrorMessage = ref('')
const receiptErrorMessage = ref('')
const lotErrorMessage = ref('')
const usageErrorMessage = ref('')

const ingredientCategoryOptions = ['malt', 'hop', 'yeast', 'adjunct', 'water_chem', 'gas', 'other']
const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']

const ingredientForm = reactive({
  name: '',
  category: '',
  default_unit: '',
  description: '',
})

const receiptForm = reactive({
  supplier_uuid: '',
  reference_code: '',
  received_at: '',
  notes: '',
})

const lotFilters = reactive({
  ingredient_id: null as number | null,
  receipt_id: null as number | null,
})

const lotForm = reactive({
  ingredient_id: null as number | null,
  receipt_id: null as number | null,
  supplier_uuid: '',
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

const usageForm = reactive({
  production_ref_uuid: '',
  used_at: '',
  notes: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

const ingredientSelectItems = computed(() =>
  ingredients.value.map((ingredient) => ({
    title: ingredient.name,
    value: ingredient.id,
  })),
)

const receiptSelectItems = computed(() =>
  receipts.value.map((receipt) => ({
    title: receipt.reference_code || `Receipt ${receipt.id}`,
    value: receipt.id,
  })),
)

onMounted(async () => {
  await refreshAll()
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function refreshAll() {
  await Promise.allSettled([loadIngredients(), loadReceipts(), loadLots(), loadUsage()])
}

async function loadIngredients() {
  ingredientLoading.value = true
  ingredientErrorMessage.value = ''
  try {
    ingredients.value = await request<Ingredient[]>('/ingredients')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load ingredients'
    ingredientErrorMessage.value = message
  } finally {
    ingredientLoading.value = false
  }
}

async function loadReceipts() {
  receiptLoading.value = true
  receiptErrorMessage.value = ''
  try {
    receipts.value = await request<InventoryReceipt[]>('/inventory-receipts')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load receipts'
    receiptErrorMessage.value = message
  } finally {
    receiptLoading.value = false
  }
}

async function loadLots() {
  lotLoading.value = true
  lotErrorMessage.value = ''
  try {
    const query = new URLSearchParams()
    if (lotFilters.ingredient_id) {
      query.set('ingredient_id', String(lotFilters.ingredient_id))
    }
    if (lotFilters.receipt_id) {
      query.set('receipt_id', String(lotFilters.receipt_id))
    }
    const path = query.toString() ? `/ingredient-lots?${query.toString()}` : '/ingredient-lots'
    lots.value = await request<IngredientLot[]>(path)
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load lots'
    lotErrorMessage.value = message
  } finally {
    lotLoading.value = false
  }
}

async function loadUsage() {
  usageLoading.value = true
  usageErrorMessage.value = ''
  try {
    usages.value = await request<InventoryUsage[]>('/inventory-usage')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load usage'
    usageErrorMessage.value = message
  } finally {
    usageLoading.value = false
  }
}

async function createIngredient() {
  try {
    const payload = {
      name: ingredientForm.name.trim(),
      category: ingredientForm.category,
      default_unit: ingredientForm.default_unit,
      description: normalizeText(ingredientForm.description),
    }
    await request<Ingredient>('/ingredients', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    ingredientForm.name = ''
    ingredientForm.category = ''
    ingredientForm.default_unit = ''
    ingredientForm.description = ''
    await loadIngredients()
    showNotice('Ingredient created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create ingredient'
    ingredientErrorMessage.value = message
    showNotice(message, 'error')
  }
}

async function createReceipt() {
  try {
    const payload = {
      supplier_uuid: normalizeText(receiptForm.supplier_uuid),
      reference_code: normalizeText(receiptForm.reference_code),
      received_at: normalizeDateTime(receiptForm.received_at),
      notes: normalizeText(receiptForm.notes),
    }
    await request<InventoryReceipt>('/inventory-receipts', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    receiptForm.supplier_uuid = ''
    receiptForm.reference_code = ''
    receiptForm.received_at = ''
    receiptForm.notes = ''
    await loadReceipts()
    showNotice('Receipt created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create receipt'
    receiptErrorMessage.value = message
    showNotice(message, 'error')
  }
}

async function createLot() {
  try {
    const payload = {
      ingredient_id: lotForm.ingredient_id,
      receipt_id: lotForm.receipt_id,
      supplier_uuid: normalizeText(lotForm.supplier_uuid),
      brewery_lot_code: normalizeText(lotForm.brewery_lot_code),
      originator_lot_code: normalizeText(lotForm.originator_lot_code),
      originator_name: normalizeText(lotForm.originator_name),
      originator_type: normalizeText(lotForm.originator_type),
      received_at: normalizeDateTime(lotForm.received_at),
      received_amount: toNumber(lotForm.received_amount),
      received_unit: lotForm.received_unit,
      best_by_at: normalizeDateTime(lotForm.best_by_at),
      expires_at: normalizeDateTime(lotForm.expires_at),
      notes: normalizeText(lotForm.notes),
    }
    await request<IngredientLot>('/ingredient-lots', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    lotForm.ingredient_id = null
    lotForm.receipt_id = null
    lotForm.supplier_uuid = ''
    lotForm.brewery_lot_code = ''
    lotForm.originator_lot_code = ''
    lotForm.originator_name = ''
    lotForm.originator_type = ''
    lotForm.received_at = ''
    lotForm.received_amount = ''
    lotForm.received_unit = ''
    lotForm.best_by_at = ''
    lotForm.expires_at = ''
    lotForm.notes = ''
    await loadLots()
    showNotice('Lot created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create lot'
    lotErrorMessage.value = message
    showNotice(message, 'error')
  }
}

async function createUsage() {
  try {
    const payload = {
      production_ref_uuid: normalizeText(usageForm.production_ref_uuid),
      used_at: normalizeDateTime(usageForm.used_at),
      notes: normalizeText(usageForm.notes),
    }
    await request<InventoryUsage>('/inventory-usage', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    usageForm.production_ref_uuid = ''
    usageForm.used_at = ''
    usageForm.notes = ''
    await loadUsage()
    showNotice('Usage created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create usage'
    usageErrorMessage.value = message
    showNotice(message, 'error')
  }
}

function ingredientName(ingredientId: number) {
  return ingredients.value.find((ingredient) => ingredient.id === ingredientId)?.name ?? `Ingredient ${ingredientId}`
}

function openLotDetails(lotId: number) {
  router.push({
    path: '/inventory/lot-details',
    query: { lot_id: String(lotId) },
  })
}
</script>

<style scoped>
.inventory-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
}

.data-table :deep(td) {
  font-size: 0.85rem;
}
</style>
