<template>
  <v-bottom-sheet
    :inset="!xs"
    :model-value="modelValue"
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card :class="{ 'mx-auto': !xs }" :max-width="xs ? undefined : 500">
      <!-- Drag handle -->
      <div class="d-flex justify-center pt-2 pb-1">
        <div class="drag-handle" />
      </div>

      <!-- Header -->
      <v-card-title class="text-h6 pb-0">Log Reading</v-card-title>
      <v-card-subtitle v-if="subtitle" class="pt-1">
        {{ subtitle }}
      </v-card-subtitle>

      <v-card-text class="pt-4">
        <v-form :disabled="saving" @submit.prevent="handleSave">
          <!-- Gravity -->
          <v-text-field
            v-model="form.gravity"
            class="mb-2"
            density="comfortable"
            inputmode="decimal"
            label="Gravity"
            :placeholder="gravityPlaceholder"
            :rules="form.gravity ? [rules.gravity] : []"
            :suffix="gravityLabel"
            @blur="autoFormatGravity"
          />

          <!-- Temperature -->
          <v-text-field
            v-model="form.temperature"
            class="mb-2"
            density="comfortable"
            inputmode="decimal"
            label="Temperature"
            :placeholder="temperaturePlaceholder"
            :rules="form.temperature ? [rules.temperature] : []"
            :suffix="temperatureLabel"
          />

          <!-- pH (collapsible) -->
          <div v-if="!showPh" class="mb-2">
            <v-btn
              density="comfortable"
              prepend-icon="mdi-plus"
              size="small"
              variant="text"
              @click="showPh = true"
            >
              Add pH
            </v-btn>
          </div>
          <v-text-field
            v-else
            v-model="form.ph"
            class="mb-2"
            density="comfortable"
            inputmode="decimal"
            label="pH"
            placeholder="4.2"
            :rules="form.ph ? [rules.ph] : []"
            suffix="pH"
          />

          <!-- Notes -->
          <v-textarea
            v-model="form.notes"
            auto-grow
            class="mb-2"
            density="comfortable"
            label="Notes"
            placeholder="Optional notes"
            rows="2"
          />

          <!-- Observation time -->
          <div v-if="!showTimePicker" class="mb-4">
            <v-btn
              density="comfortable"
              prepend-icon="mdi-clock-outline"
              size="small"
              variant="text"
              @click="showTimePicker = true"
            >
              Change time
            </v-btn>
          </div>
          <v-text-field
            v-else
            v-model="form.observedAt"
            class="mb-4"
            density="comfortable"
            label="Observation time"
            type="datetime-local"
          />

          <!-- Save button -->
          <v-btn
            block
            color="primary"
            :disabled="!canSave"
            :loading="saving"
            min-height="44"
            size="large"
            type="submit"
          >
            Save Reading
          </v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </v-bottom-sheet>
</template>

