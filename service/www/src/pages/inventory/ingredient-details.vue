<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Ingredient details
        <v-spacer />
        <v-btn size="small" variant="text" :loading="loading" @click="loadIngredients">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="4">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Ingredient lookup</v-card-title>
              <v-card-text>
                <v-select
                  v-model="detailIngredientId"
                  :items="ingredientSelectItems"
                  label="Ingredient"
                />
                <v-btn
                  block
                  color="primary"
                  class="mb-2"
                  :disabled="!detailIngredientId"
                  @click="loadIngredientDetails"
                >
                  Load details
                </v-btn>
                <v-btn block variant="text" @click="clearIngredientDetails">
                  Clear selection
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="8">
            <v-row>
              <v-col cols="12" md="6">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title>Malt detail</v-card-title>
                  <v-card-text>
                    <div class="text-body-2 text-medium-emphasis" v-if="maltDetail">
                      Maltster {{ maltDetail.maltster_name || 'n/a' }}
                    </div>
                    <div class="text-body-2 text-medium-emphasis" v-if="maltDetail">
                      Lovibond {{ maltDetail.lovibond ?? 'n/a' }}
                    </div>
                    <div class="text-body-2 text-medium-emphasis" v-if="!maltDetail">
                      No malt detail loaded.
                    </div>
                    <v-divider class="my-3" />
                    <v-text-field v-model="maltDetailForm.maltster_name" label="Maltster name" />
                    <v-text-field v-model="maltDetailForm.variety" label="Variety" />
                    <v-text-field v-model="maltDetailForm.lovibond" label="Lovibond" type="number" />
                    <v-text-field v-model="maltDetailForm.srm" label="SRM" type="number" />
                    <v-text-field v-model="maltDetailForm.diastatic_power" label="Diastatic power" type="number" />
                    <v-btn
                      block
                      color="primary"
                      :disabled="!detailIngredientId"
                      @click="createMaltDetail"
                    >
                      Save malt detail
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="6">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title>Hop detail</v-card-title>
                  <v-card-text>
                    <div class="text-body-2 text-medium-emphasis" v-if="hopDetail">
                      Producer {{ hopDetail.producer_name || 'n/a' }}
                    </div>
                    <div class="text-body-2 text-medium-emphasis" v-if="hopDetail">
                      Alpha {{ hopDetail.alpha_acid ?? 'n/a' }}
                    </div>
                    <div class="text-body-2 text-medium-emphasis" v-if="!hopDetail">
                      No hop detail loaded.
                    </div>
                    <v-divider class="my-3" />
                    <v-text-field v-model="hopDetailForm.producer_name" label="Producer name" />
                    <v-text-field v-model="hopDetailForm.variety" label="Variety" />
                    <v-text-field v-model="hopDetailForm.crop_year" label="Crop year" type="number" />
                    <v-text-field v-model="hopDetailForm.form" label="Form" />
                    <v-text-field v-model="hopDetailForm.alpha_acid" label="Alpha acid" type="number" />
                    <v-text-field v-model="hopDetailForm.beta_acid" label="Beta acid" type="number" />
                    <v-btn
                      block
                      color="primary"
                      :disabled="!detailIngredientId"
                      @click="createHopDetail"
                    >
                      Save hop detail
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title>Yeast detail</v-card-title>
                  <v-card-text>
                    <div class="text-body-2 text-medium-emphasis" v-if="yeastDetail">
                      Lab {{ yeastDetail.lab_name || 'n/a' }}
                    </div>
                    <div class="text-body-2 text-medium-emphasis" v-if="yeastDetail">
                      Strain {{ yeastDetail.strain || 'n/a' }}
                    </div>
                    <div class="text-body-2 text-medium-emphasis" v-if="!yeastDetail">
                      No yeast detail loaded.
                    </div>
                    <v-divider class="my-3" />
                    <v-text-field v-model="yeastDetailForm.lab_name" label="Lab name" />
                    <v-text-field v-model="yeastDetailForm.strain" label="Strain" />
                    <v-text-field v-model="yeastDetailForm.form" label="Form" />
                    <v-btn
                      block
                      color="primary"
                      :disabled="!detailIngredientId"
                      @click="createYeastDetail"
                    >
                      Save yeast detail
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
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
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useInventoryApi } from '@/composables/useInventoryApi'

type Ingredient = {
  id: number
  name: string
  category: string
}

type IngredientMaltDetail = {
  id: number
  ingredient_id: number
  maltster_name: string
  variety: string
  lovibond: number | null
  srm: number | null
  diastatic_power: number | null
}

type IngredientHopDetail = {
  id: number
  ingredient_id: number
  producer_name: string
  variety: string
  crop_year: number | null
  form: string
  alpha_acid: number | null
  beta_acid: number | null
}

