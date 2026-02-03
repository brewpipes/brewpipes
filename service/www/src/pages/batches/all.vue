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
          item-value="id"
          :items="sortedBatches"
          :loading="loading"
          :search="search"
          @dblclick:row="onRowDoubleClick"
        >
          <template #item.id="{ item }">
            <span class="text-medium-emphasis">#{{ item.id }}</span>
          </template>

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

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <!-- Create Batch Dialog -->
  <v-dialog v-model="createBatchDialog" :max-width="$vuetify.display.xs ? '100%' : 520" persistent>
    <v-card>
      <v-card-title class="text-h6">Create batch</v-card-title>
      <v-card-text>
        <v-form ref="formRef" @submit.prevent="createBatch">
          <v-text-field
            v-model="batchForm.short_name"
            density="comfortable"
            label="Short name"
            placeholder="IPA 24-07"
            :rules="[rules.required]"
          />
          <v-text-field
            v-model="batchForm.brew_date"
            density="comfortable"
            label="Brew date"
            type="date"
          />
          <v-autocomplete
            v-model="batchForm.recipe_id"
            clearable
            density="comfortable"
            hint="Optional - link this batch to a recipe"
            item-title="title"
            item-value="value"
            :items="recipeSelectItems"
            label="Recipe"
            :loading="recipesLoading"
            persistent-hint
          >
            <template #item="{ props, item }">
              <v-list-item v-bind="props">
                <template #subtitle>
                  <span v-if="item.raw.style">{{ item.raw.style }}</span>
                </template>
              </v-list-item>
            </template>
          </v-autocomplete>
          <v-textarea
            v-model="batchForm.notes"
            auto-grow
            density="comfortable"
            label="Notes"
            rows="2"
          />
        </v-form>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="closeCreateDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="createBatch"
        >
          Create batch
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Bulk Import Dialog -->
  <v-dialog v-model="bulkImportDialog" :max-width="$vuetify.display.xs ? '100%' : 720">
    <v-card>
      <v-card-title class="text-h6">Bulk import batches</v-card-title>
      <v-card-text>
        <v-alert class="mb-4" density="compact" type="info" variant="tonal">
          Expected columns: short_name (required), brew_date (optional), notes (optional).
        </v-alert>

        <v-file-input
          v-model="importFile"
          accept=".csv,text/csv"
          density="comfortable"
          label="CSV file"
          placeholder="Select batch import CSV"
          prepend-icon="mdi-file-delimited-outline"
          show-size
        />

        <div class="d-flex align-center justify-space-between flex-wrap ga-2 mb-4">
          <div class="text-caption text-medium-emphasis">
            Brew date format: YYYY-MM-DD.
          </div>
          <v-btn prepend-icon="mdi-download" size="small" variant="text" @click="downloadBatchTemplate">
            Download template
          </v-btn>
        </div>

        <v-alert
          v-if="importSummary"
          class="mb-3"
          density="compact"
          :type="importSummary.type"
          variant="tonal"
        >
          {{ importSummary.message }}
        </v-alert>

        <v-table v-if="importErrors.length > 0" class="data-table" density="compact">
          <thead>
            <tr>
              <th>Row</th>
              <th>Issue</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(rowError, index) in importErrors" :key="index">
              <td>{{ rowError.row ?? '-' }}</td>
              <td>{{ rowError.message }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="importUploading" variant="text" @click="bulkImportDialog = false">
          Close
        </v-btn>
        <v-btn
          color="primary"
          :disabled="!canImport"
          :loading="importUploading"
          @click="uploadBatchImport"
        >
          Upload
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

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
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { BatchDeleteDialog, BatchEditDialog, type BatchEditForm } from '@/components/batch'
  import { useApiClient } from '@/composables/useApiClient'
  import { type Recipe, type UpdateBatchRequest, useProductionApi } from '@/composables/useProductionApi'

  type Batch = {
    id: number
    uuid: string
    short_name: string
    brew_date: string | null
    recipe_id: number | null
    notes: string | null
    created_at: string
    updated_at: string
  }

  type BatchWithSummary = Batch & {
    recipe_name?: string | null
    current_phase?: string | null
  }

  type ImportRowError = {
    row: number | null
    message: string
  }

  type BatchImportRowResult = {
    row: number
    status: 'created' | 'error'
    error?: string | null
    batch?: Batch
  }

  type BatchImportResponse = {
    totals: {
      total_rows: number
      created: number
      failed: number
    }
    results: BatchImportRowResult[]
  }

  type ImportSummary = {
    message: string
    type: 'success' | 'warning' | 'error'
  }

  const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'
  const router = useRouter()
  const { request } = useApiClient(apiBase)
  const { getRecipes, getBatchSummary, updateBatch, deleteBatch, normalizeText } = useProductionApi()

  // State
  const batches = ref<BatchWithSummary[]>([])
  const recipes = ref<Recipe[]>([])
  const loading = ref(false)
  const recipesLoading = ref(false)
  const saving = ref(false)
  const errorMessage = ref('')
  const search = ref('')

  // Dialogs
  const createBatchDialog = ref(false)
  const bulkImportDialog = ref(false)

  // Edit/Delete state
  const editBatchDialog = ref(false)
  const editingBatch = ref<BatchWithSummary | null>(null)
  const savingBatchEdit = ref(false)
  const editBatchError = ref('')

  const deleteBatchDialog = ref(false)
  const deletingBatchItem = ref<BatchWithSummary | null>(null)
  const deletingBatch = ref(false)
  const deleteBatchError = ref('')

  // Form
  const formRef = ref()
  const batchForm = reactive({
    short_name: '',
    brew_date: '',
    recipe_id: null as number | null,
    notes: '',
  })

  // Import state
  const importFile = ref<File | File[] | null>(null)
  const importUploading = ref(false)
  const importResult = ref<BatchImportResponse | null>(null)
  const importErrors = ref<ImportRowError[]>([])

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
  }

  // Table configuration
  const headers = [
    { title: 'ID', key: 'id', sortable: true, width: '80px' },
    { title: 'Short Name', key: 'short_name', sortable: true },
    { title: 'Recipe', key: 'recipe_name', sortable: true },
    { title: 'Status', key: 'current_phase', sortable: true },
    { title: 'Brew Date', key: 'brew_date', sortable: true },
    { title: 'Updated', key: 'updated_at', sortable: true },
    { title: '', key: 'actions', sortable: false, align: 'end' as const, width: '100px' },
  ]

  // Computed
  const isFormValid = computed(() => {
    return batchForm.short_name.trim().length > 0
  })

  const recipeSelectItems = computed(() =>
    recipes.value.map(recipe => ({
      title: recipe.name,
      value: recipe.id,
      style: recipe.style_name,
    })),
  )

  const canImport = computed(() => Boolean(getSelectedImportFile()) && !importUploading.value)

  const importSummary = computed<ImportSummary | null>(() => {
    if (!importResult.value) {
      return null
    }
    const successCount = getImportSuccessCount(importResult.value)
    const failureCount = getImportFailureCount(importResult.value, importErrors.value)
    const total = successCount + failureCount
    const message
      = total > 0
        ? `Imported ${successCount} ${successCount === 1 ? 'batch' : 'batches'}, ${failureCount} failed.`
        : 'Import completed.'
    if (failureCount > 0 && successCount === 0) {
      return { message, type: 'error' }
    }
    if (failureCount > 0) {
      return { message, type: 'warning' }
    }
    return { message, type: 'success' }
  })

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

  // Watch for import dialog close to reset state
  watch(bulkImportDialog, isOpen => {
    if (!isOpen) {
      resetImportState()
    }
  })

  watch(importFile, value => {
    if (value) {
      importResult.value = null
      importErrors.value = []
    }
  })

  // Methods
  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function loadBatches () {
    loading.value = true
    errorMessage.value = ''
    try {
      const batchList = await request<Batch[]>('/batches')

      // Enrich batches with summary data (recipe name, current phase)
      const enrichedBatches = await Promise.all(
        batchList.map(async batch => {
          try {
            const summary = await getBatchSummary(batch.id)
            return {
              ...batch,
              recipe_name: summary.recipe_name,
              current_phase: summary.current_phase,
            }
          } catch {
            // If summary fails, return batch without enrichment
            return {
              ...batch,
              recipe_name: null,
              current_phase: null,
            }
          }
        }),
      )

      batches.value = enrichedBatches
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load batches'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      loading.value = false
    }
  }

  async function loadRecipes () {
    recipesLoading.value = true
    try {
      recipes.value = await getRecipes()
    } catch (error) {
      // Recipe loading failure is non-critical
      console.error('Failed to load recipes:', error)
    } finally {
      recipesLoading.value = false
    }
  }

  function openCreateDialog () {
    batchForm.short_name = ''
    batchForm.brew_date = ''
    batchForm.recipe_id = null
    batchForm.notes = ''
    createBatchDialog.value = true
  }

  function closeCreateDialog () {
    createBatchDialog.value = false
  }

  async function createBatch () {
    if (!isFormValid.value) {
      return
    }

    saving.value = true
    errorMessage.value = ''

    try {
      const payload = {
        short_name: batchForm.short_name.trim(),
        brew_date: normalizeDateOnly(batchForm.brew_date),
        recipe_id: batchForm.recipe_id,
        notes: normalizeText(batchForm.notes),
      }

      await request<Batch>('/batches', {
        method: 'POST',
        body: JSON.stringify(payload),
      })

      showNotice('Batch created')
      closeCreateDialog()
      await loadBatches()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create batch'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }

  function onRowDoubleClick (_event: Event, { item }: { item: BatchWithSummary }) {
    router.push(`/batches/${item.uuid}`)
  }

  // Edit/Delete functions
  function openEditDialog (batch: BatchWithSummary) {
    editingBatch.value = batch
    editBatchError.value = ''
    editBatchDialog.value = true
  }

  async function saveBatchEdit (form: BatchEditForm) {
    if (!editingBatch.value) return

    savingBatchEdit.value = true
    editBatchError.value = ''

    try {
      const payload: UpdateBatchRequest = {
        short_name: form.short_name.trim(),
        brew_date: form.brew_date ? normalizeDateOnly(form.brew_date) : null,
        recipe_id: form.recipe_id,
        notes: normalizeText(form.notes),
      }

      await updateBatch(editingBatch.value.id, payload)
      showNotice('Batch updated')
      editBatchDialog.value = false
      editingBatch.value = null
      await loadBatches()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update batch'
      editBatchError.value = message
    } finally {
      savingBatchEdit.value = false
    }
  }

  function openDeleteDialog (batch: BatchWithSummary) {
    deletingBatchItem.value = batch
    deleteBatchError.value = ''
    deleteBatchDialog.value = true
  }

  async function confirmDeleteBatch () {
    if (!deletingBatchItem.value) return

    deletingBatch.value = true
    deleteBatchError.value = ''

    try {
      await deleteBatch(deletingBatchItem.value.id)
      showNotice('Batch deleted')
      deleteBatchDialog.value = false
      deletingBatchItem.value = null
      await loadBatches()
    } catch (error) {
      // Check for 409 Conflict (batch has dependencies)
      if (error instanceof Error && error.message.includes('409')) {
        deleteBatchError.value = 'Cannot delete this batch because it has associated brew sessions, measurements, or other data. Remove those first.'
      } else if (error instanceof Error && (error.message.toLowerCase().includes('conflict') || error.message.toLowerCase().includes('dependencies'))) {
        deleteBatchError.value = 'Cannot delete this batch because it has associated brew sessions, measurements, or other data. Remove those first.'
      } else {
        const message = error instanceof Error ? error.message : 'Failed to delete batch'
        deleteBatchError.value = message
      }
    } finally {
      deletingBatch.value = false
    }
  }

  // Import functions
  function getSelectedImportFile (): File | null {
    if (!importFile.value) return null
    return Array.isArray(importFile.value) ? importFile.value[0] ?? null : importFile.value
  }

  function resetImportState () {
    importFile.value = null
    importResult.value = null
    importErrors.value = []
  }

  function getImportSuccessCount (response: BatchImportResponse): number {
    return response.totals?.created ?? response.results.filter(r => r.status === 'created').length
  }

  function getImportFailureCount (response: BatchImportResponse, errors: ImportRowError[]): number {
    const fromResponse = response.totals?.failed ?? response.results.filter(r => r.status === 'error').length
    return fromResponse > 0 ? fromResponse : errors.length
  }

  function parseImportErrors (response: BatchImportResponse): ImportRowError[] {
    return response.results
      .filter(r => r.status === 'error')
      .map(r => ({
        row: r.row,
        message: r.error ?? 'Unknown error',
      }))
  }

  async function uploadBatchImport () {
    const file = getSelectedImportFile()
    if (!file) {
      return
    }
    errorMessage.value = ''
    importUploading.value = true
    importResult.value = null
    importErrors.value = []
    try {
      const form = new FormData()
      form.append('file', file)
      const response = await request<BatchImportResponse>('/batches/import', {
        method: 'POST',
        body: form,
        headers: new Headers(),
      })
      importResult.value = response
      importErrors.value = parseImportErrors(response)
      const successCount = getImportSuccessCount(response)
      const failureCount = getImportFailureCount(response, importErrors.value)
      if (failureCount > 0) {
        const color = successCount > 0 ? 'warning' : 'error'
        showNotice(`Imported ${successCount} batches, ${failureCount} failed`, color)
      } else {
        showNotice(`Imported ${successCount} ${successCount === 1 ? 'batch' : 'batches'}`)
      }
      await loadBatches()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Import failed'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      importUploading.value = false
    }
  }

  function downloadBatchTemplate () {
    const header = 'short_name,brew_date,notes\n'
    const blob = new Blob([header], { type: 'text/csv;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'batch-import-template.csv'
    link.click()
    URL.revokeObjectURL(url)
  }

  // Formatting functions
  function normalizeDateOnly (value: string) {
    return value ? new Date(`${value}T00:00:00Z`).toISOString() : null
  }

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

  function formatPhase (phase: string) {
    const phaseLabels: Record<string, string> = {
      planning: 'Planning',
      mashing: 'Mashing',
      heating: 'Heating',
      boiling: 'Boiling',
      cooling: 'Cooling',
      fermenting: 'Fermenting',
      conditioning: 'Conditioning',
      packaging: 'Packaging',
      finished: 'Finished',
    }
    return phaseLabels[phase] ?? phase.charAt(0).toUpperCase() + phase.slice(1)
  }

  function getPhaseColor (phase: string) {
    const phaseColors: Record<string, string> = {
      planning: 'grey',
      mashing: 'orange',
      heating: 'deep-orange',
      boiling: 'red',
      cooling: 'cyan',
      fermenting: 'primary',
      conditioning: 'teal',
      packaging: 'blue',
      finished: 'success',
    }
    return phaseColors[phase] ?? 'secondary'
  }
</script>

<style scoped>
.batches-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.card-title-responsive {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.card-title-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
}

.search-field {
  width: 200px;
  min-width: 120px;
  flex-shrink: 1;
}

@media (max-width: 599px) {
  .card-title-responsive {
    flex-direction: column;
    align-items: stretch;
  }

  .card-title-actions {
    justify-content: flex-end;
  }

  .search-field {
    width: 100%;
    max-width: none;
    order: 10;
    margin-top: 8px;
  }
}

.data-table {
  overflow-x: auto;
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  white-space: nowrap;
}

.data-table :deep(td) {
  font-size: 0.85rem;
}

.batches-table :deep(tr) {
  cursor: pointer;
}

.batches-table :deep(tr:hover td) {
  background: rgba(var(--v-theme-primary), 0.04);
}

/* Ensure table scrolls horizontally on mobile */
.batches-table :deep(.v-table__wrapper) {
  overflow-x: auto;
}
</style>
