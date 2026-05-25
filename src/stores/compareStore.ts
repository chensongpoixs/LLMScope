import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useCompareStore = defineStore('compare', () => {
  const modelA = ref<string | null>(null)
  const modelB = ref<string | null>(null)
  const syncViewport = ref(true)

  const setModels = (a: string, b: string) => {
    modelA.value = a
    modelB.value = b
  }

  const toggleSync = () => {
    syncViewport.value = !syncViewport.value
  }

  return {
    modelA,
    modelB,
    syncViewport,
    setModels,
    toggleSync
  }
})
