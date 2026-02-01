<template>
  <v-container class="dashboard-page" fluid>
    <v-card class="section-card hero-card">
      <v-card-text>
        <div class="dashboard-hero">
          <div>
            <div class="text-overline text-medium-emphasis">Dashboard</div>
            <div class="text-h4 font-weight-semibold">Brew day overview</div>
            <div class="text-body-2 text-medium-emphasis">
              Track in-progress batches, upcoming brew days, and vessel readiness.
            </div>
          </div>
          <div class="hero-actions">
            <v-btn color="primary" to="/batches">Start batch</v-btn>
            <v-btn variant="tonal" to="/batches">Open workflow</v-btn>
            <v-btn :loading="loading" variant="text" @click="refreshAll">Refresh</v-btn>
          </div>
        </div>

        <v-alert v-if="errorMessage" class="mt-4" density="compact" type="error" variant="tonal">
          {{ errorMessage }}
        </v-alert>
      </v-card-text>
    </v-card>

    <v-row class="mt-4" align="stretch">
      <v-col cols="12" lg="7">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-progress-clock" />
            In-progress batches
            <v-chip class="ml-3" size="x-small" variant="tonal" color="primary">
              {{ inProgressCount }}
            </v-chip>
            <v-spacer />
            <v-btn size="small" variant="text" to="/batches">View all</v-btn>
          </v-card-title>
          <v-card-text>
            <v-list class="batch-list" lines="two">
              <v-list-item v-for="item in inProgressBatches" :key="item.batch.id" :to="item.route">
                <v-list-item-title>{{ item.batch.short_name }}</v-list-item-title>
                <v-list-item-subtitle>
                  #{{ item.batch.id }} · Brewed {{ formatBrewDate(item.batch.brew_date) }}
                </v-list-item-subtitle>
                <div class="text-caption text-medium-emphasis">
                  {{ item.phaseLabel }} · Updated {{ formatDateTime(item.phaseAt) }}
                </div>
                <template #append>
                  <div class="d-flex flex-column align-end ga-1">
                    <v-chip size="x-small" variant="tonal" :color="item.phaseTone">
                      {{ item.phaseLabel }}
                    </v-chip>
                    <div class="text-caption text-medium-emphasis">
                      {{ formatDate(item.batch.updated_at) }}
                    </div>
                  </div>
                </template>
              </v-list-item>

              <v-list-item v-if="inProgressCount === 0">
                <v-list-item-title>No active batches</v-list-item-title>
                <v-list-item-subtitle>Move a batch out of planning to start brewing.</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>

        <v-card class="section-card mt-4">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-calendar-clock" />
            Planned batches
            <v-chip class="ml-3" size="x-small" variant="tonal" color="info">
              {{ plannedCount }}
            </v-chip>
            <v-spacer />
            <v-btn size="small" variant="text" to="/batches">Plan schedule</v-btn>
          </v-card-title>
          <v-card-text>
            <v-list class="batch-list" lines="two">
              <v-list-item v-for="item in plannedBatches" :key="item.batch.id" :to="item.route">
                <v-list-item-title>{{ item.batch.short_name }}</v-list-item-title>
                <v-list-item-subtitle>
                  #{{ item.batch.id }} · {{ plannedSubtitle(item) }}
                </v-list-item-subtitle>
                <div class="text-caption text-medium-emphasis">
                  {{ plannedMeta(item) }}
                </div>
                <template #append>
                  <div class="d-flex flex-column align-end ga-1">
                    <v-chip size="x-small" variant="tonal" color="info">
                      {{ plannedBadge(item) }}
                    </v-chip>
                    <div class="text-caption text-medium-emphasis">
                      {{ formatBrewDate(item.batch.brew_date) }}
                    </div>
                  </div>
                </template>
              </v-list-item>

              <v-list-item v-if="plannedCount === 0">
                <v-list-item-title>No upcoming batches</v-list-item-title>
                <v-list-item-subtitle>Schedule the next brew to keep the cellar full.</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" lg="5">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-silo" />
            Vessels
            <v-spacer />
            <v-btn size="small" variant="text" to="/vessels">Open vessels</v-btn>
          </v-card-title>
          <v-card-text>
            <div class="vessel-summary">
              <v-chip size="x-small" variant="tonal" color="primary">
                {{ vesselSummary.occupied }} occupied
              </v-chip>
              <v-chip size="x-small" variant="tonal" color="success">
                {{ vesselSummary.available }} available
              </v-chip>
              <v-chip size="x-small" variant="tonal" color="warning">
                {{ vesselSummary.outOfService }} out of service
              </v-chip>
            </div>

            <v-list class="vessel-list" lines="two">
              <v-list-item v-for="item in vesselCards" :key="item.vessel.id">
                <v-list-item-title class="d-flex align-center flex-wrap ga-1">
                  {{ item.vessel.name }}
                  <v-chip class="ml-2" size="x-small" variant="tonal" :color="item.occupancyTone">
                    {{ item.occupancyLabel }}
                  </v-chip>
                  <v-chip
                    v-if="item.occupancyStatus"
                    size="x-small"
                    variant="flat"
                    :color="item.occupancyStatusColor"
                  >
                    {{ item.occupancyStatusLabel }}
                  </v-chip>
                </v-list-item-title>
                <v-list-item-subtitle>
                  {{ item.vessel.type }} · {{ formatCapacity(item.vessel.capacity, item.vessel.capacity_unit) }}
                </v-list-item-subtitle>
                <div class="text-caption text-medium-emphasis">
                  {{ item.occupancyDetail }}
                </div>
                <template #append>
                  <v-chip size="x-small" variant="outlined">
                    {{ item.vesselStatusLabel }}
                  </v-chip>
                </template>
              </v-list-item>

              <v-list-item v-if="vessels.length === 0">
                <v-list-item-title>No vessels yet</v-list-item-title>
                <v-list-item-subtitle>Register a vessel to track occupancy.</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'
