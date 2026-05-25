export const KeyboardShortcuts = {
  SEARCH: 'f',
  LAYER_VIEW: '1',
  MODULE_VIEW: '2',
  TENSOR_VIEW: '3',
  ELEMENT_VIEW: '4',
  RESET_CAMERA: ' ',
  UNDO: 'z'
}

export function setupKeyboardShortcuts(handlers: Record<string, () => void>) {
  const handleKeyDown = (e: KeyboardEvent) => {
    const key = e.key.toLowerCase()
    
    if (e.ctrlKey && key === KeyboardShortcuts.UNDO) {
      e.preventDefault()
      handlers.undo?.()
    } else if (key === KeyboardShortcuts.SEARCH) {
      e.preventDefault()
      handlers.search?.()
    } else if (key === KeyboardShortcuts.LAYER_VIEW) {
      handlers.layerView?.()
    } else if (key === KeyboardShortcuts.MODULE_VIEW) {
      handlers.moduleView?.()
    } else if (key === KeyboardShortcuts.TENSOR_VIEW) {
      handlers.tensorView?.()
    } else if (key === KeyboardShortcuts.ELEMENT_VIEW) {
      handlers.elementView?.()
    } else if (key === KeyboardShortcuts.RESET_CAMERA) {
      handlers.resetCamera?.()
    }
  }

  window.addEventListener('keydown', handleKeyDown)

  return () => {
    window.removeEventListener('keydown', handleKeyDown)
  }
}
