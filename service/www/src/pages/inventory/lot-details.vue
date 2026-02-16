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
                  v-model="detailLotUuid"
                  :items="lotSelectItems"
                  label="Ingredient lot"
                />
                <v-btn
                  block
                  class="mb-2"
                  color="primary"
                  :disabled="!detailLotUuid"
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
              <!-- Malt detail form: only shown for fermentable category -->
              <v-col v-if="selectedLotCategory === 'fermentable'" cols="12" md="6">
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
                      :disabled="!detailLotUuid"
                      @click="saveLotMaltDetail"
                    >
                      Save malt lot detail
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- Hop detail form: only shown for hop category -->
              <v-col v-if="selectedLotCategory === 'hop'" cols="12" md="6">
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
                      :disabled="!detailLotUuid"
                      @click="saveLotHopDetail"
                    >
                      Save hop lot detail
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- Yeast detail form: only shown for yeast category -->
              <v-col v-if="selectedLotCategory === 'yeast'" cols="12">
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
                      :disabled="!detailLotUuid"
                      @click="saveLotYeastDetail"
                    >
                      Save yeast lot detail
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- Message when no lot is selected or category has no specialized form -->
              <v-col v-if="!detailLotUuid" cols="12">
                <div class="text-body-2 text-medium-emphasis">
                  Select a lot and click "Load details" to view category-specific details.
                </div>
              </v-col>
              <v-col v-else-if="selectedLotCategory && !['fermentable', 'hop', 'yeast'].includes(selectedLotCategory)" cols="12">
                <div class="text-body-2 text-medium-emphasis">
                  No specialized detail form for category "{{ selectedLotCategory }}".
                </div>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>
  </v-container>

</template>

