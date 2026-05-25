import type { ModelInfo } from '@/stores/modelStore'
import { config } from '@/config'

const API_BASE = config.apiBaseUrl

export const modelAPI = {
  async getModels(): Promise<ModelInfo[]> {
    const response = await fetch(`${API_BASE}/api/models`, {
      signal: AbortSignal.timeout(config.timeout)
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    return response.json()
  },

  async loadModel(modelId: string) {
    const response = await fetch(`${API_BASE}/api/models/${modelId}/load`, {
      method: 'POST',
      signal: AbortSignal.timeout(config.timeout)
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    return response.json()
  },

  async getModelStructure(modelId: string) {
    const response = await fetch(`${API_BASE}/api/models/${modelId}/structure`, {
      signal: AbortSignal.timeout(config.timeout)
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    return response.json()
  },

  async getTensorInfo(modelId: string, tensorName: string) {
    const response = await fetch(`${API_BASE}/api/models/${modelId}/tensors/${tensorName}`, {
      signal: AbortSignal.timeout(config.timeout)
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    return response.json()
  }
}

export const llamaAPI = {
  async getStatus() {
    const response = await fetch(`${API_BASE}/api/llama/status`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    return response.json()
  },

  async getModels() {
    const response = await fetch(`${API_BASE}/api/llama/models`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    return response.json()
  }
}
