import { ref, onMounted, onUnmounted, type Ref, watch } from 'vue'

export interface MouseEventHandlers {
  onMouseEnter?: (e: MouseEvent) => void
  onMouseLeave?: (e: MouseEvent) => void
  onMouseMove?: (e: MouseEvent) => void
}

export interface UseMouseEventsOptions extends MouseEventHandlers {
  element?: Ref<HTMLElement | null> | HTMLElement | null
  enabled?: boolean
  selector?: string | string[]
  autoDetect?: boolean
}

export function useMouseEvents(options: UseMouseEventsOptions = {}) {
  const {
    element,
    enabled: initialEnabled = true,
    selector = '[style*="app-region:interactive"]',
    autoDetect = false,
    onMouseEnter,
    onMouseLeave,
    onMouseMove
  } = options

  // 创建元素引用
  const elementRef = ref<HTMLElement | null>(null)
  
  // 当前绑定的元素
  const currentElement = element || elementRef
  
  // MutationObserver用于监听DOM变化
  const observer = ref<MutationObserver | null>(null)
  
  // 已绑定的元素集合
  const boundElements = ref<Set<HTMLElement>>(new Set())
  
  // 启用状态
  let enabled = initialEnabled

  // 将选择器转换为数组
  const selectors = Array.isArray(selector) ? selector : [selector]
  
  // 事件处理器
  const handleMouseEnter = (e: MouseEvent) => {
    if (!enabled) return
    onMouseEnter?.(e)
  }

  const handleMouseLeave = (e: MouseEvent) => {
    if (!enabled) return
    onMouseLeave?.(e)
  }

  const handleMouseMove = (e: MouseEvent) => {
    if (!enabled) return
    onMouseMove?.(e)
  }

  // 绑定事件到元素
  const bindEvents = (el: HTMLElement) => {
    if (boundElements.value.has(el)) return
    
    el.addEventListener('mouseenter', handleMouseEnter)
    el.addEventListener('mouseleave', handleMouseLeave)
    el.addEventListener('mousemove', handleMouseMove)
    boundElements.value.add(el)
  }

  // 解绑事件
  const unbindEvents = (el: HTMLElement) => {
    if (!boundElements.value.has(el)) return
    
    el.removeEventListener('mouseenter', handleMouseEnter)
    el.removeEventListener('mouseleave', handleMouseLeave)
    el.removeEventListener('mousemove', handleMouseMove)
    boundElements.value.delete(el)
  }

  // 解绑所有事件
  const unbindAllEvents = () => {
    boundElements.value.forEach(el => {
      el.removeEventListener('mouseenter', handleMouseEnter)
      el.removeEventListener('mouseleave', handleMouseLeave)
      el.removeEventListener('mousemove', handleMouseMove)
    })
    boundElements.value.clear()
  }

  // 检查元素是否匹配任一选择器（直接使用原生选择器匹配）
  const matchesSelector = (el: HTMLElement): boolean => {
    return selectors.some((sel) => el.matches(sel))
  }

  // 扫描并绑定匹配的元素
  const scanAndBindElements = () => {
    // 使用所有选择器查询元素
    const allElements = new Set<HTMLElement>()
    
    selectors.forEach(sel => {
      const elements = document.querySelectorAll(sel) as NodeListOf<HTMLElement>
      elements.forEach(el => allElements.add(el))
    })
    
    // 绑定所有匹配的元素
    allElements.forEach(el => {
      if (matchesSelector(el)) {
        bindEvents(el)
      }
    })
  }

  // 扫描并解绑不匹配的元素
  const scanAndUnbindElements = () => {
    boundElements.value.forEach(el => {
      if (!matchesSelector(el) || !document.contains(el)) {
        unbindEvents(el)
      }
    })
  }

  // 创建MutationObserver监听DOM变化
  const createObserver = () => {
    if (observer.value) return
    
    observer.value = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        // 监听属性变化（包括style属性）
        if (mutation.type === 'attributes') {
          const target = mutation.target as HTMLElement
          if (mutation.attributeName === 'style') {
            if (matchesSelector(target)) {
              bindEvents(target)
            } else {
              unbindEvents(target)
            }
          }
        }
        
        // 监听子节点添加
        if (mutation.type === 'childList') {
          mutation.addedNodes.forEach((node) => {
            if (node.nodeType === Node.ELEMENT_NODE) {
              const element = node as HTMLElement
              if (matchesSelector(element)) {
                bindEvents(element)
              }
              // 检查子元素
              selectors.forEach(sel => {
                const children = element.querySelectorAll(sel) as NodeListOf<HTMLElement>
                children.forEach(child => {
                  if (matchesSelector(child)) {
                    bindEvents(child)
                  }
                })
              })
            }
          })
          
          // 监听子节点移除
          mutation.removedNodes.forEach((node) => {
            if (node.nodeType === Node.ELEMENT_NODE) {
              const element = node as HTMLElement
              unbindEvents(element)
              // 移除子元素
              selectors.forEach(sel => {
                const children = element.querySelectorAll(sel) as NodeListOf<HTMLElement>
                children.forEach(child => {
                  unbindEvents(child)
                })
              })
            }
          })
        }
      })
    })
    
    observer.value.observe(document.body, {
      childList: true,
      subtree: true,
      attributes: true,
      attributeFilter: ['style']
    })
  }

  // 停止观察
  const stopObserver = () => {
    if (observer.value) {
      observer.value.disconnect()
      observer.value = null
    }
  }

  // 重新绑定事件（用于元素变化时）
  const rebindEvents = () => {
    unbindAllEvents()
    scanAndBindElements()
  }

  // 启用/禁用事件监听
  const setEnabled = (value: boolean) => {
    enabled = value
  }

  // 更新事件处理器
  const updateHandlers = (newHandlers: MouseEventHandlers) => {
    Object.assign(options, newHandlers)
  }

  // 手动扫描当前DOM
  const rescan = () => {
    scanAndUnbindElements()
    scanAndBindElements()
  }

  onMounted(() => {
    if (autoDetect) {
      createObserver()
      scanAndBindElements()
    } else {
      const el = ref(currentElement).value || currentElement as HTMLElement
      if (el) {
        bindEvents(el)
      }
    }
  })

  onUnmounted(() => {
    stopObserver()
    unbindAllEvents()
  })

  return {
    elementRef,
    boundElements: boundElements.value,
    bindEvents,
    unbindEvents,
    unbindAllEvents,
    rebindEvents,
    setEnabled,
    updateHandlers,
    rescan,
    createObserver,
    stopObserver
  }
}

