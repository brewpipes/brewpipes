<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Adjustments & Transfers
        <v-spacer />
        <v-btn :loading="loading" size="small" variant="text" @click="refreshAll">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <div class="text-subtitle-1 font-weight-semibold mb-2">Adjustments</div>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Adjustment list</v-card-title>
              <v-card-text>
                <v-alert
                  v-if="adjustmentErrorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ adjustmentErrorMessage }}
                </v-alert>
                <v-table class="data-table" density="compact">
                  <thead>
                    <tr>
                      <th>Reason</th>
                      <th>Adjusted at</th>
                      <th>Notes</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="adjustment in adjustments" :key="adjustment.id">
                      <td>{{ adjustment.reason }}</td>
                      <td>{{ formatDateTime(adjustment.adjusted_at) }}</td>
                      <td>{{ adjustment.notes || '' }}</td>
                    </tr>
                    <tr v-if="adjustments.length === 0">
                      <td colspan="3">No adjustments yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create adjustment</v-card-title>
              <v-card-text>
                <v-text-field v-model="adjustmentForm.reason" label="Reason" />
                <v-text-field v-model="adjustmentForm.adjusted_at" label="Adjusted at" type="datetime-local" />
                <v-textarea
                  v-model="adjustmentForm.notes"
                  auto-grow
                  label="Notes"
                  rows="2"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="!adjustmentForm.reason.trim()"
                  @click="createAdjustment"
                >
                  Add adjustment
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <v-divider class="my-6" />

        <div class="text-subtitle-1 font-weight-semibold mb-2">Transfers</div>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Transfer list</v-card-title>
              <v-card-text>
                <v-alert
                  v-if="transferErrorMessage"
                  class="mb-3"
                  density="compact"
                  type="error"
                  variant="tonal"
                >
                  {{ transferErrorMessage }}
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

  type InventoryAdjustment = {
    id: number
    uuid: string
    reason: string
    adjusted_at: string
    notes: string
    created_at: string
    updated_at: string
  }

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

  const adjustments = ref<InventoryAdjustment[]>([])
  const locations = ref<StockLocation[]>([])
  const transfers = ref<InventoryTransfer[]>([])
  const loading = ref(false)

  const adjustmentErrorMessage = ref('')
  const transferErrorMessage = ref('')

  const adjustmentForm = reactive({
    reason: '',
    adjusted_at: '',
    notes: '',
  })

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
    locations.value.map(location => ({
      title: location.name,
      value: location.id,
    })),
  )

  onMounted(async () => {
    await refreshAll()
  })

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function refreshAll () {
    loading.value = true
    await Promise.allSettled([loadAdjustments(), loadTransfers(), loadLocations()])
    loading.value = false
  }

  async function loadAdjustments () {
    adjustmentErrorMessage.value = ''
    try {
      adjustments.value = await request<InventoryAdjustment[]>('/inventory-adjustments')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load adjustments'
      adjustmentErrorMessage.value = message
    }
  }

  async function loadLocations () {
    try {
      locations.value = await request<StockLocation[]>('/stock-locations')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load locations'
      transferErrorMessage.value = message
    }
  }

  async function loadTransfers () {
    transferErrorMessage.value = ''
    try {
      transfers.value = await request<InventoryTransfer[]>('/inventory-transfers')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load transfers'
      transferErrorMessage.value = message
    }
  }

  async function createAdjustment () {
    try {
      const payload = {
        reason: adjustmentForm.reason.trim(),
        adjusted_at: normalizeDateTime(adjustmentForm.adjusted_at),
        notes: normalizeText(adjustmentForm.notes),
      }
      await request<InventoryAdjustment>('/inventory-adjustments', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      adjustmentForm.reason = ''
      adjustmentForm.adjusted_at = ''
      adjustmentForm.notes = ''
      await loadAdjustments()
      showNotice('Adjustment created')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create adjustment'
      adjustmentErrorMessage.value = message
      showNotice(message, 'error')
    }
  }

  async function createTransfer () {
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
      transferErrorMessage.value = message
      showNotice(message, 'error')
    }
  }

  function locationName (locationId: number) {
    return locations.value.find(location => location.id === locationId)?.name ?? `Location ${locationId}`
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
