<script setup lang="ts">
import { MouseEventService } from "../bindings/pig_assistant";
import { onMounted, ref } from "vue";
import FloatingBall from "./components/FloatingBall.vue";
import { useMousePosition } from "./composables/useMousePosition";
import { useMouseEvents } from "./composables/useMouseEvents";
import CoordGrid from "./CoordGrid.vue";

// 拖拽状态管理
const dragState = ref({
  isDragging: false,
  startX: 0,
  startY: 0,
  startLeft: 0,
  startTop: 0,
  startTime:0,
  canShowCoordGrid:false
});


const appContainerRef = ref<HTMLElement | null>(null);
const coordGrid = ref<InstanceType<typeof CoordGrid> | null>(null);

// 使用拖拽composable - 自动检测多个app-region拖拽样式元素
useMousePosition({
  autoDetect: true,
  selector: [".mouse-drag", ".mouse-interactive_drag"],
  onDragStart: (element) => {
    console.log("拖拽开始", element);
  },
  onDragEnd: (element) => {
    console.log("拖拽结束", element);
  },
  onDrag: (position, element) => {
    // console.log("拖拽中", position, element)

  },
  // 自定义拖拽行为 - 实现app-container拖拽
  onCustomDrag: (e: MouseEvent, element: HTMLElement) => {
    console.log("自定义拖拽开始", element);

    // 检查是否是app-container元素

    // 初始化拖拽状态
    dragState.value.isDragging = true;
    dragState.value.startTime = Date.now();
    dragState.value.startX = e.clientX;
    dragState.value.startY = e.clientY;

    // 获取元素当前位置
    const rect = appContainerRef.value?.getBoundingClientRect();
    if (!rect) return false;
    dragState.value.startLeft = rect.left;
    dragState.value.startTop = rect.top;

    // 禁用默认拖拽行为，使用自定义实现
    return false;

    // 其他元素使用默认行为
  },
  onCustomDragMove: (e: MouseEvent, element: HTMLElement) => {
    // 只对app-container元素进行自定义拖拽处理
    if (appContainerRef.value && dragState.value.isDragging) {
      // 计算鼠标移动距离
      const deltaX = e.clientX - dragState.value.startX;
      const deltaY = e.clientY - dragState.value.startY;

      // 计算新位置
      let newX = dragState.value.startLeft + deltaX;
      let newY = dragState.value.startTop + deltaY;

      // 边界限制
      const rect = appContainerRef.value?.getBoundingClientRect();
      if (!rect) return false;
      const maxX = window.innerWidth - rect.width;
      const maxY = window.innerHeight - rect.height;

      newX = Math.max(0, Math.min(maxX, newX));
      newY = Math.max(0, Math.min(maxY, newY));

      // 添加拖拽时的视觉反馈
      appContainerRef.value.style.transform = "scale(1.02)";
      appContainerRef.value.style.boxShadow = "0 12px 40px rgba(0, 0, 0, 0.2)";
      appContainerRef.value.style.transition = "none"; // 禁用过渡动画以获得流畅拖拽

      // 更新元素位置
      appContainerRef.value.style.left = newX + "px";
      appContainerRef.value.style.top = newY + "px";

      if(!dragState.value.canShowCoordGrid){
        if(Date.now() - dragState.value.startTime > 500){

          dragState.value.canShowCoordGrid = true
        }
      }else{
        coordGrid.value?.setActiveByClientPoint(e.clientX,e.clientY);
      }

      //console.log("自定义拖拽移动", newX, newY)

      // 阻止默认拖拽行为
      return false;
    }

    // 其他元素使用默认行为
  },
  onCustomDragEnd: (e: MouseEvent, element: HTMLElement) => {
    console.log("自定义拖拽结束", element);

    if (appContainerRef.value) {
      // 重置拖拽状态
      dragState.value.isDragging = false;
      dragState.value.canShowCoordGrid = false;
      dragState.value.startTime = 0;

      // 恢复拖拽前的视觉状态
      appContainerRef.value.style.transform = "scale(1)";
      appContainerRef.value.style.boxShadow = "0 8px 32px rgba(0, 0, 0, 0.1)";
      appContainerRef.value.style.transition =
        "transform 0.2s ease, box-shadow 0.2s ease";

      // 可以在这里添加拖拽结束后的特殊处理
      // 比如：保存位置、触发动画、检查吸附等
      console.log(
        "app-container拖拽结束，最终位置:",
        appContainerRef.value.style.left,
        appContainerRef.value.style.top
      );
    }
  },
});

// 使用鼠标事件composable - 自动检测多个app-region样式元素
useMouseEvents({
  autoDetect: true,
  selector: [".mouse-interactive", ".mouse-interactive_drag"],
  onMouseEnter: (e: MouseEvent) => {
    if (!dragState.value.isDragging) {
      MouseEventService.MousePassthroughWithMove(false);
      console.log("mouse enter", e.target);
    }
  },
  onMouseLeave: (e: MouseEvent) => {
    if (!dragState.value.isDragging) {
      MouseEventService.MousePassthroughWithMove(true);
    }
    console.log("mouse leave", e.target);
  },
  onMouseMove: (e: MouseEvent) => {
    // console.log("mousemove", e.x, e.y, e.target)
  },
});

onMounted(() => {
  requestAnimationFrame(() => {
    MouseEventService.MousePassthroughWithMove(true);
  });
});
</script>

<template>
  <div class="app-container">
    <div ref="appContainerRef" class="core-content">
      <!-- 悬浮球组件 -->
      <FloatingBall />
    </div>

    <transition name="fade">
      <CoordGrid ref="coordGrid" v-show="dragState.canShowCoordGrid" />
    </transition>
  </div>
</template>

<style scoped>
.app-container {
  width: 100vw;
  height: 100vh;
}

.core-content {
  position: absolute;
  right: 0;
  bottom: 0;
  width: 200px;
  height: 100px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  pointer-events: auto;
  /* cursor: move; */
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  z-index: 1000;
  overflow: hidden;
}

.app-container:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.app-container:active {
  /* transform: scale(0.98); */
  /* cursor: grabbing; */
}

.drag-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 14px;
  font-weight: 500;
  user-select: none;
}

.drag-title {
  flex: 1;
}

.drag-indicator {
  font-size: 16px;
  opacity: 0.7;
  /* cursor: grab; */
}

.drag-indicator:active {
  /* cursor: grabbing; */
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
