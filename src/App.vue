<template>
  <n-config-provider :theme="theme">
    <n-message-provider>
      <router-view />
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NMessageProvider } from 'naive-ui'
import { useUIStore } from '@/stores/uiStore'

const uiStore = useUIStore()
const theme = computed(() => uiStore.isDarkMode ? null : null)
</script>