import { useApiClient } from '@/composables/useApiClient'

type Unit = 'ml' | 'usfloz' | 'ukfloz'

type ProcessPhase =
  | 'planning'
  | 'mashing'
  | 'heating'
  | 'boiling'
  | 'cooling'
  | 'fermenting'
  | 'conditioning'
  | 'packaging'
  | 'finished'

type Batch = {
  id: number
  uuid: string
  short_name: string
  brew_date: string | null
  updated_at: string
}

type BatchProcessPhase = {
  id: number
  batch_id: number
  process_phase: ProcessPhase
  phase_at: string
  created_at: string
}

type Vessel = {
  id: number
  uuid: string
  type: string
  name: string
  capacity: number
  capacity_unit: Unit
  status: 'active' | 'inactive' | 'retired'
  updated_at: string
}

type Volume = {
  id: number
  uuid: string
  name: string | null
  description: string | null
  amount: number
  amount_unit: Unit
  updated_at: string
}

type Occupancy = {
  id: number
  vessel_id: number
  volume_id: number
  status: string | null
  in_at: string
  out_at: string | null
  created_at: string
  updated_at: string
}

type BatchPhaseItem = {
  batch: Batch
  phase: BatchProcessPhase | null
  route: string
  phaseLabel: string
  phaseTone: string
  phaseAt: string
}

type VesselCard = {
  vessel: Vessel
  occupancy: Occupancy | null
  occupancyLabel: string
  occupancyTone: string
  occupancyDetail: string
  occupancyStatus: string | null
  occupancyStatusLabel: string
  occupancyStatusColor: string
  vesselStatusLabel: string
  sortOrder: number
}

const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

const batches = ref<Batch[]>([])
const vessels = ref<Vessel[]>([])
const volumes = ref<Volume[]>([])
const processPhases = ref<BatchProcessPhase[]>([])
const occupancies = ref<Occupancy[]>([])
const errorMessage = ref('')
const loading = ref(false)

