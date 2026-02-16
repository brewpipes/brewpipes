<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Adjustments & Transfers
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

        <!-- Search/Browse Section -->
        <v-row class="mb-4">
          <v-col cols="12" md="5">
            <v-text-field
              v-model="search"
              append-inner-icon="mdi-magnify"
              clearable
              density="compact"
              hide-details
              label="Search lots by name"
              variant="outlined"
            />
          </v-col>
          <v-col class="d-flex align-center justify-center" cols="12" md="1">
            <span class="text-medium-emphasis">or</span>
          </v-col>
          <v-col cols="12" md="4">
            <v-select
              v-model="selectedLocationUuid"
              clearable
              density="compact"
              hide-details
              :items="locationSelectItems"
              label="Browse by location"
              variant="outlined"
            />
          </v-col>
        </v-row>

        <!-- Inventory Table -->
        <v-data-table
          class="data-table"
          density="compact"
          :headers="headers"
          item-value="key"
          :items="filteredInventory"
          :loading="loading"
        >
          <template #item.type="{ item }">
            <v-chip :color="item.type === 'ingredient' ? 'blue' : 'amber'" size="small" variant="tonal">
              {{ item.type === 'ingredient' ? 'Ingredient' : 'Beer' }}
            </v-chip>
          </template>

          <template #item.quantity="{ item }">
            {{ formatAmountPreferred(item.quantity, item.unit) }}
          </template>

          <template #item.location="{ item }">
            {{ item.locationName }}
          </template>

          <template #item.actions="{ item }">
            <v-btn
              class="mr-1"
              color="primary"
              :disabled="!item.locationUuid"
              size="x-small"
              variant="tonal"
              @click="openAdjustDialog(item)"
            >
              Adjust
            </v-btn>
            <v-btn
              color="secondary"
              :disabled="!item.locationUuid"
              size="x-small"
              variant="tonal"
              @click="openTransferDialog(item)"
            >
              Transfer
            </v-btn>
          </template>

          <template #no-data>
            <div class="text-center py-4">
              <div class="text-body-2 text-medium-emphasis">
                {{ search || selectedLocationUuid ? 'No matching lots found.' : 'Search or select a location to view inventory.' }}
              </div>
            </div>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </v-container>

  <!-- Adjustment Modal -->
  <InventoryAdjustmentDialog
    v-model="adjustDialog"
    :lot="selectedLotInfo"
    :saving="saving"
    @submit="saveAdjustment"
  />

  <!-- Transfer Modal -->
  <InventoryTransferDialog
    v-model="transferDialog"
    :locations="locations"
    :lot="selectedLotInfo"
    :saving="saving"
    @submit="saveTransfer"
  />
</template>

