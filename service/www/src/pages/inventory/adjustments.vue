<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Inventory adjustments
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="loadAdjustments">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Adjustment list</v-card-title>
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

type InventoryAdjustment = {
  id: number
  uuid: string
  reason: string
  adjusted_at: string
  notes: string
  created_at: string
  updated_at: string
}

const { request, normalizeText, normalizeDateTime, formatDateTime } = useInventoryApi()

const adjustments = ref<InventoryAdjustment[]>([])
const errorMessage = ref('')
const loading = ref(false)

const adjustmentForm = reactive({
  reason: '',
  adjusted_at: '',
  notes: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

onMounted(async () => {
  await loadAdjustments()
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function loadAdjustments() {
  loading.value = true
  errorMessage.value = ''
  try {
    adjustments.value = await request<InventoryAdjustment[]>('/inventory-adjustments')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load adjustments'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function createAdjustment() {
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
