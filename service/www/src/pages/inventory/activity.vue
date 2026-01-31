<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Inventory activity
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
                Activity log
                <v-spacer />
                <v-btn size="small" variant="text" @click="loadMovements">Apply filter</v-btn>
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
                      v-model="filters.ingredient_lot_id"
                      :items="lotSelectItems"
                      label="Filter by ingredient lot"
                      clearable
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-select
                      v-model="filters.beer_lot_id"
                      :items="beerLotSelectItems"
                      label="Filter by beer lot"
                      clearable
                    />
                  </v-col>
                </v-row>
                <v-table class="data-table" density="compact">
                  <thead>
                    <tr>
                      <th>Direction</th>
                      <th>Reason</th>
                      <th>Amount</th>
                      <th>Location</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="movement in movements" :key="movement.id">
                      <td>{{ movement.direction }}</td>
                      <td>{{ movement.reason }}</td>
                      <td>{{ formatAmount(movement.amount, movement.amount_unit) }}</td>
                      <td>{{ locationName(movement.stock_location_id) }}</td>
                    </tr>
                    <tr v-if="movements.length === 0">
                      <td colspan="4">No activity yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Log activity</v-card-title>
              <v-card-text>
                <v-select
                  v-model="movementForm.ingredient_lot_id"
                  :items="lotSelectItems"
                  label="Ingredient lot (optional)"
                  clearable
                />
                <v-select
                  v-model="movementForm.beer_lot_id"
                  :items="beerLotSelectItems"
                  label="Beer lot (optional)"
                  clearable
                />
                <v-select
                  v-model="movementForm.stock_location_id"
                  :items="locationSelectItems"
                  label="Stock location"
                />
                <v-select
                  v-model="movementForm.direction"
                  :items="movementDirectionOptions"
                  label="Direction"
                />
                <v-text-field v-model="movementForm.reason" label="Reason" />
                <v-text-field v-model="movementForm.amount" label="Amount" type="number" />
                <v-combobox
                  v-model="movementForm.amount_unit"
                  :items="unitOptions"
                  label="Amount unit"
                />
                <v-text-field v-model="movementForm.occurred_at" label="Occurred at" type="datetime-local" />
                <v-text-field v-model="movementForm.receipt_id" label="Receipt ID" type="number" />
                <v-text-field v-model="movementForm.usage_id" label="Usage ID" type="number" />
                <v-text-field v-model="movementForm.adjustment_id" label="Adjustment ID" type="number" />
                <v-text-field v-model="movementForm.transfer_id" label="Transfer ID" type="number" />
                <v-textarea
                  v-model="movementForm.notes"
                  auto-grow
                  label="Notes"
                  rows="2"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="
                    !movementForm.stock_location_id ||
                    !movementForm.direction ||
                    !movementForm.reason.trim() ||
                    !movementForm.amount ||
                    !movementForm.amount_unit
                  "
                  @click="createMovement"
                >
                  Log activity
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
import { useInventoryApi } from '@/composables/useInventoryApi'

type IngredientLot = {
  id: number
  ingredient_id: number
  received_amount: number
  received_unit: string
}

type StockLocation = {
  id: number
  name: string
}

type BeerLot = {
  id: number
  lot_code: string
}

type InventoryMovement = {
  id: number
  ingredient_lot_id: number | null
  beer_lot_id: number | null
  stock_location_id: number
  direction: string
  reason: string
  amount: number
  amount_unit: string
  occurred_at: string
  receipt_id: number | null
  usage_id: number | null
  adjustment_id: number | null
  transfer_id: number | null
  notes: string
}

const { request, normalizeText, normalizeDateTime, toNumber, formatAmount } = useInventoryApi()

const lots = ref<IngredientLot[]>([])
const locations = ref<StockLocation[]>([])
const beerLots = ref<BeerLot[]>([])
const movements = ref<InventoryMovement[]>([])
const loading = ref(false)
const errorMessage = ref('')

const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']
const movementDirectionOptions = ['in', 'out']

const filters = reactive({
  ingredient_lot_id: null as number | null,
  beer_lot_id: null as number | null,
})

const movementForm = reactive({
  ingredient_lot_id: null as number | null,
  beer_lot_id: null as number | null,
  stock_location_id: null as number | null,
  direction: '',
  reason: '',
  amount: '',
  amount_unit: '',
  occurred_at: '',
  receipt_id: '',
  usage_id: '',
  adjustment_id: '',
  transfer_id: '',
  notes: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

const lotSelectItems = computed(() =>
  lots.value.map((lot) => ({
    title: `Lot ${lot.id} (${lot.received_amount} ${lot.received_unit})`,
    value: lot.id,
  })),
)

const locationSelectItems = computed(() =>
  locations.value.map((location) => ({
    title: location.name,
    value: location.id,
  })),
)

const beerLotSelectItems = computed(() =>
  beerLots.value.map((lot) => ({
    title: lot.lot_code || `Beer lot ${lot.id}`,
    value: lot.id,
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
    await Promise.all([loadLots(), loadLocations(), loadBeerLots(), loadMovements()])
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load activity'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function loadLots() {
  lots.value = await request<IngredientLot[]>('/ingredient-lots')
}

async function loadLocations() {
  locations.value = await request<StockLocation[]>('/stock-locations')
}

async function loadBeerLots() {
  beerLots.value = await request<BeerLot[]>('/beer-lots')
}

async function loadMovements() {
  const query = new URLSearchParams()
  if (filters.ingredient_lot_id) {
    query.set('ingredient_lot_id', String(filters.ingredient_lot_id))
  }
  if (filters.beer_lot_id) {
    query.set('beer_lot_id', String(filters.beer_lot_id))
  }
  const path = query.toString() ? `/inventory-movements?${query.toString()}` : '/inventory-movements'
  movements.value = await request<InventoryMovement[]>(path)
}

async function createMovement() {
  try {
    const payload = {
      ingredient_lot_id: movementForm.ingredient_lot_id,
      beer_lot_id: movementForm.beer_lot_id,
      stock_location_id: movementForm.stock_location_id,
      direction: movementForm.direction,
      reason: movementForm.reason.trim(),
      amount: toNumber(movementForm.amount),
      amount_unit: movementForm.amount_unit,
      occurred_at: normalizeDateTime(movementForm.occurred_at),
      receipt_id: toNumber(movementForm.receipt_id),
      usage_id: toNumber(movementForm.usage_id),
      adjustment_id: toNumber(movementForm.adjustment_id),
      transfer_id: toNumber(movementForm.transfer_id),
      notes: normalizeText(movementForm.notes),
    }
    await request<InventoryMovement>('/inventory-movements', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    movementForm.ingredient_lot_id = null
    movementForm.beer_lot_id = null
    movementForm.stock_location_id = null
    movementForm.direction = ''
    movementForm.reason = ''
    movementForm.amount = ''
    movementForm.amount_unit = ''
    movementForm.occurred_at = ''
    movementForm.receipt_id = ''
    movementForm.usage_id = ''
    movementForm.adjustment_id = ''
    movementForm.transfer_id = ''
    movementForm.notes = ''
    await loadMovements()
    showNotice('Activity recorded')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to log activity'
    errorMessage.value = message
    showNotice(message, 'error')
  }
}

function locationName(locationId: number) {
  return locations.value.find((location) => location.id === locationId)?.name ?? `Location ${locationId}`
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
