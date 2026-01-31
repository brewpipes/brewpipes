<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Inventory receipts
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="loadReceipts">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Receipt list</v-card-title>
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
                      <th>Reference</th>
                      <th>Supplier</th>
                      <th>Received at</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="receipt in receipts" :key="receipt.id">
                      <td>{{ receipt.reference_code || 'n/a' }}</td>
                      <td>{{ receipt.supplier_uuid || 'n/a' }}</td>
                      <td>{{ formatDateTime(receipt.received_at) }}</td>
                    </tr>
                    <tr v-if="receipts.length === 0">
                      <td colspan="3">No receipts yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create receipt</v-card-title>
              <v-card-text>
                <v-text-field v-model="receiptForm.reference_code" label="Reference code" />
                <v-text-field v-model="receiptForm.supplier_uuid" label="Supplier UUID" />
                <v-text-field v-model="receiptForm.received_at" label="Received at" type="datetime-local" />
                <v-textarea
                  v-model="receiptForm.notes"
                  auto-grow
                  label="Notes"
                  rows="2"
                />
                <v-btn block color="primary" @click="createReceipt">
                  Add receipt
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

type InventoryReceipt = {
  id: number
  uuid: string
  supplier_uuid: string
  reference_code: string
  received_at: string
  notes: string
  created_at: string
  updated_at: string
}

const { request, normalizeText, normalizeDateTime, formatDateTime } = useInventoryApi()

const receipts = ref<InventoryReceipt[]>([])
const errorMessage = ref('')
const loading = ref(false)

const receiptForm = reactive({
  supplier_uuid: '',
  reference_code: '',
  received_at: '',
  notes: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

onMounted(async () => {
  await loadReceipts()
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function loadReceipts() {
  loading.value = true
  errorMessage.value = ''
  try {
    receipts.value = await request<InventoryReceipt[]>('/inventory-receipts')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load receipts'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function createReceipt() {
  try {
    const payload = {
      supplier_uuid: normalizeText(receiptForm.supplier_uuid),
      reference_code: normalizeText(receiptForm.reference_code),
      received_at: normalizeDateTime(receiptForm.received_at),
      notes: normalizeText(receiptForm.notes),
    }
    await request<InventoryReceipt>('/inventory-receipts', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    receiptForm.supplier_uuid = ''
    receiptForm.reference_code = ''
    receiptForm.received_at = ''
    receiptForm.notes = ''
    await loadReceipts()
    showNotice('Receipt created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create receipt'
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
