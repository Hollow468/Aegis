import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Route } from '../types'
import * as routesApi from '../api/routes'

export const useRoutesStore = defineStore('routes', () => {
  const routes = ref<Route[]>([])
  const loading = ref(false)
  const total = ref(0)

  async function fetchRoutes() {
    loading.value = true
    try {
      const res = await routesApi.getRoutes()
      routes.value = res.routes
      total.value = res.total
    } finally {
      loading.value = false
    }
  }

  async function addRoute(route: Route) {
    await routesApi.createRoute(route)
    await fetchRoutes()
  }

  async function removeRoute(path: string) {
    await routesApi.deleteRoute(path)
    await fetchRoutes()
  }

  return { routes, loading, total, fetchRoutes, addRoute, removeRoute }
})
