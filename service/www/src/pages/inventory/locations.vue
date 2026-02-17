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
          <template #item.actions="{ item }">
            <v-btn
              aria-label="Edit location"
              icon="mdi-pencil"
              size="x-small"
              variant="text"
              @click.stop="openEditDialog(item)"
            />
            <v-btn
              aria-label="Delete location"
              icon="mdi-delete"
              size="x-small"
              variant="text"
              @click.stop="confirmDelete(item)"
            />
          </template>
          <template #no-data>
            <div class="text-center py-4 text-medium-emphasis">No stock locations yet.</div>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </v-container>

  <!-- Create/Edit Location Dialog -->
  <LocationCreateDialog
    v-model="locationDialog"
    :edit-location="editingLocation"
    :saving="saving"
    @submit="saveLocation"
  />

  <!-- Delete confirmation dialog -->
  <v-dialog v-model="showDeleteDialog" max-width="400" persistent>
    <v-card>
      <v-card-title class="text-h6">Delete location?</v-card-title>
      <v-card-text>
        Are you sure you want to delete <strong>{{ deletingLocation?.name }}</strong>?
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="deleting" variant="text" @click="showDeleteDialog = false">Cancel</v-btn>
        <v-btn color="error" :loading="deleting" variant="tonal" @click="handleDelete">Delete</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { StockLocation } from '@/types'
  import { onMounted, ref } from 'vue'
  import LocationCreateDialog from '@/components/inventory/LocationCreateDialog.vue'
  import type { LocationFormData } from '@/components/inventory/LocationCreateDialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'

  const { getStockLocations, createStockLocation, updateStockLocation, deleteStockLocation } = useInventoryApi()
  const { showNotice } = useSnackbar()

  // Table configuration
  const locationHeaders = [
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Type', key: 'location_type', sortable: true },
    { title: 'Description', key: 'description', sortable: false },
    { title: 'Updated', key: 'updated_at', sortable: true },
    { title: '', key: 'actions', sortable: false, align: 'end' as const, width: '100px' },
  ]

  // State
  const locations = ref<StockLocation[]>([])

  const { execute: executeLoad, loading, error: errorMessage } = useAsyncAction()
  const { execute: executeSave, loading: saving, error: saveError } = useAsyncAction()

  // Dialog state
  const locationDialog = ref(false)
  const editingLocation = ref<StockLocation | null>(null)
  const showDeleteDialog = ref(false)
  const deletingLocation = ref<StockLocation | null>(null)
  const deleting = ref(false)

  // Lifecycle
  onMounted(async () => {
    await loadLocations()
  })

  // Methods
  function openCreateDialog () {
    editingLocation.value = null
    locationDialog.value = true
  }

  function openEditDialog (location: StockLocation) {
    editingLocation.value = location
    locationDialog.value = true
  }

  async function loadLocations () {
    await executeLoad(async () => {
      locations.value = await getStockLocations()
    })
  }

  async function saveLocation (data: LocationFormData) {
    await executeSave(async () => {
      if (editingLocation.value) {
        await updateStockLocation(editingLocation.value.uuid, data)
        showNotice('Stock location updated')
      } else {
        await createStockLocation(data)
        showNotice('Stock location created')
      }

      locationDialog.value = false
      editingLocation.value = null
      await loadLocations()
    })
    if (saveError.value) {
      showNotice(saveError.value, 'error')
    }
  }

  function confirmDelete (location: StockLocation) {
    deletingLocation.value = location
    showDeleteDialog.value = true
  }

  async function handleDelete () {
    if (!deletingLocation.value) return

    deleting.value = true
    try {
      await deleteStockLocation(deletingLocation.value.uuid)
      showNotice('Stock location deleted')
      showDeleteDialog.value = false
      deletingLocation.value = null
      await loadLocations()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to delete location'
      if (message.includes('has inventory') || message.includes('409') || message.includes('Conflict')) {
        showNotice('Cannot delete this location because it has inventory. Move or adjust inventory first.', 'error', 5000)
      } else {
        showNotice(message, 'error')
      }
    } finally {
      deleting.value = false
    }
  }
</script>

<style scoped>
.inventory-page {
  position: relative;
}
</style>
