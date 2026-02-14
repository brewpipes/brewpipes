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
      :batch-uuid="batchUuid"
      :show-back-button="true"
      @back="handleBack"
    />
  </v-container>
</template>

<script lang="ts" setup>
  import type { Batch } from '@/types'
  import { computed, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'
  import BatchDetails from '@/components/BatchDetails.vue'
  import { useProductionApi } from '@/composables/useProductionApi'
  import { useRouteUuid } from '@/composables/useRouteUuid'

  const router = useRouter()

  const { getBatch } = useProductionApi()
  const { uuid: routeUuid } = useRouteUuid()

  const loading = ref(true)
  const error = ref<string | null>(null)
  const batch = ref<Batch | null>(null)

  const batchUuid = computed(() => batch.value?.uuid ?? null)

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

      batch.value = await getBatch(uuid)
    } catch (error_) {
      error.value = error_ instanceof Error && error_.message.includes('404') ? 'Batch not found' : 'Failed to load batch. Please try again.'
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
