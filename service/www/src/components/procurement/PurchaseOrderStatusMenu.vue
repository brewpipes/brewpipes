<template>
  <v-menu v-model="menuOpen" :close-on-content-click="false">
    <template #activator="{ props: menuProps }">
      <v-chip
        v-bind="menuProps"
        :color="statusColor"
        :disabled="disabled"
        size="small"
        variant="flat"
      >
        {{ formatStatus(status) }}
        <v-icon v-if="!disabled" end icon="mdi-chevron-down" size="small" />
      </v-chip>
    </template>
    <v-list density="compact">
      <v-list-subheader>Change status</v-list-subheader>
      <v-list-item
        v-for="transition in validTransitions"
        :key="transition"
        :disabled="saving"
        @click="handleStatusChange(transition)"
      >
        <template #prepend>
          <v-icon :color="getStatusColor(transition)" icon="mdi-circle" size="x-small" />
        </template>
        <v-list-item-title>{{ formatStatus(transition) }}</v-list-item-title>
      </v-list-item>
      <v-list-item v-if="validTransitions.length === 0" disabled>
        <v-list-item-title class="text-medium-emphasis">No transitions available</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue'

  type PurchaseOrderStatus = 'draft' | 'submitted' | 'confirmed' | 'partially_received' | 'received' | 'cancelled'

  const props = defineProps<{
    status: string
    disabled?: boolean
    saving?: boolean
  }>()

  const emit = defineEmits<{
    'update:status': [status: string]
  }>()

  const menuOpen = ref(false)

  // Status transition rules
  const STATUS_TRANSITIONS: Record<PurchaseOrderStatus, PurchaseOrderStatus[]> = {
    draft: ['submitted', 'cancelled'],
    submitted: ['confirmed', 'cancelled'],
    confirmed: ['partially_received', 'received', 'cancelled'],
    partially_received: ['received', 'cancelled'],
    received: [],
    cancelled: [],
  }

  const STATUS_COLORS: Record<PurchaseOrderStatus, string> = {
    draft: 'grey',
    submitted: 'blue',
    confirmed: 'green',
    partially_received: 'orange',
    received: 'green',
    cancelled: 'red',
  }

  const STATUS_LABELS: Record<PurchaseOrderStatus, string> = {
    draft: 'Draft',
    submitted: 'Submitted',
    confirmed: 'Confirmed',
    partially_received: 'Partially Received',
    received: 'Received',
    cancelled: 'Cancelled',
  }

  const statusColor = computed(() => getStatusColor(props.status))

  const validTransitions = computed(() => {
    const currentStatus = props.status as PurchaseOrderStatus
    return STATUS_TRANSITIONS[currentStatus] ?? []
  })

  function getStatusColor (status: string): string {
    return STATUS_COLORS[status as PurchaseOrderStatus] ?? 'grey'
  }

  function formatStatus (status: string): string {
    return STATUS_LABELS[status as PurchaseOrderStatus] ?? status
  }

  function handleStatusChange (newStatus: string) {
    menuOpen.value = false
    emit('update:status', newStatus)
  }
</script>
