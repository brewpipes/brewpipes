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
          <v-card-actions class="justify-end">
            <v-btn
              size="small"
              variant="text"
              prepend-icon="mdi-upload"
              @click="bulkImportDialog = true"
            >
              Bulk import
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-beaker-outline" />
            {{ selectedBatch ? selectedBatch.short_name : 'Batch details' }}
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
                          <div class="spark-card" style="--spark-color: var(--v-theme-info)">
                            <div class="spark-meta">
                              <div class="text-caption text-medium-emphasis">Temp</div>
                              <div class="text-body-2 font-weight-medium">
                                {{ temperatureSeries.latestLabel }}
                              </div>
                            </div>
                            <div class="spark-chart">
                              <svg
                                v-if="temperatureSeries.values.length"
                                :viewBox="`0 0 ${sparklineWidth} ${sparklineHeight}`"
                                preserveAspectRatio="none"
                              >
                                <path class="spark-area" :d="temperatureSeries.areaPath" />
                                <path class="spark-line" :d="temperatureSeries.linePath" />
                              </svg>
                              <div v-else class="spark-placeholder">No readings</div>
                            </div>
                          </div>
                        </v-col>
                        <v-col cols="12" md="4">
                          <div class="spark-card" style="--spark-color: var(--v-theme-secondary)">
                            <div class="spark-meta">
                              <div class="text-caption text-medium-emphasis">Gravity</div>
                              <div class="text-body-2 font-weight-medium">
                                {{ gravitySeries.latestLabel }}
                              </div>
                            </div>
                            <div class="spark-chart">
                              <svg
                                v-if="gravitySeries.values.length"
                                :viewBox="`0 0 ${sparklineWidth} ${sparklineHeight}`"
                                preserveAspectRatio="none"
                              >
                                <path class="spark-area" :d="gravitySeries.areaPath" />
                                <path class="spark-line" :d="gravitySeries.linePath" />
                              </svg>
                              <div v-else class="spark-placeholder">No readings</div>
                            </div>
                          </div>
                        </v-col>
                        <v-col cols="12" md="4">
                          <div class="spark-card" style="--spark-color: var(--v-theme-warning)">
                            <div class="spark-meta">
                              <div class="text-caption text-medium-emphasis">pH</div>
                              <div class="text-body-2 font-weight-medium">
                                {{ phSeries.latestLabel }}
                              </div>
                            </div>
                            <div class="spark-chart">
                              <svg
                                v-if="phSeries.values.length"
                                :viewBox="`0 0 ${sparklineWidth} ${sparklineHeight}`"
                                preserveAspectRatio="none"
                              >
                                <path class="spark-area" :d="phSeries.areaPath" />
                                <path class="spark-line" :d="phSeries.linePath" />
                              </svg>
                              <div v-else class="spark-placeholder">No readings</div>
                            </div>
                          </div>
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
                  <v-card class="sub-card" variant="outlined">
                    <v-card-text>
                      <v-progress-linear
                        v-if="batchSummaryLoading"
                        color="primary"
                        indeterminate
                        class="mb-4"
                      />

                      <template v-if="batchSummary">
                        <!-- Recipe & Style -->
                        <v-row class="mb-4">
                          <v-col cols="12" md="6">
                            <div class="text-overline text-medium-emphasis">Recipe</div>
                            <div class="text-h6">{{ batchSummary.recipe_name ?? 'Not specified' }}</div>
                          </v-col>
                          <v-col cols="12" md="6">
                            <div class="text-overline text-medium-emphasis">Style</div>
                            <div class="text-h6">{{ batchSummary.style_name ?? 'Not specified' }}</div>
                          </v-col>
                        </v-row>

                        <v-divider class="mb-4" />

                        <!-- Status & Location -->
                        <v-row class="mb-4">
                          <v-col cols="12" md="4">
                            <div class="text-overline text-medium-emphasis">Current Phase</div>
                            <v-chip
                              v-if="batchSummary.current_phase"
                              :color="getPhaseColor(batchSummary.current_phase)"
                              variant="tonal"
                              class="mt-1"
                            >
                              {{ formatPhase(batchSummary.current_phase) }}
                            </v-chip>
                            <div v-else class="text-body-2 text-medium-emphasis mt-1">Not set</div>
                          </v-col>
                          <v-col cols="12" md="4">
                            <div class="text-overline text-medium-emphasis">Current Vessel</div>
                            <div class="text-body-1 font-weight-medium mt-1">
                              {{ batchSummary.current_vessel ?? 'Not assigned' }}
                            </div>
                          </v-col>
                          <v-col cols="12" md="4">
                            <div class="text-overline text-medium-emphasis">Occupancy Status</div>
                            <v-menu v-if="batchSummary.current_occupancy_id" location="bottom">
                              <template #activator="{ props }">
                                <v-chip
                                  v-bind="props"
                                  :color="getOccupancyStatusColor(batchSummary.current_occupancy_status)"
                                  variant="tonal"
                                  size="small"
                                  class="mt-1 cursor-pointer"
                                  append-icon="mdi-menu-down"
                                >
                                  {{ formatOccupancyStatus(batchSummary.current_occupancy_status) }}
                                </v-chip>
                              </template>
                              <v-list density="compact" nav>
                                <v-list-subheader>Change status</v-list-subheader>
                                <v-list-item
                                  v-for="statusOption in occupancyStatusOptions"
                                  :key="statusOption.value"
                                  :active="statusOption.value === batchSummary.current_occupancy_status"
                                  @click="changeOccupancyStatus(batchSummary.current_occupancy_id!, statusOption.value)"
                                >
                                  <template #prepend>
                                    <v-avatar
                                      :color="getOccupancyStatusColor(statusOption.value)"
                                      size="24"
                                      class="mr-2"
                                    >
                                      <v-icon :icon="getOccupancyStatusIcon(statusOption.value)" size="14" />
                                    </v-avatar>
                                  </template>
                                  <v-list-item-title>{{ statusOption.title }}</v-list-item-title>
                                </v-list-item>
                              </v-list>
                            </v-menu>
                            <div v-else class="text-body-2 text-medium-emphasis mt-1">Not set</div>
                          </v-col>
                        </v-row>

                        <v-divider class="mb-4" />

                        <!-- Gravity & ABV Metrics -->
                        <div class="text-overline text-medium-emphasis mb-2">Gravity & ABV</div>
                        <v-row class="mb-4">
                          <v-col cols="6" md="3">
                            <div class="metric-card">
                              <div class="metric-label">OG</div>
                              <div class="metric-value">
                                {{ formatGravityPreferred(batchSummary.original_gravity, 'sg') }}
                              </div>
                            </div>
                          </v-col>
                          <v-col cols="6" md="3">
                            <div class="metric-card">
                              <div class="metric-label">FG</div>
                              <div class="metric-value">
                                {{ formatGravityPreferred(batchSummary.final_gravity, 'sg') }}
                              </div>
                            </div>
                          </v-col>
                          <v-col cols="6" md="3">
                            <div class="metric-card">
                              <div class="metric-label">ABV</div>
                              <div class="metric-value d-flex align-center ga-1">
                                {{ formatPercent(batchSummary.abv, 1) }}
                                <v-chip
                                  v-if="batchSummary.abv !== null && batchSummary.abv !== undefined && batchSummary.abv_calculated"
                                  size="x-small"
                                  color="info"
                                  variant="tonal"
                                >
                                  calc
                                </v-chip>
                              </div>
                            </div>
                          </v-col>
                          <v-col cols="6" md="3">
                            <div class="metric-card">
                              <div class="metric-label">IBU</div>
                              <div class="metric-value">
                                {{ batchSummary.ibu ?? '—' }}
                              </div>
                            </div>
                          </v-col>
                        </v-row>

                        <!-- Duration Metrics -->
                        <div class="text-overline text-medium-emphasis mb-2">Duration</div>
                        <v-row class="mb-4">
                          <v-col cols="6" md="4">
                            <div class="metric-card">
                              <div class="metric-label">Days in FV</div>
                              <div class="metric-value">
                                {{ formatDays(batchSummary.days_in_fermenter) }}
                              </div>
                            </div>
                          </v-col>
                          <v-col cols="6" md="4">
                            <div class="metric-card">
                              <div class="metric-label">Days in Brite</div>
                              <div class="metric-value">
                                {{ formatDays(batchSummary.days_in_brite) }}
                              </div>
                            </div>
                          </v-col>
                          <v-col cols="12" md="4">
                            <div class="metric-card metric-card--highlight">
                              <div class="metric-label">Grain to Glass</div>
                              <div class="metric-value">
                                {{ formatDays(batchSummary.days_grain_to_glass) }}
                              </div>
                            </div>
                          </v-col>
                        </v-row>

                        <!-- Volume Metrics -->
                        <div class="text-overline text-medium-emphasis mb-2">Volume & Loss</div>
                        <v-row class="mb-4">
                          <v-col cols="6" md="3">
                            <div class="metric-card">
                              <div class="metric-label">Starting</div>
                              <div class="metric-value">
                                {{ formatVolumePreferred(batchSummary.starting_volume_bbl, 'bbl') }}
                              </div>
                            </div>
                          </v-col>
                          <v-col cols="6" md="3">
                            <div class="metric-card">
                              <div class="metric-label">Current</div>
                              <div class="metric-value">
                                {{ formatVolumePreferred(batchSummary.current_volume_bbl, 'bbl') }}
                              </div>
                            </div>
                          </v-col>
                          <v-col cols="6" md="3">
                            <div class="metric-card">
                              <div class="metric-label">Total Loss</div>
                              <div class="metric-value">
                                {{ formatVolumePreferred(batchSummary.total_loss_bbl, 'bbl') }}
                              </div>
                            </div>
                          </v-col>
                          <v-col cols="6" md="3">
                            <div class="metric-card" :class="{ 'metric-card--warning': (batchSummary.loss_percentage ?? 0) > 10 }">
                              <div class="metric-label">Loss %</div>
                              <div class="metric-value">
                                {{ formatPercent(batchSummary.loss_percentage, 1) }}
                              </div>
                            </div>
                          </v-col>
                        </v-row>

                        <v-divider class="mb-4" />

                        <!-- Brew Sessions List -->
                        <div class="text-overline text-medium-emphasis mb-2">Brew Sessions</div>
                        <v-list
                          v-if="batchSummary.brew_sessions && batchSummary.brew_sessions.length > 0"
                          class="brew-session-summary-list"
                          density="compact"
                          variant="tonal"
                        >
                          <v-list-item
                            v-for="session in batchSummary.brew_sessions"
                            :key="session.id"
                          >
                            <template #prepend>
                              <v-icon icon="mdi-kettle-steam" size="small" />
                            </template>
                            <v-list-item-title>
                              {{ formatDateTime(session.brewed_at) }}
                            </v-list-item-title>
                            <v-list-item-subtitle v-if="session.notes">
                              {{ session.notes }}
                            </v-list-item-subtitle>
                          </v-list-item>
                        </v-list>
                        <v-alert
                          v-else
                          density="compact"
                          type="info"
                          variant="tonal"
                        >
                          No brew sessions recorded yet.
                        </v-alert>

                        <!-- Notes -->
                        <template v-if="batchSummary.notes">
                          <v-divider class="my-4" />
                          <div class="text-overline text-medium-emphasis mb-2">Notes</div>
                          <div class="text-body-2">{{ batchSummary.notes }}</div>
                        </template>
                      </template>

                      <v-alert
                        v-else-if="!batchSummaryLoading"
                        density="comfortable"
                        type="info"
                        variant="tonal"
                      >
                        Summary data not available.
                      </v-alert>
                    </v-card-text>
                  </v-card>
                </v-window-item>

                <v-window-item value="brew-sessions">
                  <v-row>
                    <v-col cols="12">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1 d-flex align-center">
                          Brew Sessions
                          <v-spacer />
                          <v-btn
                            icon="mdi-plus"
                            size="small"
                            variant="text"
                            aria-label="Add brew session"
                            @click="openCreateBrewSessionDialog"
                          />
                        </v-card-title>
                        <v-card-text>
                          <v-alert
                            v-if="brewSessions.length === 0"
                            density="compact"
                            type="info"
                            variant="tonal"
                            class="mb-3"
                          >
                            No brew sessions recorded. Add a brew session to track hot-side additions and measurements.
                          </v-alert>
                          <v-list v-else class="brew-session-list" lines="three">
                            <v-list-item
                              v-for="session in brewSessionsSorted"
                              :key="session.id"
                              :active="session.id === selectedBrewSessionId"
                              @click="selectBrewSession(session.id)"
                            >
                              <v-list-item-title>
                                {{ formatDateTime(session.brewed_at) }}
                              </v-list-item-title>
                              <v-list-item-subtitle>
                                <span v-if="getVesselName(session.mash_vessel_id)">
                                  Mash: {{ getVesselName(session.mash_vessel_id) }}
                                </span>
                                <span v-if="getVesselName(session.boil_vessel_id)">
                                  &bull; Boil: {{ getVesselName(session.boil_vessel_id) }}
                                </span>
                                <span v-if="getVolumeName(session.wort_volume_id)">
                                  &bull; {{ getVolumeName(session.wort_volume_id) }}
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
                                  @click.stop="openEditBrewSessionDialog(session)"
                                />
                              </template>
                            </v-list-item>
                          </v-list>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>

                  <!-- Selected Brew Session Details -->
                  <v-row v-if="selectedBrewSession" class="mt-4">
                    <v-col cols="12">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1 d-flex align-center">
                          <v-icon class="mr-2" icon="mdi-kettle-steam" size="small" />
                          {{ formatDateTime(selectedBrewSession.brewed_at) }}
                          <v-spacer />
                          <v-btn size="small" variant="text" @click="clearBrewSessionSelection">
                            Clear
                          </v-btn>
                        </v-card-title>
                        <v-card-text>
                          <v-row dense>
                            <v-col v-if="getVesselName(selectedBrewSession.mash_vessel_id)" cols="12" md="4">
                              <div class="text-caption text-medium-emphasis">Mash Vessel</div>
                              <div class="text-body-2 font-weight-medium">
                                {{ getVesselName(selectedBrewSession.mash_vessel_id) }}
                              </div>
                            </v-col>
                            <v-col v-if="getVesselName(selectedBrewSession.boil_vessel_id)" cols="12" md="4">
                              <div class="text-caption text-medium-emphasis">Boil Vessel</div>
                              <div class="text-body-2 font-weight-medium">
                                {{ getVesselName(selectedBrewSession.boil_vessel_id) }}
                              </div>
                            </v-col>
                            <v-col v-if="getVolumeName(selectedBrewSession.wort_volume_id)" cols="12" md="4">
                              <div class="text-caption text-medium-emphasis">Wort Volume</div>
                              <div class="text-body-2 font-weight-medium">
                                {{ getVolumeName(selectedBrewSession.wort_volume_id) }}
                                <span v-if="getVolumeAmount(selectedBrewSession.wort_volume_id)" class="text-medium-emphasis">
                                  ({{ getVolumeAmount(selectedBrewSession.wort_volume_id) }})
                                </span>
                              </div>
                            </v-col>
                          </v-row>
                          <div v-if="selectedBrewSession.notes" class="mt-3">
                            <div class="text-caption text-medium-emphasis">Notes</div>
                            <div class="text-body-2">{{ selectedBrewSession.notes }}</div>
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
                            icon="mdi-plus"
                            size="x-small"
                            variant="text"
                            aria-label="Add hot-side addition"
                            :disabled="!selectedBrewSession.wort_volume_id"
                            @click="openCreateHotSideAdditionDialog"
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
                              <tr v-for="addition in wortAdditionsSorted" :key="addition.id">
                                <td>
                                  <v-chip size="x-small" variant="tonal">{{ addition.addition_type }}</v-chip>
                                  <span v-if="addition.stage" class="text-medium-emphasis ml-1">{{ addition.stage }}</span>
                                </td>
                                <td>{{ formatAmount(addition.amount, addition.amount_unit) }}</td>
                                <td>{{ formatDateTime(addition.added_at) }}</td>
                              </tr>
                              <tr v-if="wortAdditionsSorted.length === 0">
                                <td colspan="3" class="text-medium-emphasis">
                                  {{ selectedBrewSession.wort_volume_id ? 'No additions recorded.' : 'Select a wort volume first.' }}
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
                            icon="mdi-plus"
                            size="x-small"
                            variant="text"
                            aria-label="Add hot-side measurement"
                            :disabled="!selectedBrewSession.wort_volume_id"
                            @click="openCreateHotSideMeasurementDialog"
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
                              <tr v-for="measurement in wortMeasurementsSorted" :key="measurement.id">
                                <td>{{ formatMeasurementKind(measurement.kind) }}</td>
                                <td>{{ formatValue(measurement.value, measurement.unit) }}</td>
                                <td>{{ formatDateTime(measurement.observed_at) }}</td>
                              </tr>
                              <tr v-if="wortMeasurementsSorted.length === 0">
                                <td colspan="3" class="text-medium-emphasis">
                                  {{ selectedBrewSession.wort_volume_id ? 'No measurements recorded.' : 'Select a wort volume first.' }}
                                </td>
                              </tr>
                            </tbody>
                          </v-table>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-window-item>

                <v-window-item value="additions">
                  <v-row>
                    <v-col cols="12">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1 d-flex align-center">
                          Addition log
                          <v-spacer />
                          <v-btn
                            icon="mdi-plus"
                            size="small"
                            variant="text"
                            aria-label="Record addition"
                            @click="createAdditionDialog = true"
                          />
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
                    <v-col cols="12">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1 d-flex align-center">
                          Measurement log
                          <v-spacer />
                          <v-btn
                            icon="mdi-plus"
                            size="small"
                            variant="text"
                            aria-label="Record measurement"
                            @click="createMeasurementDialog = true"
                          />
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
                                :placeholder="preferences.temperature === 'f' ? '67F' : '19C'"
                                inputmode="decimal"
                              />
                            </v-col>
                            <v-col cols="12" md="2">
                              <v-text-field
                                v-model="timelineReading.gravity"
                                density="compact"
                                label="Gravity"
                                :placeholder="preferences.gravity === 'sg' ? '1.056' : '13.8'"
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
        <v-autocomplete
          v-model="newBatch.recipe_id"
          :items="recipeSelectItems"
          :loading="recipesLoading"
          clearable
          density="comfortable"
          hint="Optional - link this batch to a recipe"
          item-title="title"
          item-value="value"
          label="Recipe"
          persistent-hint
        >
          <template #item="{ props, item }">
            <v-list-item v-bind="props">
              <template #subtitle>
                <span v-if="item.raw.style">{{ item.raw.style }}</span>
              </template>
            </v-list-item>
          </template>
        </v-autocomplete>
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

  <v-dialog v-model="bulkImportDialog" max-width="720">
    <v-card>
      <v-card-title class="text-h6">Bulk import batches</v-card-title>
      <v-card-text>
        <v-alert density="compact" type="info" variant="tonal" class="mb-4">
          Expected columns: short_name (required), brew_date (optional), notes (optional).
        </v-alert>

        <v-file-input
          v-model="importFile"
          accept=".csv,text/csv"
          density="comfortable"
          label="CSV file"
          placeholder="Select batch import CSV"
          prepend-icon="mdi-file-delimited-outline"
          show-size
        />

        <div class="d-flex align-center justify-space-between flex-wrap ga-2 mb-4">
          <div class="text-caption text-medium-emphasis">
            Brew date format: YYYY-MM-DD.
          </div>
          <v-btn size="small" variant="text" prepend-icon="mdi-download" @click="downloadBatchTemplate">
            Download template
          </v-btn>
        </div>

        <v-alert
          v-if="importSummary"
          class="mb-3"
          density="compact"
          :type="importSummary.type"
          variant="tonal"
        >
          {{ importSummary.message }}
        </v-alert>

        <v-table v-if="importErrors.length" class="data-table" density="compact">
          <thead>
            <tr>
              <th>Row</th>
              <th>Issue</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(rowError, index) in importErrors" :key="index">
              <td>{{ rowError.row ?? '-' }}</td>
              <td>{{ rowError.message }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" :disabled="importUploading" @click="bulkImportDialog = false">
          Close
        </v-btn>
        <v-btn
          color="primary"
          :loading="importUploading"
          :disabled="!canImport"
          @click="uploadBatchImport"
        >
          Upload
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="createAdditionDialog" max-width="520">
    <v-card>
      <v-card-title class="text-h6">Record addition</v-card-title>
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
        <v-select v-model="additionForm.amount_unit" :items="unitOptions" label="Unit" />
        <v-text-field v-model="additionForm.added_at" label="Added at" type="datetime-local" />
        <v-textarea v-model="additionForm.notes" auto-grow label="Notes" rows="2" />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="createAdditionDialog = false">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="
            !additionForm.amount ||
            (additionForm.target === 'occupancy' && !additionForm.occupancy_id)
          "
          @click="recordAddition"
        >
          Add addition
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="createMeasurementDialog" max-width="520">
    <v-card>
      <v-card-title class="text-h6">Record measurement</v-card-title>
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
        <v-textarea v-model="measurementForm.notes" auto-grow label="Notes" rows="2" />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="createMeasurementDialog = false">Cancel</v-btn>
        <v-btn
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
              :placeholder="preferences.temperature === 'f' ? 'F' : 'C'"
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
              :placeholder="preferences.gravity === 'sg' ? 'SG' : 'Plato'"
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

  <!-- Create/Edit Brew Session Dialog -->
  <v-dialog v-model="brewSessionDialog" max-width="680" persistent>
    <v-card>
      <v-card-title class="text-h6">
        {{ isEditingBrewSession ? 'Edit brew session' : 'Add brew session' }}
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="brewSessionForm.brewed_at"
              density="comfortable"
              label="Brewed at"
              type="datetime-local"
              :rules="[rules.required]"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-autocomplete
              v-model="brewSessionForm.mash_vessel_id"
              :items="mashVesselOptions"
              clearable
              density="comfortable"
              item-title="name"
              item-value="id"
              label="Mash Vessel"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-autocomplete
              v-model="brewSessionForm.boil_vessel_id"
              :items="boilVesselOptions"
              clearable
              density="comfortable"
              item-title="name"
              item-value="id"
              label="Boil Vessel"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-autocomplete
              v-model="brewSessionForm.wort_volume_id"
              :items="wortVolumeOptions"
              clearable
              density="comfortable"
              item-title="label"
              item-value="id"
              label="Wort Volume"
              hint="Select existing or create new"
              persistent-hint
            >
              <template #no-data>
                <v-list-item>
                  <v-list-item-title>No volumes available</v-list-item-title>
                </v-list-item>
              </template>
              <template #append-item>
                <v-divider class="my-2" />
                <v-list-item @click="openCreateVolumeDialog">
                  <template #prepend>
                    <v-icon icon="mdi-plus" />
                  </template>
                  <v-list-item-title>Create new wort volume</v-list-item-title>
                </v-list-item>
              </template>
            </v-autocomplete>
          </v-col>
          <v-col cols="12">
            <v-textarea
              v-model="brewSessionForm.notes"
              auto-grow
              density="comfortable"
              label="Notes"
              placeholder="Mash temps, boil notes, etc."
              rows="2"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" :disabled="savingBrewSession" @click="closeBrewSessionDialog">
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          :disabled="!isBrewSessionFormValid"
          :loading="savingBrewSession"
          @click="saveBrewSession"
        >
          {{ isEditingBrewSession ? 'Save changes' : 'Add session' }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Create Wort Volume Dialog -->
  <v-dialog v-model="createVolumeDialog" max-width="480" persistent>
    <v-card>
      <v-card-title class="text-h6">Create wort volume</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12">
            <v-text-field
              v-model="volumeForm.name"
              density="comfortable"
              label="Name"
              placeholder="IPA 24-07 Wort"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="volumeForm.amount"
              density="comfortable"
              label="Amount"
              type="number"
              :rules="[rules.required, rules.positiveNumber]"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-select
              v-model="volumeForm.amount_unit"
              :items="volumeUnitOptions"
              density="comfortable"
              label="Unit"
            />
          </v-col>
          <v-col cols="12">
            <v-textarea
              v-model="volumeForm.description"
              auto-grow
              density="comfortable"
              label="Description"
              rows="2"
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" :disabled="savingVolume" @click="createVolumeDialog = false">
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          :disabled="!isVolumeFormValid"
          :loading="savingVolume"
          @click="createWortVolume"
        >
          Create volume
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Hot-Side Addition Dialog -->
  <v-dialog v-model="hotSideAdditionDialog" max-width="520">
    <v-card>
      <v-card-title class="text-h6">Add hot-side addition</v-card-title>
      <v-card-text>
        <v-select
          v-model="hotSideAdditionForm.addition_type"
          :items="additionTypeOptions"
          density="comfortable"
          label="Addition type"
        />
        <v-text-field
          v-model="hotSideAdditionForm.stage"
          density="comfortable"
          label="Stage"
          placeholder="60 min, 15 min, whirlpool"
        />
        <v-text-field
          v-model="hotSideAdditionForm.inventory_lot_uuid"
          density="comfortable"
          label="Inventory lot UUID"
          placeholder="Optional"
        />
        <v-row>
          <v-col cols="8">
            <v-text-field
              v-model="hotSideAdditionForm.amount"
              density="comfortable"
              label="Amount"
              type="number"
            />
          </v-col>
          <v-col cols="4">
            <v-select
              v-model="hotSideAdditionForm.amount_unit"
              :items="volumeUnitOptions"
              density="comfortable"
              label="Unit"
            />
          </v-col>
        </v-row>
        <v-text-field
          v-model="hotSideAdditionForm.added_at"
          density="comfortable"
          label="Added at"
          type="datetime-local"
        />
        <v-textarea
          v-model="hotSideAdditionForm.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="hotSideAdditionDialog = false">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!hotSideAdditionForm.amount"
          :loading="savingHotSideAddition"
          @click="recordHotSideAddition"
        >
          Add addition
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Hot-Side Measurement Dialog -->
  <v-dialog v-model="hotSideMeasurementDialog" max-width="520">
    <v-card>
      <v-card-title class="text-h6">Add hot-side measurement</v-card-title>
      <v-card-text>
        <v-select
          v-model="hotSideMeasurementForm.kind"
          :items="hotSideMeasurementKinds"
          density="comfortable"
          label="Kind"
        />
        <v-row>
          <v-col cols="8">
            <v-text-field
              v-model="hotSideMeasurementForm.value"
              density="comfortable"
              label="Value"
              type="number"
            />
          </v-col>
          <v-col cols="4">
            <v-text-field
              v-model="hotSideMeasurementForm.unit"
              density="comfortable"
              label="Unit"
              :placeholder="getDefaultUnitForKind(hotSideMeasurementForm.kind)"
            />
          </v-col>
        </v-row>
        <v-text-field
          v-model="hotSideMeasurementForm.observed_at"
          density="comfortable"
          label="Observed at"
          type="datetime-local"
        />
        <v-textarea
          v-model="hotSideMeasurementForm.notes"
          auto-grow
          density="comfortable"
          label="Notes"
          rows="2"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="hotSideMeasurementDialog = false">Cancel</v-btn>
        <v-btn
          color="secondary"
          :disabled="!hotSideMeasurementForm.kind || !hotSideMeasurementForm.value"
          :loading="savingHotSideMeasurement"
          @click="recordHotSideMeasurement"
        >
          Add measurement
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useApiClient } from '@/composables/useApiClient'
import {
  useProductionApi,
  type Recipe,
  type BrewSession,
  type Vessel,
  type Volume as ProductionVolume,
  type Addition as ProductionAddition,
  type Measurement as ProductionMeasurement,
  type VolumeUnit,
  type AdditionType as ProductionAdditionType,
  type BatchSummary,
  type OccupancyStatus,
  OCCUPANCY_STATUS_VALUES,
} from '@/composables/useProductionApi'
import { useUnitPreferences } from '@/composables/useUnitPreferences'
import {
  convertTemperature,
  convertGravity,
  type TemperatureUnit,
  type GravityUnit,
} from '@/composables/useUnitConversion'

type Unit = 'ml' | 'usfloz' | 'ukfloz' | 'bbl'
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
  recipe_id: number | null
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

type ImportRowError = {
  row: number | null
  message: string
}

type BatchImportRowResult = {
  row: number
  status: 'created' | 'error'
  error?: string | null
  batch?: Batch
}

type BatchImportResponse = {
  totals: {
    total_rows: number
    created: number
    failed: number
  }
  results: BatchImportRowResult[]
}

type ImportSummary = {
  message: string
  type: 'success' | 'warning' | 'error'
}

const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'
const route = useRoute()
const { request } = useApiClient(apiBase)
const {
  getRecipes,
  getVessels,
  getVolumes: getProductionVolumes,
  createVolume: createProductionVolume,
  getBrewSessions,
  createBrewSession,
  updateBrewSession,
  getAdditionsByVolume,
  createAddition,
  getMeasurementsByVolume,
  createMeasurement,
  getBatchSummary,
  updateOccupancyStatus,
  normalizeText: apiNormalizeText,
  normalizeDateTime: apiNormalizeDateTime,
  toNumber: apiToNumber,
} = useProductionApi()

const {
  preferences,
  formatTemperaturePreferred,
  formatGravityPreferred,
  formatVolumePreferred,
} = useUnitPreferences()

const unitOptions: Unit[] = ['ml', 'usfloz', 'ukfloz', 'bbl']
const volumeUnitOptions: VolumeUnit[] = ['ml', 'usfloz', 'ukfloz', 'bbl']
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
const recipes = ref<Recipe[]>([])
const recipesLoading = ref(false)
const batchVolumes = ref<BatchVolume[]>([])
const processPhases = ref<BatchProcessPhase[]>([])
const additions = ref<Addition[]>([])
const measurements = ref<Measurement[]>([])
const volumeRelations = ref<VolumeRelation[]>([])
const batchSummary = ref<BatchSummary | null>(null)
const batchSummaryLoading = ref(false)

const selectedBatchId = ref<number | null>(null)
const activeTab = ref('summary')
const errorMessage = ref('')
const createBatchDialog = ref(false)
const createAdditionDialog = ref(false)
const createMeasurementDialog = ref(false)
const bulkImportDialog = ref(false)

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

const newBatch = reactive({
  short_name: '',
  brew_date: '',
  recipe_id: null as number | null,
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

const importFile = ref<File | File[] | null>(null)
const importUploading = ref(false)
const importResult = ref<BatchImportResponse | null>(null)
const importErrors = ref<ImportRowError[]>([])

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

const hotSideMeasurementKinds = [
  { title: 'Mash Temperature', value: 'mash_temp' },
  { title: 'Mash pH', value: 'mash_ph' },
  { title: 'Pre-Boil Gravity', value: 'pre_boil_gravity' },
  { title: 'Original Gravity', value: 'original_gravity' },
  { title: 'Boil Temperature', value: 'boil_temp' },
  { title: 'Post-Boil Volume', value: 'post_boil_volume' },
  { title: 'Other', value: 'other' },
]

// Occupancy status options - computed to ensure formatOccupancyStatus is defined
const occupancyStatusOptions = computed(() =>
  OCCUPANCY_STATUS_VALUES.map((status) => ({
    value: status,
    title: formatOccupancyStatus(status),
  }))
)

const rules = {
  required: (v: string) => !!v?.trim() || 'Required',
  positiveNumber: (v: string) => {
    const num = Number(v)
    return (Number.isFinite(num) && num > 0) || 'Must be positive'
  },
}


const selectedBatch = computed(() =>
  batches.value.find((batch) => batch.id === selectedBatchId.value) ?? null,
)

const routeBatchUuid = computed(() => {
  const params = route.params
  if ('uuid' in params) {
    const param = params.uuid
    if (typeof param === 'string' && param.trim()) {
      return param
    }
  }
  return null
})

const latestProcessPhase = computed(() => getLatest(processPhases.value, (item) => item.phase_at))
const latestLiquidPhase = computed(() => getLatest(batchVolumes.value, (item) => item.phase_at))

const sparklineWidth = 120
const sparklineHeight = 36

const temperatureSeries = computed(() =>
  buildMeasurementSeries(['temperature', 'temp'], sparklineWidth, sparklineHeight),
)
const gravitySeries = computed(() =>
  buildMeasurementSeries(['gravity', 'grav', 'sg'], sparklineWidth, sparklineHeight),
)
const phSeries = computed(() =>
  buildMeasurementSeries(['ph'], sparklineWidth, sparklineHeight),
)

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

const canImport = computed(() => Boolean(getSelectedImportFile()) && !importUploading.value)

const importSummary = computed<ImportSummary | null>(() => {
  if (!importResult.value) {
    return null
  }
  const successCount = getImportSuccessCount(importResult.value)
  const failureCount = getImportFailureCount(importResult.value, importErrors.value)
  const total = successCount + failureCount
  const message =
    total > 0
      ? `Imported ${successCount} ${successCount === 1 ? 'batch' : 'batches'}, ${failureCount} failed.`
      : 'Import completed.'
  if (failureCount > 0 && successCount === 0) {
    return { message, type: 'error' }
  }
  if (failureCount > 0) {
    return { message, type: 'warning' }
  }
  return { message, type: 'success' }
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

const recipeSelectItems = computed(() =>
  recipes.value.map((recipe) => ({
    title: recipe.name,
    value: recipe.id,
    style: recipe.style_name,
  })),
)

// Brew Session computed properties
const isEditingBrewSession = computed(() => editingBrewSessionId.value !== null)

const isBrewSessionFormValid = computed(() => {
  return brewSessionForm.brewed_at.trim().length > 0
})

const brewSessionsSorted = computed(() =>
  [...brewSessions.value].sort(
    (a, b) => new Date(b.brewed_at).getTime() - new Date(a.brewed_at).getTime()
  )
)

const selectedBrewSession = computed(() =>
  brewSessions.value.find((session) => session.id === selectedBrewSessionId.value) ?? null
)

const mashVesselOptions = computed(() =>
  vessels.value
    .filter((v) => v.status === 'active' && v.type.toLowerCase().includes('mash'))
    .map((v) => ({ id: v.id, name: v.name }))
)

const boilVesselOptions = computed(() =>
  vessels.value
    .filter((v) => v.status === 'active' && (v.type.toLowerCase().includes('kettle') || v.type.toLowerCase().includes('boil')))
    .map((v) => ({ id: v.id, name: v.name }))
)

const wortVolumeOptions = computed(() =>
  allVolumes.value.map((v) => ({
    id: v.id,
    label: v.name ? `${v.name} (${v.amount} ${v.amount_unit})` : `Volume #${v.id} (${v.amount} ${v.amount_unit})`,
  }))
)

const isVolumeFormValid = computed(() => {
  const amount = Number(volumeForm.amount)
  return Number.isFinite(amount) && amount > 0
})

const wortAdditionsSorted = computed(() =>
  [...wortAdditions.value].sort(
    (a, b) => new Date(b.added_at).getTime() - new Date(a.added_at).getTime()
  )
)

const wortMeasurementsSorted = computed(() =>
  [...wortMeasurements.value].sort(
    (a, b) => new Date(b.observed_at).getTime() - new Date(a.observed_at).getTime()
  )
)

const flowUnit = computed<Unit | null>(() => {
  const counts = new Map<Unit, number>()
  volumeRelations.value.forEach((relation) => {
    if (!relation.amount || relation.amount <= 0) {
      return
    }
    counts.set(relation.amount_unit, (counts.get(relation.amount_unit) ?? 0) + 1)
  })
  let selectedUnit: Unit | null = null
  let selectedCount = 0
  counts.forEach((count, unit) => {
    if (count > selectedCount) {
      selectedUnit = unit
      selectedCount = count
    }
  })
  return selectedUnit
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

watch([routeBatchUuid, batches], ([uuid]) => {
  if (uuid) {
    applyRouteSelection()
  }
})

watch(timelineObservedAtMenu, (isOpen) => {
  if (isOpen && !timelineReading.observed_at) {
    timelineReading.observed_at = nowInputValue()
  }
})

watch(bulkImportDialog, (isOpen) => {
  if (!isOpen) {
    resetImportState()
  }
})

watch(importFile, (value) => {
  if (value) {
    importResult.value = null
    importErrors.value = []
  }
})

// Watch for brew session selection to load wort additions/measurements
watch(selectedBrewSessionId, async (sessionId) => {
  if (!sessionId) {
    wortAdditions.value = []
    wortMeasurements.value = []
    return
  }
  const session = brewSessions.value.find((s) => s.id === sessionId)
  if (session?.wort_volume_id) {
    await loadWortData(session.wort_volume_id)
  } else {
    wortAdditions.value = []
    wortMeasurements.value = []
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
  batchSummary.value = null
}

function applyRouteSelection() {
  const uuid = routeBatchUuid.value
  if (!uuid) {
    return null
  }
  const match = batches.value.find((batch) => batch.uuid === uuid)
  if (match) {
    selectedBatchId.value = match.id
    return true
  }
  selectedBatchId.value = null
  return false
}

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

const get = <T>(path: string) => request<T>(path)
const post = <T>(path: string, payload: unknown) =>
  request<T>(path, { method: 'POST', body: JSON.stringify(payload) })
const postForm = <T>(path: string, payload: FormData) =>
  request<T>(path, { method: 'POST', body: payload, headers: new Headers() })

async function refreshAll() {
  errorMessage.value = ''
  try {
    await Promise.all([loadBatches(), loadVolumes(), loadRecipesData(), loadVesselsData(), loadAllVolumesData()])
    const routeApplied = applyRouteSelection()
    if (routeApplied === null) {
      const firstBatch = batches.value[0]
      if (!selectedBatchId.value && firstBatch) {
        selectedBatchId.value = firstBatch.id
      } else if (selectedBatchId.value) {
        await loadBatchData(selectedBatchId.value)
      }
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

async function loadRecipesData() {
  recipesLoading.value = true
  try {
    recipes.value = await getRecipes()
  } catch (error) {
    // Recipe loading failure is non-critical
    console.error('Failed to load recipes:', error)
  } finally {
    recipesLoading.value = false
  }
}

async function loadVesselsData() {
  try {
    vessels.value = await getVessels()
  } catch (error) {
    console.error('Failed to load vessels:', error)
  }
}

async function loadAllVolumesData() {
  try {
    allVolumes.value = await getProductionVolumes()
  } catch (error) {
    console.error('Failed to load volumes:', error)
  }
}

async function loadBatchData(batchId: number) {
  try {
    const [batchVolumesData, processPhasesData, additionsData, measurementsData, brewSessionsData] = await Promise.all([
      get<BatchVolume[]>(`/batch-volumes?batch_id=${batchId}`),
      get<BatchProcessPhase[]>(`/batch-process-phases?batch_id=${batchId}`),
      get<Addition[]>(`/additions?batch_id=${batchId}`),
      get<Measurement[]>(`/measurements?batch_id=${batchId}`),
      getBrewSessions(batchId),
    ])

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
    handleError(error)
  }
}

async function loadBatchSummary(batchId: number) {
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

async function loadWortData(volumeId: number) {
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
      recipe_id: newBatch.recipe_id,
      notes: normalizeText(newBatch.notes),
    }
    const created = await post<Batch>('/batches', payload)
    showNotice('Batch created')
    newBatch.short_name = ''
    newBatch.brew_date = ''
    newBatch.recipe_id = null
    newBatch.notes = ''
    await loadBatches()
    selectedBatchId.value = created.id
    createBatchDialog.value = false
  } catch (error) {
    handleError(error)
  }
}

async function uploadBatchImport() {
  const file = getSelectedImportFile()
  if (!file) {
    return
  }
  errorMessage.value = ''
  importUploading.value = true
  importResult.value = null
  importErrors.value = []
  try {
    const form = new FormData()
    form.append('file', file)
    const response = await postForm<BatchImportResponse>('/batches/import', form)
    importResult.value = response
    importErrors.value = parseImportErrors(response)
    const successCount = getImportSuccessCount(response)
    const failureCount = getImportFailureCount(response, importErrors.value)
    if (failureCount > 0) {
      const color = successCount > 0 ? 'warning' : 'error'
      showNotice(`Imported ${successCount} batches, ${failureCount} failed`, color)
    } else {
      showNotice(`Imported ${successCount} ${successCount === 1 ? 'batch' : 'batches'}`)
    }
    await loadBatches()
  } catch (error) {
    handleError(error)
  } finally {
    importUploading.value = false
  }
}

function downloadBatchTemplate() {
  const header = 'short_name,brew_date,notes\n'
  const blob = new Blob([header], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = 'batch-import-template.csv'
  link.click()
  URL.revokeObjectURL(url)
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
    createAdditionDialog.value = false
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
    createMeasurementDialog.value = false
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

function getSelectedImportFile() {
  const fileValue = importFile.value
  if (!fileValue) {
    return null
  }
  return Array.isArray(fileValue) ? fileValue[0] ?? null : fileValue
}

function resetImportState() {
  importFile.value = null
  importResult.value = null
  importErrors.value = []
  importUploading.value = false
}

function getImportSuccessCount(result: BatchImportResponse) {
  return result.totals?.created ?? 0
}

function getImportFailureCount(result: BatchImportResponse, errors: ImportRowError[]) {
  return result.totals?.failed ?? errors.length
}

function parseImportErrors(result: BatchImportResponse) {
  return (result.results ?? [])
    .filter((entry) => entry.status === 'error')
    .map((entry) => ({
      row: entry.row ?? null,
      message: entry.error ?? 'Import failed',
    }))
}

// ==================== Brew Session Functions ====================

function selectBrewSession(id: number) {
  selectedBrewSessionId.value = id
}

function clearBrewSessionSelection() {
  selectedBrewSessionId.value = null
}

function openCreateBrewSessionDialog() {
  editingBrewSessionId.value = null
  brewSessionForm.brewed_at = nowInputValue()
  brewSessionForm.mash_vessel_id = null
  brewSessionForm.boil_vessel_id = null
  brewSessionForm.wort_volume_id = null
  brewSessionForm.notes = ''
  brewSessionDialog.value = true
}

function openEditBrewSessionDialog(session: BrewSession) {
  editingBrewSessionId.value = session.id
  brewSessionForm.brewed_at = toLocalDateTimeInput(session.brewed_at)
  brewSessionForm.mash_vessel_id = session.mash_vessel_id
  brewSessionForm.boil_vessel_id = session.boil_vessel_id
  brewSessionForm.wort_volume_id = session.wort_volume_id
  brewSessionForm.notes = session.notes ?? ''
  brewSessionDialog.value = true
}

function closeBrewSessionDialog() {
  brewSessionDialog.value = false
  editingBrewSessionId.value = null
}

async function saveBrewSession() {
  if (!selectedBatchId.value || !isBrewSessionFormValid.value) {
    return
  }

  savingBrewSession.value = true
  errorMessage.value = ''

  try {
    const payload = {
      batch_id: selectedBatchId.value,
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

    closeBrewSessionDialog()
    if (selectedBatchId.value) {
      await loadBatchData(selectedBatchId.value)
    }
  } catch (error) {
    handleError(error)
  } finally {
    savingBrewSession.value = false
  }
}

function openCreateVolumeDialog() {
  volumeForm.name = selectedBatch.value ? `${selectedBatch.value.short_name} Wort` : ''
  volumeForm.description = ''
  volumeForm.amount = ''
  volumeForm.amount_unit = 'bbl'
  createVolumeDialog.value = true
}

async function createWortVolume() {
  if (!isVolumeFormValid.value) {
    return
  }

  savingVolume.value = true
  errorMessage.value = ''

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

function getVesselName(vesselId: number | null): string {
  if (!vesselId) return ''
  const vessel = vessels.value.find((v) => v.id === vesselId)
  return vessel?.name ?? `Vessel #${vesselId}`
}

function getVolumeName(volumeId: number | null): string {
  if (!volumeId) return ''
  const volume = allVolumes.value.find((v) => v.id === volumeId)
  return volume?.name ?? `Volume #${volumeId}`
}

function getVolumeAmount(volumeId: number | null): string {
  if (!volumeId) return ''
  const volume = allVolumes.value.find((v) => v.id === volumeId)
  if (!volume) return ''
  return `${volume.amount} ${volume.amount_unit}`
}

// ==================== Hot-Side Addition/Measurement Functions ====================

function openCreateHotSideAdditionDialog() {
  hotSideAdditionForm.addition_type = 'malt'
  hotSideAdditionForm.stage = ''
  hotSideAdditionForm.inventory_lot_uuid = ''
  hotSideAdditionForm.amount = ''
  hotSideAdditionForm.amount_unit = 'ml'
  hotSideAdditionForm.added_at = nowInputValue()
  hotSideAdditionForm.notes = ''
  hotSideAdditionDialog.value = true
}

async function recordHotSideAddition() {
  const session = selectedBrewSession.value
  if (!session?.wort_volume_id || !hotSideAdditionForm.amount) {
    return
  }

  savingHotSideAddition.value = true
  errorMessage.value = ''

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

function openCreateHotSideMeasurementDialog() {
  hotSideMeasurementForm.kind = 'mash_temp'
  hotSideMeasurementForm.value = ''
  hotSideMeasurementForm.unit = ''
  hotSideMeasurementForm.observed_at = nowInputValue()
  hotSideMeasurementForm.notes = ''
  hotSideMeasurementDialog.value = true
}

async function recordHotSideMeasurement() {
  const session = selectedBrewSession.value
  if (!session?.wort_volume_id || !hotSideMeasurementForm.kind || !hotSideMeasurementForm.value) {
    return
  }

  savingHotSideMeasurement.value = true
  errorMessage.value = ''

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

function getDefaultUnitForKind(kind: string): string {
  switch (kind) {
    case 'mash_temp':
    case 'boil_temp':
      return 'F'
    case 'mash_ph':
      return 'pH'
    case 'pre_boil_gravity':
    case 'original_gravity':
      return 'SG'
    case 'post_boil_volume':
      return 'bbl'
    default:
      return ''
  }
}

function toLocalDateTimeInput(isoString: string): string {
  if (!isoString) return ''
  const date = new Date(isoString)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

// ==================== End Brew Session Functions ====================

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
  if (unitMatch && unitMatch[1]) {
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

function formatDays(days: number | null | undefined) {
  if (days === null || days === undefined) {
    return '—'
  }
  if (days < 1) {
    const hours = Math.round(days * 24)
    return `${hours}h`
  }
  return `${days.toFixed(1)} days`
}

function formatDecimal(value: number | null | undefined, decimals = 1) {
  if (value === null || value === undefined) {
    return '—'
  }
  return value.toFixed(decimals)
}

function formatPercent(value: number | null | undefined, decimals = 1) {
  if (value === null || value === undefined) {
    return '—'
  }
  return `${value.toFixed(decimals)}%`
}

function formatPhase(phase: string) {
  const phaseLabels: Record<string, string> = {
    planning: 'Planning',
    mashing: 'Mashing',
    heating: 'Heating',
    boiling: 'Boiling',
    cooling: 'Cooling',
    fermenting: 'Fermenting',
    conditioning: 'Conditioning',
    packaging: 'Packaging',
    finished: 'Finished',
  }
  return phaseLabels[phase] ?? phase.charAt(0).toUpperCase() + phase.slice(1)
}

function getPhaseColor(phase: string) {
  const phaseColors: Record<string, string> = {
    planning: 'grey',
    mashing: 'orange',
    heating: 'deep-orange',
    boiling: 'red',
    cooling: 'cyan',
    fermenting: 'primary',
    conditioning: 'teal',
    packaging: 'blue',
    finished: 'success',
  }
  return phaseColors[phase] ?? 'secondary'
}

// ==================== Occupancy Status Functions ====================

function formatOccupancyStatus(status: string | null | undefined): string {
  if (!status) {
    return 'No status'
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

function getOccupancyStatusIcon(status: string | null | undefined): string {
  if (!status) {
    return 'mdi-help-circle-outline'
  }
  const statusIcons: Record<string, string> = {
    fermenting: 'mdi-molecule',
    conditioning: 'mdi-clock-outline',
    cold_crashing: 'mdi-snowflake',
    dry_hopping: 'mdi-leaf',
    carbonating: 'mdi-shimmer',
    holding: 'mdi-pause-circle-outline',
    packaging: 'mdi-package-variant',
  }
  return statusIcons[status] ?? 'mdi-circle'
}

async function changeOccupancyStatus(occupancyId: number, status: OccupancyStatus) {
  if (!selectedBatchId.value) {
    return
  }

  errorMessage.value = ''
  try {
    await updateOccupancyStatus(occupancyId, status)
    showNotice(`Status updated to ${formatOccupancyStatus(status)}`)
    // Reload batch summary to reflect the change
    await loadBatchSummary(selectedBatchId.value)
  } catch (error) {
    handleError(error)
  }
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

function buildMeasurementSeries(kinds: string[], width: number, height: number) {
  const normalizedKinds = kinds.map((kind) => normalizeMeasurementKind(kind))
  const filtered = measurements.value.filter((measurement) =>
    matchesMeasurementKind(measurement.kind, normalizedKinds),
  )
  const ordered = [...filtered].sort(
    (a, b) => measurementTimestamp(a) - measurementTimestamp(b),
  )

  // Determine if this is temperature or gravity series for unit conversion
  const isTemperature = normalizedKinds.some((k) => k === 'temperature' || k === 'temp')
  const isGravity = normalizedKinds.some((k) => k === 'gravity' || k === 'grav' || k === 'sg')

  // Convert values to preferred units for sparkline display
  const values = ordered
    .map((measurement) => {
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

  const latest = getLatest(filtered, (item) => item.observed_at ?? item.created_at)

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

  const { linePath, areaPath } = buildSparkline(values, width, height)
  return {
    values,
    latest,
    latestLabel,
    linePath,
    areaPath,
  }
}

/**
 * Normalize a temperature unit string to a TemperatureUnit type.
 * Defaults to 'f' (Fahrenheit) if not recognized.
 */
function normalizeTemperatureUnit(unit: string | null | undefined): TemperatureUnit {
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
function normalizeGravityUnit(unit: string | null | undefined): GravityUnit {
  if (!unit) return 'sg' // Default to SG
  const normalized = unit.trim().toLowerCase()
  if (normalized === 'plato' || normalized === '°p' || normalized === 'p') {
    return 'plato'
  }
  return 'sg' // Default to SG for 'sg', or unknown
}

function buildSparkline(values: number[], width: number, height: number) {
  if (values.length === 0) {
    return { linePath: '', areaPath: '' }
  }
  const min = Math.min(...values)
  const max = Math.max(...values)
  const range = max - min
  const step = values.length > 1 ? width / (values.length - 1) : width
  const points = values.map((value, index) => {
    const ratio = range === 0 ? 0.5 : (value - min) / range
    const x = index * step
    const y = height - ratio * height
    return { x, y }
  })
  const linePath = points
    .map((point, index) => `${index === 0 ? 'M' : 'L'} ${point.x} ${point.y}`)
    .join(' ')
  const lastPoint = points[points.length - 1]
  const firstPoint = points[0]
  if (!lastPoint || !firstPoint) {
    return { linePath, areaPath: '' }
  }
  const areaPath = `${linePath} L ${lastPoint.x} ${height} L ${firstPoint.x} ${height} Z`
  return { linePath, areaPath }
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
