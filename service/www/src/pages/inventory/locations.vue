<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Stock locations
        <v-spacer />
        <v-btn :loading="loading" size="small" variant="text" @click="loadLocations">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Location list</v-card-title>
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
                      <th>Updated</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="location in locations" :key="location.id">
                      <td>{{ location.name }}</td>
                      <td>{{ location.location_type || 'n/a' }}</td>
                      <td>{{ formatDateTime(location.updated_at) }}</td>
                    </tr>
                    <tr v-if="locations.length === 0">
                      <td colspan="3">No stock locations yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create stock location</v-card-title>
              <v-card-text>
                <v-text-field v-model="locationForm.name" label="Name" />
                <v-text-field v-model="locationForm.location_type" label="Location type" />
                <v-textarea
                  v-model="locationForm.description"
                  auto-grow
                  label="Description"
                  rows="2"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="!locationForm.name.trim()"
                  @click="createLocation"
                >
                  Add location
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref } from 'vue'
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

  onMounted(async () => {
    await loadLocations()
  })

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
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
      locationForm.name = ''
      locationForm.location_type = ''
      locationForm.description = ''
      await loadLocations()
      showNotice('Stock location created')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create location'
      errorMessage.value = message
      showNotice(message, 'error')
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
</style>
