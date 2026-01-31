<template>
  <v-container class="production-page" fluid>
    <v-row>
      <v-col cols="12">
        <v-card class="hero-card">
          <v-card-text>
            <v-row align="center">
              <v-col cols="12" md="7">
                <div class="kicker">Production Service</div>
                <div class="text-h3 font-weight-bold mb-2">
                  Brew day flow, end to end
                </div>
                <div class="text-body-1 text-medium-emphasis">
                  Create batches, anchor volumes in vessels, capture transfers, and log additions
                  and measurements with clear lineage.
                </div>

                <div class="d-flex flex-wrap align-center ga-2 mt-4">
                  <v-chip color="primary" size="small" variant="tonal">
                    API: {{ apiBase }}
                  </v-chip>
                  <v-chip color="secondary" size="small" variant="tonal">
                    Batches: {{ batches.length }}
                  </v-chip>
                  <v-chip v-if="latestProcessPhase" color="primary" size="small" variant="outlined">
                    Process: {{ latestProcessPhase.process_phase }}
                  </v-chip>
                  <v-chip v-if="latestLiquidPhase" color="secondary" size="small" variant="outlined">
                    Liquid: {{ latestLiquidPhase.liquid_phase }}
                  </v-chip>
                </div>
              </v-col>

              <v-col cols="12" md="5">
                <v-card class="hero-panel" variant="tonal">
                  <div class="text-overline">Active batch</div>
                  <div class="text-h5 font-weight-semibold">
                    {{ selectedBatch ? selectedBatch.short_name : 'Select a batch' }}
                  </div>
                  <div class="text-body-2 text-medium-emphasis mb-3">
                    {{ selectedBatch ? `Batch #${selectedBatch.id}` : 'Choose a batch to begin' }}
                  </div>
                  <div class="d-flex flex-wrap ga-2">
                    <v-btn color="primary" size="small" @click="refreshAll">
                      Refresh data
                    </v-btn>
                    <v-btn size="small" variant="tonal" @click="clearSelection">
                      Clear selection
                    </v-btn>
                  </div>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="mt-4" align="stretch">
      <v-col cols="12" md="4">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-barley" />
            Batches
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
                <v-list-item-subtitle>Create the first batch below.</v-list-item-subtitle>
              </v-list-item>
            </v-list>

            <v-divider class="my-4" />

            <div class="text-subtitle-1 font-weight-semibold mb-2">
              Create batch
            </div>
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
            <v-btn
              block
              color="primary"
              :disabled="!newBatch.short_name.trim()"
              @click="createBatch"
            >
              Create batch
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-beaker-outline" />
            Batch details
          </v-card-title>
          <v-card-text>
            <v-alert
              v-if="!selectedBatch"
              density="comfortable"
              type="info"
              variant="tonal"
            >
              Select a batch to manage volumes, phases, and brew day activity.
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
                      <div class="text-overline">Latest phases</div>
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
                      <div class="text-body-2 text-medium-emphasis">
                        Updated {{ formatDateTime(selectedBatch.updated_at) }}
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>

              <v-tabs v-model="activeTab" class="batch-tabs" color="primary" show-arrows>
                <v-tab value="start">Start</v-tab>
                <v-tab value="phases">Phases</v-tab>
                <v-tab value="additions">Additions</v-tab>
                <v-tab value="measurements">Measurements</v-tab>
                <v-tab value="transfers">Transfers</v-tab>
                <v-tab value="relations">Relations</v-tab>
                <v-tab value="timeline">Timeline</v-tab>
              </v-tabs>

              <v-window v-model="activeTab" class="mt-4">
                <v-window-item value="start">
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Register vessel
                        </v-card-title>
                        <v-card-text>
                          <v-text-field v-model="newVessel.type" label="Type" />
                          <v-text-field v-model="newVessel.name" label="Name" />
                          <v-text-field v-model="newVessel.capacity" label="Capacity" type="number" />
                          <v-select
                            v-model="newVessel.capacity_unit"
                            :items="unitOptions"
                            label="Capacity unit"
                          />
                          <v-select
                            v-model="newVessel.status"
                            :items="vesselStatusOptions"
                            label="Status"
                          />
                          <v-text-field v-model="newVessel.make" label="Make" />
                          <v-text-field v-model="newVessel.model" label="Model" />
                          <v-btn
                            block
                            color="secondary"
                            :disabled="!newVessel.type.trim() || !newVessel.name.trim() || !newVessel.capacity"
                            @click="createVessel"
                          >
                            Add vessel
                          </v-btn>
                        </v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Start batch volume
                        </v-card-title>
                        <v-card-text>
                          <v-text-field v-model="startVolume.name" label="Volume name" />
                          <v-text-field v-model="startVolume.description" label="Description" />
                          <v-text-field v-model="startVolume.amount" label="Amount" type="number" />
                          <v-select
                            v-model="startVolume.amount_unit"
                            :items="unitOptions"
                            label="Amount unit"
                          />
                          <v-select
                            v-model="startVolume.vessel_id"
                            :items="vesselItems"
                            label="Vessel"
                          />
                          <v-select
                            v-model="startVolume.liquid_phase"
                            :items="liquidPhaseOptions"
                            label="Liquid phase"
                          />
                          <v-text-field v-model="startVolume.phase_at" label="Phase time" type="datetime-local" />
                          <v-text-field v-model="startVolume.in_at" label="Vessel in time" type="datetime-local" />
                          <v-btn
                            block
                            color="primary"
                            :disabled="!startVolume.amount || !startVolume.vessel_id"
                            @click="createStartingVolume"
                          >
                            Create volume + occupancy
                          </v-btn>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>

                  <v-row class="mt-2">
                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1">
                          Active occupancy lookup
                        </v-card-title>
                        <v-card-text>
                          <v-select
                            v-model="occupancyLookup.kind"
                            :items="occupancyLookupOptions"
                            label="Lookup by"
                          />
                          <v-text-field v-model="occupancyLookup.id" label="ID" type="number" />
                          <v-btn
                            block
                            color="secondary"
                            :disabled="!occupancyLookup.id"
                            @click="lookupActiveOccupancy"
                          >
                            Fetch active occupancy
                          </v-btn>
                        </v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1">
                          Active occupancy
                        </v-card-title>
                        <v-card-text>
                          <div v-if="activeOccupancy">
                            <div class="text-body-2">
                              Occupancy #{{ activeOccupancy.id }}
                            </div>
                            <div class="text-body-2">
                              Vessel {{ activeOccupancy.vessel_id }}
                            </div>
                            <div class="text-body-2">
                              Volume {{ activeOccupancy.volume_id }}
                            </div>
                            <div class="text-body-2">
                              In: {{ formatDateTime(activeOccupancy.in_at) }}
                            </div>
                          </div>
                          <div v-else class="text-body-2 text-medium-emphasis">
                            No active occupancy loaded.
                          </div>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-window-item>

                <v-window-item value="phases">
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Liquid phase
                        </v-card-title>
                        <v-card-text>
                          <v-select
                            v-model="liquidPhaseForm.volume_id"
                            :items="volumeItems"
                            label="Volume"
                          />
                          <v-select
                            v-model="liquidPhaseForm.liquid_phase"
                            :items="liquidPhaseOptions"
                            label="Liquid phase"
                          />
                          <v-text-field v-model="liquidPhaseForm.phase_at" label="Phase time" type="datetime-local" />
                          <v-btn
                            block
                            color="primary"
                            :disabled="!liquidPhaseForm.volume_id"
                            @click="recordLiquidPhase"
                          >
                            Update liquid phase
                          </v-btn>

                          <v-divider class="my-4" />
                          <div class="text-subtitle-2 mb-2">History</div>
                          <v-list density="compact">
                            <v-list-item
                              v-for="phase in batchVolumesSorted"
                              :key="phase.id"
                            >
                              <v-list-item-title>
                                {{ phase.liquid_phase }}
                              </v-list-item-title>
                              <v-list-item-subtitle>
                                Volume {{ phase.volume_id }} - {{ formatDateTime(phase.phase_at) }}
                              </v-list-item-subtitle>
                            </v-list-item>
                            <v-list-item v-if="batchVolumesSorted.length === 0">
                              <v-list-item-title>No liquid phases yet.</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Process phase
                        </v-card-title>
                        <v-card-text>
                          <v-select
                            v-model="processPhaseForm.process_phase"
                            :items="processPhaseOptions"
                            label="Process phase"
                          />
                          <v-text-field v-model="processPhaseForm.phase_at" label="Phase time" type="datetime-local" />
                          <v-btn block color="secondary" @click="recordProcessPhase">
                            Add process phase
                          </v-btn>

                          <v-divider class="my-4" />
                          <div class="text-subtitle-2 mb-2">History</div>
                          <v-list density="compact">
                            <v-list-item
                              v-for="phase in processPhasesSorted"
                              :key="phase.id"
                            >
                              <v-list-item-title>
                                {{ phase.process_phase }}
                              </v-list-item-title>
                              <v-list-item-subtitle>
                                {{ formatDateTime(phase.phase_at) }}
                              </v-list-item-subtitle>
                            </v-list-item>
                            <v-list-item v-if="processPhasesSorted.length === 0">
                              <v-list-item-title>No process phases yet.</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-window-item>

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
                                <td>{{ formatValue(measurement.value, measurement.unit) }}</td>
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

                <v-window-item value="transfers">
                  <v-row>
                    <v-col cols="12" md="5">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Record transfer
                        </v-card-title>
                        <v-card-text>
                          <v-text-field
                            v-model="transferForm.source_occupancy_id"
                            label="Source occupancy ID"
                            type="number"
                          />
                          <v-select
                            v-model="transferForm.dest_vessel_id"
                            :items="vesselItems"
                            label="Destination vessel"
                          />
                          <v-select
                            v-model="transferForm.volume_id"
                            :items="volumeItems"
                            label="Volume"
                          />
                          <v-text-field v-model="transferForm.amount" label="Amount" type="number" />
                          <v-select
                            v-model="transferForm.amount_unit"
                            :items="unitOptions"
                            label="Amount unit"
                          />
                          <v-text-field v-model="transferForm.loss_amount" label="Loss amount" type="number" />
                          <v-select
                            v-model="transferForm.loss_unit"
                            :items="unitOptions"
                            label="Loss unit"
                          />
                          <v-text-field v-model="transferForm.started_at" label="Started at" type="datetime-local" />
                          <v-text-field v-model="transferForm.ended_at" label="Ended at" type="datetime-local" />
                          <v-btn
                            block
                            color="primary"
                            :disabled="!transferForm.source_occupancy_id || !transferForm.dest_vessel_id || !transferForm.volume_id || !transferForm.amount"
                            @click="recordTransfer"
                          >
                            Record transfer
                          </v-btn>
                        </v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" md="7">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1">
                          Transfer log
                        </v-card-title>
                        <v-card-text>
                          <v-table class="data-table" density="compact">
                            <thead>
                              <tr>
                                <th>Source</th>
                                <th>Destination</th>
                                <th>Amount</th>
                                <th>Start</th>
                              </tr>
                            </thead>
                            <tbody>
                              <tr v-for="transfer in transfersSorted" :key="transfer.id">
                                <td>{{ transfer.source_occupancy_id }}</td>
                                <td>{{ transfer.dest_occupancy_id }}</td>
                                <td>{{ formatAmount(transfer.amount, transfer.amount_unit) }}</td>
                                <td>{{ formatDateTime(transfer.started_at) }}</td>
                              </tr>
                              <tr v-if="transfersSorted.length === 0">
                                <td colspan="4">No transfers recorded.</td>
                              </tr>
                            </tbody>
                          </v-table>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-window-item>

                <v-window-item value="relations">
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Batch relation
                        </v-card-title>
                        <v-card-text>
                          <v-select
                            v-model="relationForm.parent_batch_id"
                            :items="batchItems"
                            label="Parent batch"
                          />
                          <v-select
                            v-model="relationForm.child_batch_id"
                            :items="batchItems"
                            label="Child batch"
                          />
                          <v-select
                            v-model="relationForm.relation_type"
                            :items="relationTypeOptions"
                            label="Relation"
                          />
                          <v-select
                            v-model="relationForm.volume_id"
                            :items="volumeItems"
                            label="Volume (optional)"
                          />
                          <v-btn
                            block
                            color="secondary"
                            :disabled="!relationForm.parent_batch_id || !relationForm.child_batch_id"
                            @click="recordBatchRelation"
                          >
                            Record batch relation
                          </v-btn>
                        </v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="tonal">
                        <v-card-title class="text-subtitle-1">
                          Volume relation
                        </v-card-title>
                        <v-card-text>
                          <v-select
                            v-model="volumeRelationForm.parent_volume_id"
                            :items="volumeItems"
                            label="Parent volume"
                          />
                          <v-select
                            v-model="volumeRelationForm.child_volume_id"
                            :items="volumeItems"
                            label="Child volume"
                          />
                          <v-select
                            v-model="volumeRelationForm.relation_type"
                            :items="relationTypeOptions"
                            label="Relation"
                          />
                          <v-text-field v-model="volumeRelationForm.amount" label="Amount" type="number" />
                          <v-select
                            v-model="volumeRelationForm.amount_unit"
                            :items="unitOptions"
                            label="Unit"
                          />
                          <v-btn
                            block
                            color="primary"
                            :disabled="!volumeRelationForm.parent_volume_id || !volumeRelationForm.child_volume_id || !volumeRelationForm.amount"
                            @click="recordVolumeRelation"
                          >
                            Record volume relation
                          </v-btn>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>

                  <v-row class="mt-2">
                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1">Batch relations</v-card-title>
                        <v-card-text>
                          <v-list density="compact">
                            <v-list-item
                              v-for="relation in batchRelationsSorted"
                              :key="relation.id"
                            >
                              <v-list-item-title>
                                {{ relation.relation_type }}
                              </v-list-item-title>
                              <v-list-item-subtitle>
                                {{ relation.parent_batch_id }} -> {{ relation.child_batch_id }}
                              </v-list-item-subtitle>
                            </v-list-item>
                            <v-list-item v-if="batchRelationsSorted.length === 0">
                              <v-list-item-title>No batch relations yet.</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-card-text>
                      </v-card>
                    </v-col>

                    <v-col cols="12" md="6">
                      <v-card class="sub-card" variant="outlined">
                        <v-card-title class="text-subtitle-1">Volume relations</v-card-title>
                        <v-card-text>
                          <v-list density="compact">
                            <v-list-item
                              v-for="relation in volumeRelationsSorted"
                              :key="relation.id"
                            >
                              <v-list-item-title>
                                {{ relation.relation_type }}
                              </v-list-item-title>
                              <v-list-item-subtitle>
                                {{ relation.parent_volume_id }} -> {{ relation.child_volume_id }}
                                - {{ formatAmount(relation.amount, relation.amount_unit) }}
                              </v-list-item-subtitle>
                            </v-list-item>
                            <v-list-item v-if="volumeRelationsSorted.length === 0">
                              <v-list-item-title>No volume relations yet.</v-list-item-title>
                            </v-list-item>
                          </v-list>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-window-item>

                <v-window-item value="timeline">
                  <v-card class="sub-card" variant="outlined">
                    <v-card-title class="text-subtitle-1">
                      Brew day timeline
                    </v-card-title>
                    <v-card-text>
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

