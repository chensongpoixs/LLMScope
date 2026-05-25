import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import './style.css'

// 全局错误处理
window.addEventListener('error', (e) => {
  console.error('Global error:', e.error)
})

window.addEventListener('unhandledrejection', (e) => {
  console.error('Unhandled rejection:', e.reason)
})

// 创建 Vue 应用
const app = createApp(App)
const pinia = createPinia()

// 注册插件
app.use(pinia)
app.use(router)

// 全局错误处理器
app.config.errorHandler = (err, instance, info) => {
  console.error('Vue error:', err)
  console.error('Component:', instance)
  console.error('Info:', info)
}

// 挂载应用
app.mount('#app')

// 调试信息
console.log('✓ LLMScope Frontend initialized')
console.log('✓ Config:', window.LLMSCOPE_CONFIG)
console.log('✓ Router ready')
console.log('✓ Pinia ready')
