<template>
  <div class="chat-view">
    <n-layout has-sider>
      <!-- 左侧边栏 - 附件菜单 -->
      <n-layout-sider
        v-if="showAttachmentMenu"
        width="240"
        bordered
        class="attachment-sidebar"
      >
        <n-menu :options="attachmentMenuOptions" @update:value="handleAttachmentSelect" />
      </n-layout-sider>

      <!-- 主聊天区域 -->
      <n-layout>
        <n-layout-header bordered class="chat-header">
          <n-space justify="space-between" align="center">
            <n-space align="center">
              <n-select
                v-model:value="selectedModel"
                :options="modelOptions"
                placeholder="选择模型"
                style="width: 300px"
                @update:value="handleModelChange"
              />
              <n-tag v-if="currentModel" type="success">
                {{ currentModel.params }}
              </n-tag>
              <n-tag v-if="currentModel">
                {{ currentModel.quantization }}
              </n-tag>
            </n-space>
            
            <!-- 实时统计 -->
            <n-space align="center">
              <n-statistic label="Token/s" :value="stats.tokensPerSecond.toFixed(1)" />
              <n-statistic label="Context" :value="`${stats.contextUsed}/${stats.contextTotal}`" />
              <n-statistic label="Output" :value="`${stats.outputTokens}/∞`" />
              <n-text>{{ stats.speed.toFixed(1) }} t/s</n-text>
            </n-space>
          </n-space>
        </n-layout-header>

        <n-layout-content class="chat-content">
          <n-scrollbar ref="scrollbarRef" style="height: 100%">
            <div class="messages-container">
              <div
                v-for="(message, index) in messages"
                :key="index"
                :class="['message', message.role]"
              >
                <div class="message-avatar">
                  <n-avatar v-if="message.role === 'user'" round>
                    👤
                  </n-avatar>
                  <n-avatar v-else round color="#18a058">
                    🤖
                  </n-avatar>
                </div>
                
                <div class="message-content">
                  <div class="message-header">
                    <n-text strong>{{ message.role === 'user' ? '你' : '助手' }}</n-text>
                    <n-text depth="3" style="font-size: 12px">{{ message.timestamp }}</n-text>
                  </div>
                  
                  <div class="message-text">
                    <n-text>{{ message.content }}</n-text>
                  </div>

                  <!-- 附件显示 -->
                  <div v-if="message.attachments && message.attachments.length > 0" class="message-attachments">
                    <n-space>
                      <n-tag
                        v-for="(file, idx) in message.attachments"
                        :key="idx"
                        closable
                        @close="removeAttachment(index, idx)"
                      >
                        {{ file.name }}
                      </n-tag>
                    </n-space>
                  </div>
                </div>
              </div>

              <!-- 正在生成的消息 -->
              <div v-if="isGenerating" class="message assistant generating">
                <div class="message-avatar">
                  <n-avatar round color="#18a058">
                    🤖
                  </n-avatar>
                </div>
                <div class="message-content">
                  <div class="message-text">
                    <n-text>{{ streamingContent }}</n-text>
                    <span class="cursor">▊</span>
                  </div>
                </div>
              </div>
            </div>
          </n-scrollbar>
        </n-layout-content>

        <n-layout-footer bordered class="chat-footer">
          <div class="input-container">
            <n-button
              text
              @click="showAttachmentMenu = !showAttachmentMenu"
              class="attachment-button"
            >
              ➕
            </n-button>

            <n-input
              v-model:value="inputMessage"
              type="textarea"
              placeholder="输入消息..."
              :autosize="{ minRows: 1, maxRows: 5 }"
              @keydown.enter.exact.prevent="handleSend"
              @keydown.enter.shift.exact="inputMessage += '\n'"
              class="message-input"
            />

            <n-button
              type="primary"
              :loading="isGenerating"
              :disabled="!inputMessage.trim() || !selectedModel"
              @click="handleSend"
              class="send-button"
            >
              发送 ➤
            </n-button>
          </div>

          <!-- 已选择的附件预览 -->
          <div v-if="pendingAttachments.length > 0" class="pending-attachments">
            <n-space>
              <n-tag
                v-for="(file, index) in pendingAttachments"
                :key="index"
                closable
                @close="removePendingAttachment(index)"
              >
                {{ file.name }}
              </n-tag>
            </n-space>
          </div>
        </n-layout-footer>
      </n-layout>
    </n-layout>

    <!-- 文件上传对话框 -->
    <input
      ref="fileInputRef"
      type="file"
      multiple
      accept="image/*,audio/*,.pdf,.txt"
      style="display: none"
      @change="handleFileSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useMessage } from 'naive-ui'
