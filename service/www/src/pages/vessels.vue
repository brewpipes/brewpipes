<template>
  <v-container class="vessels-page" fluid>
    <v-row>
      <v-col cols="12">
        <v-card class="hero-card">
          <v-card-text>
            <v-row align="center">
              <v-col cols="12" md="7">
                <div class="kicker">Production Service</div>
                <div class="text-h3 font-weight-bold mb-2">Vessel management</div>
                <div class="text-body-1 text-medium-emphasis">
                  Register tanks, kettles, and fermenters so production volumes can be anchored
                  to a physical location.
                </div>

                <div class="d-flex flex-wrap align-center ga-2 mt-4">
                  <v-chip color="primary" size="small" variant="tonal">
                    API: {{ apiBase }}
                  </v-chip>
                  <v-chip color="secondary" size="small" variant="tonal">
                    Vessels: {{ vessels.length }}
                  </v-chip>
                </div>
              </v-col>

              <v-col cols="12" md="5">
                <v-card class="hero-panel" variant="tonal">
                  <div class="text-overline">Active vessel</div>
                  <div class="text-h5 font-weight-semibold">
                    {{ selectedVessel ? selectedVessel.name : 'Select a vessel' }}
                  </div>
                  <div class="text-body-2 text-medium-emphasis mb-3">
                    {{ selectedVessel ? selectedVessel.type : 'Choose a vessel to review details' }}
                  </div>
                  <div class="d-flex flex-wrap ga-2">
                    <v-btn color="primary" size="small" :loading="loading" @click="refreshVessels">
                      Refresh list
                    </v-btn>
                    <v-btn size="small" variant="tonal" @click="clearSelection">
                      Clear selection
                    </v-btn>
                  </div>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-row class="mt-4" align="stretch">
      <v-col cols="12" md="4">
        <v-card class="section-card">
          <v-card-title class="d-flex align-center">
            <v-icon class="mr-2" icon="mdi-silo" />
            Vessel list
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

            <v-list class="vessel-list" lines="two" active-color="primary">
              <v-list-item
                v-for="vessel in vessels"
                :key="vessel.id"
                :active="vessel.id === selectedVesselId"
                @click="selectVessel(vessel.id)"
              >
                <v-list-item-title>
                  {{ vessel.name }}
                </v-list-item-title>
                <v-list-item-subtitle>
                  {{ vessel.type }} • {{ formatCapacity(vessel.capacity, vessel.capacity_unit) }}
                </v-list-item-subtitle>
                <template #append>
                  <v-chip size="x-small" variant="tonal">
                    {{ vessel.status }}
                  </v-chip>
                </template>
              </v-list-item>

              <v-list-item v-if="vessels.length === 0">
                <v-list-item-title>No vessels yet</v-list-item-title>
                <v-list-item-subtitle>Register one using the form.</v-list-item-subtitle>
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
                        {{ selectedVessel.type }} • {{ selectedVessel.status }}
                      </div>
                      <div class="text-body-2 text-medium-emphasis">
                        Capacity {{ formatCapacity(selectedVessel.capacity, selectedVessel.capacity_unit) }}
                      </div>
                      <div class="text-body-2 text-medium-emphasis" v-if="selectedVessel.make || selectedVessel.model">
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
            </div>

            <v-divider class="my-6" />

            <div class="text-subtitle-1 font-weight-semibold mb-2">Register vessel</div>
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
            <v-btn
              color="primary"
              :disabled="!newVessel.type.trim() || !newVessel.name.trim() || !newVessel.capacity"
              @click="createVessel"
            >
              Add vessel
            </v-btn>
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
import { computed, onMounted, reactive, ref } from 'vue'

type Unit = 'ml' | 'usfloz' | 'ukfloz'

type Vessel = {
  id: number
  uuid: string
  type: string
  name: string
  capacity: number
  capacity_unit: Unit
  make: string | null
  model: string | null
  status: 'active' | 'inactive' | 'retired'
  created_at: string
  updated_at: string
}

const apiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'

const unitOptions: Unit[] = ['ml', 'usfloz', 'ukfloz']
const vesselStatusOptions = ['active', 'inactive', 'retired']

const vessels = ref<Vessel[]>([])
const selectedVesselId = ref<number | null>(null)
const errorMessage = ref('')
const loading = ref(false)

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

const newVessel = reactive({
  type: '',
  name: '',
  capacity: '',
  capacity_unit: 'ml' as Unit,
  status: 'active',
  make: '',
  model: '',
})

const selectedVessel = computed(() =>
  vessels.value.find((vessel) => vessel.id === selectedVesselId.value) ?? null,
)

onMounted(async () => {
  await refreshVessels()
})

function selectVessel(id: number) {
  selectedVesselId.value = id
}

function clearSelection() {
  selectedVesselId.value = null
}

function showNotice(text: string, color = 'success') {
  snackbar.text = text
  snackbar.color = color
  snackbar.show = true
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`${apiBase}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(init.headers ?? {}),
    },
  })

  if (!response.ok) {
    const message = await response.text()
    throw new Error(message || `Request failed with ${response.status}`)
  }

  return response.json() as Promise<T>
}

async function refreshVessels() {
  loading.value = true
  errorMessage.value = ''
  try {
    vessels.value = await request<Vessel[]>('/vessels')
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to load vessels'
    errorMessage.value = message
    showNotice(message, 'error')
  } finally {
    loading.value = false
  }
}

async function createVessel() {
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
    newVessel.type = ''
    newVessel.name = ''
    newVessel.capacity = ''
    newVessel.make = ''
    newVessel.model = ''
    await refreshVessels()
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unable to create vessel'
    errorMessage.value = message
    showNotice(message, 'error')
  }
}

function normalizeText(value: string) {
  const trimmed = value.trim()
  return trimmed.length > 0 ? trimmed : null
}

function toNumber(value: string | number | null) {
  if (value === null || value === undefined || value === '') {
    return null
  }
  const parsed = Number(value)
  return Number.isFinite(parsed) ? parsed : null
}

function formatCapacity(amount: number, unit: Unit) {
  return `${amount} ${unit}`
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
.vessels-page {
  position: relative;
}

.hero-card {
  border: 1px solid rgba(47, 93, 80, 0.16);
  background:
    linear-gradient(130deg, rgba(47, 93, 80, 0.14), rgba(196, 117, 60, 0.12)),
    rgba(253, 251, 247, 0.9);
  box-shadow: 0 12px 28px rgba(28, 26, 22, 0.08);
}

.hero-panel {
  border: 1px solid rgba(196, 117, 60, 0.2);
}

.kicker {
  text-transform: uppercase;
  letter-spacing: 0.24em;
  font-size: 0.7rem;
  color: rgba(28, 26, 22, 0.6);
  margin-bottom: 6px;
}

.section-card {
  background: rgba(253, 251, 247, 0.92);
  border: 1px solid rgba(28, 26, 22, 0.08);
  box-shadow: 0 10px 24px rgba(28, 26, 22, 0.05);
}

.sub-card {
  border: 1px solid rgba(28, 26, 22, 0.08);
  background: rgba(255, 255, 255, 0.72);
}

.vessel-list {
  max-height: 420px;
  overflow: auto;
}
</style>
