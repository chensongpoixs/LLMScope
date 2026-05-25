import { ref, type Ref } from 'vue'
import { streamChatCompletions, type ChatMessage } from '@/api/chat'
import { hasThinkingStart, parseThinkingContent } from '@/utils/thinking'

export interface ChatUiMessage {
  role: 'user' | 'assistant'
  content: string
  reasoning?: string
  timestamp: string
}

export interface ChatStats {
  tokensPerSecond: number
  contextUsed: number
  contextTotal: number
  outputTokens: number
  speed: number
}

/** 按字符粗估 token（避免对每个 SSE 片段各自 ceil 导致虚高） */
function estimateTokensFromChars(charCount: number): number {
  if (charCount <= 0) return 0
  return Math.round(charCount / 4)
}

function estimateContextTokens(
  msgs: ChatUiMessage[],
  reasoningStream: string,
  contentStream: string
): number {
  let chars = 0
  for (const m of msgs) {
    chars += m.content.length
    if (m.reasoning) chars += m.reasoning.length
  }
  chars += reasoningStream.length + contentStream.length
  return estimateTokensFromChars(chars)
}

/**
 * 计算解码阶段 t/s（排除 prefill 等待首 token 的时间）
 * 优先使用 llama.cpp timings.predicted_per_second
 */
function updateSpeedStats(
  stats: Ref<ChatStats>,
  opts: {
    outputChars: number
    firstTokenAt: number | null
    now: number
    timings?: { predicted_n?: number; predicted_per_second?: number; predicted_ms?: number }
    usage?: { completion_tokens?: number }
  }
) {
  const { outputChars, firstTokenAt, now, timings, usage } = opts

  if (typeof usage?.completion_tokens === 'number') {
    stats.value.outputTokens = usage.completion_tokens
  } else if (typeof timings?.predicted_n === 'number') {
    stats.value.outputTokens = timings.predicted_n
  } else {
    stats.value.outputTokens = estimateTokensFromChars(outputChars)
  }

  if (typeof timings?.predicted_per_second === 'number' && timings.predicted_per_second > 0) {
    stats.value.speed = timings.predicted_per_second
    stats.value.tokensPerSecond = timings.predicted_per_second
    return
  }

  if (
    typeof timings?.predicted_n === 'number' &&
    typeof timings?.predicted_ms === 'number' &&
    timings.predicted_ms > 0
  ) {
    stats.value.speed = (timings.predicted_n / timings.predicted_ms) * 1000
    stats.value.tokensPerSecond = stats.value.speed
    return
  }

  if (firstTokenAt !== null && outputChars > 0) {
    const decodeSec = (now - firstTokenAt) / 1000
    if (decodeSec > 0) {
      const tokens = stats.value.outputTokens
      stats.value.speed = tokens / decodeSec
      stats.value.tokensPerSecond = stats.value.speed
    }
  }
}

export function useChatStream(
  messages: Ref<ChatUiMessage[]>,
  selectedModel: Ref<string | undefined>,
  streamingContent: Ref<string>,
  streamingReasoning: Ref<string>,
  stats: Ref<ChatStats>,
  onChunk?: () => void
) {
  const isGenerating = ref(false)

  const streamChat = async () => {
    if (!selectedModel.value) {
      throw new Error('请先选择模型')
    }

    const apiMessages: ChatMessage[] = messages.value.map(m => {
      const item: ChatMessage = {
        role: m.role,
        content: m.content
      }
      if (m.role === 'assistant' && m.reasoning?.trim()) {
        item.reasoning_content = m.reasoning
      }
      return item
    })

    let outputChars = 0
    let firstTokenAt: number | null = null
    streamingContent.value = ''
    streamingReasoning.value = ''

    await streamChatCompletions(
      {
        model: selectedModel.value,
        messages: apiMessages,
        stream: true
      },
      (parts) => {
        const hasText = Boolean(parts.reasoning || parts.content)
        if (hasText) {
          if (firstTokenAt === null) firstTokenAt = Date.now()
          if (parts.reasoning) streamingReasoning.value += parts.reasoning
          if (parts.content) streamingContent.value += parts.content
          outputChars =
            streamingReasoning.value.length + streamingContent.value.length
        }

        const now = Date.now()
        updateSpeedStats(stats, {
          outputChars,
          firstTokenAt,
          now,
          timings: parts.timings,
          usage: parts.usage
        })

        stats.value.contextUsed = Math.min(
          estimateContextTokens(
            messages.value,
            streamingReasoning.value,
            streamingContent.value
          ),
          stats.value.contextTotal
        )
        onChunk?.()
      }
    )

    let finalContent = streamingContent.value
    let finalReasoning = streamingReasoning.value

    if (!finalReasoning.trim() && hasThinkingStart(finalContent)) {
      const parsed = parseThinkingContent(finalContent)
      finalReasoning = parsed.thinking || ''
      finalContent = parsed.cleanContent
    }

    if (!finalContent.trim() && !finalReasoning.trim()) {
      throw new Error(
        '模型未返回可显示文本。请确认 llama.cpp 服务与模型已正确加载。'
      )
    }

    messages.value.push({
      role: 'assistant',
      content: finalContent.trim(),
      reasoning: finalReasoning.trim() || undefined,
      timestamp: new Date().toLocaleTimeString()
    })
    streamingContent.value = ''
    streamingReasoning.value = ''
  }

  return { isGenerating, streamChat }
}
