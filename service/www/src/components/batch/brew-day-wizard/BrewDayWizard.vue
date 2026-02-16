<template>
  <v-dialog
    :fullscreen="isXs"
    :max-width="isXs ? '100%' : 700"
    :model-value="modelValue"
    persistent
    scrollable
    @update:model-value="handleDialogUpdate"
  >
    <v-card class="d-flex flex-column" :height="isXs ? '100%' : undefined">
      <!-- Title bar -->
      <v-card-title class="d-flex align-center pa-4">
        <v-icon class="mr-2" icon="mdi-kettle-steam" />
        <span class="text-h6 text-truncate">Brew Day: {{ batch.short_name }}</span>
        <v-spacer />
        <v-btn
          aria-label="Close wizard"
          icon="mdi-close"
          size="small"
          variant="text"
          @click="handleClose"
        />
      </v-card-title>

      <v-divider />

      <!-- Step header (hidden on completion) -->
      <BrewDayWizardStepHeader
        v-if="!showCompletion"
        :current-step="currentStep"
        :steps="steps"
        @select="navigateToStep"
      />

      <v-divider v-if="!showCompletion" />

      <!-- Step content -->
      <v-card-text class="flex-grow-1 overflow-y-auto pa-4">
        <v-window v-if="!showCompletion" :key="wizardKey" v-model="currentStep" :touch="false">
          <v-window-item value="pick">
            <BrewDayWizardStepPick
              v-if="batch.recipe_uuid"
              ref="pickStepRef"
              :batch-uuid="batch.uuid"
              :recipe-uuid="batch.recipe_uuid"
              @completed="handlePickCompleted"
            />
            <v-alert
              v-else
              density="comfortable"
              type="info"
              variant="tonal"
            >
              No recipe assigned to this batch. Skip this step or assign a recipe first.
            </v-alert>
          </v-window-item>

          <v-window-item value="session">
            <BrewDayWizardStepSession
              ref="sessionStepRef"
              :batch-uuid="batch.uuid"
              :vessels="vessels"
              :volumes="volumes"
              @completed="handleSessionCompleted"
            />
          </v-window-item>

          <v-window-item value="fermenter">
            <BrewDayWizardStepFermenter
              ref="fermenterStepRef"
              :batch-uuid="batch.uuid"
              :occupancies="occupancies"
              :vessels="vessels"
              :volumes="volumes"
              @completed="handleFermenterCompleted"
            />
          </v-window-item>
        </v-window>

        <!-- Completion screen -->
        <BrewDayWizardComplete
          v-if="showCompletion"
          :completion-data="completionData"
          :steps="steps"
          @done="handleDone"
        />
      </v-card-text>

      <!-- Bottom action bar (hidden on completion) -->
      <template v-if="!showCompletion">
        <v-divider />
        <v-card-actions class="pa-4 action-bar">
          <!-- Back button -->
          <v-btn
            v-if="currentStepIndex > 0"
            min-height="44"
            prepend-icon="mdi-arrow-left"
            variant="text"
            @click="goBack"
          >
            Back
          </v-btn>
          <v-spacer />

          <!-- Skip button -->
          <v-btn
            class="mr-2"
            min-height="44"
            variant="text"
            @click="skipStep"
          >
            Skip for now
          </v-btn>

          <!-- Primary action button -->
          <v-btn
            v-if="currentStep === 'pick'"
            color="primary"
            :disabled="!canConfirmPick"
            :loading="pickConfirming"
            min-height="44"
            variant="flat"
            @click="confirmAndNext"
          >
            Confirm &amp; Next &rarr;
          </v-btn>
          <v-btn
            v-else-if="currentStep === 'session'"
            color="primary"
            :loading="sessionSaving"
            min-height="44"
            variant="flat"
            @click="saveAndNext"
          >
            Save &amp; Next &rarr;
          </v-btn>
          <v-btn
            v-else-if="currentStep === 'fermenter'"
            color="primary"
            :disabled="!canFinishBrewDay"
            :loading="fermenterSaving"
            min-height="44"
            variant="flat"
            @click="finishBrewDay"
          >
            Finish Brew Day &check;
          </v-btn>
        </v-card-actions>
      </template>
    </v-card>

    <!-- Confirm close dialog -->
    <v-dialog v-model="showConfirmClose" max-width="400" persistent>
      <v-card>
        <v-card-title class="text-h6">Leave Brew Day?</v-card-title>
        <v-card-text>
          You have unsaved progress. Completed steps are already saved, but the current step will be lost.
        </v-card-text>
        <v-card-actions class="justify-end">
          <v-btn variant="text" @click="showConfirmClose = false">Stay</v-btn>
          <v-btn color="error" variant="flat" @click="forceClose">Leave</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { WizardCompletionData, WizardStep, WizardStepId } from './types'
  import type { Batch, BrewSession, Occupancy, Vessel, Volume } from '@/types'
  import { computed, ref, watch } from 'vue'
  import { useDisplay } from 'vuetify'
  import BrewDayWizardComplete from './BrewDayWizardComplete.vue'
  import BrewDayWizardStepFermenter from './BrewDayWizardStepFermenter.vue'
  import BrewDayWizardStepHeader from './BrewDayWizardStepHeader.vue'
  import BrewDayWizardStepPick from './BrewDayWizardStepPick.vue'
  import BrewDayWizardStepSession from './BrewDayWizardStepSession.vue'

  const props = defineProps<{
    modelValue: boolean
    batch: Batch
    brewSessions: BrewSession[]
    vessels: Vessel[]
    volumes: Volume[]
    occupancies: Occupancy[]
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'completed': []
  }>()

  const { xs: isXs } = useDisplay()

  // Step refs
  const pickStepRef = ref<InstanceType<typeof BrewDayWizardStepPick> | null>(null)
  const sessionStepRef = ref<InstanceType<typeof BrewDayWizardStepSession> | null>(null)
  const fermenterStepRef = ref<InstanceType<typeof BrewDayWizardStepFermenter> | null>(null)

  // Wizard state
  const wizardKey = ref(0)
  const currentStep = ref<WizardStepId>('pick')
  const showCompletion = ref(false)
  const showConfirmClose = ref(false)

  const steps = ref<WizardStep[]>([
    { id: 'pick', label: 'Pick', icon: 'mdi-clipboard-list-outline', status: 'not_started' },
    { id: 'session', label: 'Brew', icon: 'mdi-kettle-steam', status: 'not_started' },
    { id: 'fermenter', label: 'Fermenter', icon: 'mdi-flask-round-bottom', status: 'not_started' },
  ])

  const completionData = ref<WizardCompletionData>({
    ingredientCount: 0,
    lotCount: 0,
    session: null,
    mashTemp: null,
    originalGravity: null,
    occupancy: null,
    fermenterName: null,
  })

  // Computed
  const stepOrder = ['pick', 'session', 'fermenter'] as const

  const currentStepIndex = computed(() =>
    stepOrder.indexOf(currentStep.value),
  )

  const canConfirmPick = computed(() =>
    pickStepRef.value?.hasAnyPicks ?? false,
  )

  const pickConfirming = computed(() =>
    pickStepRef.value?.confirming ?? false,
  )

  const sessionSaving = computed(() =>
    sessionStepRef.value?.saving ?? false,
  )

  const fermenterSaving = computed(() =>
    fermenterStepRef.value?.saving ?? false,
  )

  const canFinishBrewDay = computed(() =>
    fermenterStepRef.value?.isFormValid ?? false,
  )

  const hasInProgressStep = computed(() =>
    steps.value.some(s => s.status === 'in_progress'),
  )

  // Reset wizard when dialog opens
  watch(() => props.modelValue, isOpen => {
    if (isOpen) {
      resetWizard()
    }
  })

  function resetWizard () {
    wizardKey.value++
    showCompletion.value = false
    showConfirmClose.value = false

    completionData.value = {
      ingredientCount: 0,
      lotCount: 0,
      session: null,
      mashTemp: null,
      originalGravity: null,
      occupancy: null,
      fermenterName: null,
    }

    // Reset step statuses
    for (const step of steps.value) {
      step.status = 'not_started'
    }

    // Determine starting step based on existing data
    if (props.brewSessions.length > 0) {
      // Session already exists — skip pick and session
      setStepStatus('pick', 'skipped')
      setStepStatus('session', 'skipped')
      currentStep.value = 'fermenter'
    } else if (props.batch.recipe_uuid) {
      currentStep.value = 'pick'
    } else {
      // No recipe — skip pick
      setStepStatus('pick', 'skipped')
      currentStep.value = 'session'
    }

    // Mark current step as in progress
    setStepStatus(currentStep.value, 'in_progress')
  }

  function setStepStatus (stepId: WizardStepId, status: WizardStep['status']) {
    const step = steps.value.find(s => s.id === stepId)
    if (step) {
      step.status = status
    }
  }

  function navigateToStep (stepId: WizardStepId) {
    if (showCompletion.value) return
    const step = steps.value.find(s => s.id === stepId)
    if (step?.status === 'complete') return
    currentStep.value = stepId
    if (step && step.status === 'not_started') {
      step.status = 'in_progress'
    }
  }

  function goBack () {
    const prevIndex = currentStepIndex.value - 1
    if (prevIndex >= 0) {
      const prevStepId = stepOrder[prevIndex]
      if (prevStepId) {
        currentStep.value = prevStepId
      }
    }
  }

  function advanceToNextStep () {
    const nextIndex = currentStepIndex.value + 1
    if (nextIndex < stepOrder.length) {
      const nextStepId = stepOrder[nextIndex]
      if (nextStepId) {
        currentStep.value = nextStepId
        setStepStatus(nextStepId, 'in_progress')
      }
    } else {
      // All steps done
      showCompletion.value = true
    }
  }

  function skipStep () {
    setStepStatus(currentStep.value, 'skipped')
    advanceToNextStep()
  }

  // Step completion handlers
  function handlePickCompleted (data: { ingredientCount: number, lotCount: number }) {
    setStepStatus('pick', 'complete')
    completionData.value.ingredientCount = data.ingredientCount
    completionData.value.lotCount = data.lotCount
    advanceToNextStep()
  }

  function handleSessionCompleted (data: { session: BrewSession, mashTemp: string | null, originalGravity: string | null }) {
    setStepStatus('session', 'complete')
    completionData.value.session = data.session
    completionData.value.mashTemp = data.mashTemp
    completionData.value.originalGravity = data.originalGravity
    advanceToNextStep()
  }

  function handleFermenterCompleted (data: { occupancy: Occupancy, fermenterName: string }) {
    setStepStatus('fermenter', 'complete')
    completionData.value.occupancy = data.occupancy
    completionData.value.fermenterName = data.fermenterName
    advanceToNextStep()
  }

  // Action handlers
  async function confirmAndNext () {
    if (pickStepRef.value) {
      await pickStepRef.value.confirmPicks()
      // The completed event handler will advance the step
    }
  }

  async function saveAndNext () {
    if (sessionStepRef.value) {
      await sessionStepRef.value.saveSession()
      // The completed event handler will advance the step
    }
  }

  async function finishBrewDay () {
    if (fermenterStepRef.value) {
      const success = await fermenterStepRef.value.saveFermenter()
      if (!success) return
      // The completed event handler will show completion
    }
  }

  // Close handling
  function handleDialogUpdate (value: boolean) {
    if (!value) {
      handleClose()
    }
  }

  function handleClose () {
    if (hasInProgressStep.value) {
      showConfirmClose.value = true
    } else {
      forceClose()
    }
  }

  function forceClose () {
    showConfirmClose.value = false
    emit('update:modelValue', false)
    // If any steps were completed, notify parent to refresh
    if (steps.value.some(s => s.status === 'complete')) {
      emit('completed')
    }
  }

  function handleDone () {
    emit('update:modelValue', false)
    emit('completed')
  }
</script>

<style scoped>
.action-bar {
  position: sticky;
  bottom: 0;
  background: rgb(var(--v-theme-surface));
  z-index: 1;
}
</style>
