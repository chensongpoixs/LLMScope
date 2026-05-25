// 配置类型定义
export interface AppConfig {
  apiBaseUrl: string
  wsUrl: string
  timeout: number
  maxRetries: number
}

// 从 window 对象获取配置
declare global {
  interface Window {
    LLMSCOPE_CONFIG?: AppConfig
  }
}

// 默认配置（开发环境使用）
const defaultConfig: AppConfig = {
  apiBaseUrl: 'http://localhost:8080',
  wsUrl: 'ws://localhost:8080/ws',
  timeout: 30000,
  maxRetries: 3
}

// 获取配置
export function getConfig(): AppConfig {
  // 优先使用 window.LLMSCOPE_CONFIG（生产环境）
  // 如果不存在则使用默认配置（开发环境）
  return window.LLMSCOPE_CONFIG || defaultConfig
}

// 导出配置实例
export const config = getConfig()
