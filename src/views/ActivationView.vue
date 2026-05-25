<template>
  <div class="activation-view">
    <div class="prompt-panel">
      <n-card title="Prompt 输入">
        <n-space vertical>
          <n-input
            v-model:value="prompt"
            type="textarea"
            placeholder="输入 prompt..."
            :rows="3"
          />
          <n-button type="primary" @click="runInference">运行推理</n-button>
        </n-space>
      </n-card>
    </div>
    
    <div class="visualization-panel">
      <div ref="activationCanvasRef" class="activation-canvas"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NCard, NSpace, NInput, NButton } from 'naive-ui'

const prompt = ref('')
const activationCanvasRef = ref<HTMLElement>()

const runInference = () => {
  // WebSocket 连接和激活可视化将在这里实现
}
</script>

<style scoped>
.activation-view {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.prompt-panel {
  padding: 16px;
  background: white;
  border-bottom: 1px solid #e8e8e8;
}

.visualization-panel {
  flex: 1;
  overflow: hidden;
}

.activation-canvas {
  width: 100%;
  height: 100%;
  background: #fafafa;
}
</style>
