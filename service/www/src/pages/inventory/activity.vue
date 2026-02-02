<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Inventory activity
        <v-spacer />
        <v-btn :loading="loading" size="small" variant="text" @click="refreshAll">
          Refresh
        </v-btn>
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
        <v-row class="mb-4">
          <v-col cols="12" md="4">
            <v-select
              v-model="filters.ingredient_lot_id"
              clearable
              density="compact"
              hide-details
              :items="lotSelectItems"
              label="Filter by ingredient lot"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-select
              v-model="filters.beer_lot_id"
              clearable
              density="compact"
              hide-details
              :items="beerLotSelectItems"
              label="Filter by beer lot"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-btn color="primary" variant="tonal" @click="loadMovements">
              Apply filter
            </v-btn>
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
            <tr v-if="loading">
              <td class="text-center text-medium-emphasis" colspan="4">
                <v-progress-circular class="mr-2" indeterminate size="16" />
                Loading...
              </td>
            </tr>
            <tr v-else-if="movements.length === 0">
              <td class="text-medium-emphasis" colspan="4">No activity yet.</td>
            </tr>
            <tr v-for="movement in movements" v-else :key="movement.id">
              <td>{{ movement.direction }}</td>
              <td>{{ movement.reason }}</td>
              <td>{{ formatAmountPreferred(movement.amount, movement.amount_unit) }}</td>
              <td>{{ locationName(movement.stock_location_id) }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
    </v-card>
  </v-container>

</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

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

  const { request } = useInventoryApi()
  const { formatAmountPreferred } = useUnitPreferences()

  const lots = ref<IngredientLot[]>([])
  const locations = ref<StockLocation[]>([])
  const beerLots = ref<BeerLot[]>([])
  const movements = ref<InventoryMovement[]>([])
  const loading = ref(false)
  const errorMessage = ref('')

  const filters = reactive({
    ingredient_lot_id: null as number | null,
    beer_lot_id: null as number | null,
  })

  const lotSelectItems = computed(() =>
    lots.value.map(lot => ({
      title: `Lot ${lot.id} (${lot.received_amount} ${lot.received_unit})`,
      value: lot.id,
    })),
  )

  const beerLotSelectItems = computed(() =>
    beerLots.value.map(lot => ({
      title: lot.lot_code || `Beer lot ${lot.id}`,
      value: lot.id,
    })),
  )

  onMounted(async () => {
    await refreshAll()
  })

  async function refreshAll () {
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

  async function loadLots () {
    lots.value = await request<IngredientLot[]>('/ingredient-lots')
  }

  async function loadLocations () {
    locations.value = await request<StockLocation[]>('/stock-locations')
  }

  async function loadBeerLots () {
    beerLots.value = await request<BeerLot[]>('/beer-lots')
  }

  async function loadMovements () {
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

  function locationName (locationId: number) {
    return locations.value.find(location => location.id === locationId)?.name ?? `Location ${locationId}`
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
