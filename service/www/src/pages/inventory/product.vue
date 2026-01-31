<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Product
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="loadBeerLots">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-tabs v-model="activeTab" class="inventory-tabs" color="primary" show-arrows>
          <v-tab value="stock">Stock</v-tab>
        </v-tabs>

        <v-window v-model="activeTab" class="mt-4">
          <v-window-item value="stock">
            <v-row align="stretch">
              <v-col cols="12" md="7">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title>Product lot list</v-card-title>
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
                          <th>Batch UUID</th>
                          <th>Lot code</th>
                          <th>Packaged at</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="beerLot in beerLots" :key="beerLot.id">
                          <td>{{ beerLot.production_batch_uuid }}</td>
                          <td>{{ beerLot.lot_code || 'n/a' }}</td>
                          <td>{{ formatDateTime(beerLot.packaged_at) }}</td>
                        </tr>
                        <tr v-if="beerLots.length === 0">
                          <td colspan="3">No product lots yet.</td>
                        </tr>
                      </tbody>
                    </v-table>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="5">
                <v-card class="sub-card" variant="tonal">
                  <v-card-title>Create product lot</v-card-title>
                  <v-card-text>
                    <v-text-field v-model="beerLotForm.production_batch_uuid" label="Production batch UUID" />
                    <v-text-field v-model="beerLotForm.lot_code" label="Lot code" />
                    <v-text-field v-model="beerLotForm.packaged_at" label="Packaged at" type="datetime-local" />
                    <v-textarea
                      v-model="beerLotForm.notes"
                      auto-grow
                      label="Notes"
                      rows="2"
                    />
                    <v-btn
                      block
                      color="primary"
                      :disabled="!beerLotForm.production_batch_uuid.trim()"
                      @click="createBeerLot"
                    >
                      Add product lot
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-window-item>
        </v-window>
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

type BeerLot = {
  id: number
  uuid: string
  production_batch_uuid: string
  lot_code: string
  packaged_at: string
  notes: string
  created_at: string
  updated_at: string
}

const { request, normalizeText, normalizeDateTime, formatDateTime } = useInventoryApi()

const beerLots = ref<BeerLot[]>([])
const errorMessage = ref('')
const loading = ref(false)
const activeTab = ref('stock')

const beerLotForm = reactive({
  production_batch_uuid: '',
  lot_code: '',
  packaged_at: '',
  notes: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

onMounted(async () => {
  await loadBeerLots()
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function loadBeerLots() {
  loading.value = true
  errorMessage.value = ''
  try {
    beerLots.value = await request<BeerLot[]>('/beer-lots')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load product lots'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function createBeerLot() {
  try {
    const payload = {
      production_batch_uuid: beerLotForm.production_batch_uuid.trim(),
      lot_code: normalizeText(beerLotForm.lot_code),
      packaged_at: normalizeDateTime(beerLotForm.packaged_at),
      notes: normalizeText(beerLotForm.notes),
    }
    await request<BeerLot>('/beer-lots', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    beerLotForm.production_batch_uuid = ''
    beerLotForm.lot_code = ''
    beerLotForm.packaged_at = ''
    beerLotForm.notes = ''
    await loadBeerLots()
    showNotice('Product lot created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create product lot'
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
