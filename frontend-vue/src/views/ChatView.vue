<template>
  <div class="chat-page">
    <ChatInterface full-height enable-attachments />
  </div>
</template>

<script setup lang="ts">
import ChatInterface from '@/components/ChatInterface.vue'
</script>

<style scoped>
.chat-page {
  width: 100%;
  min-height: 100vh;
  background: #fff;
}
</style>
