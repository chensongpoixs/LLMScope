<template>
  <div style="padding: 20px;">
    <h1>🔬 LLMScope 测试</h1>
    <p>如果你能看到这个，说明 Vue 正常工作</p>
    <button @click="count++">点击次数: {{ count }}</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const count = ref(0)
</script>
