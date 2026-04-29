<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent } from 'echarts/components'

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent])

const props = defineProps<{ s2xx: number; s4xx: number; s5xx: number }>()

const option = computed(() => ({
  tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
  legend: { bottom: 0, textStyle: { color: '#8b8fa3', fontSize: 11 } },
  series: [{
    type: 'pie', radius: ['40%', '70%'], center: ['50%', '45%'],
    label: { show: false },
    data: [
      { value: props.s2xx, name: '2xx', itemStyle: { color: '#22c55e' } },
      { value: props.s4xx, name: '4xx', itemStyle: { color: '#f59e0b' } },
      { value: props.s5xx, name: '5xx', itemStyle: { color: '#ef4444' } }
    ]
  }]
}))
</script>

<template>
  <v-chart :option="option" autoresize style="height: 220px;" />
</template>
