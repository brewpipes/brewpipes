<template>
  <v-container class="pa-4" fluid>
    <v-alert
      v-if="loading"
      density="comfortable"
      type="info"
      variant="tonal"
    >
      Loading vessel...
    </v-alert>

    <v-alert
      v-else-if="error"
      density="comfortable"
      type="error"
      variant="tonal"
    >
      {{ error }}
    </v-alert>

    <template v-else-if="vessel">
      <!-- Header with back button -->
      <div class="d-flex align-center mb-4">
        <v-btn
          class="mr-3"
          icon="mdi-arrow-left"
          size="small"
          variant="text"
          @click="goBack"
        />
        <div>
          <div class="text-h5 font-weight-semibold">{{ vessel.name }}</div>
          <div class="text-body-2 text-medium-emphasis">
            {{ vessel.type }}
          </div>
        </div>
        <v-spacer />
        <v-btn
          class="mr-2"
          size="small"
          variant="text"
          @click="openEditDialog"
        >
          <v-icon class="mr-1" icon="mdi-pencil" size="small" />
          Edit
        </v-btn>
        <v-btn
          v-if="vessel.status !== 'retired'"
          color="warning"
          size="small"
          variant="text"
          @click="openRetireDialog"
        >
          <v-icon class="mr-1" icon="mdi-archive" size="small" />
          Retire
        </v-btn>
        <v-chip
          class="ml-2"
          :color="getVesselStatusColor(vessel.status)"
          size="small"
          variant="tonal"
        >
          {{ formatVesselStatus(vessel.status) }}
        </v-chip>
      </div>

      <v-row>
        <!-- Vessel Information Card -->
        <v-col cols="12" md="6">
          <v-card class="section-card">
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2" icon="mdi-silo" />
              Vessel Information
            </v-card-title>
            <v-card-text>
              <v-list density="compact" lines="two">
                <v-list-item>
                  <v-list-item-title>Name</v-list-item-title>
                  <v-list-item-subtitle>{{ vessel.name }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Type</v-list-item-title>
                  <v-list-item-subtitle>{{ vessel.type }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Status</v-list-item-title>
                  <v-list-item-subtitle>
                    <v-chip
                      :color="getVesselStatusColor(vessel.status)"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ formatVesselStatus(vessel.status) }}
                    </v-chip>
                  </v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Capacity</v-list-item-title>
                  <v-list-item-subtitle>
                    {{ formatVolumePreferred(vessel.capacity, vessel.capacity_unit) }}
                  </v-list-item-subtitle>
                </v-list-item>
                <v-list-item v-if="vessel.make || vessel.model">
                  <v-list-item-title>Make / Model</v-list-item-title>
                  <v-list-item-subtitle>
                    {{ vessel.make ?? '-' }} / {{ vessel.model ?? '-' }}
                  </v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Created</v-list-item-title>
                  <v-list-item-subtitle>{{ formatDateTime(vessel.created_at) }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Updated</v-list-item-title>
                  <v-list-item-subtitle>{{ formatDateTime(vessel.updated_at) }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- Current Occupancy Card -->
        <v-col cols="12" md="6">
          <v-card class="section-card">
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2" icon="mdi-clipboard-text-outline" />
              Current Occupancy
            </v-card-title>
            <v-card-text>
              <div v-if="currentOccupancy">
                <v-list density="compact" lines="two">
                  <v-list-item v-if="occupancyBatch">
                    <v-list-item-title>Batch</v-list-item-title>
                    <v-list-item-subtitle>
                      <router-link
                        class="batch-link"
                        :to="`/batches/${occupancyBatch.uuid}`"
                      >
                        {{ occupancyBatch.short_name }}
                      </router-link>
                    </v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>Status</v-list-item-title>
                    <v-list-item-subtitle>
                      <v-chip
                        :color="getOccupancyStatusColor(currentOccupancy.status)"
                        size="x-small"
                        variant="tonal"
                      >
                        {{ formatOccupancyStatus(currentOccupancy.status) }}
                      </v-chip>
                    </v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>In At</v-list-item-title>
                    <v-list-item-subtitle>{{ formatDateTime(currentOccupancy.in_at) }}</v-list-item-subtitle>
                  </v-list-item>
                </v-list>
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

        <!-- Metadata Card -->
        <v-col cols="12" md="6">
          <v-card class="section-card">
            <v-card-title class="d-flex align-center">
              <v-icon class="mr-2" icon="mdi-information-outline" />
              Metadata
            </v-card-title>
            <v-card-text>
              <v-list density="compact" lines="two">
                <v-list-item>
                  <v-list-item-title>ID</v-list-item-title>
                  <v-list-item-subtitle>{{ vessel.id }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>UUID</v-list-item-title>
                  <v-list-item-subtitle class="text-mono">{{ vessel.uuid }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>

    <!-- Edit Vessel Dialog -->
    <VesselEditDialog
      ref="editDialogRef"
      v-model="editDialogOpen"
      :vessel="vessel"
      @save="handleSaveVessel"
    />

    <!-- Snackbar for notifications -->
    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
      {{ snackbar.text }}
    </v-snackbar>
  </v-container>
</template>

<script lang="ts" setup>
  import type { UpdateVesselRequest } from '@/types'
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import VesselEditDialog from '@/components/vessel/VesselEditDialog.vue'
  import { useApiClient } from '@/composables/useApiClient'
  import {
    useFormatters,
    useOccupancyStatusFormatters,
    useVesselStatusFormatters,
  } from '@/composables/useFormatters'
  import {
    type Batch,
    type Occupancy,
    useProductionApi,
    type Vessel,
  } from '@/composables/useProductionApi'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  const route = useRoute()
  const router = useRouter()

  const productionApiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'
  const { request } = useApiClient(productionApiBase)
  const { getActiveOccupancies, getVesselByUUID, updateVessel } = useProductionApi()
  const { formatVolumePreferred } = useUnitPreferences()
  const { formatDateTime } = useFormatters()
  const { formatVesselStatus, getVesselStatusColor } = useVesselStatusFormatters()
  const { formatOccupancyStatus, getOccupancyStatusColor } = useOccupancyStatusFormatters()

  const loading = ref(true)
  const error = ref<string | null>(null)
  const vessel = ref<Vessel | null>(null)
  const occupancies = ref<Occupancy[]>([])
  const batches = ref<Batch[]>([])

  // Edit dialog state
  const editDialogOpen = ref(false)
  const editDialogRef = ref<InstanceType<typeof VesselEditDialog> | null>(null)

  // Snackbar state
  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  const routeUuid = computed(() => {
    const params = route.params
    if ('uuid' in params) {
      const param = params.uuid
      if (typeof param === 'string' && param.trim()) {
        return param
      }
    }
    return null
  })

  const currentOccupancy = computed(() => {
    if (!vessel.value) return null
    return occupancies.value.find(occ => occ.vessel_id === vessel.value!.id) ?? null
  })

  const occupancyBatch = computed(() => {
    if (!currentOccupancy.value || !currentOccupancy.value.batch_id) return null
    return batches.value.find(b => b.id === currentOccupancy.value!.batch_id) ?? null
  })

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function loadData () {
    const uuid = routeUuid.value
    if (!uuid) {
      error.value = 'Invalid vessel UUID'
      loading.value = false
      return
    }

    try {
      loading.value = true
      error.value = null

      // Fetch vessel by UUID directly, then fetch related data
      const vesselData = await getVesselByUUID(uuid)
      vessel.value = vesselData

      // Fetch occupancies and batches in parallel
      const [occupancyData, batchData] = await Promise.all([
        getActiveOccupancies(),
        request<Batch[]>('/batches'),
      ])

      occupancies.value = occupancyData
      batches.value = batchData
    } catch (error_) {
      console.error('Failed to load vessel:', error_)
      error.value = error_ instanceof Error && error_.message.includes('404') ? 'Vessel not found' : 'Failed to load vessel. Please try again.'
    } finally {
      loading.value = false
    }
  }

  function goBack () {
    router.push('/vessels/all')
  }

  function openEditDialog () {
    editDialogOpen.value = true
  }

  function openRetireDialog () {
    // Open edit dialog - user will change status to retired manually
    // The dialog will show the retirement warning when status is changed
    if (vessel.value) {
      editDialogOpen.value = true
    }
  }

  async function handleSaveVessel (data: UpdateVesselRequest) {
    if (!vessel.value) return

    editDialogRef.value?.setSaving(true)
    editDialogRef.value?.clearError()

    try {
      const updated = await updateVessel(vessel.value.id, data)
      vessel.value = updated
      editDialogOpen.value = false
      showNotice('Vessel updated successfully')
    } catch (error_) {
      console.error('Failed to update vessel:', error_)

      // Check for 409 Conflict (vessel has active occupancy)
      if (error_ instanceof Error && error_.message.includes('409')) {
        editDialogRef.value?.setError('Cannot change status: vessel has an active occupancy')
      } else {
        const message = error_ instanceof Error ? error_.message : 'Failed to update vessel'
        editDialogRef.value?.setError(message)
      }
    } finally {
      editDialogRef.value?.setSaving(false)
    }
  }

  onMounted(() => {
    loadData()
  })
</script>

<style scoped>
.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.batch-link {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
}

.batch-link:hover {
  text-decoration: underline;
}

.text-mono {
  font-family: monospace;
  font-size: 0.85em;
}
</style>
