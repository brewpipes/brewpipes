<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    max-width="600"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit' : 'Add' }} {{ ingredientTypeLabel }}
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef" @submit.prevent="handleSubmit">
          <!-- Name -->
          <v-text-field
            v-model="form.name"
            density="comfortable"
            label="Name"
            :placeholder="namePlaceholder"
            :rules="[rules.required]"
          />

          <!-- Amount and Unit -->
          <v-row dense>
            <v-col cols="8">
              <v-text-field
                v-model.number="form.amount"
                density="comfortable"
                label="Amount"
                :rules="[rules.required, rules.positive]"
                type="number"
              />
            </v-col>
            <v-col cols="4">
              <v-select
                v-model="form.amount_unit"
                density="comfortable"
                :items="unitOptions"
                label="Unit"
              />
            </v-col>
          </v-row>

          <!-- Use Stage -->
          <v-select
            v-model="form.use_stage"
            density="comfortable"
            :items="stageOptions"
            label="Use Stage"
          />

          <!-- Use Type (context-dependent) -->
          <v-select
            v-if="useTypeOptions.length > 0"
            v-model="form.use_type"
            clearable
            density="comfortable"
            :items="useTypeOptions"
            label="Use Type"
          />

          <!-- Timing (for hops) -->
          <template v-if="showTiming">
            <v-row dense>
              <v-col cols="6">
                <v-text-field
                  v-model.number="form.timing_duration_minutes"
                  density="comfortable"
                  :hint="timingHint"
                  :label="timingLabel"
                  persistent-hint
                  type="number"
                />
              </v-col>
              <v-col v-if="showTemperature" cols="6">
                <v-text-field
                  v-model.number="form.timing_temperature_c"
                  density="comfortable"
                  hint="For whirlpool additions"
                  label="Temperature (°C)"
                  persistent-hint
                  type="number"
                />
              </v-col>
            </v-row>
          </template>

          <!-- Alpha Acid (for hops) -->
          <v-text-field
            v-if="ingredientType === 'hop'"
            v-model.number="form.alpha_acid_assumed"
            density="comfortable"
            hint="Used for IBU calculations"
            label="Alpha Acid %"
            persistent-hint
            step="0.1"
            type="number"
          />

          <!-- Notes -->
          <v-textarea
            v-model="form.notes"
            auto-grow
            density="comfortable"
            label="Notes"
            rows="2"
          />
        </v-form>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSubmit"
        >
          {{ isEditing ? 'Save Changes' : 'Add Ingredient' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type {
    CreateRecipeIngredientRequest,
    RecipeIngredient,
    RecipeIngredientType,
    RecipeUseStage,
    RecipeUseType,
    UpdateRecipeIngredientRequest,
  } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'

  const props = defineProps<{
    modelValue: boolean
    ingredientType: RecipeIngredientType
    editingIngredient: RecipeIngredient | null
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateRecipeIngredientRequest | UpdateRecipeIngredientRequest]
  }>()

  const formRef = ref()

  const form = reactive({
    name: '',
    amount: null as number | null,
    amount_unit: 'lb',
    use_stage: 'mash' as RecipeUseStage,
    use_type: null as RecipeUseType | null,
    timing_duration_minutes: null as number | null,
    timing_temperature_c: null as number | null,
    alpha_acid_assumed: null as number | null,
    notes: '',
  })

  const rules = {
    required: (v: string | number | null) => {
      if (typeof v === 'number') return true
      return !!v?.toString().trim() || 'Required'
    },
    positive: (v: number | null) => {
      if (v === null || v === undefined) return 'Required'
      return v > 0 || 'Must be positive'
    },
  }

  const isEditing = computed(() => props.editingIngredient !== null)

  // Initialize form when dialog opens or editing ingredient changes
  watch([() => props.modelValue, () => props.editingIngredient], ([open, ingredient]) => {
    if (open) {
      if (ingredient) {
        // Editing existing ingredient
        form.name = ingredient.name ?? ''
        form.amount = ingredient.amount
        form.amount_unit = ingredient.amount_unit
        form.use_stage = ingredient.use_stage
        form.use_type = ingredient.use_type
        // Convert minutes to days for fermentation stage (dry hop)
        form.timing_duration_minutes = ingredient.use_stage === 'fermentation' && ingredient.timing_duration_minutes !== null ? Math.round(ingredient.timing_duration_minutes / 1440) : ingredient.timing_duration_minutes
        form.timing_temperature_c = ingredient.timing_temperature_c
        form.alpha_acid_assumed = ingredient.alpha_acid_assumed
        form.notes = ingredient.notes ?? ''
      } else {
        // Creating new ingredient
        resetForm()
      }
    }
  })

  // Update defaults when ingredient type changes
  watch(() => props.ingredientType, type => {
    if (!props.editingIngredient) {
      setDefaultsForType(type)
    }
  })

  function resetForm () {
    form.name = ''
    form.amount = null
    form.notes = ''
    form.timing_duration_minutes = null
    form.timing_temperature_c = null
    form.alpha_acid_assumed = null
    setDefaultsForType(props.ingredientType)
  }

  function setDefaultsForType (type: RecipeIngredientType) {
    switch (type) {
      case 'fermentable': {
        form.amount_unit = 'lb'
        form.use_stage = 'mash'
        form.use_type = 'base'
        break
      }
      case 'hop': {
        form.amount_unit = 'oz'
        form.use_stage = 'boil'
        form.use_type = 'bittering'
        form.timing_duration_minutes = 60
        break
      }
      case 'yeast': {
        form.amount_unit = 'pkg'
        form.use_stage = 'fermentation'
        form.use_type = 'primary'
        break
      }
      case 'adjunct': {
        form.amount_unit = 'oz'
        form.use_stage = 'boil'
        form.use_type = 'adjunct'
        break
      }
      case 'salt':
      case 'chemical': {
        form.amount_unit = 'g'
        form.use_stage = 'mash'
        form.use_type = null
        break
      }
      default: {
        form.amount_unit = 'oz'
        form.use_stage = 'boil'
        form.use_type = 'other'
      }
    }
  }

  const ingredientTypeLabel = computed(() => {
    const labels: Record<RecipeIngredientType, string> = {
      fermentable: 'Fermentable',
      hop: 'Hop',
      yeast: 'Yeast',
      adjunct: 'Adjunct',
      salt: 'Salt',
      chemical: 'Chemical',
      gas: 'Gas',
      other: 'Ingredient',
    }
    return labels[props.ingredientType] ?? 'Ingredient'
  })

  const namePlaceholder = computed(() => {
    const placeholders: Record<RecipeIngredientType, string> = {
      fermentable: 'Pale Malt 2-Row',
      hop: 'Centennial',
      yeast: 'US-05',
      adjunct: 'Coriander',
      salt: 'Gypsum',
      chemical: 'Lactic Acid',
      gas: 'CO2',
      other: 'Ingredient name',
    }
    return placeholders[props.ingredientType] ?? 'Ingredient name'
  })

  const unitOptions = computed(() => {
    switch (props.ingredientType) {
      case 'fermentable': {
        return ['lb', 'kg', 'oz', 'g']
      }
      case 'hop': {
        return ['oz', 'g', 'lb', 'kg']
      }
      case 'yeast': {
        return ['pkg', 'g', 'ml', 'L', 'cells']
      }
      case 'salt':
      case 'chemical': {
        return ['g', 'mg', 'oz', 'tsp', 'ml']
      }
      case 'gas': {
        return ['vol', 'psi', 'bar']
      }
      default: {
        return ['oz', 'g', 'lb', 'kg', 'ml', 'L', 'tsp', 'tbsp', 'cup']
      }
    }
  })

  const stageOptions = computed(() => {
    const allStages = [
      { title: 'Mash', value: 'mash' },
      { title: 'Boil', value: 'boil' },
      { title: 'Whirlpool', value: 'whirlpool' },
      { title: 'Fermentation', value: 'fermentation' },
      { title: 'Packaging', value: 'packaging' },
    ]

    // Filter based on ingredient type
    switch (props.ingredientType) {
      case 'fermentable': {
        return allStages.filter(s => ['mash', 'boil'].includes(s.value))
      }
      case 'hop': {
        return allStages
      }
      case 'yeast': {
        return allStages.filter(s => ['fermentation', 'packaging'].includes(s.value))
      }
      default: {
        return allStages
      }
    }
  })

  const useTypeOptions = computed(() => {
    switch (props.ingredientType) {
      case 'fermentable': {
        return [
          { title: 'Base Malt', value: 'base' },
          { title: 'Specialty Malt', value: 'specialty' },
          { title: 'Adjunct', value: 'adjunct' },
          { title: 'Sugar', value: 'sugar' },
          { title: 'Other', value: 'other' },
        ]
      }
      case 'hop': {
        return [
          { title: 'Bittering', value: 'bittering' },
          { title: 'Flavor', value: 'flavor' },
          { title: 'Aroma', value: 'aroma' },
          { title: 'Dry Hop', value: 'dry_hop' },
          { title: 'Other', value: 'other' },
        ]
      }
      case 'yeast': {
        return [
          { title: 'Primary', value: 'primary' },
          { title: 'Secondary', value: 'secondary' },
          { title: 'Bottle', value: 'bottle' },
          { title: 'Other', value: 'other' },
        ]
      }
      default: {
        return []
      }
    }
  })

  const showTiming = computed(() => {
    return props.ingredientType === 'hop'
  })

  const showTemperature = computed(() => {
    return props.ingredientType === 'hop' && form.use_stage === 'whirlpool'
  })

  const timingLabel = computed(() => {
    if (form.use_stage === 'fermentation') {
      return 'Duration (days)'
    }
    return 'Time (minutes)'
  })

  const timingHint = computed(() => {
    if (form.use_stage === 'fermentation') {
      return 'Days in contact with beer'
    }
    if (form.use_stage === 'boil') {
      return 'Minutes remaining in boil'
    }
    return 'Duration of addition'
  })

  const isFormValid = computed(() => {
    return (
      (form.name?.trim().length ?? 0) > 0
      && form.amount !== null
      && form.amount > 0
    )
  })

  function handleCancel () {
    emit('update:modelValue', false)
  }

  function handleSubmit () {
    if (!isFormValid.value) return

    // Convert days to minutes for dry hop timing
    let timingMinutes = form.timing_duration_minutes
    if (form.use_stage === 'fermentation' && timingMinutes !== null) {
      timingMinutes = timingMinutes * 1440 // days to minutes
    }

    const data: CreateRecipeIngredientRequest | UpdateRecipeIngredientRequest = {
      ingredient_type: props.ingredientType,
      name: (form.name ?? '').trim(),
      amount: form.amount!,
      amount_unit: form.amount_unit,
      use_stage: form.use_stage,
      use_type: form.use_type,
      timing_duration_minutes: timingMinutes,
      timing_temperature_c: form.timing_temperature_c,
      alpha_acid_assumed: form.alpha_acid_assumed,
      notes: form.notes?.trim() || null,
    }

    emit('submit', data)
  }
</script>
