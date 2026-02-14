<template>
  <div />
</template>

<script lang="ts" setup>
  import { onBeforeMount } from 'vue'
  import { useRouter } from 'vue-router'

  const router = useRouter()

  // Use onBeforeMount instead of onMounted to redirect before render
  onBeforeMount(() => {
    router.replace('/procurement/purchase-orders')
  })
</script>
