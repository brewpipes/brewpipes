<template>
  <v-container class="procurement-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Purchase order lines
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
                Line list
                <v-spacer />
                <v-btn size="small" variant="text" @click="loadLines">Apply filter</v-btn>
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
                  v-model="filters.purchase_order_id"
                  clearable
                  :items="orderSelectItems"
                  label="Filter by purchase order"
                />
                <v-table class="data-table" density="compact">
                  <thead>
                    <tr>
                      <th>Order</th>
                      <th>Line</th>
                      <th>Item</th>
                      <th>Qty</th>
                      <th>Unit cost</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="line in lines" :key="line.id">
                      <td>{{ orderNumber(line.purchase_order_id) }}</td>
                      <td>{{ line.line_number }}</td>
                      <td>{{ line.item_name }}</td>
                      <td>{{ `${line.quantity} ${line.quantity_unit}` }}</td>
                      <td>{{ formatCurrency(line.unit_cost_cents, line.currency) }}</td>
                    </tr>
                    <tr v-if="lines.length === 0">
                      <td colspan="5">No order lines yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create order line</v-card-title>
              <v-card-text>
                <v-select
                  v-model="lineForm.purchase_order_id"
                  :items="orderSelectItems"
                  label="Purchase order"
                />
                <v-text-field v-model="lineForm.line_number" label="Line number" type="number" />
                <v-select
                  v-model="lineForm.item_type"
                  :items="itemTypeOptions"
                  label="Item type"
                />
                <v-text-field v-model="lineForm.item_name" label="Item name" />
                <v-text-field v-model="lineForm.inventory_item_uuid" label="Inventory item UUID" />
                <v-text-field v-model="lineForm.quantity" label="Quantity" type="number" />
                <v-combobox
                  v-model="lineForm.quantity_unit"
                  :items="unitOptions"
                  label="Quantity unit"
                />
                <v-text-field v-model="lineForm.unit_cost_cents" label="Unit cost (cents)" type="number" />
                <v-combobox
                  v-model="lineForm.currency"
                  :items="currencyOptions"
                  label="Currency"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="
                    !lineForm.purchase_order_id ||
                      !lineForm.line_number ||
                      !lineForm.item_type ||
                      !lineForm.item_name.trim() ||
                      !lineForm.quantity ||
                      !lineForm.quantity_unit ||
                      !lineForm.unit_cost_cents ||
                      !lineForm.currency
                  "
                  @click="createLine"
                >
                  Add line
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
    id: number
    order_number: string
  }

  type PurchaseOrderLine = {
    id: number
    uuid: string
    purchase_order_id: number
    line_number: number
    item_type: string
    item_name: string
    inventory_item_uuid: string | null
    quantity: number
    quantity_unit: string
    unit_cost_cents: number
    currency: string
    created_at: string
    updated_at: string
  }

  const { request, normalizeText, toNumber, formatCurrency } = useProcurementApi()
  const route = useRoute()

  const orders = ref<PurchaseOrder[]>([])
  const lines = ref<PurchaseOrderLine[]>([])
  const loading = ref(false)
  const errorMessage = ref('')

  const itemTypeOptions = ['ingredient', 'packaging', 'service', 'equipment', 'other']
  const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']
  const currencyOptions = ['USD', 'CAD', 'EUR', 'GBP']

  const filters = reactive({
    purchase_order_id: null as number | null,
  })

  const lineForm = reactive({
    purchase_order_id: null as number | null,
    line_number: '',
    item_type: '',
    item_name: '',
    inventory_item_uuid: '',
    quantity: '',
    quantity_unit: '',
    unit_cost_cents: '',
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
      value: order.id,
    })),
  )

  onMounted(async () => {
    await refreshAll()
    const queryId = route.query.purchase_order_id
    if (typeof queryId === 'string') {
      filters.purchase_order_id = Number(queryId)
      await loadLines()
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
      await Promise.all([loadOrders(), loadLines()])
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load order lines'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }

  async function loadOrders () {
    orders.value = await request<PurchaseOrder[]>('/purchase-orders')
  }

  async function loadLines () {
    const query = new URLSearchParams()
    if (filters.purchase_order_id) {
      query.set('purchase_order_id', String(filters.purchase_order_id))
    }
    const path = query.toString() ? `/purchase-order-lines?${query.toString()}` : '/purchase-order-lines'
    lines.value = await request<PurchaseOrderLine[]>(path)
  }

  async function createLine () {
    try {
      const payload = {
        purchase_order_id: lineForm.purchase_order_id,
        line_number: toNumber(lineForm.line_number),
        item_type: lineForm.item_type,
        item_name: lineForm.item_name.trim(),
        inventory_item_uuid: normalizeText(lineForm.inventory_item_uuid),
        quantity: toNumber(lineForm.quantity),
        quantity_unit: lineForm.quantity_unit,
        unit_cost_cents: toNumber(lineForm.unit_cost_cents),
        currency: lineForm.currency,
      }
      await request<PurchaseOrderLine>('/purchase-order-lines', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      lineForm.purchase_order_id = null
      lineForm.line_number = ''
      lineForm.item_type = ''
      lineForm.item_name = ''
      lineForm.inventory_item_uuid = ''
      lineForm.quantity = ''
      lineForm.quantity_unit = ''
      lineForm.unit_cost_cents = ''
      lineForm.currency = 'USD'
      await loadLines()
      showNotice('Order line created')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create order line'
      errorMessage.value = message
      showNotice(message, 'error')
    }
  }

  function orderNumber (orderId: number) {
    return orders.value.find(order => order.id === orderId)?.order_number ?? `PO ${orderId}`
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
