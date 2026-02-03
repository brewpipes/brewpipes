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
        <v-table class="data-table" density="compact">
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Description</th>
              <th>Updated</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="location in locations" :key="location.id">
              <td>{{ location.name }}</td>
              <td>{{ location.location_type || 'n/a' }}</td>
              <td>{{ location.description || 'n/a' }}</td>
              <td>{{ formatDateTime(location.updated_at) }}</td>
            </tr>
            <tr v-if="locations.length === 0">
              <td colspan="4">No stock locations yet.</td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
    </v-card>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

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
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'

  type StockLocation = {
    id: number
    uuid: string
    name: string
    location_type: string
    description: string
    created_at: string
    updated_at: string
  }

  const { request, normalizeText, formatDateTime } = useInventoryApi()

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

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
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

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

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
      locations.value = await request<StockLocation[]>('/stock-locations')
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
      await request<StockLocation>('/stock-locations', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
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

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.card-title-responsive {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.card-title-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
}

.data-table {
  overflow-x: auto;
}

.data-table :deep(.v-table__wrapper) {
  overflow-x: auto;
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  white-space: nowrap;
}

.data-table :deep(td) {
  font-size: 0.85rem;
}
</style>
