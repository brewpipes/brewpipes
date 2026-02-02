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
              <th>Item</th>
              <th>Lot #</th>
              <th class="text-center">Direction</th>
              <th>Reason</th>
              <th>Amount</th>
              <th>Location</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading">
              <td class="text-center text-medium-emphasis" colspan="6">
                <v-progress-circular class="mr-2" indeterminate size="16" />
                Loading...
              </td>
            </tr>
            <tr v-else-if="movements.length === 0">
              <td class="text-medium-emphasis" colspan="6">No activity yet.</td>
            </tr>
            <tr v-for="movement in movements" v-else :key="movement.id">
              <td>{{ getItemName(movement) }}</td>
              <td>{{ getLotCode(movement) }}</td>
              <td class="text-center">
                <v-tooltip location="top">
                  <template #activator="{ props }">
                    <v-icon
                      v-bind="props"
                      :color="movement.direction === 'in' ? 'success' : 'warning'"
                      :icon="movement.direction === 'in' ? 'mdi-package-down' : 'mdi-package-up'"
                      size="small"
                    />
                  </template>
                  <span>{{ getDirectionTooltip(movement.direction) }}</span>
                </v-tooltip>
              </td>
              <td>
                <template v-if="movement.reason === 'use'">
                  <router-link
                    v-if="getUsageBatch(movement)"
                    class="reason-link"
                    :to="`/batches/${getUsageBatch(movement)!.uuid}`"
                  >
                    Used in {{ getUsageBatch(movement)!.short_name }}
                    <span v-if="getUsagePhase(movement)" class="text-medium-emphasis">
                      ({{ getUsagePhase(movement) }})
                    </span>
                  </router-link>
                  <span v-else class="text-medium-emphasis">Used in production</span>
                </template>
                <template v-else-if="movement.reason === 'receive'">
                  <span v-if="getReceiptSupplier(movement)">
                    Received from {{ getReceiptSupplier(movement)!.name }}
                  </span>
                  <span v-else class="text-medium-emphasis">Received</span>
                </template>
                <template v-else-if="movement.reason === 'adjust' || movement.reason === 'waste'">
                  <v-tooltip v-if="getAdjustmentNotes(movement)" location="top">
                    <template #activator="{ props }">
                      <span v-bind="props" class="adjustment-reason">
                        {{ formatAdjustmentReason(movement) }}
                        <v-icon class="ml-1" icon="mdi-information-outline" size="x-small" />
                      </span>
                    </template>
                    <span>{{ getAdjustmentNotes(movement) }}</span>
                  </v-tooltip>
                  <span v-else>{{ formatAdjustmentReason(movement) }}</span>
                </template>
                <template v-else-if="movement.reason === 'transfer'">
                  <v-tooltip v-if="getTransferNotes(movement)" location="top">
                    <template #activator="{ props }">
                      <span v-bind="props" class="transfer-reason">
                        {{ formatTransferReason(movement) }}
                        <v-icon class="ml-1" icon="mdi-information-outline" size="x-small" />
                      </span>
                    </template>
                    <span>{{ getTransferNotes(movement) }}</span>
                  </v-tooltip>
                  <span v-else>{{ formatTransferReason(movement) }}</span>
                </template>
                <template v-else>
                  <span class="text-medium-emphasis">{{ movement.reason }}</span>
                </template>
              </td>
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
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { type Batch, useProductionApi } from '@/composables/useProductionApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  // Types
  type Ingredient = {
    id: number
    uuid: string
    name: string
    category: string
    default_unit: string
  }

  type IngredientLot = {
    id: number
    ingredient_id: number
    brewery_lot_code: string | null
    received_amount: number
    received_unit: string
  }

  type StockLocation = {
    id: number
    name: string
  }

  type BeerLot = {
    id: number
    lot_code: string | null
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
    notes: string | null
  }

  type InventoryReceipt = {
    id: number
    uuid: string
    supplier_uuid: string | null
    reference_code: string | null
    received_at: string
    notes: string | null
  }

  type InventoryUsage = {
    id: number
    uuid: string
    production_ref_uuid: string | null
    used_at: string
    notes: string | null
  }

  type InventoryAdjustment = {
    id: number
    uuid: string
    reason: string
    adjusted_at: string
    notes: string | null
  }

  type InventoryTransfer = {
    id: number
    uuid: string
    source_location_id: number
    dest_location_id: number
    transferred_at: string
    notes: string | null
  }

  type Supplier = {
    id: number
    uuid: string
    name: string
  }

  // Composables
  const { request: inventoryRequest } = useInventoryApi()
  const { request: productionRequest } = useProductionApi()
  const { request: procurementRequest } = useProcurementApi()
  const { formatAmountPreferred } = useUnitPreferences()

  // Core data
  const ingredients = ref<Ingredient[]>([])
  const lots = ref<IngredientLot[]>([])
  const locations = ref<StockLocation[]>([])
  const beerLots = ref<BeerLot[]>([])
  const movements = ref<InventoryMovement[]>([])

  // Related data for rich display
  const receipts = ref<InventoryReceipt[]>([])
  const usages = ref<InventoryUsage[]>([])
  const adjustments = ref<InventoryAdjustment[]>([])
  const transfers = ref<InventoryTransfer[]>([])

  // Cross-service data
  const batches = ref<Batch[]>([])
  const suppliers = ref<Supplier[]>([])

  // UI state
  const loading = ref(false)
  const errorMessage = ref('')

  const filters = reactive({
    ingredient_lot_id: null as number | null,
    beer_lot_id: null as number | null,
  })

  // Lookup maps for efficient data resolution
  const ingredientMap = computed(() =>
    new Map(ingredients.value.map(i => [i.id, i])),
  )

  const ingredientLotMap = computed(() =>
    new Map(lots.value.map(l => [l.id, l])),
  )

  const beerLotMap = computed(() =>
    new Map(beerLots.value.map(l => [l.id, l])),
  )

  const locationMap = computed(() =>
    new Map(locations.value.map(l => [l.id, l])),
  )

  const receiptMap = computed(() =>
    new Map(receipts.value.map(r => [r.id, r])),
  )

  const usageMap = computed(() =>
    new Map(usages.value.map(u => [u.id, u])),
  )

  const adjustmentMap = computed(() =>
    new Map(adjustments.value.map(a => [a.id, a])),
  )

  const transferMap = computed(() =>
    new Map(transfers.value.map(t => [t.id, t])),
  )

  const batchByUuidMap = computed(() =>
    new Map(batches.value.map(b => [b.uuid, b])),
  )

  const supplierByUuidMap = computed(() =>
    new Map(suppliers.value.map(s => [s.uuid, s])),
  )

  // Select items for filters
  const lotSelectItems = computed(() =>
    lots.value.map(lot => {
      const ingredient = ingredientMap.value.get(lot.ingredient_id)
      const ingredientName = ingredient?.name ?? 'Unknown'
      const lotCode = lot.brewery_lot_code ?? `Lot ${lot.id}`
      return {
        title: `${ingredientName} - ${lotCode} (${lot.received_amount} ${lot.received_unit})`,
        value: lot.id,
      }
    }),
  )

  const beerLotSelectItems = computed(() =>
    beerLots.value.map(lot => ({
      title: lot.lot_code || `Beer lot ${lot.id}`,
      value: lot.id,
    })),
  )

  // Lifecycle
  onMounted(async () => {
    await refreshAll()
  })

  // Data loading
  async function refreshAll () {
    loading.value = true
    errorMessage.value = ''
    try {
      // Load core inventory data in parallel
      await Promise.all([
        loadIngredients(),
        loadLots(),
        loadLocations(),
        loadBeerLots(),
        loadMovements(),
        loadReceipts(),
        loadUsages(),
        loadAdjustments(),
        loadTransfers(),
      ])

      // Load cross-service data (non-blocking, graceful failure)
      await loadCrossServiceData()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load activity'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }

  async function loadIngredients () {
    ingredients.value = await inventoryRequest<Ingredient[]>('/ingredients')
  }

  async function loadLots () {
    lots.value = await inventoryRequest<IngredientLot[]>('/ingredient-lots')
  }

  async function loadLocations () {
    locations.value = await inventoryRequest<StockLocation[]>('/stock-locations')
  }

  async function loadBeerLots () {
    beerLots.value = await inventoryRequest<BeerLot[]>('/beer-lots')
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
    movements.value = await inventoryRequest<InventoryMovement[]>(path)
  }

  async function loadReceipts () {
    receipts.value = await inventoryRequest<InventoryReceipt[]>('/inventory-receipts')
  }

  async function loadUsages () {
    usages.value = await inventoryRequest<InventoryUsage[]>('/inventory-usage')
  }

  async function loadAdjustments () {
    adjustments.value = await inventoryRequest<InventoryAdjustment[]>('/inventory-adjustments')
  }

  async function loadTransfers () {
    transfers.value = await inventoryRequest<InventoryTransfer[]>('/inventory-transfers')
  }

  async function loadCrossServiceData () {
    // Load batches and suppliers in parallel, but don't fail if they're unavailable
    const results = await Promise.allSettled([
      productionRequest<Batch[]>('/batches'),
      procurementRequest<Supplier[]>('/suppliers'),
    ])

    if (results[0].status === 'fulfilled') {
      batches.value = results[0].value
    }

    if (results[1].status === 'fulfilled') {
      suppliers.value = results[1].value
    }
  }

  // Display helpers
  function getItemName (movement: InventoryMovement): string {
    if (movement.ingredient_lot_id) {
      const lot = ingredientLotMap.value.get(movement.ingredient_lot_id)
      if (lot) {
        const ingredient = ingredientMap.value.get(lot.ingredient_id)
        return ingredient?.name ?? `Ingredient ${lot.ingredient_id}`
      }
      return `Lot ${movement.ingredient_lot_id}`
    }

    if (movement.beer_lot_id) {
      const beerLot = beerLotMap.value.get(movement.beer_lot_id)
      return beerLot?.lot_code ?? `Beer lot ${movement.beer_lot_id}`
    }

    return 'Unknown'
  }

  function getLotCode (movement: InventoryMovement): string {
    if (movement.ingredient_lot_id) {
      const lot = ingredientLotMap.value.get(movement.ingredient_lot_id)
      return lot?.brewery_lot_code ?? `Lot ${movement.ingredient_lot_id}`
    }

    if (movement.beer_lot_id) {
      const beerLot = beerLotMap.value.get(movement.beer_lot_id)
      return beerLot?.lot_code ?? `Beer lot ${movement.beer_lot_id}`
    }

    return 'Unknown'
  }

  function getDirectionTooltip (direction: string): string {
    return direction === 'in' ? 'Received' : 'Used/Transferred/Adjusted'
  }

  function locationName (locationId: number): string {
    return locationMap.value.get(locationId)?.name ?? `Location ${locationId}`
  }

  // Reason-specific helpers
  function getUsageBatch (movement: InventoryMovement): Batch | undefined {
    if (!movement.usage_id) return undefined

    const usage = usageMap.value.get(movement.usage_id)
    if (!usage?.production_ref_uuid) return undefined

    return batchByUuidMap.value.get(usage.production_ref_uuid)
  }

  function getUsagePhase (_movement: InventoryMovement): string | null {
    // Phase information would need to be fetched from batch process phases
    // For now, return null as we don't have this data readily available
    return null
  }

  function getReceiptSupplier (movement: InventoryMovement): Supplier | undefined {
    if (!movement.receipt_id) return undefined

    const receipt = receiptMap.value.get(movement.receipt_id)
    if (!receipt?.supplier_uuid) return undefined

    return supplierByUuidMap.value.get(receipt.supplier_uuid)
  }

  function getAdjustmentNotes (movement: InventoryMovement): string | null {
    if (!movement.adjustment_id) return null

    const adjustment = adjustmentMap.value.get(movement.adjustment_id)
    return adjustment?.notes ?? null
  }

  function formatAdjustmentReason (movement: InventoryMovement): string {
    if (!movement.adjustment_id) {
      return movement.reason === 'waste' ? 'Waste' : 'Adjustment'
    }

    const adjustment = adjustmentMap.value.get(movement.adjustment_id)
    if (!adjustment) {
      return movement.reason === 'waste' ? 'Waste' : 'Adjustment'
    }

    // Format the adjustment reason nicely
    const reasonMap: Record<string, string> = {
      cycle_count: 'Cycle count',
      spoilage: 'Spoilage',
      shrink: 'Shrink',
      damage: 'Damage',
      correction: 'Correction',
      other: 'Other adjustment',
    }

    return reasonMap[adjustment.reason] ?? adjustment.reason
  }

  function getTransferNotes (movement: InventoryMovement): string | null {
    if (!movement.transfer_id) return null

    const transfer = transferMap.value.get(movement.transfer_id)
    return transfer?.notes ?? null
  }

  function formatTransferReason (movement: InventoryMovement): string {
    if (!movement.transfer_id) return 'Transferred'

    const transfer = transferMap.value.get(movement.transfer_id)
    if (!transfer) return 'Transferred'

    const sourceLoc = locationMap.value.get(transfer.source_location_id)
    const destLoc = locationMap.value.get(transfer.dest_location_id)

    if (movement.direction === 'out' && destLoc) {
      return `Transferred to ${destLoc.name}`
    }

    if (movement.direction === 'in' && sourceLoc) {
      return `Transferred from ${sourceLoc.name}`
    }

    return 'Transferred'
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

.reason-link {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
}

.reason-link:hover {
  text-decoration: underline;
}

.adjustment-reason,
.transfer-reason {
  cursor: help;
}
</style>
