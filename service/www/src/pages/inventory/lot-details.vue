<template>
  <v-container class="inventory-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        Lot details
        <v-spacer />
        <v-btn :loading="loading" size="small" variant="text" @click="loadLots">
          Refresh
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-row align="stretch">
          <v-col cols="12" md="4">
            <v-card class="sub-card" variant="tonal">
              <v-card-title>Lot lookup</v-card-title>
              <v-card-text>
                <v-select
                  v-model="detailLotId"
                  :items="lotSelectItems"
                  label="Ingredient lot"
                />
                <v-btn
                  block
                  class="mb-2"
                  color="primary"
                  :disabled="!detailLotId"
                  @click="loadLotDetails"
                >
                  Load details
                </v-btn>
                <v-btn block variant="text" @click="clearLotDetails">
                  Clear selection
                </v-btn>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" md="8">
            <v-row>
              <v-col cols="12" md="6">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title>Malt lot detail</v-card-title>
                  <v-card-text>
                    <div v-if="lotMaltDetail" class="text-body-2 text-medium-emphasis">
                      Moisture {{ lotMaltDetail.moisture_percent ?? 'n/a' }}
                    </div>
                    <div v-if="!lotMaltDetail" class="text-body-2 text-medium-emphasis">
                      No malt lot detail loaded.
                    </div>
                    <v-divider class="my-3" />
                    <v-text-field v-model="lotMaltDetailForm.moisture_percent" label="Moisture percent" type="number" />
                    <v-btn
                      block
                      color="primary"
                      :disabled="!detailLotId"
                      @click="createLotMaltDetail"
                    >
                      Save malt lot detail
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="6">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title>Hop lot detail</v-card-title>
                  <v-card-text>
                    <div v-if="lotHopDetail" class="text-body-2 text-medium-emphasis">
                      Alpha {{ lotHopDetail.alpha_acid ?? 'n/a' }}
                    </div>
                    <div v-if="!lotHopDetail" class="text-body-2 text-medium-emphasis">
                      No hop lot detail loaded.
                    </div>
                    <v-divider class="my-3" />
                    <v-text-field v-model="lotHopDetailForm.alpha_acid" label="Alpha acid" type="number" />
                    <v-text-field v-model="lotHopDetailForm.beta_acid" label="Beta acid" type="number" />
                    <v-btn
                      block
                      color="primary"
                      :disabled="!detailLotId"
                      @click="createLotHopDetail"
                    >
                      Save hop lot detail
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12">
                <v-card class="sub-card" variant="outlined">
                  <v-card-title>Yeast lot detail</v-card-title>
                  <v-card-text>
                    <div v-if="lotYeastDetail" class="text-body-2 text-medium-emphasis">
                      Viability {{ lotYeastDetail.viability_percent ?? 'n/a' }}
                    </div>
                    <div v-if="!lotYeastDetail" class="text-body-2 text-medium-emphasis">
                      No yeast lot detail loaded.
                    </div>
                    <v-divider class="my-3" />
                    <v-text-field v-model="lotYeastDetailForm.viability_percent" label="Viability percent" type="number" />
                    <v-text-field v-model="lotYeastDetailForm.generation" label="Generation" type="number" />
                    <v-btn
                      block
                      color="primary"
                      :disabled="!detailLotId"
                      @click="createLotYeastDetail"
                    >
                      Save yeast lot detail
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
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  type Ingredient = {
    id: number
    name: string
  }

  type IngredientLot = {
    id: number
    ingredient_id: number
    received_amount: number
    received_unit: string
  }

  type IngredientLotMaltDetail = {
    id: number
    ingredient_lot_id: number
    moisture_percent: number | null
  }

  type IngredientLotHopDetail = {
    id: number
    ingredient_lot_id: number
    alpha_acid: number | null
    beta_acid: number | null
  }

  type IngredientLotYeastDetail = {
    id: number
    ingredient_lot_id: number
    viability_percent: number | null
    generation: number | null
  }

  const { request, toNumber } = useInventoryApi()
  const { formatAmountPreferred } = useUnitPreferences()
  const route = useRoute()

  const ingredients = ref<Ingredient[]>([])
  const lots = ref<IngredientLot[]>([])
  const detailLotId = ref<number | null>(null)
  const loading = ref(false)

  const lotMaltDetail = ref<IngredientLotMaltDetail | null>(null)
  const lotHopDetail = ref<IngredientLotHopDetail | null>(null)
  const lotYeastDetail = ref<IngredientLotYeastDetail | null>(null)

  const lotMaltDetailForm = reactive({
    moisture_percent: '',
  })

  const lotHopDetailForm = reactive({
    alpha_acid: '',
    beta_acid: '',
  })

  const lotYeastDetailForm = reactive({
    viability_percent: '',
    generation: '',
  })

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  const lotSelectItems = computed(() =>
    lots.value.map(lot => ({
      title: `${ingredientName(lot.ingredient_id)} (${formatAmountPreferred(lot.received_amount, lot.received_unit)})`,
      value: lot.id,
    })),
  )

  watch(detailLotId, () => {
    lotMaltDetail.value = null
    lotHopDetail.value = null
    lotYeastDetail.value = null
  })

  onMounted(async () => {
    await loadLots()
    const queryId = route.query.lot_id
    if (typeof queryId === 'string') {
      detailLotId.value = Number(queryId)
      await loadLotDetails()
    }
  })

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function loadLots () {
    loading.value = true
    try {
      const [ingredientData, lotData] = await Promise.all([
        request<Ingredient[]>('/ingredients'),
        request<IngredientLot[]>('/ingredient-lots'),
      ])
      ingredients.value = ingredientData
      lots.value = lotData
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load lots'
      showNotice(message, 'error')
    } finally {
      loading.value = false
    }
  }

  async function loadLotDetails () {
    if (!detailLotId.value) {
      return
    }
    try {
      lotMaltDetail.value = await request<IngredientLotMaltDetail>(
        `/ingredient-lot-malt-details?ingredient_lot_id=${detailLotId.value}`,
      )
    } catch {
      lotMaltDetail.value = null
    }
    try {
      lotHopDetail.value = await request<IngredientLotHopDetail>(
        `/ingredient-lot-hop-details?ingredient_lot_id=${detailLotId.value}`,
      )
    } catch {
      lotHopDetail.value = null
    }
    try {
      lotYeastDetail.value = await request<IngredientLotYeastDetail>(
        `/ingredient-lot-yeast-details?ingredient_lot_id=${detailLotId.value}`,
      )
    } catch {
      lotYeastDetail.value = null
    }
  }

  function clearLotDetails () {
    detailLotId.value = null
    lotMaltDetail.value = null
    lotHopDetail.value = null
    lotYeastDetail.value = null
  }

  async function createLotMaltDetail () {
    if (!detailLotId.value) {
      return
    }
    try {
      const payload = {
        ingredient_lot_id: detailLotId.value,
        moisture_percent: toNumber(lotMaltDetailForm.moisture_percent),
      }
      await request<IngredientLotMaltDetail>('/ingredient-lot-malt-details', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      showNotice('Malt lot detail saved')
      await loadLotDetails()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save malt lot detail'
      showNotice(message, 'error')
    }
  }

  async function createLotHopDetail () {
    if (!detailLotId.value) {
      return
    }
    try {
      const payload = {
        ingredient_lot_id: detailLotId.value,
        alpha_acid: toNumber(lotHopDetailForm.alpha_acid),
        beta_acid: toNumber(lotHopDetailForm.beta_acid),
      }
      await request<IngredientLotHopDetail>('/ingredient-lot-hop-details', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      showNotice('Hop lot detail saved')
      await loadLotDetails()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save hop lot detail'
      showNotice(message, 'error')
    }
  }

  async function createLotYeastDetail () {
    if (!detailLotId.value) {
      return
    }
    try {
      const payload = {
        ingredient_lot_id: detailLotId.value,
        viability_percent: toNumber(lotYeastDetailForm.viability_percent),
        generation: toNumber(lotYeastDetailForm.generation),
      }
      await request<IngredientLotYeastDetail>('/ingredient-lot-yeast-details', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
      showNotice('Yeast lot detail saved')
      await loadLotDetails()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save yeast lot detail'
      showNotice(message, 'error')
    }
  }

  function ingredientName (ingredientId: number) {
    return ingredients.value.find(ingredient => ingredient.id === ingredientId)?.name ?? `Ingredient ${ingredientId}`
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
