<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent])

const props = defineProps<{ avg: number; p99: number }>()

const MAX_POINTS = 60
const times = ref<string[]>([])
const avgData = ref<number[]>([])
const p99Data = ref<number[]>([])

function addPoint() {
  const now = new Date()
  const label = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`
  times.value.push(label)
  avgData.value.push(props.avg)
  p99Data.value.push(props.p99)
  if (times.value.length > MAX_POINTS) {
    times.value.shift()
    avgData.value.shift()
    p99Data.value.shift()
  }
}

let timer: ReturnType<typeof setInterval>
onMounted(() => { addPoint(); timer = setInterval(addPoint, 2000) })
onUnmounted(() => { clearInterval(timer) })

const option = ref({
  grid: { top: 30, right: 16, bottom: 24, left: 50 },
  tooltip: { trigger: 'axis' },
  legend: { data: ['Avg', 'P99'], textStyle: { color: '#8b8fa3', fontSize: 11 }, top: 0 },
  xAxis: { type: 'category', data: times, axisLabel: { color: '#8b8fa3', fontSize: 10 }, axisLine: { lineStyle: { color: '#2e3148' } } },
  yAxis: { type: 'value', name: 'ms', axisLabel: { color: '#8b8fa3', fontSize: 10 }, splitLine: { lineStyle: { color: '#2e3148' } } },
  series: [
    { name: 'Avg', type: 'line', data: avgData, smooth: true, showSymbol: false, lineStyle: { color: '#22c55e', width: 2 } },
    { name: 'P99', type: 'line', data: p99Data, smooth: true, showSymbol: false, lineStyle: { color: '#f59e0b', width: 2 } }
  ]
})
</script>

<template>
  <v-chart :option="option" autoresize style="height: 200px;" />
</template>
