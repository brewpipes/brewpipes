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
  <v-dialog v-model="recipeDialog" max-width="600" persistent>
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit recipe' : 'Create recipe' }}
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef" @submit.prevent="saveRecipe">
          <v-text-field
            v-model="recipeForm.name"
            density="comfortable"
            label="Name"
            placeholder="West Coast IPA"
            :rules="[rules.required]"
          />

          <v-combobox
            v-model="recipeForm.style"
            density="comfortable"
            hint="Select an existing style or type a new one"
            item-title="name"
          item-value="uuid"
            :items="styleItems"
            label="Style"
            :loading="stylesLoading"
            persistent-hint
            return-object
            @update:search="onStyleSearch"
          >
            <template #no-data>
              <v-list-item v-if="styleSearchQuery">
                <v-list-item-title>
                  Press enter to create "{{ styleSearchQuery }}"
                </v-list-item-title>
              </v-list-item>
              <v-list-item v-else>
                <v-list-item-title>
                  Type to search or create a new style
                </v-list-item-title>
              </v-list-item>
            </template>
          </v-combobox>

          <v-textarea
            v-model="recipeForm.notes"
            auto-grow
            density="comfortable"
            label="Notes"
            placeholder="Recipe description, ingredients, process notes..."
            rows="3"
          />
        </v-form>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="closeRecipeDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="saveRecipe"
        >
          {{ isEditing ? 'Save changes' : 'Create recipe' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Delete Confirmation Dialog -->
  <v-dialog v-model="deleteDialog" max-width="400" persistent>
    <v-card>
      <v-card-title class="text-h6">Delete recipe</v-card-title>
      <v-card-text>
        <p>
          Are you sure you want to delete <strong>{{ recipeToDelete?.name }}</strong>?
        </p>
        <p class="text-medium-emphasis mt-2">This action cannot be undone.</p>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="deleting" variant="text" @click="closeDeleteDialog">Cancel</v-btn>
        <v-btn
          color="error"
          :loading="deleting"
          variant="flat"
          @click="confirmDelete"
        >
          Delete
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { formatDateTime } from '@/composables/useFormatters'
  import type { Recipe, Style } from '@/types'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { normalizeText } from '@/utils/normalize'

  const router = useRouter()

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
  const loading = ref(false)
  const stylesLoading = ref(false)
  const saving = ref(false)
  const deleting = ref(false)
  const errorMessage = ref('')
  const search = ref('')

  // Dialogs
  const recipeDialog = ref(false)
  const deleteDialog = ref(false)
  const editingRecipeUuid = ref<string | null>(null)
  const recipeToDelete = ref<Recipe | null>(null)

  // Form
  const formRef = ref()
  const styleSearchQuery = ref('')
  const recipeForm = reactive({
    name: '',
    style: null as Style | string | null,
    notes: '',
  })

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
  }

  // Table configuration
  const headers = [
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Style', key: 'style_name', sortable: true },
    { title: 'Specs', key: 'specs', sortable: false },
    { title: 'Notes', key: 'notes', sortable: false },
    { title: 'Updated', key: 'updated_at', sortable: true },
  ]

  // Computed
  const isEditing = computed(() => editingRecipeUuid.value !== null)

  const isFormValid = computed(() => {
    return recipeForm.name.trim().length > 0
  })

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

  const styleItems = computed(() => styles.value)

  // Lifecycle
  onMounted(async () => {
    await Promise.all([loadRecipes(), loadStyles()])
  })

  // Methods
  async function loadRecipes () {
    loading.value = true
    errorMessage.value = ''
    try {
      recipes.value = await getRecipes()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load recipes'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      loading.value = false
    }
  }

  async function loadStyles () {
    stylesLoading.value = true
    try {
      styles.value = await getStyles()
    } catch (error) {
      // Styles loading failure is non-critical, user can still type new styles
      console.error('Failed to load styles:', error)
    } finally {
      stylesLoading.value = false
    }
  }

  function onStyleSearch (query: string) {
    styleSearchQuery.value = query
  }

  function openCreateDialog () {
    editingRecipeUuid.value = null
    recipeForm.name = ''
    recipeForm.style = null
    recipeForm.notes = ''
    styleSearchQuery.value = ''
    recipeDialog.value = true
  }

  function onRowClick (_event: Event, { item }: { item: Recipe }) {
    router.push(`/production/recipes/${item.uuid}`)
  }

  function formatGravity (value: number | null | undefined): string {
    if (value === null || value === undefined) return ''
    return value.toFixed(3)
  }

  function formatAbv (value: number | null | undefined): string {
    if (value === null || value === undefined) return ''
    return `${value.toFixed(1)}%`
  }

  function formatIbu (value: number | null | undefined): string {
    if (value === null || value === undefined) return ''
    return String(Math.round(value))
  }

  function closeRecipeDialog () {
    recipeDialog.value = false
    editingRecipeUuid.value = null
  }

  async function saveRecipe () {
    if (!isFormValid.value) {
      return
    }

    saving.value = true
    errorMessage.value = ''

    try {
      // Determine style_uuid and style_name from the form value
      let styleUuid: string | null = null
      let styleName: string | null = null

      if (recipeForm.style) {
        if (typeof recipeForm.style === 'object' && recipeForm.style.uuid) {
          // User selected an existing style
          styleUuid = recipeForm.style.uuid
          styleName = recipeForm.style.name
        } else if (typeof recipeForm.style === 'string') {
          // User typed a new style name
          styleName = recipeForm.style.trim() || null
        }
      }

      const payload = {
        name: recipeForm.name.trim(),
        style_uuid: styleUuid,
        style_name: styleName,
        notes: normalizeText(recipeForm.notes),
      }

      if (isEditing.value && editingRecipeUuid.value) {
        await updateRecipe(editingRecipeUuid.value, payload)
        showNotice('Recipe updated')
      } else {
        await createRecipe(payload)
        showNotice('Recipe created')
      }

      closeRecipeDialog()
      await Promise.all([loadRecipes(), loadStyles()])
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save recipe'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }

  function closeDeleteDialog () {
    deleteDialog.value = false
    recipeToDelete.value = null
  }

  async function confirmDelete () {
    if (!recipeToDelete.value) {
      return
    }

    deleting.value = true
    errorMessage.value = ''

    try {
      await deleteRecipe(recipeToDelete.value.uuid)
      showNotice('Recipe deleted')
      closeDeleteDialog()
      await loadRecipes()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to delete recipe'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      deleting.value = false
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
