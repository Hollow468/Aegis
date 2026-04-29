import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as settingsApi from '../api/settings'
import type { SystemSettings } from '../api/settings'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<SystemSettings | null>(null)
  const loading = ref(false)

  async function fetchSettings() {
    loading.value = true
    try {
      settings.value = await settingsApi.getSettings()
    } finally {
      loading.value = false
    }
  }

  async function updateSettings(partial: Partial<SystemSettings>) {
    await settingsApi.updateSettings(partial)
    await fetchSettings()
  }

  return { settings, loading, fetchSettings, updateSettings }
})
