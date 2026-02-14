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
                  v-model="filters.purchase_order_uuid"
                  clearable
                  :items="orderSelectItems"
                  label="Filter by purchase order"
                />
                <v-data-table
                  class="data-table"
                  density="compact"
                  :headers="lineHeaders"
                  item-value="uuid"
                  :items="lines"
                  :loading="loading"
                >
                  <template #item.purchase_order_uuid="{ item }">
                    {{ orderNumber(item.purchase_order_uuid) }}
                  </template>
                  <template #item.quantity="{ item }">
                    {{ `${item.quantity} ${item.quantity_unit}` }}
                  </template>
                  <template #item.unit_cost_cents="{ item }">
                    {{ formatCurrency(item.unit_cost_cents, item.currency) }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No order lines yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create order line</v-card-title>
              <v-card-text>
                <v-select
                  v-model="lineForm.purchase_order_uuid"
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
                    !lineForm.purchase_order_uuid ||
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

</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRoute } from 'vue-router'
  import type { PurchaseOrder, PurchaseOrderLine } from '@/types'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useSnackbar } from '@/composables/useSnackbar'

  const {
    getPurchaseOrders,
    getPurchaseOrderLines,
    createPurchaseOrderLine,
    normalizeText,
    toNumber,
    formatCurrency,
  } = useProcurementApi()
  const route = useRoute()
  const { showNotice } = useSnackbar()

  // Table configuration
  const lineHeaders = [
    { title: 'Order', key: 'purchase_order_uuid', sortable: true },
    { title: 'Line', key: 'line_number', sortable: true },
    { title: 'Item', key: 'item_name', sortable: true },
    { title: 'Qty', key: 'quantity', sortable: true },
    { title: 'Unit cost', key: 'unit_cost_cents', sortable: true },
  ]

  const orders = ref<PurchaseOrder[]>([])
  const lines = ref<PurchaseOrderLine[]>([])
  const loading = ref(false)
  const errorMessage = ref('')

  const itemTypeOptions = ['ingredient', 'packaging', 'service', 'equipment', 'other']
  const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']
  const currencyOptions = ['USD', 'CAD', 'EUR', 'GBP']

  const filters = reactive({
    purchase_order_uuid: null as string | null,
  })

  const lineForm = reactive({
    purchase_order_uuid: null as string | null,
    line_number: '',
    item_type: '',
    item_name: '',
    inventory_item_uuid: '',
    quantity: '',
    quantity_unit: '',
    unit_cost_cents: '',
    currency: 'USD',
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
      await loadLines()
    }
  })

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
    orders.value = await getPurchaseOrders()
  }

  async function loadLines () {
    lines.value = await getPurchaseOrderLines(filters.purchase_order_uuid ?? undefined)
  }

  async function createLine () {
    try {
      const payload = {
        purchase_order_uuid: lineForm.purchase_order_uuid,
        line_number: toNumber(lineForm.line_number),
        item_type: lineForm.item_type,
        item_name: lineForm.item_name.trim(),
        inventory_item_uuid: normalizeText(lineForm.inventory_item_uuid),
        quantity: toNumber(lineForm.quantity),
        quantity_unit: lineForm.quantity_unit,
        unit_cost_cents: toNumber(lineForm.unit_cost_cents),
        currency: lineForm.currency,
      }
      await createPurchaseOrderLine(payload)
      lineForm.purchase_order_uuid = null
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

  function orderNumber (orderUuid: string) {
    const order = orders.value.find(o => o.uuid === orderUuid)
    return order?.order_number ?? 'Unknown PO'
  }
</script>

<style scoped>
.procurement-page {
  position: relative;
}
</style>
