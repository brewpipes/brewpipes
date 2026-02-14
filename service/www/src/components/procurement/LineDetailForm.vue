<template>
  <div>
    <!-- Info row -->
    <div class="d-flex flex-wrap ga-4 mb-3 text-body-2 text-medium-emphasis">
      <span><strong>Ordered:</strong> {{ line.quantity }} {{ line.quantity_unit }}</span>
      <span><strong>Previously received:</strong> {{ previouslyReceived }} {{ line.quantity_unit }}</span>
      <span :class="remaining === 0 ? 'text-success' : 'text-warning'">
        <strong>Remaining:</strong> {{ remaining }} {{ line.quantity_unit }}
      </span>
    </div>

    <!-- Form fields -->
    <v-row dense>
      <v-col cols="6" sm="4">
        <v-text-field
          density="comfortable"
          inputmode="decimal"
          label="Quantity *"
          min="0"
          :model-value="detail?.quantity ?? 0"
          :rules="[v => v > 0 || 'Required']"
          step="any"
          type="number"
          @update:model-value="update('quantity', Number($event))"
        />
      </v-col>
      <v-col cols="6" sm="4">
        <v-text-field
          density="comfortable"
          label="Unit"
          :model-value="detail?.unit ?? line.quantity_unit"
          readonly
          @update:model-value="update('unit', String($event))"
        />
      </v-col>
      <v-col cols="12" sm="4">
        <v-select
          density="comfortable"
          :items="locationItems"
          label="Storage location *"
          :model-value="detail?.locationUuid ?? ''"
          :rules="[v => !!v || 'Required']"
          @update:model-value="update('locationUuid', String($event))"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <v-text-field
          density="comfortable"
          hint="Optional internal lot code"
          label="Brewery lot code"
          :model-value="detail?.breweryLotCode ?? ''"
          persistent-hint
          @update:model-value="update('breweryLotCode', String($event))"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <v-text-field
          density="comfortable"
          hint="Optional supplier/manufacturer lot code"
          label="Supplier lot code"
          :model-value="detail?.supplierLotCode ?? ''"
          persistent-hint
          @update:model-value="update('supplierLotCode', String($event))"
        />
      </v-col>
    </v-row>

    <!-- Warning for over-receiving -->
    <v-alert
      v-if="showOverReceiveWarning"
      class="mt-2"
      density="compact"
      type="warning"
      variant="tonal"
    >
      Quantity exceeds remaining amount on order
    </v-alert>
  </div>
</template>

<script lang="ts" setup>
  import type { LineReceivingDetails, PurchaseOrderLine, StockLocation } from '@/types'
  import { computed } from 'vue'

  const props = defineProps<{
    line: PurchaseOrderLine
    detail: LineReceivingDetails | undefined
    stockLocations: StockLocation[]
    previouslyReceived: number
    ingredientName: string
  }>()

  const emit = defineEmits<{
    'update:detail': [value: Partial<LineReceivingDetails>]
  }>()

  const remaining = computed(() =>
    Math.max(0, props.line.quantity - props.previouslyReceived),
  )

  const showOverReceiveWarning = computed(() =>
    (props.detail?.quantity ?? 0) > remaining.value,
  )

  const locationItems = computed(() =>
    props.stockLocations.map(loc => ({
      title: loc.name,
      value: loc.uuid,
    })),
  )

  function update (field: keyof LineReceivingDetails, value: string | number) {
    emit('update:detail', { [field]: value })
  }
</script>
