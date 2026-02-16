<template>
  <v-container class="product-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <div class="d-flex align-center">
          <v-icon class="mr-2" icon="mdi-package-variant-closed" />
          Finished Goods
        </div>
        <div class="card-title-actions">
          <v-text-field
            v-model="searchQuery"
            class="search-field"
            clearable
            density="compact"
            hide-details
            placeholder="Search..."
            prepend-inner-icon="mdi-magnify"
            variant="outlined"
          />
          <v-btn
            :loading="loading"
            size="small"
            variant="text"
            @click="refreshAll"
          >
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

        <v-tabs v-model="activeTab" class="product-tabs" color="primary" show-arrows>
          <v-tab value="stock-levels">Stock Levels</v-tab>
          <v-tab value="all-lots">All Lots</v-tab>
        </v-tabs>

        <v-window v-model="activeTab" class="mt-4">
          <!-- Stock Levels Tab -->
          <v-window-item value="stock-levels">
            <!-- Container type filter -->
            <v-chip-group
              v-model="containerFilter"
              class="mb-4"
              color="primary"
              mandatory
              selected-class="text-primary"
            >
              <v-chip filter value="all">All</v-chip>
              <v-chip filter value="keg">
                <v-icon class="mr-1" icon="mdi-keg" size="small" />
                Kegs
              </v-chip>
              <v-chip filter value="can">
                <v-icon class="mr-1" icon="mdi-cup" size="small" />
                Cans
              </v-chip>
              <v-chip filter value="bottle">
                <v-icon class="mr-1" icon="mdi-bottle-wine" size="small" />
                Bottles
              </v-chip>
              <v-chip filter value="other">
                <v-icon class="mr-1" icon="mdi-package-variant" size="small" />
                Other
              </v-chip>
            </v-chip-group>

            <!-- Location filter -->
            <v-row v-if="locationOptions.length > 1" class="mb-4">
              <v-col cols="12" md="4" sm="6">
                <v-select
                  v-model="locationFilter"
                  clearable
                  density="compact"
                  hide-details
                  :items="locationOptions"
                  label="Filter by location"
                  variant="outlined"
                />
              </v-col>
            </v-row>

            <!-- Desktop table view -->
            <div class="d-none d-md-block">
              <v-data-table
                class="data-table stock-level-table"
                density="compact"
                :headers="stockLevelHeaders"
                hover
                item-value="beer_lot_uuid"
                :items="filteredStockLevels"
                :loading="loading"
                :search="searchQuery"
              >
                <template #item.batch_name="{ item }">
                  <router-link
                    v-if="batchMap.get(item.production_batch_uuid)"
                    class="batch-link"
                    :to="`/batches/${item.production_batch_uuid}`"
                  >
                    {{ batchName(item.production_batch_uuid) }}
                  </router-link>
                  <span v-else class="text-medium-emphasis">
                    {{ batchName(item.production_batch_uuid) }}
                  </span>
                </template>

                <template #item.lot_code="{ item }">
                  {{ item.lot_code || '---' }}
                </template>

                <template #item.package_format_name="{ item }">
                  <div class="d-flex align-center">
                    <v-icon
                      class="mr-1"
                      :icon="containerIcon(item.container)"
                      size="small"
                    />
                    {{ item.package_format_name || '---' }}
                  </div>
                </template>

                <template #item.current_quantity="{ item }">
                  {{ formatQuantity(item) }}
                </template>

                <template #item.stock_location_name="{ item }">
                  {{ item.stock_location_name }}
                </template>

                <template #item.packaged_at="{ item }">
                  {{ formatDate(item.packaged_at) }}
                </template>

                <template #item.best_by="{ item }">
                  <BestByIndicator :best-by="item.best_by" />
                </template>

                <template #no-data>
                  <div class="text-center py-8 text-medium-emphasis">
                    <v-icon class="mb-2" icon="mdi-package-variant-closed" size="48" />
                    <div class="text-body-1">No packaged goods in inventory.</div>
                    <div class="text-body-2 mt-1">
                      Package a batch to see finished goods here.
                    </div>
                  </div>
                </template>
              </v-data-table>
            </div>

            <!-- Mobile card view -->
            <div class="d-md-none">
              <v-progress-linear v-if="loading" color="primary" indeterminate />

              <div
                v-if="!loading && filteredStockLevels.length === 0"
                class="text-center py-8 text-medium-emphasis"
              >
                <v-icon class="mb-2" icon="mdi-package-variant-closed" size="48" />
                <div class="text-body-1">No packaged goods in inventory.</div>
                <div class="text-body-2 mt-1">
                  Package a batch to see finished goods here.
                </div>
              </div>

              <v-card
                v-for="item in filteredStockLevelsMobile"
                :key="`${item.beer_lot_uuid}-${item.stock_location_uuid}`"
                class="mb-3"
                variant="outlined"
              >
                <v-card-title class="d-flex align-center py-2 text-body-1">
                  <v-icon
                    class="mr-2"
                    :icon="containerIcon(item.container)"
                    size="small"
                  />
                  <router-link
                    v-if="batchMap.get(item.production_batch_uuid)"
                    class="batch-link"
                    :to="`/batches/${item.production_batch_uuid}`"
                  >
                    {{ batchName(item.production_batch_uuid) }}
                  </router-link>
                  <span v-else class="text-medium-emphasis">
                    {{ batchName(item.production_batch_uuid) }}
                  </span>
                  <v-spacer />
                  <BestByIndicator :best-by="item.best_by" chip-only />
                </v-card-title>

                <v-card-text class="pt-0">
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Lot code</span>
                    <span>{{ item.lot_code || '---' }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Format</span>
                    <span>{{ item.package_format_name || '---' }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Quantity</span>
                    <span class="font-weight-medium">{{ formatQuantity(item) }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Location</span>
                    <span>{{ item.stock_location_name }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Packaged</span>
                    <span>{{ formatDate(item.packaged_at) }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2">
                    <span class="text-medium-emphasis">Best by</span>
                    <BestByIndicator :best-by="item.best_by" />
                  </div>
                </v-card-text>
              </v-card>
            </div>
          </v-window-item>

          <!-- All Lots Tab -->
          <v-window-item value="all-lots">
            <!-- Desktop table view -->
            <div class="d-none d-md-block">
              <v-data-table
                class="data-table"
                density="compact"
                :headers="allLotsHeaders"
                hover
                item-value="uuid"
                :items="beerLots"
                :loading="loading"
                :search="searchQuery"
              >
                <template #item.batch_name="{ item }">
                  <router-link
                    v-if="batchMap.get(item.production_batch_uuid)"
                    class="batch-link"
                    :to="`/batches/${item.production_batch_uuid}`"
                  >
                    {{ batchName(item.production_batch_uuid) }}
                  </router-link>
                  <span v-else class="text-medium-emphasis">
                    {{ batchName(item.production_batch_uuid) }}
                  </span>
                </template>

                <template #item.lot_code="{ item }">
                  {{ item.lot_code || '---' }}
                </template>

                <template #item.package_format_name="{ item }">
                  <div v-if="item.package_format_name" class="d-flex align-center">
                    <v-icon
                      class="mr-1"
                      :icon="containerIcon(item.container)"
                      size="small"
                    />
                    {{ item.package_format_name }}
                  </div>
                  <span v-else class="text-medium-emphasis">---</span>
                </template>

                <template #item.quantity="{ item }">
                  {{ item.quantity != null ? item.quantity : '---' }}
                </template>

                <template #item.packaged_at="{ item }">
                  {{ formatDate(item.packaged_at) }}
                </template>

                <template #item.best_by="{ item }">
                  <BestByIndicator :best-by="item.best_by" />
                </template>

                <template #item.notes="{ item }">
                  <span class="text-medium-emphasis">{{ item.notes || '' }}</span>
                </template>

                <template #no-data>
                  <div class="text-center py-4 text-medium-emphasis">
                    No beer lots yet.
                  </div>
                </template>
              </v-data-table>
            </div>

            <!-- Mobile card view -->
            <div class="d-md-none">
              <v-progress-linear v-if="loading" color="primary" indeterminate />

              <div
                v-if="!loading && beerLots.length === 0"
                class="text-center py-4 text-medium-emphasis"
              >
                No beer lots yet.
              </div>

              <v-card
                v-for="lot in beerLots"
                :key="lot.uuid"
                class="mb-3"
                variant="outlined"
              >
                <v-card-title class="d-flex align-center py-2 text-body-1">
                  <v-icon
                    v-if="lot.container"
                    class="mr-2"
                    :icon="containerIcon(lot.container)"
                    size="small"
                  />
                  <router-link
                    v-if="batchMap.get(lot.production_batch_uuid)"
                    class="batch-link"
                    :to="`/batches/${lot.production_batch_uuid}`"
                  >
                    {{ batchName(lot.production_batch_uuid) }}
                  </router-link>
                  <span v-else class="text-medium-emphasis">
                    {{ batchName(lot.production_batch_uuid) }}
                  </span>
                  <v-spacer />
                  <BestByIndicator :best-by="lot.best_by" chip-only />
                </v-card-title>

                <v-card-text class="pt-0">
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Lot code</span>
                    <span>{{ lot.lot_code || '---' }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Format</span>
                    <span>{{ lot.package_format_name || '---' }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Quantity</span>
                    <span>{{ lot.quantity != null ? lot.quantity : '---' }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2 mb-1">
                    <span class="text-medium-emphasis">Packaged</span>
                    <span>{{ formatDate(lot.packaged_at) }}</span>
                  </div>
                  <div class="d-flex justify-space-between text-body-2">
                    <span class="text-medium-emphasis">Best by</span>
                    <BestByIndicator :best-by="lot.best_by" />
                  </div>
                  <div v-if="lot.notes" class="text-body-2 text-medium-emphasis mt-2">
                    {{ lot.notes }}
                  </div>
                </v-card-text>
              </v-card>
            </div>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script lang="ts" setup>
  import type { Batch, BeerLot, BeerLotStockLevel } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import BestByIndicator from '@/components/inventory/BestByIndicator.vue'
  import { formatDate } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { convertVolume, volumeLabels } from '@/composables/useUnitConversion'
  import { isVolumeUnit, normalizeVolumeUnit, useUnitPreferences } from '@/composables/useUnitPreferences'

  const { getBeerLots: fetchBeerLots, getBeerLotStockLevels } = useInventoryApi()
  const { getBatches } = useProductionApi()
  const { preferences } = useUnitPreferences()

  // Data
  const beerLots = ref<BeerLot[]>([])
  const stockLevels = ref<BeerLotStockLevel[]>([])
  const batches = ref<Batch[]>([])

  // UI state
  const loading = ref(false)
  const errorMessage = ref('')
  const activeTab = ref('stock-levels')
  const searchQuery = ref('')
  const containerFilter = ref('all')
  const locationFilter = ref<string | null>(null)

  // Lookup maps
  const batchMap = computed(() =>
    new Map(batches.value.map(b => [b.uuid, b])),
  )

  // Container type icon mapping
  const CONTAINER_ICONS: Record<string, string> = {
    keg: 'mdi-keg',
    can: 'mdi-cup',
    bottle: 'mdi-bottle-wine',
    cask: 'mdi-barrel',
    growler: 'mdi-bottle-tonic',
    other: 'mdi-package-variant',
  }

  const DEFAULT_CONTAINER_ICON = 'mdi-package-variant'

  function containerIcon (container: string | undefined | null): string {
    if (!container) return DEFAULT_CONTAINER_ICON
    return CONTAINER_ICONS[container] ?? DEFAULT_CONTAINER_ICON
  }

  function batchName (uuid: string): string {
    return batchMap.value.get(uuid)?.short_name ?? uuid.slice(0, 8)
  }

  // Volume formatting for quantity display
  function formatQuantity (item: BeerLotStockLevel): string {
    const qty = item.current_quantity
    const container = item.container

    // Count-based display
    if (qty != null && qty > 0 && container) {
      const label = containerLabel(container, qty)
      const volumeStr = formatVolumeInPreferred(item.current_volume, item.current_volume_unit)
      if (volumeStr) {
        return `${qty} ${label} (${volumeStr})`
      }
      return `${qty} ${label}`
    }

    // Volume-only fallback
    const volumeStr = formatVolumeInPreferred(item.current_volume, item.current_volume_unit)
    return volumeStr || '---'
  }

  function containerLabel (container: string, count: number): string {
    const labels: Record<string, [string, string]> = {
      keg: ['keg', 'kegs'],
      can: ['can', 'cans'],
      bottle: ['bottle', 'bottles'],
      cask: ['cask', 'casks'],
      growler: ['growler', 'growlers'],
      other: ['unit', 'units'],
    }
    const pair = labels[container]
    if (pair) {
      return count === 1 ? pair[0] : pair[1]
    }
    return count === 1 ? 'unit' : 'units'
  }

  function formatVolumeInPreferred (volume: number, unit: string): string | null {
    if (volume == null || !unit) return null
    const preferredUnit = preferences.value.volume
    if (!isVolumeUnit(unit)) return `${volume} ${unit}`
    const normalizedUnit = normalizeVolumeUnit(unit)
    const converted = convertVolume(volume, normalizedUnit, preferredUnit)
    if (converted === null) return `${volume} ${unit}`
    const label = volumeLabels[preferredUnit] ?? preferredUnit
    return `${converted.toFixed(2)} ${label}`
  }

  // Location options for filter
  const locationOptions = computed(() => {
    const locations = new Map<string, string>()
    for (const level of stockLevels.value) {
      locations.set(level.stock_location_uuid, level.stock_location_name)
    }
    return Array.from(locations.entries()).map(([uuid, name]) => ({
      title: name,
      value: uuid,
    }))
  })

  // Filtered stock levels
  const filteredStockLevels = computed(() => {
    let items = stockLevels.value

    // Container filter
    if (containerFilter.value !== 'all') {
      if (containerFilter.value === 'other') {
        const coreContainers = new Set(['keg', 'can', 'bottle'])
        items = items.filter(item => !item.container || !coreContainers.has(item.container))
      } else {
        items = items.filter(item => item.container === containerFilter.value)
      }
    }

    // Location filter
    if (locationFilter.value) {
      items = items.filter(item => item.stock_location_uuid === locationFilter.value)
    }

    return items
  })

  // Mobile uses the same filtered list (search handled by v-data-table on desktop,
  // but for mobile cards we apply search manually)
  const filteredStockLevelsMobile = computed(() => {
    let items = filteredStockLevels.value

    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      items = items.filter(item => {
        const batch = batchName(item.production_batch_uuid).toLowerCase()
        const lotCode = (item.lot_code ?? '').toLowerCase()
        const format = (item.package_format_name ?? '').toLowerCase()
        const location = item.stock_location_name.toLowerCase()
        return batch.includes(query) || lotCode.includes(query) || format.includes(query) || location.includes(query)
      })
    }

    return items
  })

  // Table headers
  const stockLevelHeaders = [
    { title: 'Batch', key: 'batch_name', sortable: true },
    { title: 'Lot Code', key: 'lot_code', sortable: true },
    { title: 'Format', key: 'package_format_name', sortable: true },
    { title: 'Quantity', key: 'current_quantity', sortable: true },
    { title: 'Location', key: 'stock_location_name', sortable: true },
    { title: 'Packaged', key: 'packaged_at', sortable: true },
    { title: 'Best By', key: 'best_by', sortable: true },
  ]

  const allLotsHeaders = [
    { title: 'Batch', key: 'batch_name', sortable: true },
    { title: 'Lot Code', key: 'lot_code', sortable: true },
    { title: 'Format', key: 'package_format_name', sortable: true },
    { title: 'Qty (Initial)', key: 'quantity', sortable: true },
    { title: 'Packaged', key: 'packaged_at', sortable: true },
    { title: 'Best By', key: 'best_by', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
  ]

  // Data loading
  onMounted(async () => {
    await refreshAll()
  })

  async function refreshAll () {
    loading.value = true
    errorMessage.value = ''
    try {
      // Load inventory data and cross-service data in parallel
      const [stockLevelResult, beerLotResult, batchResult] = await Promise.allSettled([
        getBeerLotStockLevels(),
        fetchBeerLots(),
        getBatches(),
      ])

      if (stockLevelResult.status === 'fulfilled') {
        stockLevels.value = stockLevelResult.value
      } else {
        errorMessage.value = 'Unable to load stock levels'
      }

      if (beerLotResult.status === 'fulfilled') {
        beerLots.value = beerLotResult.value
      }

      // Cross-service data: graceful failure
      if (batchResult.status === 'fulfilled') {
        batches.value = batchResult.value
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load finished goods'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }
</script>

<style scoped>
.product-page {
  position: relative;
}

.stock-level-table :deep(tr) {
  cursor: default;
}
</style>
