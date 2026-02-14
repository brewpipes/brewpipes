<template>
  <v-card class="section-card">
    <v-card-title class="d-flex align-center">
      <v-icon class="mr-2" icon="mdi-barley" />
      Batches
      <v-spacer />
      <v-btn
        v-if="showCreateButton"
        aria-label="Create batch"
        icon="mdi-plus"
        size="small"
        variant="text"
        @click="emit('create')"
      />
    </v-card-title>
    <v-card-text>
      <v-progress-linear v-if="loading" class="mb-3" color="primary" indeterminate />

      <v-list active-color="primary" class="batch-list" lines="two">
        <v-list-item
          v-for="batch in batches"
          :key="batch.uuid"
          :active="batch.uuid === selectedBatchUuid"
          @click="emit('select', batch.uuid)"
        >
          <v-list-item-title>
            {{ batch.short_name }}
          </v-list-item-title>
          <v-list-item-subtitle>
            {{ formatDate(batch.brew_date) }}
          </v-list-item-subtitle>
          <template #append>
            <v-chip size="x-small" variant="tonal">
              {{ formatDateTime(batch.updated_at) }}
            </v-chip>
          </template>
        </v-list-item>

        <v-list-item v-if="!loading && batches.length === 0">
          <v-list-item-title>No batches yet</v-list-item-title>
          <v-list-item-subtitle>Use + to add the first batch.</v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card-text>
    <v-card-actions v-if="showBulkImport" class="justify-end">
      <v-btn
        prepend-icon="mdi-upload"
        size="small"
        variant="text"
        @click="emit('bulk-import')"
      >
        Bulk import
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts" setup>
  import type { Batch } from '@/types'
  import { formatDate, formatDateTime } from '@/composables/useFormatters'

  withDefaults(
    defineProps<{
      batches: Batch[]
      selectedBatchUuid: string | null
      loading?: boolean
      showCreateButton?: boolean
      showBulkImport?: boolean
    }>(),
    {
      loading: false,
      showCreateButton: true,
      showBulkImport: true,
    },
  )

  const emit = defineEmits<{
    'select': [batchUuid: string]
    'create': []
    'bulk-import': []
  }>()

</script>

<style scoped>
.batch-list {
  max-height: 60vh;
  overflow-y: auto;
}
</style>