import {
  NLayout,
  NLayoutSider,
  NLayoutHeader,
  NLayoutContent,
  NLayoutFooter,
  NMenu,
  NSelect,
  NSpace,
  NTag,
  NStatistic,
  NText,
  NScrollbar,
  NAvatar,
  NInput,
  NButton
} from 'naive-ui'
import { modelAPI } from '@/api/model'
import { config } from '@/config'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: string
  attachments?: File[]
}

interface Stats {
  tokensPerSecond: number
  contextUsed: number
  contextTotal: number
  outputTokens: number
  speed: number
}

const message = useMessage()
const scrollbarRef = ref()
const fileInputRef = ref<HTMLInputElement>()

// 状态
const selectedModel = ref<string>()
const modelOptions = ref<any[]>([])
const currentModel = ref<any>(null)
const messages = ref<Message[]>([])
const inputMessage = ref('')
const isGenerating = ref(false)
const streamingContent = ref('')
const showAttachmentMenu = ref(false)
const pendingAttachments = ref<File[]>([])

// 统计数据
const stats = ref<Stats>({
  tokensPerSecond: 0,
  contextUsed: 0,
  contextTotal: 4096,
  outputTokens: 0,
  speed: 0
})

// 附件菜单选项
const attachmentMenuOptions = [
  { label: '📷 Images', key: 'images' },
  { label: '🎵 Audio Files', key: 'audio' },
  { label: '📄 Text Files', key: 'text' },
  { label: '📕 PDF Files', key: 'pdf' },
  { label: '💬 System Message', key: 'system' },
  { label: '🔧 Tools', key: 'tools' },
  { label: '🔌 MCP Servers', key: 'mcp' }
]

// 加载模型列表
const loadModels = async () => {
  try {
    const models = await modelAPI.getModels()
    modelOptions.value = models.map(m => ({
      label: `${m.name} (${m.params})`,
      value: m.id,
      model: m
    }))
    
    if (models.length > 0 && !selectedModel.value) {
      selectedModel.value = models[0].id
      currentModel.value = models[0]
      stats.value.contextTotal = models[0].contextLength * 1024
    }
  } catch (error) {
    message.error('加载模型失败: ' + (error as Error).message)
  }
}

// 模型切换
const handleModelChange = (value: string) => {
  const option = modelOptions.value.find(o => o.value === value)
  if (option) {
    currentModel.value = option.model
    stats.value.contextTotal = option.model.contextLength * 1024
    message.success(`已切换到模型: ${option.model.name}`)
  }
}

// 附件选择
const handleAttachmentSelect = (key: string) => {
  if (key === 'system' || key === 'tools' || key === 'mcp') {
    message.info(`${key} 功能开发中...`)
    return
  }

  let accept = '*/*'
  switch (key) {
    case 'images':
      accept = 'image/*'
      break
    case 'audio':
      accept = 'audio/*'
      break
    case 'text':
      accept = '.txt,.md,.json,.xml'
      break
    case 'pdf':
      accept = '.pdf'
      break
  }
  
  if (fileInputRef.value) {
    fileInputRef.value.accept = accept
    fileInputRef.value.click()
  }
}

// 文件选择处理
const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    pendingAttachments.value.push(...Array.from(target.files))
    message.success(`已添加 ${target.files.length} 个文件`)
  }
  target.value = ''
}

// 移除待发送附件
const removePendingAttachment = (index: number) => {
  pendingAttachments.value.splice(index, 1)
}

