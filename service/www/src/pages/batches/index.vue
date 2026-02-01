<script lang="ts" setup>
  import { onMounted } from 'vue'
  import { useRouter } from 'vue-router'

  const router = useRouter()

  onMounted(() => {
    router.replace('/batches/in-progress')
  })
</script>

<template>
  <div />
</template>
