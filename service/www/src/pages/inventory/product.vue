<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Product
        <v-spacer />
        <v-btn :loading="loading" size="small" variant="text" @click="loadBeerLots">
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
                        <tr v-for="beerLot in beerLots" :key="beerLot.uuid">
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

</template>

<script lang="ts" setup>
  import type { BeerLot } from '@/types'
  import { onMounted, reactive, ref } from 'vue'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { normalizeDateTime, normalizeText } from '@/utils/normalize'

  const { getBeerLots: fetchBeerLots, createBeerLot: createBeerLotApi } = useInventoryApi()
  const { showNotice } = useSnackbar()

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

  onMounted(async () => {
    await loadBeerLots()
  })

  async function loadBeerLots () {
    loading.value = true
    errorMessage.value = ''
    try {
      beerLots.value = await fetchBeerLots()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load product lots'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }

  async function createBeerLot () {
    try {
      const payload = {
        production_batch_uuid: beerLotForm.production_batch_uuid.trim(),
        lot_code: normalizeText(beerLotForm.lot_code),
        packaged_at: normalizeDateTime(beerLotForm.packaged_at),
        notes: normalizeText(beerLotForm.notes),
      }
      await createBeerLotApi(payload)
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
</style>
