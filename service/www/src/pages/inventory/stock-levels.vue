<template>
  <v-container class="stock-levels-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <div class="d-flex align-center">
          <v-icon class="mr-2" icon="mdi-chart-box-outline" />
          Stock Levels
        </div>
        <v-btn
          :loading="loading"
          size="small"
          variant="text"
          @click="loadStockLevels"
        >
          <v-icon class="mr-1" icon="mdi-refresh" size="small" />
          <span class="d-none d-sm-inline">Refresh</span>
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

        <v-tabs v-model="activeTab" class="stock-tabs" color="primary" show-arrows>
          <v-tab value="malt">Malt</v-tab>
          <v-tab value="hops">Hops</v-tab>
          <v-tab value="yeast">Yeast</v-tab>
          <v-tab value="other">Other</v-tab>
        </v-tabs>

        <v-window v-model="activeTab" class="mt-4">
          <!-- Malt Tab -->
          <v-window-item value="malt">
            <StockLevelTable
              category-label="Malt"
              :items="maltItems"
              :loading="loading"
            />
          </v-window-item>

          <!-- Hops Tab -->
          <v-window-item value="hops">
            <StockLevelTable
              category-label="Hops"
              :items="hopItems"
              :loading="loading"
            />
          </v-window-item>

          <!-- Yeast Tab -->
          <v-window-item value="yeast">
            <StockLevelTable
              category-label="Yeast"
              :items="yeastItems"
              :loading="loading"
            />
          </v-window-item>

          <!-- Other Tab -->
          <v-window-item value="other">
            <StockLevelTable
              category-label="Other"
              :items="otherItems"
              :loading="loading"
              show-category
            />
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script lang="ts" setup>
  import type { StockLevel } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { useRoute } from 'vue-router'
  import StockLevelTable from '@/components/inventory/StockLevelTable.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useInventoryApi } from '@/composables/useInventoryApi'

  const route = useRoute()
  const { getStockLevels } = useInventoryApi()

  const categoryToTab: Record<string, string> = {
    fermentable: 'malt',
    hop: 'hops',
    yeast: 'yeast',
    adjunct: 'other',
    salt: 'other',
    chemical: 'other',
    gas: 'other',
    other: 'other',
  }

  const activeTab = ref('malt')
  const stockLevels = ref<StockLevel[]>([])

  const { execute, loading, error: errorMessage } = useAsyncAction()

  // Category groupings
  const maltItems = computed(() =>
    stockLevels.value.filter(item => item.category === 'fermentable'),
  )

  const hopItems = computed(() =>
    stockLevels.value.filter(item => item.category === 'hop'),
  )

  const yeastItems = computed(() =>
    stockLevels.value.filter(item => item.category === 'yeast'),
  )

  const otherCategories = new Set(['adjunct', 'salt', 'chemical', 'gas', 'other'])
  const otherItems = computed(() =>
    stockLevels.value.filter(item => otherCategories.has(item.category)),
  )

  onMounted(async () => {
    await loadStockLevels()
  })

  async function loadStockLevels () {
    await execute(async () => {
      stockLevels.value = await getStockLevels()
      autoSelectTab()
    })
  }

  function autoSelectTab () {
    const ingredientUuid = route.query.ingredient as string | undefined
    if (!ingredientUuid) return
    const match = stockLevels.value.find(item => item.ingredient_uuid === ingredientUuid)
    if (match) {
      activeTab.value = categoryToTab[match.category] ?? 'malt'
    }
  }
</script>

<style scoped>
.stock-levels-page {
  position: relative;
}
</style>
