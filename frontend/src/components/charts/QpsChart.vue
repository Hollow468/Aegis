<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent])

const props = defineProps<{ value: number }>()

const MAX_POINTS = 60
const times = ref<string[]>([])
const values = ref<number[]>([])

function addPoint() {
  const now = new Date()
  const label = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`
  times.value.push(label)
  values.value.push(props.value)
  if (times.value.length > MAX_POINTS) {
    times.value.shift()
    values.value.shift()
  }
}

let timer: ReturnType<typeof setInterval>
onMounted(() => { addPoint(); timer = setInterval(addPoint, 2000) })
onUnmounted(() => { clearInterval(timer) })

const option = ref({
  grid: { top: 10, right: 16, bottom: 24, left: 50 },
  tooltip: { trigger: 'axis' },
  xAxis: { type: 'category', data: times, axisLabel: { color: '#8b8fa3', fontSize: 10 }, axisLine: { lineStyle: { color: '#2e3148' } } },
  yAxis: { type: 'value', axisLabel: { color: '#8b8fa3', fontSize: 10 }, splitLine: { lineStyle: { color: '#2e3148' } } },
  series: [{ type: 'line', data: values, smooth: true, showSymbol: false, lineStyle: { color: '#3b82f6', width: 2 }, areaStyle: { color: { type: 'linear', x: 0, y: 0, x2: 0, y2: 1, colorStops: [{ offset: 0, color: 'rgba(59,130,246,0.3)' }, { offset: 1, color: 'rgba(59,130,246,0)' }] } } }]
})
</script>

<template>
  <v-chart :option="option" autoresize style="height: 200px;" />
</template>
