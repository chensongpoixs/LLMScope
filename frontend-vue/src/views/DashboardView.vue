<template>
  <div class="dashboard-view">
    <n-layout has-sider>
      <n-layout-sider bordered>
        <n-menu :options="menuOptions" @update:value="handleMenuSelect" />
      </n-layout-sider>
      
      <n-layout-content class="content">
        <n-space vertical :size="24">
          <n-card title="模型参数统计">
            <n-grid :cols="3" :x-gap="16">
              <n-grid-item>
                <n-statistic label="总参数量" value="7B" />
              </n-grid-item>
              <n-grid-item>
                <n-statistic label="层数" value="32" />
              </n-grid-item>
              <n-grid-item>
                <n-statistic label="隐藏维度" value="4096" />
              </n-grid-item>
            </n-grid>
          </n-card>

          <n-card title="模型结构">
            <div class="structure-preview">
              <div class="layer">Embedding</div>
              <div class="arrow">↓</div>
              <div class="layer">Transformer × 32</div>
              <div class="arrow">↓</div>
              <div class="layer">LM Head</div>
            </div>
          </n-card>

          <n-card title="层级参数分布">
            <div ref="chartRef" style="width: 100%; height: 300px;"></div>
          </n-card>
        </n-space>
      </n-layout-content>
    </n-layout>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutContent, NMenu, NCard, NSpace, NGrid, NGridItem, NStatistic } from 'naive-ui'
import * as echarts from 'echarts'

const router = useRouter()
const route = useRoute()
const chartRef = ref<HTMLElement>()

const menuOptions = [
  { label: '总览', key: 'dashboard' },
  { label: '结构显微镜', key: 'microscope' },
  { label: 'Attention', key: 'attention' },
  { label: '激活显微镜', key: 'activation' }
]

const handleMenuSelect = (key: string) => {
  const modelId = route.params.modelId
  router.push(`/${key}/${modelId}`)
}

onMounted(() => {
  if (chartRef.value) {
    const chart = echarts.init(chartRef.value)
    chart.setOption({
      xAxis: { type: 'category', data: Array.from({ length: 32 }, (_, i) => `Layer ${i}`) },
      yAxis: { type: 'value', name: '参数量 (M)' },
      series: [{ data: Array(32).fill(220), type: 'bar' }]
    })
  }
})
</script>

<style scoped>
.dashboard-view {
  width: 100%;
  height: 100vh;
}

.content {
  padding: 24px;
}

.structure-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
}

.layer {
  padding: 16px 32px;
  background: #f0f0f0;
  border-radius: 8px;
  font-weight: 500;
}

.arrow {
  font-size: 24px;
  color: #666;
}
</style>
