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
      <div class="d-flex align-center flex-wrap ga-2 mb-4">
        <v-btn
          class="mr-1"
          icon="mdi-arrow-left"
          size="small"
          variant="text"
          @click="goBack"
        />
        <div class="mr-auto">
          <div class="text-h5 font-weight-semibold">{{ vessel.name }}</div>
          <div class="text-body-2 text-medium-emphasis">
            {{ formatVesselType(vessel.type) }}
          </div>
        </div>
        <div class="d-flex align-center ga-1">
          <v-btn
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
            :color="getVesselStatusColor(vessel.status)"
            size="small"
            variant="tonal"
          >
            {{ formatVesselStatus(vessel.status) }}
          </v-chip>
        </div>
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
                  <v-list-item-subtitle>{{ formatVesselType(vessel.type) }}</v-list-item-subtitle>
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

    <!-- Retire Vessel Dialog -->
    <VesselRetireDialog
      ref="retireDialogRef"
      v-model="retireDialogOpen"
      :vessel="vessel"
      @confirm="handleRetireVessel"
    />

  </v-container>
</template>

<script lang="ts" setup>
  import type { Batch, Occupancy, UpdateVesselRequest, Vessel } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import VesselEditDialog from '@/components/vessel/VesselEditDialog.vue'
  import VesselRetireDialog from '@/components/vessel/VesselRetireDialog.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import {
    useFormatters,
    useOccupancyStatusFormatters,
    useVesselStatusFormatters,
    useVesselTypeFormatters,
  } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useRouteUuid } from '@/composables/useRouteUuid'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { useVesselActions } from '@/composables/useVesselActions'

  const router = useRouter()

  const { getActiveOccupancies, getBatches, getVessel } = useProductionApi()
  const { formatVolumePreferred } = useUnitPreferences()
  const { formatDateTime } = useFormatters()
  const { formatVesselStatus, getVesselStatusColor } = useVesselStatusFormatters()
  const { formatVesselType } = useVesselTypeFormatters()
  const { formatOccupancyStatus, getOccupancyStatusColor } = useOccupancyStatusFormatters()
  const { uuid: routeUuid } = useRouteUuid()

  const error = ref('')
  const { execute: executeLoad, loading } = useAsyncAction({
    onError: (message) => {
      error.value = message.includes('404') ? 'Vessel not found' : 'Failed to load vessel. Please try again.'
    },
  })
  loading.value = true

  const vessel = ref<Vessel | null>(null)
  const occupancies = ref<Occupancy[]>([])
  const batches = ref<Batch[]>([])

  // Edit dialog state
  const editDialogOpen = ref(false)
  const editDialogRef = ref<InstanceType<typeof VesselEditDialog> | null>(null)

  // Retire dialog state
  const retireDialogOpen = ref(false)
  const retireDialogRef = ref<InstanceType<typeof VesselRetireDialog> | null>(null)

  const { showNotice } = useSnackbar()
  const { retireVessel, saveVessel } = useVesselActions()

  const currentOccupancy = computed(() => {
    if (!vessel.value) return null
    return occupancies.value.find(occ => occ.vessel_uuid === vessel.value!.uuid) ?? null
  })

  const occupancyBatch = computed(() => {
    if (!currentOccupancy.value || !currentOccupancy.value.batch_uuid) return null
    return batches.value.find(b => b.uuid === currentOccupancy.value!.batch_uuid) ?? null
  })

  async function loadData () {
    const uuid = routeUuid.value
    if (!uuid) {
      error.value = 'Invalid vessel UUID'
      loading.value = false
      return
    }

    await executeLoad(async () => {
      // Fetch vessel by UUID directly, then fetch related data
      const vesselData = await getVessel(uuid)
      vessel.value = vesselData

      // Fetch occupancies and batches in parallel
      const [occupancyData, batchData] = await Promise.all([
        getActiveOccupancies(),
        getBatches(),
      ])

      occupancies.value = occupancyData
      batches.value = batchData
    })
  }

  function goBack () {
    router.push('/vessels/all')
  }

  function openEditDialog () {
    editDialogOpen.value = true
  }

  function openRetireDialog () {
    if (vessel.value) {
      retireDialogOpen.value = true
    }
  }

  async function handleSaveVessel (data: UpdateVesselRequest) {
    if (!vessel.value) return

    const updated = await saveVessel(vessel.value.uuid, data, editDialogRef)
    if (updated) {
      vessel.value = updated
      editDialogOpen.value = false
    }
  }

  async function handleRetireVessel () {
    if (!vessel.value) return

    const updated = await retireVessel(vessel.value.uuid, retireDialogRef)
    if (updated) {
      vessel.value = updated
      retireDialogOpen.value = false
    }
  }

  onMounted(() => {
    loadData()
  })
</script>

<style scoped>
.text-mono {
  font-family: monospace;
  font-size: 0.85em;
}
</style>
