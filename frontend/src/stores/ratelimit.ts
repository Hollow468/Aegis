import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as rlApi from '../api/ratelimit'
import type { RateLimitRule } from '../api/ratelimit'

export const useRateLimitStore = defineStore('ratelimit', () => {
  const limits = ref<RateLimitRule[]>([])
  const loading = ref(false)

  async function fetchLimits() {
    loading.value = true
    try {
      limits.value = await rlApi.getRateLimits()
    } finally {
      loading.value = false
    }
  }

  async function addLimit(rule: RateLimitRule) {
    await rlApi.addRateLimit(rule)
    await fetchLimits()
  }

  async function removeLimit(route: string) {
    await rlApi.deleteRateLimit(route)
    await fetchLimits()
  }

  return { limits, loading, fetchLimits, addLimit, removeLimit }
})
