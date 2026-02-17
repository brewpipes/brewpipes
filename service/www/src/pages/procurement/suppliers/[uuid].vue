<template>
  <v-container class="pa-4" fluid>
    <!-- Loading state -->
    <v-alert
      v-if="loading"
      density="comfortable"
      type="info"
      variant="tonal"
    >
      Loading supplier...
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
        <v-btn size="small" variant="text" @click="goBack">
          Back to list
        </v-btn>
      </template>
    </v-alert>

    <!-- Content -->
    <template v-else-if="supplier">
      <!-- Header with back button -->
      <div class="d-flex align-center flex-wrap ga-2 mb-4">
        <v-btn
          class="mr-1"
          icon="mdi-arrow-left"
          size="small"
          variant="text"
          @click="goBack"
        />
        <div class="mr-auto">
          <div class="text-h5 font-weight-semibold">{{ supplier.name }}</div>
          <div class="text-body-2 text-medium-emphasis">
            Supplier
          </div>
        </div>
        <div class="d-flex align-center ga-1">
          <v-btn
            size="small"
            variant="text"
            @click="openEditDialog"
          >
            <v-icon class="mr-1" icon="mdi-pencil" size="small" />
            Edit
          </v-btn>
        </div>
      </div>

      <v-row>
        <!-- Supplier Information Card -->
        <v-col cols="12" md="6">
          <v-card class="section-card">
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2" icon="mdi-domain" />
              Supplier Information
            </v-card-title>
            <v-card-text>
              <v-list density="compact" lines="two">
                <v-list-item>
                  <v-list-item-title>Name</v-list-item-title>
                  <v-list-item-subtitle>{{ supplier.name }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Contact Name</v-list-item-title>
                  <v-list-item-subtitle>{{ supplier.contact_name || 'n/a' }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Email</v-list-item-title>
                  <v-list-item-subtitle>
                    <a v-if="supplier.email" :href="`mailto:${supplier.email}`">{{ supplier.email }}</a>
                    <span v-else>n/a</span>
                  </v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Phone</v-list-item-title>
                  <v-list-item-subtitle>
                    <a v-if="supplier.phone" :href="`tel:${supplier.phone}`">{{ supplier.phone }}</a>
                    <span v-else>n/a</span>
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Address Card -->
        <v-col cols="12" md="6">
          <v-card class="section-card">
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2" icon="mdi-map-marker-outline" />
              Address
            </v-card-title>
            <v-card-text>
              <v-list density="compact" lines="two">
                <v-list-item>
                  <v-list-item-title>Address</v-list-item-title>
                  <v-list-item-subtitle>
                    <template v-if="hasAddress">
                      <div v-if="supplier.address_line1">{{ supplier.address_line1 }}</div>
                      <div v-if="supplier.address_line2">{{ supplier.address_line2 }}</div>
                      <div v-if="cityRegionPostal">{{ cityRegionPostal }}</div>
                      <div v-if="supplier.country">{{ supplier.country }}</div>
                    </template>
                    <span v-else>n/a</span>
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Metadata Card -->
        <v-col cols="12" md="6">
          <v-card class="section-card">
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2" icon="mdi-information-outline" />
              Metadata
            </v-card-title>
            <v-card-text>
              <v-list density="compact" lines="two">
                <v-list-item>
                  <v-list-item-title>UUID</v-list-item-title>
                  <v-list-item-subtitle class="text-mono">{{ supplier.uuid }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Created</v-list-item-title>
                  <v-list-item-subtitle>{{ formatDateTime(supplier.created_at) }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Updated</v-list-item-title>
                  <v-list-item-subtitle>{{ formatDateTime(supplier.updated_at) }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Purchase Orders Card -->
        <v-col cols="12">
          <v-card class="section-card">
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2" icon="mdi-clipboard-list-outline" />
              Purchase Orders
              <v-spacer />
              <v-chip class="ml-2" size="small" variant="tonal">
                {{ purchaseOrders.length }}
              </v-chip>
            </v-card-title>
            <v-card-text>
              <v-alert
                v-if="ordersError"
                class="mb-3"
                density="compact"
                type="error"
                variant="tonal"
              >
                {{ ordersError }}
              </v-alert>

              <v-progress-linear v-if="ordersLoading" color="primary" indeterminate />

              <template v-else-if="purchaseOrders.length > 0">
                <!-- Desktop table view -->
                <div class="d-none d-md-block">
                  <v-data-table
                    class="data-table clickable-rows"
                    density="compact"
                    :headers="orderHeaders"
                    item-value="uuid"
                    :items="purchaseOrders"
                    @click:row="handleOrderRowClick"
                  >
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
                    <template #item.created_at="{ item }">
                      {{ formatDate(item.created_at) }}
                    </template>
                  </v-data-table>
                </div>

                <!-- Mobile card view -->
                <div class="d-md-none">
                  <v-card
                    v-for="order in purchaseOrders"
                    :key="order.uuid"
                    class="mb-3"
                    variant="outlined"
                    @click="router.push(`/procurement/purchase-orders/${order.uuid}`)"
                  >
                    <v-card-title class="d-flex align-center py-2 text-body-1">
                      <span class="font-weight-medium">{{ order.order_number }}</span>
                      <v-spacer />
                      <v-chip
                        :color="getPurchaseOrderStatusColor(order.status)"
                        size="x-small"
                        variant="flat"
                      >
                        {{ formatPurchaseOrderStatus(order.status) }}
                      </v-chip>
                    </v-card-title>

                    <v-card-text class="pt-0">
                      <div class="d-flex justify-space-between text-body-2 mb-1">
                        <span class="text-medium-emphasis">Expected</span>
                        <span>{{ formatDate(order.expected_at) }}</span>
                      </div>
                      <div class="d-flex justify-space-between text-body-2 mb-1">
                        <span class="text-medium-emphasis">Created</span>
                        <span>{{ formatDate(order.created_at) }}</span>
                      </div>
                    </v-card-text>
                  </v-card>
                </div>
              </template>

              <div
                v-else-if="!ordersLoading"
                class="text-center py-4 text-medium-emphasis"
              >
                No purchase orders for this supplier.
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>

    <!-- Edit Supplier Dialog -->
    <SupplierCreateEditDialog
      v-model="editDialogOpen"
      :edit-supplier="supplier"
      :saving="editSaving"
      @submit="handleSaveSupplier"
    />
  </v-container>
</template>

<script lang="ts" setup>
  import type { PurchaseOrder, Supplier } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import SupplierCreateEditDialog from '@/components/procurement/SupplierCreateEditDialog.vue'
  import type { SupplierFormData } from '@/components/procurement/SupplierCreateEditDialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDate, formatDateTime, usePurchaseOrderStatusFormatters } from '@/composables/useFormatters'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useRouteUuid } from '@/composables/useRouteUuid'
  import { useSnackbar } from '@/composables/useSnackbar'

  const router = useRouter()
  const { uuid: routeUuid } = useRouteUuid()
  const { showNotice } = useSnackbar()

  const {
    getSupplier,
    updateSupplier,
    getPurchaseOrders,
  } = useProcurementApi()

  const {
    formatPurchaseOrderStatus,
    getPurchaseOrderStatusColor,
  } = usePurchaseOrderStatusFormatters()

  // State
  const supplier = ref<Supplier | null>(null)
  const purchaseOrders = ref<PurchaseOrder[]>([])

  const { execute: executeLoad, loading, error: loadError } = useAsyncAction()
  // Start loading immediately to avoid flash of empty state before onMounted
  loading.value = true
  const { execute: executeOrdersLoad, loading: ordersLoading, error: ordersError } = useAsyncAction()
  const { execute: executeEdit, loading: editSaving, error: editError } = useAsyncAction()

  const error = computed(() => loadError.value || null)

  // Edit dialog state
  const editDialogOpen = ref(false)

  // Purchase order table headers
  const orderHeaders = [
    { title: 'Order', key: 'order_number', sortable: true },
    { title: 'Status', key: 'status', sortable: true },
    { title: 'Expected', key: 'expected_at', sortable: true },
    { title: 'Created', key: 'created_at', sortable: true },
  ]

  // Computed
  const hasAddress = computed(() => {
    if (!supplier.value) return false
    return !!(
      supplier.value.address_line1 ||
      supplier.value.address_line2 ||
      supplier.value.city ||
      supplier.value.region ||
      supplier.value.postal_code ||
      supplier.value.country
    )
  })

  const cityRegionPostal = computed(() => {
    if (!supplier.value) return ''
    const parts: string[] = []
    if (supplier.value.city) parts.push(supplier.value.city)
    if (supplier.value.region) parts.push(supplier.value.region)
    const cityRegion = parts.join(', ')
    if (supplier.value.postal_code) {
      return cityRegion ? `${cityRegion} ${supplier.value.postal_code}` : supplier.value.postal_code
    }
    return cityRegion
  })

  // Lifecycle
  onMounted(async () => {
    await loadSupplier()
  })

  // Methods
  async function loadSupplier () {
    const uuid = routeUuid.value
    if (!uuid) {
      loadError.value = 'Invalid supplier UUID'
      loading.value = false
      return
    }

    await executeLoad(async () => {
      supplier.value = await getSupplier(uuid)
      // Load purchase orders in parallel (non-blocking)
      await loadPurchaseOrders()
    })
    // Provide user-friendly error messages
    if (loadError.value) {
      loadError.value = loadError.value.includes('404')
        ? 'Supplier not found'
        : 'Failed to load supplier. Please try again.'
    }
  }

  async function loadPurchaseOrders () {
    if (!supplier.value) return
    await executeOrdersLoad(async () => {
      purchaseOrders.value = await getPurchaseOrders(supplier.value!.uuid)
    })
  }

  function goBack () {
    router.push('/procurement/suppliers')
  }

  function openEditDialog () {
    editDialogOpen.value = true
  }

  async function handleSaveSupplier (data: SupplierFormData) {
    if (!supplier.value) return

    await executeEdit(async () => {
      supplier.value = await updateSupplier(supplier.value!.uuid, data)
      showNotice('Supplier updated')
      editDialogOpen.value = false
    })
    if (editError.value) {
      showNotice(editError.value, 'error')
    }
  }

  function handleOrderRowClick (_event: Event, row: { item: PurchaseOrder }) {
    router.push(`/procurement/purchase-orders/${row.item.uuid}`)
  }
</script>

<style scoped>
.text-mono {
  font-family: monospace;
  font-size: 0.85em;
}

.clickable-rows :deep(tbody tr) {
  cursor: pointer;
}

.clickable-rows :deep(tbody tr:hover) {
  background-color: rgba(var(--v-theme-on-surface), 0.04);
}
</style>
