<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Ingredient lots
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="refreshAll">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title class="d-flex align-center">
                Lots list
                <v-spacer />
                <v-btn size="small" variant="text" @click="loadLots">Apply filter</v-btn>
              </v-card-title>
              <v-card-text>
                <v-alert
                  v-if="errorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ errorMessage }}
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
                      <td>{{ formatAmount(lot.received_amount, lot.received_unit) }}</td>
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

type Ingredient = {
  id: number
  name: string
}

type InventoryReceipt = {
  id: number
  reference_code: string
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

const { request, normalizeText, normalizeDateTime, toNumber, formatDateTime, formatAmount } = useInventoryApi()
const router = useRouter()

const ingredients = ref<Ingredient[]>([])
const receipts = ref<InventoryReceipt[]>([])
const lots = ref<IngredientLot[]>([])
const loading = ref(false)
const errorMessage = ref('')

const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']

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
  loading.value = true
  errorMessage.value = ''
  try {
    await Promise.all([loadIngredients(), loadReceipts(), loadLots()])
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load lots'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function loadIngredients() {
  ingredients.value = await request<Ingredient[]>('/ingredients')
}

async function loadReceipts() {
  receipts.value = await request<InventoryReceipt[]>('/inventory-receipts')
}

async function loadLots() {
  const query = new URLSearchParams()
  if (lotFilters.ingredient_id) {
    query.set('ingredient_id', String(lotFilters.ingredient_id))
  }
  if (lotFilters.receipt_id) {
    query.set('receipt_id', String(lotFilters.receipt_id))
  }
  const path = query.toString() ? `/ingredient-lots?${query.toString()}` : '/ingredient-lots'
  lots.value = await request<IngredientLot[]>(path)
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
    errorMessage.value = message
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
