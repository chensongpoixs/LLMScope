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
            <n-space vertical :size="12">
              <n-space>
                <n-button 
                  :type="showChat ? 'default' : 'primary'"
                  @click="showChat = !showChat"
                >
                  {{ showChat ? '📋 查看模型列表' : '💬 开始聊天' }}
                </n-button>
                <n-button type="primary" @click="showImportModal = true">
                  导入模型
                </n-button>
                <n-button @click="loadModels" :loading="loading">
                  刷新模型列表
                </n-button>
                <n-button @click="checkLlamaStatus">
                  检查 llama.cpp 状态
                </n-button>
              </n-space>
              
              <!-- llama.cpp 状态显示 -->
              <n-alert v-if="llamaStatus" :type="llamaStatus.connected ? 'success' : 'warning'">
                <template #header>
                  llama.cpp 服务状态
                </template>
                <n-space vertical :size="4">
                  <n-text>连接状态: {{ llamaStatus.connected ? '✓ 已连接' : '✗ 未连接' }}</n-text>
                  <n-text>服务地址: {{ llamaStatus.url }}</n-text>
                  <n-text v-if="llamaStatus.error">错误: {{ llamaStatus.error }}</n-text>
                </n-space>
              </n-alert>
            </n-space>
          </n-card>

          <!-- 聊天界面 -->
          <n-card v-if="showChat" class="chat-card">
            <ChatInterface />
          </n-card>

          <!-- 模型列表 -->
          <n-spin v-else :show="loading">
            <n-card v-if="models.length === 0 && !loading" title="开始使用">
              <n-empty description="暂无模型，请导入模型或确保 llama.cpp 服务正在运行">
                <template #extra>
                  <n-space>
                    <n-button type="primary" @click="showImportModal = true">
                      导入模型
                    </n-button>
                    <n-button @click="loadModels">
                      重新加载
                    </n-button>
                  </n-space>
                </template>
              </n-empty>
            </n-card>

            <n-grid v-else :cols="3" :x-gap="16" :y-gap="16">
              <n-grid-item v-for="model in models" :key="model.id">
                <n-card :title="model.name" hoverable @click="openModel(model.id)">
                  <n-space vertical>
                    <n-text>参数量: {{ model.params }}</n-text>
                    <n-text>架构: {{ model.architecture }}</n-text>
                    <n-text>量化: {{ model.quantization }}</n-text>
                    <n-text>Context: {{ formatContext(model.contextLength) }}</n-text>
                    <n-text v-if="model.fileSize && model.fileSize !== 'Unknown'">
                      体积: {{ model.fileSize }}
                    </n-text>
                    <n-tag :type="model.loaded ? 'success' : 'default'">
                      {{ model.loaded ? '已加载' : '未加载' }}
                    </n-tag>
                  </n-space>
                </n-card>
              </n-grid-item>
            </n-grid>
          </n-spin>
        </n-space>
      </n-layout-content>
    </n-layout>

    <n-modal v-model:show="showImportModal" preset="card" title="导入模型" style="width: 600px;">
      <n-form>
        <n-form-item label="模型名称">
          <n-input v-model:value="newModel.name" placeholder="LLaMA-2-7B" />
        </n-form-item>
        <n-form-item label="GGUF 文件路径">
          <n-input v-model:value="newModel.path" placeholder="/path/to/model.gguf" />
        </n-form-item>
        <n-space>
          <n-button type="primary" @click="importModel">导入</n-button>
          <n-button @click="showImportModal = false">取消</n-button>
        </n-space>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { 
  NLayout, 
  NLayoutHeader, 
  NLayoutContent, 
  NCard, 
  NButton, 
  NSpace, 
  NGrid, 
  NGridItem, 
  NText, 
  NModal, 
  NForm, 
  NFormItem, 
  NInput,
  NEmpty,
  NSpin,
  NTag,
  NAlert
} from 'naive-ui'
import { useModelStore, type ModelInfo } from '@/stores/modelStore'
import { modelAPI, llamaAPI } from '@/api/model'
import { config } from '@/config'
import ChatInterface from '@/components/ChatInterface.vue'

const router = useRouter()
const message = useMessage()
const modelStore = useModelStore()
const models = ref<ModelInfo[]>([])
const loading = ref(false)
const showImportModal = ref(false)
const newModel = ref({ name: '', path: '' })
const llamaStatus = ref<any>(null)
const showChat = ref(false)

/** contextLength 为 K 单位（来自后端 n_ctx/1024） */
function formatContext(contextK: number): string {
  if (!contextK || contextK <= 0) return '未知'
  if (contextK >= 1024) return `${(contextK / 1024).toFixed(1)}M`
  return `${contextK}K`
}

// 检查 llama.cpp 状态
const checkLlamaStatus = async () => {
  try {
    const status = await llamaAPI.getStatus()
    llamaStatus.value = status
    
    if (status.connected) {
      message.success('llama.cpp 服务连接成功')
    } else {
      message.warning('llama.cpp 服务未连接')
    }
  } catch (error) {
    message.error('检查 llama.cpp 状态失败: ' + (error as Error).message)
  }
}

// 从后端加载模型列表
const loadModels = async () => {
  loading.value = true
  try {
    const data = await modelAPI.getModels()
    models.value = data
    
    // 同步到 store
    modelStore.models = data
    
    message.success(`成功加载 ${data.length} 个模型`)
    console.log('Models loaded:', data)
  } catch (error) {
    message.error('加载模型失败: ' + (error as Error).message)
    console.error('Failed to load models:', error)
  } finally {
    loading.value = false
  }
}

const openModel = (modelId: string) => {
  router.push(`/dashboard/${modelId}`)
}

const importModel = () => {
  if (!newModel.value.name) {
    message.warning('请输入模型名称')
    return
  }
  
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
  
  models.value = modelStore.models
  message.success('模型导入成功')
  showImportModal.value = false
  newModel.value = { name: '', path: '' }
}

// 组件挂载时自动加载模型和检查状态
onMounted(() => {
  console.log('HomeView mounted')
  console.log('Config:', config)
  loadModels()
  checkLlamaStatus()
})
</script>

<style scoped>
.home-view {
  width: 100%;
  min-height: 100vh;
  background: #f5f5f5;
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

.chat-card {
  padding: 0;
}

.chat-card :deep(.n-card__content) {
  padding: 0;
}
</style>
