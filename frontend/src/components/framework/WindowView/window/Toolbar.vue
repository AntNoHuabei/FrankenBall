<template>
  <div
    class="window-controls"
    v-show="
      windowState.windowConfig.toolbar?.show &&
      (!windowState.windowConfig.toolbar?.hoverShow ||
        (windowState.windowConfig.toolbar?.hoverShow && props.mouseInside))
    "
  >
    <div class="title-bar">
      <!-- Teleport 目标区域 - 用于接收各个组件的工具栏 -->
      <div class="toolbar-container mouse-drag" :ref="(el: any) => setToolbarRef(el)">
        <!-- 自定义工具栏内容 -->
      </div>
      <div class="control-buttons">
        <!-- 默认控制按钮 -->
        <button
          v-if="windowState.windowConfig.toolbar?.showMinimizeButton"
          class="control-btn minimize-btn"
          @click="handleMinimize"
        >
          <svg viewBox="0 0 12 12" fill="none">
            <path d="M2 6h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
          </svg>
        </button>
        <button
          v-if="windowState.windowConfig.toolbar?.showMaximizeButton"
          class="control-btn maximize-btn"
          @click="handleMaximize"
        >
          <svg v-if="!isMaximized" viewBox="0 0 12 12" fill="none">
            <rect x="2" y="2" width="8" height="8" stroke="currentColor" stroke-width="1.5" rx="1" />
          </svg>
          <svg v-else viewBox="0 0 12 12" fill="none">
            <path
              d="M3 3.5h6M3 3.5v6M3 3.5L2.5 3M9 3.5v6M9 3.5L9.5 3M3 9.5h6M3 9.5L2.5 10M9 9.5L9.5 10"
              stroke="currentColor"
              stroke-width="1.5"
              stroke-linecap="round"
            />
          </svg>
        </button>
        <button
          class="control-btn close-btn"
          @click="handleClose"
          v-if="windowState.windowConfig.toolbar?.showCloseButton"
        >
          <svg viewBox="0 0 12 12" fill="none">
            <path d="M3 3l6 6M9 3L3 9" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useSharedToolbarRef } from "../composables/useSharedToolbarRef";
import { useWindowManager } from "../../composables-local/useWindowManager";

const props = defineProps({
  winId: {
    type: String,
    default: "",
  },
  mouseInside: {
    type: Boolean,
    default: false,
  },
});

const { setToolbarRef } = useSharedToolbarRef(props.winId);

const emit = defineEmits(["maximize"]);

const isMaximized = ref(false);

const { minimizeWindow,maximizeWindow, closeWindow, getWindowState ,setWindowSize} = useWindowManager();

const windowState = getWindowState(props.winId);

const handleMinimize = () => {
  //先检查是否有自定义的最小化逻辑

  if (windowState.value?.windowConfig.toolbar?.onMinimize) {
    const isHandled = windowState.value.windowConfig.toolbar.onMinimize(windowState.value);
    if (!isHandled) {
      //如果没有处理，默认最小化
      minimizeWindow(props.winId);
    }
    return;
  } else {
    //如果没有自定义最小化逻辑，默认最小化
    minimizeWindow(props.winId);
  }
};

const handleMaximize = () => {


  if (windowState.value?.windowConfig.toolbar?.onMaximize) {
    const isHandled = windowState.value?.windowConfig.toolbar.onMaximize(windowState.value);
    if (!isHandled) {

      maximizeWindow(props.winId)
    }
    return;
  } else {
    maximizeWindow(props.winId)
  }
};

const handleClose = () => {
  //先检查是否有自定义的关闭逻辑

  if (windowState.value?.windowConfig.toolbar?.onClose) {
    const isHandled = windowState.value?.windowConfig.toolbar.onClose(windowState.value);
    if (!isHandled) {
      //如果没有处理，默认关闭
      closeWindow(props.winId);
    }
    return;
  } else {
    //如果没有自定义关闭逻辑，默认关闭
    closeWindow(props.winId);
  }
};
</script>

<style scoped>
.window-controls {
  width: 100%;
  height: 32px;
  background: transparent;
  display: flex;
  flex-direction: column;
}

.title-bar {
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
}

.toolbar-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: left;
  height: 100%;
}

.title-text {
  flex: 1;
  font-size: 13px;
  color: #333;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  -webkit-app-region: no-drag;
}

.control-buttons {
  display: flex;
  gap: 6px;
  -webkit-app-region: no-drag;
}

.control-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s ease;
  padding: 0;
}

.control-btn svg {
  width: 12px;
  height: 12px;
  color: #ffffff80;
  transition: color 0.2s ease;
}

.control-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.control-btn:hover svg {
  color: #333;
}

.control-btn:active {
  transform: scale(0.95);
}

.close-btn:hover {
  background: #ff5f57;
}

.close-btn:hover svg {
  color: white;
}

.minimize-btn:hover {
  background: #ffbd2e;
}

.minimize-btn:hover svg {
  color: white;
}

.maximize-btn:hover {
  background: #28ca42;
}

.maximize-btn:hover svg {
  color: white;
}

.window-content {
  flex: 1;
  overflow: auto;
  overflow-x: hidden;
  background: white;
  background-size: cover;
}
:deep(.window-content .ask-tab .chat-input-area .content-class) {
  background: rgba(255, 255, 255, 0.3) !important;
}
:deep(.window-content .chat-bottom-class .content-class) {
  background: rgba(255, 255, 255, 0) !important;
}

.window-content::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}
</style>
