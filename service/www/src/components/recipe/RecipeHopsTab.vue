<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-title class="text-subtitle-1 d-flex align-center">
      Hop Schedule
      <v-spacer />
      <v-btn
        aria-label="Add hop"
        icon="mdi-plus"
        size="small"
        variant="text"
        @click="emit('create')"
      />
    </v-card-title>
    <v-card-text>
      <v-progress-linear
        v-if="loading"
        class="mb-4"
        color="primary"
        indeterminate
      />

      <!-- Desktop Table -->
      <v-table v-if="!loading && ingredients.length > 0" class="data-table d-none d-md-block" density="compact">
        <thead>
          <tr>
            <th>Name</th>
            <th>Amount</th>
            <th>AA%</th>
            <th>Use</th>
            <th>Time</th>
            <th>Notes</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="(group, stage) in groupedByStage" :key="stage">
            <tr class="stage-header">
              <td colspan="7">
                <v-chip class="my-1" color="secondary" size="small" variant="tonal">
                  {{ formatStage(stage as string) }}
                </v-chip>
              </td>
            </tr>
            <tr v-for="ingredient in group" :key="ingredient.uuid">
              <td class="font-weight-medium">{{ ingredient.name }}</td>
              <td>
                <div>{{ formatDisplayAmount(ingredient) }}</div>
                <div v-if="isScaling" class="text-caption text-medium-emphasis">
                  Recipe: {{ formatAmount(ingredient.amount, ingredient.amount_unit) }}
                </div>
              </td>
              <td>{{ formatAlphaAcid(ingredient.alpha_acid_assumed) }}</td>
              <td>
                <v-chip size="x-small" variant="tonal">
                  {{ formatUseType(ingredient.use_type) }}
                </v-chip>
              </td>
              <td>{{ formatTiming(ingredient) }}</td>
              <td class="notes-cell">{{ ingredient.notes ?? '—' }}</td>
              <td class="text-right">
                <v-btn
                  aria-label="Edit"
                  icon="mdi-pencil"
                  size="x-small"
                  variant="text"
                  @click="emit('edit', ingredient)"
                />
                <v-btn
                  aria-label="Delete"
                  color="error"
                  icon="mdi-delete"
                  size="x-small"
                  variant="text"
                  @click="emit('delete', ingredient)"
                />
              </td>
            </tr>
          </template>
        </tbody>
        <tfoot>
          <tr class="font-weight-bold">
            <td>Total</td>
            <td>{{ formatTotalWeight }}</td>
            <td colspan="5" />
          </tr>
        </tfoot>
      </v-table>

      <!-- Mobile Cards -->
      <div v-if="!loading && ingredients.length > 0" class="d-md-none">
        <template v-for="(group, stage) in groupedByStage" :key="stage">
          <v-chip class="mb-2" color="secondary" size="small" variant="tonal">
            {{ formatStage(stage as string) }}
          </v-chip>

          <v-card
            v-for="ingredient in group"
            :key="ingredient.uuid"
            class="mb-2 ingredient-card"
            variant="tonal"
          >
            <v-card-text class="pa-3">
              <div class="d-flex align-center justify-space-between mb-1">
                <span class="font-weight-medium">{{ ingredient.name }}</span>
                <v-chip size="x-small" variant="tonal">
                  {{ formatUseType(ingredient.use_type) }}
                </v-chip>
              </div>
              <div class="d-flex flex-wrap ga-3 text-body-2 mb-1">
                <div>
                  <span>{{ formatDisplayAmount(ingredient) }}</span>
                  <div v-if="isScaling" class="text-caption text-medium-emphasis">
                    Recipe: {{ formatAmount(ingredient.amount, ingredient.amount_unit) }}
                  </div>
                </div>
                <span v-if="ingredient.alpha_acid_assumed" class="text-medium-emphasis">
                  {{ formatAlphaAcid(ingredient.alpha_acid_assumed) }} AA
                </span>
                <span class="text-medium-emphasis">{{ formatTiming(ingredient) }}</span>
              </div>
              <div v-if="ingredient.notes" class="text-caption text-medium-emphasis">
                {{ ingredient.notes }}
              </div>
              <div class="d-flex justify-end mt-2">
                <v-btn
                  aria-label="Edit"
                  icon="mdi-pencil"
                  size="x-small"
                  variant="text"
                  @click="emit('edit', ingredient)"
                />
                <v-btn
                  aria-label="Delete"
                  color="error"
                  icon="mdi-delete"
                  size="x-small"
                  variant="text"
                  @click="emit('delete', ingredient)"
                />
              </div>
            </v-card-text>
          </v-card>
        </template>

        <v-card class="total-card" variant="tonal">
          <v-card-text class="pa-3 d-flex justify-space-between font-weight-bold">
            <span>Total Hops</span>
            <span>{{ formatTotalWeight }}</span>
          </v-card-text>
        </v-card>
      </div>

      <!-- Empty State -->
      <v-alert
        v-if="!loading && ingredients.length === 0"
        density="compact"
        type="info"
        variant="tonal"
      >
        <div class="d-flex align-center justify-space-between flex-wrap ga-2">
          <span>No hops added yet.</span>
          <v-btn
            color="primary"
            size="small"
            variant="text"
            @click="emit('create')"
          >
            Add hop
          </v-btn>
        </div>
      </v-alert>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { RecipeIngredient, RecipeUseStage } from '@/types'
  import { computed } from 'vue'

  const props = withDefaults(defineProps<{
    ingredients: RecipeIngredient[]
    isScaling?: boolean
    loading: boolean
    scaleAmount?: (amount: number, scalingFactor: number) => number
  }>(), {
    isScaling: false,
    scaleAmount: undefined,
  })

  const emit = defineEmits<{
    create: []
    edit: [ingredient: RecipeIngredient]
    delete: [ingredient: RecipeIngredient]
  }>()

  // Group hops by use stage
  const groupedByStage = computed(() => {
    const groups: Record<string, RecipeIngredient[]> = {}
    const stageOrder: RecipeUseStage[] = ['mash', 'boil', 'whirlpool', 'fermentation', 'packaging']

    // Initialize groups in order
    for (const stage of stageOrder) {
      groups[stage] = []
    }

    // Sort ingredients into groups
    for (const ingredient of props.ingredients) {
      const stage = ingredient.use_stage || 'boil'
      if (!groups[stage]) {
        groups[stage] = []
      }
      groups[stage].push(ingredient)
    }

    // Sort within each group by timing (longest first for boil)
    for (const stage of Object.keys(groups)) {
      const group = groups[stage]
      if (group) {
        group.sort((a, b) => {
          const aTime = a.timing_duration_minutes ?? 0
          const bTime = b.timing_duration_minutes ?? 0
          return bTime - aTime
        })
      }
    }

    // Remove empty groups
    for (const stage of Object.keys(groups)) {
      const group = groups[stage]
      if (group && group.length === 0) {
        delete groups[stage]
      }
    }

    return groups
  })

  /** Get the display amount for an ingredient (scaled or original). */
  function getDisplayAmount (ingredient: RecipeIngredient): number {
    if (props.isScaling && props.scaleAmount) {
      return props.scaleAmount(ingredient.amount, ingredient.scaling_factor)
    }
    return ingredient.amount
  }

  const totalWeight = computed(() => {
    return props.ingredients.reduce((sum, i) => sum + getDisplayAmount(i), 0)
  })

  const formatTotalWeight = computed(() => {
    if (props.ingredients.length === 0) return '—'
    const unit = props.ingredients[0]?.amount_unit ?? 'oz'
    return `${totalWeight.value.toFixed(2)} ${unit}`
  })

  function formatAmount (amount: number, unit: string): string {
    return `${amount.toFixed(2)} ${unit}`
  }

  function formatDisplayAmount (ingredient: RecipeIngredient): string {
    return formatAmount(getDisplayAmount(ingredient), ingredient.amount_unit)
  }

  function formatAlphaAcid (aa: number | null): string {
    if (aa === null || aa === undefined) return '—'
    return `${aa.toFixed(1)}%`
  }

  function formatStage (stage: string): string {
    const labels: Record<string, string> = {
      mash: 'Mash Hop',
      boil: 'Boil',
      whirlpool: 'Whirlpool',
      fermentation: 'Dry Hop',
      packaging: 'Packaging',
    }
    return labels[stage] ?? stage.charAt(0).toUpperCase() + stage.slice(1)
  }

  function formatUseType (useType: string | null): string {
    if (!useType) return 'Other'
    const labels: Record<string, string> = {
      bittering: 'Bittering',
      flavor: 'Flavor',
      aroma: 'Aroma',
      dry_hop: 'Dry Hop',
      other: 'Other',
    }
    return labels[useType] ?? useType.charAt(0).toUpperCase() + useType.slice(1).replace(/_/g, ' ')
  }

  function formatTiming (ingredient: RecipeIngredient): string {
    if (ingredient.use_stage === 'fermentation') {
      if (ingredient.timing_duration_minutes) {
        const days = Math.round(ingredient.timing_duration_minutes / 1440)
        return `${days} day${days === 1 ? '' : 's'}`
      }
      return '—'
    }

    if (ingredient.timing_duration_minutes !== null && ingredient.timing_duration_minutes !== undefined) {
      return `${ingredient.timing_duration_minutes} min`
    }

    if (ingredient.timing_temperature_c !== null && ingredient.timing_temperature_c !== undefined) {
      // Convert to F for display
      const tempF = Math.round(ingredient.timing_temperature_c * 9 / 5 + 32)
      return `@ ${tempF}°F`
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

.data-table :deep(tfoot td) {
  border-top: 2px solid rgba(var(--v-theme-on-surface), 0.12);
}

.stage-header td {
  background: rgba(var(--v-theme-surface), 0.3);
  padding-top: 8px !important;
  padding-bottom: 4px !important;
}

.notes-cell {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ingredient-card {
  background: rgba(var(--v-theme-surface), 0.5);
}

.total-card {
  background: rgba(var(--v-theme-primary), 0.08);
}
</style>
