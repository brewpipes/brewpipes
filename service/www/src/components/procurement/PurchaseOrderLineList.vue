<template>
  <v-card class="mb-4" variant="outlined">
    <v-card-title class="d-flex align-center">
      <span class="text-subtitle-1">Line Items</span>
      <v-spacer />
      <v-btn
        aria-label="Add line item"
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
        :items="lines"
        :loading="loading"
      >
        <template #item.item_type="{ item }">
          {{ formatLineItemType(item.item_type) }}
        </template>
        <template #item.quantity="{ item }">
          {{ item.quantity }} {{ item.quantity_unit }}
        </template>
        <template #item.unit_cost_cents="{ item }">
          {{ formatCurrency(item.unit_cost_cents, item.currency) }}
        </template>
        <template #item.line_total="{ item }">
          {{ formatCurrency(item.quantity * item.unit_cost_cents, item.currency) }}
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
          <div class="text-center py-4 text-medium-emphasis">No line items yet.</div>
        </template>
      </v-data-table>

      <!-- Mobile card view -->
      <div v-else>
        <v-card
          v-for="line in lines"
          :key="line.uuid"
          class="mb-2"
          variant="tonal"
        >
          <v-card-text class="pa-3">
            <div class="d-flex justify-space-between align-start">
              <div>
                <div class="text-body-2 font-weight-medium">
                  #{{ line.line_number }} - {{ line.item_name }}
                </div>
                <div class="text-caption text-medium-emphasis">
                  {{ formatLineItemType(line.item_type) }} | {{ line.quantity }} {{ line.quantity_unit }}
                </div>
                <div class="text-body-2 mt-1">
                  {{ formatCurrency(line.unit_cost_cents, line.currency) }} each |
                  <span class="font-weight-medium">
                    {{ formatCurrency(line.quantity * line.unit_cost_cents, line.currency) }}
                  </span>
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
                  <v-list-item @click="openEditDialog(line)">
                    <template #prepend>
                      <v-icon icon="mdi-pencil" size="small" />
                    </template>
                    <v-list-item-title>Edit</v-list-item-title>
                  </v-list-item>
                  <v-list-item class="text-error" @click="confirmDelete(line)">
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
        <div v-if="lines.length === 0" class="text-center py-4 text-medium-emphasis">
          No line items yet.
        </div>
      </div>
    </v-card-text>
  </v-card>

  <!-- Line Dialog -->
  <PurchaseOrderLineDialog
    v-model="lineDialogOpen"
    :ingredients="ingredients"
    :ingredients-loading="ingredientsLoading"
    :line="editingLine"
    :saving="saving"
    @cancel="closeLineDialog"
    @save="handleSaveLine"
  />

  <!-- Delete Confirmation Dialog -->
  <v-dialog v-model="deleteDialogOpen" max-width="400">
    <v-card>
      <v-card-title class="text-h6">Delete line item?</v-card-title>
      <v-card-text>
        Are you sure you want to delete line #{{ deletingLine?.line_number }} ({{ deletingLine?.item_name }})?
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
  import type { Ingredient, PurchaseOrderLine } from '@/types'
  import { onMounted, ref } from 'vue'
  import { useLineItemTypeFormatters } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import PurchaseOrderLineDialog, { type LineForm } from './PurchaseOrderLineDialog.vue'

  const props = defineProps<{
    purchaseOrderUuid: string
    lines: PurchaseOrderLine[]
    loading?: boolean
  }>()

  const emit = defineEmits<{
    refresh: []
  }>()

  const {
    createPurchaseOrderLine,
    updatePurchaseOrderLine,
    deletePurchaseOrderLine,
    formatCurrency,
  } = useProcurementApi()
  const { getIngredients } = useInventoryApi()
  const { showNotice } = useSnackbar()
  const { formatLineItemType } = useLineItemTypeFormatters()

  // Ingredients for the line item dialog autocomplete
  const ingredients = ref<Ingredient[]>([])
  const ingredientsLoading = ref(false)

  onMounted(async () => {
    ingredientsLoading.value = true
    try {
      ingredients.value = await getIngredients()
    } catch {
      // Non-critical: dialog will show empty list with helpful message
    } finally {
      ingredientsLoading.value = false
    }
  })

  const headers = [
    { title: 'Line #', key: 'line_number', sortable: true, width: '80px' },
    { title: 'Item Name', key: 'item_name', sortable: true },
    { title: 'Type', key: 'item_type', sortable: true },
    { title: 'Qty', key: 'quantity', sortable: true },
    { title: 'Unit Cost', key: 'unit_cost_cents', sortable: true },
    { title: 'Line Total', key: 'line_total', sortable: false },
    { title: '', key: 'actions', sortable: false, align: 'end' as const, width: '60px' },
  ]

  const lineDialogOpen = ref(false)
  const editingLine = ref<PurchaseOrderLine | null>(null)
  const saving = ref(false)

  const deleteDialogOpen = ref(false)
  const deletingLine = ref<PurchaseOrderLine | null>(null)
  const deleting = ref(false)

  function openCreateDialog () {
    editingLine.value = null
    lineDialogOpen.value = true
  }

  function openEditDialog (line: PurchaseOrderLine) {
    editingLine.value = line
    lineDialogOpen.value = true
  }

  function closeLineDialog () {
    lineDialogOpen.value = false
    editingLine.value = null
  }

  function confirmDelete (line: PurchaseOrderLine) {
    deletingLine.value = line
    deleteDialogOpen.value = true
  }

  async function handleSaveLine (form: LineForm) {
    saving.value = true
    try {
      if (editingLine.value) {
        await updatePurchaseOrderLine(editingLine.value.uuid, {
          line_number: form.line_number ?? undefined,
          item_type: form.item_type,
          item_name: form.item_name.trim(),
          inventory_item_uuid: form.inventory_item_uuid || null,
          quantity: form.quantity ?? undefined,
          quantity_unit: form.quantity_unit,
          unit_cost_cents: form.unit_cost_cents ?? undefined,
          currency: form.currency,
        })
        showNotice('Line item updated')
      } else {
        await createPurchaseOrderLine({
          purchase_order_uuid: props.purchaseOrderUuid,
          line_number: form.line_number,
          item_type: form.item_type,
          item_name: form.item_name.trim(),
          inventory_item_uuid: form.inventory_item_uuid || null,
          quantity: form.quantity,
          quantity_unit: form.quantity_unit,
          unit_cost_cents: form.unit_cost_cents,
          currency: form.currency,
        })
        showNotice('Line item added')
      }
      closeLineDialog()
      emit('refresh')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save line item'
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }

  async function handleDelete () {
    if (!deletingLine.value) return

    deleting.value = true
    try {
      await deletePurchaseOrderLine(deletingLine.value.uuid)
      showNotice('Line item deleted')
      deleteDialogOpen.value = false
      deletingLine.value = null
      emit('refresh')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to delete line item'
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