// 创建多个独立的鼠标事件监听器
export function createIsolatedMouseEvents(options: UseMouseEventsOptions = {}) {
  return useMouseEvents(options)
}

// 批量管理多个鼠标事件监听器
export class MouseEventManager {
  private listeners: Map<string, ReturnType<typeof useMouseEvents>> = new Map()

  addListener(id: string, options: UseMouseEventsOptions) {
    if (this.listeners.has(id)) {
      this.removeListener(id)
    }
    
    const listener = useMouseEvents(options)
    this.listeners.set(id, listener)
    return listener
  }

  removeListener(id: string) {
    const listener = this.listeners.get(id)
    if (listener) {
      listener.unbindAllEvents()
      this.listeners.delete(id)
    }
  }

  updateListener(id: string, handlers: MouseEventHandlers) {
    const listener = this.listeners.get(id)
    if (listener) {
      listener.updateHandlers(handlers)
    }
  }

  enableListener(id: string, enabled: boolean) {
    const listener = this.listeners.get(id)
    if (listener) {
      listener.setEnabled(enabled)
    }
  }

  rescanListener(id: string) {
    const listener = this.listeners.get(id)
    if (listener) {
      listener.rescan()
    }
  }

  rescanAll() {
    this.listeners.forEach(listener => listener.rescan())
  }

  clearAll() {
    this.listeners.forEach(listener => listener.unbindAllEvents())
    this.listeners.clear()
  }

  getListener(id: string) {
    return this.listeners.get(id)
  }

  getBoundElementsCount(): number {
    let total = 0
    this.listeners.forEach(listener => {
      total += listener.boundElements.size
    })
    return total
  }

  getAllBoundElements(): HTMLElement[] {
    const elements: HTMLElement[] = []
    this.listeners.forEach(listener => {
      elements.push(...listener.boundElements)
    })
    return elements
  }
}
