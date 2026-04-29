import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  async function login(username: string, password: string) {
    loading.value = true
    try {
      const res = await authApi.login({ username, password })
      token.value = res.token
      localStorage.setItem('token', res.token)
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
  }

  return { token, loading, isAuthenticated, login, logout }
})
