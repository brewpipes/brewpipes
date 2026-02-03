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
          class="data-table"
          density="compact"
          :headers="headers"
          item-value="id"
          :items="filteredRecipes"
          :loading="loading"
        >
          <template #item.style_name="{ item }">
            <v-chip v-if="item.style_name" size="small" variant="tonal">
              {{ item.style_name }}
            </v-chip>
            <span v-else class="text-medium-emphasis">—</span>
          </template>

          <template #item.notes="{ item }">
            <span v-if="item.notes" class="notes-cell">{{ item.notes }}</span>
            <span v-else class="text-medium-emphasis">—</span>
          </template>

          <template #item.updated_at="{ item }">
            {{ formatDateTime(item.updated_at) }}
          </template>

          <template #item.actions="{ item }">
            <v-btn
              icon="mdi-eye"
              size="x-small"
              variant="text"
              @click="openViewDialog(item)"
            />
            <v-btn
              icon="mdi-pencil"
              size="x-small"
              variant="text"
              @click="openEditDialog(item)"
            />
            <v-btn
              color="error"
              icon="mdi-delete-outline"
              size="x-small"
              variant="text"
              @click="openDeleteDialog(item)"
            />
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

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <!-- View Recipe Dialog -->
  <v-dialog v-model="viewDialog" max-width="600">
    <v-card v-if="selectedRecipe">
      <v-card-title class="text-h6">{{ selectedRecipe.name }}</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <div class="text-caption text-medium-emphasis">Style</div>
            <div class="text-body-1">
              <v-chip v-if="selectedRecipe.style_name" size="small" variant="tonal">
                {{ selectedRecipe.style_name }}
              </v-chip>
              <span v-else class="text-medium-emphasis">Not specified</span>
            </div>
          </v-col>
          <v-col cols="12" md="6">
            <div class="text-caption text-medium-emphasis">Last updated</div>
            <div class="text-body-1">{{ formatDateTime(selectedRecipe.updated_at) }}</div>
          </v-col>
          <v-col v-if="selectedRecipe.notes" cols="12">
            <div class="text-caption text-medium-emphasis">Notes</div>
            <div class="text-body-1">{{ selectedRecipe.notes }}</div>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="viewDialog = false">Close</v-btn>
        <v-btn color="primary" variant="text" @click="openEditDialogFromView">Edit</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

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
            item-value="id"
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
  import { type Recipe, type Style, useProductionApi } from '@/composables/useProductionApi'

  const {
    getRecipes,
    getStyles,
    createRecipe,
    updateRecipe,
    deleteRecipe,
    normalizeText,
    formatDateTime,
  } = useProductionApi()

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
  const viewDialog = ref(false)
  const recipeDialog = ref(false)
  const deleteDialog = ref(false)
  const selectedRecipe = ref<Recipe | null>(null)
  const editingRecipeId = ref<number | null>(null)
  const recipeToDelete = ref<Recipe | null>(null)

  // Form
  const formRef = ref()
  const styleSearchQuery = ref('')
  const recipeForm = reactive({
    name: '',
    style: null as Style | string | null,
    notes: '',
  })

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
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Style', key: 'style_name', sortable: true },
    { title: 'Notes', key: 'notes', sortable: false },
    { title: 'Updated', key: 'updated_at', sortable: true },
    { title: '', key: 'actions', sortable: false, align: 'end' as const },
  ]

  // Computed
  const isEditing = computed(() => editingRecipeId.value !== null)

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
  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

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
    editingRecipeId.value = null
    recipeForm.name = ''
    recipeForm.style = null
    recipeForm.notes = ''
    styleSearchQuery.value = ''
    recipeDialog.value = true
  }

  function openViewDialog (recipe: Recipe) {
    selectedRecipe.value = recipe
    viewDialog.value = true
  }

  function openEditDialog (recipe: Recipe) {
    editingRecipeId.value = recipe.id
    recipeForm.name = recipe.name
    recipeForm.notes = recipe.notes ?? ''

    // Set the style - find matching style object or use the name as string
    if (recipe.style_id) {
      const matchingStyle = styles.value.find(s => s.id === recipe.style_id)
      recipeForm.style = matchingStyle ?? recipe.style_name
    } else if (recipe.style_name) {
      recipeForm.style = recipe.style_name
    } else {
      recipeForm.style = null
    }

    styleSearchQuery.value = ''
    recipeDialog.value = true
  }

  function openEditDialogFromView () {
    if (selectedRecipe.value) {
      viewDialog.value = false
      openEditDialog(selectedRecipe.value)
    }
  }

  function closeRecipeDialog () {
    recipeDialog.value = false
    editingRecipeId.value = null
  }

  async function saveRecipe () {
    if (!isFormValid.value) {
      return
    }

    saving.value = true
    errorMessage.value = ''

    try {
      // Determine style_id and style_name from the form value
      let styleId: number | null = null
      let styleName: string | null = null

      if (recipeForm.style) {
        if (typeof recipeForm.style === 'object' && recipeForm.style.id) {
          // User selected an existing style
          styleId = recipeForm.style.id
          styleName = recipeForm.style.name
        } else if (typeof recipeForm.style === 'string') {
          // User typed a new style name
          styleName = recipeForm.style.trim() || null
        }
      }

      const payload = {
        name: recipeForm.name.trim(),
        style_id: styleId,
        style_name: styleName,
        notes: normalizeText(recipeForm.notes),
      }

      if (isEditing.value && editingRecipeId.value) {
        await updateRecipe(editingRecipeId.value, payload)
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

  function openDeleteDialog (recipe: Recipe) {
    recipeToDelete.value = recipe
    deleteDialog.value = true
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
      await deleteRecipe(recipeToDelete.value.id)
      showNotice('Recipe deleted')
      closeDeleteDialog()
      await loadRecipes()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to delete recipe'
      // Check for conflict error (recipe used by batches)
      if (message.includes('used by') || message.includes('batch')) {
        showNotice(message, 'error')
      } else {
        errorMessage.value = message
        showNotice(message, 'error')
      }
    } finally {
      deleting.value = false
    }
  }
</script>

<style scoped>
.production-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.search-field {
  max-width: 260px;
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

.notes-cell {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  max-width: 300px;
}
</style>
