<template>
  <div class="home-view">
    <n-layout>
      <n-layout-header class="header">
        <h1>🔬 LLMScope</h1>
        <p>大模型显微镜</p>
      </n-layout-header>
      
      <n-layout-content class="content">
        <n-space vertical :size="24">
          <n-card title="模型实验室">
            <n-button type="primary" @click="showImportModal = true">
              导入模型
            </n-button>
          </n-card>

          <n-grid :cols="3" :x-gap="16" :y-gap="16">
            <n-grid-item v-for="model in models" :key="model.id">
              <n-card :title="model.name" hoverable @click="openModel(model.id)">
                <n-space vertical>
                  <n-text>参数量: {{ model.params }}</n-text>
                  <n-text>架构: {{ model.architecture }}</n-text>
                  <n-text>量化: {{ model.quantization }}</n-text>
                  <n-text>Context: {{ model.contextLength }}K</n-text>
                </n-space>
              </n-card>
            </n-grid-item>
          </n-grid>
        </n-space>
      </n-layout-content>
    </n-layout>

    <n-modal v-model:show="showImportModal" preset="card" title="导入模型">
      <n-form>
        <n-form-item label="模型名称">
          <n-input v-model:value="newModel.name" placeholder="LLaMA-2-7B" />
        </n-form-item>
        <n-form-item label="GGUF 文件路径">
          <n-input v-model:value="newModel.path" placeholder="/path/to/model.gguf" />
        </n-form-item>
        <n-button type="primary" @click="importModel">导入</n-button>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NLayout, NLayoutHeader, NLayoutContent, NCard, NButton, NSpace, NGrid, NGridItem, NText, NModal, NForm, NFormItem, NInput } from 'naive-ui'
import { useModelStore } from '@/stores/modelStore'

const router = useRouter()
const modelStore = useModelStore()
const models = ref(modelStore.models)
const showImportModal = ref(false)
const newModel = ref({ name: '', path: '' })

const openModel = (modelId: string) => {
  router.push(`/dashboard/${modelId}`)
}

const importModel = () => {
  modelStore.addModel({
    id: Date.now().toString(),
    name: newModel.value.name,
    params: '7B',
    architecture: 'LLaMA',
    contextLength: 4,
    quantization: 'Q4_K_M',
    fileSize: '3.8GB',
    loaded: false
  })
  showImportModal.value = false
}
</script>

<style scoped>
.home-view {
  width: 100%;
  height: 100vh;
}

.header {
  padding: 24px;
  background: white;
  border-bottom: 1px solid #e8e8e8;
}

.header h1 {
  margin: 0;
  font-size: 32px;
}

.header p {
  margin: 4px 0 0;
  color: #666;
}

.content {
  padding: 24px;
}
</style>
