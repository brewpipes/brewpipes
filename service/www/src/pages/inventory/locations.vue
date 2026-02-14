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
            {{ item.location_type || 'n/a' }}
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
  <v-dialog v-model="createDialog" :max-width="$vuetify.display.xs ? '100%' : 480" persistent>
    <v-card>
      <v-card-title class="text-h6">Create stock location</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="locationForm.name"
          density="comfortable"
          label="Name"
          placeholder="Main warehouse"
          :rules="[rules.required]"
        />
        <v-text-field
          v-model="locationForm.location_type"
          density="comfortable"
          label="Location type"
          placeholder="Warehouse, Cold storage, etc."
        />
        <v-textarea
          v-model="locationForm.description"
          auto-grow
          density="comfortable"
          label="Description"
          placeholder="Additional details about this location..."
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="closeCreateDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="createLocation"
        >
          Create
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { StockLocation } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { normalizeText } from '@/utils/normalize'

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
  const errorMessage = ref('')
  const loading = ref(false)
  const saving = ref(false)
  const createDialog = ref(false)

  const locationForm = reactive({
    name: '',
    location_type: '',
    description: '',
  })

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
  }

  const isFormValid = computed(() => {
    return locationForm.name.trim().length > 0
  })

  onMounted(async () => {
    await loadLocations()
  })

  function openCreateDialog () {
    locationForm.name = ''
    locationForm.location_type = ''
    locationForm.description = ''
    createDialog.value = true
  }

  function closeCreateDialog () {
    createDialog.value = false
  }

  async function loadLocations () {
    loading.value = true
    errorMessage.value = ''
    try {
      locations.value = await getStockLocations()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load locations'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }

  async function createLocation () {
    if (!isFormValid.value) {
      return
    }

    saving.value = true
    errorMessage.value = ''

    try {
      const payload = {
        name: locationForm.name.trim(),
        location_type: normalizeText(locationForm.location_type),
        description: normalizeText(locationForm.description),
      }
      await createStockLocation(payload)
      closeCreateDialog()
      await loadLocations()
      showNotice('Stock location created')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create location'
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
.inventory-page {
  position: relative;
}
</style>
