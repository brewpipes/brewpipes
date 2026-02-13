<template>
  <v-container class="procurement-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Purchase order fees
        <v-spacer />
        <v-btn :loading="loading" size="small" variant="text" @click="refreshAll">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title class="d-flex align-center">
                Fee list
                <v-spacer />
                <v-btn size="small" variant="text" @click="loadFees">Apply filter</v-btn>
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
                <v-select
                  v-model="filters.purchase_order_uuid"
                  clearable
                  :items="orderSelectItems"
                  label="Filter by purchase order"
                />
                <v-table class="data-table" density="compact">
                  <thead>
                    <tr>
                      <th>Order</th>
                      <th>Fee type</th>
                      <th>Amount</th>
                      <th>Updated</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="fee in fees" :key="fee.uuid">
                      <td>{{ orderNumber(fee.purchase_order_uuid) }}</td>
                      <td>{{ fee.fee_type }}</td>
                      <td>{{ formatCurrency(fee.amount_cents, fee.currency) }}</td>
                      <td>{{ formatDateTime(fee.updated_at) }}</td>
                    </tr>
                    <tr v-if="fees.length === 0">
                      <td colspan="4">No fees yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create fee</v-card-title>
              <v-card-text>
                <v-select
                  v-model="feeForm.purchase_order_uuid"
                  :items="orderSelectItems"
                  label="Purchase order"
                />
                <v-text-field v-model="feeForm.fee_type" label="Fee type" />
                <v-text-field v-model="feeForm.amount_cents" label="Amount (cents)" type="number" />
                <v-combobox
                  v-model="feeForm.currency"
                  :items="currencyOptions"
                  label="Currency"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="!feeForm.purchase_order_uuid || !feeForm.fee_type.trim() || !feeForm.amount_cents || !feeForm.currency"
                  @click="createFee"
                >
                  Add fee
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
  import { useRoute } from 'vue-router'
  import { useProcurementApi } from '@/composables/useProcurementApi'

  type PurchaseOrder = {
    uuid: string
    order_number: string
  }

  type PurchaseOrderFee = {
    uuid: string
    purchase_order_uuid: string
    fee_type: string
    amount_cents: number
    currency: string
    created_at: string
    updated_at: string
  }

  const { request, toNumber, formatDateTime, formatCurrency } = useProcurementApi()
  const route = useRoute()

  const orders = ref<PurchaseOrder[]>([])
  const fees = ref<PurchaseOrderFee[]>([])
  const loading = ref(false)
  const errorMessage = ref('')

  const currencyOptions = ['USD', 'CAD', 'EUR', 'GBP']

  const filters = reactive({
    purchase_order_uuid: null as string | null,
  })

  const feeForm = reactive({
    purchase_order_uuid: null as string | null,
    fee_type: '',
    amount_cents: '',
    currency: 'USD',
  })

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  const orderSelectItems = computed(() =>
    orders.value.map(order => ({
      title: order.order_number,
      value: order.uuid,
    })),
  )

  onMounted(async () => {
    await refreshAll()
    const queryUuid = route.query.purchase_order_uuid
    if (typeof queryUuid === 'string') {
      filters.purchase_order_uuid = queryUuid
      await loadFees()
    }
  })

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function refreshAll () {
    loading.value = true
    errorMessage.value = ''
    try {
      await Promise.all([loadOrders(), loadFees()])
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load fees'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }

  async function loadOrders () {
    orders.value = await request<PurchaseOrder[]>('/purchase-orders')
  }

  async function loadFees () {
    const query = new URLSearchParams()
    if (filters.purchase_order_uuid) {
      query.set('purchase_order_uuid', filters.purchase_order_uuid)
    }
    const path = query.toString() ? `/purchase-order-fees?${query.toString()}` : '/purchase-order-fees'
    fees.value = await request<PurchaseOrderFee[]>(path)
  }

  async function createFee () {
    try {
      const payload = {
        purchase_order_uuid: feeForm.purchase_order_uuid,
        fee_type: feeForm.fee_type.trim(),
        amount_cents: toNumber(feeForm.amount_cents),
        currency: feeForm.currency,
      }
      await request<PurchaseOrderFee>('/purchase-order-fees', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      feeForm.purchase_order_uuid = null
      feeForm.fee_type = ''
      feeForm.amount_cents = ''
      feeForm.currency = 'USD'
      await loadFees()
      showNotice('Fee created')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create fee'
      errorMessage.value = message
      showNotice(message, 'error')
    }
  }

  function orderNumber (orderUuid: string) {
    const order = orders.value.find(o => o.uuid === orderUuid)
    return order?.order_number ?? 'Unknown PO'
  }
</script>

<style scoped>
.procurement-page {
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
