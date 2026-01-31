<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Ingredients
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="loadIngredients">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="7">
            <v-card class="sub-card" variant="outlined">
              <v-card-title>Ingredient list</v-card-title>
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
                <v-table class="data-table" density="compact">
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Category</th>
                      <th>Unit</th>
                      <th>Updated</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="ingredient in ingredients" :key="ingredient.id">
                      <td>{{ ingredient.name }}</td>
                      <td>{{ ingredient.category }}</td>
                      <td>{{ ingredient.default_unit }}</td>
                      <td>{{ formatDateTime(ingredient.updated_at) }}</td>
                      <td>
                        <v-btn
                          size="x-small"
                          variant="text"
                          @click="openDetails(ingredient.id)"
                        >
                          Details
                        </v-btn>
                      </td>
                    </tr>
                    <tr v-if="ingredients.length === 0">
                      <td colspan="5">No ingredients yet.</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="5">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Create ingredient</v-card-title>
              <v-card-text>
                <v-text-field v-model="ingredientForm.name" label="Name" />
                <v-combobox
                  v-model="ingredientForm.category"
                  :items="ingredientCategoryOptions"
                  label="Category"
                />
                <v-combobox
                  v-model="ingredientForm.default_unit"
                  :items="unitOptions"
                  label="Default unit"
                />
                <v-textarea
                  v-model="ingredientForm.description"
                  auto-grow
                  label="Description"
                  rows="2"
                />
                <v-btn
                  block
                  color="primary"
                  :disabled="!ingredientForm.name.trim() || !ingredientForm.category || !ingredientForm.default_unit"
                  @click="createIngredient"
                >
                  Add ingredient
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useInventoryApi } from '@/composables/useInventoryApi'

type Ingredient = {
  id: number
  uuid: string
  name: string
  category: string
  default_unit: string
  description: string
  created_at: string
  updated_at: string
}

const { request, normalizeText, formatDateTime } = useInventoryApi()
const router = useRouter()

const ingredients = ref<Ingredient[]>([])
const errorMessage = ref('')
const loading = ref(false)

const ingredientCategoryOptions = ['malt', 'hop', 'yeast', 'adjunct', 'water_chem', 'gas', 'other']
const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']

const ingredientForm = reactive({
  name: '',
  category: '',
  default_unit: '',
  description: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

onMounted(async () => {
  await loadIngredients()
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function loadIngredients() {
  loading.value = true
  errorMessage.value = ''
  try {
    ingredients.value = await request<Ingredient[]>('/ingredients')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load ingredients'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function createIngredient() {
  try {
    const payload = {
      name: ingredientForm.name.trim(),
      category: ingredientForm.category,
      default_unit: ingredientForm.default_unit,
      description: normalizeText(ingredientForm.description),
    }
    await request<Ingredient>('/ingredients', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    ingredientForm.name = ''
    ingredientForm.category = ''
    ingredientForm.default_unit = ''
    ingredientForm.description = ''
    await loadIngredients()
    showNotice('Ingredient created')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create ingredient'
    errorMessage.value = message
    showNotice(message, 'error')
  }
}

function openDetails(ingredientId: number) {
  router.push({
    path: '/inventory/ingredient-details',
    query: { ingredient_id: String(ingredientId) },
  })
}
</script>

<style scoped>
.inventory-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
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
