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
      <v-btn
        v-if="selectedBatch"
        aria-label="Edit batch"
        icon="mdi-pencil"
        size="small"
        variant="text"
        @click="openEditDialog"
      />
      <v-btn
        v-if="selectedBatch"
        aria-label="Delete batch"
        color="error"
        icon="mdi-delete"
        size="small"
        variant="text"
        @click="openDeleteDialog"
      />
      <v-btn v-if="!showBackButton" size="small" variant="text" @click="clearSelection">Clear</v-btn>
    </v-card-title>
    <v-card-text>
      <v-alert
        v-if="!batchUuid"
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
          <v-tab value="brew-day">Brew Day</v-tab>
          <v-tab value="brew-sessions">Brew Sessions</v-tab>
          <v-tab value="fermentation">Fermentation</v-tab>
          <v-tab value="timeline">Timeline</v-tab>
          <v-tab value="flow">Flow</v-tab>
          <v-tab value="measurements">Measurements</v-tab>
          <v-tab value="additions">Additions</v-tab>
        </v-tabs>

        <v-window v-model="activeTab" class="mt-4">
          <v-window-item value="summary">
            <!-- Brew Day Wizard button -->
            <v-btn
              v-if="showStartBrewDay"
              block
              class="mb-4 d-sm-none"
              color="primary"
              min-height="44"
              prepend-icon="mdi-kettle-steam"
              size="large"
              @click="openBrewDayWizard"
            >
              {{ brewDayButtonLabel }}
            </v-btn>
            <v-btn
              v-if="showStartBrewDay"
              class="mb-4 d-none d-sm-inline-flex"
              color="primary"
              min-height="44"
              prepend-icon="mdi-kettle-steam"
              size="large"
              @click="openBrewDayWizard"
            >
              {{ brewDayButtonLabel }}
            </v-btn>

            <BatchSummaryTab
              :has-volumes="hasBatchVolumes"
              :is-finished="isFinished"
              :loading="batchSummaryLoading"
              :summary="batchSummary"
              @assign-fermenter="openAssignFermenterDialog"
              @mark-empty="openMarkEmptyDialog"
              @occupancy-status-change="changeOccupancyStatus"
              @package="openPackagingDialog"
              @transfer="openTransferDialog"
            />
          </v-window-item>

          <v-window-item value="brew-day">
            <BatchBrewDayTab
              v-if="activeTab === 'brew-day'"
              :batch-uuid="selectedBatch?.uuid ?? null"
              :recipe-name="selectedBatch?.recipe_name ?? null"
              :recipe-uuid="selectedBatch?.recipe_uuid ?? null"
            />
          </v-window-item>

          <v-window-item value="brew-sessions">
            <BatchBrewSessionsTab
              v-if="activeTab === 'brew-sessions'"
              :additions="wortAdditions"
              :measurements="wortMeasurements"
              :selected-session-uuid="selectedBrewSessionUuid"
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

          <v-window-item value="fermentation">
            <FermentationCurve
              v-if="activeTab === 'fermentation'"
              :batch-summary="batchSummary"
              :measurements="measurements"
              :target-fg="recipeTargetFg"
              :target-og="recipeTargetOg"
              @go-to-timeline="activeTab = 'timeline'"
            />
          </v-window-item>

          <v-window-item value="timeline">
            <BatchTimelineTab
              v-if="activeTab === 'timeline'"
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
              v-if="activeTab === 'flow'"
              :links="flowLinks"
              :nodes="flowNodes"
              :notice="flowNotice"
            />
          </v-window-item>

          <v-window-item value="measurements">
            <BatchMeasurementsTab
              v-if="activeTab === 'measurements'"
              :measurements="measurements"
              @create="createMeasurementDialog = true"
            />
          </v-window-item>

          <v-window-item value="additions">
            <BatchAdditionsTab
              v-if="activeTab === 'additions'"
              :additions="additions"
              @create="createAdditionDialog = true"
            />
          </v-window-item>
        </v-window>
      </div>
    </v-card-text>
  </v-card>

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

  <BatchAssignFermenterDialog
    v-model="assignFermenterDialog"
    :active-occupancies="activeOccupancies"
    :batch-volumes="batchProductionVolumes"
    :vessels="vessels"
    @assigned="handleFermenterAssigned"
  />

  <BatchMarkEmptyDialog
    v-model="markEmptyDialog"
    :batch-name="selectedBatch?.short_name ?? ''"
    :occupancy="markEmptyOccupancy"
    :vessel-name="batchSummary?.current_vessel ?? ''"
    @emptied="handleVesselEmptied"
  />

  <BrewDayWizard
    v-if="selectedBatch"
    v-model="brewDayWizardDialog"
    :batch="selectedBatch"
    :brew-sessions="brewSessions"
    :occupancies="activeOccupancies"
    :vessels="vessels"
    :volumes="batchProductionVolumes"
    @completed="handleBrewDayWizardCompleted"
  />

  <TransferDialog
    v-model="transferDialog"
    :source-batch="selectedBatch"
    :source-occupancy="transferOccupancy"
    :source-vessel="transferVessel"
    :source-volume="transferVolume"
    @transferred="handleTransferCompleted"
  />

  <PackagingDialog
    v-model="packagingDialog"
    :source-batch="selectedBatch"
    :source-occupancy="packagingOccupancy"
    :source-vessel="packagingVessel"
    :source-volume="packagingVolume"
    @packaged="handlePackagingCompleted"
  />
