<template>
  <v-container class="production-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        <v-icon class="mr-2" icon="mdi-book-open-page-variant" />
        Recipes
        <v-spacer />
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
          class="ml-2"
          :loading="loading"
          size="small"
          variant="text"
          @click="loadRecipes"
        >
          Refresh
        </v-btn>
        <v-btn
          class="ml-2"
          color="primary"
          size="small"
          variant="text"
          @click="openCreateDialog"
        >
          New recipe
        </v-btn>
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
          class="data-table recipes-table"
          density="compact"
          :headers="headers"
          hover
          item-value="uuid"
          :items="filteredRecipes"
          :loading="loading"
          @click:row="onRowClick"
        >
          <template #item.name="{ item }">
            <span class="font-weight-medium">{{ item.name }}</span>
          </template>

          <template #item.style_name="{ item }">
            <v-chip v-if="item.style_name" size="small" variant="tonal">
              {{ item.style_name }}
            </v-chip>
            <span v-else class="text-medium-emphasis">—</span>
          </template>

          <template #item.specs="{ item }">
            <div class="d-flex flex-wrap ga-1">
              <v-chip
                v-if="item.target_og"
                color="info"
                size="x-small"
                variant="tonal"
              >
                OG {{ formatGravity(item.target_og) }}
              </v-chip>
              <v-chip
                v-if="item.target_ibu"
                color="warning"
                size="x-small"
                variant="tonal"
              >
                {{ formatIbu(item.target_ibu) }} IBU
              </v-chip>
              <v-chip
                v-if="item.target_abv"
                color="success"
                size="x-small"
                variant="tonal"
              >
                {{ formatAbv(item.target_abv) }}
              </v-chip>
              <span v-if="!item.target_og && !item.target_ibu && !item.target_abv" class="text-medium-emphasis">—</span>
            </div>
          </template>

          <template #item.notes="{ item }">
            <span v-if="item.notes" class="notes-cell">{{ item.notes }}</span>
            <span v-else class="text-medium-emphasis">—</span>
          </template>

          <template #item.updated_at="{ item }">
            {{ formatDateTime(item.updated_at) }}
          </template>

          <template #no-data>
            <div class="text-center py-4">
              <div class="text-body-2 text-medium-emphasis">No recipes yet.</div>
              <v-btn
                class="mt-2"
                color="primary"
                size="small"
                variant="text"
                @click="openCreateDialog"
              >
                Create your first recipe
              </v-btn>
            </div>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </v-container>

  <!-- Create/Edit Recipe Dialog -->
  <RecipeCreateEditDialog
    v-model="recipeDialog"
    :edit-recipe="editingRecipe"
    :saving="saving"
    :styles="styles"
    @submit="saveRecipe"
  />

  <!-- Delete Confirmation Dialog -->
  <RecipeDeleteDialog
    v-if="recipeToDelete"
    v-model="deleteDialog"
    :deleting="deleting"
    :recipe="recipeToDelete"
    @confirm="confirmDelete"
  />
</template>

<script lang="ts" setup>
  import type { Recipe, Style } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import RecipeCreateEditDialog from '@/components/recipe/RecipeCreateEditDialog.vue'
  import type { RecipeCreateEditSubmitData } from '@/components/recipe/RecipeCreateEditDialog.vue'
  import { RecipeDeleteDialog } from '@/components/recipe'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { formatDateTime, useBrewingFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'

  const router = useRouter()

  const { formatGravity, formatPercent: formatAbv, formatWholeNumber: formatIbu } = useBrewingFormatters()

  const {
    getRecipes,
    getStyles,
    createRecipe,
    updateRecipe,
    deleteRecipe,
  } = useProductionApi()
  const { showNotice } = useSnackbar()

  // State
  const recipes = ref<Recipe[]>([])
  const styles = ref<Style[]>([])
  const search = ref('')

  const { execute: executeLoad, loading, error: errorMessage } = useAsyncAction()
  const { execute: executeSave, loading: saving, error: saveError } = useAsyncAction()
  const { execute: executeDelete, loading: deleting, error: deleteError } = useAsyncAction()
  const { execute: executeLoadStyles } = useAsyncAction({ onError: () => {} })

  // Dialogs
  const recipeDialog = ref(false)
  const deleteDialog = ref(false)
  const editingRecipe = ref<Recipe | null>(null)
  const recipeToDelete = ref<Recipe | null>(null)

  // Table configuration
  const headers = [
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Style', key: 'style_name', sortable: true },
    { title: 'Specs', key: 'specs', sortable: false },
    { title: 'Notes', key: 'notes', sortable: false },
    { title: 'Updated', key: 'updated_at', sortable: true },
  ]

  // Computed
  const filteredRecipes = computed(() => {
    if (!search.value) {
      return recipes.value
    }
    const query = search.value.toLowerCase()
    return recipes.value.filter(
      recipe =>
        recipe.name.toLowerCase().includes(query)
        || (recipe.style_name?.toLowerCase().includes(query) ?? false)
        || (recipe.notes?.toLowerCase().includes(query) ?? false),
    )
  })

  // Lifecycle
  onMounted(async () => {
    await Promise.all([loadRecipes(), loadStyles()])
  })

  // Methods
  async function loadRecipes () {
    await executeLoad(async () => {
      recipes.value = await getRecipes()
    })
    if (errorMessage.value) {
      showNotice(errorMessage.value, 'error')
    }
  }

  async function loadStyles () {
    // Styles loading failure is non-critical, user can still type new styles
    await executeLoadStyles(async () => {
      styles.value = await getStyles()
    })
  }

  function openCreateDialog () {
    editingRecipe.value = null
    recipeDialog.value = true
  }

  function onRowClick (_event: Event, { item }: { item: Recipe }) {
    router.push(`/production/recipes/${item.uuid}`)
  }

  async function saveRecipe (data: RecipeCreateEditSubmitData) {
    await executeSave(async () => {
      if (editingRecipe.value) {
        await updateRecipe(editingRecipe.value.uuid, data)
        showNotice('Recipe updated')
      } else {
        await createRecipe(data)
        showNotice('Recipe created')
      }

      recipeDialog.value = false
      editingRecipe.value = null
      await Promise.all([loadRecipes(), loadStyles()])
    })
    if (saveError.value) {
      errorMessage.value = saveError.value
      showNotice(saveError.value, 'error')
    }
  }

  async function confirmDelete () {
    if (!recipeToDelete.value) {
      return
    }

    await executeDelete(async () => {
      await deleteRecipe(recipeToDelete.value!.uuid)
      showNotice('Recipe deleted')
      deleteDialog.value = false
      recipeToDelete.value = null
      await loadRecipes()
    })
    if (deleteError.value) {
      errorMessage.value = deleteError.value
      showNotice(deleteError.value, 'error')
    }
  }
</script>

<style scoped>
.production-page {
  position: relative;
}

.notes-cell {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  max-width: 300px;
}

.recipes-table :deep(tr) {
  cursor: pointer;
}

.recipes-table :deep(tr:hover td) {
  background: rgba(var(--v-theme-primary), 0.04);
}
</style>
