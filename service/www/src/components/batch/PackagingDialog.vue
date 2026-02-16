<template>
  <v-dialog
    :fullscreen="xs"
    :max-width="xs ? '100%' : 640"
    :model-value="modelValue"
    persistent
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon class="mr-2" color="teal" icon="mdi-package-variant" />
        <span class="text-h6">Package Beer</span>
        <v-spacer />
        <v-btn
          aria-label="Close"
          :disabled="saving"
          icon="mdi-close"
          size="small"
          variant="text"
          @click="handleClose"
        />
      </v-card-title>

      <v-divider />

      <v-card-text class="pa-0" style="overflow-y: auto;">
        <!-- Loading state -->
        <template v-if="loadingData">
          <v-container class="pa-6">
            <div class="d-flex flex-column align-center">
              <v-progress-circular color="teal" indeterminate size="48" />
              <p class="text-body-2 text-medium-emphasis mt-4">Loading packaging data...</p>
            </div>
          </v-container>
        </template>

        <!-- Load error state -->
        <template v-else-if="loadError">
          <v-container class="pa-6">
            <v-alert density="compact" type="error" variant="tonal">
              {{ loadError }}
            </v-alert>
          </v-container>
        </template>

        <!-- Stepper -->
        <v-stepper
          v-if="!loadingData && !loadError"
          v-model="currentStep"
          alt-labels
          flat
          hide-actions
          :items="stepItems"
        >
          <!-- Step 1: Package Details -->
          <template #item.1>
            <v-container class="pa-4">
              <!-- Source vessel (read-only) -->
              <div class="text-overline text-medium-emphasis mb-2">Source</div>
              <v-card class="mb-4" density="compact" variant="outlined">
                <v-card-text class="pa-3">
                  <div class="d-flex align-center flex-wrap ga-2">
                    <v-icon icon="mdi-flask-round-bottom" size="small" />
                    <span class="text-body-1 font-weight-medium">
                      {{ sourceVesselName }}
                    </span>
                    <span class="text-body-2 text-medium-emphasis">
                      &middot; {{ sourceBatchName }}
                    </span>
                  </div>
                  <div class="d-flex align-center flex-wrap ga-2 mt-1">
                    <v-chip
                      v-if="sourceOccupancy?.status"
                      :color="getOccupancyStatusColor(sourceOccupancy.status)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatOccupancyStatus(sourceOccupancy.status) }}
                    </v-chip>
                    <span v-if="sourceVolume" class="text-body-2">
                      {{ formatVolumePreferred(sourceVolume.amount, sourceVolume.amount_unit) }}
                    </span>
                  </div>
                </v-card-text>
              </v-card>

              <v-divider class="my-3" />

              <!-- Packaging date/time -->
              <div class="text-overline text-medium-emphasis mb-2">When</div>
              <v-text-field
                v-model="form.startedAt"
                density="comfortable"
                label="Packaging date/time"
                type="datetime-local"
              />

              <v-divider class="my-3" />

              <!-- Stock location -->
              <div class="text-overline text-medium-emphasis mb-2">Finished Goods Destination</div>
              <v-select
                v-model="form.stockLocationUuid"
                clearable
                density="comfortable"
                hint="Where packaged beer will be stored in inventory"
                item-title="title"
                item-value="value"
                :items="stockLocationOptions"
                label="Stock location"
                persistent-hint
              >
                <template #no-data>
                  <v-list-item>
                    <v-list-item-title>No stock locations</v-list-item-title>
                    <v-list-item-subtitle>Create stock locations in Inventory settings</v-list-item-subtitle>
                  </v-list-item>
                </template>
              </v-select>

              <!-- Lot code prefix -->
              <v-text-field
                v-model="form.lotCodePrefix"
                class="mt-2"
                density="comfortable"
                hint="Used to generate lot codes (e.g. IPA24-07-KEG-01)"
                label="Lot code prefix"
                persistent-hint
              />

              <v-divider class="my-3" />

              <!-- Close source checkbox -->
              <v-checkbox
                v-model="form.closeSource"
                density="compact"
                label="Close source vessel after packaging"
              >
                <template #label>
                  <div>
                    <span>Close source vessel after packaging</span>
                    <div class="text-caption text-medium-emphasis">
                      Remaining volume will be treated as loss
                    </div>
                  </div>
                </template>
              </v-checkbox>
            </v-container>
          </template>

          <!-- Step 2: Package Lines -->
          <template #item.2>
            <v-container class="pa-4">
              <div class="text-overline text-medium-emphasis mb-2">Package Lines</div>

              <div
                v-for="(line, index) in form.lines"
                :key="index"
                class="mb-3"
              >
                <div class="d-flex align-center ga-2 mb-1">
                  <span class="text-caption text-medium-emphasis font-weight-medium">
                    Line {{ index + 1 }}
                  </span>
                  <v-spacer />
                  <v-btn
                    v-if="form.lines.length > 1"
                    aria-label="Remove line"
                    density="compact"
                    icon="mdi-close"
                    size="x-small"
                    variant="text"
                    @click="removeLine(index)"
                  />
                </div>

                <v-select
                  v-model="line.packageFormatUuid"
                  density="comfortable"
                  item-title="title"
                  item-value="value"
                  :items="activeFormatOptions"
                  label="Package format"
                  :rules="[rules.required]"
                >
                  <template #item="{ props: itemProps, item }">
                    <v-list-item v-bind="itemProps">
                      <template #subtitle>
                        <span>{{ item.raw.subtitle }}</span>
                      </template>
                    </v-list-item>
                  </template>
                  <template #no-data>
                    <v-list-item>
                      <v-list-item-title>No package formats</v-list-item-title>
                      <v-list-item-subtitle>Create formats in Production settings</v-list-item-subtitle>
                    </v-list-item>
                  </template>
                </v-select>

                <v-text-field
                  v-model="line.quantity"
                  density="comfortable"
                  inputmode="numeric"
                  label="Quantity"
                  min="1"
                  :rules="line.quantity ? [rules.positiveInteger] : [rules.required]"
                  type="number"
                />

                <!-- Calculated volume for this line -->
                <div v-if="lineVolume(index)" class="text-body-2 text-medium-emphasis mb-1">
                  {{ lineVolume(index) }}
                </div>

                <v-divider v-if="index < form.lines.length - 1" class="mt-1" />
              </div>

              <v-btn
                class="mb-4"
                density="comfortable"
                prepend-icon="mdi-plus"
                variant="text"
                @click="addLine"
              >
                Add Format
              </v-btn>

              <v-divider class="my-3" />

              <!-- Running totals -->
              <div class="text-overline text-medium-emphasis mb-2">Volume Summary</div>

              <v-row dense>
                <v-col cols="6">
                  <div class="metric-card">
                    <div class="metric-label">Total Packaged</div>
                    <div class="metric-value">{{ totalPackagedLabel }}</div>
                  </div>
                </v-col>
                <v-col cols="6">
                  <div class="metric-card" :class="{ 'metric-card--warning': estimatedLossPercent > 10 }">
                    <div class="metric-label">Est. Loss</div>
                    <div class="metric-value">{{ estimatedLossLabel }}</div>
                  </div>
                </v-col>
              </v-row>

              <v-divider class="my-3" />

              <!-- Optional loss override -->
              <v-expansion-panels variant="accordion">
                <v-expansion-panel>
                  <v-expansion-panel-title class="text-body-2">
                    Override loss amount
                  </v-expansion-panel-title>
                  <v-expansion-panel-text>
                    <v-row dense>
                      <v-col cols="7">
                        <v-text-field
                          v-model="form.lossAmount"
                          density="comfortable"
                          hint="Explicit loss to record"
                          inputmode="decimal"
                          label="Loss amount"
                          persistent-hint
                          :rules="form.lossAmount ? [rules.positiveNumber] : []"
                        />
                      </v-col>
                      <v-col cols="5">
                        <v-select
                          v-model="form.lossUnit"
                          density="comfortable"
                          item-title="label"
                          item-value="value"
                          :items="volumeUnitOptions"
                          label="Unit"
                        />
                      </v-col>
                    </v-row>
                  </v-expansion-panel-text>
                </v-expansion-panel>
              </v-expansion-panels>

              <v-divider class="my-3" />

              <!-- Notes -->
              <v-textarea
                v-model="form.notes"
                auto-grow
                density="comfortable"
                label="Notes"
                rows="2"
              />
            </v-container>
          </template>

          <!-- Step 3: Review & Confirm -->
          <template #item.3>
            <v-container class="pa-4">
              <div class="text-overline text-medium-emphasis mb-3">Review Packaging Run</div>

              <!-- Source info -->
              <div class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Source</div>
                <div class="text-body-1 font-weight-medium">
                  {{ sourceVesselName }}
                  <span class="text-medium-emphasis font-weight-regular">
                    ({{ sourceBatchName }})
                  </span>
                </div>
              </div>

              <!-- Date -->
              <div class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Date</div>
                <div class="text-body-2">{{ reviewDateLabel }}</div>
              </div>

              <!-- Stock location -->
              <div v-if="form.stockLocationUuid" class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Stock Location</div>
                <div class="text-body-2">{{ reviewStockLocationName }}</div>
              </div>

              <!-- Lot code prefix -->
              <div v-if="form.lotCodePrefix" class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Lot Code Prefix</div>
                <div class="text-body-2">{{ form.lotCodePrefix }}</div>
              </div>

              <v-divider class="my-3" />

              <!-- Package lines -->
              <div class="text-caption text-medium-emphasis mb-2">Package Lines</div>
              <v-table density="compact">
                <thead>
                  <tr>
                    <th>Format</th>
                    <th class="text-right">Qty</th>
                    <th class="text-right">Volume</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(line, index) in form.lines" :key="index">
                    <td>{{ reviewFormatName(line.packageFormatUuid) }}</td>
                    <td class="text-right">{{ line.quantity }}</td>
                    <td class="text-right">{{ lineVolume(index) || '—' }}</td>
                  </tr>
                </tbody>
                <tfoot>
                  <tr class="font-weight-bold">
                    <td>Total</td>
                    <td class="text-right">{{ totalQuantity }}</td>
                    <td class="text-right">{{ totalPackagedLabel }}</td>
                  </tr>
                </tfoot>
              </v-table>

              <v-divider class="my-3" />

              <!-- Summary bullets -->
              <div class="text-body-2 font-weight-medium mb-2">After this packaging run:</div>
              <ul class="review-summary-list text-body-2 pl-4 mb-4">
                <li>{{ totalQuantity }} packages across {{ form.lines.length }} format{{ form.lines.length > 1 ? 's' : '' }}</li>
                <li>{{ totalPackagedLabel }} total packaged volume</li>
                <li v-if="form.closeSource">
                  {{ sourceVesselName }} will be marked empty
                </li>
                <li v-else>
                  {{ sourceVesselName }} will remain active
                </li>
                <li v-if="form.stockLocationUuid">
                  Beer lots will be created in inventory at {{ reviewStockLocationName }}
                </li>
                <li v-if="form.lossAmount">
                  {{ form.lossAmount }} {{ lossUnitLabel }} recorded as loss (manual override)
                </li>
                <li v-else-if="estimatedLossLabel !== '—'">
                  {{ estimatedLossLabel }} estimated loss
                </li>
              </ul>

              <!-- Notes -->
              <div v-if="form.notes" class="review-section mb-3">
                <div class="text-caption text-medium-emphasis">Notes</div>
                <div class="text-body-2">{{ form.notes }}</div>
              </div>

              <!-- Error display -->
              <v-alert
                v-if="saveError"
                class="mb-4"
                density="compact"
                type="error"
                variant="tonal"
              >
                {{ saveError }}
              </v-alert>
            </v-container>
          </template>
        </v-stepper>
      </v-card-text>

      <v-divider />

      <!-- Navigation actions -->
      <v-card-actions class="justify-space-between pa-4">
        <v-btn
          v-if="currentStep > 1"
          :disabled="saving"
          variant="text"
          @click="currentStep -= 1"
        >
          &larr; Back
        </v-btn>
        <v-spacer v-else />

        <div>
          <v-btn
            :disabled="saving"
            variant="text"
            @click="handleClose"
          >
            Cancel
          </v-btn>
          <v-btn
            v-if="currentStep < 3"
            color="teal"
            :disabled="!canProceed"
            @click="currentStep += 1"
          >
            Next &rarr;
          </v-btn>
          <v-btn
            v-else
            color="teal"
            :disabled="!canProceed"
            :loading="saving"
            @click="handleConfirm"
          >
            Confirm Packaging
          </v-btn>
        </div>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Batch, Occupancy, PackageFormat, StockLocation, Vessel, Volume, VolumeUnit } from '@/types'
  import { computed, ref, watch } from 'vue'
  import { useDisplay } from 'vuetify'
  import { formatDateTime, useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import { useInventoryApi } from '@/composables/useInventoryApi'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { convertVolume, volumeLabels } from '@/composables/useUnitConversion'
  import { volumeOptions as allVolumeOptions, useUnitPreferences } from '@/composables/useUnitPreferences'
  import { nowInputValue } from '@/utils/normalize'

  /** A line in the packaging form */
  interface PackagingLine {
    packageFormatUuid: string
    quantity: string
  }

  const props = defineProps<{
    modelValue: boolean
    sourceOccupancy?: Occupancy | null
    sourceVessel?: Vessel | null
    sourceBatch?: Batch | null
    sourceVolume?: Volume | null
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'packaged': []
  }>()

  const { xs } = useDisplay()
  const { createPackagingRun, getPackageFormats } = useProductionApi()
  const { getStockLocations } = useInventoryApi()
  const { showNotice } = useSnackbar()
  const { preferences, formatVolumePreferred } = useUnitPreferences()
  const { formatOccupancyStatus, getOccupancyStatusColor } = useOccupancyStatusFormatters()

  // Stepper state
  const currentStep = ref(1)
  const stepItems = ['Package Details', 'Package Lines', 'Review & Confirm']

  // Loading state
  const loadingData = ref(false)
  const loadError = ref('')
  const saving = ref(false)
  const saveError = ref('')

  // Reference data
  const packageFormats = ref<PackageFormat[]>([])
  const stockLocations = ref<StockLocation[]>([])

  // Form state
  const form = ref({
    startedAt: '',
    stockLocationUuid: '' as string | null,
    lotCodePrefix: '',
    closeSource: true,
    lines: [{ packageFormatUuid: '', quantity: '' }] as PackagingLine[],
    lossAmount: '',
    lossUnit: 'bbl' as VolumeUnit,
    notes: '',
  })

  // Validation rules
  const rules = {
    required: (v: string) => !!v || 'Required',
    positiveNumber: (v: string) => {
      const num = Number.parseFloat(v)
      if (isNaN(num)) return 'Enter a valid number'
      if (num <= 0) return 'Must be greater than 0'
      return true
    },
    positiveInteger: (v: string) => {
      const num = Number.parseInt(v, 10)
      if (isNaN(num)) return 'Enter a valid number'
      if (num <= 0) return 'Must be greater than 0'
      if (!Number.isInteger(Number(v))) return 'Must be a whole number'
      return true
    },
  }

  // Volume unit options
  const volumeUnitOptions = allVolumeOptions

  // ==================== Computed: Display helpers ====================

  const sourceVesselName = computed(() =>
    props.sourceVessel?.name ?? 'Unknown Vessel',
  )

  const sourceBatchName = computed(() =>
    props.sourceBatch?.short_name ?? `Batch ${props.sourceOccupancy?.batch_uuid?.slice(0, 8) ?? '—'}`,
  )

  // ==================== Computed: Format options ====================

  const formatMap = computed(() =>
    new Map(packageFormats.value.map(f => [f.uuid, f])),
  )

  /** Format volume for package format display, using higher precision for small values */
  function formatFormatVolume (amount: number, unit: VolumeUnit): string {
    const label = formatVolumePreferred(amount, unit)
    // If the formatted value shows "0.00", use "< 0.01" prefix instead
    if (/^0\.0+\s/.test(label)) {
      const unitLabel = label.replace(/^[\d.]+\s*/, '')
      return `< 0.01 ${unitLabel}`
    }
    return label
  }

  const activeFormatOptions = computed(() =>
    packageFormats.value
      .filter(f => f.is_active)
      .map(f => {
        const volLabel = formatFormatVolume(f.volume_per_unit, f.volume_per_unit_unit as VolumeUnit)
        return {
          title: f.name,
          value: f.uuid,
          subtitle: `${f.container} · ${volLabel}`,
        }
      }),
  )

  const stockLocationOptions = computed(() =>
    stockLocations.value.map(loc => ({
      title: loc.name,
      value: loc.uuid,
    })),
  )

  // ==================== Computed: Line volume calculations ====================

  function lineVolume (index: number): string {
    const line = form.value.lines[index]
    if (!line) return ''
    const format = formatMap.value.get(line.packageFormatUuid)
    if (!format) return ''
    const qty = Number.parseInt(line.quantity, 10)
    if (isNaN(qty) || qty <= 0) return ''

    const totalInFormatUnit = qty * format.volume_per_unit
    const fromUnit = format.volume_per_unit_unit as VolumeUnit
    return formatVolumePreferred(totalInFormatUnit, fromUnit)
  }

  /** Total packaged volume in user's preferred unit */
  const totalPackagedMl = computed(() => {
    let total = 0
    for (const line of form.value.lines) {
      const format = formatMap.value.get(line.packageFormatUuid)
      if (!format) continue
      const qty = Number.parseInt(line.quantity, 10)
      if (isNaN(qty) || qty <= 0) continue

      // Convert format volume to mL for consistent summing
      const lineMl = convertVolume(
        qty * format.volume_per_unit,
        format.volume_per_unit_unit as VolumeUnit,
        'ml',
      )
      if (lineMl !== null) {
        total += lineMl
      }
    }
    return total
  })

  const totalPackagedLabel = computed(() => {
    if (totalPackagedMl.value === 0) return '—'
    return formatVolumePreferred(totalPackagedMl.value, 'ml')
  })

  const totalQuantity = computed(() =>
    form.value.lines.reduce((sum, line) => {
      const qty = Number.parseInt(line.quantity, 10)
      return sum + (isNaN(qty) ? 0 : qty)
    }, 0),
  )

  /** Estimated loss: source volume - total packaged */
  const estimatedLossMl = computed(() => {
    if (!props.sourceVolume) return 0
    const sourceMl = convertVolume(
      props.sourceVolume.amount,
      props.sourceVolume.amount_unit,
      'ml',
    )
    if (sourceMl === null) return 0
    return Math.max(0, sourceMl - totalPackagedMl.value)
  })

  const estimatedLossPercent = computed(() => {
    if (!props.sourceVolume) return 0
    const sourceMl = convertVolume(
      props.sourceVolume.amount,
      props.sourceVolume.amount_unit,
      'ml',
    )
    if (sourceMl === null || sourceMl === 0) return 0
    return (estimatedLossMl.value / sourceMl) * 100
  })

  const estimatedLossLabel = computed(() => {
    if (totalPackagedMl.value === 0) return '—'
    if (!props.sourceVolume) return '—'
    const lossFormatted = formatVolumePreferred(estimatedLossMl.value, 'ml')
    return `${lossFormatted} (${estimatedLossPercent.value.toFixed(1)}%)`
  })

  const lossUnitLabel = computed(() =>
    volumeLabels[form.value.lossUnit] ?? form.value.lossUnit,
  )

  // ==================== Computed: Review step ====================

  const reviewDateLabel = computed(() =>
    form.value.startedAt
      ? formatDateTime(new Date(form.value.startedAt).toISOString())
      : 'Now',
  )

  const reviewStockLocationName = computed(() => {
    if (!form.value.stockLocationUuid) return ''
    const loc = stockLocations.value.find(l => l.uuid === form.value.stockLocationUuid)
    return loc?.name ?? 'Unknown Location'
  })

  function reviewFormatName (formatUuid: string): string {
    const format = formatMap.value.get(formatUuid)
    return format?.name ?? 'Unknown Format'
  }

  // ==================== Computed: Validation ====================

  const canProceed = computed(() => {
    if (currentStep.value === 1) return canProceedStep1.value
    if (currentStep.value === 2) return canProceedStep2.value
    if (currentStep.value === 3) return canProceedStep2.value // same validation for confirm
    return false
  })

  const canProceedStep1 = computed(() => {
    // Source occupancy is required (always pre-filled from props)
    return !!props.sourceOccupancy
  })

  const canProceedStep2 = computed(() => {
    // At least one line with valid format and quantity
    if (form.value.lines.length === 0) return false
    return form.value.lines.every(line => {
      const qty = Number.parseInt(line.quantity, 10)
      return !!line.packageFormatUuid && !isNaN(qty) && qty > 0
    })
  })

  // ==================== Line management ====================

  function addLine () {
    form.value.lines.push({ packageFormatUuid: '', quantity: '' })
  }

  function removeLine (index: number) {
    if (form.value.lines.length <= 1) return
    form.value.lines.splice(index, 1)
  }

  // ==================== Watch: Dialog open/close ====================

  watch(
    () => props.modelValue,
    async isOpen => {
      if (isOpen) {
        await resetAndLoad()
      }
    },
  )

  async function resetAndLoad () {
    // Reset state
    currentStep.value = 1
    saving.value = false
    saveError.value = ''
    loadError.value = ''

    // Reset form
    form.value = {
      startedAt: nowInputValue(),
      stockLocationUuid: null,
      lotCodePrefix: props.sourceBatch?.short_name ?? '',
      closeSource: true,
      lines: [{ packageFormatUuid: '', quantity: '' }],
      lossAmount: '',
      lossUnit: preferences.value.volume,
      notes: '',
    }

    // Load reference data
    loadingData.value = true
    try {
      const [formatsData, locationsData] = await Promise.all([
        getPackageFormats(),
        getStockLocations(),
      ])
      packageFormats.value = formatsData
      stockLocations.value = locationsData
    } catch {
      loadError.value = 'Failed to load packaging data. Please close and try again.'
    } finally {
      loadingData.value = false
    }
  }

  // ==================== Actions ====================

  function handleClose () {
    emit('update:modelValue', false)
  }

  async function handleConfirm () {
    if (!props.sourceOccupancy || !props.sourceBatch) {
      saveError.value = 'Missing source occupancy or batch data'
      return
    }

    saving.value = true
    saveError.value = ''

    try {
      const startedAt = form.value.startedAt
        ? new Date(form.value.startedAt).toISOString()
        : new Date().toISOString()

      const lossAmount = form.value.lossAmount
        ? Number.parseFloat(form.value.lossAmount)
        : undefined

      await createPackagingRun({
        batch_uuid: props.sourceBatch.uuid,
        occupancy_uuid: props.sourceOccupancy.uuid,
        started_at: startedAt,
        lines: form.value.lines.map(line => ({
          package_format_uuid: line.packageFormatUuid,
          quantity: Number.parseInt(line.quantity, 10),
        })),
        close_source: form.value.closeSource,
        stock_location_uuid: form.value.stockLocationUuid || undefined,
        lot_code_prefix: form.value.lotCodePrefix || undefined,
        loss_amount: lossAmount && lossAmount > 0 ? lossAmount : undefined,
        loss_unit: lossAmount && lossAmount > 0 ? form.value.lossUnit : undefined,
        notes: form.value.notes || undefined,
      })

      showNotice('Packaging run recorded')
      emit('packaged')
      emit('update:modelValue', false)
    } catch (error) {
      saveError.value = error instanceof Error ? error.message : 'Failed to create packaging run'
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
.metric-card {
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(var(--v-theme-surface), 0.5);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  text-align: center;
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
  font-size: 1.1rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.87);
}

.review-section {
  padding: 8px 0;
}

.review-summary-list {
  list-style-type: disc;
}

.review-summary-list li {
  margin-bottom: 4px;
}
</style>
