<template>
  <v-container class="production-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="4">
        <BatchList
          :batches="inProgressBatches"
          :loading="loading"
          :selected-batch-id="selectedBatchId"
          :show-bulk-import="false"
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

</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue'
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

  const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'
  const { request } = useApiClient(apiBase)
  const { getRecipes, getBatchSummary } = useProductionApi()

  const loading = ref(false)
  const batches = ref<Batch[]>([])
  const batchCurrentPhases = ref<Map<number, ProcessPhase | null>>(new Map())
  const recipes = ref<Recipe[]>([])
  const recipesLoading = ref(false)

  const selectedBatchId = ref<number | null>(null)
  const createBatchDialog = ref(false)

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

  // Filter batches to only show in-progress ones
  // In-progress = not finished (phase !== 'finished')
  const inProgressBatches = computed(() => {
    return batches.value.filter(batch => {
      const phase = batchCurrentPhases.value.get(batch.id)

      // If the batch is finished, exclude it
      if (phase === 'finished') {
        return false
      }

      // Include all other batches (planning, active brewing, fermenting, etc.)
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
    const batchesData = await get<Batch[]>('/batches')
    batches.value = batchesData

    // Load batch summaries to get current phases (in parallel)
    const phaseMap = new Map<number, ProcessPhase | null>()
    await Promise.all(
      batchesData.map(async batch => {
        try {
          const summary = await getBatchSummary(batch.id)
          phaseMap.set(batch.id, (summary.current_phase as ProcessPhase) ?? null)
        } catch {
          // If summary fails, treat as no phase
          phaseMap.set(batch.id, null)
        }
      }),
    )

    batchCurrentPhases.value = phaseMap
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
</style>
