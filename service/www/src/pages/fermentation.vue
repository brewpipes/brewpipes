<template>
  <v-container class="fermentation-page" fluid>
    <!-- Page header -->
    <v-card class="section-card hero-card">
      <v-card-text>
        <div class="page-hero">
          <div>
            <div class="text-overline text-medium-emphasis">Production</div>
            <div class="text-h4 font-weight-semibold">Fermentation</div>
            <div class="text-body-2 text-medium-emphasis">
              {{ subtitleText }}
            </div>
          </div>
          <div class="hero-actions">
            <v-btn
              :loading="loading"
              prepend-icon="mdi-refresh"
              variant="text"
              @click="refreshAll"
            >
              <span class="d-none d-sm-inline">Refresh</span>
            </v-btn>
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

    <!-- Loading skeleton -->
    <v-row v-if="!dataReady && loading" class="mt-4">
      <v-col v-for="n in 3" :key="n" cols="12" lg="4" sm="6">
        <v-skeleton-loader type="card" />
      </v-col>
    </v-row>

    <!-- Empty state -->
    <v-card
      v-if="dataReady && sortedCards.length === 0"
      class="section-card mt-4 text-center pa-8"
    >
      <v-icon
        class="mb-4"
        color="medium-emphasis"
        icon="mdi-flask-round-bottom-empty"
        size="64"
      />
      <div class="text-h6 text-medium-emphasis mb-2">No active fermentations</div>
      <div class="text-body-2 text-medium-emphasis mb-4">
        Assign a batch to a fermenter to start tracking.
      </div>
      <v-btn color="primary" to="/batches" variant="tonal">
        View Batches
      </v-btn>
    </v-card>

    <!-- Card grid -->
    <v-row v-if="dataReady && sortedCards.length > 0" class="mt-4">
      <v-col
        v-for="card in sortedCards"
        :key="card.occupancy.uuid"
        cols="12"
        lg="4"
        sm="6"
      >
        <FermentationCard
          :batch-summary="card.batchSummary"
          :measurements="card.measurements"
          :occupancy="card.occupancy"
          :vessel="card.vessel"
          @blend="openBlend(card)"
          @log-reading="openReadingSheet(card)"
          @mark-empty="openMarkEmpty(card)"
          @split="openSplit(card)"
          @status-changed="refreshAll"
          @transfer="openTransfer(card)"
        />
      </v-col>
    </v-row>

    <!-- QuickReadingSheet (single instance) -->
    <QuickReadingSheet
      v-model="showReadingSheet"
      :batch-name="selectedBatchName"
      :batch-uuid="selectedBatchUuid"
      :occupancy-uuid="selectedOccupancyUuid"
      :vessel-name="selectedVesselName"
      @saved="handleReadingSaved"
    />

    <!-- BatchMarkEmptyDialog (single instance) -->
    <BatchMarkEmptyDialog
      v-model="showMarkEmpty"
      :batch-name="selectedBatchName"
      :occupancy="selectedOccupancy"
      :vessel-name="selectedVesselName"
      @emptied="refreshAll"
    />

    <!-- TransferDialog (single instance, supports transfer/split/blend modes) -->
    <TransferDialog
      v-model="showTransfer"
      :mode="transferMode"
      :source-batch="transferSourceBatch"
      :source-occupancy="transferSourceOccupancy"
      :source-vessel="transferSourceVessel"
      :source-volume="transferSourceVolume"
      @transferred="refreshAll"
    />
  </v-container>
</template>

