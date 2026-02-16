<template>
  <v-dialog
    max-width="600"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditing ? 'Edit recipe' : 'Create recipe' }}
      </v-card-title>
      <v-card-text>
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
            item-value="uuid"
            :items="styles"
            label="Style"
            persistent-hint
            return-object
            @update:search="onStyleSearch"
          >
            <template #no-data>
              <v-list-item v-if="styleSearchQuery">
                <v-list-item-title>
                  Press enter to create "{{ styleSearchQuery }}"
                </v-list-item-title>
              </v-list-item>
              <v-list-item v-else>
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
            placeholder="Recipe description, ingredients, process notes..."
            rows="3"
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
          {{ isEditing ? 'Save changes' : 'Create recipe' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateRecipeRequest, Recipe, Style } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { normalizeText } from '@/utils/normalize'

  export interface RecipeCreateEditSubmitData {
    name: string
    style_uuid: string | null
    style_name: string | null
    notes: string | null
  }

  const props = defineProps<{
    modelValue: boolean
    editRecipe?: Recipe | null
    styles: Style[]
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: RecipeCreateEditSubmitData]
  }>()

  const formRef = ref()
  const styleSearchQuery = ref('')

  const form = reactive({
    name: '',
    style: null as Style | string | null,
    notes: '',
  })

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
  }

  const isEditing = computed(() => !!props.editRecipe)

  const isFormValid = computed(() => {
    return form.name.trim().length > 0
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
    styleSearchQuery.value = ''
    if (props.editRecipe) {
      form.name = props.editRecipe.name
      form.notes = props.editRecipe.notes ?? ''
      if (props.editRecipe.style_uuid) {
        const matchingStyle = props.styles.find(s => s.uuid === props.editRecipe!.style_uuid)
        form.style = matchingStyle ?? props.editRecipe.style_name
      } else if (props.editRecipe.style_name) {
        form.style = props.editRecipe.style_name
      } else {
        form.style = null
      }
    } else {
      form.name = ''
      form.style = null
      form.notes = ''
    }
  }

  function onStyleSearch (query: string) {
    styleSearchQuery.value = query
  }

  function handleSubmit () {
    if (!isFormValid.value) return

    let styleUuid: string | null = null
    let styleName: string | null = null

    if (form.style) {
      if (typeof form.style === 'object' && form.style.uuid) {
        styleUuid = form.style.uuid
        styleName = form.style.name
      } else if (typeof form.style === 'string') {
        styleName = form.style.trim() || null
      }
    }

    const payload: RecipeCreateEditSubmitData = {
      name: form.name.trim(),
      style_uuid: styleUuid,
      style_name: styleName,
      notes: normalizeText(form.notes),
    }
    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
