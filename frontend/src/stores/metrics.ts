import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { MetricsSummary } from '../types'
import * as metricsApi from '../api/metrics'

export const useMetricsStore = defineStore('metrics', () => {
  const summary = ref<MetricsSummary>({
    qps: 0,
    avg_latency: 0,
    p99_latency: 0,
    total_requests: 0,
    in_flight: 0,
    status_2xx: 0,
    status_4xx: 0,
    status_5xx: 0
  })
  const loading = ref(false)

  async function fetchSummary() {
    loading.value = true
    try {
      summary.value = await metricsApi.getMetricsSummary()
    } finally {
      loading.value = false
    }
  }

  return { summary, loading, fetchSummary }
})