// 移除消息附件
const removeAttachment = (messageIndex: number, attachmentIndex: number) => {
  if (messages.value[messageIndex].attachments) {
    messages.value[messageIndex].attachments!.splice(attachmentIndex, 1)
  }
}

// 发送消息
const handleSend = async () => {
  if (!inputMessage.value.trim() || !selectedModel.value) return

  const userMessage: Message = {
    role: 'user',
    content: inputMessage.value,
    timestamp: new Date().toLocaleTimeString(),
    attachments: pendingAttachments.value.length > 0 ? [...pendingAttachments.value] : undefined
  }

  messages.value.push(userMessage)
  const userInput = inputMessage.value
  inputMessage.value = ''
  pendingAttachments.value = []
  
  await nextTick()
  scrollToBottom()

  // 开始生成
  isGenerating.value = true
  streamingContent.value = ''
  
  try {
    await streamChat(userInput)
  } catch (error) {
    message.error('生成失败: ' + (error as Error).message)
    messages.value.push({
      role: 'assistant',
      content: '抱歉，生成失败了。请检查后端服务是否正常运行。',
      timestamp: new Date().toLocaleTimeString()
    })
  } finally {
    isGenerating.value = false
  }
}

// 流式聊天（模拟）
const streamChat = async (prompt: string) => {
  const startTime = Date.now()
  let tokenCount = 0

  const response = `这是一个模拟的回复。你说: "${prompt}"\n\n当前使用的模型是 ${currentModel.value?.name}。\n\n实际使用时，这里会连接到 llama.cpp 进行真实的推理。`
  
  for (let i = 0; i < response.length; i++) {
    streamingContent.value += response[i]
    tokenCount++
    
    const elapsed = (Date.now() - startTime) / 1000
    stats.value.tokensPerSecond = elapsed > 0 ? tokenCount / elapsed : 0
    stats.value.speed = stats.value.tokensPerSecond
    stats.value.outputTokens = tokenCount
    stats.value.contextUsed = Math.min(stats.value.contextUsed + 1, stats.value.contextTotal)

    await new Promise(resolve => setTimeout(resolve, 20))
    await nextTick()
    scrollToBottom()
  }

  messages.value.push({
    role: 'assistant',
    content: streamingContent.value,
    timestamp: new Date().toLocaleTimeString()
  })

  streamingContent.value = ''
}

// 滚动到底部
const scrollToBottom = () => {
  if (scrollbarRef.value) {
    nextTick(() => {
      const container = scrollbarRef.value.$el.querySelector('.n-scrollbar-container')
      if (container) {
        container.scrollTop = container.scrollHeight
      }
    })
  }
}

onMounted(() => {
  loadModels()
})
</script>

<style scoped>
.chat-view {
  width: 100%;
  height: 100vh;
  background: #f5f5f5;
}

.attachment-sidebar {
  background: white;
}

.chat-header {
  padding: 16px 24px;
  background: white;
}

.chat-content {
  height: calc(100vh - 140px);
  background: #f5f5f5;
}

.messages-container {
  padding: 24px;
  max-width: 900px;
  margin: 0 auto;
}

.message {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  flex: 1;
  max-width: 70%;
}

.message.user .message-content {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.message-header {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}

.message.user .message-header {
  flex-direction: row-reverse;
}

.message-text {
  padding: 12px 16px;
  border-radius: 12px;
  background: white;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  word-wrap: break-word;
  white-space: pre-wrap;
}

.message.user .message-text {
  background: #18a058;
  color: white;
}

.message.generating .message-text {
  background: #f0f0f0;
}

.cursor {
  animation: blink 1s infinite;
  margin-left: 2px;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

.message-attachments {
  margin-top: 8px;
}

.chat-footer {
  padding: 16px 24px;
  background: white;
}

.input-container {
  display: flex;
  gap: 12px;
  align-items: flex-end;
  max-width: 900px;
  margin: 0 auto;
}

.attachment-button {
  font-size: 24px;
  padding: 8px;
}

.message-input {
  flex: 1;
}

.send-button {
  height: 40px;
}

.pending-attachments {
  margin-top: 12px;
  max-width: 900px;
  margin-left: auto;
  margin-right: auto;
}
</style>
