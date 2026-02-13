<template>
  <v-container class="procurement-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <span>Purchase orders</span>
        <div class="card-title-actions">
          <v-btn
            :icon="$vuetify.display.xs"
            :loading="loading"
            size="small"
            variant="text"
            @click="refreshAll"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-refresh" />
            <span v-else>Refresh</span>
          </v-btn>
          <v-btn
            color="primary"
            :icon="$vuetify.display.xs"
            size="small"
            variant="text"
            @click="openCreateDialog"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-plus" />
            <span v-else>New order</span>
          </v-btn>
        </div>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12">
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
                  v-model="filters.supplier_uuid"
                  clearable
                  :items="supplierSelectItems"
                  label="Filter by supplier"
                />
                <v-table class="data-table" density="compact">
                  <thead>
                    <tr>
                      <th>Order</th>
                      <th>Supplier</th>
                      <th>Status</th>
                      <th>Expected</th>
                      <th />
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="order in orders" :key="order.uuid">
                      <td>{{ order.order_number }}</td>
                      <td>{{ supplierName(order.supplier_uuid) }}</td>
                      <td>{{ order.status }}</td>
                      <td>{{ formatDateTime(order.expected_at) }}</td>
                      <td class="text-right">
                        <v-btn
                          icon="mdi-pencil"
                          size="x-small"
                          variant="text"
                          @click.stop="openEditDialog(order)"
                        />
                        <v-btn
                          size="x-small"
                          variant="text"
                          @click="openLines(order.uuid)"
                        >
                          Lines
                        </v-btn>
                        <v-btn
                          size="x-small"
                          variant="text"
                          @click="openFees(order.uuid)"
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
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <v-dialog v-model="orderDialog" :max-width="$vuetify.display.xs ? '100%' : 640">
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditMode ? 'Edit purchase order' : 'Create purchase order' }}
      </v-card-title>
      <v-card-text>
        <v-alert
          v-if="dialogError"
          class="mb-3"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ dialogError }}
        </v-alert>
        <v-row>
          <v-col cols="12">
            <v-select
              v-model="orderForm.supplier_uuid"
              :disabled="isEditMode"
              :hint="isEditMode ? 'Supplier cannot be changed after creation' : ''"
              :items="supplierSelectItems"
              label="Supplier"
              :persistent-hint="isEditMode"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="orderForm.order_number" label="Order number" />
          </v-col>
          <v-col cols="12" md="6">
            <v-select v-model="orderForm.status" clearable :items="statusOptions" label="Status" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="orderForm.ordered_at" label="Ordered at" type="datetime-local" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="orderForm.expected_at" label="Expected at" type="datetime-local" />
          </v-col>
          <v-col cols="12">
            <v-textarea v-model="orderForm.notes" auto-grow label="Notes" rows="2" />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="closeDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="saveOrder"
        >
          {{ isEditMode ? 'Save changes' : 'Add purchase order' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import {
    type PurchaseOrder,
    type Supplier,
    useProcurementApi,
  } from '@/composables/useProcurementApi'

  const {
    getSuppliers,
    getPurchaseOrders,
    createPurchaseOrder,
    updatePurchaseOrder,
    normalizeText,
    normalizeDateTime,
    formatDateTime,
  } = useProcurementApi()
  const router = useRouter()

  const suppliers = ref<Supplier[]>([])
  const orders = ref<PurchaseOrder[]>([])
  const loading = ref(false)
  const errorMessage = ref('')

  // Dialog state
  const orderDialog = ref(false)
  const editingOrder = ref<PurchaseOrder | null>(null)
  const saving = ref(false)
  const dialogError = ref('')

  const isEditMode = computed(() => editingOrder.value !== null)

  const statusOptions = [
    'draft',
    'submitted',
    'confirmed',
    'partially_received',
    'received',
    'cancelled',
  ]

  const filters = reactive({
    supplier_uuid: null as string | null,
  })

  const orderForm = reactive({
    supplier_uuid: null as string | null,
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
    suppliers.value.map(supplier => ({
      title: supplier.name,
      value: supplier.uuid,
    })),
  )

  const isFormValid = computed(() => {
    return orderForm.supplier_uuid !== null && orderForm.order_number.trim().length > 0
  })

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

  async function loadSuppliers () {
    suppliers.value = await getSuppliers()
  }

  async function loadOrders () {
    orders.value = await getPurchaseOrders(filters.supplier_uuid ?? undefined)
  }

  function resetForm () {
    orderForm.supplier_uuid = null
    orderForm.order_number = ''
    orderForm.status = ''
    orderForm.ordered_at = ''
    orderForm.expected_at = ''
    orderForm.notes = ''
  }

  function openCreateDialog () {
    editingOrder.value = null
    dialogError.value = ''
    resetForm()
    orderDialog.value = true
  }

  function openEditDialog (order: PurchaseOrder) {
    editingOrder.value = order
    dialogError.value = ''
    orderForm.supplier_uuid = order.supplier_uuid
    orderForm.order_number = order.order_number
    orderForm.status = order.status || ''
    orderForm.ordered_at = order.ordered_at ? toDateTimeLocal(order.ordered_at) : ''
    orderForm.expected_at = order.expected_at ? toDateTimeLocal(order.expected_at) : ''
    orderForm.notes = order.notes || ''
    orderDialog.value = true
  }

  function closeDialog () {
    orderDialog.value = false
    editingOrder.value = null
    dialogError.value = ''
    resetForm()
  }

  function toDateTimeLocal (isoString: string): string {
    const date = new Date(isoString)
    const offset = date.getTimezoneOffset()
    const local = new Date(date.getTime() - offset * 60 * 1000)
    return local.toISOString().slice(0, 16)
  }

  async function saveOrder () {
    saving.value = true
    dialogError.value = ''

    try {
      if (isEditMode.value && editingOrder.value) {
        const payload = {
          order_number: orderForm.order_number.trim(),
          status: normalizeText(orderForm.status) ?? undefined,
          ordered_at: normalizeDateTime(orderForm.ordered_at),
          expected_at: normalizeDateTime(orderForm.expected_at),
          notes: normalizeText(orderForm.notes),
        }
        await updatePurchaseOrder(editingOrder.value.uuid, payload)
        showNotice('Purchase order updated')
      } else {
        const payload = {
          supplier_uuid: orderForm.supplier_uuid!,
          order_number: orderForm.order_number.trim(),
          status: normalizeText(orderForm.status),
          ordered_at: normalizeDateTime(orderForm.ordered_at),
          expected_at: normalizeDateTime(orderForm.expected_at),
          notes: normalizeText(orderForm.notes),
        }
        await createPurchaseOrder(payload)
        showNotice('Purchase order created')
      }
      await loadOrders()
      closeDialog()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save purchase order'
      dialogError.value = message
    } finally {
      saving.value = false
    }
  }

  function supplierName (supplierUuid: string) {
    const supplier = suppliers.value.find(s => s.uuid === supplierUuid)
    return supplier?.name ?? 'Unknown Supplier'
  }

  function openLines (orderUuid: string) {
    router.push({
      path: '/procurement/purchase-order-lines',
      query: { purchase_order_uuid: orderUuid },
    })
  }

  function openFees (orderUuid: string) {
    router.push({
      path: '/procurement/purchase-order-fees',
      query: { purchase_order_uuid: orderUuid },
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

.card-title-responsive {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.card-title-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
}

.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.data-table {
  overflow-x: auto;
}

.data-table :deep(.v-table__wrapper) {
  overflow-x: auto;
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  white-space: nowrap;
}

.data-table :deep(td) {
  font-size: 0.85rem;
}
</style>
