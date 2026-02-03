<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 520"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Edit batch</v-card-title>
      <v-card-text>
        <v-alert
          v-if="errorMessage"
          class="mb-4"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ errorMessage }}
        </v-alert>

        <v-form ref="formRef" @submit.prevent="handleSubmit">
          <v-text-field
            v-model="form.short_name"
            density="comfortable"
            label="Name"
            placeholder="IPA 24-07"
            :rules="[rules.required]"
          />
          <v-text-field
            v-model="form.brew_date"
            density="comfortable"
            label="Brew date"
            type="date"
          />
          <v-autocomplete
            v-model="form.recipe_id"
            clearable
            density="comfortable"
            hint="Optional - link this batch to a recipe"
            item-title="title"
            item-value="value"
            :items="recipeSelectItems"
            label="Recipe"
            :loading="recipesLoading"
            persistent-hint
          >
            <template #item="{ props: itemProps, item }">
              <v-list-item v-bind="itemProps">
                <template #subtitle>
                  <span v-if="item.raw.style">{{ item.raw.style }}</span>
                </template>
              </v-list-item>
            </template>
          </v-autocomplete>
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
        <v-btn :disabled="saving" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="handleSubmit"
        >
          Save changes
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed, reactive, ref, watch } from 'vue'
  import type { Batch, Recipe } from '@/types'

  export type BatchEditForm = {
    short_name: string
    brew_date: string
    recipe_id: number | null
    notes: string
  }

  const props = defineProps<{
    modelValue: boolean
    batch: Batch | null
    recipes: Recipe[]
    recipesLoading: boolean
    saving: boolean
    errorMessage: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [form: BatchEditForm]
  }>()

  const formRef = ref()

  const form = reactive<BatchEditForm>({
    short_name: '',
    brew_date: '',
    recipe_id: null,
    notes: '',
  })

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
  }

  const isFormValid = computed(() => {
    return form.short_name.trim().length > 0
  })

  const recipeSelectItems = computed(() =>
    props.recipes.map(recipe => ({
      title: recipe.name,
      value: recipe.id,
      style: recipe.style_name,
    })),
  )

  // Reset form when dialog opens with batch data
  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen && props.batch) {
        form.short_name = props.batch.short_name
        form.brew_date = props.batch.brew_date
          ? formatDateForInput(props.batch.brew_date)
          : ''
        form.recipe_id = props.batch.recipe_id
        form.notes = props.batch.notes ?? ''
      }
    },
  )

  function formatDateForInput (isoString: string): string {
    if (!isoString) return ''
    const date = new Date(isoString)
    const pad = (n: number) => String(n).padStart(2, '0')
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
  }

  function handleSubmit () {
    if (!isFormValid.value) return
    emit('submit', { ...form })
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
