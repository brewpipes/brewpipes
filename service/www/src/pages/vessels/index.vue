<script lang="ts" setup>
  import { onMounted } from 'vue'
  import { useRouter } from 'vue-router'

  const router = useRouter()

  onMounted(() => {
    router.replace('/vessels/active')
  })
</script>

<template>
  <div />
</template>