<script lang="ts" setup>
  import type { BeerLotStockLevel, CreateInventoryAdjustmentRequest, CreateInventoryTransferRequest, Ingredient, IngredientLot, InventoryMovement, StockLocation } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import InventoryAdjustmentDialog from '@/components/inventory/InventoryAdjustmentDialog.vue'
  import type { InventoryLotInfo } from '@/components/inventory/InventoryAdjustmentDialog.vue'
  import InventoryTransferDialog from '@/components/inventory/InventoryTransferDialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  type InventoryItem = {
    key: string
    type: 'ingredient' | 'beer'
    lotUuid: string
    name: string
    quantity: number
    unit: string
    locationUuid: string
    locationName: string
  }

  const {
    getStockLocations,
    getIngredients: fetchIngredients,
    getIngredientLots: fetchIngredientLots,
    getInventoryMovements: fetchInventoryMovements,
    getBeerLotStockLevels: fetchBeerLotStockLevels,
    createInventoryAdjustment,
    createInventoryTransfer,
  } = useInventoryApi()
  const { showNotice } = useSnackbar()
  const { formatAmountPreferred } = useUnitPreferences()

  // Data
  const locations = ref<StockLocation[]>([])
  const ingredients = ref<Ingredient[]>([])
  const ingredientLots = ref<IngredientLot[]>([])
  const ingredientMovements = ref<InventoryMovement[]>([])
  const beerLotStockLevels = ref<BeerLotStockLevel[]>([])

  // UI state
  const { execute: executeLoad, loading, error: errorMessage } = useAsyncAction()
  const { execute: executeSave, loading: saving } = useAsyncAction({
    onError: (message) => showNotice(message, 'error'),
  })
  const search = ref('')
  const selectedLocationUuid = ref<string | null>(null)

  // Dialogs
  const adjustDialog = ref(false)
  const transferDialog = ref(false)
  const selectedItem = ref<InventoryItem | null>(null)

  // Table headers
  const headers = [
    { title: 'Type', key: 'type', sortable: true },
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Quantity', key: 'quantity', sortable: true },
    { title: 'Location', key: 'location', sortable: true },
    { title: '', key: 'actions', sortable: false, align: 'end' as const },
  ]

  // Computed
  const locationSelectItems = computed(() =>
    locations.value.map(loc => ({
      title: loc.name,
      value: loc.uuid,
    })),
  )

  /** Build a map of ingredient lot UUID → location UUID → net quantity from movements */
  const ingredientLotLocationMap = computed(() => {
    const map = new Map<string, Map<string, { quantity: number, unit: string }>>()
    for (const movement of ingredientMovements.value) {
      if (!movement.ingredient_lot_uuid) continue
      let lotMap = map.get(movement.ingredient_lot_uuid)
      if (!lotMap) {
        lotMap = new Map()
        map.set(movement.ingredient_lot_uuid, lotMap)
      }
      const existing = lotMap.get(movement.stock_location_uuid) ?? { quantity: 0, unit: movement.amount_unit }
      const delta = movement.direction === 'in' ? movement.amount : -movement.amount
      existing.quantity += delta
      lotMap.set(movement.stock_location_uuid, existing)
    }
    return map
  })

  const allInventory = computed<InventoryItem[]>(() => {
    const items: InventoryItem[] = []
    const locationMap = new Map(locations.value.map(l => [l.uuid, l]))

    // Add ingredient lots — one row per lot per location with positive stock
    for (const lot of ingredientLots.value) {
      const ingredient = ingredients.value.find(i => i.uuid === lot.ingredient_uuid)
      const lotLocations = ingredientLotLocationMap.value.get(lot.uuid)

      if (lotLocations && lotLocations.size > 0) {
        for (const [locUuid, stock] of lotLocations) {
          if (stock.quantity <= 0) continue
          const location = locationMap.get(locUuid)
          items.push({
            key: `ingredient-${lot.uuid}-${locUuid}`,
            type: 'ingredient',
            lotUuid: lot.uuid,
            name: ingredient?.name ?? 'Unknown Ingredient',
            quantity: stock.quantity,
            unit: stock.unit,
            locationUuid: locUuid,
            locationName: location?.name ?? 'Unknown Location',
          })
        }
      } else {
        // No movements found — show lot with unknown location
        items.push({
          key: `ingredient-${lot.uuid}`,
          type: 'ingredient',
          lotUuid: lot.uuid,
          name: ingredient?.name ?? 'Unknown Ingredient',
          quantity: lot.current_amount ?? lot.received_amount,
          unit: lot.current_unit ?? lot.received_unit,
          locationUuid: '',
          locationName: '—',
        })
      }
    }

    // Add beer lot stock levels
    for (const level of beerLotStockLevels.value) {
      items.push({
        key: `beer-${level.beer_lot_uuid}-${level.stock_location_uuid}`,
        type: 'beer',
        lotUuid: level.beer_lot_uuid,
        name: level.lot_code || 'Unknown Beer Lot',
        quantity: level.current_volume,
        unit: level.current_volume_unit,
        locationUuid: level.stock_location_uuid,
        locationName: level.stock_location_name,
      })
    }

    return items
  })

  const filteredInventory = computed(() => {
    let items = allInventory.value

    // Filter by location if selected
    if (selectedLocationUuid.value) {
      items = items.filter(item => item.locationUuid === selectedLocationUuid.value)
    }

    // Filter by search term
    if (search.value) {
      const query = search.value.toLowerCase()
      items = items.filter(item =>
        item.name.toLowerCase().includes(query),
      )
    }

    return items
  })

  const selectedLotInfo = computed<InventoryLotInfo | null>(() => {
    if (!selectedItem.value) return null
    return {
      type: selectedItem.value.type,
      lotUuid: selectedItem.value.lotUuid,
      name: selectedItem.value.name,
      quantity: selectedItem.value.quantity,
      unit: selectedItem.value.unit,
      locationUuid: selectedItem.value.locationUuid,
      locationName: selectedItem.value.locationName,
    }
  })

  // Lifecycle
  onMounted(async () => {
    await refreshAll()
  })

  // Methods
  async function refreshAll () {
    await executeLoad(async () => {
      await Promise.all([
        loadLocations(),
        loadIngredients(),
        loadIngredientLots(),
        loadIngredientMovements(),
        loadBeerLotStockLevels(),
      ])
    })
  }

  async function loadLocations () {
    locations.value = await getStockLocations()
  }

  async function loadIngredients () {
    ingredients.value = await fetchIngredients()
  }

  async function loadIngredientLots () {
    ingredientLots.value = await fetchIngredientLots()
  }

  async function loadIngredientMovements () {
    ingredientMovements.value = await fetchInventoryMovements()
  }

  async function loadBeerLotStockLevels () {
    beerLotStockLevels.value = await fetchBeerLotStockLevels()
  }

  // Adjust dialog
  function openAdjustDialog (item: InventoryItem) {
    selectedItem.value = item
    adjustDialog.value = true
  }

  async function saveAdjustment (payload: CreateInventoryAdjustmentRequest) {
    await executeSave(async () => {
      await createInventoryAdjustment(payload)
      adjustDialog.value = false
      selectedItem.value = null
      await refreshAll()
      showNotice('Adjustment saved')
    })
  }

  // Transfer dialog
  function openTransferDialog (item: InventoryItem) {
    selectedItem.value = item
    transferDialog.value = true
  }

  async function saveTransfer (payload: CreateInventoryTransferRequest) {
    await executeSave(async () => {
      await createInventoryTransfer(payload)
      transferDialog.value = false
      selectedItem.value = null
      await refreshAll()
      showNotice('Transfer completed')
    })
  }
</script>

<style scoped>
.inventory-page {
  position: relative;
}
</style>
