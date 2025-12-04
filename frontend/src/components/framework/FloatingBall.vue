<script setup lang="ts">

import {useSharedStatus} from "./composables-local/useSharedStatus";
import  "@wailsio/runtime";
const {ballSize} = useSharedStatus();




</script>

<template>
  <div
      ref="ballRef"
      :style="{width: `${ballSize}px`, height: `${ballSize}px`}"
      class="floating-ball mouse-drag"
  >
    <div class="ball-content">
      <img src="/logo.png" style="width: 60%;height: 60%;-webkit-user-drag: none">
    </div>
  </div>
</template>

<style scoped>
.floating-ball {
  position: absolute;
  border-radius: 50%;
  /* cursor: move; */
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  transition: transform 0.2s ease;
  z-index: 9999;
  user-select: none;
}

.ball-content {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  user-select: none;
}

.ball-icon {
  font-size: 48px;
  color: white;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}
</style>

