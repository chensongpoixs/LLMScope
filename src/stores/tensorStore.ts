import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface TensorInfo {
  name: string
  shape: number[]
  dtype: string
  params: number
  device: string
  quantization?: string
  min?: number
  max?: number
  mean?: number
  std?: number
}

export const useTensorStore = defineStore('tensor', () => {
  const tensors = ref<TensorInfo[]>([])
  const selectedTensor = ref<TensorInfo | null>(null)

  const setTensors = (data: TensorInfo[]) => {
    tensors.value = data
  }

  const selectTensor = (name: string) => {
    selectedTensor.value = tensors.value.find(t => t.name === name) || null
  }

  return {
    tensors,
    selectedTensor,
    setTensors,
    selectTensor
  }
})
