<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-text>
      <v-progress-linear
        v-if="loading"
        class="mb-4"
        color="primary"
        indeterminate
      />

      <v-alert
        v-if="error"
        class="mb-4"
        density="comfortable"
        type="error"
        variant="tonal"
      >
        {{ error }}
      </v-alert>

      <template v-if="costs && !loading">
        <!-- Summary Cards -->
        <v-row class="mb-4">
          <v-col cols="12" sm="4">
            <div class="metric-card">
              <div class="metric-label">Total Ingredient Cost</div>
              <div class="metric-value">
                {{ formatCents(costs.totals.total_cost_cents) }}
              </div>
              <div v-if="costs.currency === 'MIXED'" class="text-caption text-warning mt-1">
                Mixed currencies
              </div>
            </div>
          </v-col>
          <v-col cols="12" sm="4">
            <div class="metric-card">
              <div class="metric-label">Cost per BBL</div>
              <div class="metric-value">
                {{ costs.totals.cost_per_bbl_cents !== null
                  ? `${formatCents(costs.totals.cost_per_bbl_cents)}/bbl`
                  : 'N/A' }}
              </div>
              <div v-if="costs.totals.batch_volume_bbl !== null" class="text-caption text-medium-emphasis mt-1">
                {{ formatVolumePreferred(costs.totals.batch_volume_bbl, 'bbl') }} batch
              </div>
            </div>
          </v-col>
          <v-col cols="12" sm="4">
            <div class="metric-card" :class="costStatusClass">
              <div class="metric-label">Costing Status</div>
              <div class="metric-value d-flex align-center justify-center ga-2">
                <v-icon
                  :color="costs.totals.cost_complete ? 'success' : 'warning'"
                  :icon="costs.totals.cost_complete ? 'mdi-check-circle' : 'mdi-alert-circle'"
                  size="20"
                />
                {{ costStatusLabel }}
              </div>
            </div>
          </v-col>
        </v-row>

        <v-divider class="mb-4" />

        <!-- Actual vs. Target Comparison -->
        <template v-if="recipe">
          <div class="text-overline text-medium-emphasis mb-2">Actual vs. Target</div>
          <v-row class="mb-4">
            <v-col cols="12">
              <!-- Desktop table (md+) -->
              <v-table class="data-table d-none d-md-block" density="compact">
                <thead>
                  <tr>
                    <th>Metric</th>
                    <th>Target</th>
                    <th>Actual</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="row in comparisonRows" :key="row.metric">
                    <td class="font-weight-medium">{{ row.metric }}</td>
                    <td>{{ row.target }}</td>
                    <td>{{ row.actual }}</td>
                    <td>
                      <v-chip
                        v-if="row.status !== 'unavailable'"
                        :color="row.status === 'in_range' ? 'success' : 'warning'"
                        size="small"
                        variant="tonal"
                      >
                        <v-icon
                          :icon="row.status === 'in_range' ? 'mdi-check' : 'mdi-alert'"
                          size="14"
                          start
                        />
                        {{ row.statusLabel }}
                      </v-chip>
                      <span v-else class="text-medium-emphasis">—</span>
                    </td>
                  </tr>
                </tbody>
              </v-table>

              <!-- Mobile cards (< md) -->
              <div class="d-md-none">
                <v-card
                  v-for="row in comparisonRows"
                  :key="row.metric"
                  class="mb-2"
                  density="compact"
                  variant="tonal"
                >
                  <v-card-text class="pa-3">
                    <div class="d-flex align-center justify-space-between mb-1">
                      <span class="text-subtitle-2 font-weight-medium">{{ row.metric }}</span>
                      <v-chip
                        v-if="row.status !== 'unavailable'"
                        :color="row.status === 'in_range' ? 'success' : 'warning'"
                        size="x-small"
                        variant="tonal"
                      >
                        <v-icon
                          :icon="row.status === 'in_range' ? 'mdi-check' : 'mdi-alert'"
                          size="12"
                          start
                        />
                        {{ row.statusLabel }}
                      </v-chip>
                      <span v-else class="text-caption text-medium-emphasis">—</span>
                    </div>
                    <div class="d-flex ga-4">
                      <div>
                        <div class="text-caption text-medium-emphasis">Target</div>
                        <div class="text-body-2">{{ row.target }}</div>
                      </div>
                      <div>
                        <div class="text-caption text-medium-emphasis">Actual</div>
                        <div class="text-body-2">{{ row.actual }}</div>
                      </div>
                    </div>
                  </v-card-text>
                </v-card>
              </div>
            </v-col>
          </v-row>

          <v-divider class="mb-4" />
        </template>

        <v-alert
          v-else-if="!recipe && !recipeLoading && summary?.recipe_name === null"
          class="mb-4"
          density="compact"
          type="info"
          variant="tonal"
        >
          No recipe linked — target comparison unavailable.
        </v-alert>

        <!-- Ingredient Cost Table -->
        <div class="text-overline text-medium-emphasis mb-2">Ingredient Costs</div>

        <template v-if="costs.line_items.length > 0">
          <!-- Desktop table (md+) -->
          <v-table class="data-table d-none d-md-block" density="compact">
            <thead>
              <tr>
                <th>Ingredient</th>
                <th>Category</th>
                <th>Lot Code</th>
                <th class="text-right">Amount Used</th>
                <th class="text-right">Unit Cost</th>
                <th class="text-right">Total Cost</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in sortedLineItems" :key="item.addition_uuid">
                <td>{{ item.ingredient_name ?? 'Unknown' }}</td>
                <td>
                  <v-chip
                    v-if="item.ingredient_category"
                    :color="getCategoryColor(item.ingredient_category)"
                    size="x-small"
                    variant="tonal"
                  >
                    {{ formatCategory(item.ingredient_category) }}
                  </v-chip>
                  <span v-else class="text-medium-emphasis">—</span>
                </td>
                <td>
                  <span class="text-caption">{{ item.lot_code ?? '—' }}</span>
                </td>
                <td class="text-right">
                  {{ formatAmountPreferred(item.amount_used, item.amount_unit) }}
                </td>
                <td class="text-right">
                  <template v-if="item.cost_source === 'unavailable'">
                    <v-tooltip location="top" text="Cost unavailable — no linked purchase order">
                      <template #activator="{ props: tooltipProps }">
                        <span v-bind="tooltipProps" class="text-medium-emphasis">—</span>
                      </template>
                    </v-tooltip>
                  </template>
                  <template v-else>
                    {{ formatUnitCost(item.unit_cost_cents, item.unit_cost_unit) }}
                  </template>
                </td>
                <td class="text-right">
                  <template v-if="item.cost_source === 'unavailable'">
                    <v-tooltip location="top" text="Cost unavailable — no linked purchase order">
                      <template #activator="{ props: tooltipProps }">
                        <span v-bind="tooltipProps" class="text-medium-emphasis">—</span>
                      </template>
                    </v-tooltip>
                  </template>
                  <template v-else>
                    {{ formatCents(item.cost_cents) }}
                  </template>
                </td>
              </tr>
            </tbody>
          </v-table>

          <!-- Mobile cards (< md) -->
          <div class="d-md-none">
            <v-card
              v-for="item in sortedLineItems"
              :key="item.addition_uuid"
              class="mb-2"
              density="compact"
              variant="tonal"
            >
              <v-card-text class="pa-3">
                <div class="d-flex align-center justify-space-between mb-1">
                  <span class="text-subtitle-2 font-weight-medium">
                    {{ item.ingredient_name ?? 'Unknown' }}
                  </span>
                  <v-chip
                    v-if="item.ingredient_category"
                    :color="getCategoryColor(item.ingredient_category)"
                    size="x-small"
                    variant="tonal"
                  >
                    {{ formatCategory(item.ingredient_category) }}
                  </v-chip>
                </div>
                <div class="d-flex justify-space-between text-body-2">
                  <span class="text-medium-emphasis">
                    {{ formatAmountPreferred(item.amount_used, item.amount_unit) }}
                  </span>
                  <span class="font-weight-medium">
                    <template v-if="item.cost_source === 'unavailable'">
                      <span class="text-medium-emphasis">Cost N/A</span>
                    </template>
                    <template v-else>
                      {{ formatCents(item.cost_cents) }}
                    </template>
                  </span>
                </div>
              </v-card-text>
            </v-card>
          </div>
        </template>

        <!-- Empty state: no line items at all -->
        <v-alert
          v-else-if="costs.uncosted_additions.length === 0"
          density="compact"
          type="info"
          variant="tonal"
        >
          No ingredients recorded for this batch.
        </v-alert>

        <!-- Empty state: additions exist but no costs -->
        <v-alert
          v-else
          density="compact"
          type="info"
          variant="tonal"
        >
          Cost data unavailable — ingredients are not linked to purchase orders.
        </v-alert>

        <!-- Uncosted Additions -->
        <template v-if="costs.uncosted_additions.length > 0">
          <v-divider class="my-4" />
          <div class="text-overline text-medium-emphasis mb-2">Uncosted Additions</div>
          <v-alert
            class="mb-3"
            density="compact"
            icon="mdi-information-outline"
            type="warning"
            variant="tonal"
          >
            {{ costs.uncosted_additions.length }}
            {{ costs.uncosted_additions.length === 1 ? 'addition' : 'additions' }}
            could not be costed.
          </v-alert>
          <!-- Desktop table (md+) -->
          <v-table class="data-table d-none d-md-block" density="compact">
            <thead>
              <tr>
                <th>Type</th>
                <th>Amount</th>
                <th>Reason</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in costs.uncosted_additions" :key="item.addition_uuid">
                <td>{{ formatAdditionType(item.addition_type) }}</td>
                <td>{{ formatAmountPreferred(item.amount_used, item.amount_unit) }}</td>
                <td>
                  <v-chip color="grey" size="x-small" variant="tonal">
                    {{ formatUncostedReason(item.reason) }}
                  </v-chip>
                </td>
              </tr>
            </tbody>
          </v-table>

          <!-- Mobile cards (< md) -->
          <div class="d-md-none">
            <v-card
              v-for="item in costs.uncosted_additions"
              :key="item.addition_uuid"
              class="mb-2"
              density="compact"
              variant="tonal"
            >
              <v-card-text class="pa-3">
                <div class="d-flex align-center justify-space-between mb-1">
                  <span class="text-subtitle-2 font-weight-medium">
                    {{ formatAdditionType(item.addition_type) }}
                  </span>
                  <v-chip color="grey" size="x-small" variant="tonal">
                    {{ formatUncostedReason(item.reason) }}
                  </v-chip>
                </div>
                <div class="text-body-2 text-medium-emphasis">
                  {{ formatAmountPreferred(item.amount_used, item.amount_unit) }}
                </div>
              </v-card-text>
            </v-card>
          </div>
        </template>
      </template>

      <!-- No data and not loading -->
      <v-alert
        v-else-if="!loading && !error"
        density="comfortable"
        type="info"
        variant="tonal"
      >
        Cost data not available.
      </v-alert>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type {
    BatchCostsResponse,
    BatchSummary,
    CostLineItem,
    Recipe,
    UncostedReason,
  } from '@/types'
  import { computed, ref, watch } from 'vue'
  import { useAdditionTypeFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { convertVolume } from '@/composables/useUnitConversion'
  import { useUnitPreferences, normalizeVolumeUnit } from '@/composables/useUnitPreferences'

  const props = defineProps<{
    batchUuid: string
    summary: BatchSummary | null
    recipe: Recipe | null
    recipeLoading: boolean
  }>()

  const { getBatchCosts } = useProductionApi()
  const { formatAmountPreferred, formatVolumePreferred, formatGravityPreferred } = useUnitPreferences()
  const { formatAdditionType } = useAdditionTypeFormatters()

  const costs = ref<BatchCostsResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Load costs when component mounts or batchUuid changes
  watch(
    () => props.batchUuid,
    async (uuid) => {
      if (uuid) {
        await loadCosts(uuid)
      }
    },
    { immediate: true },
  )

  async function loadCosts (batchUuid: string) {
    loading.value = true
    error.value = null
    try {
      costs.value = await getBatchCosts(batchUuid)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load cost data'
      error.value = message
    } finally {
      loading.value = false
    }
  }

  // Sort line items by cost descending (most expensive first)
  const sortedLineItems = computed<CostLineItem[]>(() => {
    if (!costs.value) return []
    return [...costs.value.line_items].sort((a, b) => {
      const costA = a.cost_cents ?? 0
      const costB = b.cost_cents ?? 0
      return costB - costA
    })
  })

  // Costing status label
  const costStatusLabel = computed(() => {
    if (!costs.value) return '—'
    if (costs.value.totals.cost_complete) return 'Complete'
    const { costed_line_count, uncosted_line_count } = costs.value.totals
    const total = costed_line_count + uncosted_line_count
    return `${costed_line_count} of ${total} items costed`
  })

  // Costing status card class
  const costStatusClass = computed(() => {
    if (!costs.value) return ''
    return costs.value.totals.cost_complete
      ? 'metric-card--success'
      : 'metric-card--warning'
  })

  // ==================== Actual vs. Target Comparison ====================

  interface ComparisonRow {
    metric: string
    target: string
    actual: string
    status: 'in_range' | 'above' | 'below' | 'unavailable'
    statusLabel: string
  }

  const comparisonRows = computed<ComparisonRow[]>(() => {
    const recipe = props.recipe
    const summary = props.summary
    if (!recipe || !summary) return []

    const rows: ComparisonRow[] = []

    // OG
    rows.push(buildGravityRow(
      'OG',
      recipe.target_og,
      recipe.target_og_min,
      recipe.target_og_max,
      summary.original_gravity,
    ))

    // FG
    rows.push(buildGravityRow(
      'FG',
      recipe.target_fg,
      recipe.target_fg_min,
      recipe.target_fg_max,
      summary.final_gravity,
    ))

    // ABV
    rows.push(buildPercentRow(
      'ABV',
      recipe.target_abv,
      null,
      null,
      summary.abv,
    ))

    // IBU
    rows.push(buildNumericRow(
      'IBU',
      recipe.target_ibu,
      recipe.target_ibu_min,
      recipe.target_ibu_max,
      summary.ibu,
    ))

    // Volume
    rows.push(buildVolumeRow(
      'Volume',
      recipe.batch_size,
      recipe.batch_size_unit,
      summary.current_volume_bbl,
    ))

    return rows
  })

  function buildGravityRow (
    metric: string,
    target: number | null,
    min: number | null,
    max: number | null,
    actual: number | null,
  ): ComparisonRow {
    const targetStr = formatGravityTarget(target, min, max)
    const actualStr = actual !== null ? formatGravityPreferred(actual, 'sg') : '—'
    const status = evaluateStatus(actual, target, min, max)
    return { metric, target: targetStr, actual: actualStr, status, statusLabel: formatStatusLabel(status) }
  }

  function buildPercentRow (
    metric: string,
    target: number | null,
    min: number | null,
    max: number | null,
    actual: number | null,
  ): ComparisonRow {
    const targetStr = target !== null ? `${target.toFixed(1)}%` : '—'
    const actualStr = actual !== null ? `${actual.toFixed(1)}%` : '—'
    const status = evaluateStatus(actual, target, min, max)
    return { metric, target: targetStr, actual: actualStr, status, statusLabel: formatStatusLabel(status) }
  }

  function buildNumericRow (
    metric: string,
    target: number | null,
    min: number | null,
    max: number | null,
    actual: number | null,
  ): ComparisonRow {
    const targetStr = formatNumericTarget(target, min, max)
    const actualStr = actual !== null ? String(Math.round(actual)) : '—'
    const status = evaluateStatus(actual, target, min, max)
    return { metric, target: targetStr, actual: actualStr, status, statusLabel: formatStatusLabel(status) }
  }

  function buildVolumeRow (
    metric: string,
    targetSize: number | null,
    targetUnit: string | null,
    actualBbl: number | null,
  ): ComparisonRow {
    const resolvedTargetUnit = targetUnit ?? 'bbl'
    const targetStr = targetSize !== null
      ? formatVolumePreferred(targetSize, normalizeVolumeUnit(resolvedTargetUnit))
      : '—'
    const actualStr = actualBbl !== null
      ? formatVolumePreferred(actualBbl, 'bbl')
      : '—'

    // Convert target to bbl for apples-to-apples comparison with actualBbl
    let status: ComparisonRow['status'] = 'unavailable'
    if (targetSize !== null && actualBbl !== null) {
      const targetBbl = convertVolume(targetSize, normalizeVolumeUnit(resolvedTargetUnit), 'bbl')
      if (targetBbl !== null) {
        const tolerance = targetBbl * 0.05
        if (actualBbl >= targetBbl - tolerance && actualBbl <= targetBbl + tolerance) {
          status = 'in_range'
        } else if (actualBbl < targetBbl - tolerance) {
          status = 'below'
        } else {
          status = 'above'
        }
      }
    }

    return { metric, target: targetStr, actual: actualStr, status, statusLabel: formatStatusLabel(status) }
  }

  function evaluateStatus (
    actual: number | null,
    target: number | null,
    min: number | null,
    max: number | null,
  ): ComparisonRow['status'] {
    if (actual === null || target === null) return 'unavailable'

    // If we have a range, use it
    if (min !== null && max !== null) {
      if (actual >= min && actual <= max) return 'in_range'
      if (actual < min) return 'below'
      return 'above'
    }

    // No range: use 5% tolerance around target
    const tolerance = Math.abs(target) * 0.05
    if (actual >= target - tolerance && actual <= target + tolerance) return 'in_range'
    if (actual < target - tolerance) return 'below'
    return 'above'
  }

  function formatStatusLabel (status: ComparisonRow['status']): string {
    switch (status) {
      case 'in_range': return 'In range'
      case 'above': return 'Above target'
      case 'below': return 'Below target'
      default: return '—'
    }
  }

  function formatGravityTarget (target: number | null, min: number | null, max: number | null): string {
    if (target === null) return '—'
    const base = formatGravityPreferred(target, 'sg')
    if (min !== null && max !== null) {
      return `${base} (${formatGravityPreferred(min, 'sg')}–${formatGravityPreferred(max, 'sg')})`
    }
    return base
  }

  function formatNumericTarget (target: number | null, min: number | null, max: number | null): string {
    if (target === null) return '—'
    const base = String(Math.round(target))
    if (min !== null && max !== null) {
      return `${base} (${Math.round(min)}–${Math.round(max)})`
    }
    return base
  }

  // ==================== Currency Formatting ====================

  function formatCents (cents: number | null): string {
    if (cents === null) return '—'
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(cents / 100)
  }

  function formatUnitCost (cents: number | null, unit: string | null): string {
    if (cents === null) return '—'
    const formatted = formatCents(cents)
    return unit ? `${formatted}/${unit}` : formatted
  }

  // ==================== Category Formatting ====================

  const CATEGORY_COLORS: Record<string, string> = {
    malt: 'amber',
    fermentable: 'amber',
    hop: 'green',
    yeast: 'purple',
    adjunct: 'blue-grey',
    water_chem: 'cyan',
    chemical: 'cyan',
    salt: 'cyan',
    gas: 'grey',
    other: 'grey',
  }

  function getCategoryColor (category: string): string {
    return CATEGORY_COLORS[category.toLowerCase()] ?? 'grey'
  }

  function formatCategory (category: string): string {
    return category.charAt(0).toUpperCase() + category.slice(1).replace(/_/g, ' ')
  }

  function formatUncostedReason (reason: UncostedReason): string {
    switch (reason) {
      case 'no_inventory_lot': return 'No inventory lot'
      case 'non_ingredient': return 'Non-ingredient'
      default: return reason
    }
  }
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.metric-card {
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(var(--v-theme-surface), 0.5);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  text-align: center;
}

.metric-card--success {
  background: rgba(var(--v-theme-success), 0.08);
  border-color: rgba(var(--v-theme-success), 0.2);
}

.metric-card--warning {
  background: rgba(var(--v-theme-warning), 0.08);
  border-color: rgba(var(--v-theme-warning), 0.2);
}

.metric-label {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-bottom: 4px;
}

.metric-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.87);
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
</style>
