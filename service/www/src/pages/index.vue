<template>
  <v-container class="dashboard-page" fluid>
    <v-card class="section-card hero-card">
      <v-card-text>
        <div class="dashboard-hero">
          <div>
            <div class="text-overline text-medium-emphasis">Dashboard</div>
            <div class="text-h4 font-weight-semibold">{{ breweryName }}</div>
            <div class="text-body-2 text-medium-emphasis">
              Track in-progress batches, upcoming brew days, and vessel readiness.
            </div>
          </div>
          <div class="hero-actions">
            <v-btn :loading="loading" variant="text" @click="refreshAll">Refresh</v-btn>
          </div>
        </div>

        <v-alert
          v-if="errorMessage"
          class="mt-4"
          density="compact"
          type="error"
          variant="tonal"
        >
          {{ errorMessage }}
        </v-alert>
      </v-card-text>
    </v-card>

    <v-row align="stretch" class="mt-4">
      <v-col cols="12" lg="7">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-progress-clock" />
            In-progress batches
            <v-chip class="ml-3" color="primary" size="x-small" variant="tonal">
              {{ inProgressCount }}
            </v-chip>
            <v-spacer />
            <v-btn size="small" to="/batches" variant="text">View all</v-btn>
          </v-card-title>
          <v-card-text>
            <v-list class="batch-list" lines="two">
              <v-list-item v-for="item in inProgressBatches" :key="item.batch.uuid" :to="item.route">
                <v-list-item-title>{{ item.batch.short_name }}</v-list-item-title>
                <v-list-item-subtitle>
                  Brewed {{ formatBrewDate(item.batch.brew_date) }}
                </v-list-item-subtitle>
                <div class="text-caption text-medium-emphasis">
                  {{ item.phaseLabel }} · Updated {{ formatDateTime(item.phaseAt) }}
                </div>
                <template #append>
                  <div class="d-flex flex-column align-end ga-1">
                    <v-chip :color="item.phaseTone" size="x-small" variant="tonal">
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
            <v-chip class="ml-3" color="info" size="x-small" variant="tonal">
              {{ plannedCount }}
            </v-chip>
            <v-spacer />
            <v-btn size="small" to="/batches" variant="text">Plan schedule</v-btn>
          </v-card-title>
          <v-card-text>
            <v-list class="batch-list" lines="two">
              <v-list-item v-for="item in plannedBatches" :key="item.batch.uuid" :to="item.route">
                <v-list-item-title>{{ item.batch.short_name }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ plannedSubtitle(item) }}
                </v-list-item-subtitle>
                <div class="text-caption text-medium-emphasis">
                  {{ plannedMeta(item) }}
                </div>
                <template #append>
                  <div class="d-flex flex-column align-end ga-1">
                    <v-chip color="info" size="x-small" variant="tonal">
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
        <v-card class="section-card" :class="{ 'low-stock-alert': hasLowStock }">
          <v-card-title class="d-flex align-center">
            <v-icon
              class="mr-2"
              :color="hasLowStock ? 'warning' : 'success'"
              :icon="hasLowStock ? 'mdi-alert' : 'mdi-check-circle'"
            />
            {{ hasLowStock ? 'Low Stock' : 'Stock Status' }}
            <v-spacer />
            <v-btn size="small" to="/inventory/stock-levels" variant="text">View Stock Levels</v-btn>
          </v-card-title>
          <v-card-text>
            <template v-if="hasLowStock">
              <div class="text-body-2 mb-3">
                {{ lowStockCount }} {{ lowStockCount === 1 ? 'ingredient' : 'ingredients' }} out of stock
              </div>
              <v-list class="low-stock-list" density="compact" lines="one">
                <v-list-item
                  v-for="item in lowStockItems"
                  :key="item.ingredient_uuid"
                  :to="`/inventory/stock-levels?ingredient=${item.ingredient_uuid}`"
                >
                  <template #prepend>
                    <v-icon color="warning" icon="mdi-circle-small" size="small" />
                  </template>
                  <v-list-item-title>{{ item.ingredient_name }}</v-list-item-title>
                  <template #append>
                    <v-chip color="secondary" size="x-small" variant="tonal">
                      {{ item.category }}
                    </v-chip>
                  </template>
                </v-list-item>
              </v-list>
            </template>
            <template v-else>
              <div class="d-flex align-center text-body-2 text-success">
                <v-icon class="mr-2" icon="mdi-check" size="small" />
                All ingredients in stock
              </div>
            </template>
          </v-card-text>
        </v-card>

        <v-card class="section-card mt-4">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-silo" />
            Vessels
            <v-spacer />
            <v-btn size="small" to="/vessels" variant="text">Open vessels</v-btn>
          </v-card-title>
          <v-card-text>
            <div class="vessel-summary">
              <v-chip color="primary" size="x-small" variant="tonal">
                {{ vesselSummary.occupied }} occupied
              </v-chip>
              <v-chip color="success" size="x-small" variant="tonal">
                {{ vesselSummary.available }} available
              </v-chip>
              <v-chip color="warning" size="x-small" variant="tonal">
                {{ vesselSummary.outOfService }} out of service
              </v-chip>
            </div>

            <v-list class="vessel-list" lines="two">
              <v-list-item v-for="item in vesselCards" :key="item.vessel.uuid">
                <v-list-item-title class="d-flex align-center flex-wrap ga-1">
                  {{ item.vessel.name }}
                  <v-chip class="ml-2" :color="item.occupancyTone" size="x-small" variant="tonal">
                    {{ item.occupancyLabel }}
                  </v-chip>
                  <v-chip
                    v-if="item.occupancyStatus"
                    :color="item.occupancyStatusColor"
                    size="x-small"
                    variant="flat"
                  >
                    {{ item.occupancyStatusLabel }}
                  </v-chip>
                </v-list-item-title>
                <v-list-item-subtitle>
                  {{ item.vessel.type }} · {{ formatVolumePreferred(item.vessel.capacity, item.vessel.capacity_unit) }}
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
  import type { BatchProcessPhase, ProcessPhase } from '@/components/batch/types'
  import type { Batch, Occupancy, StockLevel, Vessel, Volume } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import {
    formatDate,
    formatDateTime,
    useOccupancyStatusFormatters,
    usePhaseFormatters,
    useVesselStatusFormatters,
  } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { useUserSettings } from '@/composables/useUserSettings'

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

  const batches = ref<Batch[]>([])
  const vessels = ref<Vessel[]>([])
  const volumes = ref<Volume[]>([])
  const processPhases = ref<BatchProcessPhase[]>([])
  const occupancies = ref<Occupancy[]>([])
  const stockLevels = ref<StockLevel[]>([])
  const errorMessage = ref('')
  const loading = ref(false)

  const { getBatches, getVessels, getVolumes, getActiveOccupancies, request } = useProductionApi()
  const { getStockLevels } = useInventoryApi()
  const { formatVolumePreferred } = useUnitPreferences()
  const { breweryName } = useUserSettings()
  const { formatPhase, getPhaseColor } = usePhaseFormatters()
  const { formatVesselStatus } = useVesselStatusFormatters()
  const { formatOccupancyStatus: formatOccupancyStatusLabel, getOccupancyStatusColor } = useOccupancyStatusFormatters()

  const volumeNameMap = computed(
    () => new Map(volumes.value.map(volume => [volume.uuid, volume.name ?? `Volume ${volume.uuid.slice(0, 8)}`])),
  )

  const latestPhaseByBatch = computed(() => {
    const map = new Map<string, BatchProcessPhase>()
    for (const phase of processPhases.value) {
      const current = map.get(phase.batch_uuid)
      if (!current || toTimestamp(phase.phase_at || phase.created_at) > toTimestamp(current.phase_at || current.created_at)) {
        map.set(phase.batch_uuid, phase)
      }
    }
    return map
  })

  const batchPhaseItems = computed<BatchPhaseItem[]>(() =>
    batches.value.map(batch => {
      const phase = latestPhaseByBatch.value.get(batch.uuid) ?? null
      const phaseAt = phase
        ? normalizeTimestamp(phase.phase_at, phase.created_at || batch.updated_at)
        : batch.updated_at
      return {
        batch,
        phase,
        route: `/batches/${batch.uuid}`,
        phaseLabel: phase ? formatPhase(phase.process_phase) : 'No phase',
        phaseTone: phase ? getPhaseColor(phase.process_phase) : 'secondary',
        phaseAt,
      }
    }),
  )

  const inProgressAll = computed(() =>
    batchPhaseItems.value.filter(item => isBatchInProgress(item.batch, item.phase)),
  )

  const plannedAll = computed(() =>
    batchPhaseItems.value.filter(item => isBatchPlanned(item.batch, item.phase)),
  )

  const inProgressBatches = computed(() => [...inProgressAll.value].sort(sortInProgress).slice(0, 6))
  const plannedBatches = computed(() => [...plannedAll.value].sort(sortPlanned).slice(0, 6))

  const inProgressCount = computed(() => inProgressAll.value.length)
  const plannedCount = computed(() => plannedAll.value.length)

  const occupancyMap = computed(
    () => new Map(occupancies.value.map(occupancy => [occupancy.vessel_uuid, occupancy])),
  )

  const vesselCards = computed<VesselCard[]>(() => {
    return [...vessels.value]
      .map(vessel => {
        const occupancy = occupancyMap.value.get(vessel.uuid) ?? null
        const volumeName = occupancy
          ? volumeNameMap.value.get(occupancy.volume_uuid) ?? `Volume ${occupancy.volume_uuid.slice(0, 8)}`
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
    for (const item of vesselCards.value) {
      if (item.vessel.status !== 'active') {
        outOfService += 1
      } else if (item.occupancyLabel === 'Occupied') {
        occupied += 1
      } else {
        available += 1
      }
    }
    return { occupied, available, outOfService }
  })

  const lowStockItems = computed(() =>
    stockLevels.value.filter(item => item.total_on_hand <= 0),
  )

  const lowStockCount = computed(() => lowStockItems.value.length)

  const hasLowStock = computed(() => lowStockCount.value > 0)

  onMounted(async () => {
    await refreshAll()
  })

  async function refreshAll () {
    loading.value = true
    errorMessage.value = ''
    try {
      const [batchData, vesselData, volumeData, stockData] = await Promise.all([
        getBatches(),
        getVessels(),
        getVolumes(),
        getStockLevels(),
      ])
      batches.value = batchData
      vessels.value = vesselData
      volumes.value = volumeData
      stockLevels.value = stockData

      await Promise.all([loadProcessPhases(), loadOccupancies()])
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load dashboard'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }

  async function loadProcessPhases () {
    if (batches.value.length === 0) {
      processPhases.value = []
      return
    }
    const results = await Promise.allSettled(
      batches.value.map(batch => request<BatchProcessPhase[]>(`/batch-process-phases?batch_uuid=${batch.uuid}`)),
    )
    const phases: BatchProcessPhase[] = []
    const errors = results.filter(result => result.status === 'rejected') as PromiseRejectedResult[]
    for (const result of results) {
      if (result.status === 'fulfilled') {
        phases.push(...result.value)
      }
    }
    if (errors.length > 0) {
      throw errors[0]!.reason
    }
    processPhases.value = phases
  }

  async function loadOccupancies () {
    if (vessels.value.length === 0) {
      occupancies.value = []
      return
    }
    occupancies.value = await getActiveOccupancies()
  }

  function isBatchInProgress (batch: Batch, phase: BatchProcessPhase | null) {
    if (phase) {
      return phase.process_phase !== 'planning' && phase.process_phase !== 'finished'
    }
    if (!batch.brew_date) {
      return false
    }
    return !isFutureDate(batch.brew_date)
  }

  function isBatchPlanned (batch: Batch, phase: BatchProcessPhase | null) {
    if (phase) {
      return phase.process_phase === 'planning'
    }
    if (!batch.brew_date) {
      return true
    }
    return isFutureDate(batch.brew_date)
  }

  function sortInProgress (a: BatchPhaseItem, b: BatchPhaseItem) {
    return toTimestamp(b.phaseAt) - toTimestamp(a.phaseAt)
  }

  function sortPlanned (a: BatchPhaseItem, b: BatchPhaseItem) {
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

  function plannedSubtitle (item: BatchPhaseItem) {
    if (item.phase?.process_phase === 'planning') {
      return 'In planning'
    }
    if (!item.batch.brew_date) {
      return 'Brew date TBD'
    }
    return `Brew date ${formatDate(item.batch.brew_date)}`
  }

  function plannedMeta (item: BatchPhaseItem) {
    const updated = item.phaseAt || item.batch.updated_at
    return `Updated ${formatDateTime(updated)}`
  }

  function plannedBadge (item: BatchPhaseItem) {
    if (item.phase?.process_phase === 'planning') {
      return 'Planning'
    }
    return formatRelativeDate(item.batch.brew_date)
  }

  function formatBrewDate (value: string | null | undefined) {
    if (!value) {
      return 'Not set'
    }
    return formatDate(value)
  }

  function normalizeTimestamp (value: string | null | undefined, fallback: string) {
    if (!value) {
      return fallback
    }
    const parsed = new Date(value)
    if (Number.isNaN(parsed.getTime()) || parsed.getFullYear() < 2000) {
      return fallback
    }
    return value
  }

  function toTimestamp (value: string | null | undefined) {
    if (!value) {
      return 0
    }
    const parsed = new Date(value).getTime()
    return Number.isFinite(parsed) ? parsed : 0
  }

  function toDayTimestamp (value: string | null | undefined) {
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

  function isFutureDate (value: string) {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const target = new Date(value)
    if (Number.isNaN(target.getTime())) {
      return false
    }
    target.setHours(0, 0, 0, 0)
    return target.getTime() > today.getTime()
  }

  function formatRelativeDate (value: string | null | undefined) {
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

.low-stock-alert {
  border-left: 4px solid rgb(var(--v-theme-warning));
}

.low-stock-list {
  max-height: 200px;
  overflow: auto;
}
</style>