type IngredientYeastDetail = {
  id: number
  ingredient_id: number
  lab_name: string
  strain: string
  form: string
}

const { request, normalizeText, toNumber } = useInventoryApi()
const route = useRoute()

const ingredients = ref<Ingredient[]>([])
const detailIngredientId = ref<number | null>(null)
const loading = ref(false)

const maltDetail = ref<IngredientMaltDetail | null>(null)
const hopDetail = ref<IngredientHopDetail | null>(null)
const yeastDetail = ref<IngredientYeastDetail | null>(null)

const maltDetailForm = reactive({
  maltster_name: '',
  variety: '',
  lovibond: '',
  srm: '',
  diastatic_power: '',
})

const hopDetailForm = reactive({
  producer_name: '',
  variety: '',
  crop_year: '',
  form: '',
  alpha_acid: '',
  beta_acid: '',
})

const yeastDetailForm = reactive({
  lab_name: '',
  strain: '',
  form: '',
})

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

const ingredientSelectItems = computed(() =>
  ingredients.value.map((ingredient) => ({
    title: `${ingredient.name} (${ingredient.category})`,
    value: ingredient.id,
  })),
)

watch(detailIngredientId, () => {
  maltDetail.value = null
  hopDetail.value = null
  yeastDetail.value = null
})

onMounted(async () => {
  await loadIngredients()
  const queryId = route.query.ingredient_id
  if (typeof queryId === 'string') {
    detailIngredientId.value = Number(queryId)
    await loadIngredientDetails()
  }
})

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function loadIngredients() {
  loading.value = true
  try {
    ingredients.value = await request<Ingredient[]>('/ingredients')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load ingredients'
    showNotice(message, 'error')
  } finally {
    loading.value = false
  }
}

async function loadIngredientDetails() {
  if (!detailIngredientId.value) {
    return
  }
  try {
    maltDetail.value = await request<IngredientMaltDetail>(
      `/ingredient-malt-details?ingredient_id=${detailIngredientId.value}`,
    )
  } catch {
    maltDetail.value = null
  }
  try {
    hopDetail.value = await request<IngredientHopDetail>(
      `/ingredient-hop-details?ingredient_id=${detailIngredientId.value}`,
    )
  } catch {
    hopDetail.value = null
  }
  try {
    yeastDetail.value = await request<IngredientYeastDetail>(
      `/ingredient-yeast-details?ingredient_id=${detailIngredientId.value}`,
    )
  } catch {
    yeastDetail.value = null
  }
}

function clearIngredientDetails() {
  detailIngredientId.value = null
  maltDetail.value = null
  hopDetail.value = null
  yeastDetail.value = null
}

async function createMaltDetail() {
  if (!detailIngredientId.value) {
    return
  }
  try {
    const payload = {
      ingredient_id: detailIngredientId.value,
      maltster_name: normalizeText(maltDetailForm.maltster_name),
      variety: normalizeText(maltDetailForm.variety),
      lovibond: toNumber(maltDetailForm.lovibond),
      srm: toNumber(maltDetailForm.srm),
      diastatic_power: toNumber(maltDetailForm.diastatic_power),
    }
    await request<IngredientMaltDetail>('/ingredient-malt-details', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    showNotice('Malt detail saved')
    await loadIngredientDetails()
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to save malt detail'
    showNotice(message, 'error')
  }
}

async function createHopDetail() {
  if (!detailIngredientId.value) {
    return
  }
  try {
    const payload = {
      ingredient_id: detailIngredientId.value,
      producer_name: normalizeText(hopDetailForm.producer_name),
      variety: normalizeText(hopDetailForm.variety),
      crop_year: toNumber(hopDetailForm.crop_year),
      form: normalizeText(hopDetailForm.form),
      alpha_acid: toNumber(hopDetailForm.alpha_acid),
      beta_acid: toNumber(hopDetailForm.beta_acid),
    }
    await request<IngredientHopDetail>('/ingredient-hop-details', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    showNotice('Hop detail saved')
    await loadIngredientDetails()
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to save hop detail'
    showNotice(message, 'error')
  }
}

async function createYeastDetail() {
  if (!detailIngredientId.value) {
    return
  }
  try {
    const payload = {
      ingredient_id: detailIngredientId.value,
      lab_name: normalizeText(yeastDetailForm.lab_name),
      strain: normalizeText(yeastDetailForm.strain),
      form: normalizeText(yeastDetailForm.form),
    }
    await request<IngredientYeastDetail>('/ingredient-yeast-details', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    showNotice('Yeast detail saved')
    await loadIngredientDetails()
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to save yeast detail'
    showNotice(message, 'error')
  }
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
</style>
