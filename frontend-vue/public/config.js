// LLMScope 后端配置
// 此文件在构建后可以直接修改，无需重新编译前端代码
window.LLMSCOPE_CONFIG = {
  // llama.cpp 后端地址
  apiBaseUrl: 'http://127.0.0.1:8080',
  
  // WebSocket 地址
  wsUrl: 'ws://127.0.0.1:8080/ws',
  
  // 其他配置
  timeout: 30000, // 请求超时时间（毫秒）
  maxRetries: 3,  // 最大重试次数
}
