<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 720"
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
  >
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
          <v-btn prepend-icon="mdi-download" size="small" variant="text" @click="downloadTemplate">
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
        <v-btn :disabled="saving" variant="text" @click="handleClose">
          Close
        </v-btn>
        <v-btn
          color="primary"
          :disabled="!canImport"
          :loading="saving"
          @click="handleSubmit"
        >
          Upload
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Batch } from '@/types'
  import { computed, ref, watch } from 'vue'

  export type ImportRowError = {
    row: number | null
    message: string
  }

  export type BatchImportRowResult = {
    row: number
    status: 'created' | 'error'
    error?: string | null
    batch?: Batch
  }

  export type BatchImportResponse = {
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

  const props = defineProps<{
    modelValue: boolean
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [file: File]
  }>()

  // Import state
  const importFile = ref<File | File[] | null>(null)
  const importResult = ref<BatchImportResponse | null>(null)
  const importErrors = ref<ImportRowError[]>([])

  const canImport = computed(() => Boolean(getSelectedFile()) && !props.saving)

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

  // Reset state when dialog closes
  watch(
    () => props.modelValue,
    isOpen => {
      if (!isOpen) {
        resetState()
      }
    },
  )

  // Clear results when a new file is selected
  watch(importFile, value => {
    if (value) {
      importResult.value = null
      importErrors.value = []
    }
  })

  function getSelectedFile (): File | null {
    if (!importFile.value) return null
    return Array.isArray(importFile.value) ? importFile.value[0] ?? null : importFile.value
  }

  function resetState () {
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

  function downloadTemplate () {
    const header = 'short_name,brew_date,notes\n'
    const blob = new Blob([header], { type: 'text/csv;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'batch-import-template.csv'
    link.click()
    URL.revokeObjectURL(url)
  }

  function handleSubmit () {
    const file = getSelectedFile()
    if (!file) return
    emit('submit', file)
  }

  function handleClose () {
    emit('update:modelValue', false)
  }

  /** Called by parent after import completes to display results */
  function setImportResult (response: BatchImportResponse) {
    importResult.value = response
    importErrors.value = parseImportErrors(response)
  }

  defineExpose({
    setImportResult,
  })
</script>
