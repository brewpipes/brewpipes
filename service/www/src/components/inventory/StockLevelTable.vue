<template>
  <!-- Desktop view: Data table with expandable rows -->
  <div class="d-none d-md-block">
    <v-data-table
      v-model:expanded="expanded"
      class="stock-level-table"
      density="compact"
      :headers="tableHeaders"
      hover
      item-value="ingredient_uuid"
      :items="items"
      :loading="loading"
      :row-props="getRowProps"
      show-expand
    >
      <template #item.ingredient_name="{ item }">
        <span :class="{ 'text-medium-emphasis': isZeroStock(item) }">
          {{ item.ingredient_name }}
        </span>
      </template>

      <template #item.category="{ item }">
        <v-chip density="compact" size="small" variant="tonal">
          {{ formatCategory(item.category) }}
        </v-chip>
      </template>

      <template #item.total_on_hand="{ item }">
        <span :class="stockClass(item)">
          <v-icon
            v-if="isZeroStock(item)"
            class="mr-1"
            color="warning"
            icon="mdi-alert-outline"
            size="small"
          />
          {{ formatAmountPreferred(item.total_on_hand, item.default_unit) }}
        </span>
      </template>

      <template #item.locations="{ item }">
        <span class="text-medium-emphasis">
          {{ item.locations.length }} location{{ item.locations.length !== 1 ? 's' : '' }}
        </span>
      </template>

      <template #expanded-row="{ columns, item }">
        <tr>
          <td class="pa-0" :colspan="columns.length">
            <v-table class="location-breakdown-table" density="compact">
              <thead>
                <tr class="bg-grey-lighten-4">
                  <th class="pl-12">Location</th>
                  <th>Quantity</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="loc in item.locations" :key="loc.location_uuid">
                  <td class="pl-12">{{ loc.location_name }}</td>
                  <td :class="{ 'text-medium-emphasis': loc.quantity === 0 }">
                    {{ formatAmountPreferred(loc.quantity, item.default_unit) }}
                  </td>
                </tr>
                <tr v-if="item.locations.length === 0">
                  <td class="text-medium-emphasis pl-12" colspan="2">
                    No location data available
                  </td>
                </tr>
              </tbody>
            </v-table>
          </td>
        </tr>
      </template>

      <template #no-data>
        <div class="text-center py-4 text-medium-emphasis">
          No {{ categoryLabel.toLowerCase() }} in stock.
        </div>
      </template>
    </v-data-table>
  </div>

  <!-- Mobile view: Card-based layout -->
  <div class="d-md-none">
    <v-progress-linear v-if="loading" color="primary" indeterminate />

    <div v-if="!loading && items.length === 0" class="text-center py-4 text-medium-emphasis">
      No {{ categoryLabel.toLowerCase() }} in stock.
    </div>

    <v-card
      v-for="item in items"
      :key="item.ingredient_uuid"
      :class="['mb-3', { 'stock-card-highlight': highlightedUuid === item.ingredient_uuid }]"
      :data-ingredient-card="item.ingredient_uuid"
      variant="outlined"
    >
      <v-card-title class="d-flex align-center py-2">
        <span :class="{ 'text-medium-emphasis': isZeroStock(item) }">
          {{ item.ingredient_name }}
        </span>
        <v-spacer />
        <v-chip
          v-if="showCategory"
          class="ml-2"
          density="compact"
          size="small"
          variant="tonal"
        >
          {{ formatCategory(item.category) }}
        </v-chip>
      </v-card-title>

      <v-card-text class="pt-0">
        <div class="d-flex align-center mb-2">
          <span class="text-body-2 text-medium-emphasis mr-2">On Hand:</span>
          <span class="text-body-2 font-weight-medium" :class="stockClass(item)">
            <v-icon
              v-if="isZeroStock(item)"
              class="mr-1"
              color="warning"
              icon="mdi-alert-outline"
              size="small"
            />
            {{ formatAmountPreferred(item.total_on_hand, item.default_unit) }}
          </span>
        </div>

        <v-divider class="my-2" />

        <div class="text-caption text-medium-emphasis mb-1">Locations:</div>
        <div v-if="item.locations.length === 0" class="text-caption text-medium-emphasis">
          No location data available
        </div>
        <div
          v-for="loc in item.locations"
          :key="loc.location_uuid"
          class="d-flex justify-space-between text-body-2 py-1"
        >
          <span>{{ loc.location_name }}</span>
          <span :class="{ 'text-medium-emphasis': loc.quantity === 0 }">
            {{ formatAmountPreferred(loc.quantity, item.default_unit) }}
          </span>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts" setup>
  import type { StockLevel } from '@/types'
  import { computed, nextTick, ref, watch } from 'vue'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  const props = defineProps<{
    items: StockLevel[]
    loading: boolean
    categoryLabel: string
    showCategory?: boolean
    highlightIngredientUuid?: string
  }>()

  const { formatAmountPreferred } = useUnitPreferences()

  const expanded = ref<string[]>([])
  const highlightedUuid = ref<string | null>(null)

  function getRowProps ({ item }: { item: StockLevel }) {
    return {
      'data-ingredient-uuid': item.ingredient_uuid,
      class: highlightedUuid.value === item.ingredient_uuid ? 'stock-row-highlight' : undefined,
    }
  }

  // Watch for items + highlightIngredientUuid to scroll and highlight
  watch(
    () => [props.items, props.highlightIngredientUuid] as const,
    async ([items, uuid]) => {
      if (!uuid || items.length === 0) return
      const match = items.find(item => item.ingredient_uuid === uuid)
      if (!match) return

      highlightedUuid.value = uuid
      await nextTick()

      // Small delay to ensure Vuetify has rendered the table rows
      setTimeout(() => {
        // Try desktop table first
        const desktopRow = document.querySelector(
          `tr[data-ingredient-uuid="${uuid}"]`,
        ) as HTMLElement | null
        if (desktopRow) {
          desktopRow.scrollIntoView({ behavior: 'smooth', block: 'center' })
        }

        // Try mobile card
        const mobileCard = document.querySelector(
          `[data-ingredient-card="${uuid}"]`,
        ) as HTMLElement | null
        if (mobileCard) {
          mobileCard.scrollIntoView({ behavior: 'smooth', block: 'center' })
        }

        // Remove highlight after animation completes
        setTimeout(() => {
          highlightedUuid.value = null
        }, 2000)
      }, 100)
    },
    { immediate: true },
  )

  const tableHeaders = computed(() => {
    const headers = [
      { title: 'Ingredient', key: 'ingredient_name', sortable: true },
    ]

    if (props.showCategory) {
      headers.push({ title: 'Category', key: 'category', sortable: true })
    }

    headers.push(
      { title: 'On Hand', key: 'total_on_hand', sortable: true },
      { title: 'Locations', key: 'locations', sortable: false },
    )

    return headers
  })

  function isZeroStock (item: StockLevel): boolean {
    return item.total_on_hand === 0
  }

  function stockClass (item: StockLevel): string {
    if (isZeroStock(item)) {
      return 'text-medium-emphasis'
    }
    return ''
  }

  function formatCategory (category: string): string {
    const labels: Record<string, string> = {
      fermentable: 'Malt',
      hop: 'Hops',
      yeast: 'Yeast',
      adjunct: 'Adjunct',
      salt: 'Salt',
      chemical: 'Chemical',
      gas: 'Gas',
      other: 'Other',
    }
    return labels[category] ?? category
  }
</script>

<style scoped>
.stock-level-table :deep(tr) {
  cursor: pointer;
}

.location-breakdown-table {
  background: rgba(var(--v-theme-surface-variant), 0.3);
}

.location-breakdown-table th {
  font-weight: 500;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

/* Highlight animation for desktop table rows */
.stock-level-table :deep(tr.stock-row-highlight) {
  animation: row-highlight-pulse 2s ease-in-out;
}

/* Highlight animation for mobile cards */
.stock-card-highlight {
  animation: card-highlight-pulse 2s ease-in-out;
}

@keyframes row-highlight-pulse {
  0% { background-color: rgba(var(--v-theme-warning), 0.3); }
  50% { background-color: rgba(var(--v-theme-warning), 0.15); }
  100% { background-color: transparent; }
}

@keyframes card-highlight-pulse {
  0% { border-color: rgb(var(--v-theme-warning)); box-shadow: 0 0 8px rgba(var(--v-theme-warning), 0.4); }
  50% { border-color: rgb(var(--v-theme-warning)); box-shadow: 0 0 4px rgba(var(--v-theme-warning), 0.2); }
  100% { border-color: inherit; box-shadow: none; }
}
</style>
