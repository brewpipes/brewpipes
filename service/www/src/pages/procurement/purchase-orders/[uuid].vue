<template>
  <v-container class="pa-4" fluid>
    <!-- Loading state -->
    <v-alert
      v-if="loading"
      density="comfortable"
      type="info"
      variant="tonal"
    >
      Loading purchase order...
    </v-alert>

    <!-- Error state -->
    <v-alert
      v-else-if="error"
      density="comfortable"
      type="error"
      variant="tonal"
    >
      {{ error }}
      <template #append>
        <v-btn size="small" variant="text" @click="router.push('/procurement/purchase-orders')">
          Back to list
        </v-btn>
      </template>
    </v-alert>

    <!-- Content -->
    <template v-else-if="order">
      <!-- Header -->
      <PurchaseOrderHeader
        :currency="primaryCurrency"
        :lines-total="linesTotal"
        :order="order"
        :order-total="orderTotal"
        :supplier-name="supplierName"
        @back="handleBack"
        @status-changed="handleStatusChanged"
      />

      <!-- Action buttons -->
      <div class="d-flex flex-wrap ga-2 mb-4">
        <v-btn
          color="primary"
          prepend-icon="mdi-pencil"
          size="small"
          variant="tonal"
          @click="openEditDialog"
        >
          <span class="d-none d-sm-inline">Edit Order</span>
          <span class="d-sm-none">Edit</span>
        </v-btn>
        <v-btn
          v-if="canReceive"
          color="success"
          prepend-icon="mdi-package-down"
          size="small"
          variant="tonal"
          @click="handleReceiveShipment"
        >
          <span class="d-none d-sm-inline">Receive Shipment</span>
          <span class="d-sm-none">Receive</span>
        </v-btn>
      </div>

      <!-- Line Items -->
      <PurchaseOrderLineList
        :lines="lines"
        :loading="linesLoading"
        :purchase-order-uuid="order.uuid"
        @refresh="loadLines"
      />

      <!-- Fees -->
      <PurchaseOrderFeeList
        :fees="fees"
        :loading="feesLoading"
        :purchase-order-uuid="order.uuid"
        @refresh="loadFees"
      />

      <!-- Summary -->
      <v-card class="mb-4" variant="outlined">
        <v-card-title class="text-subtitle-1">Order Summary</v-card-title>
        <v-card-text>
          <v-row dense>
            <v-col class="text-right" cols="8" sm="10">
              <span class="text-medium-emphasis">Subtotal (Lines):</span>
            </v-col>
            <v-col class="text-right" cols="4" sm="2">
              {{ formatCurrency(linesTotal, primaryCurrency) }}
            </v-col>
          </v-row>
          <v-row dense>
            <v-col class="text-right" cols="8" sm="10">
              <span class="text-medium-emphasis">Fees:</span>
            </v-col>
            <v-col class="text-right" cols="4" sm="2">
              {{ formatCurrency(feesTotal, primaryCurrency) }}
            </v-col>
          </v-row>
          <v-divider class="my-2" />
          <v-row dense>
            <v-col class="text-right" cols="8" sm="10">
              <span class="font-weight-medium">Order Total:</span>
            </v-col>
            <v-col class="text-right font-weight-bold" cols="4" sm="2">
              {{ formatCurrency(orderTotal, primaryCurrency) }}
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>

      <!-- Notes -->
      <v-card v-if="order.notes" class="mb-4" variant="outlined">
        <v-card-title class="text-subtitle-1">Notes</v-card-title>
        <v-card-text>
          <p class="text-body-2" style="white-space: pre-wrap;">{{ order.notes }}</p>
        </v-card-text>
      </v-card>
    </template>
  </v-container>

  <!-- Receive Shipment Wizard -->
  <ReceiveShipmentWizard
    v-if="order"
    v-model="receiveWizardOpen"
    :lines="lines"
    :purchase-order="order"
    :supplier-name="supplierName"
    @received="handleReceived"
  />

  <!-- Edit Order Dialog -->
  <PurchaseOrderEditDialog
    v-if="order"
    v-model="editDialogOpen"
    :error-message="editError"
    :purchase-order="order"
    :saving="editSaving"
    @submit="saveEdit"
  />
</template>

