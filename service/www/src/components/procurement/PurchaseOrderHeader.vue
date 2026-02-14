<template>
  <v-card class="mb-4" variant="flat">
    <v-card-title class="d-flex align-center flex-wrap ga-2">
      <v-btn
        class="mr-2"
        icon="mdi-arrow-left"
        size="small"
        variant="text"
        @click="emit('back')"
      />
      <span class="text-h6">{{ order.order_number }}</span>
      <v-chip class="ml-2" color="primary" size="small" variant="tonal">
        {{ supplierName }}
      </v-chip>
      <v-spacer />
      <PurchaseOrderStatusMenu
        :disabled="statusChanging"
        :saving="statusChanging"
        :status="order.status"
        @update:status="handleStatusChange"
      />
    </v-card-title>
    <v-card-text>
      <v-row dense>
        <v-col cols="6" sm="3">
          <div class="text-caption text-medium-emphasis">Ordered</div>
          <div class="text-body-2">{{ formatDate(order.ordered_at) }}</div>
        </v-col>
        <v-col cols="6" sm="3">
          <div class="text-caption text-medium-emphasis">Expected</div>
          <div class="text-body-2">{{ formatDate(order.expected_at) }}</div>
        </v-col>
        <v-col cols="6" sm="3">
          <div class="text-caption text-medium-emphasis">Lines Total</div>
          <div class="text-body-2">{{ formatCurrency(linesTotal, currency) }}</div>
        </v-col>
        <v-col cols="6" sm="3">
          <div class="text-caption text-medium-emphasis">Order Total</div>
          <div class="text-body-2 font-weight-medium">{{ formatCurrency(orderTotal, currency) }}</div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { PurchaseOrder } from '@/types'
  import { ref } from 'vue'
  import { formatDate } from '@/composables/useFormatters'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import PurchaseOrderStatusMenu from './PurchaseOrderStatusMenu.vue'

  const props = defineProps<{
    order: PurchaseOrder
    supplierName: string
    linesTotal: number
    orderTotal: number
    currency: string
  }>()

  const emit = defineEmits<{
    'back': []
    'status-changed': [order: PurchaseOrder]
  }>()

  const { updatePurchaseOrder, formatCurrency } = useProcurementApi()
  const { showNotice } = useSnackbar()

  const statusChanging = ref(false)

  async function handleStatusChange (newStatus: string) {
    statusChanging.value = true
    try {
      const updated = await updatePurchaseOrder(props.order.uuid, { status: newStatus })
      emit('status-changed', updated)
      showNotice('Status updated')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update status'
      showNotice(message, 'error')
    } finally {
      statusChanging.value = false
    }
  }
</script>
