<template>
  <div ref="componentElement" id="componentElement" class="absolute-component-window">
    <template v-for="ws in attachWindows" :key="ws.id">
      <Window
        :style="computedStyle(ws)"
        :ref="(el: any) => bindWindowEl(ws.id, el?.$el)"
        @close="closeWindow(ws.id)"
        @minimize="minimizeWindow(ws.id)"
        @maximize="handleMaximize"
        :win-id="ws.id"
        :title="ws.windowConfig.title"
        v-show="ws.isVisible"
      >
        <component :is="ws.windowConfig.component" v-bind="{ winId: ws.id }" />
      </Window>
    </template>

    <Teleport to="body">
      <template v-for="(ws, index) in notificationWindows" :key="ws.id">
        <Window
          class="mouse-interactive"
          :style="getNotificationStyle(ws, index)"
          :ref="(el: any) => bindWindowEl(ws.id, el?.$el)"
          @close="closeWindow(ws.id)"
          @minimize="minimizeWindow(ws.id)"
          @maximize="handleMaximize"
          :win-id="ws.id"
          :title="ws.windowConfig.title"
          v-show="ws.isVisible"
        >
          <component :is="ws.windowConfig.component" v-bind="{ winId: ws.id }" />
        </Window>
      </template>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import {WindowState} from "../types/window";
import {useWindowManager} from "../composables-local/useWindowManager";
import Window from "./window/Window.vue";

const { attachWindows, notificationWindows, closeWindow, minimizeWindow, bindWindowEl } = useWindowManager({});

const handleMaximize = () => {};

const computedStyle = (ws: WindowState) => {
  return {
    width: `${ws.size.width}px`,
    height: `${ws.size.height}px`,
  };
};
// 计算组件样式 - 支持多个组件垂直排列
const getNotificationStyle = (ws: WindowState, index: number) => {
  const baseStyle = {
    width: `${ws.size.width}px`,
    height: `${ws.size.height}px`,
    position: "fixed" as const,
    right: "10px",
    bottom: `${10 + index * 15}px`, // 每个组件间隔15px
    left: "unset",
    top: "unset",
    zIndex: 10001 + index, // 后面的组件层级更高
    borderRadius: "16px",
    boxShadow: "0 8px 32px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.08)",
    backdropFilter: "blur(20px)",
    overflow: "hidden",
    pointerEvents: "auto" as const,
    background: "rgba(255, 255, 255, 0.95)",
    opacity: "1",
    transition: "all 0.3s ease",
  };
  return baseStyle;
};

// 暴露方法
defineExpose({});
</script>

<style scoped>
.absolute-component-window {
  z-index: 10000;
  position: absolute;
  left: 1px;
  top: 1px;
  display: flex;
  flex-direction: column;
  background-color: var( --primary-color);
  border-radius: 16px;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.12),
    0 2px 8px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(20px);
  box-shadow: 0px 0px 1px 0px rgba(95, 108, 144, 0.3);
  overflow: hidden;
  pointer-events: auto;
  /* 根据展开方向设置变换原点 */
  transform-origin: var(--transform-origin-x) var(--transform-origin-y);

  background-size: cover;
  width: calc(100% - 2px);
  height: calc(100% - 2px);
  :deep(.ask-tab .chat-input-area .content-class) {
    background: rgba(255, 255, 255, 0.3) !important;
  }
  :deep(.chat-bottom-class .content-class) {
    background: rgba(255, 255, 255, 0) !important;
  }
  :deep(.notes-area) {
    background: rgba(255, 255, 255, 0) !important;
  }
}

.notification-windows {
  display: flex;
  flex-direction: column-reverse; /* 关键：反向排列 */
  position: absolute;
  right: 20px;
  bottom: 20px;
  z-index: 10000;
  width: 360px;
  height: fit-content;
}

/* 组件内容样式 */
.absolute-component-window > *:not(.window-controls) {
  flex: 1;
  overflow: auto;
  background: transparent;
}

.absolute-component-window > *:not(.window-controls).with-controls {
  /* 当有控制栏时，内容区域需要减去控制栏的高度 */
  height: calc(100% - 32px);
}

/* 根据展开方向设置变换原点 */
.absolute-component-window[data-direction="right-bottom"] {
  --transform-origin-x: 0%;
  --transform-origin-y: 0%;
}

.absolute-component-window[data-direction="left-top"] {
  --transform-origin-x: 100%;
  --transform-origin-y: 100%;
}

.absolute-component-window[data-direction="left-bottom"] {
  --transform-origin-x: 100%;
  --transform-origin-y: 0%;
}

.absolute-component-window[data-direction="right-top"] {
  --transform-origin-x: 0%;
  --transform-origin-y: 100%;
}

/* 组件样式保持不变，变换动画由 transitionUtils 处理 */

.component-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  font-size: 14px;
  color: #999;
}
</style>
