<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 800"
    :model-value="modelValue"
    persistent
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <span class="text-h6">Receive Shipment</span>
        <v-spacer />
        <v-btn
          :disabled="saving"
          icon="mdi-close"
          size="small"
          variant="text"
          @click="handleClose"
        />
      </v-card-title>

      <v-divider />

      <v-card-text class="pa-0" style="overflow-y: auto;">
        <!-- Loading state for reference data -->
        <template v-if="loadingReferenceData">
          <v-container class="pa-6">
            <div class="d-flex flex-column align-center">
              <v-progress-circular color="primary" indeterminate size="48" />
              <p class="text-body-2 text-medium-emphasis mt-4">Loading reference data...</p>
            </div>
          </v-container>
        </template>

        <!-- Reference data load warning -->
        <v-alert
          v-else-if="referenceDataWarning"
          class="ma-4"
          closable
          density="compact"
          type="warning"
          variant="tonal"
          @click:close="referenceDataWarning = ''"
        >
          {{ referenceDataWarning }}
        </v-alert>

        <!-- Stepper -->
        <v-stepper
          v-if="!loadingReferenceData"
          v-model="currentStep"
          alt-labels
          flat
          hide-actions
          :items="stepItems"
        >
          <template #item.1>
            <ReceiveShipmentStep1
              v-model:selected-line-uuids="selectedLineUuids"
              :lines="lines"
              :previously-received="previouslyReceived"
            />
          </template>

          <template #item.2>
            <ReceiveShipmentStep2
              v-model:line-details="lineDetails"
              :ingredients="ingredients"
              :lines="selectedLines"
              :previously-received="previouslyReceived"
              :stock-locations="stockLocations"
            />
          </template>

          <template #item.3>
            <ReceiveShipmentStep3
              v-model:notes="receiptNotes"
              :error="saveError"
              :ingredients="ingredients"
              :line-details="lineDetails"
              :lines="selectedLines"
              :new-status="computedNewStatus"
              :stock-locations="stockLocations"
              :supplier-name="supplierName"
            />
          </template>
        </v-stepper>
      </v-card-text>

      <v-divider />

      <v-card-actions class="justify-space-between pa-4">
        <v-btn
          :disabled="currentStep === 1 || saving || loadingReferenceData"
          variant="text"
          @click="previousStep"
        >
          Back
        </v-btn>
        <div>
          <v-btn
            :disabled="saving"
            variant="text"
            @click="handleClose"
          >
            Cancel
          </v-btn>
          <v-btn
            v-if="currentStep < 3"
            color="primary"
            :disabled="!canProceed || loadingReferenceData"
            @click="nextStep"
          >
            Next
          </v-btn>
          <v-btn
            v-else
            color="success"
            :disabled="!canProceed"
            :loading="saving"
            @click="handleConfirm"
          >
            Confirm Receipt
          </v-btn>
        </div>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type {
    CreateIngredientLotRequest,
    CreateInventoryMovementRequest,
    CreateInventoryReceiptRequest,
    Ingredient,
    IngredientLot,
    LineReceivingDetails,
    PurchaseOrder,
    PurchaseOrderLine,
    StockLocation,
  } from '@/types'
  import { computed, ref, watch } from 'vue'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProcurementApi } from '@/composables/useProcurementApi'
  import ReceiveShipmentStep1 from './ReceiveShipmentStep1.vue'
  import ReceiveShipmentStep2 from './ReceiveShipmentStep2.vue'
  import ReceiveShipmentStep3 from './ReceiveShipmentStep3.vue'

  const props = defineProps<{
    modelValue: boolean
    purchaseOrder: PurchaseOrder
    lines: PurchaseOrderLine[]
    supplierName: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'received': []
  }>()

  const {
    getIngredients,
    getIngredientLots,
    getStockLocations,
    createInventoryReceipt,
    createIngredientLot,
    createInventoryMovement,
  } = useInventoryApi()

  const { updatePurchaseOrder } = useProcurementApi()

  // Wizard state
  const currentStep = ref(1)
  const stepItems = ['Select Lines', 'Enter Details', 'Review & Confirm']

  // Step 1: Line selection
  const selectedLineUuids = ref<string[]>([])

  // Step 2: Line details
  const lineDetails = ref<LineReceivingDetails[]>([])

  // Step 3: Notes
  const receiptNotes = ref('')

  // Reference data
  const ingredients = ref<Ingredient[]>([])
  const stockLocations = ref<StockLocation[]>([])
  const previouslyReceived = ref<Map<string, number>>(new Map())
  const loadingReferenceData = ref(false)
  const referenceDataWarning = ref('')

  // Saving state
  const saving = ref(false)
  const saveError = ref('')

  // Computed
  const selectedLines = computed(() =>
    props.lines.filter(line => selectedLineUuids.value.includes(line.uuid)),
  )

  const canProceed = computed(() => {
    if (currentStep.value === 1) {
      return selectedLineUuids.value.length > 0
    }
    if (currentStep.value === 2) {
      return lineDetails.value.every(
        detail => detail.quantity > 0 && detail.locationUuid,
      )
    }
    return true
  })

  const computedNewStatus = computed(() => {
    // Check if all lines will be fully received after this receipt
    const allLinesFullyReceived = props.lines.every(line => {
      const prevReceived = previouslyReceived.value.get(line.uuid) ?? 0
      const currentReceiving = lineDetails.value.find(d => d.lineUuid === line.uuid)?.quantity ?? 0
      const totalReceived = prevReceived + currentReceiving
      return totalReceived >= line.quantity
    })
    return allLinesFullyReceived ? 'received' : 'partially_received'
  })

  // Watch for dialog open to load reference data
  watch(() => props.modelValue, async open => {
    if (open) {
      resetWizard()
      await loadReferenceData()
    }
  })

  // Watch selected lines to initialize line details
  watch(selectedLineUuids, newUuids => {
    // Add new lines
    for (const uuid of newUuids) {
      if (!lineDetails.value.some(d => d.lineUuid === uuid)) {
        const line = props.lines.find(l => l.uuid === uuid)
        if (line) {
          const prevReceived = previouslyReceived.value.get(uuid) ?? 0
          const remaining = Math.max(0, line.quantity - prevReceived)
          lineDetails.value.push({
            lineUuid: uuid,
            quantity: remaining,
            unit: line.quantity_unit,
            locationUuid: '',
            breweryLotCode: '',
            supplierLotCode: '',
          })
        }
      }
    }
    // Remove deselected lines
    lineDetails.value = lineDetails.value.filter(d =>
      newUuids.includes(d.lineUuid),
    )
  })

  async function loadReferenceData () {
    loadingReferenceData.value = true
    referenceDataWarning.value = ''

    const warnings: string[] = []

    try {
      const [ingredientsResult, locationsResult] = await Promise.allSettled([
        getIngredients(),
        getStockLocations(),
      ])

      if (ingredientsResult.status === 'fulfilled') {
        ingredients.value = ingredientsResult.value
      } else {
        warnings.push('ingredients')
      }

      if (locationsResult.status === 'fulfilled') {
        stockLocations.value = locationsResult.value
      } else {
        warnings.push('stock locations')
      }

      // Load previously received quantities for each line
      await loadPreviouslyReceived()

      if (warnings.length > 0) {
        referenceDataWarning.value = `Failed to load ${warnings.join(' and ')}. Some options may be unavailable.`
      }
    } catch {
      referenceDataWarning.value = 'Failed to load reference data. Some options may be unavailable.'
    } finally {
      loadingReferenceData.value = false
    }
  }

  async function loadPreviouslyReceived () {
    // Parallelize fetching previously received quantities for all lines
    const linesWithInventory = props.lines.filter(line => line.inventory_item_uuid)

    const results = await Promise.all(
      linesWithInventory.map(async line => {
        try {
          const lots = await getIngredientLots({ purchase_order_line_uuid: line.uuid })
          const total = lots.reduce((sum: number, lot: IngredientLot) => sum + lot.received_amount, 0)
          return { lineUuid: line.uuid, received: total }
        } catch {
          return { lineUuid: line.uuid, received: 0 }
        }
      }),
    )

    const received = new Map<string, number>()

    // Add results from parallel fetch
    for (const result of results) {
      received.set(result.lineUuid, result.received)
    }

    // Add lines without inventory_item_uuid with 0 received
    for (const line of props.lines) {
      if (!line.inventory_item_uuid) {
        received.set(line.uuid, 0)
      }
    }

    previouslyReceived.value = received
  }

  function resetWizard () {
    currentStep.value = 1
    selectedLineUuids.value = []
    lineDetails.value = []
    receiptNotes.value = ''
    saveError.value = ''
    saving.value = false
  }

  function nextStep () {
    if (currentStep.value < 3 && canProceed.value) {
      currentStep.value++
    }
  }

  function previousStep () {
    if (currentStep.value > 1) {
      currentStep.value--
    }
  }

  function handleClose () {
    emit('update:modelValue', false)
  }

  async function handleConfirm () {
    saving.value = true
    saveError.value = ''

    try {
      // 1. Create inventory receipt
      const receiptPayload: CreateInventoryReceiptRequest = {
        supplier_uuid: props.purchaseOrder.supplier_uuid,
        purchase_order_uuid: props.purchaseOrder.uuid,
        received_at: new Date().toISOString(),
        notes: receiptNotes.value.trim() || null,
      }
      const receipt = await createInventoryReceipt(receiptPayload)

      // 2. For each selected line, create lot and movement
      // Track successes and failures for partial completion reporting
      const createdLots: string[] = []
      const failedLines: string[] = []

      for (const detail of lineDetails.value) {
        const line = props.lines.find(l => l.uuid === detail.lineUuid)
        if (!line || !line.inventory_item_uuid) continue

        try {
          const lotPayload: CreateIngredientLotRequest = {
            ingredient_uuid: line.inventory_item_uuid,
            receipt_uuid: receipt.uuid,
            purchase_order_line_uuid: line.uuid,
            supplier_uuid: props.purchaseOrder.supplier_uuid,
            brewery_lot_code: detail.breweryLotCode.trim() || null,
            originator_lot_code: detail.supplierLotCode.trim() || null,
            received_at: new Date().toISOString(),
            received_amount: detail.quantity,
            received_unit: detail.unit,
          }
          const lot = await createIngredientLot(lotPayload)
          createdLots.push(lot.uuid)

          const movementPayload: CreateInventoryMovementRequest = {
            ingredient_lot_uuid: lot.uuid,
            stock_location_uuid: detail.locationUuid,
            direction: 'in',
            reason: 'receive',
            amount: detail.quantity,
            amount_unit: detail.unit,
            receipt_uuid: receipt.uuid,
          }
          await createInventoryMovement(movementPayload)
        } catch {
          failedLines.push(line.item_name)
        }
      }

      // Check for partial failures
      if (failedLines.length > 0) {
        if (createdLots.length > 0) {
          saveError.value = `Partially completed. Created ${createdLots.length} lot(s). Failed: ${failedLines.join(', ')}. Please review inventory and retry failed items.`
        } else {
          saveError.value = `Failed to create lots for: ${failedLines.join(', ')}. Receipt was created but no inventory was added.`
        }
        // Don't close dialog on partial failure - let user see the error
        return
      }

      // 3. Update PO status (only if all lots succeeded)
      await updatePurchaseOrder(props.purchaseOrder.uuid, {
        status: computedNewStatus.value,
      })

      emit('received')
      emit('update:modelValue', false)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save receipt'
      saveError.value = message
    } finally {
      saving.value = false
    }
  }
</script>
