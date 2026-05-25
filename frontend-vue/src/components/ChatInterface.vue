<template>
  <div class="chat-interface">
    <!-- 顶部控制栏 -->
    <div class="chat-header">
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
    </div>

    <!-- 聊天消息区域 -->
    <div class="chat-messages">
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
    </div>

    <!-- 输入区域 -->
    <div class="chat-input">
      <n-space vertical style="width: 100%">
        <n-input
          v-model:value="inputMessage"
          type="textarea"
          placeholder="输入消息..."
          :autosize="{ minRows: 2, maxRows: 4 }"
          @keydown.enter.exact.prevent="handleSend"
          @keydown.enter.shift.exact="inputMessage += '\n'"
        />
        <n-space justify="end">
          <n-button
            type="primary"
            :loading="isGenerating"
            :disabled="!inputMessage.trim() || !selectedModel"
            @click="handleSend"
          >
            发送 ➤
          </n-button>
        </n-space>
      </n-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useMessage } from 'naive-ui'
import {
  NSpace,
  NSelect,
  NTag,
  NStatistic,
  NText,
  NScrollbar,
  NAvatar,
  NInput,
  NButton
} from 'naive-ui'
import { modelAPI } from '@/api/model'

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: string
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

const selectedModel = ref<string>()
const modelOptions = ref<any[]>([])
const currentModel = ref<any>(null)
const messages = ref<Message[]>([])
const inputMessage = ref('')
const isGenerating = ref(false)
const streamingContent = ref('')

const stats = ref<Stats>({
  tokensPerSecond: 0,
  contextUsed: 0,
  contextTotal: 4096,
  outputTokens: 0,
  speed: 0
})

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

const handleModelChange = (value: string) => {
  const option = modelOptions.value.find(o => o.value === value)
  if (option) {
    currentModel.value = option.model
    stats.value.contextTotal = option.model.contextLength * 1024
    message.success(`已切换到模型: ${option.model.name}`)
  }
}

const handleSend = async () => {
  if (!inputMessage.value.trim() || !selectedModel.value) return

  const userMessage: Message = {
    role: 'user',
    content: inputMessage.value,
    timestamp: new Date().toLocaleTimeString()
  }

  messages.value.push(userMessage)
  const userInput = inputMessage.value
  inputMessage.value = ''
  
  await nextTick()
  scrollToBottom()

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
.chat-interface {
  display: flex;
  flex-direction: column;
  height: 600px;
  background: #f5f5f5;
  border-radius: 8px;
  overflow: hidden;
}

.chat-header {
  padding: 16px;
  background: white;
  border-bottom: 1px solid #e8e8e8;
}

.chat-messages {
  flex: 1;
  overflow: hidden;
  background: #f5f5f5;
}

.messages-container {
  padding: 16px;
}

.message {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
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
  margin-bottom: 4px;
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

.chat-input {
  padding: 16px;
  background: white;
  border-top: 1px solid #e8e8e8;
}
</style>
