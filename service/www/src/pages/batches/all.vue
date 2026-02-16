<template>
  <v-container class="batches-page" fluid>
    <v-card class="section-card">
      <v-card-title class="card-title-responsive">
        <div class="d-flex align-center">
          <v-icon class="mr-2" icon="mdi-barley" />
          <span class="d-none d-sm-inline">All Batches</span>
          <span class="d-sm-none">Batches</span>
        </div>
        <div class="card-title-actions">
          <v-text-field
            v-model="search"
            append-inner-icon="mdi-magnify"
            class="search-field"
            clearable
            density="compact"
            hide-details
            label="Search"
            single-line
            variant="outlined"
          />
          <v-btn
            :icon="$vuetify.display.xs"
            :loading="loading"
            size="small"
            variant="text"
            @click="loadBatches"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-refresh" />
            <span v-else>Refresh</span>
          </v-btn>
          <v-btn
            :icon="$vuetify.display.xs"
            :prepend-icon="$vuetify.display.xs ? undefined : 'mdi-upload'"
            size="small"
            variant="text"
            @click="bulkImportDialog = true"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-upload" />
            <span v-else>Import</span>
          </v-btn>
          <v-btn
            color="primary"
            :icon="$vuetify.display.xs"
            size="small"
            variant="text"
            @click="openCreateDialog"
          >
            <v-icon v-if="$vuetify.display.xs" icon="mdi-plus" />
            <span v-else>New batch</span>
          </v-btn>
        </div>
      </v-card-title>
      <v-card-text>
        <v-alert
          v-if="errorMessage"
          class="mb-3"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ errorMessage }}
        </v-alert>

        <v-data-table
          class="data-table batches-table"
          density="compact"
          :headers="headers"
          item-value="uuid"
          :items="sortedBatches"
          :loading="loading"
          :search="search"
          @click:row="onRowClick"
        >
          <template #item.short_name="{ item }">
            <span class="font-weight-medium">{{ item.short_name }}</span>
          </template>

          <template #item.recipe_name="{ item }">
            <span v-if="item.recipe_name">{{ item.recipe_name }}</span>
            <span v-else class="text-medium-emphasis">-</span>
          </template>

          <template #item.current_phase="{ item }">
            <v-chip
              v-if="item.current_phase"
              :color="getPhaseColor(item.current_phase)"
              size="small"
              variant="tonal"
            >
              {{ formatPhase(item.current_phase) }}
            </v-chip>
            <span v-else class="text-medium-emphasis">-</span>
          </template>

          <template #item.brew_date="{ item }">
            <span v-if="item.brew_date">{{ formatDate(item.brew_date) }}</span>
            <span v-else class="text-medium-emphasis">-</span>
          </template>

          <template #item.updated_at="{ item }">
            {{ formatDateTime(item.updated_at) }}
          </template>

          <template #item.actions="{ item }">
            <v-btn
              icon="mdi-pencil"
              size="x-small"
              variant="text"
              @click.stop="openEditDialog(item)"
            />
            <v-btn
              color="error"
              icon="mdi-delete"
              size="x-small"
              variant="text"
              @click.stop="openDeleteDialog(item)"
            />
          </template>

          <template #no-data>
            <div class="text-center py-4">
              <div class="text-body-2 text-medium-emphasis">No batches yet.</div>
              <v-btn
                class="mt-2"
                color="primary"
                size="small"
                variant="text"
                @click="openCreateDialog"
              >
                Create your first batch
              </v-btn>
            </div>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </v-container>

  <!-- Create Batch Dialog -->
  <BatchCreateDialog
    v-model="createBatchDialog"
    :recipes="recipes"
    :recipes-loading="recipesLoading"
    :saving="saving"
    @submit="handleCreateBatch"
  />

  <!-- Bulk Import Dialog -->
  <BulkImportDialog
    ref="bulkImportDialogRef"
    v-model="bulkImportDialog"
    :saving="importUploading"
    @submit="handleBulkImport"
  />

  <!-- Edit Batch Dialog -->
  <BatchEditDialog
    v-model="editBatchDialog"
    :batch="editingBatch"
    :error-message="editBatchError"
    :recipes="recipes"
    :recipes-loading="recipesLoading"
    :saving="savingBatchEdit"
    @submit="saveBatchEdit"
  />

  <!-- Delete Batch Dialog -->
  <BatchDeleteDialog
    v-model="deleteBatchDialog"
    :batch="deletingBatchItem"
    :deleting="deletingBatch"
    :error-message="deleteBatchError"
    @confirm="confirmDeleteBatch"
  />
