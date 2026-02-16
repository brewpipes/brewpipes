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
          :batch-name="selectedVesselBatchName"
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
          :batch-name="selectedVesselBatchName"
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
  import type { Batch, Occupancy, OccupancyStatus, UpdateVesselRequest, Vessel } from '@/types'
  import { computed, onMounted, ref, watch } from 'vue'
  import VesselEditDialog from '@/components/vessel/VesselEditDialog.vue'
  import VesselDetails from '@/components/VesselDetails.vue'
  import VesselList from '@/components/VesselList.vue'
  import { useAsyncAction } from '@/composables/useAsyncAction'
  import { useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useSnackbar } from '@/composables/useSnackbar'
  import { useVesselActions } from '@/composables/useVesselActions'

  const vessels = ref<Vessel[]>([])
  const occupancies = ref<Occupancy[]>([])
  const batches = ref<Batch[]>([])
  const selectedVesselUuid = ref<string | null>(null)

  const { execute, loading, error: loadError } = useAsyncAction()
  const { execute: executeLoadOccupancies } = useAsyncAction({ onError: () => {} })

  // Edit dialog state
  const editDialogOpen = ref(false)
  const editDialogRef = ref<InstanceType<typeof VesselEditDialog> | null>(null)

  const { getVessels, getActiveOccupancies, getBatches, updateOccupancyStatus } = useProductionApi()
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

  const batchMap = computed(
    () => new Map(batches.value.map(b => [b.uuid, b])),
  )

  const selectedVesselOccupancy = computed(() => {
    if (!selectedVesselUuid.value) return null
    return occupancyMap.value.get(selectedVesselUuid.value) ?? null
  })

  const selectedVesselBatchName = computed(() => {
    const occ = selectedVesselOccupancy.value
    if (!occ?.batch_uuid) return null
    return batchMap.value.get(occ.batch_uuid)?.short_name ?? null
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
    await execute(async () => {
      const [vesselData, , batchData] = await Promise.all([
        getVessels(),
        loadOccupancies(),
        getBatches(),
      ])
      vessels.value = vesselData
      batches.value = batchData

      // Auto-select first active vessel if none selected
      const firstVessel = activeVessels.value[0]
      if (!selectedVesselUuid.value && firstVessel) {
        selectedVesselUuid.value = firstVessel.uuid
      }
    })
    if (loadError.value) {
      showNotice(loadError.value, 'error')
    }
  }

  async function loadOccupancies () {
    // Occupancy loading failure is non-critical
    await executeLoadOccupancies(async () => {
      occupancies.value = await getActiveOccupancies()
    })
  }

  const { execute: executeStatusChange } = useAsyncAction({
    onError: (message) => showNotice(message, 'error'),
  })

  async function changeOccupancyStatus (occupancyUuid: string, status: OccupancyStatus) {
    await executeStatusChange(async () => {
      await updateOccupancyStatus(occupancyUuid, status)
      showNotice(`Status updated to ${formatOccupancyStatus(status)}`)
      await loadOccupancies()
    })
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
