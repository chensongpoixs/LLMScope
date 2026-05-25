// LLMScope 后端配置示例
// 复制此文件为 config.js 并修改配置
window.LLMSCOPE_CONFIG = {
  // llama.cpp 后端地址
  apiBaseUrl: 'http://localhost:8080',
  
  // WebSocket 地址
  wsUrl: 'ws://localhost:8080/ws',
  
  // 其他配置
  timeout: 30000, // 请求超时时间（毫秒）
  maxRetries: 3,  // 最大重试次数
}
