<template>
  <v-container class="vessels-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="4">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-silo" />
            Vessel list
            <v-spacer />
            <v-btn
              icon="mdi-plus"
              size="small"
              variant="text"
              aria-label="Register vessel"
              @click="createVesselDialog = true"
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

            <v-list class="vessel-list" lines="two" active-color="primary">
              <v-list-item
                v-for="vessel in vessels"
                :key="vessel.id"
                :active="vessel.id === selectedVesselId"
                @click="selectVessel(vessel.id)"
              >
                <v-list-item-title>
                  {{ vessel.name }}
                </v-list-item-title>
                <v-list-item-subtitle>
                  {{ vessel.type }} - {{ formatCapacity(vessel.capacity, vessel.capacity_unit) }}
                </v-list-item-subtitle>
                <template #append>
                  <v-chip size="x-small" variant="tonal">
                    {{ vessel.status }}
                  </v-chip>
                </template>
              </v-list-item>

              <v-list-item v-if="vessels.length === 0">
                <v-list-item-title>No vessels yet</v-list-item-title>
                <v-list-item-subtitle>Use + to register the first vessel.</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="8">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-clipboard-text-outline" />
            Vessel details
            <v-spacer />
            <v-btn size="small" variant="text" :loading="loading" @click="refreshVessels">
              Refresh
            </v-btn>
            <v-btn size="small" variant="text" @click="clearSelection">Clear</v-btn>
          </v-card-title>
          <v-card-text>
            <v-alert
              v-if="!selectedVessel"
              density="comfortable"
              type="info"
              variant="tonal"
            >
              Select a vessel to view metadata and availability.
            </v-alert>

            <div v-else>
              <v-row>
                <v-col cols="12" md="6">
                  <v-card class="sub-card" variant="tonal">
                    <v-card-text>
                      <div class="text-overline">Vessel</div>
                      <div class="text-h5 font-weight-semibold">
                        {{ selectedVessel.name }}
                      </div>
                      <div class="text-body-2 text-medium-emphasis">
                        {{ selectedVessel.type }} - {{ selectedVessel.status }}
                      </div>
                      <div class="text-body-2 text-medium-emphasis">
                        Capacity {{ formatCapacity(selectedVessel.capacity, selectedVessel.capacity_unit) }}
                      </div>
                      <div class="text-body-2 text-medium-emphasis" v-if="selectedVessel.make || selectedVessel.model">
                        {{ selectedVessel.make ?? 'Make not set' }} {{ selectedVessel.model ?? '' }}
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>

                <v-col cols="12" md="6">
                  <v-card class="sub-card" variant="tonal">
                    <v-card-text>
                      <div class="text-overline">Metadata</div>
                      <div class="text-body-2 text-medium-emphasis">ID {{ selectedVessel.id }}</div>
                      <div class="text-body-2 text-medium-emphasis">UUID {{ selectedVessel.uuid }}</div>
                      <div class="text-body-2 text-medium-emphasis">
                        Updated {{ formatDateTime(selectedVessel.updated_at) }}
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>

              <!-- Occupancy Section -->
              <v-row class="mt-2">
                <v-col cols="12">
                  <v-card class="sub-card" variant="outlined">
                    <v-card-text>
                      <div class="text-overline mb-2">Current Occupancy</div>

                      <div v-if="selectedVesselOccupancy">
                        <v-row dense align="center">
                          <v-col cols="12" md="4">
                            <div class="text-caption text-medium-emphasis">Status</div>
                            <v-menu location="bottom">
                              <template #activator="{ props }">
                                <v-chip
                                  v-bind="props"
                                  :color="getOccupancyStatusColor(selectedVesselOccupancy.status)"
                                  variant="tonal"
                                  size="small"
                                  class="mt-1 cursor-pointer"
                                  append-icon="mdi-menu-down"
                                  :loading="updatingOccupancyStatus"
                                >
                                  {{ formatOccupancyStatus(selectedVesselOccupancy.status) }}
                                </v-chip>
                              </template>
                              <v-list density="compact" nav>
                                <v-list-subheader>Change status</v-list-subheader>
                                <v-list-item
                                  v-for="statusOption in occupancyStatusOptions"
                                  :key="statusOption.value"
                                  :active="statusOption.value === selectedVesselOccupancy.status"
                                  @click="changeOccupancyStatus(selectedVesselOccupancy.id, statusOption.value)"
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
                          </v-col>
                          <v-col cols="12" md="4">
                            <div class="text-caption text-medium-emphasis">Occupied Since</div>
                            <div class="text-body-2 font-weight-medium mt-1">
                              {{ formatDateTime(selectedVesselOccupancy.in_at) }}
                            </div>
                          </v-col>
                          <v-col cols="12" md="4">
                            <div class="text-caption text-medium-emphasis">Volume ID</div>
                            <div class="text-body-2 font-weight-medium mt-1">
                              {{ selectedVesselOccupancy.volume_id }}
                            </div>
                          </v-col>
                        </v-row>
                      </div>

                      <v-alert
                        v-else
                        density="compact"
                        type="success"
                        variant="tonal"
                      >
                        This vessel is currently available.
                      </v-alert>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>
            </div>

          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <v-dialog v-model="createVesselDialog" max-width="640">
    <v-card>
      <v-card-title class="text-h6">Register vessel</v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field v-model="newVessel.type" label="Type" placeholder="Fermenter" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="newVessel.name" label="Name" placeholder="FV-01" />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field v-model="newVessel.capacity" label="Capacity" type="number" />
          </v-col>
          <v-col cols="12" md="4">
            <v-select
              v-model="newVessel.capacity_unit"
              :items="unitOptions"
              label="Capacity unit"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-select
              v-model="newVessel.status"
              :items="vesselStatusOptions"
              label="Status"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="newVessel.make" label="Make" />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field v-model="newVessel.model" label="Model" />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn variant="text" @click="createVesselDialog = false">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!newVessel.type.trim() || !newVessel.name.trim() || !newVessel.capacity"
          @click="createVessel"
        >
          Add vessel
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useApiClient } from '@/composables/useApiClient'
import {
  useProductionApi,
  type OccupancyStatus,
  type Occupancy,
  OCCUPANCY_STATUS_VALUES,
} from '@/composables/useProductionApi'

