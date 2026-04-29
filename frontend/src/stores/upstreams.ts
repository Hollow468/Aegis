import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { UpstreamStatus } from '../types'
import * as upstreamsApi from '../api/upstreams'

export const useUpstreamsStore = defineStore('upstreams', () => {
  const upstreams = ref<UpstreamStatus[]>([])
  const loading = ref(false)

  async function fetchUpstreams() {
    loading.value = true
    try {
      upstreams.value = await upstreamsApi.getUpstreams()
    } finally {
      loading.value = false
    }
  }

  async function register(service: string, address: string, weight: number) {
    await upstreamsApi.registerUpstream(service, address, weight)
    await fetchUpstreams()
  }

  async function deregister(service: string, address: string) {
    await upstreamsApi.deregisterUpstream(service, address)
    await fetchUpstreams()
  }

  return { upstreams, loading, fetchUpstreams, register, deregister }
})
