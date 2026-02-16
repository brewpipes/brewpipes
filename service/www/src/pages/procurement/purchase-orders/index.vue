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
                <v-data-table
                  class="data-table clickable-rows"
                  density="compact"
                  :headers="orderHeaders"
                  item-value="uuid"
                  :items="orders"
                  :loading="loading"
                  @click:row="handleRowClick"
                >
                  <template #item.supplier_uuid="{ item }">
                    {{ supplierName(item.supplier_uuid) }}
                  </template>
                  <template #item.status="{ item }">
                    <v-chip
                      :color="getPurchaseOrderStatusColor(item.status)"
                      size="x-small"
                      variant="flat"
                    >
                      {{ formatPurchaseOrderStatus(item.status) }}
                    </v-chip>
                  </template>
                  <template #item.expected_at="{ item }">
                    {{ formatDate(item.expected_at) }}
                  </template>
                  <template #item.actions="{ item }">
                    <v-btn
                      icon="mdi-pencil"
                      size="x-small"
                      variant="text"
                      @click.stop="openEditDialog(item)"
                    />
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No purchase orders yet.</div>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>

  <PurchaseOrderCreateEditDialog
    v-model="orderDialog"
    :edit-order="editingOrder"
    :error-message="dialogError"
    :saving="saving"
    :suppliers="supplierSelectItems"
    @submit="saveOrder"
  />
</template>

<script lang="ts" setup>
  import type { PurchaseOrder, Supplier } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import PurchaseOrderCreateEditDialog from '@/components/procurement/PurchaseOrderCreateEditDialog.vue'
  import type { PurchaseOrderCreateSubmitData, PurchaseOrderEditSubmitData } from '@/components/procurement/PurchaseOrderCreateEditDialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDate, usePurchaseOrderStatusFormatters } from '@/composables/useFormatters'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useSnackbar } from '@/composables/useSnackbar'

  const {
    getSuppliers,
    getPurchaseOrders,
    createPurchaseOrder,
    updatePurchaseOrder,
  } = useProcurementApi()
  const { showNotice } = useSnackbar()
  const router = useRouter()
  const {
    formatPurchaseOrderStatus,
    getPurchaseOrderStatusColor,
  } = usePurchaseOrderStatusFormatters()

  // Table configuration
  const orderHeaders = [
    { title: 'Order', key: 'order_number', sortable: true },
    { title: 'Supplier', key: 'supplier_uuid', sortable: true },
    { title: 'Status', key: 'status', sortable: true },
    { title: 'Expected', key: 'expected_at', sortable: true },
    { title: '', key: 'actions', sortable: false, align: 'end' as const, width: '60px' },
  ]

  const suppliers = ref<Supplier[]>([])
  const orders = ref<PurchaseOrder[]>([])

  const { execute: executeLoad, loading, error: errorMessage } = useAsyncAction()
  const { execute: executeSave, loading: saving, error: dialogError } = useAsyncAction()

  // Dialog state
  const orderDialog = ref(false)
  const editingOrder = ref<PurchaseOrder | null>(null)

  const filters = reactive({
    supplier_uuid: null as string | null,
  })

  const supplierSelectItems = computed(() =>
    suppliers.value.map(supplier => ({
      title: supplier.name,
      value: supplier.uuid,
    })),
  )

  onMounted(async () => {
    await refreshAll()
  })

  async function refreshAll () {
    await executeLoad(async () => {
      await Promise.all([loadSuppliers(), loadOrders()])
    })
  }

  async function loadSuppliers () {
    suppliers.value = await getSuppliers()
  }

  async function loadOrders () {
    orders.value = await getPurchaseOrders(filters.supplier_uuid ?? undefined)
  }

  function openCreateDialog () {
    editingOrder.value = null
    dialogError.value = ''
    orderDialog.value = true
  }

  function openEditDialog (order: PurchaseOrder) {
    editingOrder.value = order
    dialogError.value = ''
    orderDialog.value = true
  }

  async function saveOrder (data: PurchaseOrderCreateSubmitData | PurchaseOrderEditSubmitData) {
    await executeSave(async () => {
      if (editingOrder.value) {
        await updatePurchaseOrder(editingOrder.value.uuid, data as PurchaseOrderEditSubmitData)
        showNotice('Purchase order updated')
      } else {
        await createPurchaseOrder(data as PurchaseOrderCreateSubmitData)
        showNotice('Purchase order created')
      }
      await loadOrders()
      orderDialog.value = false
      editingOrder.value = null
      dialogError.value = ''
    })
  }

  function supplierName (supplierUuid: string) {
    const supplier = suppliers.value.find(s => s.uuid === supplierUuid)
    return supplier?.name ?? 'Unknown Supplier'
  }

  function handleRowClick (_event: Event, row: { item: PurchaseOrder }) {
    router.push(`/procurement/purchase-orders/${row.item.uuid}`)
  }
</script>

<style scoped>
.procurement-page {
  position: relative;
}

.clickable-rows :deep(tbody tr) {
  cursor: pointer;
}

.clickable-rows :deep(tbody tr:hover) {
  background-color: rgba(var(--v-theme-on-surface), 0.04);
}
</style>
