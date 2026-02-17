<template>
  <v-dialog
    :max-width="$vuetify.display.xs ? '100%' : 480"
    :model-value="modelValue"
    persistent
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="text-h6">Retire vessel</v-card-title>
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

        <p class="text-body-1 mb-4">
          Are you sure you want to retire
          <strong>{{ vesselName }}</strong>?
        </p>

        <v-alert
          density="compact"
          icon="mdi-alert"
          type="warning"
          variant="tonal"
        >
          Retired vessels cannot be used for new occupancies. This action can be reversed by editing the vessel and changing its status back to Active.
        </v-alert>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="retiring" variant="text" @click="handleCancel">Cancel</v-btn>
        <v-btn
          color="warning"
          :loading="retiring"
          variant="flat"
          @click="handleConfirm"
        >
          Retire vessel
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Vessel } from '@/types'
  import { computed, ref, watch } from 'vue'

  const props = defineProps<{
    modelValue: boolean
    vessel: Vessel | null
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'confirm': []
  }>()

  const retiring = ref(false)
  const errorMessage = ref('')

  const vesselName = computed(() => props.vessel?.name ?? 'this vessel')

  // Reset state when dialog opens
  watch(
    () => props.modelValue,
    open => {
      if (open) {
        errorMessage.value = ''
        retiring.value = false
      }
    },
  )

  function handleConfirm () {
    emit('confirm')
  }

  function handleCancel () {
    emit('update:modelValue', false)
  }

  // Expose methods for parent to control state
  defineExpose({
    setRetiring: (value: boolean) => {
      retiring.value = value
    },
    setError: (message: string) => {
      errorMessage.value = message
    },
  })
</script>