const { request } = useApiClient(apiBase)

const volumeNameMap = computed(
  () => new Map(volumes.value.map((volume) => [volume.id, volume.name ?? `Volume ${volume.id}`])),
)

const latestPhaseByBatch = computed(() => {
  const map = new Map<number, BatchProcessPhase>()
  processPhases.value.forEach((phase) => {
    const current = map.get(phase.batch_id)
    if (!current || toTimestamp(phase.phase_at || phase.created_at) > toTimestamp(current.phase_at || current.created_at)) {
      map.set(phase.batch_id, phase)
    }
  })
  return map
})

const batchPhaseItems = computed<BatchPhaseItem[]>(() =>
  batches.value.map((batch) => {
    const phase = latestPhaseByBatch.value.get(batch.id) ?? null
    const phaseAt = phase
      ? normalizeTimestamp(phase.phase_at, phase.created_at || batch.updated_at)
      : batch.updated_at
    return {
      batch,
      phase,
      route: `/batches/${batch.uuid}`,
      phaseLabel: phase ? formatPhaseLabel(phase.process_phase) : 'No phase',
      phaseTone: phase ? phaseTone(phase.process_phase) : 'secondary',
      phaseAt,
    }
  }),
)

const inProgressAll = computed(() =>
  batchPhaseItems.value.filter((item) => isBatchInProgress(item.batch, item.phase)),
)

const plannedAll = computed(() =>
  batchPhaseItems.value.filter((item) => isBatchPlanned(item.batch, item.phase)),
)

const inProgressBatches = computed(() => [...inProgressAll.value].sort(sortInProgress).slice(0, 6))
const plannedBatches = computed(() => [...plannedAll.value].sort(sortPlanned).slice(0, 6))

const inProgressCount = computed(() => inProgressAll.value.length)
const plannedCount = computed(() => plannedAll.value.length)

const occupancyMap = computed(
  () => new Map(occupancies.value.map((occupancy) => [occupancy.vessel_id, occupancy])),
)

const vesselCards = computed<VesselCard[]>(() => {
  return [...vessels.value]
    .map((vessel) => {
      const occupancy = occupancyMap.value.get(vessel.id) ?? null
      const volumeName = occupancy
        ? volumeNameMap.value.get(occupancy.volume_id) ?? `Volume ${occupancy.volume_id}`
        : null
      const isActive = vessel.status === 'active'
      const hasOccupancy = Boolean(occupancy)

      let occupancyLabel = 'Available'
      let occupancyTone = 'success'
      let occupancyDetail = 'Ready for transfers.'
      let sortOrder = 2

      if (!isActive) {
        occupancyLabel = vessel.status === 'retired' ? 'Retired' : 'Inactive'
        occupancyTone = 'warning'
        occupancyDetail = 'Out of service.'
        sortOrder = 3
      } else if (hasOccupancy && occupancy) {
        const occupiedAt = normalizeTimestamp(occupancy.in_at, occupancy.created_at)
        const statusLabel = formatOccupancyStatusLabel(occupancy.status)
        occupancyLabel = 'Occupied'
        occupancyTone = 'primary'
        occupancyDetail = statusLabel
          ? `${volumeName} · ${statusLabel} since ${formatDateTime(occupiedAt)}`
          : `Holding ${volumeName} since ${formatDateTime(occupiedAt)}`
        sortOrder = 1
      }

      return {
        vessel,
        occupancy,
        occupancyLabel,
        occupancyTone,
        occupancyDetail,
        occupancyStatus: occupancy?.status ?? null,
        occupancyStatusLabel: formatOccupancyStatusLabel(occupancy?.status),
        occupancyStatusColor: getOccupancyStatusColor(occupancy?.status),
        vesselStatusLabel: formatVesselStatus(vessel.status),
        sortOrder,
      }
    })
    .sort((a, b) => a.sortOrder - b.sortOrder || a.vessel.name.localeCompare(b.vessel.name))
})

