<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <div class="d-flex align-center">
          <v-icon class="mr-2" icon="mdi-barley" />
          Ingredients
        </div>
        <v-btn
          color="primary"
          :icon="$vuetify.display.xs"
          :loading="ingredientLoading"
          :prepend-icon="$vuetify.display.xs ? undefined : 'mdi-plus'"
          size="small"
          variant="text"
          @click="openIngredientDialog"
        >
          <v-icon v-if="$vuetify.display.xs" icon="mdi-plus" />
          <span v-else>New ingredient</span>
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
        <v-tabs v-model="activeTab" class="inventory-tabs" color="primary" show-arrows>
          <v-tab value="malt">Malt</v-tab>
          <v-tab value="hops">Hops</v-tab>
          <v-tab value="yeast">Yeast</v-tab>
          <v-tab value="other">Other</v-tab>
          <v-tab value="usage">Usage</v-tab>
          <v-tab value="received">Received</v-tab>
        </v-tabs>

        <v-window v-model="activeTab" class="mt-4">
          <!-- Malt Tab -->
          <v-window-item value="malt">
            <LotTable
              empty-text="No malt lots yet."
              :error-message="lotErrorMessage"
              :format-amount-preferred="formatAmountPreferred"
              :format-date-time="formatDateTime"
              :ingredient-name="ingredientName"
              :loading="lotLoading"
              :lots="maltLots"
              title="Malt lots"
              @click:create="openLotDialog('fermentable')"
              @click:refresh="loadLots"
              @click:row="openLotDetails"
            />
          </v-window-item>

          <!-- Hops Tab -->
          <v-window-item value="hops">
            <LotTable
              empty-text="No hop lots yet."
              :error-message="lotErrorMessage"
              :format-amount-preferred="formatAmountPreferred"
              :format-date-time="formatDateTime"
              :ingredient-name="ingredientName"
              :loading="lotLoading"
              :lots="hopLots"
              title="Hop lots"
              @click:create="openLotDialog('hop')"
              @click:refresh="loadLots"
              @click:row="openLotDetails"
            />
          </v-window-item>

          <!-- Yeast Tab -->
          <v-window-item value="yeast">
            <LotTable
              empty-text="No yeast lots yet."
              :error-message="lotErrorMessage"
              :format-amount-preferred="formatAmountPreferred"
              :format-date-time="formatDateTime"
              :ingredient-name="ingredientName"
              :loading="lotLoading"
              :lots="yeastLots"
              title="Yeast lots"
              @click:create="openLotDialog('yeast')"
              @click:refresh="loadLots"
              @click:row="openLotDetails"
            />
          </v-window-item>

          <!-- Other Tab -->
          <v-window-item value="other">
            <LotTable
              empty-text="No other lots yet."
              :error-message="lotErrorMessage"
              :format-amount-preferred="formatAmountPreferred"
              :format-date-time="formatDateTime"
              :ingredient-category="ingredientCategory"
              :ingredient-name="ingredientName"
              :loading="lotLoading"
              :lots="otherLots"
              show-category-column
              title="Other lots"
              @click:create="openLotDialog('other')"
              @click:refresh="loadLots"
              @click:row="openLotDetails"
            />
          </v-window-item>

          <!-- Usage Tab -->
          <v-window-item value="usage">
            <v-card variant="outlined">
              <v-card-title class="d-flex align-center">
                Usage log
                <v-spacer />
                <v-btn :loading="usageLoading" size="small" variant="text" @click="loadUsage">
                  Refresh
                </v-btn>
                <v-btn
                  class="ml-2"
                  color="primary"
                  prepend-icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="openUsageDialog"
                >
                  Log usage
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
                <v-data-table
                  class="data-table"
                  density="compact"
                  :headers="usageHeaders"
                  item-value="uuid"
                  :items="usages"
                  :loading="usageLoading"
                >
                  <template #item.production_ref_uuid="{ item }">
                    {{ item.production_ref_uuid ? (batchMap.get(item.production_ref_uuid)?.short_name ?? item.production_ref_uuid) : 'n/a' }}
                  </template>
                  <template #item.used_at="{ item }">
                    {{ formatDateTime(item.used_at) }}
                  </template>
                  <template #item.notes="{ item }">
                    {{ item.notes || '' }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No usage records yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-window-item>

          <!-- Received Tab -->
          <v-window-item value="received">
            <v-card variant="outlined">
              <v-card-title class="d-flex align-center">
                Inventory receipts
                <v-spacer />
                <v-btn :loading="receiptLoading" size="small" variant="text" @click="loadReceipts">
                  Refresh
                </v-btn>
                <v-btn
                  class="ml-2"
                  color="primary"
                  prepend-icon="mdi-plus"
                  size="small"
                  variant="text"
                  @click="receiveWithoutPODialogOpen = true"
                >
                  Receive inventory
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
                <v-data-table
                  class="data-table"
                  density="compact"
                  :headers="receiptHeaders"
                  item-value="uuid"
                  :items="receipts"
                  :loading="receiptLoading"
                >
                  <template #item.reference_code="{ item }">
                    {{ item.reference_code || 'n/a' }}
                  </template>
                  <template #item.supplier_uuid="{ item }">
                    {{ item.supplier_uuid ? (supplierMap.get(item.supplier_uuid)?.name ?? item.supplier_uuid) : 'n/a' }}
                  </template>
                  <template #item.received_at="{ item }">
                    {{ formatDateTime(item.received_at) }}
                  </template>
                  <template #item.notes="{ item }">
                    {{ item.notes || '' }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No receipts yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-container>

  <!-- Create Lot Dialog -->
  <IngredientLotCreateDialog
    v-model="lotDialog"
    :category="lotDialogCategory"
    :ingredients="ingredients"
    :receipts="receipts"
    :saving="lotSaving"
    :suppliers="suppliers"
    @submit="handleLotSubmit"
  />

  <!-- Create Usage Dialog -->
  <IngredientUsageDialog
    v-model="usageDialog"
    :saving="usageSaving"
    @submit="handleUsageSubmit"
  />

  <!-- Create Receipt Dialog -->
  <IngredientReceiptDialog
    v-model="receiptDialog"
    :saving="receiptSaving"
    @submit="handleReceiptSubmit"
  />

  <!-- Create Ingredient Dialog -->
  <IngredientCreateDialog
    v-model="ingredientDialog"
    :saving="ingredientSaving"
    @submit="handleIngredientSubmit"
  />

  <!-- Receive Without PO Dialog -->
  <ReceiveWithoutPODialog
    v-model="receiveWithoutPODialogOpen"
    @received="handleReceiveWithoutPO"
  />
</template>

<script lang="ts" setup>
  import type { Batch, CreateIngredientLotRequest, CreateIngredientRequest, CreateInventoryReceiptRequest, CreateInventoryUsageRequest, Ingredient, IngredientLot, InventoryReceipt, InventoryUsage, Supplier } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import IngredientCreateDialog from '@/components/inventory/IngredientCreateDialog.vue'
  import IngredientLotCreateDialog from '@/components/inventory/IngredientLotCreateDialog.vue'
  import IngredientReceiptDialog from '@/components/inventory/IngredientReceiptDialog.vue'
  import IngredientUsageDialog from '@/components/inventory/IngredientUsageDialog.vue'
  import LotTable from '@/components/inventory/LotTable.vue'
  import ReceiveWithoutPODialog from '@/components/procurement/ReceiveWithoutPODialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  const {
    getIngredients,
    createIngredient: createIngredientApi,
    getIngredientLots,
    createIngredientLot,
    getInventoryReceipts,
    createInventoryReceipt,
    getInventoryUsages,
    createInventoryUsage,
  } = useInventoryApi()
  const { getBatches } = useProductionApi()
  const { getSuppliers } = useProcurementApi()
  const { showNotice } = useSnackbar()
  const { formatAmountPreferred } = useUnitPreferences()
  const router = useRouter()

  const activeTab = ref('malt')

  // Data
  const ingredients = ref<Ingredient[]>([])
  const receipts = ref<InventoryReceipt[]>([])
  const lots = ref<IngredientLot[]>([])
  const usages = ref<InventoryUsage[]>([])
  const batches = ref<Batch[]>([])
  const suppliers = ref<Supplier[]>([])

  // Async actions for loading
  const { execute: executeIngredientLoad, loading: ingredientLoading, error: ingredientErrorMessage } = useAsyncAction()
  const { execute: executeReceiptLoad, loading: receiptLoading, error: receiptErrorMessage } = useAsyncAction()
  const { execute: executeLotLoad, loading: lotLoading, error: lotErrorMessage } = useAsyncAction()
  const { execute: executeUsageLoad, loading: usageLoading, error: usageErrorMessage } = useAsyncAction()

  // Async actions for saving
  const { execute: executeIngredientSave, loading: ingredientSaving, error: ingredientSaveError } = useAsyncAction()
  const { execute: executeReceiptSave, loading: receiptSaving, error: receiptSaveError } = useAsyncAction()
  const { execute: executeLotSave, loading: lotSaving, error: lotSaveError } = useAsyncAction()
  const { execute: executeUsageSave, loading: usageSaving, error: usageSaveError } = useAsyncAction()

  // Dialog states
  const ingredientDialog = ref(false)
  const receiptDialog = ref(false)
  const lotDialog = ref(false)
  const usageDialog = ref(false)
  const receiveWithoutPODialogOpen = ref(false)

  // Category filter for lot creation
  const lotDialogCategory = ref<string>('fermentable')

  // Table headers
  const usageHeaders = [
    { title: 'Batch', key: 'production_ref_uuid', sortable: true },
    { title: 'Used At', key: 'used_at', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
  ]

  const receiptHeaders = [
    { title: 'Reference Code', key: 'reference_code', sortable: true },
    { title: 'Supplier', key: 'supplier_uuid', sortable: true },
    { title: 'Received At', key: 'received_at', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
  ]

  // Pre-build a lookup map to avoid O(n²) find() calls in lot filters
  const ingredientMap = computed(() => {
    const map = new Map<string, Ingredient>()
    for (const ing of ingredients.value) {
      map.set(ing.uuid, ing)
    }
    return map
  })

  const batchMap = computed(() => new Map(batches.value.map(b => [b.uuid, b])))
  const supplierMap = computed(() => new Map(suppliers.value.map(s => [s.uuid, s])))

  // Computed: Filter lots by category
  const maltLots = computed(() => {
    return lots.value.filter(lot => {
      return ingredientMap.value.get(lot.ingredient_uuid)?.category === 'fermentable'
    })
  })

  const hopLots = computed(() => {
    return lots.value.filter(lot => {
      return ingredientMap.value.get(lot.ingredient_uuid)?.category === 'hop'
    })
  })

  const yeastLots = computed(() => {
    return lots.value.filter(lot => {
      return ingredientMap.value.get(lot.ingredient_uuid)?.category === 'yeast'
    })
  })

  const otherCategories = ['adjunct', 'salt', 'chemical', 'gas', 'other']

  const otherLots = computed(() => {
    return lots.value.filter(lot => {
      const ingredient = ingredientMap.value.get(lot.ingredient_uuid)
      return ingredient && otherCategories.includes(ingredient.category)
    })
  })

  onMounted(async () => {
    await refreshAll()
  })

  async function refreshAll () {
    const [, , , , batchResult, supplierResult] = await Promise.allSettled([
      loadIngredients(),
      loadReceipts(),
      loadLots(),
      loadUsage(),
      getBatches(),
      getSuppliers(),
    ])
    if (batchResult.status === 'fulfilled') {
      batches.value = batchResult.value
    }
    if (supplierResult.status === 'fulfilled') {
      suppliers.value = supplierResult.value
    }
  }

  async function loadIngredients () {
    await executeIngredientLoad(async () => {
      ingredients.value = await getIngredients()
    })
  }

  async function loadReceipts () {
    await executeReceiptLoad(async () => {
      receipts.value = await getInventoryReceipts()
    })
  }

  async function loadLots () {
    await executeLotLoad(async () => {
      lots.value = await getIngredientLots()
    })
  }

  async function loadUsage () {
    await executeUsageLoad(async () => {
      usages.value = await getInventoryUsages()
    })
  }

  // Dialog openers
  function openIngredientDialog () {
    ingredientDialog.value = true
  }

  function openLotDialog (category: string) {
    lotDialogCategory.value = category
    lotDialog.value = true
  }

  function openUsageDialog () {
    usageDialog.value = true
  }

  // Dialog submit handlers
  async function handleIngredientSubmit (payload: CreateIngredientRequest) {
    await executeIngredientSave(async () => {
      await createIngredientApi(payload)
      ingredientDialog.value = false
      await loadIngredients()
      showNotice('Ingredient created')
    })
    if (ingredientSaveError.value) {
      ingredientErrorMessage.value = ingredientSaveError.value
      showNotice(ingredientSaveError.value, 'error')
    }
  }

  async function handleReceiptSubmit (payload: CreateInventoryReceiptRequest) {
    await executeReceiptSave(async () => {
      await createInventoryReceipt(payload)
      receiptDialog.value = false
      await loadReceipts()
      showNotice('Receipt created')
    })
    if (receiptSaveError.value) {
      receiptErrorMessage.value = receiptSaveError.value
      showNotice(receiptSaveError.value, 'error')
    }
  }

  async function handleLotSubmit (payload: CreateIngredientLotRequest) {
    await executeLotSave(async () => {
      await createIngredientLot(payload)
      lotDialog.value = false
      await loadLots()
      showNotice('Lot created')
    })
    if (lotSaveError.value) {
      lotErrorMessage.value = lotSaveError.value
      showNotice(lotSaveError.value, 'error')
    }
  }

  async function handleUsageSubmit (payload: CreateInventoryUsageRequest) {
    await executeUsageSave(async () => {
      await createInventoryUsage(payload)
      usageDialog.value = false
      await loadUsage()
      showNotice('Usage logged')
    })
    if (usageSaveError.value) {
      usageErrorMessage.value = usageSaveError.value
      showNotice(usageSaveError.value, 'error')
    }
  }

  function ingredientName (ingredientUuid: string) {
    return ingredientMap.value.get(ingredientUuid)?.name ?? 'Unknown Ingredient'
  }

  function ingredientCategory (ingredientUuid: string) {
    const category = ingredientMap.value.get(ingredientUuid)?.category
    // Return a display-friendly label
    const labels: Record<string, string> = {
      adjunct: 'Adjunct',
      salt: 'Salt',
      chemical: 'Chemical',
      gas: 'Gas',
      other: 'Other',
    }
    return labels[category ?? ''] ?? category ?? 'Unknown'
  }

  function openLotDetails (lotUuid: string) {
    router.push({
      path: '/inventory/lot-details',
      query: { lot_uuid: lotUuid },
    })
  }

  async function handleReceiveWithoutPO () {
    showNotice('Inventory received successfully', 'success')
    await Promise.allSettled([loadReceipts(), loadLots()])
  }
</script>

<style scoped>
.inventory-page {
  position: relative;
}
</style>
