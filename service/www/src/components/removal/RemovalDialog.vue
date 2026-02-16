<template>
  <v-dialog
    :fullscreen="smAndDown"
    :max-width="smAndDown ? '100%' : 560"
    :model-value="modelValue"
    persistent
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon class="mr-2" color="error" icon="mdi-delete-variant" />
        <span class="text-h6">Record Removal</span>
        <v-spacer />
        <v-btn
          aria-label="Close"
          :disabled="saving"
          icon="mdi-close"
          size="small"
          variant="text"
          @click="handleClose"
        />
      </v-card-title>

      <v-divider />

      <v-card-text class="pa-0" style="overflow-y: auto;">
        <v-stepper
          v-model="currentStep"
          alt-labels
          flat
          hide-actions
          :items="stepItems"
        >
          <!-- Step 1: Details -->
          <template #item.1>
            <v-container class="pa-4">
              <!-- Context banner -->
              <v-card v-if="contextLabel" class="mb-4" density="compact" variant="outlined">
                <v-card-text class="pa-3">
                  <div class="d-flex align-center ga-2">
                    <v-icon icon="mdi-link-variant" size="small" />
                    <span class="text-body-2 font-weight-medium">{{ contextLabel }}</span>
                  </div>
                </v-card-text>
              </v-card>

              <!-- Category selection -->
              <div class="text-overline text-medium-emphasis mb-2">Category</div>
              <v-item-group v-model="form.category" mandatory>
                <v-row dense>
                  <v-col
                    v-for="cat in categoryOptions"
                    :key="cat.value"
                    cols="6"
                    sm="4"
                  >
                    <v-item v-slot="{ isSelected, toggle }" :value="cat.value">
                      <v-card
                        :class="{ 'border-primary': isSelected }"
                        min-height="64"
                        :variant="isSelected ? 'tonal' : 'outlined'"
                        @click="toggle"
                      >
                        <v-card-text class="d-flex flex-column align-center pa-2">
                          <v-icon :color="cat.color" :icon="cat.icon" size="24" />
                          <span class="text-caption mt-1">{{ cat.label }}</span>
                        </v-card-text>
                      </v-card>
                    </v-item>
                  </v-col>
                </v-row>
              </v-item-group>

              <v-divider class="my-4" />

              <!-- Reason -->
              <div class="text-overline text-medium-emphasis mb-2">Reason</div>
              <v-select
                v-model="form.reason"
                density="comfortable"
                :items="filteredReasonOptions"
                label="Reason"
                :rules="[rules.required]"
              />

              <v-divider class="my-3" />

              <!-- Amount -->
              <div class="text-overline text-medium-emphasis mb-2">Amount</div>
              <v-row dense>
                <v-col cols="7">
                  <v-text-field
                    v-model="form.amount"
                    density="comfortable"
                    inputmode="decimal"
                    label="Amount"
                    min="0"
                    :rules="form.amount ? [rules.positiveNumber] : [rules.required]"
                    type="number"
                  />
                </v-col>
                <v-col cols="5">
                  <v-select
                    v-model="form.amountUnit"
                    density="comfortable"
                    item-title="label"
                    item-value="value"
                    :items="volumeUnitOptions"
                    label="Unit"
                  />
                </v-col>
              </v-row>

              <v-divider class="my-3" />

              <!-- Date/Time -->
              <div class="text-overline text-medium-emphasis mb-2">When</div>
              <v-text-field
                v-model="form.removedAt"
                density="comfortable"
                label="Date/time"
                type="datetime-local"
              />

              <v-divider class="my-3" />

              <!-- Notes -->
              <v-textarea
                v-model="form.notes"
                auto-grow
                density="comfortable"
                label="Notes (optional)"
                rows="2"
              />
            </v-container>
          </template>

          <!-- Step 2: Review & Confirm -->
          <template #item.2>
            <v-container class="pa-4">
              <div class="text-overline text-medium-emphasis mb-3">Review Removal</div>

              <!-- Context -->
              <div v-if="contextLabel" class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Source</div>
                <div class="text-body-1 font-weight-medium">{{ contextLabel }}</div>
              </div>

              <!-- Category & Reason -->
              <div class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Category</div>
                <div class="d-flex align-center ga-2">
                  <v-chip
                    :color="selectedCategoryColor"
                    :prepend-icon="selectedCategoryIcon"
                    size="small"
                    variant="tonal"
                  >
                    {{ selectedCategoryLabel }}
                  </v-chip>
                </div>
              </div>

              <div class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Reason</div>
                <div class="text-body-2">{{ selectedReasonLabel }}</div>
              </div>

              <!-- Amount -->
              <div class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Amount</div>
                <div class="text-body-1 font-weight-medium">
                  {{ form.amount }} {{ volumeLabelForUnit(form.amountUnit) }}
                </div>
                <div v-if="estimatedBbl !== null" class="text-body-2 text-medium-emphasis">
                  &asymp; {{ estimatedBbl.toFixed(2) }} BBL
                </div>
              </div>

              <!-- Date -->
              <div class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Date</div>
                <div class="text-body-2">{{ reviewDateLabel }}</div>
              </div>

              <!-- Notes -->
              <div v-if="form.notes" class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Notes</div>
                <div class="text-body-2">{{ form.notes }}</div>
              </div>

              <!-- Error display -->
              <v-alert
                v-if="saveError"
                class="mb-4"
                density="compact"
                type="error"
                variant="tonal"
              >
                {{ saveError }}
              </v-alert>
            </v-container>
          </template>
        </v-stepper>
      </v-card-text>

      <v-divider />

      <!-- Navigation actions -->
      <v-card-actions class="justify-space-between pa-4">
        <v-btn
          v-if="currentStep > 1"
          :disabled="saving"
          variant="text"
          @click="currentStep -= 1"
        >
          &larr; Back
        </v-btn>
        <v-spacer v-else />

        <div>
          <v-btn
            :disabled="saving"
            variant="text"
            @click="handleClose"
          >
            Cancel
          </v-btn>
          <v-btn
            v-if="currentStep < 2"
            color="error"
            :disabled="!canProceed"
            @click="currentStep += 1"
          >
            Next &rarr;
          </v-btn>
          <v-btn
            v-else
            color="error"
            :disabled="!canProceed"
            :loading="saving"
            @click="handleConfirm"
          >
            Record Removal
          </v-btn>
        </div>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { CreateRemovalRequest, Removal, RemovalCategory, RemovalReason, VolumeUnit } from '@/types'
  import { computed, ref, watch } from 'vue'
  import { useDisplay } from 'vuetify'
  import { formatDateTime } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { convertVolume, volumeLabels } from '@/composables/useUnitConversion'
  import { volumeOptions, useUnitPreferences } from '@/composables/useUnitPreferences'
  import { nowInputValue } from '@/utils/normalize'
  import { categoryColors, categoryIcons, categoryLabels, reasonLabels, reasonsByCategory } from './constants'

  // ==================== Props & Emits ====================

  const props = withDefaults(
    defineProps<{
      modelValue: boolean
      batchUuid?: string
      batchName?: string
      beerLotUuid?: string
      beerLotCode?: string
      occupancyUuid?: string
      stockLocationUuid?: string
      defaultCategory?: RemovalCategory
    }>(),
    {
      batchUuid: undefined,
      batchName: undefined,
      beerLotUuid: undefined,
      beerLotCode: undefined,
      occupancyUuid: undefined,
      stockLocationUuid: undefined,
      defaultCategory: undefined,
    },
  )

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'created': [removal: Removal]
  }>()

  // ==================== Composables ====================

  const { smAndDown } = useDisplay()
  const { createRemoval } = useInventoryApi()
  const { showNotice } = useSnackbar()
  const { preferences } = useUnitPreferences()

  // ==================== State ====================

  const currentStep = ref(1)
  const stepItems = ['Details', 'Review & Confirm']
  const saving = ref(false)
  const saveError = ref('')

  const form = ref({
    category: 'dump' as RemovalCategory,
    reason: '' as RemovalReason | '',
    amount: '',
    amountUnit: 'bbl' as VolumeUnit,
    removedAt: '',
    notes: '',
  })

  // ==================== Options ====================

  const categoryOptions = [
    { value: 'dump' as RemovalCategory, label: 'Batch Dump', icon: categoryIcons.dump, color: categoryColors.dump },
    { value: 'waste' as RemovalCategory, label: 'Waste', icon: categoryIcons.waste, color: categoryColors.waste },
    { value: 'sample' as RemovalCategory, label: 'Sample', icon: categoryIcons.sample, color: categoryColors.sample },
    { value: 'expired' as RemovalCategory, label: 'Expired', icon: categoryIcons.expired, color: categoryColors.expired },
    { value: 'other' as RemovalCategory, label: 'Other', icon: categoryIcons.other, color: categoryColors.other },
  ]

  const volumeUnitOptions = volumeOptions

  const rules = {
    required: (v: string) => !!v || 'Required',
    positiveNumber: (v: string) => {
      const num = Number.parseFloat(v)
      if (isNaN(num)) return 'Enter a valid number'
      if (num <= 0) return 'Must be greater than 0'
      return true
    },
  }

  // ==================== Computed ====================

  const contextLabel = computed(() => {
    const parts: string[] = []
    if (props.batchName) {
      parts.push(props.batchName)
    } else if (props.batchUuid) {
      parts.push(`Batch ${props.batchUuid.slice(0, 8)}`)
    }
    if (props.beerLotCode) {
      parts.push(`Lot ${props.beerLotCode}`)
    } else if (props.beerLotUuid) {
      parts.push(`Lot ${props.beerLotUuid.slice(0, 8)}`)
    }
    return parts.join(' · ')
  })

  const filteredReasonOptions = computed(() => {
    const reasons = reasonsByCategory[form.value.category] ?? ['other']
    return reasons.map(r => ({
      title: reasonLabels[r],
      value: r,
    }))
  })

  const selectedCategoryLabel = computed(() => categoryLabels[form.value.category] ?? form.value.category)
  const selectedCategoryColor = computed(() => categoryColors[form.value.category] ?? 'grey')
  const selectedCategoryIcon = computed(() => categoryIcons[form.value.category] ?? 'mdi-dots-horizontal')
  const selectedReasonLabel = computed(() => {
    if (!form.value.reason) return '—'
    return reasonLabels[form.value.reason as RemovalReason] ?? form.value.reason
  })

  const estimatedBbl = computed(() => {
    const amount = Number.parseFloat(form.value.amount)
    if (isNaN(amount) || amount <= 0) return null
    if (form.value.amountUnit === 'bbl') return amount
    const converted = convertVolume(amount, form.value.amountUnit, 'bbl')
    return converted
  })

  const reviewDateLabel = computed(() =>
    form.value.removedAt
      ? formatDateTime(new Date(form.value.removedAt).toISOString())
      : 'Now',
  )

  const canProceed = computed(() => {
    if (currentStep.value === 1) return canProceedStep1.value
    if (currentStep.value === 2) return canProceedStep1.value
    return false
  })

  const canProceedStep1 = computed(() => {
    if (!form.value.category) return false
    if (!form.value.reason) return false
    const amount = Number.parseFloat(form.value.amount)
    if (isNaN(amount) || amount <= 0) return false
    return true
  })

  // ==================== Helpers ====================

  function volumeLabelForUnit (unit: VolumeUnit): string {
    return volumeLabels[unit] ?? unit
  }

  // ==================== Watch ====================

  // Reset reason when category changes
  watch(
    () => form.value.category,
    () => {
      const validReasons = reasonsByCategory[form.value.category] ?? ['other']
      if (!validReasons.includes(form.value.reason as RemovalReason)) {
        form.value.reason = validReasons[0] ?? 'other'
      }
    },
  )

  // Reset form when dialog opens
  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        resetForm()
      }
    },
  )

  function resetForm () {
    currentStep.value = 1
    saving.value = false
    saveError.value = ''

    const category = props.defaultCategory ?? 'dump'
    const validReasons = reasonsByCategory[category] ?? ['other']

    form.value = {
      category,
      reason: validReasons[0] ?? 'other',
      amount: '',
      amountUnit: preferences.value.volume,
      removedAt: nowInputValue(),
      notes: '',
    }
  }

  // ==================== Actions ====================

  function handleClose () {
    emit('update:modelValue', false)
  }

  async function handleConfirm () {
    saving.value = true
    saveError.value = ''

    try {
      const amount = Number.parseFloat(form.value.amount)
      if (isNaN(amount) || amount <= 0) {
        saveError.value = 'Invalid amount'
        return
      }

      // Backend expects integer amounts — round to nearest whole number
      const payload: CreateRemovalRequest = {
        category: form.value.category,
        reason: form.value.reason as RemovalReason,
        amount: Math.round(amount),
        amount_unit: form.value.amountUnit,
        removed_at: form.value.removedAt
          ? new Date(form.value.removedAt).toISOString()
          : undefined,
        batch_uuid: props.batchUuid || undefined,
        beer_lot_uuid: props.beerLotUuid || undefined,
        occupancy_uuid: props.occupancyUuid || undefined,
        stock_location_uuid: props.stockLocationUuid || undefined,
        notes: form.value.notes.trim() || undefined,
      }

      const removal = await createRemoval(payload)
      showNotice('Removal recorded')
      emit('created', removal)
      emit('update:modelValue', false)
    } catch (error) {
      saveError.value = error instanceof Error ? error.message : 'Failed to record removal'
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
.review-section {
  padding: 8px 0;
}

.border-primary {
  border-color: rgb(var(--v-theme-primary)) !important;
}
</style>