const vesselSummary = computed(() => {
  let occupied = 0
  let available = 0
  let outOfService = 0
  vesselCards.value.forEach((item) => {
    if (item.vessel.status !== 'active') {
      outOfService += 1
    } else if (item.occupancyLabel === 'Occupied') {
      occupied += 1
    } else {
      available += 1
    }
  })
  return { occupied, available, outOfService }
})

onMounted(async () => {
  await refreshAll()
})

async function refreshAll() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [batchData, vesselData, volumeData] = await Promise.all([
      request<Batch[]>('/batches'),
      request<Vessel[]>('/vessels'),
      request<Volume[]>('/volumes'),
    ])
    batches.value = batchData
    vessels.value = vesselData
    volumes.value = volumeData

    await Promise.all([loadProcessPhases(), loadOccupancies()])
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load dashboard'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

async function loadProcessPhases() {
  if (batches.value.length === 0) {
    processPhases.value = []
    return
  }
  const results = await Promise.allSettled(
    batches.value.map((batch) => request<BatchProcessPhase[]>(`/batch-process-phases?batch_id=${batch.id}`)),
  )
  const phases: BatchProcessPhase[] = []
  const errors = results.filter((result) => result.status === 'rejected') as PromiseRejectedResult[]
  results.forEach((result) => {
    if (result.status === 'fulfilled') {
      phases.push(...result.value)
    }
  })
  if (errors.length > 0) {
    throw errors[0]!.reason
  }
  processPhases.value = phases
}

async function loadOccupancies() {
  if (vessels.value.length === 0) {
    occupancies.value = []
    return
  }
  occupancies.value = await request<Occupancy[]>('/occupancies?active=true')
}

function isBatchInProgress(batch: Batch, phase: BatchProcessPhase | null) {
  if (phase) {
    return phase.process_phase !== 'planning' && phase.process_phase !== 'finished'
  }
  if (!batch.brew_date) {
    return false
  }
  return !isFutureDate(batch.brew_date)
}

function isBatchPlanned(batch: Batch, phase: BatchProcessPhase | null) {
  if (phase) {
    return phase.process_phase === 'planning'
  }
  if (!batch.brew_date) {
    return true
  }
  return isFutureDate(batch.brew_date)
}

function sortInProgress(a: BatchPhaseItem, b: BatchPhaseItem) {
  return toTimestamp(b.phaseAt) - toTimestamp(a.phaseAt)
}

function sortPlanned(a: BatchPhaseItem, b: BatchPhaseItem) {
  const dateA = toDayTimestamp(a.batch.brew_date)
  const dateB = toDayTimestamp(b.batch.brew_date)
  if (dateA !== null && dateB !== null) {
    return dateA - dateB
  }
  if (dateA !== null) {
    return -1
  }
  if (dateB !== null) {
    return 1
  }
  return toTimestamp(b.batch.updated_at) - toTimestamp(a.batch.updated_at)
}

function plannedSubtitle(item: BatchPhaseItem) {
  if (item.phase?.process_phase === 'planning') {
    return 'In planning'
  }
  if (!item.batch.brew_date) {
    return 'Brew date TBD'
  }
  return `Brew date ${formatDate(item.batch.brew_date)}`
}

function plannedMeta(item: BatchPhaseItem) {
  const updated = item.phaseAt || item.batch.updated_at
  return `Updated ${formatDateTime(updated)}`
}

function plannedBadge(item: BatchPhaseItem) {
  if (item.phase?.process_phase === 'planning') {
    return 'Planning'
  }
  return formatRelativeDate(item.batch.brew_date)
}

