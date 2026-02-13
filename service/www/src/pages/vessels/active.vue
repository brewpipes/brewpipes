<template>
  <v-container class="vessels-page" fluid>
    <!-- Mobile: Show list or detail based on selection -->
    <v-row v-if="$vuetify.display.smAndDown" align="stretch">
      <v-col v-if="!selectedVesselUuid" cols="12">
        <VesselList
          :loading="loading"
          :occupancies="occupancies"
          :selected-vessel-uuid="selectedVesselUuid"
          :vessels="activeVessels"
          @refresh="refreshVessels"
          @select="selectVessel"
        />
      </v-col>

      <v-col v-else cols="12">
        <VesselDetails
          :loading="loading"
          :occupancy="selectedVesselOccupancy"
          :vessel="selectedVessel"
          @clear="clearSelection"
          @edit="openEditDialog"
          @occupancy-status-change="changeOccupancyStatus"
          @refresh="refreshVessels"
        />
      </v-col>
    </v-row>

    <!-- Desktop: Side-by-side layout -->
    <v-row v-else align="stretch">
      <v-col cols="12" md="4">
        <VesselList
          :loading="loading"
          :occupancies="occupancies"
          :selected-vessel-uuid="selectedVesselUuid"
          :vessels="activeVessels"
          @refresh="refreshVessels"
          @select="selectVessel"
        />
      </v-col>

      <v-col cols="12" md="8">
        <VesselDetails
          :loading="loading"
          :occupancy="selectedVesselOccupancy"
          :vessel="selectedVessel"
          @clear="clearSelection"
          @edit="openEditDialog"
          @occupancy-status-change="changeOccupancyStatus"
          @refresh="refreshVessels"
        />
      </v-col>
    </v-row>

    <!-- Edit Vessel Dialog -->
    <VesselEditDialog
      ref="editDialogRef"
      v-model="editDialogOpen"
      :vessel="selectedVessel"
      @save="handleSaveVessel"
    />
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>
</template>

<script lang="ts" setup>
  import type { UpdateVesselRequest } from '@/types'
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import VesselEditDialog from '@/components/vessel/VesselEditDialog.vue'
  import VesselDetails from '@/components/VesselDetails.vue'
  import VesselList from '@/components/VesselList.vue'
  import { useApiClient } from '@/composables/useApiClient'
  import { useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import {
    type Occupancy,
    type OccupancyStatus,
    useProductionApi,
    type Vessel,
  } from '@/composables/useProductionApi'

  const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

  const vessels = ref<Vessel[]>([])
  const occupancies = ref<Occupancy[]>([])
  const selectedVesselUuid = ref<string | null>(null)
  const loading = ref(false)

  // Edit dialog state
  const editDialogOpen = ref(false)
  const editDialogRef = ref<InstanceType<typeof VesselEditDialog> | null>(null)

  const { request } = useApiClient(apiBase)
  const { getActiveOccupancies, updateOccupancyStatus, updateVessel } = useProductionApi()
  const { formatOccupancyStatus } = useOccupancyStatusFormatters()

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  // Filter to only active vessels
  const activeVessels = computed(() =>
    vessels.value.filter(vessel => vessel.status === 'active'),
  )

  const selectedVessel = computed(() =>
    vessels.value.find(vessel => vessel.uuid === selectedVesselUuid.value) ?? null,
  )

  const occupancyMap = computed(
    () => new Map(occupancies.value.map(occupancy => [occupancy.vessel_uuid, occupancy])),
  )

  const selectedVesselOccupancy = computed(() => {
    if (!selectedVesselUuid.value) return null
    return occupancyMap.value.get(selectedVesselUuid.value) ?? null
  })

  onMounted(async () => {
    await refreshVessels()
  })

  watch(selectedVesselUuid, async () => {
    // Refresh occupancies when vessel selection changes
    await loadOccupancies()
  })

  function selectVessel (uuid: string) {
    selectedVesselUuid.value = uuid
  }

  function clearSelection () {
    selectedVesselUuid.value = null
  }

  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function refreshVessels () {
    loading.value = true
    try {
      const [vesselData] = await Promise.all([
        request<Vessel[]>('/vessels'),
        loadOccupancies(),
      ])
      vessels.value = vesselData

      // Auto-select first active vessel if none selected
      const firstVessel = activeVessels.value[0]
      if (!selectedVesselUuid.value && firstVessel) {
        selectedVesselUuid.value = firstVessel.uuid
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load vessels'
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

  async function changeOccupancyStatus (occupancyUuid: string, status: OccupancyStatus) {
    try {
      await updateOccupancyStatus(occupancyUuid, status)
      showNotice(`Status updated to ${formatOccupancyStatus(status)}`)
      await loadOccupancies()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to update status'
      showNotice(message, 'error')
    }
  }

  function openEditDialog () {
    if (selectedVessel.value) {
      editDialogOpen.value = true
    }
  }

  async function handleSaveVessel (data: UpdateVesselRequest) {
    if (!selectedVessel.value) return

    editDialogRef.value?.setSaving(true)
    editDialogRef.value?.clearError()

    try {
      const updated = await updateVessel(selectedVessel.value.uuid, data)
      // Update the vessel in the list
      const index = vessels.value.findIndex(v => v.uuid === updated.uuid)
      if (index !== -1) {
        vessels.value[index] = updated
      }
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
</script>

<style scoped>
.vessels-page {
  position: relative;
}
</style>
