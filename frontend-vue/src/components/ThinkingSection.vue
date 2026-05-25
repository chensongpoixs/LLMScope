<template>
  <div v-if="thinking" class="thinking-section" :class="{ 'is-streaming': isStreaming }">
    <button type="button" class="thinking-header" @click="expanded = !expanded">
      <span class="thinking-header__left">
        <svg class="thinking-icon" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.8">
          <path d="M9.5 2A5.5 5.5 0 0 0 4 7.5c0 1.58.67 3 1.74 4.01L5 14h4l1.26-2.99A5.48 5.48 0 0 0 9.5 13 5.5 5.5 0 0 0 15 7.5 5.5 5.5 0 0 0 9.5 2z" />
          <path d="M12 14v2M10 20h5" />
        </svg>
        <span class="thinking-label">
          {{ isStreaming ? 'Thinking...' : 'Thinking summary' }}
        </span>
      </span>
      <svg
        class="thinking-chevron"
        :class="{ expanded }"
        viewBox="0 0 24 24"
        width="16"
        height="16"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <path d="M6 9l6 6 6-6" />
      </svg>
    </button>
    <div v-show="expanded || isStreaming" class="thinking-body">
      <div class="thinking-text">{{ thinking }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    thinking: string
    isStreaming?: boolean
  }>(),
  { isStreaming: false }
)

const expanded = ref(false)

watch(
  () => props.isStreaming,
  (streaming) => {
    if (streaming) expanded.value = true
  },
  { immediate: true }
)
</script>

<style scoped>
.thinking-section {
  margin-bottom: 16px;
  border: 1px solid #e8e8ec;
  border-radius: 12px;
  background: rgba(243, 243, 245, 0.65);
  overflow: hidden;
}

.thinking-section.is-streaming {
  border-color: #d8d8de;
}

.thinking-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 12px 14px;
  border: none;
  background: transparent;
  cursor: pointer;
  font: inherit;
  color: #6b6b73;
  text-align: left;
}

.thinking-header:hover {
  background: rgba(0, 0, 0, 0.02);
}

.thinking-header__left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.thinking-icon {
  flex-shrink: 0;
  opacity: 0.7;
}

.thinking-label {
  font-size: 13px;
  font-weight: 500;
}

.thinking-chevron {
  flex-shrink: 0;
  transition: transform 0.2s ease;
  opacity: 0.5;
}

.thinking-chevron.expanded {
  transform: rotate(180deg);
}

.thinking-body {
  border-top: 1px solid #e8e8ec;
  padding: 12px 14px 14px;
}

.thinking-text {
  font-size: 12px;
  line-height: 1.6;
  color: #5c5c66;
  white-space: pre-wrap;
  word-break: break-word;
  max-height: 320px;
  overflow-y: auto;
}

.is-streaming .thinking-text::after {
  content: '▊';
  animation: blink 1s step-end infinite;
  margin-left: 2px;
}

@keyframes blink {
  50% { opacity: 0; }
}
</style>
