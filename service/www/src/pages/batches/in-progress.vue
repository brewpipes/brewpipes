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

  <BatchCreateDialog
    v-model="createBatchDialog"
    :recipes="recipes"
    :recipes-loading="recipesLoading"
    :saving="saving"
    @submit="handleCreateBatch"
  />

</template>

<script lang="ts" setup>
  import type { Batch, Recipe } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import BatchDetails from '@/components/BatchDetails.vue'
  import { BatchCreateDialog, type BatchCreateForm } from '@/components/batch'
  import BatchList from '@/components/BatchList.vue'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { normalizeDateOnly, normalizeText } from '@/utils/normalize'

  const { getBatches, createBatch: createBatchApi, getRecipes } = useProductionApi()

  const loading = ref(false)
  const saving = ref(false)
  const batches = ref<Batch[]>([])
  const recipes = ref<Recipe[]>([])
  const recipesLoading = ref(false)

  const selectedBatchUuid = ref<string | null>(null)
  const createBatchDialog = ref(false)
  const { showNotice } = useSnackbar()

  // Filter batches to only show in-progress ones
  // In-progress = not finished (phase !== 'finished')
  const inProgressBatches = computed(() => {
    return batches.value.filter(batch => batch.current_phase !== 'finished')
  })

  onMounted(async () => {
    await refreshAll()
  })

  function selectBatch (uuid: string) {
    selectedBatchUuid.value = uuid
  }

  function clearSelection () {
    selectedBatchUuid.value = null
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
    batches.value = await getBatches()
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

  async function handleCreateBatch (form: BatchCreateForm) {
    saving.value = true
    try {
      const payload = {
        short_name: form.short_name.trim(),
        brew_date: normalizeDateOnly(form.brew_date),
        recipe_uuid: form.recipe_uuid,
        notes: normalizeText(form.notes),
      }
      const created = await createBatchApi(payload)
      showNotice('Batch created')
      await loadBatches()
      selectedBatchUuid.value = created.uuid
      createBatchDialog.value = false
    } catch (error) {
      handleError(error)
    } finally {
      saving.value = false
    }
  }

  function handleError (error: unknown) {
    const message = error instanceof Error ? error.message : 'Unexpected error'
    showNotice(message, 'error')
  }
</script>

<style scoped>
.production-page {
  position: relative;
}
</style>
