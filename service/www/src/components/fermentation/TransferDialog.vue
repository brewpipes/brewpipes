<template>
  <v-dialog
    :fullscreen="xs"
    :max-width="xs ? '100%' : 600"
    :model-value="modelValue"
    persistent
    scrollable
    @update:model-value="emit('update:modelValue', $event)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <span class="text-h6">{{ dialogTitle }}</span>
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
              <v-progress-circular color="primary" indeterminate size="48" />
              <p class="text-body-2 text-medium-emphasis mt-4">Loading vessel data...</p>
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
          <!-- Step 1: Transfer Details -->
          <template #item.1>
            <v-container class="pa-4">
              <!-- Mode selector -->
              <v-btn-toggle
                v-model="activeMode"
                class="mb-4"
                color="primary"
                density="compact"
                mandatory
                variant="outlined"
              >
                <v-btn value="transfer">
                  <v-icon class="d-sm-none" icon="mdi-arrow-right" />
                  <span class="d-none d-sm-inline">Transfer</span>
                  <v-tooltip activator="parent" class="d-sm-none" location="bottom">
                    Transfer
                  </v-tooltip>
                </v-btn>
                <v-btn value="split">
                  <v-icon class="d-sm-none" icon="mdi-call-split" />
                  <span class="d-none d-sm-inline">Split</span>
                  <v-tooltip activator="parent" class="d-sm-none" location="bottom">
                    Split
                  </v-tooltip>
                </v-btn>
                <v-btn value="blend">
                  <v-icon class="d-sm-none" icon="mdi-call-merge" />
                  <span class="d-none d-sm-inline">Blend</span>
                  <v-tooltip activator="parent" class="d-sm-none" location="bottom">
                    Blend
                  </v-tooltip>
                </v-btn>
              </v-btn-toggle>

              <!-- ==================== TRANSFER MODE ==================== -->
              <template v-if="activeMode === 'transfer'">
                <!-- FROM section -->
                <div class="text-overline text-medium-emphasis mb-2">From</div>

                <!-- Pre-selected source: read-only card -->
                <v-card
                  v-if="sourceOccupancy"
                  class="mb-4"
                  density="compact"
                  variant="outlined"
                >
                  <v-card-text class="pa-3">
                    <div class="d-flex align-center flex-wrap ga-2">
                      <v-icon icon="mdi-flask-round-bottom" size="small" />
                      <span class="text-body-1 font-weight-medium">
                        {{ sourceVesselName }}
                      </span>
                      <span class="text-body-2 text-medium-emphasis">
                        · {{ sourceBatchName }}
                      </span>
                    </div>
                    <div class="d-flex align-center flex-wrap ga-2 mt-1">
                      <v-chip
                        v-if="sourceOccupancy.status"
                        :color="getOccupancyStatusColor(sourceOccupancy.status)"
                        size="x-small"
                        variant="tonal"
                      >
                        {{ formatOccupancyStatus(sourceOccupancy.status) }}
                      </v-chip>
                      <span v-if="sourceVolume" class="text-body-2">
                        {{ formatVolumePreferred(sourceVolume.amount, sourceVolume.amount_unit) }}
                      </span>
                      <span class="text-body-2 text-medium-emphasis">
                        · Day {{ sourceDaysInTank }}
                      </span>
                    </div>
                  </v-card-text>
                </v-card>

                <!-- No pre-selected source: select from occupied vessels -->
                <v-select
                  v-else
                  v-model="form.sourceOccupancyUuid"
                  class="mb-2"
                  density="comfortable"
                  item-title="title"
                  item-value="value"
                  :items="occupiedVesselOptions"
                  label="Source vessel"
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
                      <v-list-item-title>No occupied vessels</v-list-item-title>
                      <v-list-item-subtitle>No vessels currently have beer in them</v-list-item-subtitle>
                    </v-list-item>
                  </template>
                </v-select>

                <v-divider class="my-3" />

                <!-- TO section -->
                <div class="text-overline text-medium-emphasis mb-2">To</div>

                <v-select
                  v-model="form.destVesselUuid"
                  class="mb-2"
                  density="comfortable"
                  item-title="title"
                  item-value="value"
                  :items="availableVesselOptions"
                  label="Destination vessel"
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
                      <v-list-item-title>No available vessels</v-list-item-title>
                      <v-list-item-subtitle>Free up a vessel first</v-list-item-subtitle>
                    </v-list-item>
                  </template>
                </v-select>

                <v-alert
                  v-if="!loadingData && availableVesselOptions.length === 0"
                  class="mb-4"
                  density="compact"
                  type="warning"
                  variant="tonal"
                >
                  No available vessels. Free up a vessel first.
                </v-alert>

                <v-divider class="my-3" />

                <!-- VOLUME -->
                <div class="text-overline text-medium-emphasis mb-2">Volume</div>

                <v-row dense>
                  <v-col cols="7">
                    <v-text-field
                      v-model="form.transferAmount"
                      density="comfortable"
                      inputmode="decimal"
                      label="Transfer volume"
                      :rules="form.transferAmount ? [rules.positiveNumber] : [rules.required]"
                    />
                  </v-col>
                  <v-col cols="5">
                    <v-select
                      v-model="form.volumeUnit"
                      density="comfortable"
                      item-title="label"
                      item-value="value"
                      :items="volumeUnitOptions"
                      label="Unit"
                    />
                  </v-col>
                </v-row>

                <!-- LOSS -->
                <v-row dense>
                  <v-col cols="7">
                    <v-text-field
                      v-model="form.lossAmount"
                      density="comfortable"
                      hint="Beer left behind in hoses, trub, etc."
                      inputmode="decimal"
                      label="Transfer loss"
                      persistent-hint
                      :rules="form.lossAmount ? [rules.nonNegativeNumber] : []"
                    />
                  </v-col>
                  <v-col cols="5">
                    <v-select
                      v-model="form.volumeUnit"
                      density="comfortable"
                      item-title="label"
                      item-value="value"
                      :items="volumeUnitOptions"
                      label="Unit"
                    />
                  </v-col>
                </v-row>

                <v-divider class="my-3" />

                <!-- DESTINATION STATUS -->
                <v-select
                  v-model="form.destStatus"
                  class="mb-2"
                  density="comfortable"
                  item-title="title"
                  item-value="value"
                  :items="statusOptions"
                  label="Destination status"
                />

                <!-- CLOSE SOURCE CHECKBOX -->
                <v-checkbox
                  v-model="form.closeSource"
                  density="compact"
                  label="Close source vessel after transfer"
                />

                <!-- TRANSFER DATE -->
                <div v-if="!showTimePicker" class="mb-2">
                  <v-btn
                    density="comfortable"
                    prepend-icon="mdi-clock-outline"
                    size="small"
                    variant="text"
                    @click="showTimePicker = true"
                  >
                    Change time
                  </v-btn>
                </div>
                <v-text-field
                  v-else
                  v-model="form.transferDate"
                  class="mb-2"
                  density="comfortable"
                  label="Transfer date/time"
                  type="datetime-local"
                />
              </template>

              <!-- ==================== SPLIT MODE ==================== -->
              <template v-if="activeMode === 'split'">
                <!-- FROM section -->
                <div class="text-overline text-medium-emphasis mb-2">From</div>

                <!-- Pre-selected source: read-only card -->
                <v-card
                  v-if="sourceOccupancy"
                  class="mb-4"
                  density="compact"
                  variant="outlined"
                >
                  <v-card-text class="pa-3">
                    <div class="d-flex align-center flex-wrap ga-2">
                      <v-icon icon="mdi-flask-round-bottom" size="small" />
                      <span class="text-body-1 font-weight-medium">
                        {{ sourceVesselName }}
                      </span>
                      <span class="text-body-2 text-medium-emphasis">
                        · {{ sourceBatchName }}
                      </span>
                    </div>
                    <div class="d-flex align-center flex-wrap ga-2 mt-1">
                      <v-chip
                        v-if="sourceOccupancy.status"
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

                <!-- No pre-selected source: select from occupied vessels -->
                <v-select
                  v-else
                  v-model="form.sourceOccupancyUuid"
                  class="mb-2"
                  density="comfortable"
                  item-title="title"
                  item-value="value"
                  :items="occupiedVesselOptions"
                  label="Source vessel"
                  :rules="[rules.required]"
                >
                  <template #item="{ props: itemProps, item }">
                    <v-list-item v-bind="itemProps">
                      <template #subtitle>
                        <span>{{ item.raw.subtitle }}</span>
                      </template>
                    </v-list-item>
                  </template>
                </v-select>

                <v-divider class="my-3" />

                <!-- SPLIT INTO section -->
                <div class="text-overline text-medium-emphasis mb-2">Split Into</div>

                <div
                  v-for="(dest, index) in splitDestinations"
                  :key="index"
                  class="mb-3"
                >
                  <div class="d-flex align-center ga-2 mb-1">
                    <span class="text-caption text-medium-emphasis font-weight-medium">
                      Destination {{ index + 1 }}
                    </span>
                    <v-spacer />
                    <v-btn
                      v-if="splitDestinations.length > 2"
                      aria-label="Remove destination"
                      density="compact"
                      icon="mdi-close"
                      size="x-small"
                      variant="text"
                      @click="removeSplitDestination(index)"
                    />
                  </div>

                  <v-select
                    v-model="dest.vesselUuid"
                    class="mb-1"
                    density="comfortable"
                    item-title="title"
                    item-value="value"
                    :items="splitAvailableVesselsFor(index)"
                    label="Vessel"
                    :rules="[rules.required]"
                  >
                    <template #item="{ props: itemProps, item }">
                      <v-list-item v-bind="itemProps">
                        <template #subtitle>
                          <span>{{ item.raw.subtitle }}</span>
                        </template>
                      </v-list-item>
                    </template>
                  </v-select>

                  <v-row dense>
                    <v-col cols="7">
                      <v-text-field
                        v-model="dest.amount"
                        density="comfortable"
                        inputmode="decimal"
                        label="Volume"
                        :rules="dest.amount ? [rules.positiveNumber] : [rules.required]"
                      />
                    </v-col>
                    <v-col cols="5">
                      <v-select
                        v-model="form.volumeUnit"
                        density="comfortable"
                        item-title="label"
                        item-value="value"
                        :items="volumeUnitOptions"
                        label="Unit"
                      />
                    </v-col>
                  </v-row>

                  <v-select
                    v-model="dest.status"
                    density="comfortable"
                    item-title="title"
                    item-value="value"
                    :items="statusOptions"
                    label="Status"
                  />

                  <v-divider v-if="index < splitDestinations.length - 1" class="mt-1" />
                </div>

                <v-btn
                  v-if="splitDestinations.length < 4"
                  class="mb-3"
                  density="comfortable"
                  prepend-icon="mdi-plus"
                  variant="text"
                  @click="addSplitDestination"
                >
                  Add Destination
                </v-btn>

                <v-divider class="my-3" />

                <!-- VOLUME MATH -->
                <div class="text-overline text-medium-emphasis mb-2">Volume</div>

                <v-row dense>
                  <v-col cols="7">
                    <v-text-field
                      v-model="form.lossAmount"
                      density="comfortable"
                      hint="Beer left behind in hoses, trub, etc."
                      inputmode="decimal"
                      label="Transfer loss"
                      persistent-hint
                      :rules="form.lossAmount ? [rules.nonNegativeNumber] : []"
                    />
                  </v-col>
                  <v-col cols="5">
                    <v-select
                      v-model="form.volumeUnit"
                      density="comfortable"
                      item-title="label"
                      item-value="value"
                      :items="volumeUnitOptions"
                      label="Unit"
                    />
                  </v-col>
                </v-row>

                <!-- Running total -->
                <div class="text-body-2 text-medium-emphasis mb-1">
                  {{ splitRunningTotalLabel }}
                </div>
                <div v-if="splitVolumeMatch === 'match'" class="text-body-2 text-success mb-3">
                  ✓ Matches source volume
                </div>
                <div v-else-if="splitVolumeMatch === 'mismatch'" class="text-body-2 text-warning mb-3">
                  ⚠ {{ splitUnaccountedLabel }}
                </div>

                <v-divider class="my-3" />

                <!-- CLOSE SOURCE CHECKBOX -->
                <v-checkbox
                  v-model="form.closeSource"
                  density="compact"
                  label="Close source vessel after split"
                />

                <!-- TRANSFER DATE -->
                <div v-if="!showTimePicker" class="mb-2">
                  <v-btn
                    density="comfortable"
                    prepend-icon="mdi-clock-outline"
                    size="small"
                    variant="text"
                    @click="showTimePicker = true"
                  >
                    Change time
                  </v-btn>
                </div>
                <v-text-field
                  v-else
                  v-model="form.transferDate"
                  class="mb-2"
                  density="comfortable"
                  label="Transfer date/time"
                  type="datetime-local"
                />
              </template>

              <!-- ==================== BLEND MODE ==================== -->
              <template v-if="activeMode === 'blend'">
                <!-- BLEND FROM section -->
                <div class="text-overline text-medium-emphasis mb-2">Blend From</div>

                <div
                  v-for="(src, index) in blendSources"
                  :key="index"
                  class="mb-3"
                >
                  <div class="d-flex align-center ga-2 mb-1">
                    <span class="text-caption text-medium-emphasis font-weight-medium">
                      Source {{ index + 1 }}
                    </span>
                    <v-spacer />
                    <v-btn
                      v-if="blendSources.length > 2"
                      aria-label="Remove source"
                      density="compact"
                      icon="mdi-close"
                      size="x-small"
                      variant="text"
                      @click="removeBlendSource(index)"
                    />
                  </div>

                  <v-select
                    v-model="src.occupancyUuid"
                    class="mb-1"
                    density="comfortable"
                    item-title="title"
                    item-value="value"
                    :items="blendOccupiedVesselsFor(index)"
                    label="Vessel"
                    :rules="[rules.required]"
                  >
                    <template #item="{ props: itemProps, item }">
                      <v-list-item v-bind="itemProps">
                        <template #subtitle>
                          <span>{{ item.raw.subtitle }}</span>
                        </template>
                      </v-list-item>
                    </template>
                  </v-select>

                  <v-row dense>
                    <v-col cols="7">
                      <v-text-field
                        v-model="src.amount"
                        density="comfortable"
                        inputmode="decimal"
                        label="Volume"
                        :rules="src.amount ? [rules.positiveNumber] : [rules.required]"
                      />
                    </v-col>
                    <v-col cols="5">
                      <v-select
                        v-model="form.volumeUnit"
                        density="comfortable"
                        item-title="label"
                        item-value="value"
                        :items="volumeUnitOptions"
                        label="Unit"
                      />
                    </v-col>
                  </v-row>

                  <v-checkbox
                    v-model="src.closeAfterBlend"
                    density="compact"
                    label="Close after blend"
                  />

                  <v-divider v-if="index < blendSources.length - 1" class="mt-1" />
                </div>

                <v-btn
                  v-if="blendSources.length < 4"
                  class="mb-3"
                  density="comfortable"
                  prepend-icon="mdi-plus"
                  variant="text"
                  @click="addBlendSource"
                >
                  Add Source
                </v-btn>

                <v-divider class="my-3" />

                <!-- INTO section -->
                <div class="text-overline text-medium-emphasis mb-2">Into</div>

                <v-select
                  v-model="form.destVesselUuid"
                  class="mb-2"
                  density="comfortable"
                  item-title="title"
                  item-value="value"
                  :items="blendAvailableVesselOptions"
                  label="Destination vessel"
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
                      <v-list-item-title>No available vessels</v-list-item-title>
                      <v-list-item-subtitle>Free up a vessel first</v-list-item-subtitle>
                    </v-list-item>
                  </template>
                </v-select>

                <v-select
                  v-model="form.destStatus"
                  class="mb-2"
                  density="comfortable"
                  item-title="title"
                  item-value="value"
                  :items="statusOptions"
                  label="Destination status"
                />

                <v-divider class="my-3" />

                <!-- VOLUME MATH -->
                <div class="text-overline text-medium-emphasis mb-2">Volume</div>

                <v-row dense>
                  <v-col cols="7">
                    <v-text-field
                      v-model="form.lossAmount"
                      density="comfortable"
                      hint="Beer left behind in hoses, trub, etc."
                      inputmode="decimal"
                      label="Transfer loss"
                      persistent-hint
                      :rules="form.lossAmount ? [rules.nonNegativeNumber] : []"
                    />
                  </v-col>
                  <v-col cols="5">
                    <v-select
                      v-model="form.volumeUnit"
                      density="comfortable"
                      item-title="label"
                      item-value="value"
                      :items="volumeUnitOptions"
                      label="Unit"
                    />
                  </v-col>
                </v-row>

                <div class="text-body-2 text-medium-emphasis mb-3">
                  {{ blendVolumeMathLabel }}
                </div>

                <!-- Batch identity (when sources are from different batches) -->
                <template v-if="blendHasMultipleBatches">
                  <v-divider class="my-3" />
                  <v-alert class="mb-3" density="compact" type="warning" variant="tonal">
                    Blending from different batches
                  </v-alert>
                  <div class="text-body-2 mb-2">Which batch should the blend belong to?</div>
                  <v-radio-group v-model="blendSelectedBatchUuid" density="compact">
                    <v-radio
                      v-for="opt in blendBatchOptions"
                      :key="opt.uuid"
                      :label="opt.label"
                      :value="opt.uuid"
                    />
                  </v-radio-group>
                </template>

                <!-- TRANSFER DATE -->
                <div v-if="!showTimePicker" class="mb-2">
                  <v-btn
                    density="comfortable"
                    prepend-icon="mdi-clock-outline"
                    size="small"
                    variant="text"
                    @click="showTimePicker = true"
                  >
                    Change time
                  </v-btn>
                </div>
                <v-text-field
                  v-else
                  v-model="form.transferDate"
                  class="mb-2"
                  density="comfortable"
                  label="Transfer date/time"
                  type="datetime-local"
                />
              </template>
            </v-container>
          </template>

          <!-- Step 2: Review & Confirm -->
          <template #item.2>
            <v-container class="pa-4">
              <!-- ==================== TRANSFER REVIEW ==================== -->
              <template v-if="activeMode === 'transfer'">
                <!-- Visual flow diagram -->
                <div class="transfer-flow-diagram mx-auto mb-4">
                  <div class="flow-source text-body-1 font-weight-medium">
                    {{ reviewSourceVesselName }}
                    <span class="text-medium-emphasis font-weight-regular">
                      ({{ reviewSourceBatchName }})
                    </span>
                  </div>
                  <div class="flow-arrow-section">
                    <div class="flow-line" />
                    <div class="flow-label text-body-2">
                      {{ reviewTransferLabel }}
                      <span v-if="reviewLossLabel" class="text-medium-emphasis">
                        ({{ reviewLossLabel }} loss)
                      </span>
                    </div>
                    <div class="flow-line" />
                    <v-icon class="flow-arrow-icon" icon="mdi-arrow-down" size="24" />
                  </div>
                  <div class="flow-dest text-body-1 font-weight-medium">
                    {{ reviewDestVesselName }}
                    <span class="text-medium-emphasis font-weight-regular">
                      → {{ formatOccupancyStatus(form.destStatus) }}
                    </span>
                  </div>
                </div>

                <v-divider class="mb-4" />

                <!-- Summary bullets -->
                <div class="text-body-2 font-weight-medium mb-2">After this transfer:</div>
                <ul class="transfer-summary-list text-body-2 pl-4 mb-4">
                  <li v-if="form.closeSource">
                    {{ reviewSourceVesselName }} will be marked empty
                  </li>
                  <li v-else>
                    {{ reviewSourceVesselName }} will remain active with
                    {{ reviewRemainingLabel }}
                  </li>
                  <li>
                    {{ reviewDestVesselName }} will hold {{ reviewDestReceivesLabel }}
                  </li>
                  <li>
                    Status: {{ formatOccupancyStatus(form.destStatus) }}
                  </li>
                </ul>
              </template>

              <!-- ==================== SPLIT REVIEW ==================== -->
              <template v-if="activeMode === 'split'">
                <div class="transfer-flow-diagram mx-auto mb-4">
                  <div class="flow-source text-body-1 font-weight-medium">
                    {{ reviewSourceVesselName }}
                    <span class="text-medium-emphasis font-weight-regular">
                      ({{ reviewSourceBatchName }}, {{ reviewSourceVolumeLabel }})
                    </span>
                  </div>
                  <div class="split-flow-branches">
                    <div
                      v-for="(dest, index) in splitDestinations"
                      :key="index"
                      class="split-flow-branch text-body-2"
                    >
                      <span class="split-branch-connector">{{ index < splitDestinations.length - 1 ? '├──' : '└──' }}</span>
                      {{ parseSplitDestAmount(dest.amount) }} {{ unitLabel }}
                      → {{ splitDestVesselName(dest.vesselUuid) }}
                      ({{ formatOccupancyStatus(dest.status) }})
                    </div>
                    <div v-if="lossAmountNum > 0" class="split-flow-branch text-body-2 text-medium-emphasis">
                      <span class="split-branch-connector">└──</span>
                      {{ lossAmountNum }} {{ unitLabel }} loss
                    </div>
                  </div>
                </div>

                <v-divider class="mb-4" />

                <div class="text-body-2 font-weight-medium mb-2">After this split:</div>
                <ul class="transfer-summary-list text-body-2 pl-4 mb-4">
                  <li v-if="form.closeSource">
                    {{ reviewSourceVesselName }} will be marked empty
                  </li>
                  <li v-else>
                    {{ reviewSourceVesselName }} will remain active
                  </li>
                  <li v-for="(dest, index) in splitDestinations" :key="index">
                    {{ splitDestVesselName(dest.vesselUuid) }} will hold
                    {{ parseSplitDestAmount(dest.amount) }} {{ unitLabel }}
                    ({{ formatOccupancyStatus(dest.status) }})
                  </li>
                </ul>
              </template>

              <!-- ==================== BLEND REVIEW ==================== -->
              <template v-if="activeMode === 'blend'">
                <div class="transfer-flow-diagram mx-auto mb-4">
                  <div
                    v-for="(src, index) in blendSources"
                    :key="index"
                    class="blend-flow-source text-body-2"
                  >
                    {{ blendSourceVesselName(src.occupancyUuid) }}
                    <span class="text-medium-emphasis">
                      ({{ blendSourceBatchName(src.occupancyUuid) }}, {{ parseBlendSourceAmount(src.amount) }} {{ unitLabel }})
                    </span>
                    {{ index < blendSources.length - 1 ? '──┐' : '──┘' }}
                  </div>
                  <div class="blend-flow-dest text-body-1 font-weight-medium mt-2">
                    <span class="text-medium-emphasis mr-1">├──</span>
                    {{ blendDestAmount.toFixed(2) }} {{ unitLabel }}
                    → {{ reviewDestVesselName }}
                    ({{ formatOccupancyStatus(form.destStatus) }})
                  </div>
                  <div v-if="lossAmountNum > 0" class="blend-flow-loss text-body-2 text-medium-emphasis mt-1 ml-6">
                    {{ lossAmountNum }} {{ unitLabel }} loss
                  </div>
                </div>

                <v-divider class="mb-4" />

                <div class="text-body-2 font-weight-medium mb-2">After this blend:</div>
                <ul class="transfer-summary-list text-body-2 pl-4 mb-4">
                  <li v-for="(src, index) in blendSources" :key="index">
                    {{ blendSourceVesselName(src.occupancyUuid) }}
                    {{ src.closeAfterBlend ? 'will be marked empty' : 'will remain active' }}
                  </li>
                  <li>
                    {{ reviewDestVesselName }} will hold
                    {{ blendDestAmount.toFixed(2) }} {{ unitLabel }}
                    ({{ formatOccupancyStatus(form.destStatus) }})
                  </li>
                  <li>
                    Batch: {{ blendReviewBatchName }}
                  </li>
                </ul>
              </template>

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
          v-if="currentStep === 2"
          :disabled="saving"
          variant="text"
          @click="currentStep = 1"
        >
          ← Back
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
            v-if="currentStep === 1"
            color="primary"
            :disabled="!canProceedToReview"
            @click="currentStep = 2"
          >
            Next →
          </v-btn>
          <v-btn
            v-else
            color="primary"
            :disabled="!canProceedToReview"
            :loading="saving"
            @click="handleConfirm"
          >
            {{ confirmButtonLabel }}
          </v-btn>
        </div>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import type { Batch, Occupancy, OccupancyStatus, Vessel, Volume, VolumeUnit } from '@/types'
  import { computed, ref, watch } from 'vue'
  import { useDisplay } from 'vuetify'
  import { useOccupancyStatusFormatters, useVesselTypeFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { convertVolume, volumeLabels } from '@/composables/useUnitConversion'
  import { volumeOptions as allVolumeOptions, useUnitPreferences } from '@/composables/useUnitPreferences'
  import { OCCUPANCY_STATUS_VALUES } from '@/types'
  import { nowInputValue } from '@/utils/normalize'

  /** Transfer dialog mode */
  type TransferMode = 'transfer' | 'split' | 'blend'

  /** A destination row in split mode */
  interface SplitDestination {
    vesselUuid: string
    amount: string
    status: OccupancyStatus
  }

  /** A source row in blend mode */
  interface BlendSource {
    occupancyUuid: string
    amount: string
    closeAfterBlend: boolean
  }

  const props = defineProps<{
    modelValue: boolean
    mode?: TransferMode
    sourceOccupancy?: Occupancy | null
    sourceVessel?: Vessel | null
    sourceBatch?: Batch | null
    sourceVolume?: Volume | null
  }>()

  const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'transferred': []
  }>()

  const { xs } = useDisplay()
  const {
    getActiveOccupancies,
    getVessels,
    getVolume,
    getBatch,
    createTransfer,
    createVolume,
    createVolumeRelation,
    createBatchVolume,
  } = useProductionApi()
  const { showNotice } = useSnackbar()
  const { preferences, formatVolumePreferred } = useUnitPreferences()
  const { formatOccupancyStatus, getOccupancyStatusColor } = useOccupancyStatusFormatters()
  const { formatVesselType } = useVesselTypeFormatters()

  // Stepper state
  const currentStep = ref(1)
  const stepItems = ['Transfer Details', 'Review & Confirm']

  // Mode state
  const activeMode = ref<TransferMode>('transfer')

  // Loading state
  const loadingData = ref(false)
  const saving = ref(false)
  const saveError = ref('')
  const loadError = ref('')
  const showTimePicker = ref(false)

  // Reference data
  const allOccupancies = ref<Occupancy[]>([])
  const allVessels = ref<Vessel[]>([])

  // Resolved data for dynamically-selected source
  const resolvedSourceOccupancy = ref<Occupancy | null>(null)
  const resolvedSourceVessel = ref<Vessel | null>(null)
  const resolvedSourceBatch = ref<Batch | null>(null)
  const resolvedSourceVolume = ref<Volume | null>(null)

  // Split mode state
  const splitDestinations = ref<SplitDestination[]>([])

  // Blend mode state
  const blendSources = ref<BlendSource[]>([])
  const blendSelectedBatchUuid = ref('')

  // Resolved blend source data (occupancy uuid → batch/volume/vessel)
  const blendResolvedBatches = ref<Map<string, Batch>>(new Map())
  const blendResolvedVolumes = ref<Map<string, Volume>>(new Map())

  // Form state
  const form = ref({
    sourceOccupancyUuid: '',
    destVesselUuid: '',
    transferAmount: '',
    lossAmount: '0',
    volumeUnit: 'bbl' as VolumeUnit,
    destStatus: 'conditioning' as OccupancyStatus,
    closeSource: true,
    transferDate: '',
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
    nonNegativeNumber: (v: string) => {
      const num = Number.parseFloat(v)
      if (isNaN(num)) return 'Enter a valid number'
      if (num < 0) return 'Cannot be negative'
      return true
    },
  }

  // Volume unit options for the select
  const volumeUnitOptions = allVolumeOptions

  // Status options for destination
  const statusOptions = computed(() =>
    OCCUPANCY_STATUS_VALUES.map(status => ({
      value: status,
      title: formatOccupancyStatus(status),
    })),
  )

  // ==================== Computed: Dialog title ====================

  const dialogTitle = computed(() => {
    switch (activeMode.value) {
      case 'split': { return 'Split Beer'
      }
      case 'blend': { return 'Blend Beer'
      }
      default: { return 'Transfer Beer'
      }
    }
  })

  const confirmButtonLabel = computed(() => {
    switch (activeMode.value) {
      case 'split': { return 'Confirm Split'
      }
      case 'blend': { return 'Confirm Blend'
      }
      default: { return 'Confirm Transfer'
      }
    }
  })

  // ==================== Computed: Effective source data ====================

  /** The effective source occupancy (from props or dynamically selected) */
  const effectiveSourceOccupancy = computed<Occupancy | null>(() => {
    if (props.sourceOccupancy) return props.sourceOccupancy
    return resolvedSourceOccupancy.value
  })

  const effectiveSourceVessel = computed<Vessel | null>(() => {
    if (props.sourceVessel) return props.sourceVessel
    return resolvedSourceVessel.value
  })

  const effectiveSourceBatch = computed<Batch | null>(() => {
    if (props.sourceBatch) return props.sourceBatch
    return resolvedSourceBatch.value
  })

  const effectiveSourceVolume = computed<Volume | null>(() => {
    if (props.sourceVolume) return props.sourceVolume
    return resolvedSourceVolume.value
  })

  // ==================== Computed: Display helpers ====================

  const sourceVesselName = computed(() =>
    props.sourceVessel?.name ?? 'Unknown Vessel',
  )

  const sourceBatchName = computed(() => {
    if (props.sourceBatch?.short_name) return props.sourceBatch.short_name
    // Try resolving from blend batch data if source batch prop is not available
    const batchUuid = props.sourceOccupancy?.batch_uuid
    if (batchUuid) {
      const resolved = blendResolvedBatches.value.get(batchUuid)
      if (resolved?.short_name) return resolved.short_name
    }
    return `Batch ${batchUuid?.slice(0, 8) ?? '—'}`
  })

  const sourceDaysInTank = computed(() => {
    const occ = props.sourceOccupancy
    if (!occ?.in_at) return 0
    const inAt = new Date(occ.in_at)
    const now = new Date()
    return Math.max(0, Math.floor((now.getTime() - inAt.getTime()) / (1000 * 60 * 60 * 24)))
  })

  // ==================== Computed: Vessel options ====================

  const vesselMap = computed(() =>
    new Map(allVessels.value.map(v => [v.uuid, v])),
  )

  const occupiedVesselUuids = computed(() =>
    new Set(allOccupancies.value.map(o => o.vessel_uuid)),
  )

  /** Options for source vessel select (occupied vessels only) */
  const occupiedVesselOptions = computed(() => {
    return allOccupancies.value
      .filter(occ => occ.batch_uuid)
      .map(occ => {
        const vessel = vesselMap.value.get(occ.vessel_uuid)
        return {
          title: vessel
            ? `${vessel.name} · ${blendResolvedBatches.value.get(occ.batch_uuid!)?.short_name ?? `Batch ${occ.batch_uuid?.slice(0, 8) ?? '—'}`}`
            : `Vessel ${occ.vessel_uuid.slice(0, 8)}`,
          value: occ.uuid,
          subtitle: vessel ? `${formatVesselType(vessel.type)}` : '',
        }
      })
  })

  /** Options for destination vessel select (unoccupied, active vessels) */
  const availableVesselOptions = computed(() => {
    // Exclude the source vessel from available destinations
    const sourceVesselUuid = effectiveSourceOccupancy.value?.vessel_uuid
    return allVessels.value
      .filter(v =>
        v.status === 'active'
        && !occupiedVesselUuids.value.has(v.uuid)
        && v.uuid !== sourceVesselUuid,
      )
      .map(v => ({
        title: v.name,
        value: v.uuid,
        subtitle: `${formatVesselType(v.type)} · ${v.capacity} ${volumeLabels[v.capacity_unit] ?? v.capacity_unit}`,
      }))
  })

  // ==================== Split mode: vessel options ====================

  /** Available vessels for a split destination row, excluding vessels already selected in other rows */
  function splitAvailableVesselsFor (rowIndex: number): Array<{ title: string, value: string, subtitle: string }> {
    const selectedInOtherRows = new Set(
      splitDestinations.value
        .filter((_, i) => i !== rowIndex)
        .map(d => d.vesselUuid)
        .filter(Boolean),
    )
    return availableVesselOptions.value.filter(v => !selectedInOtherRows.has(v.value))
  }

  function addSplitDestination () {
    if (splitDestinations.value.length >= 4) return
    splitDestinations.value.push({
      vesselUuid: '',
      amount: '',
      status: getDefaultDestStatus(effectiveSourceOccupancy.value?.status ?? null),
    })
  }

  function removeSplitDestination (index: number) {
    if (splitDestinations.value.length <= 2) return
    splitDestinations.value.splice(index, 1)
  }

  // ==================== Split mode: volume math ====================

  function parseSplitDestAmount (amount: string): number {
    const num = Number.parseFloat(amount)
    return isNaN(num) ? 0 : num
  }

  const splitTotalDestVolume = computed(() =>
    splitDestinations.value.reduce((sum, d) => sum + parseSplitDestAmount(d.amount), 0),
  )

  const splitTotalWithLoss = computed(() =>
    splitTotalDestVolume.value + lossAmountNum.value,
  )

  const splitRunningTotalLabel = computed(() => {
    const parts = splitDestinations.value.map(d =>
      `${parseSplitDestAmount(d.amount).toFixed(1)}`,
    )
    const lossStr = lossAmountNum.value > 0 ? ` + ${lossAmountNum.value} loss` : ''
    return `${parts.join(' + ')}${lossStr} = ${splitTotalWithLoss.value.toFixed(1)} ${unitLabel.value}`
  })

  const sourceVolumeInFormUnit = computed(() => {
    const vol = effectiveSourceVolume.value
    if (!vol) return null
    return convertVolume(vol.amount, vol.amount_unit, form.value.volumeUnit)
  })

  const splitVolumeMatch = computed<'match' | 'mismatch' | 'unknown'>(() => {
    const sourceAmt = sourceVolumeInFormUnit.value
    if (sourceAmt === null) return 'unknown'
    const total = splitTotalWithLoss.value
    // Within 0.1%
    const tolerance = sourceAmt * 0.001
    if (Math.abs(total - sourceAmt) <= tolerance) return 'match'
    return 'mismatch'
  })

  const splitUnaccountedLabel = computed(() => {
    const sourceAmt = sourceVolumeInFormUnit.value
    if (sourceAmt === null) return ''
    const diff = sourceAmt - splitTotalWithLoss.value
    return `${Math.abs(diff).toFixed(1)} ${unitLabel.value} unaccounted`
  })

  function splitDestVesselName (vesselUuid: string): string {
    const vessel = vesselMap.value.get(vesselUuid)
    return vessel?.name ?? 'Unknown Vessel'
  }

  // ==================== Blend mode: vessel options ====================

  /** Available occupied vessels for a blend source row, excluding vessels already selected in other rows */
  function blendOccupiedVesselsFor (rowIndex: number): Array<{ title: string, value: string, subtitle: string }> {
    const selectedInOtherRows = new Set(
      blendSources.value
        .filter((_, i) => i !== rowIndex)
        .map(s => s.occupancyUuid)
        .filter(Boolean),
    )
    return occupiedVesselOptions.value.filter(v => !selectedInOtherRows.has(v.value))
  }

  /** Available vessels for blend destination (unoccupied, active, not any source vessel) */
  const blendAvailableVesselOptions = computed(() => {
    const sourceVesselUuids = new Set(
      blendSources.value
        .map(s => {
          const occ = allOccupancies.value.find(o => o.uuid === s.occupancyUuid)
          return occ?.vessel_uuid
        })
        .filter((uuid): uuid is string => !!uuid),
    )
    return allVessels.value
      .filter(v =>
        v.status === 'active'
        && !occupiedVesselUuids.value.has(v.uuid)
        && !sourceVesselUuids.has(v.uuid),
      )
      .map(v => ({
        title: v.name,
        value: v.uuid,
        subtitle: `${formatVesselType(v.type)} · ${v.capacity} ${volumeLabels[v.capacity_unit] ?? v.capacity_unit}`,
      }))
  })

  function addBlendSource () {
    if (blendSources.value.length >= 4) return
    blendSources.value.push({
      occupancyUuid: '',
      amount: '',
      closeAfterBlend: true,
    })
  }

  function removeBlendSource (index: number) {
    if (blendSources.value.length <= 2) return
    blendSources.value.splice(index, 1)
  }

  // ==================== Blend mode: volume math ====================

  function parseBlendSourceAmount (amount: string): number {
    const num = Number.parseFloat(amount)
    return isNaN(num) ? 0 : num
  }

  const blendTotalSourceVolume = computed(() =>
    blendSources.value.reduce((sum, s) => sum + parseBlendSourceAmount(s.amount), 0),
  )

  const blendDestAmount = computed(() =>
    Math.max(0, blendTotalSourceVolume.value - lossAmountNum.value),
  )

  const blendVolumeMathLabel = computed(() => {
    const parts = blendSources.value.map(s =>
      `${parseBlendSourceAmount(s.amount).toFixed(1)}`,
    )
    const lossStr = lossAmountNum.value > 0 ? ` - ${lossAmountNum.value} loss` : ''
    return `${parts.join(' + ')}${lossStr} = ${blendDestAmount.value.toFixed(1)} ${unitLabel.value}`
  })

  // ==================== Blend mode: batch identity ====================

  /** Map occupancy UUID → batch UUID for blend sources */
  const blendSourceBatchUuids = computed(() => {
    const map = new Map<string, string>()
    for (const src of blendSources.value) {
      if (!src.occupancyUuid) continue
      const occ = allOccupancies.value.find(o => o.uuid === src.occupancyUuid)
      if (occ?.batch_uuid) {
        map.set(src.occupancyUuid, occ.batch_uuid)
      }
    }
    return map
  })

  const blendUniqueBatchUuids = computed(() => {
    const uuids = new Set(blendSourceBatchUuids.value.values())
    return [...uuids]
  })

  const blendHasMultipleBatches = computed(() =>
    blendUniqueBatchUuids.value.length > 1,
  )

  const blendBatchOptions = computed(() => {
    return blendUniqueBatchUuids.value.map(uuid => {
      const batch = blendResolvedBatches.value.get(uuid)
      const sourceOcc = blendSources.value.find(s => {
        const occ = allOccupancies.value.find(o => o.uuid === s.occupancyUuid)
        return occ?.batch_uuid === uuid
      })
      const occUuid = sourceOcc?.occupancyUuid ?? ''
      const vessel = blendSourceVesselName(occUuid)
      return {
        uuid,
        label: `${batch?.short_name ?? `Batch ${uuid.slice(0, 8)}`} (from ${vessel})`,
      }
    })
  })

  /** The effective batch UUID for the blend */
  const blendEffectiveBatchUuid = computed(() => {
    if (blendHasMultipleBatches.value && blendSelectedBatchUuid.value) {
      return blendSelectedBatchUuid.value
    }
    // Single batch or default to first
    return blendUniqueBatchUuids.value[0] ?? ''
  })

  const blendReviewBatchName = computed(() => {
    const uuid = blendEffectiveBatchUuid.value
    if (!uuid) return 'Unknown Batch'
    const batch = blendResolvedBatches.value.get(uuid)
    return batch?.short_name ?? `Batch ${uuid.slice(0, 8)}`
  })

  function blendSourceVesselName (occupancyUuid: string): string {
    const occ = allOccupancies.value.find(o => o.uuid === occupancyUuid)
    if (!occ) return 'Unknown Vessel'
    const vessel = vesselMap.value.get(occ.vessel_uuid)
    return vessel?.name ?? 'Unknown Vessel'
  }

  function blendSourceBatchName (occupancyUuid: string): string {
    const occ = allOccupancies.value.find(o => o.uuid === occupancyUuid)
    if (!occ?.batch_uuid) return 'Unknown Batch'
    const batch = blendResolvedBatches.value.get(occ.batch_uuid)
    return batch?.short_name ?? `Batch ${occ.batch_uuid.slice(0, 8)}`
  }

  // ==================== Computed: Review step data ====================

  const reviewSourceVesselName = computed(() =>
    effectiveSourceVessel.value?.name ?? 'Source Vessel',
  )

  const reviewSourceBatchName = computed(() =>
    effectiveSourceBatch.value?.short_name ?? 'Unknown Batch',
  )

  const reviewSourceVolumeLabel = computed(() => {
    const vol = effectiveSourceVolume.value
    if (!vol) return '—'
    return formatVolumePreferred(vol.amount, vol.amount_unit)
  })

  const reviewDestVesselName = computed(() => {
    const vessel = vesselMap.value.get(form.value.destVesselUuid)
    return vessel?.name ?? 'Destination Vessel'
  })

  const transferAmountNum = computed(() => {
    const num = Number.parseFloat(form.value.transferAmount)
    return isNaN(num) ? 0 : num
  })

  const lossAmountNum = computed(() => {
    const num = Number.parseFloat(form.value.lossAmount)
    return isNaN(num) ? 0 : num
  })

  const destReceivesAmount = computed(() =>
    Math.max(0, transferAmountNum.value - lossAmountNum.value),
  )

  const unitLabel = computed(() =>
    volumeLabels[form.value.volumeUnit] ?? form.value.volumeUnit,
  )

  const reviewTransferLabel = computed(() =>
    `${transferAmountNum.value} ${unitLabel.value}`,
  )

  const reviewLossLabel = computed(() => {
    if (lossAmountNum.value <= 0) return ''
    return `${lossAmountNum.value} ${unitLabel.value}`
  })

  const reviewDestReceivesLabel = computed(() =>
    `${destReceivesAmount.value.toFixed(2)} ${unitLabel.value}`,
  )

  const reviewRemainingLabel = computed(() => {
    const sourceVol = effectiveSourceVolume.value
    if (!sourceVol) return '—'
    // Convert source volume to the form's unit for display
    const sourceInFormUnit = convertVolume(sourceVol.amount, sourceVol.amount_unit, form.value.volumeUnit)
    if (sourceInFormUnit === null) return '—'
    const remaining = Math.max(0, sourceInFormUnit - transferAmountNum.value)
    return `${remaining.toFixed(2)} ${unitLabel.value}`
  })

  // ==================== Computed: Validation ====================

  const canProceedToReview = computed(() => {
    if (activeMode.value === 'transfer') {
      return canProceedTransfer.value
    }
    if (activeMode.value === 'split') {
      return canProceedSplit.value
    }
    if (activeMode.value === 'blend') {
      return canProceedBlend.value
    }
    return false
  })

  const canProceedTransfer = computed(() => {
    const hasSource = !!effectiveSourceOccupancy.value
    const hasSourceVolume = !!effectiveSourceVolume.value
    const hasDest = !!form.value.destVesselUuid
    const amount = Number.parseFloat(form.value.transferAmount)
    const hasValidAmount = !isNaN(amount) && amount > 0
    const loss = Number.parseFloat(form.value.lossAmount || '0')
    const hasValidLoss = !isNaN(loss) && loss >= 0
    const hasStatus = !!form.value.destStatus
    return hasSource && hasSourceVolume && hasDest && hasValidAmount && hasValidLoss && hasStatus
  })

  const canProceedSplit = computed(() => {
    const hasSource = !!effectiveSourceOccupancy.value
    const hasSourceVolume = !!effectiveSourceVolume.value
    const loss = Number.parseFloat(form.value.lossAmount || '0')
    const hasValidLoss = !isNaN(loss) && loss >= 0

    // All destinations must have vessel, valid amount, and status
    const allDestsValid = splitDestinations.value.length >= 2
      && splitDestinations.value.every(d => {
        const amt = Number.parseFloat(d.amount)
        return !!d.vesselUuid && !isNaN(amt) && amt > 0 && !!d.status
      })

    return hasSource && hasSourceVolume && hasValidLoss && allDestsValid
  })

  const canProceedBlend = computed(() => {
    const hasDest = !!form.value.destVesselUuid
    const hasStatus = !!form.value.destStatus
    const loss = Number.parseFloat(form.value.lossAmount || '0')
    const hasValidLoss = !isNaN(loss) && loss >= 0

    // All sources must have occupancy and valid amount
    const allSourcesValid = blendSources.value.length >= 2
      && blendSources.value.every(s => {
        const amt = Number.parseFloat(s.amount)
        return !!s.occupancyUuid && !isNaN(amt) && amt > 0
      })

    // If multiple batches, must have selected one
    const batchSelected = !blendHasMultipleBatches.value || !!blendSelectedBatchUuid.value

    return hasDest && hasStatus && hasValidLoss && allSourcesValid && batchSelected
  })

  // ==================== Smart defaults for destination status ====================

  function getDefaultDestStatus (sourceStatus: OccupancyStatus | string | null): OccupancyStatus {
    switch (sourceStatus) {
      case 'fermenting': { return 'conditioning'
      }
      case 'conditioning': { return 'carbonating'
      }
      case 'cold_crashing': { return 'conditioning'
      }
      case 'dry_hopping': { return 'conditioning'
      }
      default: { return 'holding'
      }
    }
  }

  // ==================== Watch: Dynamic source selection ====================

  watch(
    () => form.value.sourceOccupancyUuid,
    async newUuid => {
      if (!newUuid || props.sourceOccupancy) return

      // Find the occupancy from loaded data
      const occ = allOccupancies.value.find(o => o.uuid === newUuid)
      if (!occ) return

      resolvedSourceOccupancy.value = occ

      // Resolve vessel
      resolvedSourceVessel.value = vesselMap.value.get(occ.vessel_uuid) ?? null

      // Set smart default for dest status
      form.value.destStatus = getDefaultDestStatus(occ.status)

      // Resolve batch and volume
      try {
        if (occ.batch_uuid) {
          resolvedSourceBatch.value = await getBatch(occ.batch_uuid)
        }
        if (occ.volume_uuid) {
          resolvedSourceVolume.value = await getVolume(occ.volume_uuid)
          // Pre-fill transfer amount
          const vol = resolvedSourceVolume.value
          if (vol) {
            const converted = convertVolume(vol.amount, vol.amount_unit, form.value.volumeUnit)
            if (converted !== null) {
              form.value.transferAmount = converted.toFixed(2)
            }
          }
        }
      } catch {
        // Non-critical — user can still fill in manually
      }
    },
  )

  // ==================== Watch: Blend source selection (resolve batch/volume data) ====================

  watch(
    () => blendSources.value.map(s => s.occupancyUuid),
    async uuids => {
      for (const uuid of uuids) {
        if (!uuid) continue
        const occ = allOccupancies.value.find(o => o.uuid === uuid)
        if (!occ) continue

        // Resolve batch if not already resolved
        if (occ.batch_uuid && !blendResolvedBatches.value.has(occ.batch_uuid)) {
          try {
            const batch = await getBatch(occ.batch_uuid)
            const updated = new Map(blendResolvedBatches.value)
            updated.set(occ.batch_uuid, batch)
            blendResolvedBatches.value = updated
          } catch {
            // Non-critical
          }
        }

        // Resolve volume and pre-fill amount if not already set
        if (occ.volume_uuid && !blendResolvedVolumes.value.has(occ.volume_uuid)) {
          try {
            const vol = await getVolume(occ.volume_uuid)
            const updated = new Map(blendResolvedVolumes.value)
            updated.set(occ.volume_uuid, vol)
            blendResolvedVolumes.value = updated

            // Pre-fill amount for this source
            const src = blendSources.value.find(s => s.occupancyUuid === uuid)
            if (src && !src.amount) {
              const converted = convertVolume(vol.amount, vol.amount_unit, form.value.volumeUnit)
              if (converted !== null) {
                src.amount = converted.toFixed(2)
              }
            }
          } catch {
            // Non-critical
          }
        }
      }

      // Auto-select first batch for blend if not yet selected
      if (!blendSelectedBatchUuid.value && blendUniqueBatchUuids.value.length > 0) {
        blendSelectedBatchUuid.value = blendUniqueBatchUuids.value[0]!
      }
    },
    { deep: true },
  )

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
    showTimePicker.value = false
    resolvedSourceOccupancy.value = null
    resolvedSourceVessel.value = null
    resolvedSourceBatch.value = null
    resolvedSourceVolume.value = null

    // Set mode from prop
    activeMode.value = props.mode ?? 'transfer'

    // Reset form
    form.value = {
      sourceOccupancyUuid: '',
      destVesselUuid: '',
      transferAmount: '',
      lossAmount: '0',
      volumeUnit: preferences.value.volume,
      destStatus: getDefaultDestStatus(props.sourceOccupancy?.status ?? null),
      closeSource: true,
      transferDate: nowInputValue(),
    }

    // Pre-fill transfer amount from source volume
    if (props.sourceVolume) {
      const converted = convertVolume(
        props.sourceVolume.amount,
        props.sourceVolume.amount_unit,
        preferences.value.volume,
      )
      if (converted !== null) {
        form.value.transferAmount = converted.toFixed(2)
      }
    }

    // Reset split destinations
    splitDestinations.value = [
      {
        vesselUuid: '',
        amount: '',
        status: getDefaultDestStatus(props.sourceOccupancy?.status ?? null),
      },
      {
        vesselUuid: '',
        amount: '',
        status: getDefaultDestStatus(props.sourceOccupancy?.status ?? null),
      },
    ]

    // Reset blend sources
    blendResolvedBatches.value = new Map()
    blendResolvedVolumes.value = new Map()
    blendSelectedBatchUuid.value = ''

    // For blend mode, pre-populate first source if we have a source occupancy
    if (props.sourceOccupancy) {
      blendSources.value = [
        {
          occupancyUuid: props.sourceOccupancy.uuid,
          amount: form.value.transferAmount,
          closeAfterBlend: true,
        },
        {
          occupancyUuid: '',
          amount: '',
          closeAfterBlend: true,
        },
      ]

      // Pre-resolve batch for the source occupancy
      if (props.sourceBatch && props.sourceOccupancy.batch_uuid) {
        const updated = new Map(blendResolvedBatches.value)
        updated.set(props.sourceOccupancy.batch_uuid, props.sourceBatch)
        blendResolvedBatches.value = updated
        blendSelectedBatchUuid.value = props.sourceOccupancy.batch_uuid
      }
      if (props.sourceVolume && props.sourceOccupancy.volume_uuid) {
        const updated = new Map(blendResolvedVolumes.value)
        updated.set(props.sourceOccupancy.volume_uuid, props.sourceVolume)
        blendResolvedVolumes.value = updated
      }
    } else {
      blendSources.value = [
        { occupancyUuid: '', amount: '', closeAfterBlend: true },
        { occupancyUuid: '', amount: '', closeAfterBlend: true },
      ]
    }

    // Load reference data
    loadingData.value = true
    try {
      const [occupancyData, vesselData] = await Promise.all([
        getActiveOccupancies(),
        getVessels(),
      ])
      allOccupancies.value = occupancyData
      allVessels.value = vesselData
    } catch {
      loadError.value = 'Failed to load vessel data. Please close and try again.'
    } finally {
      loadingData.value = false
    }
  }

  // ==================== Actions ====================

  function handleClose () {
    emit('update:modelValue', false)
  }

  async function handleConfirm () {
    switch (activeMode.value) {
      case 'transfer': {
        await handleConfirmTransfer()

        break
      }
      case 'split': {
        await handleConfirmSplit()

        break
      }
      case 'blend': {
        await handleConfirmBlend()

        break
      }
    // No default
    }
  }

  async function handleConfirmTransfer () {
    const sourceOcc = effectiveSourceOccupancy.value
    const sourceVol = effectiveSourceVolume.value

    if (!sourceOcc || !sourceVol) {
      saveError.value = 'Missing source occupancy or volume data'
      return
    }

    saving.value = true
    saveError.value = ''

    try {
      const transferAmount = Number.parseFloat(form.value.transferAmount)
      const lossAmount = Number.parseFloat(form.value.lossAmount || '0')
      const transferDate = form.value.transferDate
        ? new Date(form.value.transferDate).toISOString()
        : new Date().toISOString()

      await createTransfer({
        source_occupancy_uuid: sourceOcc.uuid,
        dest_vessel_uuid: form.value.destVesselUuid,
        volume_uuid: sourceVol.uuid,
        amount: transferAmount,
        amount_unit: form.value.volumeUnit,
        loss_amount: lossAmount > 0 ? lossAmount : undefined,
        loss_unit: lossAmount > 0 ? form.value.volumeUnit : undefined,
        started_at: transferDate,
        close_source: form.value.closeSource,
        dest_status: form.value.destStatus,
      })

      showNotice('Transfer complete')
      emit('transferred')
      emit('update:modelValue', false)
    } catch (error) {
      saveError.value = error instanceof Error ? error.message : 'Failed to create transfer'
    } finally {
      saving.value = false
    }
  }

  async function handleConfirmSplit () {
    const sourceOcc = effectiveSourceOccupancy.value
    const sourceVol = effectiveSourceVolume.value
    const sourceBatchData = effectiveSourceBatch.value

    if (!sourceOcc || !sourceVol) {
      saveError.value = 'Missing source occupancy or volume data'
      return
    }

    saving.value = true
    saveError.value = ''

    const transferDate = form.value.transferDate
      ? new Date(form.value.transferDate).toISOString()
      : new Date().toISOString()

    try {
      // Step 1: Create child volumes for each destination
      const childVolumes = await Promise.all(
        splitDestinations.value.map((dest, index) =>
          createVolume({
            name: `${sourceVol.name ?? 'Volume'} - Split ${index + 1}`,
            description: `Split from ${sourceVol.name ?? 'source volume'}`,
            amount: parseSplitDestAmount(dest.amount),
            amount_unit: form.value.volumeUnit,
          }),
        ),
      )

      // Step 2: Create transfers for each destination
      for (let i = 0; i < splitDestinations.value.length; i++) {
        const dest = splitDestinations.value[i]!
        const childVol = childVolumes[i]!
        const isLast = i === splitDestinations.value.length - 1
        const lossAmount = Number.parseFloat(form.value.lossAmount || '0')

        try {
          await createTransfer({
            source_occupancy_uuid: sourceOcc.uuid,
            dest_vessel_uuid: dest.vesselUuid,
            volume_uuid: childVol.uuid,
            amount: parseSplitDestAmount(dest.amount),
            amount_unit: form.value.volumeUnit,
            loss_amount: isLast && lossAmount > 0 ? lossAmount : undefined,
            loss_unit: isLast && lossAmount > 0 ? form.value.volumeUnit : undefined,
            started_at: transferDate,
            close_source: isLast ? form.value.closeSource : false,
            dest_status: dest.status,
          })
        } catch (error) {
          const vesselName = splitDestVesselName(dest.vesselUuid)
          const msg = error instanceof Error ? error.message : 'Unknown error'
          saveError.value = `Transfer to ${vesselName} failed: ${msg}. Previous transfers may have succeeded.`
          saving.value = false
          return
        }
      }

      // Step 3: Create volume relations
      for (const [i, childVolume] of childVolumes.entries()) {
        const childVol = childVolume!
        const dest = splitDestinations.value[i]!
        try {
          await createVolumeRelation({
            parent_volume_uuid: sourceVol.uuid,
            child_volume_uuid: childVol.uuid,
            relation_type: 'split',
            amount: parseSplitDestAmount(dest.amount),
            amount_unit: form.value.volumeUnit,
          })
        } catch {
          // Non-critical — relation is metadata, transfers already succeeded
        }
      }

      // Step 4: Create batch-volume records
      if (sourceBatchData) {
        for (const childVol of childVolumes) {
          try {
            await createBatchVolume({
              batch_uuid: sourceBatchData.uuid,
              volume_uuid: childVol.uuid,
              liquid_phase: 'beer',
            })
          } catch {
            // Non-critical
          }
        }
      }

      showNotice('Split complete')
      emit('transferred')
      emit('update:modelValue', false)
    } catch (error) {
      saveError.value = error instanceof Error ? error.message : 'Failed to create split'
    } finally {
      saving.value = false
    }
  }

  async function handleConfirmBlend () {
    if (blendSources.value.length < 2) {
      saveError.value = 'At least 2 sources are required for a blend'
      return
    }

    saving.value = true
    saveError.value = ''

    const transferDate = form.value.transferDate
      ? new Date(form.value.transferDate).toISOString()
      : new Date().toISOString()

    const batchUuid = blendEffectiveBatchUuid.value

    try {
      // Step 1: Create blended child volume
      const blendedVolume = await createVolume({
        name: `Blend - ${blendReviewBatchName.value}`,
        description: `Blended from ${blendSources.value.map(s => blendSourceVesselName(s.occupancyUuid)).join(', ')}`,
        amount: blendDestAmount.value,
        amount_unit: form.value.volumeUnit,
      })

      // Step 2: Create transfers from each source to the destination
      for (let i = 0; i < blendSources.value.length; i++) {
        const src = blendSources.value[i]!
        const lossAmount = Number.parseFloat(form.value.lossAmount || '0')
        const isFirst = i === 0

        try {
          await createTransfer({
            source_occupancy_uuid: src.occupancyUuid,
            dest_vessel_uuid: form.value.destVesselUuid,
            volume_uuid: blendedVolume.uuid,
            amount: parseBlendSourceAmount(src.amount),
            amount_unit: form.value.volumeUnit,
            loss_amount: isFirst && lossAmount > 0 ? lossAmount : undefined,
            loss_unit: isFirst && lossAmount > 0 ? form.value.volumeUnit : undefined,
            started_at: transferDate,
            close_source: src.closeAfterBlend,
            dest_status: form.value.destStatus,
          })
        } catch (error) {
          const vesselName = blendSourceVesselName(src.occupancyUuid)
          const msg = error instanceof Error ? error.message : 'Unknown error'
          saveError.value = `Transfer from ${vesselName} failed: ${msg}. Previous transfers may have succeeded.`
          saving.value = false
          return
        }
      }

      // Step 3: Create volume relations (each source → blended child)
      for (const src of blendSources.value) {
        const occ = allOccupancies.value.find(o => o.uuid === src.occupancyUuid)
        if (!occ?.volume_uuid) continue
        try {
          await createVolumeRelation({
            parent_volume_uuid: occ.volume_uuid,
            child_volume_uuid: blendedVolume.uuid,
            relation_type: 'blend',
            amount: parseBlendSourceAmount(src.amount),
            amount_unit: form.value.volumeUnit,
          })
        } catch {
          // Non-critical
        }
      }

      // Step 4: Create batch-volume record
      if (batchUuid) {
        try {
          await createBatchVolume({
            batch_uuid: batchUuid,
            volume_uuid: blendedVolume.uuid,
            liquid_phase: 'beer',
          })
        } catch {
          // Non-critical
        }
      }

      showNotice('Blend complete')
      emit('transferred')
      emit('update:modelValue', false)
    } catch (error) {
      saveError.value = error instanceof Error ? error.message : 'Failed to create blend'
    } finally {
      saving.value = false
    }
  }
</script>

<style scoped>
.transfer-flow-diagram {
  max-width: 360px;
  text-align: center;
  padding: 16px;
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  border-radius: 8px;
  background: rgba(var(--v-theme-surface), 0.5);
}

.flow-source,
.flow-dest {
  padding: 8px 0;
}

.flow-arrow-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 4px 0;
}

.flow-line {
  width: 2px;
  height: 12px;
  background: rgba(var(--v-theme-on-surface), 0.3);
}

.flow-label {
  padding: 4px 0;
}

.flow-arrow-icon {
  color: rgba(var(--v-theme-primary), 1);
}

.transfer-summary-list {
  list-style-type: disc;
}

.transfer-summary-list li {
  margin-bottom: 4px;
}

/* Split review styles */
.split-flow-branches {
  text-align: left;
  padding: 8px 0 8px 16px;
}

.split-flow-branch {
  padding: 2px 0;
}

.split-branch-connector {
  font-family: monospace;
  margin-right: 4px;
}

/* Blend review styles */
.blend-flow-source {
  text-align: left;
  padding: 2px 0;
}

.blend-flow-dest {
  text-align: left;
}

.blend-flow-loss {
  text-align: left;
}
</style>
