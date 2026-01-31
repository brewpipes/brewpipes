<template>
  <v-container class="procurement-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Purchase orders
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="refreshAll">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title class="d-flex align-center">
                Order list
                <v-spacer />
                <v-btn size="small" variant="text" @click="loadOrders">Apply filter</v-btn>
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
                  v-model="filters.supplier_id"
                  :items="supplierSelectItems"
                  label="Filter by supplier"
                  clearable
                />
                <v-table class="data-table" density="compact">
                  <thead>
                    <tr>
                      <th>Order</th>
                      <th>Supplier</th>
                      <th>Status</th>
                      <th>Expected</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="order in orders" :key="order.id">
                      <td>{{ order.order_number }}</td>
                      <td>{{ supplierName(order.supplier_id) }}</td>
                      <td>{{ order.status }}</td>
                      <td>{{ formatDateTime(order.expected_at) }}</td>
                      <td class="text-right">
                        <v-btn
                          size="x-small"
                          variant="text"
                          @click="openLines(order.id)"
                        >
                          Lines
                        </v-btn>
                        <v-btn
                          size="x-small"
                          variant="text"
                          @click="openFees(order.id)"
                        >
                          Fees
                        </v-btn>
                      </td>
                    </tr>
                    <tr v-if="orders.length === 0">
                      <td colspan="5">No purchase orders yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create purchase order</v-card-title>
              <v-card-text>
                <v-select
                  v-model="orderForm.supplier_id"
                  :items="supplierSelectItems"
                  label="Supplier"
                />
                <v-text-field v-model="orderForm.order_number" label="Order number" />
                <v-select
                  v-model="orderForm.status"
                  :items="statusOptions"
                  label="Status"
                  clearable
                />
                <v-text-field v-model="orderForm.ordered_at" label="Ordered at" type="datetime-local" />
                <v-text-field v-model="orderForm.expected_at" label="Expected at" type="datetime-local" />
                <v-textarea
                  v-model="orderForm.notes"
                  auto-grow
                  label="Notes"
                  rows="2"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="!orderForm.supplier_id || !orderForm.order_number.trim()"
                  @click="createOrder"
                >
                  Add purchase order
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
import { useRouter } from 'vue-router'
import { useProcurementApi } from '@/composables/useProcurementApi'

type Supplier = {
  id: number
  name: string
}

type PurchaseOrder = {
  id: number
  uuid: string
  supplier_id: number
  order_number: string
  status: string
  ordered_at: string | null
  expected_at: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

const { request, normalizeText, normalizeDateTime, formatDateTime } = useProcurementApi()
const router = useRouter()

const suppliers = ref<Supplier[]>([])
const orders = ref<PurchaseOrder[]>([])
const loading = ref(false)
const errorMessage = ref('')

const statusOptions = [
  'draft',
  'submitted',
  'confirmed',
  'partially_received',
  'received',
  'cancelled',
]

const filters = reactive({
  supplier_id: null as number | null,
})

const orderForm = reactive({
  supplier_id: null as number | null,
  order_number: '',
  status: '',
  ordered_at: '',
  expected_at: '',
  notes: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

const supplierSelectItems = computed(() =>
  suppliers.value.map((supplier) => ({
    title: supplier.name,
    value: supplier.id,
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
    await Promise.all([loadSuppliers(), loadOrders()])
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load purchase orders'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function loadSuppliers() {
  suppliers.value = await request<Supplier[]>('/suppliers')
}

async function loadOrders() {
  const query = new URLSearchParams()
  if (filters.supplier_id) {
    query.set('supplier_id', String(filters.supplier_id))
  }
  const path = query.toString() ? `/purchase-orders?${query.toString()}` : '/purchase-orders'
  orders.value = await request<PurchaseOrder[]>(path)
}

async function createOrder() {
  try {
    const payload = {
      supplier_id: orderForm.supplier_id,
      order_number: orderForm.order_number.trim(),
      status: normalizeText(orderForm.status),
      ordered_at: normalizeDateTime(orderForm.ordered_at),
      expected_at: normalizeDateTime(orderForm.expected_at),
      notes: normalizeText(orderForm.notes),
    }
    await request<PurchaseOrder>('/purchase-orders', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    orderForm.supplier_id = null
    orderForm.order_number = ''
    orderForm.status = ''
    orderForm.ordered_at = ''
    orderForm.expected_at = ''
    orderForm.notes = ''
    await loadOrders()
    showNotice('Purchase order created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create purchase order'
    errorMessage.value = message
    showNotice(message, 'error')
  }
}

function supplierName(supplierId: number) {
  return suppliers.value.find((supplier) => supplier.id === supplierId)?.name ?? `Supplier ${supplierId}`
}

function openLines(orderId: number) {
  router.push({
    path: '/procurement/purchase-order-lines',
    query: { purchase_order_id: String(orderId) },
  })
}

function openFees(orderId: number) {
  router.push({
    path: '/procurement/purchase-order-fees',
    query: { purchase_order_id: String(orderId) },
  })
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