<script lang="ts" setup>
  import type { Batch, BatchSummary, Measurement, Occupancy, Vessel, Volume } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { BatchMarkEmptyDialog } from '@/components/batch'
  import { FermentationCard, QuickReadingSheet, TransferDialog } from '@/components/fermentation'
  import { useProductionApi } from '@/composables/useProductionApi'

  /** Data for a single fermentation card */
  interface FermentationCardData {
    occupancy: Occupancy
    vessel: Vessel
    batchSummary: BatchSummary | null
    measurements: Measurement[]
    attentionLevel: 'warning' | 'info' | null
    daysInTank: number
  }

  const {
    getActiveOccupancies,
    getVessels,
    getBatch,
    getVolume,
    getBatchSummary,
    getMeasurementsByBatch,
  } = useProductionApi()

  const loading = ref(false)
  const dataReady = ref(false)
  const errorMessage = ref('')

  // Raw data
  const occupancies = ref<Occupancy[]>([])
  const vessels = ref<Vessel[]>([])
  const batchSummaries = ref<Map<string, BatchSummary>>(new Map())
  const batchMeasurements = ref<Map<string, Measurement[]>>(new Map())

  // QuickReadingSheet state
  const showReadingSheet = ref(false)
  const selectedBatchUuid = ref('')
  const selectedOccupancyUuid = ref<string | undefined>(undefined)
  const selectedVesselName = ref('')
  const selectedBatchName = ref('')

  // BatchMarkEmptyDialog state
  const showMarkEmpty = ref(false)
  const selectedOccupancy = ref<Occupancy | null>(null)

  // TransferDialog state
  type TransferMode = 'transfer' | 'split' | 'blend'
  const showTransfer = ref(false)
  const transferMode = ref<TransferMode>('transfer')
  const transferSourceOccupancy = ref<Occupancy | null>(null)
  const transferSourceVessel = ref<Vessel | null>(null)
  const transferSourceBatch = ref<Batch | null>(null)
  const transferSourceVolume = ref<Volume | null>(null)

  // Vessel lookup map
  const vesselMap = computed(
    () => new Map(vessels.value.map(v => [v.uuid, v])),
  )

  // Build card data
  const cardData = computed<FermentationCardData[]>(() => {
    return occupancies.value
      .filter(occ => occ.batch_uuid)
      .map(occ => {
        const vessel = vesselMap.value.get(occ.vessel_uuid)
        if (!vessel) return null

        const batchUuid = occ.batch_uuid!
        const summary = batchSummaries.value.get(batchUuid) ?? null
        const measurements = batchMeasurements.value.get(batchUuid) ?? []

        // Compute attention level
        const gravityReadings = measurements
          .filter(m => m.kind === 'gravity')
          .sort((a, b) => new Date(a.observed_at).getTime() - new Date(b.observed_at).getTime())

        let attentionLevel: 'warning' | 'info' | null = null

        // Stale: no gravity reading in 24+ hours
        if (gravityReadings.length === 0) {
          attentionLevel = 'warning'
        } else {
          const lastReading = gravityReadings.at(-1)!
          const hoursSince = (Date.now() - new Date(lastReading.observed_at).getTime()) / (1000 * 60 * 60)
          if (hoursSince >= 24) {
            attentionLevel = 'warning'
          } else if (gravityReadings.length >= 3) {
            // Stable: last 3 readings unchanged
            const recent = gravityReadings.slice(-3)
            const first = recent[0]!.value
            const isStable = recent.every(r => Math.abs(r.value - first) < 0.001)
            if (isStable) {
              attentionLevel = 'info'
            }
          }
        }

        // Days in tank
        const inAt = new Date(occ.in_at)
        const daysInTank = Math.max(0, Math.floor((Date.now() - inAt.getTime()) / (1000 * 60 * 60 * 24)))

        return {
          occupancy: occ,
          vessel,
          batchSummary: summary,
          measurements,
          attentionLevel,
          daysInTank,
        }
      })
      .filter((item): item is FermentationCardData => item !== null)
  })

  // Sorted: attention items first (warning > info > null), then by days descending
  const sortedCards = computed(() => {
    const attentionOrder = (level: 'warning' | 'info' | null): number => {
      if (level === 'warning') return 0
      if (level === 'info') return 1
      return 2
    }

    return [...cardData.value].sort((a, b) => {
      const aOrder = attentionOrder(a.attentionLevel)
      const bOrder = attentionOrder(b.attentionLevel)
      if (aOrder !== bOrder) return aOrder - bOrder
      return b.daysInTank - a.daysInTank
    })
  })

  // Subtitle
  const activeCount = computed(() => sortedCards.value.length)
  const attentionCount = computed(() =>
    sortedCards.value.filter(c => c.attentionLevel !== null).length,
  )

  const subtitleText = computed(() => {
    const parts: string[] = []
    parts.push(`${activeCount.value} active`)
    if (attentionCount.value > 0) {
      parts.push(`${attentionCount.value} need attention`)
    }
    return parts.join(' · ')
  })

  // Data loading
  onMounted(async () => {
    await refreshAll()
  })

  async function refreshAll () {
    loading.value = true
    errorMessage.value = ''

    try {
      const [occupancyData, vesselData] = await Promise.all([
        getActiveOccupancies(),
        getVessels(),
      ])

      occupancies.value = occupancyData
      vessels.value = vesselData

      // Get unique batch UUIDs from occupancies
      const batchUuids = [...new Set(
        occupancyData
          .map(occ => occ.batch_uuid)
          .filter((uuid): uuid is string => uuid !== null),
      )]

      // Load batch summaries and measurements in parallel
      if (batchUuids.length > 0) {
        const [summaryResults, measurementResults] = await Promise.all([
          Promise.allSettled(batchUuids.map(uuid => getBatchSummary(uuid))),
          Promise.allSettled(batchUuids.map(uuid => getMeasurementsByBatch(uuid))),
        ])

        const newSummaries = new Map<string, BatchSummary>()
        summaryResults.forEach((result, index) => {
          if (result.status === 'fulfilled') {
            newSummaries.set(batchUuids[index]!, result.value)
          }
        })
        batchSummaries.value = newSummaries

        const newMeasurements = new Map<string, Measurement[]>()
        measurementResults.forEach((result, index) => {
          if (result.status === 'fulfilled') {
            newMeasurements.set(batchUuids[index]!, result.value)
          }
        })
        batchMeasurements.value = newMeasurements
      } else {
        batchSummaries.value = new Map()
        batchMeasurements.value = new Map()
      }

      dataReady.value = true
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load fermentation data'
      errorMessage.value = message
    } finally {
      loading.value = false
    }
  }

  // Refresh measurements for a specific batch
  async function refreshBatchMeasurements (batchUuid: string) {
    try {
      const measurements = await getMeasurementsByBatch(batchUuid)
      const updated = new Map(batchMeasurements.value)
      updated.set(batchUuid, measurements)
      batchMeasurements.value = updated
    } catch {
      // Silently fail — card will show stale data
    }
  }

  // QuickReadingSheet handlers
  function openReadingSheet (card: FermentationCardData) {
    selectedBatchUuid.value = card.occupancy.batch_uuid ?? ''
    selectedOccupancyUuid.value = card.occupancy.uuid
    selectedVesselName.value = card.vessel.name
    selectedBatchName.value = card.batchSummary?.short_name ?? ''
    showReadingSheet.value = true
  }

  function handleReadingSaved () {
    if (selectedBatchUuid.value) {
      refreshBatchMeasurements(selectedBatchUuid.value)
    }
  }

  // BatchMarkEmptyDialog handlers
  function openMarkEmpty (card: FermentationCardData) {
    selectedOccupancy.value = card.occupancy
    selectedVesselName.value = card.vessel.name
    selectedBatchName.value = card.batchSummary?.short_name ?? `Batch ${card.occupancy.batch_uuid?.slice(0, 8) ?? '—'}`
    selectedBatchUuid.value = card.occupancy.batch_uuid ?? ''
    showMarkEmpty.value = true
  }

  // TransferDialog handlers
  async function openTransferWithMode (card: FermentationCardData, mode: TransferMode) {
    transferMode.value = mode
    transferSourceOccupancy.value = card.occupancy
    transferSourceVessel.value = card.vessel
    transferSourceBatch.value = null
    transferSourceVolume.value = null

    // Resolve batch and volume for the transfer dialog (non-critical — dialog opens regardless)
    const [batchData, volumeData] = await Promise.allSettled([
      card.occupancy.batch_uuid ? getBatch(card.occupancy.batch_uuid) : Promise.reject(new Error('no batch')),
      card.occupancy.volume_uuid ? getVolume(card.occupancy.volume_uuid) : Promise.reject(new Error('no volume')),
    ])
    if (batchData.status === 'fulfilled') {
      transferSourceBatch.value = batchData.value
    }
    if (volumeData.status === 'fulfilled') {
      transferSourceVolume.value = volumeData.value
    }

    showTransfer.value = true
  }

  function openTransfer (card: FermentationCardData) {
    openTransferWithMode(card, 'transfer')
  }

  function openSplit (card: FermentationCardData) {
    openTransferWithMode(card, 'split')
  }

  function openBlend (card: FermentationCardData) {
    openTransferWithMode(card, 'blend')
  }
</script>

<style scoped>
.fermentation-page {
  position: relative;
}

.hero-card {
  background: linear-gradient(
    135deg,
    rgba(var(--v-theme-primary), 0.14),
    rgba(var(--v-theme-secondary), 0.08)
  );
}

.page-hero {
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
</style>
