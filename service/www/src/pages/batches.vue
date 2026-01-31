<template>
  <v-container class="production-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="4">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-barley" />
            Batches
            <v-spacer />
            <v-btn
              icon="mdi-plus"
              size="small"
              variant="text"
              aria-label="Create batch"
              @click="createBatchDialog = true"
            />
          </v-card-title>
          <v-card-text>
            <v-alert
              v-if="errorMessage"
              class="mb-3"
              density="compact"
              type="error"
              variant="tonal"
            >
              {{ errorMessage }}
            </v-alert>

            <v-list class="batch-list" lines="two" active-color="primary">
              <v-list-item
                v-for="batch in batches"
                :key="batch.id"
                :active="batch.id === selectedBatchId"
                @click="selectBatch(batch.id)"
              >
                <v-list-item-title>
                  {{ batch.short_name }}
                </v-list-item-title>
                <v-list-item-subtitle>
                  #{{ batch.id }} - {{ formatDate(batch.brew_date) }}
                </v-list-item-subtitle>
                <template #append>
                  <v-chip size="x-small" variant="tonal">
                    {{ formatDateTime(batch.updated_at) }}
                  </v-chip>
                </template>
              </v-list-item>

              <v-list-item v-if="batches.length === 0">
                <v-list-item-title>No batches yet</v-list-item-title>
                <v-list-item-subtitle>Use + to add the first batch.</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-beaker-outline" />
            Batch details
            <v-spacer />
            <v-btn size="small" variant="text" @click="refreshAll">Refresh</v-btn>
            <v-btn size="small" variant="text" @click="clearSelection">Clear</v-btn>
          </v-card-title>
          <v-card-text>
            <v-alert
              v-if="!selectedBatch"
              density="comfortable"
              type="info"
              variant="tonal"
            >
              Select a batch to review timeline, flow, measurements, and additions.
            </v-alert>

            <div v-else>
              <v-row class="mb-4" align="stretch">
                <v-col cols="12" md="6">
                  <v-card class="mini-card" variant="tonal">
                    <v-card-text>
                      <div class="text-overline">Batch</div>
                      <div class="text-h5 font-weight-semibold">
                        {{ selectedBatch.short_name }}
                      </div>
                      <div class="text-body-2 text-medium-emphasis">
                        ID {{ selectedBatch.id }} - {{ selectedBatch.uuid }}
                      </div>
                      <div class="text-body-2 text-medium-emphasis">
                        Brew date {{ formatDate(selectedBatch.brew_date) }}
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>
                <v-col cols="12" md="6">
                  <v-card class="mini-card" variant="tonal">
                    <v-card-text>
                      <div class="text-overline">Status</div>
                      <div class="d-flex flex-wrap ga-2 mb-2">
                        <v-chip v-if="latestProcessPhase" color="primary" size="small" variant="tonal">
                          {{ latestProcessPhase.process_phase }}
                        </v-chip>
                        <v-chip v-if="latestLiquidPhase" color="secondary" size="small" variant="tonal">
                          {{ latestLiquidPhase.liquid_phase }}
                        </v-chip>
                        <v-chip v-if="!latestLiquidPhase && !latestProcessPhase" size="small" variant="outlined">
                          No phases yet
                        </v-chip>
                      </div>
                      <div class="d-flex flex-wrap ga-2 mb-2">
                        <v-chip
                          size="small"
                          color="info"
                          :variant="latestTemperatureMeasurement ? 'tonal' : 'outlined'"
                        >
                          Temp
                          {{
                            latestTemperatureMeasurement
                              ? formatValue(
                                  latestTemperatureMeasurement.value,
                                  latestTemperatureMeasurement.unit,
                                )
                              : 'n/a'
                          }}
                        </v-chip>
                        <v-chip
                          size="small"
                          color="secondary"
                          :variant="latestGravityMeasurement ? 'tonal' : 'outlined'"
                        >
                          Gravity
                          {{
                            latestGravityMeasurement
                              ? formatValue(latestGravityMeasurement.value, latestGravityMeasurement.unit)
                              : 'n/a'
                          }}
                        </v-chip>
                        <v-chip
                          size="small"
                          color="warning"
                          :variant="latestPhMeasurement ? 'tonal' : 'outlined'"
                        >
                          pH
                          {{
                            latestPhMeasurement
                              ? formatValue(latestPhMeasurement.value, latestPhMeasurement.unit)
                              : 'n/a'
                          }}
                        </v-chip>
                      </div>
                      <div class="text-body-2 text-medium-emphasis">
                        Updated {{ formatDateTime(selectedBatch.updated_at) }}
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>

              <v-tabs v-model="activeTab" class="batch-tabs" color="primary" show-arrows>
                <v-tab value="timeline">Timeline</v-tab>
                <v-tab value="flow">Flow</v-tab>
                <v-tab value="measurements">Measurements</v-tab>
                <v-tab value="additions">Additions</v-tab>
              </v-tabs>

              <v-window v-model="activeTab" class="mt-4">

                <v-window-item value="additions">
                  <v-row>
                    <v-col cols="12" md="5">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Record addition
                        </v-card-title>
                        <v-card-text>
                          <v-select
                            v-model="additionForm.target"
                            :items="additionTargetOptions"
                            label="Target"
                          />
                          <v-text-field
                            v-if="additionForm.target === 'occupancy'"
                            v-model="additionForm.occupancy_id"
                            label="Occupancy ID"
                            type="number"
                          />
                          <v-select
                            v-model="additionForm.addition_type"
                            :items="additionTypeOptions"
                            label="Addition type"
                          />
                          <v-text-field v-model="additionForm.stage" label="Stage" />
                          <v-text-field v-model="additionForm.inventory_lot_uuid" label="Inventory lot UUID" />
                          <v-text-field v-model="additionForm.amount" label="Amount" type="number" />
                          <v-select
                            v-model="additionForm.amount_unit"
                            :items="unitOptions"
                            label="Unit"
                          />
                          <v-text-field v-model="additionForm.added_at" label="Added at" type="datetime-local" />
                          <v-textarea
                            v-model="additionForm.notes"
                            auto-grow
                            label="Notes"
                            rows="2"
                          />
                          <v-btn
                            block
                            color="primary"
                            :disabled="
                              !additionForm.amount ||
                              (additionForm.target === 'occupancy' && !additionForm.occupancy_id)
                            "
                            @click="recordAddition"
                          >
                            Add addition
                          </v-btn>
                        </v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" md="7">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1">
                          Addition log
                        </v-card-title>
                        <v-card-text>
                          <v-table class="data-table" density="compact">
                            <thead>
                              <tr>
                                <th>Type</th>
                                <th>Amount</th>
                                <th>Target</th>
                                <th>Time</th>
                              </tr>
                            </thead>
                            <tbody>
                              <tr v-for="addition in additionsSorted" :key="addition.id">
                                <td>{{ addition.addition_type }}</td>
                                <td>{{ formatAmount(addition.amount, addition.amount_unit) }}</td>
                                <td>{{ addition.occupancy_id ?? addition.batch_id }}</td>
                                <td>{{ formatDateTime(addition.added_at) }}</td>
                              </tr>
                              <tr v-if="additionsSorted.length === 0">
                                <td colspan="4">No additions recorded.</td>
                              </tr>
                            </tbody>
                          </v-table>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-window-item>

                <v-window-item value="measurements">
                  <v-row>
                    <v-col cols="12" md="5">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Record measurement
                        </v-card-title>
                        <v-card-text>
                          <v-select
                            v-model="measurementForm.target"
                            :items="additionTargetOptions"
                            label="Target"
                          />
                          <v-text-field
                            v-if="measurementForm.target === 'occupancy'"
                            v-model="measurementForm.occupancy_id"
                            label="Occupancy ID"
                            type="number"
                          />
                          <v-text-field v-model="measurementForm.kind" label="Kind" placeholder="gravity" />
                          <v-text-field v-model="measurementForm.value" label="Value" type="number" />
                          <v-text-field v-model="measurementForm.unit" label="Unit" />
                          <v-text-field v-model="measurementForm.observed_at" label="Observed at" type="datetime-local" />
                          <v-textarea
                            v-model="measurementForm.notes"
                            auto-grow
                            label="Notes"
                            rows="2"
                          />
                          <v-btn
                            block
                            color="secondary"
                            :disabled="
                              !measurementForm.kind.trim() ||
                              !measurementForm.value ||
                              (measurementForm.target === 'occupancy' && !measurementForm.occupancy_id)
                            "
                            @click="recordMeasurement"
                          >
                            Add measurement
                          </v-btn>
                        </v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" md="7">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1">
                          Measurement log
                        </v-card-title>
                        <v-card-text>
                          <v-table class="data-table" density="compact">
                            <thead>
                              <tr>
                                <th>Kind</th>
                                <th>Value</th>
                                <th>Target</th>
                                <th>Time</th>
                              </tr>
                            </thead>
                            <tbody>
                              <tr v-for="measurement in measurementsSorted" :key="measurement.id">
                                <td>{{ measurement.kind }}</td>
                                <td>
                                  {{
                                    isNoteMeasurement(measurement)
                                      ? measurement.notes ?? 'Note'
                                      : formatValue(measurement.value, measurement.unit)
                                  }}
                                </td>
                                <td>{{ measurement.occupancy_id ?? measurement.batch_id }}</td>
                                <td>{{ formatDateTime(measurement.observed_at) }}</td>
                              </tr>
                              <tr v-if="measurementsSorted.length === 0">
                                <td colspan="4">No measurements recorded.</td>
                              </tr>
                            </tbody>
                          </v-table>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-window-item>

                <v-window-item value="timeline">
                  <v-card class="sub-card" variant="outlined">
                    <v-card-text>
                      <v-card class="sub-card mb-4" variant="tonal">
                        <v-card-text>
                          <div class="d-flex flex-wrap align-center justify-space-between ga-2">
                            <div class="text-subtitle-2 font-weight-semibold">Quick Update</div>
                            <v-btn
                              size="small"
                              variant="text"
                              append-icon="mdi-arrow-right"
                              @click="openTimelineExtendedDialog"
                            >
                              More
                            </v-btn>
                          </div>
                          <v-divider class="my-3" />

                          <v-row dense align="center">
                            <v-col cols="12" md="1" class="d-flex align-center justify-center">
                              <v-menu
                                v-model="timelineObservedAtMenu"
                                :close-on-content-click="false"
                                location="bottom"
                              >
                                <template #activator="{ props }">
                                  <v-tooltip location="top">
                                    <template #activator="{ props: tooltipProps }">
                                      <v-btn
                                        v-bind="{ ...props, ...tooltipProps }"
                                        icon="mdi-clock-outline"
                                        size="default"
                                        variant="text"
                                        :color="timelineReading.observed_at ? 'secondary' : 'primary'"
                                        aria-label="Set observation time"
                                      />
                                    </template>
                                    <span>{{ timelineObservedAtLabel }}</span>
                                  </v-tooltip>
                                </template>
                                <v-card>
                                  <v-card-text>
                                    <v-text-field
                                      v-model="timelineReading.observed_at"
                                      density="compact"
                                      label="Observed at"
                                      type="datetime-local"
                                    />
                                  </v-card-text>
                                  <v-card-actions class="justify-end">
                                    <v-btn variant="text" @click="clearTimelineObservedAt">
                                      Use now
                                    </v-btn>
                                    <v-btn variant="text" @click="timelineObservedAtMenu = false">Done</v-btn>
                                  </v-card-actions>
                                </v-card>
                              </v-menu>
                            </v-col>
                            <v-col cols="12" md="2">
                              <v-text-field
                                v-model="timelineReading.temperature"
                                density="compact"
                                label="Temp"
                                placeholder="67F"
                                inputmode="decimal"
                              />
                            </v-col>
                            <v-col cols="12" md="2">
                              <v-text-field
                                v-model="timelineReading.gravity"
                                density="compact"
                                label="Gravity"
                                placeholder="1.056"
                                inputmode="decimal"
                              />
                            </v-col>
                            <v-col cols="12" md="5">
                              <v-text-field
                                v-model="timelineReading.notes"
                                density="compact"
                                label="Notes"
                                placeholder="Aroma, flavor, observations"
                              />
                            </v-col>
                            <v-col cols="12" md="1" class="d-flex align-center justify-end">
                              <v-btn
                                color="primary"
                                :disabled="!timelineReadingReady"
                                @click="recordTimelineReading"
                              >
                                Add
                              </v-btn>
                            </v-col>
                          </v-row>
                        </v-card-text>
                      </v-card>

                      <v-timeline align="start" density="compact" side="end">
                        <v-timeline-item
                          v-for="event in timelineItems"
                          :key="event.id"
                          :dot-color="event.color"
                          :icon="event.icon"
                        >
                          <template #opposite>
                            <div class="text-caption text-medium-emphasis">
                              {{ formatDateTime(event.at) }}
                            </div>
                          </template>
                          <div class="text-subtitle-2 font-weight-semibold">
                            {{ event.title }}
                          </div>
                          <div class="text-body-2 text-medium-emphasis">
                            {{ event.subtitle }}
                          </div>
                        </v-timeline-item>

                        <v-timeline-item v-if="timelineItems.length === 0" dot-color="grey">
                          <div class="text-body-2 text-medium-emphasis">
                            No timeline events yet.
                          </div>
                        </v-timeline-item>
                      </v-timeline>
                    </v-card-text>
                  </v-card>
                </v-window-item>

                <v-window-item value="flow">
                  <v-card class="sub-card" variant="outlined">
                    <v-card-title class="text-subtitle-1">Liquid flow</v-card-title>
                    <v-card-text>
                      <v-alert
                        v-if="flowNotice"
                        class="mb-3"
                        density="compact"
                        type="info"
                        variant="tonal"
                      >
                        {{ flowNotice }}
                      </v-alert>

                      <SankeyDiagram v-if="flowLinks.length" :nodes="flowNodes" :links="flowLinks" />
                      <div v-else class="text-body-2 text-medium-emphasis">
                        No flow relations yet. Record a split or blend to visualize liquid movement.
                      </div>

                      <div class="text-caption text-medium-emphasis mt-3">
                        Flow is derived from volume relations (splits and blends).
                      </div>
                    </v-card-text>
                  </v-card>
                </v-window-item>
              </v-window>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <v-dialog v-model="createBatchDialog" max-width="520">
    <v-card>
      <v-card-title class="text-h6">Create batch</v-card-title>
      <v-card-text>
        <v-text-field
          v-model="newBatch.short_name"
          density="comfortable"
          label="Short name"
          placeholder="IPA 24-07"
        />
        <v-text-field
          v-model="newBatch.brew_date"
          density="comfortable"
          label="Brew date"
          type="date"
        />
        <v-textarea
          v-model="newBatch.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="createBatchDialog = false">Cancel</v-btn>
        <v-btn color="primary" :disabled="!newBatch.short_name.trim()" @click="createBatch">
          Create batch
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="timelineExtendedDialog" max-width="720">
    <v-card>
      <v-card-title class="text-h6">Extended reading</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.observed_at"
              density="comfortable"
              label="Observed at"
              type="datetime-local"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.temperature"
              density="comfortable"
              label="Temperature"
              type="number"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.temperature_unit"
              density="comfortable"
              label="Temp unit"
              placeholder="F"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.gravity"
              density="comfortable"
              label="Gravity"
              type="number"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.gravity_unit"
              density="comfortable"
              label="Gravity unit"
              placeholder="SG"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.ph"
              density="comfortable"
              label="pH"
              type="number"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.ph_unit"
              density="comfortable"
              label="pH unit"
              placeholder="pH"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.extra_kind"
              density="comfortable"
              label="Other kind"
              placeholder="CO2"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.extra_value"
              density="comfortable"
              label="Other value"
              type="number"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="timelineExtended.extra_unit"
              density="comfortable"
              label="Other unit"
            />
          </v-col>
          <v-col cols="12" md="8">
            <v-text-field
              v-model="timelineExtended.notes"
              density="comfortable"
              label="Notes"
              placeholder="Aroma, flavor, observations"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="timelineExtendedDialog = false">Cancel</v-btn>
        <v-btn color="primary" :disabled="!timelineExtendedReady" @click="recordTimelineExtended">
          Record
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, computed, watch } from 'vue'

