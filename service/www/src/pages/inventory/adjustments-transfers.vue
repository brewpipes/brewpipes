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
              size="x-small"
              variant="tonal"
              @click="openAdjustDialog(item)"
            >
              Adjust
            </v-btn>
            <v-btn
              color="secondary"
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
  <v-dialog v-model="adjustDialog" max-width="500" persistent>
    <v-card>
      <v-card-title class="text-h6">Adjust inventory</v-card-title>
      <v-card-text>
        <div v-if="selectedItem" class="mb-4 pa-3 selected-item-summary rounded">
          <div class="text-subtitle-2 font-weight-bold">{{ selectedItem.name }}</div>
          <div class="text-caption text-medium-emphasis">
            {{ selectedItem.type === 'ingredient' ? 'Ingredient lot' : 'Beer lot' }}
            <span class="mx-1">|</span>
            {{ selectedItem.locationName }}
          </div>
          <div class="text-body-2 mt-1">
            Current quantity: <strong>{{ formatAmountPreferred(selectedItem.quantity, selectedItem.unit) }}</strong>
          </div>
        </div>

        <v-text-field
          v-model="adjustForm.amount"
          density="comfortable"
          hint="Use negative values to decrease inventory"
          label="Adjustment amount"
          persistent-hint
          type="number"
        />
        <v-text-field
          v-model="adjustForm.reason"
          class="mt-2"
          density="comfortable"
          label="Reason"
          placeholder="Damaged, expired, count correction, etc."
          :rules="[rules.required]"
        />
        <v-textarea
          v-model="adjustForm.notes"
          auto-grow
          class="mt-2"
          density="comfortable"
          label="Notes (optional)"
          rows="2"
        />
        <v-text-field
          v-model="adjustForm.adjusted_at"
          class="mt-2"
          density="comfortable"
          label="Adjusted at"
          type="datetime-local"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="closeAdjustDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isAdjustFormValid"
          :loading="saving"
          @click="saveAdjustment"
        >
          Save Adjustment
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Transfer Modal -->
  <v-dialog v-model="transferDialog" max-width="500" persistent>
    <v-card>
      <v-card-title class="text-h6">Transfer inventory</v-card-title>
      <v-card-text>
        <div v-if="selectedItem" class="mb-4 pa-3 selected-item-summary rounded">
          <div class="text-subtitle-2 font-weight-bold">{{ selectedItem.name }}</div>
          <div class="text-caption text-medium-emphasis">
            {{ selectedItem.type === 'ingredient' ? 'Ingredient lot' : 'Beer lot' }}
          </div>
          <div class="text-body-2 mt-1">
            Available: <strong>{{ formatAmountPreferred(selectedItem.quantity, selectedItem.unit) }}</strong>
          </div>
        </div>

        <v-text-field
          v-model="transferForm.from_location"
          density="comfortable"
          disabled
          label="From location"
        />
        <v-select
          v-model="transferForm.to_location_uuid"
          class="mt-2"
          density="comfortable"
          :items="transferDestinationItems"
          label="To location"
          :rules="[rules.required]"
        />
        <v-text-field
          v-model="transferForm.quantity"
          class="mt-2"
          density="comfortable"
          :hint="selectedItem ? `Max: ${selectedItem.quantity} ${selectedItem.unit}` : ''"
          label="Quantity to transfer"
          persistent-hint
          :rules="[rules.required, rules.positiveNumber]"
          type="number"
        />
        <v-textarea
          v-model="transferForm.notes"
          auto-grow
          class="mt-2"
          density="comfortable"
          label="Notes (optional)"
          rows="2"
        />
        <v-text-field
          v-model="transferForm.transferred_at"
          class="mt-2"
          density="comfortable"
          label="Transferred at"
          type="datetime-local"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="closeTransferDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isTransferFormValid"
          :loading="saving"
          @click="saveTransfer"
        >
          Transfer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { BeerLot, Ingredient, IngredientLot, StockLocation } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { normalizeDateTime, normalizeText } from '@/utils/normalize'

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
    getBeerLots: fetchBeerLots,
    createInventoryAdjustment,
    createInventoryTransfer,
  } = useInventoryApi()
  const { showNotice } = useSnackbar()
  const { formatAmountPreferred } = useUnitPreferences()

  // Data
  const locations = ref<StockLocation[]>([])
  const ingredients = ref<Ingredient[]>([])
  const ingredientLots = ref<IngredientLot[]>([])
  const beerLots = ref<BeerLot[]>([])

  // UI state
  const loading = ref(false)
  const saving = ref(false)
  const errorMessage = ref('')
  const search = ref('')
  const selectedLocationUuid = ref<string | null>(null)

  // Dialogs
  const adjustDialog = ref(false)
  const transferDialog = ref(false)
  const selectedItem = ref<InventoryItem | null>(null)

  // Forms
  const adjustForm = reactive({
    amount: '',
    reason: '',
    notes: '',
    adjusted_at: '',
  })

  const transferForm = reactive({
    from_location: '',
    to_location_uuid: null as string | null,
    quantity: '',
    notes: '',
    transferred_at: '',
  })

  const rules = {
    required: (v: string | null) => (v !== null && v !== '' && String(v).trim() !== '') || 'Required',
    positiveNumber: (v: string | null) => {
      if (v === null || v === '') return true // Let required handle empty
      const num = Number(v)
      return (Number.isFinite(num) && num > 0) || 'Must be a positive number'
    },
  }

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

  const transferDestinationItems = computed(() => {
    if (!selectedItem.value) return locationSelectItems.value
    return locations.value
      .filter(loc => loc.uuid !== selectedItem.value?.locationUuid)
      .map(loc => ({
        title: loc.name,
        value: loc.uuid,
      }))
  })

  const allInventory = computed<InventoryItem[]>(() => {
    const items: InventoryItem[] = []

    // Add ingredient lots
    for (const lot of ingredientLots.value) {
      const ingredient = ingredients.value.find(i => i.uuid === lot.ingredient_uuid)
      const location = locations.value.find(l => l.uuid === lot.stock_location_uuid)
      items.push({
        key: `ingredient-${lot.uuid}`,
        type: 'ingredient',
        lotUuid: lot.uuid,
        name: ingredient?.name ?? 'Unknown Ingredient',
        quantity: lot.current_amount ?? lot.received_amount,
        unit: lot.current_unit ?? lot.received_unit,
        locationUuid: lot.stock_location_uuid,
        locationName: location?.name ?? 'Unknown Location',
      })
    }

    // Add beer lots
    for (const lot of beerLots.value) {
      const location = locations.value.find(l => l.uuid === lot.stock_location_uuid)
      items.push({
        key: `beer-${lot.uuid}`,
        type: 'beer',
        lotUuid: lot.uuid,
        name: lot.lot_code || 'Unknown Beer Lot',
        quantity: lot.volume,
        unit: lot.volume_unit,
        locationUuid: lot.stock_location_uuid,
        locationName: location?.name ?? 'Unknown Location',
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

  const isAdjustFormValid = computed(() => {
    return adjustForm.amount !== '' && adjustForm.reason.trim().length > 0
  })

  const isTransferFormValid = computed(() => {
    if (transferForm.to_location_uuid === null || transferForm.quantity === '') {
      return false
    }
    const qty = Number(transferForm.quantity)
    if (!Number.isFinite(qty) || qty <= 0) {
      return false
    }
    // Ensure transfer quantity doesn't exceed available
    if (selectedItem.value && qty > selectedItem.value.quantity) {
      return false
    }
    return true
  })

  // Lifecycle
  onMounted(async () => {
    await refreshAll()
  })

  // Methods
  function getDefaultDateTime () {
    const now = new Date()
    const offset = now.getTimezoneOffset()
    const local = new Date(now.getTime() - offset * 60 * 1000)
    return local.toISOString().slice(0, 16)
  }

  async function refreshAll () {
    loading.value = true
    errorMessage.value = ''
    try {
      await Promise.all([
        loadLocations(),
        loadIngredients(),
        loadIngredientLots(),
        loadBeerLots(),
      ])
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load data'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
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

  async function loadBeerLots () {
    beerLots.value = await fetchBeerLots()
  }

  // Adjust dialog
  function openAdjustDialog (item: InventoryItem) {
    selectedItem.value = item
    adjustForm.amount = ''
    adjustForm.reason = ''
    adjustForm.notes = ''
    adjustForm.adjusted_at = getDefaultDateTime()
    adjustDialog.value = true
  }

  function closeAdjustDialog () {
    adjustDialog.value = false
    selectedItem.value = null
  }

  async function saveAdjustment () {
    if (!isAdjustFormValid.value || !selectedItem.value) {
      return
    }

    saving.value = true
    errorMessage.value = ''

    try {
      const payload = {
        ingredient_lot_uuid: selectedItem.value.type === 'ingredient' ? selectedItem.value.lotUuid : null,
        beer_lot_uuid: selectedItem.value.type === 'beer' ? selectedItem.value.lotUuid : null,
        stock_location_uuid: selectedItem.value.locationUuid,
        amount: Number(adjustForm.amount),
        amount_unit: selectedItem.value.unit,
        reason: adjustForm.reason.trim(),
        notes: normalizeText(adjustForm.notes),
        adjusted_at: normalizeDateTime(adjustForm.adjusted_at),
      }
      await createInventoryAdjustment(payload)
      closeAdjustDialog()
      await refreshAll()
      showNotice('Adjustment saved')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save adjustment'
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }

  // Transfer dialog
  function openTransferDialog (item: InventoryItem) {
    selectedItem.value = item
    transferForm.from_location = item.locationName
    transferForm.to_location_uuid = null
    transferForm.quantity = ''
    transferForm.notes = ''
    transferForm.transferred_at = getDefaultDateTime()
    transferDialog.value = true
  }

  function closeTransferDialog () {
    transferDialog.value = false
    selectedItem.value = null
  }

  async function saveTransfer () {
    if (!isTransferFormValid.value || !selectedItem.value) {
      return
    }

    saving.value = true
    errorMessage.value = ''

    try {
      const payload = {
        ingredient_lot_uuid: selectedItem.value.type === 'ingredient' ? selectedItem.value.lotUuid : null,
        beer_lot_uuid: selectedItem.value.type === 'beer' ? selectedItem.value.lotUuid : null,
        source_location_uuid: selectedItem.value.locationUuid,
        dest_location_uuid: transferForm.to_location_uuid,
        quantity: Number(transferForm.quantity),
        quantity_unit: selectedItem.value.unit,
        notes: normalizeText(transferForm.notes),
        transferred_at: normalizeDateTime(transferForm.transferred_at),
      }
      await createInventoryTransfer(payload)
      closeTransferDialog()
      await refreshAll()
      showNotice('Transfer completed')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to complete transfer'
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
.inventory-page {
  position: relative;
}

.selected-item-summary {
  background: rgba(var(--v-theme-on-surface), 0.05);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
}
</style>
