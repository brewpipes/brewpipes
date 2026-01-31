<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Inventory usage
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="loadUsage">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Usage list</v-card-title>
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
                      <th>Production ref</th>
                      <th>Used at</th>
                      <th>Notes</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="usage in usages" :key="usage.id">
                      <td>{{ usage.production_ref_uuid || 'n/a' }}</td>
                      <td>{{ formatDateTime(usage.used_at) }}</td>
                      <td>{{ usage.notes || '' }}</td>
                    </tr>
                    <tr v-if="usages.length === 0">
                      <td colspan="3">No usage records yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create usage</v-card-title>
              <v-card-text>
                <v-text-field v-model="usageForm.production_ref_uuid" label="Production ref UUID" />
                <v-text-field v-model="usageForm.used_at" label="Used at" type="datetime-local" />
                <v-textarea
                  v-model="usageForm.notes"
                  auto-grow
                  label="Notes"
                  rows="2"
                />
                <v-btn block color="primary" @click="createUsage">
                  Add usage
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

type InventoryUsage = {
  id: number
  uuid: string
  production_ref_uuid: string
  used_at: string
  notes: string
  created_at: string
  updated_at: string
}

const { request, normalizeText, normalizeDateTime, formatDateTime } = useInventoryApi()

const usages = ref<InventoryUsage[]>([])
const errorMessage = ref('')
const loading = ref(false)

const usageForm = reactive({
  production_ref_uuid: '',
  used_at: '',
  notes: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

onMounted(async () => {
  await loadUsage()
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function loadUsage() {
  loading.value = true
  errorMessage.value = ''
  try {
    usages.value = await request<InventoryUsage[]>('/inventory-usage')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load usage'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function createUsage() {
  try {
    const payload = {
      production_ref_uuid: normalizeText(usageForm.production_ref_uuid),
      used_at: normalizeDateTime(usageForm.used_at),
      notes: normalizeText(usageForm.notes),
    }
    await request<InventoryUsage>('/inventory-usage', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    usageForm.production_ref_uuid = ''
    usageForm.used_at = ''
    usageForm.notes = ''
    await loadUsage()
    showNotice('Usage created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create usage'
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
