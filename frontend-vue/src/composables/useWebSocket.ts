import { ref, onUnmounted } from 'vue'
import { config } from '@/config'

export function useWebSocket(url?: string) {
  const ws = ref<WebSocket | null>(null)
  const connected = ref(false)
  const messages = ref<any[]>([])

  const connect = () => {
    // 如果没有传入 url，使用配置文件中的默认 WebSocket 地址
    const wsUrl = url || config.wsUrl
    ws.value = new WebSocket(wsUrl)

    ws.value.onopen = () => {
      connected.value = true
    }

    ws.value.onmessage = (event) => {
      const data = JSON.parse(event.data)
      messages.value.push(data)
    }

    ws.value.onclose = () => {
      connected.value = false
    }

    ws.value.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  const send = (data: any) => {
    if (ws.value && connected.value) {
      ws.value.send(JSON.stringify(data))
    }
  }

  const disconnect = () => {
    if (ws.value) {
      ws.value.close()
    }
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    messages,
    connect,
    send,
    disconnect
  }
}
