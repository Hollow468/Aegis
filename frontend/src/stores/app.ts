import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const currentPage = ref('dashboard')

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setPage(page: string) {
    currentPage.value = page
  }

  return { sidebarCollapsed, currentPage, toggleSidebar, setPage }
})