</template>

<script lang="ts" setup>
  import type { BatchSummary,
                BrewSession,
                GravityUnit,
                Occupancy,
                OccupancyStatus,
                Addition as ProductionAddition,
                AdditionType as ProductionAdditionType,
                Measurement as ProductionMeasurement,
                Volume as ProductionVolume,
                Recipe,
                TemperatureUnit,
                UpdateBatchRequest,
                Vessel,
                VolumeUnit } from '@/types'
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'
  import { formatDate, formatDateTime, useAdditionTypeFormatters, useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import {
    convertGravity,
    convertTemperature,
  } from '@/composables/useUnitConversion'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { normalizeDateOnly, normalizeDateTime, normalizeText, nowInputValue, toLocalDateTimeInput, toNumber } from '@/utils/normalize'
  import {
    type Addition,
    type AdditionType,
    type Batch,
    BatchAdditionDialog,
    BatchAdditionsTab,
    BatchAssignFermenterDialog,
    BatchBrewDayTab,
    BatchBrewSessionDialog,
    BatchBrewSessionsTab,
    BatchDeleteDialog,
    BatchEditDialog,
    type BatchEditForm,
    BatchFlowTab,
    BatchHotSideAdditionDialog,
    BatchHotSideMeasurementDialog,
    BatchMarkEmptyDialog,
    BatchMeasurementDialog,
    BatchMeasurementsTab,
    type BatchProcessPhase,
    BatchSparklineCard,
    BatchSummaryTab,
    BatchTimelineExtendedDialog,
    BatchTimelineTab,
    type BatchVolume,
    BatchVolumeDialog,
    BrewDayWizard,
    type FlowLink,
    type FlowNode,
    type Measurement,
    PackagingDialog,
    type TimelineEvent,
    type Unit,
    type Volume,
    type VolumeRelation,
  } from './batch'
  import { FermentationCurve, TransferDialog } from './fermentation'

  // Props
  const props = withDefaults(
    defineProps<{
      batchUuid: string | null
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
  const {
    getVessels,
    getVolumes: getProductionVolumes,
    createVolume: createProductionVolume,
    getBrewSessions,
    createBrewSession,
    updateBrewSession,
    getAdditionsByVolume,
    getAdditionsByBatch,
    getMeasurementsByVolume,
    getMeasurementsByBatch,
    createAddition,
    createMeasurement,
    updateBatch,
    deleteBatch,
    getBatchSummary,
    getActiveOccupancies,
    getOccupancy,
    updateOccupancyStatus,
    getRecipes,
    getRecipe,
    getBatchProcessPhases,
    getBatchVolumes,
    getVolumeRelations,
    getBatch,
    request,
  } = useProductionApi()

  const {
    preferences,
    formatTemperaturePreferred,
    formatGravityPreferred,
  } = useUnitPreferences()

  const { formatOccupancyStatus } = useOccupancyStatusFormatters()
  const { formatAdditionType } = useAdditionTypeFormatters()

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
  const batchRecipe = ref<Recipe | null>(null)

  const activeTab = ref('summary')
  const createAdditionDialog = ref(false)
  const createMeasurementDialog = ref(false)

  const { showNotice } = useSnackbar()

  const additionForm = reactive({
    target: 'batch' as 'batch' | 'occupancy',
    occupancy_uuid: '',
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
    occupancy_uuid: '',
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
  const selectedBrewSessionUuid = ref<string | null>(null)
  const wortAdditions = ref<ProductionAddition[]>([])
  const wortMeasurements = ref<ProductionMeasurement[]>([])

  // Brew Session dialogs and forms
  const brewSessionDialog = ref(false)
  const editingBrewSessionUuid = ref<string | null>(null)
  const savingBrewSession = ref(false)

  const brewSessionForm = reactive({
    brewed_at: '',
    mash_vessel_uuid: null as string | null,
    boil_vessel_uuid: null as string | null,
    wort_volume_uuid: null as string | null,
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

  // Assign fermenter dialog state
  const assignFermenterDialog = ref(false)
  const activeOccupancies = ref<Occupancy[]>([])
  const batchProductionVolumes = ref<ProductionVolume[]>([])

  // Mark empty dialog state
  const markEmptyDialog = ref(false)
  const markEmptyOccupancy = ref<Occupancy | null>(null)

  // Brew Day Wizard state
  const brewDayWizardDialog = ref(false)

  // Transfer dialog state
  const transferDialog = ref(false)
  const transferOccupancy = ref<Occupancy | null>(null)
  const transferVessel = ref<Vessel | null>(null)
  const transferVolume = ref<ProductionVolume | null>(null)

  // Packaging dialog state
  const packagingDialog = ref(false)
  const packagingOccupancy = ref<Occupancy | null>(null)
  const packagingVessel = ref<Vessel | null>(null)
  const packagingVolume = ref<ProductionVolume | null>(null)

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

  const hasBatchVolumes = computed(() => batchProductionVolumes.value.length > 0)

  // Recipe target specs for fermentation curve
  const recipeTargetOg = computed(() => batchRecipe.value?.target_og ?? null)
  const recipeTargetFg = computed(() => batchRecipe.value?.target_fg ?? null)

  // Whether the batch is in the 'finished' process phase
  const isFinished = computed(() =>
    batchSummary.value?.current_phase === 'finished',
  )

  // Brew Day Wizard computed
  const showStartBrewDay = computed(() => {
    if (!selectedBatch.value) return false
    // Hide on finished batches
    if (isFinished.value) return false
    // Show when batch has a recipe
    if (!selectedBatch.value.recipe_uuid) return false
    // Show when batch has no occupancy (fermenter not yet assigned)
    if (batchSummary.value?.current_vessel) return false
    return true
  })

  const brewDayButtonLabel = computed(() => {
    if (brewSessions.value.length > 0) return 'Continue Brew Day'
    return 'Start Brew Day'
  })

  const isEditingBrewSession = computed(() => editingBrewSessionUuid.value !== null)

  const selectedBrewSession = computed(() =>
    brewSessions.value.find(session => session.uuid === selectedBrewSessionUuid.value) ?? null,
  )

  const volumeNameMap = computed(
    () =>
      new Map(
        volumes.value.map(volume => [volume.uuid, volume.name ?? 'Unnamed Volume']),
      ),
  )

  const flowUnit = computed<string | null>(() => {
    const counts = new Map<string, number>()
    for (const relation of volumeRelations.value) {
      if (!relation.amount || relation.amount <= 0) {
        continue
      }
      counts.set(relation.amount_unit, (counts.get(relation.amount_unit) ?? 0) + 1)
    }
    let selectedUnit: string | null = null
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
    const labelFor = (volumeUuid: string) => volumeNameMap.value.get(volumeUuid) ?? 'Unnamed Volume'

    for (const relation of flowRelations.value) {
      const parentId = `volume-${relation.parent_volume_uuid}`
      const childId = `volume-${relation.child_volume_uuid}`
      if (!nodes.has(parentId)) {
        nodes.set(parentId, { id: parentId, label: labelFor(relation.parent_volume_uuid) })
      }
      if (!nodes.has(childId)) {
        nodes.set(childId, { id: childId, label: labelFor(relation.child_volume_uuid) })
      }
    }

    return Array.from(nodes.values())
  })

  const flowLinks = computed<FlowLink[]>(() => {
    const links = new Map<string, FlowLink>()

    for (const relation of flowRelations.value) {
      const source = `volume-${relation.parent_volume_uuid}`
      const target = `volume-${relation.child_volume_uuid}`
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
        id: `addition-${addition.uuid}`,
        title: `Addition: ${formatAdditionType(addition.addition_type)}`,
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
        id: `process-${phase.uuid}`,
        title: `Process phase: ${phase.process_phase}`,
        subtitle: selectedBatch.value?.short_name ?? 'Batch',
        at: phase.phase_at ?? phase.created_at,
        color: 'success',
        icon: 'mdi-progress-check',
      })
    }

    for (const phase of batchVolumes.value) {
      items.push({
        id: `liquid-${phase.uuid}`,
        title: `Liquid phase: ${phase.liquid_phase}`,
        subtitle: volumeNameMap.value.get(phase.volume_uuid) ?? 'Unnamed Volume',
        at: phase.phase_at ?? phase.created_at,
        color: 'warning',
        icon: 'mdi-water',
      })
    }

    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...items].sort((a, b) => toTimestamp(b.at) - toTimestamp(a.at))
  })

  // Watch for batchUuid changes
  watch(() => props.batchUuid, async newUuid => {
    if (newUuid) {
      await loadBatchData(newUuid)
    } else {
      clearData()
    }
  }, { immediate: true })

  // Watch for brew session selection to load wort additions/measurements
  watch(selectedBrewSessionUuid, async sessionUuid => {
    if (!sessionUuid) {
      wortAdditions.value = []
      wortMeasurements.value = []
      return
    }
    const session = brewSessions.value.find(s => s.uuid === sessionUuid)
    if (session?.wort_volume_uuid) {
      await loadWortData(session.wort_volume_uuid)
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
    if (props.batchUuid) {
      loadBatchData(props.batchUuid)
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
    batchRecipe.value = null
    batchProductionVolumes.value = []
    activeOccupancies.value = []
    markEmptyOccupancy.value = null
    brewSessions.value = []
    selectedBrewSessionUuid.value = null
    wortAdditions.value = []
    wortMeasurements.value = []
  }

  async function loadReferenceData () {
    try {
      await Promise.all([loadVolumes(), loadVesselsData(), loadAllVolumesData()])
    } catch (error) {
      console.error('Failed to load reference data:', error)
    }
  }

  async function loadVolumes () {
    volumes.value = await getProductionVolumes()
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

  async function loadBatchData (batchUuid: string) {
    loading.value = true
    try {
      const [batchData, batchVolumesData, processPhasesData, additionsData, measurementsData, brewSessionsData, batchProdVolumesData] = await Promise.all([
        getBatch(batchUuid),
        getBatchVolumes(`batch_uuid=${batchUuid}`),
        getBatchProcessPhases(batchUuid),
        getAdditionsByBatch(batchUuid),
        getMeasurementsByBatch(batchUuid),
        getBrewSessions(batchUuid),
        request<ProductionVolume[]>(`/volumes?batch_uuid=${batchUuid}`),
      ])

      selectedBatch.value = batchData
      batchVolumes.value = batchVolumesData
      processPhases.value = processPhasesData
      additions.value = additionsData
      measurements.value = measurementsData
      brewSessions.value = brewSessionsData
      batchProductionVolumes.value = batchProdVolumesData

      // Clear brew session selection when batch changes
      selectedBrewSessionUuid.value = null
      wortAdditions.value = []
      wortMeasurements.value = []

      // Load batch summary and recipe in parallel (non-blocking)
      loadBatchSummary(batchUuid)
      loadBatchRecipe(batchData.recipe_uuid)

      await loadVolumeRelations(batchVolumesData)
    } catch (error) {
      console.error('Failed to load batch data:', error)
      selectedBatch.value = null
    } finally {
      loading.value = false
    }
  }

  async function loadBatchRecipe (recipeUuid: string | null) {
    batchRecipe.value = null
    if (!recipeUuid) return
    try {
      batchRecipe.value = await getRecipe(recipeUuid)
    } catch (error) {
      console.error('Failed to load batch recipe:', error)
    }
  }

  async function loadBatchSummary (batchUuid: string) {
    batchSummaryLoading.value = true
    batchSummary.value = null
    try {
      batchSummary.value = await getBatchSummary(batchUuid)
    } catch (error) {
      console.error('Failed to load batch summary:', error)
    } finally {
      batchSummaryLoading.value = false
    }
  }

  async function loadWortData (volumeUuid: string) {
    try {
      const [additionsData, measurementsData] = await Promise.all([
        getAdditionsByVolume(volumeUuid),
        getMeasurementsByVolume(volumeUuid),
      ])
      wortAdditions.value = additionsData
      wortMeasurements.value = measurementsData
    } catch (error) {
      console.error('Failed to load wort data:', error)
    }
  }

  async function loadVolumeRelations (batchVolumeData: BatchVolume[]) {
    const volumeUuids = Array.from(
      new Set(batchVolumeData.map(item => item.volume_uuid)),
    )
    if (volumeUuids.length === 0) {
      volumeRelations.value = []
      return
    }

    const results = await Promise.allSettled(
      volumeUuids.map(uuid => getVolumeRelations(`volume_uuid=${uuid}`)),
    )

    volumeRelations.value = results.flatMap(result =>
      result.status === 'fulfilled' ? result.value : [],
    )
  }

  function openCreateBrewSessionDialog () {
    editingBrewSessionUuid.value = null
    brewSessionForm.brewed_at = nowInputValue()
    brewSessionForm.mash_vessel_uuid = null
    brewSessionForm.boil_vessel_uuid = null
    brewSessionForm.wort_volume_uuid = null
    brewSessionForm.notes = ''
    brewSessionDialog.value = true
  }

  function openEditBrewSessionDialog (session: BrewSession) {
    editingBrewSessionUuid.value = session.uuid
    brewSessionForm.brewed_at = toLocalDateTimeInput(session.brewed_at)
    brewSessionForm.mash_vessel_uuid = session.mash_vessel_uuid
    brewSessionForm.boil_vessel_uuid = session.boil_vessel_uuid
    brewSessionForm.wort_volume_uuid = session.wort_volume_uuid
    brewSessionForm.notes = session.notes ?? ''
    brewSessionDialog.value = true
  }

  async function saveBrewSession () {
    if (!props.batchUuid || !brewSessionForm.brewed_at.trim()) {
      return
    }

    savingBrewSession.value = true

    try {
      const payload = {
        batch_uuid: props.batchUuid,
        wort_volume_uuid: brewSessionForm.wort_volume_uuid,
        mash_vessel_uuid: brewSessionForm.mash_vessel_uuid,
        boil_vessel_uuid: brewSessionForm.boil_vessel_uuid,
        brewed_at: new Date(brewSessionForm.brewed_at).toISOString(),
        notes: normalizeText(brewSessionForm.notes),
      }

      if (isEditingBrewSession.value && editingBrewSessionUuid.value) {
        await updateBrewSession(editingBrewSessionUuid.value, payload)
        showNotice('Brew session updated')
      } else {
        await createBrewSession(payload)
        showNotice('Brew session added')
      }

      brewSessionDialog.value = false
      editingBrewSessionUuid.value = null
      if (props.batchUuid) {
        await loadBatchData(props.batchUuid)
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
      brewSessionForm.wort_volume_uuid = created.uuid

      createVolumeDialog.value = false
    } catch (error) {
      handleError(error)
    } finally {
      savingVolume.value = false
    }
  }

  async function recordAddition () {
    if (!props.batchUuid) {
      return
    }
    if (!additionForm.amount) {
      return
    }
    if (additionForm.target === 'occupancy' && !additionForm.occupancy_uuid) {
      return
    }
    try {
      const amount = toNumber(additionForm.amount)
      if (amount === null) return
      const payload = {
        batch_uuid: additionForm.target === 'batch' ? props.batchUuid : null,
        occupancy_uuid: additionForm.target === 'occupancy' ? additionForm.occupancy_uuid : null,
        addition_type: additionForm.addition_type,
        stage: normalizeText(additionForm.stage),
        inventory_lot_uuid: normalizeText(additionForm.inventory_lot_uuid),
        amount,
        amount_unit: additionForm.amount_unit,
        added_at: normalizeDateTime(additionForm.added_at),
        notes: normalizeText(additionForm.notes),
      }
      await createAddition(payload)
      showNotice('Addition recorded')
      additionForm.stage = ''
      additionForm.inventory_lot_uuid = ''
      additionForm.amount = ''
      additionForm.added_at = ''
      additionForm.notes = ''
      createAdditionDialog.value = false
      await loadBatchData(props.batchUuid)
    } catch (error) {
      handleError(error)
    }
  }

  async function recordMeasurement () {
    if (!props.batchUuid) {
      return
    }
    if (!measurementForm.kind.trim() || !measurementForm.value) {
      return
    }
    if (measurementForm.target === 'occupancy' && !measurementForm.occupancy_uuid) {
      return
    }
    try {
      const value = toNumber(measurementForm.value)
      if (value === null) return
      const payload = {
        batch_uuid: measurementForm.target === 'batch' ? props.batchUuid : null,
        occupancy_uuid: measurementForm.target === 'occupancy' ? measurementForm.occupancy_uuid : null,
        kind: measurementForm.kind.trim(),
        value,
        unit: normalizeText(measurementForm.unit),
        observed_at: normalizeDateTime(measurementForm.observed_at),
        notes: normalizeText(measurementForm.notes),
      }
      await createMeasurement(payload)
      showNotice('Measurement recorded')
      measurementForm.kind = ''
      measurementForm.value = ''
      measurementForm.unit = ''
      measurementForm.observed_at = ''
      measurementForm.notes = ''
      createMeasurementDialog.value = false
      await loadBatchData(props.batchUuid)
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
    if (!props.batchUuid) {
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
        batch_uuid: string
        occupancy_uuid: null
        kind: string
        value: number
        unit: string | null
        observed_at: string | null
        notes: string | null
      }> = []

      if (temperature.value !== null) {
        payloads.push({
          batch_uuid: props.batchUuid,
          occupancy_uuid: null,
          kind: 'temperature',
          value: temperature.value,
          unit: temperature.unit,
          observed_at: observedAt,
          notes: null,
        })
      }
      if (gravityValue !== null) {
        payloads.push({
          batch_uuid: props.batchUuid,
          occupancy_uuid: null,
          kind: 'gravity',
          value: gravityValue,
          unit: null,
          observed_at: observedAt,
          notes: null,
        })
      }
      if (noteText) {
        payloads.push({
          batch_uuid: props.batchUuid,
          occupancy_uuid: null,
          kind: 'note',
          value: 0,
          unit: null,
          observed_at: observedAt,
          notes: noteText,
        })
      }

      await Promise.all(payloads.map(payload => createMeasurement(payload)))
      showNotice(`Recorded ${payloads.length} timeline ${payloads.length === 1 ? 'entry' : 'entries'}`)
      resetTimelineReading()
      await loadBatchData(props.batchUuid)
    } catch (error) {
      handleError(error)
    }
  }

  async function recordTimelineExtended () {
    if (!props.batchUuid) {
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
        batch_uuid: string
        occupancy_uuid: null
        kind: string
        value: number
        unit: string | null
        observed_at: string | null
        notes: string | null
      }> = []

      if (temperatureValue !== null) {
        payloads.push({
          batch_uuid: props.batchUuid,
          occupancy_uuid: null,
          kind: 'temperature',
          value: temperatureValue,
          unit: normalizeText(timelineExtended.temperature_unit),
          observed_at: observedAt,
          notes: null,
        })
      }
      if (gravityValue !== null) {
        payloads.push({
          batch_uuid: props.batchUuid,
          occupancy_uuid: null,
          kind: 'gravity',
          value: gravityValue,
          unit: normalizeText(timelineExtended.gravity_unit),
          observed_at: observedAt,
          notes: null,
        })
      }
      if (phValue !== null) {
        payloads.push({
          batch_uuid: props.batchUuid,
          occupancy_uuid: null,
          kind: 'ph',
          value: phValue,
          unit: normalizeText(timelineExtended.ph_unit),
          observed_at: observedAt,
          notes: null,
        })
      }
      if (extraKind && extraValue !== null) {
        payloads.push({
          batch_uuid: props.batchUuid,
          occupancy_uuid: null,
          kind: extraKind,
          value: extraValue,
          unit: normalizeText(timelineExtended.extra_unit),
          observed_at: observedAt,
          notes: null,
        })
      }
      if (noteText) {
        payloads.push({
          batch_uuid: props.batchUuid,
          occupancy_uuid: null,
          kind: 'note',
          value: 0,
          unit: null,
          observed_at: observedAt,
          notes: noteText,
        })
      }

      await Promise.all(payloads.map(payload => createMeasurement(payload)))
      showNotice(`Recorded ${payloads.length} timeline ${payloads.length === 1 ? 'entry' : 'entries'}`)
      resetTimelineReading()
      resetTimelineExtended()
      timelineExtendedDialog.value = false
      await loadBatchData(props.batchUuid)
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

  function selectBrewSession (uuid: string) {
    selectedBrewSessionUuid.value = uuid
  }

  function clearBrewSessionSelection () {
    selectedBrewSessionUuid.value = null
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
    if (!session?.wort_volume_uuid || !hotSideAdditionForm.amount) {
      return
    }

    savingHotSideAddition.value = true

    try {
      const payload = {
        volume_uuid: session.wort_volume_uuid,
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
      await loadWortData(session.wort_volume_uuid)
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
    if (!session?.wort_volume_uuid || !hotSideMeasurementForm.kind || !hotSideMeasurementForm.value) {
      return
    }

    savingHotSideMeasurement.value = true

    try {
      const payload = {
        volume_uuid: session.wort_volume_uuid,
        kind: hotSideMeasurementForm.kind,
        value: Number(hotSideMeasurementForm.value),
        unit: normalizeText(hotSideMeasurementForm.unit) ?? getDefaultUnitForKind(hotSideMeasurementForm.kind),
        observed_at: hotSideMeasurementForm.observed_at ? new Date(hotSideMeasurementForm.observed_at).toISOString() : null,
        notes: normalizeText(hotSideMeasurementForm.notes),
      }

      await createMeasurement(payload)
      showNotice('Hot-side measurement recorded')

      hotSideMeasurementDialog.value = false
      await loadWortData(session.wort_volume_uuid)
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

  // ==================== Assign Fermenter Functions ====================

  async function openAssignFermenterDialog () {
    try {
      // Load active occupancies and batch volumes in parallel
      const [occupanciesData, volumesData] = await Promise.all([
        getActiveOccupancies(),
        props.batchUuid
          ? request<ProductionVolume[]>(`/volumes?batch_uuid=${props.batchUuid}`)
          : Promise.resolve([]),
      ])
      activeOccupancies.value = occupanciesData
      batchProductionVolumes.value = volumesData
      assignFermenterDialog.value = true
    } catch (error) {
      handleError(error)
    }
  }

  async function handleFermenterAssigned () {
    showNotice('Batch assigned to fermenter')
    if (props.batchUuid) {
      await loadBatchData(props.batchUuid)
    }
  }

  // ==================== Brew Day Wizard Functions ====================

  async function openBrewDayWizard () {
    try {
      // Ensure we have fresh occupancy data
      const [occupanciesData, volumesData] = await Promise.all([
        getActiveOccupancies(),
        props.batchUuid
          ? request<ProductionVolume[]>(`/volumes?batch_uuid=${props.batchUuid}`)
          : Promise.resolve([]),
      ])
      activeOccupancies.value = occupanciesData
      batchProductionVolumes.value = volumesData
      brewDayWizardDialog.value = true
    } catch (error) {
      handleError(error)
    }
  }

  async function handleBrewDayWizardCompleted () {
    if (props.batchUuid) {
      await loadBatchData(props.batchUuid)
    }
  }

  // ==================== Mark Empty Functions ====================

  async function openMarkEmptyDialog () {
    if (!batchSummary.value?.current_occupancy_uuid) return

    try {
      const occupancy = await getOccupancy(batchSummary.value.current_occupancy_uuid)
      markEmptyOccupancy.value = occupancy
      markEmptyDialog.value = true
    } catch (error) {
      handleError(error)
    }
  }

  async function handleVesselEmptied () {
    showNotice('Vessel marked as empty')
    if (props.batchUuid) {
      await loadBatchData(props.batchUuid)
    }
  }

  // ==================== Transfer Functions ====================

  async function openTransferDialog () {
    if (!batchSummary.value?.current_occupancy_uuid) return

    try {
      const occupancy = await getOccupancy(batchSummary.value.current_occupancy_uuid)
      transferOccupancy.value = occupancy

      // Resolve vessel and volume
      const vessel = vessels.value.find(v => v.uuid === occupancy.vessel_uuid) ?? null
      transferVessel.value = vessel

      if (occupancy.volume_uuid) {
        try {
          transferVolume.value = await getProductionVolumes().then(
            vols => vols.find(v => v.uuid === occupancy.volume_uuid) ?? null,
          )
        } catch {
          transferVolume.value = null
        }
      } else {
        transferVolume.value = null
      }

      transferDialog.value = true
    } catch (error) {
      handleError(error)
    }
  }

  async function handleTransferCompleted () {
    showNotice('Transfer complete')
    if (props.batchUuid) {
      await loadBatchData(props.batchUuid)
    }
  }

  // ==================== Packaging Functions ====================

  async function openPackagingDialog () {
    if (!batchSummary.value?.current_occupancy_uuid) return

    try {
      const occupancy = await getOccupancy(batchSummary.value.current_occupancy_uuid)
      packagingOccupancy.value = occupancy

      // Resolve vessel and volume
      const vessel = vessels.value.find(v => v.uuid === occupancy.vessel_uuid) ?? null
      packagingVessel.value = vessel

      if (occupancy.volume_uuid) {
        try {
          packagingVolume.value = await getProductionVolumes().then(
            vols => vols.find(v => v.uuid === occupancy.volume_uuid) ?? null,
          )
        } catch {
          packagingVolume.value = null
        }
      } else {
        packagingVolume.value = null
      }

      packagingDialog.value = true
    } catch (error) {
      handleError(error)
    }
  }

  async function handlePackagingCompleted () {
    if (props.batchUuid) {
      await loadBatchData(props.batchUuid)
    }
  }

  // ==================== Occupancy Status Functions ====================

  async function changeOccupancyStatus (occupancyUuid: string, status: OccupancyStatus) {
    if (!props.batchUuid) {
      return
    }

    try {
      await updateOccupancyStatus(occupancyUuid, status)
      showNotice(`Status updated to ${formatOccupancyStatus(status)}`)
      // Reload batch summary to reflect the change
      await loadBatchSummary(props.batchUuid)
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
    if (!props.batchUuid || !selectedBatch.value) return

    savingBatch.value = true
    editBatchError.value = ''

    try {
      const payload: UpdateBatchRequest = {
        short_name: form.short_name.trim(),
        brew_date: form.brew_date ? normalizeDateOnly(form.brew_date) : null,
        recipe_uuid: form.recipe_uuid,
        notes: normalizeText(form.notes),
      }

      await updateBatch(props.batchUuid, payload)
      showNotice('Batch updated')
      editBatchDialog.value = false
      await loadBatchData(props.batchUuid)
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
    if (!props.batchUuid || !selectedBatch.value) return

    deletingBatch.value = true
    deleteBatchError.value = ''

    try {
      await deleteBatch(props.batchUuid)
      showNotice('Batch deleted')
      deleteBatchDialog.value = false
      // Navigate back to batch list
      router.push('/batches/all')
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to delete batch'
      deleteBatchError.value = message
    } finally {
      deletingBatch.value = false
    }
  }

  // ==================== Helper Functions ====================

  function handleError (error: unknown) {
    const message = error instanceof Error ? error.message : 'Unexpected error'
    showNotice(message, 'error')
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
.mini-card {
  height: 100%;
}

.batch-tabs :deep(.v-tab) {
  text-transform: none;
  font-weight: 600;
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
