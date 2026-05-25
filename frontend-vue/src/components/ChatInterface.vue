<template>
  <div class="chat-shell" :class="{ 'chat-shell--full': fullHeight }">
    <!-- 顶部状态栏（有对话时显示） -->
    <div v-if="showTopBar" class="chat-top-bar">
      <div class="top-bar-left">
        <div class="model-pill model-pill--header">
          <span class="model-pill__text">{{ displayModelName }}</span>
          <button type="button" class="icon-btn" title="模型信息" @click="showModelInfo">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8">
              <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
              <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" />
            </svg>
          </button>
          <button type="button" class="icon-btn" title="新对话" @click="handleNewChat">
            <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8">
              <path d="M12 3l1.5 4.5L18 9l-4.5 1.5L12 15l-1.5-4.5L6 9l4.5-1.5L12 3z" />
              <path d="M5 19h14" />
            </svg>
          </button>
        </div>
      </div>
      <div class="top-bar-metrics">
        <span class="metric">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8">
            <path d="M4 6h16M4 12h10M4 18h6" />
          </svg>
          {{ totalTokensDisplay }} tokens
        </span>
        <span class="metric">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8">
            <circle cx="12" cy="12" r="9" />
            <path d="M12 7v5l3 2" />
          </svg>
          {{ formattedDuration }}
        </span>
        <span class="metric">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8">
            <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83" />
          </svg>
          {{ stats.speed.toFixed(2) }} t/s
        </span>
      </div>
    </div>

    <!-- 消息区 -->
    <div class="chat-body">
      <n-scrollbar ref="scrollbarRef" class="chat-scroll">
        <div class="messages-wrap">
          <div v-if="messages.length === 0 && !isGenerating" class="chat-empty">
            <p class="chat-empty__title">开始对话</p>
            <p class="chat-empty__hint">在下方输入消息，与当前模型聊天</p>
          </div>

          <article
            v-for="(msg, index) in messages"
            :key="index"
            :class="['msg-block', msg.role]"
          >
            <template v-if="msg.role === 'user'">
              <div class="msg-block__content msg-block__content--user">{{ msg.content }}</div>
            </template>
            <template v-else>
              <ThinkingSection
                v-if="msg.reasoning"
                :thinking="msg.reasoning"
              />
              <div v-if="msg.content" class="msg-block__content msg-block__content--answer">
                {{ msg.content }}
              </div>
            </template>
          </article>

          <article v-if="isGenerating" class="msg-block assistant generating">
            <ThinkingSection
              v-if="streamingReasoning"
              :thinking="streamingReasoning"
              is-streaming
            />
            <div v-if="streamingContent" class="msg-block__content msg-block__content--answer">
              {{ streamingContent }}<span class="cursor">▊</span>
            </div>
            <div
              v-else-if="!streamingReasoning && !streamingContent"
              class="msg-block__content msg-block__content--pending"
            >
              等待模型响应…
            </div>
          </article>

          <!-- 消息操作栏 -->
          <div
            v-if="lastAssistantIndex >= 0 && !isGenerating"
            class="message-toolbar"
          >
            <button
              type="button"
              class="toolbar-btn"
              title="复制"
              @click="copyAssistantMessage(messages[lastAssistantIndex])"
            >
              <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6">
                <rect x="9" y="9" width="13" height="13" rx="2" />
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              </svg>
            </button>
            <button type="button" class="toolbar-btn" title="编辑" @click="notifySoon('编辑')">
              <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6">
                <path d="M12 20h9" />
                <path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z" />
              </svg>
            </button>
            <button type="button" class="toolbar-btn" title="重新生成" @click="handleRegenerate">
              <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6">
                <path d="M3 12a9 9 0 0 1 15-6.7L21 3v6h-6" />
                <path d="M21 12a9 9 0 0 1-15 6.7L3 21v-6h6" />
              </svg>
            </button>
            <button type="button" class="toolbar-btn" title="分支" @click="notifySoon('分支')">
              <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6">
                <circle cx="6" cy="6" r="2" />
                <circle cx="18" cy="18" r="2" />
                <path d="M8 6h5a4 4 0 0 1 4 4v2M16 18h-5a4 4 0 0 1-4-4v-2" />
              </svg>
            </button>
            <button type="button" class="toolbar-btn" title="删除" @click="deleteLastAssistant">
              <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6">
                <path d="M3 6h18M8 6V4h8v2M19 6l-1 14H6L5 6" />
              </svg>
            </button>
          </div>
        </div>
      </n-scrollbar>
    </div>

    <!-- 底部输入区 -->
    <div class="chat-composer-wrap">
      <!-- 输入框上方：实时统计（llama.cpp 风格） -->
      <div v-if="showLiveStats" class="live-stats-bar">
        <span class="live-stat">
          Context: {{ stats.contextUsed.toLocaleString() }}/{{ stats.contextTotal.toLocaleString() }}
          ({{ contextPercent }}%)
        </span>
        <span class="live-stat">
          Output: {{ liveOutputDisplay }}/∞
        </span>
        <span class="live-stat">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8">
            <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4" />
          </svg>
          {{ liveSpeedDisplay }} t/s
        </span>
        <span class="live-stat">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8">
            <path d="M4 6h16M4 12h10M4 18h6" />
          </svg>
          {{ liveTokensDisplay }} tokens
        </span>
        <span class="live-stat">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8">
            <circle cx="12" cy="12" r="9" />
            <path d="M12 7v5l3 2" />
          </svg>
          {{ liveDurationDisplay }}
        </span>
      </div>

      <div v-if="pendingAttachments.length > 0" class="pending-files">
        <span
          v-for="(file, i) in pendingAttachments"
          :key="i"
          class="file-chip"
        >
          {{ file.name }}
          <button type="button" @click="pendingAttachments.splice(i, 1)">×</button>
        </span>
      </div>

      <div class="chat-composer">
        <n-popover
          v-if="enableAttachments"
          trigger="click"
          placement="top-start"
          :show="showAttachMenu"
          @update:show="showAttachMenu = $event"
        >
          <template #trigger>
            <button type="button" class="btn-circle btn-add" title="添加附件">+</button>
          </template>
          <div class="attach-menu">
            <button v-for="item in attachOptions" :key="item.key" type="button" @click="pickAttachment(item.key)">
              {{ item.label }}
            </button>
          </div>
        </n-popover>
        <button v-else type="button" class="btn-circle btn-add" title="添加" disabled>+</button>

        <textarea
          ref="textareaRef"
          v-model="inputMessage"
          class="composer-input"
          placeholder="Type a message..."
          rows="1"
          @keydown.enter.exact.prevent="handleSend"
          @keydown.enter.shift.exact.stop
          @input="autoResize"
        />

        <n-tooltip trigger="hover" placement="top">
          <template #trigger>
            <button
              type="button"
              class="model-pill model-pill--inline"
              @click="showModelPicker = true"
            >
              <span class="model-pill__text">{{ displayModelName }}</span>
            </button>
          </template>
          <span class="model-tooltip">{{ currentModel?.name || selectedModel }}</span>
        </n-tooltip>

        <button
          type="button"
          class="btn-circle btn-send"
          :class="{ active: canSend }"
          :disabled="!canSend || isGenerating"
          title="发送"
          @click="handleSend"
        >
          <svg v-if="!isGenerating" viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 19V5M6 11l6-6 6 6" />
          </svg>
          <span v-else class="send-loading" />
        </button>
      </div>
    </div>

    <!-- 模型选择 -->
    <n-modal v-model:show="showModelPicker" preset="card" title="选择模型" style="max-width: 420px">
      <n-select
        v-model:value="selectedModel"
        :options="modelOptions"
        placeholder="选择模型"
        @update:value="handleModelChange"
      />
    </n-modal>

    <input
      v-if="enableAttachments"
      ref="fileInputRef"
      type="file"
      multiple
      class="hidden-input"
      @change="handleFileSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useMessage, NScrollbar, NPopover, NTooltip, NModal, NSelect } from 'naive-ui'
