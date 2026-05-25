/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// 全局配置类型
interface Window {
  LLMSCOPE_CONFIG?: {
    apiBaseUrl: string
    wsUrl: string
    timeout: number
    maxRetries: number
  }
}
