<template>
  <n-config-provider>
    <n-message-provider>
      <n-loading-bar-provider>
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </n-loading-bar-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NLoadingBarProvider } from 'naive-ui'

// 添加错误处理
window.addEventListener('error', (e) => {
  console.error('Global error:', e.error)
})

window.addEventListener('unhandledrejection', (e) => {
  console.error('Unhandled rejection:', e.reason)
})
</script>

<style>
/* 确保应用占满整个视口 */
#app {
  width: 100%;
  height: 100%;
}
</style>