function formatPhaseLabel(phase: ProcessPhase) {
  return phase
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

function phaseTone(phase: ProcessPhase) {
  const tones: Record<ProcessPhase, string> = {
    planning: 'info',
    mashing: 'warning',
    heating: 'warning',
    boiling: 'warning',
    cooling: 'warning',
    fermenting: 'secondary',
    conditioning: 'secondary',
    packaging: 'primary',
    finished: 'success',
  }
  return tones[phase] ?? 'primary'
}

function formatVesselStatus(status: Vessel['status']) {
  if (status === 'inactive') {
    return 'Inactive'
  }
  if (status === 'retired') {
    return 'Retired'
  }
  return 'Active'
}

function formatOccupancyStatusLabel(status: string | null | undefined): string {
  if (!status) {
    return ''
  }
  const statusLabels: Record<string, string> = {
    fermenting: 'Fermenting',
    conditioning: 'Conditioning',
    cold_crashing: 'Cold Crashing',
    dry_hopping: 'Dry Hopping',
    carbonating: 'Carbonating',
    holding: 'Holding',
    packaging: 'Packaging',
  }
  return statusLabels[status] ?? status.charAt(0).toUpperCase() + status.slice(1).replace(/_/g, ' ')
}

function getOccupancyStatusColor(status: string | null | undefined): string {
  if (!status) {
    return 'grey'
  }
  const statusColors: Record<string, string> = {
    fermenting: 'orange',
    conditioning: 'blue',
    cold_crashing: 'cyan',
    dry_hopping: 'green',
    carbonating: 'purple',
    holding: 'grey',
    packaging: 'teal',
  }
  return statusColors[status] ?? 'secondary'
}

function formatCapacity(amount: number, unit: Unit) {
  return `${amount} ${unit}`
}

function formatDate(value: string | null | undefined) {
  if (!value) {
    return 'Unknown'
  }
  return new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
  }).format(new Date(value))
}

function formatBrewDate(value: string | null | undefined) {
  if (!value) {
    return 'Not set'
  }
  return formatDate(value)
}

function formatDateTime(value: string | null | undefined) {
  if (!value) {
    return 'Unknown'
  }
  return new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value))
}

function normalizeTimestamp(value: string | null | undefined, fallback: string) {
  if (!value) {
    return fallback
  }
  const parsed = new Date(value)
  if (Number.isNaN(parsed.getTime()) || parsed.getFullYear() < 2000) {
    return fallback
  }
  return value
}

function toTimestamp(value: string | null | undefined) {
  if (!value) {
    return 0
  }
  const parsed = new Date(value).getTime()
  return Number.isFinite(parsed) ? parsed : 0
}

function toDayTimestamp(value: string | null | undefined) {
  if (!value) {
    return null
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return null
  }
  date.setHours(0, 0, 0, 0)
  return date.getTime()
}

function isFutureDate(value: string) {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const target = new Date(value)
  if (Number.isNaN(target.getTime())) {
    return false
  }
  target.setHours(0, 0, 0, 0)
  return target.getTime() > today.getTime()
}

function formatRelativeDate(value: string | null | undefined) {
  if (!value) {
    return 'Unscheduled'
  }
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const target = new Date(value)
  if (Number.isNaN(target.getTime())) {
    return 'Unscheduled'
  }
  target.setHours(0, 0, 0, 0)
  const diff = Math.round((target.getTime() - today.getTime()) / (1000 * 60 * 60 * 24))
  if (diff === 0) {
    return 'Today'
  }
  if (diff === 1) {
    return 'Tomorrow'
  }
  if (diff > 1) {
    return `In ${diff} days`
  }
  return `${Math.abs(diff)} days ago`
}
</script>

<style scoped>
.dashboard-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.hero-card {
  background: linear-gradient(
    135deg,
    rgba(var(--v-theme-primary), 0.14),
    rgba(var(--v-theme-secondary), 0.08)
  );
}

.dashboard-hero {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  align-items: center;
  justify-content: space-between;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.batch-list,
.vessel-list {
  max-height: 360px;
  overflow: auto;
}

.vessel-summary {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}
</style>
