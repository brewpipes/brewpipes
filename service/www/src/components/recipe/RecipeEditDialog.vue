<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    max-width="600"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Edit Recipe</v-card-title>
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
            v-model="form.name"
            density="comfortable"
            label="Name"
            placeholder="West Coast IPA"
            :rules="[rules.required]"
          />

          <v-combobox
            v-model="form.style"
            density="comfortable"
            hint="Select an existing style or type a new one"
            item-title="name"
            item-value="id"
            :items="styles"
            label="Style"
            :loading="stylesLoading"
            persistent-hint
            return-object
          >
            <template #no-data>
              <v-list-item>
                <v-list-item-title>
                  Type to search or create a new style
                </v-list-item-title>
              </v-list-item>
            </template>
          </v-combobox>

          <v-textarea
            v-model="form.notes"
            auto-grow
            density="comfortable"
            label="Notes"
            placeholder="Recipe description, process notes..."
            rows="3"
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
          Save Changes
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Recipe, Style, UpdateRecipeRequest } from '@/composables/useProductionApi'
  import { computed, reactive, ref, watch } from 'vue'

  const props = defineProps<{
    modelValue: boolean
    recipe: Recipe
    styles: Style[]
    stylesLoading: boolean
    saving: boolean
    errorMessage?: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: UpdateRecipeRequest]
  }>()

  const formRef = ref()

  const form = reactive({
    name: '',
    style: null as Style | string | null,
    notes: '',
  })

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
  }

  // Initialize form when dialog opens
  watch(() => props.modelValue, open => {
    if (open) {
      form.name = props.recipe.name
      form.notes = props.recipe.notes ?? ''

      // Set the style
      if (props.recipe.style_id) {
        const matchingStyle = props.styles.find(s => s.id === props.recipe.style_id)
        form.style = matchingStyle ?? props.recipe.style_name
      } else if (props.recipe.style_name) {
        form.style = props.recipe.style_name
      } else {
        form.style = null
      }
    }
  })

  const isFormValid = computed(() => {
    return form.name.trim().length > 0
  })

  function handleCancel () {
    emit('update:modelValue', false)
  }

  function handleSubmit () {
    if (!isFormValid.value) return

    let styleId: number | null = null
    let styleName: string | null = null

    if (form.style) {
      if (typeof form.style === 'object' && form.style.id) {
        styleId = form.style.id
        styleName = form.style.name
      } else if (typeof form.style === 'string') {
        styleName = form.style.trim() || null
      }
    }

    const data: UpdateRecipeRequest = {
      name: form.name.trim(),
      style_id: styleId,
      style_name: styleName,
      notes: form.notes.trim() || null,
    }

    emit('submit', data)
  }
</script>
