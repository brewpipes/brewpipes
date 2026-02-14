<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 520"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Create batch</v-card-title>
      <v-card-text>
        <v-form ref="formRef" @submit.prevent="handleSubmit">
          <v-text-field
            v-model="form.short_name"
            density="comfortable"
            label="Short name"
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
            v-model="form.recipe_uuid"
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
          Create batch
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Recipe } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'

  export type BatchCreateForm = {
    short_name: string
    brew_date: string
    recipe_uuid: string | null
    notes: string
  }

  const props = defineProps<{
    modelValue: boolean
    recipes: Recipe[]
    recipesLoading: boolean
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [form: BatchCreateForm]
  }>()

  const formRef = ref()

  const form = reactive<BatchCreateForm>({
    short_name: '',
    brew_date: '',
    recipe_uuid: null,
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
      value: recipe.uuid,
      style: recipe.style_name,
    })),
  )

  // Reset form when dialog opens
  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        form.short_name = ''
        form.brew_date = ''
        form.recipe_uuid = null
        form.notes = ''
      }
    },
  )

  function handleSubmit () {
    if (!isFormValid.value) return
    emit('submit', { ...form })
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