<script lang="ts" setup>
  import type { PurchaseOrder, PurchaseOrderFee, PurchaseOrderLine, Supplier } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import PurchaseOrderEditDialog from '@/components/procurement/PurchaseOrderEditDialog.vue'
  import type { PurchaseOrderEditFormData } from '@/components/procurement/PurchaseOrderEditDialog.vue'
  import PurchaseOrderFeeList from '@/components/procurement/PurchaseOrderFeeList.vue'
  import PurchaseOrderHeader from '@/components/procurement/PurchaseOrderHeader.vue'
  import PurchaseOrderLineList from '@/components/procurement/PurchaseOrderLineList.vue'
  import ReceiveShipmentWizard from '@/components/procurement/ReceiveShipmentWizard.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useRouteUuid } from '@/composables/useRouteUuid'
  import { useSnackbar } from '@/composables/useSnackbar'

  const router = useRouter()
  const { uuid: routeUuid } = useRouteUuid()
  const { showNotice } = useSnackbar()

  const {
    getPurchaseOrder,
    updatePurchaseOrder,
    getSupplier,
    getPurchaseOrderLines,
    getPurchaseOrderFees,
    formatCurrency,
  } = useProcurementApi()

  // State
  const order = ref<PurchaseOrder | null>(null)
  const supplier = ref<Supplier | null>(null)
  const lines = ref<PurchaseOrderLine[]>([])
  const fees = ref<PurchaseOrderFee[]>([])

  const { execute: executeLoad, loading, error: loadError } = useAsyncAction()
  // Start loading immediately to avoid flash of empty state before onMounted
  loading.value = true
  const { execute: executeLinesLoad, loading: linesLoading } = useAsyncAction({
    onError: (message) => showNotice(message, 'error'),
  })
  const { execute: executeFeesLoad, loading: feesLoading } = useAsyncAction({
    onError: (message) => showNotice(message, 'error'),
  })
  const { execute: executeEdit, loading: editSaving, error: editError } = useAsyncAction()

  // The template uses `error` (not `loadError`) for the error display
  const error = computed(() => loadError.value || null)

  // Edit dialog state
  const editDialogOpen = ref(false)

  // Receive shipment wizard state
  const receiveWizardOpen = ref(false)

  // Computed
  const supplierName = computed(() => supplier.value?.name ?? 'Unknown Supplier')

  const linesTotal = computed(() =>
    lines.value.reduce((sum, line) => sum + (line.quantity * line.unit_cost_cents), 0),
  )

  const feesTotal = computed(() =>
    fees.value.reduce((sum, fee) => sum + fee.amount_cents, 0),
  )

  const orderTotal = computed(() => linesTotal.value + feesTotal.value)

  const primaryCurrency = computed(() => {
    // Use currency from first line, or first fee, or default to USD
    const firstLine = lines.value[0]
    if (firstLine) {
      return firstLine.currency
    }
    const firstFee = fees.value[0]
    if (firstFee) {
      return firstFee.currency
    }
    return 'USD'
  })

  const canReceive = computed(() => {
    return order.value?.status === 'confirmed' || order.value?.status === 'partially_received'
  })

  // Lifecycle
  onMounted(async () => {
    await loadOrder()
  })

  // Methods
  async function loadOrder () {
    const uuid = routeUuid.value
    if (!uuid) {
      loadError.value = 'Invalid purchase order UUID'
      return
    }

    await executeLoad(async () => {
      order.value = await getPurchaseOrder(uuid)

      // Load supplier and lines/fees in parallel
      await Promise.all([
        loadSupplier(),
        loadLines(),
        loadFees(),
      ])
    })
    // Provide user-friendly error messages
    if (loadError.value) {
      loadError.value = loadError.value.includes('404')
        ? 'Purchase order not found'
        : 'Failed to load purchase order. Please try again.'
    }
  }

  async function loadSupplier () {
    if (!order.value) return
    try {
      supplier.value = await getSupplier(order.value.supplier_uuid)
    } catch {
      // Supplier load failure is non-critical
      supplier.value = null
    }
  }

  async function loadLines () {
    if (!order.value) return
    await executeLinesLoad(async () => {
      lines.value = await getPurchaseOrderLines(order.value!.uuid)
    })
  }

  async function loadFees () {
    if (!order.value) return
    await executeFeesLoad(async () => {
      fees.value = await getPurchaseOrderFees(order.value!.uuid)
    })
  }

  function handleBack () {
    router.push('/procurement/purchase-orders')
  }

  function handleStatusChanged (updatedOrder: PurchaseOrder) {
    order.value = updatedOrder
  }

  function handleReceiveShipment () {
    receiveWizardOpen.value = true
  }

  async function handleReceived () {
    showNotice('Shipment received successfully', 'success')
    await loadOrder()
  }

  function openEditDialog () {
    if (!order.value) return
    editError.value = ''
    editDialogOpen.value = true
  }

  async function saveEdit (data: PurchaseOrderEditFormData) {
    if (!order.value) return

    await executeEdit(async () => {
      order.value = await updatePurchaseOrder(order.value!.uuid, data)
      showNotice('Purchase order updated')
      editDialogOpen.value = false
      editError.value = ''
    })
  }
</script>
