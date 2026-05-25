import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface ModelInfo {
  id: string
  name: string
  params: string
  architecture: string
  contextLength: number
  quantization: string
  fileSize: string
  loaded: boolean
}

export const useModelStore = defineStore('model', () => {
  const models = ref<ModelInfo[]>([])
  const currentModel = ref<ModelInfo | null>(null)
  const loadingProgress = ref(0)

  const addModel = (model: ModelInfo) => {
    models.value.push(model)
  }

  const loadModel = (modelId: string) => {
    const model = models.value.find(m => m.id === modelId)
    if (model) {
      currentModel.value = model
      model.loaded = true
    }
  }

  return {
    models,
    currentModel,
    loadingProgress,
    addModel,
    loadModel
  }
})
