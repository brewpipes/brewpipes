<template>
  <div class="procurement-redirect" />
</template>

<script lang="ts" setup>
  import { onMounted } from 'vue'
  import { useRouter } from 'vue-router'

  const router = useRouter()

  onMounted(() => {
    router.replace('/procurement/purchase-orders')
  })
</script>

<style scoped>
.procurement-redirect {
  min-height: 200px;
}
</style>
