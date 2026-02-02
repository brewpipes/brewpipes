<template>
  <div ref="container" class="sankey-diagram">
    <div v-if="!hasData" class="sankey-empty">
      <slot name="empty">No flow data available yet.</slot>
    </div>
    <div v-else-if="!layout" class="sankey-empty">Sizing diagram...</div>
    <svg
      v-else
      :height="height"
      role="img"
      :viewBox="`0 0 ${layoutWidth} ${height}`"
      :width="layoutWidth"
    >
      <g class="sankey-links">
        <path
          v-for="link in layout.links"
          :key="linkKey(link)"
          class="sankey-link"
          :d="linkPath(link) ?? ''"
          :style="{ strokeWidth: `${Math.max(link.width ?? 1, 1)}px` }"
        >
          <title>{{ linkTitle(link) }}</title>
        </path>
      </g>
      <g class="sankey-nodes">
        <g
          v-for="node in layout.nodes"
          :key="node.id"
          :transform="`translate(${node.x0 ?? 0}, ${node.y0 ?? 0})`"
        >
          <rect
            class="sankey-node"
            :fill="nodeColor(node)"
            :height="nodeHeight(node)"
            rx="6"
            ry="6"
            :width="nodeWidth(node)"
          />
          <text
            class="sankey-label"
            dominant-baseline="middle"
            :text-anchor="nodeLabelAnchor(node)"
            :x="nodeLabelX(node)"
            :y="nodeHeight(node) / 2"
          >
            {{ node.label }}
          </text>
          <title>{{ `${node.label} (${formatValue(node.value)})` }}</title>
        </g>
      </g>
    </svg>
  </div>
</template>

<script lang="ts" setup>
  import { sankey, sankeyLinkHorizontal } from 'd3-sankey'
  import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

  type SankeyNodeInput = {
    id: string
    label: string
  }

  type SankeyLinkInput = {
    source: string
    target: string
    value: number
    label?: string
  }

  type SankeyNodeDatum = SankeyNodeInput & {
    index?: number
    value?: number
    x0?: number
    x1?: number
    y0?: number
    y1?: number
  }

  type SankeyLinkDatum = Omit<SankeyLinkInput, 'source' | 'target'> & {
    index?: number
    width?: number
    source: SankeyNodeDatum
    target: SankeyNodeDatum
  }

  type SankeyLayout = {
    nodes: SankeyNodeDatum[]
    links: SankeyLinkDatum[]
  }

  const props = withDefaults(
    defineProps<{
      nodes: SankeyNodeInput[]
      links: SankeyLinkInput[]
      height?: number
      nodeWidth?: number
      nodePadding?: number
    }>(),
    {
      height: 360,
      nodeWidth: 14,
      nodePadding: 16,
    },
  )

  const container = ref<HTMLElement | null>(null)
  const width = ref(0)
  let observer: ResizeObserver | null = null

  const hasData = computed(() => props.nodes.length > 0 && props.links.length > 0)
  const layoutWidth = computed(() => Math.max(width.value, 480))

  const layout = computed(() => {
    if (!hasData.value || layoutWidth.value === 0) {
      return null
    }

    const layoutEngine = sankey<SankeyNodeDatum, SankeyLinkInput>()
      .nodeId((node: SankeyNodeDatum) => node.id)
      .nodeWidth(props.nodeWidth)
      .nodePadding(props.nodePadding)
      .extent([
        [0, 0],
        [layoutWidth.value, props.height],
      ])

    return layoutEngine({
      nodes: props.nodes.map(node => ({ ...node })),
      links: props.links.map(link => ({ ...link })),
    }) as unknown as SankeyLayout
  })

  const linkPath = sankeyLinkHorizontal<SankeyLinkDatum>()
  const nodePalette = [
    'rgb(var(--v-theme-primary))',
    'rgb(var(--v-theme-secondary))',
    'rgb(var(--v-theme-info))',
    'rgb(var(--v-theme-success))',
    'rgb(var(--v-theme-warning))',
  ]

  onMounted(() => {
    if (!container.value) {
      return
    }
    width.value = container.value.clientWidth
    observer = new ResizeObserver(entries => {
      const entry = entries[0]
      if (!entry) {
        return
      }
      width.value = entry.contentRect.width
    })
    observer.observe(container.value)
  })

  onBeforeUnmount(() => {
    observer?.disconnect()
  })

  function nodeWidth (node: SankeyNodeDatum) {
    return Math.max((node.x1 ?? 0) - (node.x0 ?? 0), 0)
  }

  function nodeHeight (node: SankeyNodeDatum) {
    return Math.max((node.y1 ?? 0) - (node.y0 ?? 0), 0)
  }

  function nodeLabelX (node: SankeyNodeDatum) {
    const widthValue = nodeWidth(node)
    return (node.x0 ?? 0) < layoutWidth.value / 2 ? widthValue + 8 : -8
  }

  function nodeLabelAnchor (node: SankeyNodeDatum) {
    return (node.x0 ?? 0) < layoutWidth.value / 2 ? 'start' : 'end'
  }

  function nodeColor (node: SankeyNodeDatum) {
    const index = node.index ?? 0
    return nodePalette[index % nodePalette.length]
  }

  function formatValue (value: number | undefined) {
    if (value === undefined) {
      return 'n/a'
    }
    return new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(value)
  }

  function linkTitle (link: SankeyLinkDatum) {
    return link.label ?? formatValue(link.value)
  }

  function linkKey (link: SankeyLinkDatum) {
    if (link.index !== undefined) {
      return `link-${link.index}`
    }
    return `link-${link.source.id}-${link.target.id}`
  }
</script>

<style scoped>
.sankey-diagram {
  width: 100%;
}

.sankey-diagram svg {
  display: block;
  max-width: 100%;
  overflow: visible;
}

.sankey-empty {
  padding: 12px 0;
  color: rgba(var(--v-theme-on-surface), 0.6);
  font-size: 0.9rem;
}

.sankey-link {
  fill: none;
  stroke: rgba(var(--v-theme-primary), 0.35);
  opacity: 0.85;
}

.sankey-link:hover {
  stroke: rgba(var(--v-theme-primary), 0.6);
  opacity: 1;
}

.sankey-node {
  stroke: rgba(var(--v-theme-on-surface), 0.16);
}

.sankey-label {
  fill: rgba(var(--v-theme-on-surface), 0.7);
  font-size: 0.75rem;
}
</style>