<script lang="ts" setup>
  import type { Ingredient, IngredientLot, IngredientLotHopDetail, IngredientLotMaltDetail, IngredientLotYeastDetail } from '@/types'
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { toNumber } from '@/utils/normalize'

  const {
    getIngredients,
    getIngredientLots,
    getIngredientLotMaltDetail,
    createIngredientLotMaltDetail,
    updateIngredientLotMaltDetail,
    getIngredientLotHopDetail,
    createIngredientLotHopDetail,
    updateIngredientLotHopDetail,
    getIngredientLotYeastDetail,
    createIngredientLotYeastDetail,
    updateIngredientLotYeastDetail,
  } = useInventoryApi()
  const { showNotice } = useSnackbar()
  const { formatAmountPreferred } = useUnitPreferences()
  const route = useRoute()

  const ingredients = ref<Ingredient[]>([])
  const lots = ref<IngredientLot[]>([])
  const detailLotUuid = ref<string | null>(null)

  const { execute, loading } = useAsyncAction({
    onError: (message) => showNotice(message, 'error'),
  })
  const { execute: executeSave } = useAsyncAction({
    onError: (message) => showNotice(message, 'error'),
  })

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

  /** Resolve the ingredient category for the currently selected lot */
  const selectedLotCategory = computed<string | null>(() => {
    if (!detailLotUuid.value) return null
    const lot = lots.value.find(l => l.uuid === detailLotUuid.value)
    if (!lot) return null
    const ingredient = ingredients.value.find(i => i.uuid === lot.ingredient_uuid)
    return ingredient?.category ?? null
  })

  const lotSelectItems = computed(() =>
    lots.value.map(lot => ({
      title: `${ingredientName(lot.ingredient_uuid)} (${formatAmountPreferred(lot.received_amount, lot.received_unit)})`,
      value: lot.uuid,
    })),
  )

  watch(detailLotUuid, () => {
    lotMaltDetail.value = null
    lotHopDetail.value = null
    lotYeastDetail.value = null
    resetForms()
  })

  onMounted(async () => {
    await loadLots()
    const queryUuid = route.query.lot_uuid
    if (typeof queryUuid === 'string') {
      detailLotUuid.value = queryUuid
      await loadLotDetails()
    }
  })

  async function loadLots () {
    await execute(async () => {
      const [ingredientData, lotData] = await Promise.all([
        getIngredients(),
        getIngredientLots(),
      ])
      ingredients.value = ingredientData
      lots.value = lotData
    })
  }

  async function loadLotDetails () {
    if (!detailLotUuid.value) {
      return
    }

    const category = selectedLotCategory.value

    // Only fetch the detail relevant to this lot's ingredient category
    if (category === 'fermentable') {
      try {
        lotMaltDetail.value = await getIngredientLotMaltDetail(detailLotUuid.value)
        populateMaltForm(lotMaltDetail.value)
      } catch {
        lotMaltDetail.value = null
        resetMaltForm()
      }
    } else if (category === 'hop') {
      try {
        lotHopDetail.value = await getIngredientLotHopDetail(detailLotUuid.value)
        populateHopForm(lotHopDetail.value)
      } catch {
        lotHopDetail.value = null
        resetHopForm()
      }
    } else if (category === 'yeast') {
      try {
        lotYeastDetail.value = await getIngredientLotYeastDetail(detailLotUuid.value)
        populateYeastForm(lotYeastDetail.value)
      } catch {
        lotYeastDetail.value = null
        resetYeastForm()
      }
    }
  }

  function clearLotDetails () {
    detailLotUuid.value = null
    lotMaltDetail.value = null
    lotHopDetail.value = null
    lotYeastDetail.value = null
    resetForms()
  }

  // --- Form population helpers ---

  function populateMaltForm (detail: IngredientLotMaltDetail) {
    lotMaltDetailForm.moisture_percent = detail.moisture_percent != null ? String(detail.moisture_percent) : ''
  }

  function populateHopForm (detail: IngredientLotHopDetail) {
    lotHopDetailForm.alpha_acid = detail.alpha_acid != null ? String(detail.alpha_acid) : ''
    lotHopDetailForm.beta_acid = detail.beta_acid != null ? String(detail.beta_acid) : ''
  }

  function populateYeastForm (detail: IngredientLotYeastDetail) {
    lotYeastDetailForm.viability_percent = detail.viability_percent != null ? String(detail.viability_percent) : ''
    lotYeastDetailForm.generation = detail.generation != null ? String(detail.generation) : ''
  }

  function resetMaltForm () {
    lotMaltDetailForm.moisture_percent = ''
  }

  function resetHopForm () {
    lotHopDetailForm.alpha_acid = ''
    lotHopDetailForm.beta_acid = ''
  }

  function resetYeastForm () {
    lotYeastDetailForm.viability_percent = ''
    lotYeastDetailForm.generation = ''
  }

  function resetForms () {
    resetMaltForm()
    resetHopForm()
    resetYeastForm()
  }

  // --- Save handlers (create or update based on existing detail) ---

  async function saveLotMaltDetail () {
    if (!detailLotUuid.value) {
      return
    }
    await executeSave(async () => {
      const payload = {
        ingredient_lot_uuid: detailLotUuid.value!,
        moisture_percent: toNumber(lotMaltDetailForm.moisture_percent),
      }
      if (lotMaltDetail.value) {
        await updateIngredientLotMaltDetail(lotMaltDetail.value.uuid, payload)
      } else {
        await createIngredientLotMaltDetail(payload)
      }
      showNotice('Malt lot detail saved')
      await loadLotDetails()
    })
  }

  async function saveLotHopDetail () {
    if (!detailLotUuid.value) {
      return
    }
    await executeSave(async () => {
      const payload = {
        ingredient_lot_uuid: detailLotUuid.value!,
        alpha_acid: toNumber(lotHopDetailForm.alpha_acid),
        beta_acid: toNumber(lotHopDetailForm.beta_acid),
      }
      if (lotHopDetail.value) {
        await updateIngredientLotHopDetail(lotHopDetail.value.uuid, payload)
      } else {
        await createIngredientLotHopDetail(payload)
      }
      showNotice('Hop lot detail saved')
      await loadLotDetails()
    })
  }

  async function saveLotYeastDetail () {
    if (!detailLotUuid.value) {
      return
    }
    await executeSave(async () => {
      const payload = {
        ingredient_lot_uuid: detailLotUuid.value!,
        viability_percent: toNumber(lotYeastDetailForm.viability_percent),
        generation: toNumber(lotYeastDetailForm.generation),
      }
      if (lotYeastDetail.value) {
        await updateIngredientLotYeastDetail(lotYeastDetail.value.uuid, payload)
      } else {
        await createIngredientLotYeastDetail(payload)
      }
      showNotice('Yeast lot detail saved')
      await loadLotDetails()
    })
  }

  function ingredientName (ingredientUuid: string) {
    return ingredients.value.find(ingredient => ingredient.uuid === ingredientUuid)?.name ?? 'Unknown Ingredient'
  }
</script>

<style scoped>
.inventory-page {
  position: relative;
}
</style>
