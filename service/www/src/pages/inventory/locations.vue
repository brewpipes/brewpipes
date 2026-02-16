<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <span>Stock locations</span>
        <div class="card-title-actions">
          <v-btn
            :icon="$vuetify.display.xs"
            :loading="loading"
            size="small"
            variant="text"
            @click="loadLocations"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-refresh" />
            <span v-else>Refresh</span>
          </v-btn>
          <v-btn
            color="primary"
            :icon="$vuetify.display.xs"
            :prepend-icon="$vuetify.display.xs ? undefined : 'mdi-plus'"
            size="small"
            variant="text"
            @click="openCreateDialog"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-plus" />
            <span v-else>Create location</span>
          </v-btn>
        </div>
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
          class="data-table"
          density="compact"
          :headers="locationHeaders"
          item-value="uuid"
          :items="locations"
          :loading="loading"
        >
          <template #item.location_type="{ item }">
            {{ item.location_type ? item.location_type.charAt(0).toUpperCase() + item.location_type.slice(1) : 'n/a' }}
          </template>
          <template #item.description="{ item }">
            {{ item.description || 'n/a' }}
          </template>
          <template #item.updated_at="{ item }">
            {{ formatDateTime(item.updated_at) }}
          </template>
          <template #no-data>
            <div class="text-center py-4 text-medium-emphasis">No stock locations yet.</div>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </v-container>

  <!-- Create Location Dialog -->
  <LocationCreateDialog
    v-model="createDialog"
    :saving="saving"
    @submit="handleCreateLocation"
  />
</template>

<script lang="ts" setup>
  import type { CreateStockLocationRequest, StockLocation } from '@/types'
  import { onMounted, ref } from 'vue'
  import LocationCreateDialog from '@/components/inventory/LocationCreateDialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'

  const { getStockLocations, createStockLocation } = useInventoryApi()
  const { showNotice } = useSnackbar()

  // Table configuration
  const locationHeaders = [
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Type', key: 'location_type', sortable: true },
    { title: 'Description', key: 'description', sortable: false },
    { title: 'Updated', key: 'updated_at', sortable: true },
  ]

  const locations = ref<StockLocation[]>([])
  const createDialog = ref(false)

  const { execute: executeLoad, loading, error: errorMessage } = useAsyncAction()
  const { execute: executeSave, loading: saving } = useAsyncAction({
    onError: (message) => showNotice(message, 'error'),
  })

  onMounted(async () => {
    await loadLocations()
  })

  function openCreateDialog () {
    createDialog.value = true
  }

  async function loadLocations () {
    await executeLoad(async () => {
      locations.value = await getStockLocations()
    })
  }

  async function handleCreateLocation (data: CreateStockLocationRequest) {
    await executeSave(async () => {
      await createStockLocation(data)
      createDialog.value = false
      await loadLocations()
      showNotice('Stock location created')
    })
  }
</script>

<style scoped>
.inventory-page {
  position: relative;
}
</style>
