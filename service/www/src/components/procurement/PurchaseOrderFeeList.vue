<template>
  <v-card class="mb-4" variant="outlined">
    <v-card-title class="d-flex align-center">
      <span class="text-subtitle-1">Fees</span>
      <v-spacer />
      <v-btn
        aria-label="Add fee"
        color="primary"
        icon="mdi-plus"
        size="small"
        variant="text"
        @click="openCreateDialog"
      />
    </v-card-title>
    <v-card-text>
      <!-- Desktop table view -->
      <v-data-table
        v-if="!$vuetify.display.xs"
        class="data-table"
        density="compact"
        :headers="headers"
        item-value="uuid"
        :items="fees"
        :loading="loading"
      >
        <template #item.fee_type="{ item }">
          {{ formatFeeType(item.fee_type) }}
        </template>
        <template #item.amount_cents="{ item }">
          {{ formatCurrency(item.amount_cents, item.currency) }}
        </template>
        <template #item.actions="{ item }">
          <v-menu>
            <template #activator="{ props: menuProps }">
              <v-btn
                v-bind="menuProps"
                icon="mdi-dots-vertical"
                size="x-small"
                variant="text"
              />
            </template>
            <v-list density="compact">
              <v-list-item @click="openEditDialog(item)">
                <template #prepend>
                  <v-icon icon="mdi-pencil" size="small" />
                </template>
                <v-list-item-title>Edit</v-list-item-title>
              </v-list-item>
              <v-list-item class="text-error" @click="confirmDelete(item)">
                <template #prepend>
                  <v-icon icon="mdi-delete" size="small" />
                </template>
                <v-list-item-title>Delete</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </template>
        <template #no-data>
          <div class="text-center py-4 text-medium-emphasis">No fees yet.</div>
        </template>
      </v-data-table>

      <!-- Mobile card view -->
      <div v-else>
        <v-card
          v-for="fee in fees"
          :key="fee.uuid"
          class="mb-2"
          variant="tonal"
        >
          <v-card-text class="pa-3">
            <div class="d-flex justify-space-between align-start">
              <div>
                <div class="text-body-2 font-weight-medium">
                  {{ formatFeeType(fee.fee_type) }}
                </div>
                <div class="text-body-2 mt-1">
                  {{ formatCurrency(fee.amount_cents, fee.currency) }}
                </div>
              </div>
              <v-menu>
                <template #activator="{ props: menuProps }">
                  <v-btn
                    v-bind="menuProps"
                    icon="mdi-dots-vertical"
                    size="small"
                    variant="text"
                  />
                </template>
                <v-list density="compact">
                  <v-list-item @click="openEditDialog(fee)">
                    <template #prepend>
                      <v-icon icon="mdi-pencil" size="small" />
                    </template>
                    <v-list-item-title>Edit</v-list-item-title>
                  </v-list-item>
                  <v-list-item class="text-error" @click="confirmDelete(fee)">
                    <template #prepend>
                      <v-icon icon="mdi-delete" size="small" />
                    </template>
                    <v-list-item-title>Delete</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </div>
          </v-card-text>
        </v-card>
        <div v-if="fees.length === 0" class="text-center py-4 text-medium-emphasis">
          No fees yet.
        </div>
      </div>
    </v-card-text>
  </v-card>

  <!-- Fee Dialog -->
  <PurchaseOrderFeeDialog
    v-model="feeDialogOpen"
    :fee="editingFee"
    :saving="saving"
    @cancel="closeFeeDialog"
    @save="handleSaveFee"
  />

  <!-- Delete Confirmation Dialog -->
  <v-dialog v-model="deleteDialogOpen" max-width="400">
    <v-card>
      <v-card-title class="text-h6">Delete fee?</v-card-title>
      <v-card-text>
        Are you sure you want to delete the {{ deletingFee ? formatFeeType(deletingFee.fee_type) : '' }} fee?
        This action cannot be undone.
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="deleting" variant="text" @click="deleteDialogOpen = false">Cancel</v-btn>
        <v-btn
          color="error"
          :loading="deleting"
          @click="handleDelete"
        >
          Delete
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { PurchaseOrderFee } from '@/types'
  import { ref } from 'vue'
  import { useFeeTypeFormatters } from '@/composables/useFormatters'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import PurchaseOrderFeeDialog, { type FeeForm } from './PurchaseOrderFeeDialog.vue'

  const props = defineProps<{
    purchaseOrderUuid: string
    fees: PurchaseOrderFee[]
    loading?: boolean
  }>()

  const emit = defineEmits<{
    refresh: []
  }>()

  const {
    createPurchaseOrderFee,
    updatePurchaseOrderFee,
    deletePurchaseOrderFee,
    formatCurrency,
  } = useProcurementApi()
  const { showNotice } = useSnackbar()
  const { formatFeeType } = useFeeTypeFormatters()

  const headers = [
    { title: 'Fee Type', key: 'fee_type', sortable: true },
    { title: 'Amount', key: 'amount_cents', sortable: true },
    { title: '', key: 'actions', sortable: false, align: 'end' as const, width: '60px' },
  ]

  const feeDialogOpen = ref(false)
  const editingFee = ref<PurchaseOrderFee | null>(null)
  const saving = ref(false)

  const deleteDialogOpen = ref(false)
  const deletingFee = ref<PurchaseOrderFee | null>(null)
  const deleting = ref(false)

  function openCreateDialog () {
    editingFee.value = null
    feeDialogOpen.value = true
  }

  function openEditDialog (fee: PurchaseOrderFee) {
    editingFee.value = fee
    feeDialogOpen.value = true
  }

  function closeFeeDialog () {
    feeDialogOpen.value = false
    editingFee.value = null
  }

  function confirmDelete (fee: PurchaseOrderFee) {
    deletingFee.value = fee
    deleteDialogOpen.value = true
  }

  async function handleSaveFee (form: FeeForm) {
    saving.value = true
    try {
      if (editingFee.value) {
        await updatePurchaseOrderFee(editingFee.value.uuid, {
          fee_type: form.fee_type.trim(),
          amount_cents: form.amount_cents ?? undefined,
          currency: form.currency,
        })
        showNotice('Fee updated')
      } else {
        await createPurchaseOrderFee({
          purchase_order_uuid: props.purchaseOrderUuid,
          fee_type: form.fee_type.trim(),
          amount_cents: form.amount_cents,
          currency: form.currency,
        })
        showNotice('Fee added')
      }
      closeFeeDialog()
      emit('refresh')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save fee'
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }

  async function handleDelete () {
    if (!deletingFee.value) return

    deleting.value = true
    try {
      await deletePurchaseOrderFee(deletingFee.value.uuid)
      showNotice('Fee deleted')
      deleteDialogOpen.value = false
      deletingFee.value = null
      emit('refresh')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to delete fee'
      showNotice(message, 'error')
    } finally {
      deleting.value = false
    }
  }
</script>

<style scoped>
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
