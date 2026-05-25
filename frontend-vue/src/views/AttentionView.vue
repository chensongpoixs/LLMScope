<template>
  <div class="attention-view">
    <n-layout has-sider>
      <n-layout-sider width="300" bordered>
        <div class="controls">
          <n-space vertical>
            <n-select v-model:value="selectedLayer" :options="layerOptions" placeholder="选择层" />
            <n-select v-model:value="selectedHead" :options="headOptions" placeholder="选择注意力头" />
            <n-button type="primary" @click="loadAttention">加载</n-button>
          </n-space>
        </div>
      </n-layout-sider>
      
      <n-layout-content>
        <div ref="heatmapRef" class="heatmap"></div>
      </n-layout-content>
    </n-layout>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NLayout, NLayoutSider, NLayoutContent, NSpace, NSelect, NButton } from 'naive-ui'

const selectedLayer = ref<number>()
const selectedHead = ref<number>()
const heatmapRef = ref<HTMLElement>()

const layerOptions = Array.from({ length: 32 }, (_, i) => ({
  label: `Layer ${i}`,
  value: i
}))

const headOptions = Array.from({ length: 32 }, (_, i) => ({
  label: `Head ${i}`,
  value: i
}))

const loadAttention = () => {
  // WebGL Shader 热力图渲染将在这里实现
}
</script>

<style scoped>
.attention-view {
  width: 100%;
  height: 100vh;
}

.controls {
  padding: 16px;
}

.heatmap {
  width: 100%;
  height: 100%;
  background: #fafafa;
}
</style>
