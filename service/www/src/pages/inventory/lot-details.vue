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
                      :disabled="!detailLotUuid"
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
                      :disabled="!detailLotUuid"
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
                      :disabled="!detailLotUuid"
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

</template>

<script lang="ts" setup>
  import type { Ingredient, IngredientLot, IngredientLotHopDetail, IngredientLotMaltDetail, IngredientLotYeastDetail } from '@/types'
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { toNumber } from '@/utils/normalize'

  const {
    getIngredients,
    getIngredientLots,
    getIngredientLotMaltDetail,
    createIngredientLotMaltDetail,
    getIngredientLotHopDetail,
    createIngredientLotHopDetail,
    getIngredientLotYeastDetail,
    createIngredientLotYeastDetail,
  } = useInventoryApi()
  const { showNotice } = useSnackbar()
  const { formatAmountPreferred } = useUnitPreferences()
  const route = useRoute()

  const ingredients = ref<Ingredient[]>([])
  const lots = ref<IngredientLot[]>([])
  const detailLotUuid = ref<string | null>(null)
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
    loading.value = true
    try {
      const [ingredientData, lotData] = await Promise.all([
        getIngredients(),
        getIngredientLots(),
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
    if (!detailLotUuid.value) {
      return
    }
    try {
      lotMaltDetail.value = await getIngredientLotMaltDetail(detailLotUuid.value)
    } catch {
      lotMaltDetail.value = null
    }
    try {
      lotHopDetail.value = await getIngredientLotHopDetail(detailLotUuid.value)
    } catch {
      lotHopDetail.value = null
    }
    try {
      lotYeastDetail.value = await getIngredientLotYeastDetail(detailLotUuid.value)
    } catch {
      lotYeastDetail.value = null
    }
  }

  function clearLotDetails () {
    detailLotUuid.value = null
    lotMaltDetail.value = null
    lotHopDetail.value = null
    lotYeastDetail.value = null
  }

  async function createLotMaltDetail () {
    if (!detailLotUuid.value) {
      return
    }
    try {
      const payload = {
        ingredient_lot_uuid: detailLotUuid.value,
        moisture_percent: toNumber(lotMaltDetailForm.moisture_percent),
      }
      await createIngredientLotMaltDetail(payload)
      showNotice('Malt lot detail saved')
      await loadLotDetails()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save malt lot detail'
      showNotice(message, 'error')
    }
  }

  async function createLotHopDetail () {
    if (!detailLotUuid.value) {
      return
    }
    try {
      const payload = {
        ingredient_lot_uuid: detailLotUuid.value,
        alpha_acid: toNumber(lotHopDetailForm.alpha_acid),
        beta_acid: toNumber(lotHopDetailForm.beta_acid),
      }
      await createIngredientLotHopDetail(payload)
      showNotice('Hop lot detail saved')
      await loadLotDetails()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save hop lot detail'
      showNotice(message, 'error')
    }
  }

  async function createLotYeastDetail () {
    if (!detailLotUuid.value) {
      return
    }
    try {
      const payload = {
        ingredient_lot_uuid: detailLotUuid.value,
        viability_percent: toNumber(lotYeastDetailForm.viability_percent),
        generation: toNumber(lotYeastDetailForm.generation),
      }
      await createIngredientLotYeastDetail(payload)
      showNotice('Yeast lot detail saved')
      await loadLotDetails()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to save yeast lot detail'
      showNotice(message, 'error')
    }
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
