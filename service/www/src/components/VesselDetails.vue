<template>
  <v-card class="section-card">
    <v-card-title class="d-flex align-center">
      <v-icon class="mr-2" icon="mdi-silo-outline" />
      {{ vessel ? vessel.name : 'Vessel details' }}
      <v-spacer />
      <v-btn v-if="vessel" size="small" variant="text" @click="emit('edit')">
        <v-icon class="mr-1" icon="mdi-pencil" size="small" />
        Edit
      </v-btn>
      <v-btn :loading="loading" size="small" variant="text" @click="emit('refresh')">Refresh</v-btn>
      <v-btn v-if="vessel" size="small" variant="text" @click="emit('clear')">Clear</v-btn>
    </v-card-title>
    <v-card-text>
      <v-progress-linear v-if="loading" class="mb-3" color="primary" indeterminate />

      <v-alert
        v-if="!vessel && !loading"
        density="comfortable"
        type="info"
        variant="tonal"
      >
        Select a vessel to view metadata and availability.
      </v-alert>

      <div v-else-if="vessel">
        <v-row>
          <v-col cols="12" md="6">
            <v-card class="sub-card" variant="tonal">
              <v-card-text>
                <div class="text-overline">Vessel</div>
                <div class="text-h5 font-weight-semibold">
                  {{ vessel.name }}
                </div>
                <div class="text-body-2 text-medium-emphasis">
                  {{ vessel.type }} - {{ vessel.status }}
                </div>
                <div class="text-body-2 text-medium-emphasis">
                  Capacity {{ formatVolumePreferred(vessel.capacity, vessel.capacity_unit) }}
                </div>
                <div v-if="vessel.make || vessel.model" class="text-body-2 text-medium-emphasis">
                  {{ vessel.make ?? 'Make not set' }} {{ vessel.model ?? '' }}
                </div>
              </v-card-text>
            </v-card>
          </v-col>

          <v-col cols="12" md="6">
            <v-card class="sub-card" variant="tonal">
              <v-card-text>
                <div class="text-overline">Metadata</div>
                <div class="text-body-2 text-medium-emphasis">UUID {{ vessel.uuid }}</div>
                <div class="text-body-2 text-medium-emphasis">
                  Updated {{ formatDateTime(vessel.updated_at) }}
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

                <div v-if="occupancy">
                  <v-row align="center" dense>
                    <v-col cols="12" md="4">
                      <div class="text-caption text-medium-emphasis">Status</div>
                      <v-menu location="bottom">
                        <template #activator="{ props }">
                          <v-chip
                            v-bind="props"
                            append-icon="mdi-menu-down"
                            class="mt-1 cursor-pointer"
                            :color="getOccupancyStatusColor(occupancy.status)"
                            size="small"
                            variant="tonal"
                          >
                            {{ formatOccupancyStatus(occupancy.status) }}
                          </v-chip>
                        </template>
                        <v-list density="compact" nav>
                          <v-list-subheader>Change status</v-list-subheader>
                          <v-list-item
                            v-for="statusOption in occupancyStatusOptions"
                            :key="statusOption.value"
                            :active="statusOption.value === occupancy.status"
                            @click="handleStatusChange(occupancy.uuid, statusOption.value)"
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
                        {{ formatDateTime(occupancy.in_at) }}
                      </div>
                    </v-col>
                    <v-col cols="12" md="4">
                      <div class="text-caption text-medium-emphasis">Volume</div>
                      <div class="text-body-2 font-weight-medium mt-1">
                        {{ occupancy.volume_uuid.substring(0, 8) }}
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
</template>

<script lang="ts" setup>
  import { computed } from 'vue'
  import { useFormatters, useOccupancyStatusFormatters } from '@/composables/useFormatters'
  import { useUnitPreferences } from '@/composables/useUnitPreferences'
  import { type Occupancy, OCCUPANCY_STATUS_VALUES, type OccupancyStatus, type Vessel } from '@/types'

  withDefaults(
    defineProps<{
      vessel: Vessel | null
      occupancy: Occupancy | null
      loading?: boolean
    }>(),
    {
      loading: false,
    },
  )

  const emit = defineEmits<{
    'occupancy-status-change': [occupancyUuid: string, status: OccupancyStatus]
    'edit': []
    'refresh': []
    'clear': []
  }>()

  const { formatVolumePreferred } = useUnitPreferences()
  const { formatDateTime } = useFormatters()
  const {
    formatOccupancyStatus,
    getOccupancyStatusColor,
    getOccupancyStatusIcon,
  } = useOccupancyStatusFormatters()

  const occupancyStatusOptions = computed(() =>
    OCCUPANCY_STATUS_VALUES.map(status => ({
      value: status,
      title: formatOccupancyStatus(status),
    })),
  )

  function handleStatusChange (occupancyUuid: string, status: OccupancyStatus) {
    emit('occupancy-status-change', occupancyUuid, status)
  }
</script>

<style scoped>
.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}

.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}

.cursor-pointer {
  cursor: pointer;
}
</style>
