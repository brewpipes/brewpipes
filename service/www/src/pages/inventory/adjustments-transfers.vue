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
  import type { BeerLotStockLevel, CreateInventoryAdjustmentRequest, CreateInventoryTransferRequest, IngredientLotStockLevel, StockLocation } from '@/types'
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
    getIngredientLotStockLevels: fetchIngredientLotStockLevels,
    getBeerLotStockLevels: fetchBeerLotStockLevels,
    createInventoryAdjustment,
    createInventoryTransfer,
  } = useInventoryApi()
  const { showNotice } = useSnackbar()
  const { formatAmountPreferred } = useUnitPreferences()

  // Data
  const locations = ref<StockLocation[]>([])
  const ingredientLotStockLevels = ref<IngredientLotStockLevel[]>([])
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

  const allInventory = computed<InventoryItem[]>(() => {
    const items: InventoryItem[] = []

    // Add ingredient lots — one row per lot per location with positive stock
    for (const level of ingredientLotStockLevels.value) {
      items.push({
        key: `ingredient-${level.ingredient_lot_uuid}-${level.stock_location_uuid}`,
        type: 'ingredient',
        lotUuid: level.ingredient_lot_uuid,
        name: level.ingredient_name,
        quantity: level.current_amount,
        unit: level.current_unit,
        locationUuid: level.stock_location_uuid,
        locationName: level.stock_location_name,
      })
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
        loadIngredientLotStockLevels(),
        loadBeerLotStockLevels(),
      ])
    })
  }

  async function loadLocations () {
    locations.value = await getStockLocations()
  }

  async function loadIngredientLotStockLevels () {
    ingredientLotStockLevels.value = await fetchIngredientLotStockLevels()
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
