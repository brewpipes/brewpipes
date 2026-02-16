<template>
  <v-container class="pa-4" fluid>
    <v-alert
      v-if="error"
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
  import { computed } from 'vue'
  import { useRouter } from 'vue-router'
  import BatchDetails from '@/components/BatchDetails.vue'
  import { useRouteUuid } from '@/composables/useRouteUuid'

  const router = useRouter()
  const { uuid: routeUuid } = useRouteUuid()

  const error = computed(() => routeUuid.value ? null : 'Invalid batch UUID')
  const batchUuid = computed(() => routeUuid.value ?? null)

  function handleBack () {
    router.push('/batches/all')
  }
</script>
