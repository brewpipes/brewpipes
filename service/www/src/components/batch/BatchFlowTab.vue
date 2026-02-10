<template>
  <v-card class="sub-card" variant="outlined">
    <v-card-title class="text-subtitle-1">Liquid flow</v-card-title>
    <v-card-text>
      <v-alert
        v-if="notice"
        class="mb-3"
        density="compact"
        type="info"
        variant="tonal"
      >
        {{ notice }}
      </v-alert>

      <SankeyDiagram v-if="links.length > 0" :links="links" :nodes="nodes" />
      <div v-else class="text-body-2 text-medium-emphasis">
        No flow relations yet. Record a split or blend to visualize liquid movement.
      </div>

      <div class="text-caption text-medium-emphasis mt-3">
        Flow is derived from volume relations (splits and blends).
      </div>
    </v-card-text>
  </v-card>
</template>

<script lang="ts" setup>
  import type { FlowLink, FlowNode } from './types'
  import SankeyDiagram from '@/components/SankeyDiagram.vue'

  defineProps<{
    nodes: FlowNode[]
    links: FlowLink[]
    notice: string
  }>()
</script>

<style scoped>
.sub-card {
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  background: rgba(var(--v-theme-surface), 0.7);
}
</style>
