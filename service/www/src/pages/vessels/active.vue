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


</template>

<script lang="ts" setup>
  import type { UpdateVesselRequest } from '@/types'
  import { computed, onMounted, ref, watch } from 'vue'
  import VesselEditDialog from '@/components/vessel/VesselEditDialog.vue'
  import VesselDetails from '@/components/VesselDetails.vue'
  import VesselList from '@/components/VesselList.vue'
  import { useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import { useSnackbar } from '@/composables/useSnackbar'
  import type { Occupancy, OccupancyStatus, Vessel } from '@/types'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useVesselActions } from '@/composables/useVesselActions'

  const vessels = ref<Vessel[]>([])
  const occupancies = ref<Occupancy[]>([])
  const selectedVesselUuid = ref<string | null>(null)
  const loading = ref(false)

  // Edit dialog state
  const editDialogOpen = ref(false)
  const editDialogRef = ref<InstanceType<typeof VesselEditDialog> | null>(null)

  const { getVessels, getActiveOccupancies, updateOccupancyStatus } = useProductionApi()
  const { formatOccupancyStatus } = useOccupancyStatusFormatters()
  const { saveVessel } = useVesselActions()

  const { showNotice } = useSnackbar()

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

  async function refreshVessels () {
    loading.value = true
    try {
      const [vesselData] = await Promise.all([
        getVessels(),
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

    const updated = await saveVessel(selectedVessel.value.uuid, data, editDialogRef)
    if (updated) {
      const index = vessels.value.findIndex(v => v.uuid === updated.uuid)
      if (index !== -1) {
        vessels.value[index] = updated
      }
      editDialogOpen.value = false
    }
  }
</script>

<style scoped>
.vessels-page {
  position: relative;
}
</style>