import { modelAPI } from '@/api/model'
import { useChatStream, type ChatUiMessage, type ChatStats } from '@/composables/useChatStream'
import ThinkingSection from '@/components/ThinkingSection.vue'

withDefaults(
  defineProps<{
    fullHeight?: boolean
    enableAttachments?: boolean
  }>(),
  {
    fullHeight: false,
    enableAttachments: false
  }
)

type Message = ChatUiMessage & { attachments?: File[] }
type Stats = ChatStats

const message = useMessage()
const scrollbarRef = ref()
const textareaRef = ref<HTMLTextAreaElement>()
const fileInputRef = ref<HTMLInputElement>()

const selectedModel = ref<string>()
const modelOptions = ref<{ label: string; value: string; model: unknown }[]>([])
const currentModel = ref<{ id: string; name: string; contextLength: number } | null>(null)
const messages = ref<Message[]>([])
const inputMessage = ref('')
const streamingContent = ref('')
const streamingReasoning = ref('')
const showModelPicker = ref(false)
const showAttachMenu = ref(false)
const pendingAttachments = ref<File[]>([])
const sessionStartTime = ref<number | null>(null)
const elapsedSeconds = ref(0)
const generationStartTime = ref<number | null>(null)
const generationElapsedSeconds = ref(0)
const lastTurnStats = ref({ tokens: 0, seconds: 0, speed: 0 })
let elapsedTimer: ReturnType<typeof setInterval> | null = null
let generationTimer: ReturnType<typeof setInterval> | null = null

