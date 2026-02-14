<template>
  <div>
    <v-progress-linear
      v-if="loading"
      class="mb-4"
      color="primary"
      indeterminate
    />

    <v-alert
      v-if="fetchError"
      class="mb-4"
      density="comfortable"
      type="warning"
      variant="tonal"
    >
      {{ fetchError }}
    </v-alert>

    <template v-if="!loading && !fetchError">
      <!-- Auto-fill button -->
      <div v-if="linkedIngredients.length > 0" class="d-flex justify-end mb-3">
        <v-btn
          color="primary"
          :disabled="lotsLoading"
          prepend-icon="mdi-auto-fix"
          size="small"
          variant="tonal"
          @click="autoFillAll"
        >
          Auto-fill All
        </v-btn>
      </div>

      <!-- Lots loading -->
      <v-progress-linear
        v-if="lotsLoading"
        class="mb-4"
        color="primary"
        indeterminate
      />

      <v-alert
        v-if="lotsError"
        class="mb-4"
        density="comfortable"
        type="warning"
        variant="tonal"
      >
        {{ lotsError }}
      </v-alert>

      <!-- Ingredient cards -->
      <template v-if="!lotsLoading && !lotsError">
        <v-card
          v-for="item in linkedIngredients"
          :key="item.uuid"
          class="mb-3"
          variant="outlined"
        >
          <v-card-text class="pa-3">
            <!-- Ingredient header -->
            <div class="d-flex align-center justify-space-between mb-2">
              <div class="d-flex align-center ga-2">
                <span class="font-weight-medium">{{ item.name }}</span>
                <v-chip
                  :color="getStageColor(item.use_stage)"
                  size="x-small"
                  variant="tonal"
                >
                  {{ formatStage(item.use_stage) }}
                </v-chip>
              </div>
              <div class="d-flex align-center ga-1">
                <v-icon
                  v-if="getPickStatus(item) === 'exact'"
                  color="success"
                  icon="mdi-check-circle"
                  size="small"
                />
                <v-icon
                  v-else-if="getPickStatus(item) === 'partial'"
                  color="warning"
                  icon="mdi-alert-circle-outline"
                  size="small"
                />
                <v-icon
                  v-else-if="getPickStatus(item) === 'none' && getLotsForIngredient(item.ingredient_uuid).length === 0"
                  color="error"
                  icon="mdi-close-circle-outline"
                  size="small"
                />
              </div>
            </div>

            <!-- Needed amount -->
            <div class="text-body-2 text-medium-emphasis mb-2">
              Need: {{ formatAmount(item.amount, item.amount_unit) }}
            </div>

            <!-- Pick total -->
            <div
              v-if="getIngredientPickTotal(item) > 0"
              class="text-body-2 mb-2"
              :class="getPickSummaryColor(item)"
            >
              Picked: {{ formatPickTotal(item) }}
            </div>

            <!-- No lots available -->
            <div
              v-if="getLotsForIngredient(item.ingredient_uuid).length === 0"
              class="text-body-2 text-medium-emphasis"
            >
              No lots available.
            </div>

            <!-- Lot rows -->
            <div
              v-for="lot in getLotsForIngredient(item.ingredient_uuid)"
              :key="lot.uuid"
              class="lot-row pa-2 mb-1 rounded"
            >
              <div class="d-flex align-center justify-space-between mb-1">
                <div class="text-body-2 font-weight-medium">
                  {{ lot.brewery_lot_code || lot.originator_lot_code || lot.uuid.slice(0, 8) }}
                </div>
                <div class="text-body-2 text-medium-emphasis">
                  {{ formatAmount(lot.current_amount, lot.current_unit) }} avail
                </div>
              </div>
              <v-text-field
                density="compact"
                hide-details
                inputmode="decimal"
                :max="lot.current_amount"
                min="0"
                :model-value="getPickAmount(item.ingredient_uuid, lot.uuid)"
                :placeholder="`0 ${lot.current_unit}`"
                :suffix="lot.current_unit"
                type="number"
                variant="outlined"
                @update:model-value="setPickAmount(item.ingredient_uuid, lot.uuid, $event)"
              />
            </div>
          </v-card-text>
        </v-card>

        <!-- Unlinked ingredients notice -->
        <v-alert
          v-if="unlinkedIngredients.length > 0"
          class="mb-3"
          density="comfortable"
          type="info"
          variant="tonal"
        >
          {{ unlinkedIngredients.length }} ingredient{{ unlinkedIngredients.length > 1 ? 's are' : ' is' }}
          not linked to inventory:
          {{ unlinkedIngredients.map(i => i.name).join(', ') }}.
        </v-alert>
      </template>
    </template>
  </div>
