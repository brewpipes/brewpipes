<template>
  <v-dialog
    :fullscreen="$vuetify.display.xs"
    :max-width="$vuetify.display.xs ? '100%' : 400"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Delete Recipe</v-card-title>
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

        <p>
          Are you sure you want to delete <strong>{{ recipe.name }}</strong>?
        </p>
        <v-alert
          class="mt-4"
          density="compact"
          type="warning"
          variant="tonal"
        >
          This action cannot be undone. All ingredients associated with this recipe will also be deleted.
        </v-alert>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="deleting" variant="text" @click="handleCancel">
          Cancel
        </v-btn>
        <v-btn
          color="error"
          :loading="deleting"
          variant="flat"
          @click="handleConfirm"
        >
          Delete
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Recipe } from '@/composables/useProductionApi'

  defineProps<{
    modelValue: boolean
    recipe: Recipe
    deleting: boolean
    errorMessage?: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'confirm': []
  }>()

  function handleCancel () {
    emit('update:modelValue', false)
  }

  function handleConfirm () {
    emit('confirm')
  }
</script>
