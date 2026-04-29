import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { CircuitBreakerStatus } from '../types'
import * as circuitApi from '../api/circuit'

export const useCircuitStore = defineStore('circuit', () => {
  const breakers = ref<CircuitBreakerStatus[]>([])
  const loading = ref(false)

  async function fetchBreakers() {
    loading.value = true
    try {
      breakers.value = await circuitApi.getCircuitBreakers()
    } finally {
      loading.value = false
    }
  }

  async function resetBreaker(route: string) {
    await circuitApi.resetCircuitBreaker(route)
    await fetchBreakers()
  }

  return { breakers, loading, fetchBreakers, resetBreaker }
})
