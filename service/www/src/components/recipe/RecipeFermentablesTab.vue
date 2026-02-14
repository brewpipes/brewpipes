<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-title class="text-subtitle-1 d-flex align-center">
      Grain Bill
      <v-spacer />
      <v-btn
        aria-label="Add fermentable"
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
            <th class="text-right">%</th>
            <th>Use</th>
            <th>Notes</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ingredient in sortedIngredients" :key="ingredient.uuid">
            <td class="font-weight-medium">{{ ingredient.name }}</td>
            <td>
              <div>{{ formatDisplayAmount(ingredient) }}</div>
              <div v-if="isScaling" class="text-caption text-medium-emphasis">
                Recipe: {{ formatAmount(ingredient.amount, ingredient.amount_unit) }}
              </div>
            </td>
            <td class="text-right">{{ formatPercent(ingredient, totalWeight) }}</td>
            <td>
              <v-chip size="x-small" variant="tonal">
                {{ formatUseType(ingredient.use_type) }}
              </v-chip>
            </td>
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
        </tbody>
        <tfoot>
          <tr class="font-weight-bold">
            <td>Total</td>
            <td>{{ formatTotalWeight }}</td>
            <td class="text-right">100%</td>
            <td colspan="3" />
          </tr>
        </tfoot>
      </v-table>

      <!-- Mobile Cards -->
      <div v-if="!loading && ingredients.length > 0" class="d-md-none">
        <v-card
          v-for="ingredient in sortedIngredients"
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
            <div class="d-flex align-center justify-space-between text-body-2">
              <div>
                <span>{{ formatDisplayAmount(ingredient) }}</span>
                <div v-if="isScaling" class="text-caption text-medium-emphasis">
                  Recipe: {{ formatAmount(ingredient.amount, ingredient.amount_unit) }}
                </div>
              </div>
              <span class="text-medium-emphasis">{{ formatPercent(ingredient, totalWeight) }}</span>
            </div>
            <div v-if="ingredient.notes" class="text-caption text-medium-emphasis mt-1">
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

        <v-card class="total-card" variant="tonal">
          <v-card-text class="pa-3 d-flex justify-space-between font-weight-bold">
            <span>Total</span>
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
          <span>No fermentables added yet.</span>
          <v-btn
            color="primary"
            size="small"
            variant="text"
            @click="emit('create')"
          >
            Add fermentable
          </v-btn>
        </div>
      </v-alert>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { RecipeIngredient } from '@/types'
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

  const sortedIngredients = computed(() =>
    [...props.ingredients].toSorted((a, b) => a.sort_order - b.sort_order),
  )

  /** Get the display amount for an ingredient (scaled or original). */
  function getDisplayAmount (ingredient: RecipeIngredient): number {
    if (props.isScaling && props.scaleAmount) {
      return props.scaleAmount(ingredient.amount, ingredient.scaling_factor)
    }
    return ingredient.amount
  }

  const totalWeight = computed(() => {
    return props.ingredients.reduce((sum, i) => {
      return sum + getDisplayAmount(i)
    }, 0)
  })

  const formatTotalWeight = computed(() => {
    if (props.ingredients.length === 0) return '—'
    const unit = props.ingredients[0]?.amount_unit ?? 'lb'
    return `${totalWeight.value.toFixed(2)} ${unit}`
  })

  function formatAmount (amount: number, unit: string): string {
    return `${amount.toFixed(2)} ${unit}`
  }

  function formatDisplayAmount (ingredient: RecipeIngredient): string {
    return formatAmount(getDisplayAmount(ingredient), ingredient.amount_unit)
  }

  function formatPercent (ingredient: RecipeIngredient, total: number): string {
    if (total === 0) return '—'
    const displayAmount = getDisplayAmount(ingredient)
    const percent = (displayAmount / total) * 100
    return `${percent.toFixed(1)}%`
  }

  function formatUseType (useType: string | null): string {
    if (!useType) return 'Other'
    const labels: Record<string, string> = {
      base: 'Base',
      specialty: 'Specialty',
      adjunct: 'Adjunct',
      sugar: 'Sugar',
      other: 'Other',
    }
    return labels[useType] ?? useType.charAt(0).toUpperCase() + useType.slice(1)
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

.notes-cell {
  max-width: 200px;
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
