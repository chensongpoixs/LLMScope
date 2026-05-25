<template>
  <div class="microscope-view">
    <div class="toolbar">
      <n-space>
        <n-button-group>
          <n-button @click="setLevel(1)">Layer</n-button>
          <n-button @click="setLevel(2)">Module</n-button>
          <n-button @click="setLevel(3)">Tensor</n-button>
          <n-button @click="setLevel(4)">Element</n-button>
        </n-button-group>
        <n-button @click="resetView">重置视图</n-button>
      </n-space>
    </div>
    
    <div ref="canvasRef" class="canvas"></div>
    
    <div class="minimap"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NSpace, NButton, NButtonGroup } from 'naive-ui'
import { useViewportStore } from '@/stores/viewportStore'

const viewportStore = useViewportStore()
const canvasRef = ref<HTMLElement>()

const setLevel = (level: number) => {
  viewportStore.setFocusLevel(level)
}

const resetView = () => {
  viewportStore.reset()
}

onMounted(() => {
  // PixiJS 初始化将在这里实现
})
</script>

<style scoped>
.microscope-view {
  width: 100%;
  height: 100vh;
  position: relative;
}

.toolbar {
  position: absolute;
  top: 16px;
  left: 16px;
  z-index: 10;
}

.canvas {
  width: 100%;
  height: 100%;
  background: #fafafa;
}

.minimap {
  position: absolute;
  bottom: 16px;
  right: 16px;
  width: 200px;
  height: 150px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>
