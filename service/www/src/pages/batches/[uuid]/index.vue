<template>
  <v-container class="pa-4" fluid>
    <v-alert
      v-if="loading"
      density="comfortable"
      type="info"
      variant="tonal"
    >
      Loading batch...
    </v-alert>

    <v-alert
      v-else-if="error"
      density="comfortable"
      type="error"
      variant="tonal"
    >
      {{ error }}
    </v-alert>

    <BatchDetails
      v-else
      back-button-route="/batches/all"
      back-button-text="Back to All Batches"
      :batch-id="batchId"
      :show-back-button="true"
      @back="handleBack"
    />
  </v-container>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import BatchDetails from '@/components/BatchDetails.vue'
  import { useApiClient } from '@/composables/useApiClient'

  type Batch = {
    id: number
    uuid: string
    short_name: string
    brew_date: string | null
    recipe_id: number | null
    recipe_uuid: string | null
    notes: string | null
    created_at: string
    updated_at: string
  }

  const route = useRoute()
  const router = useRouter()

  const productionApiBase = import.meta.env.VITE_PRODUCTION_API_URL ?? '/api'
  const { request } = useApiClient(productionApiBase)

  const loading = ref(true)
  const error = ref<string | null>(null)
  const batch = ref<Batch | null>(null)

  const batchId = computed(() => batch.value?.id ?? null)

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

  async function loadBatch () {
    const uuid = routeUuid.value
    if (!uuid) {
      error.value = 'Invalid batch UUID'
      loading.value = false
      return
    }

    try {
      loading.value = true
      error.value = null

      // Fetch all batches and find the one matching the UUID
      const batches = await request<Batch[]>('/batches')
      const found = batches.find(b => b.uuid === uuid)

      if (found) {
        batch.value = found
      } else {
        error.value = 'Batch not found'
      }
    } catch (error_) {
      console.error('Failed to load batch:', error_)
      error.value = 'Failed to load batch. Please try again.'
    } finally {
      loading.value = false
    }
  }

  function handleBack () {
    router.push('/batches/all')
  }

  onMounted(() => {
    loadBatch()
  })
</script>
