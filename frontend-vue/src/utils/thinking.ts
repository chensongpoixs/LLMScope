/** 与 llama.cpp WebUI 一致的思考标签解析 */

export const THINK_TAG_START = '<think>'
export const THINK_TAG_END = '</think>'

export function parseThinkingContent(content: string): {
  thinking: string | null
  cleanContent: string
} {
  const startIdx = content.indexOf(THINK_TAG_START)
  if (startIdx === -1) {
    return { thinking: null, cleanContent: content }
  }

  const endIdx = content.indexOf(THINK_TAG_END)
  if (endIdx === -1) {
    const thinking = content.slice(startIdx + THINK_TAG_START.length).trim()
    const cleanContent = content.slice(0, startIdx).trim()
    return { thinking: thinking || null, cleanContent }
  }

  const thinking = content
    .slice(startIdx + THINK_TAG_START.length, endIdx)
    .trim()
  const before = content.slice(0, startIdx).trim()
  const after = content.slice(endIdx + THINK_TAG_END.length).trim()
  const cleanContent = [before, after].filter(Boolean).join('\n\n')

  return {
    thinking: thinking || null,
    cleanContent
  }
}

export function extractPartialThinking(content: string): {
  thinking: string | null
  remainingContent: string
} {
  const startIndex = content.indexOf(THINK_TAG_START)
  if (startIndex === -1) {
    return { thinking: null, remainingContent: content }
  }

  const endIndex = content.indexOf(THINK_TAG_END)
  if (endIndex === -1) {
    return {
      thinking: content.slice(startIndex + THINK_TAG_START.length),
      remainingContent: content.slice(0, startIndex)
    }
  }

  const parsed = parseThinkingContent(content)
  return {
    thinking: parsed.thinking,
    remainingContent: parsed.cleanContent
  }
}

export function hasThinkingStart(content: string): boolean {
  return content.includes(THINK_TAG_START)
}
