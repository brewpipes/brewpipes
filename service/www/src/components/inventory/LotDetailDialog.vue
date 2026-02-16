<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 560"
    :model-value="modelValue"
    persistent
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <span class="text-h6">Lot Details</span>
        <v-spacer />
        <v-btn
          :disabled="saving"
          icon="mdi-close"
          size="small"
          variant="text"
          @click="handleClose"
        />
      </v-card-title>

      <v-divider />

      <v-card-text class="pa-4">
        <!-- Loading state -->
        <template v-if="detailLoading">
          <div class="d-flex flex-column align-center py-8">
            <v-progress-circular color="primary" indeterminate size="48" />
            <p class="text-body-2 text-medium-emphasis mt-4">Loading lot details...</p>
          </div>
        </template>

        <template v-else-if="lot != null">
          <!-- Lot info summary -->
          <div class="text-body-2 text-medium-emphasis mb-4">
            <strong>{{ lotIngredientName }}</strong>
            <span v-if="lot.brewery_lot_code"> &mdash; Lot {{ lot.brewery_lot_code }}</span>
          </div>

          <!-- Error alert -->
          <v-alert
            v-if="errorMessage"
            class="mb-4"
            closable
            density="compact"
            type="error"
            variant="tonal"
            @click:close="errorMessage = ''"
          >
            {{ errorMessage }}
          </v-alert>

          <!-- Malt detail form: only shown for fermentable category -->
          <template v-if="lotCategory === 'fermentable'">
            <div class="text-subtitle-2 mb-2">Malt lot detail</div>
            <div v-if="lotMaltDetail != null" class="text-body-2 text-medium-emphasis mb-2">
              Moisture {{ lotMaltDetail.moisture_percent != null ? lotMaltDetail.moisture_percent + '%' : 'n/a' }}
            </div>
            <div v-else class="text-body-2 text-medium-emphasis mb-2">
              No malt lot detail loaded.
            </div>
            <v-divider class="my-3" />
            <v-text-field
              v-model="lotMaltDetailForm.moisture_percent"
              density="comfortable"
              label="Moisture percent"
              type="number"
            />
            <v-btn
              block
              color="primary"
              :loading="saving"
              @click="saveLotMaltDetail"
            >
              Save malt lot detail
            </v-btn>
          </template>

          <!-- Hop detail form: only shown for hop category -->
          <template v-else-if="lotCategory === 'hop'">
            <div class="text-subtitle-2 mb-2">Hop lot detail</div>
            <div v-if="lotHopDetail != null" class="text-body-2 text-medium-emphasis mb-2">
              Alpha {{ lotHopDetail.alpha_acid != null ? lotHopDetail.alpha_acid : 'n/a' }}
            </div>
            <div v-else class="text-body-2 text-medium-emphasis mb-2">
              No hop lot detail loaded.
            </div>
            <v-divider class="my-3" />
            <v-text-field
              v-model="lotHopDetailForm.alpha_acid"
              density="comfortable"
              label="Alpha acid"
              type="number"
            />
            <v-text-field
              v-model="lotHopDetailForm.beta_acid"
              density="comfortable"
              label="Beta acid"
              type="number"
            />
            <v-btn
              block
              color="primary"
              :loading="saving"
              @click="saveLotHopDetail"
            >
              Save hop lot detail
            </v-btn>
          </template>

          <!-- Yeast detail form: only shown for yeast category -->
          <template v-else-if="lotCategory === 'yeast'">
            <div class="text-subtitle-2 mb-2">Yeast lot detail</div>
            <div v-if="lotYeastDetail != null" class="text-body-2 text-medium-emphasis mb-2">
              Viability {{ lotYeastDetail.viability_percent != null ? lotYeastDetail.viability_percent + '%' : 'n/a' }}
            </div>
            <div v-else class="text-body-2 text-medium-emphasis mb-2">
              No yeast lot detail loaded.
            </div>
            <v-divider class="my-3" />
            <v-text-field
              v-model="lotYeastDetailForm.viability_percent"
              density="comfortable"
              label="Viability percent"
              type="number"
            />
            <v-text-field
              v-model="lotYeastDetailForm.generation"
              density="comfortable"
              label="Generation"
              type="number"
            />
            <v-btn
              block
              color="primary"
              :loading="saving"
              @click="saveLotYeastDetail"
            >
              Save yeast lot detail
            </v-btn>
          </template>

          <!-- No specialized details for this category -->
          <template v-else>
            <div class="text-body-2 text-medium-emphasis py-4">
              No specialized details for this category.
            </div>
          </template>
        </template>
      </v-card-text>

      <v-divider />

      <v-card-actions class="justify-end pa-4">
        <v-btn :disabled="saving" variant="text" @click="handleClose">
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type {
    Ingredient,
    IngredientLot,
    IngredientLotHopDetail,
    IngredientLotMaltDetail,
    IngredientLotYeastDetail,
  } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { toNumber } from '@/utils/normalize'

  const props = defineProps<{
    modelValue: boolean
    lot: IngredientLot | null
    ingredients: Ingredient[]
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'saved': []
  }>()

  const {
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

  const { execute: executeLoad, loading: detailLoading } = useAsyncAction()
  const { execute: executeSave, loading: saving } = useAsyncAction({
    onError: (message) => { errorMessage.value = message },
  })

  const errorMessage = ref('')

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

  /** Resolve the ingredient for the current lot */
  const lotIngredient = computed(() => {
    if (props.lot == null) return null
    return props.ingredients.find(i => i.uuid === props.lot!.ingredient_uuid) ?? null
  })

  const lotCategory = computed<string | null>(() => lotIngredient.value?.category ?? null)
  const lotIngredientName = computed(() => lotIngredient.value?.name ?? 'Unknown Ingredient')

  // Watch for dialog open to load details
  watch(() => props.modelValue, async open => {
    if (open && props.lot != null) {
      resetAll()
      await loadLotDetails()
    }
  })

  async function loadLotDetails () {
    if (props.lot == null) return

    const lotUuid = props.lot.uuid
    const category = lotCategory.value

    await executeLoad(async () => {
      if (category === 'fermentable') {
        try {
          lotMaltDetail.value = await getIngredientLotMaltDetail(lotUuid)
          populateMaltForm(lotMaltDetail.value)
        } catch {
          lotMaltDetail.value = null
          resetMaltForm()
        }
      } else if (category === 'hop') {
        try {
          lotHopDetail.value = await getIngredientLotHopDetail(lotUuid)
          populateHopForm(lotHopDetail.value)
        } catch {
          lotHopDetail.value = null
          resetHopForm()
        }
      } else if (category === 'yeast') {
        try {
          lotYeastDetail.value = await getIngredientLotYeastDetail(lotUuid)
          populateYeastForm(lotYeastDetail.value)
        } catch {
          lotYeastDetail.value = null
          resetYeastForm()
        }
      }
    })
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

  function resetAll () {
    lotMaltDetail.value = null
    lotHopDetail.value = null
    lotYeastDetail.value = null
    errorMessage.value = ''
    resetMaltForm()
    resetHopForm()
    resetYeastForm()
  }

  // --- Save handlers (create or update based on existing detail) ---

  async function saveLotMaltDetail () {
    if (props.lot == null) return
    await executeSave(async () => {
      const payload = {
        ingredient_lot_uuid: props.lot!.uuid,
        moisture_percent: toNumber(lotMaltDetailForm.moisture_percent),
      }
      if (lotMaltDetail.value != null) {
        await updateIngredientLotMaltDetail(lotMaltDetail.value.uuid, payload)
      } else {
        await createIngredientLotMaltDetail(payload)
      }
      showNotice('Malt lot detail saved')
      await loadLotDetails()
      emit('saved')
    })
  }

  async function saveLotHopDetail () {
    if (props.lot == null) return
    await executeSave(async () => {
      const payload = {
        ingredient_lot_uuid: props.lot!.uuid,
        alpha_acid: toNumber(lotHopDetailForm.alpha_acid),
        beta_acid: toNumber(lotHopDetailForm.beta_acid),
      }
      if (lotHopDetail.value != null) {
        await updateIngredientLotHopDetail(lotHopDetail.value.uuid, payload)
      } else {
        await createIngredientLotHopDetail(payload)
      }
      showNotice('Hop lot detail saved')
      await loadLotDetails()
      emit('saved')
    })
  }

  async function saveLotYeastDetail () {
    if (props.lot == null) return
    await executeSave(async () => {
      const payload = {
        ingredient_lot_uuid: props.lot!.uuid,
        viability_percent: toNumber(lotYeastDetailForm.viability_percent),
        generation: toNumber(lotYeastDetailForm.generation),
      }
      if (lotYeastDetail.value != null) {
        await updateIngredientLotYeastDetail(lotYeastDetail.value.uuid, payload)
      } else {
        await createIngredientLotYeastDetail(payload)
      }
      showNotice('Yeast lot detail saved')
      await loadLotDetails()
      emit('saved')
    })
  }

  function handleClose () {
    emit('update:modelValue', false)
  }
</script>
