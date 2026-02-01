<template>
  <v-container class="vessels-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="4">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-silo" />
            Active Vessels
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

            <v-list active-color="primary" class="vessel-list" lines="two">
              <v-list-item
                v-for="vessel in sortedActiveVessels"
                :key="vessel.id"
                :active="vessel.id === selectedVesselId"
                @click="selectVessel(vessel.id)"
              >
                <v-list-item-title>
                  {{ vessel.name }}
                </v-list-item-title>
                <v-list-item-subtitle>
                  {{ vessel.type }} - {{ formatVolumePreferred(vessel.capacity, vessel.capacity_unit) }}
                </v-list-item-subtitle>
                <template #append>
                  <v-chip
                    :color="isVesselOccupied(vessel.id) ? 'primary' : 'grey'"
                    size="x-small"
                    variant="tonal"
                  >
                    {{ isVesselOccupied(vessel.id) ? 'Occupied' : 'Available' }}
                  </v-chip>
                </template>
              </v-list-item>

              <v-list-item v-if="sortedActiveVessels.length === 0 && !loading">
                <v-list-item-title>No active vessels</v-list-item-title>
                <v-list-item-subtitle>Register vessels in All Vessels to see them here.</v-list-item-subtitle>
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
            <v-btn :loading="loading" size="small" variant="text" @click="refreshVessels">
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
                        Capacity {{ formatVolumePreferred(selectedVessel.capacity, selectedVessel.capacity_unit) }}
                      </div>
                      <div v-if="selectedVessel.make || selectedVessel.model" class="text-body-2 text-medium-emphasis">
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
                        <v-row align="center" dense>
                          <v-col cols="12" md="4">
                            <div class="text-caption text-medium-emphasis">Status</div>
                            <v-menu location="bottom">
                              <template #activator="{ props }">
                                <v-chip
                                  v-bind="props"
                                  append-icon="mdi-menu-down"
                                  class="mt-1 cursor-pointer"
                                  :color="getOccupancyStatusColor(selectedVesselOccupancy.status)"
                                  :loading="updatingOccupancyStatus"
                                  size="small"
                                  variant="tonal"
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
                                      class="mr-2"
                                      :color="getOccupancyStatusColor(statusOption.value)"
                                      size="24"
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
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useApiClient } from '@/composables/useApiClient'
  import {
    useFormatters,
    useOccupancyStatusFormatters,
  } from '@/composables/useFormatters'
  import {
    type Occupancy,
    OCCUPANCY_STATUS_VALUES,
    type OccupancyStatus,
    useProductionApi,
    type Vessel,
  } from '@/composables/useProductionApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

  const vessels = ref<Vessel[]>([])
  const occupancies = ref<Occupancy[]>([])
  const selectedVesselId = ref<number | null>(null)
  const errorMessage = ref('')
  const loading = ref(false)
  const updatingOccupancyStatus = ref(false)

  const { request } = useApiClient(apiBase)
  const { getActiveOccupancies, updateOccupancyStatus } = useProductionApi()
  const { formatVolumePreferred } = useUnitPreferences()
  const { formatDateTime } = useFormatters()
  const {
    formatOccupancyStatus,
    getOccupancyStatusColor,
    getOccupancyStatusIcon,
  } = useOccupancyStatusFormatters()

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  // Filter to only active vessels
  const activeVessels = computed(() =>
    vessels.value.filter(vessel => vessel.status === 'active'),
  )

  // Sort: occupied vessels first, then alphabetically by name
  const sortedActiveVessels = computed(() => {
    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...activeVessels.value].sort((a, b) => {
      const aOccupied = isVesselOccupied(a.id)
      const bOccupied = isVesselOccupied(b.id)

      // Occupied vessels first
      if (aOccupied && !bOccupied) return -1
      if (!aOccupied && bOccupied) return 1

      // Within same occupancy group, sort alphabetically by name
      return a.name.localeCompare(b.name)
    })
  })

  const selectedVessel = computed(() =>
    vessels.value.find(vessel => vessel.id === selectedVesselId.value) ?? null,
  )

  const occupancyMap = computed(
    () => new Map(occupancies.value.map(occupancy => [occupancy.vessel_id, occupancy])),
  )

  const selectedVesselOccupancy = computed(() => {
    if (!selectedVesselId.value) return null
    return occupancyMap.value.get(selectedVesselId.value) ?? null
  })

  const occupancyStatusOptions = computed(() =>
    OCCUPANCY_STATUS_VALUES.map(status => ({
      value: status,
      title: formatOccupancyStatus(status),
    })),
  )

  function isVesselOccupied (vesselId: number): boolean {
    return occupancyMap.value.has(vesselId)
  }

  onMounted(async () => {
    await refreshVessels()
  })

  watch(selectedVesselId, async () => {
    // Refresh occupancies when vessel selection changes
    await loadOccupancies()
  })

  function selectVessel (id: number) {
    selectedVesselId.value = id
  }

  function clearSelection () {
    selectedVesselId.value = null
  }

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function refreshVessels () {
    loading.value = true
    errorMessage.value = ''
    try {
      const [vesselData] = await Promise.all([
        request<Vessel[]>('/vessels'),
        loadOccupancies(),
      ])
      vessels.value = vesselData

      // Auto-select first vessel if none selected
      const firstVessel = sortedActiveVessels.value[0]
      if (!selectedVesselId.value && firstVessel) {
        selectedVesselId.value = firstVessel.id
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load vessels'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      loading.value = false
    }
  }

  async function loadOccupancies () {
    try {
      occupancies.value = await getActiveOccupancies()
    } catch (error) {
      console.error('Failed to load occupancies:', error)
    }
  }

  async function changeOccupancyStatus (occupancyId: number, status: OccupancyStatus) {
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