const stats = ref<Stats>({
  tokensPerSecond: 0,
  contextUsed: 0,
  contextTotal: 4096,
  outputTokens: 0,
  speed: 0
})

const { isGenerating, streamChat } = useChatStream(
  messages,
  selectedModel,
  streamingContent,
  streamingReasoning,
  stats,
  () => {
    tickElapsed()
    tickGenerationElapsed()
    scrollToBottom()
  }
)

const attachOptions = [
  { label: '📷 图片', key: 'images' },
  { label: '🎵 音频', key: 'audio' },
  { label: '📄 文本', key: 'text' },
  { label: '📕 PDF', key: 'pdf' }
]

const displayModelName = computed(() => {
  const raw = currentModel.value?.name || selectedModel.value || '选择模型'
  return raw.replace(/\.gguf$/i, '').replace(/-/g, ' ')
})

function formatDuration(sec: number): string {
  if (sec < 60) return `${sec}s`
  const min = Math.floor(sec / 60)
  const s = sec % 60
  return s > 0 ? `${min}min ${s}s` : `${min}min`
}

const totalTokensDisplay = computed(() =>
  (isGenerating.value ? stats.value.outputTokens : lastTurnStats.value.tokens).toLocaleString()
)

const formattedDuration = computed(() =>
  formatDuration(isGenerating.value ? generationElapsedSeconds.value : lastTurnStats.value.seconds)
)

const liveTokensDisplay = computed(() =>
  (isGenerating.value ? stats.value.outputTokens : lastTurnStats.value.tokens).toLocaleString()
)

const liveDurationDisplay = computed(() =>
  formatDuration(isGenerating.value ? generationElapsedSeconds.value : lastTurnStats.value.seconds)
)

const liveSpeedDisplay = computed(() =>
  (isGenerating.value ? stats.value.speed : lastTurnStats.value.speed).toFixed(2)
)

const liveOutputDisplay = computed(() =>
  (isGenerating.value ? stats.value.outputTokens : lastTurnStats.value.tokens).toLocaleString()
)

const contextPercent = computed(() => {
  if (!stats.value.contextTotal) return 0
  return Math.min(100, Math.round((stats.value.contextUsed / stats.value.contextTotal) * 100))
})

const showLiveStats = computed(
  () => isGenerating.value || messages.value.length > 0
)

const showTopBar = computed(
  () => messages.value.length > 0 || isGenerating.value
)

const canSend = computed(
  () => Boolean(inputMessage.value.trim() && selectedModel.value)
)

const lastAssistantIndex = computed(() => {
  for (let i = messages.value.length - 1; i >= 0; i--) {
    if (messages.value[i].role === 'assistant') return i
  }
  return -1
})

function tickElapsed() {
  if (sessionStartTime.value) {
    elapsedSeconds.value = Math.floor((Date.now() - sessionStartTime.value) / 1000)
  }
}

function startElapsedTimer() {
  if (elapsedTimer) return
  elapsedTimer = setInterval(tickElapsed, 1000)
}

function stopElapsedTimer() {
  if (elapsedTimer) {
    clearInterval(elapsedTimer)
    elapsedTimer = null
  }
}

function tickGenerationElapsed() {
  if (generationStartTime.value) {
    generationElapsedSeconds.value = Math.floor(
      (Date.now() - generationStartTime.value) / 1000
    )
  }
}

function startGenerationTimer() {
  generationStartTime.value = Date.now()
  generationElapsedSeconds.value = 0
  if (generationTimer) clearInterval(generationTimer)
  generationTimer = setInterval(tickGenerationElapsed, 200)
}

function stopGenerationTimer() {
  if (generationTimer) {
    clearInterval(generationTimer)
    generationTimer = null
  }
  tickGenerationElapsed()
}

function snapshotTurnStats() {
  lastTurnStats.value = {
    tokens: stats.value.outputTokens,
    seconds: generationElapsedSeconds.value,
    speed: stats.value.speed
  }
}

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
    currentModel.value = option.model as typeof currentModel.value
    stats.value.contextTotal = (option.model as { contextLength: number }).contextLength * 1024
    showModelPicker.value = false
  }
}

const autoResize = () => {
  const el = textareaRef.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = `${Math.min(el.scrollHeight, 160)}px`
}

