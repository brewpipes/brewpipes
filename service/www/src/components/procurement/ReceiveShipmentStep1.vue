<template>
  <v-container class="pa-4" fluid>
    <v-row>
      <v-col cols="12">
        <p class="text-body-2 text-medium-emphasis mb-4">
          Select the line items you are receiving in this shipment.
        </p>
      </v-col>
    </v-row>

    <!-- Select All -->
    <v-row>
      <v-col cols="12">
        <v-checkbox
          v-model="selectAll"
          density="compact"
          hide-details
          :indeterminate="isIndeterminate"
          label="Select all lines"
        />
      </v-col>
    </v-row>

    <v-divider class="my-3" />

    <!-- Line Items -->
    <v-row v-for="line in lines" :key="line.uuid" class="line-row">
      <v-col cols="12">
        <v-card
          :class="{ 'border-primary': isSelected(line.uuid) }"
          variant="outlined"
          @click="toggleLine(line.uuid)"
        >
          <v-card-text class="d-flex align-center pa-3">
            <v-checkbox
              density="compact"
              hide-details
              :model-value="isSelected(line.uuid)"
              @click.stop
              @update:model-value="toggleLine(line.uuid)"
            />
            <div class="flex-grow-1 ml-2">
              <div class="d-flex flex-wrap align-center ga-2">
                <span class="font-weight-medium">{{ line.item_name }}</span>
                <v-chip density="compact" size="small" variant="tonal">
                  {{ line.item_type }}
                </v-chip>
              </div>
              <div class="d-flex flex-wrap ga-4 mt-1 text-body-2 text-medium-emphasis">
                <span>
                  <strong>Ordered:</strong> {{ line.quantity }} {{ line.quantity_unit }}
                </span>
                <span>
                  <strong>Received:</strong> {{ getPreviouslyReceived(line.uuid) }} {{ line.quantity_unit }}
                </span>
                <span :class="{ 'text-success': getRemaining(line) === 0, 'text-warning': getRemaining(line) > 0 }">
                  <strong>Remaining:</strong> {{ getRemaining(line) }} {{ line.quantity_unit }}
                </span>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row v-if="lines.length === 0">
      <v-col cols="12">
        <v-alert density="compact" type="info" variant="tonal">
          No line items on this purchase order.
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
  import type { PurchaseOrderLine } from '@/types'
  import { computed } from 'vue'

  const props = defineProps<{
    lines: PurchaseOrderLine[]
    selectedLineUuids: string[]
    previouslyReceived: Map<string, number>
  }>()

  const emit = defineEmits<{
    'update:selectedLineUuids': [value: string[]]
  }>()

  const selectAll = computed({
    get: () => props.selectedLineUuids.length === props.lines.length && props.lines.length > 0,
    set: (value: boolean) => {
      if (value) {
        emit('update:selectedLineUuids', props.lines.map(l => l.uuid))
      } else {
        emit('update:selectedLineUuids', [])
      }
    },
  })

  const isIndeterminate = computed(() =>
    props.selectedLineUuids.length > 0 && props.selectedLineUuids.length < props.lines.length,
  )

  function isSelected (uuid: string): boolean {
    return props.selectedLineUuids.includes(uuid)
  }

  function toggleLine (uuid: string) {
    const newSelection = isSelected(uuid)
      ? props.selectedLineUuids.filter(id => id !== uuid)
      : [...props.selectedLineUuids, uuid]
    emit('update:selectedLineUuids', newSelection)
  }

  function getPreviouslyReceived (lineUuid: string): number {
    return props.previouslyReceived.get(lineUuid) ?? 0
  }

  function getRemaining (line: PurchaseOrderLine): number {
    const received = getPreviouslyReceived(line.uuid)
    return Math.max(0, line.quantity - received)
  }
</script>

<style scoped>
.line-row .v-card {
  cursor: pointer;
  transition: border-color 0.2s;
}

.line-row .v-card:hover {
  border-color: rgb(var(--v-theme-primary));
}

.border-primary {
  border-color: rgb(var(--v-theme-primary)) !important;
  border-width: 2px;
}
</style>
