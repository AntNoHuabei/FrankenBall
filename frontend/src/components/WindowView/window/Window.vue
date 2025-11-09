<template>
  <div class="window">
    <Toolbar
      :win-id="props.winId"
      :mouse-inside="mouseInside || toolbarMouseInside"
      @mouseenter="toolbarMouseInside = true"
      @mouseleave="toolbarMouseInside = false"
    ></Toolbar>

    <div style="flex: 1; overflow-y: auto" @mouseenter="mouseInside = true" @mouseleave="mouseInside = false">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { provide, ref } from "vue";
import Toolbar from "./Toolbar.vue";

interface Props {
  winId: string;
  title: string;
}

const props = withDefaults(defineProps<Props>(), {});

provide("winId", props.winId);

const mouseInside = ref(false);
const toolbarMouseInside = ref(false);
</script>

<style lang="less" scoped>
.window {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  position: absolute;
  left: 0;
  top: 0;
  :deep(.section-title) {
    font-size: 12px;
  }
  :deep(.todo-content .todo-text) {
    font-size: 12px;
  }
}
</style>
