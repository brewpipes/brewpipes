<template>
  <div class="spark-card" :style="{ '--spark-color': `var(--v-theme-${color})` }">
    <div class="spark-meta">
      <div class="text-caption text-medium-emphasis">{{ label }}</div>
      <div class="text-body-2 font-weight-medium">
        {{ latestLabel }}
      </div>
    </div>
    <div class="spark-chart">
      <svg
        v-if="values.length > 0"
        preserveAspectRatio="none"
        :viewBox="`0 0 ${width} ${height}`"
      >
        <path class="spark-area" :d="areaPath" />
        <path class="spark-line" :d="linePath" />
      </svg>
      <div v-else class="spark-placeholder">No readings</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { computed } from 'vue'

  const props = defineProps<{
    label: string
    latestLabel: string
    values: number[]
    color: string
    width?: number
    height?: number
  }>()

  const width = computed(() => props.width ?? 120)
  const height = computed(() => props.height ?? 36)

  const sparkline = computed(() => buildSparkline(props.values, width.value, height.value))
  const linePath = computed(() => sparkline.value.linePath)
  const areaPath = computed(() => sparkline.value.areaPath)

  function buildSparkline (values: number[], w: number, h: number) {
    if (values.length === 0) {
      return { linePath: '', areaPath: '' }
    }
    const min = Math.min(...values)
    const max = Math.max(...values)
    const range = max - min
    const step = values.length > 1 ? w / (values.length - 1) : w
    const points = values.map((value, index) => {
      const ratio = range === 0 ? 0.5 : (value - min) / range
      const x = index * step
      const y = h - ratio * h
      return { x, y }
    })
    const linePath = points
      .map((point, index) => `${index === 0 ? 'M' : 'L'} ${point.x} ${point.y}`)
      .join(' ')
    const lastPoint = points.at(-1)
    const firstPoint = points[0]
    if (!lastPoint || !firstPoint) {
      return { linePath, areaPath: '' }
    }
    const areaPath = `${linePath} L ${lastPoint.x} ${h} L ${firstPoint.x} ${h} Z`
    return { linePath, areaPath }
  }
</script>

<style scoped>
.spark-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 12px;
  border-radius: 12px;
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  background: rgba(var(--v-theme-surface), 0.4);
}

.spark-meta {
  min-width: 86px;
}

.spark-chart {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.spark-chart svg {
  width: 100%;
  height: 36px;
}

.spark-placeholder {
  font-size: 0.75rem;
  color: rgba(var(--v-theme-on-surface), 0.5);
}

.spark-line {
  fill: none;
  stroke: rgb(var(--spark-color));
  stroke-width: 2;
}

.spark-area {
  fill: rgba(var(--spark-color), 0.2);
}
</style>