const runGeneration = async () => {
  isGenerating.value = true
  streamingContent.value = ''
  streamingReasoning.value = ''
  stats.value.outputTokens = 0
  stats.value.speed = 0
  startGenerationTimer()
  try {
    await streamChat()
  } catch (error) {
    message.error('生成失败: ' + (error as Error).message)
    messages.value.push({
      role: 'assistant',
      content: `生成失败：${(error as Error).message}`,
      timestamp: new Date().toLocaleTimeString()
    })
  } finally {
    isGenerating.value = false
    stopGenerationTimer()
    snapshotTurnStats()
    stopElapsedTimer()
    tickElapsed()
    await nextTick()
    scrollToBottom()
  }
}

const handleSend = async () => {
  if (!canSend.value || isGenerating.value) return

  if (!sessionStartTime.value) {
    sessionStartTime.value = Date.now()
    startElapsedTimer()
  }

  messages.value.push({
    role: 'user',
    content: inputMessage.value.trim(),
    timestamp: new Date().toLocaleTimeString(),
    attachments: pendingAttachments.value.length ? [...pendingAttachments.value] : undefined
  })
  inputMessage.value = ''
  pendingAttachments.value = []
  await nextTick()
  autoResize()
  scrollToBottom()
  await runGeneration()
}

const handleRegenerate = async () => {
  const idx = lastAssistantIndex.value
  if (idx < 0 || isGenerating.value) return
  messages.value.splice(idx, 1)
  await runGeneration()
}

const deleteLastAssistant = () => {
  const idx = lastAssistantIndex.value
  if (idx >= 0) messages.value.splice(idx, 1)
}

function formatAssistantForCopy(msg: ChatUiMessage): string {
  if (!msg.reasoning?.trim()) return msg.content
  return `【思考】\n${msg.reasoning}\n\n【回复】\n${msg.content}`
}

const copyAssistantMessage = async (msg: ChatUiMessage) => {
  try {
    await navigator.clipboard.writeText(formatAssistantForCopy(msg))
    message.success('已复制')
  } catch {
    message.error('复制失败')
  }
}

const handleNewChat = () => {
  messages.value = []
  streamingContent.value = ''
  streamingReasoning.value = ''
  sessionStartTime.value = null
  elapsedSeconds.value = 0
  generationStartTime.value = null
  generationElapsedSeconds.value = 0
  lastTurnStats.value = { tokens: 0, seconds: 0, speed: 0 }
  stopElapsedTimer()
  stopGenerationTimer()
  stats.value = {
    tokensPerSecond: 0,
    contextUsed: 0,
    contextTotal: (currentModel.value?.contextLength ?? 4) * 1024,
    outputTokens: 0,
    speed: 0
  }
}

const showModelInfo = () => {
  if (currentModel.value) {
    message.info(`模型: ${currentModel.value.name}`)
  } else {
    showModelPicker.value = true
  }
}

const notifySoon = (name: string) => message.info(`${name} 功能开发中`)

const pickAttachment = (key: string) => {
  showAttachMenu.value = false
  if (!fileInputRef.value) return
  const map: Record<string, string> = {
    images: 'image/*',
    audio: 'audio/*',
    text: '.txt,.md,.json',
    pdf: '.pdf'
  }
  fileInputRef.value.accept = map[key] || '*/*'
  fileInputRef.value.click()
}

const handleFileSelect = (e: Event) => {
  const input = e.target as HTMLInputElement
  if (input.files?.length) {
    pendingAttachments.value.push(...Array.from(input.files))
    message.success(`已添加 ${input.files.length} 个文件`)
  }
  input.value = ''
}

const scrollToBottom = () => {
  if (!scrollbarRef.value) return
  nextTick(() => {
    const container = scrollbarRef.value.$el?.querySelector('.n-scrollbar-container')
    if (container) container.scrollTop = container.scrollHeight
  })
}

onMounted(() => {
  loadModels()
})
</script>

<style scoped>
.chat-shell {
  display: flex;
  flex-direction: column;
  height: 600px;
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid #e8e8ec;
}

.chat-shell--full {
  height: 100vh;
  border-radius: 0;
  border: none;
}

/* 顶部状态栏 */
.chat-top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 20px;
  border-bottom: 1px solid #f0f0f2;
  flex-shrink: 0;
}

.top-bar-left {
  min-width: 0;
}

.model-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f3f3f5;
  border-radius: 999px;
  border: 1px solid #ebebed;
  max-width: 100%;
}

.model-pill--header {
  padding: 8px 12px;
}