type Unit = 'ml' | 'usfloz' | 'ukfloz'

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

const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

const unitOptions: Unit[] = ['ml', 'usfloz', 'ukfloz']
const vesselStatusOptions = ['active', 'inactive', 'retired']

const vessels = ref<Vessel[]>([])
const occupancies = ref<Occupancy[]>([])
const selectedVesselId = ref<number | null>(null)
const errorMessage = ref('')
const loading = ref(false)
const createVesselDialog = ref(false)
const updatingOccupancyStatus = ref(false)

const { request } = useApiClient(apiBase)
const { getActiveOccupancies, updateOccupancyStatus } = useProductionApi()

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
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

const selectedVessel = computed(() =>
  vessels.value.find((vessel) => vessel.id === selectedVesselId.value) ?? null,
)

const occupancyMap = computed(
  () => new Map(occupancies.value.map((occupancy) => [occupancy.vessel_id, occupancy])),
)

const selectedVesselOccupancy = computed(() => {
  if (!selectedVesselId.value) return null
  return occupancyMap.value.get(selectedVesselId.value) ?? null
})

const occupancyStatusOptions = computed(() =>
  OCCUPANCY_STATUS_VALUES.map((status) => ({
    value: status,
    title: formatOccupancyStatus(status),
  }))
)

onMounted(async () => {
  await refreshVessels()
})

watch(selectedVesselId, async () => {
  // Refresh occupancies when vessel selection changes
  await loadOccupancies()
})

function selectVessel(id: number) {
  selectedVesselId.value = id
}

function clearSelection() {
  selectedVesselId.value = null
}

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function refreshVessels() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [vesselData] = await Promise.all([
      request<Vessel[]>('/vessels'),
      loadOccupancies(),
    ])
    vessels.value = vesselData
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load vessels'
    errorMessage.value = message
    showNotice(message, 'error')
  } finally {
    loading.value = false
  }
}

async function loadOccupancies() {
  try {
    occupancies.value = await getActiveOccupancies()
  } catch (error) {
    console.error('Failed to load occupancies:', error)
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
    await request<Vessel>('/vessels', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    showNotice('Vessel registered')
    newVessel.type = ''
    newVessel.name = ''
    newVessel.capacity = ''
    newVessel.make = ''
    newVessel.model = ''
    await refreshVessels()
    createVesselDialog.value = false
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create vessel'
    errorMessage.value = message
    showNotice(message, 'error')
  }
}

function normalizeText(value: string) {
  const trimmed = value.trim()
  return trimmed.length > 0 ? trimmed : null
}

function toNumber(value: string | number | null) {
  if (value === null || value === undefined || value === '') {
    return null
  }
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : null
}

function formatCapacity(amount: number, unit: Unit) {
  return `${amount} ${unit}`
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
  updatingOccupancyStatus.value = true
  errorMessage.value = ''
  try {
    await updateOccupancyStatus(occupancyId, status)
    showNotice(`Status updated to ${formatOccupancyStatus(status)}`)
    await loadOccupancies()
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Failed to update status'
    errorMessage.value = message
    showNotice(message, 'error')
  } finally {
    updatingOccupancyStatus.value = false
  }
}
</script>

<style scoped>
.vessels-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.vessel-list {
  max-height: 420px;
  overflow: auto;
}

.cursor-pointer {
  cursor: pointer;
}
</style>
