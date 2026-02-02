<template>
  <v-container class="vessels-page" fluid>
    <v-card class="section-card">
      <v-card-title class="d-flex align-center">
        <v-icon class="mr-2" icon="mdi-silo" />
        All Vessels
        <v-spacer />
        <v-text-field
          v-model="search"
          append-inner-icon="mdi-magnify"
          class="search-field"
          clearable
          density="compact"
          hide-details
          label="Search"
          single-line
          variant="outlined"
        />
        <v-btn
          class="ml-2"
          :loading="loading"
          size="small"
          variant="text"
          @click="refreshData"
        >
          Refresh
        </v-btn>
        <v-btn
          class="ml-2"
          color="primary"
          size="small"
          variant="text"
          @click="openCreateDialog"
        >
          New vessel
        </v-btn>
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

        <v-data-table
          class="data-table vessels-table"
          density="compact"
          :headers="headers"
          item-value="id"
          :items="sortedVessels"
          :loading="loading"
          :search="search"
          @dblclick:row="onRowDoubleClick"
        >
          <template #item.id="{ item }">
            <span class="text-medium-emphasis">#{{ item.id }}</span>
          </template>

          <template #item.name="{ item }">
            <span class="font-weight-medium">{{ item.name }}</span>
          </template>

          <template #item.capacity="{ item }">
            {{ formatVolumePreferred(item.capacity, item.capacity_unit) }}
          </template>

          <template #item.status="{ item }">
            <v-chip
              :color="getVesselStatusColor(item.status)"
              size="small"
              variant="tonal"
            >
              {{ formatVesselStatus(item.status) }}
            </v-chip>
          </template>

          <template #item.occupancy="{ item }">
            <router-link
              v-if="getOccupancyBatchInfo(item.id)"
              class="batch-link"
              :to="`/batches/${getOccupancyBatchInfo(item.id)!.uuid}`"
            >
              {{ getOccupancyBatchInfo(item.id)!.short_name }}
            </router-link>
            <span v-else class="text-medium-emphasis">Unoccupied</span>
          </template>

          <template #item.updated_at="{ item }">
            {{ formatRelativeTime(item.updated_at) }}
          </template>

          <template #no-data>
            <div class="text-center py-4">
              <div class="text-body-2 text-medium-emphasis">No vessels yet.</div>
              <v-btn
                class="mt-2"
                color="primary"
                size="small"
                variant="text"
                @click="openCreateDialog"
              >
                Register your first vessel
              </v-btn>
            </div>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </v-container>

  <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="3000">
    {{ snackbar.text }}
  </v-snackbar>

  <!-- Create Vessel Dialog -->
  <v-dialog v-model="createVesselDialog" max-width="640" persistent>
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
        <v-btn :disabled="saving" variant="text" @click="closeCreateDialog">Cancel</v-btn>
        <v-btn
          color="primary"
          :disabled="!isFormValid"
          :loading="saving"
          @click="createVessel"
        >
          Add vessel
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import { useApiClient } from '@/composables/useApiClient'
  import {
    useFormatters,
    useVesselStatusFormatters,
    type VesselStatus,
  } from '@/composables/useFormatters'
  import {
    type Occupancy,
    useProductionApi,
    type Vessel,
  } from '@/composables/useProductionApi'
  import { useUnitPreferences, volumeOptions, type VolumeUnit } from '@/composables/useUnitPreferences'

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

  type BatchInfo = {
    uuid: string
    short_name: string
  }

  const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'
  const router = useRouter()
  const { request } = useApiClient(apiBase)
  const { getActiveOccupancies } = useProductionApi()
  const { formatVolumePreferred } = useUnitPreferences()
  const { formatRelativeTime } = useFormatters()
  const { formatVesselStatus, getVesselStatusColor } = useVesselStatusFormatters()

  const unitOptions = volumeOptions.map(opt => opt.value)
  const vesselStatusOptions = ['active', 'inactive', 'retired']

  // State
  const vessels = ref<Vessel[]>([])
  const occupancies = ref<Occupancy[]>([])
  const batches = ref<Batch[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const errorMessage = ref('')
  const search = ref('')

  // Dialog state
  const createVesselDialog = ref(false)

  // Form state
  const newVessel = reactive({
    type: '',
    name: '',
    capacity: '',
    capacity_unit: 'ml' as VolumeUnit,
    status: 'active',
    make: '',
    model: '',
  })

  const snackbar = reactive({
    show: false,
    text: '',
    color: 'success',
  })

  // Table configuration
  const headers = [
    { title: 'ID', key: 'id', sortable: true, width: '80px' },
    { title: 'Name', key: 'name', sortable: true },
    { title: 'Type', key: 'type', sortable: true },
    { title: 'Capacity', key: 'capacity', sortable: true },
    { title: 'Status', key: 'status', sortable: true },
    { title: 'Occupancy', key: 'occupancy', sortable: true },
    { title: 'Updated', key: 'updated_at', sortable: true },
  ]

  // Computed
  const isFormValid = computed(() => {
    return newVessel.type.trim().length > 0
      && newVessel.name.trim().length > 0
      && newVessel.capacity !== ''
  })

  // Map vessel_id -> occupancy for quick lookup
  const occupancyMap = computed(
    () => new Map(occupancies.value.map(occ => [occ.vessel_id, occ])),
  )

  // Map batch_id -> batch for quick lookup
  const batchMap = computed(
    () => new Map(batches.value.map(batch => [batch.id, batch])),
  )

  /**
   * Sort vessels by:
   * 1. Active vessels first
   * 2. Inactive vessels second
   * 3. Retired vessels last
   * 4. Within each status group: occupied before unoccupied
   * 5. Within each occupancy group: alphabetically by name
   */
  const sortedVessels = computed(() => {
    const statusOrder: Record<VesselStatus, number> = {
      active: 0,
      inactive: 1,
      retired: 2,
    }

    // eslint-disable-next-line unicorn/no-array-sort -- toSorted requires ES2023+
    return [...vessels.value].sort((a, b) => {
      // Sort by status first
      const statusDiff = statusOrder[a.status] - statusOrder[b.status]
      if (statusDiff !== 0) return statusDiff

      // Within same status, occupied vessels first
      const aOccupied = occupancyMap.value.has(a.id)
      const bOccupied = occupancyMap.value.has(b.id)
      if (aOccupied && !bOccupied) return -1
      if (!aOccupied && bOccupied) return 1

      // Within same occupancy group, sort alphabetically by name
      return a.name.localeCompare(b.name)
    })
  })

  // Lifecycle
  onMounted(async () => {
    await refreshData()
  })

  // Methods
  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  async function refreshData () {
    loading.value = true
    errorMessage.value = ''
    try {
      const [vesselData, occupancyData, batchData] = await Promise.all([
        request<Vessel[]>('/vessels'),
        getActiveOccupancies(),
        request<Batch[]>('/batches'),
      ])
      vessels.value = vesselData
      occupancies.value = occupancyData
      batches.value = batchData
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to load data'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      loading.value = false
    }
  }

  function getOccupancyBatchInfo (vesselId: number): BatchInfo | null {
    const occupancy = occupancyMap.value.get(vesselId)
    if (!occupancy || !occupancy.batch_id) return null

    const batch = batchMap.value.get(occupancy.batch_id)
    if (!batch) return null

    return {
      uuid: batch.uuid,
      short_name: batch.short_name,
    }
  }

  function openCreateDialog () {
    newVessel.type = ''
    newVessel.name = ''
    newVessel.capacity = ''
    newVessel.capacity_unit = 'ml'
    newVessel.status = 'active'
    newVessel.make = ''
    newVessel.model = ''
    createVesselDialog.value = true
  }

  function closeCreateDialog () {
    createVesselDialog.value = false
  }

  async function createVessel () {
    if (!isFormValid.value) return

    saving.value = true
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
      closeCreateDialog()
      await refreshData()
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Unable to create vessel'
      errorMessage.value = message
      showNotice(message, 'error')
    } finally {
      saving.value = false
    }
  }

  function onRowDoubleClick (_event: Event, { item }: { item: Vessel }) {
    router.push(`/vessels/${item.uuid}`)
  }

  // Formatting functions
  function normalizeText (value: string) {
    const trimmed = value.trim()
    return trimmed.length > 0 ? trimmed : null
  }

  function toNumber (value: string | number | null) {
    if (value === null || value === undefined || value === '') {
      return null
    }
    const parsed = Number(value)
    return Number.isFinite(parsed) ? parsed : null
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

.search-field {
  max-width: 260px;
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

.vessels-table :deep(tr) {
  cursor: pointer;
}

.vessels-table :deep(tr:hover td) {
  background: rgba(var(--v-theme-primary), 0.04);
}

.batch-link {
  color: rgb(var(--v-theme-primary));
  text-decoration: none;
}

.batch-link:hover {
  text-decoration: underline;
}
</style>
