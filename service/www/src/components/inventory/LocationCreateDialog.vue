<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 480"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Create stock location</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="form.name"
          density="comfortable"
          label="Name"
          placeholder="Main warehouse"
          :rules="[rules.required]"
        />
        <v-text-field
          v-model="form.location_type"
          density="comfortable"
          label="Location type"
          placeholder="Warehouse, Cold storage, etc."
        />
        <v-textarea
          v-model="form.description"
          auto-grow
          density="comfortable"
          label="Description"
          placeholder="Additional details about this location..."
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
          Create
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateStockLocationRequest } from '@/types'
  import { computed, reactive, watch } from 'vue'
  import { normalizeText } from '@/utils/normalize'

  const props = defineProps<{
    modelValue: boolean
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: CreateStockLocationRequest]
  }>()

  const form = reactive({
    name: '',
    location_type: '',
    description: '',
  })

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
  }

  const isFormValid = computed(() => {
    return form.name.trim().length > 0
  })

  // Reset form when dialog opens
  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        form.name = ''
        form.location_type = ''
        form.description = ''
      }
    },
  )

  function handleSubmit () {
    if (!isFormValid.value) return

    const payload: CreateStockLocationRequest = {
      name: form.name.trim(),
      location_type: normalizeText(form.location_type),
      description: normalizeText(form.description),
    }

    emit('submit', payload)
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }
</script>