type Vessel = {
  id: number
  uuid: string
  type: string
  name: string
  capacity: number
  capacity_unit: Unit
  make: string | null
  model: string | null
  status: 'active' | 'inactive' | 'retired'
  created_at: string
  updated_at: string
}

type Occupancy = {
  id: number
  uuid: string
  vessel_id: number
  volume_id: number
  in_at: string
  out_at: string | null
  created_at: string
  updated_at: string
}

type Transfer = {
  id: number
  uuid: string
  source_occupancy_id: number
  dest_occupancy_id: number
  amount: number
  amount_unit: Unit
  loss_amount: number | null
  loss_unit: Unit | null
  started_at: string
  ended_at: string | null
  created_at: string
  updated_at: string
}

type TransferRecord = {
  transfer: Transfer
  dest_occupancy: Occupancy
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

type BatchRelation = {
  id: number
  uuid: string
  parent_batch_id: number
  child_batch_id: number
  relation_type: RelationType
  volume_id: number | null
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

const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

const unitOptions: Unit[] = ['ml', 'usfloz', 'ukfloz']
const liquidPhaseOptions: LiquidPhase[] = ['water', 'wort', 'beer']
const processPhaseOptions: ProcessPhase[] = [
  'planning',
  'mashing',
  'heating',
  'boiling',
  'cooling',
  'fermenting',
  'conditioning',
  'packaging',
  'finished',
]
const relationTypeOptions: RelationType[] = ['split', 'blend']
const additionTypeOptions: AdditionType[] = [
  'malt',
  'hop',
  'yeast',
  'adjunct',
  'water_chem',
  'gas',
  'other',
]
const vesselStatusOptions = ['active', 'inactive', 'retired']
const additionTargetOptions = [
  { title: 'Batch', value: 'batch' },
  { title: 'Occupancy', value: 'occupancy' },
]
const occupancyLookupOptions = [
  { title: 'Vessel', value: 'vessel' },
  { title: 'Volume', value: 'volume' },
]

const batches = ref<Batch[]>([])
const vessels = ref<Vessel[]>([])
const volumes = ref<Volume[]>([])
const batchVolumes = ref<BatchVolume[]>([])
const processPhases = ref<BatchProcessPhase[]>([])
const transfers = ref<Transfer[]>([])
const additions = ref<Addition[]>([])
const measurements = ref<Measurement[]>([])
const batchRelations = ref<BatchRelation[]>([])
const volumeRelations = ref<VolumeRelation[]>([])
const activeOccupancy = ref<Occupancy | null>(null)

const selectedBatchId = ref<number | null>(null)
const activeTab = ref('start')
const errorMessage = ref('')

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

const newVessel = reactive({
  type: '',
  name: '',
  capacity: '',
  capacity_unit: 'ml' as Unit,
  status: 'active',
  make: '',
  model: '',
})

const startVolume = reactive({
  name: '',
  description: '',
  amount: '',
  amount_unit: 'ml' as Unit,
  vessel_id: null as number | null,
  liquid_phase: 'water' as LiquidPhase,
  phase_at: '',
  in_at: '',
})

const liquidPhaseForm = reactive({
  volume_id: null as number | null,
  liquid_phase: 'water' as LiquidPhase,
  phase_at: '',
})

const processPhaseForm = reactive({
  process_phase: 'planning' as ProcessPhase,
  phase_at: '',
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

const transferForm = reactive({
  source_occupancy_id: '',
  dest_vessel_id: null as number | null,
  volume_id: null as number | null,
  amount: '',
  amount_unit: 'ml' as Unit,
  loss_amount: '',
  loss_unit: 'ml' as Unit,
  started_at: '',
  ended_at: '',
})

const relationForm = reactive({
  parent_batch_id: null as number | null,
  child_batch_id: null as number | null,
  relation_type: 'split' as RelationType,
  volume_id: null as number | null,
})

const volumeRelationForm = reactive({
  parent_volume_id: null as number | null,
  child_volume_id: null as number | null,
  relation_type: 'split' as RelationType,
  amount: '',
  amount_unit: 'ml' as Unit,
})

const occupancyLookup = reactive({
  kind: 'vessel',
  id: '',
})

const selectedBatch = computed(() =>
  batches.value.find((batch) => batch.id === selectedBatchId.value) ?? null,
)

const latestProcessPhase = computed(() => getLatest(processPhases.value, (item) => item.phase_at))
const latestLiquidPhase = computed(() => getLatest(batchVolumes.value, (item) => item.phase_at))

const vesselItems = computed(() =>
  vessels.value.map((vessel) => ({
    title: `${vessel.name} (${vessel.type})`,
    value: vessel.id,
  })),
)

const volumeItems = computed(() =>
  volumes.value.map((volume) => ({
    title: `${volume.name ?? 'Volume'} #${volume.id}`,
    value: volume.id,
  })),
)

const batchItems = computed(() =>
  batches.value.map((batch) => ({
    title: `${batch.short_name} (#${batch.id})`,
    value: batch.id,
  })),
)

const batchVolumesSorted = computed(() =>
  sortByTime(batchVolumes.value, (item) => item.phase_at),
)
const processPhasesSorted = computed(() =>
  sortByTime(processPhases.value, (item) => item.phase_at),
)
const additionsSorted = computed(() =>
  sortByTime(additions.value, (item) => item.added_at),
)
const measurementsSorted = computed(() =>
  sortByTime(measurements.value, (item) => item.observed_at),
)
const transfersSorted = computed(() =>
  sortByTime(transfers.value, (item) => item.started_at),
)
const batchRelationsSorted = computed(() =>
  sortByTime(batchRelations.value, (item) => item.created_at),
)
const volumeRelationsSorted = computed(() =>
  sortByTime(volumeRelations.value, (item) => item.created_at),
)

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

  measurements.value.forEach((measurement) => {
    items.push({
      id: `measurement-${measurement.id}`,
      title: `Measurement: ${measurement.kind}`,
      subtitle: formatValue(measurement.value, measurement.unit),
      at: measurement.observed_at ?? measurement.created_at,
      color: 'secondary',
      icon: 'mdi-thermometer',
    })
  })

  transfers.value.forEach((transfer) => {
    items.push({
      id: `transfer-${transfer.id}`,
      title: 'Transfer',
      subtitle: `Occ ${transfer.source_occupancy_id} to ${transfer.dest_occupancy_id}`,
      at: transfer.started_at ?? transfer.created_at,
      color: 'info',
      icon: 'mdi-truck-fast-outline',
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

  return items.sort((a, b) => toTimestamp(a.at) - toTimestamp(b.at))
})

watch(selectedBatchId, (value) => {
  if (value) {
    loadBatchData(value)
  }
})

watch(batchVolumes, () => {
  if (!liquidPhaseForm.volume_id && batchVolumes.value.length > 0) {
    const latest = batchVolumesSorted.value[0]
    liquidPhaseForm.volume_id = latest?.volume_id ?? null
  }
})

watch(activeOccupancy, (value) => {
  if (value && !transferForm.source_occupancy_id) {
    transferForm.source_occupancy_id = String(value.id)
  }
  if (value && !transferForm.volume_id) {
    transferForm.volume_id = value.volume_id
  }
})

watch(selectedBatchId, (value) => {
  relationForm.parent_batch_id = value
})

onMounted(async () => {
  await refreshAll()
})

function selectBatch(id: number) {
  selectedBatchId.value = id
}

function clearSelection() {
  selectedBatchId.value = null
  activeOccupancy.value = null
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
    await Promise.all([loadBatches(), loadVessels(), loadVolumes()])
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

async function loadVessels() {
  vessels.value = await get<Vessel[]>('/vessels')
}

async function loadVolumes() {
  volumes.value = await get<Volume[]>('/volumes')
}

async function loadBatchData(batchId: number) {
  try {
    const [batchVolumesData, processPhasesData, transfersData, additionsData, measurementsData, batchRelationsData] =
      await Promise.all([
        get<BatchVolume[]>(`/batch-volumes?batch_id=${batchId}`),
        get<BatchProcessPhase[]>(`/batch-process-phases?batch_id=${batchId}`),
        get<Transfer[]>(`/transfers?batch_id=${batchId}`),
        get<Addition[]>(`/additions?batch_id=${batchId}`),
        get<Measurement[]>(`/measurements?batch_id=${batchId}`),
        get<BatchRelation[]>(`/batch-relations?batch_id=${batchId}`),
      ])

    batchVolumes.value = batchVolumesData
    processPhases.value = processPhasesData
    transfers.value = transfersData
    additions.value = additionsData
    measurements.value = measurementsData
    batchRelations.value = batchRelationsData

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
  } catch (error) {
    handleError(error)
  }
}

async function createVessel() {
  errorMessage.value = ''
  try {
    const payload = {
      type: newVessel.type.trim(),
      name: newVessel.name.trim(),
      capacity: toNumber(newVessel.capacity),
      capacity_unit: newVessel.capacity_unit,
      status: newVessel.status,
      make: normalizeText(newVessel.make),
      model: normalizeText(newVessel.model),
    }
    await post<Vessel>('/vessels', payload)
    showNotice('Vessel added')
    newVessel.type = ''
    newVessel.name = ''
    newVessel.capacity = ''
    newVessel.make = ''
    newVessel.model = ''
    await loadVessels()
  } catch (error) {
    handleError(error)
  }
}

async function createStartingVolume() {
  if (!selectedBatchId.value) {
    return
  }

  errorMessage.value = ''
  try {
    const volumePayload = {
      name: normalizeText(startVolume.name),
      description: normalizeText(startVolume.description),
      amount: toNumber(startVolume.amount),
      amount_unit: startVolume.amount_unit,
    }
    const volume = await post<Volume>('/volumes', volumePayload)

    const occupancyPayload = {
      vessel_id: startVolume.vessel_id,
      volume_id: volume.id,
      in_at: normalizeDateTime(startVolume.in_at),
    }
    const occupancy = await post<Occupancy>('/occupancies', occupancyPayload)
    activeOccupancy.value = occupancy

    const batchVolumePayload = {
      batch_id: selectedBatchId.value,
      volume_id: volume.id,
      liquid_phase: startVolume.liquid_phase,
      phase_at: normalizeDateTime(startVolume.phase_at),
    }
    await post<BatchVolume>('/batch-volumes', batchVolumePayload)

    showNotice('Starting volume created')
    startVolume.name = ''
    startVolume.description = ''
    startVolume.amount = ''
    startVolume.phase_at = ''
    startVolume.in_at = ''
    await Promise.all([loadVolumes(), loadBatchData(selectedBatchId.value)])
  } catch (error) {
    handleError(error)
  }
}

async function recordLiquidPhase() {
  if (!selectedBatchId.value || !liquidPhaseForm.volume_id) {
    return
  }
  errorMessage.value = ''
  try {
    const payload = {
      batch_id: selectedBatchId.value,
      volume_id: liquidPhaseForm.volume_id,
      liquid_phase: liquidPhaseForm.liquid_phase,
      phase_at: normalizeDateTime(liquidPhaseForm.phase_at),
    }
    await post<BatchVolume>('/batch-volumes', payload)
    showNotice('Liquid phase updated')
    liquidPhaseForm.phase_at = ''
    await loadBatchData(selectedBatchId.value)
  } catch (error) {
    handleError(error)
  }
}

async function recordProcessPhase() {
  if (!selectedBatchId.value) {
    return
  }
  errorMessage.value = ''
  try {
    const payload = {
      batch_id: selectedBatchId.value,
      process_phase: processPhaseForm.process_phase,
      phase_at: normalizeDateTime(processPhaseForm.phase_at),
    }
    await post<BatchProcessPhase>('/batch-process-phases', payload)
    showNotice('Process phase added')
    processPhaseForm.phase_at = ''
    await loadBatchData(selectedBatchId.value)
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

async function recordTransfer() {
  errorMessage.value = ''
  try {
    const payload = {
      source_occupancy_id: toNumber(transferForm.source_occupancy_id),
      dest_vessel_id: transferForm.dest_vessel_id,
      volume_id: transferForm.volume_id,
      amount: toNumber(transferForm.amount),
      amount_unit: transferForm.amount_unit,
      loss_amount: toNumber(transferForm.loss_amount),
      loss_unit: transferForm.loss_amount ? transferForm.loss_unit : null,
      started_at: normalizeDateTime(transferForm.started_at),
      ended_at: normalizeDateTime(transferForm.ended_at),
    }
    const record = await post<TransferRecord>('/transfers', payload)
    activeOccupancy.value = record.dest_occupancy
    showNotice('Transfer recorded')
    transferForm.amount = ''
    transferForm.loss_amount = ''
    transferForm.started_at = ''
    transferForm.ended_at = ''
    if (selectedBatchId.value) {
      await loadBatchData(selectedBatchId.value)
    }
  } catch (error) {
    handleError(error)
  }
}

async function recordBatchRelation() {
  if (!relationForm.parent_batch_id || !relationForm.child_batch_id) {
    return
  }
  errorMessage.value = ''
  try {
    const payload = {
      parent_batch_id: relationForm.parent_batch_id,
      child_batch_id: relationForm.child_batch_id,
      relation_type: relationForm.relation_type,
      volume_id: relationForm.volume_id,
    }
    await post<BatchRelation>('/batch-relations', payload)
    showNotice('Batch relation recorded')
    relationForm.child_batch_id = null
    relationForm.volume_id = null
    if (selectedBatchId.value) {
      await loadBatchData(selectedBatchId.value)
    }
  } catch (error) {
    handleError(error)
  }
}

async function recordVolumeRelation() {
  if (!volumeRelationForm.parent_volume_id || !volumeRelationForm.child_volume_id) {
    return
  }
  errorMessage.value = ''
  try {
    const payload = {
      parent_volume_id: volumeRelationForm.parent_volume_id,
      child_volume_id: volumeRelationForm.child_volume_id,
      relation_type: volumeRelationForm.relation_type,
      amount: toNumber(volumeRelationForm.amount),
      amount_unit: volumeRelationForm.amount_unit,
    }
    await post<VolumeRelation>('/volume-relations', payload)
    showNotice('Volume relation recorded')
    volumeRelationForm.amount = ''
    if (selectedBatchId.value) {
      await loadBatchData(selectedBatchId.value)
    }
  } catch (error) {
    handleError(error)
  }
}

async function lookupActiveOccupancy() {
  errorMessage.value = ''
  try {
    const id = toNumber(occupancyLookup.id)
    if (!id) {
      return
    }
    const query =
      occupancyLookup.kind === 'vessel'
        ? `active_vessel_id=${id}`
        : `active_volume_id=${id}`
    activeOccupancy.value = await get<Occupancy>(`/occupancies/active?${query}`)
    showNotice('Active occupancy loaded')
  } catch (error) {
    handleError(error)
  }
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
</script>

<style scoped>
.production-page {
  position: relative;
}

.hero-card {
  border: 1px solid rgba(var(--v-theme-primary), 0.25);
  background:
    linear-gradient(130deg, rgba(var(--v-theme-primary), 0.2), rgba(var(--v-theme-secondary), 0.18)),
    rgba(var(--v-theme-surface), 0.9);
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.25);
}

.hero-panel {
  border: 1px solid rgba(var(--v-theme-secondary), 0.35);
}

.kicker {
  text-transform: uppercase;
  letter-spacing: 0.28em;
  font-size: 0.7rem;
  color: rgba(var(--v-theme-on-surface), 0.6);
  margin-bottom: 6px;
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
