import { defineStore } from 'pinia'
import { ref } from 'vue'
import { config } from '@/config'

export const useWebSocketStore = defineStore('websocket', () => {
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)

  const connect = (url?: string) => {
    // 如果没有传入 url，使用配置文件中的默认 WebSocket 地址
    const wsUrl = url || config.wsUrl
    ws.value = new WebSocket(wsUrl)
    
    ws.value.onopen = () => {
      connected.value = true
    }

    ws.value.onclose = () => {
      connected.value = false
    }

    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
  }

  return {
    ws,
    connected,
    connect,
    disconnect
  }
})
