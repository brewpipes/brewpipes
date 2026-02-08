<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-title class="text-subtitle-1">
      Target Specifications
    </v-card-title>
    <v-card-text>
      <v-form ref="formRef" @submit.prevent="handleSubmit">
        <!-- Batch Size -->
        <div class="text-overline text-medium-emphasis mb-2">Batch Size</div>
        <v-row dense>
          <v-col cols="8" md="4">
            <v-text-field
              v-model.number="form.batch_size"
              density="compact"
              hide-details
              label="Size"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="4" md="2">
            <v-select
              v-model="form.batch_size_unit"
              density="compact"
              hide-details
              :items="volumeUnits"
              label="Unit"
              variant="outlined"
            />
          </v-col>
        </v-row>

        <v-divider class="my-4" />

        <!-- Gravity -->
        <div class="text-overline text-medium-emphasis mb-2">Gravity</div>
        <v-row dense>
          <v-col cols="12" md="4">
            <v-text-field
              v-model.number="form.target_og"
              density="compact"
              hide-details
              label="Target OG"
              placeholder="1.050"
              step="0.001"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="6" md="4">
            <v-text-field
              v-model.number="form.target_og_min"
              density="compact"
              hide-details
              label="OG Min"
              placeholder="1.048"
              step="0.001"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="6" md="4">
            <v-text-field
              v-model.number="form.target_og_max"
              density="compact"
              hide-details
              label="OG Max"
              placeholder="1.052"
              step="0.001"
              type="number"
              variant="outlined"
            />
          </v-col>
        </v-row>

        <v-row class="mt-2" dense>
          <v-col cols="12" md="4">
            <v-text-field
              v-model.number="form.target_fg"
              density="compact"
              hide-details
              label="Target FG"
              placeholder="1.010"
              step="0.001"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="6" md="4">
            <v-text-field
              v-model.number="form.target_fg_min"
              density="compact"
              hide-details
              label="FG Min"
              placeholder="1.008"
              step="0.001"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="6" md="4">
            <v-text-field
              v-model.number="form.target_fg_max"
              density="compact"
              hide-details
              label="FG Max"
              placeholder="1.012"
              step="0.001"
              type="number"
              variant="outlined"
            />
          </v-col>
        </v-row>

        <v-divider class="my-4" />

        <!-- IBU -->
        <div class="text-overline text-medium-emphasis mb-2">Bitterness (IBU)</div>
        <v-row dense>
          <v-col cols="12" md="3">
            <v-text-field
              v-model.number="form.target_ibu"
              density="compact"
              hide-details
              label="Target IBU"
              placeholder="40"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="6" md="3">
            <v-text-field
              v-model.number="form.target_ibu_min"
              density="compact"
              hide-details
              label="IBU Min"
              placeholder="35"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="6" md="3">
            <v-text-field
              v-model.number="form.target_ibu_max"
              density="compact"
              hide-details
              label="IBU Max"
              placeholder="45"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-select
              v-model="form.ibu_method"
              clearable
              density="compact"
              hide-details
              :items="ibuMethods"
              label="Calculation Method"
              variant="outlined"
            />
          </v-col>
        </v-row>

        <v-divider class="my-4" />

        <!-- Color -->
        <div class="text-overline text-medium-emphasis mb-2">Color (SRM)</div>
        <v-row dense>
          <v-col cols="12" md="4">
            <v-text-field
              v-model.number="form.target_srm"
              density="compact"
              hide-details
              label="Target SRM"
              placeholder="10"
              type="number"
              variant="outlined"
            >
              <template #append-inner>
                <div
                  v-if="form.target_srm"
                  class="srm-swatch"
                  :style="{ backgroundColor: srmToColor(form.target_srm) }"
                />
              </template>
            </v-text-field>
          </v-col>
          <v-col cols="6" md="4">
            <v-text-field
              v-model.number="form.target_srm_min"
              density="compact"
              hide-details
              label="SRM Min"
              placeholder="8"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="6" md="4">
            <v-text-field
              v-model.number="form.target_srm_max"
              density="compact"
              hide-details
              label="SRM Max"
              placeholder="12"
              type="number"
              variant="outlined"
            />
          </v-col>
        </v-row>

        <v-divider class="my-4" />

        <!-- Other -->
        <div class="text-overline text-medium-emphasis mb-2">Other</div>
        <v-row dense>
          <v-col cols="6" md="4">
            <v-text-field
              v-model.number="form.target_carbonation"
              density="compact"
              hide-details
              label="Carbonation (vol CO2)"
              placeholder="2.4"
              step="0.1"
              type="number"
              variant="outlined"
            />
          </v-col>
          <v-col cols="6" md="4">
            <v-text-field
              v-model.number="form.brewhouse_efficiency"
              density="compact"
              hide-details
              label="Brewhouse Efficiency %"
              placeholder="75"
              type="number"
              variant="outlined"
            />
          </v-col>
        </v-row>

        <!-- Calculated ABV Display -->
        <v-alert
          v-if="calculatedAbv !== null"
          class="mt-4"
          density="compact"
          type="info"
          variant="tonal"
        >
          <strong>Calculated ABV:</strong> {{ calculatedAbv.toFixed(1) }}%
          <span class="text-medium-emphasis">(based on target OG and FG)</span>
        </v-alert>

        <div class="d-flex justify-end mt-4">
          <v-btn
            color="primary"
            :disabled="!hasChanges"
            :loading="saving"
            type="submit"
          >
            Save Specifications
          </v-btn>
        </div>
      </v-form>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { Recipe, UpdateRecipeRequest } from '@/composables/useProductionApi'
  import { computed, reactive, ref, watch } from 'vue'
  import { srmToColor } from '@/composables/useFormatters'

  const props = defineProps<{
    recipe: Recipe
    saving: boolean
  }>()

  const emit = defineEmits<{
    save: [data: UpdateRecipeRequest]
  }>()

  const formRef = ref()

  const form = reactive({
    batch_size: props.recipe.batch_size,
    batch_size_unit: props.recipe.batch_size_unit ?? 'bbl',
    target_og: props.recipe.target_og,
    target_og_min: props.recipe.target_og_min,
    target_og_max: props.recipe.target_og_max,
    target_fg: props.recipe.target_fg,
    target_fg_min: props.recipe.target_fg_min,
    target_fg_max: props.recipe.target_fg_max,
    target_ibu: props.recipe.target_ibu,
    target_ibu_min: props.recipe.target_ibu_min,
    target_ibu_max: props.recipe.target_ibu_max,
    target_srm: props.recipe.target_srm,
    target_srm_min: props.recipe.target_srm_min,
    target_srm_max: props.recipe.target_srm_max,
    target_carbonation: props.recipe.target_carbonation,
    ibu_method: props.recipe.ibu_method,
    brewhouse_efficiency: props.recipe.brewhouse_efficiency,
  })

  // Reset form when recipe changes
  watch(() => props.recipe, newRecipe => {
    form.batch_size = newRecipe.batch_size
    form.batch_size_unit = newRecipe.batch_size_unit ?? 'bbl'
    form.target_og = newRecipe.target_og
    form.target_og_min = newRecipe.target_og_min
    form.target_og_max = newRecipe.target_og_max
    form.target_fg = newRecipe.target_fg
    form.target_fg_min = newRecipe.target_fg_min
    form.target_fg_max = newRecipe.target_fg_max
    form.target_ibu = newRecipe.target_ibu
    form.target_ibu_min = newRecipe.target_ibu_min
    form.target_ibu_max = newRecipe.target_ibu_max
    form.target_srm = newRecipe.target_srm
    form.target_srm_min = newRecipe.target_srm_min
    form.target_srm_max = newRecipe.target_srm_max
    form.target_carbonation = newRecipe.target_carbonation
    form.ibu_method = newRecipe.ibu_method
    form.brewhouse_efficiency = newRecipe.brewhouse_efficiency
  }, { deep: true })

  const volumeUnits = ['bbl', 'gal', 'L', 'hL']

  const ibuMethods = [
    { title: 'Tinseth', value: 'tinseth' },
    { title: 'Rager', value: 'rager' },
    { title: 'Garetz', value: 'garetz' },
    { title: 'Daniels', value: 'daniels' },
  ]

  const calculatedAbv = computed(() => {
    if (!form.target_og || !form.target_fg) return null
    if (form.target_og <= form.target_fg) return null
    // Standard ABV formula: (OG - FG) * 131.25
    return (form.target_og - form.target_fg) * 131.25
  })

  const hasChanges = computed(() => {
    return (
      form.batch_size !== props.recipe.batch_size
      || form.batch_size_unit !== (props.recipe.batch_size_unit ?? 'bbl')
      || form.target_og !== props.recipe.target_og
      || form.target_og_min !== props.recipe.target_og_min
      || form.target_og_max !== props.recipe.target_og_max
      || form.target_fg !== props.recipe.target_fg
      || form.target_fg_min !== props.recipe.target_fg_min
      || form.target_fg_max !== props.recipe.target_fg_max
      || form.target_ibu !== props.recipe.target_ibu
      || form.target_ibu_min !== props.recipe.target_ibu_min
      || form.target_ibu_max !== props.recipe.target_ibu_max
      || form.target_srm !== props.recipe.target_srm
      || form.target_srm_min !== props.recipe.target_srm_min
      || form.target_srm_max !== props.recipe.target_srm_max
      || form.target_carbonation !== props.recipe.target_carbonation
      || form.ibu_method !== props.recipe.ibu_method
      || form.brewhouse_efficiency !== props.recipe.brewhouse_efficiency
    )
  })

  function handleSubmit () {
    const data: UpdateRecipeRequest = {
      name: props.recipe.name,
      batch_size: form.batch_size || null,
      batch_size_unit: form.batch_size ? form.batch_size_unit : null,
      target_og: form.target_og || null,
      target_og_min: form.target_og_min || null,
      target_og_max: form.target_og_max || null,
      target_fg: form.target_fg || null,
      target_fg_min: form.target_fg_min || null,
      target_fg_max: form.target_fg_max || null,
      target_ibu: form.target_ibu || null,
      target_ibu_min: form.target_ibu_min || null,
      target_ibu_max: form.target_ibu_max || null,
      target_srm: form.target_srm || null,
      target_srm_min: form.target_srm_min || null,
      target_srm_max: form.target_srm_max || null,
      target_carbonation: form.target_carbonation || null,
      ibu_method: form.ibu_method || null,
      brewhouse_efficiency: form.brewhouse_efficiency || null,
    }
    emit('save', data)
  }
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.srm-swatch {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid rgba(var(--v-theme-on-surface), 0.2);
}
</style>
