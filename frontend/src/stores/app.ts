import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const currentPage = ref('dashboard')
  const isMobile = ref(false)

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setPage(page: string) {
    currentPage.value = page
  }

  function checkMobile() {
    isMobile.value = window.innerWidth < 768
    if (isMobile.value) sidebarCollapsed.value = true
  }

  return { sidebarCollapsed, currentPage, isMobile, toggleSidebar, setPage, checkMobile }
})
