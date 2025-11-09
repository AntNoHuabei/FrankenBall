<script setup lang="ts">
import { ref } from "vue";
import FloatingBall from "./components/FloatingBall.vue";
import CoordGrid from "./CoordGrid.vue";
import {useFloatingMouseEvents} from "./composables-local/useFloatingBallMouseEvents";
import {useFloatingBallAutoAttachEdge} from "./composables-local/useFloatingBallAutoAttachEdge";
import MenuView from "./components/MenuView.vue";
import WindowView from "./components/WindowView/WindowView.vue";
import {useViewManager} from "./composables-local/useViewManager";
import {useSharedStatus} from "./composables-local/useSharedStatus";

const appContainerRef = ref<HTMLElement | null>(null);

const floatingBallRef = ref<InstanceType<typeof FloatingBall> | null>(null);
const menuViewRef = ref<InstanceType<typeof MenuView> | null>(null);
const componentViewRef = ref<InstanceType<typeof WindowView> | null>(null);


const { viewVisibility } = useViewManager({
  ballViewEle: floatingBallRef,
  menuViewEle: menuViewRef,
  windowViewEle: componentViewRef,
  rootElement: appContainerRef,
});


/**
 * 悬浮球鼠标事件
 */
useFloatingMouseEvents(appContainerRef)
/**
 * 悬浮球自动吸附侧边
 */
useFloatingBallAutoAttachEdge(appContainerRef)


const {ballSize} = useSharedStatus();

</script>

<template>
  <div class="app-container">
    <div ref="appContainerRef" class="core-content   mouse-interactive"  :style="{'border-radius':`${ballSize/2 +1}px`}">
      <!-- 悬浮球组件 -->
      <!-- 绝对定位的悬浮球 -->
      <FloatingBall ref="floatingBallRef" v-show="viewVisibility.ball|| viewVisibility.menu" />

      <!-- 菜单视图 -->
      <MenuView ref="menuViewRef" v-show="viewVisibility.ball|| viewVisibility.menu"/>

      <!-- 组件视图 -->
      <WindowView ref="componentViewRef" v-show="viewVisibility.window" />
    </div>

    <transition name="fade">
      <CoordGrid />
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
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  pointer-events: auto;
  /* cursor: move; */
  z-index: 1000;
  overflow: hidden;
  border: #0f0f0f  solid 1px;
  background: #28284e;
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
