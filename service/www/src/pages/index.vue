<template>
  <v-container class="dashboard-page" fluid>
    <v-row align="stretch">
      <v-col cols="12" md="7">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-barley" />
            Recent batches
          </v-card-title>
          <v-card-text>
            <v-list lines="two">
              <v-list-item
                v-for="batch in recentBatches"
                :key="batch.id"
                to="/batches"
              >
                <v-list-item-title>{{ batch.short_name }}</v-list-item-title>
                <v-list-item-subtitle>
                  #{{ batch.id }} - {{ formatDate(batch.brew_date) }}
                </v-list-item-subtitle>
                <template #append>
                  <v-chip size="x-small" variant="tonal">
                    {{ formatDateTime(batch.updated_at) }}
                  </v-chip>
                </template>
              </v-list-item>

              <v-list-item v-if="recentBatches.length === 0">
                <v-list-item-title>No batches yet</v-list-item-title>
                <v-list-item-subtitle>Create the first batch to get started.</v-list-item-subtitle>
              </v-list-item>
            </v-list>

            <v-btn class="mt-2" variant="text" to="/batches">View all batches</v-btn>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="5">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-lightning-bolt-outline" />
            Quick actions
          </v-card-title>
          <v-card-text class="d-flex flex-column ga-2">
            <v-btn block color="primary" to="/batches">Create batch</v-btn>
            <v-btn block variant="tonal" to="/batches">Manage batch workflow</v-btn>
            <v-btn :loading="loading" variant="text" @click="refreshAll">Refresh</v-btn>
            <v-alert
              v-if="errorMessage"
              density="compact"
              type="error"
              variant="tonal"
            >
              {{ errorMessage }}
            </v-alert>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="mt-4" align="stretch">
      <v-col cols="12" md="6">
        <v-card class="section-card">
          <v-card-title>Vessels</v-card-title>
          <v-card-text>
            <div class="text-body-2 text-medium-emphasis">
              Manage vessel inventory and capacity planning.
            </div>
            <v-btn class="mt-3" variant="tonal" to="/vessels">Open vessels</v-btn>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card class="section-card">
          <v-card-title>Workflow</v-card-title>
          <v-card-text>
            <div class="text-body-2 text-medium-emphasis">
              Capture transfers, additions, and measurements in one place.
            </div>
            <v-btn class="mt-3" variant="tonal" to="/batches">Go to workflow</v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'

type Batch = {
  id: number
  short_name: string
  brew_date: string | null
  updated_at: string
}

type Vessel = {
  id: number
}

type Volume = {
  id: number
}

const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

const batches = ref<Batch[]>([])
const vessels = ref<Vessel[]>([])
const volumes = ref<Volume[]>([])
const errorMessage = ref('')
const loading = ref(false)

const recentBatches = computed(() =>
  [...batches.value]
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, 5),
)

onMounted(async () => {
  await refreshAll()
})

async function request<T>(path: string): Promise<T> {
  const response = await fetch(`${apiBase}${path}`)
  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || `Request failed with ${response.status}`)
  }
  return response.json() as Promise<T>
}

async function refreshAll() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [batchData, vesselData, volumeData] = await Promise.all([
      request<Batch[]>('/batches'),
      request<Vessel[]>('/vessels'),
      request<Volume[]>('/volumes'),
    ])
    batches.value = batchData
    vessels.value = vesselData
    volumes.value = volumeData
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load dashboard'
    errorMessage.value = message
  } finally {
    loading.value = false
  }
}

function formatDate(value: string | null | undefined) {
  if (!value) {
    return 'Unknown'
  }
  return new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
  }).format(new Date(value))
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
</script>

<style scoped>
.dashboard-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}
</style>
