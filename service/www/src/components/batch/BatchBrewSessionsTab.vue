<template>
  <v-row>
    <v-col cols="12">
      <v-card class="sub-card" variant="outlined">
        <v-card-title class="text-subtitle-1 d-flex align-center">
          Brew Sessions
          <v-spacer />
          <v-btn
            aria-label="Add brew session"
            icon="mdi-plus"
            size="small"
            variant="text"
            @click="emit('create-session')"
          />
        </v-card-title>
        <v-card-text>
          <v-alert
            v-if="sessions.length === 0"
            class="mb-3"
            density="compact"
            type="info"
            variant="tonal"
          >
            No brew sessions recorded. Add a brew session to track hot-side additions and measurements.
          </v-alert>
          <v-list v-else class="brew-session-list" lines="three">
            <v-list-item
              v-for="session in sessionsSorted"
              :key="session.uuid"
              :active="session.uuid === selectedSessionUuid"
              @click="emit('select-session', session.uuid)"
            >
              <v-list-item-title>
                {{ formatDateTime(session.brewed_at) }}
              </v-list-item-title>
              <v-list-item-subtitle>
                <span v-if="getVesselName(session.mash_vessel_uuid)">
                  Mash: {{ getVesselName(session.mash_vessel_uuid) }}
                </span>
                <span v-if="getVesselName(session.boil_vessel_uuid)">
                  &bull; Boil: {{ getVesselName(session.boil_vessel_uuid) }}
                </span>
                <span v-if="getVolumeName(session.wort_volume_uuid)">
                  &bull; {{ getVolumeName(session.wort_volume_uuid) }}
                </span>
              </v-list-item-subtitle>
              <v-list-item-subtitle v-if="session.notes" class="text-medium-emphasis">
                {{ session.notes }}
              </v-list-item-subtitle>
              <template #append>
                <v-btn
                  icon="mdi-pencil"
                  size="x-small"
                  variant="text"
                  @click.stop="emit('edit-session', session)"
                />
              </template>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>

  <!-- Selected Brew Session Details -->
  <v-row v-if="selectedSession" class="mt-4">
    <v-col cols="12">
      <v-card class="sub-card" variant="tonal">
        <v-card-title class="text-subtitle-1 d-flex align-center">
          <v-icon class="mr-2" icon="mdi-kettle-steam" size="small" />
          {{ formatDateTime(selectedSession.brewed_at) }}
          <v-spacer />
          <v-btn size="small" variant="text" @click="emit('clear-session')">
            Clear
          </v-btn>
        </v-card-title>
        <v-card-text>
          <v-row dense>
            <v-col v-if="getVesselName(selectedSession.mash_vessel_uuid)" cols="12" md="4">
              <div class="text-caption text-medium-emphasis">Mash Vessel</div>
              <div class="text-body-2 font-weight-medium">
                {{ getVesselName(selectedSession.mash_vessel_uuid) }}
              </div>
            </v-col>
            <v-col v-if="getVesselName(selectedSession.boil_vessel_uuid)" cols="12" md="4">
              <div class="text-caption text-medium-emphasis">Boil Vessel</div>
              <div class="text-body-2 font-weight-medium">
                {{ getVesselName(selectedSession.boil_vessel_uuid) }}
              </div>
            </v-col>
            <v-col v-if="getVolumeName(selectedSession.wort_volume_uuid)" cols="12" md="4">
              <div class="text-caption text-medium-emphasis">Wort Volume</div>
              <div class="text-body-2 font-weight-medium">
                {{ getVolumeName(selectedSession.wort_volume_uuid) }}
                <span v-if="getVolumeAmount(selectedSession.wort_volume_uuid)" class="text-medium-emphasis">
                  ({{ getVolumeAmount(selectedSession.wort_volume_uuid) }})
                </span>
              </div>
            </v-col>
          </v-row>
          <div v-if="selectedSession.notes" class="mt-3">
            <div class="text-caption text-medium-emphasis">Notes</div>
            <div class="text-body-2">{{ selectedSession.notes }}</div>
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <!-- Hot-Side Additions for Selected Brew Session -->
    <v-col cols="12" md="6">
      <v-card class="sub-card" variant="outlined">
        <v-card-title class="text-subtitle-2 d-flex align-center">
          Hot-Side Additions
          <v-spacer />
          <v-btn
            aria-label="Add hot-side addition"
            :disabled="!selectedSession.wort_volume_uuid"
            icon="mdi-plus"
            size="x-small"
            variant="text"
            @click="emit('create-addition')"
          />
        </v-card-title>
        <v-card-text>
          <v-table class="data-table" density="compact">
            <thead>
              <tr>
                <th>Type</th>
                <th>Amount</th>
                <th>Time</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="addition in additionsSorted" :key="addition.uuid">
                <td>
                  <v-chip size="x-small" variant="tonal">{{ formatAdditionType(addition.addition_type) }}</v-chip>
                  <span v-if="addition.stage" class="text-medium-emphasis ml-1">{{ addition.stage }}</span>
                </td>
                <td>{{ formatAmount(addition.amount, addition.amount_unit) }}</td>
                <td>{{ formatDateTime(addition.added_at) }}</td>
              </tr>
              <tr v-if="additionsSorted.length === 0">
                <td class="text-medium-emphasis" colspan="3">
                  {{ selectedSession.wort_volume_uuid ? 'No additions recorded.' : 'Select a wort volume first.' }}
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-card-text>
      </v-card>
    </v-col>

    <!-- Hot-Side Measurements for Selected Brew Session -->
    <v-col cols="12" md="6">
      <v-card class="sub-card" variant="outlined">
        <v-card-title class="text-subtitle-2 d-flex align-center">
          Hot-Side Measurements
          <v-spacer />
          <v-btn
            aria-label="Add hot-side measurement"
            :disabled="!selectedSession.wort_volume_uuid"
            icon="mdi-plus"
            size="x-small"
            variant="text"
            @click="emit('create-measurement')"
          />
        </v-card-title>
        <v-card-text>
          <v-table class="data-table" density="compact">
            <thead>
              <tr>
                <th>Kind</th>
                <th>Value</th>
                <th>Time</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="measurement in measurementsSorted" :key="measurement.uuid">
                <td>{{ formatMeasurementKind(measurement.kind) }}</td>
                <td>{{ formatValue(measurement.value, measurement.unit) }}</td>
                <td>{{ formatDateTime(measurement.observed_at) }}</td>
              </tr>
              <tr v-if="measurementsSorted.length === 0">
                <td class="text-medium-emphasis" colspan="3">
                  {{ selectedSession.wort_volume_uuid ? 'No measurements recorded.' : 'Select a wort volume first.' }}
                </td>
              </tr>
            </tbody>
          </v-table>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts" setup>
  import type {
    BrewSession,
    Addition as ProductionAddition,
    Measurement as ProductionMeasurement,
    Volume as ProductionVolume,
    Vessel,
  } from '@/types'
  import { computed } from 'vue'
  import { useAdditionTypeFormatters, useFormatters } from '@/composables/useFormatters'

  const props = defineProps<{
    sessions: BrewSession[]
    selectedSessionUuid: string | null
    vessels: Vessel[]
    volumes: ProductionVolume[]
    additions: ProductionAddition[]
    measurements: ProductionMeasurement[]
  }>()

  const emit = defineEmits<{
    'create-session': []
    'edit-session': [session: BrewSession]
    'select-session': [uuid: string]
    'clear-session': []
    'create-addition': []
    'create-measurement': []
  }>()

  const { formatDateTime } = useFormatters()
  const { formatAdditionType } = useAdditionTypeFormatters()

  const sessionsSorted = computed(() => {
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...props.sessions].sort(
      (a, b) => new Date(b.brewed_at).getTime() - new Date(a.brewed_at).getTime(),
    )
  })

  const selectedSession = computed(() =>
    props.sessions.find(session => session.uuid === props.selectedSessionUuid) ?? null,
  )

  const additionsSorted = computed(() => {
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...props.additions].sort(
      (a, b) => new Date(b.added_at).getTime() - new Date(a.added_at).getTime(),
    )
  })

  const measurementsSorted = computed(() => {
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...props.measurements].sort(
      (a, b) => new Date(b.observed_at).getTime() - new Date(a.observed_at).getTime(),
    )
  })

  function getVesselName (vesselUuid: string | null): string {
    if (!vesselUuid) return ''
    const vessel = props.vessels.find(v => v.uuid === vesselUuid)
    return vessel?.name ?? 'Unknown Vessel'
  }

  function getVolumeName (volumeUuid: string | null): string {
    if (!volumeUuid) return ''
    const volume = props.volumes.find(v => v.uuid === volumeUuid)
    return volume?.name ?? 'Unknown Volume'
  }

  function getVolumeAmount (volumeUuid: string | null): string {
    if (!volumeUuid) return ''
    const volume = props.volumes.find(v => v.uuid === volumeUuid)
    if (!volume) return ''
    return `${volume.amount} ${volume.amount_unit}`
  }

  function formatAmount (amount: number | null, unit: string | null | undefined) {
    if (amount === null || amount === undefined) {
      return 'Unknown'
    }
    return `${amount} ${unit ?? ''}`.trim()
  }

  function formatValue (value: number | null, unit: string | null | undefined) {
    if (value === null || value === undefined) {
      return 'Unknown'
    }
    return `${value}${unit ? ` ${unit}` : ''}`
  }

  function formatMeasurementKind (kind: string) {
    const normalized = kind.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
    if (normalized === 'ph') {
      return 'pH'
    }
    if (normalized === 'sg') {
      return 'SG'
    }
    if (normalized === 'temperature' || normalized === 'temp') {
      return 'Temp'
    }
    if (normalized === 'gravity' || normalized === 'grav') {
      return 'Gravity'
    }
    const trimmed = kind.trim()
    if (!trimmed) {
      return 'Measurement'
    }
    return trimmed.charAt(0).toUpperCase() + trimmed.slice(1)
  }
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.brew-session-list {
  max-height: 280px;
  overflow: auto;
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
}

.data-table :deep(td) {
  font-size: 0.85rem;
}
</style>
