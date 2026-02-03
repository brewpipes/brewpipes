<template>
  <v-card class="section-card">
    <v-card-title class="d-flex align-center">
      <v-btn
        v-if="showBackButton"
        aria-label="Back"
        class="mr-2"
        icon="mdi-arrow-left"
        size="small"
        variant="text"
        @click="handleBack"
      />
      <v-icon class="mr-2" icon="mdi-beaker-outline" />
      {{ selectedBatch ? selectedBatch.short_name : 'Batch details' }}
      <v-spacer />
      <v-btn size="small" variant="text" @click="refresh">Refresh</v-btn>
      <v-btn v-if="selectedBatch" aria-label="Edit batch" icon="mdi-pencil" size="small" variant="text" @click="openEditDialog" />
      <v-btn v-if="selectedBatch" aria-label="Delete batch" color="error" icon="mdi-delete" size="small" variant="text" @click="openDeleteDialog" />
      <v-btn v-if="!showBackButton" size="small" variant="text" @click="clearSelection">Clear</v-btn>
    </v-card-title>
    <v-card-text>
      <v-alert
        v-if="!batchId"
        density="comfortable"
        type="info"
        variant="tonal"
      >
        Select a batch to review timeline, flow, measurements, and additions.
      </v-alert>

      <v-alert
        v-else-if="loading"
        density="comfortable"
        type="info"
        variant="tonal"
      >
        Loading batch details...
      </v-alert>

      <v-alert
        v-else-if="!selectedBatch"
        density="comfortable"
        type="warning"
        variant="tonal"
      >
        Batch not found.
      </v-alert>

      <div v-else>
        <v-row align="stretch" class="mb-4">
          <v-col cols="12">
            <v-card class="mini-card" variant="tonal">
              <v-card-text>
                <div class="text-overline">Details</div>
                <v-row class="mt-1" dense>
                  <v-col cols="12" md="6">
                    <div class="text-caption text-medium-emphasis">Brew date</div>
                    <div class="text-body-2 font-weight-medium">
                      {{ formatDate(selectedBatch.brew_date) }}
                    </div>
                  </v-col>
                  <v-col cols="12" md="6">
                    <div class="text-caption text-medium-emphasis">Current status</div>
                    <div class="d-flex flex-wrap ga-2">
                      <v-chip v-if="latestProcessPhase" color="primary" size="small" variant="tonal">
                        {{ latestProcessPhase.process_phase }}
                      </v-chip>
                      <v-chip v-if="latestLiquidPhase" color="secondary" size="small" variant="tonal">
                        {{ latestLiquidPhase.liquid_phase }}
                      </v-chip>
                      <v-chip v-if="!latestLiquidPhase && !latestProcessPhase" size="small" variant="outlined">
                        No status yet
                      </v-chip>
                    </div>
                  </v-col>
                </v-row>
                <v-divider class="my-3" />
                <div class="text-caption text-medium-emphasis mb-2">Latest measurements</div>
                <v-row dense>
                  <v-col cols="12" md="4">
                    <BatchSparklineCard
                      color="info"
                      label="Temp"
                      :latest-label="temperatureSeries.latestLabel"
                      :values="temperatureSeries.values"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <BatchSparklineCard
                      color="secondary"
                      label="Gravity"
                      :latest-label="gravitySeries.latestLabel"
                      :values="gravitySeries.values"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <BatchSparklineCard
                      color="warning"
                      label="pH"
                      :latest-label="phSeries.latestLabel"
                      :values="phSeries.values"
                    />
                  </v-col>
                </v-row>
                <div class="text-body-2 text-medium-emphasis">
                  Last updated {{ formatDateTime(selectedBatch.updated_at) }}
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <v-tabs v-model="activeTab" class="batch-tabs" color="primary" show-arrows>
          <v-tab value="summary">Summary</v-tab>
          <v-tab value="brew-sessions">Brew Sessions</v-tab>
          <v-tab value="timeline">Timeline</v-tab>
          <v-tab value="flow">Flow</v-tab>
          <v-tab value="measurements">Measurements</v-tab>
          <v-tab value="additions">Additions</v-tab>
        </v-tabs>

        <v-window v-model="activeTab" class="mt-4">
          <v-window-item value="summary">
            <BatchSummaryTab
              :loading="batchSummaryLoading"
              :summary="batchSummary"
              @occupancy-status-change="changeOccupancyStatus"
            />
          </v-window-item>

          <v-window-item value="brew-sessions">
            <BatchBrewSessionsTab
              :additions="wortAdditions"
              :measurements="wortMeasurements"
              :selected-session-id="selectedBrewSessionId"
              :sessions="brewSessions"
              :vessels="vessels"
              :volumes="allVolumes"
              @clear-session="clearBrewSessionSelection"
              @create-addition="openCreateHotSideAdditionDialog"
              @create-measurement="openCreateHotSideMeasurementDialog"
              @create-session="openCreateBrewSessionDialog"
              @edit-session="openEditBrewSessionDialog"
              @select-session="selectBrewSession"
            />
          </v-window-item>

          <v-window-item value="timeline">
            <BatchTimelineTab
              :events="timelineItems"
              :gravity-unit="preferences.gravity"
              :reading="timelineReading"
              :reading-ready="timelineReadingReady"
              :temperature-unit="preferences.temperature"
              @open-extended="openTimelineExtendedDialog"
              @record="recordTimelineReading"
              @update:reading="updateTimelineReading"
            />
          </v-window-item>

          <v-window-item value="flow">
            <BatchFlowTab
              :links="flowLinks"
              :nodes="flowNodes"
              :notice="flowNotice"
            />
          </v-window-item>

          <v-window-item value="measurements">
            <BatchMeasurementsTab
              :measurements="measurements"
              @create="createMeasurementDialog = true"
            />
          </v-window-item>

          <v-window-item value="additions">
            <BatchAdditionsTab
              :additions="additions"
              @create="createAdditionDialog = true"
            />
          </v-window-item>
        </v-window>
      </div>
    </v-card-text>
  </v-card>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <!-- Dialogs -->
  <BatchAdditionDialog
    v-model="createAdditionDialog"
    :form="additionForm"
    @submit="recordAddition"
    @update:form="Object.assign(additionForm, $event)"
  />

  <BatchMeasurementDialog
    v-model="createMeasurementDialog"
    :form="measurementForm"
    @submit="recordMeasurement"
    @update:form="Object.assign(measurementForm, $event)"
  />

  <BatchTimelineExtendedDialog
    v-model="timelineExtendedDialog"
    :form="timelineExtended"
    :gravity-unit="preferences.gravity"
    :temperature-unit="preferences.temperature"
    @submit="recordTimelineExtended"
    @update:form="Object.assign(timelineExtended, $event)"
  />

  <BatchBrewSessionDialog
    v-model="brewSessionDialog"
    :form="brewSessionForm"
    :is-editing="isEditingBrewSession"
    :saving="savingBrewSession"
    :vessels="vessels"
    :volumes="allVolumes"
    @create-volume="openCreateVolumeDialog"
    @submit="saveBrewSession"
    @update:form="Object.assign(brewSessionForm, $event)"
  />

  <BatchVolumeDialog
    v-model="createVolumeDialog"
    :form="volumeForm"
    :saving="savingVolume"
    @submit="createWortVolume"
    @update:form="Object.assign(volumeForm, $event)"
  />

  <BatchHotSideAdditionDialog
    v-model="hotSideAdditionDialog"
    :form="hotSideAdditionForm"
    :saving="savingHotSideAddition"
    @submit="recordHotSideAddition"
    @update:form="Object.assign(hotSideAdditionForm, $event)"
  />

  <BatchHotSideMeasurementDialog
    v-model="hotSideMeasurementDialog"
    :form="hotSideMeasurementForm"
    :saving="savingHotSideMeasurement"
    @submit="recordHotSideMeasurement"
    @update:form="Object.assign(hotSideMeasurementForm, $event)"
  />

  <BatchEditDialog
    v-model="editBatchDialog"
    :batch="selectedBatch"
    :error-message="editBatchError"
    :recipes="recipes"
    :recipes-loading="recipesLoading"
    :saving="savingBatch"
    @submit="saveBatchEdit"
  />

  <BatchDeleteDialog
    v-model="deleteBatchDialog"
    :batch="selectedBatch"
    :deleting="deletingBatch"
    :error-message="deleteBatchError"
    @confirm="confirmDeleteBatch"
  />
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { useApiClient } from '@/composables/useApiClient'
  import {
    type BatchSummary,
    type BrewSession,
    type OccupancyStatus,
    type Addition as ProductionAddition,
    type AdditionType as ProductionAdditionType,
    type Measurement as ProductionMeasurement,
    type Volume as ProductionVolume,
    type UpdateBatchRequest,
    useProductionApi,
    type Vessel,
    type VolumeUnit,
  } from '@/composables/useProductionApi'
  import {
    convertGravity,
    convertTemperature,
    type GravityUnit,
    type TemperatureUnit,
  } from '@/composables/useUnitConversion'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import type { Recipe } from '@/types'
  import {
    type Addition,
    type AdditionType,
    type Batch,
    BatchAdditionDialog,
    BatchAdditionsTab,
    BatchBrewSessionDialog,
    BatchBrewSessionsTab,
    BatchDeleteDialog,
    BatchEditDialog,
    type BatchEditForm,
    BatchFlowTab,
    BatchHotSideAdditionDialog,
    BatchHotSideMeasurementDialog,
    BatchMeasurementDialog,
    BatchMeasurementsTab,
    type BatchProcessPhase,
    BatchSparklineCard,
    BatchSummaryTab,
    BatchTimelineExtendedDialog,
    BatchTimelineTab,
    type BatchVolume,
    BatchVolumeDialog,
    type FlowLink,
    type FlowNode,
    type Measurement,
    type TimelineEvent,
    type Unit,
    type Volume,
    type VolumeRelation,
  } from './batch'

  // Props
  const props = withDefaults(
    defineProps<{
      batchId: number | null
      showBackButton?: boolean
      backButtonText?: string
      backButtonRoute?: string
    }>(),
    {
      showBackButton: false,
      backButtonText: 'Back to All Batches',
      backButtonRoute: '/batches/all',
    },
  )

  // Emits
  const emit = defineEmits<{
    back: []
    cleared: []
  }>()

  const router = useRouter()
  const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'
  const { request } = useApiClient(apiBase)
  const {
    getVessels,
    getVolumes: getProductionVolumes,
    createVolume: createProductionVolume,
    getBrewSessions,
    createBrewSession,
    updateBrewSession,
    getAdditionsByVolume,
    getMeasurementsByVolume,
    createAddition,
    createMeasurement,
    updateBatch,
    deleteBatch,
    getBatchSummary,
    updateOccupancyStatus,
    getRecipes,
  } = useProductionApi()

  const {
    preferences,
    formatTemperaturePreferred,
    formatGravityPreferred,
  } = useUnitPreferences()

  // State
  const loading = ref(false)
  const selectedBatch = ref<Batch | null>(null)
  const volumes = ref<Volume[]>([])
  const batchVolumes = ref<BatchVolume[]>([])
  const processPhases = ref<BatchProcessPhase[]>([])
  const additions = ref<Addition[]>([])
  const measurements = ref<Measurement[]>([])
  const volumeRelations = ref<VolumeRelation[]>([])
  const batchSummary = ref<BatchSummary | null>(null)
  const batchSummaryLoading = ref(false)

  const activeTab = ref('summary')
  const createAdditionDialog = ref(false)
  const createMeasurementDialog = ref(false)

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  const additionForm = reactive({
    target: 'batch' as 'batch' | 'occupancy',
    occupancy_id: '',
    addition_type: 'malt' as AdditionType,
    stage: '',
    inventory_lot_uuid: '',
    amount: '',
    amount_unit: 'ml' as Unit,
    added_at: '',
    notes: '',
  })

  const measurementForm = reactive({
    target: 'batch' as 'batch' | 'occupancy',
    occupancy_id: '',
    kind: '',
    value: '',
    unit: '',
    observed_at: '',
    notes: '',
  })

  const timelineReading = reactive({
    observed_at: '',
    temperature: '',
    gravity: '',
    notes: '',
  })

  const timelineExtendedDialog = ref(false)

  const timelineExtended = reactive({
    observed_at: '',
    temperature: '',
    temperature_unit: '',
    gravity: '',
    gravity_unit: '',
    ph: '',
    ph_unit: '',
    extra_kind: '',
    extra_value: '',
    extra_unit: '',
    notes: '',
  })

  // Brew Session state
  const brewSessions = ref<BrewSession[]>([])
  const vessels = ref<Vessel[]>([])
  const allVolumes = ref<ProductionVolume[]>([])
  const selectedBrewSessionId = ref<number | null>(null)
  const wortAdditions = ref<ProductionAddition[]>([])
  const wortMeasurements = ref<ProductionMeasurement[]>([])

  // Brew Session dialogs and forms
  const brewSessionDialog = ref(false)
  const editingBrewSessionId = ref<number | null>(null)
  const savingBrewSession = ref(false)

  const brewSessionForm = reactive({
    brewed_at: '',
    mash_vessel_id: null as number | null,
    boil_vessel_id: null as number | null,
    wort_volume_id: null as number | null,
    notes: '',
  })

  // Volume creation state
  const createVolumeDialog = ref(false)
  const savingVolume = ref(false)

  const volumeForm = reactive({
    name: '',
    description: '',
    amount: '',
    amount_unit: 'bbl' as VolumeUnit,
  })

  // Hot-side addition/measurement dialogs
  const hotSideAdditionDialog = ref(false)
  const savingHotSideAddition = ref(false)

  const hotSideAdditionForm = reactive({
    addition_type: 'malt' as ProductionAdditionType,
    stage: '',
    inventory_lot_uuid: '',
    amount: '',
    amount_unit: 'ml' as VolumeUnit,
    added_at: '',
    notes: '',
  })

  const hotSideMeasurementDialog = ref(false)
  const savingHotSideMeasurement = ref(false)

  const hotSideMeasurementForm = reactive({
    kind: 'mash_temp',
    value: '',
    unit: '',
    observed_at: '',
    notes: '',
  })

  // Batch edit/delete state
  const editBatchDialog = ref(false)
  const savingBatch = ref(false)
  const editBatchError = ref('')
  const recipes = ref<Recipe[]>([])
  const recipesLoading = ref(false)

  const deleteBatchDialog = ref(false)
  const deletingBatch = ref(false)
  const deleteBatchError = ref('')

  // Computed properties
  const latestProcessPhase = computed(() => getLatest(processPhases.value, item => item.phase_at))
  const latestLiquidPhase = computed(() => getLatest(batchVolumes.value, item => item.phase_at))

  const temperatureSeries = computed(() =>
    buildMeasurementSeries(['temperature', 'temp']),
  )
  const gravitySeries = computed(() =>
    buildMeasurementSeries(['gravity', 'grav', 'sg']),
  )
  const phSeries = computed(() =>
    buildMeasurementSeries(['ph']),
  )

  const timelineReadingReady = computed(() => {
    const hasTemperature = parseTemperatureInput(timelineReading.temperature).value !== null
    const hasGravity = parseNumericInput(timelineReading.gravity) !== null
    const hasNotes = timelineReading.notes.trim().length > 0
    return hasTemperature || hasGravity || hasNotes
  })

  const isEditingBrewSession = computed(() => editingBrewSessionId.value !== null)

  const selectedBrewSession = computed(() =>
    brewSessions.value.find(session => session.id === selectedBrewSessionId.value) ?? null,
  )

  const volumeNameMap = computed(
    () =>
      new Map(
        volumes.value.map(volume => [volume.id, volume.name ?? `Volume ${volume.id}`]),
      ),
  )

  const flowUnit = computed<Unit | null>(() => {
    const counts = new Map<Unit, number>()
    for (const relation of volumeRelations.value) {
      if (!relation.amount || relation.amount <= 0) {
        continue
      }
      counts.set(relation.amount_unit, (counts.get(relation.amount_unit) ?? 0) + 1)
    }
    let selectedUnit: Unit | null = null
    let selectedCount = 0
    for (const [unit, count] of counts.entries()) {
      if (count > selectedCount) {
        selectedUnit = unit
        selectedCount = count
      }
    }
    return selectedUnit
  })

  const flowRelations = computed(() => {
    const relations = volumeRelations.value.filter(relation => relation.amount > 0)
    if (!flowUnit.value) {
      return relations
    }
    return relations.filter(relation => relation.amount_unit === flowUnit.value)
  })

  const flowNotice = computed(() => {
    if (!flowUnit.value) {
      return ''
    }
    const total = volumeRelations.value.filter(relation => relation.amount > 0).length
    const shown = flowRelations.value.length
    if (total > shown) {
      return `Showing ${shown} of ${total} relations measured in ${flowUnit.value} for consistent weights.`
    }
    return ''
  })

  const flowNodes = computed<FlowNode[]>(() => {
    const nodes = new Map<string, FlowNode>()
    const labelFor = (volumeId: number) => volumeNameMap.value.get(volumeId) ?? `Volume ${volumeId}`

    for (const relation of flowRelations.value) {
      const parentId = `volume-${relation.parent_volume_id}`
      const childId = `volume-${relation.child_volume_id}`
      if (!nodes.has(parentId)) {
        nodes.set(parentId, { id: parentId, label: labelFor(relation.parent_volume_id) })
      }
      if (!nodes.has(childId)) {
        nodes.set(childId, { id: childId, label: labelFor(relation.child_volume_id) })
      }
    }

    return Array.from(nodes.values())
  })

  const flowLinks = computed<FlowLink[]>(() => {
    const links = new Map<string, FlowLink>()

    for (const relation of flowRelations.value) {
      const source = `volume-${relation.parent_volume_id}`
      const target = `volume-${relation.child_volume_id}`
      const key = `${source}-${target}-${relation.amount_unit ?? ''}`
      const existing = links.get(key)
      if (existing) {
        existing.value += relation.amount
        existing.label = formatAmount(existing.value, relation.amount_unit)
        continue
      }
      links.set(key, {
        source,
        target,
        value: relation.amount,
        label: formatAmount(relation.amount, relation.amount_unit),
      })
    }

    return Array.from(links.values())
  })

  const timelineItems = computed(() => {
    const items: TimelineEvent[] = []

    for (const addition of additions.value) {
      items.push({
        id: `addition-${addition.id}`,
        title: `Addition: ${addition.addition_type}`,
        subtitle: `${formatAmount(addition.amount, addition.amount_unit)} ${addition.stage ?? ''}`.trim(),
        at: addition.added_at ?? addition.created_at,
        color: 'primary',
        icon: 'mdi-flask-outline',
      })
    }

    const groupedMeasurements = groupMeasurements(measurements.value)
    for (const [index, group] of groupedMeasurements.entries()) {
      const subtitle = orderMeasurementGroup(group)
        .map(measurement => formatMeasurementEntry(measurement))
        .filter(Boolean)
        .join(' | ')
      items.push({
        id: `measurement-group-${index}`,
        title: 'Reading',
        subtitle: subtitle || 'Measurements recorded',
        at: group[0]?.observed_at ?? group[0]?.created_at ?? new Date().toISOString(),
        color: 'secondary',
        icon: 'mdi-thermometer',
      })
    }

    for (const phase of processPhases.value) {
      items.push({
        id: `process-${phase.id}`,
        title: `Process phase: ${phase.process_phase}`,
        subtitle: `Batch ${phase.batch_id}`,
        at: phase.phase_at ?? phase.created_at,
        color: 'success',
        icon: 'mdi-progress-check',
      })
    }

    for (const phase of batchVolumes.value) {
      items.push({
        id: `liquid-${phase.id}`,
        title: `Liquid phase: ${phase.liquid_phase}`,
        subtitle: `Volume ${phase.volume_id}`,
        at: phase.phase_at ?? phase.created_at,
        color: 'warning',
        icon: 'mdi-water',
      })
    }

    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...items].sort((a, b) => toTimestamp(b.at) - toTimestamp(a.at))
  })

  // Watch for batchId changes
  watch(() => props.batchId, async newId => {
    if (newId) {
      await loadBatchData(newId)
    } else {
      clearData()
    }
  }, { immediate: true })

  // Watch for brew session selection to load wort additions/measurements
  watch(selectedBrewSessionId, async sessionId => {
    if (!sessionId) {
      wortAdditions.value = []
      wortMeasurements.value = []
      return
    }
    const session = brewSessions.value.find(s => s.id === sessionId)
    if (session?.wort_volume_id) {
      await loadWortData(session.wort_volume_id)
    } else {
      wortAdditions.value = []
      wortMeasurements.value = []
    }
  })

  onMounted(async () => {
    await loadReferenceData()
  })

  // Exposed methods
  function refresh () {
    if (props.batchId) {
      loadBatchData(props.batchId)
    }
  }

  defineExpose({ refresh })

  function handleBack () {
    emit('back')
    if (props.backButtonRoute) {
      router.push(props.backButtonRoute)
    }
  }

  function clearSelection () {
    emit('cleared')
  }

  function clearData () {
    selectedBatch.value = null
    batchVolumes.value = []
    processPhases.value = []
    additions.value = []
    measurements.value = []
    volumeRelations.value = []
    batchSummary.value = null
    brewSessions.value = []
    selectedBrewSessionId.value = null
    wortAdditions.value = []
    wortMeasurements.value = []
  }

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  function get<T> (path: string) {
    return request<T>(path)
  }

  function post<T> (path: string, payload: unknown) {
    return request<T>(path, { method: 'POST', body: JSON.stringify(payload) })
  }

  async function loadReferenceData () {
    try {
      await Promise.all([loadVolumes(), loadVesselsData(), loadAllVolumesData()])
    } catch (error) {
      console.error('Failed to load reference data:', error)
    }
  }

  async function loadVolumes () {
    volumes.value = await get<Volume[]>('/volumes')
  }

  async function loadVesselsData () {
    try {
      vessels.value = await getVessels()
    } catch (error) {
      console.error('Failed to load vessels:', error)
    }
  }

  async function loadAllVolumesData () {
    try {
      allVolumes.value = await getProductionVolumes()
    } catch (error) {
      console.error('Failed to load volumes:', error)
    }
  }

  async function loadBatchData (batchId: number) {
    loading.value = true
    try {
      const [batchData, batchVolumesData, processPhasesData, additionsData, measurementsData, brewSessionsData] = await Promise.all([
        get<Batch>(`/batches/${batchId}`),
        get<BatchVolume[]>(`/batch-volumes?batch_id=${batchId}`),
        get<BatchProcessPhase[]>(`/batch-process-phases?batch_id=${batchId}`),
        get<Addition[]>(`/additions?batch_id=${batchId}`),
        get<Measurement[]>(`/measurements?batch_id=${batchId}`),
        getBrewSessions(batchId),
      ])

      selectedBatch.value = batchData
      batchVolumes.value = batchVolumesData
      processPhases.value = processPhasesData
      additions.value = additionsData
      measurements.value = measurementsData
      brewSessions.value = brewSessionsData

      // Clear brew session selection when batch changes
      selectedBrewSessionId.value = null
      wortAdditions.value = []
      wortMeasurements.value = []

      // Load batch summary in parallel (non-blocking)
      loadBatchSummary(batchId)

      await loadVolumeRelations(batchVolumesData)
    } catch (error) {
      console.error('Failed to load batch data:', error)
      selectedBatch.value = null
    } finally {
      loading.value = false
    }
  }

  async function loadBatchSummary (batchId: number) {
    batchSummaryLoading.value = true
    batchSummary.value = null
    try {
      batchSummary.value = await getBatchSummary(batchId)
    } catch (error) {
      console.error('Failed to load batch summary:', error)
    } finally {
      batchSummaryLoading.value = false
    }
  }

  async function loadWortData (volumeId: number) {
    try {
      const [additionsData, measurementsData] = await Promise.all([
        getAdditionsByVolume(volumeId),
        getMeasurementsByVolume(volumeId),
      ])
      wortAdditions.value = additionsData
      wortMeasurements.value = measurementsData
    } catch (error) {
      console.error('Failed to load wort data:', error)
    }
  }

  async function loadVolumeRelations (batchVolumeData: BatchVolume[]) {
    const volumeIds = Array.from(
      new Set(batchVolumeData.map(item => item.volume_id)),
    )
    if (volumeIds.length === 0) {
      volumeRelations.value = []
      return
    }

    const results = await Promise.allSettled(
      volumeIds.map(id => get<VolumeRelation[]>(`/volume-relations?volume_id=${id}`)),
    )

    volumeRelations.value = results.flatMap(result =>
      result.status === 'fulfilled' ? result.value : [],
    )
  }

  async function recordAddition () {
    if (!props.batchId) {
      return
    }
    if (!additionForm.amount) {
      return
    }
    if (additionForm.target === 'occupancy' && !additionForm.occupancy_id) {
      return
    }
    try {
      const payload = {
        batch_id: additionForm.target === 'batch' ? props.batchId : null,
        occupancy_id: additionForm.target === 'occupancy' ? toNumber(additionForm.occupancy_id) : null,
        addition_type: additionForm.addition_type,
        stage: normalizeText(additionForm.stage),
        inventory_lot_uuid: normalizeText(additionForm.inventory_lot_uuid),
        amount: toNumber(additionForm.amount),
        amount_unit: additionForm.amount_unit,
        added_at: normalizeDateTime(additionForm.added_at),
        notes: normalizeText(additionForm.notes),
      }
      await post<Addition>('/additions', payload)
      showNotice('Addition recorded')
      additionForm.stage = ''
      additionForm.inventory_lot_uuid = ''
      additionForm.amount = ''
      additionForm.added_at = ''
      additionForm.notes = ''
      createAdditionDialog.value = false
      await loadBatchData(props.batchId)
    } catch (error) {
      handleError(error)
    }
  }

  async function recordMeasurement () {
    if (!props.batchId) {
      return
    }
    if (!measurementForm.kind.trim() || !measurementForm.value) {
      return
    }
    if (measurementForm.target === 'occupancy' && !measurementForm.occupancy_id) {
      return
    }
    try {
      const payload = {
        batch_id: measurementForm.target === 'batch' ? props.batchId : null,
        occupancy_id: measurementForm.target === 'occupancy' ? toNumber(measurementForm.occupancy_id) : null,
        kind: measurementForm.kind.trim(),
        value: toNumber(measurementForm.value),
        unit: normalizeText(measurementForm.unit),
        observed_at: normalizeDateTime(measurementForm.observed_at),
        notes: normalizeText(measurementForm.notes),
      }
      await post<Measurement>('/measurements', payload)
      showNotice('Measurement recorded')
      measurementForm.kind = ''
      measurementForm.value = ''
      measurementForm.unit = ''
      measurementForm.observed_at = ''
      measurementForm.notes = ''
      createMeasurementDialog.value = false
      await loadBatchData(props.batchId)
    } catch (error) {
      handleError(error)
    }
  }

  function openTimelineExtendedDialog () {
    resetTimelineExtended()
    const temperature = parseTemperatureInput(timelineReading.temperature)
    const gravityValue = parseNumericInput(timelineReading.gravity)
    timelineExtended.observed_at = timelineReading.observed_at || nowInputValue()
    timelineExtended.temperature = temperature.value === null ? '' : String(temperature.value)
    timelineExtended.temperature_unit = temperature.unit ?? ''
    timelineExtended.gravity = gravityValue === null ? '' : String(gravityValue)
    timelineExtended.notes = timelineReading.notes
    timelineExtendedDialog.value = true
  }

  async function recordTimelineReading () {
    if (!props.batchId) {
      return
    }

    const temperature = parseTemperatureInput(timelineReading.temperature)
    const gravityValue = parseNumericInput(timelineReading.gravity)
    const noteText = normalizeText(timelineReading.notes)

    if (temperature.value === null && gravityValue === null && !noteText) {
      return
    }

    try {
      const observedAt = timelineReading.observed_at
        ? normalizeDateTime(timelineReading.observed_at)
        : new Date().toISOString()
      const payloads: Array<{
        batch_id: number
        occupancy_id: null
        kind: string
        value: number
        unit: string | null
        observed_at: string | null
        notes: string | null
      }> = []

      if (temperature.value !== null) {
        payloads.push({
          batch_id: props.batchId,
          occupancy_id: null,
          kind: 'temperature',
          value: temperature.value,
          unit: temperature.unit,
          observed_at: observedAt,
          notes: null,
        })
      }
      if (gravityValue !== null) {
        payloads.push({
          batch_id: props.batchId,
          occupancy_id: null,
          kind: 'gravity',
          value: gravityValue,
          unit: null,
          observed_at: observedAt,
          notes: null,
        })
      }
      if (noteText) {
        payloads.push({
          batch_id: props.batchId,
          occupancy_id: null,
          kind: 'note',
          value: 0,
          unit: null,
          observed_at: observedAt,
          notes: noteText,
        })
      }

      await Promise.all(payloads.map(payload => post<Measurement>('/measurements', payload)))
      showNotice(`Recorded ${payloads.length} timeline ${payloads.length === 1 ? 'entry' : 'entries'}`)
      resetTimelineReading()
      await loadBatchData(props.batchId)
    } catch (error) {
      handleError(error)
    }
  }

  async function recordTimelineExtended () {
    if (!props.batchId) {
      return
    }

    const temperatureValue = toNumber(timelineExtended.temperature)
    const gravityValue = toNumber(timelineExtended.gravity)
    const phValue = toNumber(timelineExtended.ph)
    const extraKind = timelineExtended.extra_kind.trim()
    const extraValue = toNumber(timelineExtended.extra_value)
    const noteText = normalizeText(timelineExtended.notes)

    if (extraKind && extraValue === null) {
      return
    }

    if (
      temperatureValue === null
      && gravityValue === null
      && phValue === null
      && !extraKind
      && !noteText
    ) {
      return
    }

    try {
      const observedAt = timelineExtended.observed_at
        ? normalizeDateTime(timelineExtended.observed_at)
        : new Date().toISOString()
      const payloads: Array<{
        batch_id: number
        occupancy_id: null
        kind: string
        value: number
        unit: string | null
        observed_at: string | null
        notes: string | null
      }> = []

      if (temperatureValue !== null) {
        payloads.push({
          batch_id: props.batchId,
          occupancy_id: null,
          kind: 'temperature',
          value: temperatureValue,
          unit: normalizeText(timelineExtended.temperature_unit),
          observed_at: observedAt,
          notes: null,
        })
      }
      if (gravityValue !== null) {
        payloads.push({
          batch_id: props.batchId,
          occupancy_id: null,
          kind: 'gravity',
          value: gravityValue,
          unit: normalizeText(timelineExtended.gravity_unit),
          observed_at: observedAt,
          notes: null,
        })
      }
      if (phValue !== null) {
        payloads.push({
          batch_id: props.batchId,
          occupancy_id: null,
          kind: 'ph',
          value: phValue,
          unit: normalizeText(timelineExtended.ph_unit),
          observed_at: observedAt,
          notes: null,
        })
      }
      if (extraKind && extraValue !== null) {
        payloads.push({
          batch_id: props.batchId,
          occupancy_id: null,
          kind: extraKind,
          value: extraValue,
          unit: normalizeText(timelineExtended.extra_unit),
          observed_at: observedAt,
          notes: null,
        })
      }
      if (noteText) {
        payloads.push({
          batch_id: props.batchId,
          occupancy_id: null,
          kind: 'note',
          value: 0,
          unit: null,
          observed_at: observedAt,
          notes: noteText,
        })
      }

      await Promise.all(payloads.map(payload => post<Measurement>('/measurements', payload)))
      showNotice(`Recorded ${payloads.length} timeline ${payloads.length === 1 ? 'entry' : 'entries'}`)
      resetTimelineReading()
      resetTimelineExtended()
      timelineExtendedDialog.value = false
      await loadBatchData(props.batchId)
    } catch (error) {
      handleError(error)
    }
  }

  function resetTimelineReading () {
    timelineReading.observed_at = ''
    timelineReading.temperature = ''
    timelineReading.gravity = ''
    timelineReading.notes = ''
  }

  function updateTimelineReading (reading: { observed_at: string, temperature: string, gravity: string, notes: string }) {
    timelineReading.observed_at = reading.observed_at
    timelineReading.temperature = reading.temperature
    timelineReading.gravity = reading.gravity
    timelineReading.notes = reading.notes
  }

  function resetTimelineExtended () {
    timelineExtended.observed_at = ''
    timelineExtended.temperature = ''
    timelineExtended.temperature_unit = ''
    timelineExtended.gravity = ''
    timelineExtended.gravity_unit = ''
    timelineExtended.ph = ''
    timelineExtended.ph_unit = ''
    timelineExtended.extra_kind = ''
    timelineExtended.extra_value = ''
    timelineExtended.extra_unit = ''
    timelineExtended.notes = ''
  }

  // ==================== Brew Session Functions ====================

  function selectBrewSession (id: number) {
    selectedBrewSessionId.value = id
  }

  function clearBrewSessionSelection () {
    selectedBrewSessionId.value = null
  }

  function openCreateBrewSessionDialog () {
    editingBrewSessionId.value = null
    brewSessionForm.brewed_at = nowInputValue()
    brewSessionForm.mash_vessel_id = null
    brewSessionForm.boil_vessel_id = null
    brewSessionForm.wort_volume_id = null
    brewSessionForm.notes = ''
    brewSessionDialog.value = true
  }

  function openEditBrewSessionDialog (session: BrewSession) {
    editingBrewSessionId.value = session.id
    brewSessionForm.brewed_at = toLocalDateTimeInput(session.brewed_at)
    brewSessionForm.mash_vessel_id = session.mash_vessel_id
    brewSessionForm.boil_vessel_id = session.boil_vessel_id
    brewSessionForm.wort_volume_id = session.wort_volume_id
    brewSessionForm.notes = session.notes ?? ''
    brewSessionDialog.value = true
  }

  async function saveBrewSession () {
    if (!props.batchId || !brewSessionForm.brewed_at.trim()) {
      return
    }

    savingBrewSession.value = true

    try {
      const payload = {
        batch_id: props.batchId,
        wort_volume_id: brewSessionForm.wort_volume_id,
        mash_vessel_id: brewSessionForm.mash_vessel_id,
        boil_vessel_id: brewSessionForm.boil_vessel_id,
        brewed_at: new Date(brewSessionForm.brewed_at).toISOString(),
        notes: normalizeText(brewSessionForm.notes),
      }

      if (isEditingBrewSession.value && editingBrewSessionId.value) {
        await updateBrewSession(editingBrewSessionId.value, payload)
        showNotice('Brew session updated')
      } else {
        await createBrewSession(payload)
        showNotice('Brew session added')
      }

      brewSessionDialog.value = false
      editingBrewSessionId.value = null
      if (props.batchId) {
        await loadBatchData(props.batchId)
      }
    } catch (error) {
      handleError(error)
    } finally {
      savingBrewSession.value = false
    }
  }

  function openCreateVolumeDialog () {
    volumeForm.name = selectedBatch.value ? `${selectedBatch.value.short_name} Wort` : ''
    volumeForm.description = ''
    volumeForm.amount = ''
    volumeForm.amount_unit = 'bbl'
    createVolumeDialog.value = true
  }

  async function createWortVolume () {
    const amount = Number(volumeForm.amount)
    if (!Number.isFinite(amount) || amount <= 0) {
      return
    }

    savingVolume.value = true

    try {
      const payload = {
        name: normalizeText(volumeForm.name),
        description: normalizeText(volumeForm.description),
        amount: Number(volumeForm.amount),
        amount_unit: volumeForm.amount_unit,
      }

      const created = await createProductionVolume(payload)
      showNotice('Wort volume created')

      // Update volumes list and select the new volume
      await loadAllVolumesData()
      brewSessionForm.wort_volume_id = created.id

      createVolumeDialog.value = false
    } catch (error) {
      handleError(error)
    } finally {
      savingVolume.value = false
    }
  }

  // ==================== Hot-Side Addition/Measurement Functions ====================

  function openCreateHotSideAdditionDialog () {
    hotSideAdditionForm.addition_type = 'malt'
    hotSideAdditionForm.stage = ''
    hotSideAdditionForm.inventory_lot_uuid = ''
    hotSideAdditionForm.amount = ''
    hotSideAdditionForm.amount_unit = 'ml'
    hotSideAdditionForm.added_at = nowInputValue()
    hotSideAdditionForm.notes = ''
    hotSideAdditionDialog.value = true
  }

  async function recordHotSideAddition () {
    const session = selectedBrewSession.value
    if (!session?.wort_volume_id || !hotSideAdditionForm.amount) {
      return
    }

    savingHotSideAddition.value = true

    try {
      const payload = {
        volume_id: session.wort_volume_id,
        addition_type: hotSideAdditionForm.addition_type,
        stage: normalizeText(hotSideAdditionForm.stage),
        inventory_lot_uuid: normalizeText(hotSideAdditionForm.inventory_lot_uuid),
        amount: Number(hotSideAdditionForm.amount),
        amount_unit: hotSideAdditionForm.amount_unit,
        added_at: hotSideAdditionForm.added_at ? new Date(hotSideAdditionForm.added_at).toISOString() : null,
        notes: normalizeText(hotSideAdditionForm.notes),
      }

      await createAddition(payload)
      showNotice('Hot-side addition recorded')

      hotSideAdditionDialog.value = false
      await loadWortData(session.wort_volume_id)
    } catch (error) {
      handleError(error)
    } finally {
      savingHotSideAddition.value = false
    }
  }

  function openCreateHotSideMeasurementDialog () {
    hotSideMeasurementForm.kind = 'mash_temp'
    hotSideMeasurementForm.value = ''
    hotSideMeasurementForm.unit = ''
    hotSideMeasurementForm.observed_at = nowInputValue()
    hotSideMeasurementForm.notes = ''
    hotSideMeasurementDialog.value = true
  }

  async function recordHotSideMeasurement () {
    const session = selectedBrewSession.value
    if (!session?.wort_volume_id || !hotSideMeasurementForm.kind || !hotSideMeasurementForm.value) {
      return
    }

    savingHotSideMeasurement.value = true

    try {
      const payload = {
        volume_id: session.wort_volume_id,
        kind: hotSideMeasurementForm.kind,
        value: Number(hotSideMeasurementForm.value),
        unit: normalizeText(hotSideMeasurementForm.unit) ?? getDefaultUnitForKind(hotSideMeasurementForm.kind),
        observed_at: hotSideMeasurementForm.observed_at ? new Date(hotSideMeasurementForm.observed_at).toISOString() : null,
        notes: normalizeText(hotSideMeasurementForm.notes),
      }

      await createMeasurement(payload)
      showNotice('Hot-side measurement recorded')

      hotSideMeasurementDialog.value = false
      await loadWortData(session.wort_volume_id)
    } catch (error) {
      handleError(error)
    } finally {
      savingHotSideMeasurement.value = false
    }
  }

  function getDefaultUnitForKind (kind: string): string {
    switch (kind) {
      case 'mash_temp':
      case 'boil_temp': {
        return 'F'
      }
      case 'mash_ph': {
        return 'pH'
      }
      case 'pre_boil_gravity':
      case 'original_gravity': {
        return 'SG'
      }
      case 'post_boil_volume': {
        return 'bbl'
      }
      default: {
        return ''
      }
    }
  }

  function toLocalDateTimeInput (isoString: string): string {
    if (!isoString) return ''
    const date = new Date(isoString)
    const pad = (n: number) => String(n).padStart(2, '0')
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
  }

  // ==================== Occupancy Status Functions ====================

  async function changeOccupancyStatus (occupancyId: number, status: OccupancyStatus) {
    if (!props.batchId) {
      return
    }

    try {
      await updateOccupancyStatus(occupancyId, status)
      const statusLabels: Record<string, string> = {
        fermenting: 'Fermenting',
        conditioning: 'Conditioning',
        cold_crashing: 'Cold Crashing',
        dry_hopping: 'Dry Hopping',
        carbonating: 'Carbonating',
        holding: 'Holding',
        packaging: 'Packaging',
      }
      const label = statusLabels[status] ?? status.charAt(0).toUpperCase() + status.slice(1).replace(/_/g, ' ')
      showNotice(`Status updated to ${label}`)
      // Reload batch summary to reflect the change
      await loadBatchSummary(props.batchId)
    } catch (error) {
      handleError(error)
    }
  }

  // ==================== Batch Edit/Delete Functions ====================

  async function loadRecipesData () {
    recipesLoading.value = true
    try {
      recipes.value = await getRecipes()
    } catch (error) {
      console.error('Failed to load recipes:', error)
    } finally {
      recipesLoading.value = false
    }
  }

  function openEditDialog () {
    if (!selectedBatch.value) return
    editBatchError.value = ''
    loadRecipesData()
    editBatchDialog.value = true
  }

  async function saveBatchEdit (form: BatchEditForm) {
    if (!props.batchId || !selectedBatch.value) return

    savingBatch.value = true
    editBatchError.value = ''

    try {
      const payload: UpdateBatchRequest = {
        short_name: form.short_name.trim(),
        brew_date: form.brew_date ? normalizeDateOnly(form.brew_date) : null,
        recipe_id: form.recipe_id,
        notes: normalizeText(form.notes),
      }

      await updateBatch(props.batchId, payload)
      showNotice('Batch updated')
      editBatchDialog.value = false
      await loadBatchData(props.batchId)
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update batch'
      editBatchError.value = message
    } finally {
      savingBatch.value = false
    }
  }

  function openDeleteDialog () {
    if (!selectedBatch.value) return
    deleteBatchError.value = ''
    deleteBatchDialog.value = true
  }

  async function confirmDeleteBatch () {
    if (!props.batchId || !selectedBatch.value) return

    deletingBatch.value = true
    deleteBatchError.value = ''

    try {
      await deleteBatch(props.batchId)
      showNotice('Batch deleted')
      deleteBatchDialog.value = false
      // Navigate back to batch list
      router.push('/batches/all')
    } catch (error) {
      // Check for 409 Conflict (batch has dependencies)
      if (error instanceof Error && error.message.includes('409')) {
        deleteBatchError.value = 'Cannot delete this batch because it has associated brew sessions, measurements, or other data. Remove those first.'
      } else if (error instanceof Error && (error.message.toLowerCase().includes('conflict') || error.message.toLowerCase().includes('dependencies'))) {
        deleteBatchError.value = 'Cannot delete this batch because it has associated brew sessions, measurements, or other data. Remove those first.'
      } else {
        const message = error instanceof Error ? error.message : 'Failed to delete batch'
        deleteBatchError.value = message
      }
    } finally {
      deletingBatch.value = false
    }
  }

  function normalizeDateOnly (value: string) {
    return value ? new Date(`${value}T00:00:00Z`).toISOString() : null
  }

  // ==================== Helper Functions ====================

  function handleError (error: unknown) {
    const message = error instanceof Error ? error.message : 'Unexpected error'
    showNotice(message, 'error')
  }

  function normalizeText (value: string) {
    const trimmed = value.trim()
    return trimmed.length > 0 ? trimmed : null
  }

  function normalizeDateTime (value: string) {
    return value ? new Date(value).toISOString() : null
  }

  function toNumber (value: string | number | null) {
    if (value === null || value === undefined || value === '') {
      return null
    }
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : null
  }

  function formatDate (value: string | null | undefined) {
    if (!value) {
      return 'Unknown'
    }
    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
    }).format(new Date(value))
  }

  function formatDateTime (value: string | null | undefined) {
    if (!value) {
      return 'Unknown'
    }
    return new Intl.DateTimeFormat('en-US', {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(new Date(value))
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

  function toTimestamp (value: string | null | undefined) {
    if (!value) {
      return 0
    }
    return new Date(value).getTime()
  }

  function sortByTime<T> (items: T[], selector: (item: T) => string | null | undefined) {
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...items].sort((a, b) => toTimestamp(selector(b)) - toTimestamp(selector(a)))
  }

  function getLatest<T> (items: T[], selector: (item: T) => string | null | undefined) {
    const sorted = sortByTime(items, selector)
    return sorted.length > 0 ? sorted[0] : null
  }

  function nowInputValue () {
    const now = new Date()
    const pad = (value: number) => String(value).padStart(2, '0')
    const year = now.getFullYear()
    const month = pad(now.getMonth() + 1)
    const day = pad(now.getDate())
    const hours = pad(now.getHours())
    const minutes = pad(now.getMinutes())
    return `${year}-${month}-${day}T${hours}:${minutes}`
  }

  function parseNumericInput (value: string) {
    const normalized = value.trim().replace(/,/g, '')
    if (!normalized) {
      return null
    }
    const match = normalized.match(/-?\d*\.?\d+/)
    if (!match) {
      return null
    }
    const parsed = Number(match[0])
    return Number.isFinite(parsed) ? parsed : null
  }

  function parseTemperatureInput (value: string) {
    const parsedValue = parseNumericInput(value)
    if (parsedValue === null) {
      return { value: null, unit: null }
    }
    const unitMatch = value.match(/([cf])\s*$/i)
    if (unitMatch && unitMatch[1]) {
      return { value: parsedValue, unit: unitMatch[1].toUpperCase() }
    }
    return { value: parsedValue, unit: null }
  }

  function groupMeasurements (items: Measurement[]) {
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    const sorted = [...items].sort(
      (a, b) => measurementTimestamp(b) - measurementTimestamp(a),
    )
    const groups: Measurement[][] = []
    const thresholdMs = 2 * 60 * 1000
    let current: Measurement[] = []
    let anchor: number | null = null

    for (const measurement of sorted) {
      const timestamp = measurementTimestamp(measurement)
      if (current.length === 0) {
        current = [measurement]
        anchor = timestamp
        continue
      }
      if (anchor !== null && Math.abs(anchor - timestamp) <= thresholdMs) {
        current.push(measurement)
        continue
      }
      groups.push(current)
      current = [measurement]
      anchor = timestamp
    }

    if (current.length > 0) {
      groups.push(current)
    }

    return groups
  }

  function orderMeasurementGroup (group: Measurement[]) {
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...group].sort((a, b) => measurementPriority(a) - measurementPriority(b))
  }

  function measurementPriority (measurement: Measurement) {
    const normalized = normalizeMeasurementKind(measurement.kind)
    if (isNoteMeasurement(measurement)) {
      return 90
    }
    if (normalized === 'temperature' || normalized === 'temp') {
      return 10
    }
    if (normalized === 'gravity' || normalized === 'grav' || normalized === 'sg') {
      return 20
    }
    if (normalized === 'ph') {
      return 30
    }
    return 50
  }

  function measurementTimestamp (measurement: Measurement) {
    return toTimestamp(measurement.observed_at ?? measurement.created_at)
  }

  function formatMeasurementEntry (measurement: Measurement) {
    if (isNoteMeasurement(measurement)) {
      return measurement.notes ? `Note: ${measurement.notes}` : 'Note'
    }
    const label = formatMeasurementKind(measurement.kind)
    return `${label} ${formatValue(measurement.value, measurement.unit)}`
  }

  function formatMeasurementKind (kind: string) {
    const normalized = normalizeMeasurementKind(kind)
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

  function buildMeasurementSeries (kinds: string[]) {
    const normalizedKinds = kinds.map(kind => normalizeMeasurementKind(kind))
    const filtered = measurements.value.filter(measurement =>
      matchesMeasurementKind(measurement.kind, normalizedKinds),
    )
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    const ordered = [...filtered].sort(
      (a, b) => measurementTimestamp(a) - measurementTimestamp(b),
    )

    // Determine if this is temperature or gravity series for unit conversion
    const isTemperature = normalizedKinds.some(k => k === 'temperature' || k === 'temp')
    const isGravity = normalizedKinds.some(k => k === 'gravity' || k === 'grav' || k === 'sg')

    // Convert values to preferred units for sparkline display
    const values = ordered
      .map(measurement => {
        if (!Number.isFinite(measurement.value)) {
          return null
        }

        if (isTemperature) {
          // Determine source unit from measurement - default to Fahrenheit if not specified
          const sourceUnit = normalizeTemperatureUnit(measurement.unit)
          return convertTemperature(measurement.value, sourceUnit, preferences.value.temperature)
        }

        if (isGravity) {
          // Determine source unit from measurement - default to SG if not specified
          const sourceUnit = normalizeGravityUnit(measurement.unit)
          return convertGravity(measurement.value, sourceUnit, preferences.value.gravity)
        }

        // For other measurements (like pH), no conversion needed
        return measurement.value
      })
      .filter((value): value is number => value !== null && Number.isFinite(value))

    const latest = getLatest(filtered, item => item.observed_at ?? item.created_at)

    // Format latest label using preferred units
    let latestLabel = 'n/a'
    if (latest && Number.isFinite(latest.value)) {
      if (isTemperature) {
        const sourceUnit = normalizeTemperatureUnit(latest.unit)
        latestLabel = formatTemperaturePreferred(latest.value, sourceUnit)
      } else if (isGravity) {
        const sourceUnit = normalizeGravityUnit(latest.unit)
        latestLabel = formatGravityPreferred(latest.value, sourceUnit)
      } else {
        latestLabel = formatValue(latest.value, latest.unit)
      }
    }

    return {
      values,
      latest,
      latestLabel,
    }
  }

  /**
   * Normalize a temperature unit string to a TemperatureUnit type.
   * Defaults to 'f' (Fahrenheit) if not recognized.
   */
  function normalizeTemperatureUnit (unit: string | null | undefined): TemperatureUnit {
    if (!unit) return 'f' // Default to Fahrenheit
    const normalized = unit.trim().toLowerCase()
    if (normalized === 'c' || normalized === 'celsius' || normalized === '°c') {
      return 'c'
    }
    return 'f' // Default to Fahrenheit for 'f', 'fahrenheit', '°f', or unknown
  }

  /**
   * Normalize a gravity unit string to a GravityUnit type.
   * Defaults to 'sg' if not recognized.
   */
  function normalizeGravityUnit (unit: string | null | undefined): GravityUnit {
    if (!unit) return 'sg' // Default to SG
    const normalized = unit.trim().toLowerCase()
    if (normalized === 'plato' || normalized === '°p' || normalized === 'p') {
      return 'plato'
    }
    return 'sg' // Default to SG for 'sg', or unknown
  }

  function matchesMeasurementKind (value: string, kinds: string[]) {
    const normalized = normalizeMeasurementKind(value)
    if (!normalized) {
      return false
    }
    return kinds.some(kind => normalized.includes(kind))
  }

  function normalizeMeasurementKind (value: string) {
    return value.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
  }

  function isNoteMeasurement (measurement: Measurement) {
    return matchesMeasurementKind(measurement.kind, ['note'])
  }
</script>

<style scoped>
.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.mini-card {
  height: 100%;
}

.batch-tabs :deep(.v-tab) {
  text-transform: none;
  font-weight: 600;
}

/* Ensure tables scroll horizontally on mobile */
.data-table {
  overflow-x: auto;
}

.data-table :deep(.v-table__wrapper) {
  overflow-x: auto;
}

.data-table :deep(th) {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  white-space: nowrap;
}

.data-table :deep(td) {
  font-size: 0.85rem;
}

/* Timeline time picker responsive width */
.timeline-time-picker {
  min-width: 280px;
  max-width: 90vw;
}

.spark-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  border-radius: 12px;
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  background: rgba(var(--v-theme-surface), 0.4);
}

.spark-meta {
  min-width: 86px;
}

.spark-chart {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.spark-chart svg {
  width: 100%;
  height: 36px;
}

.spark-placeholder {
  font-size: 0.75rem;
  color: rgba(var(--v-theme-on-surface), 0.5);
}

.spark-line {
  fill: none;
  stroke: rgb(var(--spark-color));
  stroke-width: 2;
}

.spark-area {
  fill: rgba(var(--spark-color), 0.2);
}

.brew-session-list {
  max-height: 280px;
  overflow: auto;
}

.metric-card {
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(var(--v-theme-surface), 0.5);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  text-align: center;
}

.metric-card--highlight {
  background: rgba(var(--v-theme-primary), 0.08);
  border-color: rgba(var(--v-theme-primary), 0.2);
}

.metric-card--warning {
  background: rgba(var(--v-theme-warning), 0.08);
  border-color: rgba(var(--v-theme-warning), 0.2);
}

.metric-label {
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-bottom: 4px;
}

.metric-value {
  font-size: 1.25rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.87);
}

.brew-session-summary-list {
  border-radius: 8px;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
