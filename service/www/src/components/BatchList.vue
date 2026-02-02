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
          :key="batch.id"
          :active="batch.id === selectedBatchId"
          @click="emit('select', batch.id)"
        >
          <v-list-item-title>
            {{ batch.short_name }}
          </v-list-item-title>
          <v-list-item-subtitle>
            #{{ batch.id }} - {{ formatDate(batch.brew_date) }}
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
  export type Batch = {
    id: number
    uuid: string
    short_name: string
    brew_date: string | null
    recipe_id: number | null
    notes: string | null
    created_at: string
    updated_at: string
  }

  withDefaults(
    defineProps<{
      batches: Batch[]
      selectedBatchId: number | null
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
    'select': [batchId: number]
    'create': []
    'bulk-import': []
  }>()

  function formatDate (value: string | null | undefined) {
    if (!value) {
      return 'Unknown'
    }
    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
    }).format(new Date(value))
  }

  function formatDateTime (value: string | null | undefined) {
    if (!value) {
      return 'Unknown'
    }
    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(new Date(value))
  }
</script>

<style scoped>
.batch-list {
  max-height: 60vh;
  overflow-y: auto;
}
</style>
