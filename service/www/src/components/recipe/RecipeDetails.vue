<template>
  <v-card class="section-card">
    <!-- Header -->
    <v-card-title class="d-flex align-center flex-wrap ga-2">
      <v-btn
        aria-label="Back to recipes"
        class="mr-2"
        icon="mdi-arrow-left"
        size="small"
        variant="text"
        @click="emit('back')"
      />
      <v-icon class="mr-2" icon="mdi-book-open-page-variant" />
      <span class="text-h6">{{ recipe.name }}</span>
      <v-chip
        v-if="recipe.style_name"
        class="ml-2"
        color="secondary"
        size="small"
        variant="tonal"
      >
        {{ recipe.style_name }}
      </v-chip>
      <v-spacer />
      <v-btn
        :loading="refreshing"
        size="small"
        variant="text"
        @click="refresh"
      >
        <span class="d-none d-sm-inline">Refresh</span>
        <v-icon class="d-sm-none" icon="mdi-refresh" />
      </v-btn>
      <v-btn
        aria-label="Edit recipe"
        icon="mdi-pencil"
        size="small"
        variant="text"
        @click="openEditDialog"
      />
      <v-btn
        aria-label="Delete recipe"
        color="error"
        icon="mdi-delete"
        size="small"
        variant="text"
        @click="openDeleteDialog"
      />
    </v-card-title>

    <v-card-text>
      <!-- Metric Cards -->
      <v-row class="mb-4" dense>
        <v-col cols="6" md="2">
          <div class="metric-card">
            <div class="metric-label">OG</div>
            <div class="metric-value">{{ formatGravity(recipe.target_og) }}</div>
          </div>
        </v-col>
        <v-col cols="6" md="2">
          <div class="metric-card">
            <div class="metric-label">FG</div>
            <div class="metric-value">{{ formatGravity(recipe.target_fg) }}</div>
          </div>
        </v-col>
        <v-col cols="6" md="2">
          <div class="metric-card">
            <div class="metric-label">ABV</div>
            <div class="metric-value">{{ formatPercent(recipe.target_abv) }}</div>
          </div>
        </v-col>
        <v-col cols="6" md="2">
          <div class="metric-card">
            <div class="metric-label">IBU</div>
            <div class="metric-value">{{ formatNumber(recipe.target_ibu) }}</div>
          </div>
        </v-col>
        <v-col cols="6" md="2">
          <div class="metric-card">
            <div class="metric-label">SRM</div>
            <div class="metric-value d-flex align-center justify-center ga-2">
              {{ formatNumber(recipe.target_srm) }}
              <div
                v-if="recipe.target_srm"
                class="srm-swatch"
                :style="{ backgroundColor: srmToColor(recipe.target_srm) }"
              />
            </div>
          </div>
        </v-col>
        <v-col cols="6" md="2">
          <div class="metric-card">
            <div class="metric-label">Batch Size</div>
            <div class="metric-value">
              {{ displayBatchSize }}
            </div>
            <v-btn
              v-if="recipe.batch_size"
              class="mt-1"
              :color="isScaling ? 'primary' : undefined"
              density="compact"
              :prepend-icon="isScaling ? 'mdi-close' : 'mdi-resize'"
              size="x-small"
              :variant="isScaling ? 'tonal' : 'text'"
              @click="toggleScaling"
            >
              {{ isScaling ? 'Reset' : 'Scale' }}
            </v-btn>
            <div v-if="!recipe.batch_size" class="text-caption text-medium-emphasis mt-1">
              Set in Specs to scale
            </div>
            <v-btn
              v-if="recipe.batch_size"
              class="mt-1"
              :color="isScaling ? 'primary' : undefined"
              density="compact"
              :prepend-icon="isScaling ? 'mdi-close' : 'mdi-resize'"
              size="x-small"
              :variant="isScaling ? 'tonal' : 'text'"
              @click="toggleScaling"
            >
              {{ isScaling ? 'Reset' : 'Scale' }}
            </v-btn>
            <div v-if="!recipe.batch_size" class="text-caption text-medium-emphasis mt-1">
              Set in Specs to scale
            </div>
          </div>
        </v-col>
      </v-row>

      <!-- Scaling Controls -->
      <v-expand-transition>
        <div v-if="showScalingControls" class="scaling-controls mb-4">
          <v-row align="center" dense>
            <v-col cols="12" sm="auto">
              <span class="text-body-2 font-weight-medium">Scale to:</span>
            </v-col>
            <v-col cols="5" sm="3" md="2">
              <v-text-field
                v-model.number="targetBatchSize"
                density="compact"
                hide-details
                inputmode="decimal"
                label="Size"
                min="0.01"
                step="0.1"
                type="number"
                variant="outlined"
              />
            </v-col>
            <v-col cols="4" sm="2" md="2">
              <v-select
                v-model="targetBatchSizeUnit"
                density="compact"
                hide-details
                :items="scalingVolumeUnits"
                label="Unit"
                variant="outlined"
              />
            </v-col>
            <v-col cols="3" sm="auto">
              <v-chip
                v-if="isScaling"
                color="primary"
                size="small"
                variant="tonal"
              >
                {{ scaleFactorDisplay }}
              </v-chip>
            </v-col>
            <v-col class="d-none d-sm-flex" cols="auto">
              <v-btn
                color="secondary"
                density="compact"
                size="small"
                variant="text"
                @click="resetScaling"
              >
                Reset to Recipe
              </v-btn>
            </v-col>
          </v-row>
        </div>
      </v-expand-transition>

      <!-- Tabs -->
      <v-tabs v-model="activeTab" class="recipe-tabs" color="primary" show-arrows>
        <v-tab value="overview">
          <v-icon class="d-sm-none" icon="mdi-information-outline" />
          <span class="d-none d-sm-inline">Overview</span>
        </v-tab>
        <v-tab value="fermentables">
          <v-icon class="d-sm-none" icon="mdi-barley" />
          <span class="d-none d-sm-inline">Fermentables</span>
        </v-tab>
        <v-tab value="hops">
          <v-icon class="d-sm-none" icon="mdi-leaf" />
          <span class="d-none d-sm-inline">Hops</span>
        </v-tab>
        <v-tab value="yeast">
          <v-icon class="d-sm-none" icon="mdi-flask-outline" />
          <span class="d-none d-sm-inline">Yeast & Other</span>
        </v-tab>
        <v-tab value="specs">
          <v-icon class="d-sm-none" icon="mdi-tune" />
          <span class="d-none d-sm-inline">Specs</span>
        </v-tab>
      </v-tabs>

      <v-window v-model="activeTab" class="mt-4">
        <v-window-item value="overview">
          <RecipeOverviewTab
            :recipe="recipe"
          />
        </v-window-item>

        <v-window-item value="fermentables">
          <RecipeFermentablesTab
            :ingredients="fermentables"
            :is-scaling="isScaling"
            :loading="ingredientsLoading"
            :scale-amount="scaleAmount"
            @create="openIngredientDialog('fermentable')"
            @delete="deleteIngredient"
            @edit="openEditIngredientDialog"
          />
        </v-window-item>

        <v-window-item value="hops">
          <RecipeHopsTab
            :ingredients="hops"
            :is-scaling="isScaling"
            :loading="ingredientsLoading"
            :scale-amount="scaleAmount"
            @create="openIngredientDialog('hop')"
            @delete="deleteIngredient"
            @edit="openEditIngredientDialog"
          />
        </v-window-item>

        <v-window-item value="yeast">
          <RecipeYeastTab
            :adjuncts="adjuncts"
            :is-scaling="isScaling"
            :loading="ingredientsLoading"
            :scale-amount="scaleAmount"
            :water-chemistry="waterChemistry"
            :yeasts="yeasts"
            @create="openIngredientDialog"
            @delete="deleteIngredient"
            @edit="openEditIngredientDialog"
          />
        </v-window-item>

        <v-window-item value="specs">
          <RecipeSpecsTab
            :recipe="recipe"
            :saving="savingSpecs"
            @save="saveSpecs"
          />
        </v-window-item>
      </v-window>
    </v-card-text>
  </v-card>

  <!-- Edit Recipe Dialog -->
  <RecipeEditDialog
    v-model="editDialog"
    :error-message="editRecipeError"
    :recipe="recipe"
    :saving="savingRecipe"
    :styles="styles"
    :styles-loading="stylesLoading"
    @submit="saveRecipe"
  />

  <!-- Delete Recipe Dialog -->
  <RecipeDeleteDialog
    v-model="deleteDialog"
    :deleting="deletingRecipe"
    :error-message="deleteRecipeError"
    :recipe="recipe"
    @confirm="confirmDelete"
  />

  <!-- Ingredient Dialog -->
  <RecipeIngredientDialog
    v-model="ingredientDialog"
    :editing-ingredient="editingIngredient"
    :ingredient-type="ingredientDialogType"
    :saving="savingIngredient"
    @submit="saveIngredient"
  />

  <!-- Delete Ingredient Confirmation Dialog -->
  <v-dialog
    v-model="deleteIngredientDialog"
    :fullscreen="false"
    max-width="400"
    persistent
  >
    <v-card>
      <v-card-title class="text-h6">Delete Ingredient</v-card-title>
      <v-card-text>
        <p>
          Are you sure you want to delete
          <strong>{{ ingredientToDelete?.name ?? 'this ingredient' }}</strong>?
        </p>
        <v-alert
          class="mt-4"
          density="compact"
          type="warning"
          variant="tonal"
        >
          This action cannot be undone.
        </v-alert>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="deletingIngredient" variant="text" @click="cancelDeleteIngredient">
          Cancel
        </v-btn>
        <v-btn
          color="error"
          :loading="deletingIngredient"
          variant="flat"
          @click="confirmDeleteIngredient"
        >
          Delete
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type {
    CreateRecipeIngredientRequest,
    Recipe,
    RecipeIngredient,
    RecipeIngredientType,
    Style,
    UpdateRecipeIngredientRequest,
    UpdateRecipeRequest,
  } from '@/types'
  import { computed, onMounted, ref, watch } from 'vue'
  import { srmToColor, useBrewingFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useRecipeScaling } from '@/composables/useRecipeScaling'
  import { useSnackbar } from '@/composables/useSnackbar'
  import RecipeDeleteDialog from './RecipeDeleteDialog.vue'
  import RecipeEditDialog from './RecipeEditDialog.vue'
  import RecipeFermentablesTab from './RecipeFermentablesTab.vue'
  import RecipeHopsTab from './RecipeHopsTab.vue'
  import RecipeIngredientDialog from './RecipeIngredientDialog.vue'
  import RecipeOverviewTab from './RecipeOverviewTab.vue'
  import RecipeSpecsTab from './RecipeSpecsTab.vue'
  import RecipeYeastTab from './RecipeYeastTab.vue'

  const props = defineProps<{
    recipe: Recipe
  }>()

  const emit = defineEmits<{
    back: []
    deleted: []
    updated: [recipe: Recipe]
  }>()

  const {
    getStyles,
    getRecipe,
    updateRecipe,
    deleteRecipe,
    getRecipeIngredients,
    createRecipeIngredient,
    updateRecipeIngredient,
    deleteRecipeIngredient,
  } = useProductionApi()

  const { showNotice } = useSnackbar()

  const {
    formatGravity,
    formatPercent,
    formatWholeNumber: formatNumber,
    formatBatchSize,
  } = useBrewingFormatters()

  // Scaling
  const recipeBatchSize = computed(() => props.recipe.batch_size)
  const recipeBatchSizeUnit = computed(() => props.recipe.batch_size_unit)
  const {
    targetBatchSize,
    targetBatchSizeUnit,
    isScaling,
    scaleFactorDisplay,
    scaleAmount,
    resetScaling,
  } = useRecipeScaling(recipeBatchSize, recipeBatchSizeUnit)

  const showScalingControls = ref(false)
  const scalingVolumeUnits = ['bbl', 'gal', 'L', 'hL']

  function toggleScaling () {
    if (isScaling.value) {
      resetScaling()
      showScalingControls.value = false
    } else {
      // Initialize target to recipe's batch size
      targetBatchSize.value = props.recipe.batch_size
      targetBatchSizeUnit.value = props.recipe.batch_size_unit ?? 'bbl'
      showScalingControls.value = true
    }
  }

  // State
  const activeTab = ref('overview')
  const refreshing = ref(false)
  const ingredients = ref<RecipeIngredient[]>([])
  const ingredientsLoading = ref(false)
  const styles = ref<Style[]>([])
  const stylesLoading = ref(false)

  // Dialogs
  const editDialog = ref(false)
  const deleteDialog = ref(false)
  const ingredientDialog = ref(false)
  const ingredientDialogType = ref<RecipeIngredientType>('fermentable')
  const editingIngredient = ref<RecipeIngredient | null>(null)
  const deleteIngredientDialog = ref(false)
  const ingredientToDelete = ref<RecipeIngredient | null>(null)

  // Saving states
  const savingRecipe = ref(false)
  const deletingRecipe = ref(false)
  const savingIngredient = ref(false)
  const deletingIngredient = ref(false)
  const savingSpecs = ref(false)

  // Error states for dialogs
  const editRecipeError = ref('')
  const deleteRecipeError = ref('')

  // Display batch size: show scaled target when scaling is active, otherwise recipe value
  const displayBatchSize = computed(() => {
    if (isScaling.value && targetBatchSize.value !== null) {
      return formatBatchSize(targetBatchSize.value, targetBatchSizeUnit.value)
    }
    return formatBatchSize(props.recipe.batch_size, props.recipe.batch_size_unit)
  })

  // Computed ingredient lists
  const fermentables = computed(() =>
    ingredients.value.filter(i => i.ingredient_type === 'fermentable'),
  )

  const hops = computed(() =>
    ingredients.value.filter(i => i.ingredient_type === 'hop'),
  )

  const yeasts = computed(() =>
    ingredients.value.filter(i => i.ingredient_type === 'yeast'),
  )

  const adjuncts = computed(() =>
    ingredients.value.filter(i => i.ingredient_type === 'adjunct'),
  )

  const waterChemistry = computed(() =>
    ingredients.value.filter(i =>
      i.ingredient_type === 'salt' || i.ingredient_type === 'chemical',
    ),
  )

  // Lifecycle
  onMounted(async () => {
    await Promise.all([loadIngredients(), loadStyles()])
  })

  watch(() => props.recipe.uuid, async () => {
    await loadIngredients()
  })

  // Methods
  async function refresh () {
    refreshing.value = true
    try {
      const [updatedRecipe] = await Promise.all([
        getRecipe(props.recipe.uuid),
        loadIngredients(),
      ])
      emit('updated', updatedRecipe)
    } catch (error) {
      console.error('Failed to refresh:', error)
      showNotice('Failed to refresh recipe', 'error')
    } finally {
      refreshing.value = false
    }
  }

  async function loadIngredients () {
    ingredientsLoading.value = true
    try {
      ingredients.value = await getRecipeIngredients(props.recipe.uuid)
    } catch (error) {
      console.error('Failed to load ingredients:', error)
      showNotice('Failed to load ingredients', 'error')
    } finally {
      ingredientsLoading.value = false
    }
  }

  async function loadStyles () {
    stylesLoading.value = true
    try {
      styles.value = await getStyles()
    } catch (error) {
      console.error('Failed to load styles:', error)
    } finally {
      stylesLoading.value = false
    }
  }

  function openEditDialog () {
    editRecipeError.value = ''
    loadStyles()
    editDialog.value = true
  }

  function openDeleteDialog () {
    deleteRecipeError.value = ''
    deleteDialog.value = true
  }

  async function saveRecipe (data: UpdateRecipeRequest) {
    savingRecipe.value = true
    editRecipeError.value = ''
    try {
      const updated = await updateRecipe(props.recipe.uuid, data)
      showNotice('Recipe updated')
      editDialog.value = false
      emit('updated', updated)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update recipe'
      editRecipeError.value = message
    } finally {
      savingRecipe.value = false
    }
  }

  async function confirmDelete () {
    deletingRecipe.value = true
    deleteRecipeError.value = ''
    try {
      await deleteRecipe(props.recipe.uuid)
      showNotice('Recipe deleted')
      deleteDialog.value = false
      emit('deleted')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to delete recipe'
      deleteRecipeError.value = message
    } finally {
      deletingRecipe.value = false
    }
  }

  function openIngredientDialog (type: RecipeIngredientType) {
    ingredientDialogType.value = type
    editingIngredient.value = null
    ingredientDialog.value = true
  }

  function openEditIngredientDialog (ingredient: RecipeIngredient) {
    ingredientDialogType.value = ingredient.ingredient_type
    editingIngredient.value = ingredient
    ingredientDialog.value = true
  }

  async function saveIngredient (data: CreateRecipeIngredientRequest | UpdateRecipeIngredientRequest) {
    savingIngredient.value = true
    try {
      if (editingIngredient.value) {
        await updateRecipeIngredient(
          props.recipe.uuid,
          editingIngredient.value.uuid,
          data as UpdateRecipeIngredientRequest,
        )
        showNotice('Ingredient updated')
      } else {
        await createRecipeIngredient(props.recipe.uuid, data as CreateRecipeIngredientRequest)
        showNotice('Ingredient added')
      }
      ingredientDialog.value = false
      editingIngredient.value = null
      await loadIngredients()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save ingredient'
      showNotice(message, 'error')
    } finally {
      savingIngredient.value = false
    }
  }

  function deleteIngredient (ingredient: RecipeIngredient) {
    ingredientToDelete.value = ingredient
    deleteIngredientDialog.value = true
  }

  function cancelDeleteIngredient () {
    deleteIngredientDialog.value = false
    ingredientToDelete.value = null
  }

  async function confirmDeleteIngredient () {
    if (!ingredientToDelete.value) return
    deletingIngredient.value = true
    try {
      await deleteRecipeIngredient(props.recipe.uuid, ingredientToDelete.value.uuid)
      showNotice('Ingredient removed')
      deleteIngredientDialog.value = false
      ingredientToDelete.value = null
      await loadIngredients()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to remove ingredient'
      showNotice(message, 'error')
    } finally {
      deletingIngredient.value = false
    }
  }

  async function saveSpecs (data: UpdateRecipeRequest) {
    savingSpecs.value = true
    try {
      // Merge with recipe name, but data.name takes precedence if provided
      const payload: UpdateRecipeRequest = {
        ...data,
        name: data.name ?? props.recipe.name,
      }
      const updated = await updateRecipe(props.recipe.uuid, payload)
      showNotice('Specifications saved')
      emit('updated', updated)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save specifications'
      showNotice(message, 'error')
    } finally {
      savingSpecs.value = false
    }
  }

</script>

<style scoped>
.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.metric-card {
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(var(--v-theme-surface), 0.5);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  text-align: center;
}

.metric-label {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-bottom: 4px;
}

.metric-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.87);
}

.srm-swatch {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  border: 1px solid rgba(var(--v-theme-on-surface), 0.2);
}

.recipe-tabs {
  border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.12);
}

.scaling-controls {
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(var(--v-theme-primary), 0.06);
  border: 1px solid rgba(var(--v-theme-primary), 0.15);
}
</style>