<script lang="ts" setup>
  import { computed, ref, watch } from 'vue'
  import { useDisplay } from 'vuetify'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { convertGravity, convertTemperature, gravityLabels, temperatureLabels } from '@/composables/useUnitConversion'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { nowInputValue } from '@/utils/normalize'

  const props = defineProps<{
    modelValue: boolean
    batchUuid: string
    occupancyUuid?: string
    vesselName?: string
    batchName?: string
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'saved': []
  }>()

  const { xs } = useDisplay()
  const { createMeasurement } = useProductionApi()
  const { showNotice } = useSnackbar()
  const { preferences } = useUnitPreferences()

  const saving = ref(false)
  const showPh = ref(false)
  const showTimePicker = ref(false)

  const form = ref({
    gravity: '',
    temperature: '',
    ph: '',
    notes: '',
    observedAt: '',
  })

  // Computed display helpers
  const subtitle = computed(() => {
    const parts: string[] = []
    if (props.vesselName) parts.push(props.vesselName)
    if (props.batchName) parts.push(props.batchName)
    return parts.join(' \u00B7 ')
  })

  const gravityUnit = computed(() => preferences.value.gravity)
  const temperatureUnit = computed(() => preferences.value.temperature)

  const gravityLabel = computed(() => gravityLabels[gravityUnit.value])
  const temperatureLabel = computed(() => temperatureLabels[temperatureUnit.value])

  const gravityPlaceholder = computed(() =>
    gravityUnit.value === 'sg' ? '1.050' : '12.5',
  )

  const temperaturePlaceholder = computed(() =>
    temperatureUnit.value === 'f' ? '67' : '18',
  )

  // Validation rules
  const rules = {
    gravity: (v: string): boolean | string => {
      const num = Number.parseFloat(v)
      if (isNaN(num)) return 'Enter a valid number'
      if (gravityUnit.value === 'sg') {
        if (num < 0.99 || num > 1.2) return 'SG must be between 0.990 and 1.200'
      } else {
        if (num < 0 || num > 45) return 'Plato must be between 0 and 45'
      }
      return true
    },
    temperature: (v: string): boolean | string => {
      const num = Number.parseFloat(v)
      if (isNaN(num)) return 'Enter a valid number'
      return true
    },
    ph: (v: string): boolean | string => {
      const num = Number.parseFloat(v)
      if (isNaN(num)) return 'Enter a valid number'
      if (num < 2 || num > 14) return 'pH must be between 2.0 and 14.0'
      return true
    },
  }

  // Auto-format gravity on blur
  function autoFormatGravity () {
    const raw = form.value.gravity.trim()
    if (!raw || gravityUnit.value !== 'sg') return

    // If it already has a decimal point, leave it
    if (raw.includes('.')) return

    const num = Number.parseFloat(raw)
    if (isNaN(num)) return

    // "1050" → 1.050 (value between 900 and 1200 with no decimal)
    if (num >= 900 && num <= 1200) {
      form.value.gravity = (num / 1000).toFixed(3)
      return
    }

    // "50" → 1.050 (value between 0 and 200 with no decimal, SG mode)
    if (num >= 0 && num <= 200) {
      form.value.gravity = (1 + num / 1000).toFixed(3)
    }
  }

  // Can save if at least one measurement field has a valid value or notes are provided
  const canSave = computed(() => {
    if (saving.value) return false
    const hasGravity = form.value.gravity.trim() !== '' && rules.gravity(form.value.gravity) === true
    const hasTemp = form.value.temperature.trim() !== '' && rules.temperature(form.value.temperature) === true
    const hasPh = showPh.value && form.value.ph.trim() !== '' && rules.ph(form.value.ph) === true
    const hasNotes = form.value.notes.trim() !== ''
    return hasGravity || hasTemp || hasPh || hasNotes
  })

  // Reset form when sheet opens
  watch(
    () => props.modelValue,
    isOpen => {
      if (isOpen) {
        form.value = {
          gravity: '',
          temperature: '',
          ph: '',
          notes: '',
          observedAt: nowInputValue(),
        }
        showPh.value = false
        showTimePicker.value = false
      }
    },
  )

  async function handleSave () {
    if (!canSave.value) return

    saving.value = true

    try {
      const observedAt = form.value.observedAt
        ? new Date(form.value.observedAt).toISOString()
        : new Date().toISOString()

      const notes = form.value.notes.trim() || undefined

      // Build the target fields
      const target: { batch_uuid?: string, occupancy_uuid?: string } = {}
      if (props.occupancyUuid) {
        target.occupancy_uuid = props.occupancyUuid
      } else {
        target.batch_uuid = props.batchUuid
      }

      const promises: Promise<unknown>[] = []

      // Gravity measurement
      if (form.value.gravity.trim() && rules.gravity(form.value.gravity) === true) {
        const userValue = Number.parseFloat(form.value.gravity)
        // Convert from user's preferred unit to SG (storage unit)
        const sgValue = convertGravity(userValue, gravityUnit.value, 'sg')
        if (sgValue !== null) {
          promises.push(createMeasurement({
            ...target,
            kind: 'gravity',
            value: sgValue,
            unit: 'sg',
            observed_at: observedAt,
            notes,
          }))
        }
      }

      // Temperature measurement
      if (form.value.temperature.trim() && rules.temperature(form.value.temperature) === true) {
        const userValue = Number.parseFloat(form.value.temperature)
        // Convert from user's preferred unit to Celsius (storage unit)
        const celsiusValue = convertTemperature(userValue, temperatureUnit.value, 'c')
        if (celsiusValue !== null) {
          promises.push(createMeasurement({
            ...target,
            kind: 'temperature',
            value: celsiusValue,
            unit: 'c',
            observed_at: observedAt,
            notes,
          }))
        }
      }

      // pH measurement
      if (showPh.value && form.value.ph.trim() && rules.ph(form.value.ph) === true) {
        const phValue = Number.parseFloat(form.value.ph)
        promises.push(createMeasurement({
          ...target,
          kind: 'ph',
          value: phValue,
          unit: 'ph',
          observed_at: observedAt,
          notes,
        }))
      }

      // If no measurement promises but we have notes, create a note measurement
      if (promises.length === 0 && notes) {
        promises.push(createMeasurement({
          ...target,
          kind: 'note',
          value: 0,
          unit: null,
          observed_at: observedAt,
          notes,
        }))
      }

      if (promises.length === 0) return

      // Wait for all measurements to save
      const results = await Promise.allSettled(promises)
      const failures = results.filter(r => r.status === 'rejected')

      if (failures.length > 0) {
        const reason = (failures[0] as PromiseRejectedResult).reason
        const message = reason instanceof Error ? reason.message : 'Failed to save reading'
        showNotice(message, 'error')
        return
      }

      showNotice('Reading recorded')
      emit('saved')
      emit('update:modelValue', false)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to save reading'
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
  .drag-handle {
    width: 40px;
    height: 4px;
    border-radius: 2px;
    background-color: rgba(var(--v-theme-on-surface), 0.3);
  }
</style>
