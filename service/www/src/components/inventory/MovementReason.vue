<template>
  <template v-if="movement.reason === 'use'">
    <router-link
      v-if="usageBatch"
      class="reason-link"
      :to="`/batches/${usageBatch.uuid}`"
    >
      Used in {{ usageBatch.short_name }}
    </router-link>
    <span v-else class="text-medium-emphasis">Used in production</span>
  </template>
  <template v-else-if="movement.reason === 'receive'">
    <span v-if="receiptSupplier">
      Received from {{ receiptSupplier.name }}
    </span>
    <span v-else class="text-medium-emphasis">Received</span>
  </template>
  <template v-else-if="movement.reason === 'adjust' || movement.reason === 'waste'">
    <v-tooltip v-if="adjustmentNotes" location="top">
      <template #activator="{ props }">
        <span v-bind="props" class="adjustment-reason">
          {{ adjustmentReasonLabel }}
          <v-icon class="ml-1" icon="mdi-information-outline" size="x-small" />
        </span>
      </template>
      <span>{{ adjustmentNotes }}</span>
    </v-tooltip>
    <span v-else>{{ adjustmentReasonLabel }}</span>
  </template>
  <template v-else-if="movement.reason === 'transfer'">
    <v-tooltip v-if="transferNotes" location="top">
      <template #activator="{ props }">
        <span v-bind="props" class="transfer-reason">
          {{ transferReasonLabel }}
          <v-icon class="ml-1" icon="mdi-information-outline" size="x-small" />
        </span>
      </template>
      <span>{{ transferNotes }}</span>
    </v-tooltip>
    <span v-else>{{ transferReasonLabel }}</span>
  </template>
  <template v-else>
    <span class="text-medium-emphasis">{{ movement.reason }}</span>
  </template>
</template>

<script lang="ts" setup>
  import type { Batch, InventoryMovement, Supplier } from '@/types'

  defineProps<{
    movement: InventoryMovement
    usageBatch?: Batch
    receiptSupplier?: Supplier
    adjustmentNotes?: string | null
    adjustmentReasonLabel?: string
    transferNotes?: string | null
    transferReasonLabel?: string
  }>()
</script>

<style scoped>
.reason-link {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
}

.reason-link:hover {
  text-decoration: underline;
}

.adjustment-reason,
.transfer-reason {
  cursor: help;
}
</style>
