<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-text>
      <v-progress-linear
        v-if="loading"
        class="mb-4"
        color="primary"
        indeterminate
      />

      <!-- No recipe assigned -->
      <v-alert
        v-if="!loading && !recipeUuid"
        density="comfortable"
        type="info"
        variant="tonal"
      >
        No recipe assigned to this batch. Assign a recipe to generate a pick list.
      </v-alert>

      <!-- Recipe assigned but no ingredients -->
      <template v-else-if="!loading && recipeIngredients.length === 0 && !fetchError">
        <div class="text-overline text-medium-emphasis mb-2">Recipe</div>
        <div class="text-h6 mb-4">{{ displayRecipeName }}</div>
        <v-alert
          density="comfortable"
          type="info"
          variant="tonal"
        >
          This recipe has no ingredients yet. Add ingredients to the recipe to generate a pick list.
        </v-alert>
      </template>

      <!-- Fetch error -->
      <v-alert
        v-else-if="!loading && fetchError"
        density="comfortable"
        type="warning"
        variant="tonal"
      >
        {{ fetchError }}
      </v-alert>

      <!-- Pick list -->
      <template v-else-if="!loading && recipeIngredients.length > 0">
        <div class="d-flex align-center justify-space-between mb-4">
          <div>
            <div class="text-overline text-medium-emphasis mb-1">Recipe</div>
            <div class="text-h6">{{ displayRecipeName }}</div>
          </div>
          <div v-if="!pickingMode && canPick" class="d-flex align-center ga-2">
            <v-chip
              v-if="ingredientsPicked"
              color="success"
              prepend-icon="mdi-check-circle"
              size="small"
              variant="tonal"
            >
              Ingredients picked
            </v-chip>
            <v-btn
              color="primary"
              prepend-icon="mdi-clipboard-list-outline"
              size="small"
              variant="tonal"
              @click="enterPickingMode"
            >
              <span class="d-none d-sm-inline">Pick Ingredients</span>
              <span class="d-sm-none">Pick</span>
            </v-btn>
          </div>
          <v-btn
            v-if="pickingMode"
            prepend-icon="mdi-close"
            size="small"
            variant="text"
            @click="exitPickingMode"
          >
            Cancel
          </v-btn>
        </div>

        <!-- Picking Mode -->
        <template v-if="pickingMode">
          <v-alert
            class="mb-4"
            density="comfortable"
            type="info"
            variant="tonal"
          >
            Select lots and amounts for each ingredient. Lots are shown oldest first (FIFO).
          </v-alert>

          <!-- Loading lots indicator -->
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

          <!-- Per-ingredient lot selection -->
          <template v-if="!lotsLoading && !lotsError">
            <div
              v-for="item in linkedIngredients"
              :key="item.uuid"
              class="mb-4"
            >
              <v-card variant="outlined">
                <v-card-text class="pa-3">
                  <!-- Ingredient header -->
                  <div class="d-flex align-center justify-space-between mb-2">
                    <div>
                      <span class="font-weight-medium">{{ item.name }}</span>
                      <v-chip
                        class="ml-2"
                        :color="getStageColor(item.use_stage)"
                        size="x-small"
                        variant="tonal"
                      >
                        {{ formatStage(item.use_stage) }}
                      </v-chip>
                    </div>
                    <div class="d-flex align-center ga-2">
                      <span class="text-body-2 text-medium-emphasis">
                        Need: {{ formatAmount(getDisplayAmount(item), item.amount_unit) }}
                        <span v-if="isBatchScaling" class="text-caption">
                          (Recipe: {{ formatAmount(item.amount, item.amount_unit) }})
                        </span>
                      </span>
                      <v-icon
                        v-if="getPickStatus(item) === 'exact'"
                        color="success"
                        icon="mdi-check-circle"
                        size="small"
                      />
                      <v-icon
                        v-else-if="getPickStatus(item) === 'over'"
                        color="warning"
                        icon="mdi-alert"
                        size="small"
                      />
                      <v-icon
                        v-else-if="getPickStatus(item) === 'partial'"
                        color="orange"
                        icon="mdi-alert-circle-outline"
                        size="small"
                      />
                    </div>
                  </div>

                  <!-- Pick total summary -->
                  <div
                    v-if="getIngredientPickTotal(item) > 0"
                    class="text-body-2 mb-2"
                    :class="getPickSummaryColor(item)"
                  >
                    Picked: {{ formatPickBreakdown(item) }}
                  </div>

                  <!-- Available lots -->
                  <div v-if="getLotsForIngredient(item.ingredient_uuid!).length === 0" class="text-body-2 text-medium-emphasis">
                    No lots available for this ingredient.
                  </div>

                  <!-- Lot cards (mobile-friendly) -->
                  <div
                    v-for="lot in getLotsForIngredient(item.ingredient_uuid!)"
                    :key="lot.uuid"
                    class="lot-row pa-2 mb-1 rounded"
                  >
                    <v-row align="center" dense>
                      <v-col cols="12" sm="5">
                        <div class="text-body-2 font-weight-medium">
                          Lot: {{ lot.brewery_lot_code || lot.originator_lot_code || lot.uuid.slice(0, 8) }}
                        </div>
                        <div class="text-caption text-medium-emphasis">
                          Received {{ formatDate(lot.received_at) }}
                          <span v-if="lot.best_by_at"> · Best by {{ formatDate(lot.best_by_at) }}</span>
                        </div>
                      </v-col>
                      <v-col cols="6" sm="3">
                        <div class="text-body-2">
                          {{ formatAmount(lot.current_amount, lot.current_unit) }} available
                        </div>
                      </v-col>
                      <v-col cols="6" sm="4">
                        <v-text-field
                          density="compact"
                          hide-details
                          inputmode="numeric"
                          :max="lot.current_amount"
                          min="0"
                          :model-value="getPickAmount(item.ingredient_uuid!, lot.uuid)"
                          :placeholder="`0 ${lot.current_unit}`"
                          :suffix="lot.current_unit"
                          type="number"
                          variant="outlined"
                          @update:model-value="setPickAmount(item.ingredient_uuid!, lot.uuid, $event)"
                        />
                      </v-col>
                    </v-row>
                  </div>
                </v-card-text>
              </v-card>
            </div>

            <!-- Unlinked ingredients notice -->
            <v-alert
              v-if="unlinkedIngredients.length > 0"
              class="mb-4"
              density="comfortable"
              type="info"
              variant="tonal"
            >
              {{ unlinkedIngredients.length }} ingredient{{ unlinkedIngredients.length > 1 ? 's are' : ' is' }}
              not linked to inventory and cannot be picked:
              {{ unlinkedIngredients.map(i => i.name || i.ingredient_type || 'Unknown ingredient').join(', ') }}.
            </v-alert>

            <!-- Confirm button -->
            <div class="d-flex justify-end mt-4">
              <v-btn
                color="primary"
                :disabled="!hasAnyPicks"
                :loading="confirming"
                min-height="44"
                prepend-icon="mdi-check"
                size="large"
                variant="flat"
                @click="showConfirmDialog = true"
              >
                Confirm &amp; Deduct Inventory
              </v-btn>
            </div>
          </template>
        </template>

        <!-- Normal (non-picking) view -->
        <template v-if="!pickingMode">
          <!-- Needed Today Section -->
          <template v-if="neededToday.length > 0">
            <div class="d-flex align-center mb-2">
              <v-icon class="mr-2" color="primary" icon="mdi-clipboard-check-outline" size="small" />
              <span class="text-subtitle-1 font-weight-medium">Needed Today</span>
              <v-chip class="ml-2" color="primary" size="x-small" variant="tonal">
                {{ neededToday.length }}
              </v-chip>
            </div>

            <!-- Desktop Table -->
            <v-table class="data-table d-none d-md-block mb-4" density="compact">
              <thead>
                <tr>
                  <th>Ingredient</th>
                  <th>Amount</th>
                  <th>Stage</th>
                  <th>Timing</th>
                  <th>Available</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in neededToday" :key="item.uuid">
                  <td>
                    <div class="font-weight-medium">{{ item.name }}</div>
                    <div v-if="item.use_type" class="text-caption text-medium-emphasis">
                      {{ formatUseType(item.use_type) }}
                    </div>
                    <div v-if="item.ingredient_type === 'hop' && item.alpha_acid_assumed" class="text-caption text-medium-emphasis">
                      {{ formatAlphaAcid(item.alpha_acid_assumed) }} AA
                    </div>
                  </td>
                  <td>
                    <div>{{ formatAmount(getDisplayAmount(item), item.amount_unit) }}</div>
                    <div v-if="isBatchScaling" class="text-caption text-medium-emphasis">
                      Recipe: {{ formatAmount(item.amount, item.amount_unit) }}
                    </div>
                  </td>
                  <td>
                    <v-chip
                      :color="getStageColor(item.use_stage)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatStage(item.use_stage) }}
                    </v-chip>
                  </td>
                  <td>{{ formatTiming(item) }}</td>
                  <td>
                    <!-- Not linked to inventory -->
                    <v-chip v-if="!item.ingredient_uuid" color="grey" size="x-small" variant="outlined">
                      Not linked
                    </v-chip>
                    <!-- Stock unavailable -->
                    <v-chip v-else-if="stockUnavailable" color="grey" size="x-small" variant="outlined">
                      Stock unavailable
                    </v-chip>
                    <!-- Stock level -->
                    <v-chip
                      v-else
                      :color="getStockColor(item.ingredient_uuid, getDisplayAmount(item), item.amount_unit)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatStock(item.ingredient_uuid, item.amount_unit) }}
                    </v-chip>
                  </td>
                </tr>
              </tbody>
            </v-table>

            <!-- Mobile Cards -->
            <div class="d-md-none mb-4">
              <v-card
                v-for="item in neededToday"
                :key="item.uuid"
                class="mb-2 ingredient-card"
                variant="tonal"
              >
                <v-card-text class="pa-3">
                  <div class="d-flex align-center justify-space-between mb-1">
                    <span class="font-weight-medium">{{ item.name }}</span>
                    <v-chip
                      :color="getStageColor(item.use_stage)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatStage(item.use_stage) }}
                    </v-chip>
                  </div>
                  <div v-if="item.use_type" class="text-caption text-medium-emphasis mb-1">
                    {{ formatUseType(item.use_type) }}
                    <span v-if="item.ingredient_type === 'hop' && item.alpha_acid_assumed">
                      &bull; {{ formatAlphaAcid(item.alpha_acid_assumed) }} AA
                    </span>
                  </div>
                  <div class="d-flex align-center justify-space-between text-body-2">
                    <div>
                      <span>{{ formatAmount(getDisplayAmount(item), item.amount_unit) }}</span>
                      <div v-if="isBatchScaling" class="text-caption text-medium-emphasis">
                        Recipe: {{ formatAmount(item.amount, item.amount_unit) }}
                      </div>
                    </div>
                    <span class="text-medium-emphasis">{{ formatTiming(item) }}</span>
                  </div>
                  <div class="d-flex align-center justify-end mt-2">
                    <v-chip v-if="!item.ingredient_uuid" color="grey" size="x-small" variant="outlined">
                      Not linked
                    </v-chip>
                    <v-chip v-else-if="stockUnavailable" color="grey" size="x-small" variant="outlined">
                      Stock unavailable
                    </v-chip>
                    <v-chip
                      v-else
                      :color="getStockColor(item.ingredient_uuid, getDisplayAmount(item), item.amount_unit)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatStock(item.ingredient_uuid, item.amount_unit) }}
                    </v-chip>
                  </div>
                </v-card-text>
              </v-card>
            </div>
          </template>

          <!-- Needed Later Section -->
          <template v-if="neededLater.length > 0">
            <div class="d-flex align-center mb-2">
              <v-icon class="mr-2" color="secondary" icon="mdi-clock-outline" size="small" />
              <span class="text-subtitle-1 font-weight-medium text-medium-emphasis">Needed Later</span>
              <v-chip class="ml-2" color="secondary" size="x-small" variant="tonal">
                {{ neededLater.length }}
              </v-chip>
            </div>

            <!-- Desktop Table -->
            <v-table class="data-table later-table d-none d-md-block" density="compact">
              <thead>
                <tr>
                  <th>Ingredient</th>
                  <th>Amount</th>
                  <th>Stage</th>
                  <th>Timing</th>
                  <th>Available</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in neededLater" :key="item.uuid">
                  <td>
                    <div class="font-weight-medium">{{ item.name }}</div>
                    <div v-if="item.use_type" class="text-caption text-medium-emphasis">
                      {{ formatUseType(item.use_type) }}
                    </div>
                    <div v-if="item.ingredient_type === 'hop' && item.alpha_acid_assumed" class="text-caption text-medium-emphasis">
                      {{ formatAlphaAcid(item.alpha_acid_assumed) }} AA
                    </div>
                  </td>
                  <td>
                    <div>{{ formatAmount(getDisplayAmount(item), item.amount_unit) }}</div>
                    <div v-if="isBatchScaling" class="text-caption text-medium-emphasis">
                      Recipe: {{ formatAmount(item.amount, item.amount_unit) }}
                    </div>
                  </td>
                  <td>
                    <v-chip
                      :color="getStageColor(item.use_stage)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatStage(item.use_stage) }}
                    </v-chip>
                  </td>
                  <td>{{ formatTiming(item) }}</td>
                  <td>
                    <v-chip v-if="!item.ingredient_uuid" color="grey" size="x-small" variant="outlined">
                      Not linked
                    </v-chip>
                    <v-chip v-else-if="stockUnavailable" color="grey" size="x-small" variant="outlined">
                      Stock unavailable
                    </v-chip>
                    <v-chip
                      v-else
                      :color="getStockColor(item.ingredient_uuid, getDisplayAmount(item), item.amount_unit)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatStock(item.ingredient_uuid, item.amount_unit) }}
                    </v-chip>
                  </td>
                </tr>
              </tbody>
            </v-table>

            <!-- Mobile Cards -->
            <div class="d-md-none">
              <v-card
                v-for="item in neededLater"
                :key="item.uuid"
                class="mb-2 ingredient-card ingredient-card--later"
                variant="tonal"
              >
                <v-card-text class="pa-3">
                  <div class="d-flex align-center justify-space-between mb-1">
                    <span class="font-weight-medium text-medium-emphasis">{{ item.name }}</span>
                    <v-chip
                      :color="getStageColor(item.use_stage)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatStage(item.use_stage) }}
                    </v-chip>
                  </div>
                  <div v-if="item.use_type" class="text-caption text-medium-emphasis mb-1">
                    {{ formatUseType(item.use_type) }}
                    <span v-if="item.ingredient_type === 'hop' && item.alpha_acid_assumed">
                      &bull; {{ formatAlphaAcid(item.alpha_acid_assumed) }} AA
                    </span>
                  </div>
                  <div class="d-flex align-center justify-space-between text-body-2">
                    <div>
                      <span>{{ formatAmount(getDisplayAmount(item), item.amount_unit) }}</span>
                      <div v-if="isBatchScaling" class="text-caption text-medium-emphasis">
                        Recipe: {{ formatAmount(item.amount, item.amount_unit) }}
                      </div>
                    </div>
                    <span class="text-medium-emphasis">{{ formatTiming(item) }}</span>
                  </div>
                  <div class="d-flex align-center justify-end mt-2">
                    <v-chip v-if="!item.ingredient_uuid" color="grey" size="x-small" variant="outlined">
                      Not linked
                    </v-chip>
                    <v-chip v-else-if="stockUnavailable" color="grey" size="x-small" variant="outlined">
                      Stock unavailable
                    </v-chip>
                    <v-chip
                      v-else
                      :color="getStockColor(item.ingredient_uuid, getDisplayAmount(item), item.amount_unit)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatStock(item.ingredient_uuid, item.amount_unit) }}
                    </v-chip>
                  </div>
                </v-card-text>
              </v-card>
            </div>
          </template>
        </template>
      </template>
    </v-card-text>
  </v-card>

  <!-- Confirmation Dialog -->
  <v-dialog
    v-model="showConfirmDialog"
    :max-width="500"
    persistent
  >
    <v-card>
      <v-card-title class="text-h6">Confirm Ingredient Pick</v-card-title>
      <v-card-text>
        <v-alert
          v-if="confirmError"
          class="mb-4"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ confirmError }}
        </v-alert>
        <p class="text-body-1 mb-4">
          This will deduct the selected ingredients from inventory. This cannot be undone.
        </p>
        <div class="text-body-2 text-medium-emphasis">
          <div v-for="item in linkedIngredients" :key="item.uuid">
            <template v-if="getIngredientPickTotal(item) > 0">
              <strong>{{ item.name }}:</strong> {{ formatPickBreakdown(item) }}
            </template>
          </div>
        </div>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="confirming" variant="text" @click="cancelConfirm">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="confirming"
          variant="flat"
          @click="confirmPicks"
        >
          Confirm &amp; Deduct Inventory
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { BatchUsagePick, CreateBatchUsageRequest, IngredientLot, RecipeIngredient, RecipeUseStage, RecipeUseType, StockLevel } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  const props = withDefaults(defineProps<{
    batchUuid: string | null
    batchVolume?: number | null
    batchVolumeUnit?: string | null
    recipeUuid: string | null
    recipeName: string | null
  }>(), {
    batchVolume: null,
    batchVolumeUnit: null,
  })

  const { getRecipeIngredients, getRecipe } = useProductionApi()
  const { getStockLevels, getIngredientLots, createBatchUsage } = useInventoryApi()
  const { formatTemperaturePreferred } = useUnitPreferences()
  const { showNotice } = useSnackbar()

  // State
  const loading = ref(false)
  const fetchError = ref<string | null>(null)
  const recipeIngredients = ref<RecipeIngredient[]>([])
  const stockLevels = ref<StockLevel[]>([])
  const stockUnavailable = ref(false)
  const resolvedRecipeName = ref<string | null>(null)
  const resolvedRecipeBatchSize = ref<number | null>(null)
  const resolvedRecipeBatchSizeUnit = ref<string | null>(null)

  // Picking mode state
  const pickingMode = ref(false)
  const lotsLoading = ref(false)
  const lotsError = ref<string | null>(null)
  const ingredientLots = ref<Map<string, IngredientLot[]>>(new Map())
  const pickAmounts = reactive<Record<string, Record<string, number>>>({})
  // key: ingredientUuid -> lotUuid -> amount

  // Confirmation state
  const showConfirmDialog = ref(false)
  const confirming = ref(false)
  const confirmError = ref<string | null>(null)
  const ingredientsPicked = ref(false)

  // Computed: recipe display name
  const displayRecipeName = computed(() =>
    props.recipeName ?? resolvedRecipeName.value ?? 'Recipe',
  )

  // Batch-to-recipe scaling: if the batch has a volume that differs from the recipe's batch_size,
  // compute a scale factor and show scaled amounts in the pick list.
  const batchScaleFactor = computed(() => {
    const batchVol = props.batchVolume
    const recipeBatchSize = resolvedRecipeBatchSize.value
    if (!batchVol || !recipeBatchSize || recipeBatchSize === 0) return 1

    // For V1, assume same units if both are set. Cross-unit conversion can be added later.
    const batchUnit = props.batchVolumeUnit?.toLowerCase()
    const recipeUnit = resolvedRecipeBatchSizeUnit.value?.toLowerCase()
    if (batchUnit && recipeUnit && batchUnit !== recipeUnit) return 1

    return batchVol / recipeBatchSize
  })

  const isBatchScaling = computed(() =>
    batchScaleFactor.value !== 1,
  )

  /** Scale an ingredient amount for the batch. */
  function scaleBatchAmount (amount: number, ingredientScalingFactor: number = 1.0): number {
    if (ingredientScalingFactor === 0) return amount
    const fullyScaled = amount * batchScaleFactor.value
    return amount + (fullyScaled - amount) * ingredientScalingFactor
  }

  /** Get the display amount for an ingredient (scaled or original). */
  function getDisplayAmount (ingredient: RecipeIngredient): number {
    if (isBatchScaling.value) {
      return scaleBatchAmount(ingredient.amount, ingredient.scaling_factor)
    }
    return ingredient.amount
  }

  // Computed: stock lookup map (ingredient UUID → StockLevel)
  const stockMap = computed(() => {
    const map = new Map<string, StockLevel>()
    for (const level of stockLevels.value) {
      map.set(level.ingredient_uuid, level)
    }
    return map
  })

  // Computed: can enter picking mode
  const canPick = computed(() =>
    recipeIngredients.value.length > 0
    && !stockUnavailable.value
    && linkedIngredients.value.length > 0,
  )

  // Computed: ingredients that are linked to inventory (have ingredient_uuid)
  const linkedIngredients = computed(() =>
    sortIngredients(recipeIngredients.value.filter(i => i.ingredient_uuid)),
  )

  // Computed: ingredients not linked to inventory
  const unlinkedIngredients = computed(() =>
    recipeIngredients.value.filter(i => !i.ingredient_uuid),
  )

  // Computed: whether any picks have been entered
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

  // Computed: split ingredients into "today" and "later"
  const TODAY_STAGES: RecipeUseStage[] = ['mash', 'boil', 'whirlpool']

  const neededToday = computed(() =>
    sortIngredients(
      recipeIngredients.value.filter(i => TODAY_STAGES.includes(i.use_stage)),
    ),
  )

  const neededLater = computed(() =>
    sortIngredients(
      recipeIngredients.value.filter(i => !TODAY_STAGES.includes(i.use_stage)),
    ),
  )

  // Watch recipe UUID changes
  watch(() => props.recipeUuid, async (newUuid) => {
    if (newUuid) {
      await loadPickList(newUuid)
    } else {
      recipeIngredients.value = []
      stockLevels.value = []
      stockUnavailable.value = false
      fetchError.value = null
      resolvedRecipeName.value = null
    }
  }, { immediate: true })

  // Data loading
  async function loadPickList (recipeUuid: string) {
    loading.value = true
    fetchError.value = null
    stockUnavailable.value = false

    try {
      // Fetch recipe ingredients, stock levels, and recipe name in parallel.
      // Stock levels and recipe name are non-critical, so use allSettled.
      const needsRecipeName = !props.recipeName
      const [ingredientsResult, stockResult, recipeResult] = await Promise.allSettled([
        getRecipeIngredients(recipeUuid),
        getStockLevels(),
        needsRecipeName ? getRecipe(recipeUuid) : Promise.resolve(null),
      ])

      if (ingredientsResult.status === 'fulfilled') {
        recipeIngredients.value = ingredientsResult.value
      } else {
        recipeIngredients.value = []
        fetchError.value = 'Failed to load recipe ingredients. Please try refreshing.'
        return
      }

      if (stockResult.status === 'fulfilled') {
        stockLevels.value = stockResult.value
      } else {
        stockLevels.value = []
        stockUnavailable.value = true
      }

      if (recipeResult.status === 'fulfilled' && recipeResult.value) {
        resolvedRecipeName.value = recipeResult.value.name
        resolvedRecipeBatchSize.value = recipeResult.value.batch_size
        resolvedRecipeBatchSizeUnit.value = recipeResult.value.batch_size_unit
      }
      // Recipe name fetch failure is non-critical — we fall back to "Recipe"
    } finally {
      loading.value = false
    }
  }

  // Picking mode
  async function enterPickingMode () {
    pickingMode.value = true
    confirmError.value = null
    await loadLots()
  }

  function exitPickingMode () {
    pickingMode.value = false
    ingredientLots.value = new Map()
    clearPickAmounts()
    lotsError.value = null
  }

  async function loadLots () {
    lotsLoading.value = true
    lotsError.value = null

    // Collect unique ingredient UUIDs from linked ingredients
    const ingredientUuids = new Set<string>()
    for (const item of linkedIngredients.value) {
      if (item.ingredient_uuid) {
        ingredientUuids.add(item.ingredient_uuid)
      }
    }

    try {
      // Fetch lots for all ingredients in parallel
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
          // Filter to lots with stock remaining
          const availableLots = result.value.filter((lot: IngredientLot) => lot.current_amount > 0)
          lotsMap.set(entryUuid, availableLots)
        } else {
          anyFailed = true
          lotsMap.set(entryUuid, [])
        }
      }

      ingredientLots.value = lotsMap

      if (anyFailed) {
        lotsError.value = 'Some lot data could not be loaded. Available lots may be incomplete.'
      }
    } catch {
      lotsError.value = 'Failed to load lot data. Please try again.'
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

    // Clamp to [0, lot.current_amount]
    const clamped = Math.max(0, Math.min(Number.isNaN(numValue) ? 0 : numValue, maxAmount))
    pickAmounts[ingredientUuid][lotUuid] = clamped
  }

  function clearPickAmounts () {
    for (const key of Object.keys(pickAmounts)) {
      delete pickAmounts[key]
    }
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
    // Use scaled amount for comparison
    const needed = getDisplayAmount(item)
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

  function formatPickBreakdown (item: RecipeIngredient): string {
    if (!item.ingredient_uuid) return ''
    const picks = pickAmounts[item.ingredient_uuid]
    if (!picks) return ''

    const lots = getLotsForIngredient(item.ingredient_uuid)
    const parts: string[] = []
    for (const [lotUuid, amount] of Object.entries(picks)) {
      if (amount > 0) {
        if (Object.values(picks).filter(a => a > 0).length > 1) {
          // Multi-lot: include lot identifier for clarity
          const lot = lots.find(l => l.uuid === lotUuid)
          const lotLabel = lot?.brewery_lot_code || lot?.originator_lot_code || lotUuid.slice(0, 8)
          parts.push(`${amount} (${lotLabel})`)
        } else {
          parts.push(`${amount}`)
        }
      }
    }

    const total = getIngredientPickTotal(item)
    const unit = item.amount_unit
    const breakdown = parts.length > 1
      ? `${parts.join(' + ')} = ${total} ${unit}`
      : `${total} ${unit}`

    const needed = getDisplayAmount(item)
    return `${breakdown} / ${needed.toFixed(2)} ${unit} needed`
  }

  // Resolve stock location for an ingredient pick
  function resolveStockLocationUuid (ingredientUuid: string): string | null {
    const level = stockMap.value.get(ingredientUuid)
    if (!level || level.locations.length === 0) return null
    // V1: use the first location that has stock
    const locationWithStock = level.locations.find(loc => loc.quantity > 0)
    const fallback = level.locations[0]
    return locationWithStock?.location_uuid ?? fallback?.location_uuid ?? null
  }

  // Confirmation flow
  function cancelConfirm () {
    showConfirmDialog.value = false
    confirmError.value = null
  }

  async function confirmPicks () {
    if (!props.batchUuid) return

    confirming.value = true
    confirmError.value = null

    try {
      const picks: BatchUsagePick[] = []
      const skippedIngredients: string[] = []

      for (const ingredientUuid of Object.keys(pickAmounts)) {
        const lotPicks = pickAmounts[ingredientUuid]
        if (!lotPicks) continue

        // Check if any picks exist for this ingredient
        const hasPositivePick = Object.values(lotPicks).some(amt => amt > 0)
        if (!hasPositivePick) continue

        const locationUuid = resolveStockLocationUuid(ingredientUuid)
        if (!locationUuid) {
          // Track ingredients with picks but no resolvable stock location
          const ingredientName = recipeIngredients.value.find(
            i => i.ingredient_uuid === ingredientUuid,
          )?.name ?? ingredientUuid
          skippedIngredients.push(ingredientName)
          continue
        }

        // Find the recipe ingredient to get the unit
        const recipeItem = recipeIngredients.value.find(
          i => i.ingredient_uuid === ingredientUuid,
        )
        if (!recipeItem) continue

        for (const [lotUuid, amount] of Object.entries(lotPicks)) {
          if (amount <= 0) continue

          picks.push({
            ingredient_lot_uuid: lotUuid,
            stock_location_uuid: locationUuid,
            amount,
            amount_unit: recipeItem.amount_unit,
          })
        }
      }

      if (skippedIngredients.length > 0) {
        confirmError.value = `Could not resolve stock location for: ${skippedIngredients.join(', ')}. These picks will not be submitted.`
        return
      }

      if (picks.length === 0) {
        confirmError.value = 'No valid picks to submit.'
        return
      }

      const request: CreateBatchUsageRequest = {
        production_ref_uuid: props.batchUuid,
        used_at: new Date().toISOString(),
        picks,
        notes: `Brew day pick for ${displayRecipeName.value}`,
      }

      await createBatchUsage(request)

      // Success
      showConfirmDialog.value = false
      ingredientsPicked.value = true
      exitPickingMode()
      showNotice('Ingredients picked successfully')

      // Refresh stock levels
      try {
        stockLevels.value = await getStockLevels()
      } catch {
        // Non-critical — stock display may be stale
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'An unexpected error occurred'
      confirmError.value = message
    } finally {
      confirming.value = false
    }
  }

  // Sorting: group by stage order, then by sort_order within stage
  function sortIngredients (ingredients: RecipeIngredient[]): RecipeIngredient[] {
    const stageOrder: Record<string, number> = {
      mash: 0,
      boil: 1,
      whirlpool: 2,
      fermentation: 3,
      packaging: 4,
    }
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...ingredients].sort((a, b) => {
      const stageA = stageOrder[a.use_stage] ?? 99
      const stageB = stageOrder[b.use_stage] ?? 99
      if (stageA !== stageB) return stageA - stageB
      return a.sort_order - b.sort_order
    })
  }

  // Stock display helpers
  function getStockColor (ingredientUuid: string, neededAmount: number, recipeUnit: string): string {
    const level = stockMap.value.get(ingredientUuid)
    if (!level) return 'grey'
    // Only compare quantities when units match; mismatched units can't be compared
    const stockUnit = level.default_unit ?? ''
    if (stockUnit.toLowerCase() !== recipeUnit.toLowerCase()) return 'grey'
    return level.total_on_hand >= neededAmount ? 'success' : 'error'
  }

  function formatStock (ingredientUuid: string, fallbackUnit: string): string {
    const level = stockMap.value.get(ingredientUuid)
    if (!level) return 'No stock'
    const unit = level.default_unit ?? fallbackUnit
    return `${level.total_on_hand.toFixed(2)} ${unit}`
  }

  // Formatting helpers
  function formatAmount (amount: number, unit: string): string {
    return `${amount.toFixed(2)} ${unit}`
  }

  function formatDate (dateStr: string): string {
    if (!dateStr) return '—'
    try {
      return new Date(dateStr).toLocaleDateString()
    } catch {
      return '—'
    }
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

  function formatUseType (useType: RecipeUseType | null): string {
    if (!useType) return ''
    const labels: Record<string, string> = {
      bittering: 'Bittering',
      flavor: 'Flavor',
      aroma: 'Aroma',
      dry_hop: 'Dry Hop',
      base: 'Base',
      specialty: 'Specialty',
      adjunct: 'Adjunct',
      sugar: 'Sugar',
      primary: 'Primary',
      secondary: 'Secondary',
      bottle: 'Bottle',
      other: 'Other',
    }
    return labels[useType] ?? useType.charAt(0).toUpperCase() + useType.slice(1).replace(/_/g, ' ')
  }

  function formatAlphaAcid (aa: number | null): string {
    if (aa === null || aa === undefined) return '—'
    return `${aa.toFixed(1)}%`
  }

  function formatTiming (ingredient: RecipeIngredient): string {
    if (ingredient.use_stage === 'fermentation') {
      if (ingredient.timing_duration_minutes !== null && ingredient.timing_duration_minutes !== undefined) {
        const days = Math.round(ingredient.timing_duration_minutes / 1440)
        return `${days} day${days === 1 ? '' : 's'}`
      }
      return '—'
    }

    if (ingredient.timing_duration_minutes !== null && ingredient.timing_duration_minutes !== undefined) {
      return `${ingredient.timing_duration_minutes} min`
    }

    if (ingredient.timing_temperature_c !== null && ingredient.timing_temperature_c !== undefined) {
      return `@ ${formatTemperaturePreferred(ingredient.timing_temperature_c, 'c')}`
    }

    return '—'
  }
</script>

<style scoped>
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

.later-table {
  opacity: 0.85;
}

.ingredient-card {
  background: rgba(var(--v-theme-surface), 0.5);
}

.ingredient-card--later {
  opacity: 0.85;
}

.lot-row {
  background: rgba(var(--v-theme-on-surface), 0.03);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.06);
}
</style>