</template>

<script lang="ts" setup>
  import type {
    BatchUsagePick,
    CreateBatchUsageRequest,
    IngredientLot,
    RecipeIngredient,
    RecipeUseStage,
    StockLevel,
  } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'

  const props = defineProps<{
    recipeUuid: string
    batchUuid: string
  }>()

  const emit = defineEmits<{
    completed: [data: { ingredientCount: number, lotCount: number }]
  }>()

  const { getRecipeIngredients } = useProductionApi()
  const { getStockLevels, getIngredientLots, createBatchUsage } = useInventoryApi()
  const { showNotice } = useSnackbar()

  // State
  const loading = ref(false)
  const fetchError = ref<string | null>(null)
  const recipeIngredients = ref<RecipeIngredient[]>([])
  const stockLevels = ref<StockLevel[]>([])

  const lotsLoading = ref(false)
  const lotsError = ref<string | null>(null)
  const ingredientLots = ref<Map<string, IngredientLot[]>>(new Map())
  const pickAmounts = reactive<Record<string, Record<string, number>>>({})

  const confirming = ref(false)

  // Computed
  interface LinkedIngredient extends RecipeIngredient {
    ingredient_uuid: string
  }

  const linkedIngredients = computed<LinkedIngredient[]>(() =>
    sortIngredients(
      recipeIngredients.value.filter(
        (i): i is LinkedIngredient => !!i.ingredient_uuid,
      ),
    ),
  )

  const unlinkedIngredients = computed(() =>
    recipeIngredients.value.filter(i => !i.ingredient_uuid),
  )

  const hasAnyPicks = computed(() => {
    for (const ingredientUuid of Object.keys(pickAmounts)) {
      const lotPicks = pickAmounts[ingredientUuid]
      if (!lotPicks) continue
      for (const lotUuid of Object.keys(lotPicks)) {
        if ((lotPicks[lotUuid] ?? 0) > 0) return true
      }
    }
    return false
  })

  const stockMap = computed(() => {
    const map = new Map<string, StockLevel>()
    for (const level of stockLevels.value) {
      map.set(level.ingredient_uuid, level)
    }
    return map
  })

  // Expose for parent
  defineExpose({
    hasAnyPicks,
    confirmPicks,
    confirming,
  })

  // Load data on mount
  onMounted(async () => {
    await loadData()
  })

  async function loadData () {
    loading.value = true
    fetchError.value = null

    try {
      const [ingredientsResult, stockResult] = await Promise.allSettled([
        getRecipeIngredients(props.recipeUuid),
        getStockLevels(),
      ])

      if (ingredientsResult.status === 'fulfilled') {
        recipeIngredients.value = ingredientsResult.value
      } else {
        fetchError.value = 'Failed to load recipe ingredients.'
        return
      }

      if (stockResult.status === 'fulfilled') {
        stockLevels.value = stockResult.value
      }

      // Load lots for linked ingredients
      await loadLots()
    } finally {
      loading.value = false
    }
  }

  async function loadLots () {
    const ingredientUuids = new Set<string>()
    for (const item of linkedIngredients.value) {
      if (item.ingredient_uuid) {
        ingredientUuids.add(item.ingredient_uuid)
      }
    }

    if (ingredientUuids.size === 0) return

    lotsLoading.value = true
    lotsError.value = null

    try {
      const entries = [...ingredientUuids]
      const results = await Promise.allSettled(
        entries.map(uuid => getIngredientLots({ ingredient_uuid: uuid })),
      )

      const lotsMap = new Map<string, IngredientLot[]>()
      let anyFailed = false

      for (let i = 0; i < entries.length; i++) {
        const result = results[i]
        const entryUuid = entries[i]
        if (!result || !entryUuid) continue
        if (result.status === 'fulfilled') {
          const availableLots = result.value.filter((lot: IngredientLot) => lot.current_amount > 0)
          lotsMap.set(entryUuid, availableLots)
        } else {
          anyFailed = true
          lotsMap.set(entryUuid, [])
        }
      }

      ingredientLots.value = lotsMap

      if (anyFailed) {
        lotsError.value = 'Some lot data could not be loaded.'
      }
    } catch {
      lotsError.value = 'Failed to load lot data.'
    } finally {
      lotsLoading.value = false
    }
  }

  function getLotsForIngredient (ingredientUuid: string): IngredientLot[] {
    return ingredientLots.value.get(ingredientUuid) ?? []
  }

  // Pick amount management
  function getPickAmount (ingredientUuid: string, lotUuid: string): number {
    return pickAmounts[ingredientUuid]?.[lotUuid] ?? 0
  }

  function setPickAmount (ingredientUuid: string, lotUuid: string, value: string | number | null) {
    if (!pickAmounts[ingredientUuid]) {
      pickAmounts[ingredientUuid] = {}
    }

    const numValue = typeof value === 'string' ? Number.parseFloat(value) : (value ?? 0)
    const lot = getLotsForIngredient(ingredientUuid).find(l => l.uuid === lotUuid)
    const maxAmount = lot?.current_amount ?? 0
    const clamped = Math.max(0, Math.min(Number.isNaN(numValue) ? 0 : numValue, maxAmount))
    pickAmounts[ingredientUuid][lotUuid] = clamped
  }

  function getIngredientPickTotal (item: RecipeIngredient): number {
    if (!item.ingredient_uuid) return 0
    const picks = pickAmounts[item.ingredient_uuid]
    if (!picks) return 0
    return Object.values(picks).reduce((sum, amt) => sum + amt, 0)
  }

  type PickStatus = 'none' | 'partial' | 'exact' | 'over'

  function getPickStatus (item: RecipeIngredient): PickStatus {
    const total = getIngredientPickTotal(item)
    if (total === 0) return 'none'
    const needed = item.amount
    if (Math.abs(total - needed) < 0.01) return 'exact'
    if (total > needed) return 'over'
    return 'partial'
  }

  function getPickSummaryColor (item: RecipeIngredient): string {
    const status = getPickStatus(item)
    if (status === 'exact') return 'text-success'
    if (status === 'over') return 'text-warning'
    return 'text-orange'
  }

  function formatPickTotal (item: RecipeIngredient): string {
    const total = getIngredientPickTotal(item)
    return `${total.toFixed(2)} / ${item.amount.toFixed(2)} ${item.amount_unit}`
  }

  // Auto-fill: for each ingredient, if one lot has enough stock, fill it
  function autoFillAll () {
    for (const item of linkedIngredients.value) {
      if (!item.ingredient_uuid) continue
      const lots = getLotsForIngredient(item.ingredient_uuid)
      if (lots.length === 0) continue

      // If already picked, skip
      if (getIngredientPickTotal(item) > 0) continue

      const needed = item.amount

      // Try to find a single lot with enough
      const sufficientLot = lots.find(lot => lot.current_amount >= needed)
      if (sufficientLot) {
        setPickAmount(item.ingredient_uuid, sufficientLot.uuid, needed)
      } else if (lots.length === 1 && lots[0]) {
        // Only one lot, fill what's available
        setPickAmount(item.ingredient_uuid, lots[0].uuid, lots[0].current_amount)
      }
    }
  }

  // Resolve stock location for an ingredient
  function resolveStockLocationUuid (ingredientUuid: string): string | null {
    const level = stockMap.value.get(ingredientUuid)
    if (!level || level.locations.length === 0) return null
    const locationWithStock = level.locations.find(loc => loc.quantity > 0)
    const fallback = level.locations[0]
    return locationWithStock?.location_uuid ?? fallback?.location_uuid ?? null
  }

  // Confirm picks
  async function confirmPicks (): Promise<boolean> {
    confirming.value = true

    try {
      const picks: BatchUsagePick[] = []
      const ingredientNames = new Set<string>()
      const lotUuids = new Set<string>()

      for (const ingredientUuid of Object.keys(pickAmounts)) {
        const lotPicks = pickAmounts[ingredientUuid]
        if (!lotPicks) continue

        const hasPositivePick = Object.values(lotPicks).some(amt => amt > 0)
        if (!hasPositivePick) continue

        const locationUuid = resolveStockLocationUuid(ingredientUuid)
        if (!locationUuid) continue

        const recipeItem = recipeIngredients.value.find(
          i => i.ingredient_uuid === ingredientUuid,
        )
        if (!recipeItem) continue

        ingredientNames.add(recipeItem.name)

        for (const [lotUuid, amount] of Object.entries(lotPicks)) {
          if (amount <= 0) continue
          lotUuids.add(lotUuid)
          picks.push({
            ingredient_lot_uuid: lotUuid,
            stock_location_uuid: locationUuid,
            amount,
            amount_unit: recipeItem.amount_unit,
          })
        }
      }

      if (picks.length === 0) {
        showNotice('No picks to submit', 'warning')
        return false
      }

      const request: CreateBatchUsageRequest = {
        production_ref_uuid: props.batchUuid,
        used_at: new Date().toISOString(),
        picks,
        notes: 'Brew day ingredient pick',
      }

      await createBatchUsage(request)
      showNotice('Ingredients picked successfully')

      emit('completed', {
        ingredientCount: ingredientNames.size,
        lotCount: lotUuids.size,
      })

      return true
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to pick ingredients'
      showNotice(message, 'error')
      return false
    } finally {
      confirming.value = false
    }
  }

  // Sorting
  function sortIngredients<T extends RecipeIngredient> (ingredients: T[]): T[] {
    const stageOrder: Record<string, number> = {
      mash: 0,
      boil: 1,
      whirlpool: 2,
      fermentation: 3,
      packaging: 4,
    }
    return [...ingredients].sort((a, b) => {
      const stageA = stageOrder[a.use_stage] ?? 99
      const stageB = stageOrder[b.use_stage] ?? 99
      if (stageA !== stageB) return stageA - stageB
      return a.sort_order - b.sort_order
    })
  }

  // Formatting helpers
  function formatAmount (amount: number, unit: string): string {
    return `${amount.toFixed(2)} ${unit}`
  }

  function formatStage (stage: RecipeUseStage): string {
    const labels: Record<RecipeUseStage, string> = {
      mash: 'Mash',
      boil: 'Boil',
      whirlpool: 'Whirlpool',
      fermentation: 'Fermentation',
      packaging: 'Packaging',
    }
    return labels[stage] ?? stage.charAt(0).toUpperCase() + stage.slice(1)
  }

  function getStageColor (stage: RecipeUseStage): string {
    const colors: Record<RecipeUseStage, string> = {
      mash: 'brown',
      boil: 'orange',
      whirlpool: 'blue',
      fermentation: 'purple',
      packaging: 'grey',
    }
    return colors[stage] ?? 'secondary'
  }
</script>

<style scoped>
.lot-row {
  background: rgba(var(--v-theme-on-surface), 0.03);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.06);
}
</style>
