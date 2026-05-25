import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useViewportStore = defineStore('viewport', () => {
  const zoom = ref(1)
  const panX = ref(0)
  const panY = ref(0)
  const focusLevel = ref(1) // 1: Layer, 2: Module, 3: Tensor, 4: Element

  const setZoom = (value: number) => {
    zoom.value = Math.max(0.1, Math.min(10, value))
  }

  const setPan = (x: number, y: number) => {
    panX.value = x
    panY.value = y
  }

  const setFocusLevel = (level: number) => {
    focusLevel.value = Math.max(1, Math.min(4, level))
  }

  const reset = () => {
    zoom.value = 1
    panX.value = 0
    panY.value = 0
  }

  return {
    zoom,
    panX,
    panY,
    focusLevel,
    setZoom,
    setPan,
    setFocusLevel,
    reset
  }
})
