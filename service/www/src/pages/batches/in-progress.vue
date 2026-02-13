<template>
  <v-container class="production-page" fluid>
    <!-- Mobile: Show list or detail based on selection -->
    <v-row v-if="$vuetify.display.smAndDown" align="stretch">
      <v-col v-if="!selectedBatchUuid" cols="12">
        <BatchList
          :batches="inProgressBatches"
          :loading="loading"
          :selected-batch-uuid="selectedBatchUuid"
          :show-bulk-import="false"
          @create="createBatchDialog = true"
          @select="selectBatch"
        />
      </v-col>

      <v-col v-else cols="12">
        <BatchDetails
          :batch-uuid="selectedBatchUuid"
          @cleared="clearSelection"
        />
      </v-col>
    </v-row>

    <!-- Desktop: Side-by-side layout -->
    <v-row v-else align="stretch">
      <v-col cols="12" md="4">
        <BatchList
          :batches="inProgressBatches"
          :loading="loading"
          :selected-batch-uuid="selectedBatchUuid"
          :show-bulk-import="false"
          @create="createBatchDialog = true"
          @select="selectBatch"
        />
      </v-col>

      <v-col cols="12" md="8">
        <BatchDetails
          :batch-uuid="selectedBatchUuid"
          @cleared="clearSelection"
        />
      </v-col>
    </v-row>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <v-dialog v-model="createBatchDialog" :max-width="$vuetify.display.xs ? '100%' : 520">
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
          v-model="newBatch.recipe_uuid"
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
  import type { Batch } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import BatchDetails from '@/components/BatchDetails.vue'
  import BatchList from '@/components/BatchList.vue'
  import { useApiClient } from '@/composables/useApiClient'
  import { type Recipe, useProductionApi } from '@/composables/useProductionApi'

  const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'
  const { request } = useApiClient(apiBase)
  const { getRecipes } = useProductionApi()

  const loading = ref(false)
  const batches = ref<Batch[]>([])
  const recipes = ref<Recipe[]>([])
  const recipesLoading = ref(false)

  const selectedBatchUuid = ref<string | null>(null)
  const createBatchDialog = ref(false)

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  const newBatch = reactive({
    short_name: '',
    brew_date: '',
    recipe_uuid: null as string | null,
    notes: '',
  })

  // Filter batches to only show in-progress ones
  // In-progress = not finished (phase !== 'finished')
  const inProgressBatches = computed(() => {
    return batches.value.filter(batch => batch.current_phase !== 'finished')
  })

  const recipeSelectItems = computed(() =>
    recipes.value.map(recipe => ({
      title: recipe.name,
      value: recipe.uuid,
      style: recipe.style_name,
    })),
  )

  onMounted(async () => {
    await refreshAll()
  })

  function selectBatch (uuid: string) {
    selectedBatchUuid.value = uuid
  }

  function clearSelection () {
    selectedBatchUuid.value = null
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
      if (!selectedBatchUuid.value && firstBatch) {
        selectedBatchUuid.value = firstBatch.uuid
      }
    } catch (error) {
      handleError(error)
    } finally {
      loading.value = false
    }
  }

  async function loadBatches () {
    batches.value = await get<Batch[]>('/batches')
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
        recipe_uuid: newBatch.recipe_uuid,
        notes: normalizeText(newBatch.notes),
      }
      const created = await post<Batch>('/batches', payload)
      showNotice('Batch created')
      newBatch.short_name = ''
      newBatch.brew_date = ''
      newBatch.recipe_uuid = null
      newBatch.notes = ''
      await loadBatches()
      selectedBatchUuid.value = created.uuid
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
