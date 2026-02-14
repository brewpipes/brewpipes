<template>
  <div class="step-header d-flex align-center justify-center pa-2">
    <template v-for="(step, index) in steps" :key="step.id">
      <div
        v-if="index > 0"
        class="step-connector"
        :class="{ 'step-connector--complete': isBeforeOrComplete(index) }"
      />
      <div class="step-item">
        <button
          :aria-label="`Go to ${step.label}`"
          class="step-circle"
          :class="stepCircleClass(step)"
          type="button"
          @click="emit('select', step.id)"
        >
          <v-icon v-if="step.status === 'complete'" size="18">mdi-check</v-icon>
          <v-icon v-else-if="step.status === 'skipped'" size="18">mdi-minus</v-icon>
          <v-icon v-else :icon="step.icon" size="18" />
        </button>
        <span
          class="step-label d-none d-sm-block"
          :class="{ 'text-primary font-weight-medium': step.id === currentStep }"
        >
          {{ step.label }}
        </span>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
  import type { WizardStep, WizardStepId } from './types'

  const props = defineProps<{
    steps: WizardStep[]
    currentStep: WizardStepId
  }>()

  const emit = defineEmits<{
    select: [stepId: WizardStepId]
  }>()

  function stepCircleClass (step: WizardStep): Record<string, boolean> {
    return {
      'step-circle--active': step.id === props.currentStep && step.status !== 'complete',
      'step-circle--complete': step.status === 'complete',
      'step-circle--skipped': step.status === 'skipped',
      'step-circle--not-started': step.status === 'not_started' && step.id !== props.currentStep,
    }
  }

  function isBeforeOrComplete (index: number): boolean {
    // The connector before step[index] is "complete" if the previous step is complete
    const prevStep = props.steps[index - 1]
    return prevStep?.status === 'complete' || prevStep?.status === 'skipped'
  }
</script>

<style scoped>
.step-header {
  gap: 0;
  flex-wrap: nowrap;
}

.step-connector {
  flex: 1;
  height: 2px;
  max-width: 80px;
  background: rgba(var(--v-theme-on-surface), 0.2);
  transition: background 0.2s ease;
}

.step-connector--complete {
  background: rgb(var(--v-theme-success));
}

.step-circle {
  width: 44px;
  height: 44px;
  min-width: 44px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid rgba(var(--v-theme-on-surface), 0.3);
  background: transparent;
  color: rgba(var(--v-theme-on-surface), 0.5);
  cursor: pointer;
  transition: all 0.2s ease;
  outline: none;
  padding: 0;
}

.step-circle:focus-visible {
  box-shadow: 0 0 0 3px rgba(var(--v-theme-primary), 0.3);
}

.step-circle--active {
  border-color: rgb(var(--v-theme-primary));
  background: rgb(var(--v-theme-primary));
  color: rgb(var(--v-theme-on-primary));
}

.step-circle--complete {
  border-color: rgb(var(--v-theme-success));
  background: rgb(var(--v-theme-success));
  color: white;
}

.step-circle--skipped {
  border-color: rgba(var(--v-theme-on-surface), 0.3);
  background: transparent;
  color: rgba(var(--v-theme-on-surface), 0.4);
  border-style: dashed;
}

.step-circle--not-started {
  border-color: rgba(var(--v-theme-on-surface), 0.2);
  color: rgba(var(--v-theme-on-surface), 0.4);
}

.step-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.step-label {
  margin-top: 4px;
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  white-space: nowrap;
}
</style>
