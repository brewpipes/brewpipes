<template>
  <v-container class="production-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="4">
        <BatchList
          :batches="inProgressBatches"
          :loading="loading"
          :selected-batch-id="selectedBatchId"
          @bulk-import="bulkImportDialog = true"
          @create="createBatchDialog = true"
          @select="selectBatch"
        />
      </v-col>

      <v-col cols="12" md="8">
        <BatchDetails
          :batch-id="selectedBatchId"
          @cleared="clearSelection"
        />
      </v-col>
    </v-row>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <v-dialog v-model="createBatchDialog" max-width="520">
    <v-card>
      <v-card-title class="text-h6">Create batch</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="newBatch.short_name"
          density="comfortable"
          label="Short name"
          placeholder="IPA 24-07"
        />
        <v-text-field
          v-model="newBatch.brew_date"
          density="comfortable"
          label="Brew date"
          type="date"
        />
        <v-autocomplete
          v-model="newBatch.recipe_id"
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
          v-model="newBatch.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="createBatchDialog = false">Cancel</v-btn>
        <v-btn color="primary" :disabled="!newBatch.short_name.trim()" @click="createBatch">
          Create batch
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="bulkImportDialog" max-width="720">
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
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import BatchDetails from '@/components/BatchDetails.vue'
  import BatchList, { type Batch } from '@/components/BatchList.vue'
  import { useApiClient } from '@/composables/useApiClient'
  import { type Recipe, useProductionApi } from '@/composables/useProductionApi'

  type ProcessPhase
    = | 'planning'
      | 'mashing'
      | 'heating'
      | 'boiling'
      | 'cooling'
      | 'fermenting'
      | 'conditioning'
      | 'packaging'
      | 'finished'

  type BatchProcessPhase = {
    id: number
    uuid: string
    batch_id: number
    process_phase: ProcessPhase
    phase_at: string
    created_at: string
    updated_at: string
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
  const { request } = useApiClient(apiBase)
  const { getRecipes } = useProductionApi()

  const loading = ref(false)
  const batches = ref<Batch[]>([])
  const batchProcessPhases = ref<Map<number, ProcessPhase>>(new Map())
  const recipes = ref<Recipe[]>([])
  const recipesLoading = ref(false)

  const selectedBatchId = ref<number | null>(null)
  const createBatchDialog = ref(false)
  const bulkImportDialog = ref(false)

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  const newBatch = reactive({
    short_name: '',
    brew_date: '',
    recipe_id: null as number | null,
    notes: '',
  })

  const importFile = ref<File | File[] | null>(null)
  const importUploading = ref(false)
  const importResult = ref<BatchImportResponse | null>(null)
  const importErrors = ref<ImportRowError[]>([])

  // Filter batches to only show in-progress ones
  // In-progress = not finished AND (brew_date is null OR brew_date is in the past or today OR has no process phase)
  const inProgressBatches = computed(() => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)

    return batches.value.filter(batch => {
      const phase = batchProcessPhases.value.get(batch.id)

      // If the batch is finished, exclude it
      if (phase === 'finished') {
        return false
      }

      // Include batches that:
      // 1. Have no brew date (planning phase)
      // 2. Have a brew date in the future (upcoming)
      // 3. Have a brew date in the past or today (active)
      // All of these are considered "in progress" as long as not finished
      return true
    })
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

  onMounted(async () => {
    await refreshAll()
  })

  function selectBatch (id: number) {
    selectedBatchId.value = id
  }

  function clearSelection () {
    selectedBatchId.value = null
  }

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  function get<T> (path: string) {
    return request<T>(path)
  }

  function post<T> (path: string, payload: unknown) {
    return request<T>(path, { method: 'POST', body: JSON.stringify(payload) })
  }

  function postForm<T> (path: string, payload: FormData) {
    return request<T>(path, { method: 'POST', body: payload, headers: new Headers() })
  }

  async function refreshAll () {
    loading.value = true
    try {
      await Promise.all([loadBatches(), loadRecipesData()])
      // Auto-select first batch if none selected
      const firstBatch = inProgressBatches.value[0]
      if (!selectedBatchId.value && firstBatch) {
        selectedBatchId.value = firstBatch.id
      }
    } catch (error) {
      handleError(error)
    } finally {
      loading.value = false
    }
  }

  async function loadBatches () {
    const [batchesData, phasesData] = await Promise.all([
      get<Batch[]>('/batches'),
      get<BatchProcessPhase[]>('/batch-process-phases'),
    ])

    batches.value = batchesData

    // Build a map of batch ID to latest process phase
    const phaseMap = new Map<number, { phase: ProcessPhase, phase_at: string }>()
    for (const phase of phasesData) {
      const existing = phaseMap.get(phase.batch_id)
      if (!existing || new Date(phase.phase_at) > new Date(existing.phase_at)) {
        phaseMap.set(phase.batch_id, { phase: phase.process_phase, phase_at: phase.phase_at })
      }
    }

    batchProcessPhases.value = new Map(
      Array.from(phaseMap.entries()).map(([id, { phase }]) => [id, phase]),
    )
  }

  async function loadRecipesData () {
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

  async function createBatch () {
    try {
      const payload = {
        short_name: newBatch.short_name.trim(),
        brew_date: normalizeDateOnly(newBatch.brew_date),
        recipe_id: newBatch.recipe_id,
        notes: normalizeText(newBatch.notes),
      }
      const created = await post<Batch>('/batches', payload)
      showNotice('Batch created')
      newBatch.short_name = ''
      newBatch.brew_date = ''
      newBatch.recipe_id = null
      newBatch.notes = ''
      await loadBatches()
      selectedBatchId.value = created.id
      createBatchDialog.value = false
    } catch (error) {
      handleError(error)
    }
  }

  async function uploadBatchImport () {
    const file = getSelectedImportFile()
    if (!file) {
      return
    }
    importUploading.value = true
    importResult.value = null
    importErrors.value = []
    try {
      const form = new FormData()
      form.append('file', file)
      const response = await postForm<BatchImportResponse>('/batches/import', form)
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
      handleError(error)
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

  function getSelectedImportFile () {
    const fileValue = importFile.value
    if (!fileValue) {
      return null
    }
    return Array.isArray(fileValue) ? fileValue[0] ?? null : fileValue
  }

  function resetImportState () {
    importFile.value = null
    importResult.value = null
    importErrors.value = []
    importUploading.value = false
  }

  function getImportSuccessCount (result: BatchImportResponse) {
    return result.totals?.created ?? 0
  }

  function getImportFailureCount (result: BatchImportResponse, errors: ImportRowError[]) {
    return result.totals?.failed ?? errors.length
  }

  function parseImportErrors (result: BatchImportResponse) {
    return (result.results ?? [])
      .filter(entry => entry.status === 'error')
      .map(entry => ({
        row: entry.row ?? null,
        message: entry.error ?? 'Import failed',
      }))
  }

  function handleError (error: unknown) {
    const message = error instanceof Error ? error.message : 'Unexpected error'
    showNotice(message, 'error')
  }

  function normalizeText (value: string) {
    const trimmed = value.trim()
    return trimmed.length > 0 ? trimmed : null
  }

  function normalizeDateOnly (value: string) {
    return value ? new Date(`${value}T00:00:00Z`).toISOString() : null
  }
</script>

<style scoped>
.production-page {
  position: relative;
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
}

.data-table :deep(td) {
  font-size: 0.85rem;
}
</style>
