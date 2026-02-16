<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 500"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Create ingredient</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="form.name"
          density="comfortable"
          label="Name"
          :rules="[rules.required]"
        />
        <v-select
          v-model="form.category"
          density="comfortable"
          :items="ingredientCategoryOptions"
          label="Category"
          :rules="[rules.required]"
        />
        <v-combobox
          v-model="form.default_unit"
          density="comfortable"
          :items="unitOptions"
          label="Default unit"
          :rules="[rules.required]"
        />
        <v-textarea
          v-model="form.description"
          auto-grow
          density="comfortable"
          label="Description"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSubmit"
        >
          Create ingredient
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateIngredientRequest } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { normalizeText } from '@/utils/normalize'

  const ingredientCategoryOptions = [
    { title: 'Fermentable', value: 'fermentable' },
    { title: 'Hop', value: 'hop' },
    { title: 'Yeast', value: 'yeast' },
    { title: 'Adjunct', value: 'adjunct' },
    { title: 'Salt', value: 'salt' },
    { title: 'Chemical', value: 'chemical' },
    { title: 'Gas', value: 'gas' },
    { title: 'Other', value: 'other' },
  ]

  const unitOptions = ['kg', 'g', 'lb', 'oz', 'l', 'ml', 'gal', 'bbl']

  const props = defineProps<{
    modelValue: boolean
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateIngredientRequest]
  }>()

  const form = reactive({
    name: '',
    category: '',
    default_unit: '',
    description: '',
  })

  const rules = {
    required: (v: string | null) => (v !== null && v !== '' && String(v).trim() !== '') || 'Required',
  }

  const isFormValid = computed(() => {
    return form.name.trim() && form.category && form.default_unit
  })

  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        resetForm()
      }
    },
  )

  function resetForm () {
    form.name = ''
    form.category = ''
    form.default_unit = ''
    form.description = ''
  }

  function handleSubmit () {
    if (!isFormValid.value) return

    const payload: CreateIngredientRequest = {
      name: form.name.trim(),
      category: form.category,
      default_unit: form.default_unit,
      description: normalizeText(form.description),
    }
    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
