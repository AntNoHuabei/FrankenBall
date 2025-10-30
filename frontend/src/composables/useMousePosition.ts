import { reactive, onMounted, onUnmounted, ref, type Ref } from 'vue'

export interface DragState {
  isDragging: boolean
  position: { x: number; y: number }
  dragStart: { x: number; y: number }
  containerSize: { width: number; height: number }
}

export interface UseMousePositionOptions {
  initialPosition?: { x: number; y: number }
  containerSize?: { width: number; height: number }
  selector?: string | string[]
  autoDetect?: boolean
  onDragStart?: (element: HTMLElement) => void
  onDragEnd?: (element: HTMLElement) => void
  onDrag?: (position: { x: number; y: number }, element: HTMLElement) => void
  onCustomDrag?: (e: MouseEvent, element: HTMLElement) => boolean | void
  onCustomDragMove?: (e: MouseEvent, element: HTMLElement) => boolean | void
  onCustomDragEnd?: (e: MouseEvent, element: HTMLElement) => void
  disableDefaultDrag?: boolean
}

export function useMousePosition(options: UseMousePositionOptions = {}) {
  const {
    initialPosition = { x: 100, y: 100 },
    containerSize = { width: 300, height: 200 },
    selector = '[style*="app-region:drag"]',
    autoDetect = false,
    onDragStart,
    onDragEnd,
    onDrag,
    onCustomDrag,
    onCustomDragMove,
    onCustomDragEnd,
    disableDefaultDrag = false
  } = options

  // 使用reactive管理拖拽状态
  const dragState = reactive<DragState>({
    isDragging: false,
    position: { ...initialPosition },
    dragStart: { x: 0, y: 0 },
    containerSize: { ...containerSize }
  })

  // MutationObserver用于监听DOM变化
  const observer = ref<MutationObserver | null>(null)
  
  // 已绑定的元素集合
  const boundElements = ref<Set<HTMLElement>>(new Set())
  
  // 当前拖拽的元素
  const currentDragElement = ref<HTMLElement | null>(null)

  // 将选择器转换为数组
  const selectors = Array.isArray(selector) ? selector : [selector]

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
    

    console.log("allElements", allElements)

    // 绑定所有匹配的元素
    allElements.forEach(el => {
      if (matchesSelector(el)) {
        bindDragEvents(el)
      }
    })
  }

  // 扫描并解绑不匹配的元素
  const scanAndUnbindElements = () => {
    boundElements.value.forEach(el => {
      if (!matchesSelector(el) || !document.contains(el)) {
        unbindDragEvents(el)
      }
    })
  }

  // 绑定拖拽事件到元素
  const bindDragEvents = (el: HTMLElement) => {
    if (boundElements.value.has(el)) return
    
    el.addEventListener('mousedown', (e) => startDrag(e, el))
    boundElements.value.add(el)
  }

  // 解绑拖拽事件
  const unbindDragEvents = (el: HTMLElement) => {
    if (!boundElements.value.has(el)) return
    
    el.removeEventListener('mousedown', (e) => startDrag(e, el))
    boundElements.value.delete(el)
  }

  // 解绑所有事件
  const unbindAllEvents = () => {
    boundElements.value.forEach(el => {
      el.removeEventListener('mousedown', (e) => startDrag(e, el))
    })
    boundElements.value.clear()
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
              bindDragEvents(target)
            } else {
              unbindDragEvents(target)
            }
          }
        }
        
        // 监听子节点添加
        if (mutation.type === 'childList') {
          mutation.addedNodes.forEach((node) => {
            if (node.nodeType === Node.ELEMENT_NODE) {
              const element = node as HTMLElement
              if (matchesSelector(element)) {
                bindDragEvents(element)
              }
              // 检查子元素
              selectors.forEach(sel => {
                const children = element.querySelectorAll(sel) as NodeListOf<HTMLElement>
                children.forEach(child => {
                  if (matchesSelector(child)) {
                    bindDragEvents(child)
                  }
                })
              })
            }
          })
          
          // 监听子节点移除
          mutation.removedNodes.forEach((node) => {
            if (node.nodeType === Node.ELEMENT_NODE) {
              const element = node as HTMLElement
              unbindDragEvents(element)
              // 移除子元素
              selectors.forEach(sel => {
                const children = element.querySelectorAll(sel) as NodeListOf<HTMLElement>
                children.forEach(child => {
                  unbindDragEvents(child)
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

  // 拖拽开始
  const startDrag = (e: MouseEvent, element: HTMLElement) => {

    do {
      if (onCustomDrag) {
        const result = onCustomDrag(e, element)
        // 如果返回false，则阻止默认拖拽行为
        if (result === false) {
          break;
        }
      }

      // 如果禁用了默认拖拽行为，则不执行后续逻辑
      if (disableDefaultDrag) {
        return
      }


      // 获取元素的位置和尺寸
      const rect = element.getBoundingClientRect()
      dragState.position = { x: rect.left, y: rect.top }
      dragState.containerSize = { width: rect.width, height: rect.height }

      dragState.dragStart = {
        x: e.clientX - dragState.position.x,
        y: e.clientY - dragState.position.y
      }
    }while (false)
    // 如果有自定义拖拽回调，先执行它

    dragState.isDragging = true
    currentDragElement.value = element


    document.addEventListener('mousemove', onDragMove)
    document.addEventListener('mouseup', stopDrag)

    onDragStart?.(element)

  }

  // 拖拽中
  const onDragMove = (e: MouseEvent) => {
    if (!dragState.isDragging || !currentDragElement.value) return

    do {
      // 如果有自定义拖拽移动回调，先执行它
      if (onCustomDragMove) {
        const result = onCustomDragMove(e, currentDragElement.value)
        // 如果返回false，则阻止默认拖拽行为
        if (result === false) {
          break
        }
      }

      const newX = e.clientX - dragState.dragStart.x
      const newY = e.clientY - dragState.dragStart.y

      // 限制在屏幕范围内
      dragState.position.x = Math.max(0, Math.min(window.innerWidth - dragState.containerSize.width, newX))
      dragState.position.y = Math.max(0, Math.min(window.innerHeight - dragState.containerSize.height, newY))

      // 更新元素位置
      if (currentDragElement.value) {
        currentDragElement.value.style.left = dragState.position.x + 'px'
        currentDragElement.value.style.top = dragState.position.y + 'px'
      }
    }while (false);

    
    onDrag?.(dragState.position, currentDragElement.value)
  }

  // 拖拽结束
  const stopDrag = (e?: MouseEvent) => {
    dragState.isDragging = false
    const element = currentDragElement.value
    currentDragElement.value = null
    
    document.removeEventListener('mousemove', onDragMove)
    document.removeEventListener('mouseup', stopDrag)
    
    // 如果有自定义拖拽结束回调，执行它
    if (onCustomDragEnd && element && e) {
      onCustomDragEnd(e, element)
    }
    
    if (element) {
      onDragEnd?.(element)
    }
  }

  // 更新容器尺寸
  const updateContainerSize = (size: { width: number; height: number }) => {
    dragState.containerSize = { ...size }
  }

  // 设置位置
  const setPosition = (position: { x: number; y: number }) => {
    dragState.position = { ...position }
  }

  // 重新绑定事件（用于元素变化时）
  const rebindEvents = () => {
    unbindAllEvents()
    scanAndBindElements()
  }

  // 手动扫描当前DOM
  const rescan = () => {
    scanAndUnbindElements()
    scanAndBindElements()
  }

  // 清理事件监听器
  const cleanup = () => {
    document.removeEventListener('mousemove', onDragMove)
    document.removeEventListener('mouseup', stopDrag)
    stopObserver()
    unbindAllEvents()
  }

  onMounted(() => {
    if (autoDetect) {
      createObserver()
      scanAndBindElements()
    }
  })

  onUnmounted(() => {
    cleanup()
  })

  return {
    dragState,
    boundElements: boundElements.value,
    currentDragElement: currentDragElement.value,
    startDrag,
    stopDrag,
    updateContainerSize,
    setPosition,
    cleanup,
    rescan,
    createObserver,
    stopObserver,
    rebindEvents
  }
}