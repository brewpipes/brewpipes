<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Inventory transfers
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="refreshAll">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Transfer list</v-card-title>
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
                      <th>Source</th>
                      <th>Destination</th>
                      <th>Transferred at</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="transfer in transfers" :key="transfer.id">
                      <td>{{ locationName(transfer.source_location_id) }}</td>
                      <td>{{ locationName(transfer.dest_location_id) }}</td>
                      <td>{{ formatDateTime(transfer.transferred_at) }}</td>
                    </tr>
                    <tr v-if="transfers.length === 0">
                      <td colspan="3">No transfers yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create transfer</v-card-title>
              <v-card-text>
                <v-select
                  v-model="transferForm.source_location_id"
                  :items="locationSelectItems"
                  label="Source location"
                />
                <v-select
                  v-model="transferForm.dest_location_id"
                  :items="locationSelectItems"
                  label="Destination location"
                />
                <v-text-field v-model="transferForm.transferred_at" label="Transferred at" type="datetime-local" />
                <v-textarea
                  v-model="transferForm.notes"
                  auto-grow
                  label="Notes"
                  rows="2"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="!transferForm.source_location_id || !transferForm.dest_location_id"
                  @click="createTransfer"
                >
                  Add transfer
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
import { computed, onMounted, reactive, ref } from 'vue'
import { useInventoryApi } from '@/composables/useInventoryApi'

type StockLocation = {
  id: number
  name: string
}

type InventoryTransfer = {
  id: number
  uuid: string
  source_location_id: number
  dest_location_id: number
  transferred_at: string
  notes: string
  created_at: string
  updated_at: string
}

const { request, normalizeText, normalizeDateTime, formatDateTime } = useInventoryApi()

const locations = ref<StockLocation[]>([])
const transfers = ref<InventoryTransfer[]>([])
const loading = ref(false)
const errorMessage = ref('')

const transferForm = reactive({
  source_location_id: null as number | null,
  dest_location_id: null as number | null,
  transferred_at: '',
  notes: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

const locationSelectItems = computed(() =>
  locations.value.map((location) => ({
    title: location.name,
    value: location.id,
  })),
)

onMounted(async () => {
  await refreshAll()
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function refreshAll() {
  loading.value = true
  errorMessage.value = ''
  try {
    await Promise.all([loadLocations(), loadTransfers()])
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load transfers'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function loadLocations() {
  locations.value = await request<StockLocation[]>('/stock-locations')
}

async function loadTransfers() {
  transfers.value = await request<InventoryTransfer[]>('/inventory-transfers')
}

async function createTransfer() {
  try {
    const payload = {
      source_location_id: transferForm.source_location_id,
      dest_location_id: transferForm.dest_location_id,
      transferred_at: normalizeDateTime(transferForm.transferred_at),
      notes: normalizeText(transferForm.notes),
    }
    await request<InventoryTransfer>('/inventory-transfers', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    transferForm.source_location_id = null
    transferForm.dest_location_id = null
    transferForm.transferred_at = ''
    transferForm.notes = ''
    await loadTransfers()
    showNotice('Transfer created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create transfer'
    errorMessage.value = message
    showNotice(message, 'error')
  }
}

function locationName(locationId: number) {
  return locations.value.find((location) => location.id === locationId)?.name ?? `Location ${locationId}`
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
