import type { ModelInfo } from '@/stores/modelStore'
import { config } from '@/config'

const API_BASE = config.apiBaseUrl

export const modelAPI = {
  async getModels(): Promise<ModelInfo[]> {
    const response = await fetch(`${API_BASE}/models`, {
      signal: AbortSignal.timeout(config.timeout)
    })
    return response.json()
  },

  async loadModel(modelId: string) {
    const response = await fetch(`${API_BASE}/models/${modelId}/load`, {
      method: 'POST',
      signal: AbortSignal.timeout(config.timeout)
    })
    return response.json()
  },

  async getModelStructure(modelId: string) {
    const response = await fetch(`${API_BASE}/models/${modelId}/structure`, {
      signal: AbortSignal.timeout(config.timeout)
    })
    return response.json()
  },

  async getTensorInfo(modelId: string, tensorName: string) {
    const response = await fetch(`${API_BASE}/models/${modelId}/tensors/${tensorName}`, {
      signal: AbortSignal.timeout(config.timeout)
    })
    return response.json()
  }
}
