<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <div class="d-flex align-center">
          Inventory activity
        </div>
        <div class="card-title-actions">
          <v-btn :loading="loading" size="small" variant="text" @click="refreshAll">
            <v-icon class="mr-1" icon="mdi-refresh" size="small" />
            <span class="d-none d-sm-inline">Refresh</span>
          </v-btn>
        </div>
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
              v-model="filters.ingredient_lot_uuid"
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
              v-model="filters.beer_lot_uuid"
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

        <!-- Desktop table view -->
        <div class="d-none d-md-block">
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
              <tr v-for="movement in movements" v-else :key="movement.uuid">
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
                  <MovementReason
                    :adjustment-notes="getAdjustmentNotes(movement)"
                    :adjustment-reason-label="formatAdjustmentReason(movement)"
                    :movement="movement"
                    :receipt-supplier="getReceiptSupplier(movement)"
                    :transfer-notes="getTransferNotes(movement)"
                    :transfer-reason-label="formatTransferReason(movement)"
                    :usage-batch="getUsageBatch(movement)"
                  />
                </td>
                <td>{{ formatAmountPreferred(movement.amount, movement.amount_unit) }}</td>
                <td>{{ locationName(movement.stock_location_uuid) }}</td>
              </tr>
            </tbody>
          </v-table>
        </div>

        <!-- Mobile card view -->
        <div class="d-md-none">
          <v-progress-linear v-if="loading" color="primary" indeterminate />

          <div
            v-if="!loading && movements.length === 0"
            class="text-center py-8 text-medium-emphasis"
          >
            <v-icon class="mb-2" icon="mdi-swap-horizontal" size="48" />
            <div class="text-body-1">No activity yet.</div>
          </div>

          <v-card
            v-for="movement in movements"
            :key="movement.uuid"
            class="mb-3"
            variant="outlined"
          >
            <v-card-title class="d-flex align-center py-2 text-body-1">
              <v-icon
                class="mr-2"
                :color="movement.direction === 'in' ? 'success' : 'warning'"
                :icon="movement.direction === 'in' ? 'mdi-package-down' : 'mdi-package-up'"
                size="small"
              />
              <span class="text-truncate">{{ getItemName(movement) }}</span>
              <v-spacer />
              <span class="text-body-2 font-weight-medium text-no-wrap">
                {{ formatAmountPreferred(movement.amount, movement.amount_unit) }}
              </span>
            </v-card-title>

            <v-card-text class="pt-0">
              <div class="d-flex justify-space-between text-body-2 mb-1">
                <span class="text-medium-emphasis">Direction</span>
                <span>{{ movement.direction === 'in' ? 'In' : 'Out' }}</span>
              </div>
              <div class="d-flex justify-space-between text-body-2 mb-1">
                <span class="text-medium-emphasis">Reason</span>
                <span class="text-right" style="max-width: 65%;">
                  <MovementReason
                    :adjustment-notes="getAdjustmentNotes(movement)"
                    :adjustment-reason-label="formatAdjustmentReason(movement)"
                    :movement="movement"
                    :receipt-supplier="getReceiptSupplier(movement)"
                    :transfer-notes="getTransferNotes(movement)"
                    :transfer-reason-label="formatTransferReason(movement)"
                    :usage-batch="getUsageBatch(movement)"
                  />
                </span>
              </div>
              <div class="d-flex justify-space-between text-body-2 mb-1">
                <span class="text-medium-emphasis">Lot #</span>
                <span>{{ getLotCode(movement) }}</span>
              </div>
              <div class="d-flex justify-space-between text-body-2">
                <span class="text-medium-emphasis">Location</span>
                <span>{{ locationName(movement.stock_location_uuid) }}</span>
              </div>
            </v-card-text>
          </v-card>
        </div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script lang="ts" setup>
  import type {
    Batch,
    BeerLot,
    Ingredient,
    IngredientLot,
    InventoryAdjustment,
    InventoryMovement,
    InventoryReceipt,
    InventoryTransfer,
    InventoryUsage,
    StockLocation,
    Supplier,
  } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import MovementReason from '@/components/inventory/MovementReason.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  // Composables
  const {
    getIngredients: fetchIngredients,
    getIngredientLots: fetchIngredientLots,
    getStockLocations: fetchStockLocations,
    getBeerLots: fetchBeerLots,
    getInventoryMovements,
    getInventoryReceipts: fetchReceipts,
    getInventoryUsages: fetchUsages,
    getInventoryAdjustments: fetchAdjustments,
    getInventoryTransfers: fetchTransfers,
  } = useInventoryApi()
  const { getBatches: fetchBatches } = useProductionApi()
  const { getSuppliers: fetchSuppliers } = useProcurementApi()
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
  const { execute, loading, error: errorMessage } = useAsyncAction()

  const filters = reactive({
    ingredient_lot_uuid: null as string | null,
    beer_lot_uuid: null as string | null,
  })

  // Lookup maps for efficient data resolution
  const ingredientMap = computed(() =>
    new Map(ingredients.value.map(i => [i.uuid, i])),
  )

  const ingredientLotMap = computed(() =>
    new Map(lots.value.map(l => [l.uuid, l])),
  )

  const beerLotMap = computed(() =>
    new Map(beerLots.value.map(l => [l.uuid, l])),
  )

  const locationMap = computed(() =>
    new Map(locations.value.map(l => [l.uuid, l])),
  )

  const receiptMap = computed(() =>
    new Map(receipts.value.map(r => [r.uuid, r])),
  )

  const usageMap = computed(() =>
    new Map(usages.value.map(u => [u.uuid, u])),
  )

  const adjustmentMap = computed(() =>
    new Map(adjustments.value.map(a => [a.uuid, a])),
  )

  const transferMap = computed(() =>
    new Map(transfers.value.map(t => [t.uuid, t])),
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
      const ingredient = ingredientMap.value.get(lot.ingredient_uuid)
      const ingredientName = ingredient?.name ?? 'Unknown'
      const lotCode = lot.brewery_lot_code ?? 'Unknown Lot'
      return {
        title: `${ingredientName} - ${lotCode} (${lot.received_amount} ${lot.received_unit})`,
        value: lot.uuid,
      }
    }),
  )

  const beerLotSelectItems = computed(() =>
    beerLots.value.map(lot => ({
      title: lot.lot_code || 'Unknown Beer Lot',
      value: lot.uuid,
    })),
  )

  // Lifecycle
  onMounted(async () => {
    await refreshAll()
  })

  // Data loading
  async function refreshAll () {
    await execute(async () => {
      // Load core inventory data in parallel (allSettled so one failure doesn't kill all data)
      const results = await Promise.allSettled([
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

      const failures = results.filter(r => r.status === 'rejected')
      if (failures.length > 0 && failures.length < results.length) {
        throw new Error(`Some data failed to load (${failures.length} of ${results.length} requests). Displayed data may be incomplete.`)
      } else if (failures.length === results.length) {
        throw new Error('Unable to load activity')
      }

      // Load cross-service data (non-blocking, graceful failure)
      await loadCrossServiceData()
    })
  }

  async function loadIngredients () {
    ingredients.value = await fetchIngredients() ?? []
  }

  async function loadLots () {
    lots.value = await fetchIngredientLots() ?? []
  }

  async function loadLocations () {
    locations.value = await fetchStockLocations() ?? []
  }

  async function loadBeerLots () {
    beerLots.value = await fetchBeerLots() ?? []
  }

  async function loadMovements () {
    const movementFilters: { ingredient_lot_uuid?: string, beer_lot_uuid?: string } = {}
    if (filters.ingredient_lot_uuid) {
      movementFilters.ingredient_lot_uuid = filters.ingredient_lot_uuid
    }
    if (filters.beer_lot_uuid) {
      movementFilters.beer_lot_uuid = filters.beer_lot_uuid
    }
    movements.value = await getInventoryMovements(movementFilters) ?? []
  }

  async function loadReceipts () {
    receipts.value = await fetchReceipts() ?? []
  }

  async function loadUsages () {
    usages.value = await fetchUsages() ?? []
  }

  async function loadAdjustments () {
    adjustments.value = await fetchAdjustments() ?? []
  }

  async function loadTransfers () {
    transfers.value = await fetchTransfers() ?? []
  }

  async function loadCrossServiceData () {
    // Load batches and suppliers in parallel, but don't fail if they're unavailable
    const results = await Promise.allSettled([
      fetchBatches(),
      fetchSuppliers(),
    ])

    if (results[0].status === 'fulfilled') {
      batches.value = results[0].value ?? []
    }

    if (results[1].status === 'fulfilled') {
      suppliers.value = results[1].value ?? []
    }
  }

  // Display helpers
  function getItemName (movement: InventoryMovement): string {
    if (movement.ingredient_lot_uuid) {
      const lot = ingredientLotMap.value.get(movement.ingredient_lot_uuid)
      if (lot) {
        const ingredient = ingredientMap.value.get(lot.ingredient_uuid)
        return ingredient?.name ?? 'Unknown Ingredient'
      }
      return 'Unknown Lot'
    }

    if (movement.beer_lot_uuid) {
      const beerLot = beerLotMap.value.get(movement.beer_lot_uuid)
      return beerLot?.lot_code ?? 'Unknown Beer Lot'
    }

    return 'Unknown'
  }

  function getLotCode (movement: InventoryMovement): string {
    if (movement.ingredient_lot_uuid) {
      const lot = ingredientLotMap.value.get(movement.ingredient_lot_uuid)
      return lot?.brewery_lot_code ?? 'Unknown Lot'
    }

    if (movement.beer_lot_uuid) {
      const beerLot = beerLotMap.value.get(movement.beer_lot_uuid)
      return beerLot?.lot_code ?? 'Unknown Beer Lot'
    }

    return 'Unknown'
  }

  function getDirectionTooltip (direction: string): string {
    return direction === 'in' ? 'Received' : 'Used/Transferred/Adjusted'
  }

  function locationName (locationUuid: string): string {
    return locationMap.value.get(locationUuid)?.name ?? 'Unknown Location'
  }

  // Reason-specific helpers
  function getUsageBatch (movement: InventoryMovement): Batch | undefined {
    if (!movement.usage_uuid) return undefined

    const usage = usageMap.value.get(movement.usage_uuid)
    if (!usage?.production_ref_uuid) return undefined

    return batchByUuidMap.value.get(usage.production_ref_uuid)
  }

  function getReceiptSupplier (movement: InventoryMovement): Supplier | undefined {
    if (!movement.receipt_uuid) return undefined

    const receipt = receiptMap.value.get(movement.receipt_uuid)
    if (!receipt?.supplier_uuid) return undefined

    return supplierByUuidMap.value.get(receipt.supplier_uuid)
  }

  function getAdjustmentNotes (movement: InventoryMovement): string | null {
    if (!movement.adjustment_uuid) return null

    const adjustment = adjustmentMap.value.get(movement.adjustment_uuid)
    return adjustment?.notes ?? null
  }

  function formatAdjustmentReason (movement: InventoryMovement): string {
    if (!movement.adjustment_uuid) {
      return movement.reason === 'waste' ? 'Waste' : 'Adjustment'
    }

    const adjustment = adjustmentMap.value.get(movement.adjustment_uuid)
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
    if (!movement.transfer_uuid) return null

    const transfer = transferMap.value.get(movement.transfer_uuid)
    return transfer?.notes ?? null
  }

  function formatTransferReason (movement: InventoryMovement): string {
    if (!movement.transfer_uuid) return 'Transferred'

    const transfer = transferMap.value.get(movement.transfer_uuid)
    if (!transfer) return 'Transferred'

    const sourceLoc = locationMap.value.get(transfer.source_location_uuid)
    const destLoc = locationMap.value.get(transfer.dest_location_uuid)

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
</style>
