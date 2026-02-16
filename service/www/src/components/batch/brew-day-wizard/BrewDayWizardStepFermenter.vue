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

    <!-- Vessel picker -->
    <div class="text-overline text-medium-emphasis mb-2">Select Fermenter</div>

    <div v-if="fermenterOptions.length > 0" aria-label="Select fermenter" role="radiogroup">
      <v-row dense>
        <v-col
          v-for="vessel in fermenterOptions"
          :key="vessel.uuid"
          cols="6"
          sm="4"
        >
          <v-card
            :aria-checked="form.vesselUuid === vessel.uuid"
            :aria-label="`${vessel.name}`"
            class="vessel-card"
            :class="{
              'vessel-card--selected': form.vesselUuid === vessel.uuid,
              'vessel-card--occupied': vessel.occupied,
            }"
            :disabled="vessel.occupied"
            :ripple="!vessel.occupied"
            role="radio"
            variant="outlined"
            @click="!vessel.occupied && selectVessel(vessel.uuid)"
          >
            <v-card-text class="pa-3 text-center">
              <v-icon
                v-if="form.vesselUuid === vessel.uuid"
                class="mb-1"
                color="primary"
                icon="mdi-check-circle"
                size="24"
              />
              <v-icon
                v-else-if="vessel.occupied"
                class="mb-1"
                color="grey"
                icon="mdi-flask-round-bottom"
                size="24"
              />
              <v-icon
                v-else
                class="mb-1"
                icon="mdi-flask-round-bottom-empty-outline"
                size="24"
              />
              <div class="text-body-2 font-weight-medium">{{ vessel.name }}</div>
              <div class="text-caption text-medium-emphasis">
                {{ vessel.capacity }} {{ vessel.capacityUnit }}
              </div>
              <v-chip
                v-if="vessel.occupied"
                class="mt-1"
                color="grey"
                size="x-small"
                variant="tonal"
              >
                In use
              </v-chip>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </div>

    <v-alert
      v-else
      class="mb-3"
      density="comfortable"
      type="info"
      variant="tonal"
    >
      No fermenters configured. Add vessels in the Vessels section.
    </v-alert>

    <!-- Volume select -->
    <template v-if="volumeOptions.length > 1">
      <div class="text-overline text-medium-emphasis mt-4 mb-2">Volume</div>
      <v-select
        v-model="form.volumeUuid"
        density="comfortable"
        hide-details
        item-title="title"
        item-value="value"
        :items="volumeOptions"
        placeholder="Select volume"
      />
    </template>

    <!-- Status chips -->
    <div class="text-overline text-medium-emphasis mt-4 mb-2">Initial Status</div>
    <v-chip-group
      v-model="form.status"
      color="primary"
      mandatory
    >
      <v-chip
        v-for="option in statusOptions"
        :key="option.value"
        filter
        min-height="44"
        :value="option.value"
      >
        {{ option.title }}
      </v-chip>
    </v-chip-group>
  </div>
</template>

<script lang="ts" setup>
  import type { Occupancy, OccupancyStatus, Vessel, Volume } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'

  interface FermenterOption {
    uuid: string
    name: string
    capacity: number
    capacityUnit: string
    occupied: boolean
  }

  const props = defineProps<{
    batchUuid: string
    vessels: Vessel[]
    volumes: Volume[]
    occupancies: Occupancy[]
  }>()

  const emit = defineEmits<{
    completed: [data: { occupancy: Occupancy, fermenterName: string }]
  }>()

  const { createOccupancy } = useProductionApi()
  const { showNotice } = useSnackbar()

  const saving = ref(false)
  const errorMessage = ref('')

  const form = reactive({
    vesselUuid: '' as string,
    volumeUuid: '' as string,
    status: 'fermenting' as OccupancyStatus,
  })

  const statusOptions: { title: string, value: OccupancyStatus }[] = [
    { title: 'Fermenting', value: 'fermenting' },
    { title: 'Conditioning', value: 'conditioning' },
    { title: 'Cold Crashing', value: 'cold_crashing' },
    { title: 'Dry Hopping', value: 'dry_hopping' },
    { title: 'Carbonating', value: 'carbonating' },
    { title: 'Holding', value: 'holding' },
    { title: 'Packaging', value: 'packaging' },
  ]

  const occupiedVesselUuids = computed(() =>
    new Set(props.occupancies.map(o => o.vessel_uuid)),
  )

  const fermenterOptions = computed<FermenterOption[]>(() =>
    props.vessels
      .filter(v => v.status === 'active' && (v.type === 'fermenter' || v.type === 'brite_tank'))
      .map(v => ({
        uuid: v.uuid,
        name: v.name,
        capacity: v.capacity,
        capacityUnit: v.capacity_unit,
        occupied: occupiedVesselUuids.value.has(v.uuid),
      })),
  )

  const volumeOptions = computed(() =>
    props.volumes.map(v => ({
      title: v.name
        ? `${v.name} (${v.amount} ${v.amount_unit})`
        : `Volume (${v.amount} ${v.amount_unit})`,
      value: v.uuid,
    })),
  )

  const isFormValid = computed(() =>
    !!form.vesselUuid && !!form.volumeUuid,
  )

  // Auto-select volume if only one
  onMounted(() => {
    if (props.volumes.length === 1 && props.volumes[0]) {
      form.volumeUuid = props.volumes[0].uuid
    }
  })

  // Expose for parent
  defineExpose({
    saveFermenter,
    saving,
    isFormValid,
  })

  function selectVessel (uuid: string) {
    form.vesselUuid = uuid
  }

  async function saveFermenter (): Promise<boolean> {
    if (!isFormValid.value) {
      errorMessage.value = 'Please select a fermenter and volume.'
      return false
    }

    saving.value = true
    errorMessage.value = ''

    try {
      const payload = {
        vessel_uuid: form.vesselUuid,
        volume_uuid: form.volumeUuid,
        in_at: new Date().toISOString(),
        status: form.status || undefined,
      }

      const occupancy = await createOccupancy(payload)

      const vessel = props.vessels.find(v => v.uuid === form.vesselUuid)
      const fermenterName = vessel?.name ?? 'Fermenter'

      showNotice(`Assigned to ${fermenterName}`)

      emit('completed', {
        occupancy,
        fermenterName,
      })

      return true
    } catch (error: unknown) {
      const message = error instanceof Error ? error.message : 'Failed to assign fermenter'
      errorMessage.value = message
      showNotice(message, 'error')
      return false
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
.vessel-card {
  cursor: pointer;
  transition: all 0.2s ease;
  min-height: 100px;
}

.vessel-card:hover:not(.vessel-card--occupied) {
  border-color: rgb(var(--v-theme-primary));
}

.vessel-card--selected {
  border-color: rgb(var(--v-theme-primary)) !important;
  border-width: 2px;
  background: rgba(var(--v-theme-primary), 0.05);
}

.vessel-card--occupied {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