type Unit = 'ml' | 'usfloz' | 'ukfloz'
type LiquidPhase = 'water' | 'wort' | 'beer'
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
type RelationType = 'split' | 'blend'
type AdditionType =
  | 'malt'
  | 'hop'
  | 'yeast'
  | 'adjunct'
  | 'water_chem'
  | 'gas'
  | 'other'

type Batch = {
  id: number
  uuid: string
  short_name: string
  brew_date: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

type Volume = {
  id: number
  uuid: string
  name: string | null
  description: string | null
  amount: number
  amount_unit: Unit
  created_at: string
  updated_at: string
}

type VolumeRelation = {
  id: number
  uuid: string
  parent_volume_id: number
  child_volume_id: number
  relation_type: RelationType
  amount: number
  amount_unit: Unit
  created_at: string
  updated_at: string
}

type BatchVolume = {
  id: number
  uuid: string
  batch_id: number
  volume_id: number
  liquid_phase: LiquidPhase
  phase_at: string
  created_at: string
  updated_at: string
}

type BatchProcessPhase = {
  id: number
  uuid: string
  batch_id: number
  process_phase: ProcessPhase
  phase_at: string
  created_at: string
  updated_at: string
}

type Addition = {
  id: number
  uuid: string
  batch_id: number | null
  occupancy_id: number | null
  addition_type: AdditionType
  stage: string | null
  inventory_lot_uuid: string | null
  amount: number
  amount_unit: Unit
  added_at: string
  notes: string | null
  created_at: string
  updated_at: string
}

type Measurement = {
  id: number
  uuid: string
  batch_id: number | null
  occupancy_id: number | null
  kind: string
  value: number
  unit: string | null
  observed_at: string
  notes: string | null
  created_at: string
  updated_at: string
}

type TimelineEvent = {
  id: string
  title: string
  subtitle: string
  at: string
  color: string
  icon: string
}

type FlowNode = {
  id: string
  label: string
}

type FlowLink = {
  source: string
  target: string
  value: number
  label: string
}

const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

const unitOptions: Unit[] = ['ml', 'usfloz', 'ukfloz']
const additionTypeOptions: AdditionType[] = [
  'malt',
  'hop',
  'yeast',
  'adjunct',
  'water_chem',
  'gas',
  'other',
]
const additionTargetOptions = [
  { title: 'Batch', value: 'batch' },
  { title: 'Occupancy', value: 'occupancy' },
]

const batches = ref<Batch[]>([])
const volumes = ref<Volume[]>([])
const batchVolumes = ref<BatchVolume[]>([])
const processPhases = ref<BatchProcessPhase[]>([])
const additions = ref<Addition[]>([])
const measurements = ref<Measurement[]>([])
const volumeRelations = ref<VolumeRelation[]>([])

const selectedBatchId = ref<number | null>(null)
const activeTab = ref('timeline')
const errorMessage = ref('')
const createBatchDialog = ref(false)

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

const newBatch = reactive({
  short_name: '',
  brew_date: '',
  notes: '',
})

const additionForm = reactive({
  target: 'batch',
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
  target: 'batch',
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


const selectedBatch = computed(() =>
  batches.value.find((batch) => batch.id === selectedBatchId.value) ?? null,
)

const latestProcessPhase = computed(() => getLatest(processPhases.value, (item) => item.phase_at))
const latestLiquidPhase = computed(() => getLatest(batchVolumes.value, (item) => item.phase_at))

const latestTemperatureMeasurement = computed(() =>
  getLatestMeasurement(['temperature', 'temp']),
)
const latestGravityMeasurement = computed(() =>
  getLatestMeasurement(['gravity', 'grav', 'sg']),
)
const latestPhMeasurement = computed(() => getLatestMeasurement(['ph']))

const timelineObservedAtMenu = ref(false)
const timelineObservedAtLabel = computed(() =>
  timelineReading.observed_at ? formatDateTime(timelineReading.observed_at) : 'Now',
)

const timelineReadingReady = computed(() => {
  const hasTemperature = parseTemperatureInput(timelineReading.temperature).value !== null
  const hasGravity = parseNumericInput(timelineReading.gravity) !== null
  const hasNotes = timelineReading.notes.trim().length > 0
  return hasTemperature || hasGravity || hasNotes
})

const timelineExtendedReady = computed(() => {
  const hasTemperature = toNumber(timelineExtended.temperature) !== null
  const hasGravity = toNumber(timelineExtended.gravity) !== null
  const hasPh = toNumber(timelineExtended.ph) !== null
  const hasNotes = timelineExtended.notes.trim().length > 0
  const hasExtraKind = timelineExtended.extra_kind.trim().length > 0
  const extraValue = toNumber(timelineExtended.extra_value)
  if (hasExtraKind && extraValue === null) {
    return false
  }
  const hasExtra = hasExtraKind && extraValue !== null
  return hasTemperature || hasGravity || hasPh || hasExtra || hasNotes
})
const additionsSorted = computed(() =>
  sortByTime(additions.value, (item) => item.added_at),
)
const measurementsSorted = computed(() =>
  sortByTime(measurements.value, (item) => item.observed_at),
)

const volumeNameMap = computed(
  () =>
    new Map(
      volumes.value.map((volume) => [volume.id, volume.name ?? `Volume ${volume.id}`]),
    ),
)

const flowUnit = computed<Unit | null>(() => {
  const counts = new Map<Unit, number>()
  volumeRelations.value.forEach((relation) => {
    if (!relation.amount || relation.amount <= 0) {
      return
    }
    counts.set(relation.amount_unit, (counts.get(relation.amount_unit) ?? 0) + 1)
  })
  let selected: { unit: Unit; count: number } | null = null
  counts.forEach((count, unit) => {
    if (!selected || count > selected.count) {
      selected = { unit, count }
    }
  })
  return selected?.unit ?? null
})

const flowRelations = computed(() => {
  const relations = volumeRelations.value.filter((relation) => relation.amount > 0)
  if (!flowUnit.value) {
    return relations
  }
  return relations.filter((relation) => relation.amount_unit === flowUnit.value)
})

const flowNotice = computed(() => {
  if (!flowUnit.value) {
    return ''
  }
  const total = volumeRelations.value.filter((relation) => relation.amount > 0).length
  const shown = flowRelations.value.length
  if (total > shown) {
    return `Showing ${shown} of ${total} relations measured in ${flowUnit.value} for consistent weights.`
  }
  return ''
})

const flowNodes = computed<FlowNode[]>(() => {
  const nodes = new Map<string, FlowNode>()
  const labelFor = (volumeId: number) => volumeNameMap.value.get(volumeId) ?? `Volume ${volumeId}`

  flowRelations.value.forEach((relation) => {
    const parentId = `volume-${relation.parent_volume_id}`
    const childId = `volume-${relation.child_volume_id}`
    if (!nodes.has(parentId)) {
      nodes.set(parentId, { id: parentId, label: labelFor(relation.parent_volume_id) })
    }
    if (!nodes.has(childId)) {
      nodes.set(childId, { id: childId, label: labelFor(relation.child_volume_id) })
    }
  })

  return Array.from(nodes.values())
})

const flowLinks = computed<FlowLink[]>(() => {
  const links = new Map<string, FlowLink>()

  flowRelations.value.forEach((relation) => {
    const source = `volume-${relation.parent_volume_id}`
    const target = `volume-${relation.child_volume_id}`
    const key = `${source}-${target}-${relation.amount_unit ?? ''}`
    const existing = links.get(key)
    if (existing) {
      existing.value += relation.amount
      existing.label = formatAmount(existing.value, relation.amount_unit)
      return
    }
    links.set(key, {
      source,
      target,
      value: relation.amount,
      label: formatAmount(relation.amount, relation.amount_unit),
    })
  })

  return Array.from(links.values())
})

const timelineItems = computed(() => {
  const items: TimelineEvent[] = []

  additions.value.forEach((addition) => {
    items.push({
      id: `addition-${addition.id}`,
      title: `Addition: ${addition.addition_type}`,
      subtitle: `${formatAmount(addition.amount, addition.amount_unit)} ${addition.stage ?? ''}`.trim(),
      at: addition.added_at ?? addition.created_at,
      color: 'primary',
      icon: 'mdi-flask-outline',
    })
  })

  const groupedMeasurements = groupMeasurements(measurements.value)
  groupedMeasurements.forEach((group, index) => {
    const subtitle = orderMeasurementGroup(group)
      .map((measurement) => formatMeasurementEntry(measurement))
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
  })

  processPhases.value.forEach((phase) => {
    items.push({
      id: `process-${phase.id}`,
      title: `Process phase: ${phase.process_phase}`,
      subtitle: `Batch ${phase.batch_id}`,
      at: phase.phase_at ?? phase.created_at,
      color: 'success',
      icon: 'mdi-progress-check',
    })
  })

  batchVolumes.value.forEach((phase) => {
    items.push({
      id: `liquid-${phase.id}`,
      title: `Liquid phase: ${phase.liquid_phase}`,
      subtitle: `Volume ${phase.volume_id}`,
      at: phase.phase_at ?? phase.created_at,
      color: 'warning',
      icon: 'mdi-water',
    })
  })

  return items.sort((a, b) => toTimestamp(b.at) - toTimestamp(a.at))
})

watch(selectedBatchId, (value) => {
  if (value) {
    loadBatchData(value)
  }
})

watch(timelineObservedAtMenu, (isOpen) => {
  if (isOpen && !timelineReading.observed_at) {
    timelineReading.observed_at = nowInputValue()
  }
})

onMounted(async () => {
  await refreshAll()
})

function selectBatch(id: number) {
  selectedBatchId.value = id
}

function clearSelection() {
  selectedBatchId.value = null
}

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`${apiBase}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init.headers ?? {}),
    },
  })

  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || `Request failed with ${response.status}`)
  }

  if (response.status === 204) {
    return null as T
  }

  const contentType = response.headers.get('content-type') ?? ''
  if (contentType.includes('application/json')) {
    return response.json() as Promise<T>
  }

  return (await response.text()) as T
}

const get = <T>(path: string) => request<T>(path)
const post = <T>(path: string, payload: unknown) =>
  request<T>(path, { method: 'POST', body: JSON.stringify(payload) })

async function refreshAll() {
  errorMessage.value = ''
  try {
    await Promise.all([loadBatches(), loadVolumes()])
    if (!selectedBatchId.value && batches.value.length > 0) {
      selectedBatchId.value = batches.value[0].id
    } else if (selectedBatchId.value) {
      await loadBatchData(selectedBatchId.value)
    }
  } catch (error) {
    handleError(error)
  }
}

async function loadBatches() {
  batches.value = await get<Batch[]>('/batches')
}

async function loadVolumes() {
  volumes.value = await get<Volume[]>('/volumes')
}

async function loadBatchData(batchId: number) {
  try {
    const [batchVolumesData, processPhasesData, additionsData, measurementsData] = await Promise.all([
      get<BatchVolume[]>(`/batch-volumes?batch_id=${batchId}`),
      get<BatchProcessPhase[]>(`/batch-process-phases?batch_id=${batchId}`),
      get<Addition[]>(`/additions?batch_id=${batchId}`),
      get<Measurement[]>(`/measurements?batch_id=${batchId}`),
    ])

    batchVolumes.value = batchVolumesData
    processPhases.value = processPhasesData
    additions.value = additionsData
    measurements.value = measurementsData

    await loadVolumeRelations(batchVolumesData)
  } catch (error) {
    handleError(error)
  }
}

async function loadVolumeRelations(batchVolumeData: BatchVolume[]) {
  const volumeIds = Array.from(
    new Set(batchVolumeData.map((item) => item.volume_id)),
  )
  if (volumeIds.length === 0) {
    volumeRelations.value = []
    return
  }

  const results = await Promise.allSettled(
    volumeIds.map((id) => get<VolumeRelation[]>(`/volume-relations?volume_id=${id}`)),
  )

  volumeRelations.value = results.flatMap((result) =>
    result.status === 'fulfilled' ? result.value : [],
  )
}

async function createBatch() {
  errorMessage.value = ''
  try {
    const payload = {
      short_name: newBatch.short_name.trim(),
      brew_date: normalizeDateOnly(newBatch.brew_date),
      notes: normalizeText(newBatch.notes),
    }
    const created = await post<Batch>('/batches', payload)
    showNotice('Batch created')
    newBatch.short_name = ''
    newBatch.brew_date = ''
    newBatch.notes = ''
    await loadBatches()
    selectedBatchId.value = created.id
    createBatchDialog.value = false
  } catch (error) {
    handleError(error)
  }
}

async function recordAddition() {
  if (!selectedBatchId.value) {
    return
  }
  if (!additionForm.amount) {
    return
  }
  if (additionForm.target === 'occupancy' && !additionForm.occupancy_id) {
    return
  }
  errorMessage.value = ''
  try {
    const payload = {
      batch_id: additionForm.target === 'batch' ? selectedBatchId.value : null,
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
    await loadBatchData(selectedBatchId.value)
  } catch (error) {
    handleError(error)
  }
}

async function recordMeasurement() {
  if (!selectedBatchId.value) {
    return
  }
  if (!measurementForm.kind.trim() || !measurementForm.value) {
    return
  }
  if (measurementForm.target === 'occupancy' && !measurementForm.occupancy_id) {
    return
  }
  errorMessage.value = ''
  try {
    const payload = {
      batch_id: measurementForm.target === 'batch' ? selectedBatchId.value : null,
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
    await loadBatchData(selectedBatchId.value)
  } catch (error) {
    handleError(error)
  }
}

function openTimelineExtendedDialog() {
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

async function recordTimelineReading() {
  if (!selectedBatchId.value) {
    return
  }

  const temperature = parseTemperatureInput(timelineReading.temperature)
  const gravityValue = parseNumericInput(timelineReading.gravity)
  const noteText = normalizeText(timelineReading.notes)

  if (temperature.value === null && gravityValue === null && !noteText) {
    return
  }

  errorMessage.value = ''
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
        batch_id: selectedBatchId.value,
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
        batch_id: selectedBatchId.value,
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
        batch_id: selectedBatchId.value,
        occupancy_id: null,
        kind: 'note',
        value: 0,
        unit: null,
        observed_at: observedAt,
        notes: noteText,
      })
    }

    await Promise.all(payloads.map((payload) => post<Measurement>('/measurements', payload)))
    showNotice(`Recorded ${payloads.length} timeline ${payloads.length === 1 ? 'entry' : 'entries'}`)
    resetTimelineReading()
    await loadBatchData(selectedBatchId.value)
  } catch (error) {
    handleError(error)
  }
}

async function recordTimelineExtended() {
  if (!selectedBatchId.value) {
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
    temperatureValue === null &&
    gravityValue === null &&
    phValue === null &&
    !extraKind &&
    !noteText
  ) {
    return
  }

  errorMessage.value = ''
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
        batch_id: selectedBatchId.value,
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
        batch_id: selectedBatchId.value,
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
        batch_id: selectedBatchId.value,
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
        batch_id: selectedBatchId.value,
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
        batch_id: selectedBatchId.value,
        occupancy_id: null,
        kind: 'note',
        value: 0,
        unit: null,
        observed_at: observedAt,
        notes: noteText,
      })
    }

    await Promise.all(payloads.map((payload) => post<Measurement>('/measurements', payload)))
    showNotice(`Recorded ${payloads.length} timeline ${payloads.length === 1 ? 'entry' : 'entries'}`)
    resetTimelineReading()
    resetTimelineExtended()
    timelineExtendedDialog.value = false
    await loadBatchData(selectedBatchId.value)
  } catch (error) {
    handleError(error)
  }
}

function resetTimelineReading() {
  timelineReading.observed_at = ''
  timelineReading.temperature = ''
  timelineReading.gravity = ''
  timelineReading.notes = ''
}

function resetTimelineExtended() {
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

function clearTimelineObservedAt() {
  timelineReading.observed_at = ''
  timelineObservedAtMenu.value = false
}

function handleError(error: unknown) {
  const message = error instanceof Error ? error.message : 'Unexpected error'
  errorMessage.value = message
  showNotice(message, 'error')
}

function normalizeText(value: string) {
  const trimmed = value.trim()
  return trimmed.length > 0 ? trimmed : null
}

function groupMeasurements(items: Measurement[]) {
  const sorted = [...items].sort(
    (a, b) => measurementTimestamp(b) - measurementTimestamp(a),
  )
  const groups: Measurement[][] = []
  const thresholdMs = 2 * 60 * 1000
  let current: Measurement[] = []
  let anchor: number | null = null

  sorted.forEach((measurement) => {
    const timestamp = measurementTimestamp(measurement)
    if (!current.length) {
      current = [measurement]
      anchor = timestamp
      return
    }
    if (anchor !== null && Math.abs(anchor - timestamp) <= thresholdMs) {
      current.push(measurement)
      return
    }
    groups.push(current)
    current = [measurement]
    anchor = timestamp
  })

  if (current.length) {
    groups.push(current)
  }

  return groups
}

function orderMeasurementGroup(group: Measurement[]) {
  return [...group].sort((a, b) => measurementPriority(a) - measurementPriority(b))
}

function measurementPriority(measurement: Measurement) {
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

function measurementTimestamp(measurement: Measurement) {
  return toTimestamp(measurement.observed_at ?? measurement.created_at)
}

function formatMeasurementEntry(measurement: Measurement) {
  if (isNoteMeasurement(measurement)) {
    return measurement.notes ? `Note: ${measurement.notes}` : 'Note'
  }
  const label = formatMeasurementKind(measurement.kind)
  return `${label} ${formatValue(measurement.value, measurement.unit)}`
}

function formatMeasurementKind(kind: string) {
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

function nowInputValue() {
  const now = new Date()
  const pad = (value: number) => String(value).padStart(2, '0')
  const year = now.getFullYear()
  const month = pad(now.getMonth() + 1)
  const day = pad(now.getDate())
  const hours = pad(now.getHours())
  const minutes = pad(now.getMinutes())
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

function parseNumericInput(value: string) {
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

function parseTemperatureInput(value: string) {
  const parsedValue = parseNumericInput(value)
  if (parsedValue === null) {
    return { value: null, unit: null }
  }
  const unitMatch = value.match(/([cf])\s*$/i)
  if (unitMatch) {
    return { value: parsedValue, unit: unitMatch[1].toUpperCase() }
  }
  return { value: parsedValue, unit: null }
}

function normalizeDateOnly(value: string) {
  return value ? new Date(`${value}T00:00:00Z`).toISOString() : null
}

function normalizeDateTime(value: string) {
  return value ? new Date(value).toISOString() : null
}

function toNumber(value: string | number | null) {
  if (value === null || value === undefined || value === '') {
    return null
  }
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : null
}

function formatDate(value: string | null | undefined) {
  if (!value) {
    return 'Unknown'
  }
  return new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
  }).format(new Date(value))
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

function formatAmount(amount: number | null, unit: string | null | undefined) {
  if (amount === null || amount === undefined) {
    return 'Unknown'
  }
  return `${amount} ${unit ?? ''}`.trim()
}

function formatValue(value: number | null, unit: string | null | undefined) {
  if (value === null || value === undefined) {
    return 'Unknown'
  }
  return `${value}${unit ? ` ${unit}` : ''}`
}

function toTimestamp(value: string | null | undefined) {
  if (!value) {
    return 0
  }
  return new Date(value).getTime()
}

function sortByTime<T>(items: T[], selector: (item: T) => string | null | undefined) {
  return [...items].sort((a, b) => toTimestamp(selector(b)) - toTimestamp(selector(a)))
}

function getLatest<T>(items: T[], selector: (item: T) => string | null | undefined) {
  const sorted = sortByTime(items, selector)
  return sorted.length > 0 ? sorted[0] : null
}

function getLatestMeasurement(kinds: string[]) {
  const normalizedKinds = kinds.map((kind) => normalizeMeasurementKind(kind))
  const filtered = measurements.value.filter((measurement) =>
    matchesMeasurementKind(measurement.kind, normalizedKinds),
  )
  return getLatest(filtered, (item) => item.observed_at ?? item.created_at)
}

function matchesMeasurementKind(value: string, kinds: string[]) {
  const normalized = normalizeMeasurementKind(value)
  if (!normalized) {
    return false
  }
  return kinds.some((kind) => normalized.includes(kind))
}

function normalizeMeasurementKind(value: string) {
  return value.trim().toLowerCase().replace(/[^a-z0-9]/g, '')
}

function isNoteMeasurement(measurement: Measurement) {
  return matchesMeasurementKind(measurement.kind, ['note'])
}
</script>

<style scoped>
.production-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.mini-card {
  height: 100%;
}

.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.batch-list {
  max-height: 320px;
  overflow: auto;
}

.batch-tabs :deep(.v-tab) {
  text-transform: none;
  font-weight: 600;
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
