import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUIStore = defineStore('ui', () => {
  const isDarkMode = ref(false)
  const sidebarCollapsed = ref(false)

  const toggleDarkMode = () => {
    isDarkMode.value = !isDarkMode.value
  }

  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  return {
    isDarkMode,
    sidebarCollapsed,
    toggleDarkMode,
    toggleSidebar
  }
})
