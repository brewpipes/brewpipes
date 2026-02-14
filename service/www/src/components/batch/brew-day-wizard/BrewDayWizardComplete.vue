<template>
  <div class="text-center pa-4">
    <v-icon
      class="mb-4"
      color="success"
      icon="mdi-check-circle-outline"
      size="64"
    />

    <div class="text-h6 mb-4">Brew Day Complete</div>

    <v-list class="text-left mb-6" density="compact">
      <!-- Pick step -->
      <v-list-item v-if="pickStep">
        <template #prepend>
          <v-icon
            :color="pickStep.status === 'complete' ? 'success' : 'grey'"
            :icon="pickStep.status === 'complete' ? 'mdi-check-circle' : 'mdi-skip-forward'"
            size="small"
          />
        </template>
        <v-list-item-title v-if="pickStep.status === 'complete'">
          Ingredients picked ({{ completionData.ingredientCount }} ingredients from {{ completionData.lotCount }} lots)
        </v-list-item-title>
        <v-list-item-title v-else class="text-medium-emphasis">
          Ingredient picking skipped
        </v-list-item-title>
      </v-list-item>

      <!-- Session step -->
      <v-list-item v-if="sessionStep">
        <template #prepend>
          <v-icon
            :color="sessionStep.status === 'complete' ? 'success' : 'grey'"
            :icon="sessionStep.status === 'complete' ? 'mdi-check-circle' : 'mdi-skip-forward'"
            size="small"
          />
        </template>
        <v-list-item-title v-if="sessionStep.status === 'complete'">
          Brew session recorded
          <span v-if="completionData.mashTemp || completionData.originalGravity" class="text-medium-emphasis">
            ({{ [
              completionData.mashTemp ? `Mash: ${completionData.mashTemp}` : null,
              completionData.originalGravity ? `OG: ${completionData.originalGravity}` : null,
            ].filter(Boolean).join(', ') }})
          </span>
        </v-list-item-title>
        <v-list-item-title v-else class="text-medium-emphasis">
          Brew session skipped
        </v-list-item-title>
      </v-list-item>

      <!-- Fermenter step -->
      <v-list-item v-if="fermenterStep">
        <template #prepend>
          <v-icon
            :color="fermenterStep.status === 'complete' ? 'success' : 'grey'"
            :icon="fermenterStep.status === 'complete' ? 'mdi-check-circle' : 'mdi-skip-forward'"
            size="small"
          />
        </template>
        <v-list-item-title v-if="fermenterStep.status === 'complete' && completionData.fermenterName">
          Assigned to {{ completionData.fermenterName }}
        </v-list-item-title>
        <v-list-item-title v-else class="text-medium-emphasis">
          Fermenter assignment skipped
        </v-list-item-title>
      </v-list-item>
    </v-list>

    <v-btn
      block
      color="primary"
      min-height="44"
      size="large"
      variant="flat"
      @click="emit('done')"
    >
      Done &mdash; Back to Batch
    </v-btn>
  </div>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import type { WizardCompletionData, WizardStep } from './types'

  const props = defineProps<{
    steps: WizardStep[]
    completionData: WizardCompletionData
  }>()

  const emit = defineEmits<{
    done: []
  }>()

  const pickStep = computed(() => props.steps.find(s => s.id === 'pick'))
  const sessionStep = computed(() => props.steps.find(s => s.id === 'session'))
  const fermenterStep = computed(() => props.steps.find(s => s.id === 'fermenter'))
</script>
