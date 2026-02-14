<template>
  <div>
    <v-alert
      v-if="errorMessage"
      class="mb-4"
      density="compact"
      type="error"
      variant="tonal"
    >
      {{ errorMessage }}
    </v-alert>

    <!-- Mash Vessel -->
    <v-card class="mb-3" variant="outlined">
      <v-card-text class="pa-3">
        <div class="text-overline text-medium-emphasis mb-2">Mash Vessel</div>
        <v-select
          v-model="form.mashVesselUuid"
          density="comfortable"
          hide-details
          :items="mashVesselOptions"
          item-title="name"
          item-value="uuid"
          placeholder="Select mash vessel"
        >
          <template #no-data>
            <v-list-item>
              <v-list-item-title>No mash vessels available</v-list-item-title>
            </v-list-item>
          </template>
        </v-select>
      </v-card-text>
    </v-card>

    <!-- Boil Vessel -->
    <v-card class="mb-3" variant="outlined">
      <v-card-text class="pa-3">
        <div class="text-overline text-medium-emphasis mb-2">Boil Vessel</div>
        <v-select
          v-model="form.boilVesselUuid"
          density="comfortable"
          hide-details
          :items="boilVesselOptions"
          item-title="name"
          item-value="uuid"
          placeholder="Select boil vessel"
        >
          <template #no-data>
            <v-list-item>
              <v-list-item-title>No kettles available</v-list-item-title>
            </v-list-item>
          </template>
        </v-select>
      </v-card-text>
    </v-card>

    <!-- Key Measurements -->
    <v-card class="mb-3" variant="outlined">
      <v-card-text class="pa-3">
        <div class="text-overline text-medium-emphasis mb-3">Key Measurements</div>
        <v-row dense>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="form.mashTemp"
              density="comfortable"
              hide-details
              inputmode="decimal"
              label="Mash Temperature"
              :suffix="temperatureSuffix"
              type="number"
            />
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="form.preBoilGravity"
              density="comfortable"
              hide-details
              inputmode="decimal"
              label="Pre-Boil Gravity"
              placeholder="1.040"
              suffix="SG"
              type="number"
            />
          </v-col>
          <v-col cols="12" sm="4">
            <v-text-field
              v-model="form.originalGravity"
              density="comfortable"
              hide-details
              inputmode="decimal"
              label="Original Gravity"
              placeholder="1.050"
              suffix="SG"
              type="number"
            />
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <!-- Notes -->
    <v-card variant="outlined">
      <v-card-text class="pa-3">
        <v-textarea
          v-model="form.notes"
          auto-grow
          density="comfortable"
          hide-details
          label="Notes"
          placeholder="Mash temps, boil notes, etc."
          rows="2"
        />
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts" setup>
  import type { BrewSession, Vessel, Volume } from '@/types'
  import { computed, reactive, ref, watch } from 'vue'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { normalizeText, toNumber } from '@/utils/normalize'

  const props = defineProps<{
    batchUuid: string
    vessels: Vessel[]
    volumes: Volume[]
  }>()

  const emit = defineEmits<{
    completed: [data: { session: BrewSession, mashTemp: string | null, originalGravity: string | null }]
  }>()

  const { createBrewSession, createMeasurement } = useProductionApi()
  const { showNotice } = useSnackbar()
  const { preferences } = useUnitPreferences()

  const saving = ref(false)
  const errorMessage = ref('')

  const form = reactive({
    mashVesselUuid: null as string | null,
    boilVesselUuid: null as string | null,
    mashTemp: '',
    preBoilGravity: '',
    originalGravity: '',
    notes: '',
  })

  const temperatureSuffix = computed(() =>
    preferences.value.temperature === 'c' ? '°C' : '°F',
  )

  const mashVesselOptions = computed(() =>
    props.vessels
      .filter(v => v.status === 'active' && (v.type === 'mash_tun' || v.type === 'lauter_tun'))
      .map(v => ({ uuid: v.uuid, name: v.name })),
  )

  const boilVesselOptions = computed(() =>
    props.vessels
      .filter(v => v.status === 'active' && (v.type === 'kettle' || v.type === 'whirlpool'))
      .map(v => ({ uuid: v.uuid, name: v.name })),
  )

  // Auto-select if only one option
  watch(mashVesselOptions, (options) => {
    if (options.length === 1 && options[0] && !form.mashVesselUuid) {
      form.mashVesselUuid = options[0].uuid
    }
  }, { immediate: true })

  watch(boilVesselOptions, (options) => {
    if (options.length === 1 && options[0] && !form.boilVesselUuid) {
      form.boilVesselUuid = options[0].uuid
    }
  }, { immediate: true })

  // Expose for parent
  defineExpose({
    saveSession,
    saving,
  })

  async function saveSession (): Promise<boolean> {
    saving.value = true
    errorMessage.value = ''

    try {
      // Create brew session
      const sessionPayload = {
        batch_uuid: props.batchUuid,
        mash_vessel_uuid: form.mashVesselUuid,
        boil_vessel_uuid: form.boilVesselUuid,
        brewed_at: new Date().toISOString(),
        notes: normalizeText(form.notes),
      }

      const session = await createBrewSession(sessionPayload)

      // Create measurements for any entered values
      const observedAt = new Date().toISOString()
      const measurementPromises: Promise<unknown>[] = []

      const mashTempValue = toNumber(form.mashTemp)
      if (mashTempValue !== null) {
        measurementPromises.push(
          createMeasurement({
            batch_uuid: props.batchUuid,
            kind: 'mash_temp',
            value: mashTempValue,
            unit: preferences.value.temperature === 'c' ? 'C' : 'F',
            observed_at: observedAt,
          }),
        )
      }

      const preBoilGravityValue = toNumber(form.preBoilGravity)
      if (preBoilGravityValue !== null) {
        measurementPromises.push(
          createMeasurement({
            batch_uuid: props.batchUuid,
            kind: 'pre_boil_gravity',
            value: preBoilGravityValue,
            unit: 'SG',
            observed_at: observedAt,
          }),
        )
      }

      const ogValue = toNumber(form.originalGravity)
      if (ogValue !== null) {
        measurementPromises.push(
          createMeasurement({
            batch_uuid: props.batchUuid,
            kind: 'original_gravity',
            value: ogValue,
            unit: 'SG',
            observed_at: observedAt,
          }),
        )
      }

      // Non-blocking: measurements are nice-to-have
      if (measurementPromises.length > 0) {
        await Promise.allSettled(measurementPromises)
      }

      showNotice('Brew session recorded')

      emit('completed', {
        session,
        mashTemp: mashTempValue !== null ? `${mashTempValue}${temperatureSuffix.value}` : null,
        originalGravity: ogValue !== null ? ogValue.toFixed(3) : null,
      })

      return true
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to save brew session'
      errorMessage.value = message
      showNotice(message, 'error')
      return false
    } finally {
      saving.value = false
    }
  }
</script>
