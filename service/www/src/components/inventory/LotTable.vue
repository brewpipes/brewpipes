<template>
  <v-card variant="outlined">
    <v-card-title class="d-flex align-center">
      {{ title }}
      <v-spacer />
      <v-btn :loading="loading" size="small" variant="text" @click="emit('click:refresh')">
        Refresh
      </v-btn>
      <v-btn
        class="ml-2"
        color="primary"
        prepend-icon="mdi-plus"
        size="small"
        variant="text"
        @click="emit('click:create')"
      >
        Create lot
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
      <v-data-table
        class="data-table lot-table"
        density="compact"
        :headers="tableHeaders"
        hover
        item-value="uuid"
        :items="lots"
        :loading="loading"
        @click:row="onRowClick"
      >
        <template #item.brewery_lot_code="{ item }">
          {{ item.brewery_lot_code || '—' }}
        </template>
        <template #item.ingredient_uuid="{ item }">
          {{ ingredientName(item.ingredient_uuid) }}
        </template>
        <template v-if="showCategoryColumn" #item.category="{ item }">
          <v-chip density="compact" size="small" variant="tonal">
            {{ ingredientCategory(item.ingredient_uuid) }}
          </v-chip>
        </template>
        <template #item.received_amount="{ item }">
          {{ formatAmountPreferred(item.received_amount, item.received_unit) }}
        </template>
        <template #item.received_at="{ item }">
          {{ formatDateTime(item.received_at) }}
        </template>
        <template #item.best_by_at="{ item }">
          {{ item.best_by_at ? formatDateTime(item.best_by_at) : '—' }}
        </template>
        <template #item.expires_at="{ item }">
          {{ item.expires_at ? formatDateTime(item.expires_at) : '—' }}
        </template>
        <template #no-data>
          <div class="text-center py-4 text-medium-emphasis">{{ emptyText }}</div>
        </template>
      </v-data-table>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { IngredientLot } from '@/types'
  import { computed } from 'vue'

  const props = defineProps<{
    title: string
    lots: IngredientLot[]
    loading: boolean
    errorMessage: string
    emptyText: string
    showCategoryColumn?: boolean
    ingredientName: (uuid: string) => string
    ingredientCategory?: (uuid: string) => string
    formatAmountPreferred: (amount: number, unit: string) => string
    formatDateTime: (value: string | null | undefined) => string
  }>()

  const emit = defineEmits<{
    'click:row': [uuid: string]
    'click:create': []
    'click:refresh': []
  }>()

  const baseLotHeaders = [
    { title: 'Lot Code', key: 'brewery_lot_code', sortable: true },
    { title: 'Ingredient', key: 'ingredient_uuid', sortable: true },
  ]

  const categoryHeader = { title: 'Category', key: 'category', sortable: true }

  const trailingLotHeaders = [
    { title: 'Received Amount', key: 'received_amount', sortable: true },
    { title: 'Received Date', key: 'received_at', sortable: true },
    { title: 'Best By', key: 'best_by_at', sortable: true },
    { title: 'Expires', key: 'expires_at', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
  ]

  const tableHeaders = computed(() => {
    if (props.showCategoryColumn) {
      return [...baseLotHeaders, categoryHeader, ...trailingLotHeaders]
    }
    return [...baseLotHeaders, ...trailingLotHeaders]
  })

  function onRowClick (_event: Event, { item }: { item: IngredientLot }) {
    emit('click:row', item.uuid)
  }
</script>

<style scoped>
.lot-table :deep(tr) {
  cursor: pointer;
}
</style>
