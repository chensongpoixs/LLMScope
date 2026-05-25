import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface ActivationData {
  layer: number
  mean: number
  max: number
  hiddenStates?: number[][]
}

export const useActivationStore = defineStore('activation', () => {
  const activations = ref<ActivationData[]>([])
  const currentToken = ref('')
  const tokenPosition = ref(0)

  const setActivations = (data: ActivationData[]) => {
    activations.value = data
  }

  const updateToken = (token: string, position: number) => {
    currentToken.value = token
    tokenPosition.value = position
  }

  return {
    activations,
    currentToken,
    tokenPosition,
    setActivations,
    updateToken
  }
})
