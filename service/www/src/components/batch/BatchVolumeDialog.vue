<template>
  <v-dialog max-width="480" :model-value="modelValue" persistent @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title class="text-h6">Create wort volume</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12">
            <v-text-field
              density="comfortable"
              label="Name"
              :model-value="form.name"
              placeholder="IPA 24-07 Wort"
              @update:model-value="updateForm('name', $event)"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              density="comfortable"
              label="Amount"
              :model-value="form.amount"
              :rules="[rules.required, rules.positiveNumber]"
              type="number"
              @update:model-value="updateForm('amount', $event)"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-select
              density="comfortable"
              :items="volumeUnitOptions"
              label="Unit"
              :model-value="form.amount_unit"
              @update:model-value="updateForm('amount_unit', $event)"
            />
          </v-col>
          <v-col cols="12">
            <v-textarea
              auto-grow
              density="comfortable"
              label="Description"
              :model-value="form.description"
              rows="2"
              @update:model-value="updateForm('description', $event)"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="saving" variant="text" @click="emit('update:modelValue', false)">
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          :disabled="!isValid"
          :loading="saving"
          @click="emit('submit')"
        >
          Create volume
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { VolumeUnit } from '@/types'
  import { computed } from 'vue'

  export type VolumeForm = {
    name: string
    description: string
    amount: string
    amount_unit: VolumeUnit
  }

  const props = defineProps<{
    modelValue: boolean
    form: VolumeForm
    saving: boolean
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'update:form': [form: VolumeForm]
    'submit': []
  }>()

  const volumeUnitOptions: VolumeUnit[] = ['ml', 'usfloz', 'ukfloz', 'bbl']

  const rules = {
    required: (v: string) => !!v?.trim() || 'Required',
    positiveNumber: (v: string) => {
      const num = Number(v)
      return (Number.isFinite(num) && num > 0) || 'Must be positive'
    },
  }

  const isValid = computed(() => {
    const amount = Number(props.form.amount)
    return Number.isFinite(amount) && amount > 0
  })

  function updateForm<K extends keyof VolumeForm> (key: K, value: VolumeForm[K]) {
    emit('update:form', { ...props.form, [key]: value })
  }
</script>
