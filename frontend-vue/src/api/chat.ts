import { config } from '@/config'
import { parseThinkingContent } from '@/utils/thinking'

const API_BASE = config.apiBaseUrl

export interface ChatMessage {
  role: 'user' | 'assistant' | 'system'
  content: string
  reasoning_content?: string
}

export interface ChatCompletionRequest {
  model: string
  messages: ChatMessage[]
  stream?: boolean
  temperature?: number
  max_tokens?: number
}

/** llama.cpp 流式 chunk 中的 timings（与 server 一致） */
export interface LlamaTimings {
  predicted_n?: number
  predicted_ms?: number
  predicted_per_second?: number
  predicted_per_token_ms?: number
  prompt_n?: number
  prompt_ms?: number
  prompt_per_second?: number
}

export interface LlamaUsage {
  completion_tokens?: number
  prompt_tokens?: number
  total_tokens?: number
}

/** 流式增量：思考、回复及服务端统计 */
export interface StreamChunkParts {
  reasoning?: string
  content?: string
  timings?: LlamaTimings
  usage?: LlamaUsage
}

export function extractDeltaParts(
  delta: Record<string, unknown> | undefined
): Pick<StreamChunkParts, 'reasoning' | 'content'> {
  if (!delta) return {}

  const parts: Pick<StreamChunkParts, 'reasoning' | 'content'> = {}
  if (typeof delta.reasoning_content === 'string' && delta.reasoning_content.length > 0) {
    parts.reasoning = delta.reasoning_content
  }
  if (typeof delta.content === 'string' && delta.content.length > 0) {
    parts.content = delta.content
  }
  if (typeof delta.text === 'string' && delta.text.length > 0) {
    parts.content = (parts.content || '') + delta.text
  }
  return parts
}

function parseTimings(raw: unknown): LlamaTimings | undefined {
  if (!raw || typeof raw !== 'object') return undefined
  const t = raw as Record<string, unknown>
  const timings: LlamaTimings = {}
  if (typeof t.predicted_n === 'number') timings.predicted_n = t.predicted_n
  if (typeof t.predicted_ms === 'number') timings.predicted_ms = t.predicted_ms
  if (typeof t.predicted_per_second === 'number') {
    timings.predicted_per_second = t.predicted_per_second
  }
  if (typeof t.predicted_per_token_ms === 'number') {
    timings.predicted_per_token_ms = t.predicted_per_token_ms
  }
  if (typeof t.prompt_n === 'number') timings.prompt_n = t.prompt_n
  if (typeof t.prompt_ms === 'number') timings.prompt_ms = t.prompt_ms
  if (typeof t.prompt_per_second === 'number') timings.prompt_per_second = t.prompt_per_second
  return Object.keys(timings).length > 0 ? timings : undefined
}

function parseUsage(raw: unknown): LlamaUsage | undefined {
  if (!raw || typeof raw !== 'object') return undefined
  const u = raw as Record<string, unknown>
  const usage: LlamaUsage = {}
  if (typeof u.completion_tokens === 'number') usage.completion_tokens = u.completion_tokens
  if (typeof u.prompt_tokens === 'number') usage.prompt_tokens = u.prompt_tokens
  if (typeof u.total_tokens === 'number') usage.total_tokens = u.total_tokens
  return Object.keys(usage).length > 0 ? usage : undefined
}

export function extractPartsFromChunk(parsed: Record<string, unknown>): StreamChunkParts {
  const result: StreamChunkParts = {}

  const rootTimings = parseTimings(parsed.timings)
  const rootUsage = parseUsage(parsed.usage)
  if (rootTimings) result.timings = rootTimings
  if (rootUsage) result.usage = rootUsage

  const choices = parsed.choices as Record<string, unknown>[] | undefined
  const choice = choices?.[0]
  if (!choice) return result

  Object.assign(result, extractDeltaParts(choice.delta as Record<string, unknown> | undefined))

  const message = choice.message as Record<string, unknown> | undefined
  if (message) {
    if (typeof message.reasoning_content === 'string') {
      result.reasoning = message.reasoning_content
    }
    if (typeof message.content === 'string') {
      result.content = message.content
    }
  } else if (typeof choice.text === 'string') {
    result.content = choice.text
  }

  return result
}

function parseSSELine(line: string, onParts: (parts: StreamChunkParts) => void): boolean {
  const trimmed = line.trim()
  if (!trimmed) return false

  let payload = trimmed
  if (trimmed.startsWith('data:')) {
    payload = trimmed.slice(5).trim()
  }

  if (payload === '[DONE]') return true

  if (!payload.startsWith('{')) return false

  try {
    const parsed = JSON.parse(payload) as Record<string, unknown>
    const parts = extractPartsFromChunk(parsed)
    if (
      parts.reasoning ||
      parts.content ||
      parts.timings ||
      parts.usage
    ) {
      onParts(parts)
    }
  } catch {
    // 忽略不完整 JSON
  }
  return false
}

function parseStreamBody(text: string, onParts: (parts: StreamChunkParts) => void): void {
  for (const line of text.split('\n')) {
    if (parseSSELine(line, onParts)) break
  }
}

/** 流式聊天：经后端转发到 llama.cpp */
export async function streamChatCompletions(
  request: ChatCompletionRequest,
  onParts: (parts: StreamChunkParts) => void,
  signal?: AbortSignal
): Promise<void> {
  const response = await fetch(`${API_BASE}/api/inference/chat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ...request, stream: true }),
    signal
  })

  if (!response.ok) {
    const errText = await response.text()
    let msg = `HTTP ${response.status}`
    try {
      const errJson = JSON.parse(errText)
      msg = (errJson.message as string) || (errJson.error as string) || msg
    } catch {
      if (errText) msg = errText
    }
    throw new Error(msg)
  }

  const contentType = response.headers.get('Content-Type') || ''

  if (!contentType.includes('text/event-stream')) {
    const data = (await response.json()) as Record<string, unknown>
    const parts = extractPartsFromChunk(data)
    if (parts.reasoning || parts.content || parts.timings || parts.usage) {
      onParts(parts)
    } else {
      const choices = data.choices as Record<string, unknown>[] | undefined
      const choice = choices?.[0]
      const raw = (choice?.message as Record<string, unknown> | undefined)?.content
      if (typeof raw === 'string') {
        const parsed = parseThinkingContent(raw)
        if (parsed.thinking) onParts({ reasoning: parsed.thinking })
        if (parsed.cleanContent) onParts({ content: parsed.cleanContent })
      }
    }
    return
  }

  const reader = response.body?.getReader()
  if (!reader) throw new Error('无法读取响应流')

  const decoder = new TextDecoder()
  let buffer = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })
    const lines = buffer.split('\n')
    buffer = lines.pop() || ''

    for (const line of lines) {
      if (parseSSELine(line, onParts)) return
    }
  }

  if (buffer.trim()) parseStreamBody(buffer, onParts)
}
