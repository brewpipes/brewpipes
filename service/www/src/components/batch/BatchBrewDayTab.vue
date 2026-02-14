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
        <div class="text-overline text-medium-emphasis mb-1">Recipe</div>
        <div class="text-h6 mb-4">{{ displayRecipeName }}</div>

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
                <td>{{ formatAmount(item.amount, item.amount_unit) }}</td>
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
                    :color="getStockColor(item.ingredient_uuid, item.amount, item.amount_unit)"
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
                  <span>{{ formatAmount(item.amount, item.amount_unit) }}</span>
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
                    :color="getStockColor(item.ingredient_uuid, item.amount, item.amount_unit)"
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
                <td>{{ formatAmount(item.amount, item.amount_unit) }}</td>
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
                    :color="getStockColor(item.ingredient_uuid, item.amount, item.amount_unit)"
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
                  <span>{{ formatAmount(item.amount, item.amount_unit) }}</span>
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
                    :color="getStockColor(item.ingredient_uuid, item.amount, item.amount_unit)"
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
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { RecipeIngredient, RecipeUseStage, RecipeUseType, StockLevel } from '@/types'
  import { computed, ref, watch } from 'vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  const props = defineProps<{
    recipeUuid: string | null
    recipeName: string | null
  }>()

  const { getRecipeIngredients, getRecipe } = useProductionApi()
  const { getStockLevels } = useInventoryApi()
  const { formatTemperaturePreferred } = useUnitPreferences()

  // State
  const loading = ref(false)
  const fetchError = ref<string | null>(null)
  const recipeIngredients = ref<RecipeIngredient[]>([])
  const stockLevels = ref<StockLevel[]>([])
  const stockUnavailable = ref(false)
  const resolvedRecipeName = ref<string | null>(null)

  // Computed: recipe display name
  const displayRecipeName = computed(() =>
    props.recipeName ?? resolvedRecipeName.value ?? 'Recipe',
  )

  // Computed: stock lookup map (ingredient UUID → StockLevel)
  const stockMap = computed(() => {
    const map = new Map<string, StockLevel>()
    for (const level of stockLevels.value) {
      map.set(level.ingredient_uuid, level)
    }
    return map
  })

  // Computed: split ingredients into "today" and "later"
  const TODAY_STAGES: RecipeUseStage[] = ['mash', 'boil', 'whirlpool']
  const LATER_STAGES: RecipeUseStage[] = ['fermentation', 'packaging']

  const neededToday = computed(() =>
    sortIngredients(
      recipeIngredients.value.filter(i => TODAY_STAGES.includes(i.use_stage)),
    ),
  )

  const neededLater = computed(() =>
    sortIngredients(
      recipeIngredients.value.filter(i => LATER_STAGES.includes(i.use_stage)),
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
      }
      // Recipe name fetch failure is non-critical — we fall back to "Recipe"
    } catch {
      fetchError.value = 'Failed to load pick list data. Please try refreshing.'
    } finally {
      loading.value = false
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
</style>