.model-pill--inline {
  flex-shrink: 0;
  cursor: pointer;
  border: none;
  font: inherit;
  transition: background 0.15s;
}

.model-pill--inline:hover {
  background: #eaeaec;
}

.model-pill__text {
  font-size: 13px;
  font-weight: 500;
  color: #3d3d43;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}

.model-pill--header .model-pill__text {
  max-width: 280px;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: #6b6b73;
  border-radius: 8px;
  cursor: pointer;
  padding: 0;
}

.icon-btn:hover {
  background: #e8e8ec;
  color: #333;
}

.top-bar-metrics {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.metric {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #5c5c66;
  white-space: nowrap;
}

.metric svg {
  opacity: 0.55;
}

/* 消息区 */
.chat-body {
  flex: 1;
  min-height: 0;
  background: #fafafa;
}

.chat-scroll {
  height: 100%;
}

.messages-wrap {
  max-width: 820px;
  margin: 0 auto;
  padding: 20px 24px 12px;
}

.chat-empty {
  text-align: center;
  padding: 48px 16px;
  color: #9a9aa3;
}

.chat-empty__title {
  font-size: 15px;
  font-weight: 500;
  color: #6b6b73;
  margin: 0 0 6px;
}

.chat-empty__hint {
  font-size: 13px;
  margin: 0;
}

.msg-block {
  margin-bottom: 20px;
}

.msg-block__content {
  font-size: 15px;
  line-height: 1.65;
  color: #1a1a1e;
  white-space: pre-wrap;
  word-break: break-word;
}

.msg-block.user .msg-block__content {
  color: #1a1a1e;
}

.msg-block__content--user {
  color: #1a1a1e;
}

.msg-block__content--answer {
  font-size: 15px;
  line-height: 1.65;
  color: #1a1a1e;
}

.msg-block__content--pending {
  font-size: 14px;
  color: #9a9aa3;
  font-style: italic;
}

.msg-block.generating .msg-block__content--answer {
  color: #2a2a30;
}

.cursor {
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  50% { opacity: 0; }
}

/* 消息工具栏 */
.message-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 0 16px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  color: #1a1a1e;
  border-radius: 8px;
  cursor: pointer;
  padding: 0;
}

.toolbar-btn:hover {
  background: #efeff2;
}

/* 底部输入 */
.chat-composer-wrap {
  flex-shrink: 0;
  padding: 12px 20px 20px;
  background: #fff;
}

/* 输入框上方实时统计条 */
.live-stats-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px 20px;
  max-width: 820px;
  margin: 0 auto 10px;
  padding: 8px 12px;
  font-family: ui-monospace, 'Cascadia Code', Consolas, monospace;
  font-size: 13px;
  color: #6b6b73;
  line-height: 1.4;
}

.live-stat {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.live-stat svg {
  flex-shrink: 0;
  opacity: 0.55;
}

.pending-files {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  max-width: 820px;
  margin: 0 auto 8px;
}

.file-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: #f3f3f5;
  border-radius: 8px;
  font-size: 12px;
}

.file-chip button {
  border: none;
  background: none;
  cursor: pointer;
  color: #888;
  font-size: 14px;
  line-height: 1;
}

.chat-composer {
  display: flex;
  align-items: flex-end;
  gap: 10px;
  max-width: 820px;
  margin: 0 auto;
  padding: 12px 14px;
  background: #fff;
  border: 1px solid #e4e4e8;
  border-radius: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.btn-circle {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 1px solid #e0e0e4;
  background: #fff;
  color: #5c5c66;
  font-size: 20px;
  line-height: 1;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  transition: background 0.15s, border-color 0.15s;
}

.btn-add:hover:not(:disabled) {
  background: #f5f5f7;
}

.btn-add:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.composer-input {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  resize: none;
  font-size: 15px;
  line-height: 1.5;
  font-family: inherit;
  color: #1a1a1e;
  background: transparent;
  padding: 6px 0;
  max-height: 160px;
}

.composer-input::placeholder {
  color: #b0b0b8;
}

.btn-send {
  background: #d8d8dc;
  border-color: #d8d8dc;
  color: #fff;
}

.btn-send.active:not(:disabled) {
  background: #3d3d43;
  border-color: #3d3d43;
}

.btn-send:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-send svg {
  display: block;
}

.send-loading {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.attach-menu {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 140px;
}

.attach-menu button {
  text-align: left;
  padding: 8px 12px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
}

.attach-menu button:hover {
  background: #f3f3f5;
}

.hidden-input {
  display: none;
}

:global(.model-tooltip) {
  font-family: ui-monospace, 'Cascadia Code', Consolas, monospace;
  font-size: 12px;
}
</style>
