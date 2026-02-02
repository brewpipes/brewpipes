<template>
  <v-card class="section-card">
    <v-card-title class="d-flex align-center">
      <v-icon class="mr-2" icon="mdi-silo" />
      Active Vessels
      <v-spacer />
      <v-btn
        :loading="loading"
        size="small"
        variant="text"
        @click="emit('refresh')"
      >
        Refresh
      </v-btn>
    </v-card-title>
    <v-card-text>
      <v-progress-linear v-if="loading" class="mb-3" color="primary" indeterminate />

      <v-list active-color="primary" class="vessel-list" lines="two">
        <v-list-item
          v-for="vessel in sortedVessels"
          :key="vessel.id"
          :active="vessel.id === selectedVesselId"
          @click="emit('select', vessel.id)"
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

        <v-list-item v-if="!loading && vessels.length === 0">
          <v-list-item-title>No active vessels</v-list-item-title>
          <v-list-item-subtitle>Register vessels in All Vessels to see them here.</v-list-item-subtitle>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { Occupancy, Vessel } from '@/composables/useProductionApi'
  import { computed } from 'vue'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'

  const props = withDefaults(
    defineProps<{
      vessels: Vessel[]
      occupancies: Occupancy[]
      selectedVesselId: number | null
      loading?: boolean
    }>(),
    {
      loading: false,
    },
  )

  const emit = defineEmits<{
    select: [vesselId: number]
    refresh: []
  }>()

  const { formatVolumePreferred } = useUnitPreferences()

  const occupancyMap = computed(
    () => new Map(props.occupancies.map(occupancy => [occupancy.vessel_id, occupancy])),
  )

  function isVesselOccupied (vesselId: number): boolean {
    return occupancyMap.value.has(vesselId)
  }

  // Sort: occupied vessels first, then alphabetically by name
  const sortedVessels = computed(() => {
    return props.vessels.toSorted((a, b) => {
      const aOccupied = isVesselOccupied(a.id)
      const bOccupied = isVesselOccupied(b.id)

      // Occupied vessels first
      if (aOccupied && !bOccupied) return -1
      if (!aOccupied && bOccupied) return 1

      // Within same occupancy group, sort alphabetically by name
      return a.name.localeCompare(b.name)
    })
  })
</script>

<style scoped>
.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.vessel-list {
  max-height: 60vh;
  overflow-y: auto;
}
</style>
