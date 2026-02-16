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
                <v-data-table
                  class="data-table"
                  density="compact"
                  :headers="feeHeaders"
                  item-value="uuid"
                  :items="fees"
                  :loading="loading"
                >
                  <template #item.purchase_order_uuid="{ item }">
                    {{ orderNumber(item.purchase_order_uuid) }}
                  </template>
                  <template #item.fee_type="{ item }">
                    {{ formatFeeType(item.fee_type) }}
                  </template>
                  <template #item.amount_cents="{ item }">
                    {{ formatCurrency(item.amount_cents, item.currency) }}
                  </template>
                  <template #item.updated_at="{ item }">
                    {{ formatDateTime(item.updated_at) }}
                  </template>
                  <template #no-data>
                    <div class="text-center py-4 text-medium-emphasis">No fees yet.</div>
                  </template>
                </v-data-table>
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
                <v-combobox
                  v-model="feeForm.fee_type"
                  :items="feeTypeOptions"
                  label="Fee type"
                />
                <v-text-field v-model="feeForm.amount_cents" label="Amount (cents)" type="number" />
                <v-combobox
                  v-model="feeForm.currency"
                  :items="currencyOptions"
                  label="Currency"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="saving || !feeForm.purchase_order_uuid || !feeForm.fee_type.trim() || !feeForm.amount_cents || !feeForm.currency"
                  :loading="saving"
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

</template>

<script lang="ts" setup>
  import type { PurchaseOrder, PurchaseOrderFee } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRoute } from 'vue-router'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDateTime, useFeeTypeFormatters } from '@/composables/useFormatters'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { toNumber } from '@/utils/normalize'

  const {
    getPurchaseOrders,
    getPurchaseOrderFees,
    createPurchaseOrderFee,
    formatCurrency,
  } = useProcurementApi()
  const route = useRoute()
  const { showNotice } = useSnackbar()
  const { formatFeeType } = useFeeTypeFormatters()
  // v-combobox requires plain string items — {title,value} objects cause the
  // model to be set to the full object on selection, breaking .trim() and API payloads.
  const feeTypeOptions = ['shipping', 'handling', 'tax', 'insurance', 'customs', 'freight', 'hazmat', 'other']

  // Table configuration
  const feeHeaders = [
    { title: 'Order', key: 'purchase_order_uuid', sortable: true },
    { title: 'Fee type', key: 'fee_type', sortable: true },
    { title: 'Amount', key: 'amount_cents', sortable: true },
    { title: 'Updated', key: 'updated_at', sortable: true },
  ]

  const orders = ref<PurchaseOrder[]>([])
  const fees = ref<PurchaseOrderFee[]>([])

  const { execute: executeLoad, loading, error: errorMessage } = useAsyncAction()
  const { execute: executeSave, loading: saving, error: saveError } = useAsyncAction()

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

  async function refreshAll () {
    await executeLoad(async () => {
      await Promise.all([loadOrders(), loadFees()])
    })
  }

  async function loadOrders () {
    orders.value = await getPurchaseOrders()
  }

  async function loadFees () {
    fees.value = await getPurchaseOrderFees(filters.purchase_order_uuid ?? undefined)
  }

  async function createFee () {
    await executeSave(async () => {
      const payload = {
        purchase_order_uuid: feeForm.purchase_order_uuid,
        fee_type: feeForm.fee_type.trim(),
        amount_cents: toNumber(feeForm.amount_cents),
        currency: feeForm.currency,
      }
      await createPurchaseOrderFee(payload)
      feeForm.purchase_order_uuid = null
      feeForm.fee_type = ''
      feeForm.amount_cents = ''
      feeForm.currency = 'USD'
      await loadFees()
      showNotice('Fee created')
    })
    if (saveError.value) {
      errorMessage.value = saveError.value
      showNotice(saveError.value, 'error')
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
