<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 480"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Delete batch</v-card-title>
      <v-card-text>
        <v-alert
          v-if="errorMessage"
          class="mb-4"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ errorMessage }}
        </v-alert>

        <p class="text-body-1 mb-4">
          Are you sure you want to delete
          <strong>{{ batchName }}</strong>?
        </p>

        <v-alert
          density="compact"
          type="warning"
          variant="tonal"
        >
          This action cannot be undone. The batch and its associated data will be permanently removed.
        </v-alert>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="deleting" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="error"
          :loading="deleting"
          variant="flat"
          @click="handleConfirm"
        >
          Delete batch
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import type { Batch } from '@/types'

  const props = defineProps<{
    modelValue: boolean
    batch: Batch | null
    deleting: boolean
    errorMessage: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'confirm': []
  }>()

  const batchName = computed(() => props.batch?.short_name ?? 'this batch')

  function handleConfirm () {
    emit('confirm')
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