</template>

<script lang="ts" setup>
  import type { Batch, Recipe, UpdateBatchRequest } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { BatchCreateDialog, type BatchCreateForm, BatchDeleteDialog, BatchEditDialog, type BatchEditForm } from '@/components/batch'
  import BulkImportDialog from '@/components/batch/BulkImportDialog.vue'
  import type { BatchImportResponse } from '@/components/batch/BulkImportDialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDate, formatDateTime, usePhaseFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { normalizeDateOnly, normalizeText } from '@/utils/normalize'

  const router = useRouter()
  const { getBatches, createBatch: createBatchApi, getRecipes, updateBatch, deleteBatch, request } = useProductionApi()
  const { showNotice } = useSnackbar()
  const { formatPhase, getPhaseColor } = usePhaseFormatters()

  // State
  const batches = ref<Batch[]>([])
  const recipes = ref<Recipe[]>([])
  const search = ref('')

  const { execute: executeLoad, loading, error: errorMessage } = useAsyncAction()
  const { execute: executeSave, loading: saving, error: saveError } = useAsyncAction()
  const { execute: executeEdit, loading: savingBatchEdit, error: editBatchError } = useAsyncAction()
  const { execute: executeDelete, loading: deletingBatch, error: deleteBatchError } = useAsyncAction()
  const { execute: executeImport, loading: importUploading, error: importError } = useAsyncAction()
  const { execute: executeLoadRecipes, loading: recipesLoading } = useAsyncAction({ onError: () => {} })

  // Dialogs
  const createBatchDialog = ref(false)
  const bulkImportDialog = ref(false)
  const bulkImportDialogRef = ref<InstanceType<typeof BulkImportDialog> | null>(null)

  // Edit/Delete state
  const editBatchDialog = ref(false)
  const editingBatch = ref<Batch | null>(null)

  const deleteBatchDialog = ref(false)
  const deletingBatchItem = ref<Batch | null>(null)

  // Table configuration
  const headers = [
    { title: 'Short Name', key: 'short_name', sortable: true },
    { title: 'Recipe', key: 'recipe_name', sortable: true },
    { title: 'Status', key: 'current_phase', sortable: true },
    { title: 'Brew Date', key: 'brew_date', sortable: true },
    { title: 'Updated', key: 'updated_at', sortable: true },
    { title: '', key: 'actions', sortable: false, align: 'end' as const, width: '100px' },
  ]

  /**
   * Sort batches by:
   * 1. Upcoming batches (future brew_date) first, sorted by nearest date
   * 2. In-progress batches (non-finished, non-planning phases)
   * 3. Completed batches (finished phase)
   * 4. Within each group, sort by updated_at descending (most recent first)
   */
  const sortedBatches = computed(() => {
    const now = new Date()

    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...batches.value].sort((a, b) => {
      const aBrewDate = a.brew_date ? new Date(a.brew_date) : null
      const bBrewDate = b.brew_date ? new Date(b.brew_date) : null

      const aIsUpcoming = aBrewDate && aBrewDate > now
      const bIsUpcoming = bBrewDate && bBrewDate > now

      const aIsFinished = a.current_phase === 'finished'
      const bIsFinished = b.current_phase === 'finished'

      const aIsPlanning = a.current_phase === 'planning' || !a.current_phase
      const bIsPlanning = b.current_phase === 'planning' || !b.current_phase

      // Upcoming batches first
      if (aIsUpcoming && !bIsUpcoming) return -1
      if (!aIsUpcoming && bIsUpcoming) return 1

      // If both are upcoming, sort by nearest brew date
      if (aIsUpcoming && bIsUpcoming && aBrewDate && bBrewDate) {
        return aBrewDate.getTime() - bBrewDate.getTime()
      }

      // Finished batches last
      if (aIsFinished && !bIsFinished) return 1
      if (!aIsFinished && bIsFinished) return -1

      // Planning batches after in-progress but before finished
      if (aIsPlanning && !bIsPlanning && !bIsFinished) return 1
      if (!aIsPlanning && !aIsFinished && bIsPlanning) return -1

      // Within same group, sort by updated_at descending
      const aUpdated = new Date(a.updated_at).getTime()
      const bUpdated = new Date(b.updated_at).getTime()
      return bUpdated - aUpdated
    })
  })

  // Lifecycle
  onMounted(async () => {
    await Promise.all([loadBatches(), loadRecipes()])
  })

  // Methods
  async function loadBatches () {
    await executeLoad(async () => {
      batches.value = await getBatches()
    })
    if (errorMessage.value) {
      showNotice(errorMessage.value, 'error')
    }
  }

  async function loadRecipes () {
    // Recipe loading failure is non-critical
    await executeLoadRecipes(async () => {
      recipes.value = await getRecipes()
    })
  }

  function openCreateDialog () {
    createBatchDialog.value = true
  }

  async function handleCreateBatch (form: BatchCreateForm) {
    await executeSave(async () => {
      const payload = {
        short_name: form.short_name.trim(),
        brew_date: normalizeDateOnly(form.brew_date),
        recipe_uuid: form.recipe_uuid,
        notes: normalizeText(form.notes),
      }

      await createBatchApi(payload)

      showNotice('Batch created')
      createBatchDialog.value = false
      await loadBatches()
    })
    if (saveError.value) {
      errorMessage.value = saveError.value
      showNotice(saveError.value, 'error')
    }
  }

  function onRowClick (_event: Event, { item }: { item: Batch }) {
    router.push(`/batches/${item.uuid}`)
  }

  // Edit/Delete functions
  function openEditDialog (batch: Batch) {
    editingBatch.value = batch
    editBatchError.value = ''
    editBatchDialog.value = true
  }

  async function saveBatchEdit (form: BatchEditForm) {
    if (!editingBatch.value) return

    await executeEdit(async () => {
      const payload: UpdateBatchRequest = {
        short_name: form.short_name.trim(),
        brew_date: form.brew_date ? normalizeDateOnly(form.brew_date) : null,
        recipe_uuid: form.recipe_uuid,
        notes: normalizeText(form.notes),
      }

      await updateBatch(editingBatch.value!.uuid, payload)
      showNotice('Batch updated')
      editBatchDialog.value = false
      editingBatch.value = null
      await loadBatches()
    })
  }

  function openDeleteDialog (batch: Batch) {
    deletingBatchItem.value = batch
    deleteBatchError.value = ''
    deleteBatchDialog.value = true
  }

  async function confirmDeleteBatch () {
    if (!deletingBatchItem.value) return

    await executeDelete(async () => {
      await deleteBatch(deletingBatchItem.value!.uuid)
      showNotice('Batch deleted')
      deleteBatchDialog.value = false
      deletingBatchItem.value = null
      await loadBatches()
    })
  }

  // Import functions
  async function handleBulkImport (file: File) {
    errorMessage.value = ''
    await executeImport(async () => {
      const formData = new FormData()
      formData.append('file', file)
      const response = await request<BatchImportResponse>('/batches/import', {
        method: 'POST',
        body: formData,
        headers: new Headers(),
      })
      bulkImportDialogRef.value?.setImportResult(response)
      const successCount = response.totals?.created ?? 0
      const failureCount = response.totals?.failed ?? 0
      if (failureCount > 0) {
        const color = successCount > 0 ? 'warning' : 'error'
        showNotice(`Imported ${successCount} batches, ${failureCount} failed`, color)
      } else {
        showNotice(`Imported ${successCount} ${successCount === 1 ? 'batch' : 'batches'}`)
      }
      await loadBatches()
    })
    if (importError.value) {
      errorMessage.value = importError.value
      showNotice(importError.value, 'error')
    }
  }

</script>

<style scoped>
.batches-page {
  position: relative;
}

.batches-table :deep(tr) {
  cursor: pointer;
}

.batches-table :deep(tr:hover td) {
  background: rgba(var(--v-theme-primary), 0.04);
}

.batches-table :deep(.v-table__wrapper) {
  overflow-x: auto;
}
</style>
