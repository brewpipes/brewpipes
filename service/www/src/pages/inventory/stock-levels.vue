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
  import StockLevelTable from '@/components/inventory/StockLevelTable.vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'

  const { getStockLevels } = useInventoryApi()

  const activeTab = ref('malt')
  const stockLevels = ref<StockLevel[]>([])
  const loading = ref(false)
  const errorMessage = ref('')

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
    loading.value = true
    errorMessage.value = ''
    try {
      stockLevels.value = await getStockLevels()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load stock levels'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }
</script>

<style scoped>
.stock-levels-page {
  position: relative;
}
</style>
